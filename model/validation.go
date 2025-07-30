package model

import (
	"fmt"
	"reflect"
	"strings"

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
	validation.Type = ReflectKindToValidatorType(fieldValue.Type().Kind())
	validation.Requirement = "-"

	tagIndex := 0
	tagSplit := strings.Split(fieldType.Tag.Get(string(tagType)), ", ")

	if len(tagSplit) > tagIndex {
		if len(tagSplit[tagIndex]) <= 3 && tagSplit[tagIndex] != "-" {
			return Validation{}, fmt.Errorf("invalid requirement %v for field %v", tagSplit[tagIndex], fieldType.Name)
		}
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
		innerStruct := reflect.New(fieldValue.Type().Elem()).Interface()
		innerValidation, err := GetValidationsFromStruct(innerStruct, string(tagType))
		if err != nil {
			return Validation{}, fmt.Errorf("error getting inner validation from array: %v", err)
		}
		validation.InnerValidation = append(validation.InnerValidation, innerValidation...)
	} else if helper.IsStruct(fieldValue.Interface()) {
		innerStruct := reflect.New(fieldValue.Type()).Interface()
		innerValidation, err := GetValidationsFromStruct(innerStruct, string(tagType))
		if err != nil {
			return Validation{}, fmt.Errorf("error getting inner validation from struct: %v", err)
		}
		validation.InnerValidation = append(validation.InnerValidation, innerValidation...)
	}

	return validation, nil
}

// ValidatorMap is a map of validation keys to Validation objects.
// It is used to store and manage multiple validation rules for different struct fields.
type ValidatorMap map[string]Validation
