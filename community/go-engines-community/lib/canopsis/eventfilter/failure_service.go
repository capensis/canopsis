package eventfilter

import (
	"context"
	"math"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/usernotification"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	FailureTypeInvalidPattern = iota
	FailureTypeInvalidTemplate
	FailureTypeExternalDataTable
	FailureTypeExternalDataAPI
	FailureTypeOther
)

type Failure struct {
	ID          string           `bson:"_id" json:"_id"`
	Rule        string           `bson:"rule" json:"rule"`
	RuleUpdated datetime.CpsTime `bson:"rule_updated" json:"rule_updated"`
	Type        int64            `bson:"type" json:"type"`
	Timestamp   datetime.CpsTime `bson:"t" json:"t"`
	Message     string           `bson:"message" json:"message"`
	Event       *types.Event     `bson:"event,omitempty" json:"event"`
	Unread      bool             `bson:"unread,omitempty" json:"unread"`
}

type FailureService interface {
	Run(ctx context.Context)
	Add(ruleID, ruleDesc string, ruleUpdated datetime.CpsTime, failureType int64, message string, event *types.Event)
}

func NewFailureService(
	client mongo.DbClient,
	notificationStore usernotification.Store,
	interval time.Duration,
	permToNotify string,
	logger zerolog.Logger,
) FailureService {
	return &failureService{
		collection:        client.Collection(mongo.EventFilterFailureCollection),
		ruleCollection:    client.Collection(mongo.EventFilterRuleCollection),
		roleCollection:    client.Collection(mongo.RoleCollection),
		notificationStore: notificationStore,
		interval:          interval,
		permToNotify:      permToNotify,
		logger:            logger,
		countsByRule:      make(map[string]map[int64]int64),
		failedRules:       make(map[string]failedRule),
	}
}

type failureService struct {
	collection        mongo.DbCollection
	ruleCollection    mongo.DbCollection
	notificationStore usernotification.Store
	roleCollection    mongo.DbCollection
	interval          time.Duration
	permToNotify      string
	logger            zerolog.Logger

	dataMx       sync.Mutex
	inserts      []any
	countsByRule map[string]map[int64]int64
	failedRules  map[string]failedRule
}

type failedRule struct {
	Timestamp   datetime.CpsTime
	RuleUpdated datetime.CpsTime
	Description string
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

func (s *failureService) Add(ruleID, ruleDesc string, ruleUpdated datetime.CpsTime, failureType int64, message string, event *types.Event) {
	s.dataMx.Lock()
	defer s.dataMx.Unlock()
	now := datetime.NewCpsTime()
	s.inserts = append(s.inserts, Failure{
		ID:          utils.NewID(),
		Rule:        ruleID,
		RuleUpdated: ruleUpdated,
		Type:        failureType,
		Timestamp:   now,
		Message:     message,
		Event:       event,
		Unread:      true,
	})
	if _, ok := s.countsByRule[ruleID]; !ok {
		s.countsByRule[ruleID] = make(map[int64]int64, 1)
	}

	s.countsByRule[ruleID][ruleUpdated.Unix()]++
	s.failedRules[ruleID] = failedRule{
		Timestamp:   now,
		RuleUpdated: ruleUpdated,
		Description: ruleDesc,
	}
}

func (s *failureService) flush(ctx context.Context) error {
	inserts, countsByRule, failedRules := s.flushData()
	bulkSize := canopsis.DefaultBulkSize
	err := s.insertFailures(ctx, inserts, bulkSize)
	if err != nil {
		return err
	}

	err = s.updateRuleCounts(ctx, countsByRule, bulkSize)
	if err != nil {
		return err
	}

	err = s.updateRuleNotifications(ctx, failedRules)
	if err != nil {
		return err
	}

	return nil
}

func (s *failureService) insertFailures(ctx context.Context, inserts []any, bulkSize int) error {
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

func (s *failureService) updateRuleCounts(ctx context.Context, countsByRule map[string]map[int64]int64, bulkSize int) error {
	writeModels := make([]mongodriver.WriteModel, 0, bulkSize)
	for ruleID, c := range countsByRule {
		for ruleUpdated, inc := range c {
			writeModels = append(writeModels, mongodriver.NewUpdateOneModel().
				SetFilter(bson.M{
					"_id":     ruleID,
					"updated": ruleUpdated,
				}).
				SetUpdate(bson.M{"$inc": bson.M{
					"failures_count":        inc,
					"unread_failures_count": inc,
				}}))
			if len(writeModels) == bulkSize {
				_, err := s.ruleCollection.BulkWrite(ctx, writeModels)
				if err != nil {
					return err
				}

				writeModels = writeModels[:0]
			}
		}
	}

	if len(writeModels) > 0 {
		_, err := s.ruleCollection.BulkWrite(ctx, writeModels)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *failureService) updateRuleNotifications(ctx context.Context, failedRules map[string]failedRule) error {
	if len(failedRules) == 0 {
		return nil
	}

	roleIDs, err := s.findRoles(ctx)
	if err != nil || len(roleIDs) == 0 {
		return err
	}

	for ruleID, r := range failedRules {
		s.notificationStore.AddForEventFilterFailure(r.Timestamp, ruleID, r.Description, r.RuleUpdated, roleIDs)
	}

	err = s.notificationStore.Flush(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *failureService) findRoles(ctx context.Context) ([]string, error) {
	cursor, err := s.roleCollection.Find(ctx, bson.M{
		"permissions." + s.permToNotify: bson.M{"$ne": nil},
	}, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return nil, err
	}

	roleIDs := make([]string, 0)
	for cursor.Next(ctx) {
		role := struct {
			ID string `bson:"_id"`
		}{}
		err = cursor.Decode(&role)
		if err != nil {
			return nil, err
		}

		roleIDs = append(roleIDs, role.ID)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	if err = cursor.Close(ctx); err != nil {
		return nil, err
	}

	return roleIDs, nil
}

func (s *failureService) flushData() ([]any, map[string]map[int64]int64, map[string]failedRule) {
	s.dataMx.Lock()
	defer s.dataMx.Unlock()
	inserts := s.inserts
	countsByRule := s.countsByRule
	failedRules := s.failedRules
	s.inserts = make([]any, 0, len(inserts))
	s.countsByRule = make(map[string]map[int64]int64, len(countsByRule))
	s.failedRules = make(map[string]failedRule, len(failedRules))

	return inserts, countsByRule, failedRules
}
