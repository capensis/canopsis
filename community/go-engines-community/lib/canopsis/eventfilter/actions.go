package eventfilter

import (
	"encoding/json"
	"fmt"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// ActionParameters is a type containing the parameters that can be used by an
// action.
type ActionParameters struct {
	// Event is the event for which the data needs to be fetched.
	Event types.Event

	// RegexMatch contains the values of the sub-expressions of the regular
	// expressions used in the pattern.
	RegexMatch pattern.EventRegexMatches

	// ExternalData contains the data of external sources.
	ExternalData map[string]interface{}
}

// ActionBase is a type containing the fields that are common to all actions.
type ActionBase struct {
	Type ActionType `bson:"type"`
}

// ActionProcessor is an interface for types that apply an action to an event.
type ActionProcessor interface {
	// Apply applies the action to the event using the provided parameters. It
	// returns the modified event, and can return an error, swhich will be used
	// to determine the outcome of the enrichment rule.
	Apply(event types.Event, parameters ActionParameters, report *Report) (types.Event, error)

	// Validate returns an error if the action is invalid.
	Validate() error
}

// Action is a type that represents an action that can be applied to an event
// by an enrichment rule.
type Action struct {
	ActionBase
	ActionProcessor `swaggerignore:"true"`

	timezoneConfig *config.TimezoneConfig
}

func (a *Action) SetTimezoneConfig(cfg *config.TimezoneConfig) {
	a.timezoneConfig = cfg
}

func (a *Action) UnmarshalJSON(b []byte) error {
	var jsonPatterns interface{}
	err := json.Unmarshal(b, &jsonPatterns)
	if err != nil {
		return err
	}

	marshalled, err := bson.Marshal(bson.M{
		"action": jsonPatterns,
	})
	if err != nil {
		return err
	}

	var wrapper struct {
		Action Action `bson:"action"`
	}

	err = bson.Unmarshal(marshalled, &wrapper)
	if err != nil {
		return err
	}

	*a = wrapper.Action
	return nil
}

func (a Action) MarshalJSON() ([]byte, error) {
	bsonType, bsonBytes, err := bson.MarshalValue(a)
	if err != nil {
		return nil, err
	}

	if bsonType == bsontype.Undefined {
		return nil, nil
	}

	if bsonType == bsontype.Null {
		res, err := json.Marshal(nil)
		if err != nil {
			return nil, err
		}

		return res, nil
	}

	var unmarshalledBson map[string]interface{}
	raw := bson.RawValue{
		Type:  bsontype.EmbeddedDocument,
		Value: bsonBytes,
	}
	err = raw.Unmarshal(&unmarshalledBson)
	if err != nil {
		return nil, err
	}

	res, err := json.Marshal(unmarshalledBson)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *Action) UnmarshalBSONValue(_ bsontype.Type, b []byte) error {
	err := bson.Unmarshal(b, &a.ActionBase)
	if err != nil {
		return err
	}

	var unexpectedFields map[string]interface{}
	switch a.Type {
	case SetField:
		var p SetFieldProcessor
		err = bson.Unmarshal(b, &p)
		unexpectedFields = p.UnexpectedFields
		p.UnexpectedFields = nil
		a.ActionProcessor = p
	case SetFieldFromTemplate:
		var p SetFieldFromTemplateProcessor
		p.Value.SetTimezoneConfig(a.timezoneConfig)
		err = bson.Unmarshal(b, &p)
		unexpectedFields = p.UnexpectedFields
		p.UnexpectedFields = nil
		a.ActionProcessor = p
	case SetEntityInfoFromTemplate:
		var p SetEntityInfoFromTemplateProcessor
		err = bson.Unmarshal(b, &p)
		unexpectedFields = p.UnexpectedFields
		p.UnexpectedFields = nil
		a.ActionProcessor = p
	case Copy:
		var p CopyProcessor
		err = bson.Unmarshal(b, &p)
		unexpectedFields = p.UnexpectedFields
		p.UnexpectedFields = nil
		a.ActionProcessor = p
	case SetEntityInfo:
		var p SetEntityInfoProcessor
		err = bson.Unmarshal(b, &p)
		unexpectedFields = p.UnexpectedFields
		p.UnexpectedFields = nil
		a.ActionProcessor = p
	case CopyToEntityInfo:
		var p CopyToEntityInfoProcessor
		err = bson.Unmarshal(b, &p)
		unexpectedFields = p.UnexpectedFields
		p.UnexpectedFields = nil
		a.ActionProcessor = p
	case UnsetAction:
		return fmt.Errorf("missing action field")
	default:
		return fmt.Errorf("unknown action type: %s", a.Type)
	}

	if err != nil {
		return err
	}

	err = a.ActionProcessor.Validate()
	if err != nil {
		return err
	}

	// Get the unexpected fields of the ActionProcessor that are not fields of
	// the ActionBase.
	unexpectedFieldNames := make([]string, 0, len(unexpectedFields))
	for key := range unexpectedFields {
		if key != "type" {
			unexpectedFieldNames = append(unexpectedFieldNames, key)
		}
	}
	if len(unexpectedFieldNames) != 0 {
		return fmt.Errorf("Unexpected action fields: %s", strings.Join(unexpectedFieldNames, ", "))
	}

	return nil
}

func (a Action) MarshalBSONValue() (bsontype.Type, []byte, error) {
	t, p, err := bson.MarshalValue(a.ActionProcessor)
	if err != nil {
		return bsontype.Undefined, nil, err
	}

	if t == bsontype.Undefined || t == bsontype.Null {
		return t, nil, nil
	}

	m := make(map[string]interface{})
	err = bson.Unmarshal(p, &m)
	if err != nil {
		return t, nil, err
	}

	m["type"] = a.Type

	return bson.MarshalValue(m)
}

// SetFieldProcessor is an ActionProcessor that sets of a field of an event to
// a constant.
type SetFieldProcessor struct {
	// Name is the name of the field that is modified.
	Name types.OptionalString `bson:"name"`

	// Value is the new value of the field.
	Value types.OptionalInterface `bson:"value"`

	// When unmarshalling a BSON document, the fields of this document that are
	// not defined in this struct are added to UnexpectedFields.
	UnexpectedFields map[string]interface{} `bson:",inline"`
}

// Validate returns an error if the action is invalid.
func (p SetFieldProcessor) Validate() error {
	if !p.Name.Set {
		return fmt.Errorf("the name field is required")
	}
	if !p.Value.Set {
		return fmt.Errorf("the value field is required")
	}

	return nil
}

// Apply applies the action to the event using the provided parameters. It
// returns the modified event, and returns an error if the field does not exist
// or if the type of the value is not compatible with the field.
func (p SetFieldProcessor) Apply(event types.Event, parameters ActionParameters, report *Report) (types.Event, error) {
	err := event.SetField(p.Name.Value, p.Value.Value)
	return event, err
}

// SetFieldFromTemplateProcessor is an ActionProcessor that sets of a string
// field of an event using a template.
type SetFieldFromTemplateProcessor struct {
	// Name is the name of the field that is modified.
	Name types.OptionalString `bson:"name"`

	// Value is the template used to set the new value of the field.
	Value types.OptionalTemplate `bson:"value"`

	// buffer is a buffer used when executing the templates.
	buffer strings.Builder

	// When unmarshalling a BSON document, the fields of this document that are
	// not defined in this struct are added to UnexpectedFields.
	UnexpectedFields map[string]interface{} `bson:",inline"`
}

// Validate returns an error if the action is invalid.
func (p SetFieldFromTemplateProcessor) Validate() error {
	if !p.Name.Set {
		return fmt.Errorf("the name field is required")
	}
	if !p.Value.Set {
		return fmt.Errorf("the value field is required")
	}

	return nil
}

// Apply applies the action to the event using the provided parameters. It
// returns the modified event, and returns an error if the field does not exist
// or if the type of the value is not compatible with the field.
func (p SetFieldFromTemplateProcessor) Apply(event types.Event, parameters ActionParameters, report *Report) (types.Event, error) {
	p.buffer.Reset()
	err := p.Value.Value.Execute(&p.buffer, parameters)
	if err != nil {
		return event, err
	}
	err = event.SetField(p.Name.Value, p.buffer.String())
	return event, err
}

// SetEntityInfoFromTemplateProcessor is an ActionProcessor that sets of a
// information of an entity.
type SetEntityInfoFromTemplateProcessor struct {
	// Name is the name of the information that is modified.
	Name types.OptionalString `bson:"name"`

	// Value is the new value of the information.
	Value types.OptionalTemplate `bson:"value"`

	// Description is the description of the information.
	Description types.OptionalString `bson:"description"`

	// buffer is a buffer used when executing the templates.
	buffer strings.Builder

	// When unmarshalling a BSON document, the fields of this document that are
	// not defined in this struct are added to UnexpectedFields.
	UnexpectedFields map[string]interface{} `bson:",inline"`
}

// Validate returns an error if the action is invalid.
func (p SetEntityInfoFromTemplateProcessor) Validate() error {
	if !p.Name.Set {
		return fmt.Errorf("the name field is required")
	}
	if !p.Value.Set {
		return fmt.Errorf("the value field is required")
	}

	return nil
}

// Apply applies the action to the event using the provided parameters. It
// returns the modified event, and returns an error if the field does not exist
// or if the type of the value is not compatible with the field.
func (p SetEntityInfoFromTemplateProcessor) Apply(event types.Event, parameters ActionParameters, report *Report) (types.Event, error) {
	if event.Entity == nil {
		return event, fmt.Errorf("cannot set information before enrichment with entity")
	}

	p.buffer.Reset()
	err := p.Value.Value.Execute(&p.buffer, parameters)
	if err != nil {
		return event, err
	}
	value := p.buffer.String()
	entityUpdated := setEntityInfo(event.Entity, value, p.Name, p.Description)
	if report != nil && entityUpdated {
		report.EntityUpdated = entityUpdated
	}

	return event, nil
}

// CopyProcessor is an ActionProcessor that copies a value from a field to
// another.
type CopyProcessor struct {
	// From is the name of the field of the ActionParameters the value will be
	// copied from.
	From types.OptionalString `bson:"from"`

	// To is the name of the field of the event the value will be copied to.
	To types.OptionalString `bson:"to"`

	// When unmarshalling a BSON document, the fields of this document that are
	// not defined in this struct are added to UnexpectedFields.
	UnexpectedFields map[string]interface{} `bson:",inline"`
}

// Validate returns an error if the action is invalid.
func (p CopyProcessor) Validate() error {
	if !p.From.Set {
		return fmt.Errorf("the from field is required")
	}
	if !p.To.Set {
		return fmt.Errorf("the to field is required")
	}

	return nil
}

// Apply applies the action to the event using the provided parameters. It
// returns the modified event, and returns an error if the field does not exist
// or if the type of the value is not compatible with the field.
func (p CopyProcessor) Apply(event types.Event, parameters ActionParameters, report *Report) (types.Event, error) {
	value, err := utils.GetField(parameters, p.From.Value)
	if err != nil {
		return event, err
	}

	err = event.SetField(p.To.Value, value)
	if err != nil {
		return event, err
	}

	return event, nil
}

// SetEntityInfoProcessor is an ActionProcessor that sets of a
// information of an entity.
type SetEntityInfoProcessor struct {
	// Name is the name of the information that is modified.
	Name types.OptionalString `bson:"name"`

	// Description is the description of the information.
	Description types.OptionalString `bson:"description"`

	// Value is the new value of the field.
	Value types.OptionalInterface `bson:"value"`

	// When unmarshalling a BSON document, the fields of this document that are
	// not defined in this struct are added to UnexpectedFields.
	UnexpectedFields map[string]interface{} `bson:",inline"`
}

func (p SetEntityInfoProcessor) Validate() error {
	if !p.Name.Set {
		return fmt.Errorf("the name field is required")
	}
	if !p.Value.Set {
		return fmt.Errorf("the value field is required")
	}

	return nil
}

func (p SetEntityInfoProcessor) Apply(event types.Event, parameters ActionParameters, report *Report) (types.Event, error) {
	if event.Entity == nil {
		return event, fmt.Errorf("cannot set information before enrichment with entity")
	}

	value := p.Value.Value
	entityUpdated := setEntityInfo(event.Entity, value, p.Name, p.Description)
	if report != nil && entityUpdated {
		report.EntityUpdated = entityUpdated
	}

	return event, nil
}

// CopyToEntityInfoProcessor is an ActionProcessor that sets of a
// information of an entity.
type CopyToEntityInfoProcessor struct {
	// Name is the name of the information that is modified.
	Name types.OptionalString `bson:"name"`

	// Description is the description of the information.
	Description types.OptionalString `bson:"description"`

	// From is the name of the field of the ActionParameters the value will be
	// copied from.
	From types.OptionalString `bson:"from"`

	// When unmarshalling a BSON document, the fields of this document that are
	// not defined in this struct are added to UnexpectedFields.
	UnexpectedFields map[string]interface{} `bson:",inline"`
}

func (p CopyToEntityInfoProcessor) Validate() error {
	if !p.Name.Set {
		return fmt.Errorf("the name field is required")
	}
	if !p.From.Set {
		return fmt.Errorf("the from field is required")
	}

	return nil
}

func (p CopyToEntityInfoProcessor) Apply(event types.Event, parameters ActionParameters, report *Report) (types.Event, error) {
	if event.Entity == nil {
		return event, fmt.Errorf("cannot set information before enrichment with entity")
	}

	value, err := utils.GetField(parameters, p.From.Value)
	if err != nil {
		return event, err
	}

	entityUpdated := setEntityInfo(event.Entity, value, p.Name, p.Description)
	if report != nil && entityUpdated {
		report.EntityUpdated = entityUpdated
	}

	return event, nil
}

func setEntityInfo(entity *types.Entity, value interface{}, name, description types.OptionalString) bool {
	info, ok := entity.Infos[name.Value]
	if !ok {
		info = types.Info{}
	}

	entityUpdated := false
	valueChanged := !ok || info.Value != value
	if valueChanged {
		entityUpdated = true
	}

	info.Name = name.Value
	info.Value = value
	if description.Set {
		info.Description = description.Value
	}

	if entity.Infos == nil {
		entity.Infos = make(map[string]types.Info)
	}

	entity.Infos[name.Value] = info

	return entityUpdated
}
