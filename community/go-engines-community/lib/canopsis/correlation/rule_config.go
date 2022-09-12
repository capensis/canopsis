package correlation

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type RuleConfig struct {
	TimeInterval *types.DurationWithUnit `bson:"time_interval,omitempty" json:"time_interval,omitempty"`
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
