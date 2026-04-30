package component

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type CreateStrategy struct{}

func (s CreateStrategy) CanSkip(_ entitycounters.ComponentCountersCalcData) bool {
	return false
}

func (s CreateStrategy) Calculate(calcData entitycounters.ComponentCountersCalcData) entitycounters.EntityCounters {
	if calcData.Info.ComponentStateSettingsToRemove {
		calcData.Counters.DecrementState(types.AlarmStateOK, entitycounters.InheritedNone)
	} else if calcData.Info.ComponentStateSettingsToAdd {
		calcData.Counters.IncrementState(calcData.CurState, entitycounters.InheritedNone)
	} else {
		calcData.Counters.DecrementState(types.AlarmStateOK, entitycounters.InheritedNone)
		calcData.Counters.IncrementState(calcData.CurState, entitycounters.InheritedNone)
	}

	return calcData.Counters
}
