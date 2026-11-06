package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/siherrmann/validator/helper"
	"github.com/siherrmann/validator/model"
	"github.com/siherrmann/validator/parser"
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

		_, _, isInline := helper.FieldKey(fieldType)
		if isInline {
			validations = append(validations, validation.InnerValidation...)
		} else {
			if len(validation.Requirement) > 0 || len(validation.InnerValidation) > 0 {
				validations = append(validations, *validation)
			}
		}
	}
	return validations, nil
}

// GetValidationFromStructField extracts validation rules from a struct field based on the provided tag type.
// It checks the field's tag for the specified tag type and constructs a Validation object.
// The key is taken from the json tag, falls back to the yaml tag and then to the field name.
func GetValidationFromStructField(tagType string, fieldValue reflect.Value, fieldType reflect.StructField) (*model.Validation, error) {
	validation := &model.Validation{}
	validation.Key, validation.OmitEmpty, _ = helper.FieldKey(fieldType)
	validation.Type = model.ReflectTypeToValidatorType(fieldValue.Type())
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

		// Extract default value & cache AST
		p := parser.NewParser()
		root, err := p.ParseValidation(validation.Requirement)
		if err != nil {
			return nil, fmt.Errorf("error parsing validation: %v", err)
		}
		validation.Default = root.ExtractDefault()
		validation.RequirementAST = root.RootValue
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
	} else if helper.IsStruct(fieldValue.Interface()) || (fieldValue.Type().Kind() == reflect.Ptr && fieldValue.Type().Elem().Kind() == reflect.Struct) {
		fieldTypeElem := fieldValue.Type()
		if fieldTypeElem.Kind() == reflect.Ptr {
			fieldTypeElem = fieldTypeElem.Elem()
		}
		innerStruct := reflect.New(fieldTypeElem).Interface()
		innerValidation, err := GetValidationsFromStruct(innerStruct, string(tagType))
		if err != nil {
			return nil, fmt.Errorf("error getting inner validation from struct: %v", err)
		}
		validation.InnerValidation = append(validation.InnerValidation, innerValidation...)
	}

	return validation, nil
}
