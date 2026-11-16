package validators

import (
	"fmt"
	"reflect"

	"github.com/siherrmann/validator/helper"
	"github.com/siherrmann/validator/model"
)

// uniqueMapThreshold is the number of items from which duplicates are found with a map
// instead of comparing every pair. Small lists are faster without the map allocation.
const uniqueMapThreshold = 16

// ValidateUnique checks that every item of a list is unique.
// With a condition value (eg. `unqname`), the field with this name of every item (a JsonMap) has to be unique.
// Items without this field are ignored. A null value and a map (keys are always unique) are valid.
func ValidateUnique(v any, ast *model.AstValue) error {
	if v == nil {
		return nil
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map:
		return nil
	case reflect.Slice, reflect.Array:
	default:
		return fmt.Errorf("value to validate has to be an array, slice or map, was %v", rv.Kind())
	}

	key := ast.ConditionValue
	items := make([]any, 0, rv.Len())
	comparable := true
	for i := 0; i < rv.Len(); i++ {
		item := rv.Index(i).Interface()
		if key != "" {
			field, ok, err := uniqueField(item, key)
			if err != nil {
				return err
			} else if !ok {
				continue
			}
			item = field
		}
		if item != nil && !reflect.TypeOf(item).Comparable() {
			comparable = false
		}
		items = append(items, item)
	}

	if comparable && len(items) > uniqueMapThreshold {
		seen := make(map[any]struct{}, len(items))
		for _, item := range items {
			if _, ok := seen[item]; ok {
				return uniqueError(item, key)
			}
			seen[item] = struct{}{}
		}
		return nil
	}

	for i := 1; i < len(items); i++ {
		for j := 0; j < i; j++ {
			if equalValue(items[i], items[j]) {
				return uniqueError(items[i], key)
			}
		}
	}
	return nil
}

// uniqueField returns the field with the given name of a list item.
func uniqueField(item any, key string) (any, bool, error) {
	if rv := reflect.ValueOf(item); rv.Kind() != reflect.Map || rv.Type().Key().Kind() != reflect.String {
		return nil, false, fmt.Errorf("unq%v needs a list of objects, item was %T", key, item)
	}
	field, ok := helper.MapValue(item, key)
	return field, ok, nil
}

func uniqueError(item any, key string) error {
	if key != "" {
		return fmt.Errorf("value %v of %v is not unique", item, key)
	}
	return fmt.Errorf("value %v is not unique", item)
}
