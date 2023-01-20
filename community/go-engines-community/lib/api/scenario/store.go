package scenario

import (
	"context"
	"errors"
	"math"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

// Store is an interface for scenarios storage
type Store interface {
	Insert(ctx context.Context, r CreateRequest) (*Scenario, error)
	Find(ctx context.Context, q FilteredQuery) (*AggregationResult, error)
	GetOneBy(ctx context.Context, id string) (*Scenario, error)
	IsPriorityValid(ctx context.Context, priority int) (bool, error)
	Update(ctx context.Context, r UpdateRequest) (*Scenario, error)
	Delete(ctx context.Context, id string) (bool, error)

	BulkInsert(ctx context.Context, requests []CreateRequest) error
	BulkUpdate(ctx context.Context, requests []BulkUpdateRequestItem) error
	BulkDelete(ctx context.Context, ids []string) error
}

type store struct {
	dbClient              mongo.DbClient
	collection            mongo.DbCollection
	transformer           ModelTransformer
	defaultSearchByFields []string
	defaultSortBy         string
}

// NewStore instantiates scenario store.
func NewStore(db mongo.DbClient) Store {
	return &store{
		dbClient:              db,
		collection:            db.Collection(mongo.ScenarioMongoCollection),
		transformer:           NewModelTransformer(),
		defaultSearchByFields: []string{"_id", "name", "author.name"},
		defaultSortBy:         "created",
	}
}

// Find scenarios according to query.
func (s *store) Find(ctx context.Context, r FilteredQuery) (*AggregationResult, error) {
	pipeline := author.Pipeline()
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
	pipeline = append(pipeline, author.Pipeline()...)

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
	err := s.collection.FindOne(ctx, bson.M{action.PriorityField: priority}).Err()
	if errors.Is(err, mongodriver.ErrNoDocuments) {
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
	var result *Scenario

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		result = nil
		_, err := s.collection.InsertOne(ctx, model)
		if err != nil {
			return err
		}

		result, err = s.GetOneBy(ctx, model.ID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Update scenario.
func (s *store) Update(ctx context.Context, r UpdateRequest) (*Scenario, error) {
	now := types.CpsTime{Time: time.Now()}
	model := s.transformer.TransformEditRequestToModel(r.EditRequest)
	model.Updated = now

	var result *Scenario

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		result = nil
		res, err := s.collection.UpdateOne(ctx, bson.M{"_id": r.ID}, bson.M{"$set": model})
		if err != nil {
			return err
		}

		if res.MatchedCount == 0 {
			return nil
		}
		result, err = s.GetOneBy(ctx, r.ID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Delete scenario by id
func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	deleted, err := s.collection.DeleteMany(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	return deleted > 0, nil
}

func (s *store) BulkInsert(ctx context.Context, requests []CreateRequest) error {
	var err error
	writeModels := make([]mongodriver.WriteModel, 0, int(math.Min(float64(canopsis.DefaultBulkSize), float64(len(requests)))))
	now := types.CpsTime{Time: time.Now()}

	for _, r := range requests {
		model := s.transformer.TransformEditRequestToModel(r.EditRequest)
		if r.ID == "" {
			r.ID = utils.NewID()
		}

		model.ID = r.ID

		model.Created = now
		model.Updated = now

		writeModels = append(
			writeModels,
			mongodriver.NewInsertOneModel().SetDocument(model),
		)

		if len(writeModels) == canopsis.DefaultBulkSize {
			_, err = s.collection.BulkWrite(ctx, writeModels)
			if err != nil {
				return err
			}

			writeModels = writeModels[:0]
		}
	}

	if len(writeModels) > 0 {
		_, err = s.collection.BulkWrite(ctx, writeModels)
	}

	return err
}

func (s *store) BulkUpdate(ctx context.Context, requests []BulkUpdateRequestItem) error {
	var err error
	writeModels := make([]mongodriver.WriteModel, 0, int(math.Min(float64(canopsis.DefaultBulkSize), float64(len(requests)))))
	now := types.CpsTime{Time: time.Now()}

	for _, r := range requests {
		model := s.transformer.TransformEditRequestToModel(r.EditRequest)
		model.Updated = now

		writeModels = append(
			writeModels,
			mongodriver.
				NewUpdateOneModel().
				SetFilter(bson.M{"_id": r.ID}).
				SetUpdate(bson.M{"$set": model}),
		)

		if len(writeModels) == canopsis.DefaultBulkSize {
			_, err = s.collection.BulkWrite(ctx, writeModels)
			if err != nil {
				return err
			}

			writeModels = writeModels[:0]
		}
	}

	if len(writeModels) > 0 {
		_, err = s.collection.BulkWrite(ctx, writeModels)
	}

	return err
}

func (s *store) BulkDelete(ctx context.Context, ids []string) error {
	_, err := s.collection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})

	return err
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
		sortBy = "delay.value"
	}

	return common.GetSortQuery(sortBy, r.Sort)
}
