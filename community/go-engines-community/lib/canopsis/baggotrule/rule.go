package baggotrule

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type Rule struct {
	ID             string                    `bson:"_id"`
	Duration       types.DurationWithUnit    `bson:"duration"`
	AlarmPatterns  pattern.AlarmPatternList  `bson:"alarm_patterns"`
	EntityPatterns pattern.EntityPatternList `bson:"entity_patterns"`
	Updated        *types.CpsTime            `bson:"updated"`
	Priority       int                       `bson:"priority"`
}

// Matches returns true if alrs and entity match time condition and field patterns.
func (r *Rule) Matches(alarm *types.AlarmWithEntity) bool {
	return r.AlarmPatterns.Matches(&alarm.Alarm) && r.EntityPatterns.Matches(&alarm.Entity)
}
