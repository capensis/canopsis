package oldpattern

import (
	"fmt"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// AlarmTicketRegexMatches is a type that contains the values of the
// sub-expressions of regular expressions for each of the fields of an
// AlarmTicket that contain strings.
type AlarmTicketRegexMatches struct {
	Type    RegexMatches
	Author  RegexMatches
	Message RegexMatches
	Value   RegexMatches
	Data    map[string]RegexMatches
}

func NewAlarmTicketRegexMatches() AlarmTicketRegexMatches {
	return AlarmTicketRegexMatches{
		Data: make(map[string]RegexMatches),
	}
}

// AlarmTicketFields is a type representing a pattern that can be applied to an
// alarm ticket step.
// The fields are not defined directly in the AlarmTicketRefPattern struct to
// make the unmarshalling easier.
type AlarmTicketFields struct {
	Type      StringPattern            `bson:"_t"`
	Timestamp TimePattern              `bson:"t"`
	Author    StringPattern            `bson:"a"`
	Message   StringPattern            `bson:"m"`
	Value     StringPattern            `bson:"val"`
	Data      map[string]StringPattern `bson:"data"`

	// When unmarshalling a BSON document, the fields of this document that are
	// not defined in this struct are added to UnexpectedFields.
	UnexpectedFields map[string]interface{} `bson:",inline"`
}

func (f AlarmTicketFields) IsSet() bool {
	return f.Type.IsSet() ||
		f.Timestamp.IsSet() ||
		f.Author.IsSet() ||
		f.Message.IsSet() ||
		f.Value.IsSet() ||
		len(f.Data) > 0
}

// AlarmTicketRefPattern is a type representing a pattern that can be applied
// to a reference to an alarm ticket step.
type AlarmTicketRefPattern struct {
	// ShouldNotBeNil is a boolean indicating that the alarm step should not be
	// nil, and ShouldBeNil is a boolean indicating that the alarm step should
	// be nil.
	// The two booleans are needed to be able to make the difference between
	// the case where no pattern was defined (in which case the alarm step can
	// be nil or not), the case where a nil pattern was defined (in which case
	// the alarm step should be nil), and the case where a non-nil pattern was
	// defined (in which case the alarm step should not be nil).
	ShouldNotBeNil bool
	ShouldBeNil    bool

	AlarmTicketFields
}

func (p AlarmTicketRefPattern) IsSet() bool {
	return p.ShouldBeNil || p.AlarmTicketFields.IsSet()
}

// Empty returns true if the pattern has not been set
func (p AlarmTicketRefPattern) Empty() bool {
	return (!p.ShouldNotBeNil && !p.ShouldBeNil)
}

func (p AlarmTicketRefPattern) AsMongoDriverQuery(prefix string, query bson.M) {
	if p.ShouldBeNil {
		query[prefix] = nil
		return
	}

	if !p.Type.Empty() {
		query[fmt.Sprintf("%s._t", prefix)] = p.Type.AsMongoDriverQuery()
	}

	if !p.Timestamp.Empty() {
		query[fmt.Sprintf("%s.t", prefix)] = p.Timestamp.AsMongoDriverQuery()
	}

	if !p.Author.Empty() {
		query[fmt.Sprintf("%s.a", prefix)] = p.Author.AsMongoDriverQuery()
	}

	if !p.Message.Empty() {
		query[fmt.Sprintf("%s.m", prefix)] = p.Message.AsMongoDriverQuery()
	}

	if !p.Value.Empty() {
		query[fmt.Sprintf("%s.ticket", prefix)] = p.Value.AsMongoDriverQuery()
	}

	if len(p.Data) > 0 {
		for k, v := range p.Data {
			query[fmt.Sprintf("%s.ticket_data.%s", prefix, k)] = v.AsMongoDriverQuery()
		}
	}
}

// Matches returns true if an alarm ticket step is matched by a pattern. If the
// pattern contains regular expressions with sub-expressions, the values of the
// sub-expressions are written in the matches argument.
func (p AlarmTicketRefPattern) Matches(step *types.AlarmStep, matches *AlarmTicketRegexMatches) bool {
	if p.ShouldBeNil {
		return step == nil
	}

	if p.ShouldNotBeNil {
		if step == nil {
			return false
		}

		fieldsMatch :=
			p.Type.Matches(step.Type, &matches.Type) &&
				p.Timestamp.Matches(step.Timestamp) &&
				p.Author.Matches(step.Author, &matches.Author) &&
				p.Message.Matches(step.Message, &matches.Message) &&
				p.Value.Matches(step.TicketInfo.Ticket, &matches.Value)

		dataMatch := true

		for dataKey, stringPattern := range p.Data {
			var regexMatches RegexMatches

			if stringPattern.Matches(step.TicketInfo.TicketData[dataKey], &regexMatches) {
				matches.Data[dataKey] = regexMatches
			} else {
				dataMatch = false
				break
			}
		}

		return fieldsMatch && dataMatch
	}

	return true
}

func (p AlarmTicketRefPattern) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if p.ShouldBeNil {
		return bsontype.Null, []byte{}, nil
	}

	resultBson := bson.M{}

	if p.Type.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "Type", "type")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Type
	}

	if p.Timestamp.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "Timestamp", "timestamp")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Timestamp
	}

	if p.Author.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "Author", "author")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Author
	}

	if p.Message.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "Message", "message")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Message
	}

	if p.Value.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "Value", "value")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Value
	}

	if len(p.Data) > 0 {
		bsonFieldName, err := GetFieldBsonName(p, "Data", "data")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Data
	}

	if len(resultBson) > 0 {
		return bson.MarshalValue(resultBson)
	}

	return bsontype.Undefined, nil, nil
}

func (p *AlarmTicketRefPattern) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bsontype.Null:
		// The BSON value is null. The field should not be set.
		p.ShouldBeNil = true
		p.ShouldNotBeNil = false
		return nil
	case bsontype.EmbeddedDocument:
		err := bson.Unmarshal(b, &p.AlarmTicketFields)
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
	default:
		return fmt.Errorf("alarm ticket pattern should be a document or nil")
	}

	p.ShouldBeNil = false
	p.ShouldNotBeNil = true
	return nil
}
