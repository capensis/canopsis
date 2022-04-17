package statecounters

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type StateCountersService interface {
	// RecomputeEntityServiceCounters accepts recomputecounters event and computes service's counters from scratch
	RecomputeEntityServiceCounters(ctx context.Context, event types.Event) (map[string]UpdatedServicesInfo, error)
	// UpdateServiceCounters updates services' counters regarding the alarmChange value
	UpdateServiceCounters(ctx context.Context, entity types.Entity, alarm *types.Alarm, alarmChange types.AlarmChange) (map[string]UpdatedServicesInfo, error)
	// UpdateServiceState sends check service event with updated output and state from serviceInfo
	UpdateServiceState(serviceID string, serviceInfo UpdatedServicesInfo) error
	// RecomputeAllServices sends recomputecounters events for all enabled entity services
	RecomputeAllServices(ctx context.Context) error
}
