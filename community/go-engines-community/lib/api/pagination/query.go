package pagination

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	DefaultPage     = 1
	DefaultLimit    = 10
	DefaultPaginate = true
)

// Query is base request for pagination.
type Query struct {
	Page     int64 `form:"page" json:"page" binding:"numeric,gt=0"`
	Limit    int64 `form:"limit" json:"limit" binding:"numeric,gt=0"`
	Paginate bool  `form:"paginate" json:"paginate"`
}

type FilteredQuery struct {
	Query
	Search    string `form:"search"`
	Sort      string `form:"sort"`
	WithFlags bool   `form:"with_flags"`
}

func GetDefaultQuery() Query {
	return Query{
		Page:     DefaultPage,
		Limit:    DefaultLimit,
		Paginate: DefaultPaginate,
	}
}

func BindQuery(c *gin.Context, query *Query) error {
	*query = GetDefaultQuery()

	return c.ShouldBindQuery(query)
}

func BindFilteredQuery(c *gin.Context, query *FilteredQuery) error {
	query.Query = GetDefaultQuery()

	return c.ShouldBindQuery(query)
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
	if query.Paginate {
		pipeline = append(
			pipeline,
			bson.M{"$skip": (query.Page - 1) * query.Limit},
			bson.M{"$limit": query.Limit},
		)
	}

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
