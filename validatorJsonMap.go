package validator

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"strconv"

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
			log.Printf("Setting field %v (json key: %v) of %v to value %v", fieldType.Name, jsonKey, reflect.TypeOf(structToUpdate), jsonValue)
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

func InterfaceToType(in any, expected reflect.Type) (out any, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error converting interface to type %v: %v", expected, r)
		}
	}()

	log.Printf("Converting %v with type %T to %v", in, in, expected)
	switch expected.Kind() {
	case reflect.String:
		if v, ok := in.(string); ok {
			return v, nil
		}
	case reflect.Bool:
		if v, ok := in.(bool); ok {
			return v, nil
		} else if v, ok := in.(string); ok {
			switch v {
			// Case on and off are for form values.
			case "on":
				return any(true), nil
			case "off":
				return any(false), nil
			default:
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf("error parsing string to bool: %v", err)
				}
				return any(b), nil
			}
		}
	case reflect.Int:
		log.Printf("Converting %v with type %T to int", in, in)
		if v, ok := in.(float64); ok {
			return any(int(v)), nil
		} else if v, ok := in.(string); ok {
			i, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing string to int64: %v", err)
			}
			return any(int(i)), nil
		} else if v, ok := in.(int); ok {
			return v, nil
		} else if reflect.TypeOf(in).ConvertibleTo(expected) {
			return reflect.ValueOf(in).Convert(expected).Interface(), nil
		}
	case reflect.Int8:
		if v, ok := in.(float64); ok {
			return any(int8(v)), nil
		} else if v, ok := in.(string); ok {
			i, err := strconv.ParseInt(v, 10, 8)
			if err != nil {
				return nil, fmt.Errorf("error parsing string to int8: %v", err)
			}
			return any(int8(i)), nil
		} else if v, ok := in.(int8); ok {
			return v, nil
		} else if reflect.TypeOf(in).ConvertibleTo(expected) {
			return reflect.ValueOf(in).Convert(expected).Interface(), nil
		}
	case reflect.Int16:
		if v, ok := in.(float64); ok {
			return any(int16(v)), nil
		} else if v, ok := in.(string); ok {
			i, err := strconv.ParseInt(v, 10, 16)
			if err != nil {
				return nil, fmt.Errorf("error parsing string to int16: %v", err)
			}
			return any(int16(i)), nil
		} else if v, ok := in.(int16); ok {
			return v, nil
		} else if reflect.TypeOf(in).ConvertibleTo(expected) {
			return reflect.ValueOf(in).Convert(expected).Interface(), nil
		}
	case reflect.Int32:
		if v, ok := in.(float64); ok {
			return any(int32(v)), nil
		} else if v, ok := in.(string); ok {
			i, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("error parsing string to int32: %v", err)
			}
			return any(int32(i)), nil
		} else if v, ok := in.(int32); ok {
			return v, nil
		} else if reflect.TypeOf(in).ConvertibleTo(expected) {
			return reflect.ValueOf(in).Convert(expected).Interface(), nil
		}
	case reflect.Int64:
		if v, ok := in.(float64); ok {
			return any(int64(v)), nil
		} else if v, ok := in.(string); ok {
			i, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing string to int64: %v", err)
			}
			return any(int64(i)), nil
		} else if v, ok := in.(int64); ok {
			return v, nil
		} else if reflect.TypeOf(in).ConvertibleTo(expected) {
			return reflect.ValueOf(in).Convert(expected).Interface(), nil
		}
	case reflect.Uint:
		if v, ok := in.(float64); ok {
			return any(uint(v)), nil
		} else if v, ok := in.(string); ok {
			u, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing string to uint64: %v", err)
			}
			return any(uint(u)), nil
		} else if v, ok := in.(uint); ok {
			return v, nil
		} else if reflect.TypeOf(in).ConvertibleTo(expected) {
			return reflect.ValueOf(in).Convert(expected).Interface(), nil
		}
	case reflect.Uint8:
		if v, ok := in.(float64); ok {
			return any(uint8(v)), nil
		} else if v, ok := in.(string); ok {
			u, err := strconv.ParseUint(v, 10, 8)
			if err != nil {
				return nil, fmt.Errorf("error parsing string to uint8: %v", err)
			}
			return any(uint8(u)), nil
		} else if v, ok := in.(uint8); ok {
			return v, nil
		} else if reflect.TypeOf(in).ConvertibleTo(expected) {
			return reflect.ValueOf(in).Convert(expected).Interface(), nil
		}
	case reflect.Uint16:
		if v, ok := in.(float64); ok {
			return any(uint16(v)), nil
		} else if v, ok := in.(string); ok {
			u, err := strconv.ParseUint(v, 10, 16)
			if err != nil {
				return nil, fmt.Errorf("error parsing string to uint16: %v", err)
			}
			return any(uint16(u)), nil
		} else if v, ok := in.(uint16); ok {
			return v, nil
		} else if reflect.TypeOf(in).ConvertibleTo(expected) {
			return reflect.ValueOf(in).Convert(expected).Interface(), nil
		}
	case reflect.Uint32:
		if v, ok := in.(float64); ok {
			return any(uint32(v)), nil
		} else if v, ok := in.(string); ok {
			u, err := strconv.ParseUint(v, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("error parsing string to uint32: %v", err)
			}
			return any(uint32(u)), nil
		} else if v, ok := in.(uint32); ok {
			return v, nil
		} else if reflect.TypeOf(in).ConvertibleTo(expected) {
			return reflect.ValueOf(in).Convert(expected).Interface(), nil
		}
	case reflect.Uint64:
		if v, ok := in.(float64); ok {
			return any(uint64(v)), nil
		} else if v, ok := in.(string); ok {
			u, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing string to uint64: %v", err)
			}
			return u, nil
		} else if v, ok := in.(uint64); ok {
			return v, nil
		} else if reflect.TypeOf(in).ConvertibleTo(expected) {
			return reflect.ValueOf(in).Convert(expected).Interface(), nil
		}
	case reflect.Float32:
		if v, ok := in.(float64); ok {
			return any(float32(v)), nil
		} else if v, ok := in.(string); ok {
			f, err := strconv.ParseFloat(v, 32)
			if err != nil {
				return nil, fmt.Errorf("error parsing string to float32: %v", err)
			}
			return any(float32(f)), nil
		} else if v, ok := in.(float32); ok {
			return v, nil
		} else if reflect.TypeOf(in).ConvertibleTo(expected) {
			return reflect.ValueOf(in).Convert(expected).Interface(), nil
		}
	case reflect.Float64:
		if v, ok := in.(float64); ok {
			return v, nil
		} else if v, ok := in.(string); ok {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing string to float64: %v", err)
			}
			return any(float64(f)), nil
		} else if reflect.TypeOf(in).ConvertibleTo(expected) {
			return reflect.ValueOf(in).Convert(expected).Interface(), nil
		}
	case reflect.Interface:
		return in, nil
	default:
		return nil, fmt.Errorf("unsupported type %T", expected)
	}
	return nil, fmt.Errorf("unsupported type %T", expected)
}

func JsonMapToMapKV(m map[string]any, expectedKey reflect.Type, expectedValue reflect.Type) (reflect.Value, error) {
	targetMapValue := reflect.MakeMap(reflect.MapOf(expectedKey, expectedValue))
	for key, value := range m {
		valueConverted, err := InterfaceToType(value, expectedValue)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("error converting value for key %s: %v", key, err)
		}
		keyConverted, err := InterfaceToType(key, expectedKey)
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
		itemConverted, err := InterfaceToType(item, expectedValue)
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
			b, err := InterfaceToType(jsonValue, reflect.TypeOf(newString))
			if err != nil {
				return fmt.Errorf("error converting value to string: %v", err)
			}
			newString = b.(string)
			fv.SetString(newString)
		case reflect.Bool:
			var newBool bool = true
			b, err := InterfaceToType(jsonValue, reflect.TypeOf(newBool))
			if err != nil {
				return fmt.Errorf("error converting value to bool: %v", err)
			}
			newBool = b.(bool)
			fv.SetBool(newBool)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			log.Printf("Setting field %v to int value %v", fv.Type(), jsonValue)
			var newInt int64 = 0
			i, err := InterfaceToType(jsonValue, reflect.TypeOf(newInt))
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
			i, err := InterfaceToType(jsonValue, reflect.TypeOf(newUint))
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
			i, err := InterfaceToType(jsonValue, reflect.TypeOf(newFloat))
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
			var mapReflect reflect.Value
			if v, ok := jsonValue.(map[string]any); ok {
				log.Printf("Converting json map %v to mapKV with key type %v and value type %v", v, fv.Type().Key(), fv.Type().Elem())
				mapReflect, err = JsonMapToMapKV(v, fv.Type().Key(), fv.Type().Elem())
				if err != nil {
					return fmt.Errorf("error converting json map to mapKV: %v", err)
				}
				log.Printf("Converted json map to mapKV: %v, %v", mapReflect.Type(), mapReflect.Type().Elem())
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
