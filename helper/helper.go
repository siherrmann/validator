package helper

import (
	"fmt"
	"log"
	"reflect"
)

// Checks if the given value is a string.
func IsString(in interface{}) bool {
	return reflect.TypeOf(in).Kind() == reflect.String
}

// Checks if the given value is an array or slice.
func IsArray(in interface{}) bool {
	return reflect.TypeOf(in).Kind() == reflect.Array || reflect.TypeOf(in).Kind() == reflect.Slice
}

// Checks if the given value is an array or slice of structs.
func IsArrayOfStruct(in interface{}) bool {
	return IsArray(in) && reflect.TypeOf(in).Elem().Kind() == reflect.Struct
}

// Checks if the given value is an array of maps.
func IsArrayOfMap(in interface{}) bool {
	return IsArray(in) && reflect.TypeOf(in).Elem().Kind() == reflect.Map
}

// Checks if the given value is a pointer to a struct.
func CheckValidPointerToStruct(in interface{}) error {
	value := reflect.ValueOf(in)
	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("value has to be of kind pointer, was %T", value)
	}
	if value.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("value has to be of kind struct, was %T", value)
	}
	return nil
}

// Converts an array of interface{} to an array of a specific type T.
// The function checks if the elements in the input array can be converted to the specified type T.
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
				log.Printf("valueOfType: %v, index: %v, arrayOfType: %T", inReflect.Len(), i, arrayOfType)
				return []T{}, fmt.Errorf("invalid input array element type: %v, expected: %v", reflect.TypeOf(inReflect.Index(i).Interface()).Kind(), reflect.TypeOf(arrayOfType).Elem().Kind())
			}
			arrayOfType = append(arrayOfType, valueOfType)
		}
	}
	return arrayOfType, nil
}
