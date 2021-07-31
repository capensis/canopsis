package flappingrule

//go:generate mockgen -destination=../../../mocks/lib/canopsis/flappingrule/adapter.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/flappingrule Adapter

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Adapter interface is used to implement an storage ruleAdapter.
type Adapter interface {
	// Get returns all flappingRules.
	Get(ctx context.Context) ([]Rule, error)
}

type mongoAdapter struct {
	collection mongo.DbCollection
}

// NewAdapter creates new rules ruleAdapter.
func NewAdapter(client mongo.DbClient) Adapter {
	return &mongoAdapter{
		collection: client.Collection(mongo.FlappingRuleMongoCollection),
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
