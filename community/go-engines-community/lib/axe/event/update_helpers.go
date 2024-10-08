package event

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
)

func updateInactiveStart(
	ts types.CpsTime,
	withSnoozeCond bool,
	withPbhCond bool,
	withAutoInstructionCond bool,
) bson.M {
	conds := make([]bson.M, 0)
	if withSnoozeCond {
		conds = append(conds, bson.M{"$eq": bson.A{
			bson.M{"$cond": bson.M{
				"if":   "$v.snooze",
				"then": "$v.snooze",
				"else": nil,
			}},
			nil,
		}})
	}

	if withPbhCond {
		conds = append(conds, bson.M{"$in": bson.A{
			bson.M{"$cond": bson.M{
				"if":   "$v.pbehavior_info",
				"then": "$v.pbehavior_info.canonical_type",
				"else": nil,
			}},
			bson.A{nil, "", pbehavior.TypeActive},
		}})
	}

	if withAutoInstructionCond {
		conds = append(conds, bson.M{"$ne": bson.A{"$auto_instruction_in_progress", true}})
	}

	return bson.M{"$cond": bson.M{
		"if":   bson.M{"$and": conds},
		"then": nil,
		"else": ts,
	}}
}

func RemoveMetaAlarmState(
	ctx context.Context,
	metaAlarm types.Alarm,
	rule correlation.Rule,
	metaAlarmStatesService correlation.MetaAlarmStateService,
) error {
	if rule.IsManual() {
		return nil
	}

	stateID := rule.GetStateID(metaAlarm.Value.MetaValuePath)
	metaAlarmState, err := metaAlarmStatesService.GetMetaAlarmState(ctx, stateID)
	if err != nil {
		return fmt.Errorf("cannot get meta alarm state: %w", err)
	}

	if metaAlarmState.ID == "" {
		return nil
	}

	_, err = metaAlarmStatesService.ArchiveState(ctx, metaAlarmState)
	if err != nil {
		return fmt.Errorf("cannot archive meta alarm state: %w", err)
	}

	_, err = metaAlarmStatesService.DeleteState(ctx, stateID)
	if err != nil {
		return fmt.Errorf("cannot delete meta alarm state: %w", err)
	}

	return nil
}
