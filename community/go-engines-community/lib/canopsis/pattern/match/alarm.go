package match

import (
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func MatchAlarmPattern(p pattern.Alarm, alarm *types.Alarm) (bool, error) {
	if len(p) == 0 {
		return true, nil
	}

	now := time.Now() // to compute relative time values

	for idx := range p {
		matched := false

		for _, v := range p[idx] {
			f := v.Field

			var err error
			matched = false

			if infoName := pattern.GetAlarmInfoName(f); infoName != "" {
				infoVal, infoExists := alarm.GetInfoVal(infoName)

				matched, err = v.MatchAlarmInfoCondition(infoVal, infoExists)
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
				matched, err = v.Condition.MatchString(str)
			}
			if !foundField || err != nil {
				if i, ok := alarm.GetIntField(f); ok {
					foundField = true
					matched, err = v.Condition.MatchInt(i)
				}
			}
			if !foundField || err != nil {
				if r, ok := alarm.GetRefField(f); ok {
					foundField = true
					matched, err = v.Condition.MatchRef(r)
				}
			}
			if !foundField || err != nil {
				if t, ok := alarm.GetTimeField(f); ok {
					foundField = true
					matched, err = v.Condition.MatchTime(t, now)
				}
			}
			if !foundField || err != nil {
				if d, ok := alarm.GetDurationField(f); ok {
					foundField = true
					matched, err = v.Condition.MatchDuration(d)
				}
			}
			if !foundField || err != nil {
				if a, ok := alarm.GetTagsField(f); ok {
					foundField = true
					matched, err = v.Condition.MatchTags(a)
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
	return len(AlarmPatternErrors(p, forbiddenFields, onlyTimeAbsoluteFields)) == 0
}

func AlarmPatternErrors(p pattern.Alarm, forbiddenFields, onlyTimeAbsoluteFields []string) []ConditionError {
	emptyAlarm := types.Alarm{}
	forbiddenFieldsMap := make(map[string]bool, len(forbiddenFields))
	for _, field := range forbiddenFields {
		forbiddenFieldsMap[field] = true
	}

	timeAbsoluteFieldsMap := make(map[string]bool, len(onlyTimeAbsoluteFields))
	for _, field := range onlyTimeAbsoluteFields {
		timeAbsoluteFieldsMap[field] = true
	}

	now := time.Now() // to compute relative time values
	var errs []ConditionError

	for gidx := range p {
		if len(p[gidx]) == 0 {
			errs = append(errs, ConditionError{GroupIdx: gidx, CondIdx: -1, Err: pattern.ErrEmptyGroup})
			continue
		}

		for cidx, v := range p[gidx] {
			f := v.Field

			if pattern.IsForbiddenAlarmField(v, forbiddenFieldsMap, timeAbsoluteFieldsMap) {
				errs = append(errs, ConditionError{GroupIdx: gidx, CondIdx: cidx, Err: pattern.ErrForbiddenField})
				continue
			}

			if infoName := pattern.GetAlarmInfoName(f); infoName != "" {
				if err := v.InfoConditionError(); err != nil {
					errs = append(errs, ConditionError{GroupIdx: gidx, CondIdx: cidx, Err: err})
				}

				continue
			}

			var err error
			foundField := false
			if str, ok := emptyAlarm.GetStringField(f); ok {
				foundField = true
				_, err = v.Condition.MatchString(str)
			}
			if !foundField || err != nil {
				if i, ok := emptyAlarm.GetIntField(f); ok {
					foundField = true
					_, err = v.Condition.MatchInt(i)
				}
			}
			if !foundField || err != nil {
				if r, ok := emptyAlarm.GetRefField(f); ok {
					foundField = true
					_, err = v.Condition.MatchRef(r)
				}
			}
			if !foundField || err != nil {
				if t, ok := emptyAlarm.GetTimeField(f); ok {
					foundField = true
					_, err = v.Condition.MatchTime(t, now)
				}
			}
			if !foundField || err != nil {
				if d, ok := emptyAlarm.GetDurationField(f); ok {
					foundField = true
					_, err = v.Condition.MatchDuration(d)
				}
			}
			if !foundField || err != nil {
				if a, ok := emptyAlarm.GetTagsField(f); ok {
					foundField = true
					_, err = v.Condition.MatchTags(a)
				}
			}

			if !foundField {
				err = pattern.ErrUnsupportedField
			}

			if err != nil {
				errs = append(errs, ConditionError{GroupIdx: gidx, CondIdx: cidx, Err: err})
			}
		}
	}

	return errs
}
