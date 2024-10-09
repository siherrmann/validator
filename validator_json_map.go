package validator

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"

	"github.com/siherrmann/validator/model"
)

func UnmarshalJsonToJsonMapAndValidate(jsonInput []byte, validation ...model.Validation) (model.JsonMap, error) {
	mapOut, err := UnmarshalJsonToJsonMap(jsonInput)
	if err != nil {
		return nil, err
	}
	err = ValidateJsonMap(mapOut, validation...)
	if err != nil {
		return nil, err
	}
	return mapOut, nil
}

func UnmapUrlValuesAndValidateWithValidation(values url.Values, validation ...model.Validation) (model.JsonMap, error) {
	mapOut, err := UnmapUrlValuesToJsonMap(values, validation...)
	if err != nil {
		return nil, err
	}
	err = ValidateJsonMap(mapOut, validation...)
	if err != nil {
		return nil, err
	}
	return mapOut, nil
}

func UnmarshalJsonToJsonMap(jsonInput []byte) (model.JsonMap, error) {
	mapOut := model.JsonMap{}
	err := json.Unmarshal(jsonInput, &mapOut)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling: %v", err)
	}
	return mapOut, nil
}

func UnmapUrlValuesToJsonMap(values url.Values, validation ...model.Validation) (model.JsonMap, error) {
	mapOut := model.JsonMap{}
	for _, v := range validation {
		if values.Has(v.Key) {
			value := values.Get(v.Key)
			out, err := model.InterfaceFromString(value, v.Type)
			if err != nil {
				return nil, err
			}
			mapOut[v.Key] = out
		}
	}
	return mapOut, nil
}

func ValidateJsonMap(jsonMap model.JsonMap, validations ...model.Validation) error {
	for _, v := range validations {
		if jsonMap.Has(v.Key) {
			if len(v.Type) == 0 {
				v.Type = model.TypeFromInterface(v)
			}

			_, err := ValidateValueWithParser(reflect.ValueOf(v), &v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
