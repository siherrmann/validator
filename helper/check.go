package helper

import (
	"fmt"
	"reflect"
)

// Checks if the given value is a string.
func IsString(in any) bool {
	return reflect.TypeOf(in).Kind() == reflect.String
}

// Checks if the given value is a struct.
func IsStruct(in any) bool {
	return reflect.TypeOf(in).Kind() == reflect.Struct
}

// Checks if the given value is an array or slice.
func IsArray(in any) bool {
	return reflect.TypeOf(in).Kind() == reflect.Array || reflect.TypeOf(in).Kind() == reflect.Slice
}

// Checks if the given value is an array or slice of structs.
func IsArrayOfStruct(in any) bool {
	return IsArray(in) && reflect.TypeOf(in).Elem().Kind() == reflect.Struct
}

// Checks if the given value is an array of maps.
func IsArrayOfMap(in any) bool {
	return IsArray(in) && reflect.TypeOf(in).Elem().Kind() == reflect.Map
}

// Checks if the given value is a pointer to a struct.
func CheckValidPointerToStruct(in any) error {
	value := reflect.ValueOf(in)
	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("value has to be of kind pointer, was %T", value)
	}
	if value.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("value has to be of kind struct, was %T", value)
	}
	return nil
}
