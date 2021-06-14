package common

import "go.mongodb.org/mongo-driver/bson"

func GetSortQuery(sortBy, sort string) bson.M {
	sortDir := 1
	if sort == SortDesc {
		sortDir = -1
	}

	q := bson.D{{Key: sortBy, Value: sortDir}}
	if sortBy != "_id" {
		q = append(q, bson.E{Key: "_id", Value: 1})
	}

	return bson.M{"$sort": q}
}
