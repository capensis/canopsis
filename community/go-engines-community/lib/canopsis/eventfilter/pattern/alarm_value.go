package pattern

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"strings"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	mgobson "github.com/globalsign/mgo/bson"
	mongobson "go.mongodb.org/mongo-driver/bson"
)

// AlarmValueRegexMatches is a type that contains the values of the
// sub-expressions of regular expressions for each of the fields of an
// AlarmValue that contain strings.
type AlarmValueRegexMatches struct {
	ACK               AlarmStepRegexMatches
	Canceled          AlarmStepRegexMatches
	Done              AlarmStepRegexMatches
	Snooze            AlarmStepRegexMatches
	State             AlarmStepRegexMatches
	Status            AlarmStepRegexMatches
	Ticket            AlarmTicketRegexMatches
	Component         RegexMatches
	Connector         RegexMatches
	ConnectorName     RegexMatches
	DisplayName       RegexMatches
	Extra             map[string]RegexMatches
	InitialOutput     RegexMatches
	Output            RegexMatches
	InitialLongOutput RegexMatches
	LongOutput        RegexMatches
	Resource          RegexMatches
	Parents           RegexMatches
	Children          RegexMatches
}

// NewAlarmValueRegexMatches creates an AlarmValueRegexMatches, with the Extra
// field initialized to an empty map.
func NewAlarmValueRegexMatches() AlarmValueRegexMatches {
	return AlarmValueRegexMatches{
		Extra: map[string]RegexMatches{},
		Ticket: NewAlarmTicketRegexMatches(),
	}
}

// AlarmValueFields is a type representing a pattern that can be applied to an
// alarm value.
// The fields are not defined directly in the AlarmValuePattern struct to
// make the unmarshalling easier.
type AlarmValueFields struct {
	ACK                           AlarmStepRefPattern         `bson:"ack,omitempty"`
	Canceled                      AlarmStepRefPattern         `bson:"canceled,omitempty"`
	Done                          AlarmStepRefPattern         `bson:"done,omitempty"`
	Snooze                        AlarmStepRefPattern         `bson:"snooze,omitempty"`
	State                         AlarmStepRefPattern         `bson:"state,omitempty"`
	Status                        AlarmStepRefPattern         `bson:"status,omitempty"`
	Ticket                        AlarmTicketRefPattern       `bson:"ticket,omitempty"`
	Component                     StringPattern               `bson:"component"`
	Connector                     StringPattern               `bson:"connector"`
	ConnectorName                 StringPattern               `bson:"connector_name"`
	CreationDate                  TimePattern                 `bson:"creation_date"`
	ActivationDate                TimeRefPattern              `bson:"activation_date"`
	DisplayName                   StringPattern               `bson:"display_name"`
	Extra                         map[string]InterfacePattern `bson:"extra"`
	HardLimit                     IntegerRefPattern           `bson:"hard_limit,omitempty"`
	InitialOutput                 StringPattern               `bson:"initial_output"`
	Output                        StringPattern               `bson:"output"`
	InitialLongOutput             StringPattern               `bson:"initial_long_output"`
	LongOutput                    StringPattern               `bson:"long_output"`
	LastUpdateDate                TimePattern                 `bson:"last_update_date"`
	LastEventDate                 TimePattern                 `bson:"last_event_date"`
	Resource                      StringPattern               `bson:"resource,omitempty"`
	Resolved                      TimeRefPattern              `bson:"resolved,omitempty"`
	StateChangesSinceStatusUpdate IntegerPattern              `bson:"state_changes_since_status_update,omitempty"`
	TotalStateChanges             IntegerPattern              `bson:"total_state_changes,omitempty"`
	Parents                       StringArrayPattern          `bson:"parents,omitempty"`
	Children                      StringArrayPattern          `bson:"children,omitempty"`

	// When unmarshalling a BSON document, the fields of this document that are
	// not defined in this struct are added to UnexpectedFields.
	UnexpectedFields map[string]interface{} `bson:",inline"`
}

// Empty returns true if the pattern has not been set
func (p AlarmValueFields) Empty() bool {
	return p.ACK.Empty() &&
		p.Done.Empty() &&
		p.Snooze.Empty() &&
		p.State.Empty() &&
		p.Status.Empty() &&
		p.Ticket.Empty() &&
		p.Component.Empty() &&
		p.Connector.Empty() &&
		p.ConnectorName.Empty() &&
		p.CreationDate.Empty() &&
		p.DisplayName.Empty() &&
		p.HardLimit.Empty() &&
		p.InitialOutput.Empty() &&
		p.Output.Empty() &&
		p.InitialLongOutput.Empty() &&
		p.LongOutput.Empty() &&
		p.LastUpdateDate.Empty() &&
		p.LastEventDate.Empty() &&
		p.Resource.Empty() &&
		p.Resolved.Empty() &&
		p.StateChangesSinceStatusUpdate.Empty() &&
		p.TotalStateChanges.Empty() &&
		p.Parents.Empty() &&
		p.Children.Empty()
}

// AlarmValuePattern is a type representing a pattern that can be applied to an
// alarm value
type AlarmValuePattern struct {
	AlarmValueFields
}

func (p AlarmValuePattern) IsSet() bool {
	return p.AlarmValueFields.Output.IsSet() ||
		p.AlarmValueFields.Snooze.IsSet() ||
		p.AlarmValueFields.LastUpdateDate.IsSet() ||
		p.AlarmValueFields.Ticket.IsSet() ||
		p.AlarmValueFields.ACK.IsSet() ||
		p.AlarmValueFields.Done.IsSet() ||
		p.AlarmValueFields.ActivationDate.IsSet() ||
		p.AlarmValueFields.Canceled.IsSet() ||
		p.AlarmValueFields.Component.IsSet() ||
		p.AlarmValueFields.Connector.IsSet() ||
		p.AlarmValueFields.ConnectorName.IsSet() ||
		p.AlarmValueFields.CreationDate.IsSet() ||
		p.AlarmValueFields.DisplayName.IsSet() ||
		len(p.AlarmValueFields.Extra) > 0 ||
		p.AlarmValueFields.HardLimit.IsSet() ||
		p.AlarmValueFields.InitialLongOutput.IsSet() ||
		p.AlarmValueFields.InitialOutput.IsSet() ||
		p.AlarmValueFields.LastEventDate.IsSet() ||
		p.AlarmValueFields.LongOutput.IsSet() ||
		p.AlarmValueFields.Resolved.IsSet() ||
		p.AlarmValueFields.Resource.IsSet() ||
		p.AlarmValueFields.State.IsSet() ||
		p.AlarmValueFields.StateChangesSinceStatusUpdate.IsSet() ||
		p.AlarmValueFields.Status.IsSet() ||
		p.AlarmValueFields.TotalStateChanges.IsSet() ||
		!p.Parents.Empty() || !p.Children.Empty()
}

// AsMongoQuery returns a mongodb filter from the AlarmValuePattern
func (p AlarmValuePattern) AsMongoQuery(prefix string, query mgobson.M) {
	p.ACK.AsMongoQuery(fmt.Sprintf("%s.ack", prefix), query)
	p.Canceled.AsMongoQuery(fmt.Sprintf("%s.canceled", prefix), query)
	p.Done.AsMongoQuery(fmt.Sprintf("%s.done", prefix), query)
	p.Snooze.AsMongoQuery(fmt.Sprintf("%s.snooze", prefix), query)
	p.State.AsMongoQuery(fmt.Sprintf("%s.state", prefix), query)
	p.Status.AsMongoQuery(fmt.Sprintf("%s.status", prefix), query)
	p.Ticket.AsMongoQuery(fmt.Sprintf("%s.ticket", prefix), query)

	if !p.Component.Empty() {
		query[fmt.Sprintf("%s.component", prefix)] = p.Component.AsMongoQuery()
	}
	if !p.Connector.Empty() {
		query[fmt.Sprintf("%s.connector", prefix)] = p.Connector.AsMongoQuery()
	}
	if !p.ConnectorName.Empty() {
		query[fmt.Sprintf("%s.connector_name", prefix)] = p.ConnectorName.AsMongoQuery()
	}
	if !p.CreationDate.Empty() {
		query[fmt.Sprintf("%s.creation_date", prefix)] = p.CreationDate.AsMongoQuery()
	}
	if !p.DisplayName.Empty() {
		query[fmt.Sprintf("%s.display_name", prefix)] = p.DisplayName.AsMongoQuery()
	}
	if !p.HardLimit.Empty() {
		query[fmt.Sprintf("%s.hard_limit", prefix)] = p.HardLimit.AsMongoQuery()
	}
	if !p.InitialOutput.Empty() {
		query[fmt.Sprintf("%s.initial_output", prefix)] = p.InitialOutput.AsMongoQuery()
	}
	if !p.Output.Empty() {
		query[fmt.Sprintf("%s.output", prefix)] = p.Output.AsMongoQuery()
	}
	if !p.InitialLongOutput.Empty() {
		query[fmt.Sprintf("%s.initial_long_output", prefix)] = p.InitialLongOutput.AsMongoQuery()
	}
	if !p.LongOutput.Empty() {
		query[fmt.Sprintf("%s.long_output", prefix)] = p.LongOutput.AsMongoQuery()
	}
	if !p.LastUpdateDate.Empty() {
		query[fmt.Sprintf("%s.last_update_date", prefix)] = p.LastUpdateDate.AsMongoQuery()
	}
	if !p.LastEventDate.Empty() {
		query[fmt.Sprintf("%s.last_event_date", prefix)] = p.LastEventDate.AsMongoQuery()
	}
	if !p.Resource.Empty() {
		query[fmt.Sprintf("%s.resource", prefix)] = p.Resource.AsMongoQuery()
	}
	if !p.Resolved.Empty() {
		query[fmt.Sprintf("%s.resolved", prefix)] = p.Resolved.AsMongoQuery()
	}
	if !p.StateChangesSinceStatusUpdate.Empty() {
		query[fmt.Sprintf("%s.state_changes_since_status_update", prefix)] = p.StateChangesSinceStatusUpdate.AsMongoQuery()
	}
	if !p.TotalStateChanges.Empty() {
		query[fmt.Sprintf("%s.total_state_changes", prefix)] = p.TotalStateChanges.AsMongoQuery()
	}

	for key, value := range p.Extra {
		query[fmt.Sprintf("%s.extra.%s", prefix, key)] = value.AsMongoQuery()
	}
}

func (p AlarmValuePattern) AsMongoDriverQuery(prefix string, query mongobson.M) {
	p.ACK.AsMongoDriverQuery(fmt.Sprintf("%s.ack", prefix), query)
	p.Canceled.AsMongoDriverQuery(fmt.Sprintf("%s.canceled", prefix), query)
	p.Done.AsMongoDriverQuery(fmt.Sprintf("%s.done", prefix), query)
	p.Snooze.AsMongoDriverQuery(fmt.Sprintf("%s.snooze", prefix), query)
	p.State.AsMongoDriverQuery(fmt.Sprintf("%s.state", prefix), query)
	p.Status.AsMongoDriverQuery(fmt.Sprintf("%s.status", prefix), query)
	p.Ticket.AsMongoDriverQuery(fmt.Sprintf("%s.ticket", prefix), query)

	if !p.Component.Empty() {
		query[fmt.Sprintf("%s.component", prefix)] = p.Component.AsMongoDriverQuery()
	}
	if !p.Connector.Empty() {
		query[fmt.Sprintf("%s.connector", prefix)] = p.Connector.AsMongoDriverQuery()
	}
	if !p.ConnectorName.Empty() {
		query[fmt.Sprintf("%s.connector_name", prefix)] = p.ConnectorName.AsMongoDriverQuery()
	}
	if !p.CreationDate.Empty() {
		query[fmt.Sprintf("%s.creation_date", prefix)] = p.CreationDate.AsMongoDriverQuery()
	}
	if !p.DisplayName.Empty() {
		query[fmt.Sprintf("%s.display_name", prefix)] = p.DisplayName.AsMongoDriverQuery()
	}
	if !p.HardLimit.Empty() {
		query[fmt.Sprintf("%s.hard_limit", prefix)] = p.HardLimit.AsMongoDriverQuery()
	}
	if !p.InitialOutput.Empty() {
		query[fmt.Sprintf("%s.initial_output", prefix)] = p.InitialOutput.AsMongoDriverQuery()
	}
	if !p.Output.Empty() {
		query[fmt.Sprintf("%s.output", prefix)] = p.Output.AsMongoDriverQuery()
	}
	if !p.InitialLongOutput.Empty() {
		query[fmt.Sprintf("%s.initial_long_output", prefix)] = p.InitialLongOutput.AsMongoDriverQuery()
	}
	if !p.LongOutput.Empty() {
		query[fmt.Sprintf("%s.long_output", prefix)] = p.LongOutput.AsMongoDriverQuery()
	}
	if !p.LastUpdateDate.Empty() {
		query[fmt.Sprintf("%s.last_update_date", prefix)] = p.LastUpdateDate.AsMongoDriverQuery()
	}
	if !p.LastEventDate.Empty() {
		query[fmt.Sprintf("%s.last_event_date", prefix)] = p.LastEventDate.AsMongoDriverQuery()
	}
	if !p.Resource.Empty() {
		query[fmt.Sprintf("%s.resource", prefix)] = p.Resource.AsMongoDriverQuery()
	}
	if !p.Resolved.Empty() {
		query[fmt.Sprintf("%s.resolved", prefix)] = p.Resolved.AsMongoDriverQuery()
	}
	if !p.StateChangesSinceStatusUpdate.Empty() {
		query[fmt.Sprintf("%s.state_changes_since_status_update", prefix)] = p.StateChangesSinceStatusUpdate.AsMongoDriverQuery()
	}
	if !p.TotalStateChanges.Empty() {
		query[fmt.Sprintf("%s.total_state_changes", prefix)] = p.TotalStateChanges.AsMongoDriverQuery()
	}

	for key, value := range p.Extra {
		query[fmt.Sprintf("%s.extra.%s", prefix, key)] = value.AsMongoDriverQuery()
	}

	q := p.Parents.AsMongoDriverQuery()
	if q != nil && len(q) != 0 {
		query[fmt.Sprintf("%s.parents", prefix)] = q
	}

	q = p.Children.AsMongoDriverQuery()
	if q != nil && len(q) != 0 {
		query[fmt.Sprintf("%s.children", prefix)] = q
	}
}

// Matches returns true if an alarm value is matched by a pattern. If the
// pattern contains regular expressions with sub-expressions, the values of the
// sub-expressions are written in the matches argument.
func (p AlarmValuePattern) Matches(value types.AlarmValue, matches *AlarmValueRegexMatches) bool {
	match := p.ACK.Matches(value.ACK, &matches.ACK) &&
		p.Done.Matches(value.Done, &matches.Done) &&
		p.Snooze.Matches(value.Snooze, &matches.Snooze) &&
		p.State.Matches(value.State, &matches.State) &&
		p.Status.Matches(value.Status, &matches.Status) &&
		p.Ticket.Matches(value.Ticket, &matches.Ticket) &&
		p.Component.Matches(value.Component, &matches.Component) &&
		p.Connector.Matches(value.Connector, &matches.Connector) &&
		p.ConnectorName.Matches(value.ConnectorName, &matches.ConnectorName) &&
		p.CreationDate.Matches(value.CreationDate) &&
		p.ActivationDate.Matches(value.ActivationDate) &&
		p.DisplayName.Matches(value.DisplayName, &matches.DisplayName) &&
		p.HardLimit.Matches(value.HardLimit) &&
		p.InitialOutput.Matches(value.InitialOutput, &matches.InitialOutput) &&
		p.Output.Matches(value.Output, &matches.Output) &&
		p.InitialLongOutput.Matches(value.InitialLongOutput, &matches.InitialLongOutput) &&
		p.LongOutput.Matches(value.LongOutput, &matches.LongOutput) &&
		p.LastUpdateDate.Matches(value.LastUpdateDate) &&
		p.LastEventDate.Matches(value.LastEventDate) &&
		p.Resource.Matches(value.Resource, &matches.Resource) &&
		p.Resolved.Matches(value.Resolved) &&
		p.StateChangesSinceStatusUpdate.Matches(value.StateChangesSinceStatusUpdate) &&
		p.TotalStateChanges.Matches(value.TotalStateChanges) &&
		p.Parents.Matches(value.Parents) &&
		p.Children.Matches(value.Children)
	if !match {
		return false
	}

	for extraName, extraPattern := range p.Extra {
		var regexMatches RegexMatches
		match = extraPattern.Matches(value.Extra[extraName], &regexMatches)
		if match {
			matches.Extra[extraName] = regexMatches
		} else {
			return false
		}
	}

	return true
}

// SetBSON unmarshals a BSON value into an AlarmValuePattern.
func (p *AlarmValuePattern) SetBSON(raw mgobson.Raw) error {
	err := raw.Unmarshal(&p.AlarmValueFields)
	if err != nil {
		return err
	}

	if len(p.UnexpectedFields) != 0 {
		unexpectedFieldNames := make([]string, 0, len(p.UnexpectedFields))
		for key := range p.UnexpectedFields {
			unexpectedFieldNames = append(unexpectedFieldNames, key)
		}
		return UnexpectedFieldsError{
			Err: fmt.Errorf("Unexpected pattern fields: %s", strings.Join(unexpectedFieldNames, ", ")),
		}
	}

	return nil
}

func (p AlarmValuePattern) MarshalBSONValue() (bsontype.Type, []byte, error) {
	resultBson := mongobson.M{}

	if p.Output.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "Output", "output")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Output
	}

	if p.Snooze.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "Snooze", "snooze")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Snooze
	}

	if p.LastUpdateDate.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "LastUpdateDate", "lastupdatedate")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.LastUpdateDate
	}

	if p.Ticket.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "Ticket", "ticket")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Ticket
	}

	if p.ACK.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "ACK", "ack")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.ACK
	}

	if p.Done.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "Done", "done")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Done
	}

	if p.ActivationDate.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "ActivationDate", "activationdate")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.ActivationDate
	}

	if p.Canceled.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "Canceled", "canceled")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Canceled
	}

	if p.Component.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "Component", "component")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Component
	}

	if p.Connector.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "Connector", "connector")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Connector
	}

	if p.ConnectorName.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "ConnectorName", "connectorname")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.ConnectorName
	}

	if p.CreationDate.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "CreationDate", "creationdate")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.CreationDate
	}

	if p.DisplayName.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "DisplayName", "displayname")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.DisplayName
	}

	if len(p.Extra) > 0 {
		bsonFieldName, err := GetFieldBsonName(p, "Extra", "extra")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Extra
	}

	if p.HardLimit.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "HardLimit", "hardlimit")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.HardLimit
	}

	if p.InitialLongOutput.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "InitialLongOutput", "initiallongoutput")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.InitialLongOutput
	}

	if p.InitialOutput.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "InitialOutput", "initialoutput")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.InitialOutput
	}

	if p.LastEventDate.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "LastEventDate", "lasteventdate")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.LastEventDate
	}

	if p.LongOutput.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "LongOutput", "longoutput")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.LongOutput
	}

	if p.Resolved.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "Resolved", "resolved")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Resolved
	}

	if p.Resource.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "Resource", "resource")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Resource
	}

	if p.State.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "State", "state")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.State
	}

	if p.StateChangesSinceStatusUpdate.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "StateChangesSinceStatusUpdate", "statechangessincestatusupdate")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.StateChangesSinceStatusUpdate
	}

	if p.Status.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "Status", "status")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Status
	}

	if p.TotalStateChanges.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "TotalStateChanges", "totalstatechanges")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.TotalStateChanges
	}

	if !p.Parents.Empty() {
		bsonFieldName, err := GetFieldBsonName(p, "Parents", "parents")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Parents
	}

	if !p.Children.Empty() {
		bsonFieldName, err := GetFieldBsonName(p, "Children", "children")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Children
	}

	if len(resultBson) > 0 {
		return mongobson.MarshalValue(resultBson)
	}

	return bsontype.Undefined, nil, nil
}

func (p *AlarmValuePattern) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bsontype.EmbeddedDocument:
		err := mongobson.Unmarshal(b, &p.AlarmValueFields)
		if err != nil {
			return err
		}
		if len(p.UnexpectedFields) != 0 {
			unexpectedFieldNames := make([]string, 0, len(p.UnexpectedFields))
			for key := range p.UnexpectedFields {
				unexpectedFieldNames = append(unexpectedFieldNames, key)
			}

			return UnexpectedFieldsError{
				Err: fmt.Errorf("unexpected pattern fields: %s", strings.Join(unexpectedFieldNames, ", ")),
			}
		}
	default:
		return fmt.Errorf("alarm pattern should be a document")
	}

	return nil
}
