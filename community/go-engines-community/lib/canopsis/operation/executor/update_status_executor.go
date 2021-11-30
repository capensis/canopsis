package executor

import (
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	operationlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func NewUpdateStatusExecutor(configProvider config.AlarmConfigProvider) operationlib.Executor {
	return &updateStatusExecutor{configProvider: configProvider}
}

type updateStatusExecutor struct {
	configProvider config.AlarmConfigProvider
}

func (e *updateStatusExecutor) Exec(
	operation types.Operation,
	alarm *types.Alarm,
	time types.CpsTime,
	_, _, _ string,
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
		e.configProvider.Get(),
	)
	if err != nil {
		return "", err
	}

	if oldStatus == alarm.Value.Status {
		return "", nil
	}

	return types.AlarmChangeTypeUpdateStatus, nil
}
