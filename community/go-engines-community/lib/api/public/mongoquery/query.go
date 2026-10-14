package mongoquery

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/mongoquery"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// GetSearchQuery returns search query, it returns parsed search filter if it can be parsed
// or regex search filter by searchBy otherwise.
func GetSearchQuery(search string, searchBy []string) bson.M {
	return mongoquery.GetSearchQuery(search, searchBy)
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
		{"$unwind": bson.M{
			"path":                       "$" + field,
			"preserveNullAndEmptyArrays": true,
			"includeArrayIndex":          field + "_index",
		}},
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
		{"$sort": bson.M{field + "_index": 1}},
		{"$group": bson.M{
			"_id":  "$_id",
			"data": bson.M{"$first": "$$ROOT"},
			field: bson.M{"$push": bson.M{
				"$cond": bson.M{
					"if":   "$" + field + "._id",
					"then": "$" + field,
					"else": "$$REMOVE",
				},
			}},
		}},
		{"$replaceRoot": bson.M{"newRoot": bson.M{"$mergeObjects": bson.A{
			"$data",
			bson.M{field: "$" + field},
		}}}},
	}
}
