package validator

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/siherrmann/validator/helper"
	"github.com/siherrmann/validator/model"
)

type StructValue struct {
	Error  error
	Groups []string
}

// UnmarshalAndValidate unmarshals given json ([]byte) into pointer v
// and validates it with `Validate(v any, tagType ...string) error`.
func UnmarshalAndValidate(data []byte, v any, tagType ...string) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("value has to be of kind pointer, was %T", value)
	}
	if value.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("value has to be of kind struct, was %T", value)
	}

	err := json.Unmarshal(data, v)
	if err != nil {
		return fmt.Errorf("error unmarshalling %T: %v", value, err)
	}

	err = Validate(v, tagType...)
	if err != nil {
		return err
	}

	return nil
}

// Validate validates a given struct by the given tagType.
// It checks if the keys are in the struct, validates the values
// and returns an error if the validation fails.
func Validate(v any, tagType ...string) error {
	tagTypeSet := model.VLD
	if len(tagType) > 0 {
		tagTypeSet = tagType[0]
	}

	err := helper.CheckValidPointerToStruct(v)
	if err != nil {
		return err
	}

	jsonMap := model.JsonMap{}
	err = UnmapStructToJsonMap(v, &jsonMap)
	if err != nil {
		return fmt.Errorf("error unmapping struct to json map: %v", err)
	}

	validations, err := model.GetValidationsFromStruct(v, tagTypeSet)
	if err != nil {
		return fmt.Errorf("error getting validations from struct: %v", err)
	}

	_, err = ValidateWithValidation(jsonMap, validations)
	if err != nil {
		return fmt.Errorf("error validating struct: %v", err)
	}

	return nil
}
