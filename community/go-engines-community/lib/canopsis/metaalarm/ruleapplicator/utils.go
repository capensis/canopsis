package ruleapplicator

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/go-engines/lib/errt"
)

/**
@todo: Basically AlreadyBelongsToMetaalarm function should be used in rule_service, but if include an alarmAdapter there, it causes import cycle error alarm->metaalarm->alarm.
@todo: AlreadyBelongsToMetaalarm function is used in every applicator, which uses time intervals. This is a temporary solution, until we refactor our alarm/metaalarm packages.
 */
func AlreadyBelongsToMetaalarm(alarmAdapter alarm.Adapter, alarmEID string, ruleId, valuePath string) (bool, error) {
	metaAlarm, err := alarmAdapter.GetOpenedMetaAlarm(ruleId, valuePath)
	switch err.(type) {
	case errt.NotFound:
		return false, nil
	case nil:
		if metaAlarm.HasChildByEID(alarmEID) {
			return true, nil
		}
	}

	return false, err
}
