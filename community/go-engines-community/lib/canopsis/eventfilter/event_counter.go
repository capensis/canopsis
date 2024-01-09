package eventfilter

import (
	"context"
	"math"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type EventCounter interface {
	Run(ctx context.Context)
	Add(id string, lastUpdated datetime.CpsTime)
}

func NewEventCounter(
	client mongo.DbClient,
	interval time.Duration,
	logger zerolog.Logger,
) EventCounter {
	return &eventCounter{
		collection: client.Collection(mongo.EventFilterRuleCollection),
		interval:   interval,
		logger:     logger,
		counts:     make(map[string]count),
	}
}

type eventCounter struct {
	collection mongo.DbCollection
	interval   time.Duration
	logger     zerolog.Logger
	countsMx   sync.Mutex
	counts     map[string]count
}

type count struct {
	LastUpdated datetime.CpsTime
	Count       int64
}

func (s *eventCounter) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := s.flush(ctx)
			if err != nil {
				s.logger.Err(err).Msgf("cannot flush events counters")
			}
		}
	}
}

func (s *eventCounter) Add(id string, lastUpdated datetime.CpsTime) {
	s.countsMx.Lock()
	defer s.countsMx.Unlock()
	v, exists := s.counts[id]
	if exists && v.LastUpdated.Unix() == lastUpdated.Unix() {
		v.Count++
	} else if !exists || v.LastUpdated.Unix() < lastUpdated.Unix() {
		v = count{
			LastUpdated: lastUpdated,
			Count:       1,
		}
	}

	s.counts[id] = v
}

func (s *eventCounter) flush(ctx context.Context) error {
	counts := s.flushCounts()
	bulkSize := canopsis.DefaultBulkSize
	writeModels := make([]mongodriver.WriteModel, 0, int(math.Min(float64(bulkSize), float64(len(counts)))))
	for id, c := range counts {
		writeModels = append(writeModels, mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{"_id": id, "updated": c.LastUpdated}).
			SetUpdate(bson.M{"$inc": bson.M{"events_count": c.Count}}))
		if len(writeModels) == bulkSize {
			_, err := s.collection.BulkWrite(ctx, writeModels)
			if err != nil {
				return err
			}

			writeModels = writeModels[:0]
		}
	}

	if len(writeModels) > 0 {
		_, err := s.collection.BulkWrite(ctx, writeModels)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *eventCounter) flushCounts() map[string]count {
	s.countsMx.Lock()
	defer s.countsMx.Unlock()
	counts := s.counts
	s.counts = make(map[string]count, len(counts))
	return counts
}
