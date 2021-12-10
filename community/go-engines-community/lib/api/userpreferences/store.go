package userpreferences

import (
	"context"
	"errors"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Store interface {
	Find(ctx context.Context, userId, widgetId string) (*Response, error)
	Update(ctx context.Context, userId string, request EditRequest) (*Response, bool, error)
}

type store struct {
	collection mongo.DbCollection
}

func NewStore(
	dbClient mongo.DbClient,
) Store {
	return &store{
		collection: dbClient.Collection(mongo.UserPreferencesMongoCollection),
	}
}

func (s *store) Find(ctx context.Context, userId, widgetId string) (*Response, error) {
	res := Response{
		Widget:  widgetId,
		Content: map[string]interface{}{},
	}
	err := s.collection.FindOne(ctx, bson.M{
		"user":   userId,
		"widget": widgetId,
	}).Decode(&res)
	if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
		return nil, err
	}

	return &res, nil
}

func (s *store) Update(ctx context.Context, userId string, request EditRequest) (*Response, bool, error) {
	res, err := s.collection.UpdateOne(ctx, bson.M{
		"user":   userId,
		"widget": request.Widget,
	}, bson.M{
		"$set": bson.M{
			"content": request.Content,
			"updated": types.CpsTime{Time: time.Now()},
		},
		"$setOnInsert": bson.M{
			"_id":    utils.NewID(),
			"user":   userId,
			"widget": request.Widget,
		},
	}, options.Update().SetUpsert(true))

	if err != nil {
		return nil, false, err
	}

	isNew := res.UpsertedCount > 0

	return &Response{
		Widget:  request.Widget,
		Content: request.Content,
	}, isNew, nil
}
