package validators

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/siherrmann/validator/model"
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

func From(v any, ast *model.AstValue, not bool) (bool, error) {
	switch v := v.(type) {
	case string:
		from, err := ConditionValueToArrayOfString(ast.ConditionValue)
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(from, v), nil
	case int:
		i := int(0)
		from, err := ConditionValueToArrayOfT(i, ast)
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(from, v), nil
	case int8:
		i := int8(0)
		from, err := ConditionValueToArrayOfT(i, ast)
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(from, v), nil
	case int16:
		i := int16(0)
		from, err := ConditionValueToArrayOfT(i, ast)
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(from, v), nil
	case int32:
		i := int32(0)
		from, err := ConditionValueToArrayOfT(i, ast)
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(from, v), nil
	case int64:
		i := int64(0)
		from, err := ConditionValueToArrayOfT(i, ast)
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(from, v), nil
	case uint:
		i := uint(0)
		from, err := ConditionValueToArrayOfT(i, ast)
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(from, v), nil
	case uint8:
		i := uint8(0)
		from, err := ConditionValueToArrayOfT(i, ast)
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(from, v), nil
	case uint16:
		i := uint16(0)
		from, err := ConditionValueToArrayOfT(i, ast)
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(from, v), nil
	case uint32:
		i := uint32(0)
		from, err := ConditionValueToArrayOfT(i, ast)
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(from, v), nil
	case uint64:
		i := uint64(0)
		from, err := ConditionValueToArrayOfT(i, ast)
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(from, v), nil
	case float32:
		i := float32(0)
		from, err := ConditionValueToArrayOfT(i, ast)
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(from, v), nil
	case float64:
		i := float64(0)
		from, err := ConditionValueToArrayOfT(i, ast)
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(from, v), nil
	case bool:
		b, err := ConditionValueToArrayOfT(v, ast)
		if err != nil {
			return false, err
		}
		return !not == slices.Contains(b, v), nil
	case []string:
		from, err := ArrayConditionValueToArrayOfT(v, ast)
		if err != nil {
			return false, err
		}
		return FromArray(v, from, not)
	case []int:
		from, err := ArrayConditionValueToArrayOfT(v, ast)
		if err != nil {
			return false, err
		}
		return FromArray(v, from, not)
	case []int8:
		from, err := ArrayConditionValueToArrayOfT(v, ast)
		if err != nil {
			return false, err
		}
		return FromArray(v, from, not)
	case []int16:
		from, err := ArrayConditionValueToArrayOfT(v, ast)
		if err != nil {
			return false, err
		}
		return FromArray(v, from, not)
	case []int32:
		from, err := ArrayConditionValueToArrayOfT(v, ast)
		if err != nil {
			return false, err
		}
		return FromArray(v, from, not)
	case []int64:
		from, err := ArrayConditionValueToArrayOfT(v, ast)
		if err != nil {
			return false, err
		}
		return FromArray(v, from, not)
	case []uint:
		from, err := ArrayConditionValueToArrayOfT(v, ast)
		if err != nil {
			return false, err
		}
		return FromArray(v, from, not)
	case []uint8:
		from, err := ArrayConditionValueToArrayOfT(v, ast)
		if err != nil {
			return false, err
		}
		return FromArray(v, from, not)
	case []uint16:
		from, err := ArrayConditionValueToArrayOfT(v, ast)
		if err != nil {
			return false, err
		}
		return FromArray(v, from, not)
	case []uint32:
		from, err := ArrayConditionValueToArrayOfT(v, ast)
		if err != nil {
			return false, err
		}
		return FromArray(v, from, not)
	case []uint64:
		from, err := ArrayConditionValueToArrayOfT(v, ast)
		if err != nil {
			return false, err
		}
		return FromArray(v, from, not)
	case []float32:
		from, err := ArrayConditionValueToArrayOfT(v, ast)
		if err != nil {
			return false, err
		}
		return FromArray(v, from, not)
	case []float64:
		from, err := ArrayConditionValueToArrayOfT(v, ast)
		if err != nil {
			return false, err
		}
		return FromArray(v, from, not)
	case []bool:
		from, err := ArrayConditionValueToArrayOfT(v, ast)
		if err != nil {
			return false, err
		}
		return FromArray(v, from, not)
	case map[string]string:
		from, err := ConditionValueToArrayOfString(ast.ConditionValue)
		if err != nil {
			return false, err
		}
		return FromMap(v, from, not)
	case map[string]int:
		from, err := ConditionValueToArrayOfString(ast.ConditionValue)
		if err != nil {
			return false, err
		}
		return FromMap(v, from, not)
	case map[string]int8:
		from, err := ConditionValueToArrayOfString(ast.ConditionValue)
		if err != nil {
			return false, err
		}
		return FromMap(v, from, not)
	case map[string]int16:
		from, err := ConditionValueToArrayOfString(ast.ConditionValue)
		if err != nil {
			return false, err
		}
		return FromMap(v, from, not)
	case map[string]int32:
		from, err := ConditionValueToArrayOfString(ast.ConditionValue)
		if err != nil {
			return false, err
		}
		return FromMap(v, from, not)
	case map[string]int64:
		from, err := ConditionValueToArrayOfString(ast.ConditionValue)
		if err != nil {
			return false, err
		}
		return FromMap(v, from, not)
	case map[string]uint:
		from, err := ConditionValueToArrayOfString(ast.ConditionValue)
		if err != nil {
			return false, err
		}
		return FromMap(v, from, not)
	case map[string]uint8:
		from, err := ConditionValueToArrayOfString(ast.ConditionValue)
		if err != nil {
			return false, err
		}
		return FromMap(v, from, not)
	case map[string]uint16:
		from, err := ConditionValueToArrayOfString(ast.ConditionValue)
		if err != nil {
			return false, err
		}
		return FromMap(v, from, not)
	case map[string]uint32:
		from, err := ConditionValueToArrayOfString(ast.ConditionValue)
		if err != nil {
			return false, err
		}
		return FromMap(v, from, not)
	case map[string]uint64:
		from, err := ConditionValueToArrayOfString(ast.ConditionValue)
		if err != nil {
			return false, err
		}
		return FromMap(v, from, not)
	case map[string]float32:
		from, err := ConditionValueToArrayOfString(ast.ConditionValue)
		if err != nil {
			return false, err
		}
		return FromMap(v, from, not)
	case map[string]float64:
		from, err := ConditionValueToArrayOfString(ast.ConditionValue)
		if err != nil {
			return false, err
		}
		return FromMap(v, from, not)
	case map[string]bool:
		from, err := ConditionValueToArrayOfString(ast.ConditionValue)
		if err != nil {
			return false, err
		}
		return FromMap(v, from, not)
	case map[string]any:
		from, err := ConditionValueToArrayOfString(ast.ConditionValue)
		if err != nil {
			return false, err
		}
		return FromMap(v, from, not)
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
