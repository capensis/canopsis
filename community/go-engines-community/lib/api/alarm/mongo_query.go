package alarm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

const (
	defaultTimeFieldOpened   = "t"
	defaultTimeFieldResolved = "v.resolved"
)

type MongoQueryBuilder struct {
	filterCollection      mongo.DbCollection
	instructionCollection mongo.DbCollection

	defaultSearchByFields []string
	defaultSortBy         string
	defaultSort           string

	// todo remove after remove OldMongoQuery
	fieldsAliases        map[string]string
	fieldsAliasesByRegex map[string]string

	// alarmMatch is match before all lookups
	alarmMatch []bson.M
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
	lookups                     []lookupWithKey
	sort                        bson.M
	computedFieldsForAlarmMatch map[string]bool
	computedFieldsForSort       map[string]bool
	computedFields              bson.M
	// excludedFields is used to remove redundant data from result
	excludedFields []string
}

type lookupWithKey struct {
	key      string
	pipeline []bson.M
}

func NewMongoQueryBuilder(client mongo.DbClient) *MongoQueryBuilder {
	return &MongoQueryBuilder{
		filterCollection:      client.Collection(mongo.WidgetFiltersMongoCollection),
		instructionCollection: client.Collection(mongo.InstructionMongoCollection),

		defaultSearchByFields: []string{
			"v.connector",
			"v.connector_name",
			"v.component",
			"v.resource",
		},
		defaultSortBy: "t",
		defaultSort:   common.SortDesc,

		fieldsAliases: map[string]string{
			"uid":            "_id",
			"connector":      "v.connector",
			"connector_name": "v.connector_name",
			"component":      "v.component",
			"resource":       "v.resource",
			"entity_id":      "d",
			"state":          "v.state.val",
			"status":         "v.status.val",
			"snooze":         "v.snooze",
			"ack":            "v.ack",
			"cancel":         "v.cancel",
			"ticket":         "v.ticket.val",
			"output":         "v.output",
			"opened":         "t",
			"resolved":       "v.resolved",
		},
		fieldsAliasesByRegex: map[string]string{
			"^infos\\.(.+)":           "entity.infos.$1",
			"^v\\.infos\\.\\*\\.(.+)": "v.infos_array.v.$1",
		},
	}
}

func (q *MongoQueryBuilder) clear(now types.CpsTime) {
	q.alarmMatch = make([]bson.M, 0)
	q.additionalMatch = []bson.M{{"$match": bson.M{"entity.enabled": true}}}

	q.lookupsForAdditionalMatch = map[string]bool{"entity": true}
	q.lookupsForSort = make(map[string]bool)
	q.excludeLookupsBeforeSort = make([]string, 0)
	q.lookups = []lookupWithKey{
		{key: "entity", pipeline: getEntityLookup()},
		{key: "entity.category", pipeline: getEntityCategoryLookup()},
		{key: "pbehavior", pipeline: getPbehaviorLookup()},
		{key: "pbehavior.type", pipeline: getPbehaviorTypeLookup()},
		{key: "v.pbehavior_info.icon_name", pipeline: getPbehaviorInfoTypeLookup()},
	}

	q.sort = bson.M{}

	q.computedFieldsForAlarmMatch = make(map[string]bool)
	q.computedFieldsForSort = make(map[string]bool)
	q.computedFields = getComputedFields(now)
	q.excludedFields = []string{"v.steps", "pbehavior.comments", "pbehavior_info_type"}
}

func (q *MongoQueryBuilder) CreateListAggregationPipeline(ctx context.Context, r ListRequestWithPagination, now types.CpsTime) ([]bson.M, error) {
	q.clear(now)

	err := q.handleWidgetFilter(ctx, r.FilterRequest)
	if err != nil {
		return nil, err
	}
	err = q.handleFilter(ctx, r.FilterRequest)
	if err != nil {
		return nil, err
	}
	err = q.handleSort(r.SortRequest)
	if err != nil {
		return nil, err
	}

	return q.createPaginationAggregationPipeline(r.Query), nil
}

func (q *MongoQueryBuilder) CreateCountAggregationPipeline(ctx context.Context, r FilterRequest, now types.CpsTime) ([]bson.M, error) {
	q.clear(now)

	err := q.handleWidgetFilter(ctx, r)
	if err != nil {
		return nil, err
	}
	err = q.handleFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	beforeLimit, _ := q.createAggregationPipeline()

	return beforeLimit, nil
}

func (q *MongoQueryBuilder) CreateGetAggregationPipeline(
	id string,
	now types.CpsTime,
) ([]bson.M, error) {
	q.clear(now)

	q.alarmMatch = append(q.alarmMatch,
		bson.M{"$match": bson.M{"_id": id}},
	)

	query := pagination.Query{
		Page:  1,
		Limit: 1,
	}
	return q.createPaginationAggregationPipeline(query), nil
}

func (q *MongoQueryBuilder) CreateAggregationPipelineByMatch(
	match bson.M,
	r SimpleListRequest,
	now types.CpsTime,
) ([]bson.M, error) {
	q.clear(now)
	q.alarmMatch = append(q.alarmMatch, bson.M{"$match": match})
	err := q.handleSort(SortRequest{
		Sort:   r.Sort,
		SortBy: r.SortBy,
	})
	if err != nil {
		return nil, err
	}

	return q.createPaginationAggregationPipeline(r.Query), nil
}

func (q *MongoQueryBuilder) CreateChildrenAggregationPipeline(
	r ChildDetailsRequest,
	opened int,
	parentId string,
	now types.CpsTime,
) ([]bson.M, error) {
	q.clear(now)

	match := bson.M{
		"v.parents": parentId,
	}
	if opened == OnlyOpened {
		match["v.resolved"] = nil
	}

	q.alarmMatch = append(q.alarmMatch,
		bson.M{"$match": match},
	)

	q.lookups = append(q.lookups, lookupWithKey{key: "parents", pipeline: []bson.M{
		{"$graphLookup": bson.M{
			"from":                    mongo.AlarmMongoCollection,
			"startWith":               "$v.parents",
			"connectFromField":        "v.parents",
			"connectToField":          "d",
			"restrictSearchWithMatch": bson.M{"v.resolved": nil},
			"as":                      "parents",
			"maxDepth":                0,
		}},
		{"$graphLookup": bson.M{
			"from":             mongo.ResolvedAlarmMongoCollection,
			"startWith":        "$v.parents",
			"connectFromField": "v.parents",
			"connectToField":   "d",
			"as":               "resolved_parents",
			"maxDepth":         0,
		}},
		{"$lookup": bson.M{
			"from":         mongo.MetaAlarmRulesMongoCollection,
			"localField":   "parents.v.meta",
			"foreignField": "_id",
			"as":           "meta_alarm_rules",
		}},
		{"$lookup": bson.M{
			"from":         mongo.MetaAlarmRulesMongoCollection,
			"localField":   "resolved_parents.v.meta",
			"foreignField": "_id",
			"as":           "resolved_meta_alarm_rules",
		}},
		{"$addFields": bson.M{
			"meta_alarm_rules": bson.M{"$concatArrays": bson.A{"$meta_alarm_rules", "$resolved_meta_alarm_rules"}},
			"parents": bson.M{"$sum": bson.A{
				bson.M{"$size": "$parents"},
				bson.M{"$size": "$resolved_parents"},
			}},
		}},
	}})
	q.excludedFields = append(q.excludedFields, "resolved_parents", "resolved_meta_alarm_rules")

	err := q.handleSort(r.SortRequest)
	if err != nil {
		return nil, err
	}

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
	addedComputedFields := make(map[string]bool)
	beforeLimit := make([]bson.M, 0)

	q.addFieldsToPipeline(q.computedFieldsForAlarmMatch, addedComputedFields, &beforeLimit)
	beforeLimit = append(beforeLimit, q.alarmMatch...)

	q.addLookupsToPipeline(q.lookupsForAdditionalMatch, addedLookups, &beforeLimit)
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

func (q *MongoQueryBuilder) handleFilter(ctx context.Context, r FilterRequest) error {
	alarmMatch := make([]bson.M, 0)

	q.addOpenedFilter(r, &alarmMatch)
	q.addStartFromFilter(r, &alarmMatch)
	q.addStartToFilter(r, &alarmMatch)
	q.addOnlyParentsFilter(r, &alarmMatch)
	searchMarch, withLookups := q.addSearchFilter(r)
	if len(searchMarch) > 0 {
		if withLookups {
			q.additionalMatch = append(q.additionalMatch, bson.M{"$match": searchMarch})
		} else {
			alarmMatch = append(alarmMatch, searchMarch)
		}
	}

	if len(alarmMatch) > 0 {
		q.alarmMatch = append(q.alarmMatch, bson.M{"$match": bson.M{"$and": alarmMatch}})
	}

	entityMatch := make([]bson.M, 0)
	q.addCategoryFilter(r, &entityMatch)
	err := q.addInstructionsFilter(ctx, r, &entityMatch)
	if err != nil {
		return err
	}
	if len(entityMatch) > 0 {
		q.lookupsForAdditionalMatch["entity"] = true
		q.additionalMatch = append(q.additionalMatch, bson.M{"$match": bson.M{"$and": entityMatch}})
	}

	return nil
}

func (q *MongoQueryBuilder) handleWidgetFilter(ctx context.Context, r FilterRequest) error {
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

	if len(filter.OldMongoQuery) > 0 {
		var query map[string]interface{}
		err := json.Unmarshal([]byte(filter.OldMongoQuery), &query)
		if err != nil {
			return fmt.Errorf("cannot unmarshal old mongo query: %w", err)
		}

		q.computedFieldsForAlarmMatch["v.infos_array"] = true
		q.computedFields["v.infos_array"] = bson.M{"$objectToArray": "$v.infos"}
		resolvedQuery := q.resolveAliasesInQuery(query)
		extraLookups := false

		for _, lookup := range q.lookups {
			if strings.Contains(filter.OldMongoQuery, lookup.key+".") {
				extraLookups = true
				q.lookupsForAdditionalMatch[lookup.key] = true
			}
		}

		if extraLookups {
			q.additionalMatch = append(q.additionalMatch, bson.M{"$match": resolvedQuery})
		} else {
			q.alarmMatch = append(q.alarmMatch, bson.M{"$match": resolvedQuery})
		}

		for field := range q.computedFields {
			if strings.Contains(filter.OldMongoQuery, field) {
				q.computedFieldsForAlarmMatch[field] = true
			}
		}

		return nil
	}

	alarmPatternQuery, err := filter.AlarmPattern.ToMongoQuery("")
	if err != nil {
		return fmt.Errorf("invalid alarm pattern in widget filter id=%q: %w", filter.ID, err)
	}

	if len(alarmPatternQuery) > 0 {
		q.alarmMatch = append(q.alarmMatch, bson.M{"$match": alarmPatternQuery})

		if filter.AlarmPattern.HasInfosField() {
			q.computedFieldsForAlarmMatch["v.infos_array"] = true
			q.computedFields["v.infos_array"] = bson.M{"$objectToArray": "$v.infos"}
		}

		for field := range q.computedFields {
			if filter.AlarmPattern.HasField(field) {
				q.computedFieldsForAlarmMatch[field] = true
			}
		}
	}

	pbhPatternQuery, err := filter.PbehaviorPattern.ToMongoQuery("v.pbehavior_info")
	if err != nil {
		return fmt.Errorf("invalid pbehavior pattern in widget filter id=%q: %w", filter.ID, err)
	}

	if len(pbhPatternQuery) > 0 {
		q.alarmMatch = append(q.alarmMatch, bson.M{"$match": pbhPatternQuery})
	}

	entityPatternQuery, err := filter.EntityPattern.ToMongoQuery("entity")
	if err != nil {
		return fmt.Errorf("invalid entity pattern in widget filter id=%q: %w", filter.ID, err)
	}

	if len(entityPatternQuery) > 0 {
		q.lookupsForAdditionalMatch["entity"] = true
		q.additionalMatch = append(q.additionalMatch, bson.M{"$match": entityPatternQuery})
	}

	return nil
}

func (q *MongoQueryBuilder) addSearchFilter(r FilterRequest) (match bson.M, withLookups bool) {
	if r.Search == "" {
		return
	}

	searchRegexp := primitive.Regex{
		Pattern: fmt.Sprintf(".*%s.*", r.Search),
		Options: "i",
	}

	searchBy := r.SearchBy
	if len(searchBy) == 0 {
		searchBy = q.defaultSearchByFields
	}
	searchMatch := make([]bson.M, len(searchBy))
	for i := range searchBy {
		searchMatch[i] = bson.M{searchBy[i]: searchRegexp}
	}

	if !r.OnlyParents {
		return bson.M{"$or": searchMatch}, false
	}

	childrenMatch := bson.M{"$or": searchMatch}
	childrenCollection := mongo.AlarmMongoCollection
	switch r.GetOpenedFilter() {
	case OnlyOpened:
		childrenMatch = bson.M{"$and": []bson.M{
			{"v.resolved": nil},
			{"$or": searchMatch},
		}}
	case OnlyResolved:
		childrenCollection = mongo.ResolvedAlarmMongoCollection
	}

	q.lookupsForAdditionalMatch["filtered_children"] = true
	q.lookups = append(q.lookups, lookupWithKey{key: "filtered_children", pipeline: getFilteredChildrenLookup(childrenCollection, childrenMatch)})

	searchMatchWithChildren := make([]bson.M, len(searchBy))
	copy(searchMatchWithChildren, searchMatch)
	searchMatchWithChildren = append(searchMatchWithChildren, bson.M{"filtered_children": bson.M{"$ne": []string{}}})

	return bson.M{"$or": searchMatchWithChildren}, true
}

func (q *MongoQueryBuilder) addStartFromFilter(r FilterRequest, match *[]bson.M) {
	if r.StartFrom == nil {
		return
	}

	*match = append(*match, bson.M{q.getTimeField(r): bson.M{"$gte": r.StartFrom}})
}

func (q *MongoQueryBuilder) addStartToFilter(r FilterRequest, match *[]bson.M) {
	if r.StartTo == nil {
		return
	}

	*match = append(*match, bson.M{q.getTimeField(r): bson.M{"$lte": r.StartTo}})
}

func (q *MongoQueryBuilder) getTimeField(r FilterRequest) string {
	if r.TimeField == defaultTimeFieldOpened {
		return r.TimeField
	}

	if r.TimeField == "" {
		if r.GetOpenedFilter() == OnlyResolved {
			return defaultTimeFieldResolved
		}

		return defaultTimeFieldOpened
	}

	return r.TimeField
}

func (q *MongoQueryBuilder) addOpenedFilter(r FilterRequest, match *[]bson.M) {
	if r.GetOpenedFilter() != OnlyOpened {
		return
	}

	*match = append(*match, bson.M{"v.resolved": nil})
}

func (q *MongoQueryBuilder) addCategoryFilter(r FilterRequest, match *[]bson.M) {
	if r.Category == "" {
		return
	}

	*match = append(*match, bson.M{"entity.category": bson.M{"$eq": r.Category}})
}

func (q *MongoQueryBuilder) addOnlyParentsFilter(r FilterRequest, match *[]bson.M) {
	if !r.OnlyParents {
		*match = append(*match, bson.M{"v.meta": nil})
		return
	}

	*match = append(*match, bson.M{"$or": []bson.M{
		{"v.parents": nil},
		{"v.parents": bson.M{"$eq": bson.A{}}},
		{"v.meta": bson.M{"$ne": nil}},
	}})

	q.computedFields["is_meta_alarm"] = getIsMetaAlarmField()
	q.lookups = append(q.lookups, lookupWithKey{key: "meta_alarm_rule", pipeline: getMetaAlarmRuleLookup()})
	q.lookups = append(q.lookups, lookupWithKey{key: "children", pipeline: getChildrenCountLookup()})
	q.excludedFields = append(q.excludedFields, "resolved_children")
}

func (q *MongoQueryBuilder) addInstructionsFilter(ctx context.Context, r FilterRequest, match *[]bson.M) error {
	added := false

	if len(r.ExcludeInstructions) > 0 {
		filters, err := q.getInstructionsFilters(ctx, bson.M{"_id": bson.M{"$in": r.ExcludeInstructions}})
		if err != nil {
			return err
		}
		if len(filters) > 0 {
			added = true
			*match = append(*match, bson.M{"$nor": filters})
		}
	}

	if len(r.ExcludeInstructionTypes) > 0 {
		filters, err := q.getInstructionsFilters(ctx, bson.M{"type": bson.M{"$in": r.ExcludeInstructionTypes}})
		if err != nil {
			return err
		}
		if len(filters) > 0 {
			added = true
			*match = append(*match, bson.M{"$nor": filters})
		}
	}

	if len(r.IncludeInstructions) > 0 {
		filters, err := q.getInstructionsFilters(ctx, bson.M{"_id": bson.M{"$in": r.IncludeInstructions}})
		if err != nil {
			return err
		}
		if len(filters) > 0 {
			added = true
			*match = append(*match, bson.M{"$or": filters})
		} else {
			*match = append(*match, bson.M{"$nor": []bson.M{{}}})
		}
	}

	if len(r.IncludeInstructionTypes) > 0 {
		filters, err := q.getInstructionsFilters(ctx, bson.M{"type": bson.M{"$in": r.IncludeInstructionTypes}})
		if err != nil {
			return err
		}
		if len(filters) > 0 {
			added = true
			*match = append(*match, bson.M{"$or": filters})
		} else {
			*match = append(*match, bson.M{"$nor": []bson.M{{}}})
		}
	}

	if added {
		q.computedFieldsForAlarmMatch["v.duration"] = true
		q.computedFieldsForAlarmMatch["v.infos_array"] = true
		q.computedFields["v.infos_array"] = bson.M{"$objectToArray": "$v.infos"}
	}

	return nil
}

func (q *MongoQueryBuilder) getInstructionsFilters(ctx context.Context, filter bson.M) ([]bson.M, error) {
	filter["status"] = bson.M{"$in": bson.A{InstructionStatusApproved, nil}}
	cursor, err := q.instructionCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch instructions: %w", err)
	}
	defer cursor.Close(ctx)

	var filters []bson.M

	for cursor.Next(ctx) {
		var instruction Instruction
		err := cursor.Decode(&instruction)
		if err != nil {
			return nil, fmt.Errorf("cannot decode instruction: %w", err)
		}

		q, err := getInstructionQuery(instruction)
		if err != nil {
			return nil, err
		}

		if q != nil {
			filters = append(filters, q)
		}
	}

	return filters, nil
}

func (q *MongoQueryBuilder) handleSort(r SortRequest) error {
	if len(r.MultiSort) > 0 {
		idExist := false
		sortQuery := bson.D{}
		sortFields := make([]string, 0)

		for _, v := range r.MultiSort {
			split := strings.Split(v, ",")
			if len(split) != 2 {
				return errors.New("length of multi_sort value should be equal 2")
			}

			sortBy := split[0]
			sortDir := 1
			if split[1] == common.SortDesc {
				sortDir = -1
			}

			if sortBy == "_id" {
				idExist = true
			}

			sortFields = append(sortFields, sortBy)
			sortQuery = append(sortQuery, bson.E{Key: sortBy, Value: sortDir})
		}

		q.adjustLookupsForSort(sortFields)

		if !idExist {
			sortQuery = append(sortQuery, bson.E{Key: "_id", Value: 1})
		}

		q.sort = bson.M{"$sort": sortQuery}
		return nil
	}

	sortBy := r.SortBy
	if sortBy == "" {
		sortBy = q.defaultSortBy
	}
	sort := r.Sort
	if sort == "" {
		sort = q.defaultSort
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

func (q *MongoQueryBuilder) resolveAliasesInQuery(query interface{}) interface{} {
	res := query
	val := reflect.ValueOf(res)

	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			newVal := q.resolveAliasesInQuery(val.Index(i).Interface())
			val.Index(i).Set(reflect.ValueOf(newVal))
		}
	case reflect.Map:
		for _, key := range val.MapKeys() {
			newVal := q.resolveAliasesInQuery(val.MapIndex(key).Interface())
			newKey := q.resolveAlias(key.String())

			var mapVal reflect.Value
			if newVal == nil {
				mapVal = reflect.ValueOf(&newVal).Elem()
			} else {
				mapVal = reflect.ValueOf(newVal)
			}

			val.SetMapIndex(key, reflect.Value{})
			val.SetMapIndex(reflect.ValueOf(newKey), mapVal)
		}
	}

	return res
}

func (q *MongoQueryBuilder) resolveAlias(v string) string {
	if v == "" {
		return v
	}

	for alias, field := range q.fieldsAliases {
		if alias == v {
			return field
		}
	}

	for expr, repl := range q.fieldsAliasesByRegex {
		r, err := regexp.Compile(expr)
		if err == nil {
			replace := r.ReplaceAllString(v, repl)
			if v != replace {
				return replace
			}
		}
	}

	return v
}

func getEntityLookup() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "d",
			"foreignField": "_id",
			"as":           "entity",
		}},
		{"$unwind": "$entity"},
	}
}

func getEntityCategoryLookup() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.EntityCategoryMongoCollection,
			"localField":   "entity.category",
			"foreignField": "_id",
			"as":           "entity.category",
		}},
		{"$unwind": bson.M{"path": "$entity.category", "preserveNullAndEmptyArrays": true}},
	}
}

func getPbehaviorLookup() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorMongoCollection,
			"foreignField": "_id",
			"localField":   "v.pbehavior_info.id",
			"as":           "pbehavior",
		}},
		{"$unwind": bson.M{"path": "$pbehavior", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"pbehavior.last_comment": bson.M{"$arrayElemAt": bson.A{"$pbehavior.comments", -1}},
		}},
	}
}

func getPbehaviorTypeLookup() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorTypeMongoCollection,
			"foreignField": "_id",
			"localField":   "pbehavior.type_",
			"as":           "pbehavior.type",
		}},
		{"$unwind": bson.M{"path": "$pbehavior.type", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"pbehavior": bson.M{
				"$cond": bson.M{
					"if":   "$pbehavior._id",
					"then": "$pbehavior",
					"else": nil,
				},
			},
		}},
	}
}

func getPbehaviorInfoTypeLookup() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorTypeMongoCollection,
			"foreignField": "_id",
			"localField":   "v.pbehavior_info.type",
			"as":           "pbehavior_info_type",
		}},
		{"$unwind": bson.M{"path": "$pbehavior_info_type", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"v.pbehavior_info": bson.M{"$cond": bson.M{
				"if": "$v.pbehavior_info",
				"then": bson.M{"$mergeObjects": bson.A{
					"$v.pbehavior_info",
					bson.M{"icon_name": "$pbehavior_info_type.icon_name"},
				}},
				"else": nil,
			}},
		}},
	}
}

func getMetaAlarmRuleLookup() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.MetaAlarmRulesMongoCollection,
			"localField":   "v.meta",
			"foreignField": "_id",
			"as":           "meta_alarm_rule",
		}},
		{"$unwind": bson.M{"path": "$meta_alarm_rule", "preserveNullAndEmptyArrays": true}},
	}
}

func getChildrenCountLookup() []bson.M {
	return []bson.M{
		{"$graphLookup": bson.M{
			"from":                    mongo.AlarmMongoCollection,
			"startWith":               "$d",
			"connectFromField":        "d",
			"connectToField":          "v.parents",
			"restrictSearchWithMatch": bson.M{"v.resolved": nil},
			"as":                      "children",
			"maxDepth":                0,
		}},
		{"$graphLookup": bson.M{
			"from":             mongo.ResolvedAlarmMongoCollection,
			"startWith":        "$d",
			"connectFromField": "d",
			"connectToField":   "v.parents",
			"as":               "resolved_children",
			"maxDepth":         0,
		}},
		{"$addFields": bson.M{
			"children": bson.M{"$sum": bson.A{
				bson.M{"$size": "$children"},
				bson.M{"$size": "$resolved_children"},
			}},
		}},
	}
}

func getFilteredChildrenLookup(childrenCollection string, childrenMatch bson.M) []bson.M {
	return []bson.M{
		{"$graphLookup": bson.M{
			"from":                    childrenCollection,
			"startWith":               "$d",
			"connectFromField":        "d",
			"connectToField":          "v.parents",
			"restrictSearchWithMatch": childrenMatch,
			"as":                      "filtered_children",
			"maxDepth":                0,
		}},
		{"$addFields": bson.M{
			"filtered_children": "$filtered_children._id",
		}},
	}
}

func getComputedFields(now types.CpsTime) bson.M {
	return bson.M{
		"infos":        "$v.infos",
		"impact_state": bson.M{"$multiply": bson.A{"$v.state.val", "$entity.impact_level"}},
		"v.duration": bson.M{"$subtract": bson.A{
			bson.M{"$cond": bson.M{
				"if":   "$v.resolved",
				"then": "$v.resolved",
				"else": now,
			}},
			"$v.creation_date",
		}},
		"v.current_state_duration": bson.M{"$subtract": bson.A{
			bson.M{"$cond": bson.M{
				"if":   "$v.resolved",
				"then": "$v.resolved",
				"else": now,
			}},
			"$v.state.t",
		}},
		"v.active_duration": bson.M{"$subtract": bson.A{
			bson.M{"$cond": bson.M{
				"if":   "$v.resolved",
				"then": "$v.resolved",
				"else": now,
			}},
			bson.M{"$sum": bson.A{
				"$v.creation_date",
				"$v.snooze_duration",
				"$v.pbh_inactive_duration",
			}},
		}},
		"v.last_comment": bson.M{"$arrayElemAt": bson.A{
			bson.M{"$filter": bson.M{
				"input": "$v.steps",
				"cond":  bson.M{"$eq": bson.A{"$$this._t", types.AlarmStepComment}},
			}},
			-1,
		}},
	}
}

func getIsMetaAlarmField() bson.M {
	return bson.M{"$cond": bson.A{bson.M{"$not": bson.A{"$v.meta"}}, false, true}}
}

func getInstructionQuery(instruction Instruction) (bson.M, error) {
	var and []bson.M

	if len(instruction.AlarmPattern) == 0 && len(instruction.EntityPattern) == 0 &&
		(!instruction.OldAlarmPatterns.IsSet() || !instruction.OldAlarmPatterns.IsValid()) &&
		(!instruction.OldEntityPatterns.IsSet() || !instruction.OldEntityPatterns.IsValid()) {
		return nil, nil
	}

	alarmPatternQuery, err := instruction.AlarmPattern.ToMongoQuery("")
	if err != nil {
		return nil, fmt.Errorf("invalid alarm pattern in instruction id=%q: %w", instruction.ID, err)
	}
	if len(alarmPatternQuery) > 0 {
		and = append(and, alarmPatternQuery)
	} else {
		oldAlarmPatternQuery := instruction.OldAlarmPatterns.AsMongoDriverQuery()
		if len(oldAlarmPatternQuery) > 0 {
			and = append(and, oldAlarmPatternQuery)
		}
	}

	entityPatternQuery, err := instruction.EntityPattern.ToMongoQuery("entity")
	if err != nil {
		return nil, fmt.Errorf("invalid entuty pattern in instruction id=%q: %w", instruction.ID, err)
	}
	if len(entityPatternQuery) > 0 {
		and = append(and, entityPatternQuery)
	} else {
		oldEntityPatternQuery := instruction.OldEntityPatterns.AsMongoDriverQuery()
		if len(oldEntityPatternQuery) > 0 {
			and = append(and, getEntityPatternsForEntity(oldEntityPatternQuery))
		}
	}

	if len(instruction.ActiveOnPbh) > 0 {
		and = append(and, bson.M{"v.pbehavior_info.type": bson.M{"$in": instruction.ActiveOnPbh}})
	}

	if len(instruction.DisabledOnPbh) > 0 {
		and = append(and, bson.M{"v.pbehavior_info.type": bson.M{"$nin": instruction.DisabledOnPbh}})
	}

	return bson.M{"$and": and}, nil
}

func getEntityPatternsForEntity(patternBson bson.M) bson.M {
	newBson := make(bson.M)
	patternListInterface, ok := patternBson["$or"]
	if ok {
		patternList := patternListInterface.([]bson.M)
		newPatternsList := make([]bson.M, len(patternList))
		for i, pattern := range patternList {
			newPattern := make(bson.M)
			for k, vv := range pattern {
				newPattern["entity."+k] = vv
			}

			newPatternsList[i] = newPattern
		}

		// Just in case when an entity pattern's function AsMongoDriverQuery returns an empty bson array
		// that might happen if the pattern has bad format in mongo, but after unmarshalling it has isSet = true
		// since it's not possible for $or has empty array we just fill it with an empty bson
		if len(newPatternsList) == 0 {
			newPatternsList = append(newPatternsList, bson.M{})
		}

		newBson["$or"] = newPatternsList
	} else {
		return patternBson
	}

	return newBson
}
