package mongoquery

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/mongoquery"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/pagination"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// GetSortQuery returns sort query which support consistent sort : sort by "_id" is added
// if sortBy is not "_id".
func GetSortQuery(sortBy, sort string) bson.M {
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

// GetSearchQuery returns search query, it returns parsed search filter if it can be parsed
// or regex search filter by searchBy otherwise.
func GetSearchQuery(search string, searchBy []string) bson.M {
	return mongoquery.GetSearchQuery(search, searchBy)
}
