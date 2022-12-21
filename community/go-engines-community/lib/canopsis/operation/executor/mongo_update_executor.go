package executor

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func NewMongoUpdateExecutor(
	executor operation.Executor,
	adapter alarm.Adapter,
	entityAdapter entity.Adapter,
) operation.Executor {
	return &mongoUpdateExecutor{
		executor:      executor,
		adapter:       adapter,
		entityAdapter: entityAdapter,
	}
}

type mongoUpdateExecutor struct {
	executor      operation.Executor
	adapter       alarm.Adapter
	entityAdapter entity.Adapter
}

// Exec finds executor by operation and calls it.
func (e *mongoUpdateExecutor) Exec(
	ctx context.Context,
	op types.Operation,
	alarm *types.Alarm,
	entity *types.Entity,
	time types.CpsTime,
	userID, role, initiator string,
) (types.AlarmChangeType, error) {
	changeType, err := e.executor.Exec(ctx, op, alarm, entity, time, userID, role, initiator)
	if err != nil {
		return "", err
	}

	err = e.adapter.PartialUpdateOpen(ctx, alarm)
	if err != nil {
		return "", fmt.Errorf("cannot update alarm: %w", err)
	}

	switch changeType {
	case types.AlarmChangeTypeCreateAndPbhEnter, types.AlarmChangeTypePbhEnter,
		types.AlarmChangeTypePbhLeave, types.AlarmChangeTypePbhLeaveAndEnter:
		err := e.entityAdapter.UpdatePbehaviorInfo(ctx, entity.ID, entity.PbehaviorInfo)
		if err != nil {
			return "", fmt.Errorf("cannot update entity: %w", err)
		}
	}

	return changeType, nil
}
