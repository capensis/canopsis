package executor

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func NewResolveDisabledExecutor() operation.Executor {
	return &resolveDisabledExecutor{}
}

type resolveDisabledExecutor struct {
}

func (e *resolveDisabledExecutor) Exec(
	_ context.Context,
	_ types.Operation,
	alarm *types.Alarm,
	entity *types.Entity,
	_ types.CpsTime,
	_, _, _ string,
) (types.AlarmChangeType, error) {
	entity.IdleSince = nil
	entity.LastIdleRuleApply = ""

	err := alarm.PartialUpdateResolve(types.CpsTime{Time: time.Now()})
	if err != nil {
		return "", err
	}

	return types.AlarmChangeTypeResolve, nil
}
