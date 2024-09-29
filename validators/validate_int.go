package validators

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"

	"github.com/siherrmann/validator/model"
)

func CheckInt(i int, c []string, or bool) error {
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
				equal, err := strconv.Atoi(condition)
				if err != nil {
					if or {
						errors = append(errors, err)
					} else {
						return err
					}
				} else if i != equal {
					if or {
						errors = append(errors, fmt.Errorf("value must be equal to %v", equal))
					} else {
						return fmt.Errorf("value must be equal to %v", equal)
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
				} else if i == notEqual {
					if or {
						errors = append(errors, fmt.Errorf("value can't be equal to %v", notEqual))
					} else {
						return fmt.Errorf("value can't be equal to %v", notEqual)
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
				} else if i < minValue {
					if or {
						errors = append(errors, fmt.Errorf("value smaller than %v", minValue))
					} else {
						return fmt.Errorf("value smaller than %v", minValue)
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
				} else if i > maxValue {
					if or {
						errors = append(errors, fmt.Errorf("value greater than %v", maxValue))
					} else {
						return fmt.Errorf("value greater than %v", maxValue)
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
				foundInFromValues := slices.ContainsFunc(fromValues, func(fromValue string) bool {
					from, err := strconv.Atoi(fromValue)
					if err != nil {
						if or {
							errors = append(errors, err)
						} else {
							return false
						}
					}
					return i == from
				})
				if !foundInFromValues {
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
				foundInFromValues := slices.ContainsFunc(notFromValues, func(notFromValue string) bool {
					notFrom, err := strconv.Atoi(notFromValue)
					if err != nil {
						if or {
							errors = append(errors, err)
						} else {
							return false
						}
					}
					return i == notFrom
				})
				if foundInFromValues {
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
				match, err := regexp.MatchString(condition, strconv.Itoa(i))
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
