package executor

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func NewUnsnoozeExecutor() operation.Executor {
	return &unsnoozeExecutor{}
}

type unsnoozeExecutor struct {
}

func (e *unsnoozeExecutor) Exec(
	_ context.Context,
	_ types.Operation,
	alarm *types.Alarm,
	_ *types.Entity,
	time types.CpsTime,
	_, _, _ string,
) (types.AlarmChangeType, error) {
	if alarm.Value.Snooze == nil {
		return "", nil
	}

	err := alarm.PartialUpdateUnsnooze(time)
	if err != nil {
		return "", err
	}

	return types.AlarmChangeTypeUnsnooze, nil
}
