package validator

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/siherrmann/validator/helper"
	"github.com/siherrmann/validator/model"
)

// All unpack functions validate with the context of the request (request.Context()),
// so values injected with Inject into the request context can be used in requirements and functions.
//
// All unpack functions read request.Body only once and only as far as the reader allows,
// so the body can be limited with http.MaxBytesReader.

var (
	// ErrUnsupportedMediaType is returned by UnpackAndValidate (wrapped) if the Content-Type of the request
	// is not supported. A handler can answer it with 415.
	ErrUnsupportedMediaType = errors.New("unsupported media type")
	// ErrMalformedInput is returned by all unpack functions (wrapped) if the input can not be read or parsed,
	// eg. invalid json or yaml, an invalid query or a body that is larger than allowed by http.MaxBytesReader
	// (the *http.MaxBytesError is wrapped too). A handler can answer it with 400 (or 413).
	// All other errors are errors of the validation or of the struct.
	ErrMalformedInput = errors.New("malformed input")
)

// mediaTypeSource is the source of the input of a request.
type mediaTypeSource int

const (
	sourceForm mediaTypeSource = iota
	sourceJson
	sourceYaml
)

// mediaTypeSources maps the supported media types to their source.
// Media types with the suffix `+json` (eg. `application/problem+json`) are json too.
var mediaTypeSources = map[string]mediaTypeSource{
	"application/x-www-form-urlencoded": sourceForm,
	"multipart/form-data":               sourceForm,
	"application/json":                  sourceJson,
	"application/yaml":                  sourceYaml,
	"application/x-yaml":                sourceYaml,
	"text/yaml":                         sourceYaml,
}

// UnpackAndValidate unpacks the input of the request by its Content-Type, validates it and updates the given struct:
//
//   - `application/x-www-form-urlencoded` and `multipart/form-data`: the form values and the query (see UnmapAndValidate),
//   - `application/json` (also with parameters like charset, or a media type with the suffix `+json`): the json body (see UnmarshalAndValidate),
//   - `application/yaml`, `application/x-yaml` and `text/yaml`: the yaml body,
//   - no Content-Type: the body as json or, if the body is empty, the query.
//
// Any other Content-Type returns ErrUnsupportedMediaType,
// input that can not be read or parsed returns ErrMalformedInput. Both are wrapped, check them with errors.Is.
// It returns an error if the unpacking, validation or update fails.
//
// For more information look at ValidateAndUpdate.
func (r *Validator) UnpackAndValidate(request *http.Request, structToUpdate any, tagType ...string) error {
	mapOut, err := unpackRequest(request)
	if err != nil {
		return err
	}

	err = r.ValidateAndUpdateContext(request.Context(), mapOut, structToUpdate, tagType...)
	if err != nil {
		return fmt.Errorf("error updating struct: %w", err)
	}

	return nil
}

// UnpackAndValidateWithValidation unpacks the input of the request by its Content-Type (see UnpackAndValidate).
// It validates the map by the given validations and updates it.
// It returns an error if the unpacking, validation or update fails.
func (r *Validator) UnpackAndValidateWithValidation(request *http.Request, mapToUpdate *map[string]any, validations []model.Validation) error {
	mapOut, err := unpackRequest(request)
	if err != nil {
		return err
	}

	err = r.ValidateAndUpdateWithValidationContext(request.Context(), mapOut, mapToUpdate, validations)
	if err != nil {
		return fmt.Errorf("error updating map: %w", err)
	}

	return nil
}

// unpackRequest reads the input of the request by its Content-Type, see UnpackAndValidate.
func unpackRequest(request *http.Request) (map[string]any, error) {
	if request == nil {
		return nil, fmt.Errorf("request is nil")
	}

	mediaType, err := helper.MediaType(request)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnsupportedMediaType, err)
	}

	if mediaType == "" {
		return unpackWithoutMediaType(request)
	}

	source, ok := mediaTypeSources[mediaType]
	if !ok && strings.HasSuffix(mediaType, "+json") {
		source, ok = sourceJson, true
	}
	if !ok {
		return nil, fmt.Errorf("%w: %v", ErrUnsupportedMediaType, mediaType)
	}

	switch source {
	case sourceJson:
		return unmarshalRequest(request)
	case sourceYaml:
		mapOut, err := helper.UnmarshalYamlRequestToJsonMap(request)
		if err != nil {
			return nil, fmt.Errorf("%w: error unmarshaling request body: %w", ErrMalformedInput, err)
		}
		return mapOut, nil
	default:
		return unmapRequest(request)
	}
}

// unpackWithoutMediaType reads the body of a request without Content-Type as json.
// If there is no body or it is empty, the query is read instead.
func unpackWithoutMediaType(request *http.Request) (map[string]any, error) {
	if request.Body == nil || request.Body == http.NoBody {
		return unmapRequest(request)
	}
	defer request.Body.Close()

	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: error reading request body: %w", ErrMalformedInput, err)
	} else if len(bodyBytes) == 0 {
		return unmapRequest(request)
	}

	mapOut, err := helper.UnmarshalJsonToJsonMap(bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("%w: error unmarshaling request body: %w", ErrMalformedInput, err)
	}
	return mapOut, nil
}

// unmarshalRequest reads the json body of the request.
func unmarshalRequest(request *http.Request) (map[string]any, error) {
	mapOut, err := helper.UnmarshalRequestToJsonMap(request)
	if err != nil {
		return nil, fmt.Errorf("%w: error unmarshaling request body: %w", ErrMalformedInput, err)
	}
	return mapOut, nil
}

// unmapRequest reads the form values of the request.
func unmapRequest(request *http.Request) (map[string]any, error) {
	mapOut, err := helper.UnmapRequestToJsonMap(request)
	if err != nil {
		return nil, fmt.Errorf("%w: error unmapping form values: %w", ErrMalformedInput, err)
	}
	return mapOut, nil
}

// UnmarshalAndValidate unmarshals the json body of the request, validates it and updates the given struct.
// It returns an error if the unmarshaling, validation or update fails.
//
// For more information look at ValidateAndUpdate.
func (r *Validator) UnmarshalAndValidate(request *http.Request, structToUpdate any, tagType ...string) error {
	mapOut, err := unmarshalRequest(request)
	if err != nil {
		return err
	}

	err = r.ValidateAndUpdateContext(request.Context(), mapOut, structToUpdate, tagType...)
	if err != nil {
		return fmt.Errorf("error updating struct: %w", err)
	}

	return nil
}

// UnmapAndValidate unmaps the form values of the request (the query and, for a url encoded or multipart form, the body),
// validates them and updates the given struct. Files of a multipart form are ignored.
// It returns an error if the unmapping, validation or update fails.
//
// For more information look at ValidateAndUpdate.
func (r *Validator) UnmapAndValidate(request *http.Request, structToUpdate any, tagType ...string) error {
	mapOut, err := unmapRequest(request)
	if err != nil {
		return err
	}

	err = r.ValidateAndUpdateContext(request.Context(), mapOut, structToUpdate, tagType...)
	if err != nil {
		return fmt.Errorf("error updating struct: %w", err)
	}

	return nil
}

// UnmarshalAndValidateWithValidation unmarshals the json body of the request.
// It validates the map by the given validations and updates it.
// It returns an error if the unmarshaling, validation or update fails.
func (r *Validator) UnmarshalAndValidateWithValidation(request *http.Request, mapToUpdate *map[string]any, validations []model.Validation) error {
	mapOut, err := unmarshalRequest(request)
	if err != nil {
		return err
	}

	err = r.ValidateAndUpdateWithValidationContext(request.Context(), mapOut, mapToUpdate, validations)
	if err != nil {
		return fmt.Errorf("error updating struct: %w", err)
	}

	return nil
}

// UnmapAndValidateWithValidation unmaps the form values of the request (see UnmapAndValidate).
// It validates the map by the given validations and updates it.
// It returns an error if the unmapping, validation or update fails.
func (r *Validator) UnmapAndValidateWithValidation(request *http.Request, mapToUpdate *map[string]any, validations []model.Validation) error {
	mapOut, err := unmapRequest(request)
	if err != nil {
		return err
	}

	err = r.ValidateAndUpdateWithValidationContext(request.Context(), mapOut, mapToUpdate, validations)
	if err != nil {
		return fmt.Errorf("error updating struct: %w", err)
	}

	return nil
}
