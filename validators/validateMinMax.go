package validators

import (
	"fmt"

	"github.com/siherrmann/validator/model"
)

func ValidateMin[T Comparable | Array | Map](v T, ast *model.AstValue) error {
	check, err := ValueToFloat(v)
	if err != nil {
		return fmt.Errorf("invalid value for min validation: %v", err)
	}
	compare, err := ConditionValueToT(check, ast)
	if err != nil {
		return err
	}

	if !Min(check, compare) {
		return fmt.Errorf("value less than minimum condition %v", ast.ConditionValue)
	}
	return nil
}

func ValidateMax[T Comparable | Array | Map](v T, ast *model.AstValue) error {
	check, err := ValueToFloat(v)
	if err != nil {
		return fmt.Errorf("invalid value for max validation: %v", err)
	}
	compare, err := ConditionValueToT(check, ast)
	if err != nil {
		return err
	}

	if !Max(check, compare) {
		return fmt.Errorf("value greater than maximum condition %v", ast.ConditionValue)
	}
	return nil
}
