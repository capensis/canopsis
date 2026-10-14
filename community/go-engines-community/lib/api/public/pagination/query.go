package pagination

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Query is base request for pagination.
type Query struct {
	Limit  int `form:"limit,default:10" json:"limit" binding:"numeric,gt=0,lt=10000"`
	Offset int `form:"offset,default:0" json:"offset" binding:"numeric,gte=0"`
}

func (q *Query) SetDefaults() {
	if q == nil {
		return
	}

	if q.Limit == 0 {
		q.Limit = 10
	}
}

func CreateAggregationPipeline(
	query Query,
	beforeLimit []bson.M,
	sort bson.M,
	afterLimit ...[]bson.M,
) []bson.M {
	pipeline := make([]bson.M, len(beforeLimit), len(beforeLimit)+len(afterLimit)+3)
	copy(pipeline, beforeLimit)
	if len(sort) == 1 {
		pipeline = append(pipeline, sort)
	}

	pipeline = append(
		pipeline,
		bson.M{"$skip": query.Offset},
		bson.M{"$limit": query.Limit + 1},
	)

	if len(afterLimit) == 1 {
		pipeline = append(pipeline, afterLimit[0]...)
	} else if len(afterLimit) > 1 {
		panic("too much arguments")
	}

	return pipeline
}
