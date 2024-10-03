package types

import (
	"errors"
	"fmt"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
	"go.mongodb.org/mongo-driver/bson"
)

// PartialUpdateAck add ack step to alarm. It saves mongo updates.
func (a *Alarm) PartialUpdateAck(timestamp CpsTime, author, output, userID, role, initiator string, allowDouble bool) error {
	if !allowDouble && a.Value.ACK != nil {
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
	a.AddUpdate("$unset", bson.M{
		"not_acked_metric_type":      "",
		"not_acked_metric_send_time": "",
		"not_acked_since":            "",
	})

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
	a.AddUpdate("$set", bson.M{
		"not_acked_since": timestamp,
	})

	return nil
}

// PartialUpdateAssocTicket add ticket to alarm. It saves mongo updates.
func (a *Alarm) PartialUpdateAssocTicket(timestamp CpsTime, author, userID, role, initiator string, ticketInfo TicketInfo) error {
	ticketStep := NewTicketStep(AlarmStepAssocTicket, timestamp, author, ticketInfo.GetStepMessage(), userID, role, initiator, ticketInfo)
	err := a.Value.Steps.Add(ticketStep)
	if err != nil {
		return err
	}

	a.Value.Tickets = append(a.Value.Tickets, ticketStep)
	a.Value.Ticket = &ticketStep

	a.AddUpdate("$set", bson.M{"v.ticket": ticketStep})
	a.AddUpdate("$push", bson.M{"v.tickets": ticketStep, "v.steps": ticketStep})

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

	a.startInactiveInterval(timestamp)

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

	a.stopInactiveInterval(timestamp)

	return nil
}

func (a *Alarm) PartialUpdatePbhEnter(timestamp CpsTime, pbehaviorInfo PbehaviorInfo, author, output, userID, role, initiator string) error {
	if a.Value.PbehaviorInfo.Same(pbehaviorInfo) {
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

	if !a.Value.PbehaviorInfo.IsActive() {
		a.startInactiveInterval(timestamp)
		a.AddUpdate("$unset", bson.M{
			"not_acked_metric_type":      "",
			"not_acked_metric_send_time": "",
			"not_acked_since":            "",
		})
	}

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
				break
			}
		}

		if !enterTimestamp.IsZero() {
			d := int64(timestamp.Sub(enterTimestamp.Time).Seconds())
			a.Value.PbehaviorInactiveDuration += d
			a.AddUpdate("$inc", bson.M{"v.pbh_inactive_duration": d})
		}

		a.stopInactiveInterval(timestamp)
		a.AddUpdate("$set", bson.M{
			"not_acked_since": timestamp,
		})
	}

	return nil
}

func (a *Alarm) PartialUpdatePbhLeaveAndEnter(timestamp CpsTime, pbehaviorInfo PbehaviorInfo, author, output, userID, role, initiator string) error {
	if a.Value.PbehaviorInfo.Same(pbehaviorInfo) {
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
		a.Value.PbehaviorInfo.ReasonName,
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

	if leaveStep.PbehaviorCanonicalType == "active" {
		if enterStep.PbehaviorCanonicalType != "active" {
			a.startInactiveInterval(timestamp)
		}
	} else {
		enterTimestamp := CpsTime{}
		for i := len(a.Value.Steps) - 3; i >= 0; i-- {
			if a.Value.Steps[i].Type == AlarmStepPbhEnter {
				enterTimestamp = a.Value.Steps[i].Timestamp
				break
			}
		}

		if !enterTimestamp.IsZero() {
			d := int64(timestamp.Sub(enterTimestamp.Time).Seconds())
			a.Value.PbehaviorInactiveDuration += d
			a.AddUpdate("$inc", bson.M{"v.pbh_inactive_duration": d})
		}

		a.stopInactiveInterval(timestamp)
	}

	return nil
}

// PartialUpdateDeclareTicket add ticket to alarm. It saves mongo updates.
func (a *Alarm) PartialUpdateDeclareTicket(timestamp CpsTime, author, userID, role, initiator string, ticketInfo TicketInfo) error {
	ticketStep := NewTicketStep(AlarmStepDeclareTicket, timestamp, author, ticketInfo.GetStepMessage(), userID, role, initiator, ticketInfo)
	err := a.Value.Steps.Add(ticketStep)
	if err != nil {
		return err
	}

	a.Value.Tickets = append(a.Value.Tickets, ticketStep)
	a.Value.Ticket = &ticketStep

	a.AddUpdate("$set", bson.M{"v.ticket": ticketStep})
	a.AddUpdate("$push", bson.M{"v.tickets": ticketStep, "v.steps": ticketStep})

	return nil
}

func (a *Alarm) PartialUpdateWebhookDeclareTicket(timestamp CpsTime, execution, author, output, userID, role, initiator string, ticketInfo TicketInfo) error {
	newStep := NewAlarmStep(AlarmStepWebhookComplete, timestamp, author, output, userID, role, initiator)
	newStep.Execution = execution
	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}

	ticketStep := NewTicketStep(AlarmStepDeclareTicket, timestamp, author, ticketInfo.GetStepMessage(), userID, role, initiator, ticketInfo)
	ticketStep.Execution = execution
	err = a.Value.Steps.Add(ticketStep)
	if err != nil {
		return err
	}

	a.Value.Tickets = append(a.Value.Tickets, ticketStep)
	a.Value.Ticket = &ticketStep

	a.AddUpdate("$set", bson.M{"v.ticket": ticketStep})
	a.AddUpdate("$push", bson.M{"v.tickets": ticketStep, "v.steps": bson.M{"$each": bson.A{newStep, ticketStep}}})

	return nil
}

func (a *Alarm) PartialUpdateWebhookDeclareTicketFail(request bool, timestamp CpsTime, execution, author, output, failReason, userID, role, initiator string, ticketInfo TicketInfo) error {
	outputBuilder := strings.Builder{}
	outputBuilder.WriteString(output)
	if failReason != "" {
		outputBuilder.WriteString(". Fail reason: ")
		outputBuilder.WriteString(failReason)
		outputBuilder.WriteRune('.')
	}

	ticketOutput := outputBuilder.String()
	requestOutput := ticketOutput
	stepType := AlarmStepWebhookFail
	if request {
		requestOutput = output
		stepType = AlarmStepWebhookComplete
	}

	newStep := NewAlarmStep(stepType, timestamp, author, requestOutput, userID, role, initiator)
	newStep.Execution = execution
	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}

	newTicketStep := NewTicketStep(AlarmStepDeclareTicketFail, timestamp, author, ticketOutput, userID, role, initiator, ticketInfo)
	newTicketStep.Execution = execution
	err = a.Value.Steps.Add(newTicketStep)
	if err != nil {
		return err
	}

	a.Value.Tickets = append(a.Value.Tickets, newTicketStep)

	a.AddUpdate("$push", bson.M{"v.tickets": newTicketStep, "v.steps": bson.M{"$each": bson.A{newStep, newTicketStep}}})

	return nil
}

func (a *Alarm) PartialUpdateWebhookStart(timestamp CpsTime, execution, author, output, userID, role, initiator string) error {
	newStep := NewAlarmStep(AlarmStepWebhookStart, timestamp, author, output, userID, role, initiator)
	newStep.Execution = execution
	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}

	a.AddUpdate("$push", bson.M{"v.steps": newStep})

	return nil
}

func (a *Alarm) PartialUpdateWebhookFail(timestamp CpsTime, execution, author, output, failReason, userID, role, initiator string) error {
	outputBuilder := strings.Builder{}
	outputBuilder.WriteString(output)
	if failReason != "" {
		outputBuilder.WriteString(". Fail reason: ")
		outputBuilder.WriteString(failReason)
		outputBuilder.WriteRune('.')
	}

	newStep := NewAlarmStep(AlarmStepWebhookFail, timestamp, author, outputBuilder.String(), userID, role, initiator)
	newStep.Execution = execution
	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}

	a.AddUpdate("$push", bson.M{"v.steps": newStep})

	return nil
}

func (a *Alarm) PartialUpdateWebhookComplete(timestamp CpsTime, execution, author, output, userID, role, initiator string) error {
	newStep := NewAlarmStep(AlarmStepWebhookComplete, timestamp, author, output, userID, role, initiator)
	newStep.Execution = execution
	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}

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
	a.Value.Duration = int64(timestamp.Sub(a.Value.CreationDate.Time).Seconds())
	a.Value.CurrentStateDuration = int64(timestamp.Sub(a.Value.State.Timestamp.Time).Seconds())

	if a.Value.Snooze != nil {
		snoozeDuration := int64(timestamp.Sub(a.Value.Snooze.Timestamp.Time).Seconds())
		a.Value.SnoozeDuration += snoozeDuration
		a.AddUpdate("$inc", bson.M{"v.snooze_duration": snoozeDuration})
	}
	if !a.Value.PbehaviorInfo.IsActive() {
		enterTimestamp := CpsTime{}
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

	if (a.Value.Snooze != nil || !a.Value.PbehaviorInfo.IsActive() || a.InactiveAutoInstructionInProgress) && a.Value.InactiveStart != nil {
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

func (a *Alarm) PartialUpdateComment(timestamp CpsTime, author, output, userID, role, initiator string) error {
	newStep := NewAlarmStep(AlarmStepComment, timestamp, author, output, userID, role, initiator)
	a.Value.LastComment = &newStep
	err := a.Value.Steps.Add(newStep)
	if err != nil {
		return err
	}
	a.AddUpdate("$set", bson.M{"v.last_comment": a.Value.LastComment})
	a.AddUpdate("$push", bson.M{"v.steps": a.Value.LastComment})

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
	if stepType == AlarmStepAutoInstructionStart && a.InactiveAutoInstructionInProgress {
		a.startInactiveInterval(newStep.Timestamp)
	}

	return nil
}

func (a *Alarm) PartialUpdateAddExecutedInstruction(instructionID string) {
	a.AddUpdate("$addToSet", bson.M{"kpi_executed_instructions": instructionID})
}

func (a *Alarm) PartialUpdateAddExecutedAutoInstruction(instructionID string) {
	a.AddUpdate("$addToSet", bson.M{"kpi_executed_auto_instructions": instructionID})
}

func (a *Alarm) PartialUpdateUnsetAutoInstructionInProgress(timestamp CpsTime) {
	a.AddUpdate("$unset", bson.M{"auto_instruction_in_progress": ""})
	a.stopInactiveInterval(timestamp)
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

func (a *Alarm) PartialUpdateTags(eventTags map[string]string) {
	exists := make(map[string]struct{}, len(a.Tags))
	for _, tag := range a.Tags {
		exists[tag] = struct{}{}
	}

	tags := TransformEventTags(eventTags)
	var k = 0
	for _, tag := range tags {
		if _, ok := exists[tag]; !ok {
			tags[k] = tag
			k++
		}
	}
	if k == 0 {
		return
	}
	tags = tags[:k]
	a.Tags = append(a.Tags, tags...)
	a.AddUpdate("$addToSet", bson.M{
		"tags": bson.M{"$each": tags},
	})
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

func (a *Alarm) startInactiveInterval(timestamp CpsTime) {
	if a.Value.InactiveStart != nil {
		inactiveDuration := int64(timestamp.Sub(a.Value.InactiveStart.Time).Seconds())
		a.Value.InactiveDuration += inactiveDuration
		a.AddUpdate("$inc", bson.M{"v.inactive_duration": inactiveDuration})
	}

	a.Value.InactiveStart = &timestamp
	a.AddUpdate("$set", bson.M{"v.inactive_start": a.Value.InactiveStart})
}

func (a *Alarm) stopInactiveInterval(timestamp CpsTime) {
	if a.Value.InactiveStart == nil {
		return
	}

	inactiveDuration := int64(timestamp.Sub(a.Value.InactiveStart.Time).Seconds())
	a.Value.InactiveDuration += inactiveDuration
	a.AddUpdate("$inc", bson.M{"v.inactive_duration": inactiveDuration})

	if a.Value.PbehaviorInfo.IsActive() && a.Value.Snooze == nil && !a.InactiveAutoInstructionInProgress {
		a.Value.InactiveStart = nil
		a.AddUpdate("$unset", bson.M{"v.inactive_start": ""})
	} else {
		a.Value.InactiveStart = &timestamp
		a.AddUpdate("$set", bson.M{"v.inactive_start": a.Value.InactiveStart})
	}
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
				// this means, that snooze step is happened after pbh_enter step,
				// it's possible to do with a scenario feature, so if it happens,
				// then elapsed time = 0
				if lastEnterTime == 0 {
					snoozeElapsed = 0
				} else {
					snoozeElapsed += lastEnterTime - step.Timestamp.Unix()
				}

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
