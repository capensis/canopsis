package eventfilter

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoAdapter struct {
	dbCollection mongo.DbCollection
}

func NewRuleAdapter(dbClient mongo.DbClient) RuleAdapter {
	return mongoAdapter{
		dbCollection: dbClient.Collection(mongo.EventFilterRulesMongoCollection),
	}
}

func (a mongoAdapter) GetAll(ctx context.Context) ([]Rule, error) {
	return a.find(ctx, bson.M{
		"$or": []bson.M{
			{"enabled": true},
			{"enabled": bson.M{"$exists": false}},
		},
	})
}

func (a mongoAdapter) GetByTypes(ctx context.Context, types []string) ([]Rule, error) {
	return a.find(ctx, bson.M{
		"$or": []bson.M{
			{"enabled": true},
			{"enabled": bson.M{"$exists": false}},
		},
		"type": bson.M{"$in": types},
	})
}

func (a mongoAdapter) find(ctx context.Context, filter bson.M) ([]Rule, error) {
	cursor, err := a.dbCollection.Find(ctx, filter, options.Find().SetSort(bson.M{"priority": 1}))
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
