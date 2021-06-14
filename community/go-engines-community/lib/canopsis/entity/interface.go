package entity

//go:generate mockgen -destination=../../../mocks/lib/canopsis/entity/adapter.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity Adapter

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

// Adapter ...
type Adapter interface {
	Insert(entity types.Entity) error

	Update(entity types.Entity) error

	Remove(entity types.Entity) error

	BulkInsert(types.Entity) error

	AddToBulkUpdate(entityId string, data interface{}) error

	BulkUpsert(types.Entity, []string, []string) error

	FlushBulk() error

	FlushBulkInsert() error

	FlushBulkUpdate() error

	Get(id string) (types.Entity, bool)

	GetIDs(filter map[string]interface{}, ids *[]interface{}) error

	GetEntityByID(id string) (types.Entity, error)

	Count() (int, error)

	RemoveAll() error

	UpsertMany([]types.Entity) (map[string]bool, error)

	AddImpacts(ids []string, impacts []string) error

	RemoveImpacts(ids []string, impacts []string) error

	AddImpactByQuery(query interface{}, impact string) ([]string, error)

	RemoveImpactByQuery(query interface{}, impact string) ([]string, error)

	AddInfos(id string, infos map[string]types.Info) (bool, error)

	UpdateComponentInfos(id, componentID string) (map[string]types.Info, error)

	UpdateComponentInfosByComponent(componentID string) ([]string, error)

	UpdateLastEventDate(ids []string, time types.CpsTime) error
	UpdateIdleFields(ctx context.Context, id string, idleSince *types.CpsTime, lastIdleRuleApply string) error

	FindByIDs(ids []string) ([]types.Entity, error)

	GetAllWithLastUpdateDateBefore(ctx context.Context, time types.CpsTime, exclude []string) (mongo.Cursor, error)
	FindConnectorForComponent(ctx context.Context, id string) (*types.Entity, error)

	FindConnectorForResource(ctx context.Context, id string) (*types.Entity, error)
	FindComponentForResource(ctx context.Context, id string) (*types.Entity, error)

	GetWithIdleSince(ctx context.Context) (mongo.Cursor, error)

	GetImpactedServicesInfo(ctx context.Context) (mongo.Cursor, error)
	
	Bulk(ctx context.Context, models []mongodriver.WriteModel) error
}

// Service glue Adapter and Cache together
type Service interface {
	Adapter

	FlushCache() error
}
