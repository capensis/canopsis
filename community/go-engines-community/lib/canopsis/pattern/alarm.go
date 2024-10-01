package pattern

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Alarm [][]FieldCondition

func (p Alarm) GetMongoFields(prefix string) bson.M {
	if len(p) == 0 {
		return nil
	}

	if prefix != "" {
		prefix += "."
	}

	withDuration := false
	withInfos := false

	for _, group := range p {
		for _, cond := range group {
			f := cond.Field

			if infoName := GetAlarmInfoName(f); infoName != "" {
				withInfos = true

				continue
			}

			if f == "v.duration" {
				withDuration = true
			}
		}
	}

	addFields := bson.M{}
	if withDuration {
		addFields[prefix+"v.duration"] = bson.M{"$ifNull": bson.A{
			"$" + prefix + "v.duration",
			bson.M{"$subtract": bson.A{
				bson.M{"$cond": bson.M{
					"if":   "$" + prefix + "v.resolved",
					"then": "$" + prefix + "v.resolved",
					"else": time.Now().Unix(),
				}},
				"$" + prefix + "v.creation_date",
			}},
		}}
	}

	if withInfos {
		addFields[prefix+"v.infos_array"] = bson.M{"$objectToArray": "$" + prefix + "v.infos"}
	}

	return addFields
}

func (p Alarm) HasField(field string) bool {
	for _, group := range p {
		for _, condition := range group {
			if condition.Field == field {
				return true
			}
		}
	}

	return false
}

func (p Alarm) HasInfosField() bool {
	for _, group := range p {
		for _, condition := range group {
			if infoName := GetAlarmInfoName(condition.Field); infoName != "" {
				return true
			}
		}
	}

	return false
}

func (p Alarm) RemoveFields(fields, onlyTimeAbsoluteFields []string) Alarm {
	forbiddenFieldsMap := make(map[string]bool, len(fields))
	for _, field := range fields {
		forbiddenFieldsMap[field] = true
	}
	timeAbsoluteFieldsMap := make(map[string]bool, len(onlyTimeAbsoluteFields))
	for _, field := range onlyTimeAbsoluteFields {
		timeAbsoluteFieldsMap[field] = true
	}

	newGroups := make(Alarm, 0, len(p))
	for _, group := range p {
		newGroup := make([]FieldCondition, 0, len(group))
		for _, condition := range group {
			if IsForbiddenAlarmField(condition, forbiddenFieldsMap, timeAbsoluteFieldsMap) {
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

func GetAlarmInfoName(f string) string {
	if n, ok := strings.CutPrefix(f, "v.infos."); ok {
		return n
	}

	return ""
}

func IsForbiddenAlarmField(condition FieldCondition, forbiddenFieldsMap map[string]bool, timeAbsoluteFieldsMap map[string]bool) bool {
	return forbiddenFieldsMap[condition.Field] ||
		forbiddenFieldsMap["v.infos"] && strings.HasPrefix(condition.Field, "v.infos") ||
		timeAbsoluteFieldsMap[condition.Field] && condition.Condition.Type == ConditionTimeRelative
}
