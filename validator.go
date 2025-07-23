package validator

import (
	"github.com/siherrmann/validator/model"
)

type ValidationFunc func(input interface{}, validation *model.Validation) error

type Validator struct {
	ValidationFuncs map[string]ValidationFunc
}

// NewValidator creates a new Validator instance with the default validation functions.
func NewValidator() *Validator {
	return &Validator{
		ValidationFuncs: make(map[string]ValidationFunc),
	}
}

func (r *Validator) AddValidationFunc(fn ValidationFunc, name string) {
	r.ValidationFuncs[name] = fn
}

func (r *Validator) Validate(v any, tagType ...string) error {
	return Validate(v, tagType...)
}

func (r *Validator) ValidateWithValidation(jsonInput model.JsonMap, validations []model.Validation) (model.JsonMap, error) {
	return ValidateWithValidation(jsonInput, validations)
}
