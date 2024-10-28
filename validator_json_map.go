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
	mapOut, err := UnmapUrlValuesToJsonMap(values)
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

func UnmapUrlValuesToJsonMap(values url.Values) (model.JsonMap, error) {
	mapOut := model.JsonMap{}
	for k := range values {
		mapOut[k] = values.Get(k)
	}
	return mapOut, nil
}

func ValidateJsonMap(jsonMap model.JsonMap, validations ...model.Validation) error {
	for _, v := range validations {
		if jsonMap.Has(v.Key) {
			if len(v.Type) == 0 {
				v.Type = model.TypeFromInterface(v)
			}

			_, err := ValidateValueWithParser(reflect.ValueOf(jsonMap[v.Key]), &v)
			if err != nil {
				return fmt.Errorf("invalid value %s: %v", v.Key, err)
			}
		}
	}
	return nil
}
