package metaalarm

import "git.canopsis.net/canopsis/go-engines/lib/canopsis/eventfilter/pattern"

type RuleConfig struct {
	TimeInterval int `bson:"time_interval,omitempty" json:"time_interval,omitempty"`
	// AlarmPatterns represents Alarm's attribute pattern list
	AlarmPatterns pattern.AlarmPatternList `bson:"alarm_patterns,omitempty" json:"alarm_patterns,omitempty"`
	// EntityPatterns represents Entity's attribute pattern list
	EntityPatterns pattern.EntityPatternList `bson:"entity_patterns,omitempty" json:"entity_patterns,omitempty"`
	// TotalEntityPatterns represents Entity's attribute pattern list to find total entity count.
	TotalEntityPatterns *pattern.EntityPatternList `bson:"total_entity_patterns,omitempty" json:"total_entity_patterns,omitempty"`
	// EventPatterns represents Events's attribute pattern list
	EventPatterns pattern.EventPatternList `bson:"event_patterns,omitempty" json:"event_patterns,omitempty"`
	// ThresholdRate is malfunctioning entities rate threshold to trigger the rule
	ThresholdRate *float64 `bson:"threshold_rate,omitempty" json:"threshold_rate,omitempty"`
	// ThresholdCount is malfunctioning entities count threshold to trigger the rule
	ThresholdCount *int64 `bson:"threshold_count,omitempty" json:"threshold_count,omitempty"`
	//
	ValuePaths []string `bson:"value_paths,omitempty" json:"value_paths,omitempty"`
}
