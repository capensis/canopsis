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
	SendAutoInstructionStart(alarm types.Alarm, timestamp time.Time)
	SendCreate(alarm types.Alarm, timestamp time.Time)
	SendCreateAndPbhEnter(alarm types.Alarm, timestamp time.Time)
	SendCorrelation(timestamp time.Time, child types.Alarm)
	SendUserActivity(timestamp time.Time, username string, value int64)
	SendPbhEnter(alarm *types.Alarm, entity types.Entity)
	SendPbhLeave(entity types.Entity, timestamp time.Time, prevCanonicalType string, prevTimestamp time.Time)
	SendPbhLeaveAndEnter(alarm *types.Alarm, entity types.Entity, prevCanonicalType string, prevTimestamp time.Time)
	SendUpdateState(alarm types.Alarm, entity types.Entity, previousState types.CpsNumber)
}

type nullSender struct{}

func NewNullSender() Sender {
	return &nullSender{}
}

func (s *nullSender) Run(_ context.Context) {

}

func (s *nullSender) SendAck(_ types.Alarm, _ string, _ time.Time) {
}

func (s *nullSender) SendCancelAck(_ types.Alarm, _ time.Time) {
}

func (s *nullSender) SendTicket(_ types.Alarm, _ string, _ time.Time) {
}

func (s *nullSender) SendResolve(_ types.Alarm, _ types.Entity, _ time.Time) {
}

func (s *nullSender) SendAutoInstructionStart(_ types.Alarm, _ time.Time) {
}

func (s *nullSender) SendCreate(_ types.Alarm, _ time.Time) {
}

func (s *nullSender) SendCreateAndPbhEnter(_ types.Alarm, _ time.Time) {
}

func (s *nullSender) SendCorrelation(_ time.Time, _ types.Alarm) {
}

func (s *nullSender) SendUserActivity(_ time.Time, _ string, _ int64) {
}

func (s *nullSender) SendPbhEnter(_ *types.Alarm, _ types.Entity) {

}

func (s *nullSender) SendPbhLeave(_ types.Entity, _ time.Time, _ string, _ time.Time) {

}

func (s *nullSender) SendPbhLeaveAndEnter(_ *types.Alarm, _ types.Entity, _ string, _ time.Time) {

}

func (s *nullSender) SendUpdateState(_ types.Alarm, _ types.Entity, _ types.CpsNumber) {

}
