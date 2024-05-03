package pbehavior

import (
	"context"
	"fmt"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/expression/parser"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoQuery struct {
	typeCollection        mongo.DbCollection
	reasonCollection      mongo.DbCollection
	defaultSearchByFields []string
	defaultSortBy         string

	match             bson.M
	sort              bson.M
	lookupBeforeMatch map[string][]bson.M
	lookupBeforeLimit map[string][]bson.M
	lookupAfterLimit  map[string][]bson.M
	project           []bson.M
}

func CreateMongoQuery(client mongo.DbClient, authorProvider author.Provider) MongoQuery {
	return MongoQuery{
		typeCollection:        client.Collection(mongo.PbehaviorTypeMongoCollection),
		reasonCollection:      client.Collection(mongo.PbehaviorReasonMongoCollection),
		defaultSearchByFields: []string{"_id", "name", "author.name", "comments.author.name", "comments.message", "filter"},
		defaultSortBy:         "created",
		lookupBeforeMatch:     map[string][]bson.M{},
		lookupBeforeLimit:     map[string][]bson.M{},
		lookupAfterLimit: map[string][]bson.M{
			"type":            GetNestedTypePipeline(),
			"reason":          GetNestedReasonPipeline(),
			"exdate":          GetNestedExdatesPipeline(),
			"author":          authorProvider.Pipeline(),
			"comments.author": getCommentAuthorPipeline(authorProvider),
		},
		project: []bson.M{
			{"$addFields": bson.M{
				"comments": bson.M{"$cond": bson.M{
					"if":   "$comments",
					"then": "$comments",
					"else": []string{},
				}},
			}},
		},
	}
}

func (q *MongoQuery) CreateAggregationPipeline(ctx context.Context, r ListRequest) ([]bson.M, error) {
	err := q.handleFilter(ctx, r)
	if err != nil {
		return nil, err
	}
	q.handleSort(r)
	if r.WithFlags {
		q.project = append(q.project, bson.M{"$addFields": bson.M{
			"editable": bson.M{"$cond": bson.M{
				"if":   "$origin",
				"then": false,
				"else": true,
			}},
		}})
	}

	beforeLimit := make([]bson.M, 0)
	for _, m := range q.lookupBeforeMatch {
		beforeLimit = append(beforeLimit, m...)
	}
	if len(q.match) > 0 {
		beforeLimit = append(beforeLimit, bson.M{"$match": q.match})
	}
	for _, m := range q.lookupBeforeLimit {
		beforeLimit = append(beforeLimit, m...)
	}
	afterLimit := make([]bson.M, 0)
	for _, m := range q.lookupAfterLimit {
		afterLimit = append(afterLimit, m...)
	}
	// Sort one more after $group.
	afterLimit = append(afterLimit, q.sort)
	afterLimit = append(afterLimit, q.project...)

	return pagination.CreateAggregationPipeline(
		r.Query,
		beforeLimit,
		q.sort,
		afterLimit,
	), nil
}

func (q *MongoQuery) handleFilter(ctx context.Context, r ListRequest) error {
	filter, err := q.getSearchFilter(ctx, r.Search)
	if err != nil {
		return err
	}

	q.match = filter
	return nil
}

func (q *MongoQuery) handleSort(r ListRequest) {
	sortBy := r.SortBy
	if sortBy == "" {
		sortBy = q.defaultSortBy
	}

	q.adjustLookupsForSort(sortBy)
	q.sort = common.GetSortQuery(sortBy, r.Sort)
}

// getSearchFilter returns mongo query for search filter.
// To search by pbehavior type and reason its ids are found by search filter
// from corresponding mongo collections and query by ids is added to result query.
func (q *MongoQuery) getSearchFilter(ctx context.Context, search string) (bson.M, error) {
	if search == "" {
		return bson.M{}, nil
	}

	p := parser.NewParser()
	expr, err := p.Parse(search, nil)
	if err == nil {
		q.adjustLookupsForFilter(expr.GetFields())

		return expr.MongoQuery(), nil
	}

	searchRegexp := primitive.Regex{
		Pattern: fmt.Sprintf(".*%s.*", search),
		Options: "i",
	}

	res, err := q.typeCollection.Find(ctx, bson.M{"name": searchRegexp},
		options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return nil, err
	}

	types := make([]struct {
		ID string `bson:"_id"`
	}, 0)
	err = res.All(ctx, &types)
	if err != nil {
		return nil, err
	}

	res, err = q.reasonCollection.Find(ctx, bson.M{"name": searchRegexp},
		options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return nil, err
	}

	reasons := make([]struct {
		ID string `bson:"_id"`
	}, 0)
	err = res.All(ctx, &reasons)
	if err != nil {
		return nil, err
	}

	conditions := make([]bson.M, len(q.defaultSearchByFields))
	for i, field := range q.defaultSearchByFields {
		conditions[i] = bson.M{field: searchRegexp}
	}

	if len(types) > 0 {
		ids := make([]string, len(types))
		for i := range ids {
			ids[i] = types[i].ID
		}
		conditions = append(conditions, bson.M{"type_": bson.M{"$in": ids}})
	}

	if len(reasons) > 0 {
		ids := make([]string, len(reasons))
		for i := range ids {
			ids[i] = reasons[i].ID
		}
		conditions = append(conditions, bson.M{"reason": bson.M{"$in": ids}})
	}

	q.adjustLookupsForFilter([]string{"author", "comments.author"})

	return bson.M{"$or": conditions}, nil
}

func (q *MongoQuery) adjustLookupsForFilter(fields []string) {
	for k, v := range q.lookupAfterLimit {
		for _, field := range fields {
			if strings.HasPrefix(field, k) {
				delete(q.lookupAfterLimit, k)
				q.lookupBeforeMatch[k] = v
			}
		}
	}
	for k, v := range q.lookupBeforeLimit {
		for _, field := range fields {
			if strings.HasPrefix(field, k) {
				delete(q.lookupBeforeLimit, k)
				q.lookupBeforeMatch[k] = v
			}
		}
	}
}

func (q *MongoQuery) adjustLookupsForSort(sortBy string) {
	for k, v := range q.lookupAfterLimit {
		if strings.HasPrefix(sortBy, k) {
			delete(q.lookupAfterLimit, k)
			q.lookupBeforeLimit[k] = v
		}
	}
}

func GetNestedObjectsPipeline(authorProvider author.Provider) []bson.M {
	pipeline := append(GetNestedReasonPipeline(), GetNestedTypePipeline()...)
	pipeline = append(pipeline, GetNestedExdatesPipeline()...)
	pipeline = append(pipeline, authorProvider.Pipeline()...)
	pipeline = append(pipeline, getCommentAuthorPipeline(authorProvider)...)

	return pipeline
}

func GetNestedReasonPipeline() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         pbehavior.ReasonCollectionName,
			"localField":   "reason",
			"foreignField": "_id",
			"as":           "reason",
		}},
		{"$unwind": "$reason"},
	}
}

func GetNestedTypePipeline() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorTypeMongoCollection,
			"localField":   "type_",
			"foreignField": "_id",
			"as":           "type",
		}},
		{"$unwind": "$type"},
	}
}

func GetNestedExdatesPipeline() []bson.M {
	return []bson.M{
		// Lookup exdate type
		{"$unwind": bson.M{
			"path":                       "$exdates",
			"preserveNullAndEmptyArrays": true,
			"includeArrayIndex":          "exdate_index",
		}},
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorTypeMongoCollection,
			"localField":   "exdates.type",
			"foreignField": "_id",
			"as":           "exdates.type",
		}},
		{"$unwind": bson.M{"path": "$exdates.type", "preserveNullAndEmptyArrays": true}},
		{"$sort": bson.M{"exdate_index": 1}},
		{"$group": bson.M{
			"_id":     "$_id",
			"data":    bson.M{"$first": "$$ROOT"},
			"exdates": bson.M{"$push": "$exdates"},
		}},
		{"$replaceRoot": bson.M{
			"newRoot": bson.M{"$mergeObjects": bson.A{
				"$data",
				bson.M{"exdates": bson.M{"$filter": bson.M{
					"input": "$exdates",
					"cond":  bson.M{"$ne": bson.A{"$$this", bson.M{}}},
				}}},
			}},
		}},
		// Lookup exceptions
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorExceptionMongoCollection,
			"localField":   "exceptions",
			"foreignField": "_id",
			"as":           "exceptions",
		}},
		{"$unwind": bson.M{
			"path":                       "$exceptions",
			"preserveNullAndEmptyArrays": true,
			"includeArrayIndex":          "exception_index",
		}},
		{"$unwind": bson.M{
			"path":                       "$exceptions.exdates",
			"preserveNullAndEmptyArrays": true,
			"includeArrayIndex":          "exdate_index",
		}},
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorTypeMongoCollection,
			"localField":   "exceptions.exdates.type",
			"foreignField": "_id",
			"as":           "exceptions.exdates.type",
		}},
		{"$unwind": bson.M{"path": "$exceptions.exdates.type", "preserveNullAndEmptyArrays": true}},
		{"$sort": bson.M{"exdate_index": 1}},
		{"$group": bson.M{
			"_id": bson.M{
				"_id":          "$_id",
				"exception_id": "$exceptions._id",
			},
			"data":            bson.M{"$first": "$$ROOT"},
			"exceptions":      bson.M{"$first": "$exceptions"},
			"exception_index": bson.M{"$first": "$exception_index"},
			"exdates":         bson.M{"$push": "$exceptions.exdates"},
		}},
		{"$sort": bson.M{"exception_index": 1}},
		{"$group": bson.M{
			"_id":  "$_id._id",
			"data": bson.M{"$first": "$data"},
			"exceptions": bson.M{"$push": bson.M{"$mergeObjects": bson.A{
				"$exceptions",
				bson.M{
					"_id": bson.M{"$cond": bson.M{
						"if":   "$exceptions._id",
						"then": "$exceptions._id",
						"else": nil,
					}},
					"exdates": bson.M{"$filter": bson.M{
						"input": "$exdates",
						"cond":  bson.M{"$ne": bson.A{"$$this", bson.M{}}},
					}},
				},
			}}},
		}},
		{"$replaceRoot": bson.M{
			"newRoot": bson.M{"$mergeObjects": bson.A{
				"$data",
				bson.M{"exceptions": bson.M{"$filter": bson.M{
					"input": "$exceptions",
					"cond":  bson.M{"$ne": bson.A{"$$this._id", nil}},
				}}},
			}},
		}},
	}
}

func getCommentAuthorPipeline(authorProvider author.Provider) []bson.M {
	pipeline := []bson.M{
		{"$unwind": bson.M{
			"path":                       "$comments",
			"preserveNullAndEmptyArrays": true,
			"includeArrayIndex":          "comment_index",
		}},
	}
	pipeline = append(pipeline, authorProvider.PipelineForField("comments.author")...)
	pipeline = append(pipeline, bson.M{
		"$addFields": bson.M{
			"comments.author": bson.M{"$cond": bson.M{
				"if": "$comments.origin",
				"then": bson.M{
					"name":         "$comments.origin",
					"display_name": "$comments.origin",
				},
				"else": "$comments.author",
			}},
		},
	})
	pipeline = append(pipeline,
		bson.M{"$group": bson.M{
			"_id":  "$_id",
			"data": bson.M{"$first": "$$ROOT"},
			"comments": bson.M{"$push": bson.M{"$cond": bson.M{
				"if":   "$comments._id",
				"then": "$comments",
				"else": "$$REMOVE",
			}}},
		}},
		bson.M{"$replaceRoot": bson.M{"newRoot": bson.M{"$mergeObjects": bson.A{
			"$data",
			bson.M{"comments": "$comments"},
		}}}},
	)
	return pipeline
}
