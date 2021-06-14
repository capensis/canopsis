package heartbeat

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"github.com/globalsign/mgo"
)

const (
	HeartbeatCollectionName = "heartbeat"
)

type mongoAdapter struct {
	collection mongo.Collection
}

func (a *mongoAdapter) Get() ([]Heartbeat, error) {
	hl := make([]Heartbeat, 0)
	err := a.collection.Get(nil, &hl)
	return hl, err
}

func NewAdapter(collection mongo.Collection) Adapter {
	return &mongoAdapter{
		collection: collection,
	}
}

func DefaultCollection(session *mgo.Session) mongo.Collection {
	collection := session.DB(canopsis.DbName).C(HeartbeatCollectionName)
	return mongo.FromMgo(collection)
}
