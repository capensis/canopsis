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

type pbhEnterExecutor struct {
	configProvider config.AlarmConfigProvider

	metricsSender metrics.Sender
}

func NewPbhEnterExecutor(configProvider config.AlarmConfigProvider, metricsSender metrics.Sender) operationlib.Executor {
	return &pbhEnterExecutor{configProvider: configProvider, metricsSender: metricsSender}
}

func (e *pbhEnterExecutor) Exec(
	_ context.Context,
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

	go e.metricsSender.SendPbhEnter(context.Background(), alarm, *entity)

	return types.AlarmChangeTypePbhEnter, nil
}
