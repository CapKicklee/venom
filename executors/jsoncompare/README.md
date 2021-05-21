# Venom - Executor JSON Comparison

Step to compare 2 JSON files:

It uses the `jsondiff` package: https://github.com/nsf/jsondiff to semantically compare two json files

## Input

```yaml
name: JSON Comparison testsuite
testcases:
- name: perfectmatch json comparison testcase
  steps:
  - type: jsoncompare
    expectedjson: testa.json
    actualjson: testb.json
    assertions:
    - result.perfectmatch ShouldBeTrue
    - result.difference ShouldBeBlank

- name: supersetmatch json comparison testcase
  steps: 
  - type: jsoncompare
    expectedjson: testa.json
    actualjson: testc.json
    assertions:
    - result.supersetmatch ShouldBeTrue
    - result.difference ShouldNotBeBlank

- name: nomatch json comparison testcase
  steps:
  - type: jsoncompare
    expectedjson: testa.json
    actualjson: testd.json
    assertions:
    - result.nomatch ShouldBeTrue
    - result.difference ShouldNotBeBlank
```

## Output

```
result.difference
result.expected
result.result
result.systemout
result.systemerr
result.perfectmatch
result.supersetmatch
result.nomatch
result.err
```

- result.difference: expected file annotated with differences found in actual file (empty if no differences)
- result.expected: expected file as string
- result.result: actual file as string
- result.systemout: comparison explanation
- result.systemerr: error message
- result.err: if exists, this field contains error
- result.perfectmatch: indicates whether it's a perfect match or not
- result.supersectmatch: indicates whether it's a superset match or not
- result.nomatch: indicates whether it's a no match or not

## Default assertion

```yaml
result.perfectmatch ShouldBeTrue
```
