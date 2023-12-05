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
		collection: client.Collection(mongo.EventFilterFailureCollection),
		interval:   interval,
		logger:     logger,
	}
}

type failureService struct {
	collection mongo.DbCollection
	interval   time.Duration
	logger     zerolog.Logger
	insertsMx  sync.Mutex
	inserts    []any
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
	s.insertsMx.Lock()
	defer s.insertsMx.Unlock()
	s.inserts = append(s.inserts, Failure{
		ID:        utils.NewID(),
		Rule:      ruleID,
		Type:      failureType,
		Timestamp: datetime.NewCpsTime(),
		Message:   message,
		Event:     event,
		Unread:    true,
	})
}

func (s *failureService) flush(ctx context.Context) error {
	inserts := s.flushInserts()
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

	return nil
}

func (s *failureService) flushInserts() []any {
	s.insertsMx.Lock()
	defer s.insertsMx.Unlock()
	inserts := s.inserts
	s.inserts = nil
	return inserts
}
