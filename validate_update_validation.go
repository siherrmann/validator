package validator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/siherrmann/validator/model"
)

// UnmapOrAnmarshalValidateAndUpdate unmarshals given json ([]byte) or given url.Values (from request.Form),
// validates them and updates the given map.
func UnmapOrUnmarshalRequestValidateAndUpdateWithValidation(request *http.Request, mapToUpdate *model.JsonMap, validations []model.Validation) error {
	err := request.ParseForm()
	if err != nil {
		return err
	}

	if len(request.Form.Encode()) > 0 {
		err = UnmapValidateAndUpdateWithValidation(request.Form, mapToUpdate, validations)
	} else {
		var bodyBytes []byte
		bodyBytes, err = io.ReadAll(request.Body)
		if err != nil {
			return err
		}
		err = UnmarshalValidateAndUpdateWithValidation(bodyBytes, mapToUpdate, validations)
	}

	return err
}

// UnmarshalValidateAndUpdateWithValidation unmarshals given json ([]byte) into pointer mapToUpdate.
// For more information to ValidateAndUpdate look at ValidateAndUpdate(jsonInput model.JsonMap, structToUpdate interface{}) error.
func UnmarshalValidateAndUpdateWithValidation(jsonInput []byte, mapToUpdate *model.JsonMap, validations []model.Validation) error {
	jsonUnmarshaled := model.JsonMap{}

	err := json.Unmarshal(jsonInput, &jsonUnmarshaled)
	if err != nil {
		return fmt.Errorf("error unmarshaling: %v", err)
	}

	err = ValidateAndUpdateWithValidation(jsonUnmarshaled, mapToUpdate, validations)
	if err != nil {
		return fmt.Errorf("error updating struct: %v", err)
	}

	return nil
}

// UnmapValidateAndUpdateWithValidation unmaps given url.Values into pointer jsonMap.
// For more information to ValidateAndUpdate look at ValidateAndUpdate(jsonInput model.JsonMap, structToUpdate interface{}) error.
func UnmapValidateAndUpdateWithValidation(values url.Values, mapToUpdate *model.JsonMap, validations []model.Validation) error {
	mapOut, err := UnmapUrlValuesToJsonMap(values)
	if err != nil {
		return err
	}

	err = ValidateAndUpdateWithValidation(mapOut, mapToUpdate, validations)
	if err != nil {
		return fmt.Errorf("error updating struct: %v", err)
	}

	return nil
}

// ValidateAndUpdateWithValidation validates a given struct by upd tags.
// ValidateAndUpdateWithValidation needs a struct pointer, a json map as input and a validation map.
// The given struct is updated by the values in the json map.
func ValidateAndUpdateWithValidation(jsonInput model.JsonMap, mapToUpdate *model.JsonMap, validations []model.Validation) error {
	validatedValues, err := ValidateWithValidation(jsonInput, validations)
	if err != nil {
		return fmt.Errorf("error validating json map in ValidateAndUpdateWithValidation: %v", err)
	}

	err = UpdateJsonMap(validatedValues, mapToUpdate)
	if err != nil {
		return fmt.Errorf("error updating json map in ValidateAndUpdateWithValidation: %v", err)
	}

	return nil
}
