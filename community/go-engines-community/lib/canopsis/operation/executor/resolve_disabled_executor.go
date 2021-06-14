package executor

import (
	operationlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"time"
)

func NewResolveDisabledExecutor() operationlib.Executor {
	return &resolveDisabledExecutor{}
}

type resolveDisabledExecutor struct {
}

func (e *resolveDisabledExecutor) Exec(
	_ types.Operation,
	alarm *types.Alarm,
	_ types.CpsTime,
	_, _ string,
) (types.AlarmChangeType, error) {
	err := alarm.PartialUpdateResolve(types.CpsTime{Time: time.Now()})
	if err != nil {
		return "", err
	}

	return types.AlarmChangeTypeResolve, nil
}
