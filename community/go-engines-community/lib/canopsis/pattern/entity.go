package pattern

import (
	"strings"
)

const (
	entityInfosPrefix          = "infos."
	componentEntityInfosPrefix = "component_infos."
)

type Entity [][]FieldCondition

func (p Entity) RemoveFields(forbiddenFieldsMap map[string]bool) Entity {
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

func (p Entity) GetInfosNames() []string {
	var keys []string
	keysMap := make(map[string]bool)

	for _, group := range p {
		for _, cond := range group {
			if n, ok := strings.CutPrefix(cond.Field, entityInfosPrefix); ok && !keysMap[n] {
				keys = append(keys, n)
				keysMap[n] = true
			}
		}
	}

	return keys
}

func (p Entity) GetComponentInfosNames() []string {
	var keys []string
	keysMap := make(map[string]bool)

	for _, group := range p {
		for _, cond := range group {
			if n, ok := strings.CutPrefix(cond.Field, componentEntityInfosPrefix); ok && !keysMap[n] {
				keys = append(keys, n)
				keysMap[n] = true
			}
		}
	}

	return keys
}

func GetEntityInfoName(f string) string {
	if n, ok := strings.CutPrefix(f, entityInfosPrefix); ok {
		return n
	}

	return ""
}

func GetEntityComponentInfoName(f string) string {
	if n, ok := strings.CutPrefix(f, componentEntityInfosPrefix); ok {
		return n
	}

	return ""
}

func IsForbiddenEntityField(condition FieldCondition, forbiddenFieldsMap map[string]bool) bool {
	return forbiddenFieldsMap[condition.Field] ||
		forbiddenFieldsMap["infos"] && (strings.HasPrefix(condition.Field, entityInfosPrefix) || condition.Alias != "") ||
		forbiddenFieldsMap["component_infos"] && strings.HasPrefix(condition.Field, componentEntityInfosPrefix)
}
