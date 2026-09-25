package entityservice

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
)

type PauseStrategy struct{}

func (s PauseStrategy) CanSkip(_ entitycounters.EntityServiceCountersCalcData) bool {
	return false
}

func (s PauseStrategy) Calculate(calcData entitycounters.EntityServiceCountersCalcData) entitycounters.EntityCounters {
	if calcData.ServicesToRemove[calcData.Counters.ID] {
		inheritedMode := entitycounters.InheritedNone
		if calcData.PrevInherited {
			inheritedMode = entitycounters.InheritedWith
		}

		calcData.Counters.Depends--
		calcData.Counters.DecrementState(calcData.PrevState, inheritedMode)

		if calcData.AlarmExists {
			calcData.Counters.DecrementAlarmCounters(calcData.IsAcked, calcData.PrevActive, calcData.PrevPbhActive)
		}

		if !calcData.PrevActive {
			calcData.Counters.InactiveDepends--
		}

		if !calcData.CurPbhActive {
			calcData.Counters.DecrementPbhCounters(calcData.CurPbhTypeID)
		}
	} else if calcData.ServicesToAdd[calcData.Counters.ID] {
		inheritedMode := entitycounters.InheritedNone
		if calcData.CurInherited {
			inheritedMode = entitycounters.InheritedWith
		}

		calcData.Counters.Depends++
		calcData.Counters.IncrementState(calcData.CurState, inheritedMode)

		if calcData.AlarmExists {
			calcData.Counters.IncrementAlarmCounters(calcData.IsAcked, calcData.CurActive, calcData.CurPbhActive)
		}

		if !calcData.CurPbhActive {
			calcData.Counters.IncrementPbhCounters(calcData.CurPbhTypeID)
		}

		calcData.Counters.Paused++

		if !calcData.CurActive {
			calcData.Counters.InactiveDepends++
		}
	} else {
		prevInheritedMode := entitycounters.InheritedNone
		if calcData.PrevInherited {
			prevInheritedMode = entitycounters.InheritedWith
		}
		curInheritedMode := entitycounters.InheritedNone
		if calcData.CurInherited {
			curInheritedMode = entitycounters.InheritedWith
		}

		calcData.Counters.DecrementState(calcData.PrevState, prevInheritedMode)
		calcData.Counters.IncrementState(calcData.CurState, curInheritedMode)

		calcData.Counters.Paused++

		if calcData.PrevActive && !calcData.CurActive {
			calcData.Counters.InactiveDepends++
		}

		if calcData.AlarmExists {
			calcData.Counters.DecrementAlarmCounters(calcData.IsAcked, calcData.PrevActive, calcData.PrevPbhActive)
			calcData.Counters.IncrementAlarmCounters(calcData.IsAcked, calcData.CurActive, calcData.CurPbhActive)
		}
	}

	return calcData.Counters
}
