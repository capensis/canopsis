package executor

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type autoInstructionActivateExecutor struct{}

func NewAutoInstructionActivateExecutor() operation.Executor {
	return &autoInstructionActivateExecutor{}
}

func (e *autoInstructionActivateExecutor) Exec(
	_ context.Context,
	_ types.Operation,
	alarm *types.Alarm,
	_ *types.Entity,
	timestamp types.CpsTime,
	_, _, _ string,
) (types.AlarmChangeType, error) {
	if alarm.IsActivated() || !alarm.InactiveAutoInstructionInProgress {
		return types.AlarmChangeTypeNone, nil
	}

	alarm.InactiveAutoInstructionInProgress = false
	alarm.PartialUpdateUnsetAutoInstructionInProgress(timestamp)

	return types.AlarmChangeTypeAutoInstructionActivate, nil
}
