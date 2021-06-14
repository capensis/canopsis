package entity

import (
	"fmt"

	"git.canopsis.net/canopsis/go-engines/lib/errt"

	"git.canopsis.net/canopsis/go-engines/lib/mongo"

	cps "git.canopsis.net/canopsis/go-engines/lib/canopsis"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/mongo/bulk"
	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Entities constants
const (
	EntityCollectionName = mongo.EntityMongoCollection
)

// mongoAdapter ...
type mongoAdapter struct {
	Collection mongo.Collection
	mbulk      bulk.Bulk
}

// DefaultCollection returns the default mongo collection for entities
func DefaultCollection(session *mgo.Session) mongo.Collection {
	return mongo.FromMgo(session.DB(cps.DbName).C(EntityCollectionName))
}

// NewAdapter gives the correct entity adapter. Give nil to the redis client
// and it will create a new redis.Client with the dedicated redis database for entities.
func NewAdapter(collection mongo.Collection) Adapter {
	mbulk := bulk.NewSafe(collection.NewBulk(bulk.BulkSizeMax))

	return &mongoAdapter{
		Collection: collection,
		mbulk:      mbulk,
	}
}

// Insert add a new entity.
func (e *mongoAdapter) Insert(entity types.Entity) error {
	return e.Collection.Insert(entity)
}

// Update updates an existing entity or creates a new one in db.
func (e *mongoAdapter) Update(entity types.Entity) error {
	return e.Collection.Update(entity.ID, entity)
}

// Remove an existing entity in db.
func (e *mongoAdapter) Remove(entity types.Entity) error {
	return e.Collection.Remove(entity.ID)
}

// BulkInsert insert entities in bulk.
// Thread safe.
func (e *mongoAdapter) BulkInsert(entity types.Entity) error {
	_, err := e.mbulk.AddInsert(
		bulk.NewInsert(entity),
	)
	return err
}

// BulkUpdate update entities in bulk
// Thread safe.
func (e *mongoAdapter) BulkUpdate(entity types.Entity) error {
	_, err := e.mbulk.AddUniqUpdate(
		entity.ID,
		// We use $set here, else when we update watchers,
		// which are specialized entities,
		// we lose the watcher-specific fields
		bulk.NewUpdate(bson.M{"_id": entity.ID}, bson.M{"$set": entity}),
	)
	return err
}

func (e *mongoAdapter) AddToBulkUpdate(entityId string, data bson.M) error {
	_, err := e.mbulk.AddUpdate(bulk.NewUpdate(bson.M{"_id": entityId}, data))

	return err
}

func (e *mongoAdapter) BulkUpsert(entity types.Entity, newImpacts []string, newDepends []string) error {
	_, err := e.mbulk.AddUpsert(
		bulk.NewUpdate(bson.M{"_id": entity.ID}, entity.GetUpsertMongoBson(newImpacts, newDepends)),
	)

	return err
}

// FlushBulk force all bulk caches to be written.
func (e *mongoAdapter) FlushBulk() error {
	errUpdate := e.FlushBulkUpdate()

	if errUpdate != nil {
		return fmt.Errorf("entity adapter flushbulk update: %v", errUpdate)
	}
	return nil
}

// FlushBulkInsert force insert cache to be sent to Mongo.
// This function is thread safe.
// TODO: fusion with FlushBulkUpdate() in one single algorythm ?
func (e *mongoAdapter) FlushBulkInsert() error {
	_, err := e.mbulk.PerformInserts()
	if err != nil {
		return fmt.Errorf("entities bulk flush inserts: %v", err)
	}
	return nil
}

// FlushBulkUpdate force update cache to be sent to Mongo.
// This function is thread safe.
// TODO: fusion with FlushBulkInsert() ?
func (e *mongoAdapter) FlushBulkUpdate() error {
	_, err := e.mbulk.PerformUpdates()
	if err != nil {
		return fmt.Errorf("entities bulk flush updates: %v", err)
	}
	return nil
}

// Get is the same as GetEntityByID
// Return True if the document has been found
func (e *mongoAdapter) Get(id string) (types.Entity, bool) {
	entity, err := e.GetEntityByID(id)
	entity.EnsureInitialized()

	_, notfound := err.(errt.NotFound)
	if notfound {
		return entity, false
	} else if err != nil {
		return entity, false
	}

	return entity, true
}

// GetEntityByID finds an Entity from is eid
func (e *mongoAdapter) GetEntityByID(id string) (types.Entity, error) {
	var ent types.Entity
	err := e.Collection.GetOne(bson.M{"_id": id}, nil, &ent)
	return ent, err
}

func (e *mongoAdapter) Count() (int, error) {
	return e.Collection.Count()
}

func (e *mongoAdapter) RemoveAll() error {
	_, err := e.Collection.RemoveAll()
	return err
}

func (e *mongoAdapter) GetIDs(filter bson.M, ids *[]interface{}) error {
	return e.Collection.Get(filter, ids)
}
