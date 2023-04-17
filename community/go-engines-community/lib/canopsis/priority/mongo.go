package priority

import "go.mongodb.org/mongo-driver/bson"

func GetSortPipeline() []bson.M {
	return []bson.M{
		{"$addFields": bson.M{
			"has_priority": bson.M{"$gt": bson.A{"$priority", 0}},
		}},
		{"$sort": bson.D{
			{Key: "has_priority", Value: -1},
			{Key: "priority", Value: 1},
			{Key: "_id", Value: 1},
		}},
	}
}
