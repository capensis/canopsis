package idlerule

import (
	"context"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

// RuleAdapter interface is used to implement an storage adapter.
type RuleAdapter interface {
	// Get returns all rules.
	Get() ([]Rule, error)
}

type mongoAdapter struct {
	collection mongo.DbCollection
}

// NewRuleAdapter creates new rule adapter.
func NewRuleAdapter(client mongo.DbClient) RuleAdapter {
	return &mongoAdapter{
		collection: client.Collection(mongo.IdleRuleMongoCollection),
	}
}

func (a *mongoAdapter) Get() ([]Rule, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cursor, err := a.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var rules []Rule
	err = cursor.All(ctx, &rules)

	return rules, err
}
