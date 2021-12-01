package executor

import (
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
	operation types.Operation,
	alarm *types.Alarm,
	time types.CpsTime,
	userID, role, initiator string,
) (types.AlarmChangeType, error) {
	executor, ok := e.container.Get(operation.Type)
	if !ok {
		return "", nil
	}

	return executor.Exec(operation, alarm, time, userID, role, initiator)
}
