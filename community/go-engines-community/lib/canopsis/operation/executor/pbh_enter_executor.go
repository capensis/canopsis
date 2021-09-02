package executor

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	operationlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
)

type pbhEnterExecutor struct {
	configProvider config.AlarmConfigProvider
}

// NewAckExecutor creates new executor.
func NewPbhEnterExecutor(configProvider config.AlarmConfigProvider) operationlib.Executor {
	return &pbhEnterExecutor{configProvider: configProvider}
}

func (e *pbhEnterExecutor) Exec(
	_ context.Context,
	operation types.Operation,
	alarm *types.Alarm,
	time types.CpsTime,
	role, initiator string,
) (types.AlarmChangeType, error) {
	var params types.OperationPbhParameters
	var ok bool
	if params, ok = operation.Parameters.(types.OperationPbhParameters); !ok {
		return "", fmt.Errorf("invalid parameters")
	}

	if alarm.Value.PbehaviorInfo == params.PbehaviorInfo {
		return "", nil
	}

	err := alarm.PartialUpdatePbhEnter(
		time,
		params.PbehaviorInfo,
		params.Author,
		utils.TruncateString(params.Output, e.configProvider.Get().OutputLength),
		role,
		initiator,
	)
	if err != nil {
		return "", err
	}

	return types.AlarmChangeTypePbhEnter, nil
}
