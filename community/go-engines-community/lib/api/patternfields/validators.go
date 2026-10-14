package patternfields

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
)

func ValidateAlarmPattern(p pattern.Alarm, collection string) bool {
	disabledFields, ok := disabledFieldsInAlarmPattern[collection]
	if !ok {
		return false
	}

	return match.ValidateAlarmPattern(p,
		disabledFields,
		onlyAbsoluteTimeCondFieldsInAlarmPattern[collection],
	)
}

func ValidateEntityPattern(p pattern.Entity, collection string) bool {
	disabledFields, ok := disabledFieldsInEntityPattern[collection]
	if !ok {
		return false
	}

	return match.ValidateEntityPattern(p, disabledFields)
}
