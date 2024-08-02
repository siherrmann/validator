package validators

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"

	"github.com/siherrmann/validator/model"
)

func CheckMap(a reflect.Value, c []string, or bool) error {
	if slices.Contains(c, model.NONE) || len(c) == 0 {
		return nil
	}

	if a.Type().Kind() != reflect.Map {
		return fmt.Errorf("value to validate has to be a map, was %v", a.Type().Kind())
	}

	var errors []error
	for _, conFull := range c {
		conType := model.GetConditionType(conFull)

		switch conType {
		case model.EQUAL:
			condition, err := model.GetConditionByType(conFull, model.EQUAL)
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
		case model.NOT_EQUAL:
			condition, err := model.GetConditionByType(conFull, model.NOT_EQUAL)
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
		case model.MIN_VALUE:
			condition, err := model.GetConditionByType(conFull, model.MIN_VALUE)
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
		case model.MAX_VLAUE:
			condition, err := model.GetConditionByType(conFull, model.MAX_VLAUE)
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
		case model.CONTAINS:
			condition, err := model.GetConditionByType(conFull, model.CONTAINS)
			if err != nil {
				if or {
					errors = append(errors, err)
				} else {
					return err
				}
			}
			var jsonMap model.JsonMap
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
		case model.NOT_CONTAINS:
			condition, err := model.GetConditionByType(conFull, model.NOT_CONTAINS)
			if err != nil {
				if or {
					errors = append(errors, err)
				} else {
					return err
				}
			}
			var jsonMap model.JsonMap
			var ok bool
			if jsonMap, ok = a.Interface().(map[string]any); !ok {
				return fmt.Errorf("value has to be of type map, was %v", reflect.TypeOf(a))
			}
			if len(condition) != 0 {
				if _, ok := jsonMap[condition]; ok {
					if or {
						errors = append(errors, fmt.Errorf("value does contain key %v", condition))
					} else {
						return fmt.Errorf("value does contain key %v", condition)
					}
				}
			}
		case model.FROM:
			condition, err := model.GetConditionByType(conFull, model.FROM)
			if err != nil {
				if or {
					errors = append(errors, err)
				} else {
					return err
				}
			}
			var jsonMap model.JsonMap
			var ok bool
			if jsonMap, ok = a.Interface().(map[string]any); !ok {
				return fmt.Errorf("value has to be of type map, was %v", reflect.TypeOf(a))
			}
			if len(condition) != 0 {
				fromValues, err := model.GetArrayFromCondition(condition)
				if err != nil {
					if or {
						errors = append(errors, err)
					} else {
						return err
					}
				}
				for k := range jsonMap {
					if !slices.Contains(fromValues, k) {
						if or {
							errors = append(errors, fmt.Errorf("key not found in %v", fromValues))
						} else {
							return fmt.Errorf("key not found in %v", fromValues)
						}
					}
				}
			}
		case model.NOT_FROM:
			condition, err := model.GetConditionByType(conFull, model.NOT_FROM)
			if err != nil {
				if or {
					errors = append(errors, err)
				} else {
					return err
				}
			}
			var jsonMap model.JsonMap
			var ok bool
			if jsonMap, ok = a.Interface().(map[string]any); !ok {
				return fmt.Errorf("value has to be of type map, was %v", reflect.TypeOf(a))
			}
			if len(condition) != 0 {
				fromValues, err := model.GetArrayFromCondition(condition)
				if err != nil {
					if or {
						errors = append(errors, err)
					} else {
						return err
					}
				}
				for k := range jsonMap {
					if slices.Contains(fromValues, k) {
						if or {
							errors = append(errors, fmt.Errorf("key found in %v", fromValues))
						} else {
							return fmt.Errorf("key found in %v", fromValues)
						}
					}
				}
			}
		case model.NONE:
			return nil
		case model.OR:
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
