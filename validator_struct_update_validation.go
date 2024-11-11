package validator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"slices"
	"strings"

	"github.com/siherrmann/validator/model"
	"github.com/siherrmann/validator/validators"
)

// UnmapOrAnmarshalValidateAndUpdate unmarshals given json ([]byte) or given url.Values (from request.Form),
// validates them and updates the given map.
func UnmapOrUnmarshalRequestValidateAndUpdateWithValidation(request *http.Request, mapToUpdate *model.JsonMap, validations map[string]model.Validation) error {
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
func UnmarshalValidateAndUpdateWithValidation(jsonInput []byte, mapToUpdate *model.JsonMap, validations map[string]model.Validation) error {
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
func UnmapValidateAndUpdateWithValidation(values url.Values, mapToUpdate *model.JsonMap, validations map[string]model.Validation) error {
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
func ValidateAndUpdateWithValidation(jsonInput model.JsonMap, mapToUpdate *model.JsonMap, validations map[string]model.Validation) error {
	keys := []string{}
	groups := map[string]*model.Group{}
	groupSize := map[string]int{}
	groupErrors := map[string][]error{}

	for key, validation := range validations {
		if len(validation.Key) > 0 && slices.Contains(keys, validation.Key) {
			return fmt.Errorf("duplicate validation key: %v", validation.Key)
		} else {
			keys = append(keys, validation.Key)
		}

		for _, g := range validation.Groups {
			groups[g.Name] = g
			groupSize[g.Name]++
		}

		var ok bool
		var jsonValue interface{}
		if jsonValue, ok = jsonInput[validation.Key]; !ok {
			if strings.TrimSpace(validation.Requirement) == string(model.NONE) {
				continue
			} else if len(validation.Groups) == 0 {
				return fmt.Errorf("json %v key not in map", validation.Key)
			} else {
				for _, group := range validation.Groups {
					groupErrors[group.Name] = append(groupErrors[group.Name], fmt.Errorf("json %v key not in map", validation.Key))
				}
				continue
			}
		}

		var validatedValue interface{}
		var err error
		if validation.Type == model.Struct {
			return fmt.Errorf("field %v unsupported type %v", key, validation.Type)
		} else {
			validatedValue, err = ValidateValueWithParser(reflect.ValueOf(jsonValue), &validation)
		}

		if err != nil && len(validation.Groups) == 0 {
			return fmt.Errorf("field %v invalid: %v", key, err.Error())
		} else if err != nil {
			for _, group := range validation.Groups {
				groupErrors[group.Name] = append(groupErrors[group.Name], fmt.Errorf("field %v invalid: %v", key, err.Error()))
			}
			continue
		}

		(*mapToUpdate)[key] = validatedValue
	}

	err := validators.ValidateGroups(groups, groupSize, groupErrors)
	if err != nil {
		return err
	}

	return nil
}
