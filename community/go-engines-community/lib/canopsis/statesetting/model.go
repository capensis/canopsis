package statesetting

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

const (
	TypeNumber = iota
	TypePercentage
)

const (
	JUnitID   = "junit"
	ServiceID = "service"

	MethodWorst        = "worst"
	MethodWorstOfShare = "worst_of_share"
	MethodInherited    = "inherited"
	MethodDependencies = "dependencies"

	CalculationMethodNumber = "number"
	CalculationMethodShare  = "share"

	CalculationCondGT = "gt"
	CalculationCondLT = "lt"

	RuleTypeComponent = "component"
	RuleTypeService   = "service"

	StateSettingsNotificationID = "state_settings_notify_id"
)

type StateSetting struct {
	ID       string `bson:"_id"`
	Title    string `bson:"title"`
	Method   string `bson:"method"`
	Enabled  *bool  `bson:"enabled,omitempty"`
	Priority int64  `bson:"priority"`
	Type     string `bson:"type,omitempty"`

	EntityPattern          *pattern.Entity `bson:"entity_pattern,omitempty"`
	InheritedEntityPattern *pattern.Entity `bson:"inherited_entity_pattern,omitempty"`

	JUnitThresholds *JUnitThresholds `bson:"junit_thresholds,omitempty"`
	StateThresholds *StateThresholds `bson:"state_thresholds,omitempty"`
}

type StateThresholds struct {
	Critical *StateThreshold `bson:"critical,omitempty"`
	Major    *StateThreshold `bson:"major,omitempty"`
	Minor    *StateThreshold `bson:"minor,omitempty"`
	OK       *StateThreshold `bson:"ok,omitempty"`
}

type StateThreshold struct {
	Method string `bson:"method"`
	State  string `bson:"state"`
	Cond   string `bson:"cond"`
	Value  int    `bson:"value"`
}

type Counters struct {
	OK       int
	Minor    int
	Major    int
	Critical int
}

func (c Counters) GetTotal() int {
	return c.OK + c.Minor + c.Major + c.Critical
}

func (s *StateThreshold) IsReached(c Counters) bool {
	if s == nil {
		return false
	}

	if s.Method == CalculationMethodNumber {
		switch s.State {
		case types.AlarmStateTitleOK:
			return s.matchNumberCondition(c.OK)
		case types.AlarmStateTitleMinor:
			return s.matchNumberCondition(c.Minor)
		case types.AlarmStateTitleMajor:
			return s.matchNumberCondition(c.Major)
		case types.AlarmStateTitleCritical:
			return s.matchNumberCondition(c.Critical)
		}
	} else if s.Method == CalculationMethodShare {
		switch s.State {
		case types.AlarmStateTitleOK:
			return s.matchShareCondition(c.OK, c.GetTotal())
		case types.AlarmStateTitleMinor:
			return s.matchShareCondition(c.Minor, c.GetTotal())
		case types.AlarmStateTitleMajor:
			return s.matchShareCondition(c.Major, c.GetTotal())
		case types.AlarmStateTitleCritical:
			return s.matchShareCondition(c.Critical, c.GetTotal())
		}
	}

	return false
}

func (s *StateThreshold) matchNumberCondition(val int) bool {
	switch s.Cond {
	case CalculationCondGT:
		return val > s.Value
	case CalculationCondLT:
		return val < s.Value
	}

	return false
}

func (s *StateThreshold) matchShareCondition(val, all int) bool {
	return s.matchNumberCondition(int(float64(val) / float64(all) * 100))
}

type JUnitThresholds struct {
	Skipped  JUnitThreshold `bson:"skipped"`
	Errors   JUnitThreshold `bson:"errors"`
	Failures JUnitThreshold `bson:"failures"`
}

type JUnitThreshold struct {
	Minor    float64 `bson:"minor"`
	Major    float64 `bson:"major"`
	Critical float64 `bson:"critical"`
	Type     int     `bson:"type"`
}

func (s JUnitThreshold) GetState(value, total int64) int {
	if total > 0 {
		comparableValue := float64(value)
		if s.Type == TypePercentage {
			comparableValue = comparableValue / float64(total) * 100
		}

		if comparableValue > s.Critical {
			return types.AlarmStateCritical
		}

		if comparableValue > s.Major {
			return types.AlarmStateMajor
		}

		if comparableValue > s.Minor {
			return types.AlarmStateMinor
		}
	}

	return types.AlarmStateOK
}

func (s JUnitThresholds) GetState(skipped, errors, failures, total int64) int {
	worstState := s.Skipped.GetState(skipped, total)
	errorsState := s.Errors.GetState(errors, total)
	if errorsState > worstState {
		worstState = errorsState
	}

	failuresState := s.Failures.GetState(failures, total)
	if failuresState > worstState {
		worstState = failuresState
	}

	return worstState
}
