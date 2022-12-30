package executor

import (
	"context"

	operationlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func NewAutoWebhookStartExecutor() operationlib.Executor {
	return &autoWebhookStartExecutor{}
}

type autoWebhookStartExecutor struct {
}

func (e *autoWebhookStartExecutor) Exec(
	_ context.Context,
	operation types.Operation,
	alarm *types.Alarm,
	_ *types.Entity,
	time types.CpsTime,
	userID, role, initiator string,
) (types.AlarmChangeType, error) {
	params := operation.Parameters

	if userID == "" {
		userID = params.User
	}

	err := alarm.PartialUpdateAddStep(
		types.AlarmStepWebhookStart,
		time,
		params.Author,
		params.Output,
		userID,
		role,
		initiator,
	)
	if err != nil {
		return "", err
	}

	return types.AlarmChangeTypeAutoWebhookStart, nil
}
