package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
)

const (
	SortAsc  = "asc"
	SortDesc = "desc"
)

func GetSort(sortBy, sort string) bson.D {
	sortDir := 1
	if sort == SortDesc {
		sortDir = -1
	}

	if sortBy == "" {
		sortBy = "_id"
	}

	q := bson.D{{Key: sortBy, Value: sortDir}}
	if sortBy != "_id" {
		q = append(q, bson.E{Key: "_id", Value: 1})
	}

	return q
}
