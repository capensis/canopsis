package entity

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entity/dbquery"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/db"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statesetting"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type MongoQueryBuilder struct {
	filterCollection mongo.DbCollection

	defaultSearchByFields []string
	defaultSortBy         string

	// entityMatch is match before all lookups
	entityMatch []bson.M
	// additionalMatch is match after some lookups
	additionalMatch []bson.M
	// lookupsForAdditionalMatch is required for match and for result
	lookupsForAdditionalMatch map[string]bool
	// lookupsForSort is required for sort and for result
	lookupsForSort map[string]bool
	// excludeLookupsBeforeSort is used to remove data from pipeline to decrease sorted data.
	// Excluded data is added again in lookups after sort.
	excludeLookupsBeforeSort []string
	// lookups is a slice to define lookups' order since following lookups can depend on previous lookups
	lookups                          []lookupWithKey
	sort                             bson.M
	computedFieldsForAdditionalMatch map[string]bool
	computedFieldsForSort            map[string]bool
	computedFields                   bson.M
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

		defaultSearchByFields: []string{
			"_id", "name", "type",
		},
		defaultSortBy: "_id",
	}
}

func (q *MongoQueryBuilder) clear(now datetime.CpsTime) {
	q.entityMatch = []bson.M{{"$match": bson.M{
		"soft_deleted": bson.M{"$exists": false},
		"healthcheck":  bson.M{"$in": bson.A{nil, false}},
	}}}
	q.additionalMatch = make([]bson.M, 0)

	q.lookupsForAdditionalMatch = make(map[string]bool)
	q.lookupsForSort = make(map[string]bool)
	q.excludeLookupsBeforeSort = make([]string, 0)
	q.lookups = []lookupWithKey{
		{key: "alarm", pipeline: getAlarmLookup()},
		{key: "category", pipeline: dbquery.GetCategoryLookup()},
		{key: "pbehavior_info.icon_name", pipeline: dbquery.GetPbehaviorInfoTypeLookup()},
		{key: "event_stats", pipeline: getEventStatsLookup(now)},
	}

	q.sort = bson.M{}

	q.computedFieldsForAdditionalMatch = make(map[string]bool)
	q.computedFieldsForSort = make(map[string]bool)
	q.computedFields = getComputedFields()
	q.excludedFields = []string{"services", "alarm", "event_stats", "pbehavior_info_type"}
}

func (q *MongoQueryBuilder) CreateListAggregationPipeline(ctx context.Context, r ListRequestWithPagination, now datetime.CpsTime) ([]bson.M, error) {
	q.clear(now)

	err := q.handleWidgetFilter(ctx, r.ListRequest, now)
	if err != nil {
		return nil, err
	}
	err = q.handleEntityPattern(r.ListRequest)
	if err != nil {
		return nil, err
	}
	q.handleFilter(r.ListRequest)
	q.handleSort(r.SortRequest)

	if r.WithFlags {
		q.addFlags()
		q.lookups = append(q.lookups, lookupWithKey{key: "depends_count", pipeline: dbquery.GetDependsCountPipeline()})
		q.lookups = append(q.lookups, lookupWithKey{key: "impacts_count", pipeline: dbquery.GetImpactsCountPipeline()})
	}

	beforeLimit, afterLimit := q.createAggregationPipeline()

	return pagination.CreateAggregationPipeline(
		r.Query,
		beforeLimit,
		q.sort,
		afterLimit,
	), nil
}

func (q *MongoQueryBuilder) CreateTreeOfDepsAggregationPipeline(
	match bson.M,
	paginationQuery pagination.Query,
	sortRequest SortRequest,
	category, search string,
	withFlags bool,
	withStateDependsCount bool,
	now datetime.CpsTime,
) []bson.M {
	q.clear(now)
	and := []bson.M{match}
	if category != "" {
		and = append(and, bson.M{"category": bson.M{"$eq": category}})
	}

	if search != "" {
		and = append(and, common.GetSearchQuery(search, q.defaultSearchByFields))
	}

	q.entityMatch = append(q.entityMatch, bson.M{"$match": bson.M{"$and": and}})
	q.handleSort(sortRequest)

	if withFlags {
		q.lookups = append(q.lookups, lookupWithKey{key: "depends_count", pipeline: dbquery.GetDependsCountPipeline()})
		q.lookups = append(q.lookups, lookupWithKey{key: "impacts_count", pipeline: dbquery.GetImpactsCountPipeline()})
		if withStateDependsCount {
			q.lookups = append(q.lookups, lookupWithKey{key: "state_setting", pipeline: dbquery.GetStateSettingPipeline()})
			q.lookups = append(q.lookups, lookupWithKey{key: "state_depends_count", pipeline: getStateDependsCountPipeline()})
		}
	}

	beforeLimit, afterLimit := q.createAggregationPipeline()

	return pagination.CreateAggregationPipeline(
		paginationQuery,
		beforeLimit,
		q.sort,
		afterLimit,
	)
}

func (q *MongoQueryBuilder) CreateCountAggregationPipeline(ctx context.Context, r ListRequestWithPagination, now datetime.CpsTime) ([]bson.M, error) {
	q.clear(now)

	err := q.handleWidgetFilter(ctx, r.ListRequest, now)
	if err != nil {
		return nil, err
	}
	q.handleFilter(r.ListRequest)
	beforeLimit, _ := q.createAggregationPipeline()

	return beforeLimit, nil
}

func (q *MongoQueryBuilder) CreateOnlyListAggregationPipeline(ctx context.Context, r ListRequest, now datetime.CpsTime) ([]bson.M, error) {
	q.clear(now)

	err := q.handleWidgetFilter(ctx, r, now)
	if err != nil {
		return nil, err
	}
	q.handleFilter(r)
	q.handleSort(r.SortRequest)

	beforeLimit, afterLimit := q.createAggregationPipeline()

	pipeline := append(beforeLimit, q.sort)
	pipeline = append(pipeline, afterLimit...)
	return pipeline, nil
}

func (q *MongoQueryBuilder) createAggregationPipeline() ([]bson.M, []bson.M) {
	addedLookups := make(map[string]bool)
	addedComputedFields := make(map[string]bool)
	beforeLimit := make([]bson.M, len(q.entityMatch))
	copy(beforeLimit, q.entityMatch)

	q.addLookupsToPipeline(q.lookupsForAdditionalMatch, addedLookups, &beforeLimit)
	q.addFieldsToPipeline(q.computedFieldsForAdditionalMatch, addedComputedFields, &beforeLimit)
	beforeLimit = append(beforeLimit, q.additionalMatch...)

	if len(q.excludeLookupsBeforeSort) > 0 {
		project := bson.M{}
		for _, k := range q.excludeLookupsBeforeSort {
			addedLookups[k] = false
			project[k] = 0
		}
		beforeLimit = append(beforeLimit, bson.M{"$project": project})
	}

	q.addLookupsToPipeline(q.lookupsForSort, addedLookups, &beforeLimit)
	q.addFieldsToPipeline(q.computedFieldsForSort, addedComputedFields, &beforeLimit)

	afterLimit := make([]bson.M, 0)

	for _, lookup := range q.lookups {
		if !addedLookups[lookup.key] {
			afterLimit = append(afterLimit, lookup.pipeline...)
		}
	}

	addFields := bson.M{}
	for field, v := range q.computedFields {
		if !addedComputedFields[field] {
			addFields[field] = v
		}
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

func (q *MongoQueryBuilder) addFieldsToPipeline(fieldsMap, addedFields map[string]bool, pipeline *[]bson.M) {
	if len(fieldsMap) == 0 {
		return
	}

	query := bson.M{}
	for field, v := range q.computedFields {
		if fieldsMap[field] && !addedFields[field] {
			addedFields[field] = true
			query[field] = v
		}
	}

	*pipeline = append(*pipeline, bson.M{"$addFields": query})
}

func (q *MongoQueryBuilder) handleFilter(r ListRequest) {
	entityMatch := make([]bson.M, 0)
	q.addSearchFilter(r, &entityMatch)
	q.addCategoryFilter(r, &entityMatch)
	q.addTypeFilter(r, &entityMatch)
	q.addNoEventsFilter(r, &entityMatch)

	if len(entityMatch) > 0 {
		q.entityMatch = append(q.entityMatch, bson.M{"$match": bson.M{"$and": entityMatch}})
	}
}

func (q *MongoQueryBuilder) handleWidgetFilter(ctx context.Context, r ListRequest, now datetime.CpsTime) error {
	for i, id := range r.Filters {
		filter := view.WidgetFilter{}
		err := q.filterCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&filter)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return common.NewValidationError("filters."+strconv.Itoa(i), "Filter doesn't exist.")
			}

			return fmt.Errorf("cannot fetch widget filter: %w", err)
		}

		if len(filter.EntityPattern) == 0 && len(filter.PbehaviorPattern) == 0 && len(filter.AlarmPattern) == 0 ||
			len(filter.WeatherServicePattern) > 0 {
			return common.NewValidationError("filters."+strconv.Itoa(i), "Filter cannot be applied.")
		}

		entityPatternQuery, err := db.EntityPatternToMongoQuery(filter.EntityPattern, "")
		if err != nil {
			return fmt.Errorf("invalid entity pattern in widget filter id=%q: %w", filter.ID, err)
		}

		if len(entityPatternQuery) > 0 {
			q.entityMatch = append(q.entityMatch, bson.M{"$match": entityPatternQuery})
		}

		pbhPatternQuery, err := db.PbehaviorInfoPatternToMongoQuery(filter.PbehaviorPattern, "")
		if err != nil {
			return fmt.Errorf("invalid pbehavior pattern in widget filter id=%q: %w", filter.ID, err)
		}

		if len(pbhPatternQuery) > 0 {
			q.entityMatch = append(q.entityMatch, bson.M{"$match": pbhPatternQuery})
		}

		alarmPatternQuery, err := db.AlarmPatternToMongoQuery(filter.AlarmPattern, "alarm")
		if err != nil {
			return fmt.Errorf("invalid alarm pattern in widget filter id=%q: %w", filter.ID, err)
		}

		if len(alarmPatternQuery) > 0 {
			q.lookupsForAdditionalMatch["alarm"] = true
			q.additionalMatch = append(q.additionalMatch, bson.M{"$match": alarmPatternQuery})

			if filter.AlarmPattern.HasInfosField() {
				q.computedFieldsForAdditionalMatch["alarm.v.infos_array"] = true
				q.computedFields["alarm.v.infos_array"] = bson.M{"$objectToArray": "$alarm.v.infos"}
			}

			if filter.AlarmPattern.HasField("v.duration") {
				q.computedFieldsForAdditionalMatch["alarm.v.duration"] = true
				q.computedFields["alarm.v.duration"] = getDurationField(now)
			}
		}
	}

	return nil
}

func (q *MongoQueryBuilder) handleEntityPattern(r ListRequest) error {
	if r.EntityPattern == "" {
		return nil
	}

	var entityPattern pattern.Entity
	err := json.Unmarshal([]byte(r.EntityPattern), &entityPattern)
	if err != nil {
		return common.NewValidationError("entity_pattern", "EntityPattern is invalid.")
	}

	entityPatternQuery, err := db.EntityPatternToMongoQuery(entityPattern, "")
	if err != nil {
		return common.NewValidationError("entity_pattern", "EntityPattern is invalid.")
	}

	if len(entityPatternQuery) > 0 {
		q.entityMatch = append(q.entityMatch, bson.M{"$match": entityPatternQuery})
	}

	return nil
}

func (q *MongoQueryBuilder) addSearchFilter(r ListRequest, match *[]bson.M) {
	if r.Search == "" {
		return
	}

	searchBy := r.SearchBy
	if len(searchBy) == 0 {
		searchBy = q.defaultSearchByFields
	}

	*match = append(*match, common.GetSearchQuery(r.Search, searchBy))
}

func (q *MongoQueryBuilder) addCategoryFilter(r ListRequest, match *[]bson.M) {
	if r.Category == "" {
		return
	}

	*match = append(*match, bson.M{"category": bson.M{"$eq": r.Category}})
}

func (q *MongoQueryBuilder) addTypeFilter(r ListRequest, match *[]bson.M) {
	if len(r.Type) == 0 {
		return
	}

	*match = append(*match, bson.M{"type": bson.M{"$in": r.Type}})
}

func (q *MongoQueryBuilder) addNoEventsFilter(r ListRequest, match *[]bson.M) {
	if !r.NoEvents {
		return
	}

	*match = append(*match, bson.M{"idle_since": bson.M{"$gt": 0}})
}

func (q *MongoQueryBuilder) handleSort(r SortRequest) {
	sortBy := r.SortBy
	if sortBy == "" {
		sortBy = q.defaultSortBy
	}
	sort := r.Sort
	if sort == "" {
		sort = common.SortAsc
	}

	q.adjustLookupsForSort([]string{sortBy})
	q.sort = common.GetSortQuery(sortBy, sort)
}

func (q *MongoQueryBuilder) adjustLookupsForSort(sortFields []string) {
	lookupByComputedField := map[string]string{
		"state":        "alarm",
		"impact_state": "alarm",
	}

	for field := range q.computedFields {
		for _, sortField := range sortFields {
			if sortField == field {
				q.computedFieldsForSort[field] = true
				if lookup := lookupByComputedField[sortField]; lookup != "" {
					q.lookupsForSort[lookup] = true
				}
				break
			}
		}
	}

	for lookup := range q.lookupsForAdditionalMatch {
		found := false
		for _, sortField := range sortFields {
			if strings.HasPrefix(sortField, lookup) || lookup == lookupByComputedField[sortField] {
				found = true
				break
			}
		}

		if !found {
			q.excludeLookupsBeforeSort = append(q.excludeLookupsBeforeSort, lookup)
		}
	}

	for _, lookup := range q.lookups {
		for _, sortField := range sortFields {
			if strings.HasPrefix(sortField, lookup.key) {
				q.lookupsForSort[lookup.key] = true
				break
			}
		}
	}
}

func (q *MongoQueryBuilder) addFlags() {
	q.lookups = append(q.lookups, lookupWithKey{key: "deletable", pipeline: getDeletablePipeline()})
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
	}
}

func getEventStatsLookup(now datetime.CpsTime) []bson.M {
	year, month, day := now.Date()
	truncatedInLocation := datetime.CpsTime{Time: time.Date(year, month, day, 0, 0, 0, 0, now.Location())}

	return []bson.M{
		{"$lookup": bson.M{
			"from": mongo.EventStatistics,
			"let":  bson.M{"id": "$_id"},
			"pipeline": []bson.M{
				{"$match": bson.M{"$and": []bson.M{
					{"$expr": bson.M{"$eq": bson.A{"$_id", "$$id"}}},
					{"last_event": bson.M{"$gt": truncatedInLocation}},
				}}},
			},
			"as": "event_stats",
		}},
		{"$unwind": bson.M{"path": "$event_stats", "preserveNullAndEmptyArrays": true}},
	}
}

func getDeletablePipeline() []bson.M {
	return []bson.M{
		// Entity can be deleted if entity is service or if there aren't any alarm which is related to entity.
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
			"as": "alarms",
		}},
		{"$addFields": bson.M{
			"deletable": bson.M{"$cond": bson.M{
				"if": bson.M{"$or": []bson.M{
					{"$eq": bson.A{"$type", types.EntityTypeService}},
					{"$eq": bson.A{"$alarms", bson.A{}}},
				}},
				"then": true,
				"else": false,
			}},
		}},
		{"$project": bson.M{"alarms": 0}},
	}
}

func getComputedFields() bson.M {
	return bson.M{
		"ok_events": bson.M{"$ifNull": bson.A{
			"$event_stats.ok",
			0,
		}},
		"ko_events": bson.M{"$ifNull": bson.A{
			"$event_stats.ko",
			0,
		}},
		"state": bson.M{"$ifNull": bson.A{
			"$alarm.v.state.val",
			0,
		}},
		"impact_state": bson.M{"$cond": bson.M{
			"if":   "$alarm.v.state.val",
			"then": bson.M{"$multiply": bson.A{"$alarm.v.state.val", "$impact_level"}},
			"else": 0,
		}},
		"status": bson.M{"$ifNull": bson.A{
			"$alarm.v.status.val",
			0,
		}},
		"ack":                    "$alarm.v.ack",
		"snooze":                 "$alarm.v.snooze",
		"alarm_last_update_date": "$alarm.v.last_update_date",
	}
}

func getDurationField(now datetime.CpsTime) bson.M {
	return bson.M{"$ifNull": bson.A{
		"$alarm.v.duration",
		bson.M{"$subtract": bson.A{
			bson.M{"$cond": bson.M{
				"if":   "$alarm.v.resolved",
				"then": "$alarm.v.resolved",
				"else": now,
			}},
			"$alarm.v.creation_date",
		}},
	}}
}

func getStateDependsCountPipeline() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.EntityCountersCollection,
			"localField":   "_id",
			"foreignField": "_id",
			"as":           "counters",
		}},
		{"$unwind": bson.M{"path": "$counters", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"state_depends_count": bson.M{"$switch": bson.M{
				"branches": []bson.M{
					{
						"case": bson.M{"$and": []bson.M{
							{"$eq": bson.A{"$state_setting.method", statesetting.MethodInherited}},
							{"$eq": bson.A{"$type", types.EntityTypeService}},
						}},
						"then": bson.M{"$sum": bson.A{
							"$counters.inherited_state.ok",
							"$counters.inherited_state.minor",
							"$counters.inherited_state.major",
							"$counters.inherited_state.critical",
						}},
					},
					{
						"case": bson.M{"$and": []bson.M{
							{"$eq": bson.A{"$state_setting.method", statesetting.MethodInherited}},
							{"$eq": bson.A{"$type", types.EntityTypeComponent}},
						}},
						"then": bson.M{"$sum": bson.A{
							"$counters.state.ok",
							"$counters.state.minor",
							"$counters.state.major",
							"$counters.state.critical",
						}},
					},
				},
				"default": "$depends_count",
			}},
		}},
		{"$project": bson.M{"counters": 0}},
	}
}
