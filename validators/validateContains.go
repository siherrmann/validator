package validators

import (
	"fmt"
	"reflect"

	"github.com/siherrmann/validator/model"
)

func ValidateContains[T any](v T, ast *model.AstValue) error {
	switch reflect.TypeOf(v).Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		contains, err := Contains(v, ast.ConditionValue)
		if err != nil {
			return fmt.Errorf("error checking contains: %v", err)
		}
		if !contains {
			return fmt.Errorf("value does not contain %v", ast.ConditionValue)
		}
	default:
		return fmt.Errorf("value to validate has to be a string, array, slice or map, was %v", reflect.TypeOf(v).Kind())
	}

	return nil
}

func ValidateNotContains[T any](v T, ast *model.AstValue) error {
	switch reflect.TypeOf(v).Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		contains, err := Contains(v, ast.ConditionValue)
		if err != nil {
			return fmt.Errorf("error checking not contains: %v", err)
		}
		if contains {
			return fmt.Errorf("value contains condition %v", ast.ConditionValue)
		}
	default:
		return fmt.Errorf("value to validate has to be a string, array, slice or map, was %v", reflect.TypeOf(v).Kind())
	}

	return nil
}
