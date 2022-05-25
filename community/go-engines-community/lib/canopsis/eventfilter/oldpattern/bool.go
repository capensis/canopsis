package oldpattern

import (
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// BoolPattern is a type representing a pattern that can be applied to the value
// of a field of an event that contains a boolean.
type BoolPattern struct {
	types.OptionalBool
}

// Matches returns true if the value is matched by the pattern.
func (p BoolPattern) Matches(value bool) bool {
	if p.Set {
		if !(value == p.Value) {
			return false
		}
	}

	return true
}

func (p BoolPattern) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if p.OptionalBool.Set {
		return bson.MarshalValue(p.OptionalBool.Value)
	}

	return bsontype.Undefined, nil, nil
}

func (p BoolPattern) IsSet() bool {
	return p.Set
}

// Empty returns true if the condition has not been set.
func (p BoolPattern) Empty() bool {
	return !p.Set
}

// AsMongoDriverQuery returns a mongodb filter from the BoolPattern for mongo-driver
func (p BoolPattern) AsMongoDriverQuery() bson.M {
	query := make(bson.M)

	if p.Set {
		query["$eq"] = p.Value
	}

	return query
}

func (p BoolPattern) AsSqlQuery() string {
	if p.Set {
		return fmt.Sprintf("= %v", p.Value)
	}

	return ""
}
