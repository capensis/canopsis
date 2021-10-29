package flappingrule

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type Rule struct {
	ID             string                    `bson:"_id"`
	Name           string                    `bson:"name"`
	Description    string                    `bson:"description"`
	FreqLimit      int                       `bson:"freq_limit"`
	Duration       types.DurationWithUnit    `bson:"duration"`
	AlarmPatterns  pattern.AlarmPatternList  `bson:"alarm_patterns"`
	EntityPatterns pattern.EntityPatternList `bson:"entity_patterns"`
	Priority       int                       `bson:"priority"`
	Author         string                    `bson:"author"`
	Created        types.CpsTime             `bson:"created"`
	Updated        types.CpsTime             `bson:"updated"`
}

// Matches returns true if alarm and entity match field patterns.
func (r *Rule) Matches(alarm types.AlarmWithEntity) bool {
	return r.AlarmPatterns.Matches(&alarm.Alarm) && r.EntityPatterns.Matches(&alarm.Entity)
}
