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
		// This case is for the case that json.Unmarshal unmarshals an int value into a float64 value.
		if (reflect.TypeOf(arrayOfType).Elem().Kind() == reflect.Int || reflect.TypeOf(arrayOfType).Elem().Kind() == reflect.Int8 || reflect.TypeOf(arrayOfType).Elem().Kind() == reflect.Int16 || reflect.TypeOf(arrayOfType).Elem().Kind() == reflect.Int32 || reflect.TypeOf(arrayOfType).Elem().Kind() == reflect.Int64) && reflect.TypeOf(inReflect.Index(i).Interface()).Kind() == reflect.Float64 {
			valueOfType, ok := inReflect.Index(i).Interface().(float64)
			if !ok {
				return []T{}, fmt.Errorf("invalid input array element type: %v, expected: %v", reflect.TypeOf(inReflect.Index(i).Interface()).Kind(), reflect.TypeOf(arrayOfType).Elem().Kind())
			}
			var newValueInterface interface{}
			if reflect.TypeOf(arrayOfType).Elem().Kind() == reflect.Int {
				newValueInterface = int(valueOfType)
				arrayOfType = append(arrayOfType, newValueInterface.(T))
			} else if reflect.TypeOf(arrayOfType).Elem().Kind() == reflect.Int8 {
				newValueInterface = int8(valueOfType)
				arrayOfType = append(arrayOfType, newValueInterface.(T))
			} else if reflect.TypeOf(arrayOfType).Elem().Kind() == reflect.Int16 {
				newValueInterface = int16(valueOfType)
				arrayOfType = append(arrayOfType, newValueInterface.(T))
			} else if reflect.TypeOf(arrayOfType).Elem().Kind() == reflect.Int32 {
				newValueInterface = int32(valueOfType)
				arrayOfType = append(arrayOfType, newValueInterface.(T))
			} else if reflect.TypeOf(arrayOfType).Elem().Kind() == reflect.Int64 {
				newValueInterface = int64(valueOfType)
				arrayOfType = append(arrayOfType, newValueInterface.(T))
			}
		} else {
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
