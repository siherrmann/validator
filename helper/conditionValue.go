package helper

import (
	"fmt"
	"reflect"
	"strings"
)

func ConditionValueToT[T comparable](v T, condition string) (T, error) {
	rt := reflect.TypeOf(v)
	out, err := AnyToType(condition, rt)
	if err != nil {
		return v, fmt.Errorf("error converting condition value: %v", err)
	}
	return out.(T), nil
}

func ConditionValueToArrayOfAny(condition string, expected reflect.Type) ([]any, error) {
	conditionList := strings.Split(condition, ",")
	if len(conditionList) == 0 || (len(conditionList) == 1 && len(strings.TrimSpace(conditionList[0])) == 0) {
		return []any{}, fmt.Errorf("empty condition list %s value", condition)
	}

	values := []any{}
	for _, c := range conditionList {
		ct, err := AnyToType(any(c), expected)
		if err != nil {
			return nil, fmt.Errorf("error converting map key to string: %v", err)
		}
		values = append(values, any(ct))
	}

	return values, nil
}
