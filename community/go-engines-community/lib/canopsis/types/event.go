package types

//go:generate easyjson -no_std_marshalers

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
)

// event initiators
const (
	InitiatorUser     = "user"
	InitiatorSystem   = "system"
	InitiatorExternal = "external"
)

// Source types
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
	EventTypeChangestate = "changestate"
	EventTypeSnooze      = "snooze"
	EventTypeUnsnooze    = "unsnooze"
	EventTypeUncancel    = "uncancel"

	EventTypeDeclareTicketWebhook = "declareticketwebhook"
	EventTypeWebhookStarted       = "webhookstarted"
	EventTypeWebhookCompleted     = "webhookcompleted"
	EventTypeWebhookFailed        = "webhookfailed"
	EventTypeAutoWebhookStarted   = "autowebhookstarted"
	EventTypeAutoWebhookCompleted = "autowebhookcompleted"
	EventTypeAutoWebhookFailed    = "autowebhookfailed"

	EventTypePbhEnter           = "pbhenter"
	EventTypePbhLeaveAndEnter   = "pbhleaveandenter"
	EventTypePbhLeave           = "pbhleave"
	EventTypeResolveCancel      = "resolve_cancel"
	EventTypeResolveClose       = "resolve_close"
	EventTypeResolveDeleted     = "resolve_deleted"
	EventTypeUpdateStatus       = "updatestatus"
	EventTypeActivate           = "activate"
	EventTypeRunDelayedScenario = "run_delayed_scenario"

	EventTypeMetaAlarm               = "metaalarm"
	EventTypeMetaAlarmAttachChildren = "metaalarmattachchildren"
	EventTypeMetaAlarmDetachChildren = "metaalarmdetachchildren"
	EventTypeMetaAlarmUngroup        = "metaalarm_ungroup"
	EventTypeManualMetaAlarmGroup    = "manual_metaalarm_group"
	EventTypeManualMetaAlarmUngroup  = "manual_metaalarm_ungroup"
	EventTypeManualMetaAlarmUpdate   = "manual_metaalarm_update"

	// Following event types are used to add manual instruction execution to alarm steps.
	EventTypeInstructionStarted   = "instructionstarted"
	EventTypeInstructionPaused    = "instructionpaused"
	EventTypeInstructionResumed   = "instructionresumed"
	EventTypeInstructionCompleted = "instructioncompleted"
	EventTypeInstructionFailed    = "instructionfailed"
	// EventTypeInstructionAborted is the same for manual and auto instructions.
	EventTypeInstructionAborted = "instructionaborted"
	// Following event types are used to add auto instruction execution to alarm steps.
	EventTypeAutoInstructionStarted   = "autoinstructionstarted"
	EventTypeAutoInstructionCompleted = "autoinstructioncompleted"
	EventTypeAutoInstructionFailed    = "autoinstructionfailed"
	// Following event types are used to add job execution to alarm steps. Events are
	// the same for manual and auto instructions.
	EventTypeInstructionJobStarted   = "instructionjobstarted"
	EventTypeInstructionJobCompleted = "instructionjobcompleted"
	EventTypeInstructionJobFailed    = "instructionjobfailed"

	// EventTypeRecomputeEntityService is used to recompute service context graph and state.
	EventTypeRecomputeEntityService = "recomputeentityservice"
	// EventTypeEntityUpdated is used to notify engines that entity is updated out of
	// event flow.
	EventTypeEntityUpdated = "entityupdated"
	// EventTypeEntityToggled is used to notify engines that entity is enabled/disabled.
	EventTypeEntityToggled = "entitytoggled"

	EventTypeUpdateCounters = "updatecounters"
	// EventTypeJunitTestSuiteUpdated is used to notify that test suite is updated but state is not changed.
	EventTypeJunitTestSuiteUpdated = "junittestsuiteupdated"
	// EventTypeJunitTestCaseUpdated is used to notify that test case is updated but state is not changed.
	EventTypeJunitTestCaseUpdated = "junittestcaseeupdated"
	// EventTypeNoEvents is used to create alarm for entity by idle rule.
	EventTypeNoEvents = "noevents"
	// EventTypeTrigger is used in axe rpc to send auto and manual instruction triggers
	EventTypeTrigger = "trigger"
	// EventTypeAutoInstructionActivate is used to activate alarm when an autoremediation triggered by create trigger is completed
	EventTypeAutoInstructionActivate = "autoinstructionactivate"
)

const (
	ConnectorEngineService = "service"
	ConnectorJunit         = "junit"
	ConnectorApi           = "api"
)

const MaxEventTimestampVariation = 24 * time.Hour

// Event represents a canopsis event.
//
//easyjson:json
type Event struct {
	ID            *string    `bson:"_id" json:"_id"`
	Connector     string     `bson:"connector" json:"connector"`
	ConnectorName string     `bson:"connector_name" json:"connector_name"`
	EventType     string     `bson:"event_type" json:"event_type"`
	Component     string     `bson:"component" json:"component"`
	Resource      string     `bson:"resource" json:"resource"`
	PerfData      string     `bson:"perf_data" json:"perf_data"`
	Status        *CpsNumber `bson:"status" json:"status"`
	SourceType    string     `bson:"source_type" json:"source_type"`
	LongOutput    string     `bson:"long_output" json:"long_output"`
	State         CpsNumber  `bson:"state" json:"state"`
	Output        string     `bson:"output" json:"output"`
	Alarm         *Alarm     `bson:"current_alarm" json:"current_alarm"`
	Entity        *Entity    `bson:"current_entity" json:"current_entity"`

	// AlarmID is used if an event is emitted for the specific alarm.
	AlarmID string `bson:"aid,omitempty" json:"aid,omitempty"`

	Author string `bson:"author" json:"author"`
	UserID string `bson:"user_id" json:"user_id"`

	Timestamp         datetime.CpsTime   `bson:"timestamp" json:"timestamp"`
	ReceivedTimestamp datetime.MicroTime `bson:"rt" json:"rt"`

	RK          string                 `bson:"routing_key" json:"routing_key"`
	Duration    CpsNumber              `bson:"duration,omitempty" json:"duration,omitempty"`
	StatName    string                 `bson:"stat_name" json:"stat_name"`
	Debug       bool                   `bson:"debug" json:"debug"`
	Role        string                 `bson:"role,omitempty" json:"role,omitempty"`
	ExtraInfos  map[string]interface{} `bson:"extra_infos" json:"extra"`
	AlarmChange *AlarmChange           `bson:"alarm_change" json:"alarm_change"`

	// Ticket related fields
	TicketInfo `bson:",inline"`

	// Tags contains external tags for alarm.
	Tags map[string]string `bson:"tags" json:"tags"`

	MetaAlarmRuleID    string `bson:"metaalarm_rule_id,omitempty" json:"metaalarm_rule_id,omitempty"`
	MetaAlarmValuePath string `bson:"metaalarm_value_path,omitempty" json:"metaalarm_value_path,omitempty"`

	MetaAlarmParents  []string `bson:"ma_parents,omitempty" json:"ma_parents,omitempty"`
	MetaAlarmChildren []string `bson:"ma_children,omitempty" json:"ma_children,omitempty"`

	// ManualMetaAlarmAutoResolve is used for manual meta alarms.
	ManualMetaAlarmAutoResolve bool `bson:"manual_meta_alarm_auto_resolve,omitempty" json:"manual_meta_alarm_auto_resolve,omitempty"`
	// DisplayName is used for manual meta alarms.
	DisplayName string `bson:"display_name,omitempty" json:"display_name,omitempty"`

	PbehaviorInfo PbehaviorInfo `bson:"pbehavior_info" json:"pbehavior_info"`

	// Initiator is used to detect who emits event.
	// InitiatorUser - UI
	// InitiatorSystem - engines
	// InitiatorExternal - third tool
	Initiator string `bson:"initiator" json:"initiator"`

	// Only for EventTypeRunDelayedScenario
	DelayedScenarioID   string `bson:"delayed_scenario_id,omitempty" json:"delayed_scenario_id,omitempty"`
	DelayedScenarioData string `bson:"delayed_scenario_data,omitempty" json:"delayed_scenario_data,omitempty"`

	// AddedToServices contains ids of entity services to which entity has been added as dependency.
	AddedToServices []string `bson:"added_to_services,omitempty" json:"added_to_services,omitempty"`
	// RemovedFromServices contains ids of entity services from which entity has been removed as dependency.
	RemovedFromServices []string `bson:"removed_from_services,omitempty" json:"removed_from_services,omitempty"`

	// IdleRuleApply is used if event is emitted by idle rule.
	IdleRuleApply string `bson:"idle_rule_apply,omitempty" json:"idle_rule_apply,omitempty"`

	// Execution is used only for instruction events: EventTypeInstructionStarted, EventTypeInstructionCompleted, etc..
	Execution string `bson:"execution,omitempty" json:"execution,omitempty"`

	// Instruction is used only for manual instructions kpi metrics
	Instruction string `bson:"instruction,omitempty" json:"instruction,omitempty"`

	// IsMetaAlarmUpdated is true if an alarm is added to a meta alarm on an event.
	IsMetaAlarmUpdated bool `bson:"ma_updated,omitempty" json:"ma_updated,omitempty"`
	// IsInstructionMatched is true if an alarm is matched to an auto instruction on an event.
	IsInstructionMatched bool `bson:"instr_matched,omitempty" json:"instr_matched,omitempty"`

	Healthcheck bool `bson:"healthcheck,omitempty" json:"healthcheck,omitempty"`

	StateSettingUpdated bool `bson:"state_setting_updated,omitempty" json:"state_setting_updated,omitempty"`
}

// Format an event
//
//	"timestamp" is fill with time.Now()
//	"event_type" is fill with EventTypeCheck
func (e *Event) Format() {
	//events can't be later or earlier than MaxEventTimestampVariation
	now := datetime.NewCpsTime()
	if e.Timestamp.IsZero() || e.Timestamp.Time.Before(now.Add(-MaxEventTimestampVariation)) || e.Timestamp.Time.After(now.Add(MaxEventTimestampVariation)) {
		e.Timestamp = now
	}
	e.ReceivedTimestamp = datetime.NewMicroTime()
	if e.EventType == "" {
		e.EventType = EventTypeCheck
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
		return fmt.Errorf("Event.InjectExtraInfos json decode: %w", err)
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
func (e *Event) IsContextable() bool {
	switch e.EventType {
	case EventTypeCheck, EventTypeMetaAlarm,
		EventTypeEntityToggled, EventTypeEntityUpdated, EventTypeResolveDeleted:
		return true
	default:
		return false
	}
}

func (e *Event) IsOnlyServiceUpdate() bool {
	switch e.EventType {
	case EventTypeEntityToggled, EventTypeEntityUpdated, EventTypeResolveDeleted:
		return true
	default:
		return false
	}
}

// IsMatched tell if an event is catched by a regex
func (e *Event) IsMatched(regex string, fields []string) bool {
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
func (e *Event) IsValid() error {
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
		if e.Component == "" {
			return errt.NewUnknownError(errors.New("missing component"))
		}
		if e.Resource == "" {
			return errt.NewUnknownError(errors.New("missing resource"))
		}
	default:
		return errt.NewUnknownError(fmt.Errorf("wrong source type: %v", e.SourceType))
	}

	if !isValidEventType(e.EventType) {
		return errt.NewUnknownError(fmt.Errorf("wrong event type: %v", e.EventType))
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

var cpsNumberType = reflect.TypeOf(CpsNumber(0))
var cpsNumberPtrType = reflect.PtrTo(cpsNumberType)
var cpsTimeType = reflect.TypeOf(datetime.CpsTime{})
var stringType = reflect.TypeOf("")
var stringPtrType = reflect.PtrTo(stringType)
var boolType = reflect.TypeOf(false)
var mapStringStringType = reflect.TypeOf(map[string]string{})

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
			return fmt.Errorf("%[1]T value cannot be converted to an integer: %+[1]v", value)
		}
		field.Set(reflect.ValueOf(CpsNumber(integerValue)))

	case cpsNumberPtrType:
		integerValue, success := AsInteger(value)
		if !success {
			return fmt.Errorf("%[1]T value cannot be converted to an integer: %+[1]v", value)
		}
		cpsNumberValue := CpsNumber(integerValue)
		field.Set(reflect.ValueOf(&cpsNumberValue))

	case cpsTimeType:
		integerValue, success := AsInteger(value)
		if !success {
			return fmt.Errorf("%[1]T value cannot be converted to an integer: %+[1]v", value)
		}
		cpsTimeValue := datetime.CpsTime{Time: time.Unix(integerValue, 0)}
		field.Set(reflect.ValueOf(cpsTimeValue))

	case stringType:
		stringValue, success := utils.AsString(value)
		if !success {
			return fmt.Errorf("%[1]T value cannot be assigned to a string: %+[1]v", value)
		}
		field.Set(reflect.ValueOf(stringValue))

	case stringPtrType:
		stringValue, success := utils.AsString(value)
		if !success {
			return fmt.Errorf("%[1]T value cannot be assigned to a string: %+[1]v", value)
		}
		field.Set(reflect.ValueOf(&stringValue))

	case boolType:
		boolValue, success := value.(bool)
		if !success {
			return fmt.Errorf("%[1]T value cannot be assigned to a bool: %+[1]v", value)
		}
		field.Set(reflect.ValueOf(boolValue))

	case mapStringStringType:
		var err error
		if m1, ok := value.(map[string]any); ok {
			err = setMapStringStringField(field, m1)
		} else if m2, ok := value.(map[string]string); ok {
			err = setMapStringStringField(field, m2)
		} else {
			return fmt.Errorf("%[1]T value cannot be assigned to a map[string]string: %+[1]v", value)
		}
		if err != nil {
			return err
		}

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

// GetStringField is a magic getter for string fields for easier field retrieving when matching event pattern
func (e *Event) GetStringField(f string) (string, bool) {
	switch f {
	case "connector":
		return e.Connector, true
	case "connector_name":
		return e.ConnectorName, true
	case "component":
		return e.Component, true
	case "resource":
		return e.Resource, true
	case "output":
		return e.Output, true
	case "long_output":
		return e.LongOutput, true
	case "event_type":
		return e.EventType, true
	case "source_type":
		return e.SourceType, true
	case "author":
		return e.Author, true
	case "initiator":
		return e.Initiator, true
	default:
		return "", false
	}
}

// GetIntField is a magic getter for int fields for easier field retrieving when matching event pattern
func (e *Event) GetIntField(f string) (int64, bool) {
	switch f {
	case "state":
		return int64(e.State), true
	default:
		return 0, false
	}
}

// GetExtraInfoVal is a magic getter for extra infos fields for easier field retrieving when matching event pattern
func (e *Event) GetExtraInfoVal(f string) (interface{}, bool) {
	if v, ok := e.ExtraInfos[f]; ok {
		return v, true
	}

	return nil, false
}

// setMapStringStringField sets the value of a field of type map[string]string
func setMapStringStringField[T any | string](field reflect.Value, value map[string]T) error {
	if field.IsNil() {
		field.Set(reflect.MakeMap(field.Type()))
	}

	for key, value := range value {
		stringValue, success := utils.AsString(value)
		if !success {
			return fmt.Errorf("value cannot be assigned to a map[string]string under key %q: %+v", key, value)
		}
		field.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(stringValue))
	}
	return nil
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
		EventTypeDeclareTicketWebhook,
		EventTypeChangestate,
		EventTypeSnooze,
		EventTypeUnsnooze,
		EventTypeUncancel,
		EventTypeResolveCancel,
		EventTypeResolveClose,
		EventTypeResolveDeleted,
		EventTypePbhEnter,
		EventTypePbhLeaveAndEnter,
		EventTypePbhLeave,
		EventTypeUpdateStatus,
		EventTypeMetaAlarm,
		EventTypeMetaAlarmAttachChildren,
		EventTypeMetaAlarmDetachChildren,
		EventTypeMetaAlarmUngroup,
		EventTypeManualMetaAlarmGroup,
		EventTypeManualMetaAlarmUngroup,
		EventTypeManualMetaAlarmUpdate,
		EventTypeRecomputeEntityService,
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
		EventTypeInstructionJobStarted,
		EventTypeInstructionJobCompleted,
		EventTypeInstructionJobFailed,
		EventTypeJunitTestSuiteUpdated,
		EventTypeJunitTestCaseUpdated,
		EventTypeTrigger,
		EventTypeAutoInstructionActivate:
		return true
	}

	return false
}
