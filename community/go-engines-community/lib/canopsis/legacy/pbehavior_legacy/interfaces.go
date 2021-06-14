package pbehavior_legacy

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"github.com/globalsign/mgo/bson"
)

type Adapter interface {
	// Get read pbehaviors from db
	Get(bson.M) ([]types.PBehaviorLegacy, error)

	// Insert insert pbehavior in db
	Insert(types.PBehaviorLegacy) error

	// Update insert pbehavior in db
	Update(types.PBehaviorLegacy) error

	// RemoveId remove pbehavior from db by id
	RemoveId(string) error

	// GetByEntityIds get pbehaviors from db by Entity IDs
	GetByEntityIds(eids []string, enabled bool) ([]types.PBehaviorLegacy, error)
}

// Service allows you to manipulate actions.
type Service interface {
	// Create declare a new pbheavior
	Insert(types.PBehaviorLegacy) error

	// Get find pbehaviors from db
	Get(bson.M) ([]types.PBehaviorLegacy, error)

	// Remove remove a pbehavior
	Remove(types.PBehaviorLegacy) error

	// AlarmHasPBehavior checks if an alarm has an enabled pbehavior on it. Be
	// careful, this method does not check if a pbehavior is active.
	AlarmHasPBehavior(types.Alarm) bool

	// GetByEntityIds retrieves every pbehaviors on a list of entities
	GetByEntityIds(eids []string, enabled bool) ([]types.PBehaviorLegacy, error)

	// HasActivePBehavior returns true if the entity has an active pbehavior.
	HasActivePBehavior(entityID string) bool
}
