package executor

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func NewUncancelExecutor(configProvider config.AlarmConfigProvider, alarmStatusService alarmstatus.Service) operation.Executor {
	return &uncancelExecutor{configProvider: configProvider, alarmStatusService: alarmStatusService}
}

type uncancelExecutor struct {
	configProvider     config.AlarmConfigProvider
	alarmStatusService alarmstatus.Service
}

func (e *uncancelExecutor) Exec(
	_ context.Context,
	op types.Operation,
	alarm *types.Alarm,
	entity *types.Entity,
	time types.CpsTime,
	userID, role, initiator string,
) (types.AlarmChangeType, error) {
	params := op.Parameters

	if userID == "" {
		userID = params.User
	}

	if alarm.Value.Canceled == nil {
		return "", nil
	}

	alarmConfig := e.configProvider.Get()
	output := utils.TruncateString(params.Output, alarmConfig.OutputLength)
	newStep := types.NewAlarmStep(types.AlarmStepUncancel, time, params.Author, output, userID, role, initiator)
	alarm.Value.Canceled = nil

	if err := alarm.Value.Steps.Add(newStep); err != nil {
		return "", err
	}

	alarm.AddUpdate("$set", bson.M{"v.canceled": alarm.Value.Canceled})
	alarm.AddUpdate("$push", bson.M{"v.steps": newStep})

	currentStatus := alarm.Value.Status.Value
	newStatus := e.alarmStatusService.ComputeStatus(*alarm, *entity)

	if newStatus == currentStatus {
		alarm.AddUpdate("$set", bson.M{"v.canceled": alarm.Value.Canceled})
		alarm.AddUpdate("$push", bson.M{"v.steps": newStep})
		return types.AlarmChangeTypeUncancel, nil
	}

	newStepStatus := types.NewAlarmStep(types.AlarmStepStatusIncrease, time, alarm.Value.Connector+"."+alarm.Value.ConnectorName, output, userID, role, initiator)
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
		"v.canceled":                          alarm.Value.Canceled,
		"v.status":                            alarm.Value.Status,
		"v.state_changes_since_status_update": alarm.Value.StateChangesSinceStatusUpdate,
		"v.last_update_date":                  alarm.Value.LastUpdateDate,
	})
	alarm.AddUpdate("$push", bson.M{"v.steps": bson.M{"$each": bson.A{newStep, alarm.Value.Status}}})

	return types.AlarmChangeTypeUncancel, nil
}
