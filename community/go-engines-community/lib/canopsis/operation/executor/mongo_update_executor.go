package executor

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/operation"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
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
	operation types.Operation,
	alarm *types.Alarm,
	time types.CpsTime,
	role, initiator string,
) (types.AlarmChangeType, error) {
	changeType, err := e.executor.Exec(operation, alarm, time, role, initiator)
	if err != nil {
		return "", err
	}

	err = e.adapter.PartialUpdateOpen(alarm)
	if err != nil {
		return "", err
	}

	return changeType, nil
}
