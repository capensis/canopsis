package resolverule

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/oldpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type Rule struct {
	ID             string                       `bson:"_id"`
	Name           string                       `bson:"name"`
	Description    string                       `bson:"description"`
	Duration       types.DurationWithUnit       `bson:"duration"`
	AlarmPatterns  oldpattern.AlarmPatternList  `bson:"alarm_patterns"`
	EntityPatterns oldpattern.EntityPatternList `bson:"entity_patterns"`
	Priority       int                          `bson:"priority"`
	Author         string                       `bson:"author"`
	Created        types.CpsTime                `bson:"created"`
	Updated        types.CpsTime                `bson:"updated"`
}

// Matches returns true if alarm and entity match patterns.
func (r *Rule) Matches(alarm types.AlarmWithEntity) bool {
	return r.AlarmPatterns.Matches(&alarm.Alarm) && r.EntityPatterns.Matches(&alarm.Entity)
}
