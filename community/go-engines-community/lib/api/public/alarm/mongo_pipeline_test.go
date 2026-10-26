package alarm

import (
	"errors"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	patterndb "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/db"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
	"github.com/kylelemons/godebug/pretty"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var (
	testFrom = datetime.NewCpsTime(1600000000)
	testTo   = datetime.NewCpsTime(1700000000)
)

type pipelineResult struct {
	BeforeLimit []bson.M
	AfterLimit  []bson.M
	Sort        bson.M
}

func TestListPipeline_GiveOpenResolution_ShouldCheckFullPipeline(t *testing.T) {
	// Spelled out in full, so a change inside any of the stage helpers fails
	// here instead of silently passing everywhere else.
	now := datetime.NewCpsTime()
	expected := pipelineResult{
		BeforeLimit: []bson.M{
			{"$match": bson.M{"healthcheck": bson.M{"$in": bson.A{nil, false}}}},
			{"$match": bson.M{"v.resolved": nil}},
		},
		AfterLimit: []bson.M{
			{"$addFields": bson.M{
				"raw_alarm":  "$$ROOT",
				"raw_entity": "$e",
			}},
			{"$project": bson.M{
				"raw_alarm.v.steps": 0,
			}},
			{"$lookup": bson.M{
				"from":         mongo.EntityCategoryMongoCollection,
				"localField":   "e.category",
				"foreignField": "_id",
				"as":           "e.category",
			}},
			{"$unwind": bson.M{"path": "$e.category", "preserveNullAndEmptyArrays": true}},
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
				"e.pbehavior_info": bson.M{"$cond": bson.M{
					"if": "$e.pbehavior_info.id",
					"then": bson.M{
						"timestamp": "$e.pbehavior_info.timestamp",
						"rule": bson.M{
							"_id":  "$e.pbehavior_info.id",
							"name": "$e.pbehavior_info.name",
						},
						"reason": bson.M{
							"_id":  "$e.pbehavior_info.reason",
							"name": "$e.pbehavior_info.reason_name",
						},
						"type": bson.M{
							"_id":            "$e.pbehavior_info.type",
							"name":           "$e.pbehavior_info.type_name",
							"canonical_type": "$e.pbehavior_info.canonical_type",
						},
					},
					"else": nil,
				}},
			}},
			{"$addFields": bson.M{
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
			}},
			{"$addFields": bson.M{
				"type": bson.M{"$cond": bson.A{
					bson.M{"$ne": bson.A{bson.M{"$ifNull": bson.A{"$v.meta", ""}}, ""}},
					TypeMeta,
					bson.M{"$cond": bson.A{
						bson.M{"$ne": bson.A{bson.M{"$ifNull": bson.A{"$v.parents", bson.A{}}}, bson.A{}}},
						TypeGrouped,
						TypeStandalone,
					}},
				}},
			}},
			{"$lookup": bson.M{
				"from":         mongo.MetaAlarmRulesMongoCollection,
				"localField":   "v.meta",
				"foreignField": "_id",
				"as":           "meta_alarm_rule",
			}},
			{"$unwind": bson.M{"path": "$meta_alarm_rule", "preserveNullAndEmptyArrays": true}},
			{"$replaceRoot": bson.M{"newRoot": bson.M{
				"$mergeObjects": bson.A{
					"$v",
					"$$ROOT",
				},
			}}},
			{"$addFields": bson.M{"entity": "$e"}},
			{"$project": bson.M{
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
				"entity.latest_pause": 0,
				"entity.pause.path":   0,
			}},
		},
		Sort: bson.M{"$sort": bson.D{
			{Key: "v.creation_date", Value: -1},
			{Key: "_id", Value: 1},
		}},
	}

	beforeLimit, afterLimit, sort, err := ListPipeline(ListRequest{Resolution: ResolutionOpen}, now)
	if err != nil {
		t.Fatalf("expected no error; got %v", err)
	}

	assertPipeline(t, pipelineResult{beforeLimit, afterLimit, sort}, expected)
}

func TestListPipeline_GiveResolvedResolution_ShouldCheckFullPipeline(t *testing.T) {
	// Spelled out in full, so a change inside any of the stage helpers fails
	// here instead of silently passing everywhere else.
	now := datetime.NewCpsTime()
	expected := pipelineResult{
		BeforeLimit: []bson.M{
			{"$match": bson.M{"healthcheck": bson.M{"$in": bson.A{nil, false}}}},
		},
		AfterLimit: []bson.M{
			{"$lookup": bson.M{
				"from":         mongo.EntityMongoCollection,
				"localField":   "d",
				"foreignField": "_id",
				"as":           "entity",
			}},
			{"$unwind": "$entity"},
			{"$addFields": bson.M{
				"raw_alarm":  "$$ROOT",
				"raw_entity": "$entity",
			}},
			{"$project": bson.M{
				"raw_alarm.v.steps": 0,
			}},
			{"$lookup": bson.M{
				"from":         mongo.EntityCategoryMongoCollection,
				"localField":   "entity.category",
				"foreignField": "_id",
				"as":           "entity.category",
			}},
			{"$unwind": bson.M{"path": "$entity.category", "preserveNullAndEmptyArrays": true}},
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
				"entity.pbehavior_info": bson.M{"$cond": bson.M{
					"if": "$entity.pbehavior_info.id",
					"then": bson.M{
						"timestamp": "$entity.pbehavior_info.timestamp",
						"rule": bson.M{
							"_id":  "$entity.pbehavior_info.id",
							"name": "$entity.pbehavior_info.name",
						},
						"reason": bson.M{
							"_id":  "$entity.pbehavior_info.reason",
							"name": "$entity.pbehavior_info.reason_name",
						},
						"type": bson.M{
							"_id":            "$entity.pbehavior_info.type",
							"name":           "$entity.pbehavior_info.type_name",
							"canonical_type": "$entity.pbehavior_info.canonical_type",
						},
					},
					"else": nil,
				}},
			}},
			{"$addFields": bson.M{
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
			}},
			{"$addFields": bson.M{
				"type": bson.M{"$cond": bson.A{
					bson.M{"$ne": bson.A{bson.M{"$ifNull": bson.A{"$v.meta", ""}}, ""}},
					TypeMeta,
					bson.M{"$cond": bson.A{
						bson.M{"$ne": bson.A{bson.M{"$ifNull": bson.A{"$v.parents", bson.A{}}}, bson.A{}}},
						TypeGrouped,
						TypeStandalone,
					}},
				}},
			}},
			{"$lookup": bson.M{
				"from":         mongo.MetaAlarmRulesMongoCollection,
				"localField":   "v.meta",
				"foreignField": "_id",
				"as":           "meta_alarm_rule",
			}},
			{"$unwind": bson.M{"path": "$meta_alarm_rule", "preserveNullAndEmptyArrays": true}},
			{"$replaceRoot": bson.M{"newRoot": bson.M{
				"$mergeObjects": bson.A{
					"$v",
					"$$ROOT",
				},
			}}},
			{"$project": bson.M{
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
				"entity.latest_pause": 0,
				"entity.pause.path":   0,
			}},
		},
		Sort: bson.M{"$sort": bson.D{
			{Key: "v.creation_date", Value: -1},
			{Key: "_id", Value: 1},
		}},
	}

	beforeLimit, afterLimit, sort, err := ListPipeline(ListRequest{Resolution: ResolutionResolved}, now)
	if err != nil {
		t.Fatalf("expected no error; got %v", err)
	}

	assertPipeline(t, pipelineResult{beforeLimit, afterLimit, sort}, expected)
}

func TestListPipeline(t *testing.T) {
	now := datetime.NewCpsTime()
	dataSet := map[string]struct {
		Request           ListRequest
		Expected          pipelineResult
		ExpectedErr       string
		ExpectedErrFields []string
	}{
		"given search by alarm fields": {
			Request: ListRequest{
				Resolution: ResolutionOpen,
				Search:     "test search",
				SearchBy:   []string{"_id", "display_name", "infos.rule1.info1"},
			},
			Expected: pipelineResult{
				BeforeLimit: []bson.M{
					defaultFilter(),
					match(bson.M{"v.resolved": nil}),
					match(bson.M{"$or": []bson.M{
						{"_id": regex("test search")},
						{"v.display_name": regex("test search")},
						{"v.infos.rule1.info1": regex("test search")},
					}}),
				},
				AfterLimit: openAfterLimit(now),
				Sort:       defaultSort(),
			},
		},
		"given no search_by should search default fields": {
			Request: ListRequest{
				Resolution: ResolutionOpen,
				Search:     "test search",
			},
			Expected: pipelineResult{
				BeforeLimit: []bson.M{
					defaultFilter(),
					match(bson.M{"v.resolved": nil}),
					match(bson.M{"$or": []bson.M{
						{"_id": regex("test search")},
						{"v.display_name": regex("test search")},
						{"v.connector": regex("test search")},
						{"v.connector_name": regex("test search")},
						{"v.component": regex("test search")},
						{"v.resource": regex("test search")},
					}}),
				},
				AfterLimit: openAfterLimit(now),
				Sort:       defaultSort(),
			},
		},
		"given sort by alarm fields": {
			Request: ListRequest{
				Resolution: ResolutionOpen,
				Sort:       "_id,display_name:asc,infos.rule1.info1:desc",
			},
			Expected: pipelineResult{
				BeforeLimit: []bson.M{
					defaultFilter(),
					match(bson.M{"v.resolved": nil}),
				},
				AfterLimit: openAfterLimit(now),
				Sort: sortBy(
					bson.E{Key: "_id", Value: 1},
					bson.E{Key: "v.display_name", Value: 1},
					bson.E{Key: "v.infos.rule1.info1", Value: -1},
				),
			},
		},
		"given time range": {
			Request: ListRequest{
				Resolution: ResolutionOpen,
				TimeField:  "last_event_date",
				From:       &testFrom,
				To:         &testTo,
			},
			Expected: pipelineResult{
				BeforeLimit: []bson.M{
					defaultFilter(),
					match(bson.M{
						"v.resolved":        nil,
						"v.last_event_date": bson.M{"$gte": &testFrom, "$lte": &testTo},
					}),
				},
				AfterLimit: openAfterLimit(now),
				Sort:       defaultSort(),
			},
		},
		"given tags": {
			Request: ListRequest{
				Resolution: ResolutionOpen,
				Tags:       []string{"tag1", "tag2"},
			},
			Expected: pipelineResult{
				BeforeLimit: []bson.M{
					defaultFilter(),
					match(bson.M{
						"v.resolved": nil,
						"tags":       bson.M{"$in": []string{"tag1", "tag2"}},
					}),
				},
				AfterLimit: openAfterLimit(now),
				Sort:       defaultSort(),
			},
		},
		"given single standalone type should merge conditions into filter": {
			Request: ListRequest{
				Resolution: ResolutionOpen,
				Types:      []string{TypeStandalone},
			},
			Expected: pipelineResult{
				BeforeLimit: []bson.M{
					defaultFilter(),
					match(bson.M{
						"v.resolved": nil,
						"v.meta":     bson.M{"$in": bson.A{nil, ""}},
						"v.parents":  bson.M{"$in": bson.A{nil, bson.A{}}},
					}),
				},
				AfterLimit: openAfterLimit(now),
				Sort:       defaultSort(),
			},
		},
		"given single grouped type should merge conditions into filter": {
			Request: ListRequest{
				Resolution: ResolutionOpen,
				Types:      []string{TypeGrouped},
			},
			Expected: pipelineResult{
				BeforeLimit: []bson.M{
					defaultFilter(),
					match(bson.M{
						"v.resolved": nil,
						"v.parents":  bson.M{"$nin": bson.A{nil, bson.A{}}},
					}),
				},
				AfterLimit: openAfterLimit(now),
				Sort:       defaultSort(),
			},
		},
		"given single meta type should merge conditions into filter": {
			Request: ListRequest{
				Resolution: ResolutionOpen,
				Types:      []string{TypeMeta},
			},
			Expected: pipelineResult{
				BeforeLimit: []bson.M{
					defaultFilter(),
					match(bson.M{
						"v.resolved": nil,
						"v.meta":     bson.M{"$nin": bson.A{nil, ""}},
					}),
				},
				AfterLimit: openAfterLimit(now),
				Sort:       defaultSort(),
			},
		},
		"given multiple types should match any of them": {
			Request: ListRequest{
				Resolution: ResolutionOpen,
				Types:      []string{TypeStandalone, TypeMeta},
			},
			Expected: pipelineResult{
				BeforeLimit: []bson.M{
					defaultFilter(),
					match(bson.M{
						"v.resolved": nil,
						"$or": []bson.M{
							{"v.meta": bson.M{"$in": bson.A{nil, ""}}, "v.parents": bson.M{"$in": bson.A{nil, bson.A{}}}},
							{"v.meta": bson.M{"$nin": bson.A{nil, ""}}},
						},
					}),
				},
				AfterLimit: openAfterLimit(now),
				Sort:       defaultSort(),
			},
		},
		"given open and sort by entity field should not lookup": {
			Request: ListRequest{
				Resolution: ResolutionOpen,
				Sort:       "entity.name",
			},
			Expected: pipelineResult{
				BeforeLimit: []bson.M{
					defaultFilter(),
					match(bson.M{"v.resolved": nil}),
				},
				AfterLimit: openAfterLimit(now),
				Sort:       sortBy(bson.E{Key: "e.name", Value: 1}, bson.E{Key: "_id", Value: 1}),
			},
		},
		"given resolved and sort by entity field should keep lookup before limit": {
			Request: ListRequest{
				Resolution: ResolutionResolved,
				Sort:       "entity.name",
			},
			Expected: pipelineResult{
				BeforeLimit: append(
					[]bson.M{defaultFilter()},
					entityLookup(ResolutionResolved)...,
				),
				AfterLimit: resolvedAfterLimit(now),
				Sort:       sortBy(bson.E{Key: "entity.name", Value: 1}, bson.E{Key: "_id", Value: 1}),
			},
		},
		"given resolved and entity filter should move lookup after limit": {
			Request: ListRequest{
				Resolution:  ResolutionResolved,
				EntityTypes: []string{"resource"},
			},
			Expected: pipelineResult{
				BeforeLimit: append(
					append([]bson.M{defaultFilter()}, entityLookup(ResolutionResolved)...),
					match(bson.M{"entity.type": bson.M{"$in": []string{"resource"}}}),
					entityRemoveProject(),
				),
				AfterLimit: resolvedAfterLimitWithLookup(now),
				Sort:       defaultSort(),
			},
		},
		"given resolved, entity filter and entity sort should lookup once": {
			Request: ListRequest{
				Resolution:       ResolutionResolved,
				EntityCategories: []string{"cat1"},
				Sort:             "entity._id:desc",
			},
			Expected: pipelineResult{
				BeforeLimit: append(
					append([]bson.M{defaultFilter()}, entityLookup(ResolutionResolved)...),
					match(bson.M{"entity.category": bson.M{"$in": []string{"cat1"}}}),
				),
				AfterLimit: resolvedAfterLimit(now),
				Sort:       sortBy(bson.E{Key: "entity._id", Value: -1}, bson.E{Key: "_id", Value: 1}),
			},
		},
		"given resolved, entity filter and entity search should lookup once": {
			Request: ListRequest{
				Resolution:      ResolutionResolved,
				EntityUpstreams: []string{"upstream1"},
				Search:          "srv",
				SearchBy:        []string{"entity.upstream"},
			},
			Expected: pipelineResult{
				BeforeLimit: append(
					append([]bson.M{defaultFilter()}, entityLookup(ResolutionResolved)...),
					match(bson.M{"entity.upstream": bson.M{"$in": []string{"upstream1"}}}),
					match(bson.M{"$or": []bson.M{{"entity.upstream": regex("srv")}}}),
					entityRemoveProject(),
				),
				AfterLimit: resolvedAfterLimitWithLookup(now),
				Sort:       defaultSort(),
			},
		},
		"given unknown type": {
			Request: ListRequest{
				Resolution: ResolutionOpen,
				Types:      []string{"unknown"},
			},
			ExpectedErr: "unknown type: unknown",
		},
		"given not searchable field": {
			Request: ListRequest{
				Resolution: ResolutionOpen,
				Search:     "test search",
				SearchBy:   []string{"display_name", "v.steps", "entity.unknown"},
			},
			ExpectedErrFields: []string{"search_by.1", "search_by.2"},
		},
		"given not sortable field": {
			Request: ListRequest{
				Resolution: ResolutionOpen,
				Sort:       "v.steps",
			},
			ExpectedErrFields: []string{"sort.0"},
		},
	}

	for name, data := range dataSet {
		t.Run(name, func(t *testing.T) {
			beforeLimit, afterLimit, sort, err := ListPipeline(data.Request, now)
			if assertErr(t, err, data.ExpectedErr, data.ExpectedErrFields) {
				return
			}

			assertPipeline(t, pipelineResult{beforeLimit, afterLimit, sort}, data.Expected)
		})
	}
}

func TestSearchPipeline(t *testing.T) {
	now := datetime.NewCpsTime()
	alarmPattern := pattern.Alarm{{{
		Field:     "v.display_name",
		Condition: pattern.NewStringCondition(pattern.ConditionEqual, "alarm1"),
	}}}
	entityPattern := pattern.Entity{{{
		Field:     "name",
		Condition: pattern.NewStringCondition(pattern.ConditionEqual, "entity1"),
	}}}

	dataSet := map[string]struct {
		Request           SearchRequest
		Expected          pipelineResult
		ExpectedErr       string
		ExpectedErrFields []string
	}{
		"given no patterns": {
			Request: SearchRequest{Resolution: ResolutionOpen},
			Expected: pipelineResult{
				BeforeLimit: []bson.M{
					defaultFilter(),
					match(bson.M{"v.resolved": nil}),
				},
				AfterLimit: openAfterLimit(now),
				Sort:       defaultSort(),
			},
		},
		"given alarm pattern": {
			Request: SearchRequest{
				Resolution:   ResolutionOpen,
				AlarmPattern: alarmPattern,
			},
			Expected: pipelineResult{
				BeforeLimit: []bson.M{
					defaultFilter(),
					match(bson.M{"v.resolved": nil}),
					match(mustAlarmQuery(t, alarmPattern)),
				},
				AfterLimit: openAfterLimit(now),
				Sort:       defaultSort(),
			},
		},
		"given open and entity pattern should move lookup after limit": {
			Request: SearchRequest{
				Resolution:    ResolutionOpen,
				EntityPattern: entityPattern,
			},
			Expected: pipelineResult{
				BeforeLimit: append(
					append([]bson.M{defaultFilter(), match(bson.M{"v.resolved": nil})}, entityLookup(ResolutionOpen)...),
					match(mustEntityQuery(t, entityPattern, "e")),
					entityRemoveProject(),
				),
				AfterLimit: openAfterLimit(now),
				Sort:       defaultSort(),
			},
		},
		"given resolved and entity pattern should move lookup after limit": {
			Request: SearchRequest{
				Resolution:    ResolutionResolved,
				EntityPattern: entityPattern,
			},
			Expected: pipelineResult{
				BeforeLimit: append(
					append([]bson.M{defaultFilter()}, entityLookup(ResolutionResolved)...),
					match(mustEntityQuery(t, entityPattern, "entity")),
					entityRemoveProject(),
				),
				AfterLimit: resolvedAfterLimitWithLookup(now),
				Sort:       defaultSort(),
			},
		},
		"given resolved, entity pattern and entity sort should keep lookup before limit": {
			Request: SearchRequest{
				Resolution:    ResolutionResolved,
				EntityPattern: entityPattern,
				Sort:          "entity.name:desc",
			},
			Expected: pipelineResult{
				BeforeLimit: append(
					append([]bson.M{defaultFilter()}, entityLookup(ResolutionResolved)...),
					match(mustEntityQuery(t, entityPattern, "entity")),
				),
				AfterLimit: resolvedAfterLimit(now),
				Sort:       sortBy(bson.E{Key: "entity.name", Value: -1}, bson.E{Key: "_id", Value: 1}),
			},
		},
		"given multiple types should match any of them": {
			Request: SearchRequest{
				Resolution: ResolutionOpen,
				Types:      []string{TypeMeta, TypeGrouped},
			},
			Expected: pipelineResult{
				BeforeLimit: []bson.M{
					defaultFilter(),
					match(bson.M{
						"v.resolved": nil,
						"$or": []bson.M{
							{"v.meta": bson.M{"$nin": bson.A{nil, ""}}},
							{"v.parents": bson.M{"$nin": bson.A{nil, bson.A{}}}},
						},
					}),
				},
				AfterLimit: openAfterLimit(now),
				Sort:       defaultSort(),
			},
		},
		"given unknown type": {
			Request: SearchRequest{
				Resolution: ResolutionOpen,
				Types:      []string{"unknown"},
			},
			ExpectedErr: "unknown type: unknown",
		},
		"given not sortable field": {
			Request: SearchRequest{
				Resolution: ResolutionOpen,
				Sort:       "v.steps",
			},
			ExpectedErrFields: []string{"sort.0"},
		},
	}

	for name, data := range dataSet {
		t.Run(name, func(t *testing.T) {
			beforeLimit, afterLimit, sort, err := SearchPipeline(data.Request, now)
			if assertErr(t, err, data.ExpectedErr, data.ExpectedErrFields) {
				return
			}

			assertPipeline(t, pipelineResult{beforeLimit, afterLimit, sort}, data.Expected)
		})
	}
}

func match(q bson.M) bson.M {
	return bson.M{"$match": q}
}

func regex(v string) bson.Regex {
	return bson.Regex{Pattern: v, Options: "i"}
}

func sortBy(fields ...bson.E) bson.M {
	return bson.M{"$sort": bson.D(fields)}
}

func defaultSort() bson.M {
	return sortBy(bson.E{Key: "v.creation_date", Value: -1}, bson.E{Key: "_id", Value: 1})
}

func openAfterLimit(now datetime.CpsTime) []bson.M {
	return append(
		append(
			append(
				append(addRawFields("e"), nestedObjectsPipeline("e.")...),
				durationAddFields(ResolutionOpen, now), typeAddField(),
			),
			metaAlarmRuleLookup()...,
		),
		replaceRoot(), entityFieldAddField(), unusedFieldsProject(),
	)
}

func resolvedAfterLimit(now datetime.CpsTime) []bson.M {
	return append(
		append(
			append(
				append(addRawFields("entity"), nestedObjectsPipeline("entity.")...),
				durationAddFields(ResolutionResolved, now), typeAddField(),
			),
			metaAlarmRuleLookup()...,
		),
		replaceRoot(), unusedFieldsProject(),
	)
}

func resolvedAfterLimitWithLookup(now datetime.CpsTime) []bson.M {
	return append(entityLookup(ResolutionResolved), resolvedAfterLimit(now)...)
}

func assertPipeline(t *testing.T, got, expected pipelineResult) {
	t.Helper()

	if diff := pretty.Compare(expected, got); diff != "" {
		t.Errorf("unexpected pipeline (-want +got)\n%s", diff)
	}
}

func assertErr(t *testing.T, err error, expectedErr string, expectedErrFields []string) bool {
	t.Helper()

	switch {
	case expectedErr != "":
		if err == nil || err.Error() != expectedErr {
			t.Fatalf("expected error %q; got %v", expectedErr, err)
		}

		return true
	case len(expectedErrFields) > 0:
		var valErrs validator.ValidationErrors
		if !errors.As(err, &valErrs) {
			t.Fatalf("expected validation errors; got %v", err)
		}

		fields := make([]string, len(valErrs))
		for i, valErr := range valErrs {
			fields[i] = valErr.Namespace()
		}

		if diff := pretty.Compare(expectedErrFields, fields); diff != "" {
			t.Fatalf("unexpected error fields (-expected +got)\n%s", diff)
		}

		return true
	default:
		if err != nil {
			t.Fatalf("expected no error; got %v", err)
		}

		return false
	}
}

func mustAlarmQuery(t *testing.T, p pattern.Alarm) bson.M {
	t.Helper()

	q, err := patterndb.AlarmPatternToMongoQuery(p, "")
	if err != nil {
		t.Fatalf("cannot convert alarm pattern: %v", err)
	}

	return q
}

func mustEntityQuery(t *testing.T, p pattern.Entity, prefix string) bson.M {
	t.Helper()

	q, err := patterndb.EntityPatternToMongoQuery(p, prefix)
	if err != nil {
		t.Fatalf("cannot convert entity pattern: %v", err)
	}

	return q
}
