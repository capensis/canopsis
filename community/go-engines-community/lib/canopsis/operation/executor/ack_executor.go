/*
Package executor contains operation executors.
*/
package executor

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	operationlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
)

// NewAckExecutor creates new executor.
func NewAckExecutor(configProvider config.AlarmConfigProvider) operationlib.Executor {
	return &ackExecutor{
		configProvider: configProvider,
	}
}

type ackExecutor struct {
	configProvider config.AlarmConfigProvider
}

// Exec creates new ack step for alarm.
func (e *ackExecutor) Exec(
	_ context.Context,
	operation types.Operation,
	alarm *types.Alarm,
	_ *types.Entity,
	time types.CpsTime,
	userID, role, initiator string,
) (types.AlarmChangeType, error) {
	params := operation.Parameters

	if userID == "" {
		userID = params.User
	}

	allowDoubleAck := e.configProvider.Get().AllowDoubleAck
	doubleAck := alarm.Value.ACK != nil
	if doubleAck && !allowDoubleAck {
		return "", nil
	}

	err := alarm.PartialUpdateAck(
		time,
		params.Author,
		utils.TruncateString(params.Output, e.configProvider.Get().OutputLength),
		userID,
		role,
		initiator,
		allowDoubleAck,
	)

	if err != nil {
		return "", err
	}

	if !doubleAck {
		return types.AlarmChangeTypeAck, nil
	}

	return types.AlarmChangeTypeDoubleAck, nil
}
