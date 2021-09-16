package executor

import (
	"context"
	"fmt"
	libentity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	operationlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statsng"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func NewResolveStatExecutor(
	executor operationlib.Executor,
	entityAdapter libentity.Adapter,
	statService statsng.Service,
) operationlib.Executor {
	return &resolveStatExecutor{
		executor:      executor,
		entityAdapter: entityAdapter,
		statService:   statService,
	}
}

type resolveStatExecutor struct {
	executor      operationlib.Executor
	entityAdapter libentity.Adapter
	statService   statsng.Service
}

func (e *resolveStatExecutor) Exec(
	ctx context.Context,
	operation types.Operation,
	alarm *types.Alarm,
	timestamp types.CpsTime,
	role, initiator string,
) (types.AlarmChangeType, error) {
	changeType, err := e.executor.Exec(ctx, operation, alarm, timestamp, role, initiator)
	if err != nil {
		return "", err
	}

	if changeType != "" {
		entity, ok := e.entityAdapter.Get(ctx, alarm.EntityID)
		if !ok {
			return "", fmt.Errorf("cannot found entity")
		}

		err = e.statService.ProcessResolvedAlarm(*alarm, entity)
		if err != nil {
			return "", err
		}
	}

	return changeType, nil
}
