package scenario

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Store is an interface for scenarios storage
type Store interface {
	Insert(EditRequest) (*Scenario, error)
	Find(FilteredQuery) (*AggregationResult, error)
	GetOneBy(id string) (*Scenario, error)
	Update(EditRequest) (*Scenario, error)
	Delete(id string) (bool, error)
}

type store struct {
	db          mongo.DbClient
	collection  mongo.DbCollection
	transformer ModelTransformer
}

// NewStore instantiates scenario store.
func NewStore(db mongo.DbClient) Store {
	return &store{
		db:          db,
		collection:  db.Collection(mongo.ScenarioMongoCollection),
		transformer: NewModelTransformer(),
	}
}

// Find scenarios according to query.
func (s *store) Find(r FilteredQuery) (*AggregationResult, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filter := bson.M{}

	if r.Search != "" {
		searchRegexp := primitive.Regex{
			Pattern: fmt.Sprintf(".*%s.*", r.Search),
			Options: "i",
		}

		filter["$or"] = []bson.M{
			{"name": searchRegexp},
			{"author": searchRegexp},
		}
	}

	pipeline := []bson.M{{"$match": filter}}
	pipeline = append(pipeline, getNestedObjectsPipeline()...)
	cursor, err := s.collection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		getSort(r),
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
func (s *store) GetOneBy(id string) (*Scenario, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

// Create new scenario.
func (s *store) Insert(r EditRequest) (*Scenario, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	now := types.CpsTime{Time: time.Now()}
	model := s.transformer.TransformEditRequestToModel(r)
	model.ID = utils.NewID()
	model.Created = now
	model.Updated = now

	_, err := s.collection.InsertOne(ctx, model)
	if err != nil {
		return nil, err
	}

	return s.GetOneBy(model.ID)
}

// Update scenario.
func (s *store) Update(r EditRequest) (*Scenario, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	now := types.CpsTime{Time: time.Now()}
	model := s.transformer.TransformEditRequestToModel(r)
	model.Updated = now

	res, err := s.collection.UpdateOne(ctx, bson.M{"_id": r.ID}, bson.M{"$set": model})
	if err != nil {
		return nil, err
	}

	if res.MatchedCount == 0 {
		return nil, nil
	}

	return s.GetOneBy(r.ID)
}

// Delete scenario by id
func (s *store) Delete(id string) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

func getSort(r FilteredQuery) bson.M {
	sortBy := "created"
	if r.SortBy != "" {
		sortBy = r.SortBy
	}

	if sortBy == "delay" {
		sortBy = "delay.seconds"
	}

	return common.GetSortQuery(sortBy, r.Sort)
}
