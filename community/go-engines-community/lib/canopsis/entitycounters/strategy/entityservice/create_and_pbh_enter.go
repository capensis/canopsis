package entityservice

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type CreateAndPbhEnterStrategy struct{}

func (s CreateAndPbhEnterStrategy) CanSkip(_ entitycounters.EntityServiceCountersCalcData) bool {
	return false
}

func (s CreateAndPbhEnterStrategy) Calculate(calcData entitycounters.EntityServiceCountersCalcData) entitycounters.EntityCounters {
	if calcData.ServicesToRemove[calcData.Counters.ID] {
		calcData.Counters.Depends--
		calcData.Counters.DecrementState(types.AlarmStateOK, calcData.Inherited)
	} else if calcData.ServicesToAdd[calcData.Counters.ID] {
		calcData.Counters.Depends++
		calcData.Counters.IncrementState(calcData.CurState, calcData.Inherited)

		calcData.Counters.IncrementAlarmCounters(false, calcData.CurActive)

		if !calcData.CurActive {
			calcData.Counters.IncrementPbhCounters(calcData.CurPbhTypeID)
		}
	} else {
		calcData.Counters.DecrementState(types.AlarmStateOK, calcData.Inherited)
		calcData.Counters.IncrementState(calcData.CurState, calcData.Inherited)

		calcData.Counters.IncrementAlarmCounters(false, calcData.CurActive)

		if !calcData.CurActive && calcData.PrevPbhTypeID != calcData.CurPbhTypeID {
			calcData.Counters.IncrementPbhCounters(calcData.CurPbhTypeID)
		}
	}

	return calcData.Counters
}
