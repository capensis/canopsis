package common

import (
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/expression/parser"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetSortQuery returns sort query which support consistent sort : sort by "_id" is added
// if sortBy is not "_id".
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

// GetSearchQuery returns search query, it returns parsed search filter if it can be parsed
// or regex search filter by searchBy otherwise.
func GetSearchQuery(search string, searchBy []string) bson.M {
	if search == "" {
		return nil
	}

	p := parser.NewParser()
	expr, err := p.Parse(search)
	if err == nil {
		return expr.Query()
	}

	searchRegexp := primitive.Regex{
		Pattern: fmt.Sprintf(".*%s.*", search),
		Options: "i",
	}

	searchMatch := make([]bson.M, len(searchBy))
	for i := range searchBy {
		searchMatch[i] = bson.M{searchBy[i]: searchRegexp}
	}

	return bson.M{
		"$or": searchMatch,
	}
}
