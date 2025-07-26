package validators

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/siherrmann/validator/model"
)

func ConditionValueToT[T comparable](v T, ast *model.AstValue) (T, error) {
	switch any(v).(type) {
	case string:
		return any(ast.ConditionValue).(T), nil
	case bool:
		boolValue, err := strconv.ParseBool(ast.ConditionValue)
		if err != nil {
			return v, fmt.Errorf("invalid condition value for bool: %v", ast.ConditionValue)
		}
		return any(boolValue).(T), nil
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
	default:
		return v, fmt.Errorf("unsupported type for condition value: %T", v)
	}
}

func ConditionValueToArrayOfT[T comparable](v T, condition string) ([]T, error) {
	conditionList := strings.Split(condition, ",")
	if len(conditionList) == 0 || (len(conditionList) == 1 && len(strings.TrimSpace(conditionList[0])) == 0) {
		return []T{}, fmt.Errorf("empty condition list")
	}

	switch any(v).(type) {
	case string:
		return any(conditionList).([]T), nil
	case bool:
		values, err := ArrayOfStringToArrayOfT(v, conditionList)
		if err != nil {
			return []T{}, fmt.Errorf("invalid condition value for bool: %v", condition)
		}
		return any(values).([]T), nil
	case time.Time:
		var date time.Time
		var err error
		date, err = model.UnixStringToTime(condition)
		if err != nil {
			date, err = model.ISO8601StringToTime(condition)
			if err != nil {
				return nil, err
			}
		}
		return any(date).([]T), nil
	case int:
		values, err := ArrayOfStringToArrayOfT(v, conditionList)
		if err != nil {
			return []T{}, fmt.Errorf("invalid condition value for int: %v", condition)
		}
		return any(values).([]T), nil
	case int8:
		values, err := ArrayOfStringToArrayOfT(v, conditionList)
		if err != nil {
			return []T{}, fmt.Errorf("invalid condition value for int8: %v", condition)
		}
		return any(values).([]T), nil
	case int16:
		values, err := ArrayOfStringToArrayOfT(v, conditionList)
		if err != nil {
			return []T{}, fmt.Errorf("invalid condition value for int16: %v", condition)
		}
		return any(values).([]T), nil
	case int32:
		values, err := ArrayOfStringToArrayOfT(v, conditionList)
		if err != nil {
			return []T{}, fmt.Errorf("invalid condition value for int32: %v", condition)
		}
		return any(values).([]T), nil
	case int64:
		values, err := ArrayOfStringToArrayOfT(v, conditionList)
		if err != nil {
			return []T{}, fmt.Errorf("invalid condition value for int64: %v", condition)
		}
		return any(values).([]T), nil
	case uint:
		values, err := ArrayOfStringToArrayOfT(v, conditionList)
		if err != nil {
			return []T{}, fmt.Errorf("invalid condition value for uint: %v", condition)
		}
		return any(values).([]T), nil
	case uint8:
		values, err := ArrayOfStringToArrayOfT(v, conditionList)
		if err != nil {
			return []T{}, fmt.Errorf("invalid condition value for uint8: %v", condition)
		}
		return any(values).([]T), nil
	case uint16:
		values, err := ArrayOfStringToArrayOfT(v, conditionList)
		if err != nil {
			return []T{}, fmt.Errorf("invalid condition value for uint16: %v", condition)
		}
		return any(values).([]T), nil
	case uint32:
		values, err := ArrayOfStringToArrayOfT(v, conditionList)
		if err != nil {
			return []T{}, fmt.Errorf("invalid condition value for uint32: %v", condition)
		}
		return any(values).([]T), nil
	case uint64:
		values, err := ArrayOfStringToArrayOfT(v, conditionList)
		if err != nil {
			return []T{}, fmt.Errorf("invalid condition value for uint64: %v", condition)
		}
		return any(values).([]T), nil
	case float32:
		values, err := ArrayOfStringToArrayOfT(v, conditionList)
		if err != nil {
			return []T{}, fmt.Errorf("invalid condition value for float32: %v", condition)
		}
		return any(values).([]T), nil
	case float64:
		values, err := ArrayOfStringToArrayOfT(v, conditionList)
		if err != nil {
			return []T{}, fmt.Errorf("invalid condition value for float64: %v", condition)
		}
		return any(values).([]T), nil
	default:
		return []T{}, fmt.Errorf("unsupported type for condition value: %T", v)
	}
}

func ArrayConditionValueToArrayOfT[T comparable](v []T, condition string) ([]T, error) {
	conditionList := strings.Split(condition, ",")
	if len(conditionList) == 0 || (len(conditionList) == 1 && len(strings.TrimSpace(conditionList[0])) == 0) {
		return []T{}, fmt.Errorf("empty condition list")
	}

	switch any(v).(type) {
	case []string:
		return any(conditionList).([]T), nil
	case []bool:
		b := false
		values, err := ArrayOfStringToArrayOfT(b, conditionList)
		if err != nil {
			return v, fmt.Errorf("invalid condition value for bool: %v", condition)
		}
		return any(values).([]T), nil
	case []int:
		i := int(0)
		values, err := ArrayOfStringToArrayOfT(i, conditionList)
		if err != nil {
			return v, fmt.Errorf("invalid condition value for int: %v", condition)
		}
		return any(values).([]T), nil
	case []int8:
		i := int8(0)
		values, err := ArrayOfStringToArrayOfT(i, conditionList)
		if err != nil {
			return v, fmt.Errorf("invalid condition value for int8: %v", condition)
		}
		return any(values).([]T), nil
	case []int16:
		i := int16(0)
		values, err := ArrayOfStringToArrayOfT(i, conditionList)
		if err != nil {
			return v, fmt.Errorf("invalid condition value for int16: %v", condition)
		}
		return any(values).([]T), nil
	case []int32:
		i := int32(0)
		values, err := ArrayOfStringToArrayOfT(i, conditionList)
		if err != nil {
			return v, fmt.Errorf("invalid condition value for int32: %v", condition)
		}
		return any(values).([]T), nil
	case []int64:
		i := int64(0)
		values, err := ArrayOfStringToArrayOfT(i, conditionList)
		if err != nil {
			return v, fmt.Errorf("invalid condition value for int64: %v", condition)
		}
		return any(values).([]T), nil
	case []uint:
		u := uint(0)
		values, err := ArrayOfStringToArrayOfT(u, conditionList)
		if err != nil {
			return v, fmt.Errorf("invalid condition value for uint: %v", condition)
		}
		return any(values).([]T), nil
	case []uint8:
		u := uint8(0)
		values, err := ArrayOfStringToArrayOfT(u, conditionList)
		if err != nil {
			return v, fmt.Errorf("invalid condition value for uint8: %v", condition)
		}
		return any(values).([]T), nil
	case []uint16:
		u := uint16(0)
		values, err := ArrayOfStringToArrayOfT(u, conditionList)
		if err != nil {
			return v, fmt.Errorf("invalid condition value for uint16: %v", condition)
		}
		return any(values).([]T), nil
	case []uint32:
		u := uint32(0)
		values, err := ArrayOfStringToArrayOfT(u, conditionList)
		if err != nil {
			return v, fmt.Errorf("invalid condition value for uint32: %v", condition)
		}
		return any(values).([]T), nil
	case []uint64:
		u := uint64(0)
		values, err := ArrayOfStringToArrayOfT(u, conditionList)
		if err != nil {
			return v, fmt.Errorf("invalid condition value for uint64: %v", condition)
		}
		return any(values).([]T), nil
	case []float32:
		f := float32(0)
		values, err := ArrayOfStringToArrayOfT(f, conditionList)
		if err != nil {
			return v, fmt.Errorf("invalid condition value for float32: %v", condition)
		}
		return any(values).([]T), nil
	case []float64:
		f := float64(0)
		values, err := ArrayOfStringToArrayOfT(f, conditionList)
		if err != nil {
			return v, fmt.Errorf("invalid condition value for float64: %v", condition)
		}
		return any(values).([]T), nil
	default:
		return v, fmt.Errorf("unsupported type for condition value: %T", v)
	}
}

func ArrayOfStringToArrayOfT[T comparable](v T, sa []string) ([]T, error) {
	array := make([]T, len(sa))
	switch any(v).(type) {
	case string:
		for i, item := range sa {
			array[i] = any(item).(T)
		}
	case bool:
		for i, item := range sa {
			boolValue, err := strconv.ParseBool(item)
			if err != nil {
				return nil, fmt.Errorf("invalid condition value for bool: %v", item)
			}
			array[i] = any(boolValue).(T)
		}
	case int:
		for i, item := range sa {
			intValue, err := strconv.Atoi(item)
			if err != nil {
				return nil, fmt.Errorf("invalid condition value for int: %v", item)
			}
			array[i] = any(intValue).(T)
		}
	case int8:
		for i, item := range sa {
			int8Value, err := strconv.ParseInt(item, 10, 8)
			if err != nil {
				return nil, fmt.Errorf("invalid condition value for int8: %v", item)
			}
			array[i] = any(int8(int8Value)).(T)
		}
	case int16:
		for i, item := range sa {
			int16Value, err := strconv.ParseInt(item, 10, 16)
			if err != nil {
				return nil, fmt.Errorf("invalid condition value for int16: %v", item)
			}
			array[i] = any(int16(int16Value)).(T)
		}
	case int32:
		for i, item := range sa {
			int32Value, err := strconv.ParseInt(item, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid condition value for int32: %v", item)
			}
			array[i] = any(int32(int32Value)).(T)
		}
	case int64:
		for i, item := range sa {
			int64Value, err := strconv.ParseInt(item, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid condition value for int64: %v", item)
			}
			array[i] = any(int64Value).(T)
		}
	case uint:
		for i, item := range sa {
			uintValue, err := strconv.ParseUint(item, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid condition value for uint: %v", item)
			}
			array[i] = any(uint(uintValue)).(T)
		}
	case uint8:
		for i, item := range sa {
			uint8Value, err := strconv.ParseUint(item, 10, 8)
			if err != nil {
				return nil, fmt.Errorf("invalid condition value for uint8: %v", item)
			}
			array[i] = any(uint8(uint8Value)).(T)
		}
	case uint16:
		for i, item := range sa {
			uint16Value, err := strconv.ParseUint(item, 10, 16)
			if err != nil {
				return nil, fmt.Errorf("invalid condition value for uint16: %v", item)
			}
			array[i] = any(uint16(uint16Value)).(T)
		}
	case uint32:
		for i, item := range sa {
			uint32Value, err := strconv.ParseUint(item, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid condition value for uint32: %v", item)
			}
			array[i] = any(uint32(uint32Value)).(T)
		}
	case uint64:
		for i, item := range sa {
			uint64Value, err := strconv.ParseUint(item, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid condition value for uint64: %v", item)
			}
			array[i] = any(uint64Value).(T)
		}
	case float32:
		for i, item := range sa {
			floatValue, err := strconv.ParseFloat(item, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid condition value for float32: %v", item)
			}
			array[i] = any(float32(floatValue)).(T)
		}
	case float64:
		for i, item := range sa {
			floatValue, err := strconv.ParseFloat(item, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid condition value for float64: %v", item)
			}
			array[i] = any(floatValue).(T)
		}
	default:
		return nil, fmt.Errorf("unsupported type for condition value: %T", sa)
	}
	return array, nil
}

func ConditionValueToArrayOfString(condition string) ([]string, error) {
	conditionList := strings.Split(condition, ",")
	if len(conditionList) == 0 || (len(conditionList) == 1 && len(strings.TrimSpace(conditionList[0])) == 0) {
		return []string{}, fmt.Errorf("empty condition list %s value", condition)
	}
	return conditionList, nil
}

func ValueToString(v any) (string, error) {
	switch v := v.(type) {
	case string:
		return v, nil
	case bool:
		return fmt.Sprintf("%t", v), nil
	case time.Time:
		return fmt.Sprintf("%d", v.Unix()), nil
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v), nil
	case float32, float64:
		return fmt.Sprintf("%f", v), nil
	default:
		return "", fmt.Errorf("unsupported type for value: %T", v)
	}
}

func ValueToArrayOfString(v any) ([]string, error) {
	switch v := v.(type) {
	case []string:
		return v, nil
	case []bool:
		sa, err := ArrayToArrayOfString(v)
		return sa, err
	case []int:
		sa, err := ArrayToArrayOfString(v)
		return sa, err
	case []int8:
		sa, err := ArrayToArrayOfString(v)
		return sa, err
	case []int16:
		sa, err := ArrayToArrayOfString(v)
		return sa, err
	case []int32:
		sa, err := ArrayToArrayOfString(v)
		return sa, err
	case []int64:
		sa, err := ArrayToArrayOfString(v)
		return sa, err
	case []uint:
		sa, err := ArrayToArrayOfString(v)
		return sa, err
	case []uint8:
		sa, err := ArrayToArrayOfString(v)
		return sa, err
	case []uint16:
		sa, err := ArrayToArrayOfString(v)
		return sa, err
	case []uint32:
		sa, err := ArrayToArrayOfString(v)
		return sa, err
	case []uint64:
		sa, err := ArrayToArrayOfString(v)
		return sa, err
	case []float32:
		sa, err := ArrayToArrayOfString(v)
		return sa, err
	case []float64:
		sa, err := ArrayToArrayOfString(v)
		return sa, err
	case map[string]string:
		sa, err := MapToArrayOfString(v)
		return sa, err
	case map[string]bool:
		sa, err := MapToArrayOfString(v)
		return sa, err
	case map[string]int:
		sa, err := MapToArrayOfString(v)
		return sa, err
	case map[string]int8:
		sa, err := MapToArrayOfString(v)
		return sa, err
	case map[string]int16:
		sa, err := MapToArrayOfString(v)
		return sa, err
	case map[string]int32:
		sa, err := MapToArrayOfString(v)
		return sa, err
	case map[string]int64:
		sa, err := MapToArrayOfString(v)
		return sa, err
	case map[string]uint:
		sa, err := MapToArrayOfString(v)
		return sa, err
	case map[string]uint8:
		sa, err := MapToArrayOfString(v)
		return sa, err
	case map[string]uint16:
		sa, err := MapToArrayOfString(v)
		return sa, err
	case map[string]uint32:
		sa, err := MapToArrayOfString(v)
		return sa, err
	case map[string]uint64:
		sa, err := MapToArrayOfString(v)
		return sa, err
	case map[string]float32:
		sa, err := MapToArrayOfString(v)
		return sa, err
	case map[string]float64:
		sa, err := MapToArrayOfString(v)
		return sa, err
	case map[string]any:
		sa, err := MapToArrayOfString(v)
		return sa, err
	default:
		return nil, fmt.Errorf("unsupported type for value: %T", v)
	}
}

func ArrayToArrayOfString[T comparable](v []T) ([]string, error) {
	array := make([]string, len(v))
	for i, item := range v {
		str, err := ValueToString(item)
		if err != nil {
			return nil, fmt.Errorf("error converting value to string: %v", err)
		}
		array[i] = str
	}
	return array, nil
}

func MapToArrayOfString[T comparable](m map[string]T) ([]string, error) {
	array := make([]string, 0, len(m))
	for key := range m {
		str, err := ValueToString(key)
		if err != nil {
			return nil, fmt.Errorf("error converting map value to string: %v", err)
		}
		array = append(array, str)
	}
	return array, nil
}

type Number interface {
	int | float64
}

func ValueToFloat(v any) (float64, error) {
	switch v := v.(type) {
	case string:
		return float64(len(v)), nil
	case time.Time:
		return float64(v.Unix()), nil
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
		return v, nil
	case []string:
		return float64(len(v)), nil
	case []bool:
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
	case []any:
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
