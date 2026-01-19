package helper

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/google/uuid"
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

	// Nil input
	if in == nil {
		// For maps and slices, return empty collections instead of nil to prevent panics
		if expected.Kind() == reflect.Map {
			return reflect.MakeMap(expected).Interface(), nil
		}
		if expected.Kind() == reflect.Slice {
			return reflect.MakeSlice(expected, 0, 0).Interface(), nil
		}
		return reflect.Zero(expected).Interface(), nil
	}

	// Direct type match or assignable
	inType := reflect.TypeOf(in)
	if inType == expected {
		return in, nil
	}
	if inType.AssignableTo(expected) {
		return in, nil
	}

	// Handle pointer types by recursively converting to the element type,
	// then returning a pointer to the result
	if expected.Kind() == reflect.Ptr {
		elemType := expected.Elem()
		elemValue, err := AnyToType(in, elemType)
		if err != nil {
			return nil, err
		}

		// Create a pointer to the converted value
		ptrValue := reflect.New(elemType)
		ptrValue.Elem().Set(reflect.ValueOf(elemValue))
		return ptrValue.Interface(), nil
	}

	switch expKind := expected.Kind(); expKind {
	case reflect.String:
		if v, ok := in.(string); ok {
			// Check if we need to convert to a custom string type
			if reflect.TypeOf(in).ConvertibleTo(expected) {
				return reflect.ValueOf(in).Convert(expected).Interface(), nil
			}
			return v, nil
		}
	case reflect.Bool:
		if v, ok := in.(bool); ok {
			// Check if we need to convert to a custom bool type
			if reflect.TypeOf(in).ConvertibleTo(expected) {
				return reflect.ValueOf(in).Convert(expected).Interface(), nil
			}
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
	case reflect.Struct:
		if v, ok := in.(time.Time); ok {
			return v, nil
		} else if v, ok := in.(float64); ok {
			date, err := UnixStringToTime(strconv.Itoa(int(v)))
			if err != nil {
				return nil, fmt.Errorf("error parsing float to time.Time: %v", err)
			}
			return date, nil
		} else if v, ok := in.(string); ok {
			structTempt := reflect.New(expected).Interface()
			if err := json.Unmarshal([]byte(v), structTempt); err != nil {
				date, err := UnixStringToTime(v)
				if err != nil {
					date, err = ISO8601StringToTime(v)
					if err != nil {
						return nil, fmt.Errorf("error parsing string to struct: %v", err)
					}
				}
				return date, nil
			}
			return reflect.ValueOf(structTempt).Elem().Interface(), nil
		} else if jsonValueMap, ok := in.(map[string]any); ok {
			structTemp := reflect.New(expected).Interface()
			err = MapJsonMapToStruct(jsonValueMap, structTemp)
			if err != nil {
				return nil, fmt.Errorf("error parsing map to struct: %v", err)
			}
			return reflect.ValueOf(structTemp).Elem().Interface(), nil
		} else if reflect.TypeOf(in).ConvertibleTo(expected) {
			return reflect.ValueOf(in).Convert(expected).Interface(), nil
		}
	case reflect.Map:
		if v, ok := in.(map[string]any); ok {
			mapValue, err := JsonMapToMapKV(v, expected.Key(), expected.Elem())
			if err != nil {
				return nil, err
			}
			return mapValue.Interface(), nil
		} else if reflect.TypeOf(in).ConvertibleTo(expected) {
			return reflect.ValueOf(in).Convert(expected).Interface(), nil
		}
	case reflect.Array, reflect.Slice:
		if inArray, ok := in.([]any); ok {
			typedArray, err := ArrayToArrayOfType(inArray, expected.Elem())
			if err != nil {
				return nil, err
			}
			return typedArray.Interface(), nil
		} else if strValue, ok := in.(string); ok {
			// Special case: if expected is []byte or [N]byte, convert string to byte array/slice
			if expected.Elem().Kind() == reflect.Uint8 {
				if expected.Kind() == reflect.Slice {
					return []byte(strValue), nil
				} else if expected.Kind() == reflect.Array {
					// First try to parse as UUID as most common use case for [16]byte
					parsedUUID, err := uuid.Parse(strValue)
					if err == nil {
						return parsedUUID, nil
					}

					byteSlice := []byte(strValue)
					arrayValue := reflect.New(expected).Elem()
					for i := 0; i < expected.Len() && i < len(byteSlice); i++ {
						arrayValue.Index(i).SetUint(uint64(byteSlice[i]))
					}
					return arrayValue.Interface(), nil
				}
			}
		}

		// Handle array/slice of other types
		if reflect.TypeOf(in).ConvertibleTo(expected) {
			return reflect.ValueOf(in).Convert(expected).Interface(), nil
		}
	case reflect.Interface:
		return in, nil
	default:
		return nil, fmt.Errorf("unsupported type %T", expected)
	}
	return nil, fmt.Errorf("unsupported type %T", expected)
}
