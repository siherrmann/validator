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
