package pattern

import (
	"fmt"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
)

type Alarm [][]FieldCondition

func (p Alarm) Match(alarm types.Alarm) (bool, error) {
	if len(p) == 0 {
		return true, nil
	}

	for _, group := range p {
		matched := false

		for _, v := range group {
			f := v.Field
			cond := v.Condition
			var err error
			matched = false

			if infoName := getAlarmInfoName(f); infoName != "" {
				infoVal, ok := getAlarmInfoVal(alarm, infoName)
				if v.FieldType == "" {
					matched, err = cond.MatchRef(infoVal)
				} else if ok {
					switch v.FieldType {
					case FieldTypeString:
						var s string
						if s, err = getStringValue(infoVal); err == nil {
							matched, _, err = cond.MatchString(s)
						}
					case FieldTypeInt:
						var i int64
						if i, err = getIntValue(infoVal); err == nil {
							matched, err = cond.MatchInt(i)
						}
					case FieldTypeBool:
						var b bool
						if b, err = getBoolValue(infoVal); err == nil {
							matched, err = cond.MatchBool(b)
						}
					case FieldTypeStringArray:
						var a []string
						if a, err = getStringArrayValue(infoVal); err == nil {
							matched, err = cond.MatchStringArray(a)
						}
					default:
						return false, fmt.Errorf("invalid field type for %q field: %s", f, v.FieldType)
					}
				}

				if err != nil {
					return false, fmt.Errorf("invalid condition for %q field: %w", f, err)
				}

				if !matched {
					break
				}

				continue
			}

			if str, ok := getAlarmStringField(alarm, f); ok {
				matched, _, err = cond.MatchString(str)
			} else if i, ok := getAlarmIntField(alarm, f); ok {
				matched, err = cond.MatchInt(i)
			} else if r, ok := getAlarmRefField(alarm, f); ok {
				matched, err = cond.MatchRef(r)
			} else if t, ok := getAlarmTimeField(alarm, f); ok {
				matched, err = cond.MatchTime(t)
			} else if d, ok := getAlarmDurationField(alarm, f); ok {
				matched, err = cond.MatchDuration(d)
			} else {
				err = ErrUnsupportedField
			}

			if err != nil {
				return false, fmt.Errorf("invalid condition for %q field: %w", f, err)
			}

			if !matched {
				break
			}
		}

		if matched {
			return true, nil
		}
	}

	return false, nil
}

func (p Alarm) Validate(forbiddenFields, onlyTimeAbsoluteFields []string) bool {
	emptyAlarm := types.Alarm{}
	forbiddenFieldsMap := make(map[string]bool, len(forbiddenFields))
	for _, field := range forbiddenFields {
		forbiddenFieldsMap[field] = true
	}
	timeAbsoluteFieldsMap := make(map[string]bool, len(onlyTimeAbsoluteFields))
	for _, field := range onlyTimeAbsoluteFields {
		timeAbsoluteFieldsMap[field] = true
	}

	for _, group := range p {
		if len(group) == 0 {
			return false
		}

		for _, v := range group {
			f := v.Field
			cond := v.Condition
			var err error

			if isForbiddenAlarmField(v, forbiddenFieldsMap, timeAbsoluteFieldsMap) {
				return false
			}

			if infoName := getAlarmInfoName(f); infoName != "" {
				switch v.FieldType {
				case FieldTypeString:
					_, _, err = cond.MatchString("")
				case FieldTypeInt:
					_, err = cond.MatchInt(0)
				case FieldTypeBool:
					_, err = cond.MatchBool(false)
				case FieldTypeStringArray:
					_, err = cond.MatchStringArray([]string{})
				default:
					_, err = cond.MatchRef(nil)
				}

				if err != nil {
					return false
				}

				continue
			}

			if str, ok := getAlarmStringField(emptyAlarm, f); ok {
				_, _, err = cond.MatchString(str)
			} else if i, ok := getAlarmIntField(emptyAlarm, f); ok {
				_, err = cond.MatchInt(i)
			} else if r, ok := getAlarmRefField(emptyAlarm, f); ok {
				_, err = cond.MatchRef(r)
			} else if t, ok := getAlarmTimeField(emptyAlarm, f); ok {
				_, err = cond.MatchTime(t)
			} else if d, ok := getAlarmDurationField(emptyAlarm, f); ok {
				_, err = cond.MatchDuration(d)
			} else {
				err = ErrUnsupportedField
			}

			if err != nil {
				return false
			}
		}
	}

	return true
}

func (p Alarm) ToMongoQuery(prefix string) (bson.M, error) {
	if len(p) == 0 {
		return nil, nil
	}

	if prefix != "" {
		prefix += "."
	}

	groupQueries := make([]bson.M, len(p))
	var err error

	for i, group := range p {
		condQueries := make([]bson.M, len(group))
		for j, cond := range group {
			f := cond.Field

			if infoName := getAlarmInfoName(f); infoName != "" {
				f = prefix + "v.infos_array.v." + infoName

				condQueries[j], err = cond.Condition.ToMongoQuery(f)
				if err != nil {
					return nil, fmt.Errorf("invalid condition for %q field: %w", f, err)
				}

				conds := getTypeMongoQuery(f, cond.FieldType)

				if len(conds) > 0 {
					conds = append(conds, condQueries[j])
					condQueries[j] = bson.M{"$and": conds}
				}

				continue
			}

			f = prefix + f
			condQueries[j], err = cond.Condition.ToMongoQuery(f)
			if err != nil {
				return nil, fmt.Errorf("invalid condition for %q field: %w", f, err)
			}
		}

		groupQueries[i] = bson.M{"$and": condQueries}
	}

	return bson.M{"$or": groupQueries}, nil
}

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

			if infoName := getAlarmInfoName(f); infoName != "" {
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
			if infoName := getAlarmInfoName(condition.Field); infoName != "" {
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
			if isForbiddenAlarmField(condition, forbiddenFieldsMap, timeAbsoluteFieldsMap) {
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

func getAlarmStringField(alarm types.Alarm, f string) (string, bool) {
	switch f {
	case "v.display_name":
		return alarm.Value.DisplayName, true
	case "v.output":
		return alarm.Value.Output, true
	case "v.connector":
		return alarm.Value.Connector, true
	case "v.connector_name":
		return alarm.Value.ConnectorName, true
	case "v.component":
		return alarm.Value.Component, true
	case "v.resource":
		return alarm.Value.Resource, true
	case "v.last_comment.m":
		if alarm.Value.LastComment == nil {
			return "", true
		}
		return alarm.Value.LastComment.Message, true
	case "v.ack.a":
		if alarm.Value.ACK == nil {
			return "", true
		}

		return alarm.Value.ACK.Author, true
	default:
		return "", false
	}
}

func getAlarmIntField(alarm types.Alarm, f string) (int64, bool) {
	switch f {
	case "v.state.val":
		if alarm.Value.State == nil {
			return 0, true
		}
		return int64(alarm.Value.State.Value), true
	case "v.status.val":
		if alarm.Value.Status == nil {
			return 0, true
		}
		return int64(alarm.Value.Status.Value), true
	default:
		return 0, false
	}
}

func getAlarmRefField(alarm types.Alarm, f string) (interface{}, bool) {
	switch f {
	case "v.ack":
		if alarm.Value.ACK == nil {
			return nil, true
		}
		return alarm.Value.ACK, true
	case "v.ticket":
		if alarm.Value.Ticket == nil {
			return nil, true
		}
		return alarm.Value.Ticket, true
	case "v.canceled":
		if alarm.Value.Canceled == nil {
			return nil, true
		}
		return alarm.Value.Canceled, true
	case "v.snooze":
		if alarm.Value.Snooze == nil {
			return nil, true
		}
		return alarm.Value.Snooze, true
	default:
		return nil, false
	}
}

func getAlarmTimeField(alarm types.Alarm, field string) (time.Time, bool) {
	switch field {
	case "v.creation_date":
		return alarm.Value.CreationDate.Time, true
	case "v.last_event_date":
		return alarm.Value.LastEventDate.Time, true
	case "v.last_update_date":
		return alarm.Value.LastUpdateDate.Time, true
	case "v.ack.t":
		if alarm.Value.ACK != nil {
			return alarm.Value.ACK.Timestamp.Time, true
		}

		return time.Time{}, true
	case "v.resolved":
		if alarm.Value.Resolved != nil {
			return alarm.Value.Resolved.Time, true
		}

		return time.Time{}, true
	default:
		return time.Time{}, false
	}
}

func getAlarmDurationField(alarm types.Alarm, field string) (int64, bool) {
	switch field {
	case "v.duration":
		if alarm.Value.Duration > 0 {
			return alarm.Value.Duration, true
		}

		if alarm.Value.Resolved != nil {
			return int64(alarm.Value.Resolved.Sub(alarm.Time.Time).Seconds()), true
		}

		return int64(time.Since(alarm.Time.Time).Seconds()), true
	default:
		return 0, false
	}
}

func getAlarmInfoVal(alarm types.Alarm, f string) (interface{}, bool) {
	for _, infosByRule := range alarm.Value.Infos {
		if v, ok := infosByRule[f]; ok {
			return v, true
		}
	}

	return nil, false
}

func getAlarmInfoName(f string) string {
	if n := strings.TrimPrefix(f, "v.infos."); n != f {
		return n
	}

	return ""
}

func getTypeMongoQuery(f, ft string) []bson.M {
	var conds []bson.M

	switch ft {
	case FieldTypeString:
		conds = []bson.M{{f: bson.M{"$type": "string"}}}
	case FieldTypeInt:
		conds = []bson.M{{f: bson.M{"$type": bson.A{"long", "int", "decimal"}}}}
	case FieldTypeBool:
		conds = []bson.M{{f: bson.M{"$type": "bool"}}}
	case FieldTypeStringArray:
		// Cond {"$type": "string"} checks only if an array contains at least one string element,
		// other elements can be any type.
		conds = []bson.M{{f: bson.M{"$type": "array"}}, {f: bson.M{"$type": "string"}}}
	}

	return conds
}

func isForbiddenAlarmField(condition FieldCondition, forbiddenFieldsMap map[string]bool, timeAbsoluteFieldsMap map[string]bool) bool {
	return forbiddenFieldsMap[condition.Field] ||
		forbiddenFieldsMap["v.infos"] && strings.HasPrefix(condition.Field, "v.infos") ||
		timeAbsoluteFieldsMap[condition.Field] && condition.Condition.Type == ConditionTimeRelative
}
