package notification

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/mongoquery"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/usernotification"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const SettingsID = "notification"

type Store interface {
	Find(ctx context.Context, r pagination.Query, userID string, roleIDs []string) (AggregationResult, error)
	GetSettings(ctx context.Context) (SettingsResponse, error)
	UpdateSettings(ctx context.Context, r UpdateSettingsRequest) (SettingsResponse, error)
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
		{"$lookup": bson.M{
			"from":         mongo.EventFilterRuleCollection,
			"localField":   "rule._id",
			"foreignField": "_id",
			"let":          bson.M{"updated": "$rule.updated"},
			"as":           "eventfilter",
			"pipeline": []bson.M{
				{"$match": bson.M{"$expr": bson.M{"$eq": bson.A{"$updated", "$$updated"}}}},
				{"$limit": 1},
			},
		}},
		{"$unwind": bson.M{"path": "$eventfilter", "preserveNullAndEmptyArrays": true}},
		{"$match": bson.M{
			"$or": []bson.M{
				{"type": bson.M{"$ne": usernotification.TypeEventFilterFailure}},
				{"eventfilter": bson.M{"$ne": nil}},
			},
		}},
	}
	afterLimit := s.authorProvider.Pipeline()
	cursor, err := s.collection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r,
		beforeLimit,
		mongoquery.GetSortQuery("time", mongo.SortDesc),
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

func (s *store) GetSettings(ctx context.Context) (SettingsResponse, error) {
	res := SettingsResponse{}
	err := s.settingsCollection.FindOne(ctx, bson.M{"_id": SettingsID}).Decode(&res)

	return res, err
}

func (s *store) UpdateSettings(ctx context.Context, request UpdateSettingsRequest) (SettingsResponse, error) {
	res := SettingsResponse{}
	err := s.settingsCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": SettingsID},
		bson.M{"$set": Settings{
			Instruction: request.Instruction,
			Author:      request.Author,
			Updated:     datetime.NewCpsTime(),
		}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&res)

	return res, err
}
