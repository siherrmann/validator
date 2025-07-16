package validators

import (
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
			name: "Valid regex",
			args: args{
				v:   "hello123",
				ast: &model.AstValue{ConditionValue: "^[a-z]+[0-9]+$"},
			},
			wantErr: false,
		},
		{
			name: "Invalid regex",
			args: args{
				v:   "hello@123",
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
			name: "Invalid condition value type",
			args: args{
				v:   "test",
				ast: &model.AstValue{ConditionValue: "banana"},
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
}
