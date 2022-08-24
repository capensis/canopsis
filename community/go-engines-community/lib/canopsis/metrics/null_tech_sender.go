package metrics

import (
	"context"
	"time"
)

type nullTechSender struct{}

func NewNullTechMetricsSender() TechSender {
	return &nullTechSender{}
}

func (s *nullTechSender) SendFifoQueue(_ context.Context, _ time.Time, _ int64) {

}

func (s *nullTechSender) SendFifoEventBatch(_ context.Context, _ []FifoEventMetric) {

}

func (s *nullTechSender) SendCheEventBatch(_ context.Context, _ []CheEventMetric) {

}

func (s *nullTechSender) SendAxeEventBatch(_ context.Context, _ []AxeEventMetric) {

}

func (s *nullTechSender) SendAxePeriodical(_ context.Context, _ time.Time, _ int64) {

}

func (s *nullTechSender) SendPBehaviorPeriodical(_ context.Context, _ time.Time, _ int64) {

}

func (s *nullTechSender) SendCheEntityInfo(_ context.Context, _ time.Time, name string) {

}
