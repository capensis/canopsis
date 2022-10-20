package types

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"sort"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
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
	AlarmStepAutoInstructionStart    = "autoinstructionstart"
	AlarmStepAutoInstructionComplete = "autoinstructioncomplete"
	AlarmStepAutoInstructionFail     = "autoinstructionfail"
	// Following alarm steps are used for job execution.
	AlarmStepInstructionJobStart    = "instructionjobstart"
	AlarmStepInstructionJobComplete = "instructionjobcomplete"
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
			LastUpdateDate:    event.Timestamp,
			LastEventDate:     now,
			DisplayName:       GenDisplayName(alarmConfig.DisplayNameScheme),
			Infos:             make(map[string]map[string]interface{}),
			RuleVersion:       make(map[string]string),
			State: &AlarmStep{
				Type:      AlarmStepStateIncrease,
				Timestamp: event.Timestamp,
				Author:    event.Connector + "." + event.ConnectorName,
				Message:   event.Output,
				Value:     event.State,
			},
			Status: &AlarmStep{
				Type:      AlarmStepStatusIncrease,
				Timestamp: event.Timestamp,
				Author:    event.Connector + "." + event.ConnectorName,
				Message:   event.Output,
				Value:     AlarmStatusOngoing,
			},
			TotalStateChanges: 1,
		},
	}
	alarm.Value.LongOutputHistory = append(alarm.Value.LongOutputHistory, event.LongOutput)
	alarm.Value.Steps = AlarmSteps{*alarm.Value.State, *alarm.Value.Status}

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
		steps = append(steps, NewAlarmStep(ticket.Type, ticket.Timestamp, ticket.Author, ticket.Message, ticket.UserID, ticket.Role, ""))
	}
	if a.IsSnoozed() {
		steps = append(steps, *a.Value.Snooze)
	}
	sort.Sort(ByTimestamp{steps})
	return steps, ticket
}

// Apply actions (ACK, Snooze, AssocTicket, DeclareTicket) from steps to alarm
func (a *Alarm) ApplyActions(steps AlarmSteps, ticket *AlarmTicket, allowDoubleAck bool) error {
	ts := NewCpsTime(time.Now().Unix())

	for j := 0; j < len(steps); j++ {
		step := steps[j]
		step.Author = "correlation"
		step.Timestamp = ts
		switch step.Type {
		case AlarmStepAck:
			err := a.PartialUpdateAck(ts, step.Author, step.Message, step.UserID, step.Role, step.Initiator, allowDoubleAck)
			if err != nil {
				return err
			}
		case AlarmStepSnooze:
			d := DurationWithUnit{
				Value: int64(step.Value) - step.Timestamp.Unix(),
				Unit:  "s",
			}
			err := a.PartialUpdateSnooze(ts, d, step.Author, step.Message, step.UserID, step.Role, step.Initiator)
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

			err := a.PartialUpdateAssocTicket(ts, ticket.Data, step.Author, ticket.Value, step.UserID, step.Role, step.Initiator)
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

			err := a.PartialUpdateDeclareTicket(ts, step.Author, step.Message, ticket.Value, ticket.Data, step.UserID, step.Role, step.Initiator)
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

// UpdateOutput updates an alarm output field
func (a *Alarm) UpdateOutput(newOutput string) {
	a.Value.Output = newOutput

	a.AddUpdate("$set", bson.M{
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

		a.AddUpdate("$set", bson.M{
			"v.long_output":         a.Value.LongOutput,
			"v.long_output_history": a.Value.LongOutputHistory,
		})
	}
}

// Resolve mark as resolved an Alarm with a timestamp [sic]
func (a *Alarm) Resolve(timestamp *CpsTime) {
	a.Value.Resolved = timestamp
}

// Closable checks the last step for it's state to be OK for at least d interval.
// Reference time is time.Now() when this function is called.
func (a Alarm) Closable(d time.Duration) bool {
	// prevent some silly crash
	if a.Value.State == nil {
		return false
	}

	alarmState := a.Value.State.Value
	ls, err := a.Value.Steps.Last()
	if err == nil && alarmState == AlarmStateOK && ls.Timestamp.Time.Before(time.Now().Add(-d)) {
		return true
	}

	if err != nil && alarmState != AlarmStateOK {
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
	return snoozeEnd.After(NewCpsTime())
}

// IsStateLocked checks that the Alarm is not Locked (by manual intervention for example)
func (a *Alarm) IsStateLocked() bool {
	return a.Value.State != nil && a.Value.State.Type == AlarmStepChangeState
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

	a.AddUpdate("$addToSet", bson.M{"v.children": bson.M{"$each": a.childrenUpdate}})
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
		a.AddUpdate("$pull", bson.M{"v.children": bson.M{"$in": a.childrenRemove}})
	}
}

func (a *Alarm) AddParent(parentEID string) {
	if a.HasParentByEID(parentEID) {
		return
	}

	a.Value.Parents = append(a.Value.Parents, parentEID)
	a.parentsUpdate = append(a.parentsUpdate, parentEID)
	a.AddUpdate("$addToSet", bson.M{"v.parents": bson.M{"$each": a.parentsUpdate}})
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
		a.AddUpdate("$pull", bson.M{"v.parents": bson.M{"$in": a.parentsRemove}})
	}
}

func (a *Alarm) SetMeta(meta string) {
	a.Value.Meta = meta
	a.AddUpdate("$set", bson.M{"v.meta": meta})
}

func (a *Alarm) SetMetaValuePath(path string) {
	a.Value.MetaValuePath = path
	a.AddUpdate("$set", bson.M{"v.meta_value_path": path})
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
