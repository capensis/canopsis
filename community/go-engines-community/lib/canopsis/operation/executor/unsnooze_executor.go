package executor

import (
	operationlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func NewUnsnoozeExecutor() operationlib.Executor {
	return &unsnoozeExecutor{}
}

type unsnoozeExecutor struct {
}

func (e *unsnoozeExecutor) Exec(
	_ types.Operation,
	alarm *types.Alarm,
	_ types.CpsTime,
	_, _ string,
) (types.AlarmChangeType, error) {
	if alarm.Value.Snooze == nil {
		return "", nil
	}

	err := alarm.PartialUpdateUnsnooze()
	if err != nil {
		return "", err
	}

	return types.AlarmChangeTypeUnsnooze, nil
}
