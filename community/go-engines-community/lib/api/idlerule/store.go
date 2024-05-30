package idlerule

import (
	"cmp"
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/priority"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlerule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type Store interface {
	Insert(context.Context, CreateRequest) (*Rule, error)
	Find(context.Context, FilteredQuery) (*AggregationResult, error)
	GetOneBy(ctx context.Context, id string) (*Rule, error)
	Update(context.Context, UpdateRequest) (*Rule, error)
	Delete(ctx context.Context, id, userId string) (bool, error)
}

type store struct {
	dbClient              mongo.DbClient
	collection            mongo.DbCollection
	authorProvider        author.Provider
	defaultSearchByFields []string
	defaultSortBy         string
}

func NewStore(db mongo.DbClient, authorProvider author.Provider) Store {
	return &store{
		dbClient:              db,
		collection:            db.Collection(mongo.IdleRuleMongoCollection),
		authorProvider:        authorProvider,
		defaultSearchByFields: []string{"_id", "name", "description", "author.name"},
		defaultSortBy:         "created",
	}
}

func (s *store) Find(ctx context.Context, r FilteredQuery) (*AggregationResult, error) {
	pipeline := s.authorProvider.Pipeline()
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
	pipeline = append(pipeline, s.authorProvider.Pipeline()...)
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
	now := datetime.NewCpsTime()
	rule := transformRequestToModel(r.EditRequest)

	rule.ID = cmp.Or(r.ID, utils.NewID())
	rule.Created = now
	rule.Updated = now

	var idleRule *Rule
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		idleRule = nil
		_, err := s.collection.InsertOne(ctx, rule)
		if err != nil {
			return err
		}

		err = priority.UpdateFollowing(ctx, s.collection, rule.ID, rule.Priority)
		if err != nil {
			return err
		}
		idleRule, err = s.GetOneBy(ctx, rule.ID)
		return err
	})

	if err != nil {
		return nil, err
	}

	return idleRule, nil
}

func (s *store) Update(ctx context.Context, r UpdateRequest) (*Rule, error) {
	model := transformRequestToModel(r.EditRequest)
	model.ID = r.ID
	model.Updated = datetime.NewCpsTime()

	var idleRule *Rule
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		idleRule = nil

		_, err := s.collection.UpdateOne(ctx, bson.M{"_id": model.ID}, bson.M{"$set": model})
		if err != nil {
			return err
		}

		err = priority.UpdateFollowing(ctx, s.collection, model.ID, model.Priority)
		if err != nil {
			return err
		}

		idleRule, err = s.GetOneBy(ctx, r.ID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return idleRule, nil
}

func (s *store) Delete(ctx context.Context, id, userId string) (bool, error) {
	var deleted int64

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		deleted = 0

		// required to get the author in action log listener.
		res, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"author": userId}})
		if err != nil || res.MatchedCount == 0 {
			return err
		}

		deleted, err = s.collection.DeleteOne(ctx, bson.M{"_id": id})
		return err
	})

	return deleted > 0, err
}

func (s *store) getSort(r FilteredQuery) bson.M {
	sortBy := cmp.Or(r.SortBy, s.defaultSortBy)
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
		Priority:             r.Priority,
		Duration:             r.Duration,
		Comment:              r.Comment,
		DisableDuringPeriods: r.DisableDuringPeriods,
		AlarmCondition:       r.AlarmCondition,
		Operation:            operation,
		AlarmPatternFields: r.AlarmPatternFieldsRequest.ToModelWithoutFields(
			common.GetForbiddenFieldsInAlarmPattern(mongo.IdleRuleMongoCollection),
			common.GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(mongo.IdleRuleMongoCollection),
		),
		EntityPatternFields: r.EntityPatternFieldsRequest.ToModelWithoutFields(
			common.GetForbiddenFieldsInEntityPattern(mongo.IdleRuleMongoCollection),
		),
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
