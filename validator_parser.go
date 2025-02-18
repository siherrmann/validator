package validator

import (
	"fmt"
	"reflect"

	"github.com/siherrmann/validator/model"
	"github.com/siherrmann/validator/parser"
	"github.com/siherrmann/validator/validators"
)

func ValidateValueWithParser(input reflect.Value, validation *model.Validation) (interface{}, error) {
	validValue, err := validation.GetValidValue(input.Interface())
	if err != nil {
		return nil, err
	}

	lexer := parser.NewLexer(validation.Requirement)
	p := parser.NewParser(lexer)
	r, err := p.ParseValidation()
	if err != nil {
		return nil, err
	}

	var checkFunction func(reflect.Value, *model.AstValue) error
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
		return nil, fmt.Errorf("invalid validation type: %v", validation.Type)
	}

	err = r.RootValue.RunFuncOnConditionGroup(reflect.ValueOf(validValue), checkFunction)
	if err != nil {
		return nil, err
	}
	return validValue, nil
}
