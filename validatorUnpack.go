package validator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/siherrmann/validator/helper"
	"github.com/siherrmann/validator/model"
)

// UnmapOrUnmarshalAndValidate unmarshals given json ([]byte) or given url.Values (from request.Form),
// validates them and updates the given struct.
// It returns an error if the unmapping or validation fails.
//
// It is actually doing the same as UnmapOrUnmarshalValidateAndUpdate, but in another order.
// It does directly update the struct and validates afterwards.
// Normally you would either only use Validate or use UnmapOrUnmarshalValidateAndUpdate for early return on error.
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

// UnmarshalAndValidate unmarshals given json ([]byte) into pointer v.
// It validates the struct by the given tagType.
// It returns an error if the unmarshaling or validation fails.
//
// It is actually doing the same as UnmarshalValidateAndUpdate, but in another order.
// It does directly update the struct and validates afterwards.
// Normally you would either only use Validate or use UnmarshalValidateAndUpdate for early return on error.
func (r *Validator) UnmarshalAndValidate(request *http.Request, structToValidate any, tagType ...string) error {
	err := helper.CheckValidPointerToStruct(structToValidate)
	if err != nil {
		return err
	}

	var bodyBytes []byte
	bodyBytes, err = io.ReadAll(request.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyBytes, structToValidate)
	if err != nil {
		return fmt.Errorf("error unmarshaling json: %v", err)
	}

	err = r.Validate(structToValidate, tagType...)
	if err != nil {
		return fmt.Errorf("error validating struct: %v", err)
	}

	return nil
}

// UnmapAndValidate unmaps the url.Values from the request.Form into a JsonMap and puts it into the given struct.
// It validates the struct by the given tagType.
// It returns an error if the unmapping or validation fails.
//
// It is actually doing the same as UnmapValidateAndUpdate, but in another order.
// It does directly update the struct and validates afterwards.
// Normally you would either only use Validate or use UnmapValidateAndUpdate for early return on error.
func (r *Validator) UnmapAndValidate(request *http.Request, structToValidate any, tagType ...string) error {
	mapOut, err := helper.UnmapRequestToJsonMap(request)
	if err != nil {
		return fmt.Errorf("error unmapping form values: %v", err)
	}

	err = helper.MapJsonMapToStruct(mapOut, structToValidate)
	if err != nil {
		return fmt.Errorf("error mapping json map to struct: %v", err)
	}

	err = r.Validate(structToValidate, tagType...)
	if err != nil {
		return fmt.Errorf("error validating url values: %v", err)
	}

	return nil
}

// UnmapOrUnmarshalValidateAndUpdate unmarshals given json ([]byte) or given url.Values (from request.Form),
// validates them and updates the given struct.
// It returns an error if the unmapping, validation or update fails.
//
// For more information look at ValidateAndUpdate.
func (r *Validator) UnmapOrUnmarshalValidateAndUpdate(request *http.Request, structToUpdate any, tagType ...string) error {
	err := request.ParseForm()
	if err != nil {
		return err
	}

	if len(request.Form.Encode()) > 0 {
		err = r.UnmapValidateAndUpdate(request, structToUpdate, tagType...)
	} else {
		err = r.UnmarshalValidateAndUpdate(request, structToUpdate, tagType...)
	}

	return err
}

// UnmarshalValidateAndUpdate unmarshals given json ([]byte) into pointer v.
// It returns an error if the unmapping, validation or update fails.
//
// For more information look at ValidateAndUpdate.
func (r *Validator) UnmarshalValidateAndUpdate(request *http.Request, structToUpdate any, tagType ...string) error {
	mapOut, err := helper.UnmarshalRequestToJsonMap(request)
	if err != nil {
		return fmt.Errorf("error unmarshaling request body: %v", err)
	}

	err = r.ValidateAndUpdate(mapOut, structToUpdate, tagType...)
	if err != nil {
		return fmt.Errorf("error updating struct: %v", err)
	}

	return nil
}

// UnmapValidateAndUpdate unmaps given url.Values into pointer jsonMap.
// It returns an error if the unmapping, validation or update fails.
//
// For more information look at ValidateAndUpdate.
func (r *Validator) UnmapValidateAndUpdate(request *http.Request, structToUpdate any, tagType ...string) error {
	mapOut, err := helper.UnmapRequestToJsonMap(request)
	if err != nil {
		return fmt.Errorf("error unmapping form values: %v", err)
	}

	err = r.ValidateAndUpdate(mapOut, structToUpdate, tagType...)
	if err != nil {
		return fmt.Errorf("error updating struct: %v", err)
	}

	return nil
}

// UnmapOrUnmarshalValidateAndUpdateWithValidation unmarshals given json ([]byte) or given url.Values (from request.Form).
// It validates the map with the given validations and updates the given map.
// It returns an error if the unmapping, validation or update fails.
func (r *Validator) UnmapOrUnmarshalValidateAndUpdateWithValidation(request *http.Request, mapToUpdate *map[string]any, validations []model.Validation) error {
	err := request.ParseForm()
	if err != nil {
		return err
	}

	if len(request.Form.Encode()) > 0 {
		err = r.UnmapValidateAndUpdateWithValidation(request, mapToUpdate, validations)
	} else {
		err = r.UnmarshalValidateAndUpdateWithValidation(request, mapToUpdate, validations)
	}

	return err
}

// UnmarshalValidateAndUpdateWithValidation unmarshals given json ([]byte) into pointer mapToUpdate.
// It validates the map by the given validations and updates it.
// It returns an error if the unmarshaling, validation or update fails.
func (r *Validator) UnmarshalValidateAndUpdateWithValidation(request *http.Request, mapToUpdate *map[string]any, validations []model.Validation) error {
	mapOut, err := helper.UnmarshalRequestToJsonMap(request)
	if err != nil {
		return fmt.Errorf("error unmarshaling request body: %v", err)
	}

	err = r.ValidateAndUpdateWithValidation(mapOut, mapToUpdate, validations)
	if err != nil {
		return fmt.Errorf("error updating struct: %v", err)
	}

	return nil
}

// UnmapValidateAndUpdateWithValidation unmaps given url.Values into pointer jsonMap.
// It validates the map by the given validations and updates it.
// It returns an error if the unmapping, validation or update fails.
func (r *Validator) UnmapValidateAndUpdateWithValidation(request *http.Request, mapToUpdate *map[string]any, validations []model.Validation) error {
	mapOut, err := helper.UnmapRequestToJsonMap(request)
	if err != nil {
		return fmt.Errorf("error unmapping form values: %v", err)
	}

	err = r.ValidateAndUpdateWithValidation(mapOut, mapToUpdate, validations)
	if err != nil {
		return fmt.Errorf("error updating struct: %v", err)
	}

	return nil
}
