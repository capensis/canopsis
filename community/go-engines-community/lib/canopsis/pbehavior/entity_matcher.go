package pbehavior

import (
	"context"
	"math"
	"sort"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

const defaultAggregationStep = 100

// EntityMatcher checks if an entity is matched to filter.
type EntityMatcher interface {
	MatchAll(ctx context.Context, entityID string, filters map[string]interface{}) ([]string, error)
}

// NewEntityMatcher creates new matcher.
func NewEntityMatcher(dbClient mongo.DbClient, aggregationStep ...int) EntityMatcher {
	aggregationStepVal := defaultAggregationStep
	if len(aggregationStep) > 0 {
		if len(aggregationStep) > 1 {
			panic("too many arguments")
		}
		aggregationStepVal = aggregationStep[0]
	}

	return &entityMatcher{
		aggregationStep: aggregationStepVal,
		dbCollection:    dbClient.Collection(mongo.EntityMongoCollection),
	}
}

// entityMatcher parses filter and executes parsed mongo query to check if entity is matched.
type entityMatcher struct {
	aggregationStep int
	dbCollection    mongo.DbCollection
}

// MatchAll matches entity by mongo expression. It uses $facet to reduce amount
// of requests to mongo. It makes query for "aggregationStep" filters per mongo request
// to not reach $facet expression limit.
func (m *entityMatcher) MatchAll(
	ctx context.Context,
	entityID string,
	filters map[string]interface{},
) ([]string, error) {
	matched := make(map[string]bool)
	filtersLen := len(filters)
	offsetCount := int(math.Ceil(float64(filtersLen) / float64(m.aggregationStep)))
	filtersArr := make([]keyValue, filtersLen)
	i := 0
	for k, v := range filters {
		filtersArr[i].Key = k
		filtersArr[i].Value = v
		i++
	}
	sort.Slice(filtersArr, func(i, j int) bool {
		return filtersArr[i].Key < filtersArr[j].Key
	})

	for offset := 0; offset < offsetCount; offset++ {
		b := offset * m.aggregationStep
		e := int(math.Min(float64(filtersLen), float64(b+m.aggregationStep)))
		subFilters := filtersArr[b:e]
		pipeline, err := getEntityAggregatePipeline(entityID, subFilters)
		if err != nil {
			return nil, err
		}

		cursor, err := m.dbCollection.Aggregate(ctx, pipeline)
		if err != nil {
			return nil, err
		}

		// Check context done.
		select {
		case <-ctx.Done():
			return nil, nil
		default:
		}

		// Transform mongo doc to result.
		cursor.Next(ctx)
		doc := make(map[string][]string)
		err = cursor.Decode(doc)
		if err != nil {
			return nil, err
		}

		err = cursor.Close(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range subFilters {
			if r, ok := doc[v.Key]; ok {
				matched[v.Key] = len(r) > 0
			} else {
				matched[v.Key] = true
			}
		}
	}

	keys := make([]string, 0)
	for k, ok := range matched {
		if ok {
			keys = append(keys, k)
		}
	}

	return keys, nil
}

// getEntityAggregatePipeline returns doc where key is filter key and value is 1 or 0.
func getEntityAggregatePipeline(
	entityID string,
	filters []keyValue,
) ([]bson.M, error) {
	facetPipeline := bson.M{}

	for _, v := range filters {
		facetPipeline[v.Key] = []bson.M{{"$match": v.Value}}
	}

	return []bson.M{
		{"$match": bson.M{"_id": entityID}},
		{"$facet": facetPipeline},
		{"$addFields": bson.M{
			"ids": bson.M{
				"$arrayToObject": bson.M{
					"$map": bson.M{
						"input": bson.M{"$objectToArray": "$$ROOT"},
						"as":    "each",
						"in": bson.M{
							"k": "$$each.k",
							"v": bson.M{"$map": bson.M{
								"input": "$$each.v",
								"as":    "e",
								"in":    "$$e._id",
							}},
						},
					},
				},
			},
		}},
		{"$replaceRoot": bson.M{"newRoot": "$ids"}},
	}, nil
}

type keyValue struct {
	Key   string
	Value interface{}
}
