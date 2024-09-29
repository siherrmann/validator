package validators

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/siherrmann/validator/model"
)

func CheckString(s string, c []string, or bool) error {
	if slices.Contains(c, string(model.NONE)) || len(c) == 0 {
		return nil
	}

	var errors []error
	for _, conFull := range c {
		conType, err := model.GetConditionType(conFull)
		if err != nil {
			if or {
				errors = append(errors, err)
			} else {
				return err
			}
		}

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
				if s != condition {
					if or {
						errors = append(errors, fmt.Errorf("value must be equal to %v", condition))
					} else {
						return fmt.Errorf("value must be equal to %v", condition)
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
				if s == condition {
					if or {
						errors = append(errors, fmt.Errorf("value can't be equal to %v", condition))
					} else {
						return fmt.Errorf("value can't be equal to %v", condition)
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
				} else if len(strings.TrimSpace(s)) < minValue {
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
				} else if len(strings.TrimSpace(s)) > maxValue {
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
				if !strings.Contains(s, condition) {
					if or {
						errors = append(errors, fmt.Errorf("value does not contain %v", condition))
					} else {
						return fmt.Errorf("value does not contain %v", condition)
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
				if strings.Contains(s, condition) {
					if or {
						errors = append(errors, fmt.Errorf("value does contain %v", condition))
					} else {
						return fmt.Errorf("value does contain %v", condition)
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
				if !slices.Contains(fromValues, s) {
					if or {
						errors = append(errors, fmt.Errorf("value not found in %v", fromValues))
					} else {
						return fmt.Errorf("value not found in %v", fromValues)
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
				if slices.Contains(notFromValues, s) {
					if or {
						errors = append(errors, fmt.Errorf("value found in %v", notFromValues))
					} else {
						return fmt.Errorf("value found in %v", notFromValues)
					}
				}
			}
		case model.REGX:
			condition, err := model.GetConditionByType(conFull, model.REGX)
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
