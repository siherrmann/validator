package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type JsonMap map[string]interface{}

func (c JsonMap) Value() (driver.Value, error) {
	return c.Marshal()
}

func (c *JsonMap) Scan(value interface{}) error {
	return c.Unmarshal(value)
}

func (r JsonMap) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *JsonMap) Unmarshal(value interface{}) error {
	if s, ok := value.(map[string]interface{}); ok {
		*r = JsonMap(s)
	} else {
		b, ok := value.([]byte)
		if !ok {
			return errors.New("type assertion to []byte failed")
		}
		return json.Unmarshal(b, r)
	}
	return nil
}

func (r *JsonMap) Has(key string) bool {
	if _, ok := (*r)[key]; ok {
		return true
	}
	return false
}

func TypeFromInterface(in interface{}) (ValidatorType, error) {
	switch in.(type) {
	case string:
		return String, nil
	case int:
		return Int, nil
	case float64:
		return Float, nil
	case bool:
		return Bool, nil
	case map[string]interface{}:
		return Map, nil
	case time.Time:
		return Time, nil
	default:
		return "", fmt.Errorf("input type unsupported: %T", in)
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
	case Map:
		jsonMarshalled, err := json.Marshal(in)
		if err != nil {
			return nil, err
		}
		out := JsonMap{}
		err = json.Unmarshal(jsonMarshalled, &out)
		if err != nil {
			return nil, err
		}
		return out, nil
	case Time:
		// Unix seconds date string
		match, _ := regexp.MatchString("^[0-9]{1,}$", in)
		if match {
			seconds, err := strconv.ParseInt(in, 10, 64)
			if err != nil {
				panic(err)
			}
			return time.Unix(seconds, 0), nil
		}

		layout := ""
		// Iso8601 date string in local time (yyyy-MM-ddTHH:mm:ss.mmmuuu)
		match, _ = regexp.MatchString("^[-:.T0-9]{26}$", in)
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
			return nil, fmt.Errorf("input string could not be converted to time.Time with err: %v", err.Error())
		}
		return date, nil
	default:
		return nil, fmt.Errorf("receiver unsupported: %v", inType)
	}
}
