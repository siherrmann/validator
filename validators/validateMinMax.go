package validators

import (
	"fmt"
	"strconv"

	"github.com/siherrmann/validator/helper"
	"github.com/siherrmann/validator/model"
)

func ValidateMin(v any, ast *model.AstValue) error {
	check, err := helper.AnyToFloat(v)
	if err != nil {
		return fmt.Errorf("invalid value for min validation: %v", err)
	}
	compare, err := minMaxCompare(check, ast)
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
	compare, err := minMaxCompare(check, ast)
	if err != nil {
		return err
	}

	if !Max(check, compare) {
		return fmt.Errorf("value greater than maximum condition %v", ast.ConditionValue)
	}
	return nil
}

// minMaxCompare returns the condition value of min and max as float64.
// It uses the value parsed by the parser and only parses it, if the AST was built without the parser.
func minMaxCompare(check float64, ast *model.AstValue) (float64, error) {
	if ast.HasNumber {
		return ast.Number, nil
	}
	return helper.ConditionValueToT(check, ast.ConditionValue)
}

// ValidateCorrectMin checks if at least the condition value of inner fields are correct.
func ValidateCorrectMin(correct int, ast *model.AstValue) error {
	compare, err := strconv.Atoi(ast.ConditionValue)
	if err != nil {
		return fmt.Errorf("invalid condition value for cmn validation: %v", err)
	}

	if correct < compare {
		return fmt.Errorf("less than %v correct inner fields, got %v", ast.ConditionValue, correct)
	}
	return nil
}

// ValidateCorrectMax checks if at most the condition value of inner fields are correct.
func ValidateCorrectMax(correct int, ast *model.AstValue) error {
	compare, err := strconv.Atoi(ast.ConditionValue)
	if err != nil {
		return fmt.Errorf("invalid condition value for cmx validation: %v", err)
	}

	if correct > compare {
		return fmt.Errorf("more than %v correct inner fields, got %v", ast.ConditionValue, correct)
	}
	return nil
}
