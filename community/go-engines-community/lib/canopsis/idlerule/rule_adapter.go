package idlerule

//go:generate mockgen -destination=../../../mocks/lib/canopsis/idlerule/adapter.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlerule RuleAdapter

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RuleAdapter interface is used to implement an storage adapter.
type RuleAdapter interface {
	// GetEnabled returns all enabled rules.
	GetEnabled(ctx context.Context) ([]Rule, error)
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

func (a *mongoAdapter) GetEnabled(ctx context.Context) ([]Rule, error) {
	cursor, err := a.collection.Find(ctx, bson.M{"enabled": true},
		options.Find().SetSort(bson.M{"priority": 1}))
	if err != nil {
		return nil, err
	}

	var rules []Rule
	err = cursor.All(ctx, &rules)

	return rules, err
}
