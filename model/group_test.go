package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGroups(t *testing.T) {
	t.Run("Valid groups", func(t *testing.T) {
		groups, err := GetGroups("gr1min1 gr2max2")
		assert.NoError(t, err, "Expected no error when getting groups")
		assert.NotEmpty(t, groups, "Expected groups to be not empty")
		assert.Equal(t, 2, len(groups), "Expected 2 groups to be parsed")
		assert.Equal(t, &Group{Name: "gr1", ConditionType: "min", ConditionValue: "1"}, groups[0], "Expected first group to match")
		assert.Equal(t, &Group{Name: "gr2", ConditionType: "max", ConditionValue: "2"}, groups[1], "Expected second group to match")
	})

	t.Run("Invalid group", func(t *testing.T) {
		groups, err := GetGroups("invalid")
		assert.Error(t, err, "Expected error when getting groups")
		assert.Empty(t, groups, "Expected groups to be empty")
	})

	t.Run("Invalid group too short", func(t *testing.T) {
		groups, err := GetGroups("gr")
		assert.Error(t, err, "Expected error when getting groups")
		assert.Empty(t, groups, "Expected groups to be empty")
	})

	t.Run("Invalid group missing condition", func(t *testing.T) {
		groups, err := GetGroups("gr1")
		assert.Error(t, err, "Expected error when getting groups")
		assert.Empty(t, groups, "Expected groups to be empty")
	})

	t.Run("Invalid group invalid condition type", func(t *testing.T) {
		groups, err := GetGroups("gr1inv1")
		assert.Error(t, err, "Expected error when getting groups")
		assert.Empty(t, groups, "Expected groups to be empty")
	})

	t.Run("Invalid group empty condition value", func(t *testing.T) {
		groups, err := GetGroups("gr1equ")
		assert.Error(t, err, "Expected error when getting groups")
		assert.Empty(t, groups, "Expected groups to be empty")
	})
}
