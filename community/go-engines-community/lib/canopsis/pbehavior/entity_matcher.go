package pbehavior

import (
	"context"
	"encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"math"
	"sort"
)

const defaultAggregationStep = 100

// EntityMatcher checks if an entity is matched to filter.
type EntityMatcher interface {
	Match(ctx context.Context, entityID string, filter string) (bool, error)
	MatchAll(ctx context.Context, entityID string, filters map[string]string) (map[string]bool, error)
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

// Match matches entity by mongo expression.
func (m *entityMatcher) Match(
	ctx context.Context,
	entityID string,
	filter string,
) (res bool, resErr error) {
	if filter == "" {
		return true, nil
	}

	pipeline, err := transformFilter(filter)
	if err != nil {
		return false, err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	cursor, err := m.dbCollection.Find(ctx, []bson.M{
		{"$match": bson.M{"_id": entityID}},
		pipeline,
	})
	if err != nil {
		return false, err
	}

	defer func() {
		err := cursor.Close(ctx)
		if resErr == nil && err != nil {
			resErr = err
		}
	}()

	for cursor.Next(ctx) {
		return true, nil
	}

	return false, nil
}

// MatchAll matches entity by mongo expression. It uses $facet to reduce amount
// of requests to mongo. It makes query for "aggregationStep" filters per mongo request
// to not reach $facet expression limit.
func (m *entityMatcher) MatchAll(
	ctx context.Context,
	entityID string,
	filters map[string]string,
) (map[string]bool, error) {
	res := make(map[string]bool)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
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
		pipeline, err := getAggregatePipeline([]string{entityID}, subFilters)
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
			return nil, ctx.Err()
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
				res[v.Key] = len(r) > 0
			} else {
				res[v.Key] = true
			}
		}
	}

	return res, nil
}

// transformFilter unmarshals string to bson expression.
func transformFilter(filter string) (bson.M, error) {
	var bsonFilter bson.M
	err := json.Unmarshal([]byte(filter), &bsonFilter)
	if err != nil {
		return nil, err
	}

	return bson.M{"$match": bsonFilter}, nil
}

// getAggregatePipeline returns doc where key is filter key and value is 1 or 0.
func getAggregatePipeline(
	entityIDs []string,
	filters []keyValue,
) ([]bson.M, error) {
	facetPipeline := bson.M{}

	for _, v := range filters {
		p, err := transformFilter(v.Value)
		if err != nil {
			return nil, err
		}

		facetPipeline[v.Key] = []bson.M{p}
	}

	return []bson.M{
		{"$match": bson.M{"_id": bson.M{"$in": entityIDs}}},
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
	Value string
}
