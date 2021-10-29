package metrics

//go:generate mockgen -destination=../../../mocks/lib/metrics/sender.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics Sender

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"time"
)

type Sender interface {
	SendAck(ctx context.Context, alarm types.Alarm, author string, timestamp time.Time)
	SendCancelAck(ctx context.Context, alarm types.Alarm, timestamp time.Time)
	SendAssocTicket(ctx context.Context, alarm types.Alarm, author string, timestamp time.Time)
	SendResolve(ctx context.Context, alarm types.Alarm, timestamp time.Time)
	SendAutoInstructionStart(ctx context.Context, alarm types.Alarm, timestamp time.Time)
	SendCreate(ctx context.Context, alarm types.Alarm, timestamp time.Time)
	SendCreateAndPbhEnter(ctx context.Context, alarm types.Alarm, timestamp time.Time)
	SendCorrelation(ctx context.Context, timestamp time.Time, child types.Alarm)
	SendUserLogin(ctx context.Context, timestamp time.Time, username string)
	SendUserActivity(ctx context.Context, timestamp time.Time, username string, value int64)
}

type nullSender struct{}

func NewNullSender() Sender {
	return &nullSender{}
}

func (s *nullSender) SendAck(_ context.Context, _ types.Alarm, _ string, _ time.Time) {
}

func (s *nullSender) SendCancelAck(_ context.Context, _ types.Alarm, _ time.Time) {
}

func (s *nullSender) SendAssocTicket(_ context.Context, _ types.Alarm, _ string, _ time.Time) {
}

func (s *nullSender) SendResolve(_ context.Context, _ types.Alarm, _ time.Time) {
}

func (s *nullSender) SendAutoInstructionStart(_ context.Context, _ types.Alarm, _ time.Time) {
}

func (s *nullSender) SendCreate(_ context.Context, _ types.Alarm, _ time.Time) {
}

func (s *nullSender) SendCreateAndPbhEnter(_ context.Context, _ types.Alarm, _ time.Time) {
}

func (s *nullSender) SendCorrelation(_ context.Context, _ time.Time, _ types.Alarm) {
}

func (s *nullSender) SendUserLogin(_ context.Context, _ time.Time, _ string) {
}

func (s *nullSender) SendUserActivity(_ context.Context, _ time.Time, _ string, _ int64) {
}
