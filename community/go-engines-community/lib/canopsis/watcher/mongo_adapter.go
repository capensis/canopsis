package watcher

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/entity"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"git.canopsis.net/canopsis/go-engines/lib/mongo/bulk"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/rs/zerolog"
)

// Adapter allows you to manipulate watchers in database
type mongoAdapter struct {
	collection           mongo.Collection
	mbulk                bulk.Bulk
	entityCollectionName string
	alarmCollectionName  string
	logger               zerolog.Logger
}

// NewAdapter gives the correct mongo watcher adapter.
func NewAdapter(
	collection mongo.Collection,
	safeBulk bulk.Bulk,
	entityCollectionName string,
	alarmCollectionName string,
	logger zerolog.Logger,
) Adapter {
	return mongoAdapter{
		collection:           collection,
		mbulk:                safeBulk,
		entityCollectionName: entityCollectionName,
		alarmCollectionName:  alarmCollectionName,
		logger:               logger,
	}
}

// DefaultCollection returns the collection containing the watchers.
func DefaultCollection(session *mgo.Session) mongo.Collection {
	collection := session.DB(canopsis.DbName).C(entity.EntityCollectionName)
	return mongo.FromMgo(collection)
}

func (a mongoAdapter) GetAll(watchers *[]Watcher) error {
	return a.collection.Get(bson.M{"type": types.EntityTypeWatcher}, watchers)
}

func (a mongoAdapter) GetAllValidWatchers(watchers *[]Watcher) error {
	var tmpWatchers []Watcher
	err := a.GetAll(&tmpWatchers)
	if err != nil {
		return err
	}
	for _, watcher := range tmpWatchers {
		if watcher.Entities.IsSet() && watcher.Entities.IsValid() && watcher.Entity.Enabled {
			*watchers = append(*watchers, watcher)
		} else {
			a.logger.Debug().Str("watcher", watcher.ID).Str("watcher_name", watcher.Name).Msg("Skipping watcher")
		}
	}
	return nil
}

func (a mongoAdapter) GetByID(id string, watcher *Watcher) error {
	var tmpWatcher Watcher
	err := a.collection.GetByID(id, &tmpWatcher)
	if err != nil {
		return err
	}

	if tmpWatcher.Type == types.EntityTypeWatcher {
		*watcher = tmpWatcher
	} else {
		a.logger.Warn().Str("watcher", tmpWatcher.ID).Str("type", tmpWatcher.Type).Msg("Is not a watcher")
	}
	return nil
}

func (a mongoAdapter) GetEntities(watcher Watcher, entities *[]types.Entity) error {
	mongoFilter := watcher.Entities.AsMongoQuery()
	return a.collection.Get(mongoFilter, entities)
}

func (a mongoAdapter) getAnnotatedEntitiesIter(query bson.M) *mgo.Iter {
	pipeline := []bson.M{
		{
			"$match": query,
		},

		// Move the entity away from the root of the returned document, to make
		// the unmarshaling easier
		{
			"$project": bson.M{
				"entity": "$$ROOT",
			},
		},

		// Get the entity's unresolved alarm
		// The unwind step prevents from working with a list of alarms, which
		// cannot be filtered efficiently and may reach MongoDB's document size
		// limit.
		// The preserveNullAndEmptyArrays option is necessary to keep the
		// entities with no alarms.
		{
			"$lookup": bson.M{
				"from":         a.alarmCollectionName,
				"localField":   "entity._id",
				"foreignField": "d",
				"as":           "alarm",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$alarm",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$match": bson.M{
				"alarm.v.resolved": nil,
			},
		},
	}

	return a.collection.Pipe(pipeline).Iter()
}

func (a mongoAdapter) getAnnotatedEntities(query bson.M) ([]AnnotatedEntity, error) {
	pipeline := []bson.M{
		bson.M{
			"$match": query,
		},

		// Move the entity away from the root of the returned document, to make
		// the unmarshaling easier
		bson.M{
			"$project": bson.M{
				"entity": "$$ROOT",
			},
		},

		// Get the entity's unresolved alarm
		// The unwind step prevents from working with a list of alarms, which
		// cannot be filtered efficiently and may reach MongoDB's document size
		// limit.
		// The preserveNullAndEmptyArrays option is necessary to keep the
		// entities with no alarms.
		bson.M{
			"$lookup": bson.M{
				"from":         a.alarmCollectionName,
				"localField":   "entity._id",
				"foreignField": "d",
				"as":           "alarm",
			},
		},
		bson.M{
			"$unwind": bson.M{
				"path":                       "$alarm",
				"preserveNullAndEmptyArrays": true,
			},
		},
		bson.M{
			"$match": bson.M{
				"alarm.v.resolved": nil,
			},
		},
	}

	entities := []AnnotatedEntity{}
	err := a.collection.Aggregate(pipeline, &entities)
	return entities, err
}

func (a mongoAdapter) GetAllAnnotatedEntities() ([]AnnotatedEntity, error) {
	return a.getAnnotatedEntities(bson.M{
		"enabled": true,
	})
}

func (a mongoAdapter) GetAnnotatedEntitiesIter() *mgo.Iter {
	return a.getAnnotatedEntitiesIter(bson.M{
		"enabled": true,
	})
}

func (a mongoAdapter) GetAnnotatedDependencies(watcherID string) ([]AnnotatedEntity, error) {
	return a.getAnnotatedEntities(bson.M{
		"enabled": true,
		"impact":  watcherID,
	})
}

func (a mongoAdapter) Update(watcherID string, update interface{}) error {
	return a.collection.Update(watcherID, update)
}

func (a mongoAdapter) GetAnnotatedDependenciesIter(watcherID string) *mgo.Iter {
	return a.getAnnotatedEntitiesIter(bson.M{
		"enabled": true,
		"impact":  watcherID,
	})
}
