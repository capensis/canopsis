package watcher

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
)

// StateCounters is a struct containing the number of alarms in each state that
// impact a watcher.
type StateCounters struct {
	Critical int64
	Major    int64
	Minor    int64
	Info     int64
}

// NewStateCountersFromState returns a StateCounters that corresponds to the
// state of a single dependency.
// For example, if the state of this dependency is critical (and it does not
// have any active pbehavior), the Critical counter is equal to 1, and the
// other counters are equal to 0.
// The counters corresponding to a watcher can be obtained by adding the
// counters of all its dependencies together.
func NewStateCountersFromState(state DependencyState) StateCounters {
	if state.IsEntityActive {
		switch state.AlarmState {
		case types.AlarmStateCritical:
			return StateCounters{Critical: 1}
		case types.AlarmStateMajor:
			return StateCounters{Major: 1}
		case types.AlarmStateMinor:
			return StateCounters{Minor: 1}
		}
	}

	return StateCounters{Info: 1}
}

// Negate returns a new StateCounters, with all the counters negated.
func (c StateCounters) Negate() StateCounters {
	return StateCounters{
		Critical: -c.Critical,
		Major:    -c.Major,
		Minor:    -c.Minor,
		Info:     -c.Info,
	}
}

// Add returns a new StateCounters containing the sums of two StateCounters.
func (c StateCounters) Add(other StateCounters) StateCounters {
	return StateCounters{
		Critical: c.Critical + other.Critical,
		Major:    c.Major + other.Major,
		Minor:    c.Minor + other.Minor,
		Info:     c.Info + other.Info,
	}
}

// IsZero returns true if all the counters of the StateCounters are equal to 0.
func (c StateCounters) IsZero() bool {
	return c.Critical == 0 &&
		c.Major == 0 &&
		c.Minor == 0 &&
		c.Info == 0
}

// AlarmCounters is a struct containing various counters that are used to
// determine a watcher's state and output.
type AlarmCounters struct {
	// All is count of unresolved
	All int64
	// Alarms is count of unresolved and active (by pbehavior)
	Alarms int64
	State  StateCounters
	// Acknowledged is count of unresolved and acked and active (by pbehavior)
	Acknowledged int64
	// NotAcknowledged is count of unresolved and unacked and active (by pbehavior)
	NotAcknowledged int64
	// PbehaviorCounters contains counters for each pbehavior type.
	PbehaviorCounters map[string]int64
}

// NewAlarmCountersFromState returns an AlarmCounters that corresponds to the
// state of a single dependency.
// For example, if this dependency has a critical alarm that has been
// acknowleged, the Alarms, State.Critical and Acknowledged counters are equal
// to 1, and the other counters are equal to 0.
// The counters corresponding to a watcher can be obtained by adding the
// counters of all its dependencies together.
func NewAlarmCountersFromState(state DependencyState) AlarmCounters {
	c := AlarmCounters{}

	if state.HasAlarm {
		c.All = 1
	}

	if state.HasAlarm && state.IsEntityActive {
		c.Alarms = 1

		if state.AlarmAcknowledged {
			c.Acknowledged = 1
		} else {
			c.NotAcknowledged = 1
		}
	}

	if state.HasAlarm && !state.IsEntityActive && state.PbehaviorType != "" {
		c.PbehaviorCounters = map[string]int64{
			state.PbehaviorType: 1,
		}
	}

	c.State = NewStateCountersFromState(state)

	return c
}

// Negate returns a new AlarmCounters, with all the counters negated.
func (c AlarmCounters) Negate() AlarmCounters {
	pbehaviorCounters := make(map[string]int64, len(c.PbehaviorCounters))
	for pbhType, counter := range c.PbehaviorCounters {
		pbehaviorCounters[pbhType] = -counter
	}

	return AlarmCounters{
		All:               -c.All,
		Alarms:            -c.Alarms,
		State:             c.State.Negate(),
		Acknowledged:      -c.Acknowledged,
		NotAcknowledged:   -c.NotAcknowledged,
		PbehaviorCounters: pbehaviorCounters,
	}
}

// Add returns a new AlarmCounters containing the sums of two AlarmCounters.
func (c AlarmCounters) Add(other AlarmCounters) AlarmCounters {
	pbehaviorCounters := make(map[string]int64)
	for pbhType, counter := range c.PbehaviorCounters {
		pbehaviorCounters[pbhType] = counter
	}
	for pbhType, otherCounter := range other.PbehaviorCounters {
		if counter, ok := c.PbehaviorCounters[pbhType]; ok {
			pbehaviorCounters[pbhType] = counter + otherCounter
		} else {
			pbehaviorCounters[pbhType] = otherCounter
		}
	}

	return AlarmCounters{
		All:               c.All + other.All,
		Alarms:            c.Alarms + other.Alarms,
		State:             c.State.Add(other.State),
		Acknowledged:      c.Acknowledged + other.Acknowledged,
		NotAcknowledged:   c.NotAcknowledged + other.NotAcknowledged,
		PbehaviorCounters: pbehaviorCounters,
	}
}

// IsZero returns true if all the counters of the AlarmCounters are equal to 0.
func (c AlarmCounters) IsZero() bool {
	ifPbhCountersZero := true
	for _, v := range c.PbehaviorCounters {
		if v != 0 {
			ifPbhCountersZero = false
			break
		}
	}

	return c.Alarms == 0 &&
		c.State.IsZero() &&
		c.Acknowledged == 0 &&
		c.NotAcknowledged == 0 &&
		ifPbhCountersZero
}

// GetCountersIncrementsFromStates returns a map containing, for each watcher
// that is impacted by a change that occured on an entity, the AlarmCounters
// that should be added to the watcher's counters.
//
// This function handles the changes of the entity's alarms and pbehaviors, as
// well as the changes in the context-graph.
//
// See counters_test.go for examples for this function.
func GetCountersIncrementsFromStates(
	previousState DependencyState,
	currentState DependencyState,
) map[string]AlarmCounters {
	// oldDependencyIncrement contains the values that should be added to a
	// watcher's counters when the entity was removed from its dependencies.
	oldDependencyIncrement := NewAlarmCountersFromState(previousState).Negate()

	// newDependencyIncrement contains the values that should be added to a
	// watcher's counters when the entity was added to its dependencies.
	newDependencyIncrement := NewAlarmCountersFromState(currentState)

	// updateIncrement contains the values that should be added to a watcher's
	// counters when the entity is (and already was) one of its dependencies.
	updateIncrement := oldDependencyIncrement.Add(newDependencyIncrement)

	impacts := NewImpactDiffFromImpacts(
		previousState.ImpactedWatchers, currentState.ImpactedWatchers)

	increments := map[string]AlarmCounters{}
	for _, watcherID := range impacts.All {
		increment := AlarmCounters{}
		if impacts.Previous[watcherID] && impacts.Current[watcherID] {
			increment = updateIncrement
		} else if impacts.Current[watcherID] {
			increment = newDependencyIncrement
		} else if impacts.Previous[watcherID] {
			increment = oldDependencyIncrement
		}

		if !increment.IsZero() {
			increments[watcherID] = increment
		}
	}

	return increments
}
