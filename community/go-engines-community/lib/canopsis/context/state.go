package context

import (
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
)

// Pending states:
//  OK: no modifications since last Extract
//  Insert: entity to be inserted
//  Update: entity to be updated to an existing one
const (
	PendingInsert = iota
	PendingUpdate
	PendingOK
)

// EntityState state is a Pending... state
type EntityState struct {
	Entity     types.Entity
	HasDepend  map[string]bool
	HasImpact  map[string]bool
	NewDepends []string
	NewImpacts []string
	State      int
	Updated    time.Time
}

// NewEntityState creates a new EntityState
func NewEntityState(
	entity types.Entity,
	state int,
	updated time.Time,
) *EntityState {
	var newImpacts, newDepends []string

	hasDepend := map[string]bool{}
	for _, dependID := range entity.Depends {
		if state == PendingInsert {
			newDepends = append(newDepends, dependID)
		}
		hasDepend[dependID] = true
	}

	hasImpact := map[string]bool{}
	for _, impactID := range entity.Impacts {
		if state == PendingInsert {
			newImpacts = append(newImpacts, impactID)
		}
		hasImpact[impactID] = true
	}

	return &EntityState{
		Entity:     entity,
		HasDepend:  hasDepend,
		HasImpact:  hasImpact,
		State:      state,
		Updated:    updated,
		NewImpacts: newImpacts,
		NewDepends: newDepends,
	}
}

// AddImpacts adds entity ids to the entity's impacts.
// It returns the updated EntityState, and a boolean indicating whether the
// EntityState has changed or not.
func (e *EntityState) AddImpacts(impacts ...string) bool {
	updated := false
	for _, impact := range impacts {
		if !e.HasImpact[impact] {
			e.HasImpact[impact] = true
			e.Entity.Impacts = append(e.Entity.Impacts, impact)
			e.NewImpacts = append(e.NewImpacts, impact)
			updated = true
		}
	}

	return updated
}

// AddDepends adds entity ids to the entity's depends.
// It returns the updated EntityState, and a boolean indicating whether the
// EntityState has changed or not.
func (e *EntityState) AddDepends(depends ...string) bool {
	updated := false
	for _, depend := range depends {
		if !e.HasDepend[depend] {
			e.HasDepend[depend] = true
			e.Entity.Depends = append(e.Entity.Depends, depend)
			e.NewDepends = append(e.NewDepends, depend)
			updated = true
		}
	}

	return updated
}
