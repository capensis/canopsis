package event

import (
	"context"
	"maps"

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

// CalculateCountersWithPauseChange calculates the counters of the event and, when the entity has been
// paused or resumed by the event itself, see updatePauseOnServicesChange, the counters of this pause
// change: the event is counted with the pause the entity had before, which is what the counters of its
// services hold, then the pause or the resume is applied to all of them.
func (h *countersHelper) CalculateCountersWithPauseChange(
	ctx context.Context,
	alarm *types.Alarm,
	entity *types.Entity,
	alarmChange types.AlarmChange,
	change pauseChange,
) (bool, countersResult, error) {
	if change.alarmChange != types.AlarmChangeTypePause && change.alarmChange != types.AlarmChangeTypeResume {
		return h.CalculateCounters(ctx, alarm, entity, alarmChange)
	}

	entityBefore := *entity
	entityBefore.Pause = nil
	if change.alarmChange == types.AlarmChangeTypeResume {
		// The entity was paused by the pause which has just ended.
		ended := types.Pause{}
		if entity.LatestPause != nil {
			ended = entity.LatestPause.Pause
		}

		entityBefore.Pause = &ended
	}

	updated, res, err := h.CalculateCounters(ctx, alarm, &entityBefore, alarmChange)
	if err != nil {
		return false, res, err
	}

	// Only the pause of the entity changes, its pbehavior is the same before and after it.
	pauseAlarmChange := types.NewAlarmChange()
	pauseAlarmChange.Type = change.alarmChange
	pauseAlarmChange.PreviousPbehaviorTypeID = entity.PbehaviorInfo.TypeID
	pauseAlarmChange.PreviousPbehaviorCannonicalType = entity.PbehaviorInfo.CanonicalType

	pauseUpdated, pauseRes, err := h.CalculateCounters(ctx, alarm, entity, pauseAlarmChange)
	if err != nil {
		return false, res, err
	}

	if pauseRes.UpdatedServiceStates != nil {
		if res.UpdatedServiceStates == nil {
			res.UpdatedServiceStates = pauseRes.UpdatedServiceStates
		} else {
			// The states of the pause change are the newest ones.
			maps.Copy(res.UpdatedServiceStates, pauseRes.UpdatedServiceStates)
		}
	}

	if pauseRes.IsComponentStateChanged {
		res.IsComponentStateChanged = true
		res.ComponentID = pauseRes.ComponentID
		res.NewComponentState = pauseRes.NewComponentState
	}

	return updated || pauseUpdated, res, nil
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
