package correlation

import (
	"context"
	mongodriver "go.mongodb.org/mongo-driver/mongo"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

// RulesAdapter is a type that provides access to the MongoDB collection containing
// the meta-alarm rules
type mongoAdapter struct {
	dbClient     mongo.DbClient
	dbCollection mongo.DbCollection
}

func (a mongoAdapter) Get() ([]Rule, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cursor, err := a.dbCollection.Find(ctx, bson.M{
		"type": bson.M{
			"$ne": RuleManualGroup,
		},
	})
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

func (a mongoAdapter) Save(rule Rule) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := a.dbCollection.InsertOne(ctx, rule)
	if err != nil {
		return err
	}

	return err
}

func (a mongoAdapter) GetManualRule() (Rule, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var rule Rule

	err := a.dbCollection.FindOne(ctx, bson.M{
		"type": bson.M{
			"$eq": RuleManualGroup,
		},
	}).Decode(&rule)

	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return rule, nil
		}
	}

	return rule, err
}

func (a mongoAdapter) GetRule(id string) (Rule, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	res := a.dbCollection.FindOne(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		return Rule{}, res.Err()
	}

	var rule Rule
	err := res.Decode(&rule)
	if err != nil {
		return Rule{}, err
	}

	return rule, nil
}

// NewRuleAdapter returns an rulesAdapter to a rules collection.
func NewRuleAdapter(dbClient mongo.DbClient) RulesAdapter {
	return mongoAdapter{
		dbClient:     dbClient,
		dbCollection: dbClient.Collection(mongo.MetaAlarmRulesMongoCollection),
	}
}
