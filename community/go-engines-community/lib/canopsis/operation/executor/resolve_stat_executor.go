package executor

import (
	"context"
	libentity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	operationlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statsng"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func NewResolveStatExecutor(
	executor operationlib.Executor,
	entityAdapter libentity.Adapter,
	statService statsng.Service,
	metricsSender metrics.Sender,
) operationlib.Executor {
	return &resolveStatExecutor{
		executor:      executor,
		entityAdapter: entityAdapter,
		statService:   statService,
		metricsSender: metricsSender,
	}
}

type resolveStatExecutor struct {
	executor      operationlib.Executor
	entityAdapter libentity.Adapter
	statService   statsng.Service
	metricsSender metrics.Sender
}

func (e *resolveStatExecutor) Exec(
	ctx context.Context,
	operation types.Operation,
	alarm *types.Alarm,
	entity *types.Entity,
	timestamp types.CpsTime,
	role, initiator string,
) (types.AlarmChangeType, error) {
	changeType, err := e.executor.Exec(ctx, operation, alarm, entity, timestamp, role, initiator)
	if err != nil {
		return "", err
	}

	if changeType != "" {
		err = e.statService.ProcessResolvedAlarm(*alarm, *entity)
		if err != nil {
			return "", err
		}

		go e.metricsSender.SendResolve(ctx, *alarm, *entity, timestamp.Time)
	}

	return changeType, nil
}
