package validator

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/siherrmann/validator/helper"
	"github.com/siherrmann/validator/model"
	"github.com/siherrmann/validator/parser"
	"github.com/siherrmann/validator/validators"
)

type ValidationFunc func(input any, astValue *model.AstValue) error

// ValidationFuncContext is a ValidationFunc that also gets the context of the validation call,
// eg. to read values injected with [Inject] by using [Get].
type ValidationFuncContext func(ctx context.Context, input any, astValue *model.AstValue) error

// Validator is the main struct for validation.
type Validator struct {
	// ValidationFuncs holds the functions for `fun<name>` by name.
	// Functions added with AddValidationFunc are stored wrapped as ValidationFuncContext.
	ValidationFuncs map[string]ValidationFuncContext
}

// NewValidator creates a new Validator instance with an empty validation functions map.
func NewValidator() *Validator {
	return &Validator{
		ValidationFuncs: make(map[string]ValidationFuncContext),
	}
}

// AddValidationFunc adds a custom validation function to the Validator.
// The function can be used in validation requirements with the name provided (`fun<name>`).
// It replaces a function with the same name. Use AddValidationFuncContext to read injected values.
func (r *Validator) AddValidationFunc(fn ValidationFunc, name string) {
	r.AddValidationFuncContext(func(_ context.Context, input any, astValue *model.AstValue) error {
		return fn(input, astValue)
	}, name)
}

// AddValidationFuncContext adds a custom validation function that gets the context of the validation call.
// The function can be used in validation requirements with the name provided (`fun<name>`),
// values injected with [Inject] can be read in it with [Get].
// It replaces a function with the same name.
func (r *Validator) AddValidationFuncContext(fn ValidationFuncContext, name string) {
	if r.ValidationFuncs == nil {
		r.ValidationFuncs = make(map[string]ValidationFuncContext)
	}
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
	return r.ValidateContext(context.Background(), v, tagType...)
}

// ValidateContext is Validate with the context of the validation call,
// which holds the values injected with [Inject].
func (r *Validator) ValidateContext(ctx context.Context, v any, tagType ...string) error {
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

	_, err = r.ValidateWithValidationContext(ctx, jsonMap, validations)
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
	return r.ValidateAndUpdateContext(context.Background(), jsonInput, structToUpdate, tagType...)
}

// ValidateAndUpdateContext is ValidateAndUpdate with the context of the validation call,
// which holds the values injected with [Inject].
func (r *Validator) ValidateAndUpdateContext(ctx context.Context, jsonInput map[string]any, structToUpdate any, tagType ...string) error {
	tagTypeSet := model.VLD
	if len(tagType) > 0 {
		tagTypeSet = tagType[0]
	}

	validations, err := getCachedValidationsFromStruct(structToUpdate, tagTypeSet)
	if err != nil {
		return fmt.Errorf("error getting validations from struct: %v", err)
	}

	validatedMap, err := r.ValidateWithValidationContext(ctx, jsonInput, validations)
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
	return r.ValidateAndUpdateWithValidationContext(context.Background(), jsonInput, mapToUpdate, validations)
}

// ValidateAndUpdateWithValidationContext is ValidateAndUpdateWithValidation with the context of the validation call,
// which holds the values injected with [Inject].
func (r *Validator) ValidateAndUpdateWithValidationContext(ctx context.Context, jsonInput map[string]any, mapToUpdate *map[string]any, validations []model.Validation) error {
	validatedValues, err := r.ValidateWithValidationContext(ctx, jsonInput, validations)
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
	return r.ValidateWithValidationContext(context.Background(), jsonInput, validations)
}

// ValidateWithValidationContext is ValidateWithValidation with the context of the validation call,
// which holds the values injected with [Inject].
func (r *Validator) ValidateWithValidationContext(ctx context.Context, jsonInput map[string]any, validations []model.Validation) (map[string]any, error) {
	return r.validateWithScope(ctx, jsonInput, validations, nil)
}

// validateWithScope is ValidateWithValidation for a level below the given parent scope (nil for the top level).
//
// The scope of this level is only created if a validation of this level (or below) uses other fields,
// so validations without `rel` or `$` do not pay for it.
func (r *Validator) validateWithScope(ctx context.Context, jsonInput map[string]any, validations []model.Validation, parent *scope) (map[string]any, error) {
	var sc *scope
	for i := range validations {
		if validations[i].UsesFields() {
			sc = &scope{current: jsonInput, parent: parent}
			break
		}
	}

	groups := map[string]*model.Group{}
	groupSize := map[string]int{}
	groupErrors := map[string][]error{}
	fieldErrors := model.FieldErrors{}

	validateValues := map[string]any{}

	for validationIndex := range validations {
		validation := validations[validationIndex]
		if len(validation.Key) > 0 && hasKeyBefore(validations, validationIndex) {
			return map[string]any{}, fmt.Errorf("duplicate validation key: %v", validation.Key)
		}

		for _, g := range validation.Groups {
			groups[g.Name] = g
			groupSize[g.Name]++
		}

		jsonValue, ok, err := r.validateField(ctx, jsonInput, &validation, sc)
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

// hasKeyBefore reports if a validation before index i has the same key as validation i.
func hasKeyBefore(validations []model.Validation, i int) bool {
	for j := 0; j < i; j++ {
		if validations[j].Key == validations[i].Key {
			return true
		}
	}
	return false
}

// validateField validates the value of a single validation in the given JsonMap.
// It returns the validated value and if the value should be added to the validated JsonMap.
//
// A requirement that depends on other fields (`rel` or `$`) is also checked if the key is missing
// or the value is empty with omitempty, with the value null (or ""). This makes conditional requirements
// like `rel(type equoci) => min1` possible. The field is then not added to the validated JsonMap.
func (r *Validator) validateField(ctx context.Context, jsonInput map[string]any, validation *model.Validation, sc *scope) (any, bool, error) {
	skipResult := false
	jsonValue, ok := jsonInput[validation.Key]
	if !ok {
		if validation.Default != "" {
			jsonValue = validation.Default
		} else if strings.TrimSpace(validation.Requirement) == string(model.NONE) {
			return nil, false, nil
		} else if requirementUsesFields(validation) {
			skipResult = true
		} else {
			return nil, false, fmt.Errorf("json %v key not in map", validation.Key)
		}
	}
	if validation.OmitEmpty && jsonValue == "" {
		if !requirementUsesFields(validation) {
			return nil, false, nil
		}
		skipResult = true
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
			jsonValue, err = r.validateStruct(ctx, jsonValueMap, validation, sc)
			if err != nil {
				return nil, false, fmt.Errorf("field %v invalid: %v", validation.Key, err.Error())
			}
		} else {
			err = r.validateValue(ctx, jsonValue, validation, sc)
		}
	case model.Array:
		if helper.IsArray(jsonValue) && len(validation.InnerValidation) > 0 {
			jsonArray, ok := jsonValue.([]any)
			if !ok {
				return nil, false, fmt.Errorf("field %v must be of type array, was %T", validation.Key, jsonValue)
			}

			// The requirement of the list itself (eg. `min1` or `unqname`).
			if err := r.validateValue(ctx, jsonArray, validation, sc); err != nil {
				if validation.Default != "" {
					return validation.Default, true, nil
				}
				return nil, false, fmt.Errorf("field %v invalid: %v", validation.Key, err.Error())
			}

			validatedArray := []any{}
			for _, jsonValueInner := range jsonArray {
				jsonValueInnerMap, err := helper.GetValidMap(jsonValueInner)
				if err != nil {
					return nil, false, fmt.Errorf("field %v invalid: %v", validation.Key, err.Error())
				}

				validatedInnerMap, err := r.validateWithScope(ctx, jsonValueInnerMap, validation.InnerValidation, sc)
				if err != nil {
					return nil, false, fmt.Errorf("field %v invalid: %v", validation.Key, err.Error())
				}
				validatedArray = append(validatedArray, validatedInnerMap)
			}
			jsonValue = validatedArray
		} else if helper.IsArray(jsonValue) {
			err = r.validateValue(ctx, jsonValue, validation, sc)
		} else if helper.IsString(jsonValue) {
			// Check if the value is a string from a url value.
			jsonValue = []string{jsonValue.(string)}
			err = r.validateValue(ctx, jsonValue, validation, sc)
		} else if jsonValue == nil {
			err = r.validateValue(ctx, jsonValue, validation, sc)
		} else if len(validation.InnerValidation) == 0 {
			// A lone value (e.g. a single form value) is treated as a one-element list.
			jsonValue = []any{jsonValue}
			err = r.validateValue(ctx, jsonValue, validation, sc)
		}
	default:
		err = r.validateValue(ctx, jsonValue, validation, sc)
	}

	if err != nil {
		if validation.Default != "" {
			return validation.Default, true, nil
		}
		return nil, false, fmt.Errorf("field %v invalid: %v", validation.Key, err.Error())
	}
	if skipResult {
		return nil, false, nil
	}

	return jsonValue, true, nil
}

// requirementUsesFields reports if the requirement of the validation contains `rel` or a `$` reference.
// It is only called for missing or empty values, so parsing an uncached requirement here is fine.
func requirementUsesFields(validation *model.Validation) bool {
	astValue, err := getRequirementAST(validation)
	return err == nil && astValue.UsesFields
}

// validateStruct validates the JsonMap of a struct field by its inner validations.
//
// Without a requirement on the struct field all inner fields have to be valid.
// With a requirement, the requirement is checked on the JsonMap itself (eg. min and max check the number of keys).
// If the JsonMap is empty and fulfills the requirement, the inner validations are skipped.
// If the requirement contains cmn or cmx, invalid inner fields are left out instead of failing
// and cmn/cmx check the number of correct inner fields (the fields in the validated JsonMap).
func (r *Validator) validateStruct(ctx context.Context, jsonValueMap map[string]any, validation *model.Validation, sc *scope) (map[string]any, error) {
	astValue, err := getRequirementAST(validation)
	if err != nil {
		return nil, err
	}

	if !astValue.HasCondition(model.ConditionType.IsRequirement) {
		return r.validateWithScope(ctx, jsonValueMap, validation.InnerValidation, sc)
	}

	if !astValue.HasCondition(model.ConditionType.IsCorrect) {
		err = r.runConditions(ctx, jsonValueMap, astValue, nil, sc)
		if err != nil {
			return nil, err
		} else if len(jsonValueMap) == 0 {
			return map[string]any{}, nil
		}
		return r.validateWithScope(ctx, jsonValueMap, validation.InnerValidation, sc)
	}

	validated, innerErr := r.validateWithScope(ctx, jsonValueMap, validation.InnerValidation, sc)
	var fieldErrors model.FieldErrors
	if innerErr != nil && !errors.As(innerErr, &fieldErrors) {
		return nil, innerErr
	}

	correct := len(validated)
	err = r.runConditions(ctx, jsonValueMap, astValue, &correct, sc)
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
// A requirement with `rel` or `$` references can not be validated with this function,
// because the other fields are unknown. Use ValidateWithValidation instead.
//
// It returns an error if the validation fails.
func (r *Validator) ValidateValueWithParser(input any, validation *model.Validation) error {
	return r.validateValue(context.Background(), input, validation, nil)
}

// validateValue validates a value against a given validation with the scope of its level.
func (r *Validator) validateValue(ctx context.Context, input any, validation *model.Validation, sc *scope) error {
	astValue, err := getRequirementAST(validation)
	if err != nil {
		return err
	}

	return r.runConditions(ctx, input, astValue, nil, sc)
}

// RunValidatorsOnConditionGroup runs the validators for each condition in the [astValue] on the [input].
// If the condition is a group, it recursively calls itself on the group.
// If the operator is AND, it returns an error if any condition fails.
// If the operator is OR, it collects all errors and returns them if all conditions fail.
// If the operator is IMPLIES (`=>`), the group is fulfilled if the premise fails, otherwise the consequence has to be fulfilled.
// [correct] is the number of correct inner fields for cmn and cmx, it is nil if the input has no inner validations.
//
// A requirement with `rel` or `$` references can not be validated with this function,
// because the other fields are unknown. Use ValidateWithValidation instead.
func (r *Validator) RunValidatorsOnConditionGroup(input any, astValue *model.AstValue, correct *int) error {
	return r.runConditions(context.Background(), input, astValue, correct, nil)
}

// runConditions is RunValidatorsOnConditionGroup with the scope of the validated level.
func (r *Validator) runConditions(ctx context.Context, input any, astValue *model.AstValue, correct *int, sc *scope) error {
	onlyOr := isOnlyOr(astValue.ConditionGroup)
	var errors []error
	for i, v := range astValue.ConditionGroup {
		var err error
		switch v.Type {
		case model.EMPTY:
			return nil
		case model.GROUP:
			err = r.runConditions(ctx, input, v, correct, sc)
		case model.CONDITION:
			var fatal bool
			if v.Ref != nil {
				fatal, err = r.runFieldCondition(ctx, input, v, correct, sc)
			} else {
				fatal, err = r.runCondition(ctx, input, v, correct)
			}
			if fatal {
				return err
			}
		}
		if err == nil && onlyOr {
			// One fulfilled condition is enough, the other conditions do not have to be checked.
			return nil
		}
		if err != nil {
			if v.Operator == model.IMPLIES {
				// The premise is not fulfilled, so the consequence does not have to be.
				return nil
			} else if i > 0 && astValue.ConditionGroup[i-1].Operator == model.IMPLIES {
				return fmt.Errorf("%v (required because of %v)", err, premiseToString(astValue.ConditionGroup[i-1]))
			} else if (i == 0 && v.Operator == model.OR) || (i > 0 && astValue.ConditionGroup[i-1].Operator == model.OR) {
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

// isOnlyOr reports if all conditions of the group are connected with OR.
func isOnlyOr(conditions model.ConditionGroup) bool {
	if len(conditions) < 2 {
		return false
	}
	for _, v := range conditions[:len(conditions)-1] {
		if v.Operator != model.OR {
			return false
		}
	}
	return true
}

// runCondition runs the validator of a single condition (without `rel` or `$` reference).
// fatal reports a configuration error (eg. an unknown function), that is returned also inside an OR group.
func (r *Validator) runCondition(ctx context.Context, input any, v *model.AstValue, correct *int) (fatal bool, err error) {
	switch v.ConditionType {
	case model.NONE, model.DEF:
		return false, nil
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
	case model.UNIQUE:
		err = validators.ValidateUnique(input, v)
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
			return true, fmt.Errorf("unknown validation function: %v", v.ConditionValue)
		}
		err = fun(ctx, input, v)
	default:
		return true, fmt.Errorf("unknown condition type: %v", v.ConditionType)
	}
	return false, err
}
