package executor

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func NewResolveDoneExecutor() operation.Executor {
	return &resolveDoneExecutor{}
}

type resolveDoneExecutor struct {
}

func (e *resolveDoneExecutor) Exec(
	_ context.Context,
	_ types.Operation,
	alarm *types.Alarm,
	entity *types.Entity,
	_ types.CpsTime,
	_, _, _ string,
) (types.AlarmChangeType, error) {
	if alarm.Value.Resolved != nil || alarm.Value.Done == nil {
		return "", nil
	}

	err := alarm.PartialUpdateResolve(types.NewCpsTime())
	if err != nil {
		return "", err
	}
	entity.IdleSince = nil
	entity.LastIdleRuleApply = ""

	return types.AlarmChangeTypeResolve, nil
}
