package validators

import (
	"testing"

	"github.com/siherrmann/validator/model"
	"github.com/stretchr/testify/assert"
)

func TestValidateGroups(t *testing.T) {
	tests := []struct {
		name          string
		groups        map[string]*model.Group
		groupSize     map[string]int
		groupErrors   map[string][]error
		expectedError string
	}{
		{
			name: "Valid group with min condition",
			groups: map[string]*model.Group{
				"gr1": {Name: "gr1", ConditionType: model.MIN_VALUE, ConditionValue: "2"},
			},
			groupSize:     map[string]int{"gr1": 3},
			groupErrors:   map[string][]error{"gr1": {}},
			expectedError: "",
		},
		{
			name: "Invalid group with min condition",
			groups: map[string]*model.Group{
				"gr1": {Name: "gr1", ConditionType: model.MIN_VALUE, ConditionValue: "3"},
			},
			groupSize:     map[string]int{"gr1": 2},
			groupErrors:   map[string][]error{"gr1": {}},
			expectedError: "less then 3 in group gr1 without error, all errors: []",
		},
		{
			name: "Invalid group with invalid min condition",
			groups: map[string]*model.Group{
				"gr1": {Name: "gr1", ConditionType: model.MIN_VALUE, ConditionValue: "apple"},
			},
			groupSize:     map[string]int{"gr1": 2},
			groupErrors:   map[string][]error{"gr1": {}},
			expectedError: "strconv.Atoi: parsing \"apple\": invalid syntax",
		},
		{
			name: "Valid group with max condition",
			groups: map[string]*model.Group{
				"gr1": {Name: "gr1", ConditionType: model.MAX_VALUE, ConditionValue: "3"},
			},
			groupSize:     map[string]int{"gr1": 2},
			groupErrors:   map[string][]error{"gr1": {}},
			expectedError: "",
		},
		{
			name: "Invalid group with max condition",
			groups: map[string]*model.Group{
				"gr1": {Name: "gr1", ConditionType: model.MAX_VALUE, ConditionValue: "2"},
			},
			groupSize:     map[string]int{"gr1": 3},
			groupErrors:   map[string][]error{"gr1": {}},
			expectedError: "greater then 2 in group gr1 without error, all errors: []",
		},
		{
			name: "Invalid group with invalid max condition",
			groups: map[string]*model.Group{
				"gr1": {Name: "gr1", ConditionType: model.MAX_VALUE, ConditionValue: "apple"},
			},
			groupSize:     map[string]int{"gr1": 3},
			groupErrors:   map[string][]error{"gr1": {}},
			expectedError: "strconv.Atoi: parsing \"apple\": invalid syntax",
		},
		{
			name: "Invalid group condition type",
			groups: map[string]*model.Group{
				"gr1": {Name: "gr1", ConditionType: "invalid", ConditionValue: "1"},
			},
			groupSize:     map[string]int{"gr1": 1},
			groupErrors:   map[string][]error{"gr1": {}},
			expectedError: "invalid group condition type invalid",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateGroups(test.groups, test.groupSize, test.groupErrors)
			if test.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, test.expectedError)
			}
		})
	}
}
