package executor

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	operationlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
)

func NewDoneExecutor(configProvider config.AlarmConfigProvider) operationlib.Executor {
	return &doneExecutor{configProvider: configProvider}
}

type doneExecutor struct {
	configProvider config.AlarmConfigProvider
}

func (e *doneExecutor) Exec(
	_ context.Context,
	operation types.Operation,
	alarm *types.Alarm,
	_ *types.Entity,
	time types.CpsTime,
	userID, role, initiator string,
) (types.AlarmChangeType, error) {
	var params types.OperationParameters
	var ok bool
	if params, ok = operation.Parameters.(types.OperationParameters); !ok {
		return "", fmt.Errorf("invalid parameters")
	}

	if userID == "" {
		userID = params.User
	}

	if alarm.Value.Done != nil {
		return "", nil
	}

	err := alarm.PartialUpdateDone(
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

	return types.AlarmChangeTypeDone, nil
}
