package view

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

func NewMongoAdapter(client mongo.DbClient) Adapter {
	return &mongoAdapter{
		client:           client,
		collection:       client.Collection(mongo.ViewMongoCollection),
		widgetCollection: client.Collection(mongo.WidgetMongoCollection),
	}
}

type mongoAdapter struct {
	client           mongo.DbClient
	collection       mongo.DbCollection
	widgetCollection mongo.DbCollection
}

func (a *mongoAdapter) FindJunitWidgets(ctx context.Context) ([]Widget, error) {
	cursor, err := a.widgetCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"type": WidgetTypeJunit}},
		{"$lookup": bson.M{
			"from":         mongo.ViewTabMongoCollection,
			"localField":   "tab",
			"foreignField": "_id",
			"as":           "tab_obj",
		}},
		{"$unwind": "$tab_obj"},
		{"$lookup": bson.M{
			"from":         mongo.ViewMongoCollection,
			"localField":   "tab_obj.view",
			"foreignField": "_id",
			"as":           "view",
		}},
		{"$unwind": "$view"},
		{"$match": bson.M{"view.enabled": true}},
	})
	if err != nil {
		return nil, err
	}

	widgets := make([]Widget, 0)
	err = cursor.All(ctx, &widgets)
	if err != nil {
		return nil, err
	}

	return widgets, nil
}

func (a *mongoAdapter) AddTestSuitesToJunitWidgets(
	ctx context.Context,
	widgetIDs, testSuiteIDs []string,
) error {
	if len(widgetIDs) == 0 || len(testSuiteIDs) == 0 {
		return nil
	}

	_, err := a.widgetCollection.UpdateMany(ctx,
		bson.M{
			"_id": bson.M{"$in": widgetIDs},
		},
		bson.M{"$push": bson.M{
			"internal_parameters." + WidgetInternalParamJunitTestSuites: bson.M{
				"$each": testSuiteIDs,
			},
		}},
	)
	return err
}
