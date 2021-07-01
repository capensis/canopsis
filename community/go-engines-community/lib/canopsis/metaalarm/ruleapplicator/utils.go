package ruleapplicator

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
)

/**
@todo: Basically AlreadyBelongsToMetaalarm function should be used in rule_service, but if include an alarmAdapter there, it causes import cycle error alarm->metaalarm->alarm.
@todo: AlreadyBelongsToMetaalarm function is used in every applicator, which uses time intervals. This is a temporary solution, until we refactor our alarm/metaalarm packages.
*/
func AlreadyBelongsToMetaalarm(alarmAdapter alarm.Adapter, alarmEID string, ruleId, valuePath string) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	metaAlarm, err := alarmAdapter.GetOpenedMetaAlarm(ctx, ruleId, valuePath)
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
