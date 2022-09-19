package oldpattern

import (
	"fmt"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type StringArrayConditions struct {
	HasEvery         types.OptionalStringArray `bson:"has_every,omitempty"`
	HasOneOf         types.OptionalStringArray `bson:"has_one_of,omitempty"`
	HasNot           types.OptionalStringArray `bson:"has_not,omitempty"`
	IsEmpty          types.OptionalBool        `bson:"is_empty,omitempty"`
	UnexpectedFields map[string]interface{}    `bson:",inline"`
}

// Matches returns true if the value satisfies each of the conditions defined
// in the StringArrayConditions.
func (p StringArrayConditions) Matches(value []string) bool {
	valueMap := make(map[string]bool)

	for _, v := range value {
		valueMap[v] = true
	}

	hasEveryCondition := true
	hasOneOfCondition := true
	hasNotCondition := true
	isEmptyCondition := true

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

	if p.IsEmpty.Set {
		isEmptyCondition = !p.IsEmpty.Value && len(value) != 0 || p.IsEmpty.Value && len(value) == 0
	}

	return hasEveryCondition && hasOneOfCondition && hasNotCondition && isEmptyCondition
}

func (p StringArrayConditions) Empty() bool {
	return !(p.HasEvery.Set || p.HasNot.Set || p.HasOneOf.Set || p.IsEmpty.Set)
}

func (p StringArrayConditions) OnlyHasNotCondition() bool {
	return !p.HasEvery.Set && !p.HasOneOf.Set && !p.IsEmpty.Set && p.HasNot.Set
}

func (p StringArrayConditions) OnlyIsEmpty() bool {
	return !p.HasEvery.Set && !p.HasOneOf.Set && !p.HasNot.Set && p.IsEmpty.Set && p.IsEmpty.Value
}

func (p StringArrayConditions) AsMongoDriverQuery() bson.M {
	query := make(bson.M)
	if p.HasOneOf.Set {
		query["$in"] = p.HasOneOf.Value
	}

	if p.HasNot.Set {
		query["$nin"] = p.HasNot.Value
	}

	if p.HasEvery.Set {
		query["$all"] = p.HasEvery.Value
	}

	if p.IsEmpty.Set {
		if p.IsEmpty.Value {
			query["$eq"] = bson.A{}
		} else {
			query["$exists"] = true
			query["$ne"] = bson.A{}
		}
	}

	return query
}

type StringArrayPattern struct {
	StringArrayConditions
}

type fieldBsonValue struct {
	isFieldSet bool
	fieldName  string
	bsonName   string
	value      interface{}
}

func (p StringArrayPattern) MarshalBSONValue() (bsontype.Type, []byte, error) {
	resultBson := bson.M{}

	for _, v := range []fieldBsonValue{
		{p.HasEvery.Set, "HasEvery", "hasevery", p.HasEvery.Value},
		{p.HasNot.Set, "HasNot", "hasnot", p.HasNot.Value},
		{p.HasOneOf.Set, "HasOneOf", "hasoneof", p.HasOneOf.Value},
		{p.IsEmpty.Set, "IsEmpty", "isempty", p.IsEmpty.Value},
	} {
		if v.isFieldSet {
			bsonFieldName, err := GetFieldBsonName(p, v.fieldName, v.bsonName)
			if err != nil {
				return bsontype.Undefined, nil, err
			}

			resultBson[bsonFieldName] = v.value
		}
	}

	if len(resultBson) > 0 {
		return bson.MarshalValue(resultBson)
	}

	return bsontype.Undefined, nil, nil
}

func (p *StringArrayPattern) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bsontype.EmbeddedDocument:
		err := bson.Unmarshal(b, &p.StringArrayConditions)
		if err != nil {
			return err
		}
		if len(p.UnexpectedFields) != 0 {
			unexpectedFieldNames := make([]string, 0, len(p.UnexpectedFields))
			for key := range p.UnexpectedFields {
				unexpectedFieldNames = append(unexpectedFieldNames, key)
			}
			return fmt.Errorf("unexpected pattern fields: %s", strings.Join(unexpectedFieldNames, ", "))
		}
		return nil

	default:
		return fmt.Errorf("a string array pattern should be a document")
	}
}
