package sharetoken

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/token"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type Store interface {
	Insert(ctx context.Context, userId string, r EditRequest) (*Response, error)
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type store struct {
	client     mongo.DbClient
	collection mongo.DbCollection

	tokenGenerator token.Generator

	defaultSearchByFields []string
	defaultSortBy         string
}

func NewStore(
	dbClient mongo.DbClient,
	tokenGenerator token.Generator,
) Store {
	return &store{
		client:         dbClient,
		collection:     dbClient.Collection(mongo.ShareTokenMongoCollection),
		tokenGenerator: tokenGenerator,

		defaultSearchByFields: []string{"value", "user.name", "role.name", "description"},
		defaultSortBy:         "created",
	}
}

func (s *store) Insert(ctx context.Context, userId string, r EditRequest) (*Response, error) {
	now := types.NewCpsTime()
	var expired types.CpsTime
	if r.Duration != nil && r.Duration.Value > 0 {
		expired = r.Duration.AddTo(now)
	}
	accessToken, err := s.tokenGenerator.Generate(userId, expired.Time)
	if err != nil {
		return nil, err
	}

	model := Model{
		ID:          utils.NewID(),
		Value:       accessToken,
		User:        userId,
		Description: r.Description,
		Created:     now,
		Accessed:    now,
	}
	if !expired.IsZero() {
		model.Expired = &expired
	}

	var response *Response
	err = s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		_, err = s.collection.InsertOne(ctx, model)
		if err != nil {
			return err
		}

		pipeline := []bson.M{
			{"$match": bson.M{"_id": model.ID}},
		}
		pipeline = append(pipeline, getUserPipeline()...)
		cursor, err := s.collection.Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}
		defer cursor.Close(ctx)
		if cursor.Next(ctx) {
			response = &Response{}
			err = cursor.Decode(response)
			if err != nil {
				return err
			}
		}

		return err
	})

	return response, err
}

func (s *store) Find(ctx context.Context, request ListRequest) (*AggregationResult, error) {
	pipeline := getUserPipeline()
	filter := common.GetSearchQuery(request.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	sortBy := s.defaultSortBy
	if request.SortBy != "" {
		sortBy = request.SortBy
	}

	displayedChars := 5
	project := []bson.M{
		{"$addFields": bson.M{
			"value": bson.M{"$concat": bson.A{
				"**********",
				bson.M{"$substr": bson.A{
					"$value",
					bson.M{"$subtract": bson.A{bson.M{"$strLenCP": "$value"}, displayedChars}},
					displayedChars,
				}},
			}},
		}},
	}

	cursor, err := s.collection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		request.Query,
		pipeline,
		common.GetSortQuery(sortBy, request.Sort),
		project,
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

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	deleted, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	return deleted > 0, nil
}

func getUserPipeline() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.RightsMongoCollection,
			"localField":   "user",
			"foreignField": "_id",
			"as":           "user",
		}},
		{"$unwind": bson.M{"path": "$user", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         mongo.RightsMongoCollection,
			"localField":   "user.role",
			"foreignField": "_id",
			"as":           "role",
		}},
		{"$unwind": bson.M{"path": "$role", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"user.name": "$user.crecord_name",
			"role.name": "$role.crecord_name",
		}},
	}
}
