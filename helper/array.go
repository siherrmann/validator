package helper

import (
	"fmt"
	"reflect"
)

func ArrayToArrayOfAny(v any) ([]any, error) {
	aany := []any{}
	rv := reflect.ValueOf(v)
	if rv.Type().Kind() != reflect.Array && rv.Type().Kind() != reflect.Slice {
		return nil, fmt.Errorf("invalid type %v, has to be array or slice", rv.Type().Kind())
	}

	for i := 0; i < rv.Len(); i++ {
		aany = append(aany, rv.Index(i).Interface())
	}

	return aany, nil
}
