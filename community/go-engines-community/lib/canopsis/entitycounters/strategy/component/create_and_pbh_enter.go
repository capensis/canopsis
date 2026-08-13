package component

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type CreateAndPbhEnterStrategy struct{}

func (s CreateAndPbhEnterStrategy) CanSkip(calcData entitycounters.ComponentCountersCalcData) bool {
	return !calcData.CurActive && !calcData.Info.ComponentStateSettingsToAdd && !calcData.Info.ComponentStateSettingsToRemove
}

func (s CreateAndPbhEnterStrategy) Calculate(calcData entitycounters.ComponentCountersCalcData) entitycounters.EntityCounters {
	if calcData.Info.ComponentStateSettingsToRemove {
		calcData.Counters.DecrementState(types.AlarmStateOK, entitycounters.InheritedNone)
	} else if calcData.Info.ComponentStateSettingsToAdd {
		calcData.Counters.IncrementState(calcData.CurState, entitycounters.InheritedNone)
	} else {
		// do not check for active because it's already checked above
		calcData.Counters.DecrementState(types.AlarmStateOK, entitycounters.InheritedNone)
		calcData.Counters.IncrementState(calcData.CurState, entitycounters.InheritedNone)
	}

	return calcData.Counters
}
