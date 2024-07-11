package userpreferences

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/widgetfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	Find(ctx context.Context, userID, widgetId string) (*Response, error)
	Update(ctx context.Context, userID string, request EditRequest) (*Response, error)
}

type store struct {
	client           mongo.DbClient
	collection       mongo.DbCollection
	filterCollection mongo.DbCollection
	authorProvider   author.Provider
}

func NewStore(
	dbClient mongo.DbClient,
	authorProvider author.Provider,
) Store {
	return &store{
		client:           dbClient,
		collection:       dbClient.Collection(mongo.UserPreferencesMongoCollection),
		filterCollection: dbClient.Collection(mongo.WidgetFiltersMongoCollection),
		authorProvider:   authorProvider,
	}
}

func (s *store) Find(ctx context.Context, userID, widgetId string) (*Response, error) {
	res := Response{
		Widget:  widgetId,
		Content: map[string]interface{}{},
		Filters: make([]widgetfilter.Response, 0),
	}
	pipeline := []bson.M{
		{"$match": bson.M{
			"user":   userID,
			"widget": widgetId,
		}},
		{"$lookup": bson.M{
			"from": mongo.WidgetFiltersMongoCollection,
			"let": bson.M{
				"widget": "$widget",
				"user":   "$user",
			},
			"pipeline": []bson.M{
				{"$match": bson.M{
					"$expr": bson.M{"$and": []bson.M{
						{"$eq": bson.A{"$widget", "$$widget"}},
						{"$eq": bson.A{"$author", "$$user"}},
					}},
					"is_user_preference": true,
				}},
			},
			"as": "filters",
		}},
		{"$unwind": bson.M{"path": "$filters", "preserveNullAndEmptyArrays": true}},
	}
	pipeline = append(pipeline, s.authorProvider.PipelineForField("filters.author")...)
	pipeline = append(pipeline,
		bson.M{"$sort": bson.M{"filters.position": 1}},
		bson.M{"$group": bson.M{
			"_id":     nil,
			"user":    bson.M{"$first": "$user"},
			"widget":  bson.M{"$first": "$widget"},
			"content": bson.M{"$first": "$content"},
			"filters": bson.M{"$push": bson.M{"$cond": bson.M{
				"if":   "$filters._id",
				"then": "$filters",
				"else": "$$REMOVE",
			}}},
		}},
	)
	cursor, err := s.collection.Aggregate(ctx, pipeline)
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
		pipeline := []bson.M{
			{"$match": bson.M{
				"author":             userID,
				"widget":             widgetId,
				"is_user_preference": true,
			}},
			{"$sort": bson.M{"position": 1}},
		}
		pipeline = append(pipeline, s.authorProvider.Pipeline()...)
		filterCursor, err := s.filterCollection.Aggregate(ctx, pipeline)
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

func (s *store) Update(ctx context.Context, userID string, request EditRequest) (*Response, error) {
	var response *Response

	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		_, err := s.collection.UpdateOne(ctx, bson.M{
			"user":   userID,
			"widget": request.Widget,
		}, bson.M{
			"$set": bson.M{
				"content": request.Content,
				"updated": datetime.NewCpsTime(),
			},
			"$setOnInsert": bson.M{
				"_id":    utils.NewID(),
				"user":   userID,
				"widget": request.Widget,
			},
		}, options.Update().SetUpsert(true))

		if err != nil {
			return err
		}

		response, err = s.Find(ctx, userID, request.Widget)
		return err
	})

	return response, err
}
