package validator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/siherrmann/validator/model"
	"github.com/siherrmann/validator/validators"
)

type StructValue struct {
	Error  error
	Groups []string
}

// UnmarshalAndValidate unmarshals given json ([]byte) into pointer v.
// For more information to Validate look at [Validate(v any) error].
func UnmarshalAndValidate(data []byte, v any) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("value has to be of kind pointer, was %T", value)
	}
	if value.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("value has to be of kind struct, was %T", value)
	}

	err := json.Unmarshal(data, v)
	if err != nil {
		return fmt.Errorf("error unmarshalling %T: %v", value, err)
	}

	err = Validate(v)
	if err != nil {
		return err
	}

	return nil
}

// Validate validates a given struct by vld tags.
// Validate needs a struct as input.
//
// All fields in the struct need a vld tag.
// If you want to use multiple conditions you can add them with a space in between them.
//
// A complex example for password would be:
// `vld:"min8 max30 rex^(.*[A-Z])+(.*)$ rex^(.*[a-z])+(.*)$ rex^(.*\\d)+(.*)$ rex^(.*[\x60!@#$%^&*()_+={};':\"|\\,.<>/?~-])+(.*)$"`
//
// If you want to ignore one field in the validator you can add `vld:"-"`.
// If you don't add the vld tag to every field the function will fail with an error.
//
// Conditions have different usages per variable type:
//
// equ - int/float/string/bool == condition, len(array) == condition
//
// neq - int/float/string/bool != condition, len(array) != condition
//
// min - int/float >= condition, len(string/array) >= condition
//
// max - int/float <= condition, len(string/array) <= condition
//
// con - strings.Contains(string, condition), contains(array, condition), int/float ignored
//
// rex - regexp.MatchString(condition, int/float/string), array ignored
//
// For con you need to put in a condition that is convertable to the underlying type of the arrary.
// Eg. for an array of int the condition must be convertable to int (bad: `vld:"conA"`, good: `vld:"con1"`).
//
// In the case of rex the int and float input will get converted to a string (strconv.Itoa(int) and fmt.Sprintf("%f", f)).
// If you want to check more complex cases you can obviously replace equ, neq, min, max and con with one regular expression.
func Validate(v any) error {
	// check if value is a pointer to a struct
	structValue := reflect.ValueOf(v)
	if structValue.Kind() != reflect.Ptr {
		return fmt.Errorf("structValue has to be of kind pointer, was %T", structValue)
	}
	if structValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("structValue has to be of kind struct, was %T", structValue)
	}

	// get valid reflect structValue of struct
	structFull := structValue.Elem()

	keys := []string{}
	groups := map[string]*model.Group{}
	groupSize := map[string]int{}
	groupErrors := map[string][]error{}

	for i := 0; i < structFull.Type().NumField(); i++ {
		value := structFull.Field(i)
		fieldName := structFull.Type().Field(i).Name
		tag := structFull.Type().Field(i).Tag.Get(string(model.VLD))

		validation := &model.Validation{}
		err := validation.Fill(tag, model.VLD, value)
		if err != nil {
			return err
		}

		// early return/continue for empty requirement
		if strings.TrimSpace(validation.Requirement) == string(model.NONE) {
			continue
		}

		if len(validation.Key) > 0 && slices.Contains(keys, validation.Key) {
			return fmt.Errorf("duplicate validation key: %v", validation.Key)
		} else {
			keys = append(keys, validation.Key)
		}

		for _, g := range validation.Groups {
			groups[g.Name] = g
			groupSize[g.Name]++
		}

		_, err = ValidateValueWithParser(value, validation)
		if err != nil && len(validation.Groups) == 0 {
			return fmt.Errorf("field %v of %v invalid: %v", fieldName, reflect.TypeOf(v), err.Error())
		} else if err != nil {
			for _, group := range validation.Groups {
				groupErrors[group.Name] = append(groupErrors[group.Name], fmt.Errorf("field %v of %v invalid: %v", fieldName, reflect.TypeOf(v), err.Error()))
				continue
			}
		}
	}

	err := validators.ValidateGroups(groups, groupSize, groupErrors)
	if err != nil {
		return err
	}

	return nil
}
