package eventfilter

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
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

func (o *Outcome) UnmarshalBSONValue(t bsontype.Type, b []byte) error {
	var s string
	if t == bsontype.String {
		var ok bool
		s, _, ok = bsoncore.ReadString(b)
		if !ok {
			return fmt.Errorf("unable to parse outcome")
		}
	} else {
		return fmt.Errorf("unable to parse outcome, unknown type: %v", t)
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

func (t *Type) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	var s string
	if valueType == bsontype.String {
		var ok bool
		s, _, ok = bsoncore.ReadString(b)
		if !ok {
			return fmt.Errorf("unable to parse type")
		}
	} else {
		return fmt.Errorf("unable to parse type, unknown type: %v", valueType)
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

	// SetEntityInfo is a type of action that sets an information
	// of an entity using a constant.
	SetEntityInfo ActionType = "set_entity_info"

	// CopyToEntityInfo is a type of action that copies a value from a field to
	// an information of an entity.
	CopyToEntityInfo ActionType = "copy_to_entity_info"

	// Copy is a type of action that copies a value from a field to another.
	Copy ActionType = "copy"
)

func (t *ActionType) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	var s string
	if valueType == bsontype.String {
		var ok bool
		s, _, ok = bsoncore.ReadString(b)
		if !ok {
			return fmt.Errorf("unable to parse action")
		}
	} else {
		return fmt.Errorf("unable to parse action, unknown type: %v", valueType)
	}

	switch ActionType(s) {
	case SetField, SetFieldFromTemplate, SetEntityInfoFromTemplate, Copy, SetEntityInfo, CopyToEntityInfo:
		*t = ActionType(s)
	default:
		return fmt.Errorf("unknown action: %s", s)
	}

	return nil
}
