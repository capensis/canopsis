package watcher

import "git.canopsis.net/canopsis/go-engines/lib/canopsis/types"

const (
	// MethodWorst is the method that returns the worst state among the alarms
	// that impact a watcher.
	MethodWorst = "worst"
)

// worst returns the worst state among the alarms that impact a watcher.
func worst(counters AlarmCounters) int {
	if counters.State.Critical > 0 {
		return types.AlarmStateCritical
	}
	if counters.State.Major > 0 {
		return types.AlarmStateMajor
	}
	if counters.State.Minor > 0 {
		return types.AlarmStateMinor
	}
	return types.AlarmStateOK
}
