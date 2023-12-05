package pattern

import (
	"strings"
)

type Entity [][]FieldCondition

func (p Entity) RemoveFields(fields []string) Entity {
	forbiddenFieldsMap := make(map[string]bool, len(fields))
	for _, field := range fields {
		forbiddenFieldsMap[field] = true
	}

	newGroups := make(Entity, 0, len(p))
	for _, group := range p {
		newGroup := make([]FieldCondition, 0, len(group))
		for _, condition := range group {
			if IsForbiddenEntityField(condition, forbiddenFieldsMap) {
				continue
			}

			newGroup = append(newGroup, condition)
		}
		if len(newGroup) > 0 {
			newGroups = append(newGroups, newGroup)
		}
	}

	if len(newGroups) > 0 {
		return newGroups
	}

	return nil
}

func GetEntityInfoName(f string) string {
	if n := strings.TrimPrefix(f, "infos."); n != f {
		return n
	}

	return ""
}

func GetEntityComponentInfoName(f string) string {
	if n := strings.TrimPrefix(f, "component_infos."); n != f {
		return n
	}

	return ""
}

func IsForbiddenEntityField(condition FieldCondition, forbiddenFieldsMap map[string]bool) bool {
	return forbiddenFieldsMap[condition.Field] ||
		forbiddenFieldsMap["infos"] && strings.HasPrefix(condition.Field, "infos") ||
		forbiddenFieldsMap["component_infos"] && strings.HasPrefix(condition.Field, "component_infos")
}
