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
	Find(ctx context.Context, userId string) ([]Response, error)
	Update(ctx context.Context, userId, widgetId string, request EditRequest) (*Response, bool, error)
}

type store struct {
	collection, viewCollection mongo.DbCollection
}

func NewStore(
	dbClient mongo.DbClient,
) Store {
	return &store{
		collection:     dbClient.Collection(mongo.UserPreferencesMongoCollection),
		viewCollection: dbClient.Collection(mongo.ViewMongoCollection),
	}
}

func (s *store) Find(ctx context.Context, userId string) ([]Response, error) {
	cursor, err := s.collection.Find(ctx, bson.M{"user": userId})
	if err != nil {
		return nil, err
	}

	res := make([]Response, 0)
	err = cursor.All(ctx, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *store) Update(ctx context.Context, userId, widgetId string, request EditRequest) (*Response, bool, error) {
	err := s.viewCollection.FindOne(ctx, bson.M{"tabs.widgets._id": widgetId}).Err()
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, false, nil
		}
		return nil, false, err
	}

	res, err := s.collection.UpdateOne(ctx, bson.M{
		"user":   userId,
		"widget": widgetId,
	}, bson.M{
		"$set": bson.M{
			"content": request.Content,
			"updated": types.CpsTime{Time: time.Now()},
		},
		"$setOnInsert": bson.M{
			"_id":    utils.NewID(),
			"user":   userId,
			"widget": widgetId,
		},
	}, options.Update().SetUpsert(true))

	if err != nil {
		return nil, false, err
	}

	isNew := res.UpsertedCount > 0

	return &Response{
		Widget:  widgetId,
		Content: request.Content,
	}, isNew, nil
}
