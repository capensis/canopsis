package entityservice

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
)

type PbhLeaveAndEnterStrategy struct{}

func (s PbhLeaveAndEnterStrategy) CanSkip(calcData entitycounters.EntityServiceCountersCalcData) bool {
	return calcData.PrevActive && calcData.CurActive && len(calcData.ServicesToAdd) == 0 && len(calcData.ServicesToRemove) == 0
}

func (s PbhLeaveAndEnterStrategy) Calculate(calcData entitycounters.EntityServiceCountersCalcData) entitycounters.EntityCounters {
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
			calcData.Counters.IncrementAlarmCounters(calcData.IsAcked, calcData.CurActive)
		}

		if !calcData.CurActive {
			calcData.Counters.IncrementPbhCounters(calcData.CurPbhTypeID)
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

		if calcData.AlarmExists {
			calcData.Counters.DecrementAlarmCounters(calcData.IsAcked, calcData.PrevActive)
			calcData.Counters.IncrementAlarmCounters(calcData.IsAcked, calcData.CurActive)
		}

		if !calcData.PrevActive {
			calcData.Counters.DecrementPbhCounters(calcData.PrevPbhTypeID)
		}

		if !calcData.CurActive {
			calcData.Counters.IncrementPbhCounters(calcData.CurPbhTypeID)
		}
	}

	return calcData.Counters
}
