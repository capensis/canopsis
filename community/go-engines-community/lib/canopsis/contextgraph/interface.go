// Package contextgraph contains a service, which is responsible for building canopsis context graph.
package contextgraph

//go:generate mockgen -destination=../../../mocks/lib/canopsis/contextgraph/contextgraph.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/contextgraph Manager,EntityServiceStorage

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
)

type EntityServiceStorage interface {
	GetAll(ctx context.Context) ([]entityservice.EntityService, error)
	Get(ctx context.Context, serviceID string) (entityservice.EntityService, error)
}

type Manager interface {
	// HandleResource processes resource event.
	HandleResource(ctx context.Context, event *types.Event, commRegister mongo.CommandsRegister) (Report, error)
	// HandleComponent processes component event.
	HandleComponent(ctx context.Context, event *types.Event, commRegister mongo.CommandsRegister) (Report, error)
	// HandleConnector processes connector event.
	HandleConnector(ctx context.Context, event *types.Event, commRegister mongo.CommandsRegister) (Report, error)
	// HandleService processes service event.
	HandleService(ctx context.Context, event *types.Event, commRegister mongo.CommandsRegister) (Report, error)
	// RefreshServices refreshes slice of available services. Should be used before AssignServices calls until the service cache is implemented.
	RefreshServices(ctx context.Context) error
	// AssignServices processes slice of entities to check if they're belonged to an entity service and modifies them.
	AssignServices(eventEntity *types.Entity, commRegister mongo.CommandsRegister)
	// AssignStateSetting assigns a state setting for a component or a service, returns true if new state setting is assigned.
	AssignStateSetting(ctx context.Context, entity *types.Entity, commRegister mongo.CommandsRegister) (bool, error)
	// UpdateImpactedServicesFromDependencies updates impacted services from dependencies info for connector entity
	UpdateImpactedServicesFromDependencies(ctx context.Context) error
	// RecomputeService recomputes context graph for an entity service
	RecomputeService(ctx context.Context, serviceID string, commRegister mongo.CommandsRegister) (types.Entity, error)
	// ProcessComponentDependencies processes component's dependencies to update component infos or state setting parameters.
	ProcessComponentDependencies(ctx context.Context, component *types.Entity, commRegister mongo.CommandsRegister) ([]string, error)
	// UpdateLastEventDate updates last event date field in the entity document
	UpdateLastEventDate(ctx context.Context, event *types.Event, updateConnectorLastEventDate bool) error

	InheritComponent(resource, component *types.Entity, commRegister mongo.CommandsRegister) error
}
