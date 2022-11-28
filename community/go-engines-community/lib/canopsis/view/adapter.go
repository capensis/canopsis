package view

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoAdapter(client mongo.DbClient) Adapter {
	return &mongoAdapter{
		client:     client,
		collection: client.Collection(mongo.ViewMongoCollection),
	}
}

type mongoAdapter struct {
	client     mongo.DbClient
	collection mongo.DbCollection
}

func (a *mongoAdapter) FindJunitWidgets(ctx context.Context) ([]Widget, error) {
	cursor, err := a.collection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"enabled":           true,
			"tabs.widgets.type": WidgetTypeJunit,
		}},
		{"$unwind": "$tabs"},
		{"$unwind": "$tabs.widgets"},
		{"$match": bson.M{"tabs.widgets.type": WidgetTypeJunit}},
		{"$replaceRoot": bson.M{"newRoot": "$tabs.widgets"}},
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

	_, err := a.collection.UpdateMany(ctx,
		bson.M{
			"tabs.widgets._id": bson.M{"$in": widgetIDs},
		},
		bson.M{"$push": bson.M{
			"tabs.$[tab].widgets.$[widget].internal_parameters." + WidgetInternalParamJunitTestSuites: bson.M{
				"$each": testSuiteIDs,
			},
		}},
		options.Update().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"tab.widgets._id": bson.M{"$in": widgetIDs}},
				bson.M{"widget._id": bson.M{"$in": widgetIDs}},
			},
		}),
	)
	return err
}
