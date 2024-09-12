package event

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func RemoveMetaAlarmState(
	ctx context.Context,
	metaAlarm types.Alarm,
	rule correlation.Rule,
	metaAlarmStatesService correlation.MetaAlarmStateService,
) error {
	if rule.IsManual() {
		return nil
	}

	stateID := rule.GetStateID(metaAlarm.Value.MetaValuePath)
	metaAlarmState, err := metaAlarmStatesService.GetMetaAlarmState(ctx, stateID)
	if err != nil {
		return fmt.Errorf("cannot get meta alarm state: %w", err)
	}

	if metaAlarmState.ID == "" {
		return nil
	}

	_, err = metaAlarmStatesService.ArchiveState(ctx, metaAlarmState)
	if err != nil {
		return fmt.Errorf("cannot archive meta alarm state: %w", err)
	}

	_, err = metaAlarmStatesService.DeleteState(ctx, stateID)
	if err != nil {
		return fmt.Errorf("cannot delete meta alarm state: %w", err)
	}

	return nil
}
