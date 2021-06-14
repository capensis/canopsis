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

func (a *mongoAdapter) FindJunitWidgetsTestSuiteIDs(ctx context.Context, widgetIDs []string) ([]string, error) {
	cursor, err := a.collection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"tabs.widgets._id": bson.M{"$in": widgetIDs},
		}},
		{"$unwind": "$tabs"},
		{"$unwind": "$tabs.widgets"},
		{"$match": bson.M{
			"tabs.widgets.type": WidgetTypeJunit,
			"tabs.widgets._id":  bson.M{"$in": widgetIDs},
		}},
		{"$project": bson.M{
			"test_suites": "$tabs.widgets.internal_parameters." + WidgetInternalParamJunitTestSuites,
		}},
	})

	if err != nil {
		return nil, err
	}

	testSuiteIDs := make([]string, 0)
	for cursor.Next(ctx) {
		res := &struct {
			TestSuites []string `bson:"test_suites"`
		}{}
		err := cursor.Decode(res)
		if err != nil {
			return nil, err
		}

		testSuiteIDs = append(testSuiteIDs, res.TestSuites...)
	}

	return testSuiteIDs, nil
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

func (a *mongoAdapter) RemoveTestSuitesFromJunitWidgets(
	ctx context.Context,
	testSuiteIDs []string,
) error {
	if len(testSuiteIDs) == 0 {
		return nil
	}

	_, err := a.collection.UpdateMany(ctx,
		bson.M{
			"tabs.widgets.type": WidgetTypeJunit,
		},
		bson.M{"$pullAll": bson.M{
			"tabs.$[tab].widgets.$[widget].internal_parameters." + WidgetInternalParamJunitTestSuites: testSuiteIDs,
		}},
		options.Update().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"tab.widgets.type": WidgetTypeJunit},
				bson.M{"widget.type": WidgetTypeJunit},
			},
		}),
	)
	return err
}
