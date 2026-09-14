package validator

import (
	"reflect"
	"sync"

	"github.com/siherrmann/validator/model"
)

// validationCacheKey identifies a cached extraction by both the struct type and the
// tag type it was extracted under, since the same struct type can be validated under
// different tag types (e.g. "vld" for input validation, a custom tag for updates) with
// different requirements per field.
type validationCacheKey struct {
	structType reflect.Type
	tagType    string
}

// validationCache holds the result of GetValidationsFromStruct per (struct type, tag
// type). Extraction walks the struct by reflection and parses every tagged field's
// requirement into an AST — expensive work that only ever depends on the type, never
// on the values a particular call is validating, since Go struct tags are fixed at
// compile time. A sync.Map fits this: the key set is bounded by the number of struct
// types the program validates, reads vastly outnumber writes, and nothing needs
// eviction.
//
// Cached slices are shared with callers as-is, not copied. That is only safe because
// nothing downstream mutates a Validation, its Groups, or its RequirementAST after
// extraction — ValidateWithValidation and validators.ValidateGroups only ever read
// them. If that invariant changes, callers of getCachedValidationsFromStruct must be
// given a copy instead.
var validationCache sync.Map

// getCachedValidationsFromStruct returns the cached extraction for v's type and
// tagType, populating the cache on a miss. An extraction that returns an error is
// never cached, so a transient failure doesn't poison every later call for that type.
func getCachedValidationsFromStruct(v any, tagType string) ([]model.Validation, error) {
	key := validationCacheKey{structType: reflect.TypeOf(v), tagType: tagType}
	if cached, ok := validationCache.Load(key); ok {
		return cached.([]model.Validation), nil
	}

	validations, err := GetValidationsFromStruct(v, tagType)
	if err != nil {
		return nil, err
	}

	validationCache.Store(key, validations)
	return validations, nil
}
