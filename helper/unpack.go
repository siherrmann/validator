package helper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func UnmarshalRequestToJsonMap(request *http.Request) (map[string]any, error) {
	if request == nil {
		return nil, fmt.Errorf("request is nil")
	}

	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading request body: %v", err)
	}
	defer request.Body.Close()

	return UnmarshalJsonToJsonMap(bodyBytes)
}

func UnmapRequestToJsonMap(request *http.Request) (map[string]any, error) {
	if request == nil {
		return nil, fmt.Errorf("request is nil")
	}

	err := request.ParseForm()
	if err != nil {
		return nil, fmt.Errorf("error parsing form: %v", err)
	}

	return UnmapUrlValuesToJsonMap(request.Form)
}

func UnmarshalJsonToJsonMap(jsonInput []byte) (map[string]any, error) {
	mapOut := map[string]any{}
	err := json.Unmarshal(jsonInput, &mapOut)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling: %v", err)
	}
	return mapOut, nil
}

func UnmapUrlValuesToJsonMap(values url.Values) (map[string]any, error) {
	mapOut := map[string]any{}
	for k := range values {
		if len(values[k]) > 1 {
			arrayOut := []any{}
			for _, v := range values[k] {
				var unmarshalled any
				err := json.Unmarshal([]byte(v), &unmarshalled)
				if err == nil {
					arrayOut = append(arrayOut, unmarshalled)
				} else {
					arrayOut = append(arrayOut, v)
				}
			}
			mapOut[k] = arrayOut
		} else {
			value := values.Get(k)
			var unmarshalled any
			err := json.Unmarshal([]byte(value), &unmarshalled)
			if err == nil {
				mapOut[k] = unmarshalled
			} else {
				mapOut[k] = value
			}
		}
	}
	return mapOut, nil
}
