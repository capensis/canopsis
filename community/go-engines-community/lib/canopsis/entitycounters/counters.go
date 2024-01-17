package entitycounters

import (
	"maps"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statesetting"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type StateCounters struct {
	Critical int `bson:"critical"`
	Major    int `bson:"major"`
	Minor    int `bson:"minor"`
	Ok       int `bson:"ok"`
}

type EntityCounters struct {
	ID                   string         `bson:"_id"`
	All                  int            `bson:"all"`
	Active               int            `bson:"active"`
	State                StateCounters  `bson:"state"`
	InheritedState       StateCounters  `bson:"inherited_state"`
	Acknowledged         int            `bson:"acked"`
	AcknowledgedUnderPbh int            `bson:"acked_under_pbh"`
	NotAcknowledged      int            `bson:"unacked"`
	PbehaviorCounters    map[string]int `bson:"pbehavior,omitempty"`
	UnderPbehavior       int            `bson:"under_pbh"`
	Depends              int            `bson:"depends"`
	OutputTemplate       string         `bson:"output_template,omitempty"`
	Output               string         `bson:"output,omitempty"`

	Rule *statesetting.StateSetting `bson:"rule"`
}

func (s *EntityCounters) Reset() {
	s.All = 0
	s.Active = 0
	s.State = StateCounters{}
	s.InheritedState = StateCounters{}
	s.Acknowledged = 0
	s.AcknowledgedUnderPbh = 0
	s.NotAcknowledged = 0
	clear(s.PbehaviorCounters)
	s.UnderPbehavior = 0
	s.Depends = 0
	s.Output = ""
}

func (s *EntityCounters) Copy() EntityCounters {
	pbhCounters := make(map[string]int, len(s.PbehaviorCounters))
	maps.Copy(pbhCounters, s.PbehaviorCounters)

	c := *s
	c.PbehaviorCounters = pbhCounters

	return c
}

func (s *EntityCounters) GetWorstState() int {
	if s.Rule != nil {
		if s.Rule.Method == statesetting.MethodDependencies {
			counters := statesetting.Counters{
				OK:       s.State.Ok,
				Minor:    s.State.Minor,
				Major:    s.State.Major,
				Critical: s.State.Critical,
			}

			criticalThresholds := s.Rule.StateThresholds.Critical
			if criticalThresholds != nil && criticalThresholds.IsReached(counters) {
				return types.AlarmStateCritical
			}

			majorThresholds := s.Rule.StateThresholds.Major
			if majorThresholds != nil && majorThresholds.IsReached(counters) {
				return types.AlarmStateMajor
			}

			minorThresholds := s.Rule.StateThresholds.Minor
			if minorThresholds != nil && minorThresholds.IsReached(counters) {
				return types.AlarmStateMinor
			}

			okThresholds := s.Rule.StateThresholds.OK
			if okThresholds != nil && okThresholds.IsReached(counters) {
				return types.AlarmStateOK
			}

			if criticalThresholds == nil && s.State.Critical > 0 {
				return types.AlarmStateCritical
			}

			if majorThresholds == nil && s.State.Major > 0 {
				return types.AlarmStateMajor
			}

			if minorThresholds == nil && s.State.Minor > 0 {
				return types.AlarmStateMinor
			}

			if okThresholds == nil {
				return types.AlarmStateOK
			}
		} else if s.Rule.Method == statesetting.MethodInherited && s.Rule.Type == statesetting.RuleTypeService {
			if s.InheritedState.Critical > 0 {
				return types.AlarmStateCritical
			}

			if s.InheritedState.Major > 0 {
				return types.AlarmStateMajor
			}

			if s.InheritedState.Minor > 0 {
				return types.AlarmStateMinor
			}

			return types.AlarmStateOK
		}
	}

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

func (s *EntityCounters) IncrementState(state int, withInherited bool) {
	switch state {
	case types.AlarmStateOK:
		s.State.Ok++
		if withInherited {
			s.InheritedState.Ok++
		}
	case types.AlarmStateMinor:
		s.State.Minor++
		if withInherited {
			s.InheritedState.Minor++
		}
	case types.AlarmStateMajor:
		s.State.Major++
		if withInherited {
			s.InheritedState.Major++
		}
	case types.AlarmStateCritical:
		s.State.Critical++
		if withInherited {
			s.InheritedState.Critical++
		}
	}
}

func (s *EntityCounters) DecrementState(state int, withInherited bool) {
	switch state {
	case types.AlarmStateOK:
		s.State.Ok--
		if withInherited {
			s.InheritedState.Ok--
		}
	case types.AlarmStateMinor:
		s.State.Minor--
		if withInherited {
			s.InheritedState.Minor--
		}
	case types.AlarmStateMajor:
		s.State.Major--
		if withInherited {
			s.InheritedState.Major--
		}
	case types.AlarmStateCritical:
		s.State.Critical--
		if withInherited {
			s.InheritedState.Critical--
		}
	}
}

func (s *EntityCounters) IncrementAlarmCounters(acked, isActive bool) {
	s.All++

	if isActive {
		s.Active++
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

func (s *EntityCounters) DecrementAlarmCounters(acked, isActive bool) {
	s.All--

	if isActive {
		s.Active--
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

func (s *EntityCounters) IncrementPbhCounters(typeID string) {
	s.UnderPbehavior++
	s.PbehaviorCounters[typeID]++
}

func (s *EntityCounters) DecrementPbhCounters(typeID string) {
	s.UnderPbehavior--
	s.PbehaviorCounters[typeID]--
}

func (s *EntityCounters) Sub(o EntityCounters) map[string]int {
	diff := make(map[string]int)

	setNotZero(diff, "all", s.All-o.All)
	setNotZero(diff, "active", s.Active-o.Active)
	setNotZero(diff, "depends", s.Depends-o.Depends)

	setNotZero(diff, "state.ok", s.State.Ok-o.State.Ok)
	setNotZero(diff, "state.minor", s.State.Minor-o.State.Minor)
	setNotZero(diff, "state.major", s.State.Major-o.State.Major)
	setNotZero(diff, "state.critical", s.State.Critical-o.State.Critical)

	setNotZero(diff, "inherited_state.ok", s.InheritedState.Ok-o.InheritedState.Ok)
	setNotZero(diff, "inherited_state.minor", s.InheritedState.Minor-o.InheritedState.Minor)
	setNotZero(diff, "inherited_state.major", s.InheritedState.Major-o.InheritedState.Major)
	setNotZero(diff, "inherited_state.critical", s.InheritedState.Critical-o.InheritedState.Critical)

	setNotZero(diff, "acked", s.Acknowledged-o.Acknowledged)
	setNotZero(diff, "acked_under_pbh", s.AcknowledgedUnderPbh-o.AcknowledgedUnderPbh)
	setNotZero(diff, "unacked", s.NotAcknowledged-o.NotAcknowledged)
	setNotZero(diff, "under_pbh", s.UnderPbehavior-o.UnderPbehavior)

	for k, v := range s.PbehaviorCounters {
		setNotZero(diff, "pbehavior."+k, v-o.PbehaviorCounters[k])
	}

	return diff
}

func setNotZero(m map[string]int, k string, v int) {
	if v == 0 {
		return
	}

	m[k] = v
}
