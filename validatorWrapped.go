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
func UnmarshalAndValidate(request *http.Request, v any, tagType ...string) error {
	r := NewValidator()
	return r.UnmarshalAndValidate(request, v, tagType...)
}

// UnmapOrUnmarshalValidateAndUpdate is the wrapper function for the UnmapOrUnmarshalValidateAndUpdate method of the Validator struct.
// More details can be found in the UnmapOrUnmarshalValidateAndUpdate method.
func UnmapOrUnmarshalValidateAndUpdate(request *http.Request, structToUpdate any, tagType ...string) error {
	r := NewValidator()
	return r.UnmapOrUnmarshalValidateAndUpdate(request, structToUpdate, tagType...)
}

// UnmapValidateAndUpdate is the wrapper function for the UnmapValidateAndUpdate method of the Validator struct.
// More details can be found in the UnmapValidateAndUpdate method.
func UnmapValidateAndUpdate(request *http.Request, structToUpdate any, tagType ...string) error {
	r := NewValidator()
	return r.UnmapValidateAndUpdate(request, structToUpdate, tagType...)
}

// UnmarshalValidateAndUpdate is the wrapper function for the UnmarshalValidateAndUpdate method of the Validator struct.
// More details can be found in the UnmarshalValidateAndUpdate method.
func UnmarshalValidateAndUpdate(request *http.Request, structToUpdate any, tagType ...string) error {
	r := NewValidator()
	return r.UnmarshalValidateAndUpdate(request, structToUpdate, tagType...)
}

// UnmapOrUnmarshalValidateAndUpdateWithValidation is the wrapper function for the UnmapOrUnmarshalValidateAndUpdateWithValidation method of the Validator struct.
// More details can be found in the UnmapOrUnmarshalValidateAndUpdateWithValidation method.
func UnmapOrUnmarshalValidateAndUpdateWithValidation(request *http.Request, mapToUpdate *map[string]any, validations []model.Validation) error {
	r := NewValidator()
	return r.UnmapOrUnmarshalValidateAndUpdateWithValidation(request, mapToUpdate, validations)
}

// UnmapValidateAndUpdateWithValidation is the wrapper function for the UnmapValidateAndUpdateWithValidation method of the Validator struct.
// More details can be found in the UnmapValidateAndUpdateWithValidation method.
func UnmapValidateAndUpdateWithValidation(request *http.Request, mapToUpdate *map[string]any, validations []model.Validation) error {
	r := NewValidator()
	return r.UnmapValidateAndUpdateWithValidation(request, mapToUpdate, validations)
}

// UnmarshalValidateAndUpdateWithValidation is the wrapper function for the UnmarshalValidateAndUpdateWithValidation method of the Validator struct.
// More details can be found in the UnmarshalValidateAndUpdateWithValidation method.
func UnmarshalValidateAndUpdateWithValidation(request *http.Request, mapToUpdate *map[string]any, validations []model.Validation) error {
	r := NewValidator()
	return r.UnmarshalValidateAndUpdateWithValidation(request, mapToUpdate, validations)
}

// UnmarshalYamlValidateAndUpdate is the wrapper function for the UnmarshalYamlValidateAndUpdate method of the Validator struct.
// More details can be found in the UnmarshalYamlValidateAndUpdate method.
func UnmarshalYamlValidateAndUpdate(yamlInput []byte, structToUpdate any, tagType ...string) error {
	r := NewValidator()
	return r.UnmarshalYamlValidateAndUpdate(yamlInput, structToUpdate, tagType...)
}

// UnmarshalYamlValidateAndUpdateWithValidation is the wrapper function for the UnmarshalYamlValidateAndUpdateWithValidation method of the Validator struct.
// More details can be found in the UnmarshalYamlValidateAndUpdateWithValidation method.
func UnmarshalYamlValidateAndUpdateWithValidation(yamlInput []byte, mapToUpdate *map[string]any, validations []model.Validation) error {
	r := NewValidator()
	return r.UnmarshalYamlValidateAndUpdateWithValidation(yamlInput, mapToUpdate, validations)
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

// UnmarshalYamlValidateAndUpdateContext is the wrapper function for the UnmarshalYamlValidateAndUpdateContext method of the Validator struct.
// More details can be found in the UnmarshalYamlValidateAndUpdateContext method.
func UnmarshalYamlValidateAndUpdateContext(ctx context.Context, yamlInput []byte, structToUpdate any, tagType ...string) error {
	r := NewValidator()
	return r.UnmarshalYamlValidateAndUpdateContext(ctx, yamlInput, structToUpdate, tagType...)
}

// UnmarshalYamlValidateAndUpdateWithValidationContext is the wrapper function for the UnmarshalYamlValidateAndUpdateWithValidationContext method of the Validator struct.
// More details can be found in the UnmarshalYamlValidateAndUpdateWithValidationContext method.
func UnmarshalYamlValidateAndUpdateWithValidationContext(ctx context.Context, yamlInput []byte, mapToUpdate *map[string]any, validations []model.Validation) error {
	r := NewValidator()
	return r.UnmarshalYamlValidateAndUpdateWithValidationContext(ctx, yamlInput, mapToUpdate, validations)
}
