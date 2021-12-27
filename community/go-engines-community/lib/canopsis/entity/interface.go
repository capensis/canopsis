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
	Insert(ctx context.Context, entity types.Entity) error

	Update(ctx context.Context, entity types.Entity) error

	Get(ctx context.Context, id string) (types.Entity, bool)

	GetIDs(ctx context.Context, filter map[string]interface{}, ids *[]interface{}) error

	GetEntityByID(ctx context.Context, id string) (types.Entity, error)

	Count(ctx context.Context) (int, error)

	UpsertMany(ctx context.Context, entities []types.Entity) (map[string]bool, error)

	AddImpacts(ctx context.Context, ids []string, impacts []string) error

	RemoveImpacts(ctx context.Context, ids []string, impacts []string) error

	AddImpactByQuery(ctx context.Context, query interface{}, impact string) ([]string, error)

	RemoveImpactByQuery(ctx context.Context, query interface{}, impact string) ([]string, error)

	AddInfos(ctx context.Context, id string, infos map[string]types.Info) (bool, error)

	UpdateComponentInfos(ctx context.Context, id, componentID string) (map[string]types.Info, error)

	UpdateComponentInfosByComponent(ctx context.Context, componentID string) ([]string, error)

	UpdateLastEventDate(ctx context.Context, ids []string, time types.CpsTime) error
	UpdateIdleFields(ctx context.Context, id string, idleSince *types.CpsTime, lastIdleRuleApply string) error

	FindByIDs(ctx context.Context, ids []string) ([]types.Entity, error)

	GetAllWithLastUpdateDateBefore(ctx context.Context, time types.CpsTime, exclude []string) (mongo.Cursor, error)
	FindConnectorForComponent(ctx context.Context, id string) (*types.Entity, error)

	FindConnectorForResource(ctx context.Context, id string) (*types.Entity, error)
	FindComponentForResource(ctx context.Context, id string) (*types.Entity, error)

	GetWithIdleSince(ctx context.Context) (mongo.Cursor, error)

	GetImpactedServicesInfo(ctx context.Context) (mongo.Cursor, error)

	Bulk(ctx context.Context, models []mongodriver.WriteModel) error

	FindToCheckPbehaviorInfo(ctx context.Context, idsWithPbehaviors []string, exceptIds []string) (mongo.Cursor, error)

	UpdatePbehaviorInfo(ctx context.Context, id string, info types.PbehaviorInfo) error
}

// Service glue Adapter and Cache together
type Service interface {
	Adapter

	FlushCache() error
}
