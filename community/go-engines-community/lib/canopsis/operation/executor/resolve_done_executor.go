package executor

import (
	"context"
	operationlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"time"
)

func NewResolveDoneExecutor() operationlib.Executor {
	return &resolveDoneExecutor{}
}

type resolveDoneExecutor struct {
}

func (e *resolveDoneExecutor) Exec(
	_ context.Context,
	_ types.Operation,
	alarm *types.Alarm,
	_ *types.Entity,
	_ types.CpsTime,
	_, _, _ string,
) (types.AlarmChangeType, error) {
	if alarm.Value.Resolved != nil || alarm.Value.Done == nil {
		return "", nil
	}

	err := alarm.PartialUpdateResolve(types.CpsTime{Time: time.Now()})
	if err != nil {
		return "", err
	}

	return types.AlarmChangeTypeResolve, nil
}
