package executor

import (
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	operationlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
)

// NewCancelExecutor creates new executor.
func NewCancelExecutor(configProvider config.AlarmConfigProvider) operationlib.Executor {
	return &cancelExecutor{configProvider: configProvider}
}

type cancelExecutor struct {
	configProvider config.AlarmConfigProvider
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

	alarmConfig := e.configProvider.Get()
	err := alarm.PartialUpdateCancel(
		time,
		params.Author,
		utils.TruncateString(params.Output, alarmConfig.OutputLength),
		role,
		initiator,
		alarmConfig,
	)
	if err != nil {
		return "", err
	}

	return types.AlarmChangeTypeCancel, nil
}
