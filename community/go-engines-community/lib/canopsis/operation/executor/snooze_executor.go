package executor

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	operationlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
)

// NewSnoozeExecutor creates new executor.
func NewSnoozeExecutor(configProvider config.AlarmConfigProvider) operationlib.Executor {
	return &snoozeExecutor{configProvider: configProvider}
}

type snoozeExecutor struct {
	configProvider config.AlarmConfigProvider
}

// Exec creates new snooze step for alarm.
func (e *snoozeExecutor) Exec(
	_ context.Context,
	operation types.Operation,
	alarm *types.Alarm,
	_ types.Entity,
	time types.CpsTime,
	userID, role, initiator string,
) (types.AlarmChangeType, error) {
	var params types.OperationSnoozeParameters
	var ok bool
	if params, ok = operation.Parameters.(types.OperationSnoozeParameters); !ok {
		return "", fmt.Errorf("invalid parameters")
	}

	if alarm.Value.Snooze != nil {
		return "", nil
	}

	err := alarm.PartialUpdateSnooze(
		time,
		types.CpsNumber(params.Duration.Seconds),
		params.Author,
		utils.TruncateString(params.Output, e.configProvider.Get().OutputLength),
		userID,
		role,
		initiator,
	)
	if err != nil {
		return "", err
	}

	return types.AlarmChangeTypeSnooze, nil
}
