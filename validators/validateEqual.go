package validators

import (
	"fmt"
	"reflect"

	"github.com/siherrmann/validator/helper"
	"github.com/siherrmann/validator/model"
)

func ValidateEqual(v any, ast *model.AstValue) error {
	var check any
	var compare any
	var err error
	switch reflect.TypeOf(v).Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		check = reflect.ValueOf(v).Len()
		compare, err = helper.ConditionValueToT(check, ast.ConditionValue)
		if err != nil {
			return err
		}
	default:
		check = v
		compare, err = helper.ConditionValueToT(v, ast.ConditionValue)
		if err != nil {
			return err
		}
	}

	if !Equal(check, compare) {
		return fmt.Errorf("value not equal condition %v", ast.ConditionValue)
	}
	return nil
}

func ValidateNotEqual(v any, ast *model.AstValue) error {
	var check any
	var compare any
	var err error
	switch reflect.TypeOf(v).Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		check = reflect.ValueOf(v).Len()
		compare, err = helper.ConditionValueToT(check, ast.ConditionValue)
		if err != nil {
			return err
		}
	default:
		check = v
		compare, err = helper.ConditionValueToT(v, ast.ConditionValue)
		if err != nil {
			return err
		}
	}

	if Equal(check, compare) {
		return fmt.Errorf("value equal condition %v", ast.ConditionValue)
	}
	return nil
}
