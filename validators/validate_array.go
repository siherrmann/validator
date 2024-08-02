package validators

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"

	"github.com/siherrmann/validator/model"
)

func CheckArray(a reflect.Value, c []string, or bool) error {
	if slices.Contains(c, model.NONE) || len(c) == 0 {
		return nil
	}

	if a.Type().Kind() != reflect.Array && a.Type().Kind() != reflect.Slice {
		return fmt.Errorf("value to validate has to be a array or slice, was %v", a.Type().Kind())
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
			if len(condition) != 0 {
				contains, err := ValueContains(a, condition)
				if err != nil {
					if or {
						errors = append(errors, err)
					} else {
						return err
					}
				} else if len(contains) == 0 {
					if or {
						errors = append(errors, fmt.Errorf("value does not contain %v", contains))
					} else {
						return fmt.Errorf("value does not contain %v", contains)
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
			if len(condition) != 0 {
				contains, err := ValueContains(a, condition)
				if err != nil {
					if or {
						errors = append(errors, err)
					} else {
						return err
					}
				} else if len(contains) != 0 {
					if or {
						errors = append(errors, fmt.Errorf("value does contain %v", contains))
					} else {
						return fmt.Errorf("value does contain %v", contains)
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
			if len(condition) != 0 {
				fromValues, err := model.GetArrayFromCondition(condition)
				if err != nil {
					if or {
						errors = append(errors, err)
					} else {
						return err
					}
				}
				notFound, err := ValueFrom(a, fromValues, true)
				if err != nil {
					if or {
						errors = append(errors, err)
					} else {
						return err
					}
				} else if len(notFound) != 0 {
					if or {
						errors = append(errors, fmt.Errorf("from values do not contain %v", notFound))
					} else {
						return fmt.Errorf("from values do not contain %v", notFound)
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
			if len(condition) != 0 {
				notFromValues, err := model.GetArrayFromCondition(condition)
				if err != nil {
					if or {
						errors = append(errors, err)
					} else {
						return err
					}
				}
				found, err := ValueFrom(a, notFromValues, false)
				if err != nil {
					if or {
						errors = append(errors, err)
					} else {
						return err
					}
				} else if len(found) != 0 {
					if or {
						errors = append(errors, fmt.Errorf("notFrom values do contain %v", found))
					} else {
						return fmt.Errorf("notFrom values do contain %v", found)
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

// ValueContains checks, if the [reflectValue] contains [contain].
func ValueContains(reflectValue reflect.Value, contain string) (string, error) {
	switch v := reflectValue.Interface().(type) {
	case []int:
		c, err := strconv.Atoi(contain)
		if err != nil {
			return "", err
		} else if slices.Contains(v, c) {
			return contain, nil
		}
		return "", nil
	case []float32:
		c, err := strconv.ParseFloat(contain, 32)
		if err != nil {
			return "", err
		} else if slices.Contains(v, float32(c)) {
			return contain, nil
		}
		return "", nil
	case []float64:
		c, err := strconv.ParseFloat(contain, 64)
		if err != nil {
			return "", err
		} else if slices.Contains(v, c) {
			return contain, nil
		}
		return "", nil
	case []string:
		if slices.Contains(v, contain) {
			return contain, nil
		}
		return "", nil
	default:
		return "", fmt.Errorf("type %v not supported", reflect.TypeOf(v))
	}
}

// ValueFrom checks, if the [reflectValue] consists of only ([shouldFind] == true)/of none ([shouldFind] == false) of the values from [from].
// If a value is not found with [shouldFind] == true the missing value is given back.
// If a value is found with [shouldFind] == false the value found is given back.
func ValueFrom(reflectValue reflect.Value, from []string, shouldFind bool) (string, error) {
	switch v := reflectValue.Interface().(type) {
	case []int:
		for _, value := range v {
			contains := slices.ContainsFunc(from, func(fromValue string) bool {
				c, err := strconv.Atoi(fromValue)
				if err != nil {
					return false
				}
				return c == value
			})
			if contains != shouldFind {
				return fmt.Sprint(value), nil
			}
		}
	case []float32:
		for _, value := range v {
			contains := slices.ContainsFunc(from, func(fromValue string) bool {
				c, err := strconv.ParseFloat(fromValue, 32)
				if err != nil {
					return false
				}
				return float32(c) == value
			})
			if contains != shouldFind {
				return fmt.Sprint(value), nil
			}
		}
	case []float64:
		for _, value := range v {
			contains := slices.ContainsFunc(from, func(fromValue string) bool {
				c, err := strconv.ParseFloat(fromValue, 64)
				if err != nil {
					return false
				}
				return c == value
			})
			if contains != shouldFind {
				return fmt.Sprint(value), nil
			}
		}
	case []string:
		for _, value := range v {
			contains := slices.Contains(from, value)
			if contains != shouldFind {
				return fmt.Sprint(value), nil
			}
		}
	default:
		return "", fmt.Errorf("type %v not supported", reflect.TypeOf(v))
	}
	return "", nil
}
