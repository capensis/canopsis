package pagination

import "go.mongodb.org/mongo-driver/v2/bson"

const (
	SortAsc  = "asc"
	SortDesc = "desc"
)

// Query is base request for pagination.
type Query struct {
	Page  int64 `form:"page,default:1" json:"page" binding:"numeric,gt=0"`
	Limit int64 `form:"limit,default:10" json:"limit" binding:"numeric,gt=0,lt=10000"`
}

func CreateAggregationPipeline(
	query Query,
	beforeLimit []bson.M,
	sort bson.M,
	afterLimit ...[]bson.M,
) []bson.M {
	result := make([]bson.M, len(beforeLimit))
	copy(result, beforeLimit)

	pipeline := make([]bson.M, 0, 4)
	if len(sort) == 1 {
		pipeline = append(pipeline, sort)
	}

	pipeline = append(
		pipeline,
		bson.M{"$skip": (query.Page - 1) * query.Limit},
		bson.M{"$limit": query.Limit},
	)
	if len(afterLimit) == 1 {
		pipeline = append(pipeline, afterLimit[0]...)
	} else if len(afterLimit) > 1 {
		panic("too much arguments")
	}

	totalCountPipeline := []bson.M{
		{"$count": "count"},
	}

	result = append(result,
		bson.M{"$facet": bson.M{
			"data":        pipeline,
			"total_count": totalCountPipeline,
		}},
		bson.M{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	)

	return result
}
