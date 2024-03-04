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
							matched, err = cond.MatchString(s)
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

			foundField := false
			if str, ok := getAlarmStringField(alarm, f); ok {
				foundField = true
				matched, err = cond.MatchString(str)
			}
			if !foundField || err != nil {
				if i, ok := getAlarmIntField(alarm, f); ok {
					foundField = true
					matched, err = cond.MatchInt(i)
				}
			}
			if !foundField || err != nil {
				if r, ok := getAlarmRefField(alarm, f); ok {
					foundField = true
					matched, err = cond.MatchRef(r)
				}
			}
			if !foundField || err != nil {
				if t, ok := getAlarmTimeField(alarm, f); ok {
					foundField = true
					matched, err = cond.MatchTime(t)
				}
			}
			if !foundField || err != nil {
				if d, ok := getAlarmDurationField(alarm, f); ok {
					foundField = true
					matched, err = cond.MatchDuration(d)
				}
			}
			if !foundField || err != nil {
				if a, ok := getAlarmStringArrayField(alarm, f); ok {
					foundField = true
					matched, err = cond.MatchStringArray(a)
				}
			}

			if !foundField {
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
					_, err = cond.MatchString("")
				case FieldTypeInt:
					_, err = cond.MatchInt(0)
				case FieldTypeBool:
					_, err = cond.MatchBool(false)
				case FieldTypeStringArray:
					_, err = cond.MatchStringArray([]string{})
				case "":
					_, err = cond.MatchRef(nil)
				default:
					return false
				}

				if err != nil {
					return false
				}

				continue
			}

			foundField := false
			if str, ok := getAlarmStringField(emptyAlarm, f); ok {
				foundField = true
				_, err = cond.MatchString(str)
			}
			if !foundField || err != nil {
				if i, ok := getAlarmIntField(emptyAlarm, f); ok {
					foundField = true
					_, err = cond.MatchInt(i)
				}
			}
			if !foundField || err != nil {
				if r, ok := getAlarmRefField(emptyAlarm, f); ok {
					foundField = true
					_, err = cond.MatchRef(r)
				}
			}
			if !foundField || err != nil {
				if t, ok := getAlarmTimeField(emptyAlarm, f); ok {
					foundField = true
					_, err = cond.MatchTime(t)
				}
			}
			if !foundField || err != nil {
				if d, ok := getAlarmDurationField(emptyAlarm, f); ok {
					foundField = true
					_, err = cond.MatchDuration(d)
				}
			}
			if !foundField || err != nil {
				if a, ok := getAlarmStringArrayField(emptyAlarm, f); ok {
					foundField = true
					_, err = cond.MatchStringArray(a)
				}
			}

			if !foundField {
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
	groupQueries, err := p.getGroupMongoQueries(prefix)
	if err != nil || len(groupQueries) == 0 {
		return nil, err
	}
	return bson.M{"$or": groupQueries}, nil
}

func (p Alarm) ToNegativeMongoQuery(prefix string) (bson.M, error) {
	groupQueries, err := p.getGroupMongoQueries(prefix)
	if err != nil || len(groupQueries) == 0 {
		return nil, err
	}
	return bson.M{"$nor": groupQueries}, nil
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

func (p Alarm) getGroupMongoQueries(prefix string) ([]bson.M, error) {
	if len(p) == 0 {
		return nil, nil
	}

	if prefix != "" {
		prefix += "."
	}

	emptyAlarm := types.Alarm{}
	groupQueries := make([]bson.M, len(p))
	var err error

	for i, group := range p {
		condQueries := make([]bson.M, len(group))
		for j, cond := range group {
			if infoName := getAlarmInfoName(cond.Field); infoName != "" {
				mongoField := prefix + "v.infos_array.v." + infoName

				switch cond.FieldType {
				case FieldTypeString:
					condQueries[j], err = cond.Condition.StringToMongoQuery(mongoField, true)
				case FieldTypeInt:
					condQueries[j], err = cond.Condition.IntToMongoQuery(mongoField, true)
				case FieldTypeBool:
					condQueries[j], err = cond.Condition.BoolToMongoQuery(mongoField)
				case FieldTypeStringArray:
					condQueries[j], err = cond.Condition.StringArrayToMongoQuery(mongoField, true)
				case "":
					condQueries[j], err = cond.Condition.RefToMongoQuery(mongoField)
				default:
					return nil, fmt.Errorf("invalid field type for %q field: %s", cond.Field, cond.FieldType)
				}
				if err != nil {
					return nil, fmt.Errorf("invalid condition for %q field: %w", cond.Field, err)
				}

				continue
			}

			mongoField := prefix + cond.Field
			foundField := false
			if _, ok := getAlarmStringField(emptyAlarm, cond.Field); ok {
				foundField = true
				condQueries[j], err = cond.Condition.StringToMongoQuery(mongoField, false)
			}
			if !foundField || err != nil {
				if _, ok := getAlarmIntField(emptyAlarm, cond.Field); ok {
					foundField = true
					condQueries[j], err = cond.Condition.IntToMongoQuery(mongoField, false)
				}
			}
			if !foundField || err != nil {
				if _, ok := getAlarmRefField(emptyAlarm, cond.Field); ok {
					foundField = true
					condQueries[j], err = cond.Condition.RefToMongoQuery(mongoField)
				}
			}
			if !foundField || err != nil {
				if _, ok := getAlarmTimeField(emptyAlarm, cond.Field); ok {
					foundField = true
					condQueries[j], err = cond.Condition.TimeToMongoQuery(mongoField)
				}
			}
			if !foundField || err != nil {
				if _, ok := getAlarmDurationField(emptyAlarm, cond.Field); ok {
					foundField = true
					condQueries[j], err = cond.Condition.DurationToMongoQuery(mongoField)
				}
			}
			if !foundField || err != nil {
				if _, ok := getAlarmStringArrayField(emptyAlarm, cond.Field); ok {
					foundField = true
					condQueries[j], err = cond.Condition.StringArrayToMongoQuery(mongoField, false)
				}
			}

			if !foundField {
				err = ErrUnsupportedField
			}
			if err != nil {
				return nil, fmt.Errorf("invalid condition for %q field: %w", cond.Field, err)
			}
		}

		groupQueries[i] = bson.M{"$and": condQueries}
	}

	return groupQueries, nil
}

func getAlarmStringField(alarm types.Alarm, f string) (string, bool) {
	switch f {
	case "v.display_name":
		return alarm.Value.DisplayName, true
	case "v.output":
		return alarm.Value.Output, true
	case "v.long_output":
		return alarm.Value.LongOutput, true
	case "v.initial_output":
		return alarm.Value.InitialOutput, true
	case "v.initial_long_output":
		return alarm.Value.InitialLongOutput, true
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
	case "v.last_comment.initiator":
		if alarm.Value.LastComment == nil {
			return "", true
		}
		return alarm.Value.LastComment.Initiator, true
	case "v.ticket.m":
		if alarm.Value.Ticket == nil {
			return "", true
		}

		return alarm.Value.Ticket.Message, true
	case "v.ticket.ticket":
		if alarm.Value.Ticket == nil {
			return "", true
		}

		return alarm.Value.Ticket.Ticket, true
	case "v.ticket.initiator":
		if alarm.Value.Ticket == nil {
			return "", true
		}

		return alarm.Value.Ticket.Initiator, true
	case "v.ack.a":
		if alarm.Value.ACK == nil {
			return "", true
		}

		return alarm.Value.ACK.Author, true
	case "v.ack.m":
		if alarm.Value.ACK == nil {
			return "", true
		}

		return alarm.Value.ACK.Message, true
	case "v.ack.initiator":
		if alarm.Value.ACK == nil {
			return "", true
		}

		return alarm.Value.ACK.Initiator, true
	case "v.canceled.initiator":
		if alarm.Value.Canceled == nil {
			return "", true
		}

		return alarm.Value.Canceled.Initiator, true
	default:
		if n := strings.TrimPrefix(f, "v.ticket.ticket_data."); n != f {
			if alarm.Value.Ticket == nil || alarm.Value.Ticket.TicketData == nil {
				return "", true
			}

			return alarm.Value.Ticket.TicketData[n], true
		}

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
	case "v.total_state_changes":
		return int64(alarm.Value.TotalStateChanges), true
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
	case "v.activation_date":
		if alarm.Value.ActivationDate == nil {
			return nil, true
		}
		return alarm.Value.ActivationDate, true
	case "v.change_state":
		if alarm.Value.ChangeState == nil {
			return nil, true
		}
		return alarm.Value.ChangeState, true
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
	case "v.activation_date":
		if alarm.Value.ActivationDate != nil {
			return alarm.Value.ActivationDate.Time, true
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

func getAlarmStringArrayField(alarm types.Alarm, field string) ([]string, bool) {
	switch field {
	case "tags":
		return alarm.Tags, true
	default:
		return nil, false
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

func isForbiddenAlarmField(condition FieldCondition, forbiddenFieldsMap map[string]bool, timeAbsoluteFieldsMap map[string]bool) bool {
	return forbiddenFieldsMap[condition.Field] ||
		forbiddenFieldsMap["v.infos"] && strings.HasPrefix(condition.Field, "v.infos") ||
		timeAbsoluteFieldsMap[condition.Field] && condition.Condition.Type == ConditionTimeRelative
}
