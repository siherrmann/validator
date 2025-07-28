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

func ArrayToArrayOfType(a []any, expectedValue reflect.Type) (reflect.Value, error) {
	targetArray := reflect.MakeSlice(reflect.SliceOf(expectedValue), len(a), len(a))
	for i, item := range a {
		itemConverted, err := AnyToType(item, expectedValue)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("error converting item at index %d: %v", i, err)
		}
		targetArray.Index(i).Set(reflect.ValueOf(itemConverted))
	}
	return targetArray, nil
}
