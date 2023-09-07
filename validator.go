package validator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	NONE      string = "-"
	EQUAL     string = "equ"
	NOT_EQUAL string = "neq"
	MIN_VALUE string = "min"
	MAX_VLAUE string = "max"
	CONTAINS  string = "con"
	REGX      string = "rex"
	OR        string = "||"
)

type StructValue struct {
	Error  error
	Groups []string
}

// UnmarshalAndValidate unmarshals given json ([]byte) into pointer v.
// For more information to Validate look at Validate(value any).
func UnmarshalAndValidate(data []byte, v any) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("value has to be of kind pointer, was %T", value)
	}
	if value.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("value has to be of kind struct, was %T", value)
	}

	err := json.Unmarshal(data, v)
	if err != nil {
		return fmt.Errorf("error unmarshalling %T: %v", value, err)
	}

	err = Validate(v)
	if err != nil {
		return err
	}

	return nil
}

// Validate validates a given struct by vld tags.
// Validate needs a struct as input.
//
// All fields in the struct need a vld tag.
// If you want to use multiple conditions you can add them with a space in between them.
//
// A complex example for password would be:
// `vld:"min8 max30 rex^(.*[A-Z])+(.*)$ rex^(.*[a-z])+(.*)$ rex^(.*\\d)+(.*)$ rex^(.*[\x60!@#$%^&*()_+={};':\"|\\,.<>/?~-])+(.*)$"`
//
// If you want to ignore one field in the validator you can add `vld:"-"`.
// If you don't add the vld tag to every field the function will fail with an error.
//
// Conditions have different usages per variable type:
//
// equ - int/float/string == condition, len(array) == condition
//
// neq - int/float/string != condition, len(array) != condition
//
// min - int/float >= condition, len(string/array) >= condition
//
// max - int/float <= condition, len(string/array) <= condition
//
// con - strings.Contains(string, condition), contains(array, condition), int/float ignored
//
// rex - regexp.MatchString(condition, int/float/string), array ignored
//
// For con you need to put in a condition that is convertable to the underlying type of the arrary.
// Eg. for an array of int the condition must be convertable to int (bad: `vld:"conA"`, good: `vld:"con1"`).
//
// In the case of rex the int and float input will get converted to a string (strconv.Itoa(int) and fmt.Sprintf("%f", f)).
// If you want to check more complex cases you can obviously replace equ, neq, min, max and con with one regular expression.
func Validate(v any) error {
	// check if value is a pointer to a struct
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("value has to be of kind pointer, was %T", value)
	}
	if value.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("value has to be of kind struct, was %T", value)
	}

	// get valid reflect value of struct
	structFull := value.Elem()

	groups := map[string]string{}
	groupErrors := map[string][]error{}
	groupSize := map[string]int{}

	for i := 0; i < structFull.Type().NumField(); i++ {
		tag := structFull.Type().Field(i).Tag.Get("vld")
		if len(strings.TrimSpace(tag)) == 0 {
			return fmt.Errorf("no validate tag found for field %s, to ignore the validation put a 'vld:\"-\"' into the tag", structFull.Type().Field(i).Name)
		}

		tagSplit := strings.Split(tag, ", ")
		groupsValue := []string{}
		groupsString := []string{}
		if len(tagSplit) > 1 {
			groupsString = strings.Split(tagSplit[1], " ")

			for _, g := range groupsString {
				group := getConditionType(g)
				condition, err := getConditionByType(g, group)
				if err != nil {
					return fmt.Errorf("error extracting group: %v", err)
				}

				groupsValue = append(groupsValue, group)
				groups[group] = condition
				groupSize[group]++
			}
		}

		value := structFull.Field(i)
		fieldName := structFull.Type().Field(i).Name
		conditions := strings.Split(tagSplit[0], " ")

		if contains(conditions, NONE) {
			continue
		}

		or := false
		if contains(conditions, OR) {
			or = true
		}

		switch value.Type().Kind() {
		case reflect.Int:
			valueTemp := value.Int()
			err := checkInt(int(valueTemp), conditions, or)
			if err != nil && len(groupsString) == 0 {
				return fmt.Errorf("field %v invalid: %v", fieldName, err.Error())
			} else if err != nil {
				for _, groupName := range groupsValue {
					groupErrors[groupName] = append(groupErrors[groupName], fmt.Errorf("field %v invalid: %v", fieldName, err.Error()))
				}
			}
		case reflect.Float32:
		case reflect.Float64:
			valueTemp := value.Float()
			err := checkFloat(valueTemp, conditions, or)
			if err != nil && len(groupsString) == 0 {
				return fmt.Errorf("field %v invalid: %v", fieldName, err.Error())
			} else if err != nil {
				for _, groupName := range groupsValue {
					groupErrors[groupName] = append(groupErrors[groupName], fmt.Errorf("field %v invalid: %v", fieldName, err.Error()))
				}
			}
		case reflect.String:
			valueTemp := value.String()
			err := checkString(valueTemp, conditions, or)
			if err != nil && len(groupsString) == 0 {
				return fmt.Errorf("field %v invalid: %v", fieldName, err.Error())
			} else if err != nil {
				for _, groupName := range groupsValue {
					groupErrors[groupName] = append(groupErrors[groupName], fmt.Errorf("field %v invalid: %v", fieldName, err.Error()))
				}
			}
		case reflect.Array:
		case reflect.Slice:
			valueTemp := value
			err := checkArray(valueTemp, conditions, or)
			if err != nil && len(groupsString) == 0 {
				return fmt.Errorf("field %v invalid: %v", fieldName, err.Error())
			} else if err != nil {
				for _, groupName := range groupsValue {
					groupErrors[groupName] = append(groupErrors[groupName], fmt.Errorf("field %v invalid: %v", fieldName, err.Error()))
				}
			}
		default:
			return fmt.Errorf("invalid field type: %v", value.Type().Kind())
		}
	}

	if len(groups) != 0 {
		for groupName, groupCondition := range groups {
			conType := getConditionType(groupCondition)

			switch conType {
			case MIN_VALUE:
				condition, err := getConditionByType(groupCondition, MIN_VALUE)
				if err != nil {
					return err
				}
				if len(condition) != 0 {
					minValue, err := strconv.Atoi(condition)
					if err != nil {
						return err
					} else if (groupSize[groupName] - len(groupErrors[groupName])) < minValue {
						return fmt.Errorf("less the %v in group without error, all errors: %v", minValue, groupErrors[groupName])
					}
				}
			case MAX_VLAUE:
				condition, err := getConditionByType(groupCondition, MAX_VLAUE)
				if err != nil {
					return err
				}
				if len(condition) != 0 {
					maxValue, err := strconv.Atoi(condition)
					if err != nil {
						return err
					} else if (groupSize[groupName] - len(groupErrors[groupName])) > maxValue {
						return fmt.Errorf("more the %v in group without error, all errors: %v", maxValue, groupErrors[groupName])
					}
				}
			default:
				return fmt.Errorf("invalid group condition type %s", conType)
			}
		}
	}

	return nil
}

func checkInt(i int, c []string, or bool) error {
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
				} else if i != equal {
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
				} else if i < minValue {
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

func checkFloat(f float64, c []string, or bool) error {
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
		default:
			return fmt.Errorf("invalid condition type %s", conType)
		}
	}

	if len(errors) >= len(c) {
		return fmt.Errorf("no condition fulfilled, all errors: %v", errors)
	}

	return nil
}

func checkString(s string, c []string, or bool) error {
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
		default:
			return fmt.Errorf("invalid condition type %s", conType)
		}
	}

	if len(errors) >= len(c) {
		return fmt.Errorf("no condition fulfilled, all errors: %v", errors)
	}

	return nil
}

func checkArray(a reflect.Value, c []string, or bool) error {
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
		default:
			return fmt.Errorf("invalid condition type %s", conType)
		}
	}

	if len(errors) >= len(c) {
		return fmt.Errorf("no condition fulfilled, all errors: %v", errors)
	}

	return nil
}

func getConditionType(s string) string {
	if len(s) > 2 {
		return s[:3]
	}
	return s
}

func getConditionByType(conditionFull string, conditionType string) (string, error) {
	if len(conditionType) != 3 {
		return "", fmt.Errorf("length of conditionType has to be 3: %s", conditionType)
	}
	condition := strings.TrimPrefix(conditionFull, conditionType)
	if len(condition) == 0 {
		return "", fmt.Errorf("empty %s value", conditionType)
	}

	return condition, nil
}

func contains[V comparable](list []V, v V) bool {
	for _, s := range list {
		if v == s {
			return true
		}
	}
	return false
}
