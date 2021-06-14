package executor

import (
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/config"
	operationlib "git.canopsis.net/canopsis/go-engines/lib/canopsis/operation"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/utils"
)

// NewChangeStateExecutor creates new executor.
func NewChangeStateExecutor(cfg config.CanopsisConf) operationlib.Executor {
	return &changeStateExecutor{cfg: cfg}
}

type changeStateExecutor struct {
	cfg config.CanopsisConf
}

// Exec emits change state event.
func (e *changeStateExecutor) Exec(
	operation types.Operation,
	alarm *types.Alarm,
	time types.CpsTime,
	role, initiator string,
) (types.AlarmChangeType, error) {
	var params types.OperationChangeStateParameters
	var ok bool
	if params, ok = operation.Parameters.(types.OperationChangeStateParameters); !ok {
		return "", fmt.Errorf("invalid parameters")
	}

	if alarm.Value.State == nil || alarm.Value.State.Value == types.AlarmStateOK {
		return "", fmt.Errorf("cannot change ok state")
	}

	if alarm.Value.State != nil && alarm.Value.State.Value == params.State {
		return "", nil
	}

	err := alarm.PartialUpdateChangeState(
		params.State,
		time,
		params.Author,
		utils.TruncateString(params.Output, e.cfg.Alarm.OutputLength),
		role,
		initiator,
	)
	if err != nil {
		return "", err
	}

	return types.AlarmChangeTypeChangeState, nil
}
