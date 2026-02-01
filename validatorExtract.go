package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/siherrmann/validator/helper"
	"github.com/siherrmann/validator/model"
)

// GetValidationsFromStruct extracts validation rules from a struct based on the provided tag type.
// It iterates over the struct fields, checks for the specified tag type, and constructs Validation.
func GetValidationsFromStruct(in any, tagType string) ([]model.Validation, error) {
	err := helper.CheckValidPointerToStruct(in)
	if err != nil {
		return nil, err
	}

	validations := []model.Validation{}

	structFull := reflect.ValueOf(in).Elem()
	for i := 0; i < structFull.Type().NumField(); i++ {
		field := structFull.Field(i)
		fieldType := structFull.Type().Field(i)

		validation, err := GetValidationFromStructField(tagType, field, fieldType)
		if err != nil {
			return nil, err
		} else if validation == nil {
			continue
		}

		if len(validation.Requirement) > 0 {
			validations = append(validations, *validation)
		}
	}
	return validations, nil
}

// GetValidationFromStructField extracts validation rules from a struct field based on the provided tag type.
// It checks the field's tag for the specified tag type and constructs a Validation object.
// If no json tag is found, it uses the field name as the key.
func GetValidationFromStructField(tagType string, fieldValue reflect.Value, fieldType reflect.StructField) (*model.Validation, error) {
	validation := &model.Validation{}
	validation.Key = fieldType.Name
	if len(fieldType.Tag.Get("json")) > 0 {
		validation.Key = fieldType.Tag.Get("json")
	}
	validation.Type = model.ReflectKindToValidatorType(fieldValue.Type().Kind())
	validation.Requirement = "-"

	tagIndex := 0
	tagSplit := strings.Split(fieldType.Tag.Get(string(tagType)), ", ")

	if len(tagSplit) > tagIndex {
		// Ignore if tag is empty, we do not want to validate this field at all
		if len(tagSplit[tagIndex]) <= 3 && tagSplit[tagIndex] != "-" {
			return nil, nil
		}
		validation.Requirement = tagSplit[tagIndex]
		tagIndex++
	}

	if len(tagSplit) > tagIndex {
		var err error
		validation.Groups, err = model.GetGroups(tagSplit[tagIndex])
		if err != nil {
			return nil, fmt.Errorf("error extracting group: %v", err)
		}
	}

	if helper.IsArrayOfStruct(fieldValue.Interface()) {
		innerStruct := reflect.New(fieldValue.Type().Elem()).Interface()
		innerValidation, err := GetValidationsFromStruct(innerStruct, string(tagType))
		if err != nil {
			return nil, fmt.Errorf("error getting inner validation from array: %v", err)
		}
		validation.InnerValidation = append(validation.InnerValidation, innerValidation...)
	} else if helper.IsStruct(fieldValue.Interface()) {
		innerStruct := reflect.New(fieldValue.Type()).Interface()
		innerValidation, err := GetValidationsFromStruct(innerStruct, string(tagType))
		if err != nil {
			return nil, fmt.Errorf("error getting inner validation from struct: %v", err)
		}
		validation.InnerValidation = append(validation.InnerValidation, innerValidation...)
	}

	return validation, nil
}
