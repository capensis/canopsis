// Package watcher implements the watcher service used by the watcher engine to
// compute the watchers' states.
//
// Read https://doc.canopsis.net/guide-administration/moteurs/moteur-watcher/
// for more details on watchers.
//
// A previous implementation of the watcher engine used to get all the alarms
// that impact a watcher from the database each time its state needed to be
// recomputed. This did not scale well for watchers with a lot dependencies.
//
// The current implementation was written to make the most common operation
// (updating the watcher's states after an event) fast and scalable. It works
// by storing counters (number of alarm in each state, number of acknowledged
// alarms, ...) for each watcher in redis, updating those counters when an
// event is received, and computing the watchers' states from these counters.
//
// Multiple types were defined for this implementation:
//
//  - Watcher (defined in model.go) is the type used to represent a watcher in
//    the MongoDB database.
//  - AlarmCounters contains various counters for a watcher (number of alarm
//    in each state, number of acknowledged alarms, ...), and is used to
//    compute the watchers' states and outputs.
//  - DependencyState contains all the informations on an entity that affect
//    the watchers' states. It is stored in redis, and is used to update the
//    watchers' counters.
//  - ImpactDiff is used to keep track of changes in the context-graph (when
//    an entity is added or removed from a watcher's dependencies).
//  - CountersCache is the service that stores the DependencyStates of the
//    entities and updates the AlarmCounters of the watchers in redis.
//
// When an event is received by the watcher engine:
//
//  - it is converted into a DependencyState for the corresponding entity (in
//    Service.Process);
//  - the DependencyState is compared to the previous DependencyState of this
//    entity, that is stored in redis (in CountersCache.Process);
//  - the impacted watchers' AlarmCounters are updated in redis (in
//    CountersCache.Process);
//  - the state and output of the watchers are deduced from the AlarmCounters,
//    and sent to the axe engine (in Service.Process).
package watcher

import (
	"context"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"github.com/globalsign/mgo"
	"github.com/streadway/amqp"
)

// AnnotatedEntity is a struct containing an entity, a list of pbehaviors that
// impact it, and the ongoing alarm on this entity (or nil if there isn't one).
// The pbehaviors may or may not be active.
type AnnotatedEntity struct {
	types.Entity `bson:"entity"`
	Alarm        *types.Alarm `bson:"alarm"`
}

// Adapter is an interface that provides methods for database queries regarding
// watchers and their dependencies.
type Adapter interface {
	// GetAll gets all watchers from database
	GetAll(watchers *[]Watcher) error

	// GetAllValidWatchers gets all valid and enabled watchers from database
	GetAllValidWatchers(watchers *[]Watcher) error

	// GetByID finds the watcher by its entity id
	GetByID(id string, watchers *Watcher) error

	// GetEntities gets the entities watched by the watcher
	GetEntities(watcher Watcher, entities *[]types.Entity) error

	// GetAllAnnotatedEntities returns a slice containing all the entities,
	// annotated with their pbehaviors and alarms.
	// Note that this method may allocate a lot of memory. This may be
	// optimised by returning an mgo.Iter object instead of a slice. As of this
	// writing, the code complexity required to use an iterator outweighs the
	// benefits of doing so.
	GetAllAnnotatedEntities() ([]AnnotatedEntity, error)

	// GetAnnotatedEntitiesIter returns a mgo.Iter to iterate entities,
	// annotated with their pbehaviors and alarms to avoid allocating a lot of
	// memory by GetAllAnnotatedEntities.
	GetAnnotatedEntitiesIter() *mgo.Iter

	// GetAnnotatedDependencies returns a slice containing the dependencies of
	// a watcher annotated with their pbehaviors and alarms.
	// Note that this method may allocate a lot of memory. This may be
	// optimised by returning an mgo.Iter object instead of a slice. As of this
	// writing, the code complexity required to use an iterator outweighs the
	// benefits of doing so.
	GetAnnotatedDependencies(watcherID string) ([]AnnotatedEntity, error)

	Update(watcherID string, update interface{}) error

	// GetAnnotatedDependenciesIter returns a mgo.Iter to iterate dependencies of
	// a watcher annotated with their pbehaviors and alarms to avoid allocating a lot of
	// memory by GetAnnotatedDependencies.
	GetAnnotatedDependenciesIter(watcherID string) *mgo.Iter
}

// CountersCache is a type that handles the counting of the alarms and entities
// that impact watchers.
type CountersCache interface {
	// ProcessState processes and entity's update, and update the counters of
	// the watchers that are (or used to be) impacted by it.
	// It returns a map containing the impacted watchers and the new values of
	// their counters.
	// This method should be safe to run in multiple goroutines, since it is
	// used in the PeriodicalProcess and WorkerProcess of the watcher.
	ProcessState(currentState DependencyState) (map[string]AlarmCounters, error)
}

// Service is the interface implemented by the watcher service, that processes
// events, computes the watchers' states accordingly, and sends events to the
// axe engine so that the watchers' alarms are updated.
// This service does not handle the links (impact and depends) between the
// watcher and their dependencies. This is handled by the context builder,
// which is used by the che engine.
type Service interface {
	// Process updates the watchers impacted by the event and alarmChange
	// This method should be called on each event published by the axe engine,
	// so that the watchers' states are updated in real time.
	Process(ctx context.Context, event types.Event) error

	ProcessRpc(ctx context.Context, event types.Event) error

	ProcessResolvedAlarm(ctx context.Context, alarm types.Alarm, entity types.Entity) error

	// UpdateWatcher updates the state of a watcher given its ID.
	// FIXME: the current implementation of this method does not handle
	// dependencies being removed from the watcher.
	UpdateWatcher(ctx context.Context, watcherID string) error

	// ComputeAllWatchers updates the states of all watchers.
	// As of this writing, calling this method regularly (at the beat) is
	// necessary to handle changes in the pbehaviors and in the context-graph,
	// since those changes do not trigger an event.
	// FIXME: the current implementation of this method does not handle
	// entities being disabled or removed from the database.
	ComputeAllWatchers(ctx context.Context) error

	FlushDB() error
}

// AmqpChannelPublisher is an interface that represents a non-consumable AMQP
// channel. This interface is implemented by amqp.Channel. It should be used
// in services that only publish to amqp, in order to be able to test them
// easily by mocking this interface.
type AmqpChannelPublisher interface {
	// Publish sends an amqp.Publishing from the client to an exchange on the server.
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
}
