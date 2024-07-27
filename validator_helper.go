package validator

import (
	"fmt"
	"reflect"
	"strings"
)

func ArrayOfInterfaceToArrayOf[T comparable](in []interface{}) ([]T, error) {
	inReflect := reflect.ValueOf(in)
	arrayOfType := []T{}
	for i := 0; i < inReflect.Len(); i++ {
		var newValueInterface interface{}
		switch reflect.TypeOf(arrayOfType).Elem().Kind() {
		case reflect.Int:
			valueOfType, ok := inReflect.Index(i).Interface().(float64)
			if ok {
				newValueInterface = int(valueOfType)
				arrayOfType = append(arrayOfType, newValueInterface.(T))
			} else {
				valueOfType, ok := inReflect.Index(i).Interface().(T)
				if !ok {
					return []T{}, fmt.Errorf("invalid input array element type: %v, expected: %v", reflect.TypeOf(inReflect.Index(i).Interface()).Kind(), reflect.TypeOf(arrayOfType).Elem().Kind())
				}
				arrayOfType = append(arrayOfType, valueOfType)
			}
		case reflect.Int8:
			valueOfType, ok := inReflect.Index(i).Interface().(float64)
			if ok {
				newValueInterface = int8(valueOfType)
				arrayOfType = append(arrayOfType, newValueInterface.(T))
			} else {
				valueOfType, ok := inReflect.Index(i).Interface().(T)
				if !ok {
					return []T{}, fmt.Errorf("invalid input array element type: %v, expected: %v", reflect.TypeOf(inReflect.Index(i).Interface()).Kind(), reflect.TypeOf(arrayOfType).Elem().Kind())
				}
				arrayOfType = append(arrayOfType, valueOfType)
			}
		case reflect.Int16:
			valueOfType, ok := inReflect.Index(i).Interface().(float64)
			if ok {
				newValueInterface = int16(valueOfType)
				arrayOfType = append(arrayOfType, newValueInterface.(T))
			} else {
				valueOfType, ok := inReflect.Index(i).Interface().(T)
				if !ok {
					return []T{}, fmt.Errorf("invalid input array element type: %v, expected: %v", reflect.TypeOf(inReflect.Index(i).Interface()).Kind(), reflect.TypeOf(arrayOfType).Elem().Kind())
				}
				arrayOfType = append(arrayOfType, valueOfType)
			}
		case reflect.Int32:
			valueOfType, ok := inReflect.Index(i).Interface().(float64)
			if ok {
				newValueInterface = int32(valueOfType)
				arrayOfType = append(arrayOfType, newValueInterface.(T))
			} else {
				valueOfType, ok := inReflect.Index(i).Interface().(T)
				if !ok {
					return []T{}, fmt.Errorf("invalid input array element type: %v, expected: %v", reflect.TypeOf(inReflect.Index(i).Interface()).Kind(), reflect.TypeOf(arrayOfType).Elem().Kind())
				}
				arrayOfType = append(arrayOfType, valueOfType)
			}
		case reflect.Int64:
			valueOfType, ok := inReflect.Index(i).Interface().(float64)
			if ok {
				newValueInterface = int64(valueOfType)
				arrayOfType = append(arrayOfType, newValueInterface.(T))
			} else {
				valueOfType, ok := inReflect.Index(i).Interface().(T)
				if !ok {
					return []T{}, fmt.Errorf("invalid input array element type: %v, expected: %v", reflect.TypeOf(inReflect.Index(i).Interface()).Kind(), reflect.TypeOf(arrayOfType).Elem().Kind())
				}
				arrayOfType = append(arrayOfType, valueOfType)
			}
		default:
			valueOfType, ok := inReflect.Index(i).Interface().(T)
			if !ok {
				return []T{}, fmt.Errorf("invalid input array element type: %v, expected: %v", reflect.TypeOf(inReflect.Index(i).Interface()).Kind(), reflect.TypeOf(arrayOfType).Elem().Kind())
			}
			arrayOfType = append(arrayOfType, valueOfType)
		}
	}
	return arrayOfType, nil
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

func getConditionsAndOrFromString(in string) ([]string, bool) {
	or := false
	conditions := strings.Split(in, " ")
	if Contains(conditions, OR) {
		conditions = RemoveWhere(conditions, func(v string) bool {
			return v == OR
		})
		or = true
	}
	return conditions, or
}

func Contains[V comparable](list []V, v V) bool {
	for _, s := range list {
		if v == s {
			return true
		}
	}
	return false
}

func RemoveWhere[V comparable](list []V, f func(V) bool) []V {
	var res []V
	for _, v := range list {
		if !f(v) {
			res = append(res, v)
		}
	}
	return res
}
