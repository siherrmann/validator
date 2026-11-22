package helper

import (
	"encoding"
	"encoding/json/v2"
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

// UnmapStructToJsonMap writes the fields of the given struct pointer into the JsonMap.
// The keys are taken from the StructModel of the struct, so embedded and inline structs are flattened.
// Nested structs (also behind pointers) and slices/arrays of structs are converted recursively
// to JsonMaps (and []any of JsonMaps), so their inner validations can be run like for json input.
// Fields promoted from a nil embedded pointer are not written, they are missing like in json.
// Structs implementing encoding.TextMarshaler (eg. time.Time) are kept as they are.
func UnmapStructToJsonMap(structInput any, jsonMapToUpdate *map[string]any) error {
	err := CheckValidPointerToStruct(structInput)
	if err != nil {
		return err
	}

	unmapStructValue(reflect.ValueOf(structInput).Elem(), *jsonMapToUpdate)
	return nil
}

var textMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()

// unmapStructValue writes the fields of the struct value into the JsonMap.
func unmapStructValue(structValue reflect.Value, jsonMap map[string]any) {
	for _, field := range StructModelOf(structValue.Type()).Fields {
		value, ok := fieldForRead(structValue, field.Index)
		if ok {
			jsonMap[field.Key] = unmapValue(value)
		}
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

// MapJsonMapToStruct writes the values of the JsonMap into the given struct pointer.
// The keys are taken from the StructModel of the struct, so embedded and inline structs are flattened
// and nil embedded pointers are set if one of their fields is written. Only keys in the JsonMap are written,
// nested JsonMaps are written into the existing nested struct. Nil maps and slices of missing keys are set to empty ones.
//
// Values are converted with AnyToType. Fields with an own UnmarshalJSON or the json options `string` or `format:`
// are written with encoding/json instead (see StructModelOf).
func MapJsonMapToStruct(jsonMapInput map[string]any, structToUpdate any) error {
	err := CheckValidPointerToStruct(structToUpdate)
	if err != nil {
		return err
	}
	return mapJsonMapToStructValue(jsonMapInput, reflect.ValueOf(structToUpdate).Elem())
}

func mapJsonMapToStructValue(jsonMapInput map[string]any, structValue reflect.Value) error {
	for _, field := range StructModelOf(structValue.Type()).Fields {
		jsonValue, ok := jsonMapInput[field.Key]
		if !ok {
			initEmpty(structValue, field)
			continue
		}

		target := fieldForWrite(structValue, field.Index)
		var err error
		if field.decodeType != nil {
			err = decodeJSON(target, field.decodeType, jsonValue)
		} else if jsonValueMap, isMap := jsonValue.(map[string]any); isMap && isStructOrPointer(target.Type()) {
			if target.Kind() == reflect.Ptr {
				if target.IsNil() {
					target.Set(reflect.New(target.Type().Elem()))
				}
				target = target.Elem()
			}
			err = mapJsonMapToStructValue(jsonValueMap, target)
			if err != nil {
				return fmt.Errorf("error mapping nested struct %v: %v", field.Field.Name, err)
			}
			continue
		} else if jsonArray, isArray := jsonValue.([]any); isArray && isSliceOfMaps(jsonArray, target.Type()) {
			err = mapJsonArrayToSlice(jsonArray, target)
		} else {
			err = SetStructValueByJson(target, jsonValue)
		}
		if err != nil {
			return fmt.Errorf("could not set field %v (json key: %v) of %v: %v", field.Field.Name, field.Key, structValue.Type(), err)
		}
	}
	return nil
}

// initEmpty sets a nil map or slice field to an empty one to prevent panics.
// Fields behind a nil embedded pointer are not set, so the pointer stays nil.
func initEmpty(structValue reflect.Value, field StructField) {
	kind := field.Field.Type.Kind()
	if kind != reflect.Map && kind != reflect.Slice {
		return
	}
	value, ok := fieldForRead(structValue, field.Index)
	if !ok || !value.CanSet() || !value.IsNil() {
		return
	}
	if kind == reflect.Map {
		value.Set(reflect.MakeMap(value.Type()))
	} else {
		value.Set(reflect.MakeSlice(value.Type(), 0, 0))
	}
}

// isSliceOfMaps reports if all items of the array are JsonMaps and the target is a slice of structs.
func isSliceOfMaps(jsonArray []any, target reflect.Type) bool {
	if target.Kind() != reflect.Slice || target.Elem().Kind() != reflect.Struct {
		return false
	}
	for _, item := range jsonArray {
		if _, ok := item.(map[string]any); !ok {
			return false
		}
	}
	return true
}

// mapJsonArrayToSlice writes the JsonMaps of the array into a new slice of structs, like AnyToType does,
// but without converting every item on its own.
func mapJsonArrayToSlice(jsonArray []any, target reflect.Value) error {
	slice := reflect.MakeSlice(target.Type(), len(jsonArray), len(jsonArray))
	for i, item := range jsonArray {
		err := mapJsonMapToStructValue(item.(map[string]any), slice.Index(i))
		if err != nil {
			return fmt.Errorf("error converting item at index %d: %v", i, err)
		}
	}
	target.Set(slice)
	return nil
}

func isStructOrPointer(t reflect.Type) bool {
	return t.Kind() == reflect.Struct || (t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct)
}

// decodeJSON writes the value into the field with encoding/json, by decoding it into a struct
// of the given decodeType that has the field (with its json options) as only field.
// The field keeps its current value as base, like encoding/json does for existing structs and maps.
func decodeJSON(target reflect.Value, decodeType reflect.Type, jsonValue any) error {
	encoded, err := json.Marshal(map[string]any{"V": jsonValue})
	if err != nil {
		return fmt.Errorf("error marshaling value: %v", err)
	}
	wrapper := reflect.New(decodeType)
	wrapper.Elem().Field(0).Set(target)
	err = json.Unmarshal(encoded, wrapper.Interface())
	if err != nil {
		return err
	}
	target.Set(wrapper.Elem().Field(0))
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
