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
		calcData.Counters.Depends--
		calcData.Counters.DecrementState(calcData.CurState, calcData.Inherited)
		calcData.Counters.DecrementAlarmCounters(calcData.IsAcked, calcData.CurActive)

		if !calcData.CurActive {
			calcData.Counters.DecrementPbhCounters(calcData.CurPbhTypeID)
		}
	} else if calcData.ServicesToAdd[calcData.Counters.ID] {
		calcData.Counters.Depends++
		calcData.Counters.IncrementState(types.AlarmStateOK, calcData.Inherited)

		if !calcData.CurActive {
			calcData.Counters.IncrementPbhCounters(calcData.CurPbhTypeID)
		}
	} else {
		calcData.Counters.DecrementState(calcData.CurState, calcData.Inherited)
		calcData.Counters.IncrementState(types.AlarmStateOK, calcData.Inherited)

		calcData.Counters.DecrementAlarmCounters(calcData.IsAcked, calcData.CurActive)

		// todo: why not enabled???
		if !calcData.CurActive && !calcData.EntityEnabled {
			calcData.Counters.DecrementPbhCounters(calcData.CurPbhTypeID)
		}
	}

	return calcData.Counters
}
