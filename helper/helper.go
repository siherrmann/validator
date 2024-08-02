package helper

import (
	"fmt"
	"reflect"
)

func ArrayOfInterfaceToArrayOf[T comparable](in []interface{}) ([]T, error) {
	inReflect := reflect.ValueOf(in)
	arrayOfType := []T{}
	for i := 0; i < inReflect.Len(); i++ {
		var newValueInterface interface{}
		switch reflect.TypeOf(arrayOfType).Elem().Kind() {
		case reflect.Int:
			valueOfType, ok := inReflect.Index(i).Interface().(float64)
			if ok {
				newValueInterface = int(valueOfType)
				arrayOfType = append(arrayOfType, newValueInterface.(T))
			} else {
				valueOfType, ok := inReflect.Index(i).Interface().(T)
				if !ok {
					return []T{}, fmt.Errorf("invalid input array element type: %v, expected: %v", reflect.TypeOf(inReflect.Index(i).Interface()).Kind(), reflect.TypeOf(arrayOfType).Elem().Kind())
				}
				arrayOfType = append(arrayOfType, valueOfType)
			}
		case reflect.Int8:
			valueOfType, ok := inReflect.Index(i).Interface().(float64)
			if ok {
				newValueInterface = int8(valueOfType)
				arrayOfType = append(arrayOfType, newValueInterface.(T))
			} else {
				valueOfType, ok := inReflect.Index(i).Interface().(T)
				if !ok {
					return []T{}, fmt.Errorf("invalid input array element type: %v, expected: %v", reflect.TypeOf(inReflect.Index(i).Interface()).Kind(), reflect.TypeOf(arrayOfType).Elem().Kind())
				}
				arrayOfType = append(arrayOfType, valueOfType)
			}
		case reflect.Int16:
			valueOfType, ok := inReflect.Index(i).Interface().(float64)
			if ok {
				newValueInterface = int16(valueOfType)
				arrayOfType = append(arrayOfType, newValueInterface.(T))
			} else {
				valueOfType, ok := inReflect.Index(i).Interface().(T)
				if !ok {
					return []T{}, fmt.Errorf("invalid input array element type: %v, expected: %v", reflect.TypeOf(inReflect.Index(i).Interface()).Kind(), reflect.TypeOf(arrayOfType).Elem().Kind())
				}
				arrayOfType = append(arrayOfType, valueOfType)
			}
		case reflect.Int32:
			valueOfType, ok := inReflect.Index(i).Interface().(float64)
			if ok {
				newValueInterface = int32(valueOfType)
				arrayOfType = append(arrayOfType, newValueInterface.(T))
			} else {
				valueOfType, ok := inReflect.Index(i).Interface().(T)
				if !ok {
					return []T{}, fmt.Errorf("invalid input array element type: %v, expected: %v", reflect.TypeOf(inReflect.Index(i).Interface()).Kind(), reflect.TypeOf(arrayOfType).Elem().Kind())
				}
				arrayOfType = append(arrayOfType, valueOfType)
			}
		case reflect.Int64:
			valueOfType, ok := inReflect.Index(i).Interface().(float64)
			if ok {
				newValueInterface = int64(valueOfType)
				arrayOfType = append(arrayOfType, newValueInterface.(T))
			} else {
				valueOfType, ok := inReflect.Index(i).Interface().(T)
				if !ok {
					return []T{}, fmt.Errorf("invalid input array element type: %v, expected: %v", reflect.TypeOf(inReflect.Index(i).Interface()).Kind(), reflect.TypeOf(arrayOfType).Elem().Kind())
				}
				arrayOfType = append(arrayOfType, valueOfType)
			}
		default:
			valueOfType, ok := inReflect.Index(i).Interface().(T)
			if !ok {
				return []T{}, fmt.Errorf("invalid input array element type: %v, expected: %v", reflect.TypeOf(inReflect.Index(i).Interface()).Kind(), reflect.TypeOf(arrayOfType).Elem().Kind())
			}
			arrayOfType = append(arrayOfType, valueOfType)
		}
	}
	return arrayOfType, nil
}
