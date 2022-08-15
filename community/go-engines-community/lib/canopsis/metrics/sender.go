package metrics

//go:generate mockgen -destination=../../../mocks/lib/metrics/sender.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics Sender

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type Sender interface {
	SendAck(ctx context.Context, alarm types.Alarm, userID string, timestamp time.Time)
	SendCancelAck(ctx context.Context, alarm types.Alarm, timestamp time.Time)
	SendTicket(ctx context.Context, alarm types.Alarm, userID string, timestamp time.Time)
	SendResolve(ctx context.Context, alarm types.Alarm, entity types.Entity, timestamp time.Time)
	SendAutoInstructionStart(ctx context.Context, alarm types.Alarm, timestamp time.Time)
	SendCreate(ctx context.Context, alarm types.Alarm, timestamp time.Time)
	SendCreateAndPbhEnter(ctx context.Context, alarm types.Alarm, timestamp time.Time)
	SendCorrelation(ctx context.Context, timestamp time.Time, child types.Alarm)
	SendUserActivity(ctx context.Context, timestamp time.Time, username string, value int64)
	SendPbhEnter(ctx context.Context, alarm *types.Alarm, entity types.Entity)
	SendPbhLeave(ctx context.Context, entity types.Entity, timestamp time.Time, prevCanonicalType string, prevTimestamp time.Time)
	SendPbhLeaveAndEnter(ctx context.Context, alarm *types.Alarm, entity types.Entity, prevCanonicalType string, prevTimestamp time.Time)
	SendUpdateState(ctx context.Context, alarm types.Alarm, entity types.Entity, previousState types.CpsNumber)
	SendInstructionAssignForAlarm(ctx context.Context, entityID string, timestamp time.Time)
	SendInstructionExecutionForAlarm(ctx context.Context, entityID string, timestamp time.Time)
	SendInstructionAssignForInstruction(ctx context.Context, instructionID string, timestamp time.Time)
	SendInstructionExecutionForInstruction(ctx context.Context, instructionID string, timestamp time.Time)
}

type nullSender struct{}

func NewNullSender() Sender {
	return &nullSender{}
}

func (s *nullSender) SendAck(_ context.Context, _ types.Alarm, _ string, _ time.Time) {
}

func (s *nullSender) SendCancelAck(_ context.Context, _ types.Alarm, _ time.Time) {
}

func (s *nullSender) SendTicket(_ context.Context, _ types.Alarm, _ string, _ time.Time) {
}

func (s *nullSender) SendResolve(_ context.Context, _ types.Alarm, _ types.Entity, _ time.Time) {
}

func (s *nullSender) SendAutoInstructionStart(_ context.Context, _ types.Alarm, _ time.Time) {
}

func (s *nullSender) SendCreate(_ context.Context, _ types.Alarm, _ time.Time) {
}

func (s *nullSender) SendCreateAndPbhEnter(_ context.Context, _ types.Alarm, _ time.Time) {
}

func (s *nullSender) SendCorrelation(_ context.Context, _ time.Time, _ types.Alarm) {
}

func (s *nullSender) SendUserActivity(_ context.Context, _ time.Time, _ string, _ int64) {
}

func (s *nullSender) SendPbhEnter(_ context.Context, _ *types.Alarm, _ types.Entity) {

}

func (s *nullSender) SendPbhLeave(_ context.Context, _ types.Entity, _ time.Time, _ string, _ time.Time) {

}

func (s *nullSender) SendPbhLeaveAndEnter(_ context.Context, _ *types.Alarm, _ types.Entity, _ string, _ time.Time) {

}

func (s *nullSender) SendUpdateState(_ context.Context, _ types.Alarm, _ types.Entity, _ types.CpsNumber) {

}

func (s *nullSender) SendInstructionAssignForAlarm(_ context.Context, _ string, _ time.Time) {

}

func (s *nullSender) SendInstructionExecutionForAlarm(_ context.Context, _ string, _ time.Time) {

}

func (s *nullSender) SendInstructionAssignForInstruction(_ context.Context, _ string, _ time.Time) {

}

func (s *nullSender) SendInstructionExecutionForInstruction(_ context.Context, _ string, _ time.Time) {

}
