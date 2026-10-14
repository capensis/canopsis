package mongoquery

import (
	"regexp"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// GetSearchQuery returns search query, it returns parsed search filter if it can be parsed
// or regex search filter by searchBy otherwise.
func GetSearchQuery(search string, searchBy []string) bson.M {
	if search == "" {
		return nil
	}

	searchRegexp := bson.Regex{
		Pattern: GetSearchPattern(search),
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

func GetSearchPattern(search string) string {
	if _, err := regexp.Compile(search); err == nil {
		return search
	}

	return ".*" + regexp.QuoteMeta(search) + ".*"
}
