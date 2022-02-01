package correlation

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type RuleConfig struct {
	TimeInterval *types.DurationWithUnit `bson:"time_interval,omitempty" json:"time_interval,omitempty"`
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
	// CorelID is an AlarmWithEntity template to mark that alarms are correlated
	CorelID string `bson:"corel_id,omitempty" json:"corel_id,omitempty"`
	// CorelStatus is an AlarmWithEntity template to get correlation relation value
	CorelStatus string `bson:"corel_status,omitempty" json:"corel_status,omitempty"`
	// CorelParent is the correlation relation value, which mark alarm as a parent
	CorelParent string `bson:"corel_parent,omitempty" json:"corel_parent,omitempty"`
	// CorelChild is the correlation relation value, which mark alarm as a child
	CorelChild string `bson:"corel_child,omitempty" json:"corel_child,omitempty"`
}
