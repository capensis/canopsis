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
		calcData.Counters.Depends--
		calcData.Counters.DecrementState(types.AlarmStateOK, calcData.Inherited)
	} else if calcData.ServicesToAdd[calcData.Counters.ID] {
		calcData.Counters.Depends++
		calcData.Counters.IncrementState(calcData.CurState, calcData.Inherited)

		calcData.Counters.IncrementAlarmCounters(false, true)
	} else {
		calcData.Counters.DecrementState(types.AlarmStateOK, calcData.Inherited)
		calcData.Counters.IncrementState(calcData.CurState, calcData.Inherited)

		calcData.Counters.IncrementAlarmCounters(false, true)
	}

	return calcData.Counters
}
