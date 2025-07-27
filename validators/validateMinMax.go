package validators

import (
	"fmt"

	"github.com/siherrmann/validator/helper"
	"github.com/siherrmann/validator/model"
)

func ValidateMin(v any, ast *model.AstValue) error {
	check, err := helper.AnyToFloat(v)
	if err != nil {
		return fmt.Errorf("invalid value for min validation: %v", err)
	}
	compare, err := helper.ConditionValueToT(check, ast.ConditionValue)
	if err != nil {
		return err
	}

	if !Min(check, compare) {
		return fmt.Errorf("value less than minimum condition %v", ast.ConditionValue)
	}
	return nil
}

func ValidateMax(v any, ast *model.AstValue) error {
	check, err := helper.AnyToFloat(v)
	if err != nil {
		return fmt.Errorf("invalid value for max validation: %v", err)
	}
	compare, err := helper.ConditionValueToT(check, ast.ConditionValue)
	if err != nil {
		return err
	}

	if !Max(check, compare) {
		return fmt.Errorf("value greater than maximum condition %v", ast.ConditionValue)
	}
	return nil
}
