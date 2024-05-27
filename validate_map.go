package validator

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

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
