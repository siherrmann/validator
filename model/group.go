package model

import (
	"fmt"
	"strings"
)

type Group struct {
	Name           string
	ConditionType  ConditionType
	ConditionValue string
}

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
