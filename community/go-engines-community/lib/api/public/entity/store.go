package entity

import (
	"context"
	"slices"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/db"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/patternfields"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/mongoquery"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/sortparam"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	patterndb "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/db"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const DefaultSortField = "name"

var sortableFields = []string{
	"_id",
	"name",
	"last_event_date",
	"created",
	"updated",
	"imported",
}

var defaultSearchFields = []string{
	"_id",
	"name",
}

var searchableFields = []string{
	"_id",
	"name",
	"connector",
	"component",
	"upstream",
}

var infoValueFieldsForErrParam = []string{
	"infos.*.value",
	"component_infos.*.value",
}

type Store interface {
	Find(ctx context.Context, r PaginatedListRequest) (pagination.ListResponse[Response], error)
	Search(ctx context.Context, r PaginatedSearchRequest) (pagination.ListResponse[Response], error)
}

type store struct {
	dbClient     mongo.DbClient
	dbCollection mongo.DbCollection
	transformer  patternfields.Transformer
}

func NewStore(dbClient mongo.DbClient, transformer patternfields.Transformer) Store {
	return &store{
		dbClient:     dbClient,
		dbCollection: dbClient.Collection(mongo.EntityMongoCollection),
		transformer:  transformer,
	}
}

func (s *store) Find(ctx context.Context, r PaginatedListRequest) (res pagination.ListResponse[Response], err error) {
	sort, err := sortparam.ValidateFunc(r.Sort, IsSortableField, SortErrParam(), DefaultSortField)
	if err != nil {
		return res, err
	}

	beforeLimit, afterLimit, err := ListPipeline(r.ListRequest)
	if err != nil {
		return res, err
	}

	afterLimit = append(afterLimit, sort)
	pipeline := pagination.CreateAggregationPipeline(
		r.Query,
		beforeLimit,
		sort,
		afterLimit,
	)
	ar, err := db.AggregateOne[pagination.Data[Response]](ctx, s.dbCollection, pipeline)
	if err != nil {
		return res, err
	}

	for i := range ar.Data {
		ar.Data[i].setDefaults()
	}

	res = pagination.NewResponse[Response](r.Query, ar)

	return res, nil
}

func (s *store) Search(ctx context.Context, r PaginatedSearchRequest) (res pagination.ListResponse[Response], err error) {
	sort, err := sortparam.ValidateFunc(r.Sort, IsSortableField, SortErrParam(), DefaultSortField)
	if err != nil {
		return res, err
	}

	if len(r.EntityPattern) > 0 {
		patternAliases := patternfields.GetAliases(r.EntityPattern)
		if len(patternAliases) != 0 {
			aliases, err := s.transformer.FetchAliases(ctx, patternAliases)
			if err != nil {
				return res, err
			}

			var valErrs validator.ValidationErrors
			r.EntityPattern, _, valErrs = s.transformer.ApplyAliases(r.EntityPattern, aliases, "EntityPattern")
			if len(valErrs) > 0 {
				return res, validation.NewError(valErrs, r)
			}
		}
	}

	beforeLimit, afterLimit, err := SearchPipeline(r.SearchRequest)
	if err != nil {
		return res, err
	}

	afterLimit = append(afterLimit, sort)
	pipeline := pagination.CreateAggregationPipeline(
		r.Query,
		beforeLimit,
		sort,
		afterLimit,
	)
	ar, err := db.AggregateOne[pagination.Data[Response]](ctx, s.dbCollection, pipeline)
	if err != nil {
		return res, err
	}

	for i := range ar.Data {
		ar.Data[i].setDefaults()
	}

	res = pagination.NewResponse[Response](r.Query, ar)

	return res, nil
}

func ListPipeline(r ListRequest) ([]bson.M, []bson.M, error) {
	filter := bson.M{}
	if r.Enabled != nil {
		filter["enabled"] = *r.Enabled
	}

	if len(r.Types) > 0 {
		filter["type"] = bson.M{"$in": r.Types}
	}

	if len(r.Categories) > 0 {
		filter["category"] = bson.M{"$in": r.Categories}
	}

	if len(r.Connectors) > 0 {
		filter["connector"] = bson.M{"$in": r.Connectors}
	}

	if len(r.Components) > 0 {
		filter["component"] = bson.M{"$in": r.Components}
	}

	if len(r.Upstreams) > 0 {
		filter["upstream"] = bson.M{"$in": r.Upstreams}
	}

	if r.ImportSource != "" {
		filter["import_source"] = r.ImportSource
	}

	beforeLimit := make([]bson.M, 0)
	if len(filter) != 0 {
		beforeLimit = append(beforeLimit, bson.M{"$match": filter})
	}

	if r.Search != "" {
		searchBy := r.SearchBy
		errs := validator.ValidationErrors{}
		errParam := searchErrParam()
		for i, f := range searchBy {
			if !isSearchField(f) {
				errs = append(errs, validation.NewFieldErrorWithParam("oneof", strconv.Itoa(i), "search_by."+strconv.Itoa(i), errParam))
			}
		}

		if len(errs) > 0 {
			return nil, nil, validation.NewError(errs, nil)
		}

		if len(searchBy) == 0 {
			searchBy = defaultSearchFields
		}

		if q := mongoquery.GetSearchQuery(r.Search, searchBy); len(q) > 0 {
			beforeLimit = append(beforeLimit, bson.M{"$match": q})
		}
	}

	afterLimit := nestedObjectsPipeline()
	afterLimit = append(afterLimit, bson.M{
		"$project": bson.M{
			"services":  0,
			"perf_data": 0,
			"comments":  0,
			"imtags":    0,
		},
	})

	return beforeLimit, afterLimit, nil
}

func SearchPipeline(r SearchRequest) ([]bson.M, []bson.M, error) {
	filter := bson.M{}
	if r.Enabled != nil {
		filter["enabled"] = *r.Enabled
	}

	if r.ImportSource != "" {
		filter["import_source"] = r.ImportSource
	}

	beforeLimit := make([]bson.M, 0)
	if len(filter) != 0 {
		beforeLimit = append(beforeLimit, bson.M{"$match": filter})
	}

	if len(r.EntityPattern) > 0 {
		q, err := patterndb.EntityPatternToMongoQuery(r.EntityPattern, "")
		if err != nil {
			return nil, nil, validation.NewSingleError("entity_pattern", "EntityPattern", "EntityPattern", r)
		}

		if len(q) > 0 {
			beforeLimit = append(beforeLimit, bson.M{"$match": q})
		}
	}

	if len(r.PbehaviorPattern) > 0 {
		q, err := patterndb.PbehaviorInfoPatternToMongoQuery(r.PbehaviorPattern, "")
		if err != nil {
			return nil, nil, validation.NewSingleError("pbehavior_pattern", "PbehaviorPattern", "PbehaviorPattern", r)
		}

		if len(q) > 0 {
			beforeLimit = append(beforeLimit, bson.M{"$match": q})
		}
	}

	afterLimit := nestedObjectsPipeline()
	afterLimit = append(afterLimit, bson.M{
		"$project": bson.M{
			"services":  0,
			"perf_data": 0,
			"comments":  0,
			"imtags":    0,
		},
	})

	return beforeLimit, afterLimit, nil
}

func IsSortableField(f string) bool {
	return slices.Contains(sortableFields, f) || isInfoValueField(f)
}

func SortErrParam() string {
	return strings.Join(sortableFields, " ") + " " + strings.Join(infoValueFieldsForErrParam, " ")
}

func isSearchField(f string) bool {
	return slices.Contains(searchableFields, f) || isInfoValueField(f)
}

func searchErrParam() string {
	return strings.Join(searchableFields, " ") + " " + strings.Join(infoValueFieldsForErrParam, " ")
}

func isInfoValueField(f string) bool {
	parts := strings.Split(f, ".")
	if len(parts) != 3 {
		return false
	}

	return (parts[0] == "infos" || parts[0] == "component_infos") && parts[2] == "value"
}

func nestedObjectsPipeline() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.EntityCategoryMongoCollection,
			"localField":   "category",
			"foreignField": "_id",
			"as":           "category",
		}},
		{"$unwind": bson.M{"path": "$category", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"pbehavior_info": bson.M{"$cond": bson.M{
				"if": "$pbehavior_info.id",
				"then": bson.M{
					"timestamp": "$pbehavior_info.timestamp",
					"rule": bson.M{
						"_id":  "$pbehavior_info.id",
						"name": "$pbehavior_info.name",
					},
					"reason": bson.M{
						"_id":  "$pbehavior_info.reason",
						"name": "$pbehavior_info.reason_name",
					},
					"type": bson.M{
						"_id":            "$pbehavior_info.type",
						"name":           "$pbehavior_info.type_name",
						"canonical_type": "$pbehavior_info.canonical_type",
					},
				},
				"else": nil,
			}},
		}},
	}
}
