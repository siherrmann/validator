package helper

import (
	"fmt"
	"reflect"
)

func MapKeysToArrayOfAny(v any) ([]any, error) {
	aany := []any{}
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return nil, fmt.Errorf("invalid nil value, has to be map")
	}
	if rv.Type().Kind() != reflect.Map {
		return nil, fmt.Errorf("invalid type %v, has to be map", rv.Type().Kind())
	}

	for _, mk := range rv.MapKeys() {
		aany = append(aany, mk.Interface())
	}

	return aany, nil
}

// MapValue returns the value of the given key in a map with string keys (eg. a JsonMap or map[string]string).
// ok is false if the key does not exist or v is not a map with string keys.
func MapValue(v any, key string) (any, bool) {
	if m, ok := v.(map[string]any); ok {
		value, ok := m[key]
		return value, ok
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Map || rv.Type().Key().Kind() != reflect.String {
		return nil, false
	}
	value := rv.MapIndex(reflect.ValueOf(key).Convert(rv.Type().Key()))
	if !value.IsValid() {
		return nil, false
	}
	return value.Interface(), true
}
