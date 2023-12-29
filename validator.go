package validator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

const (
	NONE      string = "-"
	EQUAL     string = "equ"
	NOT_EQUAL string = "neq"
	MIN_VALUE string = "min"
	MAX_VLAUE string = "max"
	CONTAINS  string = "con"
	REGX      string = "rex"
	OR        string = "||"
)

type StructValue struct {
	Error  error
	Groups []string
}

// UnmarshalAndValidate unmarshals given json ([]byte) into pointer v.
// For more information to Validate look at Validate(value any).
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
		if len(strings.TrimSpace(tag)) == 0 || strings.TrimSpace(tag) == NONE {
			continue
		}

		tagSplit := strings.Split(tag, ", ")

		value := structFull.Field(i)
		fieldName := structFull.Type().Field(i).Name

		or := false
		conditions := strings.Split(tagSplit[0], " ")
		if Contains(conditions, OR) {
			conditions = RemoveWhere[string](conditions, func(v string) bool {
				return v == OR
			})
			or = true
		}

		groupsValue := []string{}
		groupsString := []string{}
		if len(tagSplit) > 1 {
			groupsString = strings.Split(tagSplit[1], " ")

			for _, g := range groupsString {
				group := getConditionType(g)
				condition, err := getConditionByType(g, group)
				if err != nil {
					return fmt.Errorf("error extracting group: %v", err)
				}

				groupsValue = append(groupsValue, group)
				groups[group] = condition
				groupSize[group]++
			}
		}

		switch value.Type().Kind() {
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			valueTemp := value.Int()
			err := checkInt(int(valueTemp), conditions, or)
			if err != nil && len(groupsString) == 0 {
				return fmt.Errorf("field %v of %v invalid: %v", fieldName, reflect.TypeOf(v), err.Error())
			} else if err != nil {
				for _, groupName := range groupsValue {
					groupErrors[groupName] = append(groupErrors[groupName], fmt.Errorf("field %v invalid: %v", fieldName, err.Error()))
				}
			}
		case reflect.Float64, reflect.Float32:
			valueTemp := value.Float()
			err := checkFloat(valueTemp, conditions, or)
			if err != nil && len(groupsString) == 0 {
				return fmt.Errorf("field %v of %v invalid: %v", fieldName, reflect.TypeOf(v), err.Error())
			} else if err != nil {
				for _, groupName := range groupsValue {
					groupErrors[groupName] = append(groupErrors[groupName], fmt.Errorf("field %v invalid: %v", fieldName, err.Error()))
				}
			}
		case reflect.String:
			valueTemp := value.String()
			err := checkString(valueTemp, conditions, or)
			if err != nil && len(groupsString) == 0 {
				return fmt.Errorf("field %v of %v invalid: %v", fieldName, reflect.TypeOf(v), err.Error())
			} else if err != nil {
				for _, groupName := range groupsValue {
					groupErrors[groupName] = append(groupErrors[groupName], fmt.Errorf("field %v invalid: %v", fieldName, err.Error()))
				}
			}
		case reflect.Array, reflect.Slice:
			valueTemp := value
			err := checkArray(valueTemp, conditions, or)
			if err != nil && len(groupsString) == 0 {
				return fmt.Errorf("field %v of %v invalid: %v", fieldName, reflect.TypeOf(v), err.Error())
			} else if err != nil {
				for _, groupName := range groupsValue {
					groupErrors[groupName] = append(groupErrors[groupName], fmt.Errorf("field %v invalid: %v", fieldName, err.Error()))
				}
			}
		default:
			return fmt.Errorf("invalid field type for %v: %v", reflect.TypeOf(v), value.Type().Kind())
		}
	}

	err := validateGroup(groups, groupSize, groupErrors)
	if err != nil {
		return err
	}

	return nil
}

func UnmarshalValidateAndUpdate(jsonInput []byte, structToUpdate interface{}) error {
	jsonUnmarshaled := map[string]interface{}{}

	err := json.Unmarshal(jsonInput, &jsonUnmarshaled)
	if err != nil {
		return fmt.Errorf("error unmarshaling: %v", err)
	}

	err = ValidateAndUpdate(jsonUnmarshaled, structToUpdate)
	if err != nil {
		return fmt.Errorf("error updating struct: %v", err)
	}

	return nil
}

// ValidateAndUpdate validates a given struct by upd tags.
// ValidateAndUpdate needs a struct pointer and a json map as input.
// The given struct is updated by the values in the json map.
//
// All fields in the struct need a upd tag.
// The tag has to contain the key value for the json struct.
// If no tag is present the field in the struct is ignored and does not get updated.
//
// The second part of the tag contains the conditions for the validation.
//
// If you want to use multiple conditions you can add them with a space in between them.
//
// A complex example for password would be:
// `upd:"password, min8 max30 rex^(.*[A-Z])+(.*)$ rex^(.*[a-z])+(.*)$ rex^(.*\\d)+(.*)$ rex^(.*[\x60!@#$%^&*()_+={};':\"|\\,.<>/?~-])+(.*)$"`
//
// If you want don't want to validate the field you can add `upd:"json_key, -"`.
// If you don't add the upd tag to every field the function will fail with an error.
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
// Eg. for an array of int the condition must be convertable to int (bad: `upd:"array, conA"`, good: `upd:"array, con1"`).
//
// In the case of rex the int and float input will get converted to a string (strconv.Itoa(int) and fmt.Sprintf("%f", f)).
// If you want to check more complex cases you can obviously replace equ, neq, min, max and con with one regular expression.
func ValidateAndUpdate(jsonInput map[string]interface{}, structToUpdate interface{}) error {
	// check if value is a pointer to a struct
	value := reflect.ValueOf(structToUpdate)
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
		tag := structFull.Type().Field(i).Tag.Get("upd")
		if len(strings.TrimSpace(tag)) == 0 || strings.TrimSpace(tag) == NONE {
			continue
		}

		tagSplit := strings.Split(tag, ", ")

		value := structFull.Field(i)
		fieldName := structFull.Type().Field(i).Name

		jsonKey := tagSplit[0]

		or := false
		conditions := []string{}
		if len(tagSplit) > 1 {
			conditions = strings.Split(tagSplit[1], " ")
			if Contains(conditions, OR) {
				conditions = RemoveWhere[string](conditions, func(v string) bool {
					return v == OR
				})
				or = true
			}
		}

		groupsValue := []string{}
		groupsString := []string{}
		if len(tagSplit) > 2 {
			groupsString = strings.Split(tagSplit[2], " ")

			for _, g := range groupsString {
				group := getConditionType(g)
				condition, err := getConditionByType(g, group)
				if err != nil {
					return fmt.Errorf("error extracting group: %v", err)
				}

				groupsValue = append(groupsValue, group)
				groups[group] = condition
				groupSize[group]++
			}
		}

		var ok bool
		var err error
		var jsonValue interface{}
		if jsonValue, ok = jsonInput[jsonKey]; !ok {
			for _, groupName := range groupsValue {
				groupErrors[groupName] = append(groupErrors[groupName], fmt.Errorf("json %v key not in map", jsonKey))
			}
			continue
		}

		switch structValueType := value.Type().Kind(); structValueType {
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			var newInt int64
			if fl, ok := jsonValue.(float64); ok {
				// This case is for the case that json.Unmarshal unmarshals an int value into a float64 value.
				newInt = int64(fl)
			} else if _, ok := jsonValue.(int); ok {
				newInt = int64(jsonValue.(int))
			} else if _, ok := jsonValue.(int64); ok {
				newInt = int64(jsonValue.(int64))
			} else if _, ok := jsonValue.(int32); ok {
				newInt = int64(jsonValue.(int32))
			} else if _, ok := jsonValue.(int16); ok {
				newInt = int64(jsonValue.(int16))
			} else if _, ok := jsonValue.(int8); ok {
				newInt = int64(jsonValue.(int8))
			} else {
				return fmt.Errorf("input value for %v has to be of type %v, was %v", reflect.TypeOf(structToUpdate), structValueType, reflect.ValueOf(jsonValue).Kind())
			}

			err = checkInt(int(newInt), conditions, or)
			if err != nil && len(groupsString) == 0 {
				return fmt.Errorf("field %v of %v invalid: %v", fieldName, reflect.TypeOf(structToUpdate), err.Error())
			} else if err != nil {
				for _, groupName := range groupsValue {
					groupErrors[groupName] = append(groupErrors[groupName], fmt.Errorf("field %v invalid: %v", fieldName, err.Error()))
				}
				continue
			}

			fieldValue := reflect.ValueOf(structToUpdate).Elem().FieldByName(fieldName)
			err = setStructValueByJson(fieldValue, jsonKey, jsonValue)
			if err != nil && len(groupsString) == 0 {
				return fmt.Errorf("could not set field %v of %v: %v", fieldName, reflect.TypeOf(structToUpdate), err.Error())
			} else if err != nil {
				for _, groupName := range groupsValue {
					groupErrors[groupName] = append(groupErrors[groupName], fmt.Errorf("could not set field %v: %v", fieldName, err.Error()))
				}
			}
		case reflect.Float64, reflect.Float32:
			var newFloat float64
			if fl, ok := jsonValue.(float64); ok {
				newFloat = float64(fl)
			} else if fl, ok := jsonValue.(float32); ok {
				newFloat = float64(fl)
			} else {
				return fmt.Errorf("input value for %v has to be of type %v, was %v", reflect.TypeOf(structToUpdate), structValueType, reflect.ValueOf(jsonValue).Kind())
			}

			err = checkFloat(float64(newFloat), conditions, or)
			if err != nil && len(groupsString) == 0 {
				return fmt.Errorf("field %v of %v invalid: %v", fieldName, reflect.TypeOf(structToUpdate), err.Error())
			} else if err != nil {
				for _, groupName := range groupsValue {
					groupErrors[groupName] = append(groupErrors[groupName], fmt.Errorf("field %v invalid: %v", fieldName, err.Error()))
				}
				continue
			}

			fieldValue := reflect.ValueOf(structToUpdate).Elem().FieldByName(fieldName)
			err = setStructValueByJson(fieldValue, jsonKey, jsonValue)
			if err != nil && len(groupsString) == 0 {
				return fmt.Errorf("could not set field %v of %v: %v", fieldName, reflect.TypeOf(structToUpdate), err.Error())
			} else if err != nil {
				for _, groupName := range groupsValue {
					groupErrors[groupName] = append(groupErrors[groupName], fmt.Errorf("could not set field %v: %v", fieldName, err.Error()))
				}
			}
		case reflect.String:
			if _, ok := jsonValue.(string); !ok {
				return fmt.Errorf("input value for %v has to be of type %v, was %v", reflect.TypeOf(structToUpdate), structValueType, reflect.ValueOf(jsonValue).Kind())
			}

			err = checkString(jsonValue.(string), conditions, or)
			if err != nil && len(groupsString) == 0 {
				return fmt.Errorf("field %v of %v invalid: %v", fieldName, reflect.TypeOf(structToUpdate), err.Error())
			} else if err != nil {
				for _, groupName := range groupsValue {
					groupErrors[groupName] = append(groupErrors[groupName], fmt.Errorf("field %v invalid: %v", fieldName, err.Error()))
				}
				continue
			}

			fieldValue := reflect.ValueOf(structToUpdate).Elem().FieldByName(fieldName)
			err = setStructValueByJson(fieldValue, jsonKey, jsonValue)
			if err != nil && len(groupsString) == 0 {
				return fmt.Errorf("could not set field %v of %v: %v", fieldName, reflect.TypeOf(structToUpdate), err.Error())
			} else if err != nil {
				for _, groupName := range groupsValue {
					groupErrors[groupName] = append(groupErrors[groupName], fmt.Errorf("could not set field %v: %v", fieldName, err.Error()))
				}
			}
		case reflect.Bool:
			if _, ok := jsonValue.(bool); !ok {
				return fmt.Errorf("input value for %v has to be of type %v, was %v", reflect.TypeOf(structToUpdate), structValueType, reflect.ValueOf(jsonValue).Kind())
			}

			fieldValue := reflect.ValueOf(structToUpdate).Elem().FieldByName(fieldName)
			err = setStructValueByJson(fieldValue, jsonKey, jsonValue)
			if err != nil && len(groupsString) == 0 {
				return fmt.Errorf("could not set field %v of %v: %v", fieldName, reflect.TypeOf(structToUpdate), err.Error())
			} else if err != nil {
				for _, groupName := range groupsValue {
					groupErrors[groupName] = append(groupErrors[groupName], fmt.Errorf("could not set field %v: %v", fieldName, err.Error()))
				}
				continue
			}
		case reflect.Array, reflect.Slice:
			err = checkArray(reflect.ValueOf(jsonValue), conditions, or)
			if err != nil && len(groupsString) == 0 {
				return fmt.Errorf("field %v of %v invalid: %v", fieldName, reflect.TypeOf(structToUpdate), err.Error())
			} else if err != nil {
				for _, groupName := range groupsValue {
					groupErrors[groupName] = append(groupErrors[groupName], fmt.Errorf("field %v invalid: %v", fieldName, err.Error()))
				}
				continue
			}

			fieldValue := reflect.ValueOf(structToUpdate).Elem().FieldByName(fieldName)
			err = setStructValueByJson(fieldValue, jsonKey, jsonInput[jsonKey])
			if err != nil && len(groupsString) == 0 {
				return fmt.Errorf("could not set field %v of %v: %v", fieldName, reflect.TypeOf(structToUpdate), err.Error())
			} else if err != nil {
				for _, groupName := range groupsValue {
					groupErrors[groupName] = append(groupErrors[groupName], fmt.Errorf("could not set field %v: %v", fieldName, err.Error()))
				}
			}
		default:
			return fmt.Errorf("invalid field type for %v: %v", reflect.TypeOf(structToUpdate), value.Type().Kind())
		}
	}

	err := validateGroup(groups, groupSize, groupErrors)
	if err != nil {
		return err
	}

	return nil
}

func setStructValueByJson(fv reflect.Value, jsonKey string, jsonValue interface{}) error {
	if fv.IsValid() && fv.CanSet() {
		switch fv.Kind() {
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			var newInt int64
			if fl, ok := jsonValue.(float64); ok {
				// This case is for the case that json.Unmarshal unmarshals an int value into a float64 value.
				newInt = int64(fl)
			} else if _, ok := jsonValue.(int); ok {
				newInt = int64(jsonValue.(int))
			} else if _, ok := jsonValue.(int64); ok {
				newInt = int64(jsonValue.(int64))
			} else if _, ok := jsonValue.(int32); ok {
				newInt = int64(jsonValue.(int32))
			} else if _, ok := jsonValue.(int16); ok {
				newInt = int64(jsonValue.(int16))
			} else if _, ok := jsonValue.(int8); ok {
				newInt = int64(jsonValue.(int8))
			} else {
				return fmt.Errorf("input value has to be of type %v, was %v", fv.Kind(), reflect.ValueOf(jsonValue).Kind())
			}

			if fv.OverflowInt(newInt) {
				return fmt.Errorf("cannot set overflowing int for field %v", jsonKey)
			}
			fv.SetInt(newInt)
		case reflect.Float64, reflect.Float32:
			var newFloat float64
			if fl, ok := jsonValue.(float64); ok {
				newFloat = float64(fl)
			} else if fl, ok := jsonValue.(float32); ok {
				newFloat = float64(fl)
			} else {
				return fmt.Errorf("input value has to be of type %v, was %v", fv.Kind(), reflect.ValueOf(jsonValue).Kind())
			}

			if fv.OverflowFloat(newFloat) {
				return fmt.Errorf("cannot set overflowing float for field %v", jsonKey)
			}
			fv.SetFloat(newFloat)
		case reflect.String:
			if _, ok := jsonValue.(string); !ok {
				return fmt.Errorf("input value has to be of type %v, was %v", fv.Kind(), reflect.ValueOf(jsonValue).Kind())
			}

			fv.SetString(string(jsonValue.(string)))
		case reflect.Bool:
			if _, ok := jsonValue.(bool); !ok {
				return fmt.Errorf("input value has to be of type %v, was %v", fv.Kind(), reflect.ValueOf(jsonValue).Kind())
			}

			fv.SetBool(bool(jsonValue.(bool)))
		case reflect.Array, reflect.Slice:
			if reflect.TypeOf(jsonValue).Kind() != reflect.Array && reflect.TypeOf(jsonValue).Kind() != reflect.Slice {
				return fmt.Errorf("input value has to be of type %v or %v, was %v of %v", reflect.Array, reflect.Slice, reflect.ValueOf(jsonValue).Kind(), reflect.TypeOf(jsonValue).Elem().Kind())
			}

			switch t := reflect.TypeOf(fv.Interface()).Elem().Kind(); t {
			case reflect.Int:
				if _, ok := jsonValue.([]interface{}); ok {
					typedArray, err := ArrayOfInterfaceToArrayOf[int](jsonValue.([]interface{}))
					if err != nil {
						return err
					}
					reflect.Copy(fv, reflect.ValueOf(typedArray))
				} else if _, ok := jsonValue.([]int); ok {
					reflect.Copy(fv, reflect.ValueOf(jsonValue.([]int)))
				} else {
					return fmt.Errorf("input value has to be of type %v, was %v", t, reflect.TypeOf(jsonValue).Elem().Kind())
				}
			case reflect.Int64:
				if _, ok := jsonValue.([]interface{}); ok {
					typedArray, err := ArrayOfInterfaceToArrayOf[int64](jsonValue.([]interface{}))
					if err != nil {
						return err
					}
					reflect.Copy(fv, reflect.ValueOf(typedArray))
				} else if _, ok := jsonValue.([]int64); ok {
					reflect.Copy(fv, reflect.ValueOf(jsonValue.([]int64)))
				} else {
					return fmt.Errorf("input value has to be of type %v, was %v", t, reflect.TypeOf(jsonValue).Elem().Kind())
				}
			case reflect.Int32:
				if _, ok := jsonValue.([]interface{}); ok {
					typedArray, err := ArrayOfInterfaceToArrayOf[int32](jsonValue.([]interface{}))
					if err != nil {
						return err
					}
					reflect.Copy(fv, reflect.ValueOf(typedArray))
				} else if _, ok := jsonValue.([]int32); ok {
					reflect.Copy(fv, reflect.ValueOf(jsonValue.([]int32)))
				} else {
					return fmt.Errorf("input value has to be of type %v, was %v", t, reflect.TypeOf(jsonValue).Elem().Kind())
				}
			case reflect.Int16:
				if _, ok := jsonValue.([]interface{}); ok {
					typedArray, err := ArrayOfInterfaceToArrayOf[int16](jsonValue.([]interface{}))
					if err != nil {
						return err
					}
					reflect.Copy(fv, reflect.ValueOf(typedArray))
				} else if _, ok := jsonValue.([]int16); ok {
					reflect.Copy(fv, reflect.ValueOf(jsonValue.([]int16)))
				} else {
					return fmt.Errorf("input value has to be of type %v, was %v", t, reflect.TypeOf(jsonValue).Elem().Kind())
				}
			case reflect.Int8:
				if _, ok := jsonValue.([]interface{}); ok {
					typedArray, err := ArrayOfInterfaceToArrayOf[int8](jsonValue.([]interface{}))
					if err != nil {
						return err
					}
					reflect.Copy(fv, reflect.ValueOf(typedArray))
				} else if _, ok := jsonValue.([]int8); ok {
					reflect.Copy(fv, reflect.ValueOf(jsonValue.([]int8)))
				} else {
					return fmt.Errorf("input value has to be of type %v, was %v", t, reflect.TypeOf(jsonValue).Elem().Kind())
				}
			case reflect.Float64:
				if _, ok := jsonValue.([]interface{}); ok {
					typedArray, err := ArrayOfInterfaceToArrayOf[float64](jsonValue.([]interface{}))
					if err != nil {
						return err
					}
					reflect.Copy(fv, reflect.ValueOf(typedArray))
				} else if _, ok := jsonValue.([]float64); ok {
					reflect.Copy(fv, reflect.ValueOf(jsonValue.([]float64)))
				} else {
					return fmt.Errorf("input value has to be of type %v, was %v", t, reflect.TypeOf(jsonValue).Elem().Kind())
				}
			case reflect.Float32:
				if _, ok := jsonValue.([]interface{}); ok {
					typedArray, err := ArrayOfInterfaceToArrayOf[float32](jsonValue.([]interface{}))
					if err != nil {
						return err
					}
					reflect.Copy(fv, reflect.ValueOf(typedArray))
				} else if _, ok := jsonValue.([]float32); ok {
					reflect.Copy(fv, reflect.ValueOf(jsonValue.([]float32)))
				} else {
					return fmt.Errorf("input value has to be of type %v, was %v", t, reflect.TypeOf(jsonValue).Elem().Kind())
				}
			case reflect.String:
				if _, ok := jsonValue.([]interface{}); ok {
					typedArray, err := ArrayOfInterfaceToArrayOf[string](jsonValue.([]interface{}))
					if err != nil {
						return err
					}
					reflect.Copy(fv, reflect.ValueOf(typedArray))
				} else if _, ok := jsonValue.([]string); ok {
					reflect.Copy(fv, reflect.ValueOf(jsonValue.([]string)))
				} else {
					return fmt.Errorf("input value has to be of type %v, was %v", t, reflect.TypeOf(jsonValue).Elem().Kind())
				}
			case reflect.Bool:
				if _, ok := jsonValue.([]interface{}); ok {
					typedArray, err := ArrayOfInterfaceToArrayOf[bool](jsonValue.([]interface{}))
					if err != nil {
						return err
					}
					reflect.Copy(fv, reflect.ValueOf(typedArray))
				} else if _, ok := jsonValue.([]bool); ok {
					reflect.Copy(fv, reflect.ValueOf(jsonValue.([]bool)))
				} else {
					return fmt.Errorf("input value has to be of type %v, was %v", t, reflect.TypeOf(jsonValue).Elem().Kind())
				}
			default:
				return fmt.Errorf("invalid array element type: %v", reflect.TypeOf(fv.Interface()).Elem().Kind())
			}
		default:
			return fmt.Errorf("invalid field type: %v", reflect.TypeOf(jsonValue).Elem().Kind())
		}
	}
	return nil
}
