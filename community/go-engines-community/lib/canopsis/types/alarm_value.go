package types

import (
	"fmt"
	"strconv"
	"strings"

	cps "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
)

// Alarm consts
const (
	AlarmStepsHardLimit         = 2000
	AlarmLongOutputHistoryLimit = 100
)

const (
	RuleNameScenarioPrefix          = "Scenario: "
	RuleNameIdleRulePrefix          = "Idle rule: "
	RuleNameDeclareTicketRulePrefix = "Ticket declaration rule: "

	OutputCommentPrefix   = "Comment: "
	OutputComponentPrefix = "Component: "
)

// PbhCanonicalTypeActive is duplicate of pbehavior.TypeActive because of package cycle.
// todo: move type constants from pbehavior package to types package
const PbhCanonicalTypeActive = "active"

// AlarmStep represents a generic step used in an alarm.
type AlarmStep struct {
	Type                   string           `bson:"_t" json:"_t"`
	Timestamp              datetime.CpsTime `bson:"t" json:"t"`
	Author                 string           `bson:"a" json:"a"`
	UserID                 string           `bson:"user_id,omitempty" json:"user_id,omitempty"`
	Message                string           `bson:"m" json:"m"`
	Role                   string           `bson:"role,omitempty" json:"role,omitempty"`
	Value                  CpsNumber        `bson:"val" json:"val"`
	StateCounter           CropCounter      `bson:"statecounter,omitempty" json:"statecounter,omitempty"`
	PbehaviorCanonicalType string           `bson:"pbehavior_canonical_type,omitempty" json:"pbehavior_canonical_type,omitempty"`
	Initiator              string           `bson:"initiator,omitempty" json:"initiator,omitempty"`
	// Execution contains id
	// - of instruction execution for instruction steps
	// - of webhook execution for webhook steps
	Execution string `bson:"exec,omitempty" json:"exec,omitempty"`

	TicketInfo `bson:",inline"`

	DisplayGroup        string `bson:"dgroup,omitempty" json:"dgroup,omitempty"`
	InPbehaviorInterval bool   `bson:"in_pbh,omitempty" json:"in_pbh,omitempty"`
}

func (s *AlarmStep) GetInitiator() string {
	if s == nil {
		return ""
	}

	return s.Initiator
}

type TicketInfo struct {
	Ticket            string            `bson:"ticket,omitempty" json:"ticket,omitempty"`
	TicketURL         string            `bson:"ticket_url,omitempty" json:"ticket_url,omitempty"`
	TicketComment     string            `bson:"ticket_comment,omitempty" json:"ticket_comment,omitempty"`
	TicketSystemName  string            `bson:"ticket_system_name,omitempty" json:"ticket_system_name,omitempty"`
	TicketMetaAlarmID string            `bson:"ticket_meta_alarm_id,omitempty" json:"ticket_meta_alarm_id,omitempty"`
	TicketRuleID      string            `bson:"ticket_rule_id,omitempty" json:"ticket_rule_id,omitempty"`
	TicketRuleName    string            `bson:"ticket_rule_name,omitempty" json:"ticket_rule_name,omitempty"`
	TicketData        map[string]string `bson:"ticket_data,omitempty" json:"ticket_data,omitempty"`
}

func (t TicketInfo) GetStepMessage() string {
	builder := strings.Builder{}
	builder.WriteString(t.TicketRuleName)
	extraData := false
	if t.Ticket != "" {
		if t.TicketRuleName != "" {
			builder.WriteString(". ")
		}
		builder.WriteString("Ticket ID: ")
		builder.WriteString(t.Ticket)
		extraData = true
	}
	if t.TicketURL != "" {
		if t.TicketRuleName != "" || extraData {
			builder.WriteString(". ")
		}
		builder.WriteString("Ticket URL: ")
		builder.WriteString(t.TicketURL)
		extraData = true
	}
	for k, v := range t.TicketData {
		if t.TicketRuleName != "" || extraData {
			builder.WriteString(". ")
		}
		builder.WriteString("Ticket ")
		builder.WriteString(k)
		builder.WriteString(": ")
		builder.WriteString(v)
		extraData = true
	}
	if extraData {
		builder.WriteRune('.')
	}
	return builder.String()
}

// NewAlarmStep returns an AlarmStep.
// If the timestamp or author are empty, default values will be used to create an AlarmStep.
func NewAlarmStep(
	stepType string,
	timestamp datetime.CpsTime,
	author, msg, userID, role, initiator string,
	inPbehaviorInterval bool,
) AlarmStep {
	authorAlarmStep := author
	if authorAlarmStep == "" {
		authorAlarmStep = cps.DefaultEventAuthor
	}

	timestampAlarmStep := timestamp
	if timestampAlarmStep.IsZero() {
		timestampAlarmStep = datetime.NewCpsTime()
	}

	return AlarmStep{
		Author:              authorAlarmStep,
		UserID:              userID,
		Message:             msg,
		Timestamp:           timestampAlarmStep,
		Type:                stepType,
		Role:                role,
		Initiator:           initiator,
		InPbehaviorInterval: inPbehaviorInterval,
	}
}

// CropCounter provides an explicit way of counting the steps that were cropped.
type CropCounter struct {
	StateChanges  int `bson:"statechanges" json:"statechanges"`
	Stateinc      int `bson:"stateinc" json:"stateinc"`
	Statedec      int `bson:"statedec" json:"statedec"`
	StateInfo     int `bson:"state:0,omitempty" json:"state:0,omitempty"`
	StateMinor    int `bson:"state:1,omitempty" json:"state:1,omitempty"`
	StateMajor    int `bson:"state:2,omitempty" json:"state:2,omitempty"`
	StateCritical int `bson:"state:3,omitempty" json:"state:3,omitempty"`
}

// MergeCounter merges the current counter with the provided counter and returns the merged counter.
func (counter *CropCounter) MergeCounter(secondCounter CropCounter) {
	counter.StateChanges += secondCounter.StateChanges
	counter.Stateinc += secondCounter.Stateinc
	counter.Statedec += secondCounter.Statedec
	counter.StateInfo += secondCounter.StateInfo
	counter.StateMinor += secondCounter.StateMinor
	counter.StateMajor += secondCounter.StateMajor
	counter.StateCritical += secondCounter.StateCritical
}

// UpdateWithStep updates the CropCounter with the provided step informations.
func (counter *CropCounter) UpdateWithStep(step AlarmStep) {
	// Count the step types
	switch sType := step.Type; sType {
	case AlarmStepStateIncrease:
		counter.Stateinc += 1
	case AlarmStepStateDecrease:
		counter.Statedec += 1
	}

	// Count the step states
	switch state := step.Value; state {
	case 0:
		counter.StateInfo += 1
	case 1:
		counter.StateMinor += 1
	case 2:
		counter.StateMajor += 1
	case 3:
		counter.StateCritical += 1
	}

	counter.StateChanges++
}

func (counter CropCounter) IsZero() bool {
	return counter == CropCounter{}
}

func (counter *CropCounter) getMsg() string {
	msgBuilder := strings.Builder{}
	msgBuilder.WriteString("State increased: ")
	msgBuilder.WriteString(strconv.Itoa(counter.Stateinc))
	msgBuilder.WriteString("\n")
	msgBuilder.WriteString("State decreased: ")
	msgBuilder.WriteString(strconv.Itoa(counter.Statedec))
	msgBuilder.WriteString("\n")
	msgBuilder.WriteString("State changes: ")
	msgBuilder.WriteString(strconv.Itoa(counter.StateChanges))
	if counter.StateInfo > 0 {
		msgBuilder.WriteString("\n")
		msgBuilder.WriteString("State changes to ok: ")
		msgBuilder.WriteString(strconv.Itoa(counter.StateInfo))
	}

	if counter.StateMinor > 0 {
		msgBuilder.WriteString("\n")
		msgBuilder.WriteString("State changes to minor: ")
		msgBuilder.WriteString(strconv.Itoa(counter.StateMinor))
	}

	if counter.StateMajor > 0 {
		msgBuilder.WriteString("\n")
		msgBuilder.WriteString("State changes to major: ")
		msgBuilder.WriteString(strconv.Itoa(counter.StateMajor))
	}

	if counter.StateCritical > 0 {
		msgBuilder.WriteString("\n")
		msgBuilder.WriteString("State changes to critical: ")
		msgBuilder.WriteString(strconv.Itoa(counter.StateCritical))
	}

	return msgBuilder.String()
}

// AlarmSteps is a sortable implementation of []*AlarmStep. Used for sorting
// steps in some functions. Implements sort.Interface
type AlarmSteps []AlarmStep

// Add handle adding a step to the list
func (s *AlarmSteps) Add(step AlarmStep) error {
	if len(*s) < AlarmStepsHardLimit ||
		step.Type == AlarmStepStateDecrease && step.Value == AlarmStateOK ||
		step.Type == AlarmStepStatusDecrease && step.Value == AlarmStatusOff ||
		step.Type == AlarmStepStatusIncrease && step.Value == AlarmStatusCancelled ||
		step.Type == AlarmStepStatusDecrease && step.Value == AlarmStatusCancelled {
		*s = append(*s, step)
		return nil
	}

	return fmt.Errorf("max number of steps reached: %v", AlarmStepsHardLimit)
}

// Crop steps by replacing stateinc and statedec steps after the current status with a statecounter step
// Returns :
//   - the updated alarm steps
//   - True if it was updated, false else
//
// param currentStatus: the current status of the alarm. The steps will be cropped from this status
// param cropNum: crop only if we have at least cropNum steps with type AlarmStepStateIncrease or AlarmStepStateDecrease
func (s AlarmSteps) Crop(currentStatus *AlarmStep, cropNum int) (AlarmSteps, bool) {
	nbStepsToCrop := 0
	currentStatusIdx := -1

	// finding current status index
	// starting from the end as the slice is sorted from oldest to newest step
	for i := len(s) - 1; i >= 0 && currentStatusIdx < 0; i-- {
		step := s[i]
		if step.Type == AlarmStepStateIncrease || step.Type == AlarmStepStateDecrease {
			nbStepsToCrop += 1
		}
		if step.Type == currentStatus.Type && step.Timestamp.Time.Equal(currentStatus.Timestamp.Time) {
			currentStatusIdx = i
		}
	}

	if currentStatusIdx < 0 || nbStepsToCrop <= cropNum {
		return s, false
	}

	// Contains all the steps before those to be cropped
	cleanedSteps := s[:currentStatusIdx+1]

	counter := CropCounter{}
	for i := currentStatusIdx + 1; i < len(s); i++ {
		step := s[i]
		if nbStepsToCrop-counter.Stateinc-counter.Statedec <= cropNum {
			// Keeps cropNum most recent stateinc/dec steps
			// ! rewrites on s, but only rewrites on already processed steps !
			cleanedSteps = append(cleanedSteps, step)
		} else if step.Type == AlarmStepStateIncrease || step.Type == AlarmStepStateDecrease {
			counter.UpdateWithStep(step)
		} else {
			// We only add the step when it isn't a statedec/inc
			// ! rewrites on s, but only rewrites on already processed steps !
			cleanedSteps = append(cleanedSteps, step)
		}
	}

	cleanedSteps = cleanedSteps.UpdateStateCounter(currentStatus, currentStatusIdx, counter)

	return cleanedSteps, true
}

// UpdateStateCounter updates the alarm steps with the statecounter step
// Returns the updated AlarmSteps and the updated (or newly created) statecounter step.
// param currentStatus: the current status of the alarm. The statecounter step infos will come from it.
// param currentStatusIdx: the alarm current status' index. It is used to insert or update the statecounter step right after it.
// param counter: the crop counter to update or create the statecounter step from.
func (s AlarmSteps) UpdateStateCounter(currentStatus *AlarmStep, currentStatusIdx int, counter CropCounter) AlarmSteps {
	// index of where the statecounter is supposed to be, right after currStatus
	counterIdx := currentStatusIdx + 1

	if len(s) == counterIdx {
		// If last step is last change of status
		// create and append new statecounter
		newStep := AlarmStep{
			Author:              currentStatus.Author,
			UserID:              currentStatus.UserID,
			Initiator:           currentStatus.Initiator,
			Value:               currentStatus.Value,
			InPbehaviorInterval: currentStatus.InPbehaviorInterval,
			Timestamp:           datetime.NewCpsTime(),
			Type:                AlarmStepStateCounter,
			StateCounter:        CropCounter{},
		}

		s = append(s, newStep)
	} else if s[counterIdx].Type != AlarmStepStateCounter {
		// Else if the step just after the status isn't statecounter
		// create and insert new statecounter right after status
		newStep := AlarmStep{
			Author:              currentStatus.Author,
			UserID:              currentStatus.UserID,
			Initiator:           currentStatus.Initiator,
			Value:               currentStatus.Value,
			InPbehaviorInterval: currentStatus.InPbehaviorInterval,
			Timestamp:           datetime.NewCpsTime(),
			Type:                AlarmStepStateCounter,
			StateCounter:        CropCounter{},
		}

		// insert
		s = append(s, AlarmStep{})
		copy(s[counterIdx+1:], s[counterIdx:])
		s[counterIdx] = newStep
	}
	// Else the step just after current status is statecounter
	// And we don't have to do anything

	// update the existing statecounter
	s[counterIdx].StateCounter.MergeCounter(counter)
	s[counterIdx].Message = s[counterIdx].StateCounter.getMsg()

	return s
}

// Last returns the last step, if any, or returns an error.
func (s AlarmSteps) Last() (AlarmStep, error) {
	if len(s) > 0 {
		return s[len(s)-1], nil
	}
	return AlarmStep{}, fmt.Errorf("no step")
}

func (s AlarmSteps) Len() int      { return len(s) }
func (s AlarmSteps) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ByTimestamp struct {
	AlarmSteps
}

func (s ByTimestamp) Less(i, j int) bool {
	return s.AlarmSteps[i].Timestamp.Before(s.AlarmSteps[j].Timestamp)
}

// PbehaviorInfo represents current state of entity.
type PbehaviorInfo struct {
	// Timestamp is time when entity enters pbehavior.
	// Use pointer of CpsTime to unmarshal null and undefined to nil pointer instead of zero CpsTime.
	Timestamp *datetime.CpsTime `bson:"timestamp" json:"timestamp" swaggertype:"integer"`
	// ID is ID of pbehavior.PBehavior.
	ID string `bson:"id" json:"id"`
	// Name is Name of pbehavior.PBehavior.
	Name string `bson:"name" json:"name"`
	// ReasonName is Name of pbehavior.Reason.
	ReasonName string `bson:"reason_name" json:"reason_name"`
	// ReasonID is ID of pbehavior.Reason.
	ReasonID string `bson:"reason" json:"reason"`
	// TypeID is ID of pbehavior.Type.
	TypeID string `bson:"type" json:"type"`
	// TypeName is Name of pbehavior.Type.
	TypeName string `bson:"type_name" json:"type_name"`
	// CanonicalType is Type of pbehavior.Type.
	CanonicalType string `bson:"canonical_type" json:"canonical_type"`
}

func (i *PbehaviorInfo) IsDefaultActive() bool {
	return i.ID == ""
}

func (i *PbehaviorInfo) IsActive() bool {
	return i.CanonicalType == "" || i.CanonicalType == PbhCanonicalTypeActive
}

func (i *PbehaviorInfo) Is(t string) bool {
	if i.CanonicalType == t {
		return true
	}

	if i.CanonicalType == "" && t == PbhCanonicalTypeActive {
		return true
	}

	return false
}

func (i *PbehaviorInfo) OneOf(t []string) bool {
	for _, v := range t {
		if i.Is(v) {
			return true
		}
	}

	return false
}

func (i PbehaviorInfo) IsZero() bool {
	return i == PbehaviorInfo{}
}

func (i PbehaviorInfo) Same(v PbehaviorInfo) bool {
	v.Timestamp = i.Timestamp

	return i == v
}

// GetStringField is a magic getter for string fields for easier field retrieving when matching pbehavior info pattern
func (i *PbehaviorInfo) GetStringField(f string) (string, bool) {
	switch f {
	case "pbehavior_info.id":
		return i.ID, true
	case "pbehavior_info.type":
		return i.TypeID, true
	case "pbehavior_info.canonical_type":
		if i.CanonicalType == "" {
			return PbhCanonicalTypeActive, true
		}
		return i.CanonicalType, true
	case "pbehavior_info.reason":
		return i.ReasonID, true
	default:
		return "", false
	}
}

// AlarmValue represents a full description of an alarm.
type AlarmValue struct {
	ACK         *AlarmStep  `bson:"ack,omitempty" json:"ack,omitempty"`
	Canceled    *AlarmStep  `bson:"canceled,omitempty" json:"canceled,omitempty"`
	Snooze      *AlarmStep  `bson:"snooze,omitempty" json:"snooze,omitempty"`
	State       *AlarmStep  `bson:"state,omitempty" json:"state,omitempty"`
	Status      *AlarmStep  `bson:"status,omitempty" json:"status,omitempty"`
	LastComment *AlarmStep  `bson:"last_comment,omitempty" json:"last_comment,omitempty"`
	ChangeState *AlarmStep  `bson:"change_state,omitempty" json:"change_state,omitempty"`
	Tickets     []AlarmStep `bson:"tickets,omitempty" json:"tickets,omitempty"`
	// Ticket contains the last created ticket
	Ticket *AlarmStep `bson:"ticket,omitempty" json:"ticket,omitempty"`
	Steps  AlarmSteps `bson:"steps" json:"steps"`

	Component         string            `bson:"component" json:"component"`
	Connector         string            `bson:"connector" json:"connector"`
	ConnectorName     string            `bson:"connector_name" json:"connector_name"`
	CreationDate      datetime.CpsTime  `bson:"creation_date" json:"creation_date"`
	ActivationDate    *datetime.CpsTime `bson:"activation_date,omitempty" json:"activation_date,omitempty"`
	DisplayName       string            `bson:"display_name" json:"display_name"`
	HardLimit         *CpsNumber        `bson:"hard_limit,omitempty" json:"hard_limit,omitempty"`
	InitialOutput     string            `bson:"initial_output" json:"initial_output"`
	Output            string            `bson:"output" json:"output"`
	InitialLongOutput string            `bson:"initial_long_output" json:"initial_long_output"`
	LongOutput        string            `bson:"long_output" json:"long_output"`
	LongOutputHistory []string          `bson:"long_output_history" json:"long_output_history"`
	LastUpdateDate    datetime.CpsTime  `bson:"last_update_date" json:"last_update_date"`
	LastEventDate     datetime.CpsTime  `bson:"last_event_date" json:"last_event_date"`
	Resource          string            `bson:"resource,omitempty" json:"resource,omitempty"`
	Resolved          *datetime.CpsTime `bson:"resolved,omitempty" json:"resolved,omitempty"`
	PbehaviorInfo     PbehaviorInfo     `bson:"pbehavior_info,omitempty" json:"pbehavior_info,omitempty"`
	Meta              string            `bson:"meta,omitempty" json:"meta,omitempty"`
	MetaValuePath     string            `bson:"meta_value_path,omitempty" json:"meta_value_path,omitempty"`

	Parents         []string `bson:"parents" json:"parents"`
	Children        []string `bson:"children" json:"children"`
	UnlinkedParents []string `bson:"unlinked_parents" json:"unlinked_parents"`

	StateChangesSinceStatusUpdate CpsNumber `bson:"state_changes_since_status_update,omitempty" json:"state_changes_since_status_update,omitempty"`
	TotalStateChanges             CpsNumber `bson:"total_state_changes,omitempty" json:"total_state_changes,omitempty"`
	// EventsCount accumulates count of check events.
	EventsCount CpsNumber `bson:"events_count,omitempty" json:"events_count,omitempty"`

	Infos map[string]map[string]interface{} `bson:"infos" json:"infos"`

	// store version of dynamic-infos rule
	RuleVersion map[string]string `bson:"infos_rule_version" json:"infos_rule_version"`

	// InactiveStart represents start of snooze or maintenance, pause, inactive pbehavior interval.
	// It's used only to compute InactiveDuration.
	InactiveStart *datetime.CpsTime `bson:"inactive_start,omitempty" json:"inactive_start"`
	// Duration represents a duration from creation date to resolve date.
	// Keep omitempty.
	Duration int64 `bson:"duration,omitempty" json:"duration"`
	// CurrentStateDuration represents a duration when an alarm was in current state.
	// Keep omitempty.
	CurrentStateDuration int64 `bson:"current_state_duration,omitempty" json:"current_state_duration"`
	// ActiveDuration represents a duration when an alarm wasn't in snooze or in maintenance, pause, inactive pbehavior interval.
	// Keep omitempty.
	ActiveDuration int64 `bson:"active_duration,omitempty" json:"active_duration"`
	// InactiveDuration represents a duration when an alarm was in snooze or in maintenance, pause, inactive pbehavior interval.
	InactiveDuration int64 `bson:"inactive_duration" json:"inactive_duration"`
	// SnoozeDuration represents a duration when an alarm was in snooze.
	SnoozeDuration int64 `bson:"snooze_duration" json:"snooze_duration"`
	// PbehaviorInactiveDuration represents a duration when an alarm was in maintenance, pause, inactive pbehavior interval.
	PbehaviorInactiveDuration int64 `bson:"pbh_inactive_duration" json:"pbh_inactive_duration"`
}

func (v *AlarmValue) Transform() {
	if v.Resolved != nil && v.Resolved.Unix() == 0 {
		v.Resolved = nil
	}
}
