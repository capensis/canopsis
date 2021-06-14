/*
Package executor contains operation executors.
*/
package executor

import (
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/config"
	operationlib "git.canopsis.net/canopsis/go-engines/lib/canopsis/operation"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/utils"
)

// NewAckExecutor creates new executor.
func NewAckExecutor(cfg config.CanopsisConf) operationlib.Executor {
	return &ackExecutor{cfg: cfg}
}

type ackExecutor struct {
	cfg config.CanopsisConf
}

// Exec creates new ack step for alarm.
func (e *ackExecutor) Exec(
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

	if alarm.Value.ACK != nil {
		return "", nil
	}

	err := alarm.PartialUpdateAck(
		time,
		params.Author,
		utils.TruncateString(params.Output, e.cfg.Alarm.OutputLength),
		role,
		initiator,
	)

	if err != nil {
		return "", err
	}

	return types.AlarmChangeTypeAck, nil
}
