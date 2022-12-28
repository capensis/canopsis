package entityservice

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/oldpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type EntityService struct {
	types.Entity   `bson:",inline"`
	OutputTemplate string        `bson:"output_template" json:"output_template"`
	Counters       AlarmCounters `bson:"counters" json:"counters"`

	savedpattern.EntityPatternFields `bson:",inline"`
	OldEntityPatterns                oldpattern.EntityPatternList `bson:"old_entity_patterns,omitempty" json:"old_entity_patterns,omitempty"`
}

// GetServiceState returns the state of the service.
func GetServiceState(counters AlarmCounters) int {
	if counters.State.Critical > 0 {
		return types.AlarmStateCritical
	}
	if counters.State.Major > 0 {
		return types.AlarmStateMajor
	}
	if counters.State.Minor > 0 {
		return types.AlarmStateMinor
	}

	return types.AlarmStateOK
}
