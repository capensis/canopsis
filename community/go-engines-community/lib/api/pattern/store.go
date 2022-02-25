package pattern

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type Store interface {
	Insert(ctx context.Context, r EditRequest) (*Response, error)
	GetById(ctx context.Context, id, userId string) (*Response, error)
	Find(ctx context.Context, r ListRequest, userId string) (*AggregationResult, error)
	Update(ctx context.Context, r EditRequest) (*Response, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type store struct {
	dbClient     mongo.DbClient
	dbCollection mongo.DbCollection

	defaultSearchByFields []string
	defaultSortBy         string
}

func NewStore(
	dbClient mongo.DbClient,
) Store {
	return &store{
		dbClient:     dbClient,
		dbCollection: dbClient.Collection(mongo.PatternMongoCollection),

		defaultSearchByFields: []string{"_id", "author", "title"},
		defaultSortBy:         "created",
	}
}

func (s *store) Insert(ctx context.Context, request EditRequest) (*Response, error) {
	now := types.NewCpsTime()
	model := transformRequestToModel(request)
	model.ID = utils.NewID()
	model.Created = now
	model.Updated = now

	_, err := s.dbCollection.InsertOne(ctx, model)
	if err != nil {
		return nil, err
	}

	return s.GetById(ctx, model.ID, model.Author)
}

func (s *store) GetById(ctx context.Context, id, userId string) (*Response, error) {
	pipeline := []bson.M{{"$match": bson.M{
		"_id": id,
		"$or": []bson.M{
			{"author": userId},
			{"is_corporate": true},
		},
	}}}
	pipeline = append(pipeline, getAuthorPipeline()...)
	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		response := Response{}
		err := cursor.Decode(&response)
		if err != nil {
			return nil, err
		}

		return &response, nil
	}

	return nil, nil
}

func (s *store) Find(ctx context.Context, request ListRequest, userId string) (*AggregationResult, error) {
	pipeline := make([]bson.M, 0)

	if request.Corporate == nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"$or": []bson.M{
			{"author": userId},
			{"is_corporate": true},
		}}})
	} else if *request.Corporate {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"is_corporate": true}})
	} else {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"author": userId, "is_corporate": false}})
	}

	filter := common.GetSearchQuery(request.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	sortBy := s.defaultSortBy
	if request.SortBy != "" {
		sortBy = request.SortBy
	}

	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		request.Query,
		pipeline,
		common.GetSortQuery(sortBy, request.Sort),
		getAuthorPipeline(),
	))

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	var result AggregationResult
	if cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func (s *store) Update(ctx context.Context, request EditRequest) (*Response, error) {
	now := types.NewCpsTime()
	model := transformRequestToModel(request)
	model.ID = request.ID
	model.Updated = now

	res, err := s.dbCollection.UpdateOne(
		ctx,
		bson.M{"_id": request.ID},
		bson.M{"$set": model},
	)
	if err != nil || res.MatchedCount == 0 {
		return nil, err
	}

	return s.GetById(ctx, model.ID, model.Author)
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	deleted, err := s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	return deleted > 0, nil
}

func getAuthorPipeline() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.RightsMongoCollection,
			"localField":   "author",
			"foreignField": "_id",
			"as":           "author",
		}},
		{"$unwind": bson.M{"path": "$author", "preserveNullAndEmptyArrays": true}},
	}
}

func transformRequestToModel(request EditRequest) savedpattern.SavedPattern {
	model := savedpattern.SavedPattern{
		Title:       request.Title,
		Type:        request.Type,
		IsCorporate: *request.IsCorporate,
		Author:      request.Author,
	}

	switch request.Type {
	case savedpattern.TypeAlarm:
		model.AlarmPattern = request.AlarmPattern
	case savedpattern.TypeEntity:
		model.EntityPattern = request.EntityPattern
	case savedpattern.TypePbehavior:
		model.PbehaviorPattern = request.PbehaviorPattern
	}

	return model
}
