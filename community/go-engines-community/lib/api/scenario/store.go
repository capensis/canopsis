package scenario

import (
	"cmp"
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/priority"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type Store interface {
	Insert(ctx context.Context, r CreateRequest) (*Scenario, error)
	Find(ctx context.Context, q FilteredQuery) (*AggregationResult, error)
	GetOneBy(ctx context.Context, id string) (*Scenario, error)
	Update(ctx context.Context, r UpdateRequest) (*Scenario, error)
	Delete(ctx context.Context, id, userId string) (bool, error)
}

type store struct {
	dbClient              mongo.DbClient
	collection            mongo.DbCollection
	transformer           ModelTransformer
	authorProvider        author.Provider
	defaultSearchByFields []string
	defaultSortBy         string
}

func NewStore(db mongo.DbClient, authorProvider author.Provider) Store {
	return &store{
		dbClient:              db,
		collection:            db.Collection(mongo.ScenarioMongoCollection),
		transformer:           NewModelTransformer(),
		authorProvider:        authorProvider,
		defaultSearchByFields: []string{"_id", "name", "author.name"},
		defaultSortBy:         "created",
	}
}

func (s *store) Find(ctx context.Context, r FilteredQuery) (*AggregationResult, error) {
	pipeline := s.authorProvider.Pipeline()
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

func (s *store) GetOneBy(ctx context.Context, id string) (*Scenario, error) {
	pipeline := []bson.M{{"$match": bson.M{"_id": id}}}
	pipeline = append(pipeline, getNestedObjectsPipeline()...)
	pipeline = append(pipeline, s.authorProvider.Pipeline()...)

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

func (s *store) Insert(ctx context.Context, r CreateRequest) (*Scenario, error) {
	now := datetime.NewCpsTime()
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

		err = priority.UpdateFollowing(ctx, s.collection, model.ID, model.Priority)
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

func (s *store) Update(ctx context.Context, r UpdateRequest) (*Scenario, error) {
	now := datetime.NewCpsTime()
	model := s.transformer.TransformEditRequestToModel(r.EditRequest)
	model.Updated = now

	var result *Scenario

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		result = nil
		res, err := s.collection.UpdateOne(ctx, bson.M{"_id": r.ID}, bson.M{"$set": model})
		if err != nil || res.MatchedCount == 0 {
			return err
		}

		err = priority.UpdateFollowing(ctx, s.collection, r.ID, r.Priority)
		if err != nil {
			return err
		}

		result, err = s.GetOneBy(ctx, r.ID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return result, nil
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
	if sortBy == "delay" {
		sortBy = "delay.value"
	}

	return common.GetSortQuery(sortBy, r.Sort)
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
