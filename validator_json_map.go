package validator

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/siherrmann/validator/model"
	"github.com/siherrmann/validator/validators"
)

func NewJsonMapFromUrlValues(values url.Values, validation ...model.Validation) (model.JsonMap, error) {
	mapOut := model.JsonMap{}
	for _, v := range validation {
		if values.Has(v.Key) {
			value := values.Get(v.Key)
			out, err := InterfaceFromString(value, v.Type, v.Conditions)
			if err != nil {
				return nil, err
			}
			mapOut[v.Key] = out
		}
	}
	return mapOut, nil
}

func InterfaceFromString(in string, inType model.ValidatorType, condition string) (interface{}, error) {
	or := false
	conditions, or := model.GetConditionsAndOrFromString(condition)

	switch inType {
	case model.String:
		err := validators.CheckString(in, conditions, or)
		if err != nil {
			return nil, err
		}
		return in, nil
	case model.Int:
		out, err := strconv.Atoi(in)
		if err != nil {
			return nil, err
		}
		err = validators.CheckInt(out, conditions, or)
		if err != nil {
			return nil, err
		}
		return out, nil
	case model.Float:
		out, err := strconv.ParseFloat(strings.TrimSpace(in), 64)
		if err != nil {
			return nil, err
		}
		err = validators.CheckFloat(out, conditions, or)
		if err != nil {
			return nil, err
		}
		return out, nil
	case model.Bool:
		out, err := strconv.ParseBool(in)
		if err != nil {
			return nil, err
		}
		return out, nil
	case model.Map:
		jsonMarshalled, err := json.Marshal(in)
		if err != nil {
			return nil, err
		}
		out := model.JsonMap{}
		json.Unmarshal(jsonMarshalled, &out)
		err = validators.CheckMap(reflect.ValueOf(out), conditions, or)
		if err != nil {
			return nil, err
		}
		return out, nil
	case model.Time:
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
