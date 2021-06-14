package watcher

import (
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
)

// DependencyState is a struct containing the informations about a dependency
// that are used to compute a watcher's state.
type DependencyState struct {
	EntityID          string
	ImpactedWatchers  []string
	HasAlarm          bool
	AlarmState        int
	AlarmAcknowledged bool
	LastUpdateDate    time.Time
	IsEntityActive    bool
	PbehaviorType     string
}

// NewDependencyState creates a DependencyState given an entity, an (optional)
// alarm and a list of pbehaviors (active or not) that impact the entity.
// watchers should be a map containing all the watchers, and is used to filter
// the impact of the entity to only keep the watchers.
// FIXME: LastUpdateDate is used to ensure that only the most recent
// DependencyStates are taken into account. This is necessary because the
// ComputeAllWatchers method may take a few seconds to run, which could easily
// lead to a race condition. As of this writing, LastUpdateDate is set to the
// time the data enters the engine (the time the event was received, or the
// time the MongoDB query was started). This may break if several instances of
// the engine are run in parallel.
func NewDependencyState(
	entity types.Entity,
	alarm *types.Alarm,
	isEntityActive bool,
	pbehaviorType string,
	watchers map[string]Watcher,
	date time.Time,
) DependencyState {
	var impactedWatchers []string
	for _, impact := range entity.Impacts {
		_, isWatcher := watchers[impact]
		if isWatcher {
			impactedWatchers = append(impactedWatchers, impact)
		}
	}

	if alarm == nil || alarm.Value.Resolved != nil && alarm.Value.Resolved.Unix() > 0 {
		return DependencyState{
			EntityID:         entity.ID,
			ImpactedWatchers: impactedWatchers,
			LastUpdateDate:   date,
			IsEntityActive:   isEntityActive,
			PbehaviorType:    pbehaviorType,
		}
	}

	return DependencyState{
		EntityID:          entity.ID,
		ImpactedWatchers:  impactedWatchers,
		HasAlarm:          true,
		AlarmState:        int(alarm.Value.State.Value),
		AlarmAcknowledged: alarm.Value.ACK != nil,
		LastUpdateDate:    date,
		IsEntityActive:    isEntityActive,
		PbehaviorType:     pbehaviorType,
	}
}
