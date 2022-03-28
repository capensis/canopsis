package types

//go:generate easyjson -no_std_marshalers

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
)

//event initiators
const (
	InitiatorUser     = "user"
	InitiatorSystem   = "system"
	InitiatorExternal = "external"
)

//Source types
const (
	SourceTypeResource  = "resource"
	SourceTypeComponent = "component"
	SourceTypeConnector = "connector"
	SourceTypeService   = "service"
	SourceTypeMetaAlarm = "metaalarm"
)

// Event types.
// Add each new event type to isValidEventType func.
const (
	EventTypeAck         = "ack"
	EventTypeAckremove   = "ackremove"
	EventTypeAssocTicket = "assocticket"
	EventTypeCancel      = "cancel"
	EventTypeCheck       = "check"
	EventTypeComment     = "comment"
	// EventTypeDeclareTicket is used for manual declareticket trigger which is designed
	// to trigger webhook with declare ticket parameter.
	EventTypeDeclareTicket = "declareticket"
	// EventTypeDeclareTicketWebhook is triggered after declare ticket creation by webhook.
	EventTypeDeclareTicketWebhook = "declareticketwebhook"

	EventTypeDone              = "done"
	EventTypeChangestate       = "changestate"
	EventTypeKeepstate         = "keepstate"
	EventTypePerf              = "perf"
	EventTypeSnooze            = "snooze"
	EventTypeUnsnooze          = "unsnooze"
	EventTypeStatCounterInc    = "statcounterinc"
	EventTypeStatDuration      = "statduration"
	EventTypeStatStateInterval = "statstateinterval"
	EventTypeUncancel          = "uncancel"

	EventTypeMetaAlarm          = "metaalarm"
	EventTypeMetaAlarmUpdated   = "metaalarmupdated"
	EventTypePbhEnter           = "pbhenter"
	EventTypePbhLeaveAndEnter   = "pbhleaveandenter"
	EventTypePbhLeave           = "pbhleave"
	EventTypePbhCreate          = "pbhcreate"
	EventTypeResolveDone        = "resolve_done"
	EventTypeResolveCancel      = "resolve_cancel"
	EventTypeResolveClose       = "resolve_close"
	EventTypeUpdateStatus       = "updatestatus"
	EventManualMetaAlarmGroup   = "manual_metaalarm_group"
	EventManualMetaAlarmUngroup = "manual_metaalarm_ungroup"
	EventManualMetaAlarmUpdate  = "manual_metaalarm_update"
	EventTypeActivate           = "activate"
	EventTypeRunDelayedScenario = "run_delayed_scenario"

	// Following event types are used to add manual instruction execution to alarm steps.
	EventTypeInstructionStarted   = "instructionstarted"
	EventTypeInstructionPaused    = "instructionpaused"
	EventTypeInstructionResumed   = "instructionresumed"
	EventTypeInstructionCompleted = "instructioncompleted"
	EventTypeInstructionFailed    = "instructionfailed"
	// EventTypeInstructionAborted is the same for manual and auto instructions.
	EventTypeInstructionAborted = "instructionaborted"
	// Following event types are used to add auto instruction execution to alarm steps.
	EventTypeAutoInstructionStarted        = "autoinstructionstarted"
	EventTypeAutoInstructionCompleted      = "autoinstructioncompleted"
	EventTypeAutoInstructionFailed         = "autoinstructionfailed"
	EventTypeAutoInstructionAlreadyRunning = "autoinstructionalreadyrunning"
	// Following event types are used to add job execution to alarm steps. Events are
	// the same for manual and auto instructions.
	EventTypeInstructionJobStarted   = "instructionjobstarted"
	EventTypeInstructionJobCompleted = "instructionjobcompleted"
	EventTypeInstructionJobAborted   = "instructionjobaborted"
	EventTypeInstructionJobFailed    = "instructionjobfailed"

	// EventTypeRecomputeEntityService is used to recompute service context graph and state.
	EventTypeRecomputeEntityService = "recomputeentityservice"
	// EventTypeUpdateEntityService is used to update service cache in engines.
	EventTypeUpdateEntityService = "updateentityservice"
	// EventTypeEntityUpdated is used to notify engines that entity is updated out of
	// event flow.
	EventTypeEntityUpdated = "entityupdated"
	// EventTypeEntityToggled is used to notify engines that entity is enabled/disabled.
	EventTypeEntityToggled = "entitytoggled"
	// EventTypeAlarmSkipped is used to check alarm in service counters if alarm was skipped
	// during service recompute.
	EventTypeAlarmSkipped = "alarmskipped"
	// EventTypeJunitTestSuiteUpdated is used to notify that test suite is updated but state is not changed.
	EventTypeJunitTestSuiteUpdated = "junittestsuiteupdated"
	// EventTypeJunitTestCaseUpdated is used to notify that test case is updated but state is not changed.
	EventTypeJunitTestCaseUpdated = "junittestcaseeupdated"

	EventTypeStateIncrease  = "stateinc"
	EventTypeStateDecrease  = "statedec"
	EventTypeStatusIncrease = "statusinc"
	EventTypeStatusDecrease = "statusdec"

	// EventTypeNoEvents is used to create alarm for entity by idle rule.
	EventTypeNoEvents = "noevents"
)

const (
	ConnectorEngineService = "service"
	ConnectorJunit         = "junit"
)

const MaxEventTimestampVariation = 24 * time.Hour

//PerfData represents a perf data array
type PerfData struct {
	Metric string  `bson:"metric" json:"metric"`
	Unit   string  `bson:"unit" json:"unit"`
	Value  float64 `bson:"value" json:"value"`
}

// Event represents a canopsis event.
//easyjson:json
type Event struct {
	ID            *string    `bson:"_id" json:"_id"`
	Connector     string     `bson:"connector" json:"connector"`
	ConnectorName string     `bson:"connector_name" json:"connector_name"`
	EventType     string     `bson:"event_type" json:"event_type"`
	Component     string     `bson:"component" json:"component"`
	Resource      string     `bson:"resource" json:"resource"`
	PerfData      *string    `bson:"perf_data" json:"perf_data"`
	PerfDataArray []PerfData `bson:"perf_data_array" json:"perf_data_array"`
	Status        *CpsNumber `bson:"status" json:"status"`
	Timestamp     CpsTime    `bson:"timestamp" json:"timestamp"`
	StateType     *CpsNumber `bson:"state_type" json:"state_type"`
	SourceType    string     `bson:"source_type" json:"source_type"`
	LongOutput    string     `bson:"long_output" json:"long_output"`
	State         CpsNumber  `bson:"state" json:"state"`
	Output        string     `bson:"output" json:"output"`
	Alarm         *Alarm     `bson:"current_alarm" json:"current_alarm"`
	Entity        *Entity    `bson:"current_entity" json:"current_entity"`

	Author string `bson:"author" json:"author"`
	UserID string `bson:"user_id" json:"user_id"`

	RK string `bson:"routing_key" json:"routing_key"`
	// AckResources is used to ack all resource alarms on ack component alarm.
	// It also adds declare ticket to all resource alarms on ack webhook.
	// It's still used by some old users but meta alarms must be used instead.
	AckResources bool                   `json:"ack_resources"`
	Duration     CpsNumber              `json:"duration"`
	Ticket       string                 `bson:"ticket" json:"ticket"`
	StatName     string                 `bson:"stat_name" json:"stat_name"`
	Debug        bool                   `bson:"debug" json:"debug"`
	Role         string                 `bson:"role,omitempty" json:"role,omitempty"`
	ExtraInfos   map[string]interface{} `json:"extra"`
	AlarmChange  *AlarmChange           `bson:"alarm_change" json:"alarm_change"`

	MetaAlarmRuleID    string `bson:"metaalarm_rule_id" json:"metaalarm_rule_id"`
	MetaAlarmValuePath string `bson:"metaalarm_value_path" json:"metaalarm_value_path"`

	MetaAlarmRelatedParents []string  `bson:"ma_related_parents" json:"ma_related_parents"`
	MetaAlarmParents        *[]string `bson:"ma_parents" json:"ma_parents"`
	MetaAlarmChildren       *[]string `bson:"ma_children" json:"ma_children"`

	PbehaviorInfo PbehaviorInfo `bson:"pbehavior_info" json:"pbehavior_info"`

	// Initiator is used to detect who emits event.
	// InitiatorUser - UI
	// InitiatorSystem - engines
	// InitiatorExternal - third tool
	Initiator string `bson:"initiator" json:"initiator"`

	// Only for EventTypeRunDelayedScenario
	DelayedScenarioID string `bson:"delayed_scenario_id,omitempty" json:"delayed_scenario_id,omitempty"`

	// AddedToServices contains ids of entity services to which entity has been added as dependency.
	AddedToServices []string `bson:"added_to_services,omitempty" json:"added_to_services,omitempty"`
	// RemovedFromServices contains ids of entity services from which entity has been removed as dependency.
	RemovedFromServices []string `bson:"removed_from_services,omitempty" json:"removed_from_services,omitempty"`

	// PbhParameters is used only with EventTypePbhCreate
	PbhParameters string `bson:"pbh_parameters,omitempty" json:"pbh_parameters,omitempty"`

	// IdleRuleApply is used if event is emitted by idle rule.
	IdleRuleApply string `bson:"idle_rule_apply,omitempty" json:"idle_rule_apply,omitempty"`

	// Execution is used only for instruction events: EventTypeInstructionStarted, EventTypeInstructionCompleted, etc..
	Execution string `bson:"execution,omitempty" json:"execution,omitempty"`
}

// ContextInformation regroup context values necessary for creating a new entity
type ContextInformation struct {
	ID      string // Entity ID
	Name    string
	Type    string
	Impacts []string
	Depends []string
}

// NewEventFromJSON create an Event from a raw json
func NewEventFromJSON(body []byte) (Event, error) {
	var event Event
	err := json.Unmarshal(body, &event)
	if err != nil {
		return event, fmt.Errorf("NewEvent error: %v", err)
	}
	event.Format()

	return event, nil
}

func NewEventFromAlarm(alarm Alarm) Event {
	event := Event{
		Connector:     alarm.Value.Connector,
		ConnectorName: alarm.Value.ConnectorName,
		Resource:      alarm.Value.Resource,
		Component:     alarm.Value.Component,
		Alarm:         &alarm,
		Timestamp:     CpsTime{Time: time.Now()},
	}
	event.SourceType = event.DetectSourceType()

	return event
}

// Format an event
//  "timestamp" is fill with time.Now()
//  "event_type" is fill with EventTypeCheck
//  if "entity" is not null, "impacts" and "depends" are ensured to be initialized
func (e *Event) Format() {
	//events can't be later or earlier than MaxEventTimestampVariation
	now := NewCpsTime()
	if e.Timestamp.IsZero() || e.Timestamp.Time.Before(now.Add(-MaxEventTimestampVariation)) || e.Timestamp.Time.After(now.Add(MaxEventTimestampVariation)) {
		e.Timestamp = now
	}
	if e.EventType == "" {
		e.EventType = EventTypeCheck
	}
	if e.StateType == nil {
		statetype := CpsNumber(1)
		e.StateType = &statetype
	}
	if e.Initiator == "" {
		e.Initiator = InitiatorExternal
	}

	if e.Entity != nil {
		e.Entity.EnsureInitialized()
	}
}

// GetEID generates an eid from an event
func (e *Event) GetEID() string {
	if e.Resource != "" {
		return e.Resource + "/" + e.Component
	}
	if e.Component != "" {
		return e.Component
	}
	return e.Connector + "/" + e.ConnectorName
}

// GetLockID returns lock name that used to block alarm
func (e *Event) GetLockID() string {
	return e.GetEID()
}

// InjectExtraInfos takes the raw JSON event document and puts any unknown
// field into Event.ExtraInfos
func (e *Event) InjectExtraInfos(source []byte) error {
	if len(e.ExtraInfos) > 0 {
		return nil
	}

	unmatchedParams := make(map[string]interface{})
	if err := json.Unmarshal(source, &unmatchedParams); err == nil {
		for _, jsonKey := range e.GetRequiredKeys() {
			delete(unmatchedParams, jsonKey)
		}
	} else {
		return fmt.Errorf("Event.InjectExtraInfos json decode: %v", err)
	}
	if e.ExtraInfos == nil {
		e.ExtraInfos = make(map[string]interface{})
	}

	for k, v := range unmatchedParams {
		e.ExtraInfos[k] = v
	}

	return nil
}

// IsContextable tells you if the given event can lead to context enrichment.
func (e Event) IsContextable() bool {
	switch e.EventType {
	case EventTypeCheck, EventTypePerf, EventTypeDeclareTicket, EventTypeMetaAlarm,
		EventTypeEntityToggled, EventTypeEntityUpdated:
		return true
	default:
		return false
	}
}

// IsMatched tell if an event is catched by a regex
func (e Event) IsMatched(regex string, fields []string) bool {
	for _, fieldName := range fields {
		field := utils.GetStringField(e, fieldName)
		matched, _ := regexp.MatchString(regex, field)
		if matched {
			return true
		}
	}
	return false
}

// IsValid checks if an Event is valid for Canopsis processing.
// the error returned, if any, is of type errt.UnknownError
func (e Event) IsValid() error {
	if e.Connector == "" || e.ConnectorName == "" {
		return errt.NewUnknownError(errors.New("missing connector"))
	}

	switch e.SourceType {
	case SourceTypeConnector:
		/*do nothing*/
	case SourceTypeComponent, SourceTypeMetaAlarm, SourceTypeService:
		if e.Component == "" {
			return errt.NewUnknownError(errors.New("missing component"))
		}
	case SourceTypeResource:
		if e.Resource == "" {
			return errt.NewUnknownError(errors.New("missing resource"))
		}
	default:
		return errt.NewUnknownError(fmt.Errorf("wrong source type: %v", e.SourceType))
	}

	if !isValidEventType(e.EventType) {
		return errt.NewUnknownError(fmt.Errorf("wrong event type: %v", e.EventType))
	}

	switch e.EventType {
	case EventTypePerf:
		if e.PerfData == nil || e.PerfDataArray == nil {
			return errt.NewUnknownError(errors.New("perfdata without data"))
		}
	}

	return nil
}

func (e *Event) DetectSourceType() string {
	if e.Resource != "" {
		return SourceTypeResource
	}

	if e.Component != "" {
		return SourceTypeComponent
	}

	return SourceTypeConnector
}

// GenericEvent contains an interface so you can do this:
// body := `<a json document>`
// var gevent GenericEvent
// json.Unmarshal(body, &gevent.Content)
// gevent.PartialID(<rules>)
type GenericEvent struct {
	Content interface{}
}

// JSONUnmarshal is a shortcut for this:
// var event GenericEvent
// json.Unmarshal(body, &event.Content)
func (e *GenericEvent) JSONUnmarshal(body []byte) error {
	return json.Unmarshal(body, &e.Content)
}

// GetCompatRK returns the event routing key. For compatibility only with old engines.
func (e *Event) GetCompatRK() string {
	if e.SourceType == SourceTypeResource {
		return fmt.Sprintf(
			"%s.%s.%s.%s.%s.%s",
			e.Connector,
			e.ConnectorName,
			e.EventType,
			e.SourceType,
			e.Component,
			e.Resource,
		)
	}
	return fmt.Sprintf(
		"%s.%s.%s.%s.%s",
		e.Connector,
		e.ConnectorName,
		e.EventType,
		e.SourceType,
		e.Component,
	)
}

// GetRequiredKeys read all declared json tags in the struct
func (e *Event) GetRequiredKeys() []string {
	var values []string

	typeof := reflect.TypeOf(e).Elem()

	for i := 0; i < typeof.NumField(); i++ {
		field := typeof.Field(i)
		tag := field.Tag.Get("json")

		values = append(values, tag)
	}

	return values
}

// GenerateContextInformations generate connector, component and resource entity informations.
// First element is always connector
// Second element is always component
// Third element doesnt exists if event is not SourceTypeResource, otherwise it is resource
func (e *Event) GenerateContextInformations() []ContextInformation {
	connectorID := e.Connector + "/" + e.ConnectorName
	componentID := e.Component

	connector := ContextInformation{
		ID:      connectorID,
		Name:    e.ConnectorName,
		Type:    EntityTypeConnector,
		Impacts: []string{},
		Depends: []string{componentID},
	}

	component := ContextInformation{
		ID:      componentID,
		Name:    e.Component,
		Type:    EntityTypeComponent,
		Impacts: []string{connectorID},
		Depends: []string{},
	}

	if e.SourceType != SourceTypeResource {
		return []ContextInformation{connector, component}
	}

	resourceID := e.Resource + "/" + e.Component

	connector.Impacts = append(connector.Impacts, resourceID)
	component.Depends = append(component.Depends, resourceID)

	resource := ContextInformation{
		ID:      resourceID,
		Name:    e.Resource,
		Type:    EntityTypeResource,
		Impacts: []string{componentID},
		Depends: []string{connectorID},
	}

	return []ContextInformation{connector, component, resource}
}

// NewEntity create an entity from ContextInformation struct
func (ci *ContextInformation) NewEntity() Entity {
	return NewEntity(
		ci.ID,
		ci.Name,
		ci.Type,
		map[string]Info{},
		ci.Impacts,
		ci.Depends,
	)
}

var cpsNumberType = reflect.TypeOf(CpsNumber(0))
var cpsNumberPtrType = reflect.PtrTo(cpsNumberType)
var cpsTimeType = reflect.TypeOf(CpsTime{})
var stringType = reflect.TypeOf("")
var stringPtrType = reflect.PtrTo(stringType)
var boolType = reflect.TypeOf(false)
var entityPtrType = reflect.PtrTo(reflect.TypeOf(Entity{}))

// SetField sets the value of a field of an event given its name.
func (e *Event) SetField(name string, value interface{}) (err error) {
	// Recover from panics at the end of the function and returns an error
	// instead. The code below should not panic, but this prevents the engines
	// from crashing if there is a mistake in the use of the functions of the
	// reflect packages (which panic when misused)
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %+v", r)
		}
	}()

	field := reflect.ValueOf(e).Elem().FieldByName(name)
	if !field.IsValid() {
		// The field does not exist, add the value to ExtraInfos.
		e.ExtraInfos[name] = value
		return nil
	}

	if !field.CanSet() {
		return fmt.Errorf("field %s cannot be set", name)
	}

	// For each possible type of field, try to convert the value to this type,
	// and assign it.
	switch field.Type() {
	case cpsNumberType:
		integerValue, success := AsInteger(value)
		if !success {
			return fmt.Errorf("value cannot be converted to an integer: %+v", value)
		}
		field.Set(reflect.ValueOf(CpsNumber(integerValue)))

	case cpsNumberPtrType:
		integerValue, success := AsInteger(value)
		if !success {
			return fmt.Errorf("value cannot be converted to an integer: %+v", value)
		}
		cpsNumberValue := CpsNumber(integerValue)
		field.Set(reflect.ValueOf(&cpsNumberValue))

	case cpsTimeType:
		integerValue, success := AsInteger(value)
		if !success {
			return fmt.Errorf("value cannot be converted to an integer: %+v", value)
		}
		cpsTimeValue := CpsTime{Time: time.Unix(integerValue, 0)}
		field.Set(reflect.ValueOf(cpsTimeValue))

	case stringType:
		stringValue, success := utils.AsString(value)
		if !success {
			return fmt.Errorf("value cannot be assigned to a string: %+v", value)
		}
		field.Set(reflect.ValueOf(stringValue))

	case stringPtrType:
		stringValue, success := utils.AsString(value)
		if !success {
			return fmt.Errorf("value cannot be assigned to a string: %+v", value)
		}
		field.Set(reflect.ValueOf(&stringValue))

	case boolType:
		boolValue, success := value.(bool)
		if !success {
			return fmt.Errorf("value cannot be assigned to a bool: %+v", value)
		}
		field.Set(reflect.ValueOf(boolValue))

	case entityPtrType:
		entityValue, success := value.(Entity)
		if !success {
			return fmt.Errorf("value cannot be assigned to an entity: %+v", value)
		}
		field.Set(reflect.ValueOf(&entityValue))

	default:
		return fmt.Errorf("cannot set field %s of type %v", name, field.Type())
	}

	return nil
}

func (e *Event) IsPbehaviorEvent() bool {
	return e.EventType == EventTypePbhEnter ||
		e.EventType == EventTypePbhLeave ||
		e.EventType == EventTypePbhLeaveAndEnter
}

func isValidEventType(t string) bool {
	switch t {
	case EventTypeCheck,
		EventTypeActivate,
		EventTypeAck,
		EventTypeAckremove,
		EventTypeAssocTicket,
		EventTypeCancel,
		EventTypeComment,
		EventTypeDeclareTicket,
		EventTypeDeclareTicketWebhook,
		EventTypeDone,
		EventTypeChangestate,
		EventTypeSnooze,
		EventTypeUnsnooze,
		EventTypeUncancel,
		EventTypeResolveDone,
		EventTypeResolveCancel,
		EventTypeResolveClose,
		EventTypePbhEnter,
		EventTypePbhLeaveAndEnter,
		EventTypePbhLeave,
		EventTypePbhCreate,
		EventTypeUpdateStatus,
		EventTypeMetaAlarm,
		EventTypeMetaAlarmUpdated,
		EventManualMetaAlarmGroup,
		EventManualMetaAlarmUngroup,
		EventManualMetaAlarmUpdate,
		EventTypeRecomputeEntityService,
		EventTypeUpdateEntityService,
		EventTypeEntityUpdated,
		EventTypeEntityToggled,
		EventTypeNoEvents,
		EventTypeRunDelayedScenario,
		EventTypeInstructionStarted,
		EventTypeInstructionPaused,
		EventTypeInstructionResumed,
		EventTypeInstructionCompleted,
		EventTypeInstructionFailed,
		EventTypeInstructionAborted,
		EventTypeAutoInstructionStarted,
		EventTypeAutoInstructionCompleted,
		EventTypeAutoInstructionFailed,
		EventTypeAutoInstructionAlreadyRunning,
		EventTypeInstructionJobStarted,
		EventTypeInstructionJobCompleted,
		EventTypeInstructionJobAborted,
		EventTypeInstructionJobFailed,
		EventTypeAlarmSkipped,
		EventTypeJunitTestSuiteUpdated,
		EventTypeJunitTestCaseUpdated,
		EventTypeKeepstate,
		EventTypePerf,
		EventTypeStateIncrease,
		EventTypeStateDecrease,
		EventTypeStatusIncrease,
		EventTypeStatusDecrease:
		return true
	}

	return false
}
