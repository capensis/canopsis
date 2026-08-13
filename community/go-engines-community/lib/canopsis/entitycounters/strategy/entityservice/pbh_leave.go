package entityservice

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
)

type PbhLeaveStrategy struct{}

func (s PbhLeaveStrategy) CanSkip(calcData entitycounters.EntityServiceCountersCalcData) bool {
	return calcData.PrevActive && len(calcData.ServicesToAdd) == 0 && len(calcData.ServicesToRemove) == 0
}

func (s PbhLeaveStrategy) Calculate(calcData entitycounters.EntityServiceCountersCalcData) entitycounters.EntityCounters {
	if calcData.ServicesToRemove[calcData.Counters.ID] {
		inheritedMode := entitycounters.InheritedNone
		if calcData.PrevInherited {
			inheritedMode = entitycounters.InheritedWith
		}

		calcData.Counters.Depends--
		calcData.Counters.DecrementState(calcData.PrevState, inheritedMode)

		if calcData.AlarmExists {
			calcData.Counters.DecrementAlarmCounters(calcData.IsAcked, calcData.PrevActive)
		}

		if !calcData.PrevActive {
			calcData.Counters.DecrementPbhCounters(calcData.PrevPbhTypeID)
		}
	} else if calcData.ServicesToAdd[calcData.Counters.ID] {
		inheritedMode := entitycounters.InheritedNone
		if calcData.CurInherited {
			inheritedMode = entitycounters.InheritedWith
		}

		calcData.Counters.Depends++
		calcData.Counters.IncrementState(calcData.CurState, inheritedMode)

		if calcData.AlarmExists {
			calcData.Counters.IncrementAlarmCounters(calcData.IsAcked, true)
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

		if !calcData.PrevActive {
			calcData.Counters.DecrementPbhCounters(calcData.PrevPbhTypeID)

			if calcData.AlarmExists {
				calcData.Counters.DecrementAlarmCounters(calcData.IsAcked, false)
				calcData.Counters.IncrementAlarmCounters(calcData.IsAcked, true)
			}
		}
	}

	return calcData.Counters
}
