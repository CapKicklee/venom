package jsoncompare

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mattn/go-zglob"
	"github.com/nsf/jsondiff"
	"github.com/mitchellh/mapstructure"

	"github.com/ovh/venom"
)

// Name of executor
const Name = "jsoncompare"

// New returns a new Executor
func New() venom.Executor {
	return &Executor{}
}

// Executor struct
type Executor struct {
	ExpectedJson	string	 `json:"expectedjson,omitempty" yaml:"expectedjson,omitempty"`
	ActualJson		string	 `json:"actualjson,omitempty" yaml:"actualjson,omitempty"`
}

// ZeroValueResult return an empty implementation of this executor result
func (Executor) ZeroValueResult() interface{} {
	return Result{}
}

// Result represents a step result
type Result struct {
	Difference  	string		`json:"difference,omitempty" yaml:"difference,omitempty"`
	Expected		string		`json:"expected,omitempty" yaml:"expected,omitempty"`
	Result			string		`json:"result,omitempty" yaml:"result,omitempty"`
	Systemout   	string		`json:"systemout,omitempty" yaml:"systemout,omitempty"` // put in testcase.Systemout by venom if present
	Systemerr   	string		`json:"systemerr,omitempty" yaml:"systemerr,omitempty"` // put in testcase.Systemerr by venom if present
	Err         	string      `json:"error,omitempty" yaml:"error,omitempty"`
	PerfectMatch	bool		`json:"perfectmatch,omitempty" yaml:"perfectmatch,omitempty"`
	SupersetMatch	bool		`json:"supersetmatch,omitempty" yaml:"supersetmatch,omitempty"`
	NoMatch			bool		`json:"nomatch,omitempty" yaml:"nomatch,omitempty"`
}

// GetDefaultAssertions return default assertions for this executor
// Optional
func (Executor) GetDefaultAssertions() *venom.StepAssertions {
	return &venom.StepAssertions{Assertions: []string{"result.perfectmatch ShouldBeTrue"}}
}

// Run execute TestStep
func (Executor)	Run(ctx context.Context, step venom.TestStep) (interface{}, error) {
	// transform step to Executor Instance
	var e Executor
	if err := mapstructure.Decode(step, &e); err != nil {
		return nil, err
	}

	r := Result{}

	if e.ExpectedJson == "" || e.ActualJson == "" {
		return nil, fmt.Errorf("Invalid params, 2 json strings expected")
	}

	workdir := venom.StringVarFromCtx(ctx, "venom.testsuite.workdir")
	expectedJson, errr := e.readJSONFile(workdir, e.ExpectedJson)
	if errr != nil {
		r.Err = errr.Error()
	}

	actualJson, err := e.readJSONFile(workdir, e.ActualJson)
	if err != nil {
		r.Err = errr.Error()
	}

	r.Expected = string(expectedJson[:])
	r.Result = string(actualJson[:])

	opts := jsondiff.DefaultHTMLOptions()
	diff, text := jsondiff.Compare([]byte (r.Result), []byte (r.Expected), &opts)
	switch int(diff) {
	case 0:
		r.PerfectMatch = true
		r.Systemout = "Perfect match between the 2 json files"
	case 1:
		r.SupersetMatch = true
		r.Systemout = "The actual json file is richer than the expected one, there are few fields more"
		// will return '' if perfect match, otherwise, will return the expected file annotated with the differences
		// see the example here : https://nosmileface.dev/jsondiff/
		r.Difference = text
	case 2:
		r.NoMatch = true
		r.Systemout = "The json files are different, some fields do not have the same value"
		// will return '' if perfect match, otherwise, will return the expected file annotated with the differences
		// see the example here : https://nosmileface.dev/jsondiff/
		r.Difference = text
	default:
		r.Systemerr = "There was a problem during files analysis. Please check both files and run again"
	}
	
	return r, nil
}

func (e *Executor) readJSONFile(workdir string, path string) (string, error) {
	absPath := filepath.Join(workdir, path)

	fileInfo, _ := os.Stat(absPath) // nolint
	if fileInfo != nil && fileInfo.IsDir() {
		absPath = filepath.Dir(absPath)
	}

	filesPath, errg := zglob.Glob(absPath)
	if errg != nil {
		return "", fmt.Errorf("Error reading files on path:%s :%s", absPath, errg)
	}

	if len(filesPath) == 0 {
		return "", fmt.Errorf("Invalid path '%s' or file not found", absPath)
	}

	var content string

	for _, f := range filesPath {
		f, erro := os.Open(f)
		if erro != nil {
			return "", fmt.Errorf("Error while opening file: %s", erro)
		}
		defer f.Close()

		h := md5.New()
		tee := io.TeeReader(f, h)

		b, errr := ioutil.ReadAll(tee)
		if errr != nil {
			return "", fmt.Errorf("Error while reading file: %s", errr)
		}
		content += string(b)
	}

	return content, nil
}
