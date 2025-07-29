package validators

import (
	"testing"

	"github.com/siherrmann/validator/model"
	"github.com/stretchr/testify/assert"
)

func TestValidateMin(t *testing.T) {
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
			name: "Valid min",
			args: args{
				v:   10,
				ast: &model.AstValue{ConditionValue: "5"},
			},
			wantErr: false,
		},
		{
			name: "Invalid min",
			args: args{
				v:   10,
				ast: &model.AstValue{ConditionValue: "15"},
			},
			wantErr: true,
		},
		{
			name: "Invalid value type",
			args: args{
				v:   true,
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
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateMin(test.args.v, test.args.ast)
			if test.wantErr {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
			}
		})
	}
}

func TestValidateMax(t *testing.T) {
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
			name: "Valid max",
			args: args{
				v:   10,
				ast: &model.AstValue{ConditionValue: "15"},
			},
			wantErr: false,
		},
		{
			name: "Invalid max",
			args: args{
				v:   10,
				ast: &model.AstValue{ConditionValue: "5"},
			},
			wantErr: true,
		},
		{
			name: "Invalid value type",
			args: args{
				v:   true,
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
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateMax(test.args.v, test.args.ast)
			if test.wantErr {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
			}
		})
	}
}
