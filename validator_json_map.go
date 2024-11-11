package validator

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/siherrmann/validator/model"
)

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
