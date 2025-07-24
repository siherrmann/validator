package model

import (
	"fmt"
	"strings"
)

// Group represents a validation group with a name, condition type, and condition value.
// The name is expected to start with "gr" followed by a condition type and value.
// The condition type is validated against known condition types, allowed are min and max.
// The condition value is extracted based on the condition type.
type Group struct {
	Name           string
	ConditionType  ConditionType
	ConditionValue string
}

// GetGroups parses a string to extract validation groups.
// It splits the string by spaces and creates a Group for each part.
// Each group must start with "gr" and can have a condition type and value.
// If the group name is invalid or the condition type is not recognized, an error is returned.
func GetGroups(s string) ([]*Group, error) {
	groups := []*Group{}
	groupsString := strings.Split(s, " ")
	for _, g := range groupsString {
		group := &Group{}
		if len(g) > 2 {
			group.Name = g[:3]
		} else {
			group.Name = g
		}

		if !strings.HasPrefix(group.Name, "gr") {
			return nil, fmt.Errorf("invalid group name: %s", group.Name)
		}

		var err error
		var groupCondition string
		if len(g[3:]) > 0 {
			groupCondition = g[3:]
		} else {
			groupCondition = ""
		}
		group.ConditionType, err = GetConditionType(groupCondition)
		if err != nil {
			return nil, err
		}
		group.ConditionValue, err = GetConditionByType(groupCondition, group.ConditionType)
		if err != nil {
			return nil, err
		}

		groups = append(groups, group)
	}
	return groups, nil
}
