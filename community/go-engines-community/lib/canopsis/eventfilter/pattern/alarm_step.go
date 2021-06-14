package pattern

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"strings"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	mgobson "github.com/globalsign/mgo/bson"
	mongobson "go.mongodb.org/mongo-driver/bson"
)

// AlarmStepRegexMatches is a type that contains the values of the
// sub-expressions of regular expressions for each of the fields of an
// AlarmStep that contain strings.
type AlarmStepRegexMatches struct {
	Type    RegexMatches
	Author  RegexMatches
	Message RegexMatches
}

// AlarmStepFields is a type representing a pattern that can be applied to an
// alarm step.
// The fields are not defined directly in the AlarmStepRefPattern struct to
// make the unmarshalling easier.
type AlarmStepFields struct {
	Type      StringPattern  `bson:"_t"`
	Timestamp TimePattern    `bson:"t"`
	Author    StringPattern  `bson:"a"`
	Message   StringPattern  `bson:"m"`
	Value     IntegerPattern `bson:"val"`

	// When unmarshalling a BSON document, the fields of this document that are
	// not defined in this struct are added to UnexpectedFields.
	UnexpectedFields map[string]interface{} `bson:",inline"`
}

func (f AlarmStepFields) IsSet() bool {
	return f.Type.IsSet() ||
		f.Timestamp.IsSet() ||
		f.Author.IsSet() ||
		f.Message.IsSet() ||
		f.Value.IsSet()
}

// AlarmStepRefPattern is a type representing a pattern that can be applied to
// a reference to an alarm step.
type AlarmStepRefPattern struct {
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

	AlarmStepFields
}

// Empty returns true if the pattern has not been set
func (p AlarmStepRefPattern) Empty() bool {
	return (!p.ShouldNotBeNil && !p.ShouldBeNil)
}

func (p AlarmStepRefPattern) IsSet() bool {
	return p.ShouldBeNil || p.AlarmStepFields.IsSet()
}

// AsMongoQuery returns a mongodb filter corresponding to the AlarmStepRefPattern.
func (p AlarmStepRefPattern) AsMongoQuery(prefix string, query mgobson.M) {
	if p.ShouldBeNil {
		query[prefix] = nil
		return
	}

	if !p.Type.Empty() {
		query[fmt.Sprintf("%s._t", prefix)] = p.Type.AsMongoQuery()
	}

	if !p.Timestamp.Empty() {
		query[fmt.Sprintf("%s.t", prefix)] = p.Timestamp.AsMongoQuery()
	}

	if !p.Author.Empty() {
		query[fmt.Sprintf("%s.a", prefix)] = p.Author.AsMongoQuery()
	}

	if !p.Message.Empty() {
		query[fmt.Sprintf("%s.m", prefix)] = p.Message.AsMongoQuery()
	}

	if !p.Value.Empty() {
		query[fmt.Sprintf("%s.val", prefix)] = p.Value.AsMongoQuery()
	}
}

// AsMongoQuery returns a mongodb filter corresponding to the AlarmStepRefPattern.
func (p AlarmStepRefPattern) AsMongoDriverQuery(prefix string, query mongobson.M) {
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
		query[fmt.Sprintf("%s.val", prefix)] = p.Value.AsMongoDriverQuery()
	}
}

// Matches returns true if an alarm step is matched by a pattern. If the
// pattern contains regular expressions with sub-expressions, the values of the
// sub-expressions are written in the matches argument.
func (p AlarmStepRefPattern) Matches(step *types.AlarmStep, matches *AlarmStepRegexMatches) bool {
	if p.ShouldBeNil {
		return step == nil
	}

	if p.ShouldNotBeNil {
		return step != nil &&
			p.Type.Matches(step.Type, &matches.Type) &&
			p.Timestamp.Matches(step.Timestamp) &&
			p.Author.Matches(step.Author, &matches.Author) &&
			p.Message.Matches(step.Message, &matches.Message) &&
			p.Value.Matches(step.Value)
	}

	return true
}

// SetBSON unmarshals a BSON value into an AlarmStepRefPattern.
func (p *AlarmStepRefPattern) SetBSON(raw mgobson.Raw) error {
	switch raw.Kind {
	case mgobson.ElementNil:
		// The BSON value is null. The field should not be set.
		p.ShouldBeNil = true
		p.ShouldNotBeNil = false
		return nil

	default:
		err := raw.Unmarshal(&p.AlarmStepFields)
		if err != nil {
			return err
		}

		if len(p.UnexpectedFields) != 0 {
			unexpectedFieldNames := make([]string, 0, len(p.UnexpectedFields))
			for key := range p.UnexpectedFields {
				unexpectedFieldNames = append(unexpectedFieldNames, key)
			}
			return fmt.Errorf("Unexpected pattern fields: %s", strings.Join(unexpectedFieldNames, ", "))
		}

		p.ShouldBeNil = false
		p.ShouldNotBeNil = true
		return nil
	}
}

func (p AlarmStepRefPattern) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if p.ShouldBeNil {
		return bsontype.Null, []byte{}, nil
	}

	resultBson := mongobson.M{}

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

	if len(resultBson) > 0 {
		return mongobson.MarshalValue(resultBson)
	}

	return bsontype.Undefined, nil, nil
}

func (p *AlarmStepRefPattern) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bsontype.Null:
		// The BSON value is null. The field should not be set.
		p.ShouldBeNil = true
		p.ShouldNotBeNil = false
		return nil
	case bsontype.EmbeddedDocument:
		err := mongobson.Unmarshal(b, &p.AlarmStepFields)
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
		return fmt.Errorf("alarm step pattern should be a document or nil")
	}

	p.ShouldBeNil = false
	p.ShouldNotBeNil = true
	return nil
}
