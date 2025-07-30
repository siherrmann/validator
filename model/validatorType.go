package model

import (
	"reflect"
)

// ValidatorType is the type for all available validation types.
type ValidatorType string

const (
	String ValidatorType = "string"
	Int    ValidatorType = "int"
	Float  ValidatorType = "float"
	Bool   ValidatorType = "bool"
	Array  ValidatorType = "array"
	Map    ValidatorType = "map"
	Struct ValidatorType = "struct"
)

func (v ValidatorType) ToReflectType() reflect.Type {
	switch v {
	case String:
		return reflect.TypeOf("")
	case Int:
		return reflect.TypeOf(int(0))
	case Float:
		return reflect.TypeOf(float64(0))
	case Bool:
		return reflect.TypeOf(false)
	case Array:
		return reflect.TypeOf([]string{})
	case Map:
		return reflect.TypeOf(map[string]string{})
	case Struct:
		return reflect.TypeOf(struct{}{})
	default:
		return reflect.TypeOf(struct{}{})
	}
}

// TypeFromInterface determines the ValidatorType based on the type of the input interface.
// It checks the type of the input and returns the corresponding ValidatorType.
// If the type is not recognized, it defaults to Struct.
// It handles basic types like string, int, float, bool, and complex types like JsonMap and arrays.
// It also checks for time.Time type and returns the appropriate ValidatorType.
func ReflectKindToValidatorType(reflectType reflect.Kind) ValidatorType {
	switch reflectType {
	case reflect.String:
		return String
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		return Int
	case reflect.Float64, reflect.Float32:
		return Float
	case reflect.Bool:
		return Bool
	case reflect.Map:
		return Map
	case reflect.Slice, reflect.Array:
		return Array
	case reflect.Struct:
		return Struct
	default:
		return Struct
	}
}
