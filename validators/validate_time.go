package validators

import (
	"fmt"
	"slices"
	"time"

	"github.com/siherrmann/validator/model"
)

func CheckTime(t time.Time, c []string, or bool) error {
	if slices.Contains(c, model.NONE) || len(c) == 0 {
		return nil
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
				compareTime, err := model.InterfaceFromString(condition, model.Time)
				if err != nil {
					if or {
						errors = append(errors, err)
					} else {
						return err
					}
				}
				if !t.Equal(compareTime.(time.Time)) {
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
				compareTime, err := model.InterfaceFromString(condition, model.Time)
				if err != nil {
					if or {
						errors = append(errors, err)
					} else {
						return err
					}
				}
				if t.Equal(compareTime.(time.Time)) {
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
				minValue, err := model.InterfaceFromString(condition, model.Time)
				if err != nil {
					if or {
						errors = append(errors, err)
					} else {
						return err
					}
				}
				if t.Before(minValue.(time.Time)) {
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
				maxValue, err := model.InterfaceFromString(condition, model.Time)
				if err != nil {
					if or {
						errors = append(errors, err)
					} else {
						return err
					}
				}
				if t.After(maxValue.(time.Time)) {
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
					from, err := model.InterfaceFromString(fromValue, model.Time)
					if err != nil {
						if or {
							errors = append(errors, err)
						} else {
							return false
						}
					}
					return t.Equal(from.(time.Time))
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
					notFrom, err := model.InterfaceFromString(notFromValue, model.Time)
					if err != nil {
						if or {
							errors = append(errors, err)
						} else {
							return false
						}
					}
					return t.Equal(notFrom.(time.Time))
				})
				if foundInFromValues {
					if or {
						errors = append(errors, fmt.Errorf("value found in %v", notFromValues))
					} else {
						return fmt.Errorf("value found in %v", notFromValues)
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
