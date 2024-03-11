package alarm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/expression/parser"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

const (
	defaultTimeFieldOpened   = "t"
	defaultTimeFieldResolved = "v.resolved"
	entityInfosPrefix        = "entity.infos"
)

type MongoQueryBuilder struct {
	filterCollection      mongo.DbCollection
	instructionCollection mongo.DbCollection
	authorProvider        author.Provider

	defaultSearchByFields         []string
	availableSearchByFields       map[string]struct{}
	availableSearchByEntityFields map[string]struct{}
	defaultSortBy                 string
	defaultSort                   string

	fieldsAliases        map[string]string
	fieldsAliasesByRegex map[string]string

	searchPipeline []bson.M
	// alarmMatch is match before all lookups
	alarmMatch []bson.M
	// additionalMatch is match after some lookups
	additionalMatch []bson.M
	// lookupsForAdditionalMatch is required for match and for result
	lookupsForAdditionalMatch     map[string]bool
	lookupsOnlyForAdditionalMatch map[string]bool
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

func NewMongoQueryBuilder(client mongo.DbClient, authorProvider author.Provider) *MongoQueryBuilder {
	return &MongoQueryBuilder{
		filterCollection:      client.Collection(mongo.WidgetFiltersMongoCollection),
		instructionCollection: client.Collection(mongo.InstructionMongoCollection),
		authorProvider:        authorProvider,

		defaultSearchByFields: []string{
			"v.connector",
			"v.connector_name",
			"v.component",
			"v.resource",
		},
		availableSearchByFields: map[string]struct{}{
			"v.connector":      {},
			"v.connector_name": {},
			"v.component":      {},
			"v.resource":       {},
			"v.display_name":   {},
			"v.output":         {},
			"v.ticket.ticket":  {},
		},
		availableSearchByEntityFields: map[string]struct{}{
			"entity.name":        {},
			"entity.connector":   {},
			"entity.component":   {},
			"entity.description": {},
			"entity.type":        {},
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
			"ticket":         "v.ticket.ticket",
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

func (q *MongoQueryBuilder) clear(now types.CpsTime, userID string) {
	q.searchPipeline = make([]bson.M, 0)
	q.alarmMatch = []bson.M{
		{"$match": bson.M{
			"healthcheck": bson.M{"$in": bson.A{nil, false}},
		}},
	}
	q.additionalMatch = make([]bson.M, 0)

	q.lookupsForAdditionalMatch = make(map[string]bool)
	q.lookupsOnlyForAdditionalMatch = make(map[string]bool)
	q.lookupsForSort = make(map[string]bool)
	q.excludeLookupsBeforeSort = make([]string, 0)
	q.lookups = []lookupWithKey{
		{key: "entity", pipeline: getEntityLookup()},
		{key: "entity.category", pipeline: getEntityCategoryLookup()},
		{key: "pbehavior", pipeline: getPbehaviorLookup(q.authorProvider)},
		{key: "pbehavior.type", pipeline: getPbehaviorTypeLookup()},
		{key: "v.pbehavior_info.icon_name", pipeline: getPbehaviorInfoTypeLookup()},
	}

	q.sort = bson.M{}

	q.computedFieldsForAlarmMatch = make(map[string]bool)
	q.computedFieldsForSort = make(map[string]bool)
	q.computedFields = getComputedFields(now, userID)
	q.excludedFields = []string{"bookmarks", "v.steps", "pbehavior.comments", "pbehavior_info_type", "entity.services"}
}

func (q *MongoQueryBuilder) CreateGetDisplayNamesPipeline(r GetDisplayNamesRequest, now types.CpsTime) ([]bson.M, error) {
	q.clear(now, "")

	err := q.handlePatterns(FilterRequest{
		BaseFilterRequest: BaseFilterRequest{
			AlarmPattern:     r.AlarmPattern,
			EntityPattern:    r.EntityPattern,
			PbehaviorPattern: r.PbehaviorPattern,
		},
	})
	if err != nil {
		return nil, err
	}

	match := bson.M{"v.resolved": nil}

	if r.Search != "" {
		match["v.display_name"] = primitive.Regex{
			Pattern: fmt.Sprintf(".*%s.*", r.Search),
			Options: "i",
		}
	}

	q.alarmMatch = append(q.alarmMatch, bson.M{"$match": match})

	sortDir := 1
	if r.Sort == common.SortDesc {
		sortDir = -1
	}

	q.sort = bson.M{"$sort": bson.M{"v.display_name": sortDir}}

	// maps are not used need to call functions below
	addedLookups := make(map[string]bool)
	addedComputedFields := make(map[string]bool)

	pipeline := make([]bson.M, 0)
	q.addFieldsToPipeline(q.computedFieldsForAlarmMatch, addedComputedFields, &pipeline)
	pipeline = append(pipeline, q.alarmMatch...)

	q.addLookupsToPipeline(q.lookupsForAdditionalMatch, addedLookups, &pipeline)
	pipeline = append(pipeline, q.additionalMatch...)

	return pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		q.sort,
		[]bson.M{{"$project": bson.M{"_id": 1, "display_name": "$v.display_name"}}},
	), nil
}

func (q *MongoQueryBuilder) CreateListAggregationPipeline(ctx context.Context, r ListRequestWithPagination, now types.CpsTime, userID string) ([]bson.M, error) {
	q.clear(now, userID)

	err := q.handleWidgetFilter(ctx, r.FilterRequest)
	if err != nil {
		return nil, err
	}
	err = q.handlePatterns(r.FilterRequest)
	if err != nil {
		return nil, err
	}
	err = q.handleFilter(ctx, r.FilterRequest, userID)
	if err != nil {
		return nil, err
	}
	err = q.handleSort(r.SortRequest)
	if err != nil {
		return nil, err
	}
	q.handleDependencies(r.WithDependencies)

	return q.createPaginationAggregationPipeline(r.Query), nil
}

func (q *MongoQueryBuilder) CreateCountAggregationPipeline(ctx context.Context, r FilterRequest, userID string, now types.CpsTime) ([]bson.M, error) {
	q.clear(now, userID)

	err := q.handleWidgetFilter(ctx, r)
	if err != nil {
		return nil, err
	}
	err = q.handleFilter(ctx, r, userID)
	if err != nil {
		return nil, err
	}

	beforeLimit, _ := q.createAggregationPipeline()

	return beforeLimit, nil
}

func (q *MongoQueryBuilder) CreateGetAggregationPipeline(
	match bson.M,
	now types.CpsTime,
	userID string,
	opened int,
	onlyParents bool,
) ([]bson.M, error) {
	q.clear(now, userID)
	q.handleOpened(opened)
	q.handleDependencies(true)
	q.alarmMatch = append(q.alarmMatch, bson.M{"$match": match})
	if onlyParents {
		q.computedFields["is_meta_alarm"] = getIsMetaAlarmField()
		q.lookups = append(q.lookups, lookupWithKey{key: "meta_alarm_rule", pipeline: getMetaAlarmRuleLookup()})
		q.lookups = append(q.lookups, lookupWithKey{key: "children", pipeline: getChildrenCountLookup()})
		q.excludedFields = append(q.excludedFields, "resolved_children")
	}

	query := pagination.Query{
		Page:  1,
		Limit: 1,
	}

	return q.createPaginationAggregationPipeline(query), nil
}

func (q *MongoQueryBuilder) CreateAggregationPipelineByMatch(
	ctx context.Context,
	alarmMatch bson.M,
	entityMatch bson.M,
	paginationQuery pagination.Query,
	sortRequest SortRequest,
	filterRequest FilterRequest,
	now types.CpsTime,
	userID string,
) ([]bson.M, error) {
	q.clear(now, userID)
	if len(alarmMatch) > 0 {
		q.alarmMatch = append(q.alarmMatch, bson.M{"$match": alarmMatch})
	}
	if len(entityMatch) > 0 {
		q.lookupsForAdditionalMatch["entity"] = true
		q.additionalMatch = append(q.additionalMatch, bson.M{"$match": entityMatch})
	}

	err := q.handleFilter(ctx, filterRequest, userID)
	if err != nil {
		return nil, err
	}

	err = q.handleSort(sortRequest)
	if err != nil {
		return nil, err
	}
	q.handleDependencies(true)

	return q.createPaginationAggregationPipeline(paginationQuery), nil
}

func (q *MongoQueryBuilder) CreateChildrenAggregationPipeline(
	r ChildDetailsRequest,
	opened int,
	parentId string,
	search string,
	userID string,
	searchBy []string,
	now types.CpsTime,
) ([]bson.M, error) {
	q.clear(now, userID)
	q.handleOpened(opened)
	q.handleDependencies(true)
	q.alarmMatch = append(q.alarmMatch, bson.M{"$match": bson.M{"v.parents": parentId}})
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

	if search != "" {
		p := parser.NewParser()
		expr, err := p.Parse(search)
		if err == nil {
			query := expr.ExprQuery()
			resolvedQuery := q.resolveAliasesInQuery(query).(bson.M)
			if err != nil {
				return nil, err
			}

			b, err := json.Marshal(resolvedQuery)
			if err != nil {
				return nil, fmt.Errorf("cannot marshal search expression: %w", err)
			}

			resolvedSearch := string(b)
			for field := range q.computedFields {
				if strings.Contains(resolvedSearch, field) {
					q.computedFieldsForAlarmMatch[field] = true
				}
			}

			q.computedFields["filtered"] = bson.M{"$cond": bson.M{
				"if":   resolvedQuery,
				"then": true,
				"else": false,
			}}
		} else {
			filteredSearchBy := make([]string, 0, len(searchBy))
			for _, f := range searchBy {
				if _, ok := q.availableSearchByFields[f]; ok {
					filteredSearchBy = append(filteredSearchBy, f)
					continue
				}

				if _, ok := q.availableSearchByEntityFields[f]; ok {
					filteredSearchBy = append(filteredSearchBy, f)
					continue
				}

				if strings.HasPrefix(f, entityInfosPrefix) {
					filteredSearchBy = append(filteredSearchBy, f)
					continue
				}
			}

			if len(filteredSearchBy) == 0 {
				filteredSearchBy = q.defaultSearchByFields
			}
			searchMatch := make([]bson.M, len(filteredSearchBy))
			for i := range filteredSearchBy {
				searchMatch[i] = bson.M{"$regexMatch": bson.M{
					"input":   "$" + filteredSearchBy[i],
					"regex":   fmt.Sprintf(".*%s.*", search),
					"options": "i",
				}}
			}

			q.computedFields["filtered"] = bson.M{"$cond": bson.M{
				"if":   bson.M{"$or": searchMatch},
				"then": true,
				"else": false,
			}}
		}
	}

	err := q.handleSort(r.SortRequest)
	if err != nil {
		return nil, err
	}

	return q.createPaginationAggregationPipeline(r.Query), nil
}

func (q *MongoQueryBuilder) CreateOnlyListAggregationPipeline(
	ctx context.Context,
	r ListRequest,
	now types.CpsTime,
	userID string,
) ([]bson.M, error) {
	q.clear(now, userID)

	err := q.handleWidgetFilter(ctx, r.FilterRequest)
	if err != nil {
		return nil, err
	}
	err = q.handleFilter(ctx, r.FilterRequest, userID)
	if err != nil {
		return nil, err
	}
	err = q.handleSort(r.SortRequest)
	if err != nil {
		return nil, err
	}
	q.handleDependencies(r.WithDependencies)

	beforeLimit, afterLimit := q.createAggregationPipeline()
	pipeline := append(beforeLimit, q.sort)
	pipeline = append(pipeline, afterLimit...)
	return pipeline, nil
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
	beforeLimit := make([]bson.M, len(q.searchPipeline))
	copy(beforeLimit, q.searchPipeline)

	q.addFieldsToPipeline(q.computedFieldsForAlarmMatch, addedComputedFields, &beforeLimit)
	beforeLimit = append(beforeLimit, q.alarmMatch...)

	q.addLookupsToPipeline(q.lookupsForAdditionalMatch, addedLookups, &beforeLimit)
	q.addLookupsToPipeline(q.lookupsOnlyForAdditionalMatch, addedLookups, &beforeLimit)
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
		if !addedLookups[lookup.key] && !q.lookupsOnlyForAdditionalMatch[lookup.key] {
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

func (q *MongoQueryBuilder) handleFilter(ctx context.Context, r FilterRequest, userID string) error {
	alarmMatch := make([]bson.M, 0)
	entityMatch := make([]bson.M, 0)

	q.addOpenedFilter(r.GetOpenedFilter(), &alarmMatch, &entityMatch)
	q.addStartFromFilter(r, &alarmMatch)
	q.addStartToFilter(r, &alarmMatch)
	q.addOnlyParentsFilter(r, &alarmMatch)
	q.addTagFilter(r, &alarmMatch)
	q.addBookmarkFilter(r, userID, &alarmMatch)
	searchMarch, withLookups, err := q.addSearchFilter(r)
	if err != nil {
		return err
	}
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

	q.addCategoryFilter(r, &entityMatch)
	err = q.addInstructionsFilter(ctx, r, &entityMatch)
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
	if len(r.Filters) == 0 {
		return nil
	}

	for _, v := range r.Filters {
		filter := view.WidgetFilter{}
		err := q.filterCollection.FindOne(ctx, bson.M{"_id": v}).Decode(&filter)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return common.NewValidationError("filter", "Filter doesn't exist.")
			}
			return fmt.Errorf("cannot fetch widget filter: %w", err)
		}

		err = q.handleAlarmPattern(filter.AlarmPattern)
		if err != nil {
			return fmt.Errorf("invalid alarm pattern in widget filter id=%q: %w", filter.ID, err)
		}

		err = q.handlePbehaviorPattern(filter.PbehaviorPattern)
		if err != nil {
			return fmt.Errorf("invalid pbehavior pattern in widget filter id=%q: %w", filter.ID, err)
		}

		err = q.handleEntityPattern(filter.EntityPattern)
		if err != nil {
			return fmt.Errorf("invalid entity pattern in widget filter id=%q: %w", filter.ID, err)
		}

		if len(filter.AlarmPattern) == 0 && len(filter.PbehaviorPattern) == 0 && len(filter.EntityPattern) == 0 &&
			len(filter.OldMongoQuery) > 0 {
			var query map[string]any
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
		}
	}

	return nil
}

func (q *MongoQueryBuilder) handlePatterns(r FilterRequest) error {
	if r.AlarmPattern != "" {
		var alarmPattern pattern.Alarm
		err := json.Unmarshal([]byte(r.AlarmPattern), &alarmPattern)
		if err != nil {
			return common.NewValidationError("alarm_pattern", "AlarmPattern is invalid.")
		}
		err = q.handleAlarmPattern(alarmPattern)
		if err != nil {
			return common.NewValidationError("alarm_pattern", "AlarmPattern is invalid.")
		}
	}

	if r.PbehaviorPattern != "" {
		var pbehaviorPattern pattern.PbehaviorInfo
		err := json.Unmarshal([]byte(r.PbehaviorPattern), &pbehaviorPattern)
		if err != nil {
			return common.NewValidationError("pbehavior_pattern", "PbehaviorPattern is invalid.")
		}
		err = q.handlePbehaviorPattern(pbehaviorPattern)
		if err != nil {
			return common.NewValidationError("pbehavior_pattern", "PbehaviorPattern is invalid.")
		}
	}

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

	return nil
}

func (q *MongoQueryBuilder) handleAlarmPattern(alarmPattern pattern.Alarm) error {
	alarmPatternQuery, err := alarmPattern.ToMongoQuery("")
	if err != nil {
		return err
	}

	if len(alarmPatternQuery) > 0 {
		q.alarmMatch = append(q.alarmMatch, bson.M{"$match": alarmPatternQuery})

		if alarmPattern.HasInfosField() {
			q.computedFieldsForAlarmMatch["v.infos_array"] = true
			q.computedFields["v.infos_array"] = bson.M{"$objectToArray": "$v.infos"}
		}

		for field := range q.computedFields {
			if alarmPattern.HasField(field) {
				q.computedFieldsForAlarmMatch[field] = true
			}
		}
	}

	return nil
}

func (q *MongoQueryBuilder) handlePbehaviorPattern(pbehaviorPattern pattern.PbehaviorInfo) error {
	pbhPatternQuery, err := pbehaviorPattern.ToMongoQuery("v")
	if err != nil {
		return err
	}

	if len(pbhPatternQuery) > 0 {
		q.alarmMatch = append(q.alarmMatch, bson.M{"$match": pbhPatternQuery})
	}

	return nil
}

func (q *MongoQueryBuilder) handleEntityPattern(entityPattern pattern.Entity) error {
	entityPatternQuery, err := entityPattern.ToMongoQuery("entity")
	if err != nil {
		return err
	}

	if len(entityPatternQuery) > 0 {
		q.lookupsForAdditionalMatch["entity"] = true
		q.additionalMatch = append(q.additionalMatch, bson.M{"$match": entityPatternQuery})
	}

	return nil
}

func (q *MongoQueryBuilder) addSearchFilter(r FilterRequest) (bson.M, bool, error) {
	if r.Search == "" {
		return nil, false, nil
	}

	p := parser.NewParser()
	expr, err := p.Parse(r.Search)
	if err == nil {
		query := expr.Query()
		resolvedQuery := q.resolveAliasesInQuery(query).(bson.M)
		b, err := json.Marshal(resolvedQuery)
		if err != nil {
			return nil, false, fmt.Errorf("cannot marshal search expression: %w", err)
		}
		resolvedSearch := string(b)
		extraLookups := false

		for _, lookup := range q.lookups {
			if strings.Contains(resolvedSearch, lookup.key+".") {
				extraLookups = true
				q.lookupsForAdditionalMatch[lookup.key] = true
			}
		}

		for field := range q.computedFields {
			if strings.Contains(resolvedSearch, field) {
				q.computedFieldsForAlarmMatch[field] = true
			}
		}

		return resolvedQuery, extraLookups, nil
	}

	searchRegexp := primitive.Regex{
		Pattern: fmt.Sprintf(".*%s.*", r.Search),
		Options: "i",
	}

	searchBy := make([]string, 0, len(r.SearchBy))
	searchByEntity := false
	for _, f := range r.SearchBy {
		if _, ok := q.availableSearchByFields[f]; ok {
			searchBy = append(searchBy, f)
			continue
		}

		if _, ok := q.availableSearchByEntityFields[f]; ok {
			searchBy = append(searchBy, f)
			searchByEntity = true
			continue
		}

		if strings.HasPrefix(f, entityInfosPrefix) {
			searchBy = append(searchBy, f)
			searchByEntity = true
			continue
		}
	}

	if len(searchBy) == 0 {
		searchBy = q.defaultSearchByFields
	}

	searchMatch := make([]bson.M, len(searchBy))
	for i := range searchBy {
		searchMatch[i] = bson.M{searchBy[i]: searchRegexp}
	}

	if !r.OnlyParents {
		if searchByEntity {
			q.lookupsForAdditionalMatch["entity"] = true
		}

		return bson.M{"$or": searchMatch}, searchByEntity, nil
	}

	match := bson.M{"$or": searchMatch}
	metaAlarmLookupCollection := mongo.AlarmMongoCollection
	metaAlarmLookupMatch := bson.M{}
	switch r.GetOpenedFilter() {
	case OnlyOpened:
		match = bson.M{
			"v.resolved": nil,
			"$or":        searchMatch,
		}
		metaAlarmLookupMatch = bson.M{"v.resolved": nil}
	case OnlyResolved:
		metaAlarmLookupCollection = mongo.ResolvedAlarmMongoCollection
	}

	q.searchPipeline = getOnlyParentsSearchPipeline(match, metaAlarmLookupCollection, metaAlarmLookupMatch, searchByEntity)

	return nil, false, nil
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

func (q *MongoQueryBuilder) addOpenedFilter(opened int, match *[]bson.M, entityMatch *[]bson.M) {
	if opened == OnlyOpened {
		*match = append(*match, bson.M{"v.resolved": nil})
		return
	}

	// disabled entity cannot have open alarm but can have resolved
	*entityMatch = append(*entityMatch, bson.M{"entity.enabled": true})
}

func (q *MongoQueryBuilder) addCategoryFilter(r FilterRequest, match *[]bson.M) {
	if r.Category == "" {
		return
	}

	*match = append(*match, bson.M{"entity.category": bson.M{"$eq": r.Category}})
}

func (q *MongoQueryBuilder) addTagFilter(r FilterRequest, match *[]bson.M) {
	if r.Tag == "" {
		return
	}

	*match = append(*match, bson.M{"tags": r.Tag})
}

func (q *MongoQueryBuilder) addBookmarkFilter(r FilterRequest, userID string, match *[]bson.M) {
	if !r.OnlyBookmarks {
		return
	}

	*match = append(*match, bson.M{"bookmarks": userID})
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
	withMatch := false
	withExecution := false
	withExecutionType := false

	for _, instructionFilter := range r.Instructions {
		if len(instructionFilter.Exclude) > 0 {
			if instructionFilter.Running != nil && *instructionFilter.Running {
				withExecution = true
				*match = append(*match, bson.M{"instruction_execution.instruction": bson.M{"$nin": instructionFilter.Exclude}})
				continue
			}

			filters, err := q.getInstructionsFilters(
				ctx,
				bson.M{"_id": bson.M{"$in": instructionFilter.Exclude}},
			)
			if err != nil {
				return err
			}
			if len(filters) > 0 {
				if instructionFilter.Running == nil {
					withMatch = true
					*match = append(*match, bson.M{"$nor": filters})
				} else {
					withExecution = true
					withMatch = true
					*match = append(*match, bson.M{"instruction_execution.instruction": bson.M{"$or": []bson.M{
						{"$in": instructionFilter.Exclude},
						{"$nor": filters},
					}}})
				}
			}
			continue
		}

		if len(instructionFilter.ExcludeTypes) > 0 {
			if instructionFilter.Running != nil && *instructionFilter.Running {
				withExecutionType = true
				*match = append(*match, bson.M{"instruction_execution.type": bson.M{"$nin": instructionFilter.ExcludeTypes}})
				continue
			}

			filters, err := q.getInstructionsFilters(
				ctx,
				bson.M{"type": bson.M{"$in": instructionFilter.ExcludeTypes}},
			)
			if err != nil {
				return err
			}
			if len(filters) > 0 {
				if instructionFilter.Running == nil {
					withMatch = true
					*match = append(*match, bson.M{"$nor": filters})
				} else {
					withExecutionType = true
					withMatch = true
					*match = append(*match, bson.M{"$or": []bson.M{
						{"instruction_execution.type": bson.M{"$in": instructionFilter.ExcludeTypes}},
						{"$nor": filters},
					}})
				}
			}
			continue
		}

		if len(instructionFilter.Include) > 0 {
			if instructionFilter.Running != nil && *instructionFilter.Running {
				withExecution = true
				*match = append(*match, bson.M{"instruction_execution.instruction": bson.M{"$in": instructionFilter.Include}})
				continue
			}

			filters, err := q.getInstructionsFilters(
				ctx,
				bson.M{"_id": bson.M{"$in": instructionFilter.Include}},
			)
			if err != nil {
				return err
			}
			if len(filters) > 0 {
				if instructionFilter.Running == nil {
					withMatch = true
					*match = append(*match, bson.M{"$or": filters})
				} else {
					withMatch = true
					withExecution = true
					*match = append(*match, bson.M{"$and": []bson.M{
						{"instruction_execution.instruction": bson.M{"$nin": instructionFilter.Include}},
						{"$or": filters},
					}})
				}
			} else {
				*match = append(*match, bson.M{"$nor": []bson.M{{}}})
			}
			continue
		}

		if len(instructionFilter.IncludeTypes) > 0 {
			if instructionFilter.Running != nil && *instructionFilter.Running {
				withExecutionType = true
				*match = append(*match, bson.M{"instruction_execution.type": bson.M{"$in": instructionFilter.IncludeTypes}})
				continue
			}

			filters, err := q.getInstructionsFilters(
				ctx,
				bson.M{"type": bson.M{"$in": instructionFilter.IncludeTypes}},
			)
			if err != nil {
				return err
			}
			if len(filters) > 0 {
				if instructionFilter.Running == nil {
					withMatch = true
					*match = append(*match, bson.M{"$or": filters})
				} else {
					withMatch = true
					withExecutionType = true
					*match = append(*match, bson.M{"$and": []bson.M{
						{"instruction_execution.type": bson.M{"$nin": instructionFilter.IncludeTypes}},
						{"$or": filters},
					}})
				}
			} else {
				*match = append(*match, bson.M{"$nor": []bson.M{{}}})
			}
			continue
		}
	}

	if withMatch {
		q.computedFieldsForAlarmMatch["v.duration"] = true
		q.computedFieldsForAlarmMatch["v.infos_array"] = true
		q.computedFields["v.infos_array"] = bson.M{"$objectToArray": "$v.infos"}
	}
	if withExecution || withExecutionType {
		q.lookupsOnlyForAdditionalMatch["instruction_execution"] = true
		q.lookups = append(q.lookups, lookupWithKey{key: "instruction_execution", pipeline: getInstructionExecutionLookup(withExecutionType)})
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

	for lookup := range q.lookupsOnlyForAdditionalMatch {
		q.excludeLookupsBeforeSort = append(q.excludeLookupsBeforeSort, lookup)
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

func (q *MongoQueryBuilder) resolveAliasesInQuery(query any) any {
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
	case reflect.String:
		return q.resolveAlias(val.Interface().(string))
	}

	return res
}

func (q *MongoQueryBuilder) resolveAlias(v string) string {
	if v == "" {
		return v
	}

	prefix := ""
	if v[0] == '$' {
		v = v[1:]
		prefix = "$"
	}

	for alias, field := range q.fieldsAliases {
		if alias == v {
			return prefix + field
		}
	}

	for expr, repl := range q.fieldsAliasesByRegex {
		r, err := regexp.Compile(expr)
		if err == nil {
			replace := r.ReplaceAllString(v, repl)
			if v != replace {
				return prefix + replace
			}
		}
	}

	return prefix + v
}

func (q *MongoQueryBuilder) handleOpened(opened int) {
	alarmMatch := make([]bson.M, 0)
	entityMatch := make([]bson.M, 0)
	q.addOpenedFilter(opened, &alarmMatch, &entityMatch)
	if len(alarmMatch) > 0 {
		q.alarmMatch = append(q.alarmMatch, bson.M{"$match": bson.M{"$and": alarmMatch}})
	}

	if len(entityMatch) > 0 {
		q.lookupsForAdditionalMatch["entity"] = true
		q.additionalMatch = append(q.additionalMatch, bson.M{"$match": bson.M{"$and": entityMatch}})
	}
}

func (q *MongoQueryBuilder) handleDependencies(withDependencies bool) {
	if withDependencies {
		q.lookups = append(q.lookups, lookupWithKey{key: "entity.impacts_counts", pipeline: getImpactsCountPipeline()})
	}
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

func getPbehaviorLookup(authorProvider author.Provider) []bson.M {
	pipeline := []bson.M{
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
	pipeline = append(pipeline, bson.M{"$addFields": bson.M{
		"pbehavior.last_comment": bson.M{
			"$cond": bson.M{
				"if":   "$pbehavior.last_comment._id",
				"then": "$pbehavior.last_comment",
				"else": "$$REMOVE",
			},
		},
	}})
	return pipeline
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
			"opened_children": bson.M{"$size": "$children"},
			"closed_children": bson.M{"$size": "$resolved_children"},
		}},
	}
}

func getInstructionExecutionLookup(withType bool) []bson.M {
	pipeline := []bson.M{
		{"$lookup": bson.M{
			"from": mongo.InstructionExecutionMongoCollection,
			"let":  bson.M{"alarm": "$_id"},
			"pipeline": []bson.M{
				{"$match": bson.M{"$and": []bson.M{
					{"$expr": bson.M{"$eq": bson.A{"$$alarm", "$alarm"}}},
					{"status": bson.M{"$in": bson.A{InstructionExecutionStatusRunning, InstructionExecutionStatusWaitResult}}},
				}}},
			},
			"as": "instruction_execution",
		}},
		{"$unwind": bson.M{"path": "$instruction_execution", "preserveNullAndEmptyArrays": true}},
	}
	if withType {
		pipeline = append(pipeline, []bson.M{
			{"$lookup": bson.M{
				"from":         mongo.InstructionMongoCollection,
				"localField":   "instruction_execution.instruction",
				"foreignField": "_id",
				"as":           "instruction_execution.type",
			}},
			{"$unwind": bson.M{"path": "$instruction_execution.type", "preserveNullAndEmptyArrays": true}},
			{"$addFields": bson.M{
				"instruction_execution.type": "$instruction_execution.type.type",
			}},
		}...)
	}
	return pipeline
}

func getComputedFields(now types.CpsTime, userID string) bson.M {
	computedFields := bson.M{
		"infos":        "$v.infos",
		"impact_state": bson.M{"$multiply": bson.A{"$v.state.val", "$entity.impact_level"}},
		"v.duration": bson.M{"$ifNull": bson.A{
			"$v.duration",
			bson.M{"$subtract": bson.A{
				bson.M{"$cond": bson.M{
					"if":   "$v.resolved",
					"then": "$v.resolved",
					"else": now,
				}},
				"$v.creation_date",
			}},
		}},
		"v.current_state_duration": bson.M{"$ifNull": bson.A{
			"$v.current_state_duration",
			bson.M{"$subtract": bson.A{
				bson.M{"$cond": bson.M{
					"if":   "$v.resolved",
					"then": "$v.resolved",
					"else": now,
				}},
				"$v.state.t",
			}},
		}},
		"v.active_duration": bson.M{"$ifNull": bson.A{
			"$v.active_duration",
			bson.M{"$subtract": bson.A{
				bson.M{"$cond": bson.M{
					"if":   "$v.resolved",
					"then": "$v.resolved",
					"else": now,
				}},
				bson.M{"$sum": bson.A{
					"$v.creation_date",
					"$v.inactive_duration",
				}},
			}},
		}},
	}

	if userID != "" {
		computedFields["bookmark"] = bson.M{
			"$cond": bson.M{
				"if": bson.M{
					"$and": bson.A{
						bson.M{"$isArray": "$bookmarks"},
						bson.M{"$in": bson.A{userID, "$bookmarks"}},
					},
				},
				"then": true,
				"else": false,
			},
		}
	}

	return computedFields
}

func getIsMetaAlarmField() bson.M {
	return bson.M{"$cond": bson.A{bson.M{"$not": bson.A{"$v.meta"}}, false, true}}
}

func getInstructionQuery(instruction Instruction) (bson.M, error) {
	alarmPatternQuery, err := pattern.AlarmPatternToMongoQuery("", instruction.AlarmPattern, instruction.OldAlarmPatterns)
	if err != nil {
		return nil, fmt.Errorf("invalid alarm pattern in instruction id=%q: %w", instruction.ID, err)
	}

	entityPatternQuery, err := pattern.EntityPatternToMongoQuery("entity", instruction.EntityPattern, instruction.OldEntityPatterns)
	if err != nil {
		return nil, fmt.Errorf("invalid entity pattern in instruction id=%q: %w", instruction.ID, err)
	}

	if len(alarmPatternQuery) == 0 && len(entityPatternQuery) == 0 {
		return nil, nil
	}

	var and []bson.M
	if len(alarmPatternQuery) > 0 {
		and = append(and, alarmPatternQuery)
	}

	if len(entityPatternQuery) > 0 {
		and = append(and, entityPatternQuery)
	}

	if len(instruction.ActiveOnPbh) > 0 {
		and = append(and, bson.M{"v.pbehavior_info.type": bson.M{"$in": instruction.ActiveOnPbh}})
	}

	if len(instruction.DisabledOnPbh) > 0 {
		and = append(and, bson.M{"v.pbehavior_info.type": bson.M{"$nin": instruction.DisabledOnPbh}})
	}

	return bson.M{"$and": and}, nil
}

func getImpactsCountPipeline() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "entity.services",
			"foreignField": "_id",
			"as":           "service_impacts",
			"pipeline":     []bson.M{{"$project": bson.M{"_id": 1}}},
		}},
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "entity._id",
			"foreignField": "services",
			"as":           "depends",
			"pipeline":     []bson.M{{"$project": bson.M{"_id": 1}}},
		}},
		{"$addFields": bson.M{
			"entity.depends_count": bson.M{"$size": "$depends"},
			"entity.impacts_count": bson.M{"$size": "$service_impacts"},
		}},
		{"$project": bson.M{
			"service_impacts": 0,
			"depends":         0,
		}},
	}
}

func getOnlyParentsSearchPipeline(
	match bson.M,
	metaAlarmLookupCollection string,
	metaAlarmLookupMatch bson.M,
	searchByEntity bool,
) []bson.M {
	var pipeline []bson.M
	if searchByEntity {
		pipeline = append(pipeline, getEntityLookup()...)
	}

	pipeline = append(pipeline, bson.M{"$match": match})
	if searchByEntity {
		pipeline = append(pipeline, bson.M{"$project": bson.M{"entity": 0}})
	}

	pipeline = append(pipeline, []bson.M{
		{"$unwind": bson.M{"path": "$v.parents", "preserveNullAndEmptyArrays": true}},
		{"$group": bson.M{
			"_id": "$v.parents",
			"alarms": bson.M{"$push": bson.M{"$cond": bson.M{
				"if":   "$v.parents",
				"then": "$$REMOVE",
				"else": "$$ROOT",
			}}},
		}},
		{"$graphLookup": bson.M{
			"from":                    metaAlarmLookupCollection,
			"startWith":               "$_id",
			"connectFromField":        "_id",
			"connectToField":          "d",
			"restrictSearchWithMatch": metaAlarmLookupMatch,
			"as":                      "meta_alarm",
			"maxDepth":                0,
		}},
		{"$unwind": bson.M{"path": "$meta_alarm", "preserveNullAndEmptyArrays": true}},
		{"$unwind": bson.M{"path": "$alarms", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"alarm": bson.M{"$cond": bson.M{
				"if":   "$meta_alarm",
				"then": "$meta_alarm",
				"else": "$alarms",
			}},
		}},
		{"$match": bson.M{"alarm": bson.M{"$ne": nil}}},
		{"$group": bson.M{
			"_id":   "$alarm._id",
			"alarm": bson.M{"$first": "$alarm"},
		}},
		{"$replaceRoot": bson.M{"newRoot": "$alarm"}},
	}...)

	return pipeline
}
