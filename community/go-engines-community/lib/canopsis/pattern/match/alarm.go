package match

import (
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func MatchAlarmPattern(p pattern.Alarm, alarm *types.Alarm) (bool, error) {
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

			if infoName := pattern.GetAlarmInfoName(f); infoName != "" {
				infoVal, ok := alarm.GetInfoVal(infoName)
				if v.FieldType == "" {
					matched, err = cond.MatchRef(infoVal)
				} else if ok {
					switch v.FieldType {
					case pattern.FieldTypeString:
						var s string
						if s, err = pattern.GetStringValue(infoVal); err == nil {
							matched, err = cond.MatchString(s)
						}
					case pattern.FieldTypeInt:
						var i int64
						if i, err = pattern.GetIntValue(infoVal); err == nil {
							matched, err = cond.MatchInt(i)
						}
					case pattern.FieldTypeBool:
						var b bool
						if b, err = pattern.GetBoolValue(infoVal); err == nil {
							matched, err = cond.MatchBool(b)
						}
					case pattern.FieldTypeStringArray:
						var a []string
						if a, err = pattern.GetStringArrayValue(infoVal); err == nil {
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
			if str, ok := alarm.GetStringField(f); ok {
				foundField = true
				matched, err = cond.MatchString(str)
			}
			if !foundField || err != nil {
				if i, ok := alarm.GetIntField(f); ok {
					foundField = true
					matched, err = cond.MatchInt(i)
				}
			}
			if !foundField || err != nil {
				if r, ok := alarm.GetRefField(f); ok {
					foundField = true
					matched, err = cond.MatchRef(r)
				}
			}
			if !foundField || err != nil {
				if t, ok := alarm.GetTimeField(f); ok {
					foundField = true
					matched, err = cond.MatchTime(t)
				}
			}
			if !foundField || err != nil {
				if d, ok := alarm.GetDurationField(f); ok {
					foundField = true
					matched, err = cond.MatchDuration(d)
				}
			}
			if !foundField || err != nil {
				if a, ok := alarm.GetStringArrayField(f); ok {
					foundField = true
					matched, err = cond.MatchStringArray(a)
				}
			}

			if !foundField {
				err = pattern.ErrUnsupportedField
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

func ValidateAlarmPattern(p pattern.Alarm, forbiddenFields, onlyTimeAbsoluteFields []string) bool {
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

			if pattern.IsForbiddenAlarmField(v, forbiddenFieldsMap, timeAbsoluteFieldsMap) {
				return false
			}

			if infoName := pattern.GetAlarmInfoName(f); infoName != "" {
				switch v.FieldType {
				case pattern.FieldTypeString:
					_, err = cond.MatchString("")
				case pattern.FieldTypeInt:
					_, err = cond.MatchInt(0)
				case pattern.FieldTypeBool:
					_, err = cond.MatchBool(false)
				case pattern.FieldTypeStringArray:
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
			if str, ok := emptyAlarm.GetStringField(f); ok {
				foundField = true
				_, err = cond.MatchString(str)
			}
			if !foundField || err != nil {
				if i, ok := emptyAlarm.GetIntField(f); ok {
					foundField = true
					_, err = cond.MatchInt(i)
				}
			}
			if !foundField || err != nil {
				if r, ok := emptyAlarm.GetRefField(f); ok {
					foundField = true
					_, err = cond.MatchRef(r)
				}
			}
			if !foundField || err != nil {
				if t, ok := emptyAlarm.GetTimeField(f); ok {
					foundField = true
					_, err = cond.MatchTime(t)
				}
			}
			if !foundField || err != nil {
				if d, ok := emptyAlarm.GetDurationField(f); ok {
					foundField = true
					_, err = cond.MatchDuration(d)
				}
			}
			if !foundField || err != nil {
				if a, ok := emptyAlarm.GetStringArrayField(f); ok {
					foundField = true
					_, err = cond.MatchStringArray(a)
				}
			}

			if !foundField {
				err = pattern.ErrUnsupportedField
			}

			if err != nil {
				return false
			}
		}
	}

	return true
}
