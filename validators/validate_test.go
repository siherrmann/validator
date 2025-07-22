package validators

import (
	"testing"

	"github.com/siherrmann/validator/model"
	"github.com/stretchr/testify/assert"
)

func TestRunFuncOnConditionGroup(t *testing.T) {
	type args struct {
		input    any
		astValue *model.AstValue
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid condition group with empty condition",
			args: args{
				input: "apple",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:     model.EMPTY,
							Operator: model.AND,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Valid condition group with NONE condition type",
			args: args{
				input: "apple",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:          model.CONDITION,
							ConditionType: model.NONE,
							Operator:      model.AND,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid condition group with invalid condition type",
			args: args{
				input: "apple",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:          model.CONDITION,
							ConditionType: model.ConditionType("Uff"),
							Operator:      model.AND,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid condition group with inner group",
			args: args{
				input: "apple",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type: model.GROUP,
							ConditionGroup: []*model.AstValue{
								{
									Type:           model.CONDITION,
									ConditionType:  model.EQUAL,
									ConditionValue: "apple",
									Operator:       model.AND,
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid with OR operator",
			args: args{
				input: "apple",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.EQUAL,
							ConditionValue: "banana",
							Operator:       model.OR,
						},
						{
							Type:           model.CONDITION,
							ConditionType:  model.EQUAL,
							ConditionValue: "cherry",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid value EQUAL",
			args: args{
				input: "apple",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.EQUAL,
							ConditionValue: "apple",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid value EQUAL",
			args: args{
				input: "apple",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.EQUAL,
							ConditionValue: "banana",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid value NOT_EQUAL",
			args: args{
				input: "apple",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.NOT_EQUAL,
							ConditionValue: "banana",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid value NOT_EQUAL",
			args: args{
				input: "apple",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.NOT_EQUAL,
							ConditionValue: "apple",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid value MIN_VALUE",
			args: args{
				input: 5,
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.MIN_VALUE,
							ConditionValue: "3",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid value MIN_VALUE",
			args: args{
				input: 3,
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.MIN_VALUE,
							ConditionValue: "5",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid value MAX_VALUE",
			args: args{
				input: 3,
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.MAX_VALUE,
							ConditionValue: "5",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid value MAX_VALUE",
			args: args{
				input: 5,
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.MAX_VALUE,
							ConditionValue: "3",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid value CONTAINS",
			args: args{
				input: "pineapple",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.CONTAINS,
							ConditionValue: "apple",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid value CONTAINS",
			args: args{
				input: "banana",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.CONTAINS,
							ConditionValue: "apple",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid value NOT_CONTAINS",
			args: args{
				input: "banana",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.NOT_CONTAINS,
							ConditionValue: "apple",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid value NOT_CONTAINS",
			args: args{
				input: "pineapple",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.NOT_CONTAINS,
							ConditionValue: "apple",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid value FROM",
			args: args{
				input: "apple",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.FROM,
							ConditionValue: "apple,banana",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid value FROM",
			args: args{
				input: "cherry",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.FROM,
							ConditionValue: "apple,banana",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid value NOT_FROM",
			args: args{
				input: "cherry",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.NOT_FROM,
							ConditionValue: "apple,banana",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid value NOT_FROM",
			args: args{
				input: "apple",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.NOT_FROM,
							ConditionValue: "apple,banana",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid value REGX",
			args: args{
				input: "abc123",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.REGX,
							ConditionValue: "^[a-z]+[0-9]+$",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid value REGX",
			args: args{
				input: "abc@123",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.REGX,
							ConditionValue: "^[a-z]+[0-9]+$",
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RunFuncOnConditionGroup(tt.args.input, tt.args.astValue)
			if tt.wantErr {
				assert.Error(t, err, "Expected an error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
			}
		})
	}
}
