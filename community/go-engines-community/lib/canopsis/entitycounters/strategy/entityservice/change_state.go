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
		calcData.Counters.Depends--
		calcData.Counters.DecrementState(calcData.PrevState, calcData.Inherited)
		calcData.Counters.DecrementAlarmCounters(calcData.IsAcked, calcData.CurActive)

		if !calcData.CurActive {
			calcData.Counters.DecrementPbhCounters(calcData.CurPbhTypeID)
		}
	} else if calcData.ServicesToAdd[calcData.Counters.ID] {
		calcData.Counters.Depends++
		calcData.Counters.IncrementState(calcData.CurState, calcData.Inherited)
		calcData.Counters.IncrementAlarmCounters(calcData.IsAcked, calcData.CurActive)

		if !calcData.CurActive {
			calcData.Counters.IncrementPbhCounters(calcData.CurPbhTypeID)
		}
	} else {
		calcData.Counters.DecrementState(calcData.PrevState, calcData.Inherited)
		calcData.Counters.IncrementState(calcData.CurState, calcData.Inherited)
	}

	return calcData.Counters
}
