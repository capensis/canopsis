package executor

import (
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/config"
	operationlib "git.canopsis.net/canopsis/go-engines/lib/canopsis/operation"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/utils"
)

// NewAckRemoveExecutor creates new executor.
func NewAckRemoveExecutor(configProvider config.AlarmConfigProvider) operationlib.Executor {
	return &ackRemoveExecutor{configProvider: configProvider}
}

type ackRemoveExecutor struct {
	configProvider config.AlarmConfigProvider
}

// Exec creates new ack remove step for alarm.
func (e *ackRemoveExecutor) Exec(
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

	if alarm.Value.ACK == nil {
		return "", nil
	}

	err := alarm.PartialUpdateUnack(
		time,
		params.Author,
		utils.TruncateString(params.Output, e.configProvider.Get().OutputLength),
		role,
		initiator,
	)
	if err != nil {
		return "", err
	}

	return types.AlarmChangeTypeAckremove, nil
}
