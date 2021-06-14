package executor

import (
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/config"
	operationlib "git.canopsis.net/canopsis/go-engines/lib/canopsis/operation"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/utils"
)

// NewSnoozeExecutor creates new executor.
func NewSnoozeExecutor(cfg config.CanopsisConf) operationlib.Executor {
	return &snoozeExecutor{cfg: cfg}
}

type snoozeExecutor struct {
	cfg config.CanopsisConf
}

// Exec creates new snooze step for alarm.
func (e *snoozeExecutor) Exec(
	operation types.Operation,
	alarm *types.Alarm,
	time types.CpsTime,
	role, initiator string,
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
		utils.TruncateString(params.Output, e.cfg.Alarm.OutputLength),
		role,
		initiator,
	)
	if err != nil {
		return "", err
	}

	return types.AlarmChangeTypeSnooze, nil
}
