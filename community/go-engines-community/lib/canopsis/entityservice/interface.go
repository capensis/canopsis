package entityservice

//go:generate mockgen -destination=../../../mocks/lib/canopsis/entityservice/entityservice.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice Adapter,CountersCache,Storage

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

// Adapter is an interface that provides methods for database queries regarding
// services and their dependencies.
type Adapter interface {
	GetAll(ctx context.Context) ([]EntityService, error)

	GetEnabled(ctx context.Context) ([]EntityService, error)

	GetByID(ctx context.Context, id string) (*EntityService, error)

	// UpdateCounters saves service counters to storage.
	UpdateCounters(context.Context, string, AlarmCounters) error

	// UpdateBulk bulk update
	UpdateBulk(ctx context.Context, writeModels []mongodriver.WriteModel) error

	GetServiceDependencies(ctx context.Context, serviceID string) (mongo.Cursor, error)

	GetDependenciesCount(ctx context.Context, serviceID string) (int64, error)

	AddToService(ctx context.Context, serviceId string, ids []string) error
	RemoveFromService(ctx context.Context, serviceId string, ids []string) error
	AddToServiceByQuery(ctx context.Context, serviceId string, query bson.M) (int64, error)
	RemoveFromServiceByQuery(ctx context.Context, serviceId string, query bson.M) (int64, error)
}

// Manager is used to implement context graph modifier for entity service.
type Manager interface {
	// LoadServices loads services from storage.
	LoadServices(ctx context.Context) error
	// UpdateServices adds or removes entities.
	UpdateServices(ctx context.Context, entities []types.Entity) (addedTo map[string][]string, removedFrom map[string][]string, err error)
	// UpdateService adds service to matched entities impacts and updates its depends.
	UpdateService(ctx context.Context, serviceID string) (isUpdated bool, removedFrom []string, err error)
	// ReloadService loads service from storage to cache.
	ReloadService(ctx context.Context, serviceID string) error
	// HasEntityServiceByComponentInfos returns true of at least one entity service patterns
	// contains component infos condition.
	HasEntityServiceByComponentInfos(ctx context.Context) (bool, error)
}

// CountersCache is a type that handles the counting of the alarms and entities
// that impact services.
type CountersCache interface {
	Update(context.Context, map[string]AlarmCounters) (map[string]AlarmCounters, error)

	Replace(context.Context, string, AlarmCounters) error

	Remove(context.Context, string) error

	RemoveAndGet(context.Context, string) (*AlarmCounters, error)

	ClearAll(context.Context) error

	KeepOnly(ctx context.Context, ids []string) error
}

type Service interface {
	// Process updates the services impacted by the event and alarmChange
	// This method should be called on each event published by the axe engine
	// and context graph update by the che engine,
	// so that the services' states are updated in real time.
	Process(ctx context.Context, event types.Event) error

	ProcessRpc(ctx context.Context, event types.Event) error

	// UpdateService updates the state of a service given its ID.
	UpdateService(ctx context.Context, event types.Event) error

	// ReloadService loads service from storage to cache.
	ReloadService(ctx context.Context, serviceID string) error

	// ComputeAllServices updates the states of all services.
	ComputeAllServices(ctx context.Context) error

	RecomputeIdleSince(parentCtx context.Context) error

	// ClearCache clears all counters of entity services from cache.
	ClearCache(ctx context.Context) error
}
