package validator

import (
	"fmt"
	"reflect"
	"strconv"
)

func checkArray(a reflect.Value, c []string, or bool) error {
	if contains(c, NONE) || len(c) == 0 {
		return nil
	}

	if a.Type().Kind() != reflect.Array && a.Type().Kind() != reflect.Slice {
		return fmt.Errorf("value to validate has to be a array or slice, was %v", a.Type().Kind())
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
			if len(condition) != 0 {
				switch v := a.Interface().(type) {
				case []int:
					contain, err := strconv.Atoi(condition)
					if err != nil {
						if or {
							errors = append(errors, err)
						} else {
							return err
						}
					} else if !contains(v, contain) {
						if or {
							errors = append(errors, fmt.Errorf("value does not contain %v", contain))
						} else {
							return fmt.Errorf("value does not contain %v", contain)
						}
					}
				case []float32:
					contain, err := strconv.ParseFloat(condition, 32)
					if err != nil {
						if or {
							errors = append(errors, err)
						} else {
							return err
						}
					} else if !contains(v, float32(contain)) {
						if or {
							errors = append(errors, fmt.Errorf("value does not contain %v", contain))
						} else {
							return fmt.Errorf("value does not contain %v", contain)
						}
					}
				case []float64:
					contain, err := strconv.ParseFloat(condition, 64)
					if err != nil {
						if or {
							errors = append(errors, err)
						} else {
							return err
						}
					} else if !contains(v, contain) {
						if or {
							errors = append(errors, fmt.Errorf("value does not contain %v", contain))
						} else {
							return fmt.Errorf("value does not contain %v", contain)
						}
					}
				case []string:
					if !contains(v, condition) {
						if or {
							errors = append(errors, fmt.Errorf("value does not contain %v", condition))
						} else {
							return fmt.Errorf("value does not contain %v", condition)
						}
					}
				default:
					return fmt.Errorf("type %v not supported", reflect.TypeOf(v))
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
