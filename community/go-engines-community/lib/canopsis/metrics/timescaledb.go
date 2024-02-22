package metrics

import (
	"context"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

type timescaleDBSender struct {
	pgPoolProvider postgres.PoolProvider
	configProvider config.MetricsConfigProvider
	logger         zerolog.Logger

	batchesMx sync.Mutex
	batches   map[string][]batchItem
}

type batchItem struct {
	query     string
	arguments []any
}

func NewTimescaleDBSender(
	pgPoolProvider postgres.PoolProvider,
	configProvider config.MetricsConfigProvider,
	logger zerolog.Logger,
) Sender {
	return &timescaleDBSender{
		configProvider: configProvider,
		pgPoolProvider: pgPoolProvider,
		logger:         logger,

		batches: make(map[string][]batchItem),
	}
}

func (s *timescaleDBSender) Run(ctx context.Context) {
	ticker := time.NewTicker(s.configProvider.Get().FlushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.send(ctx)
			ticker.Reset(s.configProvider.Get().FlushInterval)
		}
	}
}

func (s *timescaleDBSender) SendAck(_ types.Alarm, _ string, _ time.Time) {
}

func (s *timescaleDBSender) SendCancelAck(_ types.Alarm, _ time.Time) {
}

func (s *timescaleDBSender) SendTicket(_ types.Alarm, _ string, _ time.Time) {
}

func (s *timescaleDBSender) SendResolve(_ types.Alarm, _ types.Entity, _ time.Time) {
}

func (s *timescaleDBSender) SendAutoInstructionExecutionStart(_ types.Alarm, _ time.Time) {
}

func (s *timescaleDBSender) SendAutoInstructionExecutionForInstruction(_ string, _ time.Time) {
}

func (s *timescaleDBSender) SendAutoInstructionAssignForInstructions(_ []string, _ time.Time) {

}

func (s *timescaleDBSender) SendCreate(_ types.Alarm, _ time.Time) {
}

func (s *timescaleDBSender) SendCreateAndPbhEnter(_ types.Alarm, _ time.Time) {
}

func (s *timescaleDBSender) SendCorrelation(_ time.Time, _ types.Alarm) {
}

func (s *timescaleDBSender) SendUserActivity(_ time.Time, _ string, _ int64) {
}

func (s *timescaleDBSender) SendPbhEnter(_ types.Alarm, _ types.Entity) {

}

func (s *timescaleDBSender) SendPbhLeave(_ types.Entity, _ time.Time, _ string, _ time.Time) {

}

func (s *timescaleDBSender) SendPbhLeaveAndEnter(_ types.Alarm, _ types.Entity, _ string, _ time.Time) {

}

func (s *timescaleDBSender) SendUpdateState(_ types.Alarm, _ types.Entity, _ types.CpsNumber) {

}

func (s *timescaleDBSender) SendInstructionAssignForAlarm(_ string, _ time.Time) {

}

func (s *timescaleDBSender) SendInstructionAssignForAlarms(_ []string, _ time.Time) {

}

func (s *timescaleDBSender) SendInstructionExecutionForAlarm(_ string, _ time.Time) {

}

func (s *timescaleDBSender) SendInstructionAssignForInstruction(_ string, _ time.Time, _ int64) {

}

func (s *timescaleDBSender) SendInstructionAssignForInstructions(_ []string, _ time.Time) {

}

func (s *timescaleDBSender) SendInstructionExecutionForInstruction(_ string, _ time.Time) {

}

func (s *timescaleDBSender) SendNotAckedInHourInc(_ types.Alarm, _ time.Time) {

}

func (s *timescaleDBSender) SendNotAckedInFourHoursInc(_ types.Alarm, _ time.Time) {

}

func (s *timescaleDBSender) SendNotAckedInDayInc(_ types.Alarm, _ time.Time) {

}

func (s *timescaleDBSender) SendNotAckedInHourDec(_ types.Alarm, _ time.Time) {

}

func (s *timescaleDBSender) SendNotAckedInFourHoursDec(_ types.Alarm, _ time.Time) {

}

func (s *timescaleDBSender) SendNotAckedInDayDec(_ types.Alarm, _ time.Time) {

}

func (s *timescaleDBSender) SendRemoveNotAckedMetric(_ types.Alarm, _ time.Time, _ string) {

}

func (s *timescaleDBSender) SendPerfData(_ time.Time, _, _ string, _ float64, _ string) {

}

func (s *timescaleDBSender) SendEventMetrics(_ types.Alarm, _ types.Entity, _ types.AlarmChange, _ time.Time, _, _, _, _ string) {

}

func (s *timescaleDBSender) SendSliMetric(_ time.Time, _ types.Alarm, _ types.Entity) {

}

func (s *timescaleDBSender) SendMessageRate(timestamp time.Time) {
	query := "INSERT INTO " + MessageRate + " (time) VALUES ($1)"
	s.addBatch(MessageRate, batchItem{
		query: query,
		arguments: []any{
			timestamp.UTC(),
		},
	})
}

func (s *timescaleDBSender) send(ctx context.Context) {
	if !s.configProvider.Get().Enabled {
		s.cleanBatches()

		return
	}

	batches := s.flushBatches()
	pgPool, err := s.pgPoolProvider.Get(ctx)
	if err != nil {
		s.logger.Err(err).Msg("cannot connect to postgres")

		return
	}

	batch := &pgx.Batch{}
	numInserts := 0
	for _, items := range batches {
		lastIndex := len(items) - 1
		for i, item := range items {
			batch.Queue(item.query, item.arguments...)
			numInserts++

			if numInserts >= canopsis.DefaultBulkSize || i == lastIndex {
				err = pgPool.SendBatch(ctx, batch)
				if err != nil {
					s.logger.Err(err).Msg("cannot send metrics")

					return
				}

				batch = &pgx.Batch{}
				numInserts = 0
			}
		}
	}
}

func (s *timescaleDBSender) flushBatches() map[string][]batchItem {
	s.batchesMx.Lock()
	defer s.batchesMx.Unlock()
	if len(s.batches) == 0 {
		return nil
	}

	batches := make(map[string][]batchItem, len(s.batches))
	for metricName, items := range s.batches {
		if len(items) == 0 {
			continue
		}

		batches[metricName] = make([]batchItem, 0, len(items))
	}

	res := s.batches
	s.batches = batches

	return res
}

func (s *timescaleDBSender) addBatch(metricName string, item batchItem) {
	s.batchesMx.Lock()
	defer s.batchesMx.Unlock()

	if _, ok := s.batches[metricName]; !ok {
		s.batches[metricName] = make([]batchItem, 0, 1)
	}

	s.batches[metricName] = append(s.batches[metricName], item)
}

func (s *timescaleDBSender) cleanBatches() {
	s.batchesMx.Lock()
	defer s.batchesMx.Unlock()
	clear(s.batches)
}
