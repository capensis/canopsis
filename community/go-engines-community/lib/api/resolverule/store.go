package resolverule

import (
	"context"
	"errors"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/resolverule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	Insert(ctx context.Context, r CreateRequest) (*Response, error)
	GetById(ctx context.Context, id string) (*Response, error)
	Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error)
	Update(ctx context.Context, r UpdateRequest) (*Response, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type store struct {
	dbClient     mongo.DbClient
	dbCollection mongo.DbCollection

	defaultSearchByFields []string
}

func NewStore(
	dbClient mongo.DbClient,
) Store {
	return &store{
		dbClient:     dbClient,
		dbCollection: dbClient.Collection(mongo.ResolveRuleMongoCollection),

		defaultSearchByFields: []string{"_id", "author", "name", "description"},
	}
}

func (s *store) Insert(ctx context.Context, request CreateRequest) (*Response, error) {
	id := request.ID
	if id == "" {
		id = utils.NewID()
	}
	now := types.NewCpsTime(time.Now().Unix())

	_, err := s.dbCollection.InsertOne(ctx, resolverule.Rule{
		ID:             id,
		Name:           request.Name,
		Description:    request.Description,
		Duration:       request.Duration,
		AlarmPatterns:  request.AlarmPatterns,
		EntityPatterns: request.EntityPatterns,
		Priority:       request.Priority,
		Author:         request.Author,
		Created:        now,
		Updated:        now,
	})
	if err != nil {
		return nil, err
	}

	err = s.updateFollowingPriorities(ctx, id, request.Priority)
	if err != nil {
		return nil, err
	}

	return s.GetById(ctx, id)
}

func (s *store) GetById(ctx context.Context, id string) (*Response, error) {
	pipeline := []bson.M{{"$match": bson.M{"_id": id}}}
	pipeline = append(pipeline, s.authorPipeline()...)

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

func (s *store) Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error) {
	pipeline := make([]bson.M, 0)
	filter := common.GetSearchQuery(query.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	sortBy := "created"
	if query.SortBy != "" {
		sortBy = query.SortBy
	}

	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		query.Query,
		pipeline,
		common.GetSortQuery(sortBy, query.Sort),
		s.authorPipeline(),
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

func (s *store) Update(ctx context.Context, request UpdateRequest) (*Response, error) {
	model := resolverule.Rule{}
	err := s.dbCollection.FindOne(ctx, bson.M{"_id": request.ID}).Decode(&model)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	now := types.NewCpsTime(time.Now().Unix())
	_, err = s.dbCollection.UpdateOne(
		ctx,
		bson.M{"_id": request.ID},
		bson.M{"$set": resolverule.Rule{
			ID:             request.ID,
			Name:           request.Name,
			Description:    request.Description,
			Duration:       request.Duration,
			AlarmPatterns:  request.AlarmPatterns,
			EntityPatterns: request.EntityPatterns,
			Priority:       request.Priority,
			Author:         request.Author,
			Created:        model.Created,
			Updated:        now,
		}},
	)
	if err != nil {
		return nil, err
	}

	if model.Priority != request.Priority {
		err := s.updateFollowingPriorities(ctx, request.ID, request.Priority)
		if err != nil {
			return nil, err
		}
	}

	return s.GetById(ctx, model.ID)
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	deleted, err := s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	return deleted > 0, nil
}

func (s *store) updateFollowingPriorities(ctx context.Context, id string, priority int) error {
	err := s.dbCollection.FindOne(ctx, bson.M{
		"_id":      bson.M{"$ne": id},
		"priority": priority,
	}).Err()
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil
		}
		return err
	}

	_, err = s.dbCollection.UpdateMany(
		ctx,
		bson.M{
			"_id":      bson.M{"$ne": id},
			"priority": bson.M{"$gte": priority},
		},
		bson.M{"$inc": bson.M{"priority": 1}},
	)

	return err
}

func (s *store) authorPipeline() []bson.M {
	return []bson.M{
		// Author
		{"$lookup": bson.M{
			"from":         mongo.RightsMongoCollection,
			"localField":   "author",
			"foreignField": "_id",
			"as":           "author",
		}},
		{"$unwind": bson.M{"path": "$author", "preserveNullAndEmptyArrays": true}},
	}
}
