package executor

import (
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/config"
	operationlib "git.canopsis.net/canopsis/go-engines/lib/canopsis/operation"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/utils"
)

// NewCancelExecutor creates new executor.
func NewCancelExecutor(cfg config.CanopsisConf) operationlib.Executor {
	return &cancelExecutor{cfg: cfg}
}

type cancelExecutor struct {
	cfg config.CanopsisConf
}

// Exec creates new cancel step for alarm and update alarm status.
func (e *cancelExecutor) Exec(
	operation types.Operation,
	alarm *types.Alarm,
	time types.CpsTime,
	role, initiator string,
) (types.AlarmChangeType, error) {
	var params types.OperationParameters
	var ok bool
	if params, ok = operation.Parameters.(types.OperationParameters); !ok {
		return "", fmt.Errorf("invalid parameters")
	}

	err := alarm.PartialUpdateCancel(
		time,
		params.Author,
		utils.TruncateString(params.Output, e.cfg.Alarm.OutputLength),
		role,
		initiator,
		e.cfg,
	)
	if err != nil {
		return "", err
	}

	return types.AlarmChangeTypeCancel, nil
}
