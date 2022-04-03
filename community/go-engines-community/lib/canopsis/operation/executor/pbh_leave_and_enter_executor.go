package executor

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	operationlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
)

type pbhLeaveAndEnterExecutor struct {
	configProvider config.AlarmConfigProvider
}

// NewAckExecutor creates new executor.
func NewPbhLeaveAndEnterExecutor(configProvider config.AlarmConfigProvider) operationlib.Executor {
	return &pbhLeaveAndEnterExecutor{configProvider: configProvider}
}

func (e *pbhLeaveAndEnterExecutor) Exec(
	ctx context.Context,
	operation types.Operation,
	alarm *types.Alarm,
	entity *types.Entity,
	time types.CpsTime,
	userID, role, initiator string,
) (types.AlarmChangeType, error) {
	var params types.OperationPbhParameters
	var ok bool
	if params, ok = operation.Parameters.(types.OperationPbhParameters); !ok {
		return "", fmt.Errorf("invalid parameters")
	}

	if userID == "" {
		userID = params.User
	}

	currPbehaviorInfo := alarm.Value.PbehaviorInfo

	if currPbehaviorInfo.Same(params.PbehaviorInfo) {
		return "", nil
	}

	err := alarm.PartialUpdatePbhLeaveAndEnter(
		time,
		params.PbehaviorInfo,
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

	return types.AlarmChangeTypePbhLeaveAndEnter, nil
}
