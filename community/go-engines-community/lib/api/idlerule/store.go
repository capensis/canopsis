package idlerule

import (
	"cmp"
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/dbvalidation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/mongoquery"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/patternfields"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/priority"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlerule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
)

type Store interface {
	Insert(context.Context, CreateRequest) (*Response, error)
	Find(context.Context, FilteredQuery) (*AggregationResult, error)
	GetOneBy(ctx context.Context, id string) (*Response, error)
	Update(context.Context, UpdateRequest) (*Response, error)
	Delete(ctx context.Context, id, userID string) (bool, error)
}

type store struct {
	dbClient              mongo.DbClient
	collection            mongo.DbCollection
	pbhTypeCollection     mongo.DbCollection
	pbhReasonCollection   mongo.DbCollection
	authorProvider        author.Provider
	transformer           patternfields.Transformer
	defaultSearchByFields []string
	defaultSortBy         string
	dupErrorParser        validation.DuplicateErrorParser
}

func NewStore(db mongo.DbClient, authorProvider author.Provider, transformer patternfields.Transformer) Store {
	return &store{
		dbClient:              db,
		collection:            db.Collection(mongo.IdleRuleMongoCollection),
		pbhTypeCollection:     db.Collection(mongo.PbehaviorTypeMongoCollection),
		pbhReasonCollection:   db.Collection(mongo.PbehaviorReasonMongoCollection),
		authorProvider:        authorProvider,
		transformer:           transformer,
		defaultSearchByFields: []string{"_id", "name", "description", "author.name"},
		defaultSortBy:         "created",
		dupErrorParser:        validation.NewDuplicateErrorParser(),
	}
}

func (s *store) Find(ctx context.Context, r FilteredQuery) (*AggregationResult, error) {
	pipeline := s.authorProvider.Pipeline()
	filter := mongoquery.GetSearchQuery(r.Search, s.defaultSearchByFields)
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

func (s *store) GetOneBy(ctx context.Context, id string) (*Response, error) {
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
		rule := &Response{}
		err = cursor.Decode(rule)
		if err != nil {
			return nil, err
		}

		return rule, nil
	}

	return nil, nil
}

func (s *store) Insert(ctx context.Context, r CreateRequest) (*Response, error) {
	now := datetime.NewCpsTime()
	rule := transformRequestToModel(r.EditRequest)

	rule.ID = cmp.Or(r.ID, utils.NewID())
	rule.Created = now
	rule.Updated = now

	var idleRule *Response
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		idleRule = nil

		if r.Operation != nil {
			err := dbvalidation.ValidateExist(ctx, s.pbhTypeCollection, r, "Operation.Parameters.Type", r.Operation.Parameters.Type)
			if err != nil {
				return err
			}

			err = dbvalidation.ValidateExist(ctx, s.pbhReasonCollection, r, "Operation.Parameters.Reason", r.Operation.Parameters.Reason)
			if err != nil {
				return err
			}
		}

		err := s.transformPatternRequestsToModel(ctx, r.EditRequest, &rule)
		if err != nil {
			return err
		}

		_, err = s.collection.InsertOne(ctx, rule)
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err, rule)
			}

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

func (s *store) Update(ctx context.Context, r UpdateRequest) (*Response, error) {
	model := transformRequestToModel(r.EditRequest)
	model.ID = r.ID
	model.Updated = datetime.NewCpsTime()

	var idleRule *Response
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		idleRule = nil

		err := dbvalidation.ValidateExist(ctx, s.pbhTypeCollection, r, "Operation.Parameters.Type", r.Operation.Parameters.Type)
		if err != nil {
			return err
		}

		err = dbvalidation.ValidateExist(ctx, s.pbhReasonCollection, r, "Operation.Parameters.Reason", r.Operation.Parameters.Reason)
		if err != nil {
			return err
		}

		err = s.transformPatternRequestsToModel(ctx, r.EditRequest, &model)
		if err != nil {
			return err
		}

		_, err = s.collection.UpdateOne(ctx, bson.M{"_id": model.ID}, bson.M{"$set": model})
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err, model)
			}

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

func (s *store) Delete(ctx context.Context, id, userID string) (bool, error) {
	var deleted int64

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		deleted = 0

		// required to get the author in action log listener.
		res, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"author": userID}})
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

	return mongoquery.GetSortQuery(sortBy, r.Sort)
}

func (s *store) transformPatternRequestsToModel(ctx context.Context, r EditRequest, model *idlerule.Rule) (err error) {
	model.AlarmPatternFields, model.EntityPatternFields, model.Aliases, err = s.transformer.TransformAlarmAndEntityRequest(ctx, r.AlarmRequest, r.EntityRequest, r, s.collection.Name())

	return err
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
