package validator

import (
	"fmt"

	"github.com/siherrmann/validator/helper"
	"github.com/siherrmann/validator/model"
)

// UnmarshalYamlValidateAndUpdate unmarshals given yaml ([]byte) into a JsonMap.
// It validates the map by the tagType of the given struct and updates the struct.
// It returns an error if the unmarshaling, validation or update fails.
//
// The keys are taken from the json tag, fall back to the yaml tag and then to the field name.
// For more information look at ValidateAndUpdate.
func (r *Validator) UnmarshalYamlValidateAndUpdate(yamlInput []byte, structToUpdate any, tagType ...string) error {
	mapOut, err := helper.UnmarshalYamlToJsonMap(yamlInput)
	if err != nil {
		return fmt.Errorf("error unmarshaling yaml: %v", err)
	}

	err = r.ValidateAndUpdate(mapOut, structToUpdate, tagType...)
	if err != nil {
		return fmt.Errorf("error updating struct: %v", err)
	}

	return nil
}

// UnmarshalYamlValidateAndUpdateWithValidation unmarshals given yaml ([]byte) into a JsonMap.
// It validates the map by the given validations and updates the given map.
// It returns an error if the unmarshaling, validation or update fails.
func (r *Validator) UnmarshalYamlValidateAndUpdateWithValidation(yamlInput []byte, mapToUpdate *map[string]any, validations []model.Validation) error {
	mapOut, err := helper.UnmarshalYamlToJsonMap(yamlInput)
	if err != nil {
		return fmt.Errorf("error unmarshaling yaml: %v", err)
	}

	err = r.ValidateAndUpdateWithValidation(mapOut, mapToUpdate, validations)
	if err != nil {
		return fmt.Errorf("error updating map: %v", err)
	}

	return nil
}
