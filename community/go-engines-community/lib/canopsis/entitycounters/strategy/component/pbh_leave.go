package component

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type PbhLeaveStrategy struct{}

func (s PbhLeaveStrategy) CanSkip(calcData entitycounters.ComponentCountersCalcData) bool {
	return (calcData.PrevActive || !calcData.AlarmExists || calcData.CurState == types.AlarmStateOK) && !calcData.Info.ComponentStateSettingsToAdd && !calcData.Info.ComponentStateSettingsToRemove
}

func (s PbhLeaveStrategy) Calculate(calcData entitycounters.ComponentCountersCalcData) entitycounters.EntityCounters {
	if calcData.Info.ComponentStateSettingsToRemove {
		calcData.Counters.DecrementState(calcData.PrevState, entitycounters.InheritedNone)
	} else if calcData.Info.ComponentStateSettingsToAdd {
		calcData.Counters.IncrementState(calcData.CurState, entitycounters.InheritedNone)
	} else {
		calcData.Counters.DecrementState(calcData.PrevState, entitycounters.InheritedNone)
		calcData.Counters.IncrementState(calcData.CurState, entitycounters.InheritedNone)
	}

	return calcData.Counters
}
