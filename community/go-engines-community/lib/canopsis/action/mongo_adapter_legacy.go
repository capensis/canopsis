package action

import (
	"log"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type scenarioWrapper struct {
	wbh Scenario
}

func (u *scenarioWrapper) SetBSON(raw bson.Raw) error {
	m := make(map[string]interface{})
	err := raw.Unmarshal(&m)
	if err != nil {
		log.Printf("unable to parse action: %v", err)
		return nil
	}

	if err := raw.Unmarshal(&u.wbh); err != nil {
		log.Printf("unable to parse parse: %v", err)
		return nil
	}
	if _, ok := m["enabled"]; !ok {
		u.wbh.Enabled = true
	}

	return nil
}

// Adapter allows you to manipulate alarms in database
type mongoAdapterLegacy struct {
	collection mongo.Collection
}

// NewAdapter gives the correct mongo action adapter.
func NewAdapterLegacy(collection mongo.Collection) Adapter {
	return mongoAdapterLegacy{
		collection: collection,
	}
}

func DefaultCollection(session *mgo.Session) mongo.Collection {
	collection := session.DB(canopsis.DbName).C(mongo.ScenarioMongoCollection)
	return mongo.FromMgo(collection)
}

func (a mongoAdapterLegacy) Get() ([]Scenario, error) {
	var al []scenarioWrapper
	err := a.collection.GetSorted(nil, []string{PriorityField, IdField}, &al)
	if err != nil {
		return nil, err
	}
	scenarios := make([]Scenario, len(al))
	for i, a := range al {
		scenarios[i] = a.wbh
	}
	return scenarios, err
}

func (a mongoAdapterLegacy) GetEnabled() ([]Scenario, error) {
	var al []scenarioWrapper
	err := a.collection.GetSorted(
		bson.M{"$or": []interface{}{bson.M{"enabled": true}, bson.M{"enabled": bson.M{"$exists": false}}}},
		[]string{PriorityField, IdField},
		&al)
	if err != nil {
		return nil, err
	}
	scenarios := make([]Scenario, len(al))
	for i, a := range al {
		scenarios[i] = a.wbh
	}
	return scenarios, err
}

func (a mongoAdapterLegacy) GetEnabledById(id string) (Scenario, error) {
	var scenario Scenario
	err := a.collection.GetOne(bson.M{"$and": []bson.M{
		{"_id": id},
		{"$or": []bson.M{{"enabled": true}, {"enabled": bson.M{"$exists": false}}}},
	}}, nil, &scenario)

	return scenario, err
}

func (a mongoAdapterLegacy) GetEnabledByIDs(ids []string) ([]Scenario, error) {
	var al []scenarioWrapper
	err := a.collection.GetSorted(
		bson.M{"$and": []bson.M{
			{"_id": bson.M{"$in": ids}},
			{"$or": []bson.M{{"enabled": true}, {"enabled": bson.M{"$exists": false}}}},
		}},
		[]string{PriorityField, IdField},
		&al)
	if err != nil {
		return nil, err
	}
	scenarios := make([]Scenario, len(al))
	for i, a := range al {
		scenarios[i] = a.wbh
	}
	return scenarios, err
}
