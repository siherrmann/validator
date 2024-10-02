# validator

## Description

Validator for structs to validate fields by tag

Validate validates a given struct by `vld` tags or updates a given struct by `upd` tags.
Validate needs a struct as input and can update the struct by map or json input.

## Validate

You can add a validate tag with the syntax `vld:"[requirement], [groups]"`.
Groups are seperated by a space (eg. `gr1min1 gr2max1`).
Conditions and operators in a requirement are seperated by a space (eg. `max0 || (min10 && max30)`).

All fields in the struct need a `vld` tag.
If you want to ignore one field in the validator you can add `vld:"-"`.
If you don't add the vld tag to every field the function will fail with an error.

## Update
You can add a validate tag with the syntax `upd:"[json_key], [requirement], [groups]"`.
The json key has to be a valid key from given json for unmarshalling or a key from a given map (eg. `gr1min1 gr2max1`).
Groups are seperated by a space (eg. `gr1min1 gr2max1`).
Conditions and operators in a requirement are seperated by a space (eg. `max0 || (min10 && max30)`).

Only fields in the struct you want to update need a `upd` tag.

## Requirement

You can build complex requirements by building a query of conditions, operators (`&&` (=AND) and `||` (=OR)) and groups (with `(` and `)`).

A complex example for a password check (min length 8, max length 30, at least one capital letter, one small letter, one digit and one special character) would be:
`vld:"min8 max30 rex^(.*[A-Z])+(.*)$ rex^(.*[a-z])+(.*)$ rex^(.*\\d)+(.*)$ rex^(.*[\x60!@#$%^&*()_+={};':\"|\\,.<>/?~-])+(.*)$"`.
In this example all connections are `&&` (=AND) connections. Because behind the requirement check is a little parser you can also do more complex requirements with multiple conditions grouped and connected with AND and OR connections.
You can do for example `vld:"max0 || ((min10 && max30) || equTest)"` for a string that has to be either empty, the string `Test` or between 10 and 30 characters long. And yes, the outer brackets are not needed 😉.

## Condition types

`-` - ignores the field

`equ` - equal (value or length)

`neq` - not equal (value or length)

`min` - min (value or length)

`max` - max (value or length)

`con` - contains

`rex` - regular expression

## Usage of conditions

Conditions have different usages per variable type:

**equ** - `int/float/string == condition`, `len(array) == condition`

**neq** - `int/float/string != condition`, `len(array) != condition`

**min** - `int/float >= condition`, `len(strings.TrimSpace(string)/array) >= condition`

**max** - `int/float <= condition`, `len(strings.TrimSpace(string)/array) <= condition`

**con** - `strings.Contains(string, condition)`, `contains(array, condition)`, int/float ignored

**rex** - `regexp.MatchString(condition, strconv.Itoa(int)/strconv.FormatFloat(float, 'f', 3, 64)/string)`, array ignored

For con you need to put in a condition that is convertable to the underlying type of the arrary.
Eg. for an array of int the condition must be convertable to int (bad: `` `vld:"conA"` ``, good: `` `vld:"con1"` ``).

In the case of rex the int and float input will get converted to a string (`strconv.Itoa(int)` and `fmt.Sprintf("%f", f)`).
If you want to check more complex cases you can obviously replace `equ`, `neq`, `min`, `max` and `con` with one regular expression.

## Groups

You also have the posibillity to add groups. So if you want to check on an update, that at least one field is updated, you can add all fields to a group `upd:"field_name, min1, gr1min1"`.
Usable conditions for groups are `min` and `max`.

A small code example would be:

```go
type Error struct {
	ID                  int       `json:"id" vld:"-"`
	StatusCode          int       `json:"status_code" vld:"min100" upd:"status_code, min100, gr1min1"`
	Message             string    `json:"message" vld:"min1" vld:"min1, gr1min1" upd:"status_code, min1, gr1min1"`
	UnderlyingException string    `json:"underlying_exception" vld:"min1, gr1min1" upd:"status_code, min1, gr1min1"`
	CreatedAt           time.Time `json:"created_at" vld:"-"`
}
```

`ID` and `CreatedAt` are not getting validated, because you would probably not insert these but create these on db level.
`StatusCode` is necessary on creation and on update it is in a group with `Message` and `UnderlyingException` where one of them must be given.
One of `Message` and `UnderlyingException` is required on creation.

## Testing

To run the tests run `go test .`.

## Benchmark

To run benchmarks run `go test -bench . -count 100 > bench.txt` (with memory allocation would be `go test -bench . -benchmem -count 100 > bench.txt` but they are 0). To see the results in a nice way after the run install `go install golang.org/x/perf/cmd/benchstat@latest` and log the results to the console: `benchstat bench.txt`.

## Creating a new version

To create a new tagged version run eg. `git tag v0.1.3`. To push and publish it run eg. `git push origin v0.1.3`.