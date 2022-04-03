package statecounters

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type StateCountersService interface {
	RecomputeEntityServiceCounters(ctx context.Context, event types.Event) error
	UpdateServiceCounters(ctx context.Context, entity types.Entity, alarm *types.Alarm, alarmChange types.AlarmChange) (map[string]int, error)
	UpdateServiceState(serviceID string, serviceState int) error
}
