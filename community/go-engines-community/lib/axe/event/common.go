package event

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters/calculator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func processComponentAndServiceCounters(
	ctx context.Context,
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator,
	componentCountersCalculator calculator.ComponentCountersCalculator,
	alarm *types.Alarm,
	entity *types.Entity,
	alarmChange types.AlarmChange,
) (map[string]entitycounters.UpdatedServicesInfo, bool, int, error) {
	updatedServiceStates, err := entityServiceCountersCalculator.CalculateCounters(ctx, entity, alarm, alarmChange)
	if err != nil {
		return nil, false, 0, err
	}

	var componentStateChanged bool
	var newComponentState int

	if entity.Type == types.EntityTypeResource {
		componentStateChanged, newComponentState, err = componentCountersCalculator.CalculateCounters(ctx, entity, alarm, alarmChange)
		if err != nil {
			return nil, false, 0, err
		}
	}

	return updatedServiceStates, componentStateChanged, newComponentState, nil
}
