package validators

import (
	"testing"

	"github.com/siherrmann/validator/model"
	"github.com/stretchr/testify/assert"
)

func TestValidateContains(t *testing.T) {
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
			name: "Valid contains",
			args: args{
				v:   "banana",
				ast: &model.AstValue{ConditionValue: "ana"},
			},
			wantErr: false,
		},
		{
			name: "Invalid contains",
			args: args{
				v:   "banana",
				ast: &model.AstValue{ConditionValue: "ane"},
			},
			wantErr: true,
		},
		{
			name: "Invalid value type",
			args: args{
				v:   123,
				ast: &model.AstValue{ConditionValue: "2"},
			},
			wantErr: true,
		},
		{
			name: "Invalid condition value type",
			args: args{
				v:   []int{1, 2, 3},
				ast: &model.AstValue{ConditionValue: "banana"},
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateContains(test.args.v, test.args.ast)
			if test.wantErr {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
			}
		})
	}
}

func TestValidateNotContains(t *testing.T) {
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
			name: "Valid not contains",
			args: args{
				v:   "banana",
				ast: &model.AstValue{ConditionValue: "apple"},
			},
			wantErr: false,
		},
		{
			name: "Invalid not contains",
			args: args{
				v:   "banana",
				ast: &model.AstValue{ConditionValue: "ana"},
			},
			wantErr: true,
		},
		{
			name: "Invalid value type",
			args: args{
				v:   123,
				ast: &model.AstValue{ConditionValue: "2"},
			},
			wantErr: true,
		},
		{
			name: "Invalid condition value type",
			args: args{
				v:   []int{1, 2, 3},
				ast: &model.AstValue{ConditionValue: "banana"},
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateNotContains(test.args.v, test.args.ast)
			if test.wantErr {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
			}
		})
	}
}
