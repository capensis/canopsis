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

func (s *store) Find(ctx context.Context, userId, widgetId string) (*Response, error) {
	ok, err := s.existWidget(ctx, widgetId)
	if err != nil || !ok {
		return nil, err
	}

	res := Response{
		Widget:  widgetId,
		Content: map[string]interface{}{},
	}
	err = s.collection.FindOne(ctx, bson.M{
		"user":   userId,
		"widget": widgetId,
	}).Decode(&res)
	if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
		return nil, err
	}

	return &res, nil
}

func (s *store) Update(ctx context.Context, userId string, request EditRequest) (*Response, bool, error) {
	ok, err := s.existWidget(ctx, request.Widget)
	if err != nil || !ok {
		return nil, false, err
	}

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

func (s *store) existWidget(ctx context.Context, widgetId string) (bool, error) {
	err := s.viewCollection.FindOne(ctx, bson.M{"tabs.widgets._id": widgetId}).Err()
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
