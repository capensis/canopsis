package author

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type Author struct {
	ID   string `bson:"_id" json:"_id"`
	Name string `bson:"name" json:"name"`
}

type Role struct {
	ID   string `bson:"_id" json:"_id"`
	Name string `bson:"name" json:"name"`
}

func Pipeline() []bson.M {
	return PipelineForField("author")
}

func PipelineForField(field string) []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.RightsMongoCollection,
			"localField":   field,
			"foreignField": "_id",
			"as":           field,
		}},
		{"$unwind": bson.M{"path": "$" + field, "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			field + ".name": "$" + field + ".crecord_name",
		}},
		{"$addFields": bson.M{
			field: bson.M{"$cond": bson.M{
				"if":   "$" + field + "._id",
				"then": "$" + field,
				"else": "$$REMOVE",
			}},
		}},
	}
}
