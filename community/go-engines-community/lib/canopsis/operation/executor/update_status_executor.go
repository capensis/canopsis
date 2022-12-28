package executor

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
)

func NewUpdateStatusExecutor(configProvider config.AlarmConfigProvider, alarmStatusService alarmstatus.Service) operation.Executor {
	return &updateStatusExecutor{configProvider: configProvider, alarmStatusService: alarmStatusService}
}

type updateStatusExecutor struct {
	configProvider     config.AlarmConfigProvider
	alarmStatusService alarmstatus.Service
}

func (e *updateStatusExecutor) Exec(
	_ context.Context,
	op types.Operation,
	alarm *types.Alarm,
	entity *types.Entity,
	time types.CpsTime,
	_, _, _ string,
) (types.AlarmChangeType, error) {
	params := op.Parameters

	currentStatus := alarm.Value.Status.Value
	newStatus := e.alarmStatusService.ComputeStatus(*alarm, *entity)

	if newStatus == currentStatus {
		return "", nil
	}

	// Create new Step to keep track of the alarm history
	newStep := types.NewAlarmStep(types.AlarmStepStatusIncrease, time, alarm.Value.Connector+"."+alarm.Value.ConnectorName, params.Output, "", "", "")
	newStep.Value = newStatus

	if newStatus < currentStatus {
		newStep.Type = types.AlarmStepStatusDecrease
	}

	alarm.Value.Status = &newStep
	err := alarm.Value.Steps.Add(newStep)
	if err != nil {
		return "", err
	}

	alarm.Value.StateChangesSinceStatusUpdate = 0
	alarm.Value.LastUpdateDate = time

	alarm.AddUpdate("$set", bson.M{
		"v.status":                            alarm.Value.Status,
		"v.state_changes_since_status_update": alarm.Value.StateChangesSinceStatusUpdate,
		"v.last_update_date":                  alarm.Value.LastUpdateDate,
	})
	alarm.AddUpdate("$push", bson.M{"v.steps": alarm.Value.Status})

	return types.AlarmChangeTypeUpdateStatus, nil
}
