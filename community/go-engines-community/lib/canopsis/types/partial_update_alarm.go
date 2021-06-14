package types

import (
	"errors"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
	"go.mongodb.org/mongo-driver/bson"
)

// PartialUpdateAck add ack step to alarm. It saves mongo updates.
func (a *Alarm) PartialUpdateAck(timestamp CpsTime, author, output, role, initiator string) error {
	if a.Value.ACK != nil {
		return nil
	}

	newStep := NewAlarmStep(AlarmStepAck, timestamp, author, output, role, initiator)
	a.Value.ACK = &newStep

	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}

	a.addUpdate("$set", bson.M{"v.ack": a.Value.ACK})
	a.addUpdate("$push", bson.M{"v.steps": a.Value.ACK})

	return nil
}

// PartialUpdateUnack deletes ack step from alarm. It saves mongo updates.
func (a *Alarm) PartialUpdateUnack(timestamp CpsTime, author, output, role, initiator string) error {
	if a.Value.ACK == nil {
		return nil
	}

	newStep := NewAlarmStep(AlarmStepAckRemove, timestamp, author, output, role, initiator)
	a.Value.ACK = nil

	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}

	a.addUpdate("$unset", bson.M{"v.ack": ""})
	a.addUpdate("$push", bson.M{"v.steps": newStep})

	return nil
}

// PartialUpdateCancel add canceled and status change steps to alarm. It saves mongo updates.
func (a *Alarm) PartialUpdateCancel(timestamp CpsTime, author, output, role,
	initiator string, alarmConfig config.AlarmConfig) error {
	if a.Value.Canceled != nil {
		return nil
	}

	newStepCancel := NewAlarmStep(AlarmStepCancel, timestamp, author, output, role, initiator)
	a.Value.Canceled = &newStepCancel

	if err := a.Value.Steps.Add(newStepCancel); err != nil {
		return err
	}

	currentStatus := a.CurrentStatus(alarmConfig)
	newStatus := a.ComputeStatus(alarmConfig)

	if newStatus == currentStatus && a.Value.Status != nil {
		a.addUpdate("$set", bson.M{"v.canceled": a.Value.Canceled})
		a.addUpdate("$push", bson.M{"v.steps": a.Value.Canceled})
		return nil
	}

	newStepStatus := NewAlarmStep(AlarmStepStatusIncrease, timestamp, a.Value.Connector+"."+a.Value.ConnectorName, output, role, initiator)
	newStepStatus.Value = newStatus
	if a.Value.Status != nil && newStepStatus.Value < a.Value.Status.Value {
		newStepStatus.Type = AlarmStepStatusDecrease
	}
	a.Value.Status = &newStepStatus
	if err := a.Value.Steps.Add(newStepStatus); err != nil {
		return err
	}

	a.Value.StateChangesSinceStatusUpdate = 0
	a.Value.LastUpdateDate = timestamp

	a.addUpdate("$set", bson.M{
		"v.canceled":                          a.Value.Canceled,
		"v.status":                            a.Value.Status,
		"v.state_changes_since_status_update": a.Value.StateChangesSinceStatusUpdate,
		"v.last_update_date":                  a.Value.LastUpdateDate,
	})
	a.addUpdate("$push", bson.M{"v.steps": bson.M{"$each": bson.A{a.Value.Canceled, a.Value.Status}}})

	return nil
}

func (a *Alarm) PartialUpdateUncancel(timestamp CpsTime, author, output, role,
	initiator string, alarmConfig config.AlarmConfig) error {
	if a.Value.Canceled == nil {
		return nil
	}

	newStep := NewAlarmStep(AlarmStepUncancel, timestamp, author, output, role, initiator)
	a.Value.Canceled = nil

	if err := a.Value.Steps.Add(newStep); err != nil {
		return err
	}

	a.addUpdate("$set", bson.M{"v.canceled": a.Value.Canceled})
	a.addUpdate("$push", bson.M{"v.steps": newStep})

	currentStatus := a.CurrentStatus(alarmConfig)
	newStatus := a.ComputeStatus(alarmConfig)

	if newStatus == currentStatus && a.Value.Status != nil {
		a.addUpdate("$set", bson.M{"v.canceled": a.Value.Canceled})
		a.addUpdate("$push", bson.M{"v.steps": newStep})
		return nil
	}

	newStepStatus := NewAlarmStep(AlarmStepStatusIncrease, timestamp, a.Value.Connector+"."+a.Value.ConnectorName, output, role, initiator)
	newStepStatus.Value = newStatus
	if a.Value.Status != nil && newStepStatus.Value < a.Value.Status.Value {
		newStepStatus.Type = AlarmStepStatusDecrease
	}
	a.Value.Status = &newStepStatus
	if err := a.Value.Steps.Add(newStepStatus); err != nil {
		return err
	}

	a.Value.StateChangesSinceStatusUpdate = 0
	a.Value.LastUpdateDate = timestamp

	a.addUpdate("$set", bson.M{
		"v.canceled":                          a.Value.Canceled,
		"v.status":                            a.Value.Status,
		"v.state_changes_since_status_update": a.Value.StateChangesSinceStatusUpdate,
		"v.last_update_date":                  a.Value.LastUpdateDate,
	})
	a.addUpdate("$push", bson.M{"v.steps": bson.M{"$each": bson.A{newStep, a.Value.Status}}})

	return nil
}

// PartialUpdateChangeState add state change step to alarm. It saves mongo updates.
func (a *Alarm) PartialUpdateChangeState(state CpsNumber, timestamp CpsTime, author, output, role, initiator string) error {
	if a.Value.State != nil && a.Value.State.Value == state {
		return nil
	}

	newStep := NewAlarmStep(AlarmStepChangeState, timestamp, author, output, role, initiator)
	newStep.Value = state
	a.Value.State = &newStep

	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}

	a.addUpdate("$set", bson.M{"v.state": a.Value.State})
	a.addUpdate("$push", bson.M{"v.steps": a.Value.State})

	return nil
}

// PartialUpdateNoEvents add state step to alarm. It saves mongo updates.
func (a *Alarm) PartialUpdateNoEvents(state CpsNumber, timestamp CpsTime, author, output, role, initiator string,
	alarmConfig config.AlarmConfig) error {
	currentState := a.CurrentState()
	stateUpdated := false
	if currentState != state {
		// Create new Step to keep track of the alarm history
		newStep := NewAlarmStep(AlarmStepStateIncrease, timestamp, author, output, role, initiator)
		newStep.Value = state

		if state < currentState {
			newStep.Type = AlarmStepStateDecrease
		}

		a.Value.State = &newStep
		err := a.Value.Steps.Add(newStep)
		if err != nil {
			return err
		}

		stateUpdated = true
	}

	currentStatus := a.CurrentStatus(alarmConfig)
	newStatus := CpsNumber(AlarmStatusNoEvents)
	if state == AlarmStateOK {
		newStatus = a.ComputeStatus(alarmConfig)
	}

	if newStatus == currentStatus && a.Value.Status != nil {
		if stateUpdated {
			a.addUpdate("$set", bson.M{"v.state": a.Value.State})
			a.addUpdate("$push", bson.M{"v.steps": a.Value.State})
		}
		return nil
	}

	// Create new Step to keep track of the alarm history
	newStepStatus := NewAlarmStep(AlarmStepStatusIncrease, timestamp, author, output, role, initiator)
	newStepStatus.Value = newStatus

	if newStatus < currentStatus {
		newStepStatus.Type = AlarmStepStatusDecrease
	}

	a.Value.Status = &newStepStatus
	err := a.Value.Steps.Add(newStepStatus)
	if err != nil {
		return err
	}

	a.Value.StateChangesSinceStatusUpdate = 0
	a.Value.LastUpdateDate = timestamp

	set := bson.M{
		"v.status":                            a.Value.Status,
		"v.state_changes_since_status_update": a.Value.StateChangesSinceStatusUpdate,
		"v.last_update_date":                  a.Value.LastUpdateDate,
	}
	newSteps := bson.A{}
	if stateUpdated {
		set["v.state"] = a.Value.State
		newSteps = append(newSteps, a.Value.State)
	}
	newSteps = append(newSteps, a.Value.Status)
	a.addUpdate("$set", set)
	a.addUpdate("$push", bson.M{"v.steps": bson.M{"$each": newSteps}})

	return nil
}

// PartialUpdateAssocTicket add ticket to alarm. It saves mongo updates.
func (a *Alarm) PartialUpdateAssocTicket(timestamp CpsTime, author, ticketNumber, role, initiator string) error {
	newStep := NewAlarmStep(AlarmStepAssocTicket, timestamp, author, ticketNumber, role, initiator)
	ticketStep := newStep.NewTicket(ticketNumber, nil)
	a.Value.Ticket = &ticketStep

	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}

	a.addUpdate("$set", bson.M{"v.ticket": a.Value.Ticket})
	a.addUpdate("$push", bson.M{"v.steps": newStep})

	return nil
}

// PartialUpdateSnooze add snooze step to alarm. It saves mongo updates.
func (a *Alarm) PartialUpdateSnooze(timestamp CpsTime, duration CpsNumber, author, output, role, initiator string) error {
	if a.Value.Snooze != nil {
		return nil
	}

	newStep := NewAlarmStep(AlarmStepSnooze, timestamp, author, output, role, initiator)
	if duration == 0 {
		return errt.NewUnknownError(errors.New("no duration for snoozing"))
	}
	newStep.Value = CpsNumber(timestamp.Time.Unix()) + duration
	a.Value.Snooze = &newStep

	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}

	a.addUpdate("$set", bson.M{"v.snooze": a.Value.Snooze})
	a.addUpdate("$push", bson.M{"v.steps": a.Value.Snooze})

	return nil
}

func (a *Alarm) PartialUpdateUnsnooze() error {
	if a.Value.Snooze == nil {
		return nil
	}

	a.Value.Snooze = nil
	a.addUpdate("$set", bson.M{"v.snooze": a.Value.Snooze})

	return nil
}

func (a *Alarm) PartialUpdatePbhEnter(timestamp CpsTime, pbehaviorInfo PbehaviorInfo, author, output, role, initiator string) error {
	if a.Value.PbehaviorInfo == pbehaviorInfo {
		return nil
	}

	newStep := NewAlarmStep(AlarmStepPbhEnter, timestamp, author, output, role, initiator)
	newStep.PbehaviorCanonicalType = pbehaviorInfo.CanonicalType
	a.Value.PbehaviorInfo = pbehaviorInfo

	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}

	a.addUpdate("$set", bson.M{"v.pbehavior_info": a.Value.PbehaviorInfo})
	a.addUpdate("$push", bson.M{"v.steps": newStep})

	return nil
}

func (a *Alarm) PartialUpdatePbhLeave(timestamp CpsTime, author, output, role, initiator string) error {
	if a.Value.Snooze != nil {
		ResolveSnoozeAfterPbhLeave(timestamp, a)
		a.addUpdate("$set", bson.M{"v.snooze": a.Value.Snooze})
	}

	newStep := NewAlarmStep(AlarmStepPbhLeave, timestamp, author, output, role, initiator)
	newStep.PbehaviorCanonicalType = a.Value.PbehaviorInfo.CanonicalType
	a.Value.PbehaviorInfo = PbehaviorInfo{}

	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}

	a.addUpdate("$unset", bson.M{"v.pbehavior_info": ""})
	a.addUpdate("$push", bson.M{"v.steps": newStep})

	return nil
}

func (a *Alarm) PartialUpdatePbhLeaveAndEnter(timestamp CpsTime, pbehaviorInfo PbehaviorInfo, author, output, role, initiator string) error {
	if a.Value.PbehaviorInfo == pbehaviorInfo {
		return nil
	}

	if pbehaviorInfo.CanonicalType == "active" && a.Value.Snooze != nil {
		ResolveSnoozeAfterPbhLeave(timestamp, a)
		a.addUpdate("$set", bson.M{"v.snooze": a.Value.Snooze})
	}

	leaveOutput := fmt.Sprintf(
		"Pbehavior %s. Type: %s. Reason: %s",
		a.Value.PbehaviorInfo.Name,
		a.Value.PbehaviorInfo.TypeName,
		a.Value.PbehaviorInfo.Reason,
	)

	leaveStep := NewAlarmStep(AlarmStepPbhLeave, timestamp, author, leaveOutput, role, initiator)
	leaveStep.PbehaviorCanonicalType = a.Value.PbehaviorInfo.CanonicalType

	err := a.Value.Steps.Add(leaveStep)
	if err != nil {
		return err
	}

	enterStep := NewAlarmStep(AlarmStepPbhEnter, timestamp, author, output, role, initiator)
	enterStep.PbehaviorCanonicalType = pbehaviorInfo.CanonicalType

	a.Value.PbehaviorInfo = pbehaviorInfo
	err = a.Value.Steps.Add(enterStep)
	if err != nil {
		return err
	}

	a.addUpdate("$set", bson.M{"v.pbehavior_info": a.Value.PbehaviorInfo})
	a.addUpdate("$push", bson.M{"v.steps": bson.M{"$each": bson.A{leaveStep, enterStep}}})

	return nil
}

// PartialUpdateDeclareTicket add ticket to alarm. It saves mongo updates.
func (a *Alarm) PartialUpdateDeclareTicket(timestamp CpsTime, author, output, ticketNumber string, data map[string]string, role, initiator string) error {
	newStep := NewAlarmStep(AlarmStepDeclareTicket, timestamp, author, output, role, initiator)
	ticketStep := newStep.NewTicket(ticketNumber, data)
	a.Value.Ticket = &ticketStep

	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}

	a.addUpdate("$set", bson.M{"v.ticket": a.Value.Ticket})
	a.addUpdate("$push", bson.M{"v.steps": newStep})

	return nil
}

func (a *Alarm) PartialUpdateState(timestamp CpsTime, state CpsNumber, output string,
	alarmConfig config.AlarmConfig) error {
	currentState := a.CurrentState()

	if state != currentState {
		// Event is an OK, so the alarm should be resolved anyway
		if a.IsStateLocked() && state != AlarmStateOK {
			return nil
		}

		// Create new Step to keep track of the alarm history
		newStep := NewAlarmStep(AlarmStepStateIncrease, timestamp, a.Value.Connector+"."+a.Value.ConnectorName, output, "", "")
		newStep.Value = state

		if state < currentState {
			newStep.Type = AlarmStepStateDecrease
		}

		a.Value.State = &newStep
		err := a.Value.Steps.Add(newStep)
		if err != nil {
			return err
		}

		a.Value.TotalStateChanges++
		a.Value.LastUpdateDate = timestamp
	}

	currentStatus := a.CurrentStatus(alarmConfig)
	newStatus := a.ComputeStatus(alarmConfig)
	if state == currentState && currentStatus == newStatus {
		return nil
	}

	if newStatus == currentStatus && a.Value.Status != nil {
		if state != currentState {
			a.Value.StateChangesSinceStatusUpdate++

			a.addUpdate("$set", bson.M{
				"v.state":                             a.Value.State,
				"v.state_changes_since_status_update": a.Value.StateChangesSinceStatusUpdate,
				"v.total_state_changes":               a.Value.TotalStateChanges,
				"v.last_update_date":                  a.Value.LastUpdateDate,
			})
			a.addUpdate("$push", bson.M{"v.steps": a.Value.State})
		}
		return nil
	}

	// Create new Step to keep track of the alarm history
	newStepStatus := NewAlarmStep(AlarmStepStatusIncrease, timestamp, a.Value.Connector+"."+a.Value.ConnectorName, output, "", "")
	newStepStatus.Value = newStatus

	if newStatus < currentStatus {
		newStepStatus.Type = AlarmStepStatusDecrease
	}

	a.Value.Status = &newStepStatus
	err := a.Value.Steps.Add(newStepStatus)
	if err != nil {
		return err
	}

	a.Value.StateChangesSinceStatusUpdate = 0
	a.Value.LastUpdateDate = timestamp

	set := bson.M{
		"v.status":                            a.Value.Status,
		"v.state_changes_since_status_update": a.Value.StateChangesSinceStatusUpdate,
		"v.last_update_date":                  a.Value.LastUpdateDate,
	}
	newSteps := bson.A{}
	if state != currentState {
		set["v.total_state_changes"] = a.Value.TotalStateChanges
		set["v.state"] = a.Value.State
		newSteps = append(newSteps, a.Value.State)
	}
	newSteps = append(newSteps, a.Value.Status)

	a.addUpdate("$set", set)
	a.addUpdate("$push", bson.M{"v.steps": bson.M{"$each": newSteps}})

	return nil
}

func (a *Alarm) PartialUpdateStatus(timestamp CpsTime, output string, alarmConfig config.AlarmConfig) error {
	currentStatus := a.CurrentStatus(alarmConfig)
	newStatus := a.ComputeStatus(alarmConfig)

	if newStatus == currentStatus && a.Value.Status != nil {
		return nil
	}

	// Create new Step to keep track of the alarm history
	newStep := NewAlarmStep(AlarmStepStatusIncrease, timestamp, a.Value.Connector+"."+a.Value.ConnectorName, output, "", "")
	newStep.Value = newStatus

	if newStatus < currentStatus {
		newStep.Type = AlarmStepStatusDecrease
	}

	a.Value.Status = &newStep
	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}

	a.Value.StateChangesSinceStatusUpdate = 0
	a.Value.LastUpdateDate = timestamp

	a.addUpdate("$set", bson.M{
		"v.status":                            a.Value.Status,
		"v.state_changes_since_status_update": a.Value.StateChangesSinceStatusUpdate,
		"v.last_update_date":                  a.Value.LastUpdateDate,
	})
	a.addUpdate("$push", bson.M{"v.steps": a.Value.Status})

	return nil
}

func (a *Alarm) PartialUpdateLastEventDate(timestamp CpsTime) {
	a.Value.LastEventDate = timestamp
	a.addUpdate("$set", bson.M{
		"v.last_event_date": a.Value.LastEventDate,
	})
}

func (a *Alarm) PartialUpdateEventsCount() {
	a.Value.EventsCount++
	a.addUpdate("$set", bson.M{
		"v.events_count": a.Value.EventsCount,
	})
}

func (a *Alarm) PartialUpdateActivate(timestamp CpsTime) error {
	if a.IsActivated() {
		return nil
	}

	a.Value.ActivationDate = &timestamp
	a.addUpdate("$set", bson.M{"v.activation_date": a.Value.ActivationDate})

	return nil
}

func (a *Alarm) PartialUpdateResolve(timestamp CpsTime) error {
	a.Value.Resolved = &timestamp
	a.addUpdate("$set", bson.M{"v.resolved": a.Value.Resolved})

	return nil
}

func (a *Alarm) PartialUpdateDone(timestamp CpsTime, author, output, role, initiator string) error {
	if a.Value.Done != nil {
		return nil
	}

	newStep := NewAlarmStep(AlarmStepDone, timestamp, author, output, role, initiator)
	a.Value.Done = &newStep

	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}

	a.addUpdate("$set", bson.M{"v.done": a.Value.Done})
	a.addUpdate("$push", bson.M{"v.steps": a.Value.Done})

	return nil
}

func (a *Alarm) PartialUpdateComment(timestamp CpsTime, author, output, role, initiator string) error {
	newStep := NewAlarmStep(AlarmStepComment, timestamp, author, output, role, initiator)
	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}

	a.addUpdate("$push", bson.M{"v.steps": newStep})

	return nil
}

func (a *Alarm) PartialUpdateCropSteps() {
	if a.CropSteps() {
		a.addUpdate("$set", bson.M{"v.steps": a.Value.Steps})
	}
}

func (a *Alarm) PartialUpdateAddStep(stepType string, timestamp CpsTime, author, msg, role, initiator string) error {
	newStep := NewAlarmStep(stepType, timestamp, author, msg, role, initiator)
	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}

	a.addUpdate("$push", bson.M{"v.steps": newStep})

	return nil
}

// addUpdate adds new mongo updates.
func (a *Alarm) addUpdate(key string, update bson.M) {
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
