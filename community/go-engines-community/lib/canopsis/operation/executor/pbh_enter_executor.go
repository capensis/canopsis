package executor

import (
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/config"
	operationlib "git.canopsis.net/canopsis/go-engines/lib/canopsis/operation"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/utils"
)

type pbhEnterExecutor struct {
	cfg config.CanopsisConf
}

// NewAckExecutor creates new executor.
func NewPbhEnterExecutor(cfg config.CanopsisConf) operationlib.Executor {
	return &pbhEnterExecutor{cfg: cfg}
}

func (e *pbhEnterExecutor) Exec(
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
		utils.TruncateString(params.Output, e.cfg.Alarm.OutputLength),
		role,
		initiator,
	)
	if err != nil {
		return "", err
	}

	return types.AlarmChangeTypePbhEnter, nil
}
