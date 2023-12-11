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

	CalculationNumber = "number"
	CalculationShare  = "share"

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
