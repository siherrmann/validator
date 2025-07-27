package helper

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func AnyToFloat(v any) (float64, error) {
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
	default:
		rv := reflect.ValueOf(v)
		switch rv.Type().Kind() {
		case reflect.Array, reflect.Slice, reflect.Map:
			return float64(rv.Len()), nil
		default:
			return 0, fmt.Errorf("unsupported type for value: %T", v)
		}
	}
}

func AnyToString(v any) (string, error) {
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

func AnyToArrayOfString(v any) ([]string, error) {
	if v, ok := v.([]string); ok {
		return v, nil
	}

	rv := reflect.ValueOf(v)
	switch rv.Type().Kind() {
	case reflect.Array, reflect.Slice:
		values := []string{}
		rv := reflect.ValueOf(v)
		for i := 0; i < rv.Len(); i++ {
			avs, err := AnyToString(rv.Index(i).Interface())
			if err != nil {
				return nil, fmt.Errorf("error converting array value to string: %T", v)
			}
			values = append(values, avs)
		}
		return values, nil
	case reflect.Map:
		keys := []string{}
		for _, mk := range reflect.ValueOf(v).MapKeys() {
			mks, err := AnyToString(mk.Interface())
			if err != nil {
				return nil, fmt.Errorf("error converting map key to string: %T", v)
			}
			keys = append(keys, mks)
		}
		return keys, nil
	default:
		return nil, fmt.Errorf("unsupported type for value: %T", v)
	}
}

func AnyToType(in any, expected reflect.Type) (out any, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error converting interface to type %v: %v", expected, r)
		}
	}()

	switch expected.Kind() {
	case reflect.String:
		if v, ok := in.(string); ok {
			return v, nil
		}
	case reflect.Bool:
		if v, ok := in.(bool); ok {
			return v, nil
		} else if v, ok := in.(string); ok {
			switch v {
			// Case on and off are for form values.
			case "on":
				return any(true), nil
			case "off":
				return any(false), nil
			default:
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf("error parsing string to bool: %v", err)
				}
				return any(b), nil
			}
		}
	case reflect.Int:
		if v, ok := in.(int); ok {
			return v, nil
		} else if v, ok := in.(float64); ok {
			return any(int(v)), nil
		} else if v, ok := in.(string); ok {
			i, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing string to int64: %v", err)
			}
			return any(int(i)), nil
		} else if reflect.TypeOf(in).ConvertibleTo(expected) {
			return reflect.ValueOf(in).Convert(expected).Interface(), nil
		}
	case reflect.Int8:
		if v, ok := in.(int8); ok {
			return v, nil
		} else if v, ok := in.(float64); ok {
			return any(int8(v)), nil
		} else if v, ok := in.(string); ok {
			i, err := strconv.ParseInt(v, 10, 8)
			if err != nil {
				return nil, fmt.Errorf("error parsing string to int8: %v", err)
			}
			return any(int8(i)), nil
		} else if reflect.TypeOf(in).ConvertibleTo(expected) {
			return reflect.ValueOf(in).Convert(expected).Interface(), nil
		}
	case reflect.Int16:
		if v, ok := in.(int16); ok {
			return v, nil
		} else if v, ok := in.(float64); ok {
			return any(int16(v)), nil
		} else if v, ok := in.(string); ok {
			i, err := strconv.ParseInt(v, 10, 16)
			if err != nil {
				return nil, fmt.Errorf("error parsing string to int16: %v", err)
			}
			return any(int16(i)), nil
		} else if reflect.TypeOf(in).ConvertibleTo(expected) {
			return reflect.ValueOf(in).Convert(expected).Interface(), nil
		}
	case reflect.Int32:
		if v, ok := in.(int32); ok {
			return v, nil
		} else if v, ok := in.(float64); ok {
			return any(int32(v)), nil
		} else if v, ok := in.(string); ok {
			i, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("error parsing string to int32: %v", err)
			}
			return any(int32(i)), nil
		} else if reflect.TypeOf(in).ConvertibleTo(expected) {
			return reflect.ValueOf(in).Convert(expected).Interface(), nil
		}
	case reflect.Int64:
		if v, ok := in.(int64); ok {
			return v, nil
		} else if v, ok := in.(float64); ok {
			return any(int64(v)), nil
		} else if v, ok := in.(string); ok {
			i, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing string to int64: %v", err)
			}
			return any(int64(i)), nil
		} else if reflect.TypeOf(in).ConvertibleTo(expected) {
			return reflect.ValueOf(in).Convert(expected).Interface(), nil
		}
	case reflect.Uint:
		if v, ok := in.(uint); ok {
			return v, nil
		} else if v, ok := in.(float64); ok {
			return any(uint(v)), nil
		} else if v, ok := in.(string); ok {
			u, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing string to uint64: %v", err)
			}
			return any(uint(u)), nil
		} else if reflect.TypeOf(in).ConvertibleTo(expected) {
			return reflect.ValueOf(in).Convert(expected).Interface(), nil
		}
	case reflect.Uint8:
		if v, ok := in.(uint8); ok {
			return v, nil
		} else if v, ok := in.(float64); ok {
			return any(uint8(v)), nil
		} else if v, ok := in.(string); ok {
			u, err := strconv.ParseUint(v, 10, 8)
			if err != nil {
				return nil, fmt.Errorf("error parsing string to uint8: %v", err)
			}
			return any(uint8(u)), nil
		} else if reflect.TypeOf(in).ConvertibleTo(expected) {
			return reflect.ValueOf(in).Convert(expected).Interface(), nil
		}
	case reflect.Uint16:
		if v, ok := in.(uint16); ok {
			return v, nil
		} else if v, ok := in.(float64); ok {
			return any(uint16(v)), nil
		} else if v, ok := in.(string); ok {
			u, err := strconv.ParseUint(v, 10, 16)
			if err != nil {
				return nil, fmt.Errorf("error parsing string to uint16: %v", err)
			}
			return any(uint16(u)), nil
		} else if reflect.TypeOf(in).ConvertibleTo(expected) {
			return reflect.ValueOf(in).Convert(expected).Interface(), nil
		}
	case reflect.Uint32:
		if v, ok := in.(uint32); ok {
			return v, nil
		} else if v, ok := in.(float64); ok {
			return any(uint32(v)), nil
		} else if v, ok := in.(string); ok {
			u, err := strconv.ParseUint(v, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("error parsing string to uint32: %v", err)
			}
			return any(uint32(u)), nil
		} else if reflect.TypeOf(in).ConvertibleTo(expected) {
			return reflect.ValueOf(in).Convert(expected).Interface(), nil
		}
	case reflect.Uint64:
		if v, ok := in.(uint64); ok {
			return v, nil
		} else if v, ok := in.(float64); ok {
			return any(uint64(v)), nil
		} else if v, ok := in.(string); ok {
			u, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing string to uint64: %v", err)
			}
			return u, nil
		} else if reflect.TypeOf(in).ConvertibleTo(expected) {
			return reflect.ValueOf(in).Convert(expected).Interface(), nil
		}
	case reflect.Float32:
		if v, ok := in.(float32); ok {
			return v, nil
		} else if v, ok := in.(float64); ok {
			return any(float32(v)), nil
		} else if v, ok := in.(string); ok {
			f, err := strconv.ParseFloat(v, 32)
			if err != nil {
				return nil, fmt.Errorf("error parsing string to float32: %v", err)
			}
			return any(float32(f)), nil
		} else if reflect.TypeOf(in).ConvertibleTo(expected) {
			return reflect.ValueOf(in).Convert(expected).Interface(), nil
		}
	case reflect.Float64:
		if v, ok := in.(float64); ok {
			return v, nil
		} else if v, ok := in.(string); ok {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing string to float64: %v", err)
			}
			return any(float64(f)), nil
		} else if reflect.TypeOf(in).ConvertibleTo(expected) {
			return reflect.ValueOf(in).Convert(expected).Interface(), nil
		}
	case reflect.Interface:
		return in, nil
	default:
		return nil, fmt.Errorf("unsupported type %T", expected)
	}
	return nil, fmt.Errorf("unsupported type %T", expected)
}
