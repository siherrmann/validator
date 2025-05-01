package validator

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"reflect"

	"github.com/siherrmann/validator/helper"
	"github.com/siherrmann/validator/model"
)

func GetValidMap(in interface{}) (model.JsonMap, error) {
	v, ok := in.(model.JsonMap)
	if !ok {
		v, ok = in.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("error getting valid map from json")
		} else {
			return model.JsonMap(v), nil
		}
	}
	return v, nil
}

func UnmarshalJsonToJsonMap(jsonInput []byte) (model.JsonMap, error) {
	mapOut := model.JsonMap{}
	err := json.Unmarshal(jsonInput, &mapOut)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling: %v", err)
	}
	return mapOut, nil
}

func UnmapUrlValuesToJsonMap(values url.Values) (model.JsonMap, error) {
	mapOut := model.JsonMap{}
	for k := range values {
		if len(values[k]) > 1 {
			arrayOut := []interface{}{}
			for _, v := range values[k] {
				var unmarshalled interface{}
				err := json.Unmarshal([]byte(v), &unmarshalled)
				if err == nil {
					arrayOut = append(arrayOut, unmarshalled)
				} else {
					arrayOut = append(arrayOut, v)
				}
			}
			mapOut[k] = arrayOut
		} else {
			value := values.Get(k)
			var unmarshalled interface{}
			err := json.Unmarshal([]byte(value), &unmarshalled)
			if err == nil {
				mapOut[k] = unmarshalled
			} else {
				mapOut[k] = value
			}
		}
	}
	return mapOut, nil
}

func MapJsonMapToStruct(jsonMapInput model.JsonMap, structToUpdate interface{}) error {
	err := helper.CheckValidPointerToStruct(structToUpdate)
	if err != nil {
		return err
	}

	structFull := reflect.ValueOf(structToUpdate).Elem()
	for i := 0; i < structFull.Type().NumField(); i++ {
		field := structFull.Field(i)
		fieldType := structFull.Type().Field(i)

		fieldKey := fieldType.Name
		jsonKey := fieldType.Tag.Get("json")
		if len(jsonKey) > 0 {
			fieldKey = jsonKey
		}

		if jsonValue, ok := jsonMapInput[fieldKey]; ok {
			log.Printf("json key: %v, field key: %v", jsonKey, fieldKey)
			err := SetStructValueByJson(field, jsonValue)
			if err != nil {
				return fmt.Errorf("could not set field %v (json key: %v) of %v: %v", fieldType.Name, jsonKey, reflect.TypeOf(structToUpdate), err.Error())
			}
		}
	}
	return nil
}

func UnmapStructToJsonMap(structInput interface{}, jsonMapToUpdate *model.JsonMap) error {
	err := helper.CheckValidPointerToStruct(structInput)
	if err != nil {
		return err
	}

	structFull := reflect.ValueOf(structInput).Elem()
	for i := 0; i < structFull.Type().NumField(); i++ {
		field := structFull.Field(i)
		fieldType := structFull.Type().Field(i)

		fieldKey := fieldType.Name
		jsonKey := fieldType.Tag.Get("json")
		if len(jsonKey) > 0 {
			fieldKey = jsonKey
		}

		(*jsonMapToUpdate)[fieldKey] = field.Interface()
	}
	return nil
}

func UpdateJsonMap(validatedValues model.JsonMap, jsonMapToUpdate *model.JsonMap) error {
	for k, v := range validatedValues {
		(*jsonMapToUpdate)[k] = v
	}
	return nil
}

func SetStructValueByJson(fv reflect.Value, jsonValue interface{}) error {
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
				structTemp := reflect.New(fv.Type()).Interface()
				validMap, err := GetValidMap(jsonValue)
				if err != nil {
					return fmt.Errorf("error getting valid map for struct %v: %v", structTemp, err)
				}
				err = MapJsonMapToStruct(validMap, structTemp)
				if err != nil {
					return fmt.Errorf("error setting struct value: %v", err)
				}
				fv.Set(reflect.ValueOf(structTemp).Elem())
			}
		case reflect.Map:
			if _, ok := jsonValue.(map[string]interface{}); !ok {
				return fmt.Errorf("input value has to be of type %v, was %v", fv.Kind(), reflect.ValueOf(jsonValue).Kind())
			}
			fv.Set(reflect.ValueOf(jsonValue))
		case reflect.Array, reflect.Slice:
			if !helper.IsArray(jsonValue) {
				return fmt.Errorf("input value has to be of type %v or %v, was %v", reflect.Array, reflect.Slice, reflect.ValueOf(jsonValue).Kind())
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
					underlying := fv.Type().Elem()
					typedArray := reflect.New(reflect.SliceOf(underlying)).Elem()
					for _, v := range a {
						if m, ok := v.(map[string]interface{}); ok {
							structTempt := reflect.New(underlying).Interface()
							err := MapJsonMapToStruct(m, structTempt)
							if err != nil {
								return err
							}
							typedArray = reflect.Append(typedArray, reflect.ValueOf(structTempt).Elem())
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
