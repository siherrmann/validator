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
