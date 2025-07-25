package event

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters/calculator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
)

func newCountersHelper(
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator,
	componentCountersCalculator calculator.ComponentCountersCalculator,
	eventsSender entitycounters.EventsSender,
	logger zerolog.Logger,
) *countersHelper {
	return &countersHelper{
		entityServiceCountersCalculator: entityServiceCountersCalculator,
		componentCountersCalculator:     componentCountersCalculator,
		eventsSender:                    eventsSender,
		logger:                          logger,
	}
}

// countersHelper is used to update dependency counters for services and components.
type countersHelper struct {
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator
	componentCountersCalculator     calculator.ComponentCountersCalculator
	eventsSender                    entitycounters.EventsSender
	logger                          zerolog.Logger
}

// countersHelper contains updated dependency counters for services and components.
type countersResult struct {
	UpdatedServiceStates    map[string]entitycounters.UpdatedServicesInfo
	IsComponentStateChanged bool
	ComponentID             string
	NewComponentState       int
}

func (h *countersHelper) CalculateCounters(
	ctx context.Context,
	alarm *types.Alarm,
	entity *types.Entity,
	alarmChange types.AlarmChange,
) (bool, countersResult, error) {
	res := countersResult{}
	serviceCountersUpdated, updatedServiceStates, err := h.entityServiceCountersCalculator.CalculateCounters(ctx, entity, alarm, alarmChange)
	if err != nil {
		return false, res, err
	}

	var componentCountersUpdated bool
	var componentStateChanged bool
	var newComponentState int
	if entity.Type == types.EntityTypeResource {
		componentCountersUpdated, componentStateChanged, newComponentState, err = h.componentCountersCalculator.CalculateCounters(ctx, entity, alarm, alarmChange)
		if err != nil {
			return false, res, err
		}
	}

	res.UpdatedServiceStates = updatedServiceStates
	res.IsComponentStateChanged = componentStateChanged
	res.ComponentID = entity.Component
	res.NewComponentState = newComponentState

	return serviceCountersUpdated || componentCountersUpdated, res, nil
}

func (h *countersHelper) UpdateStates(
	ctx context.Context,
	res countersResult,
) {
	for servID, servInfo := range res.UpdatedServiceStates {
		err := h.eventsSender.UpdateServiceState(ctx, servID, servInfo)
		if err != nil {
			h.logger.Err(err).Msg("failed to update service state")
		}
	}

	if res.IsComponentStateChanged {
		err := h.eventsSender.UpdateComponentState(ctx, res.ComponentID, res.NewComponentState)
		if err != nil {
			h.logger.Err(err).Msg("failed to update component state")
		}
	}
}
