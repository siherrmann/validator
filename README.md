# validator

[![Go Reference](https://pkg.go.dev/badge/github.com/siherrmann/validator.svg)](https://pkg.go.dev/github.com/siherrmann/validator)
[![Go Coverage](https://github.com/siherrmann/validator/wiki/coverage.svg)](https://raw.githack.com/wiki/siherrmann/validator/coverage.html)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/siherrmann/validator/blob/master/LICENSE)

Validator for structs and json maps to validate fields by tag.

# 💡 Goal of this package

The goal of the validator package is to provide a fast, simple validation tool that can validate existing structs and validate url values and json bodys to update an empty struct with these values if they are valid.

Main focus is centralisation of validation to one object. You can obviously do one request object for each request type (which is sometimes still needed), but with this package you can do even something like that:

```go
type User struct {
	ID	int	`json:"id" del:"min1"`
	Name	string	`json:"name" cre:"min1" upd:"min1, gr1min1"`
	Adress	string	`json:"adress" cre:"min1" upd:"min1, gr1min1"`
}
```

You can see the 3 different tags (`cre`, `upd`, `del`). With these tags you can use the same object in 3 different handlers with 3 different validations. In this example for creation you need name and adress, for an update you need at least one of name or adress, and for deletion you need an id greater 0.

---

# 🚀 Getting started

`Validate` validates a given struct by `vld` or custom tags. `ValidateAndUpdate` does update the given struct with the given json after validating the json. `UnmarshalValidateAndUpdate` and similar functions are unpacking something (request body or url values), then validating the input and updating the given struct. `ValidateAndUpdateWithValidation` gives you the ability to update a map with the values from a json map by using an array of `Validation` (which is the equivalent for tags in a struct).

---

# Validation

You can add a validate tag with the syntax `vld:"[requirement], [groups]"`.
Groups are seperated by a space (eg. `gr1min1 gr2max1`).
Conditions and operators in a requirement are seperated by a space (eg. `max0 || (min10 && max30)`).

All fields that you want to validate in the struct need a `vld` tag (or custom tag if specified).
If you don't want to validate the field you can add `vld:"-"`. If you then use an update function it does update it without validating.

## Requirement

You can build complex requirements by building a query of conditions, operators (`&&` (=AND) and `||` (=OR)) and groups (with `(` and `)`).

A complex example for a password check (min length 8, max length 30, at least one capital letter, one small letter, one digit and one special character) would be:
`vld:"min8 max30 rex^(.*[A-Z])+(.*)$ rex^(.*[a-z])+(.*)$ rex^(.*\\d)+(.*)$ rex^(.*[\x60!@#$%^&*()_+={};':\"|\\,.<>/?~-])+(.*)$"`.
In this example all connections are `&&` (=AND) connections. Because behind the requirement check is a little parser you can also do more complex requirements with multiple conditions grouped and connected with AND and OR connections.
You can do for example `vld:"max0 || ((min10 && max30) || equTest)"` for a string that has to be either empty, the string `Test` or between 10 and 30 characters long. And yes, the outer brackets are not needed 😉.

### Condition types

Conditions have different usages per variable type:
- `-` - Not validating/update without validating.
- `equ` - `int/float/string == condition`, `len(array) == condition`
- `neq` - `int/float/string != condition`, `len(array) != condition`
- `min` - `int/float >= condition`, `len(strings.TrimSpace(string)/array) >= condition`
- `max` - `int/float <= condition`, `len(strings.TrimSpace(string)/array) <= condition`
- `con` - `strings.Contains(string, condition)`, `contains(array, condition)`, int/float ignored
- `nco` - `!strings.Contains(string, condition)`, `!contains(array, condition)`, int/float ignored
- `frm` - Checks if given comma seperated list contains value/every item in array/every key in map.
- `nfr` - Checks if given comma seperated list does not contain value/every item in array/every key in map.
- `rex` - `regexp.MatchString(condition, strconv.Itoa(int)/strconv.FormatFloat(float, 'f', 3, 64)/string)`, array ignored

For con you need to put in a condition that is convertable to the underlying type of the arrary.
Eg. for an array of int the condition must be convertable to int (bad: `vld:"conA"`, good: `vld:"con1"`).

In the case of rex the int and float input will get converted to a string (`strconv.Itoa(int)` and `fmt.Sprintf("%f", f)`).
If you want to check more complex cases you can obviously replace `equ`, `neq`, `min`, `max` and `con` with one regular expression.

## Groups

You also have the posibillity to add groups. So if you want to check on an update, that at least one field is updated, you can add all fields to a group `upd:"min1, gr1min1"`.
Usable conditions for groups are `min` and `max`.

A small code example would be:

```go
type Error struct {
	ID                  int       `json:"id"`
	StatusCode          int       `json:"status_code" vld:"min100" upd:"min100, gr1min1"`
	Message             string    `json:"message" vld:"min1" upd:"min1, gr1min1"`
	UnderlyingException string    `json:"underlying_exception" vld:"min1, gr1min1" upd:"min1, gr1min1"`
	CreatedAt           time.Time `json:"created_at" vld:"-"`
}
```

`ID` and `CreatedAt` in this example are not getting validated, because you would probably not insert these but create these on database level.
`StatusCode` is necessary on creation and on update it is in a group with `Message` and `UnderlyingException` where one of them must be given.
One of `Message` and `UnderlyingException` is required on creation.

---

# Testing

To run the tests run `go test ./...`.

## Benchmark

To run benchmarks run `go test -bench . -count 100 > bench.txt` (with memory allocation would be `go test -bench . -benchmem -count 100 > bench.txt` but they are 0). To see the results in a nice way after the run install `go install golang.org/x/perf/cmd/benchstat@latest` and log the results to the console with `benchstat bench.txt`.

---

# ⭐ Features

- Validate structs by `vld` tag.
- Validate json or form values by struct tag.
- Validate json or form values by validations.
- Validate and update a struct from `JsonMap`.
- Use custom tags.
- Use multiple custom tags in one struct for multiple validation situations.
