package executor

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
)

type pbhEnterExecutor struct {
	configProvider config.AlarmConfigProvider

	metricsSender metrics.Sender
}

func NewPbhEnterExecutor(configProvider config.AlarmConfigProvider, metricsSender metrics.Sender) operation.Executor {
	return &pbhEnterExecutor{configProvider: configProvider, metricsSender: metricsSender}
}

func (e *pbhEnterExecutor) Exec(
	_ context.Context,
	op types.Operation,
	alarm *types.Alarm,
	entity *types.Entity,
	time types.CpsTime,
	userID, role, initiator string,
) (types.AlarmChangeType, error) {
	params := op.Parameters
	if params.PbehaviorInfo == nil {
		return "", fmt.Errorf("invalid parameters")
	}

	if userID == "" {
		userID = params.User
	}

	if entity.PbehaviorInfo.Same(*params.PbehaviorInfo) {
		return "", nil
	}

	err := alarm.PartialUpdatePbhEnter(
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

	e.metricsSender.SendPbhEnter(alarm, *entity)
	if !alarm.Value.PbehaviorInfo.IsActive() {
		e.metricsSender.SendRemoveNotAckedMetric(*alarm, time.Time)
	}

	return types.AlarmChangeTypePbhEnter, nil
}
