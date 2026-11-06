package helper

import (
	"reflect"
	"strings"
)

// FieldKey resolves the map key and tag options of a struct field.
// The key is taken from the json tag, falls back to the yaml tag and then to the field name.
// A tag name of "-" or "" keeps the field name as key.
// The options omitempty and inline are read from the same tag the key is resolved from.
func FieldKey(fieldType reflect.StructField) (key string, omitEmpty bool, inline bool) {
	key = fieldType.Name

	tag := fieldType.Tag.Get("json")
	if len(tag) == 0 {
		tag = fieldType.Tag.Get("yaml")
	}
	if len(tag) == 0 {
		return key, false, false
	}

	parts := strings.Split(tag, ",")
	if parts[0] != "-" && parts[0] != "" {
		key = parts[0]
	}
	for _, part := range parts[1:] {
		switch part {
		case "omitempty":
			omitEmpty = true
		case "inline":
			inline = true
		}
	}
	return key, omitEmpty, inline
}
