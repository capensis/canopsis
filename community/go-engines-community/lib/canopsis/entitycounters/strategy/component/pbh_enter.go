package component

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type PbhEnterStrategy struct{}

func (s PbhEnterStrategy) CanSkip(calcData entitycounters.ComponentCountersCalcData) bool {
	return (calcData.CurActive || !calcData.AlarmExists || calcData.PrevState == types.AlarmStateOK) && !calcData.Info.ComponentStateSettingsToAdd && !calcData.Info.ComponentStateSettingsToRemove
}

func (s PbhEnterStrategy) Calculate(calcData entitycounters.ComponentCountersCalcData) entitycounters.EntityCounters {
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
