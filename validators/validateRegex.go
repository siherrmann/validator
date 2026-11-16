package validators

import (
	"fmt"
	"reflect"

	"github.com/siherrmann/validator/helper"
	"github.com/siherrmann/validator/model"
)

func ValidateRegex(v any, ast *model.AstValue) error {
	// A null value has nothing to match
	if v == nil {
		return fmt.Errorf("null value does not match regex %v", ast.ConditionValue)
	}

	switch reflect.TypeOf(v).Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		checks, err := helper.AnyToArrayOfString(v)
		if err != nil {
			return err
		}
		for _, check := range checks {
			match := matchRegex(ast, check)
			if !match {
				return fmt.Errorf("value %v does not match regex %v", check, ast.ConditionValue)
			}
		}
	default:
		check, err := helper.AnyToString(v)
		if err != nil {
			return fmt.Errorf("error converting value to string: %v", err)
		}
		match := matchRegex(ast, check)
		if !match {
			return fmt.Errorf("value %v does not match regex %v", check, ast.ConditionValue)
		}
	}

	return nil
}

// matchRegex uses the regex compiled by the parser and only compiles it, if the AST was built without the parser.
func matchRegex(ast *model.AstValue, s string) bool {
	if ast.Regex != nil {
		return ast.Regex.MatchString(s)
	}
	return Regex(s, ast.ConditionValue)
}
