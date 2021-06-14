package eventfilter

import (
	"context"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoAdapter struct {
	dbClient     mongo.DbClient
	dbCollection mongo.DbCollection
}

// NewAdapter returns an adapter to a rules collection.
func NewAdapter(dbClient mongo.DbClient) Adapter {
	return &mongoAdapter{
		dbClient:     dbClient,
		dbCollection: dbClient.Collection(mongo.EventFilterRulesMongoCollection),
	}
}

// List returns a list of all the rules that are enabled and valid, sorted by
// ascending priority.
func (a *mongoAdapter) List() ([]Rule, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var ruleUnpackers []RuleUnpacker

	// Get the rules that are enabled, or where enabled is not set.
	filter := bson.M{
		"$or": []bson.M{
			{"enabled": true},
			{"enabled": bson.M{"$exists": false}},
		},
	}

	cursor, err := a.dbCollection.Find(ctx, filter, options.Find().SetSort(bson.M{"priority": 1}))
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &ruleUnpackers)
	if err != nil {
		return nil, err
	}

	err = cursor.Close(ctx)
	if err != nil {
		return nil, err
	}

	// Filter-out the invalid rules, and unpack the rules
	var rules []Rule
	for _, rule := range ruleUnpackers {
		if rule.Valid {
			rules = append(rules, rule.Rule)
		}
	}

	return rules, err
}
