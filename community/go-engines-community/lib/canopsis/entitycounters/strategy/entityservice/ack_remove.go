package entityservice

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
)

type AckRemoveStrategy struct{}

func (s AckRemoveStrategy) CanSkip(_ entitycounters.EntityServiceCountersCalcData) bool {
	return false
}

func (s AckRemoveStrategy) Calculate(calcData entitycounters.EntityServiceCountersCalcData) entitycounters.EntityCounters {
	if calcData.ServicesToRemove[calcData.Counters.ID] {
		inheritedMode := entitycounters.InheritedNone
		if calcData.PrevInherited {
			inheritedMode = entitycounters.InheritedWith
		}

		calcData.Counters.Depends--
		calcData.Counters.DecrementState(calcData.CurState, inheritedMode)
		calcData.Counters.DecrementAlarmCounters(true, calcData.CurActive, calcData.CurPbhActive)

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
		calcData.Counters.IncrementState(calcData.CurState, inheritedMode)
		calcData.Counters.IncrementAlarmCounters(false, calcData.CurActive, calcData.CurPbhActive)

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
		if calcData.PrevInherited != calcData.CurInherited {
			if calcData.PrevInherited {
				calcData.Counters.DecrementState(calcData.CurState, entitycounters.InheritedOnly)
			} else {
				calcData.Counters.IncrementState(calcData.CurState, entitycounters.InheritedOnly)
			}
		}

		if calcData.CurActive {
			calcData.Counters.Acknowledged--
			calcData.Counters.NotAcknowledged++
		} else {
			calcData.Counters.AcknowledgedInactive--
		}

		if !calcData.CurPbhActive {
			calcData.Counters.AcknowledgedUnderPbh--
		}
	}

	return calcData.Counters
}
