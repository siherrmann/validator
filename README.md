# validator

## Description

Validator for structs to validate fields by tag

Validate validates a given struct by vld tags.
Validate needs a struct as input.

All fields in the struct need a vld tag.
If you want to ignore one field in the validator you can add `` `vld:"-"` ``.
If you don't add the vld tag to every field the function will fail with an error.

If you want to use multiple conditions you can add them with a space in between them.

A complex example for a password check (min length 8, max length 30, at least one capital letter, one small letter, one digit and one special character) would be:
`` `vld:"min8 max30 rex^(.*[A-Z])+(.*)$ rex^(.*[a-z])+(.*)$ rex^(.*\\d)+(.*)$ rex^(.*[\x60!@#$%^&*()_+={};':\"|\\,.<>/?~-])+(.*)$"` ``

## Condition types

**-** - ignores the field

**equ** - equal (value or length)

**neq** - not equal (value or length)

**min** - min (value or length)

**max** - max (value or length)

**con** - contains

**rex** - regular expression

## Usage

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
If you want to check more complex cases you can obviously replace **equ**, **neq**, **min**, **max** and **con** with one regular expression.
