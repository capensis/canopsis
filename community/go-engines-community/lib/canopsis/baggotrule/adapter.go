package baggotrule

//go:generate mockgen -destination=../../../mocks/lib/canopsis/baggotrule/adapter.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/baggotrule Adapter

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

// Adapter interface is used to implement an storage ruleAdapter.
type Adapter interface {
	// Get returns all baggotRules.
	Get(ctx context.Context) ([]Rule, error)
}

type mongoAdapter struct {
	collection mongo.DbCollection
}

// NewAdapter creates new rules ruleAdapter.
func NewAdapter(client mongo.DbClient) Adapter {
	return &mongoAdapter{
		collection: client.Collection(mongo.BaggotRuleMongoCollection),
	}
}

func (a *mongoAdapter) Get(ctx context.Context) ([]Rule, error) {
	cursor, err := a.collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{
		{Key: "priority", Value: 1},
		{Key: "_id", Value: 1},
	}))
	if err != nil {
		return nil, err
	}

	var rules []Rule
	err = cursor.All(ctx, &rules)

	if err != nil {
		return nil, err
	}

	return rules, err
}
