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

type pbhLeaveAndEnterExecutor struct {
	configProvider config.AlarmConfigProvider

	metricsSender metrics.Sender
}

// NewAckExecutor creates new executor.
func NewPbhLeaveAndEnterExecutor(configProvider config.AlarmConfigProvider, metricsSender metrics.Sender) operationlib.Executor {
	return &pbhLeaveAndEnterExecutor{configProvider: configProvider, metricsSender: metricsSender}
}

func (e *pbhLeaveAndEnterExecutor) Exec(
	ctx context.Context,
	operation types.Operation,
	alarm *types.Alarm,
	entity *types.Entity,
	time types.CpsTime,
	userID, role, initiator string,
) (types.AlarmChangeType, error) {
	params := operation.Parameters
	if params.PbehaviorInfo == nil {
		return "", fmt.Errorf("invalid parameters")
	}

	if userID == "" {
		userID = params.User
	}

	currPbehaviorInfo := alarm.Value.PbehaviorInfo

	if currPbehaviorInfo.Same(*params.PbehaviorInfo) {
		return "", nil
	}

	err := alarm.PartialUpdatePbhLeaveAndEnter(
		time,
		*params.PbehaviorInfo,
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

	go e.metricsSender.SendPbhLeaveAndEnter(context.Background(), alarm, *entity, currPbehaviorInfo.CanonicalType, currPbehaviorInfo.Timestamp.Time)

	return types.AlarmChangeTypePbhLeaveAndEnter, nil
}
