package validators

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/siherrmann/validator/helper"
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
	if va == nil {
		return false, nil
	}
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
		if rt.Elem().Kind() == reflect.Interface {
			// E.g. []any from url values: compare against each item in its own type.
			for i := 0; i < rv.Len(); i++ {
				item := rv.Index(i).Interface()
				if item == nil {
					continue
				}
				con, err := helper.AnyToType(contain, reflect.TypeOf(item))
				if err == nil && reflect.DeepEqual(item, con) {
					return true, nil
				}
			}
			return false, nil
		}
		con, err := helper.AnyToType(contain, rt.Elem())
		if err != nil {
			return false, fmt.Errorf("error converting condition to %T: %v", rt.Elem().Kind, err)
		}
		vaany, err := helper.ArrayToArrayOfAny(va)
		if err != nil {
			return false, fmt.Errorf("error converting value to array of any: %v", err)
		}
		return slices.Contains(vaany, con), nil
	case reflect.Map:
		if rv.Len() == 0 {
			return false, nil
		}
		con, err := helper.AnyToType(contain, rt.Key())
		if err != nil {
			return false, fmt.Errorf("error converting condition to %T: %v", rt.Elem().Kind, err)
		}
		vaany, err := helper.MapKeysToArrayOfAny(va)
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
	case nil:
		// A null value is not part of any list
		return not, nil
	case string, bool,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		b, err := helper.ConditionValueToArrayOfAny(from, reflect.TypeOf(v))
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(b, any(v)), nil
	default:
		switch reflect.TypeOf(v).Kind() {
		case reflect.Array, reflect.Slice:
			f, err := helper.ConditionValueToArrayOfAny(from, reflect.TypeOf(v).Elem())
			if err != nil {
				return false, err
			}
			conditionValues, err := helper.ArrayToArrayOfAny(f)
			if err != nil {
				return false, err
			}

			values := []any{}
			rv := reflect.ValueOf(v)
			for i := 0; i < rv.Len(); i++ {
				avi := rv.Index(i).Interface()
				values = append(values, avi)
			}

			return FromArray(values, conditionValues, not)
		case reflect.Map:
			f, err := helper.ConditionValueToArrayOfAny(from, reflect.TypeOf(v).Key())
			if err != nil {
				return false, err
			}
			conditionValues, err := helper.ArrayToArrayOfAny(f)
			if err != nil {
				return false, err
			}

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

// fromValues is From with the values of a referenced field, see ValidateFromValues.
func fromValues(v any, from any, not bool) bool {
	if v == nil {
		// A null value is not part of any list
		return not
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			if containsValue(from, rv.Index(i).Interface()) == not {
				return false
			}
		}
		return true
	case reflect.Map:
		iter := rv.MapRange()
		for iter.Next() {
			if containsValue(from, iter.Key().Interface()) == not {
				return false
			}
		}
		return true
	default:
		return containsValue(from, v) != not
	}
}

// containsValue reports if the item is the value, an item of the list or a key of the map from.
func containsValue(from any, item any) bool {
	if values, ok := from.([]any); ok {
		for _, value := range values {
			if equalValue(item, value) {
				return true
			}
		}
		return false
	}

	rv := reflect.ValueOf(from)
	switch rv.Kind() {
	case reflect.Invalid:
		return false
	case reflect.Slice, reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			if equalValue(item, rv.Index(i).Interface()) {
				return true
			}
		}
		return false
	case reflect.Map:
		iter := rv.MapRange()
		for iter.Next() {
			if equalValue(item, iter.Key().Interface()) {
				return true
			}
		}
		return false
	default:
		return equalValue(item, from)
	}
}

// equalValue compares two values, converting b to the type of a if the types differ.
func equalValue(a any, b any) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	typeA := reflect.TypeOf(a)
	if typeA == reflect.TypeOf(b) {
		if typeA.Comparable() {
			return a == b
		}
		return reflect.DeepEqual(a, b)
	}
	converted, err := helper.AnyToType(b, typeA)
	return err == nil && reflect.DeepEqual(a, converted)
}
