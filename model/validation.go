package model

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/siherrmann/validator/helper"
)

// Default tag type.
const VLD string = "vld"

// Validation represents a validation rule for a struct field.
type Validation struct {
	Key         string
	Type        ValidatorType
	Requirement string
	Groups      []*Group
	Default     string
	// Inner Struct validation
	InnerValidation []Validation
}

// GetValidationsFromStruct extracts validation rules from a struct based on the provided tag type.
// It iterates over the struct fields, checks for the specified tag type, and constructs Validation.
func GetValidationsFromStruct(in any, tagType string) ([]Validation, error) {
	err := helper.CheckValidPointerToStruct(in)
	if err != nil {
		return nil, err
	}

	validations := []Validation{}

	structFull := reflect.ValueOf(in).Elem()
	for i := 0; i < structFull.Type().NumField(); i++ {
		field := structFull.Field(i)
		fieldType := structFull.Type().Field(i)

		validation, err := GetValidationFromStructField(tagType, field, fieldType)
		if err != nil {
			return nil, err
		}

		if len(validation.Requirement) > 0 {
			validations = append(validations, validation)
		}
	}
	return validations, nil
}

// GetValidationFromStructField extracts validation rules from a struct field based on the provided tag type.
// It checks the field's tag for the specified tag type and constructs a Validation object.
// If no json tag is found, it uses the field name as the key.
func GetValidationFromStructField(tagType string, fieldValue reflect.Value, fieldType reflect.StructField) (Validation, error) {
	validation := Validation{}
	validation.Key = fieldType.Name
	if len(fieldType.Tag.Get("json")) > 0 {
		validation.Key = fieldType.Tag.Get("json")
	}
	validation.Type = TypeFromInterface(fieldValue.Interface())
	validation.Requirement = "-"

	tagIndex := 0
	tagSplit := strings.Split(fieldType.Tag.Get(string(tagType)), ", ")

	if len(tagSplit) > tagIndex {
		validation.Requirement = tagSplit[tagIndex]
		tagIndex++
	}

	if len(tagSplit) > tagIndex {
		var err error
		validation.Groups, err = GetGroups(tagSplit[tagIndex])
		if err != nil {
			return Validation{}, fmt.Errorf("error extracting group: %v", err)
		}
	}

	if helper.IsArrayOfStruct(fieldValue.Interface()) {
		for i := 0; i < fieldValue.Len(); i++ {
			innerStruct := fieldValue.Index(i).Addr().Interface()
			innerValidation, err := GetValidationsFromStruct(innerStruct, string(tagType))
			if err != nil {
				return Validation{}, fmt.Errorf("error getting inner validation: %v", err)
			}
			validation.InnerValidation = append(validation.InnerValidation, innerValidation...)
		}
	}

	return validation, nil
}

// ValidatorMap is a map of validation keys to Validation objects.
// It is used to store and manage multiple validation rules for different struct fields.
type ValidatorMap map[string]Validation

// ValidatorType is the type for all available validation types.
type ValidatorType string

const (
	String      ValidatorType = "string"
	Int         ValidatorType = "int"
	Float       ValidatorType = "float"
	Bool        ValidatorType = "bool"
	Array       ValidatorType = "array"
	Map         ValidatorType = "map"
	Struct      ValidatorType = "struct"
	Time        ValidatorType = "time"
	TimeISO8601 ValidatorType = "timeIso8601"
	TimeUnix    ValidatorType = "timeUnix"
)

// TypeFromInterface determines the ValidatorType based on the type of the input interface.
// It checks the type of the input and returns the corresponding ValidatorType.
// If the type is not recognized, it defaults to Struct.
// It handles basic types like string, int, float, bool, and complex types like JsonMap and arrays.
// It also checks for time.Time type and returns the appropriate ValidatorType.
func TypeFromInterface(in any) ValidatorType {
	switch in.(type) {
	case string:
		return String
	case int, int64, int32, int16, int8:
		return Int
	case float64, float32:
		return Float
	case bool:
		return Bool
	case JsonMap, map[string]string, map[string]int, map[string]int64, map[string]int32, map[string]int16, map[string]int8, map[string]float64, map[string]float32, map[string]bool:
		return Map
	case []string, []int, []int64, []int32, []int16, []int8, []float64, []float32, []bool:
		return Array
	case time.Time:
		return Time
	default:
		// custom types
		if reflect.TypeOf(in).Kind() == reflect.String {
			return String
		} else if reflect.TypeOf(in).Kind() == reflect.Int {
			return Int
		} else if reflect.TypeOf(in).Kind() == reflect.Float64 || reflect.TypeOf(in).Kind() == reflect.Float32 {
			return Float
		} else if reflect.TypeOf(in).Kind() == reflect.Bool {
			return Bool
			// other types
		} else if helper.IsArray(in) {
			return Array
		} else {
			return Struct
		}
	}
}
