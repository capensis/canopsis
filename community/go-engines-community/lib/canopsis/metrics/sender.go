package metrics

//go:generate mockgen -destination=../../../mocks/lib/metrics/sender.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics Sender

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type Sender interface {
	Run(ctx context.Context)
	SendAck(alarm types.Alarm, userID string, timestamp time.Time)
	SendCancelAck(alarm types.Alarm, timestamp time.Time)
	SendTicket(alarm types.Alarm, userID string, timestamp time.Time)
	SendResolve(alarm types.Alarm, entity types.Entity, timestamp time.Time)
	SendCreate(alarm types.Alarm, timestamp time.Time)
	SendCreateAndPbhEnter(alarm types.Alarm, timestamp time.Time)
	SendCorrelation(timestamp time.Time, child types.Alarm)
	SendUserActivity(timestamp time.Time, username string, value int64)
	SendPbhEnter(alarm types.Alarm, entity types.Entity)
	SendPbhLeave(entity types.Entity, timestamp time.Time, prevCanonicalType string, prevTimestamp time.Time)
	SendPbhLeaveAndEnter(alarm types.Alarm, entity types.Entity, prevCanonicalType string, prevTimestamp time.Time)
	SendUpdateState(alarm types.Alarm, entity types.Entity, previousState types.CpsNumber)

	SendAutoInstructionExecutionStart(alarm types.Alarm, timestamp time.Time)
	SendAutoInstructionExecutionForInstruction(instructionID string, timestamp time.Time)
	SendAutoInstructionAssignForInstructions(instructionIDs []string, timestamp time.Time)

	SendInstructionAssignForAlarm(entityID string, timestamp time.Time)
	SendInstructionAssignForAlarms(entityIDs []string, timestamp time.Time)
	SendInstructionExecutionForAlarm(entityID string, timestamp time.Time)
	SendInstructionAssignForInstruction(instructionID string, timestamp time.Time, value int64)
	SendInstructionAssignForInstructions(instructionIDs []string, timestamp time.Time)
	SendInstructionExecutionForInstruction(instructionID string, timestamp time.Time)

	SendNotAckedInHourInc(alarm types.Alarm, timestamp time.Time)
	SendNotAckedInFourHoursInc(alarm types.Alarm, timestamp time.Time)
	SendNotAckedInDayInc(alarm types.Alarm, timestamp time.Time)
	SendNotAckedInHourDec(alarm types.Alarm, timestamp time.Time)
	SendNotAckedInFourHoursDec(alarm types.Alarm, timestamp time.Time)
	SendNotAckedInDayDec(alarm types.Alarm, timestamp time.Time)
	SendRemoveNotAckedMetric(alarm types.Alarm, timestamp time.Time, notAckedMetricType string)

	SendPerfData(timestamp time.Time, entityID, name string, value float64, unit string)

	SendEventMetrics(alarm types.Alarm, entity types.Entity, alarmChange types.AlarmChange, timestamp time.Time, initiator, userID, instructionID, notAckedMetricType string)

	SendSliMetric(timestamp time.Time, alarm types.Alarm, entity types.Entity)

	SendMessageRate(timestamp time.Time, eventType, connectorName string)
}
