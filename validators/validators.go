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

type Comparable interface {
	string | int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64 | bool | any
}

type Array interface {
	[]string | []int | []int8 | []int16 | []int32 | []int64 |
		[]uint | []uint8 | []uint16 | []uint32 | []uint64 |
		[]float32 | []float64 | []bool | []any
}

type Map interface {
	map[string]string | map[string]int | map[string]int8 | map[string]int16 | map[string]int32 | map[string]int64 |
		map[string]uint | map[string]uint8 | map[string]uint16 | map[string]uint32 | map[string]uint64 |
		map[string]float32 | map[string]float64 | map[string]bool | map[string]any
}

func ToArrayOfComparable(v any) ([]Comparable, error) {
	vKind := reflect.TypeOf(v).Kind()
	if vKind != reflect.Slice && vKind != reflect.Array {
		return nil, fmt.Errorf("value has to be of type array or slice, was %v", vKind)
	}

	s := reflect.ValueOf(v)
	vComparable := make([]Comparable, s.Len())
	for i := 0; i < s.Len(); i++ {
		elem := s.Index(i).Interface()
		switch elem.(type) {
		case string, int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64,
			float32, float64, bool:
			vComparable[i] = elem.(Comparable)
		default:
			return nil, fmt.Errorf("element at index %d is of type %T which is not Comparable", i, elem)
		}
	}
	return vComparable, nil
}

func ConditionValueToT[T Comparable](v T, ast *model.AstValue) (T, error) {
	switch any(v).(type) {
	case string:
		return any(ast.ConditionValue).(T), nil
	case int:
		intValue, err := strconv.Atoi(ast.ConditionValue)
		if err != nil {
			return v, fmt.Errorf("invalid condition value for int: %v", ast.ConditionValue)
		}
		return any(intValue).(T), nil
	case int8:
		int8Value, err := strconv.ParseInt(ast.ConditionValue, 10, 8)
		if err != nil {
			return v, fmt.Errorf("invalid condition value for int8: %v", ast.ConditionValue)
		}
		return any(int8(int8Value)).(T), nil
	case int16:
		int16Value, err := strconv.ParseInt(ast.ConditionValue, 10, 16)
		if err != nil {
			return v, fmt.Errorf("invalid condition value for int16: %v", ast.ConditionValue)
		}
		return any(int16(int16Value)).(T), nil
	case int32:
		int32Value, err := strconv.ParseInt(ast.ConditionValue, 10, 32)
		if err != nil {
			return v, fmt.Errorf("invalid condition value for int32: %v", ast.ConditionValue)
		}
		return any(int32(int32Value)).(T), nil
	case int64:
		int64Value, err := strconv.ParseInt(ast.ConditionValue, 10, 64)
		if err != nil {
			return v, fmt.Errorf("invalid condition value for int64: %v", ast.ConditionValue)
		}
		return any(int64Value).(T), nil
	case uint:
		uintValue, err := strconv.ParseUint(ast.ConditionValue, 10, 64)
		if err != nil {
			return v, fmt.Errorf("invalid condition value for uint: %v", ast.ConditionValue)
		}
		return any(uint(uintValue)).(T), nil
	case uint8:
		uint8Value, err := strconv.ParseUint(ast.ConditionValue, 10, 8)
		if err != nil {
			return v, fmt.Errorf("invalid condition value for uint8: %v", ast.ConditionValue)
		}
		return any(uint8(uint8Value)).(T), nil
	case uint16:
		uint16Value, err := strconv.ParseUint(ast.ConditionValue, 10, 16)
		if err != nil {
			return v, fmt.Errorf("invalid condition value for uint16: %v", ast.ConditionValue)
		}
		return any(uint16(uint16Value)).(T), nil
	case uint32:
		uint32Value, err := strconv.ParseUint(ast.ConditionValue, 10, 32)
		if err != nil {
			return v, fmt.Errorf("invalid condition value for uint32: %v", ast.ConditionValue)
		}
		return any(uint32(uint32Value)).(T), nil
	case uint64:
		uint64Value, err := strconv.ParseUint(ast.ConditionValue, 10, 64)
		if err != nil {
			return v, fmt.Errorf("invalid condition value for uint64: %v", ast.ConditionValue)
		}
		return any(uint64Value).(T), nil
	case float32:
		floatValue, err := strconv.ParseFloat(ast.ConditionValue, 32)
		if err != nil {
			return v, fmt.Errorf("invalid condition value for float32: %v", ast.ConditionValue)
		}
		return any(float32(floatValue)).(T), nil
	case float64:
		floatValue, err := strconv.ParseFloat(ast.ConditionValue, 64)
		if err != nil {
			return v, fmt.Errorf("invalid condition value for float: %v", ast.ConditionValue)
		}
		return any(floatValue).(T), nil
	case bool:
		boolValue, err := strconv.ParseBool(ast.ConditionValue)
		if err != nil {
			return v, fmt.Errorf("invalid condition value for bool: %v", ast.ConditionValue)
		}
		return any(boolValue).(T), nil
	default:
		return v, fmt.Errorf("unsupported type for condition value: %T", v)
	}
}

func ArrayConditionValueToT[T []Comparable](v T, ast *model.AstValue) (Comparable, error) {
	switch any(v).(type) {
	case []string:
		return any(ast.ConditionValue).(T), nil
	case []int:
		intValue, err := strconv.Atoi(ast.ConditionValue)
		if err != nil {
			return any(0).(T), fmt.Errorf("invalid condition value for int: %v", ast.ConditionValue)
		}
		return any(intValue).(T), nil
	case []int8:
		int8Value, err := strconv.ParseInt(ast.ConditionValue, 10, 8)
		if err != nil {
			return any(0).(T), fmt.Errorf("invalid condition value for int8: %v", ast.ConditionValue)
		}
		return any(int8(int8Value)).(T), nil
	case []int16:
		int16Value, err := strconv.ParseInt(ast.ConditionValue, 10, 16)
		if err != nil {
			return any(0).(T), fmt.Errorf("invalid condition value for int16: %v", ast.ConditionValue)
		}
		return any(int16(int16Value)).(T), nil
	case []int32:
		int32Value, err := strconv.ParseInt(ast.ConditionValue, 10, 32)
		if err != nil {
			return any(0).(T), fmt.Errorf("invalid condition value for int32: %v", ast.ConditionValue)
		}
		return any(int32(int32Value)).(T), nil
	case []int64:
		int64Value, err := strconv.ParseInt(ast.ConditionValue, 10, 64)
		if err != nil {
			return any(0).(T), fmt.Errorf("invalid condition value for int64: %v", ast.ConditionValue)
		}
		return any(int64Value).(T), nil
	case []uint:
		uintValue, err := strconv.ParseUint(ast.ConditionValue, 10, 64)
		if err != nil {
			return any(0).(T), fmt.Errorf("invalid condition value for uint: %v", ast.ConditionValue)
		}
		return any(uint(uintValue)).(T), nil
	case []uint8:
		uint8Value, err := strconv.ParseUint(ast.ConditionValue, 10, 8)
		if err != nil {
			return any(0).(T), fmt.Errorf("invalid condition value for uint8: %v", ast.ConditionValue)
		}
		return any(uint8(uint8Value)).(T), nil
	case []uint16:
		uint16Value, err := strconv.ParseUint(ast.ConditionValue, 10, 16)
		if err != nil {
			return any(0).(T), fmt.Errorf("invalid condition value for uint16: %v", ast.ConditionValue)
		}
		return any(uint16(uint16Value)).(T), nil
	case []uint32:
		uint32Value, err := strconv.ParseUint(ast.ConditionValue, 10, 32)
		if err != nil {
			return any(0).(T), fmt.Errorf("invalid condition value for uint32: %v", ast.ConditionValue)
		}
		return any(uint32(uint32Value)).(T), nil
	case []uint64:
		uint64Value, err := strconv.ParseUint(ast.ConditionValue, 10, 64)
		if err != nil {
			return any(0).(T), fmt.Errorf("invalid condition value for uint64: %v", ast.ConditionValue)
		}
		return any(uint64Value).(T), nil
	case []float32:
		floatValue, err := strconv.ParseFloat(ast.ConditionValue, 32)
		if err != nil {
			return any(0).(T), fmt.Errorf("invalid condition value for float32: %v", ast.ConditionValue)
		}
		return any(float32(floatValue)).(T), nil
	case []float64:
		floatValue, err := strconv.ParseFloat(ast.ConditionValue, 64)
		if err != nil {
			return any(0).(T), fmt.Errorf("invalid condition value for float: %v", ast.ConditionValue)
		}
		return any(floatValue).(T), nil
	case []bool:
		boolValue, err := strconv.ParseBool(ast.ConditionValue)
		if err != nil {
			return any(false).(T), fmt.Errorf("invalid condition value for bool: %v", ast.ConditionValue)
		}
		return any(boolValue).(T), nil
	default:
		return any(0).(T), fmt.Errorf("unsupported type for condition value: %T", v)
	}
}

type Number interface {
	int | float64
}

func ValueToFloat[V Comparable | Array | Map](v V) (float64, error) {
	switch v := any(v).(type) {
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return float64(v), nil
	case string:
		return float64(len(v)), nil
	case []string:
		return float64(len(v)), nil
	case []int:
		return float64(len(v)), nil
	case []int8:
		return float64(len(v)), nil
	case []int16:
		return float64(len(v)), nil
	case []int32:
		return float64(len(v)), nil
	case []int64:
		return float64(len(v)), nil
	case []uint:
		return float64(len(v)), nil
	case []uint8:
		return float64(len(v)), nil
	case []uint16:
		return float64(len(v)), nil
	case []uint32:
		return float64(len(v)), nil
	case []uint64:
		return float64(len(v)), nil
	case []float32:
		return float64(len(v)), nil
	case []float64:
		return float64(len(v)), nil
	case []bool:
		return float64(len(v)), nil
	case map[string]string:
		return float64(len(v)), nil
	case map[string]int:
		return float64(len(v)), nil
	case map[string]int8:
		return float64(len(v)), nil
	case map[string]int16:
		return float64(len(v)), nil
	case map[string]int32:
		return float64(len(v)), nil
	case map[string]int64:
		return float64(len(v)), nil
	case map[string]uint:
		return float64(len(v)), nil
	case map[string]uint8:
		return float64(len(v)), nil
	case map[string]uint16:
		return float64(len(v)), nil
	case map[string]uint32:
		return float64(len(v)), nil
	case map[string]uint64:
		return float64(len(v)), nil
	case map[string]float32:
		return float64(len(v)), nil
	case map[string]float64:
		return float64(len(v)), nil
	case map[string]bool:
		return float64(len(v)), nil
	case map[string]any:
		return float64(len(v)), nil
	default:
		return 0, fmt.Errorf("unsupported type for value: %T", v)
	}
}

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

func StringContains(s, contain string) bool {
	return strings.Contains(s, contain)
}

func ArrayContains[T comparable](v []T, contain T) bool {
	return slices.Contains(v, contain)
}

func MapContains[T map[string]Comparable](m T, key string) bool {
	_, exists := m[key]
	return exists
}

func StringNotContains(s, notContain string) bool {
	return !strings.Contains(s, notContain)
}

func ArrayNotContains[T comparable](v T, notContain []T) bool {
	return !slices.Contains(notContain, v)
}

func From[T comparable](v T, from []T) bool {
	return slices.Contains(from, v)
}

func ArrayItemsFrom[T comparable](v []T, from []T) bool {
	for _, vItem := range v {
		if !slices.Contains(from, vItem) {
			return false
		}
	}
	return true
}

func NotFrom[T comparable](v T, notFrom []T) bool {
	return !slices.Contains(notFrom, v)
}

func ArrayItemsNotFrom[T comparable](v []T, notFrom []T) bool {
	for _, vItem := range v {
		if slices.Contains(notFrom, vItem) {
			return false
		}
	}
	return true
}

func RegxMatch(s, regex string) bool {
	matched, _ := regexp.MatchString(regex, s)
	return matched
}
