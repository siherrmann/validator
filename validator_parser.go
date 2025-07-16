package validator

import (
	"github.com/siherrmann/validator/model"
	"github.com/siherrmann/validator/parser"
	"github.com/siherrmann/validator/validators"
)

func ValidateValueWithParser[T comparable](input T, validation *model.Validation) error {
	lexer := parser.NewLexer(validation.Requirement)
	p := parser.NewParser(lexer)
	r, err := p.ParseValidation()
	if err != nil {
		return err
	}

	err = validators.RunFuncOnConditionGroup(input, r.RootValue)
	if err != nil {
		return err
	}

	return nil
}
