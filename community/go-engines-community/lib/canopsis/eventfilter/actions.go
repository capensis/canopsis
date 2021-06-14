package eventfilter

import (
	"fmt"
	"strings"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/utils"
	"github.com/globalsign/mgo/bson"
)

// ActionParameters is a type containing the parameters that can be used by an
// action.
type ActionParameters struct {
	DataSourceGetterParameters

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
	ActionProcessor
}

// SetBSON unmarshals a BSON value into an Action.
func (a *Action) SetBSON(raw bson.Raw) error {
	err := raw.Unmarshal(&a.ActionBase)
	if err != nil {
		return err
	}

	var unexpectedFields map[string]interface{}
	switch a.Type {
	case SetField:
		var p SetFieldProcessor
		err = raw.Unmarshal(&p)
		unexpectedFields = p.UnexpectedFields
		a.ActionProcessor = p
	case SetFieldFromTemplate:
		var p SetFieldFromTemplateProcessor
		err = raw.Unmarshal(&p)
		unexpectedFields = p.UnexpectedFields
		a.ActionProcessor = p
	case SetEntityInfoFromTemplate:
		var p SetEntityInfoFromTemplateProcessor
		err = raw.Unmarshal(&p)
		unexpectedFields = p.UnexpectedFields
		a.ActionProcessor = p
	case Copy:
		var p CopyProcessor
		err = raw.Unmarshal(&p)
		unexpectedFields = p.UnexpectedFields
		a.ActionProcessor = p
	case UnsetAction:
		return fmt.Errorf("missing action field")
	default:
		return fmt.Errorf("unknown action type: %s", a.Type)
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

// SetFieldProcessor is an ActionProcessor that sets of a field of an event to
// a constant.
type SetFieldProcessor struct {
	// Name is the name of the field that is modified.
	Name utils.OptionalString `bson:"name"`

	// Value is the new value of the field.
	Value utils.OptionalInterface `bson:"value"`

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
	Name utils.OptionalString `bson:"name"`

	// Value is the template used to set the new value of the field.
	Value utils.OptionalTemplate `bson:"value"`

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
	Name utils.OptionalString `bson:"name"`

	// Value is the new value of the information.
	Value utils.OptionalTemplate `bson:"value"`

	// Description is the description of the information.
	Description utils.OptionalString `bson:"description"`

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

	info, success := event.Entity.Infos[p.Name.Value]
	if !success {
		info = types.Info{}
	}

	valueChanged := !success || info.Value != value
	if report != nil && valueChanged {
		report.EntityUpdated = true
	}

	info.Name = p.Name.Value
	info.Value, info.RealValue = value, value
	if p.Description.Set {
		info.Description = p.Description.Value
	}

	event.Entity.Infos[p.Name.Value] = info

	return event, nil
}

// CopyProcessor is an ActionProcessor that copies a value from a field to
// another.
type CopyProcessor struct {
	// From is the name of the field of the ActionParameters the value will be
	// copied from.
	From utils.OptionalString `bson:"from"`

	// To is the name of the field of the event the value will be copied to.
	To utils.OptionalString `bson:"to"`

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
