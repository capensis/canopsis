package executor

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
)

type pbhLeaveExecutor struct {
	configProvider config.AlarmConfigProvider

	metricsSender metrics.Sender
}

func NewPbhLeaveExecutor(configProvider config.AlarmConfigProvider, metricsSender metrics.Sender) operation.Executor {
	return &pbhLeaveExecutor{configProvider: configProvider, metricsSender: metricsSender}
}

func (e *pbhLeaveExecutor) Exec(
	_ context.Context,
	op types.Operation,
	alarm *types.Alarm,
	entity *types.Entity,
	time types.CpsTime,
	userID, role, initiator string,
) (types.AlarmChangeType, error) {
	params := op.Parameters

	if userID == "" {
		userID = params.User
	}

	currPbehaviorInfo := entity.PbehaviorInfo

	if currPbehaviorInfo.IsDefaultActive() {
		return "", nil
	}

	err := alarm.PartialUpdatePbhLeave(
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

	entity.PbehaviorInfo = alarm.Value.PbehaviorInfo

	e.metricsSender.SendPbhLeave(*entity, time.Time, currPbehaviorInfo.CanonicalType, currPbehaviorInfo.Timestamp.Time)

	return types.AlarmChangeTypePbhLeave, nil
}
