package validator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

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
		err = r.UnmapAndValidate(request.Form, structToUpdate, tagType...)
	} else {
		var bodyBytes []byte
		bodyBytes, err = io.ReadAll(request.Body)
		if err != nil {
			return err
		}
		err = r.UnmarshalAndValidate(bodyBytes, structToUpdate, tagType...)
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
func (r *Validator) UnmarshalAndValidate(jsonInput []byte, structToValidate any, tagType ...string) error {
	err := helper.CheckValidPointerToStruct(structToValidate)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonInput, structToValidate)
	if err != nil {
		return fmt.Errorf("error unmarshaling json: %v", err)
	}

	err = r.Validate(structToValidate, tagType...)
	if err != nil {
		return fmt.Errorf("error validating struct: %v", err)
	}

	return nil
}

// UnmapAndValidate unmaps given url.Values into a JsonMap and puts it into the given struct.
// It validates the struct by the given tagType.
// It returns an error if the unmapping or validation fails.
//
// It is actually doing the same as UnmapValidateAndUpdate, but in another order.
// It does directly update the struct and validates afterwards.
// Normally you would either only use Validate or use UnmapValidateAndUpdate for early return on error.
func (r *Validator) UnmapAndValidate(values url.Values, structToValidate any, tagType ...string) error {
	err := helper.CheckValidPointerToStruct(structToValidate)
	if err != nil {
		return err
	}

	mapOut, err := UnmapUrlValuesToJsonMap(values)
	if err != nil {
		return fmt.Errorf("error unmapping form values: %v", err)
	}

	err = MapJsonMapToStruct(mapOut, structToValidate)
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
		err = r.UnmapValidateAndUpdate(request.Form, structToUpdate, tagType...)
	} else {
		var bodyBytes []byte
		bodyBytes, err = io.ReadAll(request.Body)
		if err != nil {
			return err
		}
		err = r.UnmarshalValidateAndUpdate(bodyBytes, structToUpdate, tagType...)
	}

	return err
}

// UnmarshalValidateAndUpdate unmarshals given json ([]byte) into pointer v.
// It returns an error if the unmapping, validation or update fails.
//
// For more information look at ValidateAndUpdate.
func (r *Validator) UnmarshalValidateAndUpdate(jsonInput []byte, structToUpdate any, tagType ...string) error {
	jsonUnmarshaled := model.JsonMap{}

	err := json.Unmarshal(jsonInput, &jsonUnmarshaled)
	if err != nil {
		return fmt.Errorf("error unmarshaling request body: %v", err)
	}

	err = r.ValidateAndUpdate(jsonUnmarshaled, structToUpdate, tagType...)
	if err != nil {
		return fmt.Errorf("error updating struct: %v", err)
	}

	return nil
}

// UnmapValidateAndUpdate unmaps given url.Values into pointer jsonMap.
// It returns an error if the unmapping, validation or update fails.
//
// For more information look at ValidateAndUpdate.
func (r *Validator) UnmapValidateAndUpdate(values url.Values, structToUpdate any, tagType ...string) error {
	mapOut, err := UnmapUrlValuesToJsonMap(values)
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
func (r *Validator) UnmapOrUnmarshalValidateAndUpdateWithValidation(request *http.Request, mapToUpdate *model.JsonMap, validations []model.Validation) error {
	err := request.ParseForm()
	if err != nil {
		return err
	}

	if len(request.Form.Encode()) > 0 {
		err = r.UnmapValidateAndUpdateWithValidation(request.Form, mapToUpdate, validations)
	} else {
		var bodyBytes []byte
		bodyBytes, err = io.ReadAll(request.Body)
		if err != nil {
			return err
		}
		err = r.UnmarshalValidateAndUpdateWithValidation(bodyBytes, mapToUpdate, validations)
	}

	return err
}

// UnmarshalValidateAndUpdateWithValidation unmarshals given json ([]byte) into pointer mapToUpdate.
// It validates the map by the given validations and updates it.
// It returns an error if the unmarshaling, validation or update fails.
func (r *Validator) UnmarshalValidateAndUpdateWithValidation(jsonInput []byte, mapToUpdate *model.JsonMap, validations []model.Validation) error {
	jsonUnmarshaled := model.JsonMap{}

	err := json.Unmarshal(jsonInput, &jsonUnmarshaled)
	if err != nil {
		return fmt.Errorf("error unmarshaling: %v", err)
	}

	err = r.ValidateAndUpdateWithValidation(jsonUnmarshaled, mapToUpdate, validations)
	if err != nil {
		return fmt.Errorf("error updating struct: %v", err)
	}

	return nil
}

// UnmapValidateAndUpdateWithValidation unmaps given url.Values into pointer jsonMap.
// It validates the map by the given validations and updates it.
// It returns an error if the unmapping, validation or update fails.
func (r *Validator) UnmapValidateAndUpdateWithValidation(values url.Values, mapToUpdate *model.JsonMap, validations []model.Validation) error {
	mapOut, err := UnmapUrlValuesToJsonMap(values)
	if err != nil {
		return err
	}

	err = r.ValidateAndUpdateWithValidation(mapOut, mapToUpdate, validations)
	if err != nil {
		return fmt.Errorf("error updating struct: %v", err)
	}

	return nil
}
