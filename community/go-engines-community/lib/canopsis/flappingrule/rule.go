package flappingrule

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/oldpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type Rule struct {
	ID                string                       `bson:"_id"`
	Name              string                       `bson:"name"`
	Description       string                       `bson:"description"`
	FreqLimit         int                          `bson:"freq_limit"`
	Duration          types.DurationWithUnit       `bson:"duration"`
	OldAlarmPatterns  oldpattern.AlarmPatternList  `bson:"old_alarm_patterns,omitempty"`
	OldEntityPatterns oldpattern.EntityPatternList `bson:"old_entity_patterns,omitempty"`
	Priority          int                          `bson:"priority"`
	Author            string                       `bson:"author"`
	Created           types.CpsTime                `bson:"created"`
	Updated           types.CpsTime                `bson:"updated"`

	savedpattern.EntityPatternFields `bson:",inline"`
	savedpattern.AlarmPatternFields  `bson:",inline"`
}

// Matches returns true if alarm and entity match field patterns.
func (r *Rule) Matches(alarmWithEntity types.AlarmWithEntity) (bool, error) {
	return pattern.Match(alarmWithEntity.Entity, alarmWithEntity.Alarm, r.EntityPattern, r.AlarmPattern, r.OldEntityPatterns, r.OldAlarmPatterns)
}
