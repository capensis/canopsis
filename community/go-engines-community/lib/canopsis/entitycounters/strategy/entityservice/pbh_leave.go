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
		calcData.Counters.Depends--
		calcData.Counters.DecrementState(calcData.PrevState, calcData.Inherited)

		if calcData.AlarmExists {
			calcData.Counters.DecrementAlarmCounters(calcData.IsAcked, calcData.PrevActive)
		}

		if !calcData.PrevActive {
			calcData.Counters.DecrementPbhCounters(calcData.PrevPbhTypeID)
		}
	} else if calcData.ServicesToAdd[calcData.Counters.ID] {
		calcData.Counters.Depends++
		calcData.Counters.IncrementState(calcData.CurState, calcData.Inherited)

		if calcData.AlarmExists {
			calcData.Counters.IncrementAlarmCounters(calcData.IsAcked, true)
		}
	} else {
		calcData.Counters.DecrementState(calcData.PrevState, calcData.Inherited)
		calcData.Counters.IncrementState(calcData.CurState, calcData.Inherited)

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
