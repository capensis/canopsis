package context

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
)

// EnrichmentCenter is the high level API for context management.
type EnrichmentCenter interface {
	// Handle context creation and update on received event.
	// param source: the original event as byte slice.
	// param ef: which fields are included/excluded from extra infos enrichment.
	//     see EnrichFields doc for more infos.
	Handle(event types.Event, ef EnrichFields) *types.Entity

	Get(event types.Event) (*types.Entity, error)

	// Update updates an entity. It does not update the context graph.

	// Update updates an entity.
	// It also recomputes the link between the entity and the watchers, since
	// those may be changed by the modifications of the entity, and returns the
	// modified entity.
	Update(entity types.Entity) types.Entity

	// Flush documents that did not reiceived any updates
	Flush() error

	LoadWatchers() error

	EnrichResourceInfoWithComponentInfo(event *types.Event, entity *types.Entity) error
}

// Builder is the low level API for context management.
type Builder interface {
	// UpdateLinkedEntities creates the entity corresponding to an event, as
	// well as the entities that are linked to it (for example, the component
	// and connector linked to a resource), if they do not exist. It also adds
	// these entities to each other depends/impact lists.
	// This method DOES NOT update the links with watchers. This part is
	// handled in the UpdateWatchersLinks method.
	UpdateLinkedEntities(event types.Event) *types.Entity

	// UpdateWatchersLinks updates the links between the entities corresponding
	// to the event and the watchers.
	// This method checks if the entities corresponding to the event (resource,
	// component, watcher and/or connector) are watched by watchers, and adds
	// them to the dependencies of those watchers.
	// If the type of the event is updatewatcher and the entity is a watcher,
	// the dependencies of the entity are updated.
	// This method should be called after the Enrich method, to ensure that all
	// the entity's informations are taken into account.
	UpdateWatchersLinks(event types.Event, eventEntity *types.Entity) *types.Entity

	// Enrich adds fields from the event to the entity's informations, and
	// returns the entity.
	// The fields that are used in Enrich are specified by the fields
	// parameter.
	// Enrich should be called after UpdateLinkedEntities, to make sure that
	// the event corresponding to the entity exists. If the entity does not
	// exist, Enrich returns nil.
	Enrich(event types.Event, fields EnrichFields) *types.Entity

	// Update updates the state of the entity in the builder's cache.
	// It also recomputes the link between the entity and the watchers, since
	// those may be changed by the modifications of the entity, and returns the
	// modified entity.
	Update(entity types.Entity) types.Entity

	Extract() map[string]*EntityState
	LoadWatchers() error
}
