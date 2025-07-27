package validators

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"
)

func Equal[T comparable](a, equal T) bool {
	return a == equal
}

func NotEqual[T comparable](a, notEqual T) bool {
	return a != notEqual
}

func Min[T int | float64](v, min T) bool {
	return v >= min
}

func Max[T int | float64](v, max T) bool {
	return v <= max
}

func Contains(va any, contain string) (bool, error) {
	if s, ok := va.(string); ok {
		return strings.Contains(s, contain), nil
	}

	rv := reflect.ValueOf(va)
	rt := rv.Type()
	switch rt.Kind() {
	case reflect.Array, reflect.Slice:
		if rv.Len() == 0 {
			return false, nil
		}
		con, err := InterfaceToType(contain, rt.Elem())
		if err != nil {
			return false, fmt.Errorf("error converting condition to %T: %v", rt.Elem().Kind, err)
		}
		vaany, err := ArrayReflectToArrayOfAny(va)
		if err != nil {
			return false, fmt.Errorf("error converting value to array of any: %v", err)
		}
		return slices.Contains(vaany, con), nil
	case reflect.Map:
		if rv.Len() == 0 {
			return false, nil
		}
		con, err := InterfaceToType(contain, rt.Key())
		if err != nil {
			return false, fmt.Errorf("error converting condition to %T: %v", rt.Elem().Kind, err)
		}
		vaany, err := MapKeysToArrayOfAny(va)
		if err != nil {
			return false, fmt.Errorf("error converting value to array of any: %v", err)
		}
		return slices.Contains(vaany, con), nil
	default:
		return false, fmt.Errorf("type %v not supported for contains", reflect.TypeOf(va))
	}
}

func From(v any, from string, not bool) (bool, error) {
	switch v := v.(type) {
	case string, bool,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		b, err := ConditionValueToArrayOfT(from, reflect.TypeOf(v))
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(b, any(v)), nil
	default:
		switch reflect.TypeOf(v).Kind() {
		case reflect.Array, reflect.Slice:
			f, err := ConditionValueToArrayOfT(from, reflect.TypeOf(v).Elem())
			if err != nil {
				return false, err
			}
			conditionValues := ArrayOfTToArrayOfAny(f)

			values := []any{}
			rv := reflect.ValueOf(v)
			for i := 0; i < rv.Len(); i++ {
				avi := rv.Index(i).Interface()
				values = append(values, avi)
			}

			return FromArray(values, conditionValues, not)
		case reflect.Map:
			f, err := ConditionValueToArrayOfT(from, reflect.TypeOf(v).Key())
			if err != nil {
				return false, err
			}
			conditionValues := ArrayOfTToArrayOfAny(f)

			values := []any{}
			rv := reflect.ValueOf(v)
			for _, mk := range rv.MapKeys() {
				values = append(values, mk.Interface())
			}

			return FromArray(values, conditionValues, not)
		}
		return false, fmt.Errorf("type %v not supported for From validation", reflect.TypeOf(v))
	}
}

func FromArray[T comparable](v []T, from []T, not bool) (bool, error) {
	for _, item := range v {
		if not == slices.Contains(from, item) {
			return false, nil
		}
	}
	return true, nil
}

func Regex(s, regex string) bool {
	matched, _ := regexp.MatchString(regex, s)
	return matched
}
