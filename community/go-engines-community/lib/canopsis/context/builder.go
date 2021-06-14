package context

import (
	"sync"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/entity"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/watcher"
	"github.com/rs/zerolog"
)

type builder struct {
	pendingOk      map[string]*EntityState
	pending        map[string]*EntityState
	entityAdapter  entity.Adapter
	watcherAdapter watcher.Adapter
	watchers       []watcher.Watcher
	watchersMutex  sync.Mutex
	logger         zerolog.Logger
}

// LoadWatchers load all the watchers into builder.watchers
func (b *builder) LoadWatchers() error {
	b.watchersMutex.Lock()
	defer b.watchersMutex.Unlock()

	// Reset the watchers list
	var newWatchers []watcher.Watcher
	err := b.watcherAdapter.GetAllValidWatchers(&newWatchers)
	b.watchers = newWatchers

	return err
}

// Extract all pending operations: new entities, updated entities.
// When calling Extract, all entities are marked as cleaned so the next
// call will not return them again.
func (b *builder) Extract() map[string]*EntityState {
	updates := b.pending

	b.reset()

	return updates
}

// reset internal states
func (b *builder) reset() {
	b.pending = make(map[string]*EntityState)
	b.pendingOk = make(map[string]*EntityState)
}

func (b *builder) getPending(id string) *EntityState {
	if v, ok := b.pending[id]; ok {
		return v
	}

	return nil
}

func (b *builder) getPendingOk(id string) *EntityState {
	if v, ok := b.pendingOk[id]; ok {
		return v
	}

	return nil
}

// get the entity if it exists
func (b *builder) get(id string) *EntityState {
	// local data that has not yet expired (extracted)
	entityState := b.getPending(id)
	if entityState != nil {
		return entityState
	}

	entityState = b.getPendingOk(id)
	if entityState != nil {
		return entityState
	}

	// remote data
	entity, exists := b.entityAdapter.Get(id)
	if exists {
		es := NewEntityState(entity, PendingOK, time.Now())
		b.setLocal(es)
		return es
	}

	return nil
}

// setLocal local cache with given entityState for further updates
func (b *builder) setLocal(entityState *EntityState) {
	id := entityState.Entity.ID

	if entityState.State != PendingOK {
		b.pending[id] = entityState
		delete(b.pendingOk, id)
	} else {
		b.pendingOk[id] = entityState
		delete(b.pending, id)
	}
}

// UpdateLinkedEntities creates the entity corresponding to an event, as well
// as the entities that are linked to it (for example, the component and
// connector linked to a resource), if they do not exist. It also adds these
// entities to each other depends/impact lists.
// This method DOES NOT update the links with watchers. This part is handled in
// the UpdateWatchersLinks method.
func (b *builder) UpdateLinkedEntities(event types.Event) *types.Entity {
	ctxInfos := event.GenerateContextInformations()
	var eventEntity *types.Entity
	for _, ctxInfo := range ctxInfos {
		var entityState *EntityState

		if ctxInfo.ID == event.GetEID() {
			entityState = b.get(ctxInfo.ID)
		}

		if entityState != nil {
			impacts := ctxInfo.Impacts
			depends := ctxInfo.Depends

			ui := entityState.AddImpacts(impacts...)
			ud := entityState.AddDepends(depends...)

			if ui || ud {
				if entityState.State == PendingOK {
					entityState.State = PendingUpdate
				}
				entityState.Updated = time.Now()
				b.setLocal(entityState)
			}

			if !ui && !ud && len(impacts) > 0 {
				impactState := b.get(impacts[0])
				impactState.State = PendingOK
				b.setLocal(impactState)
			}

			if !ud && len(depends) > 0 {
				dependState := b.get(depends[0])
				dependState.State = PendingOK
				b.setLocal(dependState)
			}
		} else {
			newEntity := ctxInfo.NewEntity()
			entityState = b.Create(newEntity)
			b.setLocal(entityState)
		}

		if entityState.Entity.ID == event.GetEID() {
			eventEntity = &entityState.Entity
		}
	}

	return eventEntity
}

// UpdateWatchersLinks updates the links between the entities corresponding to
// the event and the watchers.
// This method checks if the entities corresponding to the event (resource,
// component, watcher and/or connector) are watched by watchers, and adds them
// to the dependencies of those watchers.
// If the type of the event is updatewatcher and the entity is a watcher, the
// dependencies of the entity are updated.
// This method should be called after the Enrich method, to ensure that all the
// entity's informations are taken into account.
func (b *builder) UpdateWatchersLinks(event types.Event, eventEntity *types.Entity) *types.Entity {
	ctxInfos := event.GenerateContextInformations()
	for _, ctxInfo := range ctxInfos {
		entityState := b.get(ctxInfo.ID)
		if entityState == nil {
			entityState = &EntityState{}
		}

		// Updates the impacts and depends of both entity and watchers
		b.updateEntityWatchers(entityState)

		// We want to update only the watcher, when the current entity is the component, which is the watcher
		if event.EventType == types.EventTypeUpdateWatcher && entityState.Entity.ID == event.Component {
			b.updateWatcher(entityState)
		}

		if entityState.Entity.ID == event.GetEID() {
			eventEntity = &entityState.Entity
		}
	}

	return eventEntity
}

// updateExistingEntityState updates the provided existing EntityState, and when setLocal is true, updates it in cache
func (b *builder) updateExistingEntityState(entityState *EntityState, setLocal bool) {
	if entityState.State == PendingOK {
		entityState.State = PendingUpdate
	}
	entityState.Updated = time.Now()
	if setLocal {
		b.setLocal(entityState)
	}
}

// updateImpactsAndDepends updates impacts and depends on both provided entityStates, and updates them in cache accordingly to their respective provided updating boolean
func (b *builder) updateImpactsAndDepends(impactingEntityState, dependingEntityState *EntityState, impactingUpdating, dependingUpdating bool) {
	ui := impactingEntityState.AddImpacts(dependingEntityState.Entity.ID)
	ud := dependingEntityState.AddDepends(impactingEntityState.Entity.ID)

	if ui {
		b.updateExistingEntityState(impactingEntityState, impactingUpdating)
	}
	if ud {
		b.updateExistingEntityState(dependingEntityState, dependingUpdating)
	}
}

// updateEntityWatchers updates the links between an entity and the watchers
// that depend on it.
func (b *builder) updateEntityWatchers(entityState *EntityState) {
	b.watchersMutex.Lock()
	defer b.watchersMutex.Unlock()

	for _, aWatcher := range b.watchers {
		if aWatcher.CheckEntityInWatcher(entityState.Entity) {
			watcherState := b.get(aWatcher.ID)
			// If it doesn't exist, we don't want to re-create it
			if watcherState != nil {
				b.updateImpactsAndDepends(entityState, watcherState, true, true)
			}
		}
	}
}

func (b *builder) invalidateWatcherState(watcherState *EntityState, aWatcher watcher.Watcher) {
	for i := range watcherState.HasImpact {
		delete(watcherState.HasImpact, i)
	}

	for d := range watcherState.HasDepend {
		delete(watcherState.HasDepend, d)
	}

	for _, i := range aWatcher.Impacts {
		watcherState.HasImpact[i] = true
	}

	for _, i := range aWatcher.Depends {
		watcherState.HasDepend[i] = true
	}
}

// updateWatcher updates the watcher linked to the provided watcher EntityState
func (b *builder) updateWatcher(watcherState *EntityState) {
	var watchedEntities []types.Entity
	var aWatcher watcher.Watcher
	// If we get a updatewatcher event,
	// we'll update the watcher list
	err := b.LoadWatchers()
	if err != nil {
		b.logger.Warn().Err(err).Msg("enrich context from update watcher event, updating the watcher list")
		return
	}

	err = b.watcherAdapter.GetByID(watcherState.Entity.ID, &aWatcher)
	if err != nil {
		b.logger.Warn().Err(err).Msg("enrich context from update watcher event, getting the watcher")
		return
	}

	err = b.watcherAdapter.GetEntities(aWatcher, &watchedEntities)
	if err != nil {
		b.logger.Warn().Err(err).Msg("enrich context from update watcher event, getting watched entities")
		return
	}

	b.invalidateWatcherState(watcherState, aWatcher)

	// building links between the watcher and the watched entities
	for _, watchedEntity := range watchedEntities {
		watchedEntityState := b.get(watchedEntity.ID)
		if watchedEntityState == nil {
			watchedEntityState = &EntityState{}
		}
		b.updateImpactsAndDepends(watchedEntityState, watcherState, true, false)
	}
	// We setLocal here and not in the loop because watcherState changes often in the loop
	b.setLocal(watcherState)
}

// Create the entity without updating impacts/depends
func (b *builder) Create(entity types.Entity) *EntityState {
	return NewEntityState(entity, PendingInsert, time.Now())
}

// Enrich adds fields from the event to the entity's informations, and returns
// the entity.
// The fields that are used in Enrich are specified by the fields parameter.
// Enrich should be called after UpdateLinkedEntities, to make sure that the
// event corresponding to the entity exists. If the entity does not exist,
// Enrich returns nil.
func (b *builder) Enrich(event types.Event, fields EnrichFields) *types.Entity {
	eid := event.GetEID()
	entityState := b.get(eid)

	if entityState == nil {
		return nil
	}

	updated := false

	for key, val := range event.ExtraInfos {
		if !fields.Allow(key) {
			continue
		}

		value, err := types.InterfaceToString(val)

		if err != nil {
			b.logger.Warn().Err(err).Str("event", eid).Str("key", key).Msg("enrich context from event")
		}

		info := types.NewInfo(key, "", value)
		oldVal, exists := entityState.Entity.Infos[info.Name]

		if !exists || (exists && oldVal.Value != info.Value) {
			updated = true
			entityState.Entity.Infos[info.Name] = info
		}
	}

	if updated {
		if entityState.State == PendingOK {
			entityState.State = PendingUpdate
		}
		entityState.Updated = time.Now()
		b.setLocal(entityState)
	}

	return &entityState.Entity
}

// Update updates the state of the entity in the builder's cache.
// It also recomputes the link between the entity and the watchers, since those
// may be changed by the modifications of the entity, and returns the modified
// entity.
func (b *builder) Update(entity types.Entity) types.Entity {
	entityState := b.get(entity.ID)
	if entityState == nil {
		entityState = b.Create(entity)
	} else {
		entityState.Entity = entity
		if entityState.State == PendingOK {
			entityState.State = PendingUpdate
		}
		entityState.Updated = time.Now()
		b.setLocal(entityState)
	}

	b.updateEntityWatchers(entityState)
	return entityState.Entity
}

// NewBuilder creates a new context builder.
func NewBuilder(entityAdapter entity.Adapter, watcherAdapter watcher.Adapter, logger zerolog.Logger) Builder {
	b := builder{
		entityAdapter:  entityAdapter,
		watcherAdapter: watcherAdapter,
		logger:         logger,
	}
	b.reset()
	return &b
}
