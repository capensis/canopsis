package executor

import (
	"context"

	operationlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func NewResolveDeletedExecutor() operationlib.Executor {
	return &resolveDeletedExecutor{}
}

type resolveDeletedExecutor struct {
}

func (e *resolveDeletedExecutor) Exec(
	_ context.Context,
	_ types.Operation,
	alarm *types.Alarm,
	entity *types.Entity,
	_ types.CpsTime,
	_, _, _ string,
) (types.AlarmChangeType, error) {
	if alarm.Value.Resolved != nil || entity.SoftDeleted == nil {
		return "", nil
	}

	err := alarm.PartialUpdateResolve(types.NewCpsTime())
	if err != nil {
		return "", err
	}

	return types.AlarmChangeTypeResolve, nil
}
