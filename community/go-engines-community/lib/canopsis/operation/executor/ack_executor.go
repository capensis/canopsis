/*
Package executor contains operation executors.
*/
package executor

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	operationlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
)

// NewAckExecutor creates new executor.
func NewAckExecutor(metricsSender metrics.Sender, configProvider config.AlarmConfigProvider) operationlib.Executor {
	return &ackExecutor{
		metricsSender:  metricsSender,
		configProvider: configProvider,
	}
}

type ackExecutor struct {
	metricsSender  metrics.Sender
	configProvider config.AlarmConfigProvider
}

// Exec creates new ack step for alarm.
func (e *ackExecutor) Exec(
	ctx context.Context,
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
		utils.TruncateString(params.Output, e.configProvider.Get().OutputLength),
		role,
		initiator,
	)

	if err != nil {
		return "", err
	}

	go e.metricsSender.SendAck(ctx, *alarm, params.Author, time.Time)

	return types.AlarmChangeTypeAck, nil
}
