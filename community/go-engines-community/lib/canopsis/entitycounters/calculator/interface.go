package calculator

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type ComponentCountersCalculator interface {
	CalculateCounters(ctx context.Context, entity *types.Entity, alarm *types.Alarm, alarmChange types.AlarmChange) (bool, int, error)
	RecomputeCounters(ctx context.Context, entity *types.Entity) (int, error)
	RecomputeAll(ctx context.Context) error
}

type EntityServiceCountersCalculator interface {
	CalculateCounters(ctx context.Context, entity *types.Entity, alarm *types.Alarm, alarmChange types.AlarmChange) (map[string]entitycounters.UpdatedServicesInfo, error)
	RecomputeCounters(ctx context.Context, entity *types.Entity) (map[string]entitycounters.UpdatedServicesInfo, error)
	RecomputeAll(ctx context.Context) error
}

type ComponentCountersStrategy interface {
	CanSkip(calcData entitycounters.ComponentCountersCalcData) bool
	Calculate(calcData entitycounters.ComponentCountersCalcData) entitycounters.EntityCounters
}

type EntityServiceCountersStrategy interface {
	CanSkip(calcData entitycounters.EntityServiceCountersCalcData) bool
	Calculate(calcData entitycounters.EntityServiceCountersCalcData) entitycounters.EntityCounters
}
