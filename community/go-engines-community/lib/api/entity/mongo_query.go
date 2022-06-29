package entity

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
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

func (q *MongoQueryBuilder) clear(now types.CpsTime) {
	q.entityMatch = make([]bson.M, 0)
	q.additionalMatch = make([]bson.M, 0)

	q.lookupsForAdditionalMatch = make(map[string]bool)
	q.lookupsForSort = make(map[string]bool)
	q.excludeLookupsBeforeSort = make([]string, 0)
	q.lookups = []lookupWithKey{
		{key: "alarm", pipeline: getAlarmLookup()},
		{key: "category", pipeline: getCategoryLookup()},
		{key: "pbehavior_info.icon_name", pipeline: getPbehaviorInfoTypeLookup()},
		{key: "event_stats", pipeline: getEventStatsLookup(now)},
	}

	q.sort = bson.M{}

	q.computedFieldsForAdditionalMatch = make(map[string]bool)
	q.computedFieldsForSort = make(map[string]bool)
	q.computedFields = getComputedFields()
	q.excludedFields = []string{"alarm", "event_stats", "pbehavior_info_type"}
}

func (q *MongoQueryBuilder) CreateListAggregationPipeline(ctx context.Context, r ListRequestWithPagination, now types.CpsTime) ([]bson.M, error) {
	q.clear(now)

	err := q.handleWidgetFilter(ctx, r.ListRequest, now)
	if err != nil {
		return nil, err
	}
	err = q.handleFilter(r.ListRequest)
	if err != nil {
		return nil, err
	}
	err = q.handleSort(r.ListRequest)
	if err != nil {
		return nil, err
	}

	if r.WithFlags {
		q.lookups = append(q.lookups, lookupWithKey{key: "deletable", pipeline: getDeletablePipeline()})
	}

	beforeLimit, afterLimit := q.createAggregationPipeline()

	return pagination.CreateAggregationPipeline(
		r.Query,
		beforeLimit,
		q.sort,
		afterLimit,
	), nil
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

func (q *MongoQueryBuilder) handleFilter(r ListRequest) error {
	entityMatch := make([]bson.M, 0)
	q.addSearchFilter(r, &entityMatch)
	q.addCategoryFilter(r, &entityMatch)
	q.addTypeFilter(r, &entityMatch)
	q.addNoEventsFilter(r, &entityMatch)

	if len(entityMatch) > 0 {
		q.entityMatch = append(q.entityMatch, bson.M{"$match": bson.M{"$and": entityMatch}})
	}

	return nil
}

func (q *MongoQueryBuilder) handleWidgetFilter(ctx context.Context, r ListRequest, now types.CpsTime) error {
	if r.Filter == "" {
		return nil
	}

	filter := view.WidgetFilter{}
	err := q.filterCollection.FindOne(ctx, bson.M{"_id": r.Filter}).Decode(&filter)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return common.NewValidationError("filter", errors.New("Filter doesn't exist."))
		}
		return fmt.Errorf("cannot fetch widget filter: %w", err)
	}

	entityPatternQuery, err := filter.EntityPattern.ToMongoQuery("")
	if err != nil {
		return fmt.Errorf("invalid entity pattern in widget filter id=%q: %w", filter.ID, err)
	}

	if len(entityPatternQuery) > 0 {
		q.entityMatch = append(q.entityMatch, bson.M{"$match": entityPatternQuery})
	}

	pbhPatternQuery, err := filter.PbehaviorPattern.ToMongoQuery("pbehavior_info")
	if err != nil {
		return fmt.Errorf("invalid pbehavior pattern in widget filter id=%q: %w", filter.ID, err)
	}

	if len(pbhPatternQuery) > 0 {
		q.entityMatch = append(q.entityMatch, bson.M{"$match": pbhPatternQuery})
	}

	alarmPatternQuery, err := filter.AlarmPattern.ToMongoQuery("alarm")
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

	if len(entityPatternQuery) == 0 && len(pbhPatternQuery) == 0 && len(alarmPatternQuery) == 0 &&
		len(filter.OldMongoQuery) > 0 {
		var query map[string]interface{}
		err := json.Unmarshal([]byte(filter.OldMongoQuery), &query)
		if err != nil {
			return fmt.Errorf("cannot unmarshal old mongo query: %w", err)
		}

		q.entityMatch = append(q.entityMatch, bson.M{"$match": query})

		return nil
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

func (q *MongoQueryBuilder) handleSort(r ListRequest) error {
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

	return nil
}

func (q *MongoQueryBuilder) adjustLookupsForSort(sortFields []string) {
	for field := range q.computedFields {
		for _, sortField := range sortFields {
			if sortField == field {
				q.computedFieldsForSort[field] = true
				break
			}
		}
	}

	for lookup := range q.lookupsForAdditionalMatch {
		found := false
		for _, sortField := range sortFields {
			if strings.HasPrefix(sortField, lookup) {
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
	}
}

func getEventStatsLookup(now types.CpsTime) []bson.M {
	year, month, day := now.Date()
	truncatedInLocation := types.CpsTime{Time: time.Date(year, month, day, 0, 0, 0, 0, now.Location())}

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
			"from":         mongo.AlarmMongoCollection,
			"localField":   "_id",
			"foreignField": "d",
			"as":           "alarms",
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
		"ok_events": bson.M{"$cond": bson.M{
			"if":   "$event_stats.ok",
			"then": "$event_stats.ok",
			"else": 0,
		}},
		"ko_events": bson.M{"$cond": bson.M{
			"if":   "$event_stats.ko",
			"then": "$event_stats.ko",
			"else": 0,
		}},
		"state": bson.M{"$cond": bson.M{
			"if":   "$alarm.v.state.val",
			"then": "$alarm.v.state.val",
			"else": 0,
		}},
	}
}

func getDurationField(now types.CpsTime) bson.M {
	return bson.M{"$subtract": bson.A{
		bson.M{"$cond": bson.M{
			"if":   "$v.resolved",
			"then": "$v.resolved",
			"else": now,
		}},
		"$v.creation_date",
	}}
}
