package types

import (
	"log"
	"regexp"
	"sort"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
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
	AlarmStepComment         = "comment"
	AlarmStepSnooze          = "snooze"
	AlarmStepUnsnooze        = "unsnooze"
	AlarmStepStateCounter    = "statecounter"
	AlarmStepChangeState     = "changestate"
	AlarmStepPbhEnter        = "pbhenter"
	AlarmStepPbhLeave        = "pbhleave"
	AlarmStepMetaAlarmAttach = "metaalarmattach"
	AlarmStepMetaAlarmDetach = "metaalarmdetach"
	AlarmStepActivate        = "activate"
	AlarmStepResolve         = "resolve"

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

	// Following alarm steps are used for junit.
	AlarmStepJunitTestSuiteUpdate = "junittestsuiteupdate"
	AlarmStepJunitTestCaseUpdate  = "junittestcaseupdate"
)

func GetAlarmStepTypes() []string {
	return []string{
		AlarmStepStateIncrease,
		AlarmStepStateDecrease,
		AlarmStepStatusIncrease,
		AlarmStepStatusDecrease,
		AlarmStepAck,
		AlarmStepAckRemove,
		AlarmStepComment,
		AlarmStepSnooze,
		AlarmStepUnsnooze,
		AlarmStepStateCounter,
		AlarmStepChangeState,
		AlarmStepPbhEnter,
		AlarmStepPbhLeave,
		AlarmStepMetaAlarmAttach,
		AlarmStepMetaAlarmDetach,
		AlarmStepActivate,
		AlarmStepResolve,
		AlarmStepAssocTicket,
		AlarmStepDeclareTicket,
		AlarmStepDeclareTicketFail,
		AlarmStepWebhookStart,
		AlarmStepWebhookComplete,
		AlarmStepWebhookFail,
		AlarmStepInstructionStart,
		AlarmStepInstructionPause,
		AlarmStepInstructionResume,
		AlarmStepInstructionComplete,
		AlarmStepInstructionFail,
		AlarmStepInstructionAbort,
		AlarmStepAutoInstructionStart,
		AlarmStepAutoInstructionComplete,
		AlarmStepAutoInstructionFail,
		AlarmStepJunitTestSuiteUpdate,
		AlarmStepJunitTestCaseUpdate,
	}
}

func NeedEntityLookup(collectionName string) bool {
	return mongo.AlarmMongoCollection != collectionName
}

// Alarm represents an alarm document.
type Alarm struct {
	ID       string           `bson:"_id" json:"_id"`
	Time     datetime.CpsTime `bson:"t" json:"t"`
	EntityID string           `bson:"d" json:"d"`

	Tags                []string           `bson:"tags,omitempty" json:"tags,omitempty"`
	ExternalTags        []string           `bson:"etags,omitempty" json:"etags,omitempty"`
	InternalTags        []string           `bson:"itags,omitempty" json:"itags,omitempty"`
	InternalTagsUpdated datetime.MicroTime `bson:"itags_upd" json:"itags_upd"`
	// todo move all field from Value to Alarm
	Value AlarmValue `bson:"v" json:"v"`

	// is used only for manual instructions KPI metrics
	KpiAssignedInstructions []string `bson:"kpi_assigned_instructions,omitempty" json:"kpi_assigned_instructions,omitempty"`
	KpiExecutedInstructions []string `bson:"kpi_executed_instructions,omitempty" json:"kpi_executed_instructions,omitempty"`
	// is used only for auto instructions KPI metrics
	KpiAssignedAutoInstructions []string `bson:"kpi_assigned_auto_instructions,omitempty" json:"kpi_assigned_auto_instructions,omitempty"`
	KpiExecutedAutoInstructions []string `bson:"kpi_executed_auto_instructions,omitempty" json:"kpi_executed_auto_instructions,omitempty"`

	// is used only for not acked metrics
	NotAckedMetricType     string            `bson:"not_acked_metric_type,omitempty" json:"-"`
	NotAckedMetricSendTime *datetime.CpsTime `bson:"not_acked_metric_send_time,omitempty" json:"-"`
	NotAckedSince          *datetime.CpsTime `bson:"not_acked_since,omitempty" json:"-"`

	// InactiveAutoInstructionInProgress shows that autoremediation is launched and alarm is not active until the remediation is finished
	InactiveAutoInstructionInProgress bool `bson:"auto_instruction_in_progress,omitempty" json:"auto_instruction_in_progress,omitempty"`

	InactiveDelayMetaAlarmInProgress bool `bson:"inactive_delay_meta_alarm_in_progress,omitempty" json:"inactive_delay_meta_alarm_in_progress,omitempty"`
	// MetaAlarmInactiveDelay shows that an alarm is matched to some meta alarm rules with child_inactive_delay
	MetaAlarmInactiveDelay []MetaAlarmInactiveDelay `bson:"meta_alarm_inactive_delay,omitempty" json:"meta_alarm_inactive_delay,omitempty"`

	Healthcheck bool `bson:"healthcheck,omitempty" json:"-"`
}

// AlarmWithEntityField is used to store an alarm with it's entity in the database.
// Use Alarm in all other cases.
type AlarmWithEntityField struct {
	Alarm `bson:",inline"`

	Entity Entity `bson:"entity"`
}

type MetaAlarmInactiveDelay struct {
	ID      string           `bson:"_id"`
	Expired datetime.CpsTime `bson:"expired"`
}

// AlarmWithEntity is an encapsulated type, mostly to facilitate the alarm manipulation for the post-processors
type AlarmWithEntity struct {
	Alarm  Alarm  `bson:"alarm" json:"alarm"`
	Entity Entity `bson:"entity" json:"entity"`
}

// CropSteps calls Crop() on Alarm.Value.Steps with alarm parameters.
// returns true if the alarm was modified.
func (a *Alarm) CropSteps(cropNum int) bool {
	updated := false
	if a.Value.Status != nil {
		croppedSteps, cropUpdate := a.Value.Steps.Crop(
			a.Value.Status,
			cropNum,
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
	return snoozeEnd.After(datetime.NewCpsTime())
}

// IsStateLocked checks that the Alarm is not Locked (by manual intervention for example)
func (a *Alarm) IsStateLocked() bool {
	return a.Value.ChangeState != nil
}

// IsMalfunctioning...
func (a *Alarm) IsMalfunctioning() bool {
	return a.Value.State.Value != AlarmStateOK
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
}

func (a *Alarm) RemoveChild(childEID string) {
	for idx, child := range a.Value.Children {
		if child == childEID {
			a.Value.Children = append(a.Value.Children[:idx], a.Value.Children[idx+1:]...)

			return
		}
	}
}

func (a *Alarm) AddParent(parentEID string) bool {
	if a.HasParentByEID(parentEID) {
		return false
	}

	a.Value.Parents = append(a.Value.Parents, parentEID)

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

	a.Value.UnlinkedParents = append(a.Value.UnlinkedParents, parentEID)

	return true
}

func (a *Alarm) IsActivated() bool {
	return a.Value.ActivationDate != nil
}

func (a *Alarm) IsInActivePeriod() bool {
	return a.Value.PbehaviorInfo.IsActive()
}

func (a *Alarm) CanActivate() bool {
	return !a.IsActivated() &&
		!a.IsSnoozed() &&
		a.Value.PbehaviorInfo.IsActive() &&
		!a.InactiveAutoInstructionInProgress &&
		!a.InactiveDelayMetaAlarmInProgress
}

// GetStringField is a magic getter for string fields for easier field retrieving when matching alarm pattern
func (a *Alarm) GetStringField(f string) (string, bool) {
	switch f {
	case "v.display_name":
		return a.Value.DisplayName, true
	case "v.output":
		return a.Value.Output, true
	case "v.long_output":
		return a.Value.LongOutput, true
	case "v.initial_output":
		return a.Value.InitialOutput, true
	case "v.initial_long_output":
		return a.Value.InitialLongOutput, true
	case "v.connector":
		return a.Value.Connector, true
	case "v.connector_name":
		return a.Value.ConnectorName, true
	case "v.component":
		return a.Value.Component, true
	case "v.resource":
		return a.Value.Resource, true
	case "v.last_comment.m":
		if a.Value.LastComment == nil {
			return "", true
		}
		return a.Value.LastComment.Message, true
	case "v.last_comment.initiator":
		return a.Value.LastComment.GetInitiator(), true
	case "v.ticket.m":
		if a.Value.Ticket == nil {
			return "", true
		}

		return a.Value.Ticket.Message, true
	case "v.ticket.ticket":
		if a.Value.Ticket == nil {
			return "", true
		}

		return a.Value.Ticket.Ticket, true
	case "v.ticket.initiator":
		return a.Value.Ticket.GetInitiator(), true
	case "v.ack.a":
		if a.Value.ACK == nil {
			return "", true
		}

		return a.Value.ACK.Author, true
	case "v.ack.m":
		if a.Value.ACK == nil {
			return "", true
		}

		return a.Value.ACK.Message, true
	case "v.ack.initiator":
		return a.Value.ACK.GetInitiator(), true
	case "v.canceled.initiator":
		return a.Value.Canceled.GetInitiator(), true
	default:
		if n := strings.TrimPrefix(f, "v.ticket.ticket_data."); n != f {
			if a.Value.Ticket == nil || a.Value.Ticket.TicketData == nil {
				return "", true
			}

			return a.Value.Ticket.TicketData[n], true
		}

		return "", false
	}
}

// GetIntField is a magic getter for int fields for easier field retrieving when matching alarm pattern
func (a *Alarm) GetIntField(f string) (int64, bool) {
	switch f {
	case "v.state.val":
		if a.Value.State == nil {
			return 0, true
		}
		return int64(a.Value.State.Value), true
	case "v.status.val":
		if a.Value.Status == nil {
			return 0, true
		}
		return int64(a.Value.Status.Value), true
	case "v.total_state_changes":
		return int64(a.Value.TotalStateChanges), true
	default:
		return 0, false
	}
}

// GetRefField is a magic getter for reference fields for easier field retrieving when matching alarm pattern
func (a *Alarm) GetRefField(f string) (interface{}, bool) {
	switch f {
	case "v.ack":
		if a.Value.ACK == nil {
			return nil, true
		}
		return a.Value.ACK, true
	case "v.ticket":
		if a.Value.Ticket == nil {
			return nil, true
		}
		return a.Value.Ticket, true
	case "v.canceled":
		if a.Value.Canceled == nil {
			return nil, true
		}
		return a.Value.Canceled, true
	case "v.snooze":
		if a.Value.Snooze == nil {
			return nil, true
		}
		return a.Value.Snooze, true
	case "v.activation_date":
		if a.Value.ActivationDate == nil {
			return nil, true
		}
		return a.Value.ActivationDate, true
	case "v.change_state":
		if a.Value.ChangeState == nil {
			return nil, true
		}
		return a.Value.ChangeState, true
	case "v.meta":
		if a.Value.Meta == "" {
			return nil, true
		}
		return a.Value.Meta, true
	default:
		return nil, false
	}
}

// GetTimeField is a magic getter for time fields for easier field retrieving when matching alarm pattern
func (a *Alarm) GetTimeField(f string) (time.Time, bool) {
	switch f {
	case "v.creation_date":
		return a.Value.CreationDate.Time, true
	case "v.last_event_date":
		return a.Value.LastEventDate.Time, true
	case "v.last_update_date":
		return a.Value.LastUpdateDate.Time, true
	case "v.ack.t":
		if a.Value.ACK != nil {
			return a.Value.ACK.Timestamp.Time, true
		}

		return time.Time{}, true
	case "v.resolved":
		if a.Value.Resolved != nil {
			return a.Value.Resolved.Time, true
		}

		return time.Time{}, true
	case "v.activation_date":
		if a.Value.ActivationDate != nil {
			return a.Value.ActivationDate.Time, true
		}
		return time.Time{}, true
	default:
		return time.Time{}, false
	}
}

// GetDurationField is a magic getter for duration fields for easier field retrieving when matching alarm pattern
func (a *Alarm) GetDurationField(f string) (int64, bool) {
	switch f {
	case "v.duration":
		if a.Value.Duration > 0 {
			return a.Value.Duration, true
		}

		if a.Value.Resolved != nil {
			return int64(a.Value.Resolved.Sub(a.Time.Time).Seconds()), true
		}

		return int64(time.Since(a.Time.Time).Seconds()), true
	default:
		return 0, false
	}
}

// GetStringArrayField is a magic getter for string array fields for easier field retrieving when matching alarm pattern
func (a *Alarm) GetStringArrayField(f string) ([]string, bool) {
	switch f {
	case "tags":
		return a.Tags, true
	default:
		return nil, false
	}
}

// GetInfoVal is a magic getter for infos fields for easier field retrieving when matching alarm pattern
func (a *Alarm) GetInfoVal(f string) (interface{}, bool) {
	for _, infosByRule := range a.Value.Infos {
		if v, ok := infosByRule[f]; ok {
			return v, true
		}
	}

	return nil, false
}

func (a *Alarm) GetMetaAlarmInactiveExpiration(ruleId string) (datetime.CpsTime, bool) {
	for _, delay := range a.MetaAlarmInactiveDelay {
		if delay.ID == ruleId {
			return delay.Expired, true
		}
	}

	return datetime.CpsTime{}, false
}
