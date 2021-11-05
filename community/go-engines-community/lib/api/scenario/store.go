package scenario

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// Store is an interface for scenarios storage
type Store interface {
	Insert(ctx context.Context, r CreateRequest) (*Scenario, error)
	Find(ctx context.Context, q FilteredQuery) (*AggregationResult, error)
	GetOneBy(ctx context.Context, id string) (*Scenario, error)
	IsPriorityValid(ctx context.Context, priority int) (bool, error)
	Update(ctx context.Context, r UpdateRequest) (*Scenario, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type store struct {
	db                    mongo.DbClient
	collection            mongo.DbCollection
	transformer           ModelTransformer
	defaultSearchByFields []string
	defaultSortBy         string
}

// NewStore instantiates scenario store.
func NewStore(db mongo.DbClient) Store {
	return &store{
		db:                    db,
		collection:            db.Collection(mongo.ScenarioMongoCollection),
		transformer:           NewModelTransformer(),
		defaultSearchByFields: []string{"_id", "name", "author"},
		defaultSortBy:         "created",
	}
}

// Find scenarios according to query.
func (s *store) Find(ctx context.Context, r FilteredQuery) (*AggregationResult, error) {
	pipeline := make([]bson.M, 0)
	filter := common.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	pipeline = append(pipeline, getNestedObjectsPipeline()...)
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

// GetOneBy scenario by id.
func (s *store) GetOneBy(ctx context.Context, id string) (*Scenario, error) {
	pipeline := []bson.M{{"$match": bson.M{"_id": id}}}
	pipeline = append(pipeline, getNestedObjectsPipeline()...)

	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if cursor.Next(ctx) {
		res := &Scenario{}
		err := cursor.Decode(res)
		if err != nil {
			return nil, err
		}

		return res, nil
	}

	return nil, nil
}

func (s *store) IsPriorityValid(ctx context.Context, priority int) (bool, error) {
	err := s.collection.FindOne(ctx, bson.M{"priority": priority}).Err()
	if err == nil {
		return false, nil
	}

	if err == mongodriver.ErrNoDocuments {
		return true, nil
	}

	return false, err
}

// Create new scenario.
func (s *store) Insert(ctx context.Context, r CreateRequest) (*Scenario, error) {
	now := types.CpsTime{Time: time.Now()}
	model := s.transformer.TransformEditRequestToModel(r.EditRequest)

	if r.ID == "" {
		r.ID = utils.NewID()
	}

	model.ID = r.ID

	model.Created = now
	model.Updated = now

	_, err := s.collection.InsertOne(ctx, model)
	if err != nil {
		return nil, err
	}

	return s.GetOneBy(ctx, model.ID)
}

// Update scenario.
func (s *store) Update(ctx context.Context, r UpdateRequest) (*Scenario, error) {
	now := types.CpsTime{Time: time.Now()}
	model := s.transformer.TransformEditRequestToModel(r.EditRequest)
	model.Updated = now

	res, err := s.collection.UpdateOne(ctx, bson.M{"_id": r.ID}, bson.M{"$set": model})
	if err != nil {
		return nil, err
	}

	if res.MatchedCount == 0 {
		return nil, nil
	}

	return s.GetOneBy(ctx, r.ID)
}

// Delete scenario by id
func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	deleted, err := s.collection.DeleteMany(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	return deleted > 0, nil
}

func getNestedObjectsPipeline() []bson.M {
	return []bson.M{
		{"$unwind": bson.M{
			"path":                       "$actions",
			"preserveNullAndEmptyArrays": true,
			"includeArrayIndex":          "action_index",
		}},
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorTypeMongoCollection,
			"localField":   "actions.parameters.type",
			"foreignField": "_id",
			"as":           "actions.parameters.type",
		}},
		{"$unwind": bson.M{"path": "$actions.parameters.type", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorReasonMongoCollection,
			"localField":   "actions.parameters.reason",
			"foreignField": "_id",
			"as":           "actions.parameters.reason",
		}},
		{"$unwind": bson.M{"path": "$actions.parameters.reason", "preserveNullAndEmptyArrays": true}},
		{"$project": bson.M{
			"actions.parameters.reason.created": 0,
		}},
		{"$sort": bson.M{"action_index": 1}},
		{"$group": bson.M{
			"_id":     "$_id",
			"data":    bson.M{"$first": "$$ROOT"},
			"actions": bson.M{"$push": "$actions"},
		}},
		{"$replaceRoot": bson.M{
			"newRoot": bson.M{"$mergeObjects": bson.A{
				"$data",
				bson.M{"actions": "$actions"},
			}},
		}},
	}
}

func (s *store) getSort(r FilteredQuery) bson.M {
	sortBy := s.defaultSortBy
	if r.SortBy != "" {
		sortBy = r.SortBy
	}

	if sortBy == "delay" {
		sortBy = "delay.seconds"
	}

	return common.GetSortQuery(sortBy, r.Sort)
}
