package context

//go:generate mockgen -destination=../../../mocks/lib/canopsis/context/context.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/context EnrichmentCenter

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

// EnrichmentCenter is the high level API for context management.
type EnrichmentCenter interface {
	// Handle context creation and update on received event.
	Handle(ctx context.Context, event types.Event) (*types.Entity, UpdatedEntityServices, error)

	// HandleEntityServiceUpdate updates context graph for entity service.
	HandleEntityServiceUpdate(ctx context.Context, serviceID string) (*UpdatedEntityServices, error)

	// Get finds entity for event.
	Get(ctx context.Context, event types.Event) (*types.Entity, error)

	// UpdateEntityInfos updates an entity.
	// It also recomputes the link between the entity and the services, since
	// those may be changed by the modifications of the entity, and returns the
	// modified entity.
	UpdateEntityInfos(ctx context.Context, entity *types.Entity) (UpdatedEntityServices, error)

	// UpdateImpactedServices updates the impacted_services field in connector documents
	UpdateImpactedServices(ctx context.Context) error

	// LoadServices loads services from storage to cache.
	LoadServices(ctx context.Context) error

	// ReloadService loads service from storage to cache.
	ReloadService(ctx context.Context, serviceID string) error
}

// UpdatedEntityServices represents changes in context graph.
type UpdatedEntityServices struct {
	// AddedTo contains ids of entity services to which entity has been
	// added as dependency.
	AddedTo []string
	// RemovedFrom contains ids of entity services from which entity has been
	// removed as dependency.
	RemovedFrom []string
	// UpdatedComponentInfosResources contains ids of entities which component infos
	// were updated on component event.
	UpdatedComponentInfosResources []string
}

func (l UpdatedEntityServices) Add(r UpdatedEntityServices) UpdatedEntityServices {
	res := UpdatedEntityServices{}
	res.AddedTo = append(res.AddedTo, l.AddedTo...)
	res.AddedTo = append(res.AddedTo, r.AddedTo...)
	res.RemovedFrom = append(res.RemovedFrom, l.RemovedFrom...)
	res.RemovedFrom = append(res.RemovedFrom, r.RemovedFrom...)
	res.UpdatedComponentInfosResources = append(res.UpdatedComponentInfosResources, l.UpdatedComponentInfosResources...)
	res.UpdatedComponentInfosResources = append(res.UpdatedComponentInfosResources, r.UpdatedComponentInfosResources...)
	return res
}
