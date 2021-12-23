package idlerule

import (
	"context"
	"encoding/json"
	"errors"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlerule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"sync"
	"time"
)

type Store interface {
	Insert(context.Context, EditRequest) (*idlerule.Rule, error)
	Find(context.Context, FilteredQuery) (*AggregationResult, error)
	GetOneBy(ctx context.Context, id string) (*idlerule.Rule, error)
	Update(context.Context, EditRequest) (*idlerule.Rule, error)
	Delete(ctx context.Context, id string) (bool, error)
	CountByPatterns(ctx context.Context, filter CountByPatternRequest, timeout int, overLimit int) (*CountByPatternResult, error)
}

type store struct {
	collection            mongo.DbCollection
	entityCollection      mongo.DbCollection
	alarmCollection       mongo.DbCollection
	defaultSearchByFields []string
	defaultSortBy         string
}

type countResult struct {
	entities    int64
	alarms      int64
	entitiesErr error
	alarmsErr   error
}

func NewStore(db mongo.DbClient) Store {
	return &store{
		collection:            db.Collection(mongo.IdleRuleMongoCollection),
		entityCollection:      db.Collection(mongo.EntityMongoCollection),
		alarmCollection:       db.Collection(mongo.AlarmMongoCollection),
		defaultSearchByFields: []string{"_id", "name", "description", "author"},
		defaultSortBy:         "created",
	}
}

func (s *store) Find(ctx context.Context, r FilteredQuery) (*AggregationResult, error) {
	pipeline := make([]bson.M, 0)
	filter := common.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	cursor, err := s.collection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		s.getSort(r),
	))

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	res := AggregationResult{}

	if cursor.Next(ctx) {
		err := cursor.Decode(&res)
		if err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (s *store) GetOneBy(ctx context.Context, id string) (*idlerule.Rule, error) {
	res := s.collection.FindOne(ctx, bson.M{"_id": id})
	if err := res.Err(); err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	rule := &idlerule.Rule{}
	err := res.Decode(rule)
	if err != nil {
		return nil, err
	}

	return rule, nil
}

func (s *store) Insert(ctx context.Context, r EditRequest) (*idlerule.Rule, error) {
	now := types.CpsTime{Time: time.Now()}
	rule := transformRequestToModel(r)
	rule.ID = utils.NewID()
	rule.Created = now
	rule.Updated = now

	_, err := s.collection.InsertOne(ctx, rule)
	if err != nil {
		return nil, err
	}

	err = s.updateFollowingPriorities(ctx, rule.ID, rule.Priority)
	if err != nil {
		return nil, err
	}

	return &rule, nil
}

func (s *store) Update(ctx context.Context, r EditRequest) (*idlerule.Rule, error) {
	prevRule, err := s.GetOneBy(ctx, r.ID)
	if err != nil || prevRule == nil {
		return nil, err
	}

	now := types.CpsTime{Time: time.Now()}
	rule := transformRequestToModel(r)
	rule.ID = r.ID
	rule.Created = prevRule.Created
	rule.Updated = now

	_, err = s.collection.UpdateOne(ctx, bson.M{"_id": rule.ID}, bson.M{"$set": rule})
	if err != nil {
		return nil, err
	}

	if prevRule.Priority != rule.Priority {
		err := s.updateFollowingPriorities(ctx, rule.ID, rule.Priority)
		if err != nil {
			return nil, err
		}
	}

	return &rule, nil
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	deleted, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	return deleted > 0, nil
}

func (s store) countEntitiesCollection(ctx context.Context, entityPattern pattern.EntityPatternList) (int64, error) {
	cursor, err := s.entityCollection.Aggregate(ctx, []bson.M{
		{"$match": entityPattern.AsMongoDriverQuery()},
		{"$count": "total_count"},
	})
	if err != nil {
		return 0, err
	}

	defer cursor.Close(ctx)

	ar := AggregationResult{}
	if cursor.Next(ctx) {
		err = cursor.Decode(&ar)
	}

	return ar.GetTotal(), err
}

func (s store) countAlarmsCollection(ctx context.Context, alarmPattern pattern.AlarmPatternList) (int64, error) {
	cursor, err := s.alarmCollection.Aggregate(ctx, []bson.M{
		{"$match": alarmPattern.AsMongoDriverQuery()},
		{"$count": "total_count"},
	})
	if err != nil {
		return 0, err
	}

	defer cursor.Close(ctx)

	ar := AggregationResult{}
	if cursor.Next(ctx) {
		err = cursor.Decode(&ar)
	}

	return ar.GetTotal(), err
}

func (s store) count(ctx context.Context, request CountByPatternRequest, overLimit int) <-chan countResult {
	v := make(chan countResult)

	go func(v chan<- countResult) {
		defer close(v)

		var entitiesCount, alarmCount int64
		var entitiesErr, alarmErr error

		wg := sync.WaitGroup{}
		wg.Add(2)

		ctx, cancel := context.WithCancel(ctx)

		go func(cancel context.CancelFunc) {
			defer wg.Done()

			if !request.EntityPatterns.IsSet() || !request.EntityPatterns.IsValid() || len(request.EntityPatterns.Patterns) == 0 {
				return
			}

			entitiesCount, entitiesErr = s.countEntitiesCollection(ctx, request.EntityPatterns)
			if entitiesErr != nil || entitiesCount > int64(overLimit) {
				cancel()
			}
		}(cancel)

		go func(cancel context.CancelFunc) {
			defer wg.Done()

			if !request.AlarmPatterns.IsSet() || !request.AlarmPatterns.IsValid() || len(request.AlarmPatterns.Patterns) == 0 {
				return
			}

			alarmCount, alarmErr = s.countAlarmsCollection(ctx, request.AlarmPatterns)
			if alarmErr != nil || alarmCount > int64(overLimit) {
				cancel()
			}
		}(cancel)

		wg.Wait()

		res := countResult{
			entities:    entitiesCount,
			alarms:      alarmCount,
			entitiesErr: entitiesErr,
			alarmsErr:   alarmErr,
		}

		v <- res
	}(v)

	return v
}

func (s store) CountByPatterns(ctx context.Context, request CountByPatternRequest, timeout int, overLimit int) (*CountByPatternResult, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	res := <-s.count(ctx, request, overLimit)

	if errors.Is(res.alarmsErr, context.DeadlineExceeded) || errors.Is(res.entitiesErr, context.DeadlineExceeded) {
		return nil, context.DeadlineExceeded
	}

	if res.alarmsErr != nil && !errors.Is(res.alarmsErr, context.Canceled) {
		return nil, res.alarmsErr
	}

	if res.entitiesErr != nil && !errors.Is(res.entitiesErr, context.Canceled) {
		return nil, res.entitiesErr
	}

	return &CountByPatternResult{
		TotalCountEntities: res.entities,
		TotalCountAlarms:   res.alarms,
	}, nil
}

func (s *store) updateFollowingPriorities(ctx context.Context, id string, priority int64) error {
	err := s.collection.FindOne(ctx, bson.M{
		"_id":      bson.M{"$ne": id},
		"priority": priority,
	}).Err()
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil
		}
		return err
	}

	_, err = s.collection.UpdateMany(
		ctx,
		bson.M{
			"_id":      bson.M{"$ne": id},
			"priority": bson.M{"$gte": priority},
		},
		bson.M{"$inc": bson.M{"priority": 1}},
	)

	return err
}

func (s *store) getSort(r FilteredQuery) bson.M {
	sortBy := s.defaultSortBy
	if r.SortBy != "" {
		sortBy = r.SortBy
	}

	if sortBy == "duration" {
		sortBy = "duration.value"
	}

	return common.GetSortQuery(sortBy, r.Sort)
}

func transformRequestToModel(r EditRequest) idlerule.Rule {
	var operation *idlerule.Operation
	if r.Operation != nil {
		var params map[string]interface{}
		b, _ := json.Marshal(r.Operation.Parameters)
		_ = json.Unmarshal(b, &params)
		operation = &idlerule.Operation{
			Type:       r.Operation.Type,
			Parameters: params,
		}
	}

	return idlerule.Rule{
		Name:                 r.Name,
		Description:          r.Description,
		Author:               r.Author,
		Enabled:              *r.Enabled,
		Type:                 r.Type,
		Priority:             *r.Priority,
		Duration:             r.Duration,
		EntityPatterns:       r.EntityPatterns,
		DisableDuringPeriods: r.DisableDuringPeriods,
		AlarmPatterns:        r.AlarmPatterns,
		AlarmCondition:       r.AlarmCondition,
		Operation:            operation,
	}
}
