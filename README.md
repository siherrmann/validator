# validator

## Description

Validator for structs to validate fields by tag

`Validate` validates a given struct by `vld` tags. `ValidateAndUpdate` does update the given struct with the given json after validating the json. `UnmarshalValidateAndUpdate` and similar functions are unpacking something (request body or url values), then validating the input and updating the given struct.

## Validation

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

## Condition types

No condition does neither validate nor update the field.

`-` - not validating/update without validating

`equ` - equal (value or length)

`neq` - not equal (value or length)

`min` - min (value or length)

`max` - max (value or length)

`con` - contains

`nco` - contains not

`frm` - is contained in given list (whitelist)

`nfr` - is not contained in given list (blacklist)

`rex` - regular expression

## Usage of conditions

Conditions have different usages per variable type:

**equ** - `int/float/string == condition`, `len(array) == condition`

**neq** - `int/float/string != condition`, `len(array) != condition`

**min** - `int/float >= condition`, `len(strings.TrimSpace(string)/array) >= condition`

**max** - `int/float <= condition`, `len(strings.TrimSpace(string)/array) <= condition`

**con** - `strings.Contains(string, condition)`, `contains(array, condition)`, int/float ignored

**nco** - `!strings.Contains(string, condition)`, `!contains(array, condition)`, int/float ignored

**frm** - checks if given comma seperated list contains value/every item in array is contained in the comma seperated list, bool ignored

**nfr** - checks if given comma seperated list does not contain value/every item in array is not contained in the comma seperated list, bool ignored

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
	ID                  int       `json:"id"`
	StatusCode          int       `json:"status_code" vld:"min100" upd:"min100, gr1min1"`
	Message             string    `json:"message" vld:"min1" upd:"min1, gr1min1"`
	UnderlyingException string    `json:"underlying_exception" vld:"min1, gr1min1" upd:"min1, gr1min1"`
	CreatedAt           time.Time `json:"created_at" vld:"-"`
}
```

`ID` and `CreatedAt` are not getting validated, because you would probably not insert these but create these on db level.
`StatusCode` is necessary on creation and on update it is in a group with `Message` and `UnderlyingException` where one of them must be given.
One of `Message` and `UnderlyingException` is required on creation.

## Code example

```go
package main

import (
	"io"
	"net/http"
    "fmt"
    "time"

	"github.com/siherrmann/validator"
)

type Error struct {
	ID                  int       `json:"id" del:"min1"`
	StatusCode          int       `json:"status_code" vld:"min100" upd:"min100, gr1min1"`
	Message             string    `json:"message" vld:"min1" upd:"min1, gr1min1"`
	UnderlyingException string    `json:"underlying_exception" vld:"min1, gr1min1" upd:"min1, gr1min1"`
	CreatedAt           time.Time `json:"created_at" vld:"-"`
}

func HandleError(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("error reading request body: %v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// This uses the default `vld` tag. It unmarshals the request body, validates it
	// and updates the `newError` variable. You could use that for creating a new instance
	// of error with all needed parameters.
	var newError Error
	err = validator.UnmarshalValidateAndUpdate(body, &newError)
	if err != nil {
		http.Error(w, fmt.Sprintf("error validating new errror: %v", err), http.StatusBadRequest)
		return
	}

	// For updating an error you could use the `upd` tag (including goups) to make sure 
	// that at least one of the values is updated and if so is valid.
	var exisitingErrorFromDb Error
	err = validator.UnmarshalValidateAndUpdate(body, &exisitingErrorFromDb, "upd")
	if err != nil {
		http.Error(w, fmt.Sprintf("error validating error update: %v", err), http.StatusBadRequest)
		return
	}

	// If you would only want to validate some given struct with another tag
	// (for example if you want to check if a given error is valid to delete,
	// containing an id in this case) you could do:
	err = validator.Validate(&newError, "del")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error unmarshaling JSON: %v", err), http.StatusBadRequest)
		return
	}
}
```

## Testing

To run the tests run `go test .`.

## Benchmark

To run benchmarks run `go test -bench . -count 100 > bench.txt` (with memory allocation would be `go test -bench . -benchmem -count 100 > bench.txt` but they are 0). To see the results in a nice way after the run install `go install golang.org/x/perf/cmd/benchstat@latest` and log the results to the console: `benchstat bench.txt`.

## Creating a new version

To create a new tagged version run eg. `git tag v0.1.3`. To push and publish it run eg. `git push origin v0.1.3`. You can get the latest tags with `git log $(git describe --tags)`.