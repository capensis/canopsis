package executor

import (
	operationlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"time"
)

func NewResolveCancelExecutor() operationlib.Executor {
	return &resolveCancelExecutor{}
}

type resolveCancelExecutor struct {
}

func (e *resolveCancelExecutor) Exec(
	_ types.Operation,
	alarm *types.Alarm,
	_ types.CpsTime,
	_, _ string,
) (types.AlarmChangeType, error) {
	if alarm.Value.Resolved != nil || alarm.Value.Canceled == nil {
		return "", nil
	}

	err := alarm.PartialUpdateResolve(types.CpsTime{Time: time.Now()})
	if err != nil {
		return "", err
	}

	return types.AlarmChangeTypeResolve, nil
}
