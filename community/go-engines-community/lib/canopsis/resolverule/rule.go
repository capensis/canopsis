package resolverule

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/oldpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

const DefaultRule = "default_rule"

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

	return pattern.Match(alarmWithEntity.Entity, alarmWithEntity.Alarm, r.EntityPattern, r.AlarmPattern, r.OldEntityPatterns, r.OldAlarmPatterns)
}
