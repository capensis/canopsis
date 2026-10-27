package entityservice

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type CreateStrategy struct{}

func (s CreateStrategy) CanSkip(_ entitycounters.EntityServiceCountersCalcData) bool {
	return false
}

func (s CreateStrategy) Calculate(calcData entitycounters.EntityServiceCountersCalcData) entitycounters.EntityCounters {
	if calcData.ServicesToRemove[calcData.Counters.ID] {
		inheritedMode := entitycounters.InheritedNone
		if calcData.PrevInherited {
			inheritedMode = entitycounters.InheritedWith
		}

		calcData.Counters.Depends--
		calcData.Counters.DecrementState(types.AlarmStateOK, inheritedMode)

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
		calcData.Counters.IncrementState(calcData.CurState, inheritedMode)

		calcData.Counters.IncrementAlarmCounters(false, calcData.CurActive, calcData.CurPbhActive)

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

		calcData.Counters.DecrementState(types.AlarmStateOK, prevInheritedMode)
		calcData.Counters.IncrementState(calcData.CurState, curInheritedMode)

		calcData.Counters.IncrementAlarmCounters(false, calcData.CurActive, calcData.CurPbhActive)
	}

	return calcData.Counters
}
