package serviceweather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	pbehaviorlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type MongoQueryBuilder struct {
	filterCollection mongo.DbCollection

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

func NewMongoQueryBuilder(client mongo.DbClient) *MongoQueryBuilder {
	return &MongoQueryBuilder{
		filterCollection: client.Collection(mongo.WidgetFiltersMongoCollection),

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
	q.excludedFields = []string{"depends", "impact"}
}

func (q *MongoQueryBuilder) CreateListAggregationPipeline(ctx context.Context, r ListRequest) ([]bson.M, error) {
	q.clear()

	q.entityMatch = append(q.entityMatch, bson.M{"$match": bson.M{
		"type":         types.EntityTypeService,
		"enabled":      true,
		"soft_deleted": bson.M{"$in": bson.A{false, nil}},
	}})
	q.lookups = []lookupWithKey{
		{key: "category", pipeline: getCategoryLookup()},
		{key: "alarm", pipeline: getAlarmLookup()},
		{key: "pbehavior", pipeline: getPbehaviorLookup()},
		{key: "alarm_counters", pipeline: getPbehaviorAlarmCountersLookup()},
		{key: "pbehavior_info.icon_name", pipeline: getPbehaviorInfoTypeLookup()},
	}
	err := q.handleWidgetFilter(ctx, r)
	if err != nil {
		return nil, err
	}
	q.handleFilter(r)
	q.handleSort(r.SortBy, r.Sort)

	return q.createPaginationAggregationPipeline(r.Query), nil
}

func (q *MongoQueryBuilder) CreateListDependenciesAggregationPipeline(ids []string, r EntitiesListRequest, now types.CpsTime) ([]bson.M, error) {
	q.clear()

	q.entityMatch = append(q.entityMatch, bson.M{"$match": bson.M{
		"_id":     bson.M{"$in": ids},
		"enabled": true,
	}})
	q.lookups = []lookupWithKey{
		{key: "category", pipeline: getCategoryLookup()},
		{key: "alarm", pipeline: getAlarmLookup()},
		{key: "pbehavior", pipeline: getPbehaviorLookup()},
		{key: "pbehavior_info.icon_name", pipeline: getPbehaviorInfoTypeLookup()},
		{key: "stats", pipeline: getEventStatsLookup(now)},
	}
	q.handleSort(r.SortBy, r.Sort)
	q.computedFields = getListDependenciesComputedFields()

	return q.createPaginationAggregationPipeline(r.Query), nil
}

func (q *MongoQueryBuilder) createPaginationAggregationPipeline(query pagination.Query) []bson.M {
	beforeLimit, afterLimit := q.createAggregationPipeline()

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
}

func (q *MongoQueryBuilder) handleWidgetFilter(ctx context.Context, r ListRequest) error {
	if len(r.Filters) == 0 {
		return nil
	}

	for _, v := range r.Filters {
		filter := view.WidgetFilter{}
		err := q.filterCollection.FindOne(ctx, bson.M{"_id": v}).Decode(&filter)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return common.NewValidationError("filter", errors.New("Filter doesn't exist."))
			}
			return fmt.Errorf("cannot fetch widget filter: %w", err)
		}

		if len(filter.EntityPattern) == 0 && len(filter.WeatherServicePattern) == 0 && len(filter.OldMongoQuery) == 0 ||
			len(filter.AlarmPattern) > 0 ||
			len(filter.PbehaviorPattern) > 0 {
			return common.NewValidationError("filter", errors.New("Filter cannot be applied."))
		}

		entityPatternQuery, err := filter.EntityPattern.ToMongoQuery("")
		if err != nil {
			return fmt.Errorf("invalid entity pattern in widget filter id=%q: %w", filter.ID, err)
		}

		if len(entityPatternQuery) > 0 {
			q.entityMatch = append(q.entityMatch, bson.M{"$match": entityPatternQuery})
		}

		weatherPatternQuery, err := filter.WeatherServicePattern.ToMongoQuery("")
		if err != nil {
			return fmt.Errorf("invalid weather service pattern in widget filter id=%q: %w", filter.ID, err)
		}
		if len(weatherPatternQuery) > 0 {
			q.lookupsForAdditionalMatch["alarm"] = true

			if filter.WeatherServicePattern.HasField("is_grey") ||
				filter.WeatherServicePattern.HasField("icon") ||
				filter.WeatherServicePattern.HasField("secondary_icon") {
				q.lookupsForAdditionalMatch["alarm_counters"] = true
			}
			q.additionalMatch = append(q.additionalMatch, bson.M{"$match": weatherPatternQuery})
		}

		if len(entityPatternQuery) == 0 && len(weatherPatternQuery) == 0 &&
			len(filter.OldMongoQuery) > 0 {
			var query map[string]interface{}
			err := json.Unmarshal([]byte(filter.OldMongoQuery), &query)
			if err != nil {
				return fmt.Errorf("cannot unmarshal old mongo query: %w", err)
			}

			for _, lookup := range q.lookups {
				q.lookupsForAdditionalMatch[lookup.key] = true
			}

			q.additionalMatch = append(q.additionalMatch, bson.M{"$match": query})
		}
	}

	return nil
}

func (q *MongoQueryBuilder) addCategoryFilter(r ListRequest, match *[]bson.M) {
	if r.Category != "" {
		*match = append(*match, bson.M{"category": bson.M{"$eq": r.Category}})
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

func getCategoryLookup() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.EntityCategoryMongoCollection,
			"localField":   "category",
			"foreignField": "_id",
			"as":           "category",
		}},
		{"$unwind": bson.M{"path": "$category", "preserveNullAndEmptyArrays": true}},
	}
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
		}},
		{"$project": bson.M{"alarm": 0}},
	}
}

func getPbehaviorLookup() []bson.M {
	return []bson.M{
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
			"from":         mongo.RightsMongoCollection,
			"foreignField": "_id",
			"localField":   "pbehavior.author",
			"as":           "pbehavior.author",
		}},
		{"$unwind": bson.M{"path": "$pbehavior.author", "preserveNullAndEmptyArrays": true}},
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
		{"$addFields": bson.M{
			// todo keep array for backward compatibility
			"pbehaviors": bson.M{"$cond": bson.M{
				"if":   "$pbehavior._id",
				"then": bson.A{"$pbehavior"},
				"else": bson.A{},
			}},
		}},
	}
}

func getPbehaviorInfoTypeLookup() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorTypeMongoCollection,
			"foreignField": "_id",
			"localField":   "pbehavior_info.type",
			"as":           "pbehavior_info_type",
		}},
		{"$unwind": bson.M{"path": "$pbehavior_info_type", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"pbehavior_info": bson.M{"$cond": bson.M{
				"if": "$pbehavior_info",
				"then": bson.M{"$mergeObjects": bson.A{
					"$pbehavior_info",
					bson.M{"icon_name": "$pbehavior_info_type.icon_name"},
				}},
				"else": nil,
			}},
		}},
		{"$project": bson.M{"pbehavior_info_type": 0}},
	}
}

func getPbehaviorAlarmCountersLookup() []bson.M {
	defaultVal := types.AlarmStateTitleOK
	stateVals := []bson.M{
		{
			"case": bson.M{"$eq": bson.A{types.AlarmStateMinor, "$state.val"}},
			"then": types.AlarmStateTitleMinor,
		},
		{
			"case": bson.M{"$eq": bson.A{types.AlarmStateMajor, "$state.val"}},
			"then": types.AlarmStateTitleMajor,
		},
		{
			"case": bson.M{"$eq": bson.A{types.AlarmStateCritical, "$state.val"}},
			"then": types.AlarmStateTitleCritical,
		},
	}

	return []bson.M{
		{"$addFields": bson.M{"pbh_types": bson.M{"$ifNull": bson.A{
			bson.M{"$map": bson.M{
				"input": bson.M{"$objectToArray": "$alarms_cumulative_data.watched_pbehavior_count"},
				"in":    "$$this.k",
			}},
			[]int{-1},
		}}}},
		{"$lookup": bson.M{
			"from": mongo.PbehaviorTypeMongoCollection,
			"as":   "alarm_counters",
			"let": bson.M{
				"pbh_types":  "$pbh_types",
				"cumulative": bson.M{"$objectToArray": "$alarms_cumulative_data.watched_pbehavior_count"},
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
		{"$project": bson.M{
			"_id":            "$_id",
			"alarm_counters": "$alarm_counters",
			"data":           "$$ROOT",
		}},
		{"$replaceRoot": bson.M{
			"newRoot": bson.M{"$mergeObjects": bson.A{
				"$data",
				bson.M{"alarm_counters": bson.M{"$filter": bson.M{
					"input": "$alarm_counters",
					"cond":  bson.M{"$gt": bson.A{"$$this.count", 0}},
				}}},
			}},
		}},
		{"$addFields": bson.M{
			"depends_count": bson.M{"$size": "$depends"},
			"has_open_alarm": bson.M{"$cond": bson.M{
				"if":   bson.M{"$gt": bson.A{"$alarms_cumulative_data.watched_not_acked_count", 0}},
				"then": true,
				"else": false,
			}},
			"watched_inactive_count": bson.M{"$sum": bson.M{
				"$map": bson.M{
					"input": bson.M{
						"$filter": bson.M{
							"input": "$alarm_counters",
							"cond":  bson.M{"$in": bson.A{"$$this.type.type", bson.A{pbehaviorlib.TypeMaintenance, pbehaviorlib.TypePause, pbehaviorlib.TypeInactive}}},
						},
					},
					"in": "$$this.count",
				},
			}},
			"watched_pbehavior_type": bson.M{"$switch": bson.M{
				"branches": []bson.M{
					{
						"case": bson.M{"$ne": bson.A{bson.A{}, bson.M{
							"$filter": bson.M{
								"input": "$alarm_counters",
								"cond":  bson.M{"$eq": bson.A{"$$this.type.type", pbehaviorlib.TypeMaintenance}},
							},
						}}},
						"then": pbehaviorlib.TypeMaintenance,
					},
					{
						"case": bson.M{"$ne": bson.A{bson.A{}, bson.M{
							"$filter": bson.M{
								"input": "$alarm_counters",
								"cond":  bson.M{"$eq": bson.A{"$$this.type.type", pbehaviorlib.TypePause}},
							},
						}}},
						"then": pbehaviorlib.TypePause,
					},
					{
						"case": bson.M{"$ne": bson.A{bson.A{}, bson.M{
							"$filter": bson.M{
								"input": "$alarm_counters",
								"cond":  bson.M{"$eq": bson.A{"$$this.type.type", pbehaviorlib.TypeInactive}},
							},
						}}},
						"then": pbehaviorlib.TypeInactive,
					},
				},
				"default": "",
			}},
		}},
		{"$addFields": bson.M{
			"icon": bson.M{"$switch": bson.M{
				"branches": append(
					[]bson.M{
						// If service is not active return pbehavior type icon.
						{
							"case": bson.M{"$and": []bson.M{
								{"$ifNull": bson.A{"$pbehavior_info", false}},
								{"$ne": bson.A{"$pbehavior_info.canonical_type", pbehaviorlib.TypeActive}},
							}},
							"then": "$pbehavior_info.canonical_type",
						},
						// If all watched alarms are not active return most priority pbehavior type of watched alarms.
						{
							"case": bson.M{"$and": []bson.M{
								{"$gt": bson.A{"$watched_inactive_count", 0}},
								{"$eq": bson.A{"$watched_inactive_count", "$depends_count"}},
							}},
							"then": "$watched_pbehavior_type",
						},
					},
					// Else return state icon.
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
							{"$gt": bson.A{"$watched_inactive_count", 0}},
							{"$eq": bson.A{"$watched_inactive_count", "$depends_count"}},
						}},
						"then": true,
					},
				},
				"default": false,
			}},
			"secondary_icon": bson.M{"$switch": bson.M{
				"branches": []bson.M{
					{
						// If only some watched alarms are not active return most priority pbehavior type of watched alarms.
						"case": bson.M{"$and": []bson.M{
							{"$gt": bson.A{"$watched_inactive_count", 0}},
							{"$lt": bson.A{"$watched_inactive_count", "$depends_count"}},
						}},
						"then": "$watched_pbehavior_type",
					},
				},
				"default": "",
			}},
		}},
		{"$project": bson.M{
			"watched_inactive_count": 0,
			"watched_pbehavior_type": 0,
			"alarms_cumulative_data": 0,
		}},
	}
}

func getEventStatsLookup(now types.CpsTime) []bson.M {
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
	defaultVal := types.AlarmStateTitleOK
	stateVals := []bson.M{
		{
			"case": bson.M{"$eq": bson.A{types.AlarmStateMinor, "$state.val"}},
			"then": types.AlarmStateTitleMinor,
		},
		{
			"case": bson.M{"$eq": bson.A{types.AlarmStateMajor, "$state.val"}},
			"then": types.AlarmStateTitleMajor,
		},
		{
			"case": bson.M{"$eq": bson.A{types.AlarmStateCritical, "$state.val"}},
			"then": types.AlarmStateTitleCritical,
		},
	}

	return bson.M{
		"is_grey": bson.M{"$and": []bson.M{
			{"$ifNull": bson.A{"$pbehavior_info", false}},
			{"$ne": bson.A{"$pbehavior_info.canonical_type", pbehaviorlib.TypeActive}},
		}},
		"icon": bson.M{"$switch": bson.M{
			"branches": append(
				// If service is not active return pbehavior type icon.
				[]bson.M{{
					"case": bson.M{"$and": []bson.M{
						{"$ifNull": bson.A{"$pbehavior_info", false}},
						{"$ne": bson.A{"$pbehavior_info.canonical_type", pbehaviorlib.TypeActive}},
					}},
					"then": "$pbehavior_info.canonical_type",
				}},
				// Else return state icon.
				stateVals...,
			),
			"default": defaultVal,
		}},
		"depends_count": bson.M{"$cond": bson.M{
			"if":   bson.M{"$eq": bson.A{"$type", types.EntityTypeService}},
			"then": bson.M{"$size": "$depends"},
			"else": 0,
		}},
	}
}
