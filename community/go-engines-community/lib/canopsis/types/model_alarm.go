package types

import (
	"log"
	"regexp"
	"sort"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// Alarm states
const (
	AlarmStateOK = iota
	AlarmStateMinor
	AlarmStateMajor
	AlarmStateCritical
)

const (
	AlarmStateTitleOK       = "ok"
	AlarmStateTitleMinor    = "minor"
	AlarmStateTitleMajor    = "major"
	AlarmStateTitleCritical = "critical"
)

// Alarm statuses
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

// Alarm steps
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
	AlarmStepSnooze          = "snooze"
	AlarmStepStateCounter    = "statecounter"
	AlarmStepChangeState     = "changestate"
	AlarmStepPbhEnter        = "pbhenter"
	AlarmStepPbhLeave        = "pbhleave"
	AlarmStepMetaAlarmAttach = "metaalarmattach"

	AlarmStepAssocTicket       = "assocticket"
	AlarmStepDeclareTicket     = "declareticket"
	AlarmStepDeclareTicketFail = "declareticketfail"
	AlarmStepWebhookStart      = "webhookstart"
	AlarmStepWebhookComplete   = "webhookcomplete"
	AlarmStepWebhookFail       = "webhookfail"

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
	ID       string  `bson:"_id" json:"_id"`
	Time     CpsTime `bson:"t" json:"t"`
	EntityID string  `bson:"d" json:"d"`

	Tags                []string  `bson:"tags" json:"tags"`
	ExternalTags        []string  `bson:"etags" json:"etags"`
	InternalTags        []string  `bson:"itags" json:"itags"`
	InternalTagsUpdated MicroTime `bson:"itags_upd" json:"itags_upd"`
	// todo move all field from Value to Alarm
	Value AlarmValue `bson:"v" json:"v"`
	// update contains alarm changes after last mongo update. Use functions Update* to
	// fill it.
	update         bson.M
	childrenUpdate []string
	childrenRemove []string
	parentsUpdate  []string
	parentsRemove  []string

	// is used only for manual instructions KPI metrics
	KpiAssignedInstructions []string `bson:"kpi_assigned_instructions,omitempty" json:"kpi_assigned_instructions,omitempty"`
	KpiExecutedInstructions []string `bson:"kpi_executed_instructions,omitempty" json:"kpi_executed_instructions,omitempty"`
	// is used only for auto instructions KPI metrics
	KpiAssignedAutoInstructions []string `bson:"kpi_assigned_auto_instructions,omitempty" json:"kpi_assigned_auto_instructions,omitempty"`
	KpiExecutedAutoInstructions []string `bson:"kpi_executed_auto_instructions,omitempty" json:"kpi_executed_auto_instructions,omitempty"`

	// is used only for not acked metrics
	NotAckedMetricType     string   `bson:"not_acked_metric_type,omitempty" json:"-"`
	NotAckedMetricSendTime *CpsTime `bson:"not_acked_metric_send_time,omitempty" json:"-"`
	NotAckedSince          *CpsTime `bson:"not_acked_since,omitempty" json:"-"`

	// InactiveAutoInstructionInProgress shows that autoremediation is launched and alarm is not active until the remediation is finished
	InactiveAutoInstructionInProgress bool `bson:"auto_instruction_in_progress,omitempty" json:"auto_instruction_in_progress,omitempty"`

	Healthcheck bool `bson:"healthcheck,omitempty" json:"-"`
}

// AlarmWithEntity is an encapsulated type, mostly to facilitate the alarm manipulation for the post-processors
type AlarmWithEntity struct {
	Alarm  Alarm  `bson:"alarm" json:"alarm"`
	Entity Entity `bson:"entity" json:"entity"`
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
func (a *Alarm) GetAppliedActions() (steps AlarmSteps) {
	steps = make([]AlarmStep, 0, 3)

	if a.Value.ACK != nil {
		steps = append(steps, *a.Value.ACK)
	}

	for _, ticketStep := range a.Value.Tickets {
		if ticketStep.Type == AlarmStepDeclareTicket || ticketStep.Type == AlarmStepAssocTicket {
			steps = append(steps, ticketStep)
		}
	}
	if a.IsSnoozed() {
		steps = append(steps, *a.Value.Snooze)
	}
	if a.Value.LastComment != nil {
		steps = append(steps, *a.Value.LastComment)
	}
	sort.Sort(ByTimestamp{steps})
	return steps
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

// Closable checks the last step for it's state to be OK for at least d interval.
// Reference time is time.Now() when this function is called.
func (a *Alarm) Closable(d time.Duration) bool {
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
func (a *Alarm) IsAck() bool {
	return a.Value.ACK != nil // && !a.IsResolved()
}

// IsCanceled check if an Alarm is canceled
func (a *Alarm) IsCanceled() bool {
	return a.Value.Canceled != nil && !a.IsResolved()
}

// IsMatched tell if an alarm is catched by a regex
func (a *Alarm) IsMatched(regex string, fields []string) bool {
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
func (a *Alarm) IsSnoozed() bool {
	if a.Value.Snooze == nil {
		return false
	}

	snoozeEnd := a.Value.Snooze.Value.CpsTimestamp()
	return snoozeEnd.After(NewCpsTime())
}

// IsStateLocked checks that the Alarm is not Locked (by manual intervention for example)
func (a *Alarm) IsStateLocked() bool {
	return a.Value.ChangeState != nil
}

// IsMalfunctioning...
func (a *Alarm) IsMalfunctioning() bool {
	return a.Value.Status.Value != AlarmStateOK
}

// HasSingleAck returns true if the alarm has been acknowledged exactly once.
// Note that this method will return false if the alarm has received a first
// ack, an ackremove, and a second ack.
// It should be used to run actions on the first acknowledgement only.
func (a *Alarm) HasSingleAck() bool {
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

func (a *Alarm) IsMetaAlarm() bool {
	return a.Value.Meta != ""
}

func (a *Alarm) IsMetaChild() bool {
	return len(a.Value.Parents) > 0
}

func (a *Alarm) HasChildByEID(childEID string) bool {
	for _, child := range a.Value.Children {
		if child == childEID {
			return true
		}
	}

	return false
}

func (a *Alarm) HasParentByEID(parentEID string) bool {
	for _, parent := range a.Value.Parents {
		if parent == parentEID {
			return true
		}
	}

	return false
}

func (a *Alarm) HasUnlinkedParentByEID(parentEID string) bool {
	for _, parent := range a.Value.UnlinkedParents {
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

func (a *Alarm) AddParent(parentEID string) bool {
	if a.HasParentByEID(parentEID) {
		return false
	}

	a.Value.Parents = append(a.Value.Parents, parentEID)
	a.parentsUpdate = append(a.parentsUpdate, parentEID)
	a.AddUpdate("$addToSet", bson.M{"v.parents": bson.M{"$each": a.parentsUpdate}})
	return true
}

func (a *Alarm) RemoveParent(parentEID string) bool {
	removed := false
	for idx, parent := range a.Value.Parents {
		if parent == parentEID {
			a.Value.Parents = append(a.Value.Parents[:idx], a.Value.Parents[idx+1:]...)
			removed = true

			break
		}
	}

	if !removed {
		return false
	}

	a.parentsRemove = append(a.parentsRemove, parentEID)
	a.Value.UnlinkedParents = append(a.Value.UnlinkedParents, parentEID)
	a.AddUpdate("$pull", bson.M{"v.parents": bson.M{"$in": a.parentsRemove}})
	a.AddUpdate("$push", bson.M{"v.unlinked_parents": parentEID})
	return true
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

func (a *Alarm) IsActivated() bool {
	return a.Value.ActivationDate != nil
}

func (a *Alarm) IsInActivePeriod() bool {
	return a.Value.PbehaviorInfo.IsActive()
}

func (a *Alarm) CanActivate() bool {
	return !a.IsActivated() && !a.IsSnoozed() && a.Value.PbehaviorInfo.IsActive() && !a.InactiveAutoInstructionInProgress
}

func (a *Alarm) IncrementEventsCount(count CpsNumber) {
	a.AddUpdate("$inc", bson.M{"v.events_count": count})
}

func (a *Alarm) DecrementEventsCount(count CpsNumber) {
	a.AddUpdate("$inc", bson.M{"v.events_count": -count})
}
