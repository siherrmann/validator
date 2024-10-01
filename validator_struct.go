package validator

import (
	"encoding/json"
	"fmt"
	"reflect"
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
// equ - int/float/string == condition, len(array) == condition
//
// neq - int/float/string != condition, len(array) != condition
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
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("value has to be of kind pointer, was %T", value)
	}
	if value.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("value has to be of kind struct, was %T", value)
	}

	// get valid reflect value of struct
	structFull := value.Elem()

	groups := map[string]string{}
	groupSize := map[string]int{}
	groupErrors := map[string][]error{}

	for i := 0; i < structFull.Type().NumField(); i++ {
		tag := structFull.Type().Field(i).Tag.Get("vld")
		if len(strings.TrimSpace(tag)) == 0 || strings.TrimSpace(tag) == string(model.NONE) {
			continue
		}

		tagSplit := strings.Split(tag, ", ")
		condition := "-"
		if len(tagSplit) > 0 {
			condition = tagSplit[0]
		}

		value := structFull.Field(i)
		fieldName := structFull.Type().Field(i).Name

		groupsValue := []string{}
		groupsString := []string{}
		if len(tagSplit) > 1 {
			groupsString = strings.Split(tagSplit[1], " ")

			for _, g := range groupsString {
				group, err := model.GetGroup(g)
				if err != nil {
					return fmt.Errorf("error extracting group: %v", err)
				}

				condition, err := model.GetConditionByType(g, model.ConditionType(group))
				if err != nil {
					return fmt.Errorf("error extracting group condition: %v", err)
				}

				groupsValue = append(groupsValue, group)
				groups[group] = condition
				groupSize[group]++
			}
		}

		var err error
		switch value.Type().Kind() {
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			// TODO update to only parsing
			// err = ValidateValueWithoutParser(value, conditions, or, validators.CheckInt)
			err = ValidateValueWithParser(value, condition, validators.CheckInt)
		case reflect.Float64, reflect.Float32:
			// TODO update to only parsing
			// err = ValidateValueWithoutParser(value, conditions, or, validators.CheckFloat)
			err = ValidateValueWithParser(value, condition, validators.CheckFloat)
		case reflect.String:
			// TODO update to only parsing
			// err = ValidateValueWithoutParser(value, conditions, or, validators.CheckString)
			err = ValidateValueWithParser(value, condition, validators.CheckString)
		case reflect.Array, reflect.Slice:
			// TODO update to only parsing
			// err = ValidateValueWithoutParser(value, conditions, or, validators.CheckArray)
			err = ValidateValueWithParser(value, condition, validators.CheckArray)
		case reflect.Map:
			// TODO validate?
			// TODO update to only parsing
			// err = ValidateValueWithoutParser(value, conditions, or, validators.CheckMap)
			err = ValidateValueWithParser(value, condition, validators.CheckMap)
		case reflect.Struct:
			// TODO validate?
			// TODO update to only parsing
			// err = ValidateValueWithoutParser(value, conditions, or, validators.CheckStruct)
			err = ValidateValueWithParser(value, condition, validators.CheckTime)
		default:
			return fmt.Errorf("invalid field type for %v in %v: %v", fieldName, reflect.TypeOf(v), value.Type().Kind())
		}

		if err != nil && len(groupsString) == 0 {
			return fmt.Errorf("field %v of %v invalid: %v", fieldName, reflect.TypeOf(v), err.Error())
		} else if err != nil {
			for _, groupName := range groupsValue {
				groupErrors[groupName] = append(groupErrors[groupName], fmt.Errorf("field %v invalid: %v", fieldName, err.Error()))
			}
		}
	}

	err := validators.ValidateGroup(groups, groupSize, groupErrors)
	if err != nil {
		return err
	}

	return nil
}
