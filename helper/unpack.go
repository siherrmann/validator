package helper

import (
	"encoding/json/v2"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"go.yaml.in/yaml/v3"
)

// MultipartMaxMemory is the number of bytes of a multipart form that are kept in memory,
// the rest of the files is stored in temporary files (see http.Request.ParseMultipartForm).
// It is the same as the default of net/http.
const MultipartMaxMemory = 32 << 20

// UnmarshalRequestToJsonMap unmarshals the json body of the request into a JsonMap.
// An error while reading the body (eg. *http.MaxBytesError) is wrapped.
func UnmarshalRequestToJsonMap(request *http.Request) (map[string]any, error) {
	bodyBytes, err := readRequestBody(request)
	if err != nil {
		return nil, err
	}
	return UnmarshalJsonToJsonMap(bodyBytes)
}

// UnmarshalYamlRequestToJsonMap unmarshals the yaml body of the request into a JsonMap, see UnmarshalYamlToJsonMap.
// An error while reading the body (eg. *http.MaxBytesError) is wrapped.
func UnmarshalYamlRequestToJsonMap(request *http.Request) (map[string]any, error) {
	bodyBytes, err := readRequestBody(request)
	if err != nil {
		return nil, err
	}
	return UnmarshalYamlToJsonMap(bodyBytes)
}

func readRequestBody(request *http.Request) ([]byte, error) {
	if request == nil {
		return nil, fmt.Errorf("request is nil")
	}
	if request.Body == nil {
		return nil, nil
	}
	defer request.Body.Close()

	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading request body: %w", err)
	}
	return bodyBytes, nil
}

// UnmapRequestToJsonMap unmaps the form values of the request (the query and,
// for a url encoded or multipart form, the body) into a JsonMap. Files of a multipart form are ignored.
// An error while reading the body (eg. *http.MaxBytesError) is wrapped.
func UnmapRequestToJsonMap(request *http.Request) (map[string]any, error) {
	if request == nil {
		return nil, fmt.Errorf("request is nil")
	}

	var err error
	if IsMediaType(request, "multipart/form-data") {
		err = request.ParseMultipartForm(MultipartMaxMemory)
	} else {
		err = request.ParseForm()
	}
	if err != nil {
		return nil, fmt.Errorf("error parsing form: %w", err)
	}

	return UnmapUrlValuesToJsonMap(request.Form)
}

// MediaType returns the media type of the Content-Type header of the request without parameters
// (eg. `application/json` for `application/json; charset=utf-8`). It is empty if there is no Content-Type.
func MediaType(request *http.Request) (string, error) {
	contentType := request.Header.Get("Content-Type")
	if contentType == "" {
		return "", nil
	}
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", fmt.Errorf("invalid Content-Type %q: %v", contentType, err)
	}
	return mediaType, nil
}

// IsMediaType reports if the Content-Type of the request is one of the given media types.
func IsMediaType(request *http.Request, mediaTypes ...string) bool {
	mediaType, err := MediaType(request)
	return err == nil && slices.Contains(mediaTypes, mediaType)
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
