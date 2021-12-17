package types

import (
	"errors"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
	"go.mongodb.org/mongo-driver/bson"
)

// PartialUpdateAck add ack step to alarm. It saves mongo updates.
func (a *Alarm) PartialUpdateAck(timestamp CpsTime, author, output, userID, role, initiator string) error {
	if a.Value.ACK != nil {
		return nil
	}

	newStep := NewAlarmStep(AlarmStepAck, timestamp, author, output, userID, role, initiator)
	a.Value.ACK = &newStep

	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}

	a.AddUpdate("$set", bson.M{"v.ack": a.Value.ACK})
	a.AddUpdate("$push", bson.M{"v.steps": a.Value.ACK})

	return nil
}

// PartialUpdateUnack deletes ack step from alarm. It saves mongo updates.
func (a *Alarm) PartialUpdateUnack(timestamp CpsTime, author, output, userID, role, initiator string) error {
	if a.Value.ACK == nil {
		return nil
	}

	newStep := NewAlarmStep(AlarmStepAckRemove, timestamp, author, output, userID, role, initiator)
	a.Value.ACK = nil

	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}

	a.AddUpdate("$unset", bson.M{"v.ack": ""})
	a.AddUpdate("$push", bson.M{"v.steps": newStep})

	return nil
}

// PartialUpdateAssocTicket add ticket to alarm. It saves mongo updates.
func (a *Alarm) PartialUpdateAssocTicket(timestamp CpsTime, ticketData map[string]string, author, ticketNumber, userID, role, initiator string) error {
	newStep := NewAlarmStep(AlarmStepAssocTicket, timestamp, author, ticketNumber, userID, role, initiator)
	ticketStep := newStep.NewTicket(ticketNumber, ticketData)
	a.Value.Ticket = &ticketStep

	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}

	a.AddUpdate("$set", bson.M{"v.ticket": a.Value.Ticket})
	a.AddUpdate("$push", bson.M{"v.steps": newStep})

	return nil
}

// PartialUpdateSnooze add snooze step to alarm. It saves mongo updates.
func (a *Alarm) PartialUpdateSnooze(timestamp CpsTime, duration DurationWithUnit, author, output, userID, role, initiator string) error {
	if a.Value.Snooze != nil {
		return nil
	}

	newStep := NewAlarmStep(AlarmStepSnooze, timestamp, author, output, userID, role, initiator)
	if duration.Value == 0 {
		return errt.NewUnknownError(errors.New("no duration for snoozing"))
	}
	newStep.Value = CpsNumber(duration.AddTo(timestamp).Unix())
	a.Value.Snooze = &newStep

	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}

	a.AddUpdate("$set", bson.M{"v.snooze": a.Value.Snooze})
	a.AddUpdate("$push", bson.M{"v.steps": a.Value.Snooze})

	return nil
}

func (a *Alarm) PartialUpdateUnsnooze(timestamp CpsTime) error {
	if a.Value.Snooze == nil {
		return nil
	}

	d := int64(timestamp.Sub(a.Value.Snooze.Timestamp.Time).Seconds())
	a.Value.SnoozeDuration += d
	a.Value.Snooze = nil
	a.AddUpdate("$set", bson.M{"v.snooze": a.Value.Snooze})
	a.AddUpdate("$inc", bson.M{"v.snooze_duration": d})

	return nil
}

func (a *Alarm) PartialUpdatePbhEnter(timestamp CpsTime, pbehaviorInfo PbehaviorInfo, author, output, userID, role, initiator string) error {
	if a.Value.PbehaviorInfo == pbehaviorInfo {
		return nil
	}

	newStep := NewAlarmStep(AlarmStepPbhEnter, timestamp, author, output, userID, role, initiator)
	newStep.PbehaviorCanonicalType = pbehaviorInfo.CanonicalType
	a.Value.PbehaviorInfo = pbehaviorInfo

	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}

	a.AddUpdate("$set", bson.M{"v.pbehavior_info": a.Value.PbehaviorInfo})
	a.AddUpdate("$push", bson.M{"v.steps": newStep})

	return nil
}

func (a *Alarm) PartialUpdatePbhLeave(timestamp CpsTime, author, output, userID, role, initiator string) error {
	if a.Value.Snooze != nil {
		ResolveSnoozeAfterPbhLeave(timestamp, a)
		a.AddUpdate("$set", bson.M{"v.snooze": a.Value.Snooze})
	}

	newStep := NewAlarmStep(AlarmStepPbhLeave, timestamp, author, output, userID, role, initiator)
	newStep.PbehaviorCanonicalType = a.Value.PbehaviorInfo.CanonicalType
	a.Value.PbehaviorInfo = PbehaviorInfo{}

	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}

	a.AddUpdate("$unset", bson.M{"v.pbehavior_info": ""})
	a.AddUpdate("$push", bson.M{"v.steps": newStep})

	if newStep.PbehaviorCanonicalType != "active" {
		enterTimestamp := CpsTime{}
		for i := len(a.Value.Steps) - 2; i >= 0; i-- {
			if a.Value.Steps[i].Type == AlarmStepPbhEnter {
				enterTimestamp = a.Value.Steps[i].Timestamp
			}
		}

		d := int64(timestamp.Sub(enterTimestamp.Time).Seconds())
		a.Value.PbehaviorInactiveDuration += d
		a.AddUpdate("$inc", bson.M{"v.pbh_inactive_duration": d})
	}

	return nil
}

func (a *Alarm) PartialUpdatePbhLeaveAndEnter(timestamp CpsTime, pbehaviorInfo PbehaviorInfo, author, output, userID, role, initiator string) error {
	if a.Value.PbehaviorInfo == pbehaviorInfo {
		return nil
	}

	if pbehaviorInfo.CanonicalType == "active" && a.Value.Snooze != nil {
		ResolveSnoozeAfterPbhLeave(timestamp, a)
		a.AddUpdate("$set", bson.M{"v.snooze": a.Value.Snooze})
	}

	leaveOutput := fmt.Sprintf(
		"Pbehavior %s. Type: %s. Reason: %s.",
		a.Value.PbehaviorInfo.Name,
		a.Value.PbehaviorInfo.TypeName,
		a.Value.PbehaviorInfo.Reason,
	)

	leaveStep := NewAlarmStep(AlarmStepPbhLeave, timestamp, author, leaveOutput, userID, role, initiator)
	leaveStep.PbehaviorCanonicalType = a.Value.PbehaviorInfo.CanonicalType

	err := a.Value.Steps.Add(leaveStep)
	if err != nil {
		return err
	}

	enterStep := NewAlarmStep(AlarmStepPbhEnter, timestamp, author, output, userID, role, initiator)
	enterStep.PbehaviorCanonicalType = pbehaviorInfo.CanonicalType

	a.Value.PbehaviorInfo = pbehaviorInfo
	err = a.Value.Steps.Add(enterStep)
	if err != nil {
		return err
	}

	a.AddUpdate("$set", bson.M{"v.pbehavior_info": a.Value.PbehaviorInfo})
	a.AddUpdate("$push", bson.M{"v.steps": bson.M{"$each": bson.A{leaveStep, enterStep}}})

	if leaveStep.PbehaviorCanonicalType != "active" {
		enterTimestamp := CpsTime{}
		for i := len(a.Value.Steps) - 3; i >= 0; i-- {
			if a.Value.Steps[i].Type == AlarmStepPbhEnter {
				enterTimestamp = a.Value.Steps[i].Timestamp
			}
		}

		d := int64(timestamp.Sub(enterTimestamp.Time).Seconds())
		a.Value.PbehaviorInactiveDuration += d
		a.AddUpdate("$inc", bson.M{"v.pbh_inactive_duration": d})
	}

	return nil
}

// PartialUpdateDeclareTicket add ticket to alarm. It saves mongo updates.
func (a *Alarm) PartialUpdateDeclareTicket(timestamp CpsTime, author, output, ticketNumber string, data map[string]string, userID, role, initiator string) error {
	newStep := NewAlarmStep(AlarmStepDeclareTicket, timestamp, author, output, userID, role, initiator)
	ticketStep := newStep.NewTicket(ticketNumber, data)
	a.Value.Ticket = &ticketStep

	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}

	a.AddUpdate("$set", bson.M{"v.ticket": a.Value.Ticket})
	a.AddUpdate("$push", bson.M{"v.steps": newStep})

	return nil
}

func (a *Alarm) PartialUpdateLastEventDate(timestamp CpsTime) {
	a.Value.LastEventDate = timestamp
	a.AddUpdate("$set", bson.M{
		"v.last_event_date": a.Value.LastEventDate,
	})
}

func (a *Alarm) PartialUpdateEventsCount() {
	a.Value.EventsCount++
	a.AddUpdate("$set", bson.M{
		"v.events_count": a.Value.EventsCount,
	})
}

func (a *Alarm) PartialUpdateActivate(timestamp CpsTime) error {
	if a.IsActivated() {
		return nil
	}

	a.Value.ActivationDate = &timestamp
	a.AddUpdate("$set", bson.M{"v.activation_date": a.Value.ActivationDate})

	return nil
}

func (a *Alarm) PartialUpdateResolve(timestamp CpsTime) error {
	a.Value.Resolved = &timestamp
	a.AddUpdate("$set", bson.M{"v.resolved": a.Value.Resolved})

	return nil
}

func (a *Alarm) PartialUpdateDone(timestamp CpsTime, author, output, userID, role, initiator string) error {
	if a.Value.Done != nil {
		return nil
	}

	newStep := NewAlarmStep(AlarmStepDone, timestamp, author, output, userID, role, initiator)
	a.Value.Done = &newStep

	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}

	a.AddUpdate("$set", bson.M{"v.done": a.Value.Done})
	a.AddUpdate("$push", bson.M{"v.steps": a.Value.Done})

	return nil
}

func (a *Alarm) PartialUpdateComment(timestamp CpsTime, author, output, userID, role, initiator string) error {
	newStep := NewAlarmStep(AlarmStepComment, timestamp, author, output, userID, role, initiator)
	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}

	a.AddUpdate("$push", bson.M{"v.steps": newStep})

	return nil
}

func (a *Alarm) PartialUpdateAddInstructionStep(stepType string, timestamp CpsTime,
	execution, author, msg, userID, role, initiator string) error {
	newStep := NewAlarmStep(stepType, timestamp, author, msg, userID, role, initiator)
	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}

	newStep.Execution = execution

	a.AddUpdate("$push", bson.M{"v.steps": newStep})

	return nil
}

func (a *Alarm) PartialUpdateCropSteps() {
	if a.CropSteps() {
		a.AddUpdate("$set", bson.M{"v.steps": a.Value.Steps})
	}
}

func (a *Alarm) PartialUpdateAddStep(stepType string, timestamp CpsTime, author, msg, userID, role, initiator string) error {
	newStep := NewAlarmStep(stepType, timestamp, author, msg, userID, role, initiator)
	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}

	a.AddUpdate("$push", bson.M{"v.steps": newStep})

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
		mergedUpdate := a.update[key].(bson.M)
		for k, v := range update {
			mergedUpdate[k] = v
		}
		a.update[key] = mergedUpdate
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

func ResolveSnoozeAfterPbhLeave(timestamp CpsTime, alarm *Alarm) {
	steps := alarm.Value.Steps

	var snoozeDuration int64
	var snoozeElapsed int64
	var lastEnterTime int64

	if alarm.Value.Snooze != nil && alarm.Value.Snooze.Initiator != InitiatorUser {
	Loop:
		for i := len(steps) - 1; i >= 0; i-- {
			step := steps[i]
			switch step.Type {
			case AlarmStepSnooze:
				snoozeElapsed += lastEnterTime - step.Timestamp.Unix()
				snoozeDuration = int64(step.Value) - step.Timestamp.Unix()

				break Loop
			case AlarmStepPbhEnter:
				if step.PbehaviorCanonicalType != "active" {
					lastEnterTime = step.Timestamp.Unix()
				}
			case AlarmStepPbhLeave:
				if step.PbehaviorCanonicalType != "active" {
					snoozeElapsed += lastEnterTime - step.Timestamp.Unix()
				}
			}
		}

		alarm.Value.Snooze.Value = CpsNumber(timestamp.Unix() + snoozeDuration - snoozeElapsed)
	}
}
