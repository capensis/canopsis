package resolverule

import (
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/oldpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

const DefaultRule = "default-rule"

type Rule struct {
	ID                string                       `bson:"_id,omitempty"`
	Name              string                       `bson:"name"`
	Description       string                       `bson:"description"`
	Duration          types.DurationWithUnit       `bson:"duration"`
	OldAlarmPatterns  oldpattern.AlarmPatternList  `bson:"old_alarm_patterns,omitempty"`
	OldEntityPatterns oldpattern.EntityPatternList `bson:"old_entity_patterns,omitempty"`
	Priority          int                          `bson:"priority"`
	Author            string                       `bson:"author"`
	Created           types.CpsTime                `bson:"created,omitempty"`
	Updated           types.CpsTime                `bson:"updated,omitempty"`

	savedpattern.EntityPatternFields `bson:",inline"`
	savedpattern.AlarmPatternFields  `bson:",inline"`
}

// Matches returns true if alarm and entity match patterns.
func (r *Rule) Matches(alarmWithEntity types.AlarmWithEntity) (bool, error) {
	if r.ID == DefaultRule {
		return true, nil
	}

	if !r.OldAlarmPatterns.IsSet() && !r.OldEntityPatterns.IsSet() &&
		len(r.EntityPattern) == 0 && len(r.AlarmPattern) == 0 {
		return false, nil
	}

	var matched bool
	var err error

	if r.OldAlarmPatterns.IsSet() {
		if !r.OldAlarmPatterns.IsValid() {
			return false, pattern.InvalidOldAlarmPattern
		}

		matched = r.OldAlarmPatterns.Matches(&alarmWithEntity.Alarm)
	} else {
		matched, err = r.AlarmPattern.Match(alarmWithEntity.Alarm)
		if err != nil {
			return false, fmt.Errorf("resolve rule has an invalid alarm pattern : %w", err)
		}
	}

	if !matched {
		return false, nil
	}

	if r.OldEntityPatterns.IsSet() {
		if !r.OldEntityPatterns.IsValid() {
			return false, pattern.InvalidOldEntityPattern
		}

		matched = r.OldEntityPatterns.Matches(&alarmWithEntity.Entity)
	} else {
		matched, _, err = r.EntityPattern.Match(alarmWithEntity.Entity)
		if err != nil {
			return false, fmt.Errorf("resolve rule has an invalid entity pattern : %w", err)
		}
	}

	return matched, nil
}
