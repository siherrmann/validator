package model

import (
	"fmt"
	"reflect"
	"strings"
)

type Validation struct {
	Key         string
	Type        ValidatorType
	Requirement string
	Groups      []*Group
}

func (r *Validation) Fill(tag string, tagType TagType, value reflect.Value) error {
	r.Type = TypeFromInterface(value.Interface())

	tagSplit := strings.Split(tag, ", ")
	r.Requirement = "-"

	tagIndex := 0
	if len(tagSplit) > tagIndex && tagType == UPD {
		r.Key = strings.TrimSpace(tagSplit[tagIndex])
		tagIndex++
	}

	if len(tagSplit) > tagIndex && tagType == UPD || len(tagSplit) > tagIndex {
		r.Requirement = tagSplit[tagIndex]
		tagIndex++
	}

	var err error
	if len(tagSplit) > tagIndex {
		r.Groups, err = GetGroups(tagSplit[tagIndex])
		if err != nil {
			return fmt.Errorf("error extracting group: %v", err)
		}
	}

	return nil
}

func (r *Validation) GetValidValue(in interface{}) (interface{}, error) {
	var out interface{}
	var err error
	switch in := in.(type) {
	case string:
		out, err = InterfaceFromString(in, r.Type)
	case int, int64, int32, int16, int8:
		out, err = InterfaceFromInt(in.(int), r.Type)
	case float64, float32:
		out, err = InterfaceFromFloat(in.(float64), r.Type)
	default:
		return in, nil
	}
	return out, err
}

type TagType string

const (
	VLD TagType = "vld"
	UPD TagType = "upd"
)

type ValidatorMap map[string]Validation

type ValidatorType string

const (
	String      ValidatorType = "string"
	Int         ValidatorType = "int"
	Float       ValidatorType = "float"
	Bool        ValidatorType = "bool"
	Array       ValidatorType = "array"
	Map         ValidatorType = "map"
	Struct      ValidatorType = "struct"
	Time        ValidatorType = "time"
	TimeISO8601 ValidatorType = "timeIso8601"
	TimeUnix    ValidatorType = "timeUnix"
)
