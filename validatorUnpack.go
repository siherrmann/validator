package validator

import (
	"fmt"
	"net/http"

	"github.com/siherrmann/validator/helper"
	"github.com/siherrmann/validator/model"
)

// All unpack functions validate with the context of the request (request.Context()),
// so values injected with Inject into the request context can be used in requirements and functions.

// UnmapOrUnmarshalAndValidate unpacks the url.Values (from request.Form) or, if there are none, the json body of the request,
// validates them and updates the given struct.
// It returns an error if the unmapping, validation or update fails.
//
// For more information look at ValidateAndUpdate.
func (r *Validator) UnmapOrUnmarshalAndValidate(request *http.Request, structToUpdate any, tagType ...string) error {
	err := request.ParseForm()
	if err != nil {
		return err
	}

	if len(request.Form.Encode()) > 0 {
		err = r.UnmapAndValidate(request, structToUpdate, tagType...)
	} else {
		err = r.UnmarshalAndValidate(request, structToUpdate, tagType...)
	}

	return err
}

// UnmarshalAndValidate unmarshals the json body of the request, validates it and updates the given struct.
// It returns an error if the unmapping, validation or update fails.
//
// For more information look at ValidateAndUpdate.
func (r *Validator) UnmarshalAndValidate(request *http.Request, structToUpdate any, tagType ...string) error {
	mapOut, err := helper.UnmarshalRequestToJsonMap(request)
	if err != nil {
		return fmt.Errorf("error unmarshaling request body: %v", err)
	}

	err = r.ValidateAndUpdateContext(request.Context(), mapOut, structToUpdate, tagType...)
	if err != nil {
		return fmt.Errorf("error updating struct: %v", err)
	}

	return nil
}

// UnmapAndValidate unmaps the url.Values (from request.Form), validates them and updates the given struct.
// It returns an error if the unmapping, validation or update fails.
//
// For more information look at ValidateAndUpdate.
func (r *Validator) UnmapAndValidate(request *http.Request, structToUpdate any, tagType ...string) error {
	mapOut, err := helper.UnmapRequestToJsonMap(request)
	if err != nil {
		return fmt.Errorf("error unmapping form values: %v", err)
	}

	err = r.ValidateAndUpdateContext(request.Context(), mapOut, structToUpdate, tagType...)
	if err != nil {
		return fmt.Errorf("error updating struct: %v", err)
	}

	return nil
}

// UnmapOrUnmarshalAndValidateWithValidation unmarshals given json ([]byte) or given url.Values (from request.Form).
// It validates the map with the given validations and updates the given map.
// It returns an error if the unmapping, validation or update fails.
func (r *Validator) UnmapOrUnmarshalAndValidateWithValidation(request *http.Request, mapToUpdate *map[string]any, validations []model.Validation) error {
	err := request.ParseForm()
	if err != nil {
		return err
	}

	if len(request.Form.Encode()) > 0 {
		err = r.UnmapAndValidateWithValidation(request, mapToUpdate, validations)
	} else {
		err = r.UnmarshalAndValidateWithValidation(request, mapToUpdate, validations)
	}

	return err
}

// UnmarshalAndValidateWithValidation unmarshals the json body of the request.
// It validates the map by the given validations and updates it.
// It returns an error if the unmarshaling, validation or update fails.
func (r *Validator) UnmarshalAndValidateWithValidation(request *http.Request, mapToUpdate *map[string]any, validations []model.Validation) error {
	mapOut, err := helper.UnmarshalRequestToJsonMap(request)
	if err != nil {
		return fmt.Errorf("error unmarshaling request body: %v", err)
	}

	err = r.ValidateAndUpdateWithValidationContext(request.Context(), mapOut, mapToUpdate, validations)
	if err != nil {
		return fmt.Errorf("error updating struct: %v", err)
	}

	return nil
}

// UnmapAndValidateWithValidation unmaps the url.Values (from request.Form).
// It validates the map by the given validations and updates it.
// It returns an error if the unmapping, validation or update fails.
func (r *Validator) UnmapAndValidateWithValidation(request *http.Request, mapToUpdate *map[string]any, validations []model.Validation) error {
	mapOut, err := helper.UnmapRequestToJsonMap(request)
	if err != nil {
		return fmt.Errorf("error unmapping form values: %v", err)
	}

	err = r.ValidateAndUpdateWithValidationContext(request.Context(), mapOut, mapToUpdate, validations)
	if err != nil {
		return fmt.Errorf("error updating struct: %v", err)
	}

	return nil
}
