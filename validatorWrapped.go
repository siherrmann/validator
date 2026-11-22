package validator

import (
	"context"
	"net/http"

	"github.com/siherrmann/validator/model"
)

// Validate is the wrapper function for the Validate method of the Validator struct.
// More details can be found in the Validate method.
func Validate(v any, tagType ...string) error {
	r := NewValidator()
	return r.Validate(v, tagType...)
}

// ValidateWithValidation is the wrapper function for the ValidateWithValidation method of the Validator struct.
// More details can be found in the ValidateWithValidation method.
func ValidateWithValidation(jsonInput map[string]any, validations []model.Validation) (map[string]any, error) {
	r := NewValidator()
	return r.ValidateWithValidation(jsonInput, validations)
}

// ValidateAndUpdate is the wrapper function for the ValidateAndUpdate method of the Validator struct.
// More details can be found in the ValidateAndUpdate method.
func ValidateAndUpdate(jsonInput map[string]any, structToUpdate any, tagType ...string) error {
	r := NewValidator()
	return r.ValidateAndUpdate(jsonInput, structToUpdate, tagType...)
}

// ValidateAndUpdateWithValidation is the wrapper function for the ValidateAndUpdateWithValidation method of the Validator struct.
// More details can be found in the ValidateAndUpdateWithValidation method.
func ValidateAndUpdateWithValidation(jsonInput map[string]any, mapToUpdate *map[string]any, validations []model.Validation) error {
	r := NewValidator()
	return r.ValidateAndUpdateWithValidation(jsonInput, mapToUpdate, validations)
}

// UnmapOrUnmarshalAndValidate is the wrapper function for the UnmapOrUnmarshalAndValidate method of the Validator struct.
// More details can be found in the UnmapOrUnmarshalAndValidate method.
func UnmapOrUnmarshalAndValidate(request *http.Request, structToUpdate any, tagType ...string) error {
	r := NewValidator()
	return r.UnmapOrUnmarshalAndValidate(request, structToUpdate, tagType...)
}

// UnmapAndValidate is the wrapper function for the UnmapAndValidate method of the Validator struct.
// More details can be found in the UnmapAndValidate method.
func UnmapAndValidate(request *http.Request, structToUpdate any, tagType ...string) error {
	r := NewValidator()
	return r.UnmapAndValidate(request, structToUpdate, tagType...)
}

// UnmarshalAndValidate is the wrapper function for the UnmarshalAndValidate method of the Validator struct.
// More details can be found in the UnmarshalAndValidate method.
func UnmarshalAndValidate(request *http.Request, structToUpdate any, tagType ...string) error {
	r := NewValidator()
	return r.UnmarshalAndValidate(request, structToUpdate, tagType...)
}

// UnmapOrUnmarshalAndValidateWithValidation is the wrapper function for the UnmapOrUnmarshalAndValidateWithValidation method of the Validator struct.
// More details can be found in the UnmapOrUnmarshalAndValidateWithValidation method.
func UnmapOrUnmarshalAndValidateWithValidation(request *http.Request, mapToUpdate *map[string]any, validations []model.Validation) error {
	r := NewValidator()
	return r.UnmapOrUnmarshalAndValidateWithValidation(request, mapToUpdate, validations)
}

// UnmapAndValidateWithValidation is the wrapper function for the UnmapAndValidateWithValidation method of the Validator struct.
// More details can be found in the UnmapAndValidateWithValidation method.
func UnmapAndValidateWithValidation(request *http.Request, mapToUpdate *map[string]any, validations []model.Validation) error {
	r := NewValidator()
	return r.UnmapAndValidateWithValidation(request, mapToUpdate, validations)
}

// UnmarshalAndValidateWithValidation is the wrapper function for the UnmarshalAndValidateWithValidation method of the Validator struct.
// More details can be found in the UnmarshalAndValidateWithValidation method.
func UnmarshalAndValidateWithValidation(request *http.Request, mapToUpdate *map[string]any, validations []model.Validation) error {
	r := NewValidator()
	return r.UnmarshalAndValidateWithValidation(request, mapToUpdate, validations)
}

// UnmarshalYamlAndValidate is the wrapper function for the UnmarshalYamlAndValidate method of the Validator struct.
// More details can be found in the UnmarshalYamlAndValidate method.
func UnmarshalYamlAndValidate(yamlInput []byte, structToUpdate any, tagType ...string) error {
	r := NewValidator()
	return r.UnmarshalYamlAndValidate(yamlInput, structToUpdate, tagType...)
}

// UnmarshalYamlAndValidateWithValidation is the wrapper function for the UnmarshalYamlAndValidateWithValidation method of the Validator struct.
// More details can be found in the UnmarshalYamlAndValidateWithValidation method.
func UnmarshalYamlAndValidateWithValidation(yamlInput []byte, mapToUpdate *map[string]any, validations []model.Validation) error {
	r := NewValidator()
	return r.UnmarshalYamlAndValidateWithValidation(yamlInput, mapToUpdate, validations)
}

// ValidateContext is the wrapper function for the ValidateContext method of the Validator struct.
// More details can be found in the ValidateContext method.
func ValidateContext(ctx context.Context, v any, tagType ...string) error {
	r := NewValidator()
	return r.ValidateContext(ctx, v, tagType...)
}

// ValidateWithValidationContext is the wrapper function for the ValidateWithValidationContext method of the Validator struct.
// More details can be found in the ValidateWithValidationContext method.
func ValidateWithValidationContext(ctx context.Context, jsonInput map[string]any, validations []model.Validation) (map[string]any, error) {
	r := NewValidator()
	return r.ValidateWithValidationContext(ctx, jsonInput, validations)
}

// ValidateAndUpdateContext is the wrapper function for the ValidateAndUpdateContext method of the Validator struct.
// More details can be found in the ValidateAndUpdateContext method.
func ValidateAndUpdateContext(ctx context.Context, jsonInput map[string]any, structToUpdate any, tagType ...string) error {
	r := NewValidator()
	return r.ValidateAndUpdateContext(ctx, jsonInput, structToUpdate, tagType...)
}

// ValidateAndUpdateWithValidationContext is the wrapper function for the ValidateAndUpdateWithValidationContext method of the Validator struct.
// More details can be found in the ValidateAndUpdateWithValidationContext method.
func ValidateAndUpdateWithValidationContext(ctx context.Context, jsonInput map[string]any, mapToUpdate *map[string]any, validations []model.Validation) error {
	r := NewValidator()
	return r.ValidateAndUpdateWithValidationContext(ctx, jsonInput, mapToUpdate, validations)
}

// UnmarshalYamlAndValidateContext is the wrapper function for the UnmarshalYamlAndValidateContext method of the Validator struct.
// More details can be found in the UnmarshalYamlAndValidateContext method.
func UnmarshalYamlAndValidateContext(ctx context.Context, yamlInput []byte, structToUpdate any, tagType ...string) error {
	r := NewValidator()
	return r.UnmarshalYamlAndValidateContext(ctx, yamlInput, structToUpdate, tagType...)
}

// UnmarshalYamlAndValidateWithValidationContext is the wrapper function for the UnmarshalYamlAndValidateWithValidationContext method of the Validator struct.
// More details can be found in the UnmarshalYamlAndValidateWithValidationContext method.
func UnmarshalYamlAndValidateWithValidationContext(ctx context.Context, yamlInput []byte, mapToUpdate *map[string]any, validations []model.Validation) error {
	r := NewValidator()
	return r.UnmarshalYamlAndValidateWithValidationContext(ctx, yamlInput, mapToUpdate, validations)
}
