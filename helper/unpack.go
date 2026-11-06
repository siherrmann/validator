package helper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"go.yaml.in/yaml/v3"
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

// UnmarshalYamlToJsonMap unmarshals yaml into a JsonMap.
// Nested maps with non-string keys are converted to map[string]any, so the result
// has the same shape as an unmarshaled json object. An empty document results in an empty map.
func UnmarshalYamlToJsonMap(yamlInput []byte) (map[string]any, error) {
	var out any
	err := yaml.Unmarshal(yamlInput, &out)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling yaml: %v", err)
	}
	if out == nil {
		return map[string]any{}, nil
	}

	mapOut, ok := normalizeYamlValue(out).(map[string]any)
	if !ok {
		return nil, fmt.Errorf("error unmarshaling yaml: top level has to be a mapping, was %T", out)
	}
	return mapOut, nil
}

// normalizeYamlValue recursively converts map[any]any (yaml mappings with non-string keys) to map[string]any.
func normalizeYamlValue(in any) any {
	switch v := in.(type) {
	case map[string]any:
		for key, value := range v {
			v[key] = normalizeYamlValue(value)
		}
		return v
	case map[any]any:
		mapOut := make(map[string]any, len(v))
		for key, value := range v {
			mapOut[fmt.Sprint(key)] = normalizeYamlValue(value)
		}
		return mapOut
	case []any:
		for i, value := range v {
			v[i] = normalizeYamlValue(value)
		}
		return v
	default:
		return in
	}
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
			if strings.Contains(k, ".") {
				setNestedDotValue(mapOut, strings.Split(k, "."), arrayOut)
			} else {
				mapOut[k] = arrayOut
			}
		} else {
			value := values.Get(k)
			var unmarshalled any
			var valToSet any
			err := json.Unmarshal([]byte(value), &unmarshalled)
			if err == nil {
				valToSet = unmarshalled
			} else {
				valToSet = value
			}
			
			if strings.Contains(k, ".") {
				setNestedDotValue(mapOut, strings.Split(k, "."), valToSet)
			} else {
				mapOut[k] = valToSet
			}
		}
	}
	return mapOut, nil
}

func setNestedDotValue(m map[string]any, keys []string, value any) {
	if len(keys) == 0 {
		return
	}
	if len(keys) == 1 {
		m[keys[0]] = value
		return
	}
	key := keys[0]
	if existing, ok := m[key]; ok {
		if nestedMap, ok := existing.(map[string]any); ok {
			setNestedDotValue(nestedMap, keys[1:], value)
			return
		}
	}
	nestedMap := make(map[string]any)
	m[key] = nestedMap
	setNestedDotValue(nestedMap, keys[1:], value)
}
