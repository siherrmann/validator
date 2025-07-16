package validators

import (
	"testing"

	"github.com/siherrmann/validator/model"
	"github.com/stretchr/testify/assert"
)

func TestValidateEqual(t *testing.T) {
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
			name: "Valid equal",
			args: args{
				v:   5,
				ast: &model.AstValue{ConditionValue: "5"},
			},
			wantErr: false,
		},
		{
			name: "Invalid equal",
			args: args{
				v:   5,
				ast: &model.AstValue{ConditionValue: "6"},
			},
			wantErr: true,
		},
		{
			name: "Invalid value type",
			args: args{
				v:   "test",
				ast: &model.AstValue{ConditionValue: "5"},
			},
			wantErr: true,
		},
		{
			name: "Invalid condition value type",
			args: args{
				v:   1,
				ast: &model.AstValue{ConditionValue: "banana"},
			},
			wantErr: true,
		},
		{
			name: "Invalid condition value type array",
			args: args{
				v:   []int{1, 2, 3},
				ast: &model.AstValue{ConditionValue: "banana"},
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateEqual(test.args.v, test.args.ast)
			if test.wantErr {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
			}
		})
	}
}

func TestValidateNotEqual(t *testing.T) {
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
			name: "Valid not equal",
			args: args{
				v:   5,
				ast: &model.AstValue{ConditionValue: "6"},
			},
			wantErr: false,
		},
		{
			name: "Invalid not equal",
			args: args{
				v:   5,
				ast: &model.AstValue{ConditionValue: "5"},
			},
			wantErr: true,
		},
		{
			name: "Invalid condition value type",
			args: args{
				v:   args{},
				ast: &model.AstValue{ConditionValue: "banana"},
			},
			wantErr: true,
		},
		{
			name: "Invalid condition value type array",
			args: args{
				v:   []int{1, 2, 3},
				ast: &model.AstValue{ConditionValue: "banana"},
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateNotEqual(test.args.v, test.args.ast)
			if test.wantErr {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
			}
		})
	}
}
