package validator

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/siherrmann/validator/helper"
	"github.com/siherrmann/validator/model"
)

// UnmapOrAnmarshalValidateAndUpdate unmarshals given json ([]byte) or given url.Values (from request.Form),
// validates them and updates the given struct.
func UnmapOrUnmarshalRequestValidateAndUpdate(request *http.Request, structToUpdate interface{}) error {
	err := request.ParseForm()
	if err != nil {
		return err
	}

	if len(request.Form.Encode()) > 0 {
		err = UnmapValidateAndUpdate(request.Form, structToUpdate)
	} else {
		var bodyBytes []byte
		bodyBytes, err = io.ReadAll(request.Body)
		if err != nil {
			return err
		}
		err = UnmarshalValidateAndUpdate(bodyBytes, structToUpdate)
	}

	return err
}

// UnmarshalValidateAndUpdate unmarshals given json ([]byte) into pointer v.
// For more information to ValidateAndUpdate look at ValidateAndUpdate(jsonInput model.JsonMap, structToUpdate interface{}) error.
func UnmarshalValidateAndUpdate(jsonInput []byte, structToUpdate interface{}) error {
	jsonUnmarshaled := model.JsonMap{}

	err := json.Unmarshal(jsonInput, &jsonUnmarshaled)
	if err != nil {
		return fmt.Errorf("error unmarshaling: %v", err)
	}

	err = ValidateAndUpdate(jsonUnmarshaled, structToUpdate)
	if err != nil {
		return fmt.Errorf("error updating struct: %v", err)
	}

	return nil
}

// UnmapValidateAndUpdate unmaps given url.Values into pointer jsonMap.
// For more information to ValidateAndUpdate look at ValidateAndUpdate(jsonInput model.JsonMap, structToUpdate interface{}) error.
func UnmapValidateAndUpdate(values url.Values, structToUpdate interface{}) error {
	mapOut, err := UnmapUrlValuesToJsonMap(values)
	if err != nil {
		return err
	}

	log.Printf("unmapped values: %v", mapOut)

	err = ValidateAndUpdate(mapOut, structToUpdate)
	if err != nil {
		return fmt.Errorf("error updating struct: %v", err)
	}

	return nil
}

// ValidateAndUpdate validates a given struct by upd tags.
// ValidateAndUpdate needs a struct pointer and a json map as input.
// The given struct is updated by the values in the json map.
//
// All fields in the struct need a upd tag.
// The tag has to contain the key value for the json struct.
// If no tag is present the field in the struct is ignored and does not get updated.
//
// The second part of the tag contains the conditions for the validation.
//
// If you want to use multiple conditions you can add them with a space in between them.
//
// A complex example for password would be:
// `upd:"password, min8 max30 rex^(.*[A-Z])+(.*)$ rex^(.*[a-z])+(.*)$ rex^(.*\\d)+(.*)$ rex^(.*[\x60!@#$%^&*()_+={};':\"|\\,.<>/?~-])+(.*)$"`
//
// If you want don't want to validate the field you can add `upd:"json_key, -"`.
// If you don't add the upd tag to every field the function will fail with an error.
//
// Conditions have different usages per variable type:
//
// equ - int/float/string == condition, len(array) == condition
//
// neq - int/float/string != condition, len(array) != condition
//
// min - int/float >= condition, len(string/array) >= condition
//
// max - int/float <= condition, len(string/array) <= condition
//
// con - strings.Contains(string, condition), contains(array, condition), int/float ignored
//
// rex - regexp.MatchString(condition, int/float/string), array ignored
//
// For con you need to put in a condition that is convertable to the underlying type of the arrary.
// Eg. for an array of int the condition must be convertable to int (bad: `upd:"array, conA"`, good: `upd:"array, con1"`).
//
// In the case of rex the int and float input will get converted to a string (strconv.Itoa(int) and fmt.Sprintf("%f", f)).
// If you want to check more complex cases you can obviously replace equ, neq, min, max and con with one regular expression.
func ValidateAndUpdate(jsonInput model.JsonMap, structToUpdate interface{}) error {
	err := helper.CheckValidPointerToStruct(structToUpdate)
	if err != nil {
		return err
	}

	validations, err := model.GetValidationsFromStruct(structToUpdate, string(model.UPD))
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
