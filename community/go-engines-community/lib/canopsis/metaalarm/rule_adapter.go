package metaalarm

import (
	"errors"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis"
	"git.canopsis.net/canopsis/go-engines/lib/errt"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const (
	RulesCollectionName = "meta_alarm_rules"
)

// RulesAdapter is a type that provides access to the MongoDB collection containing
// the meta-alarm rules
type mongoAdapter struct {
	collection mongo.Collection
}

func (a mongoAdapter) Get() ([]Rule, error) {
	var rules []Rule
	err := a.collection.Get(map[string]interface{}{
		"type": bson.M{
			"$ne": RuleManualGroup,
		},
	}, &rules)
	return rules, err
}

func (a mongoAdapter) Save(rule Rule) error {
	err := a.collection.Insert(rule)
	return err
}

func (a mongoAdapter) GetManualRule() (Rule, error) {
	var rules []Rule
	err := a.collection.Get(map[string]interface{}{
		"type": bson.M{
			"$eq": RuleManualGroup,
		},
	}, &rules)
	if err != nil {
		return Rule{}, err
	}

	if len(rules) == 0 {
		return Rule{}, errt.NewNotFound(errors.New("not found existing manualrule"))
	}

	return rules[0], nil
}

func (a mongoAdapter) GetRule(id string) (Rule, error) {
	var rule Rule
	err := a.collection.GetOne(bson.M{"_id": id}, nil, &rule)
	return rule, err
}

// NewRuleAdapter returns an rulesAdapter to a rules collection.
func NewRuleAdapter(collection mongo.Collection) RulesAdapter {
	return mongoAdapter{
		collection: collection,
	}
}

// DefaultCollection returns the default collection containing the meta-alarm rules.
func DefaultRulesCollection(session *mgo.Session) mongo.Collection {
	collection := session.DB(canopsis.DbName).C(RulesCollectionName)
	return mongo.FromMgo(collection)
}
