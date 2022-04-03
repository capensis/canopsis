package common

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"

var StateTitles = map[int]string{
	types.AlarmStateOK:       types.AlarmStateTitleOK,
	types.AlarmStateMinor:    types.AlarmStateTitleMinor,
	types.AlarmStateMajor:    types.AlarmStateTitleMajor,
	types.AlarmStateCritical: types.AlarmStateTitleCritical,
}
