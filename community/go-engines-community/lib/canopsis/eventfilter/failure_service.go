package eventfilter

import (
	"context"
	"math"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

const (
	FailureTypeInvalidPattern = iota
	FailureTypeInvalidTemplate
	FailureTypeExternalDataMongo
	FailureTypeExternalDataApi
	FailureTypeOther
)

type Failure struct {
	ID        string           `bson:"_id" json:"_id"`
	Rule      string           `bson:"rule" json:"rule"`
	Type      int64            `bson:"type" json:"type"`
	Timestamp datetime.CpsTime `bson:"t" json:"t"`
	Message   string           `bson:"message" json:"message"`
	Event     *types.Event     `bson:"event,omitempty" json:"event"`
	Unread    bool             `bson:"unread,omitempty" json:"unread"`
}

type FailureService interface {
	Run(ctx context.Context)
	Add(ruleID string, failureType int64, message string, event *types.Event)
}

func NewFailureService(
	client mongo.DbClient,
	interval time.Duration,
	logger zerolog.Logger,
) FailureService {
	return &failureService{
		collection:     client.Collection(mongo.EventFilterFailureCollection),
		ruleCollection: client.Collection(mongo.EventFilterRuleCollection),
		interval:       interval,
		logger:         logger,
		countsByRule:   make(map[string]int64),
	}
}

type failureService struct {
	collection     mongo.DbCollection
	ruleCollection mongo.DbCollection
	interval       time.Duration
	logger         zerolog.Logger

	insertsAndCountsByRuleMx sync.Mutex
	inserts                  []any
	countsByRule             map[string]int64
}

func (s *failureService) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := s.flush(ctx)
			if err != nil {
				s.logger.Err(err).Msgf("cannot flush event filter failures")
			}
		}
	}
}

func (s *failureService) Add(ruleID string, failureType int64, message string, event *types.Event) {
	s.insertsAndCountsByRuleMx.Lock()
	defer s.insertsAndCountsByRuleMx.Unlock()
	s.inserts = append(s.inserts, Failure{
		ID:        utils.NewID(),
		Rule:      ruleID,
		Type:      failureType,
		Timestamp: datetime.NewCpsTime(),
		Message:   message,
		Event:     event,
		Unread:    true,
	})
	s.countsByRule[ruleID]++
}

func (s *failureService) flush(ctx context.Context) error {
	inserts, countsByRule := s.flushInserts()
	bulkSize := canopsis.DefaultBulkSize
	l := len(inserts)
	bulkCount := int(math.Ceil(float64(l) / float64(bulkSize)))
	for i := 0; i < bulkCount; i++ {
		begin := i * bulkSize
		end := begin + bulkSize
		if end > l {
			end = l
		}
		_, err := s.collection.InsertMany(ctx, inserts[begin:end])
		if err != nil {
			return err
		}
	}

	ruleWriteModels := make([]mongodriver.WriteModel, 0, bulkSize)
	for ruleID, inc := range countsByRule {
		ruleWriteModels = append(ruleWriteModels, mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{"_id": ruleID}).
			SetUpdate(bson.M{"$inc": bson.M{
				"failures_count":        inc,
				"unread_failures_count": inc,
			}}))
		if len(ruleWriteModels) == bulkSize {
			_, err := s.ruleCollection.BulkWrite(ctx, ruleWriteModels)
			if err != nil {
				return err
			}

			ruleWriteModels = ruleWriteModels[:0]
		}
	}

	if len(ruleWriteModels) > 0 {
		_, err := s.ruleCollection.BulkWrite(ctx, ruleWriteModels)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *failureService) flushInserts() ([]any, map[string]int64) {
	s.insertsAndCountsByRuleMx.Lock()
	defer s.insertsAndCountsByRuleMx.Unlock()
	inserts := s.inserts
	countsByRule := s.countsByRule
	s.inserts = make([]any, 0, len(inserts))
	s.countsByRule = make(map[string]int64, len(countsByRule))

	return inserts, countsByRule
}
