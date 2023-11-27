package validator

import (
	"fmt"
	"regexp"
	"strconv"
)

func checkFloat(f float64, c []string, or bool) error {
	if contains(c, NONE) || len(c) == 0 {
		return nil
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
				equal, err := strconv.ParseFloat(condition, 64)
				if err != nil {
					if or {
						errors = append(errors, err)
					} else {
						return err
					}
				} else if f != equal {
					if or {
						errors = append(errors, fmt.Errorf("value must be equal to %v", equal))
					} else {
						return fmt.Errorf("value must be equal to %v", equal)
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
				notEqual, err := strconv.ParseFloat(condition, 64)
				if err != nil {
					if or {
						errors = append(errors, err)
					} else {
						return err
					}
				} else if f == notEqual {
					if or {
						errors = append(errors, fmt.Errorf("value can't be equal to %v", notEqual))
					} else {
						return fmt.Errorf("value can't be equal to %v", notEqual)
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
				minValue, err := strconv.ParseFloat(condition, 64)
				if err != nil {
					if or {
						errors = append(errors, err)
					} else {
						return err
					}
				} else if f < minValue {
					if or {
						errors = append(errors, fmt.Errorf("value smaller than %v", minValue))
					} else {
						return fmt.Errorf("value smaller than %v", minValue)
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
				maxValue, err := strconv.ParseFloat(condition, 64)
				if err != nil {
					if or {
						errors = append(errors, err)
					} else {
						return err
					}
				} else if f > maxValue {
					if or {
						errors = append(errors, fmt.Errorf("value greater than %v", maxValue))
					} else {
						return fmt.Errorf("value greater than %v", maxValue)
					}
				}
			}
		case REGX:
			condition, err := getConditionByType(conFull, REGX)
			if err != nil {
				if or {
					errors = append(errors, err)
				} else {
					return err
				}
			}
			if len(condition) != 0 {
				match, err := regexp.MatchString(condition, strconv.FormatFloat(f, 'f', 3, 64))
				if err != nil {
					if or {
						errors = append(errors, err)
					} else {
						return err
					}
				} else if !match {
					if or {
						errors = append(errors, fmt.Errorf("value does match regex %v", condition))
					} else {
						return fmt.Errorf("value does match regex %v", condition)
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
