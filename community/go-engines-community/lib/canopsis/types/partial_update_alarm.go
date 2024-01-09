package types

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"go.mongodb.org/mongo-driver/bson"
)

func (a *Alarm) PartialUpdateLastEventDate(timestamp datetime.CpsTime) {
	a.Value.LastEventDate = timestamp
	a.AddUpdate("$set", bson.M{
		"v.last_event_date": a.Value.LastEventDate,
	})
}

func (a *Alarm) PartialUpdateResolve(timestamp datetime.CpsTime) error {
	a.Value.Resolved = &timestamp
	a.Value.Duration = int64(timestamp.Sub(a.Value.CreationDate.Time).Seconds())
	a.Value.CurrentStateDuration = int64(timestamp.Sub(a.Value.State.Timestamp.Time).Seconds())

	if a.Value.Snooze != nil {
		snoozeDuration := int64(timestamp.Sub(a.Value.Snooze.Timestamp.Time).Seconds())
		a.Value.SnoozeDuration += snoozeDuration
		a.AddUpdate("$inc", bson.M{"v.snooze_duration": snoozeDuration})
	}
	if !a.Value.PbehaviorInfo.IsActive() {
		enterTimestamp := datetime.CpsTime{}
		for i := len(a.Value.Steps) - 2; i >= 0; i-- {
			if a.Value.Steps[i].Type == AlarmStepPbhEnter {
				enterTimestamp = a.Value.Steps[i].Timestamp
				break
			}
		}

		if !enterTimestamp.IsZero() {
			pbhDuration := int64(timestamp.Sub(enterTimestamp.Time).Seconds())
			a.Value.PbehaviorInactiveDuration += pbhDuration
			a.AddUpdate("$inc", bson.M{"v.pbh_inactive_duration": pbhDuration})
		}
	}

	if (a.Value.Snooze != nil || !a.Value.PbehaviorInfo.IsActive()) && a.Value.InactiveStart != nil {
		inactiveDuration := int64(timestamp.Sub(a.Value.InactiveStart.Time).Seconds())
		a.Value.InactiveDuration += inactiveDuration
		a.AddUpdate("$inc", bson.M{"v.inactive_duration": inactiveDuration})
	}

	a.Value.ActiveDuration = a.Value.Duration - a.Value.InactiveDuration
	a.AddUpdate("$set", bson.M{
		"v.resolved": a.Value.Resolved,

		"v.duration":               a.Value.Duration,
		"v.current_state_duration": a.Value.CurrentStateDuration,
		"v.active_duration":        a.Value.ActiveDuration,
	})
	a.AddUpdate("$unset", bson.M{
		"not_acked_metric_type":      "",
		"not_acked_metric_send_time": "",
		"not_acked_since":            "",
	})

	return nil
}

func (a *Alarm) PartialUpdateAddStepWithStep(newStep AlarmStep) error {
	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}

	a.AddUpdate("$push", bson.M{"v.steps": newStep})

	return nil
}

// AddUpdate adds new mongo updates.
func (a *Alarm) AddUpdate(key string, update bson.M) {
	if a.update == nil {
		a.update = make(bson.M)
	}

	if _, ok := a.update[key]; ok {
		if mergedUpdate, ok := a.update[key].(bson.M); ok {
			for k, v := range update {
				mergedUpdate[k] = v
			}
			a.update[key] = mergedUpdate
		}
	} else {
		a.update[key] = update
	}
}

// GetUpdate returns mongo updates from last update.
func (a *Alarm) GetUpdate() bson.M {
	return a.update
}

// CleanUpdate removes mongo updates. Call it after succeeded update.
func (a *Alarm) CleanUpdate() {
	a.update = nil
	a.childrenUpdate = nil
	a.parentsUpdate = nil
	a.childrenRemove = nil
	a.parentsRemove = nil
}
