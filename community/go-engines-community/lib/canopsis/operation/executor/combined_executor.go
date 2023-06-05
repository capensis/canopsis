package executor

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

// NewCombinedExecutor creates new executor.
func NewCombinedExecutor(container operation.ExecutorContainer) operation.Executor {
	return &combinedExecutor{container: container}
}

type combinedExecutor struct {
	container operation.ExecutorContainer
}

// Exec finds executor by operation and calls it.
func (e *combinedExecutor) Exec(
	ctx context.Context,
	op types.Operation,
	alarm *types.Alarm,
	entity *types.Entity,
	time types.CpsTime,
	userID, role, initiator string,
) (types.AlarmChangeType, error) {
	executor, ok := e.container.Get(op.Type)
	if !ok {
		return "", nil
	}

	return executor.Exec(ctx, op, alarm, entity, time, userID, role, initiator)
}
