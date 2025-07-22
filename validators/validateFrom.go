package validators

import (
	"fmt"

	"github.com/siherrmann/validator/model"
)

func ValidateFrom[T any](v T, ast *model.AstValue) error {
	from, err := From(v, ast.ConditionValue, false)
	if err != nil {
		return fmt.Errorf("error checking from: %v", err)
	}
	if !from {
		return fmt.Errorf("from %v does not contain value", ast.ConditionValue)
	}

	return nil
}

func ValidateNotFrom[T any](v T, ast *model.AstValue) error {
	notFrom, err := From(v, ast.ConditionValue, true)
	if err != nil {
		return fmt.Errorf("error checking not from: %v", err)
	}
	if !notFrom {
		return fmt.Errorf("from %v does contain value", ast.ConditionValue)
	}

	return nil
}
