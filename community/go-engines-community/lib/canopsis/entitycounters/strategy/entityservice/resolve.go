package entityservice

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type ResolveStrategy struct{}

func (s ResolveStrategy) CanSkip(_ entitycounters.EntityServiceCountersCalcData) bool {
	return false
}

func (s ResolveStrategy) Calculate(calcData entitycounters.EntityServiceCountersCalcData) entitycounters.EntityCounters {
	if calcData.ServicesToRemove[calcData.Counters.ID] {
		inheritedMode := entitycounters.InheritedNone
		if calcData.PrevInherited {
			inheritedMode = entitycounters.InheritedWith
		}

		calcData.Counters.Depends--
		calcData.Counters.DecrementState(calcData.CurState, inheritedMode)
		calcData.Counters.DecrementAlarmCounters(calcData.IsAcked, calcData.CurActive, calcData.CurPbhActive)

		if !calcData.CurPbhActive {
			calcData.Counters.DecrementPbhCounters(calcData.CurPbhTypeID)
		}

		if calcData.CurPause {
			calcData.Counters.Paused--
		}

		if !calcData.CurActive {
			calcData.Counters.InactiveDepends--
		}
	} else if calcData.ServicesToAdd[calcData.Counters.ID] {
		inheritedMode := entitycounters.InheritedNone
		if calcData.CurInherited {
			inheritedMode = entitycounters.InheritedWith
		}

		calcData.Counters.Depends++
		calcData.Counters.IncrementState(types.AlarmStateOK, inheritedMode)

		if !calcData.CurPbhActive {
			calcData.Counters.IncrementPbhCounters(calcData.CurPbhTypeID)
		}

		if calcData.CurPause {
			calcData.Counters.Paused++
		}

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

		calcData.Counters.DecrementState(calcData.CurState, prevInheritedMode)
		calcData.Counters.IncrementState(types.AlarmStateOK, curInheritedMode)

		calcData.Counters.DecrementAlarmCounters(calcData.IsAcked, calcData.CurActive, calcData.CurPbhActive)

		// todo: why not enabled???
		if !calcData.CurPbhActive && !calcData.EntityEnabled {
			calcData.Counters.DecrementPbhCounters(calcData.CurPbhTypeID)
		}
	}

	return calcData.Counters
}
