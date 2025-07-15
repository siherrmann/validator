package validators

import (
	"fmt"
	"reflect"

	"github.com/siherrmann/validator/model"
)

func ValidateEqual[T Comparable | Array | Map](v T, ast *model.AstValue) error {
	var check interface{}
	var compare interface{}
	var err error
	switch reflect.TypeOf(v).Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		check = reflect.ValueOf(v).Len()
		compare, err = ConditionValueToT(check, ast)
		if err != nil {
			return err
		}
	default:
		check = v
		compare, err = ConditionValueToT(v, ast)
		if err != nil {
			return err
		}
	}

	if !Equal(check, compare) {
		return fmt.Errorf("value not equal condition %v", ast.ConditionValue)
	}
	return nil
}

func ValidateNotEqual[T Comparable | Array | Map](v T, ast *model.AstValue) error {
	var check interface{}
	var compare interface{}
	var err error
	switch reflect.TypeOf(v).Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		check = reflect.ValueOf(v).Len()
		compare, err = ConditionValueToT(check, ast)
		if err != nil {
			return err
		}
	default:
		check = v
		compare, err = ConditionValueToT(v, ast)
		if err != nil {
			return err
		}
	}

	if Equal(check, compare) {
		return fmt.Errorf("value equal condition %v", ast.ConditionValue)
	}
	return nil
}
