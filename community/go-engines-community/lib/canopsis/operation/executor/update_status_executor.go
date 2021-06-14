package executor

import (
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/config"
	operationlib "git.canopsis.net/canopsis/go-engines/lib/canopsis/operation"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
)

func NewUpdateStatusExecutor(cfg config.CanopsisConf) operationlib.Executor {
	return &updateStatusExecutor{cfg: cfg}
}

type updateStatusExecutor struct {
	cfg config.CanopsisConf
}

func (e *updateStatusExecutor) Exec(
	operation types.Operation,
	alarm *types.Alarm,
	time types.CpsTime,
	_, _ string,
) (types.AlarmChangeType, error) {
	var params types.OperationParameters
	var ok bool
	if params, ok = operation.Parameters.(types.OperationParameters); !ok {
		return "", fmt.Errorf("invalid parameters")
	}

	oldStatus := alarm.Value.Status
	err := alarm.PartialUpdateStatus(
		time,
		params.Output,
		e.cfg,
	)
	if err != nil {
		return "", err
	}

	if oldStatus == alarm.Value.Status {
		return "", nil
	}

	return types.AlarmChangeTypeUpdateStatus, nil
}
