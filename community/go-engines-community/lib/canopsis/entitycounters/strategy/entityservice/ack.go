package entityservice

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
)

type AckStrategy struct{}

func (s AckStrategy) CanSkip(_ entitycounters.EntityServiceCountersCalcData) bool {
	return false
}

func (s AckStrategy) Calculate(calcData entitycounters.EntityServiceCountersCalcData) entitycounters.EntityCounters {
	if calcData.ServicesToRemove[calcData.Counters.ID] {
		calcData.Counters.Depends--
		calcData.Counters.DecrementState(calcData.CurState, calcData.Inherited)
		calcData.Counters.DecrementAlarmCounters(false, calcData.CurActive)

		if !calcData.CurActive {
			calcData.Counters.DecrementPbhCounters(calcData.CurPbhTypeID)
		}
	} else if calcData.ServicesToAdd[calcData.Counters.ID] {
		calcData.Counters.Depends++
		calcData.Counters.IncrementState(calcData.CurState, calcData.Inherited)
		calcData.Counters.IncrementAlarmCounters(true, calcData.CurActive)

		if !calcData.CurActive {
			calcData.Counters.IncrementPbhCounters(calcData.CurPbhTypeID)
		}
	} else {
		if !calcData.CurActive {
			calcData.Counters.AcknowledgedUnderPbh++
		} else {
			calcData.Counters.Acknowledged++
			calcData.Counters.NotAcknowledged--
		}
	}

	return calcData.Counters
}
