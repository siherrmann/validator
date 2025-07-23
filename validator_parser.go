package validator

import (
	"fmt"
	"slices"
	"strings"

	"github.com/siherrmann/validator/helper"
	"github.com/siherrmann/validator/model"
	"github.com/siherrmann/validator/parser"
	"github.com/siherrmann/validator/validators"
)

// ValidateWithValidation validates a given JsonMap by the given validations.
// It checks if the keys are in the map, validates the values and returns a new JsonMap.
//
// If a validation has groups, it checks if the values are valid for the groups.
// If a validation has a key that is already in the map, it returns an error.
//
// It returns a new JsonMap with the validated values or an error if the validation fails.
func ValidateWithValidation(jsonInput model.JsonMap, validations []model.Validation) (model.JsonMap, error) {
	keys := []string{}
	groups := map[string]*model.Group{}
	groupSize := map[string]int{}
	groupErrors := map[string][]error{}

	validateValues := model.JsonMap{}

	for _, validation := range validations {
		if len(validation.Key) > 0 && slices.Contains(keys, validation.Key) {
			return model.JsonMap{}, fmt.Errorf("duplicate validation key: %v", validation.Key)
		} else {
			keys = append(keys, validation.Key)
		}

		for _, g := range validation.Groups {
			groups[g.Name] = g
			groupSize[g.Name]++
		}

		var ok bool
		var jsonValue interface{}
		if jsonValue, ok = jsonInput[validation.Key]; !ok {
			if strings.TrimSpace(validation.Requirement) == string(model.NONE) {
				continue
			} else if len(validation.Groups) == 0 {
				return model.JsonMap{}, fmt.Errorf("json %v key not in map", validation.Key)
			} else {
				for _, group := range validation.Groups {
					groupErrors[group.Name] = append(groupErrors[group.Name], fmt.Errorf("json %v key not in map", validation.Key))
				}
				continue
			}
		}

		var err error
		switch validation.Type {
		case model.Struct:
			jsonValueInner, err := GetValidMap(jsonValue)
			if err != nil {
				return model.JsonMap{}, fmt.Errorf("field %v invalid: %v", validation.Key, err.Error())
			}

			jsonValue, err = ValidateWithValidation(jsonValueInner, validation.InnerValidation)
			if err != nil {
				return model.JsonMap{}, fmt.Errorf("field %v invalid: %v", validation.Key, err.Error())
			}
		case model.Array:
			if helper.IsArray(jsonValue) && len(validation.InnerValidation) > 0 {
				jsonArray, ok := jsonValue.([]interface{})
				if !ok {
					return model.JsonMap{}, fmt.Errorf("field %v must be of type array, was %T", validation.Key, jsonValue)
				}

				for _, jsonValueInner := range jsonArray {
					jsonValueInnerMap, err := GetValidMap(jsonValueInner)
					if err != nil {
						return model.JsonMap{}, fmt.Errorf("field %v invalid: %v", validation.Key, err.Error())
					}

					_, err = ValidateWithValidation(jsonValueInnerMap, validation.InnerValidation)
					if err != nil {
						return model.JsonMap{}, fmt.Errorf("field %v invalid: %v", validation.Key, err.Error())
					}
				}
			} else if helper.IsArray(jsonValue) {
				err = ValidateValueWithParser(jsonValue, &validation)
			} else if helper.IsString(jsonValue) {
				// Check if the value is a string from a url value.
				jsonValue = []string{jsonValue.(string)}
				err = ValidateValueWithParser(jsonValue, &validation)
			}
		default:
			err = ValidateValueWithParser(jsonValue, &validation)
		}

		if err != nil && len(validation.Groups) == 0 {
			return model.JsonMap{}, fmt.Errorf("field %v invalid: %v", validation.Key, err.Error())
		} else if err != nil {
			for _, group := range validation.Groups {
				groupErrors[group.Name] = append(groupErrors[group.Name], fmt.Errorf("field %v invalid: %v", validation.Key, err.Error()))
			}
			continue
		}

		validateValues[validation.Key] = jsonValue
	}

	err := validators.ValidateGroups(groups, groupSize, groupErrors)
	if err != nil {
		return model.JsonMap{}, err
	}

	return validateValues, nil
}

// ValidateValueWithParser validates a value against a given validation using the parser.
// It parses the validation requirement and runs the validation function on the input value.
//
// It returns an error if the validation fails.
func ValidateValueWithParser[T comparable](input T, validation *model.Validation) error {
	p := parser.NewParser()
	r, err := p.ParseValidation(validation.Requirement)
	if err != nil {
		return err
	}

	err = validators.RunFuncOnConditionGroup(input, r.RootValue)
	if err != nil {
		return err
	}

	return nil
}
