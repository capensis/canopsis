package pattern

import (
	"strings"
)

type Event [][]FieldCondition

func (p Event) HasField(field string) bool {
	for _, group := range p {
		for _, condition := range group {
			if condition.Field == field {
				return true
			}
		}
	}

	return false
}

func (p Event) HasExtraInfosField() bool {
	for _, group := range p {
		for _, condition := range group {
			if infoName := GetEventExtraInfoName(condition.Field); infoName != "" {
				return true
			}
		}
	}

	return false
}

func GetEventExtraInfoName(f string) string {
	if n, ok := strings.CutPrefix(f, "extra."); ok {
		return n
	}

	return ""
}
