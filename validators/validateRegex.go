package validators

import (
	"fmt"
	"reflect"

	"github.com/siherrmann/validator/model"
)

func ValidateRegex[T Comparable | Array | Map](v T, ast *model.AstValue) error {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		checks, err := ValueToArrayOfString(v)
		if err != nil {
			return err
		}
		for _, check := range checks {
			match := Regex(check, ast.ConditionValue)
			if !match {
				return fmt.Errorf("value %v does not match regex %v", check, ast.ConditionValue)
			}
		}
	default:
		check, err := ValueToString(v)
		if err != nil {
			return fmt.Errorf("error converting value to string: %v", err)
		}
		match := Regex(check, ast.ConditionValue)
		if !match {
			return fmt.Errorf("value %v does not match regex %v", check, ast.ConditionValue)
		}
	}

	return nil
}
