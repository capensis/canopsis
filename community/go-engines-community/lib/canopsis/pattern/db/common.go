package db

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"go.mongodb.org/mongo-driver/bson"
)

func getTypeMongoQuery(f, ft string) []bson.M {
	var conds []bson.M

	switch ft {
	case pattern.FieldTypeString:
		conds = []bson.M{{f: bson.M{"$type": "string"}}}
	case pattern.FieldTypeInt:
		conds = []bson.M{{f: bson.M{"$type": bson.A{"long", "int", "decimal"}}}}
	case pattern.FieldTypeBool:
		conds = []bson.M{{f: bson.M{"$type": "bool"}}}
	case pattern.FieldTypeStringArray:
		// Cond {"$type": "string"} checks only if an array contains at least one string element,
		// other elements can be any type.
		conds = []bson.M{{f: bson.M{"$type": "array"}}, {f: bson.M{"$type": "string"}}}
	}

	return conds
}
