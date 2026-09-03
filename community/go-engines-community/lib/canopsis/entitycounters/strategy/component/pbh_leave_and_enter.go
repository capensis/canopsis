package component

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type PbhLeaveAndEnterStrategy struct{}

func (s PbhLeaveAndEnterStrategy) CanSkip(calcData entitycounters.ComponentCountersCalcData) bool {
	return (calcData.PrevActive == calcData.CurActive || !calcData.AlarmExists || calcData.PrevState == types.AlarmStateOK && calcData.CurState == types.AlarmStateOK) && !calcData.Info.ComponentStateSettingsToAdd && !calcData.Info.ComponentStateSettingsToRemove
}

func (s PbhLeaveAndEnterStrategy) Calculate(calcData entitycounters.ComponentCountersCalcData) entitycounters.EntityCounters {
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
