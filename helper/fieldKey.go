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

	// strings.Cut instead of strings.Split, because this runs for every field on every
	// Validate call and Split allocates a slice for every tag.
	name, options, _ := strings.Cut(tag, ",")
	if name != "-" && name != "" {
		key = name
	}
	for options != "" {
		var option string
		option, options, _ = strings.Cut(options, ",")
		switch option {
		case "omitempty":
			omitEmpty = true
		case "inline":
			inline = true
		}
	}
	return key, omitEmpty, inline
}
