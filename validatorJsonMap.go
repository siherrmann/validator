package validator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/siherrmann/validator/helper"
	"github.com/siherrmann/validator/model"
)

func GetValidMap(in any) (model.JsonMap, error) {
	v, ok := in.(model.JsonMap)
	if !ok {
		v, ok = in.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("error getting valid map from json")
		} else {
			return model.JsonMap(v), nil
		}
	}
	return v, nil
}

func UnmarshalRequestToJsonMap(request *http.Request) (model.JsonMap, error) {
	if request == nil {
		return nil, fmt.Errorf("request is nil")
	}

	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading request body: %v", err)
	}
	defer request.Body.Close()

	return UnmarshalJsonToJsonMap(bodyBytes)
}

func UnmapRequestToJsonMap(request *http.Request) (model.JsonMap, error) {
	if request == nil {
		return nil, fmt.Errorf("request is nil")
	}

	err := request.ParseForm()
	if err != nil {
		return nil, fmt.Errorf("error parsing form: %v", err)
	}

	return UnmapUrlValuesToJsonMap(request.Form)
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
			arrayOut := []any{}
			for _, v := range values[k] {
				var unmarshalled any
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
			var unmarshalled any
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

func MapJsonMapToStruct(jsonMapInput model.JsonMap, structToUpdate any) error {
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
			err := SetStructValueByJson(field, jsonValue)
			if err != nil {
				return fmt.Errorf("could not set field %v (json key: %v) of %v: %v", fieldType.Name, jsonKey, reflect.TypeOf(structToUpdate), err.Error())
			}
		}
	}
	return nil
}

func UnmapStructToJsonMap(structInput any, jsonMapToUpdate *model.JsonMap) error {
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

func UpdateJsonMap(validatedValues model.JsonMap, jsonMapToUpdate *model.JsonMap) {
	for k, v := range validatedValues {
		(*jsonMapToUpdate)[k] = v
	}
}

func JsonMapToMapKV(m map[string]any, expectedKey reflect.Type, expectedValue reflect.Type) (reflect.Value, error) {
	targetMapValue := reflect.MakeMap(reflect.MapOf(expectedKey, expectedValue))
	for key, value := range m {
		valueConverted, err := helper.AnyToType(value, expectedValue)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("error converting value for key %s: %v", key, err)
		}
		keyConverted, err := helper.AnyToType(key, expectedKey)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("error converting key %s: %v", key, err)
		}
		targetMapValue.SetMapIndex(reflect.ValueOf(keyConverted), reflect.ValueOf(valueConverted))
	}

	return targetMapValue, nil
}

func JsonArrayToArrayOf(a []any, expectedValue reflect.Type) (reflect.Value, error) {
	targetArray := reflect.MakeSlice(reflect.SliceOf(expectedValue), len(a), len(a))
	for i, item := range a {
		itemConverted, err := helper.AnyToType(item, expectedValue)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("error converting item at index %d: %v", i, err)
		}
		targetArray.Index(i).Set(reflect.ValueOf(itemConverted))
	}
	return targetArray, nil
}

func SetStructValueByJson(fv reflect.Value, jsonValue any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error setting struct value: %v", r)
		}
	}()

	if fv.IsValid() && fv.CanSet() {
		switch fv.Kind() {
		case reflect.String:
			var newString string = ""
			b, err := helper.AnyToType(jsonValue, reflect.TypeOf(newString))
			if err != nil {
				return fmt.Errorf("error converting value to string: %v", err)
			}
			newString = b.(string)
			fv.SetString(newString)
		case reflect.Bool:
			var newBool bool = true
			b, err := helper.AnyToType(jsonValue, reflect.TypeOf(newBool))
			if err != nil {
				return fmt.Errorf("error converting value to bool: %v", err)
			}
			newBool = b.(bool)
			fv.SetBool(newBool)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var newInt int64 = 0
			i, err := helper.AnyToType(jsonValue, reflect.TypeOf(newInt))
			if err != nil {
				return fmt.Errorf("error converting value to int: %v", err)
			}
			newInt = i.(int64)

			if fv.OverflowInt(newInt) {
				return fmt.Errorf("cannot set overflowing int")
			}
			fv.SetInt(newInt)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			var newUint uint64 = 0
			i, err := helper.AnyToType(jsonValue, reflect.TypeOf(newUint))
			if err != nil {
				return fmt.Errorf("error converting value to uint: %v", err)
			}
			newUint = i.(uint64)

			if fv.OverflowUint(newUint) {
				return fmt.Errorf("cannot set overflowing uint")
			}
			fv.SetUint(newUint)
		case reflect.Float32, reflect.Float64:
			var newFloat float64 = 0
			i, err := helper.AnyToType(jsonValue, reflect.TypeOf(newFloat))
			if err != nil {
				return fmt.Errorf("error converting value to float: %v", err)
			}
			newFloat = i.(float64)

			if fv.OverflowFloat(newFloat) {
				return fmt.Errorf("cannot set overflowing float")
			}
			fv.SetFloat(newFloat)
		case reflect.Struct:
			if v, ok := jsonValue.(string); ok {
				date, err := helper.AnyToType(v, reflect.TypeOf(time.Time{}))
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
			var mapReflect reflect.Value
			if v, ok := jsonValue.(map[string]any); ok {
				mapReflect, err = JsonMapToMapKV(v, fv.Type().Key(), fv.Type().Elem())
				if err != nil {
					return fmt.Errorf("error converting json map to mapKV: %v", err)
				}
			} else {
				mapReflect = reflect.ValueOf(jsonValue)
			}

			if mapReflect.Type().ConvertibleTo(fv.Type()) {
				fv.Set(mapReflect.Convert(fv.Type()))
				return nil
			} else {
				return fmt.Errorf("json map %T is not convertible to type %v", jsonValue, fv.Type())
			}
		case reflect.Array, reflect.Slice:
			if !helper.IsArray(jsonValue) {
				return fmt.Errorf("input value has to be of type %v or %v, was %v", reflect.Array, reflect.Slice, reflect.ValueOf(jsonValue).Kind())
			}

			switch t := reflect.TypeOf(fv.Interface()).Elem().Kind(); t {
			case reflect.Struct:
				if a, ok := jsonValue.([]any); ok {
					underlying := fv.Type().Elem()
					typedArray := reflect.New(reflect.SliceOf(underlying)).Elem()
					for _, v := range a {
						if m, ok := v.(map[string]any); ok {
							structTempt := reflect.New(underlying).Interface()
							err := MapJsonMapToStruct(m, structTempt)
							if err != nil {
								return err
							}
							typedArray = reflect.Append(typedArray, reflect.ValueOf(structTempt).Elem())
						} else {
							return fmt.Errorf("input value inside array has to be of type map[string]any, was %v", reflect.TypeOf(v))
						}
					}
					fv.Set(typedArray)
				} else {
					return fmt.Errorf("input value has to be of type []any, was %v", reflect.TypeOf(jsonValue))
				}
			default:
				if v, ok := jsonValue.([]any); ok {
					typedArray, err := JsonArrayToArrayOf(v, fv.Type().Elem())
					if err != nil {
						return err
					}
					fv.Set(typedArray)
				} else {
					fv.Set(reflect.ValueOf(jsonValue))
				}
			}
		default:
			return fmt.Errorf("invalid field type: %v", reflect.TypeOf(jsonValue).Elem().Kind())
		}
	}
	return nil
}
