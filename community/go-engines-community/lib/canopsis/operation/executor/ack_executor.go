/*
Package executor contains operation executors.
*/
package executor

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
)

// NewAckExecutor creates new executor.
func NewAckExecutor(metricsSender metrics.Sender, configProvider config.AlarmConfigProvider) operation.Executor {
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
	_ context.Context,
	op types.Operation,
	alarm *types.Alarm,
	_ *types.Entity,
	time types.CpsTime,
	userID, role, initiator string,
) (types.AlarmChangeType, error) {
	params := op.Parameters

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
		metricsUserID := ""
		if initiator == types.InitiatorUser {
			metricsUserID = userID
		}
		e.metricsSender.SendAck(*alarm, metricsUserID, time.Time)
		e.metricsSender.SendRemoveNotAckedMetric(*alarm, time.Time)

		return types.AlarmChangeTypeAck, nil
	}

	return types.AlarmChangeTypeDoubleAck, nil
}
