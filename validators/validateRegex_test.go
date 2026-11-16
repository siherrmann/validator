package validators

import (
	"regexp"
	"testing"

	"github.com/siherrmann/validator/model"
	"github.com/stretchr/testify/assert"
)

func TestValidateRegex(t *testing.T) {
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
			name: "Valid regex string",
			args: args{
				v:   "hello123",
				ast: &model.AstValue{ConditionValue: "^[a-z]+[0-9]+$"},
			},
			wantErr: false,
		},
		{
			name: "Invalid regex string",
			args: args{
				v:   "hello@123",
				ast: &model.AstValue{ConditionValue: "^[a-z]+[0-9]+$"},
			},
			wantErr: true,
		},
		{
			name: "Valid regex array",
			args: args{
				v:   []string{"hello123"},
				ast: &model.AstValue{ConditionValue: "^[a-z]+[0-9]+$"},
			},
			wantErr: false,
		},
		{
			name: "Invalid regex array",
			args: args{
				v:   []string{"hello@123"},
				ast: &model.AstValue{ConditionValue: "^[a-z]+[0-9]+$"},
			},
			wantErr: true,
		},
		{
			name: "Invalid value type",
			args: args{
				v:   args{},
				ast: &model.AstValue{ConditionValue: "^[a-z]+[0-9]+$"},
			},
			wantErr: true,
		},
		{
			name: "Invalid value type array",
			args: args{
				v:   []struct{}{},
				ast: &model.AstValue{ConditionValue: "^[a-z]+[0-9]+$"},
			},
			wantErr: false,
		},
		{
			name: "Invalid condition value type",
			args: args{
				v:   "test",
				ast: &model.AstValue{ConditionValue: "banana"},
			},
			wantErr: true,
		},
		{
			name: "Invalid null regex",
			args: args{
				v:   nil,
				ast: &model.AstValue{ConditionValue: "^.*$"},
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateRegex(test.args.v, test.args.ast)
			if test.wantErr {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
			}
		})
	}

	t.Run("Regex compiled by the parser", func(t *testing.T) {
		ast := &model.AstValue{ConditionValue: "invalid(", Regex: regexp.MustCompile("^[A-Z]+$")}
		assert.NoError(t, ValidateRegex("ABC", ast))
		assert.Error(t, ValidateRegex("abc", ast))
		assert.NoError(t, ValidateRegex(map[string]int{"ABC": 1}, ast))
	})
}
