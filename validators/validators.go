package validators

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
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
	switch v := va.(type) {
	case string:
		return strings.Contains(v, contain), nil
	case []bool:
		c, err := strconv.ParseBool(contain)
		if err != nil {
			return false, err
		}
		return slices.Contains(v, c), nil
	case []int:
		c, err := strconv.Atoi(contain)
		if err != nil {
			return false, err
		}
		return slices.Contains(v, c), nil
	case []int8:
		c, err := strconv.ParseInt(contain, 10, 8)
		if err != nil {
			return false, err
		}
		return slices.Contains(v, int8(c)), nil
	case []int16:
		c, err := strconv.ParseInt(contain, 10, 16)
		if err != nil {
			return false, err
		}
		return slices.Contains(v, int16(c)), nil
	case []int32:
		c, err := strconv.ParseInt(contain, 10, 32)
		if err != nil {
			return false, err
		}
		return slices.Contains(v, int32(c)), nil
	case []int64:
		c, err := strconv.ParseInt(contain, 10, 64)
		if err != nil {
			return false, err
		}
		return slices.Contains(v, int64(c)), nil
	case []uint:
		c, err := strconv.ParseUint(contain, 10, 64)
		if err != nil {
			return false, err
		}
		return slices.Contains(v, uint(c)), nil
	case []uint8:
		c, err := strconv.ParseUint(contain, 10, 8)
		if err != nil {
			return false, err
		}
		return slices.Contains(v, uint8(c)), nil
	case []uint16:
		c, err := strconv.ParseUint(contain, 10, 16)
		if err != nil {
			return false, err
		}
		return slices.Contains(v, uint16(c)), nil
	case []uint32:
		c, err := strconv.ParseUint(contain, 10, 32)
		if err != nil {
			return false, err
		}
		return slices.Contains(v, uint32(c)), nil
	case []uint64:
		c, err := strconv.ParseUint(contain, 10, 64)
		if err != nil {
			return false, err
		}
		return slices.Contains(v, uint64(c)), nil
	case []float32:
		c, err := strconv.ParseFloat(contain, 32)
		if err != nil {
			return false, err
		}
		return slices.Contains(v, float32(c)), nil
	case []float64:
		c, err := strconv.ParseFloat(contain, 64)
		if err != nil {
			return false, err
		}
		return slices.Contains(v, c), nil
	case []string:
		return slices.Contains(v, contain), nil
	case map[string]string:
		_, exists := v[contain]
		return exists, nil
	case map[string]bool:
		_, exists := v[contain]
		return exists, nil
	case map[string]int:
		_, exists := v[contain]
		return exists, nil
	case map[string]int8:
		_, exists := v[contain]
		return exists, nil
	case map[string]int16:
		_, exists := v[contain]
		return exists, nil
	case map[string]int32:
		_, exists := v[contain]
		return exists, nil
	case map[string]int64:
		_, exists := v[contain]
		return exists, nil
	case map[string]uint:
		_, exists := v[contain]
		return exists, nil
	case map[string]uint8:
		_, exists := v[contain]
		return exists, nil
	case map[string]uint16:
		_, exists := v[contain]
		return exists, nil
	case map[string]uint32:
		_, exists := v[contain]
		return exists, nil
	case map[string]uint64:
		_, exists := v[contain]
		return exists, nil
	case map[string]float32:
		_, exists := v[contain]
		return exists, nil
	case map[string]float64:
		_, exists := v[contain]
		return exists, nil
	case map[string]any:
		_, exists := v[contain]
		return exists, nil
	default:
		return false, fmt.Errorf("type %v not supported", reflect.TypeOf(v))
	}
}

func From(v any, from string, not bool) (bool, error) {
	switch v := v.(type) {
	case string:
		f, err := ConditionValueToArrayOfString(from)
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(f, v), nil
	case bool:
		b, err := ConditionValueToArrayOfT(v, from)
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(b, v), nil
	case int:
		f, err := ConditionValueToArrayOfT(v, from)
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(f, v), nil
	case int8:
		f, err := ConditionValueToArrayOfT(v, from)
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(f, v), nil
	case int16:
		f, err := ConditionValueToArrayOfT(v, from)
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(f, v), nil
	case int32:
		f, err := ConditionValueToArrayOfT(v, from)
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(f, v), nil
	case int64:
		f, err := ConditionValueToArrayOfT(v, from)
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(f, v), nil
	case uint:
		f, err := ConditionValueToArrayOfT(v, from)
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(f, v), nil
	case uint8:
		f, err := ConditionValueToArrayOfT(v, from)
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(f, v), nil
	case uint16:
		f, err := ConditionValueToArrayOfT(v, from)
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(f, v), nil
	case uint32:
		f, err := ConditionValueToArrayOfT(v, from)
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(f, v), nil
	case uint64:
		f, err := ConditionValueToArrayOfT(v, from)
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(f, v), nil
	case float32:
		f, err := ConditionValueToArrayOfT(v, from)
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(f, v), nil
	case float64:
		f, err := ConditionValueToArrayOfT(v, from)
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(f, v), nil
	case []string:
		f, err := ArrayConditionValueToArrayOfT(v, from)
		if err != nil {
			return false, err
		}
		return FromArray(v, f, not)
	case []int:
		f, err := ArrayConditionValueToArrayOfT(v, from)
		if err != nil {
			return false, err
		}
		return FromArray(v, f, not)
	case []int8:
		f, err := ArrayConditionValueToArrayOfT(v, from)
		if err != nil {
			return false, err
		}
		return FromArray(v, f, not)
	case []int16:
		f, err := ArrayConditionValueToArrayOfT(v, from)
		if err != nil {
			return false, err
		}
		return FromArray(v, f, not)
	case []int32:
		f, err := ArrayConditionValueToArrayOfT(v, from)
		if err != nil {
			return false, err
		}
		return FromArray(v, f, not)
	case []int64:
		f, err := ArrayConditionValueToArrayOfT(v, from)
		if err != nil {
			return false, err
		}
		return FromArray(v, f, not)
	case []uint:
		f, err := ArrayConditionValueToArrayOfT(v, from)
		if err != nil {
			return false, err
		}
		return FromArray(v, f, not)
	case []uint8:
		f, err := ArrayConditionValueToArrayOfT(v, from)
		if err != nil {
			return false, err
		}
		return FromArray(v, f, not)
	case []uint16:
		f, err := ArrayConditionValueToArrayOfT(v, from)
		if err != nil {
			return false, err
		}
		return FromArray(v, f, not)
	case []uint32:
		f, err := ArrayConditionValueToArrayOfT(v, from)
		if err != nil {
			return false, err
		}
		return FromArray(v, f, not)
	case []uint64:
		f, err := ArrayConditionValueToArrayOfT(v, from)
		if err != nil {
			return false, err
		}
		return FromArray(v, f, not)
	case []float32:
		f, err := ArrayConditionValueToArrayOfT(v, from)
		if err != nil {
			return false, err
		}
		return FromArray(v, f, not)
	case []float64:
		f, err := ArrayConditionValueToArrayOfT(v, from)
		if err != nil {
			return false, err
		}
		return FromArray(v, f, not)
	case []bool:
		f, err := ArrayConditionValueToArrayOfT(v, from)
		if err != nil {
			return false, err
		}
		return FromArray(v, f, not)
	case map[string]string:
		f, err := ConditionValueToArrayOfString(from)
		if err != nil {
			return false, err
		}
		return FromMap(v, f, not)
	case map[string]int:
		f, err := ConditionValueToArrayOfString(from)
		if err != nil {
			return false, err
		}
		return FromMap(v, f, not)
	case map[string]int8:
		f, err := ConditionValueToArrayOfString(from)
		if err != nil {
			return false, err
		}
		return FromMap(v, f, not)
	case map[string]int16:
		f, err := ConditionValueToArrayOfString(from)
		if err != nil {
			return false, err
		}
		return FromMap(v, f, not)
	case map[string]int32:
		f, err := ConditionValueToArrayOfString(from)
		if err != nil {
			return false, err
		}
		return FromMap(v, f, not)
	case map[string]int64:
		f, err := ConditionValueToArrayOfString(from)
		if err != nil {
			return false, err
		}
		return FromMap(v, f, not)
	case map[string]uint:
		f, err := ConditionValueToArrayOfString(from)
		if err != nil {
			return false, err
		}
		return FromMap(v, f, not)
	case map[string]uint8:
		f, err := ConditionValueToArrayOfString(from)
		if err != nil {
			return false, err
		}
		return FromMap(v, f, not)
	case map[string]uint16:
		f, err := ConditionValueToArrayOfString(from)
		if err != nil {
			return false, err
		}
		return FromMap(v, f, not)
	case map[string]uint32:
		f, err := ConditionValueToArrayOfString(from)
		if err != nil {
			return false, err
		}
		return FromMap(v, f, not)
	case map[string]uint64:
		f, err := ConditionValueToArrayOfString(from)
		if err != nil {
			return false, err
		}
		return FromMap(v, f, not)
	case map[string]float32:
		f, err := ConditionValueToArrayOfString(from)
		if err != nil {
			return false, err
		}
		return FromMap(v, f, not)
	case map[string]float64:
		f, err := ConditionValueToArrayOfString(from)
		if err != nil {
			return false, err
		}
		return FromMap(v, f, not)
	case map[string]bool:
		f, err := ConditionValueToArrayOfString(from)
		if err != nil {
			return false, err
		}
		return FromMap(v, f, not)
	case map[string]any:
		f, err := ConditionValueToArrayOfString(from)
		if err != nil {
			return false, err
		}
		return FromMap(v, f, not)
	default:
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

func FromMap[T comparable, V comparable](v map[T]V, from []T, not bool) (bool, error) {
	for key := range v {
		if not == slices.Contains(from, key) {
			return false, nil
		}
	}
	return true, nil
}

func Regex(s, regex string) bool {
	matched, _ := regexp.MatchString(regex, s)
	return matched
}
