package eventfilter

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// RulesCollectionName is the name of the MongoDB collection containing the
// event filter's rules.
const (
	RulesCollectionName = "eventfilter"
)

// Adapter is a type that provides access to the MongoDB collection containing
// the event filter's rules
type Adapter struct {
	collection mongo.Collection
}

// NewAdapter returns an adapter to a rules collection.
func NewAdapter(collection mongo.Collection) Adapter {
	return Adapter{
		collection: collection,
	}
}

// DefaultCollection returns the default collection containing the event
// filter's rules.
func DefaultCollection(session *mgo.Session) mongo.Collection {
	collection := session.DB(canopsis.DbName).C(RulesCollectionName)
	return mongo.FromMgo(collection)
}

// List returns a list of all the rules that are enabled and valid, sorted by
// ascending priority.
func (a Adapter) List() ([]Rule, error) {
	ruleUnpackers := []RuleUnpacker{}

	// Get the rules that are enabled, or where enabled is not set.
	filter := bson.M{
		"$or": []bson.M{
			{"enabled": true},
			{"enabled": bson.M{"$exists": false}},
		},
	}
	err := a.collection.GetSorted(filter, []string{"priority"}, &ruleUnpackers)

	// Filter-out the invalid rules, and unpack the rules
	rules := []Rule{}
	for _, rule := range ruleUnpackers {
		if rule.Valid {
			rules = append(rules, rule.Rule)
		}
	}

	return rules, err
}
