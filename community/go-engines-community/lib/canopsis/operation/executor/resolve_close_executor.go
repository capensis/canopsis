package executor

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func NewResolveCloseExecutor() operation.Executor {
	return &resolveCloseExecutor{}
}

type resolveCloseExecutor struct {
}

func (e *resolveCloseExecutor) Exec(
	_ context.Context,
	_ types.Operation,
	alarm *types.Alarm,
	entity *types.Entity,
	_ types.CpsTime,
	_, _, _ string,
) (types.AlarmChangeType, error) {
	if alarm.Value.Resolved != nil || !alarm.Closable(0*time.Second) {
		return "", nil
	}

	entity.IdleSince = nil
	entity.LastIdleRuleApply = ""

	err := alarm.PartialUpdateResolve(types.CpsTime{Time: time.Now()})
	if err != nil {
		return "", err
	}

	return types.AlarmChangeTypeResolve, nil
}
