package neweventfilter

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RulesAdapter is a type that provides access to the MongoDB collection containing
// the meta-alarm rules
type mongoAdapter struct {
	dbCollection mongo.DbCollection
}

// NewRuleAdapter returns an rulesAdapter to a rules collection.
func NewRuleAdapter(dbClient mongo.DbClient) RulesAdapter {
	return mongoAdapter{
		dbCollection: dbClient.Collection(mongo.EventFilterRulesMongoCollection),
	}
}

func (a mongoAdapter) GetAll(ctx context.Context) ([]Rule, error) {
	//TODO: after eventfilter revamp, it should retrieve all kind of rules
	cursor, err := a.dbCollection.Find(ctx, bson.M{
		"$or": []bson.M{
			{"enabled": true},
			{"enabled": bson.M{"$exists": false}},
		},
		"type": RuleTypeChangeEntity,
	}, options.Find().SetSort(bson.M{"priority": 1}))
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
