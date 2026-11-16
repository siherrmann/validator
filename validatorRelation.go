package validator

import (
	"fmt"
	"reflect"

	"github.com/siherrmann/validator/helper"
	"github.com/siherrmann/validator/model"
	"github.com/siherrmann/validator/validators"
)

// This file holds the validation of `rel` conditions and `$` references,
// which need the values of other fields of the validated input.

// scope is one level of the validated input, used to resolve `rel` conditions and `$` references.
// A level (struct or item of a list of structs) gets its own scope pointing to the scope of the level above,
// but only if a validation on this level or below uses other fields (see validateWithScope).
type scope struct {
	current map[string]any
	parent  *scope
}

// runFieldCondition runs a `rel` condition or a condition with a `$` reference.
//
// A `rel` condition runs its inner conditions on the value of the related field (null if missing).
// A `$` reference is resolved to the value of the referenced field. `frm` and `nfr` accept a list
// (or the keys of a map) as reference, all other conditions need a single value.
// The resolved value is used in a new condition, because the AST is shared and must not be modified.
func (r *Validator) runFieldCondition(input any, v *model.AstValue, correct *int, sc *scope) (fatal bool, err error) {
	if sc == nil {
		return true, fmt.Errorf("%v with field %v can only be validated together with the other fields", v.ConditionType, v.Ref.Raw)
	}

	value, found, err := resolveFieldRef(sc, v.Ref)
	if err != nil {
		return true, err
	}

	if v.ConditionType == model.RELATION {
		if err := r.runConditions(value, v, nil, sc); err != nil {
			if v.Operator == model.IMPLIES {
				// A premise that is not fulfilled is the normal case, the error is not shown.
				return false, err
			}
			return false, fmt.Errorf("related field %v invalid: %v", v.Ref.Raw, err)
		}
		return false, nil
	}

	if !found || value == nil {
		return false, fmt.Errorf("referenced field %v is not given", v.Ref.Raw)
	}

	switch v.ConditionType {
	case model.FROM, model.NOT_FROM:
		return false, validators.ValidateFromValues(input, value, v.ConditionType == model.NOT_FROM, v.Ref.Raw)
	}

	conditionValue, err := helper.AnyToString(value)
	if err != nil {
		return false, fmt.Errorf("referenced field %v has to be a single value for %v: %v", v.Ref.Raw, v.ConditionType, err)
	}
	// A new condition without the values precomputed by the parser, they belong to the `$` value.
	resolved := model.AstValue{
		Type:           model.CONDITION,
		ConditionType:  v.ConditionType,
		ConditionValue: conditionValue,
		Start:          v.Start,
		End:            v.End,
	}
	return r.runCondition(input, &resolved, correct)
}

// resolveFieldRef returns the value of the referenced field, starting at the given scope.
// found is false if a field of the path does not exist.
func resolveFieldRef(sc *scope, ref *model.FieldRef) (value any, found bool, err error) {
	current := sc
	if ref.Root {
		for current.parent != nil {
			current = current.parent
		}
	}
	for i := 0; i < ref.Up; i++ {
		if current.parent == nil {
			return nil, false, fmt.Errorf("field path %v goes above the top level", ref.Raw)
		}
		current = current.parent
	}

	value = current.current
	for _, segment := range ref.Path {
		value, found = lookupField(value, segment)
		if !found {
			return nil, false, nil
		}
	}
	return value, true, nil
}

// lookupField returns the field of a JsonMap. On a list it returns the field of every item.
func lookupField(value any, name string) (any, bool) {
	switch v := value.(type) {
	case map[string]any:
		field, ok := v[name]
		return field, ok
	case []any:
		fields := make([]any, 0, len(v))
		for _, item := range v {
			if field, ok := lookupField(item, name); ok {
				fields = append(fields, field)
			}
		}
		return fields, true
	case nil:
		return nil, false
	}

	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Map:
		return helper.MapValue(value, name)
	case reflect.Slice, reflect.Array:
		fields := make([]any, 0, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			if field, ok := lookupField(rv.Index(i).Interface(), name); ok {
				fields = append(fields, field)
			}
		}
		return fields, true
	}
	return nil, false
}

// premiseToString formats the premise of an implication for error messages.
func premiseToString(premise *model.AstValue) string {
	if premise.Type == model.GROUP {
		return "(" + premise.AstGroupToString() + ")"
	}
	condition := *premise
	condition.Operator = ""
	return condition.AstConditionToString()
}
