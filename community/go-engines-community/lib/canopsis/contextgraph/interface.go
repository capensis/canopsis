package contextgraph

//go:generate mockgen -destination=../../../mocks/lib/canopsis/contextgraph/contextgraph.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/contextgraph Manager,EntityServiceStorage

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type EntityServiceStorage interface {
	GetAll(ctx context.Context) ([]entityservice.EntityService, error)
	Get(ctx context.Context, serviceID string) (entityservice.EntityService, error)
}

type Manager interface {
	Handle(ctx context.Context, event types.Event) (types.Entity, []types.Entity, error)

	CheckServices(ctx context.Context, entities []types.Entity) ([]types.Entity, error)

	UpdateImpactedServices(ctx context.Context) error

	RecomputeService(ctx context.Context, serviceID string) (types.Entity, []types.Entity, error)

	UpdateEntities(ctx context.Context, eventEntityID string, entities []types.Entity) (types.Entity, error)

	FillResourcesWithInfos(ctx context.Context, component types.Entity) ([]types.Entity, error)
}
