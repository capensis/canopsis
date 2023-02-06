package executor

import (
	"context"

	libentity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func NewResolveStatExecutor(
	executor operation.Executor,
	entityAdapter libentity.Adapter,
	metricsSender metrics.Sender,
) operation.Executor {
	return &resolveStatExecutor{
		executor:      executor,
		entityAdapter: entityAdapter,
		metricsSender: metricsSender,
	}
}

type resolveStatExecutor struct {
	executor      operation.Executor
	entityAdapter libentity.Adapter
	metricsSender metrics.Sender
}

func (e *resolveStatExecutor) Exec(
	ctx context.Context,
	op types.Operation,
	alarm *types.Alarm,
	entity *types.Entity,
	timestamp types.CpsTime,
	userID, role, initiator string,
) (types.AlarmChangeType, error) {
	changeType, err := e.executor.Exec(ctx, op, alarm, entity, timestamp, userID, role, initiator)
	if err != nil {
		return "", err
	}

	if changeType != "" {
		e.metricsSender.SendResolve(*alarm, *entity, timestamp.Time)
		e.metricsSender.SendRemoveNotAckedMetric(*alarm, timestamp.Time)
	}

	return changeType, nil
}
