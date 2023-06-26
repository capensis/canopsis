package statecounters

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type StateCounters struct {
	Critical int `bson:"critical"`
	Major    int `bson:"major"`
	Minor    int `bson:"minor"`
	Ok       int `bson:"ok"`
}

type EntityServiceCounters struct {
	ID                   string         `bson:"_id"`
	All                  int            `bson:"all"`
	Active               int            `bson:"active"`
	State                StateCounters  `bson:"state"`
	Acknowledged         int            `bson:"acked"`
	AcknowledgedUnderPbh int            `bson:"acked_under_pbh"`
	NotAcknowledged      int            `bson:"unacked"`
	PbehaviorCounters    map[string]int `bson:"pbehavior,omitempty"`
	UnderPbehavior       int            `bson:"under_pbh"`
	Depends              int            `bson:"depends"`
	OutputTemplate       string         `bson:"output_template,omitempty"`
}

func (s *EntityServiceCounters) GetWorstState() int {
	if s.State.Critical > 0 {
		return types.AlarmStateCritical
	}

	if s.State.Major > 0 {
		return types.AlarmStateMajor
	}

	if s.State.Minor > 0 {
		return types.AlarmStateMinor
	}

	return types.AlarmStateOK
}

func (s *EntityServiceCounters) IncrementState(state int) {
	switch state {
	case types.AlarmStateOK:
		s.State.Ok++
	case types.AlarmStateMinor:
		s.State.Minor++
	case types.AlarmStateMajor:
		s.State.Major++
	case types.AlarmStateCritical:
		s.State.Critical++
	}
}

func (s *EntityServiceCounters) DecrementState(state int) {
	switch state {
	case types.AlarmStateOK:
		s.State.Ok--
	case types.AlarmStateMinor:
		s.State.Minor--
	case types.AlarmStateMajor:
		s.State.Major--
	case types.AlarmStateCritical:
		s.State.Critical--
	}
}

func (s *EntityServiceCounters) IncrementAlarmCounters(state int, acked, isActive bool) {
	if isActive {
		s.Active++
		s.IncrementState(state)
	} else {
		s.IncrementState(types.AlarmStateOK)
	}

	if acked && isActive {
		s.Acknowledged++
	}

	if acked && !isActive {
		s.AcknowledgedUnderPbh++
	}

	if !acked && isActive {
		s.NotAcknowledged++
	}
}

func (s *EntityServiceCounters) DecrementAlarmCounters(state int, acked, isActive bool) {
	if isActive {
		s.Active--
		s.DecrementState(state)
	} else {
		s.DecrementState(types.AlarmStateOK)
	}

	if acked && isActive {
		s.Acknowledged--
	}

	if acked && !isActive {
		s.AcknowledgedUnderPbh--
	}

	if !acked && isActive {
		s.NotAcknowledged--
	}
}

func (s *EntityServiceCounters) IncrementPbhCounters(typeID string) {
	s.UnderPbehavior++
	s.PbehaviorCounters[typeID]++
}

func (s *EntityServiceCounters) DecrementPbhCounters(typeID string) {
	s.UnderPbehavior--
	s.PbehaviorCounters[typeID]--
}

func (s *EntityServiceCounters) Sub(o EntityServiceCounters) map[string]int {
	diff := make(map[string]int)

	diff["all"] = s.All - o.All
	diff["active"] = s.Active - o.Active
	diff["state.ok"] = s.State.Ok - o.State.Ok
	diff["state.minor"] = s.State.Minor - o.State.Minor
	diff["state.major"] = s.State.Major - o.State.Major
	diff["state.critical"] = s.State.Critical - o.State.Critical
	diff["acked"] = s.Acknowledged - o.Acknowledged
	diff["acked_under_pbh"] = s.AcknowledgedUnderPbh - o.AcknowledgedUnderPbh
	diff["unacked"] = s.NotAcknowledged - o.NotAcknowledged
	diff["under_pbh"] = s.UnderPbehavior - o.UnderPbehavior
	diff["depends"] = s.Depends - o.Depends

	for k, v := range s.PbehaviorCounters {
		diff["pbehavior."+k] = v - o.PbehaviorCounters[k]
	}

	return diff
}
