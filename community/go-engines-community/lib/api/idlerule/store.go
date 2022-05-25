package idlerule

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlerule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Store interface {
	Insert(context.Context, CreateRequest) (*Rule, error)
	Find(context.Context, FilteredQuery) (*AggregationResult, error)
	GetOneBy(ctx context.Context, id string) (*Rule, error)
	Update(context.Context, UpdateRequest) (*Rule, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type store struct {
	collection            mongo.DbCollection
	entityCollection      mongo.DbCollection
	alarmCollection       mongo.DbCollection
	defaultSearchByFields []string
	defaultSortBy         string
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
		getNestedObjectsPipeline(),
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

func (s *store) GetOneBy(ctx context.Context, id string) (*Rule, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"_id": id}},
	}
	pipeline = append(pipeline, getNestedObjectsPipeline()...)
	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		rule := &Rule{}
		err = cursor.Decode(rule)
		if err != nil {
			return nil, err
		}

		return rule, nil
	}

	return nil, nil
}

func (s *store) Insert(ctx context.Context, r CreateRequest) (*Rule, error) {
	now := types.CpsTime{Time: time.Now()}
	rule := transformRequestToModel(r.EditRequest)
	if r.ID == "" {
		r.ID = utils.NewID()
	}

	rule.ID = r.ID
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

	return s.GetOneBy(ctx, r.ID)
}

func (s *store) Update(ctx context.Context, r UpdateRequest) (*Rule, error) {
	prevRule, err := s.GetOneBy(ctx, r.ID)
	if err != nil || prevRule == nil {
		return nil, err
	}

	now := types.CpsTime{Time: time.Now()}
	rule := transformRequestToModel(r.EditRequest)
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

	return s.GetOneBy(ctx, r.ID)
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	deleted, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	return deleted > 0, nil
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
		operation = &idlerule.Operation{
			Type:       r.Operation.Type,
			Parameters: r.Operation.Parameters,
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

func getNestedObjectsPipeline() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorTypeMongoCollection,
			"localField":   "operation.parameters.type",
			"foreignField": "_id",
			"as":           "operation.parameters.type",
		}},
		{"$unwind": bson.M{"path": "$operation.parameters.type", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorReasonMongoCollection,
			"localField":   "operation.parameters.reason",
			"foreignField": "_id",
			"as":           "operation.parameters.reason",
		}},
		{"$unwind": bson.M{"path": "$operation.parameters.reason", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"operation": bson.M{"$cond": bson.M{
				"if":   "$operation.type",
				"then": "$operation",
				"else": nil,
			}},
		}},
	}
}
