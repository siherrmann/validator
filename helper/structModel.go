package helper

import (
	"encoding/json/v2"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"
)

// StructModel describes how a struct type looks as a JsonMap.
// It is the single source for the keys of a struct, used for reading (UnmapStructToJsonMap),
// writing (MapJsonMapToStruct) and the extraction of validations, so all of them behave the same.
type StructModel struct {
	Type   reflect.Type
	Fields []StructField
}

// StructField is a field of a StructModel.
type StructField struct {
	// Key is the key of the field in the JsonMap (see FieldKey).
	Key string
	// Index is the index path of the field (see reflect.Value.FieldByIndex),
	// it has more than one element for fields promoted from embedded or inline structs.
	Index     []int
	OmitEmpty bool
	// Field is the struct field itself.
	Field reflect.StructField
	// decodeType is set for fields that are written by encoding/json (see needsJSONDecode).
	// It is a struct with the field as only field `V`, so the json options of the field are kept.
	decodeType reflect.Type
}

var structModelCache sync.Map

// StructModelOf returns the cached StructModel of the given struct type.
//
// The fields follow the rules of encoding/json:
//   - Unexported fields and fields tagged with `json:"-"` (or `yaml:"-"` without json tag) are left out.
//   - Embedded structs (also behind pointers) without a name in the tag and fields tagged with `inline` or `embed`
//     are flattened into the parent. Exported fields of unexported embedded structs are promoted too.
//   - If several fields have the same key, the shallowest one wins. At the same depth the field with a name
//     in the tag wins, otherwise all of them are left out.
//
// Recursive types are fine, a struct is only flattened once.
func StructModelOf(structType reflect.Type) *StructModel {
	if cached, ok := structModelCache.Load(structType); ok {
		return cached.(*StructModel)
	}
	structModel := &StructModel{Type: structType, Fields: structFields(structType)}
	cached, _ := structModelCache.LoadOrStore(structType, structModel)
	return cached.(*StructModel)
}

type structLevel struct {
	structType reflect.Type
	index      []int
}

type fieldCandidate struct {
	field StructField
	named bool
}

// structFields collects the fields of the struct type level by level (breadth first) like encoding/json,
// so the depth of a field is known for resolving conflicting keys.
func structFields(structType reflect.Type) []StructField {
	candidates := []fieldCandidate{}
	visited := map[reflect.Type]bool{}
	next := []structLevel{{structType: structType}}
	for len(next) > 0 {
		current := next
		next = nil
		for _, level := range current {
			if visited[level.structType] {
				continue
			}
			visited[level.structType] = true

			for i := 0; i < level.structType.NumField(); i++ {
				field := level.structType.Field(i)
				if !isModelField(field) {
					continue
				}

				index := make([]int, len(level.index)+1)
				copy(index, level.index)
				index[len(level.index)] = i

				key, omitEmpty, inline := FieldKey(field)
				named := hasTagName(field)
				fieldType := field.Type
				if fieldType.Kind() == reflect.Ptr {
					fieldType = fieldType.Elem()
				}
				if fieldType.Kind() == reflect.Struct && (inline || (field.Anonymous && !named)) {
					next = append(next, structLevel{structType: fieldType, index: index})
					continue
				}

				candidates = append(candidates, fieldCandidate{
					field: StructField{
						Key:        key,
						Index:      index,
						OmitEmpty:  omitEmpty,
						Field:      field,
						decodeType: decodeType(field),
					},
					named: named,
				})
			}
		}
	}

	return dominantFields(candidates)
}

// dominantFields drops the candidates hidden by others with the same key and sorts the rest in field order.
func dominantFields(candidates []fieldCandidate) []StructField {
	byKey := map[string][]fieldCandidate{}
	for _, candidate := range candidates {
		byKey[candidate.field.Key] = append(byKey[candidate.field.Key], candidate)
	}

	fields := []StructField{}
	for _, sameKey := range byKey {
		if dominant, ok := dominantField(sameKey); ok {
			fields = append(fields, dominant)
		}
	}
	sort.Slice(fields, func(i, j int) bool {
		return lessIndex(fields[i].Index, fields[j].Index)
	})
	return fields
}

// dominantField returns the field that wins between candidates with the same key.
// ok is false if there is no single winner.
func dominantField(candidates []fieldCandidate) (StructField, bool) {
	depth := len(candidates[0].field.Index)
	for _, candidate := range candidates {
		depth = min(depth, len(candidate.field.Index))
	}

	var dominant *fieldCandidate
	count, namedCount := 0, 0
	for i := range candidates {
		if len(candidates[i].field.Index) != depth {
			continue
		}
		count++
		if candidates[i].named {
			namedCount++
			dominant = &candidates[i]
		} else if dominant == nil || !dominant.named {
			dominant = &candidates[i]
		}
	}
	if count == 1 || namedCount == 1 {
		return dominant.field, true
	}
	return StructField{}, false
}

func lessIndex(a, b []int) bool {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return len(a) < len(b)
}

// isModelField reports if the field is part of the JsonMap of its struct (or promotes fields to it).
func isModelField(field reflect.StructField) bool {
	if field.Anonymous {
		fieldType := field.Type
		if !field.IsExported() {
			// Fields of an unexported embedded struct are promoted, but a nil pointer to it can not be set.
			if fieldType.Kind() == reflect.Ptr || fieldType.Kind() != reflect.Struct {
				return false
			}
		}
	} else if !field.IsExported() {
		return false
	}

	if tag, ok := field.Tag.Lookup("json"); ok {
		return tag != "-"
	}
	return field.Tag.Get("yaml") != "-"
}

// hasTagName reports if the tag the key is resolved from (see FieldKey) has a name.
func hasTagName(field reflect.StructField) bool {
	tag := field.Tag.Get("json")
	if len(tag) == 0 {
		tag = field.Tag.Get("yaml")
	}
	name, _, _ := strings.Cut(tag, ",")
	return name != "" && name != "-"
}

var (
	timeType              = reflect.TypeOf(time.Time{})
	jsonUnmarshalerType   = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
	jsonUnmarshalFromType = reflect.TypeOf((*json.UnmarshalerFrom)(nil)).Elem()
)

// decodeType returns the type used to write the field with encoding/json, or nil if it is written by AnyToType.
// See needsJSONDecode.
func decodeType(field reflect.StructField) reflect.Type {
	if !needsJSONDecode(field) {
		return nil
	}
	_, options, _ := strings.Cut(field.Tag.Get("json"), ",")
	tag := `json:"V`
	if options != "" {
		tag += "," + options
	}
	tag += `"`
	return reflect.StructOf([]reflect.StructField{{Name: "V", Type: field.Type, Tag: reflect.StructTag(tag)}})
}

// needsJSONDecode reports if the field has to be written by encoding/json, because its type has an own
// UnmarshalJSON (or UnmarshalJSONFrom) or its json tag has the options `string` or `format:`.
// time.Time is written by AnyToType, which also accepts unix timestamps.
func needsJSONDecode(field reflect.StructField) bool {
	fieldType := field.Type
	if fieldType == timeType || (fieldType.Kind() == reflect.Ptr && fieldType.Elem() == timeType) {
		return false
	}
	for _, t := range []reflect.Type{fieldType, reflect.PointerTo(fieldType)} {
		if t.Implements(jsonUnmarshalerType) || t.Implements(jsonUnmarshalFromType) {
			return true
		}
	}

	_, options, _ := strings.Cut(field.Tag.Get("json"), ",")
	for options != "" {
		var option string
		option, options, _ = strings.Cut(options, ",")
		if option == "string" || strings.HasPrefix(option, "format:") {
			return true
		}
	}
	return false
}

// fieldForRead returns the field of the index path in the struct value.
// ok is false if an embedded pointer on the path is nil, so the field does not exist.
func fieldForRead(structValue reflect.Value, index []int) (reflect.Value, bool) {
	value := structValue
	for i, fieldIndex := range index {
		if i > 0 && value.Kind() == reflect.Ptr {
			if value.IsNil() {
				return reflect.Value{}, false
			}
			value = value.Elem()
		}
		value = value.Field(fieldIndex)
	}
	return value, true
}

// fieldForWrite returns the field of the index path in the struct value, nil embedded pointers on the path are set.
func fieldForWrite(structValue reflect.Value, index []int) reflect.Value {
	value := structValue
	for i, fieldIndex := range index {
		if i > 0 && value.Kind() == reflect.Ptr {
			if value.IsNil() {
				value.Set(reflect.New(value.Type().Elem()))
			}
			value = value.Elem()
		}
		value = value.Field(fieldIndex)
	}
	return value
}
