package pbehavior_legacy

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"os"
	"time"
)

// PBehaviorLegacyCollectionName is where we store pbehaviors
const (
	PBehaviorLegacyCollectionName = "default_pbehavior"
)

// Adapter allows you to manipulate pbehaviors in database
type mongoAdapter struct {
	collection mongo.Collection
}

// NewAdapter gives the correct mongo pbehavior adapter.
func NewAdapter(collection mongo.Collection) Adapter {
	return mongoAdapter{
		collection: collection,
	}
}

func DefaultCollection(session *mgo.Session) mongo.Collection {
	collection := session.DB(canopsis.DbName).C(PBehaviorLegacyCollectionName)
	return mongo.FromMgo(collection)
}

func (a mongoAdapter) Get(filter bson.M) ([]types.PBehaviorLegacy, error) {
	pbehaviors := []types.PBehaviorLegacy{}
	err := a.collection.Get(filter, &pbehaviors)
	if err != nil {
		return pbehaviors, err
	}

	timezoneName := os.Getenv("CPS_PBH_TIMEZONE")
	if timezoneName == "" {
		timezoneName = "Europe/Paris"
	}
	defaultLocation, _ := time.LoadLocation(timezoneName)
	defaultTimezone := defaultLocation.String()
	for i, _ := range pbehaviors {
		if pbehaviors[i].Timezone == "" {
			pbehaviors[i].Timezone = defaultTimezone
		}
	}
	return pbehaviors, nil
}

func (a mongoAdapter) GetByEntityIds(eids []string, enabled bool) ([]types.PBehaviorLegacy, error) {
	pbehaviors := []types.PBehaviorLegacy{}
	err := a.collection.Get(bson.M{"enabled": enabled, "eids": bson.M{"$in": eids}}, &pbehaviors)
	return pbehaviors, err
}

func (a mongoAdapter) Insert(pbehavior types.PBehaviorLegacy) error {
	return a.collection.Insert(pbehavior)
}

func (a mongoAdapter) Update(pbehavior types.PBehaviorLegacy) error {
	return a.collection.Update(pbehavior.ID, pbehavior)
}

func (a mongoAdapter) RemoveId(id string) error {
	return a.collection.Remove(id)
}
