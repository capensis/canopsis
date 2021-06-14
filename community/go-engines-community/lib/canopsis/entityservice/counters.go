package entityservice

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

// GetAlarmCountersFromEvent returns counters for old alarm state and
// for new alarm state base on alarm change type.
func GetAlarmCountersFromEvent(event types.Event) (*AlarmCounters, *AlarmCounters, bool) {
	if event.Alarm == nil {
		return nil, nil, false
	}

	var oldCounters, currentCounters *AlarmCounters
	alarmCounters := getAlarmCounters(*event.Alarm,
		event.Alarm.Value.PbehaviorInfo.CanonicalType, event.Alarm.Value.PbehaviorInfo.TypeID)
	isChanged := true

	alarmChangeType := types.AlarmChangeTypeNone
	if event.AlarmChange != nil {
		alarmChangeType = event.AlarmChange.Type
	}

	switch alarmChangeType {
	case types.AlarmChangeTypeAck:
		currentCounters, oldCounters = &AlarmCounters{}, &AlarmCounters{}
		*currentCounters = alarmCounters
		*oldCounters = alarmCounters
		oldCounters.Acknowledged = 0
		oldCounters.NotAcknowledged = 1
	case types.AlarmChangeTypeAckremove:
		currentCounters, oldCounters = &AlarmCounters{}, &AlarmCounters{}
		*currentCounters = alarmCounters
		*oldCounters = alarmCounters
		oldCounters.Acknowledged = 1
		oldCounters.NotAcknowledged = 0
	case types.AlarmChangeTypeCreate, types.AlarmChangeTypeCreateAndPbhEnter:
		currentCounters = &AlarmCounters{}
		*currentCounters = alarmCounters
	case types.AlarmChangeTypePbhEnter, types.AlarmChangeTypePbhLeave, types.AlarmChangeTypePbhLeaveAndEnter:
		currentCounters, oldCounters = &AlarmCounters{}, &AlarmCounters{}
		*currentCounters = alarmCounters
		*oldCounters = getAlarmCounters(*event.Alarm,
			event.AlarmChange.PreviousPbehaviorCannonicalType, event.AlarmChange.PreviousPbehaviorTypeID)

		if (event.AlarmChange.PreviousPbehaviorCannonicalType == "" ||
			event.AlarmChange.PreviousPbehaviorCannonicalType == pbehavior.TypeActive) &&
			event.Alarm.Value.PbehaviorInfo.IsActive() {
			isChanged = false
		}
	case types.AlarmChangeTypeResolve:
		oldCounters = &AlarmCounters{}
		*oldCounters = alarmCounters
	case types.AlarmChangeTypeStateDecrease, types.AlarmChangeTypeStateIncrease, types.AlarmChangeTypeChangeState:
		currentCounters, oldCounters = &AlarmCounters{}, &AlarmCounters{}
		*currentCounters = alarmCounters
		*oldCounters = alarmCounters
		oldCounters.State = NewStateCounters(event.AlarmChange.PreviousState)
	default:
		isChanged = false
		oldCounters = &alarmCounters
		currentCounters = &alarmCounters
	}

	return oldCounters, currentCounters, isChanged
}

// GetAlarmCountersFromAlarm returns alarm counters based on alarm.
func GetAlarmCountersFromAlarm(alarm types.Alarm) AlarmCounters {
	return getAlarmCounters(alarm, alarm.Value.PbehaviorInfo.CanonicalType,
		alarm.Value.PbehaviorInfo.TypeID)
}

// getAlarmCounters returns counters base on alarm.
func getAlarmCounters(
	alarm types.Alarm,
	pbhCanonicalType, pbhType string,
) AlarmCounters {
	counters := AlarmCounters{
		All: 1,
	}

	if pbhCanonicalType == "" || pbhCanonicalType == pbehavior.TypeActive {
		counters.Alarms = 1
		counters.State = NewStateCounters(alarm.Value.State.Value)

		if alarm.Value.ACK == nil {
			counters.NotAcknowledged = 1
		} else {
			counters.Acknowledged = 1
		}
	} else {
		counters.State = NewStateCounters(types.AlarmStateOK)
		counters.PbehaviorCounters = map[string]int64{
			pbhType: 1,
		}
	}

	return counters
}

// GetServiceIDsFromEvent return services which entity is added to,
// services which entity is removed from, services which entity is kept in.
func GetServiceIDsFromEvent(event types.Event, serviceIDs []string) ([]string, []string, []string) {
	removedFromServices := event.RemovedFromServices
	addedToServices := event.AddedToServices
	services := make(map[string]bool, len(serviceIDs))
	changedServices := make(map[string]bool, len(removedFromServices)+len(addedToServices))

	for _, id := range serviceIDs {
		services[id] = true
	}
	for _, id := range removedFromServices {
		changedServices[id] = true
	}
	for _, id := range addedToServices {
		changedServices[id] = true
	}

	var unchangedServices []string
	if event.Entity != nil {
		for _, impact := range event.Entity.Impacts {
			if !services[impact] || changedServices[impact] {
				continue
			}

			unchangedServices = append(unchangedServices, impact)
		}
	}

	return addedToServices, removedFromServices, unchangedServices
}

// NewStateCounters create state counters.
func NewStateCounters(state types.CpsNumber) StateCounters {
	stateCounters := StateCounters{}
	switch state {
	case types.AlarmStateCritical:
		stateCounters.Critical = 1
	case types.AlarmStateMajor:
		stateCounters.Major = 1
	case types.AlarmStateMinor:
		stateCounters.Minor = 1
	default:
		stateCounters.Info = 1
	}

	return stateCounters
}

// StateCounters is a struct containing the number of alarms in each state that
// impact a service.
type StateCounters struct {
	Critical int64 `bson:"critical"`
	Major    int64 `bson:"major"`
	Minor    int64 `bson:"minor"`
	Info     int64 `bson:"info"`
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
// determine a service's state and output.
type AlarmCounters struct {
	// All is count of unresolved
	All int64 `bson:"all"`
	// Alarms is count of unresolved and active (by pbehavior)
	Alarms int64         `bson:"active"`
	State  StateCounters `bson:"state"`
	// Acknowledged is count of unresolved and acked and active (by pbehavior)
	Acknowledged int64 `bson:"acked"`
	// NotAcknowledged is count of unresolved and unacked and active (by pbehavior)
	NotAcknowledged int64 `bson:"unacked"`
	// PbehaviorCounters contains counters for each pbehavior type.
	PbehaviorCounters map[string]int64 `bson:"pbehavior"`
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
