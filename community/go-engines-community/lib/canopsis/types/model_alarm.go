package types

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"sort"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
)

//Alarm states
const (
	AlarmStateOK = iota
	AlarmStateMinor
	AlarmStateMajor
	AlarmStateCritical
	AlarmStateUnknown
)

const (
	AlarmStateTitleOK       = "ok"
	AlarmStateTitleMinor    = "minor"
	AlarmStateTitleMajor    = "major"
	AlarmStateTitleCritical = "critical"
)

//Alarm statuses
const (
	AlarmStatusOff = iota
	AlarmStatusOngoing
	AlarmStatusStealthy
	AlarmStatusFlapping
	AlarmStatusCancelled
	AlarmStatusNoEvents
)

const (
	AlarmStatusTitleOff       = "off"
	AlarmStatusTitleOngoing   = "ongoing"
	AlarmStatusTitleStealthy  = "stealthy"
	AlarmStatusTitleFlapping  = "flapping"
	AlarmStatusTitleCancelled = "cancelled"
)

//Alarm steps
const (
	AlarmStepStateIncrease   = "stateinc"
	AlarmStepStateDecrease   = "statedec"
	AlarmStepStatusIncrease  = "statusinc"
	AlarmStepStatusDecrease  = "statusdec"
	AlarmStepAck             = "ack"
	AlarmStepAckRemove       = "ackremove"
	AlarmStepCancel          = "cancel"
	AlarmStepUncancel        = "uncancel"
	AlarmStepComment         = "comment"
	AlarmStepDone            = "done"
	AlarmStepDeclareTicket   = "declareticket"
	AlarmStepAssocTicket     = "assocticket"
	AlarmStepSnooze          = "snooze"
	AlarmStepStateCounter    = "statecounter"
	AlarmStepChangeState     = "changestate"
	AlarmStepPbhEnter        = "pbhenter"
	AlarmStepPbhLeave        = "pbhleave"
	AlarmStepMetaAlarmAttach = "metaalarmattach"

	// Following alarm steps are used for manual instruction execution.
	AlarmStepInstructionStart    = "instructionstart"
	AlarmStepInstructionPause    = "instructionpause"
	AlarmStepInstructionResume   = "instructionresume"
	AlarmStepInstructionComplete = "instructioncomplete"
	AlarmStepInstructionFail     = "instructionfail"
	// AlarmStepInstructionAbort are used for manual and auto instruction execution.
	AlarmStepInstructionAbort = "instructionabort"
	// Following alarm steps are used for manual instruction execution.
	AlarmStepAutoInstructionStart          = "autoinstructionstart"
	AlarmStepAutoInstructionComplete       = "autoinstructioncomplete"
	AlarmStepAutoInstructionFail           = "autoinstructionfail"
	AlarmStepAutoInstructionAlreadyRunning = "autoinstructionalreadyrunning"
	// Following alarm steps are used for job execution.
	AlarmStepInstructionJobStart    = "instructionjobstart"
	AlarmStepInstructionJobComplete = "instructionjobcomplete"
	AlarmStepInstructionJobAbort    = "instructionjobabort"
	AlarmStepInstructionJobFail     = "instructionjobfail"

	// Following alarm steps are used for junit.
	AlarmStepJunitTestSuiteUpdate = "junittestsuiteupdate"
	AlarmStepJunitTestCaseUpdate  = "junittestcaseupdate"
)

const (
	StepEngineCorrelationAuthor = "engine.correlation"
)

// Alarm represents an alarm document.
type Alarm struct {
	ID       string     `bson:"_id" json:"_id"`
	Time     CpsTime    `bson:"t" json:"t"`
	EntityID string     `bson:"d" json:"d"`
	Value    AlarmValue `bson:"v" json:"v"`
	// update contains alarm changes after last mongo update. Use functions Update* to
	// fill it.
	update         bson.M
	childrenUpdate []string
	childrenRemove []string
	parentsUpdate  []string
	parentsRemove  []string
}

// AlarmWithEntity is an encapsulated type, mostly to facilitate the alarm manipulation for the post-processors
type AlarmWithEntity struct {
	Alarm  Alarm  `bson:"alarm" json:"alarm"`
	Entity Entity `bson:"entity" json:"entity"`
}

// NewAlarm creates en new Alarm from an Event
func NewAlarm(event Event, alarmConfig config.AlarmConfig) (Alarm, error) {
	// When the alarms are written to the database, the dates are converted to
	// timestamps, in seconds.
	// The dates used internally should be the same as the ones stored in the
	// database, so that the engines do not work on different versions of the
	// same alarm.
	now := CpsTime{time.Now().Truncate(time.Second)}

	if event.Timestamp.IsZero() {
		return Alarm{}, errors.New("field Timestamp is not set")
	}

	alarm := Alarm{
		EntityID: event.GetEID(),
		ID:       utils.NewID(),
		Time:     now,
		Value: AlarmValue{
			Connector:         event.Connector,
			ConnectorName:     event.ConnectorName,
			Component:         event.Component,
			Resource:          event.Resource,
			Output:            event.Output,
			InitialOutput:     event.Output,
			InitialLongOutput: event.LongOutput,
			LongOutput:        event.LongOutput,
			CreationDate:      now,
			LastUpdateDate:    now,
			LastEventDate:     now,
			DisplayName:       GenDisplayName(alarmConfig.DisplayNameScheme),
			Infos:             make(map[string]map[string]interface{}),
			RuleVersion:       make(map[string]string),
		},
	}
	alarm.Value.LongOutputHistory = append(alarm.Value.LongOutputHistory, event.LongOutput)
	alarm.Update(event, alarmConfig)

	return alarm, nil
}

// AlarmID build an alarmid from given parameters. Used by Alarm.AlarmID()
func AlarmID(connector, connectorName, entityID string) string {
	return fmt.Sprintf(
		"%s/%s/%s",
		connector,
		connectorName,
		entityID,
	)
}

// AlarmID returns current alarm's alarmid.
func (a Alarm) AlarmID() string {
	return AlarmID(
		a.Value.Connector,
		a.Value.ConnectorName,
		a.EntityID,
	)
}

// AlarmComponentID is like Alarm.AlarmID() but uses Alarm.Value.Component
// instead of Alarm.EntityID
func (a Alarm) AlarmComponentID() string {
	return AlarmID(
		a.Value.Connector,
		a.Value.ConnectorName,
		a.Value.Component,
	)
}

// CacheID implements cache.Cache interface
func (a Alarm) CacheID() string {
	return a.AlarmID()
}

// CropSteps calls Crop() on Alarm.Value.Steps with alarm parameters.
// returns true if the alarm was modified.
func (a *Alarm) CropSteps() bool {
	updated := false
	if a.Value.Status != nil {
		croppedSteps, cropUpdate := a.Value.Steps.Crop(
			a.Value.Status,
			AlarmStepCropMinStates,
		)
		// Updates the alarm steps
		a.Value.Steps = croppedSteps

		updated = updated || cropUpdate
	}
	return updated
}

// GetAppliedActions fetches applied to alarm actions: ACK, Snooze, AssocTicket, DeclareTicket
// Result is in a sorted by timestamp AlarmSteps, ticket data when defined
func (a *Alarm) GetAppliedActions() (steps AlarmSteps, ticket *AlarmTicket) {
	steps = make([]AlarmStep, 0, 3)

	if a.Value.ACK != nil {
		steps = append(steps, *a.Value.ACK)
	}
	if ticket = a.Value.Ticket; ticket != nil {
		steps = append(steps, NewAlarmStep(ticket.Type, ticket.Timestamp, ticket.Author, ticket.Message, ticket.Role, ""))
	}
	if a.IsSnoozed() {
		steps = append(steps, *a.Value.Snooze)
	}
	sort.Sort(ByTimestamp{steps})
	return steps, ticket
}

// Apply actions (ACK, Snooze, AssocTicket, DeclareTicket) from steps to alarm
func (a *Alarm) ApplyActions(steps AlarmSteps, ticket *AlarmTicket) error {
	ts := NewCpsTime(time.Now().Unix())

	for j := 0; j < len(steps); j++ {
		step := steps[j]
		step.Author = "correlation"
		step.Timestamp = ts
		switch step.Type {
		case AlarmStepAck:
			err := a.PartialUpdateAck(ts, step.Author, step.Message, step.Role, step.Initiator)
			if err != nil {
				return err
			}
		case AlarmStepSnooze:
			err := a.PartialUpdateSnooze(ts, step.Value, step.Author, step.Message, step.Role, step.Initiator)
			if err != nil {
				return err
			}
		case AlarmStepAssocTicket:
			if a.Value.Ticket != nil {
				continue
			}

			if ticket == nil {
				continue
			}

			err := a.PartialUpdateAssocTicket(ts, ticket.Data, step.Author, ticket.Value, step.Role, step.Initiator)
			if err != nil {
				return err
			}
		case AlarmStepDeclareTicket:
			if a.Value.Ticket != nil {
				continue
			}

			if ticket == nil {
				continue
			}

			err := a.PartialUpdateDeclareTicket(ts, step.Author, step.Message, ticket.Value, ticket.Data, step.Role, step.Initiator)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported action type: %s", step.Type)
		}
	}

	return nil
}

// CurrentState returns the Current State of the Alarm
func (a *Alarm) CurrentState() CpsNumber {
	if a.Value.State == nil {
		return AlarmStateOK
	}

	alarmState := a.Value.State.Value

	if alarmState < AlarmStateOK {
		return AlarmStateOK
	}

	return alarmState
}

// CurrentStatus returns the Current status of the alarm
func (a *Alarm) CurrentStatus(alarmConfig config.AlarmConfig) CpsNumber {
	if a.Value.Status == nil {
		return a.ComputeStatus(alarmConfig)
	}

	currentAlarmStatus := a.Value.Status.Value
	if currentAlarmStatus < AlarmStatusOff {
		return AlarmStatusOff
	}

	return currentAlarmStatus
}

// Update an alarm from an Event
func (a *Alarm) Update(e Event, alarmConfig config.AlarmConfig) bool {
	timestamp := e.Timestamp
	author := e.Connector + "." + e.ConnectorName
	message := e.Output

	updatedState := a.updateState(e)
	updatedStatus := a.UpdateStatus(timestamp, author, message, alarmConfig)
	return updatedState || updatedStatus
}

// UpdateOutput updates an alarm output field
func (a *Alarm) UpdateOutput(newOutput string) {
	a.Value.Output = newOutput

	a.addUpdate("$set", bson.M{
		"v.output": a.Value.Output,
	})
}

// UpdateLongOutput updates an alarm output field
func (a *Alarm) UpdateLongOutput(newOutput string) {
	if (len(a.Value.LongOutputHistory) == 0) || (a.Value.LongOutputHistory[len(a.Value.LongOutputHistory)-1] != newOutput) {
		a.Value.LongOutput = newOutput
		history := append(a.Value.LongOutputHistory, newOutput)
		if len(history) > 100 {
			history = history[len(history)-100:]
		}
		a.Value.LongOutputHistory = history

		a.addUpdate("$set", bson.M{
			"v.long_output":         a.Value.LongOutput,
			"v.long_output_history": a.Value.LongOutputHistory,
		})
	}
}

// updateState If Event is a check event and if the associated alarm is not
// locked, increases or decreases the State of the alarm.
func (a *Alarm) updateState(e Event) bool {
	if e.EventType != EventTypeCheck && e.EventType != EventTypeMetaAlarm {
		return false
	}

	currentAlarmState := a.CurrentState()
	if e.State == currentAlarmState {
		// State hasn't changed. Stop here
		return false
	}
	if a.IsStateLocked() {
		// Event is an OK, so the alarm should be resolved anyway
		if e.State == AlarmStateOK {
			a.Value.State.Type = ""
		} else {
			// Stop here: state should be preserved.
			return false
		}
	}

	// Create new Step to keep track of the alarm history
	newStep := AlarmStep{
		Type:      AlarmStepStateIncrease,
		Timestamp: e.Timestamp,
		Author:    e.Connector + "." + e.ConnectorName,
		Message:   e.Output,
		Value:     e.State,
	}
	if e.State < currentAlarmState {
		newStep.Type = AlarmStepStateDecrease
	}

	a.Value.State = &newStep
	a.Value.Steps.Add(newStep)

	a.Value.StateChangesSinceStatusUpdate++
	a.Value.TotalStateChanges++
	a.Value.LastUpdateDate = e.Timestamp

	return true
}

// UpdateStatus recomputes the status of the alarm. If the status changes, a
// new step is added to the alarm, with the timestamp, author and message given
// in parameters.
func (a *Alarm) UpdateStatus(timestamp CpsTime, author, message string, alarmConfig config.AlarmConfig) bool {
	currentAlarmStatus := a.CurrentStatus(alarmConfig)
	newStatus := a.ComputeStatus(alarmConfig)

	if newStatus == currentAlarmStatus && a.Value.Status != nil {
		return false
	}

	// Create new Step to keep track of the alarm history.
	newStep := AlarmStep{
		Type:      AlarmStepStatusIncrease,
		Timestamp: timestamp,
		Author:    author,
		Message:   message,
		Value:     newStatus,
	}

	// If the status is Decreasing
	if newStatus < currentAlarmStatus {
		newStep.Type = AlarmStepStatusDecrease
	}
	a.Value.Status = &newStep
	a.Value.Steps.Add(newStep)

	a.Value.StateChangesSinceStatusUpdate = CpsNumber(0)
	a.Value.LastUpdateDate = timestamp

	return true
}

// ComputeStatus of an Alarm from an Event
func (a *Alarm) ComputeStatus(alarmConfig config.AlarmConfig) CpsNumber {
	if a.IsCanceled() {
		return AlarmStatusCancelled
	}

	if IsFlapping(a) {
		return AlarmStatusFlapping
	}

	if a.IsStealthy(alarmConfig) {
		return AlarmStatusStealthy
	}

	if a.Value.State == nil {
		return AlarmStatusOff
	}

	if a.Value.State.Value != AlarmStateOK {
		return AlarmStatusOngoing
	}

	return AlarmStatusOff
}

// Ack an alarm
func (a *Alarm) Ack(event Event) AlarmStep {
	newStep := NewAlarmStepFromEvent(AlarmStepAck, event)
	a.Value.ACK = &newStep

	return newStep
}

// Unack removes an ack on an alarm
func (a *Alarm) Unack(event Event) AlarmStep {
	newStep := NewAlarmStepFromEvent(AlarmStepAckRemove, event)
	a.Value.ACK = nil

	return newStep
}

func (a *Alarm) PbhEnter(event Event) AlarmStep {
	newStep := NewAlarmStepFromEvent(AlarmStepPbhEnter, event)
	newStep.PbehaviorCanonicalType = event.PbehaviorInfo.CanonicalType

	a.Value.PbehaviorInfo = event.PbehaviorInfo

	return newStep
}

func (a *Alarm) PbhLeave(event Event) AlarmStep {
	newStep := NewAlarmStepFromEvent(AlarmStepPbhLeave, event)
	newStep.PbehaviorCanonicalType = a.Value.PbehaviorInfo.CanonicalType

	a.Value.PbehaviorInfo = PbehaviorInfo{}

	return newStep
}

// Cancel cancel an alarm
func (a *Alarm) Cancel(event Event) AlarmStep {
	newStep := NewAlarmStepFromEvent(AlarmStepCancel, event)
	a.Value.Canceled = &newStep

	return newStep
}

// Uncancel uncancel an alarm
func (a *Alarm) Uncancel(event Event) AlarmStep {
	newStep := NewAlarmStepFromEvent(AlarmStepUncancel, event)
	a.Value.Canceled = nil

	return newStep
}

// ChangeState force the state of an alarm
func (a *Alarm) ChangeState(event Event) AlarmStep {
	newStep := NewAlarmStepFromEvent(AlarmStepChangeState, event)
	newStep.Value = event.State
	a.Value.State = &newStep

	return newStep
}

// Comment comment an alarm
func (a *Alarm) Comment(event Event) AlarmStep {
	newStep := NewAlarmStepFromEvent(AlarmStepComment, event)
	// It just add a step with a comment

	return newStep
}

// Done mark an alarm as done
func (a *Alarm) Done(event Event) AlarmStep {
	newStep := NewAlarmStepFromEvent(AlarmStepDone, event)
	a.Value.Done = &newStep

	return newStep
}

// Ticket add a ticket on an alarm
func (a *Alarm) Ticket(stepType string, timestamp CpsTime, author string, ticketNumber string, role string, data map[string]string, initiator string) AlarmStep {
	newStep := NewAlarmStep(stepType, timestamp, author, ticketNumber, role, initiator)
	ticketStep := newStep.NewTicket(ticketNumber, data)
	a.Value.Ticket = &ticketStep

	return newStep
}

// AssocTicket associate a ticket number to an alarm
func (a *Alarm) AssocTicket(event Event) AlarmStep {
	data := make(map[string]string)
	return a.Ticket(AlarmStepAssocTicket, event.Timestamp, event.Author, event.Ticket, event.Role, data, event.Initiator)
}

// DeclareTicket ask for a creation
func (a *Alarm) DeclareTicket(event Event) AlarmStep {
	//TODO: declare should  not read ticket parameter, but schedule a ticket
	// creation and get back a ticket number.
	// ! BUT ! On "ack and report" action from frontend, it send a declareticket
	// event with the corresponding ticket number...
	data := make(map[string]string)
	return a.Ticket(AlarmStepDeclareTicket, event.Timestamp, event.Author, event.Ticket, event.Role, data, event.Initiator)
}

// Resolve mark as resolved an Alarm with a timestamp [sic]
func (a *Alarm) Resolve(timestamp *CpsTime) {
	a.Value.Resolved = timestamp
}

// ResolveCancel forces alarm resolution on cancel
func (a *Alarm) ResolveCancel(timestamp *CpsTime) {
	a.Resolve(timestamp)
	a.Value.Status = a.Value.Canceled
	a.Value.Status.Value = AlarmStatusCancelled
}

// Snooze apply a snooze step to an Alarm
func (a *Alarm) Snooze(timestamp CpsTime, duration CpsNumber, author, output, role, initiator string) (AlarmStep, error) {
	newStep := NewAlarmStep(AlarmStepSnooze, timestamp, author, output, role, initiator)
	if duration == 0 {
		return newStep, errt.NewUnknownError(errors.New("no duration for snoozing"))
	}
	newStep.Value = CpsNumber(timestamp.Time.Unix()) + duration
	a.Value.Snooze = &newStep

	return newStep, nil
}

// SnoozeFromEvent apply a snooze step to an Alarm
func (a *Alarm) SnoozeFromEvent(event Event) (AlarmStep, error) {
	duration := CpsNumber(0)
	if event.Duration != nil {
		duration = *event.Duration
	}
	return a.Snooze(event.Timestamp, duration, event.Author, event.Output, event.Role, event.Initiator)
}

// UnSnooze cancel a snooze
func (a *Alarm) UnSnooze() {
	a.Value.LastUpdateDate = CpsTime{time.Now()}
	a.Value.Snooze = nil
}

// Closable checks the last step for it's state to be OK for at least
// baggotTime. Reference time is time.Now() when this function is called.
func (a Alarm) Closable(baggotTime time.Duration) bool {
	// prevent some silly crash
	if a.Value.State == nil {
		return false
	}

	ls, err := a.Value.Steps.Last()
	if err == nil && a.Value.State.Value == AlarmStateOK && ls.Timestamp.Time.Before(time.Now().Add(-baggotTime)) {
		return true
	}

	if err != nil && a.Value.State.Value != AlarmStateOK {
		log.Printf("warning: alarm %s has empty steps but is not OK", a.ID)
	}

	return false
}

// IsAck check if an Alarm is acked
func (a Alarm) IsAck() bool {
	return a.Value.ACK != nil // && !a.IsResolved()
}

// IsCanceled check if an Alarm is canceled
func (a Alarm) IsCanceled() bool {
	return a.Value.Canceled != nil && !a.IsResolved()
}

// IsMatched tell if an alarm is catched by a regex
func (a Alarm) IsMatched(regex string, fields []string) bool {
	for _, fieldName := range fields {
		field := utils.GetStringField(a.Value, fieldName)
		matched, _ := regexp.MatchString(regex, field)
		if matched {
			return true
		}
	}
	return false
}

// IsResolved tell if an alarm has been resolved
func (a *Alarm) IsResolved() bool {
	return a.Value.Resolved != nil
}

// IsSnoozed check if an Alarm is snoozed
func (a Alarm) IsSnoozed() bool {
	if a.Value.Snooze == nil {
		return false
	}

	snoozeEnd := a.Value.Snooze.Value.CpsTimestamp()
	return snoozeEnd.After(time.Now())
}

// IsStateLocked checks that the Alarm is not Locked (by manual intervention for example)
func (a *Alarm) IsStateLocked() bool {
	return a.Value.State != nil && a.Value.State.Type == AlarmStepChangeState
}

// IsStealthy checks if an Alarm is currently stealthy
func (a Alarm) IsStealthy(alarmConfig config.AlarmConfig) bool {
	// FIXME: crash on watch event creating a new alarm, cause State is not initialized
	if a.Value.State.Value != AlarmStateOK {
		return false
	}

	for i := len(a.Value.Steps) - 1; i >= 0; i-- {
		s := a.Value.Steps[i]
		timeSinceThisStep := time.Since(s.Timestamp.Time)
		if timeSinceThisStep >= alarmConfig.StealthyInterval {
			break
		}

		if s.Value != AlarmStateOK {
			switch s.Type {
			case AlarmStepStatusIncrease:
				fallthrough
			case AlarmStepStateDecrease:
				return true
			default:
				break
			}
		}
	}

	return false
}

// IsMalfunctioning...
func (a Alarm) IsMalfunctioning() bool {
	return a.Value.Status.Value != AlarmStateOK
}

// HasSingleAck returns true if the alarm has been acknowledged exactly once.
// Note that this method will return false if the alarm has received a first
// ack, an ackremove, and a second ack.
// It should be used to run actions on the first acknowledgement only.
func (a Alarm) HasSingleAck() bool {
	hasAck := false
	for _, step := range a.Value.Steps {
		if step.Type == AlarmStepAck {
			if hasAck {
				// This is the second ack
				return false
			}
			hasAck = true
		}
	}
	return hasAck
}

func (a Alarm) IsMetaAlarm() bool {
	return a.Value.Meta != ""
}

func (a Alarm) IsMetaChildren() bool {
	return a.Value.Parents != nil
}

func (a Alarm) HasChildByEID(childEID string) bool {
	for _, child := range a.Value.Children {
		if child == childEID {
			return true
		}
	}

	return false
}

func (a Alarm) HasParentByEID(parentEID string) bool {
	for _, parent := range a.Value.Parents {
		if parent == parentEID {
			return true
		}
	}

	return false
}

func (a *Alarm) AddChild(childEID string) {
	if a.HasChildByEID(childEID) {
		return
	}

	a.Value.Children = append(a.Value.Children, childEID)
	a.childrenUpdate = append(a.childrenUpdate, childEID)

	a.addUpdate("$addToSet", bson.M{"v.children": bson.M{"$each": a.childrenUpdate}})
}

func (a *Alarm) RemoveChild(childEID string) {
	removed := false
	for idx, child := range a.Value.Children {
		if child == childEID {
			a.Value.Children = append(a.Value.Children[:idx], a.Value.Children[idx+1:]...)
			removed = true

			break
		}
	}

	if removed {
		a.childrenRemove = append(a.childrenRemove, childEID)
		a.addUpdate("$pull", bson.M{"v.children": bson.M{"$in": a.childrenRemove}})
	}
}

func (a *Alarm) AddParent(parentEID string) {
	if a.HasParentByEID(parentEID) {
		return
	}

	a.Value.Parents = append(a.Value.Parents, parentEID)
	a.parentsUpdate = append(a.parentsUpdate, parentEID)
	a.addUpdate("$addToSet", bson.M{"v.parents": bson.M{"$each": a.parentsUpdate}})
}

func (a *Alarm) RemoveParent(parentEID string) {
	removed := false
	for idx, parent := range a.Value.Parents {
		if parent == parentEID {
			a.Value.Parents = append(a.Value.Parents[:idx], a.Value.Parents[idx+1:]...)
			removed = true

			break
		}
	}

	if removed {
		a.parentsRemove = append(a.parentsRemove, parentEID)
		a.addUpdate("$pull", bson.M{"v.parents": bson.M{"$in": a.parentsRemove}})
	}
}

func (a *Alarm) SetMeta(meta string) {
	a.Value.Meta = meta
	a.addUpdate("$set", bson.M{"v.meta": meta})
}

func (a *Alarm) SetMetaValuePath(path string) {
	a.Value.MetaValuePath = path
	a.addUpdate("$set", bson.M{"v.meta_value_path": path})
}

// UpdateState updates alarm's state to stateValue if it worst
func (a *Alarm) UpdateState(stateValue CpsNumber, ts CpsTime) {
	if stateValue == AlarmStateUnknown || stateValue < AlarmStateMinor {
		return
	}
	if a.IsStateLocked() {
		if stateValue != AlarmStateOK {
			// Stop here: state should be preserved.
			return
		}
		// Event is an OK, so the alarm should be resolved anyway
		a.Value.State.Type = ""
	}
	if a.CurrentState() == stateValue {
		return
	}

	stepType := AlarmStepStateIncrease
	if stateValue < a.Value.State.Value {
		stepType = AlarmStepStateDecrease
	}

	a.Value.State.Value = stateValue

	// Create new Step to keep track of the alarm history
	newStep := AlarmStep{
		Type:      stepType,
		Timestamp: ts,
		Author:    a.Value.Connector + "." + a.Value.ConnectorName,
		Message:   a.Value.Output,
		Value:     stateValue,
	}

	if err := a.Value.Steps.Add(newStep); err != nil {
		return
	}
	a.Value.State = &newStep

	a.Value.StateChangesSinceStatusUpdate++
	a.Value.TotalStateChanges++
	a.Value.LastUpdateDate = ts
}

func (a *Alarm) Activate() {
	a.Value.ActivationDate = &CpsTime{time.Now()}
}

func (a Alarm) IsActivated() bool {
	return a.Value.ActivationDate != nil
}

func (a Alarm) IsInActivePeriod() bool {
	return a.Value.PbehaviorInfo.IsActive()
}
