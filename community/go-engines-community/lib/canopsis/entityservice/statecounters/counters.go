package statecounters

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"

type StateCounters struct {
	Critical int64 `bson:"critical"`
	Major    int64 `bson:"major"`
	Minor    int64 `bson:"minor"`
	Info     int64 `bson:"info"`
}

type EntityServiceCounters struct {
	ID                string           `bson:"_id"`
	All               int              `bson:"all"`
	Alarms            int              `bson:"active"`
	State             StateCounters    `bson:"state"`
	Acknowledged      int              `bson:"acked"`
	NotAcknowledged   int              `bson:"unacked"`
	PbehaviorCounters map[string]int64 `bson:"pbehavior"`
}

func (s EntityServiceCounters) GetWorstState() int {
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
		s.State.Info++
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
		s.State.Info--
	case types.AlarmStateMinor:
		s.State.Minor--
	case types.AlarmStateMajor:
		s.State.Major--
	case types.AlarmStateCritical:
		s.State.Critical--
	}
}

func (s *EntityServiceCounters) IncrementAlarmCounters(state int, acked bool) {
	s.Alarms++
	s.IncrementState(state)
	if acked {
		s.Acknowledged++
	} else {
		s.NotAcknowledged++
	}
}

func (s *EntityServiceCounters) DecrementAlarmCounters(state int, acked bool) {
	s.Alarms--
	s.DecrementState(state)
	if acked {
		s.Acknowledged--
	} else {
		s.NotAcknowledged--
	}
}
