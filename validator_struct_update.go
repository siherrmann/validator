package validator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"slices"
	"strings"

	"github.com/siherrmann/validator/helper"
	"github.com/siherrmann/validator/model"
	"github.com/siherrmann/validator/validators"
)

// UnmapOrAnmarshalValidateAndUpdate unmarshals given json ([]byte) or given url.Values (from request.Form),
// validates them and updates the given struct.
func UnmapOrUnmarshalRequestValidateAndUpdate(request *http.Request, structToUpdate interface{}) error {
	err := request.ParseForm()
	if err != nil {
		return err
	}

	if len(request.Form.Encode()) > 0 {
		err = UnmapValidateAndUpdate(request.Form, structToUpdate)
	} else {
		var bodyBytes []byte
		bodyBytes, err = io.ReadAll(request.Body)
		if err != nil {
			return err
		}
		err = UnmarshalValidateAndUpdate(bodyBytes, structToUpdate)
	}

	return err
}

// UnmarshalValidateAndUpdate unmarshals given json ([]byte) into pointer v.
// For more information to ValidateAndUpdate look at ValidateAndUpdate(jsonInput model.JsonMap, structToUpdate interface{}) error.
func UnmarshalValidateAndUpdate(jsonInput []byte, structToUpdate interface{}) error {
	jsonUnmarshaled := model.JsonMap{}

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

// UnmapValidateAndUpdate unmaps given url.Values into pointer jsonMap.
// For more information to ValidateAndUpdate look at ValidateAndUpdate(jsonInput model.JsonMap, structToUpdate interface{}) error.
func UnmapValidateAndUpdate(values url.Values, structToUpdate interface{}) error {
	mapOut, err := UnmapUrlValuesToJsonMap(values)
	if err != nil {
		return err
	}

	err = ValidateAndUpdate(mapOut, structToUpdate)
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
func ValidateAndUpdate(jsonInput model.JsonMap, structToUpdate interface{}) error {
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

	keys := []string{}
	groups := map[string]*model.Group{}
	groupSize := map[string]int{}
	groupErrors := map[string][]error{}

	for i := 0; i < structFull.Type().NumField(); i++ {
		tag := structFull.Type().Field(i).Tag.Get(string(model.UPD))
		field := structFull.Field(i)
		fieldName := structFull.Type().Field(i).Name

		validation := &model.Validation{}
		err := validation.Fill(tag, model.UPD, field)
		if err != nil {
			return err
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

		var ok bool
		var jsonValue interface{}
		if jsonValue, ok = jsonInput[validation.Key]; !ok {
			if strings.TrimSpace(validation.Requirement) == string(model.NONE) {
				continue
			} else if len(validation.Groups) == 0 {
				return fmt.Errorf("json %v key not in map", validation.Key)
			} else {
				for _, group := range validation.Groups {
					groupErrors[group.Name] = append(groupErrors[group.Name], fmt.Errorf("json %v key not in map", validation.Key))
				}
				continue
			}
		}

		var validatedValue interface{}
		if validation.Type == model.Struct {
			var validMap interface{}
			validMap, err = validation.GetValidValue(jsonValue)
			if err != nil {
				return err
			}

			err = ValidateAndUpdate(validMap.(map[string]interface{}), field.Addr().Interface())
			validatedValue = field.Interface()
		} else {
			validatedValue, err = ValidateValueWithParser(reflect.ValueOf(jsonValue), validation)
		}

		if err != nil && len(validation.Groups) == 0 {
			return fmt.Errorf("field %v of %v invalid: %v", fieldName, reflect.TypeOf(structToUpdate), err.Error())
		} else if err != nil {
			for _, group := range validation.Groups {
				groupErrors[group.Name] = append(groupErrors[group.Name], fmt.Errorf("field %v of %v invalid: %v", fieldName, reflect.TypeOf(structToUpdate), err.Error()))
			}
			continue
		}

		err = setStructValueByJson(field, validatedValue)
		if err != nil && len(validation.Groups) == 0 {
			return fmt.Errorf("could not set field %v (json key: %v) of %v: %v", fieldName, validation.Key, reflect.TypeOf(structToUpdate), err.Error())
		} else if err != nil {
			for _, group := range groups {
				groupErrors[group.Name] = append(groupErrors[group.Name], fmt.Errorf("could not set field %v: %v", fieldName, err.Error()))
			}
			continue
		}
	}

	err := validators.ValidateGroups(groups, groupSize, groupErrors)
	if err != nil {
		return err
	}

	return nil
}

func setStructValueByJson(fv reflect.Value, jsonValue interface{}) error {
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
				return fmt.Errorf("cannot set overflowing int")
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
				return fmt.Errorf("cannot set overflowing float")
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
		case reflect.Struct:
			if v, ok := jsonValue.(string); ok {
				validation := &model.Validation{Type: model.Time}
				date, err := validation.InterfaceFromString(v)
				if err != nil {
					return err
				}
				fv.Set(reflect.ValueOf(date))
			} else {
				fv.Set(reflect.ValueOf(jsonValue))
			}
		case reflect.Map:
			if _, ok := jsonValue.(map[string]interface{}); !ok {
				return fmt.Errorf("input value has to be of type %v, was %v", fv.Kind(), reflect.ValueOf(jsonValue).Kind())
			}
			fv.Set(reflect.ValueOf(jsonValue))
		case reflect.Array, reflect.Slice:
			if reflect.TypeOf(jsonValue).Kind() != reflect.Array && reflect.TypeOf(jsonValue).Kind() != reflect.Slice {
				return fmt.Errorf("input value has to be of type %v or %v, was %v of %v", reflect.Array, reflect.Slice, reflect.ValueOf(jsonValue).Kind(), reflect.TypeOf(jsonValue).Elem().Kind())
			}

			switch t := reflect.TypeOf(fv.Interface()).Elem().Kind(); t {
			case reflect.Int:
				if v, ok := jsonValue.([]interface{}); ok {
					typedArray, err := helper.ArrayOfInterfaceToArrayOf[int](v)
					if err != nil {
						return err
					}
					fv.Set(reflect.ValueOf(typedArray))
				} else if _, ok := jsonValue.([]int); ok {
					fv.Set(reflect.ValueOf(jsonValue.([]int)))
				} else {
					return fmt.Errorf("input value has to be of type %v, was %v", t, reflect.TypeOf(jsonValue).Elem().Kind())
				}
			case reflect.Int64:
				if v, ok := jsonValue.([]interface{}); ok {
					typedArray, err := helper.ArrayOfInterfaceToArrayOf[int64](v)
					if err != nil {
						return err
					}
					fv.Set(reflect.ValueOf(typedArray))
				} else if _, ok := jsonValue.([]int64); ok {
					fv.Set(reflect.ValueOf(jsonValue.([]int64)))
				} else {
					return fmt.Errorf("input value has to be of type %v, was %v", t, reflect.TypeOf(jsonValue).Elem().Kind())
				}
			case reflect.Int32:
				if v, ok := jsonValue.([]interface{}); ok {
					typedArray, err := helper.ArrayOfInterfaceToArrayOf[int32](v)
					if err != nil {
						return err
					}
					fv.Set(reflect.ValueOf(typedArray))
				} else if _, ok := jsonValue.([]int32); ok {
					fv.Set(reflect.ValueOf(jsonValue.([]int32)))
				} else {
					return fmt.Errorf("input value has to be of type %v, was %v", t, reflect.TypeOf(jsonValue).Elem().Kind())
				}
			case reflect.Int16:
				if v, ok := jsonValue.([]interface{}); ok {
					typedArray, err := helper.ArrayOfInterfaceToArrayOf[int16](v)
					if err != nil {
						return err
					}
					fv.Set(reflect.ValueOf(typedArray))
				} else if _, ok := jsonValue.([]int16); ok {
					fv.Set(reflect.ValueOf(jsonValue.([]int16)))
				} else {
					return fmt.Errorf("input value has to be of type %v, was %v", t, reflect.TypeOf(jsonValue).Elem().Kind())
				}
			case reflect.Int8:
				if v, ok := jsonValue.([]interface{}); ok {
					typedArray, err := helper.ArrayOfInterfaceToArrayOf[int8](v)
					if err != nil {
						return err
					}
					fv.Set(reflect.ValueOf(typedArray))
				} else if _, ok := jsonValue.([]int8); ok {
					fv.Set(reflect.ValueOf(jsonValue.([]int8)))
				} else {
					return fmt.Errorf("input value has to be of type %v, was %v", t, reflect.TypeOf(jsonValue).Elem().Kind())
				}
			case reflect.Float64:
				if v, ok := jsonValue.([]interface{}); ok {
					typedArray, err := helper.ArrayOfInterfaceToArrayOf[float64](v)
					if err != nil {
						return err
					}
					fv.Set(reflect.ValueOf(typedArray))
				} else if _, ok := jsonValue.([]float64); ok {
					fv.Set(reflect.ValueOf(jsonValue.([]float64)))
				} else {
					return fmt.Errorf("input value has to be of type %v, was %v", t, reflect.TypeOf(jsonValue).Elem().Kind())
				}
			case reflect.Float32:
				if v, ok := jsonValue.([]interface{}); ok {
					typedArray, err := helper.ArrayOfInterfaceToArrayOf[float32](v)
					if err != nil {
						return err
					}
					fv.Set(reflect.ValueOf(typedArray))
				} else if _, ok := jsonValue.([]float32); ok {
					fv.Set(reflect.ValueOf(jsonValue.([]float32)))
				} else {
					return fmt.Errorf("input value has to be of type %v, was %v", t, reflect.TypeOf(jsonValue).Elem().Kind())
				}
			case reflect.String:
				if v, ok := jsonValue.([]interface{}); ok {
					typedArray, err := helper.ArrayOfInterfaceToArrayOf[string](v)
					if err != nil {
						return err
					}
					fv.Set(reflect.ValueOf(typedArray))
				} else if _, ok := jsonValue.([]string); ok {
					fv.Set(reflect.ValueOf(jsonValue.([]string)))
				} else {
					return fmt.Errorf("input value has to be of type %v, was %v", t, reflect.TypeOf(jsonValue).Elem().Kind())
				}
			case reflect.Bool:
				if v, ok := jsonValue.([]interface{}); ok {
					typedArray, err := helper.ArrayOfInterfaceToArrayOf[bool](v)
					if err != nil {
						return err
					}
					fv.Set(reflect.ValueOf(typedArray))
				} else if _, ok := jsonValue.([]bool); ok {
					fv.Set(reflect.ValueOf(jsonValue.([]bool)))
				} else {
					return fmt.Errorf("input value has to be of type %v, was %v", t, reflect.TypeOf(jsonValue).Elem().Kind())
				}
			case reflect.Struct:
				if a, ok := jsonValue.([]interface{}); ok {
					typedArray := reflect.New(fv.Type())
					for _, v := range a {
						if m, ok := v.(map[string]interface{}); ok {
							structTempt := reflect.New(fv.Type().Elem()).Interface()
							err := ValidateAndUpdate(m, structTempt)
							if err != nil {
								return err
							}
							typedArray = reflect.Append(typedArray.Elem(), reflect.ValueOf(structTempt).Elem())
						} else {
							return fmt.Errorf("input value inside array has to be of type map[string]interface{}, was %v", reflect.TypeOf(v))
						}
					}
					fv.Set(typedArray)
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
