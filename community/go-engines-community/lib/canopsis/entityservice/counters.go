package entityservice

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

// GetAlarmCountersFromEvent returns counters for old alarm state and
// for new alarm state base on alarm change type.
func GetAlarmCountersFromEvent(event types.Event) (*AlarmCounters, *AlarmCounters, bool) {
	isChanged := true
	alarmChangeType := types.AlarmChangeTypeNone
	if event.AlarmChange != nil {
		alarmChangeType = event.AlarmChange.Type
	}

	var oldCounters, currentCounters *AlarmCounters

	if event.Alarm == nil {
		switch alarmChangeType {
		case types.AlarmChangeTypeEntityToggled:
			currentCounters = &AlarmCounters{}
			*currentCounters = getEntityCounters(event.Entity.PbehaviorInfo.CanonicalType, event.Entity.PbehaviorInfo.TypeID)
			oldCounters = currentCounters
		case types.AlarmChangeTypePbhEnter, types.AlarmChangeTypePbhLeave, types.AlarmChangeTypePbhLeaveAndEnter:
			currentCounters, oldCounters = &AlarmCounters{}, &AlarmCounters{}
			*currentCounters = getEntityCounters(event.Entity.PbehaviorInfo.CanonicalType, event.Entity.PbehaviorInfo.TypeID)
			*oldCounters = getEntityCounters(event.AlarmChange.PreviousPbehaviorCannonicalType, event.AlarmChange.PreviousPbehaviorTypeID)

			if (event.AlarmChange.PreviousPbehaviorCannonicalType == "" ||
				event.AlarmChange.PreviousPbehaviorCannonicalType == pbehavior.TypeActive) &&
				event.Entity.PbehaviorInfo.IsActive() {
				isChanged = false
			}
		default:
			currentCounters = &AlarmCounters{}
			*currentCounters = getEntityCounters(event.Entity.PbehaviorInfo.CanonicalType, event.Entity.PbehaviorInfo.TypeID)
			oldCounters = currentCounters
			isChanged = false
		}

		return oldCounters, currentCounters, isChanged
	}

	alarmCounters := getAlarmCounters(event.Alarm.Value.State.Value, event.Alarm.Value.ACK != nil,
		event.Alarm.Value.PbehaviorInfo.CanonicalType, event.Alarm.Value.PbehaviorInfo.TypeID)

	switch alarmChangeType {
	case types.AlarmChangeTypeAck:
		currentCounters, oldCounters = &AlarmCounters{}, &AlarmCounters{}
		*currentCounters = alarmCounters
		*oldCounters = alarmCounters
		if event.Alarm.Value.PbehaviorInfo.IsActive() {
			oldCounters.Acknowledged = 0
			oldCounters.NotAcknowledged = 1
			oldCounters.AcknowledgedUnderPbh = 0
		} else {
			oldCounters.Acknowledged = 0
			oldCounters.NotAcknowledged = 0
			oldCounters.AcknowledgedUnderPbh = 0
		}
	case types.AlarmChangeTypeAckremove:
		currentCounters, oldCounters = &AlarmCounters{}, &AlarmCounters{}
		*currentCounters = alarmCounters
		*oldCounters = alarmCounters
		if event.Alarm.Value.PbehaviorInfo.IsActive() {
			oldCounters.Acknowledged = 1
			oldCounters.NotAcknowledged = 0
			oldCounters.AcknowledgedUnderPbh = 0
		} else {
			oldCounters.Acknowledged = 0
			oldCounters.NotAcknowledged = 0
			oldCounters.AcknowledgedUnderPbh = 1
		}
	case types.AlarmChangeTypeCreate:
		currentCounters = &AlarmCounters{}
		*currentCounters = alarmCounters
	case types.AlarmChangeTypeCreateAndPbhEnter:
		currentCounters, oldCounters = &AlarmCounters{}, &AlarmCounters{}
		*oldCounters = getEntityCounters(event.AlarmChange.PreviousPbehaviorCannonicalType, event.AlarmChange.PreviousPbehaviorTypeID)
		*currentCounters = alarmCounters
	case types.AlarmChangeTypePbhEnter, types.AlarmChangeTypePbhLeave, types.AlarmChangeTypePbhLeaveAndEnter:
		currentCounters, oldCounters = &AlarmCounters{}, &AlarmCounters{}
		*currentCounters = alarmCounters
		*oldCounters = getAlarmCounters(event.Alarm.Value.State.Value, event.Alarm.Value.ACK != nil,
			event.AlarmChange.PreviousPbehaviorCannonicalType, event.AlarmChange.PreviousPbehaviorTypeID)

		if (event.AlarmChange.PreviousPbehaviorCannonicalType == "" ||
			event.AlarmChange.PreviousPbehaviorCannonicalType == pbehavior.TypeActive) &&
			event.Alarm.Value.PbehaviorInfo.IsActive() {
			isChanged = false
		}
	case types.AlarmChangeTypeResolve:
		currentCounters, oldCounters = &AlarmCounters{}, &AlarmCounters{}
		*oldCounters = alarmCounters
		*currentCounters = getEntityCounters(event.Entity.PbehaviorInfo.CanonicalType, event.Entity.PbehaviorInfo.TypeID)
	case types.AlarmChangeTypeStateDecrease, types.AlarmChangeTypeStateIncrease, types.AlarmChangeTypeChangeState:
		currentCounters, oldCounters = &AlarmCounters{}, &AlarmCounters{}
		*currentCounters = alarmCounters
		*oldCounters = getAlarmCounters(event.AlarmChange.PreviousState, event.Alarm.Value.ACK != nil,
			event.Alarm.Value.PbehaviorInfo.CanonicalType, event.Alarm.Value.PbehaviorInfo.TypeID)
	default:
		isChanged = false
		oldCounters = &alarmCounters
		currentCounters = &alarmCounters
	}

	return oldCounters, currentCounters, isChanged
}

// GetAlarmCountersFromAlarm returns alarm counters based on alarm.
func GetAlarmCountersFromAlarm(alarm types.Alarm) AlarmCounters {
	return getAlarmCounters(alarm.Value.State.Value, alarm.Value.ACK != nil, alarm.Value.PbehaviorInfo.CanonicalType,
		alarm.Value.PbehaviorInfo.TypeID)
}

func GetAlarmCountersFromEntity(entity types.Entity) AlarmCounters {
	return getEntityCounters(entity.PbehaviorInfo.CanonicalType, entity.PbehaviorInfo.TypeID)
}

// getAlarmCounters returns counters base on alarm.
func getAlarmCounters(
	state types.CpsNumber,
	acked bool,
	pbhCanonicalType, pbhType string,
) AlarmCounters {
	counters := AlarmCounters{
		All: 1,
	}

	if pbhCanonicalType == "" || pbhCanonicalType == pbehavior.TypeActive {
		counters.Active = 1
		counters.State = newStateCounters(state)

		if acked {
			counters.Acknowledged = 1
		} else {
			counters.NotAcknowledged = 1
		}
	} else {
		counters.State = newStateCounters(types.AlarmStateOK)
		counters.PbehaviorCounters = map[string]int64{
			pbhType: 1,
		}
		if acked {
			counters.AcknowledgedUnderPbh = 1
		}
	}

	return counters
}

// getEntityCounters returns counters base on entity.
func getEntityCounters(
	pbhCanonicalType, pbhType string,
) AlarmCounters {
	counters := AlarmCounters{}
	if pbhCanonicalType != "" && pbhCanonicalType != pbehavior.TypeActive {
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

	k := 0
	for _, id := range removedFromServices {
		if changedServices[id] {
			continue
		}

		changedServices[id] = true
		removedFromServices[k] = id
		k++
	}

	removedFromServices = removedFromServices[:k]
	k = 0
	for _, id := range addedToServices {
		if changedServices[id] {
			continue
		}

		changedServices[id] = true
		addedToServices[k] = id
		k++
	}

	addedToServices = addedToServices[:k]
	var unchangedServices []string
	if event.Entity != nil {
		for _, impact := range event.Entity.Services {
			if !services[impact] || changedServices[impact] {
				continue
			}

			unchangedServices = append(unchangedServices, impact)
		}
	}

	return addedToServices, removedFromServices, unchangedServices
}

// newStateCounters create state counters.
func newStateCounters(state types.CpsNumber) StateCounters {
	stateCounters := StateCounters{}
	switch state {
	case types.AlarmStateCritical:
		stateCounters.Critical = 1
	case types.AlarmStateMajor:
		stateCounters.Major = 1
	case types.AlarmStateMinor:
		stateCounters.Minor = 1
	default:
		stateCounters.Ok = 1
	}

	return stateCounters
}

// StateCounters is a struct containing the number of alarms in each state that
// impact a service.
type StateCounters struct {
	Critical int64 `bson:"critical"`
	Major    int64 `bson:"major"`
	Minor    int64 `bson:"minor"`
	Ok       int64 `bson:"ok"`
}

// Negate returns a new StateCounters, with all the counters negated.
func (c StateCounters) Negate() StateCounters {
	return StateCounters{
		Critical: -c.Critical,
		Major:    -c.Major,
		Minor:    -c.Minor,
		Ok:       -c.Ok,
	}
}

// Add returns a new StateCounters containing the sums of two StateCounters.
func (c StateCounters) Add(other StateCounters) StateCounters {
	return StateCounters{
		Critical: c.Critical + other.Critical,
		Major:    c.Major + other.Major,
		Minor:    c.Minor + other.Minor,
		Ok:       c.Ok + other.Ok,
	}
}

// IsZero returns true if all the counters of the StateCounters are equal to 0.
func (c StateCounters) IsZero() bool {
	return c.Critical == 0 &&
		c.Major == 0 &&
		c.Minor == 0 &&
		c.Ok == 0
}

// AlarmCounters is a struct containing various counters that are used to
// determine a service's state and output.
type AlarmCounters struct {
	// All is count of unresolved
	All int64 `bson:"all" json:"all"`
	// Active is count of unresolved and active (by pbehavior)
	Active int64         `bson:"active" json:"active"`
	State  StateCounters `bson:"state" json:"state"`
	// Acknowledged is count of unresolved and acked and active (by pbehavior)
	Acknowledged int64 `bson:"acked" json:"acked"`
	// NotAcknowledged is count of unresolved and unacked and active (by pbehavior)
	NotAcknowledged int64 `bson:"unacked" json:"unacked"`
	// AcknowledgedUnderPbh is count of unresolved and acked and under pbehavior.
	AcknowledgedUnderPbh int64 `bson:"acked_under_pbh" json:"acked_under_pbh"`
	// PbehaviorCounters contains counters for each pbehavior type.
	PbehaviorCounters map[string]int64 `bson:"pbehavior" json:"pbehavior"`
	UnderPbehavior    int64            `bson:"-" json:"under_pbh"`
	// Depends is used only for output_template.
	Depends int64 `bson:"-" json:"depends"`
}

// Negate returns a new AlarmCounters, with all the counters negated.
func (c AlarmCounters) Negate() AlarmCounters {
	pbehaviorCounters := make(map[string]int64, len(c.PbehaviorCounters))
	for pbhType, counter := range c.PbehaviorCounters {
		pbehaviorCounters[pbhType] = -counter
	}

	return AlarmCounters{
		All:                  -c.All,
		Active:               -c.Active,
		State:                c.State.Negate(),
		Acknowledged:         -c.Acknowledged,
		NotAcknowledged:      -c.NotAcknowledged,
		AcknowledgedUnderPbh: -c.AcknowledgedUnderPbh,
		PbehaviorCounters:    pbehaviorCounters,
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
		All:                  c.All + other.All,
		Active:               c.Active + other.Active,
		State:                c.State.Add(other.State),
		Acknowledged:         c.Acknowledged + other.Acknowledged,
		NotAcknowledged:      c.NotAcknowledged + other.NotAcknowledged,
		AcknowledgedUnderPbh: c.AcknowledgedUnderPbh + other.AcknowledgedUnderPbh,
		PbehaviorCounters:    pbehaviorCounters,
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

	return c.Active == 0 &&
		c.State.IsZero() &&
		c.Acknowledged == 0 &&
		c.NotAcknowledged == 0 &&
		c.AcknowledgedUnderPbh == 0 &&
		ifPbhCountersZero
}
