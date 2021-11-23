package executor

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func NewMongoUpdateExecutor(
	executor operation.Executor,
	adapter alarm.Adapter,
) operation.Executor {
	return &mongoUpdateExecutor{
		executor: executor,
		adapter:  adapter,
	}
}

type mongoUpdateExecutor struct {
	executor operation.Executor
	adapter  alarm.Adapter
}

// Exec finds executor by operation and calls it.
func (e *mongoUpdateExecutor) Exec(
	ctx context.Context,
	operation types.Operation,
	alarm *types.Alarm,
	entity types.Entity,
	time types.CpsTime,
	userID, role, initiator string,
) (types.AlarmChangeType, error) {
	changeType, err := e.executor.Exec(ctx, operation, alarm, entity, time, userID, role, initiator)
	if err != nil {
		return "", err
	}

	err = e.adapter.PartialUpdateOpen(ctx, alarm)
	if err != nil {
		return "", err
	}

	return changeType, nil
}
