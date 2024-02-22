package component

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type ResolveStrategy struct{}

func (s ResolveStrategy) CanSkip(calcData entitycounters.ComponentCountersCalcData) bool {
	return (!calcData.CurActive || calcData.PrevState == types.AlarmStateOK) && !calcData.Info.ComponentStateSettingsToAdd && !calcData.Info.ComponentStateSettingsToRemove
}

func (s ResolveStrategy) Calculate(calcData entitycounters.ComponentCountersCalcData) entitycounters.EntityCounters {
	if calcData.Info.ComponentStateSettingsToRemove {
		calcData.Counters.DecrementState(calcData.PrevState, false)
	} else if calcData.Info.ComponentStateSettingsToAdd {
		calcData.Counters.IncrementState(types.AlarmStateOK, false)
	} else {
		calcData.Counters.DecrementState(calcData.PrevState, false)
		calcData.Counters.IncrementState(types.AlarmStateOK, false)
	}

	return calcData.Counters
}
