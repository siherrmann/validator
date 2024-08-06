package validator

import (
	"fmt"
	"net/url"
	"reflect"
	"time"

	"github.com/siherrmann/validator/model"
	"github.com/siherrmann/validator/validators"
)

func UnmapJsonMapFromUrlValues(values url.Values, validation ...model.Validation) (model.JsonMap, error) {
	mapOut := model.JsonMap{}
	for _, v := range validation {
		if values.Has(v.Key) {
			value := values.Get(v.Key)
			out, err := model.InterfaceFromString(value, v.Type)
			if err != nil {
				return nil, err
			}
			mapOut[v.Key] = out
		}
	}
	return mapOut, nil
}

func UnmapAndValidateJsonMapFromUrlValues(values url.Values, validation ...model.Validation) (model.JsonMap, error) {
	mapOut := model.JsonMap{}
	for _, v := range validation {
		if values.Has(v.Key) {
			value := values.Get(v.Key)
			out, err := model.InterfaceFromString(value, v.Type)
			if err != nil {
				return nil, err
			}
			err = ValidateInterface(out, v.Type, v.Conditions)
			if err != nil {
				return nil, err
			}
			mapOut[v.Key] = out
		}
	}
	return mapOut, nil
}

func ValidateJsonMap(jsonMap model.JsonMap, validation ...model.Validation) error {
	for _, v := range validation {
		if jsonMap.Has(v.Key) {
			vType, err := model.TypeFromInterface(v)
			if err != nil {
				return err
			}
			err = ValidateInterface(v, vType, v.Conditions)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ValidateInterface(in interface{}, inType model.ValidatorType, condition string) error {
	or := false
	conditions, or := model.GetConditionsAndOrFromString(condition)

	switch inType {
	case model.String:
		inString, ok := in.(string)
		if !ok {
			return fmt.Errorf("invalid string: %v", in)
		}
		err := validators.CheckString(inString, conditions, or)
		if err != nil {
			return err
		}
		return nil
	case model.Int:
		inInt, ok := in.(int)
		if !ok {
			return fmt.Errorf("invalid int: %v", in)
		}
		err := validators.CheckInt(inInt, conditions, or)
		if err != nil {
			return err
		}
		return nil
	case model.Float:
		inFloat, ok := in.(float64)
		if !ok {
			return fmt.Errorf("invalid float: %v", in)
		}
		err := validators.CheckFloat(inFloat, conditions, or)
		if err != nil {
			return err
		}
		return nil
	case model.Bool:
		_, ok := in.(bool)
		if !ok {
			return fmt.Errorf("invalid bool: %v", in)
		}
		return nil
	case model.Map:
		inMap, ok := in.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid map: %v", in)
		}
		err := validators.CheckMap(reflect.ValueOf(inMap), conditions, or)
		if err != nil {
			return err
		}
		return nil
	case model.Time:
		inTime, ok := in.(time.Time)
		if !ok {
			return fmt.Errorf("invalid time: %v", in)
		}
		err := validators.CheckTime(inTime, conditions, or)
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("receiver unsupported: %v", inType)
	}
}
