package model

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/siherrmann/validator/helper"
)

// Default tag type.
const VLD string = "vld"

type Validation struct {
	Key         string
	Type        ValidatorType
	Requirement string
	Groups      []*Group
	Default     string
	// Inner Struct validation
	InnerValidation []Validation
}

func GetValidationsFromStruct(in interface{}, tagType string) ([]Validation, error) {
	err := helper.CheckValidPointerToStruct(in)
	if err != nil {
		return nil, err
	}

	validations := []Validation{}

	structFull := reflect.ValueOf(in).Elem()
	for i := 0; i < structFull.Type().NumField(); i++ {
		field := structFull.Field(i)
		fieldType := structFull.Type().Field(i)

		validation, err := GetValidationFromStructField(tagType, field, fieldType)
		if err != nil {
			return nil, err
		}

		if len(validation.Requirement) > 0 {
			validations = append(validations, validation)
		}
	}
	return validations, nil
}

func GetValidationFromStructField(tagType string, fieldValue reflect.Value, fieldType reflect.StructField) (Validation, error) {
	validation := Validation{}
	validation.Key = fieldType.Name
	if len(fieldType.Tag.Get("json")) > 0 {
		validation.Key = fieldType.Tag.Get("json")
	}
	validation.Type = TypeFromInterface(fieldValue.Interface())
	validation.Requirement = "-"

	tagIndex := 0
	tagSplit := strings.Split(fieldType.Tag.Get(string(tagType)), ", ")

	if len(tagSplit) > tagIndex {
		validation.Requirement = tagSplit[tagIndex]
		tagIndex++
	}

	if len(tagSplit) > tagIndex {
		var err error
		validation.Groups, err = GetGroups(tagSplit[tagIndex])
		if err != nil {
			return Validation{}, fmt.Errorf("error extracting group: %v", err)
		}
	}

	if helper.IsArrayOfStruct(fieldValue.Interface()) {
		for i := 0; i < fieldValue.Len(); i++ {
			innerStruct := fieldValue.Index(i).Addr().Interface()
			innerValidation, err := GetValidationsFromStruct(innerStruct, string(tagType))
			if err != nil {
				return Validation{}, fmt.Errorf("error getting inner validation: %v", err)
			}
			validation.InnerValidation = append(validation.InnerValidation, innerValidation...)
		}
	}

	return validation, nil
}

func (r *Validation) GetValidValue(in interface{}) (interface{}, error) {
	var out interface{}
	var err error
	switch in := in.(type) {
	case string:
		out, err = r.InterfaceFromString(in)
	case int:
		out, err = r.InterfaceFromInt(in)
	case int64:
		out, err = r.InterfaceFromInt(int(in))
	case int32:
		out, err = r.InterfaceFromInt(int(in))
	case int16:
		out, err = r.InterfaceFromInt(int(in))
	case int8:
		out, err = r.InterfaceFromInt(int(in))
	case float64:
		out, err = r.InterfaceFromFloat(in)
	case float32:
		out, err = r.InterfaceFromFloat(float64(in))
	case []string:
		out, err = r.InterfaceFromArrayOfString(in)
	default:
		if v, ok := in.(JsonMap); ok {
			return map[string]interface{}(v), nil
		}
		return in, nil
	}
	return out, err
}

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

func TypeFromInterface(in interface{}) ValidatorType {
	switch in.(type) {
	case string:
		return String
	case int, int64, int32, int16, int8:
		return Int
	case float64, float32:
		return Float
	case bool:
		return Bool
	case JsonMap, map[string]string, map[string]int, map[string]int64, map[string]int32, map[string]int16, map[string]int8, map[string]float64, map[string]float32, map[string]bool:
		return Map
	case []string, []int, []int64, []int32, []int16, []int8, []float64, []float32, []bool:
		return Array
	case time.Time:
		return Time
	default:
		// custom types
		if reflect.TypeOf(in).Kind() == reflect.String {
			return String
		} else if reflect.TypeOf(in).Kind() == reflect.Int {
			return Int
		} else if reflect.TypeOf(in).Kind() == reflect.Float64 || reflect.TypeOf(in).Kind() == reflect.Float32 {
			return Float
		} else if reflect.TypeOf(in).Kind() == reflect.Bool {
			return Bool
			// other types
		} else if helper.IsArray(in) {
			return Array
		} else {
			return Struct
		}
	}
}

func (r *Validation) InterfaceFromString(in string) (interface{}, error) {
	switch r.Type {
	case String:
		return in, nil
	case Int:
		out, err := strconv.Atoi(in)
		if err != nil {
			return nil, err
		}
		return out, nil
	case Float:
		out, err := strconv.ParseFloat(strings.TrimSpace(in), 64)
		if err != nil {
			return nil, err
		}
		return out, nil
	case Bool:
		// case for html forms
		if in == "on" {
			return true, nil
		}
		out, err := strconv.ParseBool(in)
		if err != nil {
			return nil, err
		}
		return out, nil
	case Array:
		out := []interface{}{}
		if strings.Contains(in, "[") {
			err := json.Unmarshal([]byte(in), &out)
			if err != nil {
				return nil, err
			}
		} else if len(in) > 0 {
			out = append(out, in)
		}
		return out, nil
	case Map, Struct:
		out := map[string]interface{}{}
		err := json.Unmarshal([]byte(in), &out)
		if err != nil {
			return nil, err
		}
		return out, nil
	case Time:
		date, err := UnixStringToTime(in)
		if err != nil {
			date, err := ISO8601StringToTime(in)
			if err != nil {
				return nil, err
			}
			return date, nil
		}
		return date, nil
	case TimeUnix:
		return UnixStringToTime(in)
	case TimeISO8601:
		return ISO8601StringToTime(in)
	default:
		return nil, fmt.Errorf("receiver unsupported: %v", r.Type)
	}
}

func (r *Validation) InterfaceFromInt(in int) (interface{}, error) {
	switch r.Type {
	case Int:
		return in, nil
	case Float:
		return float64(in), nil
	case Bool:
		switch in {
		case 1:
			return true, nil
		case 0:
			return false, nil
		}
		return nil, fmt.Errorf("bool must be 0 or 1")
	case Time, TimeUnix:
		date, err := UnixStringToTime(strconv.Itoa(in))
		return date, err
	default:
		return nil, fmt.Errorf("receiver unsupported: %v", r.Type)
	}
}

func (r *Validation) InterfaceFromFloat(in float64) (interface{}, error) {
	switch r.Type {
	case Int:
		return int(in), nil
	case Float:
		return in, nil
	case Bool:
		switch in {
		case 1:
			return true, nil
		case 0:
			return false, nil
		}
		return nil, fmt.Errorf("bool must be 0 or 1")
	case Time, TimeUnix:
		date, err := UnixStringToTime(strconv.Itoa(int(in)))
		return date, err
	default:
		return nil, fmt.Errorf("receiver unsupported: %v", r.Type)
	}
}

func (r *Validation) InterfaceFromArrayOfString(in []string) (interface{}, error) {
	switch r.Type {
	case Array:
		return in, nil
	default:
		// Case for array of strings from url values
		s := ""
		if len(in) > 0 {
			s = in[0]
		}
		return r.InterfaceFromString(s)
	}
}

func UnixStringToTime(in string) (time.Time, error) {
	// Unix seconds date string
	match, _ := regexp.MatchString("^[0-9]{1,}$", in)
	if !match {
		return time.Time{}, fmt.Errorf("invalid unix time: %v", in)
	}

	seconds, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing unix string to time: %v", err)
	}
	return time.Unix(seconds, 0), nil
}

func ISO8601StringToTime(in string) (time.Time, error) {
	layout := ""
	// Iso8601 date string in local time (yyyy-MM-ddTHH:mm:ss.mmmuuu)
	match, _ := regexp.MatchString("^[-:.T0-9]{26}$", in)
	if match {
		layout = "2006-01-02T15:04:05.000000"
	}

	// Iso8601 date string in UTC time (yyyy-MM-ddTHH:mm:ss.mmmuuuZ)
	match, _ = regexp.MatchString("^[-:.T0-9]{26}Z$", in)
	if match {
		layout = "2006-01-02T15:04:05.000000Z"
	}

	// Iso8601 date string in local time without microseconds (yyyy-MM-ddTHH:mm:ss.mmm)
	match, _ = regexp.MatchString("^[-:.T0-9]{23}$", in)
	if match {
		layout = "2006-01-02T15:04:05.000"
	}

	// Iso8601 date string in UTC time without microseconds (yyyy-MM-ddTHH:mm:ss.mmmZ)
	match, _ = regexp.MatchString("^[-:.T0-9]{23}Z$", in)
	if match {
		layout = "2006-01-02T15:04:05.000Z"
	}

	date, err := time.Parse(layout, in)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing iso8601 string to time: %v", err)
	}
	return date, nil
}
