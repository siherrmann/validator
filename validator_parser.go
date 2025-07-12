package validator

import (
	"fmt"

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

	var checkFunction func(T, *model.AstValue) error
	switch validation.Type {
	case model.String:
		checkFunction = validators.CheckString
	case model.Int:
		checkFunction = validators.CheckInt
	case model.Float:
		checkFunction = validators.CheckFloat
	case model.Bool:
		checkFunction = validators.CheckBool
	case model.Array:
		checkFunction = validators.CheckArray
	case model.Map, model.Struct:
		checkFunction = validators.CheckMap
	case model.Time, model.TimeISO8601, model.TimeUnix:
		checkFunction = validators.CheckTime
	default:
		return fmt.Errorf("invalid validation type: %v", validation.Type)
	}

	err = model.RunFuncOnConditionGroup(input, r.RootValue, checkFunction)
	if err != nil {
		return err
	}

	return nil
}
