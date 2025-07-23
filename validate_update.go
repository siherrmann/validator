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

// UnmapOrAnmarshalValidateAndUpdate unmarshals given json ([]byte) or given url.Values (from request.Form),
// validates them and updates the given struct.
func UnmapOrUnmarshalRequestValidateAndUpdate(request *http.Request, structToUpdate interface{}, tagType ...string) error {
	err := request.ParseForm()
	if err != nil {
		return err
	}

	if len(request.Form.Encode()) > 0 {
		err = UnmapValidateAndUpdate(request.Form, structToUpdate, tagType...)
	} else {
		var bodyBytes []byte
		bodyBytes, err = io.ReadAll(request.Body)
		if err != nil {
			return err
		}
		err = UnmarshalValidateAndUpdate(bodyBytes, structToUpdate, tagType...)
	}

	return err
}

// UnmarshalValidateAndUpdate unmarshals given json ([]byte) into pointer v.
// For more information to ValidateAndUpdate look at ValidateAndUpdate(jsonInput model.JsonMap, structToUpdate interface{}) error.
func UnmarshalValidateAndUpdate(jsonInput []byte, structToUpdate interface{}, tagType ...string) error {
	jsonUnmarshaled := model.JsonMap{}

	err := json.Unmarshal(jsonInput, &jsonUnmarshaled)
	if err != nil {
		return fmt.Errorf("error unmarshaling: %v", err)
	}

	err = ValidateAndUpdate(jsonUnmarshaled, structToUpdate, tagType...)
	if err != nil {
		return fmt.Errorf("error updating struct: %v", err)
	}

	return nil
}

// UnmapValidateAndUpdate unmaps given url.Values into pointer jsonMap.
// For more information to ValidateAndUpdate look at ValidateAndUpdate(jsonInput model.JsonMap, structToUpdate interface{}) error.
func UnmapValidateAndUpdate(values url.Values, structToUpdate interface{}, tagType ...string) error {
	mapOut, err := UnmapUrlValuesToJsonMap(values)
	if err != nil {
		return err
	}

	err = ValidateAndUpdate(mapOut, structToUpdate, tagType...)
	if err != nil {
		return fmt.Errorf("error updating struct: %v", err)
	}

	return nil
}

// ValidateAndUpdate validates a given JsonMap by the given validations and updates the struct.
// It checks if the keys are in the map, validates the values and returns a new JsonMap.
func ValidateAndUpdate(jsonInput model.JsonMap, structToUpdate interface{}, tagType ...string) error {
	tagTypeSet := model.VLD
	if len(tagType) > 0 {
		tagTypeSet = tagType[0]
	}

	err := helper.CheckValidPointerToStruct(structToUpdate)
	if err != nil {
		return err
	}

	validations, err := model.GetValidationsFromStruct(structToUpdate, tagTypeSet)
	if err != nil {
		return fmt.Errorf("error getting validations from struct: %v", err)
	}

	validatedMap, err := ValidateWithValidation(jsonInput, validations)
	if err != nil {
		return fmt.Errorf("error validating struct: %v", err)
	}

	err = MapJsonMapToStruct(validatedMap, structToUpdate)
	if err != nil {
		return fmt.Errorf("error mapping json map to struct: %v", err)
	}

	return nil
}
