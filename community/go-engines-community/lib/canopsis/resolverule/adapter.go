package resolverule

//go:generate go tool go.uber.org/mock/mockgen -destination=../../../mocks/lib/canopsis/resolverule/rule.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/resolverule Adapter

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/priority"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Adapter interface is used to implement a storage adapter.
type Adapter interface {
	Get(ctx context.Context) ([]Rule, error)
}

type mongoAdapter struct {
	collection mongo.DbCollection
}

// NewAdapter creates new rule adapter.
func NewAdapter(client mongo.DbClient) Adapter {
	return &mongoAdapter{
		collection: client.Collection(mongo.ResolveRuleMongoCollection),
	}
}

func (a *mongoAdapter) Get(ctx context.Context) ([]Rule, error) {
	cursor, err := a.collection.Aggregate(
		ctx,
		append([]bson.M{{"$match": bson.M{"enabled": true}}}, priority.GetSortPipeline()...),
	)
	if err != nil {
		return nil, err
	}

	var rules []Rule
	err = cursor.All(ctx, &rules)

	if err != nil {
		return nil, err
	}

	return rules, nil
}
