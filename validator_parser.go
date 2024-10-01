package validator

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/siherrmann/validator/model"
	"github.com/siherrmann/validator/parser"
)

func ValidateValueWithoutParser(input reflect.Value, conditions []string, or bool, checkFunction func(reflect.Value, *model.AstValue) error) error {
	if slices.Contains(conditions, string(model.NONE)) || len(conditions) == 0 {
		return nil
	}

	var errors []error
	for _, conFull := range conditions {
		conType, err := model.GetConditionType(conFull)
		if err != nil {
			return err
		}

		conValue, err := model.GetConditionByType(conFull, conType)
		if err != nil {
			return err
		}

		err = checkFunction(input, &model.AstValue{ConditionType: conType, ConditionValue: conValue})
		if err != nil && or {
			errors = append(errors, err)
		} else if err != nil {
			return err
		}
	}

	if len(errors) >= len(conditions) {
		return fmt.Errorf("no condition fulfilled, all errors: %v", errors)
	}

	return nil
}

func ValidateValueWithParser(input reflect.Value, requirement string, checkFunction func(reflect.Value, *model.AstValue) error) error {
	lexer := parser.NewLexer(requirement)
	p := parser.NewParser(lexer)
	r, err := p.ParseValidation()
	if err != nil {
		return err
	}

	err = r.RootValue.RunFuncOnConditionGroup(input, checkFunction)
	if err != nil {
		return err
	}
	return nil
}
