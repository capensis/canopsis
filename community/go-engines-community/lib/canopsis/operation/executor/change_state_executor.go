package executor

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	operationlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
)

// NewChangeStateExecutor creates new executor.
func NewChangeStateExecutor(configProvider config.AlarmConfigProvider) operationlib.Executor {
	return &changeStateExecutor{configProvider: configProvider}
}

type changeStateExecutor struct {
	configProvider config.AlarmConfigProvider
}

// Exec emits change state event.
func (e *changeStateExecutor) Exec(
	_ context.Context,
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
		utils.TruncateString(params.Output, e.configProvider.Get().OutputLength),
		role,
		initiator,
	)
	if err != nil {
		return "", err
	}

	return types.AlarmChangeTypeChangeState, nil
}
