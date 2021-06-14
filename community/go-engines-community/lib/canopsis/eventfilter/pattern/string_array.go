package pattern

import (
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/utils"
	mgobson "github.com/globalsign/mgo/bson"
	mongobson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"strings"
)

type StringArrayConditions struct {
	HasEvery			utils.OptionalStringArray	`bson:"has_every,omitempty"`
	HasOneOf			utils.OptionalStringArray	`bson:"has_one_of,omitempty"`
	HasNot 				utils.OptionalStringArray 	`bson:"has_not,omitempty"`
	UnexpectedFields 	map[string]interface{} 		`bson:",inline"`
}

// Matches returns true if the value satisfies each of the conditions defined
// in the IntegerConditions.
func (p StringArrayConditions) Matches(value []string) bool {
	valueMap := make(map[string]bool)

	for _, v := range value {
		valueMap[v] = true
	}

	hasEveryCondition := true
	hasOneOfCondition := true
	hasNotCondition := true

	if p.HasEvery.Set {
		for _, v := range p.HasEvery.Value {
			_, exists := valueMap[v]
			if !exists {
				hasEveryCondition = false
			}
		}
	}

	if p.HasOneOf.Set {
		hasOneOfCondition = false
		for _, v := range p.HasOneOf.Value {
			_, exists := valueMap[v]
			if exists {
				hasOneOfCondition = true
				break
			}
		}
	}

	if p.HasNot.Set {
		for _, v := range p.HasNot.Value {
			_, exists := valueMap[v]
			if exists {
				hasNotCondition = false
			}
		}
	}

	return hasEveryCondition && hasOneOfCondition && hasNotCondition
}

func (p StringArrayConditions) Empty() bool {
	return !(p.HasEvery.Set || p.HasNot.Set || p.HasOneOf.Set)
}

func (p StringArrayConditions) OnlyHasNotCondition() bool {
	return !p.HasEvery.Set && !p.HasOneOf.Set && p.HasNot.Set
}

func (p StringArrayConditions) AsMongoQuery() mgobson.M {
	var query mgobson.M
	query = make(mgobson.M)
	if p.HasOneOf.Set {
		query["$in"] = p.HasOneOf.Value
	}

	if p.HasNot.Set {
		query["$nin"] = p.HasNot.Value
	}

	if p.HasEvery.Set {
		query["$all"] = p.HasEvery.Value
	}

	return query
}

func (p StringArrayConditions) AsMongoDriverQuery() mongobson.M {
	var query mongobson.M
	query = make(mongobson.M)
	if p.HasOneOf.Set {
		query["$in"] = p.HasOneOf.Value
	}

	if p.HasNot.Set {
		query["$nin"] = p.HasNot.Value
	}

	if p.HasEvery.Set {
		query["$all"] = p.HasEvery.Value
	}

	return query
}

type StringArrayPattern struct {
	StringArrayConditions
}

func (p *StringArrayPattern) SetBSON(raw mgobson.Raw) error {
	switch raw.Kind {
	case mgobson.ElementDocument:
		err := raw.Unmarshal(&p.StringArrayConditions)
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
		return nil

	default:
		return fmt.Errorf("A string array patter should be a document")
	}
}

func (p StringArrayPattern) MarshalBSONValue() (bsontype.Type, []byte, error) {
	resultBson := mongobson.M{}

	if p.HasEvery.Set {
		bsonFieldName, err := GetFieldBsonName(p, "HasEvery", "hasevery")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.HasEvery.Value
	}

	if p.HasNot.Set {
		bsonFieldName, err := GetFieldBsonName(p, "HasNot", "hasnot")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.HasNot.Value
	}

	if p.HasOneOf.Set {
		bsonFieldName, err := GetFieldBsonName(p, "HasOneOf", "hasoneof")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.HasOneOf.Value
	}

	if len(resultBson) > 0 {
		return mongobson.MarshalValue(resultBson)
	}

	return bsontype.Undefined, nil, nil
}

func (p *StringArrayPattern) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bsontype.EmbeddedDocument:
		err := mongobson.Unmarshal(b, &p.StringArrayConditions)
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
		return nil

	default:
		return fmt.Errorf("A string array patter should be a document")
	}
}


