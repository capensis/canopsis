package component

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
)

type ChangeStateStrategy struct{}

func (s ChangeStateStrategy) CanSkip(calcData entitycounters.ComponentCountersCalcData) bool {
	return !calcData.CurActive && !calcData.Info.ComponentStateSettingsToAdd && !calcData.Info.ComponentStateSettingsToRemove
}

func (s ChangeStateStrategy) Calculate(calcData entitycounters.ComponentCountersCalcData) entitycounters.EntityCounters {
	if calcData.Info.ComponentStateSettingsToRemove {
		calcData.Counters.DecrementState(calcData.PrevState, false)
	} else if calcData.Info.ComponentStateSettingsToAdd {
		calcData.Counters.IncrementState(calcData.CurState, false)
	} else {
		calcData.Counters.DecrementState(calcData.PrevState, false)
		calcData.Counters.IncrementState(calcData.CurState, false)
	}

	return calcData.Counters
}
