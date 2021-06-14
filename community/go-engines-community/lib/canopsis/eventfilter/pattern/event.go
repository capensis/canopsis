package pattern

import (
	"fmt"
	"log"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	mgobson "github.com/globalsign/mgo/bson"
	mongobson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// EventPattern is a type representing a pattern that can be applied to an
// event
type EventPattern struct {
	ID            StringRefPattern            `bson:"_id"`
	Connector     StringPattern               `bson:"connector"`
	ConnectorName StringPattern               `bson:"connector_name"`
	EventType     StringPattern               `bson:"event_type"`
	Component     StringPattern               `bson:"component"`
	Resource      StringPattern               `bson:"resource"`
	PerfData      StringRefPattern            `bson:"perf_data"`
	Status        IntegerRefPattern           `bson:"status"`
	Timestamp     TimePattern                 `bson:"timestamp"`
	StateType     IntegerRefPattern           `bson:"state_type"`
	SourceType    StringPattern               `bson:"source_type"`
	LongOutput    StringPattern               `bson:"long_output"`
	State         IntegerPattern              `bson:"state"`
	Output        StringPattern               `bson:"output"`
	Entity        EntityPattern               `bson:"current_entity"`
	Author        StringPattern               `bson:"author"`
	RK            StringPattern               `bson:"routing_key"`
	AckResources  BoolPattern                 `bson:"ack_resources"`
	Duration      IntegerRefPattern           `bson:"duration"`
	Ticket        StringPattern               `bson:"ticket"`
	StatName      StringPattern               `bson:"stat_name"`
	Debug         BoolPattern                 `bson:"debug"`
	ExtraInfos    map[string]InterfacePattern `bson:",inline"`
}

// EventRegexMatches is a type that contains the values of the sub-expressions
// of regular expressions for each of the fields of an Event that contain
// strings.
type EventRegexMatches struct {
	ID            RegexMatches
	Connector     RegexMatches
	ConnectorName RegexMatches
	EventType     RegexMatches
	Component     RegexMatches
	Resource      RegexMatches
	PerfData      RegexMatches
	SourceType    RegexMatches
	LongOutput    RegexMatches
	Output        RegexMatches
	Entity        EntityRegexMatches
	Author        RegexMatches
	RK            RegexMatches
	Ticket        RegexMatches
	StatName      RegexMatches
	ExtraInfos    map[string]RegexMatches
}

// NewEventRegexMatches creates an EventRegexMatches, with the Entity field
// initialized.
func NewEventRegexMatches() EventRegexMatches {
	return EventRegexMatches{
		Entity:     NewEntityRegexMatches(),
		ExtraInfos: make(map[string]RegexMatches),
	}
}

// Matches returns true if an event is matched by a pattern. If the pattern
// contains regular expressions with sub-expressions, the values of the
// sub-expressions are written in the matches argument.
func (p EventPattern) Matches(event types.Event, matches *EventRegexMatches) bool {
	match := p.ID.Matches(event.ID, &matches.ID) &&
		p.Connector.Matches(event.Connector, &matches.Connector) &&
		p.ConnectorName.Matches(event.ConnectorName, &matches.ConnectorName) &&
		p.EventType.Matches(event.EventType, &matches.EventType) &&
		p.Component.Matches(event.Component, &matches.Component) &&
		p.Resource.Matches(event.Resource, &matches.Resource) &&
		p.PerfData.Matches(event.PerfData, &matches.PerfData) &&
		p.Status.Matches(event.Status) &&
		p.Timestamp.Matches(event.Timestamp) &&
		p.StateType.Matches(event.StateType) &&
		p.SourceType.Matches(event.SourceType, &matches.SourceType) &&
		p.LongOutput.Matches(event.LongOutput, &matches.LongOutput) &&
		p.State.Matches(event.State) &&
		p.Output.Matches(event.Output, &matches.Output) &&
		p.Author.Matches(event.Author, &matches.Author) &&
		p.RK.Matches(event.RK, &matches.RK) &&
		p.AckResources.Matches(event.AckResources) &&
		p.Duration.Matches(event.Duration) &&
		p.Ticket.Matches(event.Ticket, &matches.Ticket) &&
		p.StatName.Matches(event.StatName, &matches.StatName) &&
		p.Debug.Matches(event.Debug) &&
		p.Entity.Matches(event.Entity, &matches.Entity)
	if !match {
		return false
	}

	for infoName, infoPattern := range p.ExtraInfos {
		var regexMatches RegexMatches
		match = infoPattern.Matches(event.ExtraInfos[infoName], &regexMatches)
		if match {
			matches.ExtraInfos[infoName] = regexMatches
		} else {
			return false
		}
	}

	return true
}

func (p EventPattern) MarshalBSONValue() (bsontype.Type, []byte, error) {
	resultBson := mongobson.M{}

	if p.ID.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "ID", "id")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.ID
	}

	if p.Entity.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "Entity", "entity")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Entity
	}

	if len(p.ExtraInfos) > 0 {
		for k, v := range p.ExtraInfos {
			resultBson[k] = v
		}
	}

	if p.Resource.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "Resource", "resource")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Resource
	}

	if p.ConnectorName.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "ConnectorName", "connectorname")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.ConnectorName
	}

	if p.Component.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "Component", "component")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Component
	}

	if p.State.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "State", "state")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.State
	}

	if p.Status.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "Status", "status")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Status
	}

	if p.LongOutput.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "LongOutput", "longoutput")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.LongOutput
	}

	if p.ConnectorName.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "ConnectorName", "connectorname")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.ConnectorName
	}

	if p.Ticket.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "Ticket", "ticket")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Ticket
	}

	if p.Author.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "Author", "author")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Author
	}

	if p.Timestamp.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "Timestamp", "timestamp")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Timestamp
	}

	if p.Output.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "Output", "output")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Output
	}

	if p.Debug.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "Debug", "debug")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Debug
	}

	if p.EventType.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "EventType", "eventtype")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.EventType
	}

	if p.AckResources.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "AckResources", "ackresources")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.AckResources
	}

	if p.Duration.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "Duration", "duration")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Duration
	}

	if p.PerfData.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "PerfData", "perfdata")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.PerfData
	}

	if p.RK.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "RK", "rk")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.RK
	}

	if p.SourceType.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "SourceType", "sourcetype")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.SourceType
	}

	if p.StateType.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "StateType", "statetype")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.StateType
	}

	if p.StatName.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "StatName", "statname")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.StatName
	}

	if len(resultBson) > 0 {
		return mongobson.MarshalValue(resultBson)
	}

	return bsontype.Undefined, nil, nil
}

// EventPatternList is a type representing a list of event patterns.
// An event is matched by an EventPatternList if it is matched by one of its
// EventPatterns.
// The zero value of an EventPatternList (i.e. an EventPatternList that has not
// been set) is considered valid, and matches all events.
type EventPatternList struct {
	Patterns []EventPattern

	// isSet is a boolean indicating whether the EventPatternList has been set
	// explicitly or not.
	Set bool

	// isValid is a boolean indicating whether the event patterns or valid or
	// not.
	// isValid is also false if the EventPatternList has not been set.
	Valid bool
}

// Matches returns true if the event is matched by the EventPatternList.
func (l EventPatternList) GetRegexMatches(event types.Event) (EventRegexMatches, bool) {
	regexMatches := NewEventRegexMatches()

	if !l.Set {
		return regexMatches, true
	}

	for _, pattern := range l.Patterns {
		if pattern.Matches(event, &regexMatches) {
			return regexMatches, true
		}
	}

	return regexMatches, false
}

func (l EventPatternList) Matches(event types.Event) bool {
	if !l.Set {
		return true
	}

	regexMatches := NewEventRegexMatches()
	for _, pattern := range l.Patterns {
		if pattern.Matches(event, &regexMatches) {
			return true
		}
	}

	return false
}

// IsSet returns true if the EventPatternList has been set explicitly.
func (l EventPatternList) IsSet() bool {
	return l.Set
}

// IsValid returns true if the EventPatternList is valid.
func (l EventPatternList) IsValid() bool {
	return !l.Set || l.Valid
}

// SetBSON unmarshals a BSON value into a EventPatternList.
// If it cannot be unmarshalled, it is marked as invalid.
func (l *EventPatternList) SetBSON(raw mgobson.Raw) error {
	l.Set = true

	err := raw.Unmarshal(&l.Patterns)
	if err != nil {
		log.Printf("unable to parse event pattern list: %v", err)
		l.Valid = false
		return nil
	}

	l.Valid = true
	return nil
}

func (l EventPatternList) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if !l.Set {
		return bsontype.Null, nil, nil
	}

	return mongobson.MarshalValue(l.Patterns)
}

func (l *EventPatternList) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bsontype.Null:
		l.Set = false
		l.Valid = false
	case bsontype.Array:
		l.Set = true
		l.Valid = false

		var raw mongobson.Raw
		err := mongobson.Unmarshal(b, &raw)
		if err != nil {
			return err
		}

		array, err := raw.Values()
		if err != nil {
			return err
		}

		for _, v := range array {
			document, ok := v.DocumentOK()
			if !ok {
				return fmt.Errorf("unable to parse event pattern list element")
			}

			var pattern EventPattern

			err = mongobson.Unmarshal(document, &pattern)
			if err != nil {
				if _, ok = err.(UnexpectedFieldsError); ok {
					return nil
				}

				return err
			}

			l.Patterns = append(l.Patterns, pattern)
		}

		l.Valid = true
	default:
		return fmt.Errorf("event pattern list should be an array or nil")
	}

	return nil
}
