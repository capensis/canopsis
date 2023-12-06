// Package contextgraph contains a service, which is responsible for building canopsis context graph.
package contextgraph

//go:generate mockgen -destination=../../../mocks/lib/canopsis/contextgraph/contextgraph.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/contextgraph Manager,EntityServiceStorage

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type EntityServiceStorage interface {
	GetAll(ctx context.Context) ([]entityservice.EntityService, error)
	Get(ctx context.Context, serviceID string) (entityservice.EntityService, error)
}

type Manager interface {
	//HandleEvent processes canopsis contextable events, upserts new entities,
	//returns event entity as first arguments, and slice of linked entities if they're new as a second argument
	HandleEvent(ctx context.Context, event types.Event) (types.Entity, []types.Entity, error)
	//CheckServices processes slice of entities to check if they're belonged to an entity service and add this information to the impact/depends slices,
	//returns slice of the same entities, which were passed to a function plus entity service entities, which were affected.
	CheckServices(ctx context.Context, entities []types.Entity) ([]types.Entity, error)
	//UpdateImpactedServicesFromDependencies updates impacted services from dependencies info for connector entity
	UpdateImpactedServicesFromDependencies(ctx context.Context) error
	//RecomputeService recomputes context graph for an entity service
	RecomputeService(ctx context.Context, serviceID string) (types.Entity, []types.Entity, error)
	//UpdateEntities updates entities in the db, eventEntityID is needed to retrieve event entity from entities slice
	UpdateEntities(ctx context.Context, eventEntityID string, entities []types.Entity, updateServicesData bool) (types.Entity, error)
	//FillResourcesWithInfos fills all dependent component's resources with component_infos
	FillResourcesWithInfos(ctx context.Context, component types.Entity) ([]types.Entity, error)
	//UpdateLastEventDate updates last event date field in the entity document
	UpdateLastEventDate(ctx context.Context, eventType string, entityID string, timestamp datetime.CpsTime) error
}
