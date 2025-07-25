package event

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters/calculator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
)

func newComponentAndServiceCountersHelper(
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator,
	componentCountersCalculator calculator.ComponentCountersCalculator,
	eventsSender entitycounters.EventsSender,
	logger zerolog.Logger,
) *componentAndServiceCountersHelper {
	return &componentAndServiceCountersHelper{
		entityServiceCountersCalculator: entityServiceCountersCalculator,
		componentCountersCalculator:     componentCountersCalculator,
		eventsSender:                    eventsSender,
		logger:                          logger,
	}
}

type componentAndServiceCountersHelper struct {
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator
	componentCountersCalculator     calculator.ComponentCountersCalculator
	eventsSender                    entitycounters.EventsSender
	logger                          zerolog.Logger
}

type componentAndServiceCountersResult struct {
	UpdatedServiceStates    map[string]entitycounters.UpdatedServicesInfo
	IsComponentStateChanged bool
	ComponentID             string
	NewComponentState       int
}

func (h *componentAndServiceCountersHelper) Process(
	ctx context.Context,
	alarm *types.Alarm,
	entity *types.Entity,
	alarmChange types.AlarmChange,
) (bool, componentAndServiceCountersResult, error) {
	res := componentAndServiceCountersResult{}
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

func (h *componentAndServiceCountersHelper) PostProcess(
	ctx context.Context,
	res componentAndServiceCountersResult,
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
