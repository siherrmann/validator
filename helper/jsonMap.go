package helper

import (
	"encoding"
	"fmt"
	"reflect"
)

func GetValidMap(in any) (map[string]any, error) {
	v, ok := in.(map[string]any)
	if ok {
		return v, nil
	}
	return nil, fmt.Errorf("error getting valid map from json")
}

// UnmapStructToJsonMap writes the exported fields of the given struct pointer into the JsonMap.
// Nested structs (also behind pointers) and slices/arrays of structs are converted recursively
// to JsonMaps (and []any of JsonMaps), so their inner validations can be run like for json input.
// Fields tagged with `inline` are flattened into the parent map, unexported fields are skipped.
// Structs implementing encoding.TextMarshaler (eg. time.Time) are kept as they are.
func UnmapStructToJsonMap(structInput any, jsonMapToUpdate *map[string]any) error {
	err := CheckValidPointerToStruct(structInput)
	if err != nil {
		return err
	}

	unmapStructValue(reflect.ValueOf(structInput).Elem(), *jsonMapToUpdate)
	return nil
}

// unmapInline returns the JsonMap of an inline struct (or pointer to struct) field.
// ok is false if the field is not a struct, a nil pointer has no fields.
func unmapInline(field reflect.Value) (map[string]any, bool) {
	if field.Kind() == reflect.Ptr {
		if field.Type().Elem().Kind() != reflect.Struct {
			return nil, false
		}
		if field.IsNil() {
			return nil, true
		}
		field = field.Elem()
	}
	if field.Kind() != reflect.Struct {
		return nil, false
	}
	inner := map[string]any{}
	unmapStructValue(field, inner)
	return inner, true
}

var textMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()

// unmapStructValue writes the exported fields of the struct value into the JsonMap.
func unmapStructValue(structValue reflect.Value, jsonMap map[string]any) {
	structType := structValue.Type()
	for i := 0; i < structType.NumField(); i++ {
		fieldType := structType.Field(i)
		if !fieldType.IsExported() {
			continue
		}
		field := structValue.Field(i)

		fieldKey, _, isInline := FieldKey(fieldType)
		if isInline {
			if inner, ok := unmapInline(field); ok {
				for k, v := range inner {
					jsonMap[k] = v
				}
				continue
			}
		}

		jsonMap[fieldKey] = unmapValue(field)
	}
}

// unmapValue converts a struct (or pointer to struct) to a JsonMap
// and a slice/array of structs to a []any of JsonMaps. All other values are returned as they are.
func unmapValue(value reflect.Value) any {
	switch value.Kind() {
	case reflect.Ptr:
		if value.Type().Elem().Kind() == reflect.Struct && !isTextMarshalerStruct(value.Type().Elem()) {
			if value.IsNil() {
				return nil
			}
			return unmapValue(value.Elem())
		}
	case reflect.Struct:
		if !isTextMarshalerStruct(value.Type()) {
			jsonMap := map[string]any{}
			unmapStructValue(value, jsonMap)
			return jsonMap
		}
	case reflect.Slice, reflect.Array:
		if value.Type().Elem().Kind() == reflect.Struct && !isTextMarshalerStruct(value.Type().Elem()) {
			if value.Kind() == reflect.Slice && value.IsNil() {
				return nil
			}
			array := make([]any, value.Len())
			for i := 0; i < value.Len(); i++ {
				array[i] = unmapValue(value.Index(i))
			}
			return array
		}
	}
	return value.Interface()
}

// isTextMarshalerStruct reports if the struct type (or its pointer) implements encoding.TextMarshaler.
func isTextMarshalerStruct(structType reflect.Type) bool {
	return structType.Implements(textMarshalerType) || reflect.PointerTo(structType).Implements(textMarshalerType)
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

		fieldKey, _, isInline := FieldKey(fieldType)

		if isInline {
			target := field
			if field.Kind() == reflect.Ptr {
				if field.IsNil() {
					if field.CanSet() {
						field.Set(reflect.New(field.Type().Elem()))
					}
				}
				target = field.Elem()
			}
			if target.CanAddr() {
				err := MapJsonMapToStruct(jsonMapInput, target.Addr().Interface())
				if err != nil {
					return fmt.Errorf("error mapping inline struct %v: %v", fieldType.Name, err.Error())
				}
			} else if field.Kind() == reflect.Ptr {
				err := MapJsonMapToStruct(jsonMapInput, field.Interface())
				if err != nil {
					return fmt.Errorf("error mapping inline struct %v: %v", fieldType.Name, err.Error())
				}
			}
			continue
		}

		if jsonValue, ok := jsonMapInput[fieldKey]; ok {
			if jsonValueMap, isMap := jsonValue.(map[string]any); isMap && (field.Kind() == reflect.Struct || (field.Kind() == reflect.Ptr && field.Elem().Kind() == reflect.Struct)) {
				target := field
				if field.Kind() == reflect.Ptr {
					if field.IsNil() {
						if field.CanSet() {
							field.Set(reflect.New(field.Type().Elem()))
						}
					}
					target = field.Elem()
				}
				if target.CanAddr() {
					err := MapJsonMapToStruct(jsonValueMap, target.Addr().Interface())
					if err != nil {
						return fmt.Errorf("error mapping nested struct %v: %v", fieldType.Name, err.Error())
					}
					continue
				} else if field.Kind() == reflect.Ptr {
					err := MapJsonMapToStruct(jsonValueMap, field.Interface())
					if err != nil {
						return fmt.Errorf("error mapping nested struct %v: %v", fieldType.Name, err.Error())
					}
					continue
				}
			}

			err := SetStructValueByJson(field, jsonValue)
			if err != nil {
				return fmt.Errorf("could not set field %v (json key: %v) of %v: %v", fieldType.Name, fieldKey, reflect.TypeOf(structToUpdate), err.Error())
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
