// Package contextgraph contains a service, which is responsible for building canopsis context graph.
package contextgraph

//go:generate go tool go.uber.org/mock/mockgen -destination=../../../mocks/lib/canopsis/contextgraph/contextgraph.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/contextgraph Manager,EntityServiceStorage

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
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
	// LoadServices refreshes slice of available services. Should be used before AssignServices calls until the service cache is implemented.
	LoadServices(ctx context.Context) error
	// AssignServices processes an entity to check if it belonged to loaded entity services and modifies it.
	AssignServices(eventEntity *types.Entity, commRegister mongo.CommandsRegister)
	// AssignServicesByInfoNames processes only services indexed by the provided info/component info updates keys.
	AssignServicesByInfoNames(eventEntity *types.Entity, infoUpdates, componentInfoUpdates map[string]eventfilter.UpdatedValue, commRegister mongo.CommandsRegister)
	// AssignStateSetting assigns a state setting for a component or a service, returns true if new state setting is assigned.
	AssignStateSetting(ctx context.Context, entity *types.Entity, commRegister mongo.CommandsRegister) (bool, error)
	// UpdateImpactedServicesFromDependencies updates impacted services from dependencies info for connector entity
	UpdateImpactedServicesFromDependencies(ctx context.Context) error
	// RecomputeService recomputes context graph for an entity service
	RecomputeService(ctx context.Context, serviceID string, commRegister mongo.CommandsRegister) (types.Entity, error)
	// ProcessComponentDependencies processes component's dependencies to update component infos or state setting parameters.
	ProcessComponentDependencies(ctx context.Context, component *types.Entity, updatedInfos map[string]eventfilter.UpdatedValue, commRegister mongo.CommandsRegister) ([]string, []types.Entity, error)
	// InheritComponentFields fills resource with component infos and check if resource is matched by component state setting.
	InheritComponentFields(resource, component *types.Entity, commRegister mongo.CommandsRegister) error
}
