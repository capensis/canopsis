package notification

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const SettingsID = "notification"

type Store interface {
	Find(ctx context.Context, r pagination.Query, userID string, roleIDs []string) (AggregationResult, error)
	GetSettings(ctx context.Context) (Settings, error)
	UpdateSettings(ctx context.Context, r Settings) (Settings, error)
}

type store struct {
	collection         mongo.DbCollection
	settingsCollection mongo.DbCollection
	authorProvider     author.Provider
}

func NewStore(dbClient mongo.DbClient, authorProvider author.Provider) Store {
	return &store{
		collection:         dbClient.Collection(mongo.UserNotificationCollection),
		settingsCollection: dbClient.Collection(mongo.UserNotificationSettingsCollection),
		authorProvider:     authorProvider,
	}
}

func (s *store) Find(ctx context.Context, r pagination.Query, userID string, roleIDs []string) (AggregationResult, error) {
	res := AggregationResult{}
	beforeLimit := []bson.M{
		{"$match": bson.M{
			"$or": []bson.M{
				{"user": userID},
				{"roles": bson.M{"$in": roleIDs}},
			},
		}},
	}
	afterLimit := s.authorProvider.Pipeline()
	cursor, err := s.collection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r,
		beforeLimit,
		common.GetSortQuery("time", mongo.SortDesc),
		afterLimit,
	))
	if err != nil {
		return res, err
	}

	if cursor.Next(ctx) {
		err = cursor.Decode(&res)
		if err != nil {
			return res, err
		}
	}

	if err = cursor.Err(); err != nil {
		return res, err
	}

	if err = cursor.Close(ctx); err != nil {
		return res, err
	}

	return res, nil
}

func (s *store) GetSettings(ctx context.Context) (Settings, error) {
	res := Settings{}
	err := s.settingsCollection.FindOne(ctx, bson.M{"_id": SettingsID}).Decode(&res)

	return res, err
}

func (s *store) UpdateSettings(ctx context.Context, request Settings) (Settings, error) {
	res := Settings{}
	err := s.settingsCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": SettingsID},
		bson.M{"$set": request},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&res)

	return res, err
}
