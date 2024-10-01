package validator

import (
	"fmt"
	"net/url"
	"reflect"

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
	switch inType {
	case model.String:
		err := ValidateValueWithParser(reflect.ValueOf(in), condition, validators.CheckString)
		if err != nil {
			return err
		}
		return nil
	case model.Int:
		err := ValidateValueWithParser(reflect.ValueOf(in), condition, validators.CheckInt)
		if err != nil {
			return err
		}
		return nil
	case model.Float:
		err := ValidateValueWithParser(reflect.ValueOf(in), condition, validators.CheckFloat)
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
		err := ValidateValueWithParser(reflect.ValueOf(in), condition, validators.CheckMap)
		if err != nil {
			return err
		}
		return nil
	case model.Time:
		err := ValidateValueWithParser(reflect.ValueOf(in), condition, validators.CheckTime)
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("receiver unsupported: %v", inType)
	}
}
