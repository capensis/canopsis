package statesetting

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"

const (
	TypeNumber = iota
	TypePercentage
)

const (
	TypeJUnit = "junit"
)

const (
	MethodWorst        = "worst"
	MethodWorstOfShare = "worst_of_share"
)

type StateThresholds struct {
	Minor    float64 `bson:"minor"`
	Major    float64 `bson:"major"`
	Critical float64 `bson:"critical"`
	Type     int     `bson:"type"`
}

type JUnitThresholds struct {
	Skipped  StateThresholds `bson:"skipped"`
	Errors   StateThresholds `bson:"errors"`
	Failures StateThresholds `bson:"failures"`
}

type StateSetting struct {
	ID              string           `bson:"_id"`
	Type            string           `bson:"type"`
	Method          string           `bson:"method"`
	JUnitThresholds *JUnitThresholds `bson:"junit_thresholds,omitempty"`
}

func (s StateThresholds) GetState(value, total int64) int {
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
