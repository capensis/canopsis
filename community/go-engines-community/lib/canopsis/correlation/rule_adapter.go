package correlation

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

// RulesAdapter is a type that provides access to the MongoDB collection containing
// the meta-alarm rules
type mongoAdapter struct {
	dbClient     mongo.DbClient
	dbCollection mongo.DbCollection
}

func (a mongoAdapter) Get(ctx context.Context) ([]Rule, error) {
	cursor, err := a.dbCollection.Find(ctx, bson.M{
		"type": bson.M{
			"$ne": RuleTypeManualGroup,
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

func (a mongoAdapter) Save(ctx context.Context, rule Rule) error {
	_, err := a.dbCollection.InsertOne(ctx, rule)
	if err != nil {
		return err
	}

	return err
}

func (a mongoAdapter) GetManualRule(ctx context.Context, autoResolve bool) (Rule, error) {
	var rule Rule

	err := a.dbCollection.FindOne(ctx, bson.M{
		"type":         RuleTypeManualGroup,
		"auto_resolve": autoResolve,
	}).Decode(&rule)

	if errors.Is(err, mongodriver.ErrNoDocuments) {
		return rule, nil
	}

	return rule, err
}

func (a mongoAdapter) GetRule(ctx context.Context, id string) (Rule, error) {
	res := a.dbCollection.FindOne(ctx, bson.M{"_id": id})
	if err := res.Err(); err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return Rule{}, nil
		}

		return Rule{}, err
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
