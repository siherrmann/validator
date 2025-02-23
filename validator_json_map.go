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
		if len(values[k]) > 1 {
			mapOut[k] = values[k]
		} else {
			mapOut[k] = values.Get(k)
		}
	}
	return mapOut, nil
}
