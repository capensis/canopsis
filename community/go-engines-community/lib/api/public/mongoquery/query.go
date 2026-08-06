package mongoquery

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// SortQuery returns sort query which support consistent sort : sort by "_id" is added
// if sortBy is not "_id".
func SortQuery(sortBy, sort string) bson.M {
	sortDir := 1
	if sort == pagination.SortDesc {
		sortDir = -1
	}

	q := bson.D{{Key: sortBy, Value: sortDir}}
	if sortBy != "_id" {
		q = append(q, bson.E{Key: "_id", Value: 1})
	}

	return bson.M{"$sort": q}
}

func ReferencePipeline(field, collection string) []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         collection,
			"localField":   field,
			"foreignField": "_id",
			"as":           field,
			"pipeline": []bson.M{
				{"$project": bson.M{
					"_id":  1,
					"name": 1,
				}},
			},
		}},
		{"$unwind": bson.M{"path": "$" + field, "preserveNullAndEmptyArrays": true}},
	}
}

func AuthorPipeline() []bson.M {
	return ReferencePipeline("author", mongo.UserCollection)
}

func ArrayReferencePipeline(field, collection string) []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         collection,
			"localField":   field,
			"foreignField": "_id",
			"as":           field,
			"pipeline": []bson.M{
				{"$project": bson.M{
					"_id":  1,
					"name": 1,
				}},
			},
		}},
	}
}
