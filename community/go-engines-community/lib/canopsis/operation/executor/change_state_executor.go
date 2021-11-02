package executor

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	operationlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// NewChangeStateExecutor creates new executor.
func NewChangeStateExecutor(
	configProvider config.AlarmConfigProvider,
	alarmStatusService alarmstatus.Service,
	metricsSender metrics.Sender,
) operationlib.Executor {
	return &changeStateExecutor{
		configProvider:     configProvider,
		alarmStatusService: alarmStatusService,
		metricsSender:      metricsSender,
	}
}

type changeStateExecutor struct {
	configProvider     config.AlarmConfigProvider
	alarmStatusService alarmstatus.Service

	metricsSender metrics.Sender
}

// Exec emits change state event.
func (e *changeStateExecutor) Exec(
	ctx context.Context,
	operation types.Operation,
	alarm *types.Alarm,
	entity *types.Entity,
	time types.CpsTime,
	role, initiator string,
) (types.AlarmChangeType, error) {
	var params types.OperationChangeStateParameters
	var ok bool
	if params, ok = operation.Parameters.(types.OperationChangeStateParameters); !ok {
		return "", fmt.Errorf("invalid parameters")
	}

	currentState := alarm.Value.State.Value
	if currentState == types.AlarmStateOK {
		return "", fmt.Errorf("cannot change ok state")
	}

	if currentState == params.State {
		return "", nil
	}

	conf := e.configProvider.Get()
	output := utils.TruncateString(params.Output, conf.OutputLength)

	newStep := types.NewAlarmStep(types.AlarmStepChangeState, time, params.Author, output, role, initiator)
	newStep.Value = params.State
	alarm.Value.State = &newStep

	err := alarm.Value.Steps.Add(newStep)
	if err != nil {
		return "", err
	}

	currentStatus := alarm.Value.Status.Value
	newStatus := e.alarmStatusService.ComputeStatus(*alarm, *entity)

	if newStatus == currentStatus {
		alarm.AddUpdate("$set", bson.M{"v.state": alarm.Value.State})
		alarm.AddUpdate("$push", bson.M{"v.steps": alarm.Value.State})

		return types.AlarmChangeTypeChangeState, nil
	}

	newStepStatus := types.NewAlarmStep(types.AlarmStepStatusIncrease, time, params.Author, output, role, initiator)
	newStepStatus.Value = newStatus
	if alarm.Value.Status != nil && newStepStatus.Value < alarm.Value.Status.Value {
		newStepStatus.Type = types.AlarmStepStatusDecrease
	}
	alarm.Value.Status = &newStepStatus
	if err := alarm.Value.Steps.Add(newStepStatus); err != nil {
		return "", err
	}

	alarm.Value.StateChangesSinceStatusUpdate = 0
	alarm.Value.LastUpdateDate = time

	alarm.AddUpdate("$set", bson.M{
		"v.state":                             alarm.Value.State,
		"v.status":                            alarm.Value.Status,
		"v.state_changes_since_status_update": alarm.Value.StateChangesSinceStatusUpdate,
		"v.last_update_date":                  alarm.Value.LastUpdateDate,
	})
	alarm.AddUpdate("$push", bson.M{"v.steps": bson.M{"$each": bson.A{alarm.Value.State, alarm.Value.Status}}})

	go e.metricsSender.SendUpdateState(ctx, *alarm, *entity, currentState)

	return types.AlarmChangeTypeChangeState, nil
}
