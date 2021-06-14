package executor

import (
	operationlib "git.canopsis.net/canopsis/go-engines/lib/canopsis/operation"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"time"
)

func NewResolveCloseExecutor() operationlib.Executor {
	return &resolveCloseExecutor{}
}

type resolveCloseExecutor struct {
}

func (e *resolveCloseExecutor) Exec(
	_ types.Operation,
	alarm *types.Alarm,
	_ types.CpsTime,
	_, _ string,
) (types.AlarmChangeType, error) {
	if alarm.Value.Resolved != nil || !alarm.Closable(0*time.Second) {
		return "", nil
	}

	err := alarm.PartialUpdateResolve(types.CpsTime{Time: time.Now()})
	if err != nil {
		return "", err
	}

	return types.AlarmChangeTypeResolve, nil
}
