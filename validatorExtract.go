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
// It iterates over the fields of the StructModel of the struct (see helper.StructModelOf),
// so fields of embedded and inline structs are extracted on the level of the struct.
//
// A field whose type refers back to a struct that contains it (eg. `Children []Node` in Node)
// gets an InnerValidationRef to the validations of that struct instead of InnerValidation.
func GetValidationsFromStruct(in any, tagType string) ([]model.Validation, error) {
	err := helper.CheckValidPointerToStruct(in)
	if err != nil {
		return nil, err
	}

	e := &extractor{tagType: tagType, building: map[reflect.Type]*[]model.Validation{}}
	validations, err := e.structValidations(reflect.TypeOf(in).Elem())
	if err != nil {
		return nil, err
	}
	if e.hasRef {
		markFieldUsage(validations)
	}
	return validations, nil
}

// GetValidationFromStructField extracts validation rules from a struct field based on the provided tag type.
// It checks the field's tag for the specified tag type and constructs a Validation object.
// The key is taken from the json tag, falls back to the yaml tag and then to the field name.
func GetValidationFromStructField(tagType string, fieldValue reflect.Value, fieldType reflect.StructField) (*model.Validation, error) {
	key, omitEmpty, _ := helper.FieldKey(fieldType)
	e := &extractor{tagType: tagType, building: map[reflect.Type]*[]model.Validation{}}
	validation, err := e.fieldValidation(fieldType, fieldValue.Type(), key, omitEmpty)
	if err != nil || validation == nil {
		return validation, err
	}
	if e.hasRef {
		validations := []model.Validation{*validation}
		markFieldUsage(validations)
		validation = &validations[0]
	}
	return validation, nil
}

// extractor holds the state of a single extraction.
type extractor struct {
	tagType string
	// building holds the struct types that are extracted at the moment, with a pointer
	// to their validations, which is filled when the extraction of the struct is done.
	building map[reflect.Type]*[]model.Validation
	// hasRef is set if an InnerValidationRef was set, NoFieldUsage then has to be checked again.
	hasRef bool
}

// structValidations extracts the validations of all fields of the struct type.
func (e *extractor) structValidations(structType reflect.Type) ([]model.Validation, error) {
	ref := new([]model.Validation)
	e.building[structType] = ref
	defer delete(e.building, structType)

	validations := []model.Validation{}
	for _, field := range helper.StructModelOf(structType).Fields {
		validation, err := e.fieldValidation(field.Field, field.Field.Type, field.Key, field.OmitEmpty)
		if err != nil {
			return nil, err
		} else if validation == nil {
			continue
		}
		validations = append(validations, *validation)
	}

	*ref = validations
	return validations, nil
}

// fieldValidation extracts the validation of a field of the given type, it is nil if the field has no tag.
func (e *extractor) fieldValidation(fieldType reflect.StructField, valueType reflect.Type, key string, omitEmpty bool) (*model.Validation, error) {
	validation := &model.Validation{
		Key:         key,
		OmitEmpty:   omitEmpty,
		Type:        model.ReflectTypeToValidatorType(valueType),
		Requirement: "-",
	}

	tagIndex := 0
	tagSplit := strings.Split(fieldType.Tag.Get(e.tagType), ", ")

	if len(tagSplit) > tagIndex {
		// Ignore if tag is empty, we do not want to validate this field at all
		if len(tagSplit[tagIndex]) == 0 {
			return nil, nil
		}
		validation.Requirement = tagSplit[tagIndex]
		tagIndex++

		// Extract default value & cache AST
		p := parser.NewParser()
		root, err := p.ParseValidation(validation.Requirement)
		if err != nil {
			return nil, fmt.Errorf("error parsing validation of field %v: %w", key, err)
		}
		err = checkNumberConditions(root.RootValue, validation.Type)
		if err != nil {
			return nil, fmt.Errorf("error parsing validation of field %v: %w %q: %w", key, model.ErrInvalidRequirement, validation.Requirement, err)
		}
		validation.Default = root.ExtractDefault()
		validation.RequirementAST = root.RootValue
	}

	if len(tagSplit) > tagIndex {
		var err error
		validation.Groups, err = model.GetGroups(tagSplit[tagIndex])
		if err != nil {
			return nil, fmt.Errorf("error extracting group of field %v: %w %q: %w", key, model.ErrInvalidRequirement, tagSplit[tagIndex], err)
		}
	}

	if innerType, ok := innerStructType(valueType); ok {
		if ref, building := e.building[innerType]; building {
			validation.InnerValidationRef = ref
			e.hasRef = true
		} else if hasTag(innerType, e.tagType, map[reflect.Type]bool{}) {
			innerValidation, err := e.structValidations(innerType)
			if err != nil {
				return nil, fmt.Errorf("error getting inner validation from %v: %w", valueType, err)
			}
			validation.InnerValidation = innerValidation
		}
	}

	if validation.RequirementAST.HasCondition(model.ConditionType.IsCorrect) && (validation.Type != model.Struct || (validation.InnerValidation == nil && validation.InnerValidationRef == nil)) {
		// InnerValidationRef is not filled yet for a recursive type, so Inner() can not be used here.
		return nil, fmt.Errorf("error parsing validation of field %v: %w %q: cmn and cmx are only usable on struct fields with inner validations", key, model.ErrInvalidRequirement, validation.Requirement)
	}

	validation.NoFieldUsage = !requirementOrInnerUsesFields(validation)
	return validation, nil
}

// hasTag reports if a field of the struct type or of one of its inner structs has the tag.
// A struct without any tag has no inner validations, its value is taken as it is.
func hasTag(structType reflect.Type, tagType string, seen map[reflect.Type]bool) bool {
	seen[structType] = true
	for _, field := range helper.StructModelOf(structType).Fields {
		if field.Field.Tag.Get(tagType) != "" {
			return true
		}
		if innerType, ok := innerStructType(field.Field.Type); ok && !seen[innerType] && hasTag(innerType, tagType, seen) {
			return true
		}
	}
	return false
}

// checkNumberConditions returns an error if an `equ` or `neq` condition of an int or float field
// (or of a list or map, where the length is compared) has a value that can not be converted to a number.
// Conditions on other fields (`rel`) and references (`$`) are only known while validating.
func checkNumberConditions(astValue *model.AstValue, validatorType model.ValidatorType) error {
	var numberType reflect.Type
	switch validatorType {
	case model.Int, model.Array, model.Map:
		numberType = reflect.TypeOf(0)
	case model.Float:
		numberType = reflect.TypeOf(0.0)
	default:
		return nil
	}

	for _, v := range astValue.ConditionGroup {
		switch {
		case v.Type == model.GROUP:
			if err := checkNumberConditions(v, validatorType); err != nil {
				return err
			}
		case v.Ref != nil || (v.ConditionType != model.EQUAL && v.ConditionType != model.NOT_EQUAL):
			continue
		default:
			if _, err := helper.AnyToType(v.ConditionValue, numberType); err != nil {
				return fmt.Errorf("condition %v%v on a field of type %v, value %v is not a number", v.ConditionType, v.ConditionValue, validatorType, v.ConditionValue)
			}
		}
	}
	return nil
}

// innerStructType returns the struct type of a struct, pointer to struct or slice/array of structs.
func innerStructType(valueType reflect.Type) (reflect.Type, bool) {
	switch valueType.Kind() {
	case reflect.Slice, reflect.Array:
		valueType = valueType.Elem()
		return valueType, valueType.Kind() == reflect.Struct
	case reflect.Ptr:
		valueType = valueType.Elem()
	}
	return valueType, valueType.Kind() == reflect.Struct
}

// requirementOrInnerUsesFields reports if the requirement of the validation or one of its inner validations
// uses `rel` or a `$` reference. An InnerValidationRef counts as not using fields here, see markFieldUsage.
func requirementOrInnerUsesFields(validation *model.Validation) bool {
	if validation.RequirementAST != nil && validation.RequirementAST.UsesFields {
		return true
	}
	for i := range validation.InnerValidation {
		if validation.InnerValidation[i].UsesFields() {
			return true
		}
	}
	return false
}

// markFieldUsage sets NoFieldUsage to false for all validations whose inner validations use fields
// through an InnerValidationRef. It repeats until nothing changes, because a reference can be part of a cycle.
func markFieldUsage(validations []model.Validation) {
	for markFieldUsageOnce(validations, map[*model.Validation]bool{}) {
	}
}

func markFieldUsageOnce(validations []model.Validation, seen map[*model.Validation]bool) (changed bool) {
	if len(validations) == 0 || seen[&validations[0]] {
		return false
	}
	seen[&validations[0]] = true

	for i := range validations {
		validation := &validations[i]
		inner := validation.Inner()
		if markFieldUsageOnce(inner, seen) {
			changed = true
		}
		if !validation.NoFieldUsage {
			continue
		}
		for j := range inner {
			if inner[j].UsesFields() {
				validation.NoFieldUsage = false
				changed = true
				break
			}
		}
	}
	return changed
}
