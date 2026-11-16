package validators

import (
	"testing"

	"github.com/siherrmann/validator/model"
	"github.com/stretchr/testify/assert"
)

func TestValidateFrom(t *testing.T) {
	type args struct {
		v   any
		ast *model.AstValue
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid from",
			args: args{
				v:   "apple",
				ast: &model.AstValue{ConditionValue: "apple,banana"},
			},
			wantErr: false,
		},
		{
			name: "Invalid from",
			args: args{
				v:   "orange",
				ast: &model.AstValue{ConditionValue: "apple,banana"},
			},
			wantErr: true,
		},
		{
			name: "Invalid condition value type",
			args: args{
				v:   1,
				ast: &model.AstValue{ConditionValue: "apple,banana"},
			},
			wantErr: true,
		},
		{
			name: "Invalid null from",
			args: args{
				v:   nil,
				ast: &model.AstValue{ConditionValue: "apple,banana"},
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateFrom(test.args.v, test.args.ast)
			if test.wantErr {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
			}
		})
	}
}

func TestValidateNotFrom(t *testing.T) {
	type args struct {
		v   any
		ast *model.AstValue
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid not from",
			args: args{
				v:   "orange",
				ast: &model.AstValue{ConditionValue: "apple,banana"},
			},
			wantErr: false,
		},
		{
			name: "Invalid not from",
			args: args{
				v:   "apple",
				ast: &model.AstValue{ConditionValue: "apple,banana"},
			},
			wantErr: true,
		},
		{
			name: "Invalid condition value type",
			args: args{
				v:   1,
				ast: &model.AstValue{ConditionValue: "apple,banana"},
			},
			wantErr: true,
		},
		{
			name: "Valid null not from",
			args: args{
				v:   nil,
				ast: &model.AstValue{ConditionValue: "apple,banana"},
			},
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateNotFrom(test.args.v, test.args.ast)
			if test.wantErr {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
			}
		})
	}
}

func TestValidateFromValues(t *testing.T) {
	tests := []struct {
		name    string
		v       any
		from    any
		not     bool
		wantErr bool
	}{
		{name: "Value in list", v: 443, from: []any{80, 443}},
		{name: "Value not in list", v: 8080, from: []any{80, 443}, wantErr: true},
		{name: "Value in list of other type", v: 443, from: []any{80.0, 443.0}},
		{name: "Value in typed list", v: "a", from: []string{"a", "b"}},
		{name: "Value in map keys", v: "FOO", from: map[string]string{"FOO": "1"}},
		{name: "Single value", v: "a", from: "a"},
		{name: "List input", v: []string{"a", "b"}, from: []any{"a", "b", "c"}},
		{name: "List input with item not in list", v: []string{"a", "x"}, from: []any{"a", "b"}, wantErr: true},
		{name: "Map input keys", v: map[string]int{"a": 1}, from: []any{"a"}},
		{name: "Null value", v: nil, from: []any{"a"}, wantErr: true},
		{name: "Not from: value not in list", v: "BAR", from: map[string]string{"FOO": "1"}, not: true},
		{name: "Not from: value in list", v: []string{"FOO"}, from: map[string]string{"FOO": "1"}, not: true, wantErr: true},
		{name: "Not from: null value", v: nil, from: []any{"a"}, not: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateFromValues(test.v, test.from, test.not, "other")
			if test.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "other")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateFromWithParsedList(t *testing.T) {
	// The condition value is intentionally different, to show that the list split by the parser is used.
	ast := &model.AstValue{ConditionValue: "invalid", List: []string{"go", "dockerfile"}}
	assert.NoError(t, ValidateFrom("go", ast))
	assert.Error(t, ValidateFrom("oci", ast))
	assert.NoError(t, ValidateNotFrom("oci", ast))
	assert.Error(t, ValidateNotFrom("go", ast))
}
