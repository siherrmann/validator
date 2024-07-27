package validator

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Validation struct {
	Key        string
	Type       ValidatorType
	Conditions string
}

type ValidatorMap map[string]Validation

type ValidatorType string

const (
	String ValidatorType = "string"
	Int    ValidatorType = "int"
	Float  ValidatorType = "float"
	Bool   ValidatorType = "bool"
	Time   ValidatorType = "time"
)

func (r ValidatorType) InterfaceFromString(in string, condition string) (interface{}, error) {
	or := false
	conditions, or := getConditionsAndOrFromString(condition)

	switch r {
	case String:
		return in, nil
	case Int:
		out, err := strconv.Atoi(in)
		if err != nil {
			return nil, err
		}
		err = checkInt(out, conditions, or)
		if err != nil {
			return nil, err
		}
		return out, nil
	case Float:
		out, err := strconv.ParseFloat(strings.TrimSpace(in), 64)
		if err != nil {
			return nil, err
		}
		err = checkFloat(out, conditions, or)
		if err != nil {
			return nil, err
		}
		return out, nil
	case Bool:
		out, err := strconv.ParseBool(in)
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
		return nil, fmt.Errorf("receiver unsupported: %v", r)
	}
}

type JsonMap map[string]interface{}

func (a JsonMap) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *JsonMap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

func NewJsonMapFromUrlValues(values url.Values, validation ...Validation) (JsonMap, error) {
	mapOut := JsonMap{}
	for _, v := range validation {
		if values.Has(v.Key) {
			value := values.Get(v.Key)
			out, err := v.Type.InterfaceFromString(value, v.Conditions)
			if err != nil {
				return nil, err
			}
			mapOut[v.Key] = out
		}
	}
	return mapOut, nil
}

func checkMap(a reflect.Value, c []string, or bool) error {
	if Contains(c, NONE) || len(c) == 0 {
		return nil
	}

	if a.Type().Kind() != reflect.Map {
		return fmt.Errorf("value to validate has to be a map, was %v", a.Type().Kind())
	}

	var errors []error
	for _, conFull := range c {
		conType := getConditionType(conFull)

		switch conType {
		case EQUAL:
			condition, err := getConditionByType(conFull, EQUAL)
			if err != nil {
				if or {
					errors = append(errors, err)
				} else {
					return err
				}
			}
			if len(condition) != 0 {
				equal, err := strconv.Atoi(condition)
				if err != nil {
					if or {
						errors = append(errors, err)
					} else {
						return err
					}
				} else if a.Len() != equal {
					if or {
						errors = append(errors, fmt.Errorf("value shorter than %v", equal))
					} else {
						return fmt.Errorf("value shorter than %v", equal)
					}
				}
			}
		case NOT_EQUAL:
			condition, err := getConditionByType(conFull, NOT_EQUAL)
			if err != nil {
				if or {
					errors = append(errors, err)
				} else {
					return err
				}
			}
			if len(condition) != 0 {
				notEqual, err := strconv.Atoi(condition)
				if err != nil {
					if or {
						errors = append(errors, err)
					} else {
						return err
					}
				} else if a.Len() == notEqual {
					if or {
						errors = append(errors, fmt.Errorf("value longer than %v", notEqual))
					} else {
						return fmt.Errorf("value longer than %v", notEqual)
					}
				}
			}
		case MIN_VALUE:
			condition, err := getConditionByType(conFull, MIN_VALUE)
			if err != nil {
				if or {
					errors = append(errors, err)
				} else {
					return err
				}
			}
			if len(condition) != 0 {
				minValue, err := strconv.Atoi(condition)
				if err != nil {
					if or {
						errors = append(errors, err)
					} else {
						return err
					}
				} else if a.Len() < minValue {
					if or {
						errors = append(errors, fmt.Errorf("value shorter than %v", minValue))
					} else {
						return fmt.Errorf("value shorter than %v", minValue)
					}
				}
			}
		case MAX_VLAUE:
			condition, err := getConditionByType(conFull, MAX_VLAUE)
			if err != nil {
				if or {
					errors = append(errors, err)
				} else {
					return err
				}
			}
			if len(condition) != 0 {
				maxValue, err := strconv.Atoi(condition)
				if err != nil {
					if or {
						errors = append(errors, err)
					} else {
						return err
					}
				} else if a.Len() > maxValue {
					if or {
						errors = append(errors, fmt.Errorf("value longer than %v", maxValue))
					} else {
						return fmt.Errorf("value longer than %v", maxValue)
					}
				}
			}
		case CONTAINS:
			condition, err := getConditionByType(conFull, CONTAINS)
			if err != nil {
				if or {
					errors = append(errors, err)
				} else {
					return err
				}
			}
			var jsonMap JsonMap
			var ok bool
			if jsonMap, ok = a.Interface().(map[string]any); !ok {
				return fmt.Errorf("value has to be of type map, was %v", reflect.TypeOf(a))
			}
			if len(condition) != 0 {
				if _, ok := jsonMap[condition]; !ok {
					if or {
						errors = append(errors, fmt.Errorf("value does not contain key %v", condition))
					} else {
						return fmt.Errorf("value does not contain key %v", condition)
					}
				}
			}
		case NONE:
			return nil
		case OR:
			continue
		default:
			return fmt.Errorf("invalid condition type %s", conType)
		}
	}

	if len(errors) >= len(c) {
		return fmt.Errorf("no condition fulfilled, all errors: %v", errors)
	}

	return nil
}
