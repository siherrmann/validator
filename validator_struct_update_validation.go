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
// validates them and updates the given struct.
func UnmapOrUnmarshalRequestValidateAndUpdateWithValidation(request *http.Request, structToUpdate interface{}, validations map[string]model.Validation) error {
	err := request.ParseForm()
	if err != nil {
		return err
	}

	if len(request.Form.Encode()) > 0 {
		err = UnmapValidateAndUpdateWithValidation(request.Form, structToUpdate, validations)
	} else {
		var bodyBytes []byte
		bodyBytes, err = io.ReadAll(request.Body)
		if err != nil {
			return err
		}
		err = UnmarshalValidateAndUpdateWithValidation(bodyBytes, structToUpdate, validations)
	}

	return err
}

// UnmarshalValidateAndUpdateWithValidation unmarshals given json ([]byte) into pointer v.
// For more information to ValidateAndUpdate look at ValidateAndUpdate(jsonInput model.JsonMap, structToUpdate interface{}) error.
func UnmarshalValidateAndUpdateWithValidation(jsonInput []byte, structToUpdate interface{}, validations map[string]model.Validation) error {
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

// UnmapValidateAndUpdateWithValidation unmaps given url.Values into pointer jsonMap.
// For more information to ValidateAndUpdate look at ValidateAndUpdate(jsonInput model.JsonMap, structToUpdate interface{}) error.
func UnmapValidateAndUpdateWithValidation(values url.Values, structToUpdate interface{}, validations map[string]model.Validation) error {
	mapOut, err := UnmapUrlValuesToJsonMap(values)
	if err != nil {
		return err
	}

	err = ValidateAndUpdate(mapOut, structToUpdate)
	if err != nil {
		return fmt.Errorf("error updating struct: %v", err)
	}

	return nil
}

// ValidateAndUpdateWithValidation validates a given struct by upd tags.
// ValidateAndUpdateWithValidation needs a struct pointer, a json map as input and a validation map.
// The given struct is updated by the values in the json map.
func ValidateAndUpdateWithValidation(jsonInput model.JsonMap, structToUpdate interface{}, validations map[string]model.Validation) error {
	// check if value is a pointer to a struct
	value := reflect.ValueOf(structToUpdate)
	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("value has to be of kind pointer, was %T", value)
	}
	if value.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("value has to be of kind struct, was %T", value)
	}

	// get valid reflect value of struct
	structFull := value.Elem()

	keys := []string{}
	groups := map[string]*model.Group{}
	groupSize := map[string]int{}
	groupErrors := map[string][]error{}

	for i := 0; i < structFull.Type().NumField(); i++ {
		tag := structFull.Type().Field(i).Tag.Get(string(model.UPD))
		field := structFull.Field(i)
		fieldName := structFull.Type().Field(i).Name

		validation := &model.Validation{}
		err := validation.FillOnlyKey(tag, model.UPD, field)
		if err != nil {
			return err
		}

		if val, ok := validations[validation.Key]; ok {
			validation.Type = val.Type
			validation.Requirement = val.Requirement
			validation.Groups = val.Groups
		}

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
		if validation.Type == model.Struct {
			var validMap interface{}
			validMap, err = validation.GetValidValue(jsonValue)
			if err != nil {
				return err
			}

			err = ValidateAndUpdate(validMap.(map[string]interface{}), field.Addr().Interface())
			validatedValue = field.Interface()
		} else {
			validatedValue, err = ValidateValueWithParser(reflect.ValueOf(jsonValue), validation)
		}

		if err != nil && len(validation.Groups) == 0 {
			return fmt.Errorf("field %v of %v invalid: %v", fieldName, reflect.TypeOf(structToUpdate), err.Error())
		} else if err != nil {
			for _, group := range validation.Groups {
				groupErrors[group.Name] = append(groupErrors[group.Name], fmt.Errorf("field %v of %v invalid: %v", fieldName, reflect.TypeOf(structToUpdate), err.Error()))
			}
			continue
		}

		err = setStructValueByJson(field, validation.Key, validatedValue)
		if err != nil && len(validation.Groups) == 0 {
			return fmt.Errorf("could not set field %v of %v: %v", fieldName, reflect.TypeOf(structToUpdate), err.Error())
		} else if err != nil {
			for _, group := range groups {
				groupErrors[group.Name] = append(groupErrors[group.Name], fmt.Errorf("could not set field %v: %v", fieldName, err.Error()))
			}
			continue
		}
	}

	err := validators.ValidateGroups(groups, groupSize, groupErrors)
	if err != nil {
		return err
	}

	return nil
}
