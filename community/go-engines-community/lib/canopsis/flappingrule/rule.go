package flappingrule

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type Rule struct {
	ID                string                   `bson:"_id"`
	FlappingInterval  types.DurationWithUnit   `bson:"flapping_interval"`
	FlappingFreqLimit int                      `bson:"flapping_freq_limit"`
	AlarmPatterns     pattern.AlarmPatternList `bson:"alarm_patterns"`
	Updated           *types.CpsTime           `bson:"updated"`
	Priority          int                      `bson:"priority"`
}

// Matches returns true if alrs and entity match time condition and field patterns.
func (r *Rule) Matches(alarm *types.Alarm) bool {
	return r.AlarmPatterns.Matches(alarm)
}

// LoadRulesReport is a struct containing the rules that have been added,
// modified, or removed between two calls of the Service.LoadRules method.
type LoadRulesReport struct {
	Unchanged  []Rule
	Added      []Rule
	Modified   []Rule
	Removed    []Rule
	FailParsed []Rule
}
