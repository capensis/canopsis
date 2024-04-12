package serviceweather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entity/dbquery"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/db"
	pbehaviorlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

const (
	StateIconOk       = "wb_sunny"
	StateIconMinor    = "person"
	StateIconMajor    = "person"
	StateIconCritical = "wb_cloudy"
)

type MongoQueryBuilder struct {
	filterCollection mongo.DbCollection
	authorProvider   author.Provider

	defaultSortBy string

	// entityMatch is match before all lookups
	entityMatch []bson.M
	// additionalMatch is match after some lookups
	additionalMatch []bson.M
	// lookupsForAdditionalMatch is required for match and for result
	lookupsForAdditionalMatch map[string]bool
	// lookupsForSort is required for sort and for result
	lookupsForSort map[string]bool
	// lookups is a slice to define lookups' order since following lookups can depend on previous lookups
	lookups []lookupWithKey
	sort    bson.M

	computedFields bson.M
	// excludedFields is used to remove redundant data from result
	excludedFields []string
}

type lookupWithKey struct {
	key      string
	pipeline []bson.M
}

func NewMongoQueryBuilder(client mongo.DbClient, authorProvider author.Provider) *MongoQueryBuilder {
	return &MongoQueryBuilder{
		filterCollection: client.Collection(mongo.WidgetFiltersMongoCollection),
		authorProvider:   authorProvider,

		defaultSortBy: "name",
	}
}

func (q *MongoQueryBuilder) clear() {
	q.entityMatch = make([]bson.M, 0)
	q.additionalMatch = make([]bson.M, 0)

	q.lookupsForAdditionalMatch = make(map[string]bool)
	q.lookupsForSort = make(map[string]bool)
	q.lookups = make([]lookupWithKey, 0)

	q.sort = bson.M{}
	q.computedFields = bson.M{}
	q.excludedFields = []string{"services"}
}

func (q *MongoQueryBuilder) CreateListAggregationPipeline(ctx context.Context, r ListRequest) ([]bson.M, error) {
	q.clear()

	q.entityMatch = append(q.entityMatch, bson.M{"$match": bson.M{
		"type":    types.EntityTypeService,
		"enabled": true,
	}})
	q.lookups = []lookupWithKey{
		{key: "category", pipeline: dbquery.GetCategoryLookup()},
		{key: "alarm", pipeline: getAlarmLookup()},
		{key: "pbehavior", pipeline: getPbehaviorLookup(q.authorProvider)},
		{key: "pbehavior_info.last_comment", pipeline: dbquery.GetPbehaviorInfoLastCommentLookup(q.authorProvider)},
		{key: "counters", pipeline: getPbehaviorAlarmCountersLookup()},
	}
	err := q.handleWidgetFilter(ctx, r)
	if err != nil {
		return nil, err
	}
	err = q.handlePatterns(r)
	if err != nil {
		return nil, err
	}
	q.handleFilter(r)
	q.handleSort(r.SortBy, r.Sort)

	return q.createPaginationAggregationPipeline(r.Query), nil
}

func (q *MongoQueryBuilder) CreateListDependenciesAggregationPipeline(id string, r EntitiesListRequest, now datetime.CpsTime) ([]bson.M, error) {
	q.clear()

	q.entityMatch = append(q.entityMatch, bson.M{"$match": bson.M{
		"services": id,
		"enabled":  true,
	}})
	q.lookups = []lookupWithKey{
		{key: "category", pipeline: dbquery.GetCategoryLookup()},
		{key: "alarm", pipeline: getAlarmLookup()},
		{key: "pbehavior", pipeline: getPbehaviorLookup(q.authorProvider)},
		{key: "pbehavior_info.last_comment", pipeline: dbquery.GetPbehaviorInfoLastCommentLookup(q.authorProvider)},
		{key: "stats", pipeline: getEventStatsLookup(now)},
		{key: "depends_count", pipeline: dbquery.GetDependsCountPipeline()},
	}
	q.handleSort(r.SortBy, r.Sort)
	q.computedFields = getListDependenciesComputedFields()

	if r.PbhOrigin != "" {
		q.lookups = append(q.lookups, lookupWithKey{key: "pbh_origin_icon", pipeline: getPbhOriginLookup(r.PbhOrigin)})
	}

	return q.createPaginationAggregationPipeline(r.Query), nil
}

func (q *MongoQueryBuilder) createPaginationAggregationPipeline(query pagination.Query) []bson.M {
	beforeLimit, afterLimit := q.createAggregationPipeline()
	if len(afterLimit) > 0 {
		afterLimit = append(afterLimit, q.sort)
	}

	return pagination.CreateAggregationPipeline(
		query,
		beforeLimit,
		q.sort,
		afterLimit,
	)
}

func (q *MongoQueryBuilder) createAggregationPipeline() ([]bson.M, []bson.M) {
	addedLookups := make(map[string]bool)
	beforeLimit := make([]bson.M, len(q.entityMatch))
	copy(beforeLimit, q.entityMatch)

	q.addLookupsToPipeline(q.lookupsForAdditionalMatch, addedLookups, &beforeLimit)
	beforeLimit = append(beforeLimit, q.additionalMatch...)

	q.addLookupsToPipeline(q.lookupsForSort, addedLookups, &beforeLimit)

	afterLimit := make([]bson.M, 0)
	for _, lookup := range q.lookups {
		if !addedLookups[lookup.key] {
			afterLimit = append(afterLimit, lookup.pipeline...)
		}
	}

	addFields := bson.M{}
	for field, v := range q.computedFields {
		addFields[field] = v
	}

	if len(addFields) > 0 {
		afterLimit = append(afterLimit, bson.M{"$addFields": addFields})
	}

	project := bson.M{}
	for _, v := range q.excludedFields {
		project[v] = 0
	}
	afterLimit = append(afterLimit, bson.M{"$project": project})

	return beforeLimit, afterLimit
}

func (q *MongoQueryBuilder) addLookupsToPipeline(lookupsMap, addedLookups map[string]bool, pipeline *[]bson.M) {
	if len(lookupsMap) == 0 {
		return
	}

	for _, lookup := range q.lookups {
		if lookupsMap[lookup.key] && !addedLookups[lookup.key] {
			addedLookups[lookup.key] = true
			*pipeline = append(*pipeline, lookup.pipeline...)
		}
	}
}

func (q *MongoQueryBuilder) handleFilter(r ListRequest) {
	entityMatch := make([]bson.M, 0)
	q.addCategoryFilter(r, &entityMatch)
	if len(entityMatch) > 0 {
		q.entityMatch = append(q.entityMatch, bson.M{"$match": bson.M{"$and": entityMatch}})
	}

	additionalMatch := make([]bson.M, 0)
	q.addHideGreyFilter(r, &additionalMatch)
	if len(additionalMatch) > 0 {
		q.lookupsForAdditionalMatch["alarm"] = true
		q.lookupsForAdditionalMatch["counters"] = true
		q.additionalMatch = append(q.additionalMatch, bson.M{"$match": bson.M{"$and": additionalMatch}})
	}
}

func (q *MongoQueryBuilder) handleWidgetFilter(ctx context.Context, r ListRequest) error {
	for i, id := range r.Filters {
		filter := view.WidgetFilter{}
		err := q.filterCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&filter)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return common.NewValidationError("filters."+strconv.Itoa(i), "Filter doesn't exist.")
			}

			return fmt.Errorf("cannot fetch widget filter: %w", err)
		}

		if len(filter.EntityPattern) == 0 && len(filter.WeatherServicePattern) == 0 ||
			len(filter.AlarmPattern) > 0 ||
			len(filter.PbehaviorPattern) > 0 {
			return common.NewValidationError("filters."+strconv.Itoa(i), "Filter cannot be applied.")
		}

		err = q.handleEntityPattern(filter.EntityPattern)
		if err != nil {
			return fmt.Errorf("invalid entity pattern in widget filter id=%q: %w", filter.ID, err)
		}

		err = q.handleWeatherServicePattern(filter.WeatherServicePattern)
		if err != nil {
			return fmt.Errorf("invalid weather service pattern in widget filter id=%q: %w", filter.ID, err)
		}
	}

	return nil
}

func (q *MongoQueryBuilder) handlePatterns(r ListRequest) error {
	if r.EntityPattern != "" {
		var entityPattern pattern.Entity
		err := json.Unmarshal([]byte(r.EntityPattern), &entityPattern)
		if err != nil {
			return common.NewValidationError("entity_pattern", "EntityPattern is invalid.")
		}
		err = q.handleEntityPattern(entityPattern)
		if err != nil {
			return common.NewValidationError("entity_pattern", "EntityPattern is invalid.")
		}
	}

	if r.WeatherServicePattern != "" {
		var weatherPattern pattern.WeatherServicePattern
		err := json.Unmarshal([]byte(r.WeatherServicePattern), &weatherPattern)
		if err != nil {
			return common.NewValidationError("weather_service_pattern", "WeatherServicePattern is invalid.")
		}
		err = q.handleWeatherServicePattern(weatherPattern)
		if err != nil {
			return common.NewValidationError("weather_service_pattern", "WeatherServicePattern is invalid.")
		}
	}

	return nil
}

func (q *MongoQueryBuilder) handleEntityPattern(entityPattern pattern.Entity) error {
	entityPatternQuery, err := db.EntityPatternToMongoQuery(entityPattern, "")
	if err != nil {
		return err
	}

	if len(entityPatternQuery) > 0 {
		q.entityMatch = append(q.entityMatch, bson.M{"$match": entityPatternQuery})
	}

	return nil
}

func (q *MongoQueryBuilder) handleWeatherServicePattern(weatherServicePattern pattern.WeatherServicePattern) error {
	weatherPatternQuery, err := db.WeatherServicePatternToMongoQuery(weatherServicePattern, "")
	if err != nil {
		return err
	}

	if len(weatherPatternQuery) > 0 {
		q.lookupsForAdditionalMatch["alarm"] = true

		if weatherServicePattern.HasField("is_grey") ||
			weatherServicePattern.HasField("icon") ||
			weatherServicePattern.HasField("secondary_icon") {
			q.lookupsForAdditionalMatch["counters"] = true
		}
		q.additionalMatch = append(q.additionalMatch, bson.M{"$match": weatherPatternQuery})
	}

	return nil
}

func (q *MongoQueryBuilder) addCategoryFilter(r ListRequest, match *[]bson.M) {
	if r.Category != "" {
		*match = append(*match, bson.M{"category": bson.M{"$eq": r.Category}})
	}
}

func (q *MongoQueryBuilder) addHideGreyFilter(r ListRequest, match *[]bson.M) {
	if r.HideGrey {
		*match = append(*match, bson.M{"is_grey": false})
	}
}

func (q *MongoQueryBuilder) handleSort(sortBy, sort string) {
	if sortBy == "" {
		sortBy = q.defaultSortBy
	}
	switch sortBy {
	case "state":
		sortBy = "state.val"
		q.lookupsForSort["alarm"] = true
	case "impact_state":
		q.lookupsForSort["alarm"] = true
	}
	sortDir := 1
	if sort == common.SortDesc {
		sortDir = -1
	}

	sortQuery := bson.D{{Key: sortBy, Value: sortDir}}
	if sortBy != "name" {
		sortQuery = append(sortQuery, bson.E{Key: "name", Value: 1})
	}

	q.sort = bson.M{"$sort": sortQuery}
}

func getAlarmLookup() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from": mongo.AlarmMongoCollection,
			"let":  bson.M{"id": "$_id"},
			"pipeline": []bson.M{
				{"$match": bson.M{"$and": []bson.M{
					{"$expr": bson.M{"$eq": bson.A{"$d", "$$id"}}},
					{"v.resolved": nil},
				}}},
				{"$limit": 1},
			},
			"as": "alarm",
		}},
		{"$unwind": bson.M{"path": "$alarm", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"connector":        "$alarm.v.connector",
			"connector_name":   "$alarm.v.connector_name",
			"component":        "$alarm.v.component",
			"resource":         "$alarm.v.resource",
			"output":           "$alarm.v.output",
			"last_update_date": "$alarm.v.last_update_date",
			"state":            "$alarm.v.state",
			"status":           "$alarm.v.status",
			"snooze":           "$alarm.v.snooze",
			"ack":              "$alarm.v.ack",
			"impact_state":     bson.M{"$multiply": bson.A{"$alarm.v.state.val", "$impact_level"}},
			// For dependencies query
			"alarm_id":      "$alarm._id",
			"creation_date": "$alarm.v.creation_date",
			"display_name":  "$alarm.v.display_name",
			"ticket":        "$alarm.v.ticket",
			"tickets":       "$alarm.v.tickets",
		}},
		{"$project": bson.M{"alarm": 0}},
	}
}

func getPbehaviorLookup(authorProvider author.Provider) []bson.M {
	pipeline := []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorMongoCollection,
			"localField":   "pbehavior_info.id",
			"foreignField": "_id",
			"as":           "pbehavior",
		}},
		{"$unwind": bson.M{"path": "$pbehavior", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"pbehavior.last_comment": bson.M{"$arrayElemAt": bson.A{"$pbehavior.comments", -1}},
		}},
		{"$project": bson.M{
			"pbehavior.comments": 0,
		}},
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorTypeMongoCollection,
			"foreignField": "_id",
			"localField":   "pbehavior.type_",
			"as":           "pbehavior.type",
		}},
		{"$unwind": bson.M{"path": "$pbehavior.type", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorReasonMongoCollection,
			"foreignField": "_id",
			"localField":   "pbehavior.reason",
			"as":           "pbehavior.reason",
		}},
		{"$unwind": bson.M{"path": "$pbehavior.reason", "preserveNullAndEmptyArrays": true}},
	}

	pipeline = append(pipeline, authorProvider.PipelineForField("pbehavior.author")...)
	pipeline = append(pipeline, authorProvider.PipelineForField("pbehavior.last_comment.author")...)
	pipeline = append(pipeline,
		bson.M{"$addFields": bson.M{
			"pbehavior.last_comment": bson.M{
				"$cond": bson.M{
					"if":   "$pbehavior.last_comment._id",
					"then": "$pbehavior.last_comment",
					"else": "$$REMOVE",
				},
			}}},
		bson.M{"$addFields": bson.M{
			"pbehavior": bson.M{
				"$cond": bson.M{
					"if":   "$pbehavior._id",
					"then": "$pbehavior",
					"else": "$$REMOVE",
				},
			}}},
		bson.M{"$addFields": bson.M{
			// todo keep array for backward compatibility
			"pbehaviors": bson.M{"$cond": bson.M{
				"if":   "$pbehavior._id",
				"then": bson.A{"$pbehavior"},
				"else": bson.A{},
			}},
		}},
	)
	return pipeline
}

func getPbehaviorAlarmCountersLookup() []bson.M {
	defaultVal := StateIconOk
	stateVals := []bson.M{
		{
			"case": bson.M{"$eq": bson.A{types.AlarmStateMinor, "$state.val"}},
			"then": StateIconMinor,
		},
		{
			"case": bson.M{"$eq": bson.A{types.AlarmStateMajor, "$state.val"}},
			"then": StateIconMajor,
		},
		{
			"case": bson.M{"$eq": bson.A{types.AlarmStateCritical, "$state.val"}},
			"then": StateIconCritical,
		},
	}

	return []bson.M{
		{
			"$lookup": bson.M{
				"from":         mongo.EntityCountersCollection,
				"localField":   "_id",
				"foreignField": "_id",
				"as":           "counters",
			},
		},
		{
			"$unwind": bson.M{"path": "$counters", "preserveNullAndEmptyArrays": true},
		},
		{"$addFields": bson.M{"pbh_types": bson.M{"$ifNull": bson.A{
			bson.M{"$map": bson.M{
				"input": bson.M{"$objectToArray": "$counters.pbehavior"},
				"in":    "$$this.k",
			}},
			[]int{-1},
		}}}},
		{"$lookup": bson.M{
			"from": mongo.PbehaviorTypeMongoCollection,
			"as":   "pbh_types_counters",
			"let": bson.M{
				"pbh_types":  "$pbh_types",
				"cumulative": bson.M{"$objectToArray": "$counters.pbehavior"},
			},
			"pipeline": []bson.M{
				{"$match": bson.M{"$expr": bson.M{"$in": []string{"$_id", "$$pbh_types"}}}},
				{"$addFields": bson.M{
					"count": bson.M{"$mergeObjects": bson.M{
						"$filter": bson.M{
							"input": "$$cumulative",
							"cond":  bson.M{"$eq": []string{"$$this.k", "$_id"}},
						},
					}},
				}},
				{"$project": bson.M{
					"count": "$count.v",
					"type":  "$$ROOT",
				}},
				{"$match": bson.M{"$expr": bson.M{"$gt": bson.A{"$count", 0}}}},
				{"$project": bson.M{"type.loader_id": 0, "_id": 0, "type.count": 0}},
			},
		}},
		{"$project": bson.M{"pbh_types": 0}},
		{"$unwind": bson.M{"path": "$pbh_types_counters", "preserveNullAndEmptyArrays": true}},
		{"$sort": bson.M{"pbh_types_counters.priority": -1}},
		{"$group": bson.M{
			"_id":                "$_id",
			"pbh_types_counters": bson.M{"$push": "$pbh_types_counters"},
			"data":               bson.M{"$first": "$$ROOT"},
		}},
		{"$replaceRoot": bson.M{
			"newRoot": bson.M{"$mergeObjects": bson.A{
				"$data",
				bson.M{"counters": bson.M{"$mergeObjects": bson.A{
					"$data.counters",
					bson.M{"pbh_types": bson.M{"$filter": bson.M{
						"input": "$pbh_types_counters",
						"cond":  bson.M{"$gt": bson.A{"$$this.count", 0}},
					}}},
				}}},
			}},
		}},
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "_id",
			"foreignField": "services",
			"as":           "depends",
			"pipeline": []bson.M{
				{"$project": bson.M{"_id": 1}},
			},
		}},
		{"$addFields": bson.M{
			"counters.depends": bson.M{"$size": "$depends"},
			"has_open_alarm": bson.M{"$cond": bson.M{
				"if":   bson.M{"$gt": bson.A{"$counters.unacked", 0}},
				"then": true,
				"else": false,
			}},
			"counters.under_pbh": bson.M{"$sum": bson.M{
				"$map": bson.M{
					"input": "$counters.pbh_types",
					"in":    "$$this.count",
				},
			}},
		}},
		{"$project": bson.M{"depends": 0}},
		{"$addFields": bson.M{
			"icon": bson.M{"$switch": bson.M{
				"branches": append(
					[]bson.M{
						{
							"case": bson.M{"$and": []bson.M{
								{"$ifNull": bson.A{"$pbehavior_info", false}},
								{"$ne": bson.A{"$pbehavior_info.canonical_type", pbehaviorlib.TypeActive}},
							}},
							"then": "$pbehavior_info.icon_name",
						},
						{
							"case": bson.M{"$and": []bson.M{
								{"$gt": bson.A{"$counters.under_pbh", 0}},
								{"$eq": bson.A{"$counters.under_pbh", "$counters.depends"}},
							}},
							"then": "",
						},
					},
					stateVals...,
				),
				"default": defaultVal,
			}},
			"is_grey": bson.M{"$switch": bson.M{
				"branches": []bson.M{
					{
						"case": bson.M{"$and": []bson.M{
							{"$ifNull": bson.A{"$pbehavior_info", false}},
							{"$ne": bson.A{"$pbehavior_info.canonical_type", pbehaviorlib.TypeActive}},
						}},
						"then": true,
					},
					{
						"case": bson.M{"$and": []bson.M{
							{"$gt": bson.A{"$counters.under_pbh", 0}},
							{"$eq": bson.A{"$counters.under_pbh", "$counters.depends"}},
						}},
						"then": true,
					},
				},
				"default": false,
			}},
			"secondary_icon": bson.M{"$switch": bson.M{
				"branches": []bson.M{
					{
						"case": bson.M{"$gt": bson.A{"$counters.under_pbh", 0}},
						"then": bson.M{"$arrayElemAt": bson.A{
							bson.M{"$map": bson.M{
								"input": "$counters.pbh_types",
								"in":    "$$this.type.icon_name",
							}},
							0,
						}},
					},
				},
				"default": "",
			}},
		}},
	}
}

func getPbhOriginLookup(origin string) []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorMongoCollection,
			"localField":   "_id",
			"foreignField": "entity",
			"pipeline": []bson.M{
				{"$match": bson.M{"origin": origin}},
				{"$limit": 1},
			},
			"as": "pbh_origin",
		}},
		{"$unwind": bson.M{"path": "$pbh_origin", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorTypeMongoCollection,
			"localField":   "pbh_origin.type_",
			"foreignField": "_id",
			"as":           "pbh_origin.type",
		}},
		{"$unwind": bson.M{"path": "$pbh_origin.type", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"pbh_origin_icon": "$pbh_origin.type.icon_name",
		}},
		{"$project": bson.M{"pbh_origin": 0}},
	}
}

func getEventStatsLookup(now datetime.CpsTime) []bson.M {
	year, month, day := now.Date()
	truncatedInLocation := time.Date(year, month, day, 0, 0, 0, 0, now.Location()).Unix()

	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.EventStatistics,
			"localField":   "_id",
			"foreignField": "_id",
			"as":           "stats",
		}},
		{"$unwind": bson.M{"path": "$stats", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			// stats counters with "last_event" prior "truncatedInLocation" represent as 0
			"stats.ok": bson.M{"$cond": bson.M{
				"if":   bson.M{"$gt": bson.A{"$stats.last_event", truncatedInLocation}},
				"then": "$stats.ok",
				"else": 0,
			}},
			"stats.ko": bson.M{"$cond": bson.M{
				"if":   bson.M{"$gt": bson.A{"$stats.last_event", truncatedInLocation}},
				"then": "$stats.ko",
				"else": 0,
			}},
		}},
	}
}

func getListDependenciesComputedFields() bson.M {
	defaultVal := StateIconOk
	stateVals := []bson.M{
		{
			"case": bson.M{"$eq": bson.A{types.AlarmStateMinor, "$state.val"}},
			"then": StateIconMinor,
		},
		{
			"case": bson.M{"$eq": bson.A{types.AlarmStateMajor, "$state.val"}},
			"then": StateIconMajor,
		},
		{
			"case": bson.M{"$eq": bson.A{types.AlarmStateCritical, "$state.val"}},
			"then": StateIconCritical,
		},
	}

	return bson.M{
		"is_grey": bson.M{"$and": []bson.M{
			{"$ifNull": bson.A{"$pbehavior_info", false}},
			{"$ne": bson.A{"$pbehavior_info.canonical_type", pbehaviorlib.TypeActive}},
		}},
		"icon": bson.M{"$switch": bson.M{
			"branches": append(
				[]bson.M{{
					"case": bson.M{"$and": []bson.M{
						{"$ifNull": bson.A{"$pbehavior_info", false}},
						{"$ne": bson.A{"$pbehavior_info.canonical_type", pbehaviorlib.TypeActive}},
					}},
					"then": "$pbehavior_info.icon_name",
				}},
				stateVals...,
			),
			"default": defaultVal,
		}},
	}
}
