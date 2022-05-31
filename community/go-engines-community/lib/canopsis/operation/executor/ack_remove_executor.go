package executor

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	operationlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
)

// NewAckRemoveExecutor creates new executor.
func NewAckRemoveExecutor(metricsSender metrics.Sender, configProvider config.AlarmConfigProvider) operationlib.Executor {
	return &ackRemoveExecutor{
		metricsSender:  metricsSender,
		configProvider: configProvider,
	}
}

type ackRemoveExecutor struct {
	metricsSender  metrics.Sender
	configProvider config.AlarmConfigProvider
}

// Exec creates new ack remove step for alarm.
func (e *ackRemoveExecutor) Exec(
	ctx context.Context,
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

	if alarm.Value.ACK == nil {
		return "", nil
	}

	err := alarm.PartialUpdateUnack(
		time,
		params.Author,
		utils.TruncateString(params.Output, e.configProvider.Get().OutputLength),
		userID,
		role,
		initiator,
	)
	if err != nil {
		return "", err
	}

	go e.metricsSender.SendCancelAck(context.Background(), *alarm, time.Time)

	return types.AlarmChangeTypeAckremove, nil
}
