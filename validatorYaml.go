package validator

import (
	"context"
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
	return r.UnmarshalYamlValidateAndUpdateContext(context.Background(), yamlInput, structToUpdate, tagType...)
}

// UnmarshalYamlValidateAndUpdateContext is UnmarshalYamlValidateAndUpdate with the context of the validation call,
// which holds the values injected with [Inject].
func (r *Validator) UnmarshalYamlValidateAndUpdateContext(ctx context.Context, yamlInput []byte, structToUpdate any, tagType ...string) error {
	mapOut, err := helper.UnmarshalYamlToJsonMap(yamlInput)
	if err != nil {
		return fmt.Errorf("error unmarshaling yaml: %v", err)
	}

	err = r.ValidateAndUpdateContext(ctx, mapOut, structToUpdate, tagType...)
	if err != nil {
		return fmt.Errorf("error updating struct: %v", err)
	}

	return nil
}

// UnmarshalYamlValidateAndUpdateWithValidation unmarshals given yaml ([]byte) into a JsonMap.
// It validates the map by the given validations and updates the given map.
// It returns an error if the unmarshaling, validation or update fails.
func (r *Validator) UnmarshalYamlValidateAndUpdateWithValidation(yamlInput []byte, mapToUpdate *map[string]any, validations []model.Validation) error {
	return r.UnmarshalYamlValidateAndUpdateWithValidationContext(context.Background(), yamlInput, mapToUpdate, validations)
}

// UnmarshalYamlValidateAndUpdateWithValidationContext is UnmarshalYamlValidateAndUpdateWithValidation
// with the context of the validation call, which holds the values injected with [Inject].
func (r *Validator) UnmarshalYamlValidateAndUpdateWithValidationContext(ctx context.Context, yamlInput []byte, mapToUpdate *map[string]any, validations []model.Validation) error {
	mapOut, err := helper.UnmarshalYamlToJsonMap(yamlInput)
	if err != nil {
		return fmt.Errorf("error unmarshaling yaml: %v", err)
	}

	err = r.ValidateAndUpdateWithValidationContext(ctx, mapOut, mapToUpdate, validations)
	if err != nil {
		return fmt.Errorf("error updating map: %v", err)
	}

	return nil
}
