package helper

import (
	"fmt"
	"reflect"
	"strings"
)

func GetValidMap(in any) (map[string]any, error) {
	v, ok := in.(map[string]any)
	if ok {
		return v, nil
	}
	return nil, fmt.Errorf("error getting valid map from json")
}

func UnmapStructToJsonMap(structInput any, jsonMapToUpdate *map[string]any) error {
	err := CheckValidPointerToStruct(structInput)
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

func JsonMapToMapKV(m map[string]any, expectedKey reflect.Type, expectedValue reflect.Type) (reflect.Value, error) {
	targetMapValue := reflect.MakeMap(reflect.MapOf(expectedKey, expectedValue))
	for key, value := range m {
		valueConverted, err := AnyToType(value, expectedValue)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("error converting value for key %s: %v", key, err)
		}
		keyConverted, err := AnyToType(key, expectedKey)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("error converting key %s: %v", key, err)
		}
		targetMapValue.SetMapIndex(reflect.ValueOf(keyConverted), reflect.ValueOf(valueConverted))
	}

	return targetMapValue, nil
}

func MapJsonMapToStruct(jsonMapInput map[string]any, structToUpdate any) error {
	err := CheckValidPointerToStruct(structToUpdate)
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
			// Split on comma to handle omitempty and other options
			jsonKey = strings.Split(jsonKey, ",")[0]
			if jsonKey != "-" {
				fieldKey = jsonKey
			}
		}

		if jsonValue, ok := jsonMapInput[fieldKey]; ok {
			err := SetStructValueByJson(field, jsonValue)
			if err != nil {
				return fmt.Errorf("could not set field %v (json key: %v) of %v: %v", fieldType.Name, jsonKey, reflect.TypeOf(structToUpdate), err.Error())
			}
		} else {
			// Initialize nil map and slice fields with empty collections to prevent panics
			if field.CanSet() {
				if field.Kind() == reflect.Map && field.IsNil() {
					field.Set(reflect.MakeMap(field.Type()))
				} else if field.Kind() == reflect.Slice && field.IsNil() {
					field.Set(reflect.MakeSlice(field.Type(), 0, 0))
				}
			}
		}
	}
	return nil
}

func SetStructValueByJson(fv reflect.Value, jsonValue any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error setting struct value: %v", r)
		}
	}()

	if fv.IsValid() && fv.CanSet() {
		converted, err := AnyToType(jsonValue, fv.Type())
		if err != nil {
			return err
		}
		fv.Set(reflect.ValueOf(converted))
	}
	return nil
}
