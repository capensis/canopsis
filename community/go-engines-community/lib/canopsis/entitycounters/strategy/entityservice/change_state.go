package entityservice

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
)

type ChangeStateStrategy struct{}

func (s ChangeStateStrategy) CanSkip(calcData entitycounters.EntityServiceCountersCalcData) bool {
	return !calcData.CurActive && len(calcData.ServicesToAdd) == 0 && len(calcData.ServicesToRemove) == 0
}

func (s ChangeStateStrategy) Calculate(calcData entitycounters.EntityServiceCountersCalcData) entitycounters.EntityCounters {
	if calcData.ServicesToRemove[calcData.Counters.ID] {
		inheritedMode := entitycounters.InheritedNone
		if calcData.PrevInherited {
			inheritedMode = entitycounters.InheritedWith
		}

		calcData.Counters.Depends--
		calcData.Counters.DecrementState(calcData.PrevState, inheritedMode)
		calcData.Counters.DecrementAlarmCounters(calcData.IsAcked, calcData.CurActive)

		if !calcData.CurActive {
			calcData.Counters.DecrementPbhCounters(calcData.CurPbhTypeID)
		}
	} else if calcData.ServicesToAdd[calcData.Counters.ID] {
		inheritedMode := entitycounters.InheritedNone
		if calcData.CurInherited {
			inheritedMode = entitycounters.InheritedWith
		}

		calcData.Counters.Depends++
		calcData.Counters.IncrementState(calcData.CurState, inheritedMode)
		calcData.Counters.IncrementAlarmCounters(calcData.IsAcked, calcData.CurActive)

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
	}

	return calcData.Counters
}
