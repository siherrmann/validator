package validator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func checkString(s string, c []string, or bool) error {
	if Contains(c, NONE) || len(c) == 0 {
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
				if s != condition {
					if or {
						errors = append(errors, fmt.Errorf("value must be equal to %v", condition))
					} else {
						return fmt.Errorf("value must be equal to %v", condition)
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
				if s == condition {
					if or {
						errors = append(errors, fmt.Errorf("value can't be equal to %v", condition))
					} else {
						return fmt.Errorf("value can't be equal to %v", condition)
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
				} else if len(strings.TrimSpace(s)) < minValue {
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
				} else if len(strings.TrimSpace(s)) > maxValue {
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
				if !strings.Contains(s, condition) {
					if or {
						errors = append(errors, fmt.Errorf("value does not include %v", condition))
					} else {
						return fmt.Errorf("value does not include %v", condition)
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
				match, err := regexp.MatchString(condition, s)
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
