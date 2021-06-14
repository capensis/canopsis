package eventfilter

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
)

// Outcome is a type representing the outcome of an event filter's rule. A
// value of type Outcome is returned when a rule is applied to indicate what
// should be done with the event.
type Outcome string

const (
	// UnsetOutcome is a const used for outcomes that have not been set.
	UnsetOutcome Outcome = ""

	// Pass indicates that the event should be passed to the next rule of the
	// event filter.
	Pass Outcome = "pass"

	// Break indicates that the event should skip the remaining rules of the
	// event filter.
	Break Outcome = "break"

	// Drop indicates that the event should be removed.
	Drop Outcome = "drop"
)

// SetBSON unmarshals a BSON value into an Outcome.
func (o *Outcome) SetBSON(raw bson.Raw) error {
	var s string
	err := raw.Unmarshal(&s)
	if err != nil {
		return fmt.Errorf("unable to parse outcome: %v", err)
	}

	switch Outcome(s) {
	case UnsetOutcome, Pass, Break, Drop:
		*o = Outcome(s)
	default:
		return fmt.Errorf("unknown rule outcome: %s", s)
	}

	return nil
}

// Type is a type representing the type of an event filter's rule.
type Type string

const (
	// RuleTypeUnset is a const used for rules whose type is not set.
	RuleTypeUnset Type = ""

	// RuleTypeBreak is a type of rule that break events out of the event
	// filter.
	RuleTypeBreak Type = "break"

	// RuleTypeDrop is a type of rule that delete events.
	RuleTypeDrop Type = "drop"

	// RuleTypeEnrichment is a type of rule that modify an event.
	RuleTypeEnrichment Type = "enrichment"
)

// SetBSON unmarshals a BSON value into a Type.
func (t *Type) SetBSON(raw bson.Raw) error {
	var s string
	err := raw.Unmarshal(&s)
	if err != nil {
		return fmt.Errorf("unable to parse type: %v", err)
	}

	switch Type(s) {
	case RuleTypeBreak, RuleTypeDrop, RuleTypeEnrichment:
		*t = Type(s)
	default:
		return fmt.Errorf("unknown rule type: %s", s)
	}

	return nil
}

// ActionType is a type representing the type of an enrichment rule's action.
type ActionType string

const (
	// UnsetAction is a const used for actions whose type is not set.
	UnsetAction ActionType = ""

	// SetField is a type of action that sets a field of an event to a
	// constant.
	SetField ActionType = "set_field"

	// SetFieldFromTemplate is a type of action that sets a string field of an
	// event using a template.
	SetFieldFromTemplate ActionType = "set_field_from_template"

	// SetEntityInfoFromTemplate is a type of action that sets an information
	// of an entity using a template.
	SetEntityInfoFromTemplate ActionType = "set_entity_info_from_template"

	// Copy is a type of action that copies a value from a field to another.
	Copy ActionType = "copy"
)

// SetBSON unmarshals a BSON value into a ActionType.
func (t *ActionType) SetBSON(raw bson.Raw) error {
	var s string
	err := raw.Unmarshal(&s)
	if err != nil {
		return fmt.Errorf("unable to parse action: %v", err)
	}

	switch ActionType(s) {
	case SetField, SetFieldFromTemplate, SetEntityInfoFromTemplate, Copy:
		*t = ActionType(s)
	default:
		return fmt.Errorf("unknown action: %s", s)
	}

	return nil
}
