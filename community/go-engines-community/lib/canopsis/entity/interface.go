package entity

//go:generate mockgen -destination=../../../mocks/lib/canopsis/entity/adapter.go git.canopsis.net/canopsis/go-engines/lib/canopsis/entity Adapter

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"github.com/globalsign/mgo/bson"
)

// Adapter ...
type Adapter interface {
	Insert(entity types.Entity) error

	Update(entity types.Entity) error

	Remove(entity types.Entity) error

	BulkInsert(types.Entity) error

	BulkUpdate(types.Entity) error

	AddToBulkUpdate(entityId string, data bson.M) error

	BulkUpsert(types.Entity, []string, []string) error

	FlushBulk() error

	FlushBulkInsert() error

	FlushBulkUpdate() error

	Get(id string) (types.Entity, bool)

	GetIDs(filter bson.M, ids *[]interface{}) error

	GetEntityByID(id string) (types.Entity, error)

	Count() (int, error)

	RemoveAll() error
}

// Service glue Adapter and Cache together
type Service interface {
	Adapter

	FlushCache() error
}
