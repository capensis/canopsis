package executor

import (
	operationlib "git.canopsis.net/canopsis/go-engines/lib/canopsis/operation"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
)

func NewActivateExecutor() operationlib.Executor {
	return &activateExecutor{}
}

type activateExecutor struct {
}

func (e *activateExecutor) Exec(
	_ types.Operation,
	alarm *types.Alarm,
	time types.CpsTime,
	_, _ string,
) (types.AlarmChangeType, error) {
	if alarm.IsActivated() {
		return "", nil
	}

	err := alarm.PartialUpdateActivate(time)
	if err != nil {
		return "", err
	}

	return types.AlarmChangeTypeActivate, nil
}
