package userpreferences

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/widgetfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	Find(ctx context.Context, userId, widgetId string) (*Response, error)
	Update(ctx context.Context, userId string, request EditRequest) (*Response, bool, error)
}

type store struct {
	client           mongo.DbClient
	collection       mongo.DbCollection
	filterCollection mongo.DbCollection
}

func NewStore(
	dbClient mongo.DbClient,
) Store {
	return &store{
		client:           dbClient,
		collection:       dbClient.Collection(mongo.UserPreferencesMongoCollection),
		filterCollection: dbClient.Collection(mongo.WidgetFiltersMongoCollection),
	}
}

func (s *store) Find(ctx context.Context, userId, widgetId string) (*Response, error) {
	res := Response{
		Widget:  widgetId,
		Content: map[string]interface{}{},
		Filters: make([]widgetfilter.Response, 0),
	}
	cursor, err := s.collection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"user":   userId,
			"widget": widgetId,
		}},
		{"$lookup": bson.M{
			"from":         mongo.WidgetFiltersMongoCollection,
			"localField":   "widget",
			"foreignField": "widget",
			"as":           "filters",
		}},
		{"$unwind": bson.M{"path": "$filters", "preserveNullAndEmptyArrays": true}},
		{"$sort": bson.M{"filters.position": 1}},
		{"$group": bson.M{
			"_id":     nil,
			"user":    bson.M{"$first": "$user"},
			"widget":  bson.M{"$first": "$widget"},
			"content": bson.M{"$first": "$content"},
			"filters": bson.M{"$push": "$filters"},
		}},
		{"$addFields": bson.M{
			"filters": bson.M{"$filter": bson.M{
				"input": "$filters",
				"cond": bson.M{"$and": []bson.M{
					{"$eq": bson.A{"$$this.author", "$user"}},
					{"$eq": bson.A{"$$this.is_private", true}},
				}},
			}},
		}},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		err := cursor.Decode(&res)
		if err != nil {
			return nil, err
		}
	} else {
		filterCursor, err := s.filterCollection.Find(ctx, bson.M{
			"author":     userId,
			"widget":     widgetId,
			"is_private": true,
		}, options.Find().SetSort(bson.M{"position": 1}))
		if err != nil {
			return nil, err
		}
		err = filterCursor.All(ctx, &res.Filters)
		if err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (s *store) Update(ctx context.Context, userId string, request EditRequest) (*Response, bool, error) {
	var response *Response
	isNew := false
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		isNew = false

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
			return err
		}

		isNew = res.UpsertedCount > 0
		response, err = s.Find(ctx, userId, request.Widget)
		return err
	})

	return response, isNew, err
}
