package calculator

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type ComponentCountersCalculator interface {
	// CalculateCounters calculates counters for the component by dependency alarm change and returns true if component state is changed as the first argument, and a new component state as the second one.
	CalculateCounters(
		ctx context.Context,
		entity *types.Entity,
		alarm *types.Alarm,
		alarmChange types.AlarmChange,
	) (isCountersUpdated bool, isStateUpdated bool, newState int, _ error)
	// RecomputeCounters recomputes all counters for component's dependencies, return a new calculated component state.
	RecomputeCounters(ctx context.Context, component *types.Entity) (int, error)
	// RecomputeAll sends entityupdated events for all components with state settings, which should recompute counters.
	RecomputeAll(ctx context.Context) error
}

type EntityServiceCountersCalculator interface {
	// CalculateCounters calculates counters for the entity service by dependency alarm change and returns a map with updated entity service ids and a new calculated state and output.
	CalculateCounters(
		ctx context.Context,
		entity *types.Entity,
		alarm *types.Alarm,
		alarmChange types.AlarmChange,
	) (isAnyServiceCountersUpdated bool, updatedServicesInfos map[string]entitycounters.UpdatedServicesInfo, _ error)
	// RecomputeCounters recomputes all counters for service's dependencies, returns a map with updated entity service id and a new calculated state and output.
	RecomputeCounters(ctx context.Context, service *types.Entity) (map[string]entitycounters.UpdatedServicesInfo, error)
	// RecomputeAll sends recomputeentityservice events for all entity services, which should recompute counters.
	RecomputeAll(ctx context.Context) error
}

type ComponentCountersStrategy interface {
	// CanSkip defines conditions where we can say that counters won't be changed and we can skip calculation
	CanSkip(calcData entitycounters.ComponentCountersCalcData) bool
	// Calculate contains counters calculation logic.
	Calculate(calcData entitycounters.ComponentCountersCalcData) entitycounters.EntityCounters
}

type EntityServiceCountersStrategy interface {
	// CanSkip defines conditions where we can say that counters won't be changed and we can skip calculation
	CanSkip(calcData entitycounters.EntityServiceCountersCalcData) bool
	// Calculate contains counters calculation logic.
	Calculate(calcData entitycounters.EntityServiceCountersCalcData) entitycounters.EntityCounters
}
