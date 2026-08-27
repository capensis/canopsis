package alarm

import (
	"cmp"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/mongoquery"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/sortparam"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	patterndb "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/db"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const DefaultSort = "creation_date:desc"

const valuePrefix = "v."

var sortableFields = []string{
	"_id",
	"display_name",
	"connector",
	"connector_name",
	"component",
	"resource",
	"creation_date",
	"resolved",
	"activation_date",
	"last_event_date",
	"last_update_date",
	"last_st_upd_dt",
	"duration",
	"current_state_duration",
	"active_duration",
}

var entitySortableFields = []string{
	"entity._id",
	"entity.name",
}

var defaultSearchFields = []string{
	"_id",
	"v.display_name",
	"v.connector",
	"v.connector_name",
	"v.component",
	"v.resource",
}

var searchableFields = []string{
	"_id",
	"display_name",
	"connector",
	"connector_name",
	"component",
	"resource",
	"output",
	"ticket.ticket",
}

var entitySearchableFields = []string{
	"entity.upstream",
}

var infoValueFieldsForErrParam = []string{
	"infos.*.*",
}

var entityInfoValueFieldsForErrParam = []string{
	"entity.infos.*.value",
	"entity.component_infos.*.value",
}

var stepSortableFields = []string{
	"t",
	"a",
	"user_id",
	"m",
}

func ListPipeline(r ListRequest, now datetime.CpsTime) (beforeLimit []bson.M, afterLimit []bson.M, sort bson.M, err error) {
	prefix := EntityPrefix(r.Resolution)

	sort, sortNeedsEntity, err := resolveSort(r.Sort, prefix)
	if err != nil {
		return nil, nil, nil, err
	}

	filter, err := listFilter(r)
	if err != nil {
		return nil, nil, nil, err
	}

	search, searchNeedsEntity, err := searchQuery(r.Search, r.SearchBy, prefix)
	if err != nil {
		return nil, nil, nil, err
	}

	alarmMatches := []bson.M{filter}
	entityMatches := []bson.M{listEntityFilter(r, prefix)}
	if searchNeedsEntity {
		entityMatches = append(entityMatches, search)
	} else {
		alarmMatches = append(alarmMatches, search)
	}

	beforeLimit, afterLimit = buildPipeline(r.Resolution, alarmMatches, entityMatches, sortNeedsEntity, now)

	return beforeLimit, afterLimit, sort, nil
}

func SearchPipeline(r SearchRequest, now datetime.CpsTime) (beforeLimit []bson.M, afterLimit []bson.M, sort bson.M, err error) {
	prefix := EntityPrefix(r.Resolution)

	sort, sortNeedsEntity, err := resolveSort(r.Sort, prefix)
	if err != nil {
		return nil, nil, nil, err
	}

	filter, err := searchFilter(r)
	if err != nil {
		return nil, nil, nil, err
	}

	alarmMatches := []bson.M{filter}

	if len(r.AlarmPattern) > 0 {
		q, err := patterndb.AlarmPatternToMongoQuery(r.AlarmPattern, "")
		if err != nil {
			return nil, nil, nil, validation.NewSingleError("alarm_pattern", "AlarmPattern", "AlarmPattern", r)
		}

		alarmMatches = append(alarmMatches, q)
	}

	if len(r.PbehaviorPattern) > 0 {
		q, err := patterndb.PbehaviorInfoPatternToMongoQuery(r.PbehaviorPattern, "v")
		if err != nil {
			return nil, nil, nil, validation.NewSingleError("pbehavior_pattern", "PbehaviorPattern", "PbehaviorPattern", r)
		}

		alarmMatches = append(alarmMatches, q)
	}

	var entityMatches []bson.M

	if len(r.EntityPattern) > 0 {
		q, err := patterndb.EntityPatternToMongoQuery(r.EntityPattern, strings.TrimSuffix(prefix, "."))
		if err != nil {
			return nil, nil, nil, validation.NewSingleError("entity_pattern", "EntityPattern", "EntityPattern", r)
		}

		entityMatches = append(entityMatches, q)
	}

	beforeLimit, afterLimit = buildPipeline(r.Resolution, alarmMatches, entityMatches, sortNeedsEntity, now)

	return beforeLimit, afterLimit, sort, nil
}

// buildPipeline lays out the stages: alarm-only matches first, then the entity
// join followed by the matches that need it. The join stays before the limit
// only if sorting needs it, otherwise it is projected away and redone after the
// limit, so fewer documents are joined.
func buildPipeline(resolution string, alarmMatches, entityMatches []bson.M, sortNeedsEntity bool, now datetime.CpsTime) (beforeLimit, afterLimit []bson.M) {
	lookup := entityLookup(resolution)

	beforeLimit = []bson.M{defaultFilter()}
	for _, q := range alarmMatches {
		if len(q) > 0 {
			beforeLimit = append(beforeLimit, bson.M{"$match": q})
		}
	}

	joined := false
	for _, q := range entityMatches {
		if len(q) == 0 {
			continue
		}

		if !joined {
			beforeLimit = append(beforeLimit, lookup...)
			joined = true
		}

		beforeLimit = append(beforeLimit, bson.M{"$match": q})
	}

	switch {
	case sortNeedsEntity && !joined:
		beforeLimit = append(beforeLimit, lookup...)
		joined = true
	case !sortNeedsEntity && joined:
		beforeLimit = append(beforeLimit, entityRemoveProject())
		afterLimit = append(afterLimit, lookup...)
	}

	if !joined {
		afterLimit = append(afterLimit, lookup...)
	}

	afterLimit = append(afterLimit, nestedObjectsPipeline(EntityPrefix(resolution))...)
	afterLimit = append(afterLimit, durationAddFields(resolution, now))
	afterLimit = append(afterLimit, typeAddField())
	afterLimit = append(afterLimit, metaAlarmRuleLookup()...)
	afterLimit = append(afterLimit, replaceRoot())
	if resolution == ResolutionOpen {
		afterLimit = append(afterLimit, entityFieldAddField())
	}

	afterLimit = append(afterLimit, unusedFieldsProject())

	return beforeLimit, afterLimit
}

// EntityPrefix returns the path the entity is available at:
// resolved alarms are joined with the entity collection, open alarms embed it.
func EntityPrefix(resolution string) string {
	if resolution == ResolutionResolved {
		return "entity."
	}

	return "e."
}

// entityLookup returns the stages joining the entity collection, or nil if the entity is already embedded in the alarm.
func entityLookup(resolution string) []bson.M {
	if resolution != ResolutionResolved {
		return nil
	}

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

func metaAlarmRuleLookup() []bson.M {
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

func resolveSort(s, prefix string) (sort bson.M, needsEntity bool, err error) {
	sort, err = sortparam.ValidateResolveFunc(cmp.Or(s, DefaultSort), func(f string) (string, bool) {
		rf, isEntityField, isSortable := ResolveSortableField(f, prefix)
		needsEntity = needsEntity || isEntityField

		return rf, isSortable
	}, SortErrParam())

	return sort, needsEntity, err
}

func listFilter(r ListRequest) (bson.M, error) {
	filter := bson.M{}
	if r.Resolution == ResolutionOpen {
		filter["v.resolved"] = nil
	}

	if err := applyTypeFilter(filter, r.Types); err != nil {
		return nil, err
	}

	if r.From != nil || r.To != nil {
		q := bson.M{}
		if r.From != nil {
			q["$gte"] = r.From
		}

		if r.To != nil {
			q["$lte"] = r.To
		}

		filter[valuePrefix+r.TimeField] = q
	}

	if len(r.MetaAlarms) > 0 {
		filter["v.parents"] = bson.M{"$in": r.MetaAlarms}
	}

	if len(r.Tags) > 0 {
		filter["tags"] = bson.M{"$in": r.Tags}
	}

	return filter, nil
}

func listEntityFilter(r ListRequest, prefix string) bson.M {
	filter := bson.M{}
	if len(r.EntityTypes) > 0 {
		filter[prefix+"type"] = bson.M{"$in": r.EntityTypes}
	}

	if len(r.EntityCategories) > 0 {
		filter[prefix+"category"] = bson.M{"$in": r.EntityCategories}
	}

	if len(r.EntityConnectors) > 0 {
		filter[prefix+"connector"] = bson.M{"$in": r.EntityConnectors}
	}

	if len(r.EntityComponents) > 0 {
		filter[prefix+"component"] = bson.M{"$in": r.EntityComponents}
	}

	if len(r.EntityUpstreams) > 0 {
		filter[prefix+"upstream"] = bson.M{"$in": r.EntityUpstreams}
	}

	return filter
}

func searchFilter(r SearchRequest) (bson.M, error) {
	filter := bson.M{}
	if r.Resolution == ResolutionOpen {
		filter["v.resolved"] = nil
	}

	if err := applyTypeFilter(filter, r.Types); err != nil {
		return nil, err
	}

	if len(r.MetaAlarms) > 0 {
		filter["v.parents"] = bson.M{"$in": r.MetaAlarms}
	}

	return filter, nil
}

// applyTypeFilter matches the alarms of any of the requested types.
func applyTypeFilter(filter bson.M, types []string) error {
	queries := make([]bson.M, 0, len(types))
	for _, t := range types {
		switch t {
		case TypeStandalone:
			queries = append(queries, bson.M{
				"v.meta":    bson.M{"$in": bson.A{nil, ""}},
				"v.parents": bson.M{"$in": bson.A{nil, bson.A{}}},
			})
		case TypeMeta:
			queries = append(queries, bson.M{
				"v.meta": bson.M{"$nin": bson.A{nil, ""}},
			})
		case TypeGrouped:
			queries = append(queries, bson.M{
				"v.parents": bson.M{"$nin": bson.A{nil, bson.A{}}},
			})
		default:
			return fmt.Errorf("unknown type: %s", t)
		}
	}

	if len(queries) == 0 {
		return nil
	}

	if len(queries) == 1 {
		maps.Copy(filter, queries[0])

		return nil
	}

	// match all alarms
	if len(queries) == 3 {
		return nil
	}

	filter["$or"] = queries

	return nil
}

func typeAddField() bson.M {
	return bson.M{
		"$addFields": bson.M{
			"type": bson.M{"$cond": bson.A{
				bson.M{"$ne": bson.A{bson.M{"$ifNull": bson.A{"$v.meta", ""}}, ""}},
				TypeMeta,
				bson.M{"$cond": bson.A{
					bson.M{"$ne": bson.A{bson.M{"$ifNull": bson.A{"$v.parents", bson.A{}}}, bson.A{}}},
					TypeGrouped,
					TypeStandalone,
				}},
			}},
		},
	}
}

// searchQuery builds the query for the search parameter. It returns
// needsEntity if any of the searched fields belongs to the entity.
func searchQuery(search string, searchBy []string, prefix string) (q bson.M, needsEntity bool, err error) {
	if search == "" {
		return nil, false, nil
	}

	fields := make([]string, len(searchBy))
	errs := validator.ValidationErrors{}
	errParam := searchErrParam()
	for i, f := range searchBy {
		var isEntityField, isSearchable bool
		fields[i], isEntityField, isSearchable = resolveSearchField(f, prefix)
		if !isSearchable {
			errs = append(errs, validation.NewFieldErrorWithParam("oneof", strconv.Itoa(i), "search_by."+strconv.Itoa(i), errParam))
		}

		needsEntity = needsEntity || isEntityField
	}

	if len(errs) > 0 {
		return nil, false, validation.NewError(errs, nil)
	}

	if len(fields) == 0 {
		fields = defaultSearchFields
	}

	return mongoquery.GetSearchQuery(search, fields), needsEntity, nil
}

func ResolveSortableField(f, entityPrefix string) (resolvedField string, isEntityField bool, isSearchable bool) {
	if slices.Contains(sortableFields, f) || isInfoValueField(f) {
		if f == "_id" {
			return f, false, true
		}

		return valuePrefix + f, false, true
	}

	if slices.Contains(entitySortableFields, f) {
		if entityPrefix != "entity." {
			return entityPrefix + strings.TrimPrefix(f, "entity."), true, true
		}

		return f, true, true
	}

	rf, ok := resolveEntityInfoValueField(f, entityPrefix)
	if ok {
		return rf, true, true
	}

	return "", false, false
}

func SortErrParam() string {
	return strings.Join(sortableFields, " ") + " " +
		strings.Join(infoValueFieldsForErrParam, " ") + " " +
		strings.Join(entitySortableFields, " ") + " " +
		strings.Join(entityInfoValueFieldsForErrParam, " ")
}

func resolveSearchField(f, entityPrefix string) (resolvedField string, isEntityField bool, isSearchable bool) {
	if slices.Contains(searchableFields, f) || isInfoValueField(f) {
		if f == "_id" {
			return f, false, true
		}

		return valuePrefix + f, false, true
	}

	if slices.Contains(entitySearchableFields, f) {
		if entityPrefix != "entity." {
			return entityPrefix + strings.TrimPrefix(f, "entity."), true, true
		}

		return f, true, true
	}

	rf, ok := resolveEntityInfoValueField(f, entityPrefix)
	if ok {
		return rf, true, true
	}

	return "", false, false
}

func searchErrParam() string {
	return strings.Join(searchableFields, " ") + " " +
		strings.Join(infoValueFieldsForErrParam, " ") + " " +
		strings.Join(entitySearchableFields, " ") + " " +
		strings.Join(entityInfoValueFieldsForErrParam, " ")
}

func isInfoValueField(f string) bool {
	parts := strings.Split(f, ".")
	if len(parts) != 3 {
		return false
	}

	return parts[0] == "infos"
}

func resolveEntityInfoValueField(f, entityPrefix string) (string, bool) {
	parts := strings.Split(f, ".")
	if len(parts) != 4 {
		return "", false
	}

	if parts[0] == "entity" && (parts[1] == "infos" || parts[1] == "component_infos") && parts[3] == "value" {
		return entityPrefix + strings.Join(parts[1:], "."), true
	}

	return "", false
}

func nestedObjectsPipeline(entityPrefix string) []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.EntityCategoryMongoCollection,
			"localField":   entityPrefix + "category",
			"foreignField": "_id",
			"as":           entityPrefix + "category",
		}},
		{"$unwind": bson.M{"path": "$" + entityPrefix + "category", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"pbehavior_info": bson.M{"$cond": bson.M{
				"if": "$v.pbehavior_info.id",
				"then": bson.M{
					"timestamp": "$v.pbehavior_info.timestamp",
					"rule": bson.M{
						"_id":  "$v.pbehavior_info.id",
						"name": "$v.pbehavior_info.name",
					},
					"reason": bson.M{
						"_id":  "$v.pbehavior_info.reason",
						"name": "$v.pbehavior_info.reason_name",
					},
					"type": bson.M{
						"_id":            "$v.pbehavior_info.type",
						"name":           "$v.pbehavior_info.type_name",
						"canonical_type": "$v.pbehavior_info.canonical_type",
					},
				},
				"else": nil,
			}},
			entityPrefix + "pbehavior_info": bson.M{"$cond": bson.M{
				"if": "$" + entityPrefix + "pbehavior_info.id",
				"then": bson.M{
					"timestamp": "$" + entityPrefix + "pbehavior_info.timestamp",
					"rule": bson.M{
						"_id":  "$" + entityPrefix + "pbehavior_info.id",
						"name": "$" + entityPrefix + "pbehavior_info.name",
					},
					"reason": bson.M{
						"_id":  "$" + entityPrefix + "pbehavior_info.reason",
						"name": "$" + entityPrefix + "pbehavior_info.reason_name",
					},
					"type": bson.M{
						"_id":            "$" + entityPrefix + "pbehavior_info.type",
						"name":           "$" + entityPrefix + "pbehavior_info.type_name",
						"canonical_type": "$" + entityPrefix + "pbehavior_info.canonical_type",
					},
				},
				"else": nil,
			}},
		}},
	}
}

func defaultFilter() bson.M {
	return bson.M{"$match": bson.M{
		"healthcheck": bson.M{"$in": bson.A{nil, false}},
	}}
}

func replaceRoot() bson.M {
	return bson.M{"$replaceRoot": bson.M{"newRoot": bson.M{
		"$mergeObjects": bson.A{
			"$v",
			"$$ROOT",
		},
	}}}
}

func entityFieldAddField() bson.M {
	return bson.M{"$addFields": bson.M{
		"entity": "$e",
	}}
}

func entityRemoveProject() bson.M {
	return bson.M{
		"$project": bson.M{
			"entity": 0,
		},
	}
}

func unusedFieldsProject() bson.M {
	return bson.M{
		"$project": bson.M{
			"v":                   0,
			"e":                   0,
			"bookmarks":           0,
			"parents":             0,
			"children":            0,
			"unlinked_parents":    0,
			"long_output_history": 0,
			"entity.services":     0,
			"entity.perf_data":    0,
			"entity.comments":     0,
			"entity.imtags":       0,
		},
	}
}

func durationAddFields(resolution string, now datetime.CpsTime) bson.M {
	if resolution == ResolutionResolved {
		return bson.M{
			"$addFields": bson.M{
				"v.duration": bson.M{"$ifNull": bson.A{
					"$v.duration",
					bson.M{"$subtract": bson.A{
						"$v.resolved",
						"$v.creation_date",
					}},
				}},
				"v.current_state_duration": bson.M{"$ifNull": bson.A{
					"$v.current_state_duration",
					bson.M{"$subtract": bson.A{
						"$v.resolved",
						"$v.state.t",
					}},
				}},
				"v.active_duration": bson.M{"$ifNull": bson.A{
					"$v.active_duration",
					bson.M{"$subtract": bson.A{
						"$v.resolved",
						bson.M{"$sum": bson.A{
							"$v.creation_date",
							"$v.inactive_duration",
						}},
					}},
				}},
			},
		}
	}

	return bson.M{
		"$addFields": bson.M{
			"v.duration": bson.M{"$subtract": bson.A{
				now,
				"$v.creation_date",
			}},
			"v.current_state_duration": bson.M{"$subtract": bson.A{
				now,
				"$v.state.t",
			}},
			"v.active_duration": bson.M{"$subtract": bson.A{
				now,
				bson.M{"$sum": bson.A{
					"$v.creation_date",
					"$v.inactive_duration",
					bson.M{"$cond": bson.M{
						"if":   "$v.inactive_start",
						"then": bson.M{"$subtract": bson.A{now, "$v.inactive_start"}},
						"else": 0,
					}},
				}},
			}},
		},
	}
}
