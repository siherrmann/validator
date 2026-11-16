package validator

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/siherrmann/validator/helper"
	"github.com/siherrmann/validator/model"
	"github.com/siherrmann/validator/parser"
	"github.com/siherrmann/validator/validators"
)

type ValidationFunc func(input any, astValue *model.AstValue) error

// Validator is the main struct for validation.
type Validator struct {
	ValidationFuncs map[string]ValidationFunc
}

// NewValidator creates a new Validator instance with an empty validation functions map.
func NewValidator() *Validator {
	return &Validator{
		ValidationFuncs: make(map[string]ValidationFunc),
	}
}

// AddValidationFunc adds a custom validation function to the Validator.
// The function can be used in validation requirements with the name provided (`fun<name>`).
func (r *Validator) AddValidationFunc(fn ValidationFunc, name string) {
	r.ValidationFuncs[name] = fn
}

// Validate validates a given struct by the given tagType.
// It checks if the keys are in the struct and validates the values.
// It returns an error if the validation fails.
//
// Extraction is cached per (type, tagType) after the first call for that combination
// — see getCachedValidationsFromStruct. Use GetValidationsFromStruct directly if a
// fresh, uncached extraction is required.
func (r *Validator) Validate(v any, tagType ...string) error {
	tagTypeSet := model.VLD
	if len(tagType) > 0 {
		tagTypeSet = tagType[0]
	}

	jsonMap := map[string]any{}
	err := helper.UnmapStructToJsonMap(v, &jsonMap)
	if err != nil {
		return fmt.Errorf("error unmapping struct to json map: %v", err)
	}

	validations, err := getCachedValidationsFromStruct(v, tagTypeSet)
	if err != nil {
		return fmt.Errorf("error getting validations from struct: %v", err)
	}

	_, err = r.ValidateWithValidation(jsonMap, validations)
	if err != nil {
		return fmt.Errorf("error validating struct: %v", err)
	}

	return nil
}

// ValidateAndUpdate validates a given JsonMap by the given validations and updates the struct.
// It checks if the keys are in the map, validates the values and updates the struct if the validation passes.
// It returns an error if the validation fails or if the struct cannot be updated.
//
// Extraction is cached per (type, tagType), same as Validate.
func (r *Validator) ValidateAndUpdate(jsonInput map[string]any, structToUpdate any, tagType ...string) error {
	tagTypeSet := model.VLD
	if len(tagType) > 0 {
		tagTypeSet = tagType[0]
	}

	validations, err := getCachedValidationsFromStruct(structToUpdate, tagTypeSet)
	if err != nil {
		return fmt.Errorf("error getting validations from struct: %v", err)
	}

	validatedMap, err := r.ValidateWithValidation(jsonInput, validations)
	if err != nil {
		return fmt.Errorf("error validating struct: %v", err)
	}

	err = helper.MapJsonMapToStruct(validatedMap, structToUpdate)
	if err != nil {
		return fmt.Errorf("error mapping json map to struct: %v", err)
	}

	return nil
}

// ValidateAndUpdateWithValidation validates a given JsonMap by the given validations and updates the map.
// It checks if the keys are in the map, validates the values and updates the map if the validation passes.
// It returns an error if the validation fails or if the map cannot be updated.
func (r *Validator) ValidateAndUpdateWithValidation(jsonInput map[string]any, mapToUpdate *map[string]any, validations []model.Validation) error {
	validatedValues, err := r.ValidateWithValidation(jsonInput, validations)
	if err != nil {
		return fmt.Errorf("error validating json map: %v", err)
	}

	for k, v := range validatedValues {
		(*mapToUpdate)[k] = v
	}

	return nil
}

// ValidateWithValidation validates a given JsonMap by the given validations.
// This is the main validation function that is used for all validation.
// It checks if the keys are in the map, validates the values and returns a new JsonMap.
//
// If a validation has groups, it checks if the values are valid for the groups.
// If a validation has a key that is already in the map, it returns an error.
//
// All fields are validated, also if one is invalid. The errors of invalid fields (that are not in a group)
// are returned together as [model.FieldErrors], with the JsonMap of the valid fields.
// Any other error (eg. a duplicate key or a failed group) is returned with an empty JsonMap.
//
// It returns a new JsonMap with the validated values or an error if the validation fails.
func (r *Validator) ValidateWithValidation(jsonInput map[string]any, validations []model.Validation) (map[string]any, error) {
	keys := []string{}
	groups := map[string]*model.Group{}
	groupSize := map[string]int{}
	groupErrors := map[string][]error{}
	fieldErrors := model.FieldErrors{}

	validateValues := map[string]any{}

	for validationIndex := range validations {
		validation := validations[validationIndex]
		if len(validation.Key) > 0 && slices.Contains(keys, validation.Key) {
			return map[string]any{}, fmt.Errorf("duplicate validation key: %v", validation.Key)
		} else {
			keys = append(keys, validation.Key)
		}

		for _, g := range validation.Groups {
			groups[g.Name] = g
			groupSize[g.Name]++
		}

		jsonValue, ok, err := r.validateField(jsonInput, &validation)
		if err != nil {
			if len(validation.Groups) == 0 {
				fieldErrors = append(fieldErrors, err)
			}
			for _, group := range validation.Groups {
				groupErrors[group.Name] = append(groupErrors[group.Name], err)
			}
			continue
		} else if !ok {
			continue
		}

		validateValues[validation.Key] = jsonValue
	}

	err := validators.ValidateGroups(groups, groupSize, groupErrors)
	if err != nil {
		return map[string]any{}, err
	}

	if len(fieldErrors) > 0 {
		return validateValues, fieldErrors
	}

	return validateValues, nil
}

// validateField validates the value of a single validation in the given JsonMap.
// It returns the validated value and if the value should be added to the validated JsonMap.
func (r *Validator) validateField(jsonInput map[string]any, validation *model.Validation) (any, bool, error) {
	jsonValue, ok := jsonInput[validation.Key]
	if !ok {
		if validation.Default != "" {
			jsonValue = validation.Default
		} else if strings.TrimSpace(validation.Requirement) == string(model.NONE) {
			return nil, false, nil
		} else {
			return nil, false, fmt.Errorf("json %v key not in map", validation.Key)
		}
	}
	if validation.OmitEmpty && jsonValue == "" {
		return nil, false, nil
	}

	if jsonValue != nil {
		switch validation.Type {
		case model.String, model.Int, model.Float, model.Bool:
			coerced, coerceErr := helper.AnyToType(jsonValue, validation.Type.ToReflectType())
			if coerceErr == nil {
				jsonValue = coerced
			}
		}
	}

	var err error
	switch validation.Type {
	case model.Struct:
		if jsonValueMap, ok := jsonValue.(map[string]any); ok {
			jsonValue, err = r.validateStruct(jsonValueMap, validation)
			if err != nil {
				return nil, false, fmt.Errorf("field %v invalid: %v", validation.Key, err.Error())
			}
		} else {
			err = r.ValidateValueWithParser(jsonValue, validation)
		}
	case model.Array:
		if helper.IsArray(jsonValue) && len(validation.InnerValidation) > 0 {
			jsonArray, ok := jsonValue.([]any)
			if !ok {
				return nil, false, fmt.Errorf("field %v must be of type array, was %T", validation.Key, jsonValue)
			}

			validatedArray := []any{}
			for _, jsonValueInner := range jsonArray {
				jsonValueInnerMap, err := helper.GetValidMap(jsonValueInner)
				if err != nil {
					return nil, false, fmt.Errorf("field %v invalid: %v", validation.Key, err.Error())
				}

				validatedInnerMap, err := r.ValidateWithValidation(jsonValueInnerMap, validation.InnerValidation)
				if err != nil {
					return nil, false, fmt.Errorf("field %v invalid: %v", validation.Key, err.Error())
				}
				validatedArray = append(validatedArray, validatedInnerMap)
			}
			jsonValue = validatedArray
		} else if helper.IsArray(jsonValue) {
			err = r.ValidateValueWithParser(jsonValue, validation)
		} else if helper.IsString(jsonValue) {
			// Check if the value is a string from a url value.
			jsonValue = []string{jsonValue.(string)}
			err = r.ValidateValueWithParser(jsonValue, validation)
		} else if jsonValue == nil {
			err = r.ValidateValueWithParser(jsonValue, validation)
		} else if len(validation.InnerValidation) == 0 {
			// A lone value (e.g. a single form value) is treated as a one-element list.
			jsonValue = []any{jsonValue}
			err = r.ValidateValueWithParser(jsonValue, validation)
		}
	default:
		err = r.ValidateValueWithParser(jsonValue, validation)
	}

	if err != nil {
		if validation.Default != "" {
			return validation.Default, true, nil
		}
		return nil, false, fmt.Errorf("field %v invalid: %v", validation.Key, err.Error())
	}

	return jsonValue, true, nil
}

// validateStruct validates the JsonMap of a struct field by its inner validations.
//
// Without a requirement on the struct field all inner fields have to be valid.
// With a requirement, the requirement is checked on the JsonMap itself (eg. min and max check the number of keys).
// If the JsonMap is empty and fulfills the requirement, the inner validations are skipped.
// If the requirement contains cmn or cmx, invalid inner fields are left out instead of failing
// and cmn/cmx check the number of correct inner fields (the fields in the validated JsonMap).
func (r *Validator) validateStruct(jsonValueMap map[string]any, validation *model.Validation) (map[string]any, error) {
	astValue, err := getRequirementAST(validation)
	if err != nil {
		return nil, err
	}

	if !astValue.HasCondition(model.ConditionType.IsRequirement) {
		return r.ValidateWithValidation(jsonValueMap, validation.InnerValidation)
	}

	if !astValue.HasCondition(model.ConditionType.IsCorrect) {
		err = r.RunValidatorsOnConditionGroup(jsonValueMap, astValue, nil)
		if err != nil {
			return nil, err
		} else if len(jsonValueMap) == 0 {
			return map[string]any{}, nil
		}
		return r.ValidateWithValidation(jsonValueMap, validation.InnerValidation)
	}

	validated, innerErr := r.ValidateWithValidation(jsonValueMap, validation.InnerValidation)
	var fieldErrors model.FieldErrors
	if innerErr != nil && !errors.As(innerErr, &fieldErrors) {
		return nil, innerErr
	}

	correct := len(validated)
	err = r.RunValidatorsOnConditionGroup(jsonValueMap, astValue, &correct)
	if err != nil {
		if innerErr != nil {
			return nil, fmt.Errorf("%v, inner errors: %v", err, innerErr)
		}
		return nil, err
	} else if len(jsonValueMap) == 0 {
		return map[string]any{}, nil
	}

	return validated, nil
}

// getRequirementAST returns the cached AST of the validation or parses the requirement.
func getRequirementAST(validation *model.Validation) (*model.AstValue, error) {
	if validation.RequirementAST != nil {
		return validation.RequirementAST, nil
	}

	p := parser.NewParser()
	v, err := p.ParseValidation(validation.Requirement)
	if err != nil {
		return nil, err
	}
	return v.RootValue, nil
}

// ValidateValueWithParser validates a value against a given validation using the parser.
// It parses the validation requirement and runs the validation function on the input value.
//
// It returns an error if the validation fails.
func (r *Validator) ValidateValueWithParser(input any, validation *model.Validation) error {
	astValue, err := getRequirementAST(validation)
	if err != nil {
		return err
	}

	err = r.RunValidatorsOnConditionGroup(input, astValue, nil)
	if err != nil {
		return err
	}

	return nil
}

// RunValidatorsOnConditionGroup runs the validators for each condition in the [astValue] on the [input].
// If the condition is a group, it recursively calls itself on the group.
// If the operator is AND, it returns an error if any condition fails.
// If the operator is OR, it collects all errors and returns them if all conditions fail.
// [correct] is the number of correct inner fields for cmn and cmx, it is nil if the input has no inner validations.
func (r *Validator) RunValidatorsOnConditionGroup(input any, astValue *model.AstValue, correct *int) error {
	var errors []error
	for i, v := range astValue.ConditionGroup {
		var err error
		switch v.Type {
		case model.EMPTY:
			return nil
		case model.GROUP:
			err = r.RunValidatorsOnConditionGroup(input, v, correct)
		case model.CONDITION:
			switch v.ConditionType {
			case model.NONE, model.DEF:
				continue
			case model.EQUAL:
				err = validators.ValidateEqual(input, v)
			case model.NOT_EQUAL:
				err = validators.ValidateNotEqual(input, v)
			case model.MIN_VALUE:
				err = validators.ValidateMin(input, v)
			case model.MAX_VALUE:
				err = validators.ValidateMax(input, v)
			case model.CONTAINS:
				err = validators.ValidateContains(input, v)
			case model.NOT_CONTAINS:
				err = validators.ValidateNotContains(input, v)
			case model.FROM:
				err = validators.ValidateFrom(input, v)
			case model.NOT_FROM:
				err = validators.ValidateNotFrom(input, v)
			case model.REGX:
				err = validators.ValidateRegex(input, v)
			case model.IS_TYPE:
				err = validators.ValidateIsType(input, v)
			case model.CORRECT_MIN, model.CORRECT_MAX:
				if correct == nil {
					err = fmt.Errorf("%v is only usable on struct fields with inner validations", v.ConditionType)
				} else if v.ConditionType == model.CORRECT_MIN {
					err = validators.ValidateCorrectMin(*correct, v)
				} else {
					err = validators.ValidateCorrectMax(*correct, v)
				}
			case model.FUNC:
				fun, ok := r.ValidationFuncs[v.ConditionValue]
				if !ok {
					return fmt.Errorf("unknown validation function: %v", v.ConditionValue)
				}
				err = fun(input, v)
			default:
				return fmt.Errorf("unknown condition type: %v", v.ConditionType)
			}
		}
		if err != nil {
			if (i == 0 && v.Operator == model.OR) || (i > 0 && astValue.ConditionGroup[i-1].Operator == model.OR) {
				errors = append(errors, err)
			} else {
				return err
			}
		}
	}

	if len(astValue.ConditionGroup) > 0 && len(errors) >= len(astValue.ConditionGroup) {
		return fmt.Errorf("no condition fulfilled, all errors: %v", errors)
	}

	return nil
}
