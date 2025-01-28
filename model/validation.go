package model

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Validation struct {
	Key         string
	Type        ValidatorType
	Requirement string
	Groups      []*Group
	Default     string
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

	if len(tagSplit) > tagIndex {
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

func (r *Validation) FillOnlyKey(tag string, tagType TagType, value reflect.Value) error {
	r.Type = TypeFromInterface(value.Interface())

	tagSplit := strings.Split(tag, ", ")
	r.Requirement = "-"

	tagIndex := 0
	if len(tagSplit) > tagIndex && tagType == UPD {
		r.Key = strings.TrimSpace(tagSplit[tagIndex])
		tagIndex++
	}

	return nil
}

func (r *Validation) GetValidValue(in interface{}) (interface{}, error) {
	var out interface{}
	var err error
	switch in := in.(type) {
	case string:
		out, err = InterfaceFromString(in, r.Type)
	case int:
		out, err = InterfaceFromInt(in, r.Type)
	case int64:
		out, err = InterfaceFromInt(int(in), r.Type)
	case int32:
		out, err = InterfaceFromInt(int(in), r.Type)
	case int16:
		out, err = InterfaceFromInt(int(in), r.Type)
	case int8:
		out, err = InterfaceFromInt(int(in), r.Type)
	case float64:
		out, err = InterfaceFromFloat(in, r.Type)
	case float32:
		out, err = InterfaceFromFloat(float64(in), r.Type)
	case []interface{}:
		return in, nil
	default:
		if v, ok := in.(JsonMap); ok {
			return map[string]interface{}(v), nil
		}
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
		if reflect.ValueOf(in).Kind() == reflect.Array || reflect.ValueOf(in).Kind() == reflect.Slice {
			return Array
		} else {
			return Struct
		}
	}
}

func InterfaceFromString(in string, inType ValidatorType) (interface{}, error) {
	switch inType {
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
		err := json.Unmarshal([]byte(in), &out)
		if err != nil {
			return nil, err
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
		return nil, fmt.Errorf("receiver unsupported: %v", inType)
	}
}

func InterfaceFromInt(in int, inType ValidatorType) (interface{}, error) {
	switch inType {
	case Int:
		return in, nil
	case Float:
		return float64(in), nil
	case Bool:
		if in == 1 {
			return true, nil
		} else if in == 0 {
			return false, nil
		}
		return nil, fmt.Errorf("bool must be 0 or 1")
	case Time, TimeUnix:
		date, err := UnixStringToTime(strconv.Itoa(in))
		return date, err
	default:
		return nil, fmt.Errorf("receiver unsupported: %v", inType)
	}
}

func InterfaceFromFloat(in float64, inType ValidatorType) (interface{}, error) {
	switch inType {
	case Int:
		return int(in), nil
	case Float:
		return in, nil
	case Bool:
		if in == 1 {
			return true, nil
		} else if in == 0 {
			return false, nil
		}
		return nil, fmt.Errorf("bool must be 0 or 1")
	case Time, TimeUnix:
		date, err := UnixStringToTime(strconv.Itoa(int(in)))
		return date, err
	default:
		return nil, fmt.Errorf("receiver unsupported: %v", inType)
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
