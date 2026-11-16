package validators

import (
	"testing"

	"github.com/siherrmann/validator/model"
	"github.com/stretchr/testify/assert"
)

func TestValidateUnique(t *testing.T) {
	long := make([]int, 40)
	for i := range long {
		long[i] = i
	}
	longDuplicate := append(append([]int{}, long...), 7)

	tests := []struct {
		name    string
		v       any
		key     string
		wantErr bool
	}{
		{name: "Null value", v: nil},
		{name: "Empty list", v: []string{}},
		{name: "Unique strings", v: []string{"a", "b", "c"}},
		{name: "Duplicate strings", v: []string{"a", "b", "a"}, wantErr: true},
		{name: "Duplicate any", v: []any{1.0, 2.0, 1.0}, wantErr: true},
		{name: "Unique long list", v: long},
		{name: "Duplicate long list", v: longDuplicate, wantErr: true},
		{name: "Not comparable items", v: []any{[]int{1}, []int{2}}},
		{name: "Duplicate not comparable items", v: []any{[]int{1}, []int{1}}, wantErr: true},
		{name: "Map", v: map[string]int{"a": 1, "b": 1}},
		{name: "Unique key", v: []any{map[string]any{"name": "a"}, map[string]any{"name": "b"}}, key: "name"},
		{name: "Duplicate key", v: []any{map[string]any{"name": "a"}, map[string]any{"name": "a"}}, key: "name", wantErr: true},
		{name: "Missing key is ignored", v: []any{map[string]any{"name": "a"}, map[string]any{}}, key: "name"},
		{name: "Duplicate key in typed maps", v: []map[string]string{{"name": "a"}, {"name": "a"}}, key: "name", wantErr: true},
		{name: "Key on list of strings", v: []string{"a"}, key: "name", wantErr: true},
		{name: "No list", v: "abc", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateUnique(test.v, &model.AstValue{ConditionType: model.UNIQUE, ConditionValue: test.key})
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
