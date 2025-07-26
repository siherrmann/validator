package validator

import (
	"net/http"
	"net/url"

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
func ValidateWithValidation(jsonInput model.JsonMap, validations []model.Validation) (model.JsonMap, error) {
	r := NewValidator()
	return r.ValidateWithValidation(jsonInput, validations)
}

// ValidateAndUpdate is the wrapper function for the ValidateAndUpdate method of the Validator struct.
// More details can be found in the ValidateAndUpdate method.
func ValidateAndUpdate(jsonInput model.JsonMap, structToUpdate any, tagType ...string) error {
	r := NewValidator()
	return r.ValidateAndUpdate(jsonInput, structToUpdate, tagType...)
}

// ValidateAndUpdateWithValidation is the wrapper function for the ValidateAndUpdateWithValidation method of the Validator struct.
// More details can be found in the ValidateAndUpdateWithValidation method.
func ValidateAndUpdateWithValidation(jsonInput model.JsonMap, mapToUpdate *model.JsonMap, validations []model.Validation) error {
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
func UnmapAndValidate(values url.Values, structToUpdate any, tagType ...string) error {
	r := NewValidator()
	return r.UnmapAndValidate(values, structToUpdate, tagType...)
}

// UnmarshalAndValidate is the wrapper function for the UnmarshalAndValidate method of the Validator struct.
// More details can be found in the UnmarshalAndValidate method.
func UnmarshalAndValidate(data []byte, v any, tagType ...string) error {
	r := NewValidator()
	return r.UnmarshalAndValidate(data, v, tagType...)
}

// UnmapOrUnmarshalValidateAndUpdate is the wrapper function for the UnmapOrUnmarshalValidateAndUpdate method of the Validator struct.
// More details can be found in the UnmapOrUnmarshalValidateAndUpdate method.
func UnmapOrUnmarshalValidateAndUpdate(request *http.Request, structToUpdate any, tagType ...string) error {
	r := NewValidator()
	return r.UnmapOrUnmarshalValidateAndUpdate(request, structToUpdate, tagType...)
}

// UnmapValidateAndUpdate is the wrapper function for the UnmapValidateAndUpdate method of the Validator struct.
// More details can be found in the UnmapValidateAndUpdate method.
func UnmapValidateAndUpdate(values url.Values, structToUpdate any, tagType ...string) error {
	r := NewValidator()
	return r.UnmapValidateAndUpdate(values, structToUpdate, tagType...)
}

// UnmarshalValidateAndUpdate is the wrapper function for the UnmarshalValidateAndUpdate method of the Validator struct.
// More details can be found in the UnmarshalValidateAndUpdate method.
func UnmarshalValidateAndUpdate(jsonInput []byte, structToUpdate any, tagType ...string) error {
	r := NewValidator()
	return r.UnmarshalValidateAndUpdate(jsonInput, structToUpdate, tagType...)
}

// UnmapOrUnmarshalValidateAndUpdateWithValidation is the wrapper function for the UnmapOrUnmarshalValidateAndUpdateWithValidation method of the Validator struct.
// More details can be found in the UnmapOrUnmarshalValidateAndUpdateWithValidation method.
func UnmapOrUnmarshalValidateAndUpdateWithValidation(request *http.Request, mapToUpdate *model.JsonMap, validations []model.Validation) error {
	r := NewValidator()
	return r.UnmapOrUnmarshalValidateAndUpdateWithValidation(request, mapToUpdate, validations)
}

// UnmapValidateAndUpdateWithValidation is the wrapper function for the UnmapValidateAndUpdateWithValidation method of the Validator struct.
// More details can be found in the UnmapValidateAndUpdateWithValidation method.
func UnmapValidateAndUpdateWithValidation(values url.Values, mapToUpdate *model.JsonMap, validations []model.Validation) error {
	r := NewValidator()
	return r.UnmapValidateAndUpdateWithValidation(values, mapToUpdate, validations)
}

// UnmarshalValidateAndUpdateWithValidation is the wrapper function for the UnmarshalValidateAndUpdateWithValidation method of the Validator struct.
// More details can be found in the UnmarshalValidateAndUpdateWithValidation method.
func UnmarshalValidateAndUpdateWithValidation(jsonInput []byte, mapToUpdate *model.JsonMap, validations []model.Validation) error {
	r := NewValidator()
	return r.UnmarshalValidateAndUpdateWithValidation(jsonInput, mapToUpdate, validations)
}
