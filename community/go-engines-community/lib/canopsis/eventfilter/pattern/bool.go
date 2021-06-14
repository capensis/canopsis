package pattern

import (
	"git.canopsis.net/canopsis/go-engines/lib/utils"
	mgobson "github.com/globalsign/mgo/bson"
	mongobson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// BoolPattern is a type representing a pattern that can be applied to the value
// of a field of an event that contains a boolean.
type BoolPattern struct {
	utils.OptionalBool
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
		return mongobson.MarshalValue(p.OptionalBool.Value)
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

// AsMongoQuery returns a mongodb filter from the BoolPattern for mgo-driver
func (p BoolPattern) AsMongoQuery() mgobson.M {
	var query mgobson.M
	query = make(mgobson.M)

	if p.Set {
		query["$eq"] = p.Value
	}

	return query
}

// AsMongoDriverQuery returns a mongodb filter from the BoolPattern for mongo-driver
func (p BoolPattern) AsMongoDriverQuery() mongobson.M {
	var query mongobson.M
	query = make(mongobson.M)

	if p.Set {
		query["$eq"] = p.Value
	}

	return query
}
