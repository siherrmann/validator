package validators

import (
	"fmt"
	"reflect"

	"github.com/siherrmann/validator/model"
)

func ValidateContains[T any](v T, ast *model.AstValue) error {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Map:
		var vComparable map[string]Comparable
		var ok bool
		if vComparable, ok = any(v).(map[string]Comparable); !ok {
			return fmt.Errorf("value has to be of type map of Comparable, was %v", reflect.TypeOf(v).Kind())
		}

		contains := MapContains(vComparable, ast.ConditionValue)
		if !contains {
			return fmt.Errorf("value does not contain key %v", ast.ConditionValue)
		}
	case reflect.Array, reflect.Slice:
		var vComparable []Comparable
		var ok bool
		fmt.Printf("Type of v: %T\n", v)
		if vComparable, ok = any(v).([]Comparable); !ok {
			return fmt.Errorf("value has to be of type array or slice of Comparable, was %v", reflect.TypeOf(v).Kind())
		}

		compare, err := ArrayConditionValueToT(vComparable, ast)
		if err != nil {
			return err
		}

		contains := ArrayContains(vComparable, compare)
		if !contains {
			return fmt.Errorf("value does not contain condition %v", ast.ConditionValue)
		}
	case reflect.String:
		check := any(v).(string)
		if !StringContains(check, ast.ConditionValue) {
			return fmt.Errorf("value does not contain condition %v", ast.ConditionValue)
		}
	default:
		return fmt.Errorf("value to validate has to be a string, array, slice or map, was %v", reflect.TypeOf(v).Kind())
	}

	return nil
}
