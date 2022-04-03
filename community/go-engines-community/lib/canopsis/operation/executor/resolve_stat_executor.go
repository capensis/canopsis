package executor

import (
	"context"
	libentity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	operationlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func NewResolveStatExecutor(
	executor operationlib.Executor,
	entityAdapter libentity.Adapter,
) operationlib.Executor {
	return &resolveStatExecutor{
		executor:      executor,
		entityAdapter: entityAdapter,
	}
}

type resolveStatExecutor struct {
	executor      operationlib.Executor
	entityAdapter libentity.Adapter
}

func (e *resolveStatExecutor) Exec(
	ctx context.Context,
	operation types.Operation,
	alarm *types.Alarm,
	entity *types.Entity,
	timestamp types.CpsTime,
	userID, role, initiator string,
) (types.AlarmChangeType, error) {
	changeType, err := e.executor.Exec(ctx, operation, alarm, entity, timestamp, userID, role, initiator)
	if err != nil {
		return "", err
	}

	return changeType, nil
}
