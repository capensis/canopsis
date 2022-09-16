package alarm

import (
	"context"
	"fmt"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	"github.com/golang/mock/gomock"
	"github.com/kylelemons/godebug/pretty"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMongoQueryBuilder_CreateListAggregationPipeline_GivenPaginationRequest_ShouldBuildQueryWithLookupsAfterLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockDbClient := createMockDbClient(ctrl)
	request := ListRequestWithPagination{
		Query: pagination.Query{
			Page:     2,
			Limit:    10,
			Paginate: true,
		},
	}
	now := types.NewCpsTime()
	expectedDataPipeline := []bson.M{
		{"$sort": bson.D{{Key: "t", Value: -1}, {Key: "_id", Value: 1}}},
		{"$skip": 10},
		{"$limit": 10},
	}
	expectedDataPipeline = append(expectedDataPipeline, getEntityLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getEntityCategoryLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorInfoTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$addFields": getComputedFields(now),
	})
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$project": bson.M{
			"v.steps":             0,
			"pbehavior.comments":  0,
			"pbehavior_info_type": 0,
		},
	})
	expected := []bson.M{
		{"$match": bson.M{"$and": []bson.M{{"v.meta": nil}}}},
	}
	expected = append(expected, getEntityLookup()...)
	expected = append(expected, []bson.M{
		{"$match": bson.M{"entity.enabled": true}},
		{"$project": bson.M{"entity": 0}},
		{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	}...)

	b := NewMongoQueryBuilder(mockDbClient)
	result, err := b.CreateListAggregationPipeline(ctx, request, now)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("unexpected result: %s", diff)
	}
}

func TestMongoQueryBuilder_CreateListAggregationPipeline_GivenRequestWithWidgetFilterWithAlarmAndPbehaviorInfoPatterns_ShouldBuildQueryWithLookupsAfterLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filterId := "test-filter"
	filter := view.WidgetFilter{
		AlarmPatternFields: savedpattern.AlarmPatternFields{
			AlarmPattern: pattern.Alarm{
				{
					{
						Field:     "v.resource",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource"),
					},
				},
			},
		},
		PbehaviorPatternFields: savedpattern.PbehaviorPatternFields{
			PbehaviorPattern: pattern.PbehaviorInfo{
				{
					{
						Field:     "pbehavior_info.canonical_type",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "pause"),
					},
				},
			},
		},
	}
	mockDbClient := createMockDbClientWithFilterFetching(ctrl, filterId, filter)
	request := ListRequestWithPagination{
		Query: pagination.GetDefaultQuery(),
		ListRequest: ListRequest{
			FilterRequest: FilterRequest{
				BaseFilterRequest: BaseFilterRequest{
					Filter: filterId,
				},
			},
		},
	}
	now := types.NewCpsTime()
	expectedDataPipeline := []bson.M{
		{"$sort": bson.D{{Key: "t", Value: -1}, {Key: "_id", Value: 1}}},
		{"$skip": 0},
		{"$limit": 10},
	}
	expectedDataPipeline = append(expectedDataPipeline, getEntityLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getEntityCategoryLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorInfoTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$addFields": getComputedFields(now),
	})
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$project": bson.M{
			"v.steps":             0,
			"pbehavior.comments":  0,
			"pbehavior_info_type": 0,
		},
	})
	expected := []bson.M{
		{"$match": bson.M{"$or": []bson.M{{"$and": []bson.M{
			{"v.resource": bson.M{"$eq": "test-resource"}},
		}}}}},
		{"$match": bson.M{"$or": []bson.M{{"$and": []bson.M{
			{"v.pbehavior_info.canonical_type": bson.M{"$eq": "pause"}},
		}}}}},
		{"$match": bson.M{"$and": []bson.M{{"v.meta": nil}}}},
	}
	expected = append(expected, getEntityLookup()...)
	expected = append(expected, []bson.M{
		{"$match": bson.M{"entity.enabled": true}},
		{"$project": bson.M{"entity": 0}},
		{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	}...)

	b := NewMongoQueryBuilder(mockDbClient)
	result, err := b.CreateListAggregationPipeline(ctx, request, now)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("unexpected result: %s", diff)
	}
}

func TestMongoQueryBuilder_CreateListAggregationPipeline_GivenRequestWithWidgetFilterWithAlarmPatternWithDurationField_ShouldBuildQueryWithAddFieldsBeforeLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filterId := "test-filter"
	durationCond, err := pattern.NewDurationCondition(pattern.ConditionGT, types.DurationWithUnit{
		Value: 10,
		Unit:  "m",
	})
	if err != nil {
		panic(err)
	}
	filter := view.WidgetFilter{
		AlarmPatternFields: savedpattern.AlarmPatternFields{
			AlarmPattern: pattern.Alarm{
				{
					{
						Field:     "v.resource",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource"),
					},
					{
						Field:     "v.duration",
						Condition: durationCond,
					},
				},
			},
		},
	}
	mockDbClient := createMockDbClientWithFilterFetching(ctrl, filterId, filter)
	request := ListRequestWithPagination{
		Query: pagination.GetDefaultQuery(),
		ListRequest: ListRequest{
			FilterRequest: FilterRequest{
				BaseFilterRequest: BaseFilterRequest{
					Filter: filterId,
				},
			},
		},
	}
	now := types.NewCpsTime()
	expectedDataPipeline := []bson.M{
		{"$sort": bson.D{{Key: "t", Value: -1}, {Key: "_id", Value: 1}}},
		{"$skip": 0},
		{"$limit": 10},
	}
	expectedDataPipeline = append(expectedDataPipeline, getEntityLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getEntityCategoryLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorInfoTypeLookup()...)
	fields := getComputedFields(now)
	durationField := fields["v.duration"]
	delete(fields, "v.duration")
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$addFields": fields,
	})
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$project": bson.M{
			"v.steps":             0,
			"pbehavior.comments":  0,
			"pbehavior_info_type": 0,
		},
	})
	expected := []bson.M{
		{"$addFields": bson.M{"v.duration": durationField}},
		{"$match": bson.M{"$or": []bson.M{{"$and": []bson.M{
			{"v.resource": bson.M{"$eq": "test-resource"}},
			{"v.duration": bson.M{"$gt": 600}},
		}}}}},
		{"$match": bson.M{"$and": []bson.M{{"v.meta": nil}}}},
	}
	expected = append(expected, getEntityLookup()...)
	expected = append(expected, []bson.M{
		{"$match": bson.M{"entity.enabled": true}},
		{"$project": bson.M{"entity": 0}},
		{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	}...)

	b := NewMongoQueryBuilder(mockDbClient)
	result, err := b.CreateListAggregationPipeline(ctx, request, now)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("unexpected result: %s", diff)
	}
}

func TestMongoQueryBuilder_CreateListAggregationPipeline_GivenRequestWithWidgetFilterWithAlarmPatternWithInfosField_ShouldBuildQueryWithInfosTransformBeforeLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filterId := "test-filter"
	filter := view.WidgetFilter{
		AlarmPatternFields: savedpattern.AlarmPatternFields{
			AlarmPattern: pattern.Alarm{
				{
					{
						Field:     "v.resource",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource"),
					},
					{
						Field:     "v.infos.info_name",
						FieldType: pattern.FieldTypeInt,
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 3),
					},
				},
			},
		},
	}
	mockDbClient := createMockDbClientWithFilterFetching(ctrl, filterId, filter)
	request := ListRequestWithPagination{
		Query: pagination.GetDefaultQuery(),
		ListRequest: ListRequest{
			FilterRequest: FilterRequest{
				BaseFilterRequest: BaseFilterRequest{
					Filter: filterId,
				},
			},
		},
	}
	now := types.NewCpsTime()
	expectedDataPipeline := []bson.M{
		{"$sort": bson.D{{Key: "t", Value: -1}, {Key: "_id", Value: 1}}},
		{"$skip": 0},
		{"$limit": 10},
	}
	expectedDataPipeline = append(expectedDataPipeline, getEntityLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getEntityCategoryLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorInfoTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$addFields": getComputedFields(now),
	})
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$project": bson.M{
			"v.steps":             0,
			"pbehavior.comments":  0,
			"pbehavior_info_type": 0,
		},
	})
	expected := []bson.M{
		{"$addFields": bson.M{"v.infos_array": bson.M{"$objectToArray": "$v.infos"}}},
		{"$match": bson.M{"$or": []bson.M{{"$and": []bson.M{
			{"v.resource": bson.M{"$eq": "test-resource"}},
			{"$and": []bson.M{
				{"v.infos_array.v.info_name": bson.M{"$type": bson.A{"long", "int", "decimal"}}},
				{"v.infos_array.v.info_name": bson.M{"$eq": 3}},
			}},
		}}}}},
		{"$match": bson.M{"$and": []bson.M{{"v.meta": nil}}}},
	}
	expected = append(expected, getEntityLookup()...)
	expected = append(expected, []bson.M{
		{"$match": bson.M{"entity.enabled": true}},
		{"$project": bson.M{"entity": 0}},
		{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	}...)

	b := NewMongoQueryBuilder(mockDbClient)
	result, err := b.CreateListAggregationPipeline(ctx, request, now)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("unexpected result: %s", diff)
	}
}

func TestMongoQueryBuilder_CreateListAggregationPipeline_GivenRequestWithWidgetFilterWithAllPatterns_ShouldBuildQueryWithEntityLookupsBeforeLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filterId := "test-filter"
	filter := view.WidgetFilter{
		AlarmPatternFields: savedpattern.AlarmPatternFields{
			AlarmPattern: pattern.Alarm{
				{
					{
						Field:     "v.resource",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource"),
					},
				},
			},
		},
		PbehaviorPatternFields: savedpattern.PbehaviorPatternFields{
			PbehaviorPattern: pattern.PbehaviorInfo{
				{
					{
						Field:     "pbehavior_info.canonical_type",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "pause"),
					},
				},
			},
		},
		EntityPatternFields: savedpattern.EntityPatternFields{
			EntityPattern: pattern.Entity{
				{
					{
						Field:     "category",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-category"),
					},
				},
			},
		},
	}
	mockDbClient := createMockDbClientWithFilterFetching(ctrl, filterId, filter)
	request := ListRequestWithPagination{
		Query: pagination.GetDefaultQuery(),
		ListRequest: ListRequest{
			FilterRequest: FilterRequest{
				BaseFilterRequest: BaseFilterRequest{
					Filter: filterId,
				},
			},
		},
	}
	now := types.NewCpsTime()
	expectedDataPipeline := []bson.M{
		{"$sort": bson.D{{Key: "t", Value: -1}, {Key: "_id", Value: 1}}},
		{"$skip": 0},
		{"$limit": 10},
	}
	expectedDataPipeline = append(expectedDataPipeline, getEntityLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getEntityCategoryLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorInfoTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$addFields": getComputedFields(now),
	})
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$project": bson.M{
			"v.steps":             0,
			"pbehavior.comments":  0,
			"pbehavior_info_type": 0,
		},
	})
	expected := []bson.M{
		{"$match": bson.M{"$or": []bson.M{{"$and": []bson.M{
			{"v.resource": bson.M{"$eq": "test-resource"}},
		}}}}},
		{"$match": bson.M{"$or": []bson.M{{"$and": []bson.M{
			{"v.pbehavior_info.canonical_type": bson.M{"$eq": "pause"}},
		}}}}},
		{"$match": bson.M{"$and": []bson.M{{"v.meta": nil}}}},
	}
	expected = append(expected, getEntityLookup()...)
	expected = append(expected,
		bson.M{"$match": bson.M{"entity.enabled": true}},
		bson.M{"$match": bson.M{"$or": []bson.M{{"$and": []bson.M{
			{"entity.category": bson.M{"$eq": "test-category"}},
		}}}}},
		bson.M{"$project": bson.M{"entity": 0}},
		bson.M{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		bson.M{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	)

	b := NewMongoQueryBuilder(mockDbClient)
	result, err := b.CreateListAggregationPipeline(ctx, request, now)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("unexpected result: %s", diff)
	}
}

func TestMongoQueryBuilder_CreateListAggregationPipeline_GivenRequestWithWidgetFilterWithOldQuery_ShouldBuildQueryWithLookupsBeforeLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filterId := "test-filter"
	filter := view.WidgetFilter{
		OldMongoQuery: `{"$and": [
			{"v.connector": "test-connector"}
		]}`,
	}
	mockDbClient := createMockDbClientWithFilterFetching(ctrl, filterId, filter)
	request := ListRequestWithPagination{
		Query: pagination.GetDefaultQuery(),
		ListRequest: ListRequest{
			FilterRequest: FilterRequest{
				BaseFilterRequest: BaseFilterRequest{
					Filter: filterId,
				},
			},
		},
	}
	now := types.NewCpsTime()
	expectedDataPipeline := []bson.M{
		{"$sort": bson.D{{Key: "t", Value: -1}, {Key: "_id", Value: 1}}},
		{"$skip": 0},
		{"$limit": 10},
	}
	expectedDataPipeline = append(expectedDataPipeline, getEntityLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getEntityCategoryLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorInfoTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$addFields": getComputedFields(now),
	})
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$project": bson.M{
			"v.steps":             0,
			"pbehavior.comments":  0,
			"pbehavior_info_type": 0,
		},
	})
	expected := []bson.M{
		{"$addFields": bson.M{"v.infos_array": bson.M{"$objectToArray": "$v.infos"}}},
		{"$match": bson.M{"$and": []bson.M{
			{"v.connector": "test-connector"},
		}}},
		{"$match": bson.M{"$and": []bson.M{{"v.meta": nil}}}},
	}
	expected = append(expected, getEntityLookup()...)
	expected = append(expected, []bson.M{
		{"$match": bson.M{"entity.enabled": true}},
		{"$project": bson.M{"entity": 0}},
		{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	}...)

	b := NewMongoQueryBuilder(mockDbClient)
	result, err := b.CreateListAggregationPipeline(ctx, request, now)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("unexpected result: %s", diff)
	}
}

func TestMongoQueryBuilder_CreateListAggregationPipeline_GivenRequestWithWidgetFilterWithOldQuery_ShouldBuildQueryWithLookupsAfterLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filterId := "test-filter"
	filter := view.WidgetFilter{
		OldMongoQuery: `{"$and": [
			{"v.connector": "test-connector"},
			{"v.duration": {"$gt": 600}},
			{"pbehavior._id": "test-pbehavior"},
			{"entity.category.name": "test-category"},
			{"v.infos.*.info_name": 3}
		]}`,
	}
	mockDbClient := createMockDbClientWithFilterFetching(ctrl, filterId, filter)
	request := ListRequestWithPagination{
		Query: pagination.GetDefaultQuery(),
		ListRequest: ListRequest{
			FilterRequest: FilterRequest{
				BaseFilterRequest: BaseFilterRequest{
					Filter: filterId,
				},
			},
		},
	}
	now := types.NewCpsTime()
	expectedDataPipeline := []bson.M{
		{"$sort": bson.D{{Key: "t", Value: -1}, {Key: "_id", Value: 1}}},
		{"$skip": 0},
		{"$limit": 10},
	}
	expectedDataPipeline = append(expectedDataPipeline, getEntityLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getEntityCategoryLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorInfoTypeLookup()...)
	fields := getComputedFields(now)
	durationField := fields["v.duration"]
	infosField := fields["infos"]
	delete(fields, "v.duration")
	delete(fields, "infos")
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$addFields": fields,
	})
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$project": bson.M{
			"v.steps":             0,
			"pbehavior.comments":  0,
			"pbehavior_info_type": 0,
		},
	})
	expected := []bson.M{
		{"$addFields": bson.M{
			"v.duration":    durationField,
			"infos":         infosField,
			"v.infos_array": bson.M{"$objectToArray": "$v.infos"},
		}},
		{"$match": bson.M{"$and": []bson.M{{"v.meta": nil}}}},
	}
	expected = append(expected, getEntityLookup()...)
	expected = append(expected, getEntityCategoryLookup()...)
	expected = append(expected, getPbehaviorLookup()...)
	expected = append(expected,
		bson.M{"$match": bson.M{"entity.enabled": true}},
		bson.M{"$match": bson.M{"$and": []bson.M{
			{"v.connector": "test-connector"},
			{"v.duration": bson.M{"$gt": 600}},
			{"pbehavior._id": "test-pbehavior"},
			{"entity.category.name": "test-category"},
			{"v.infos_array.v.info_name": 3},
		}}},
		bson.M{"$project": bson.M{
			"entity":          0,
			"entity.category": 0,
			"pbehavior":       0,
		}},
		bson.M{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		bson.M{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	)

	b := NewMongoQueryBuilder(mockDbClient)
	result, err := b.CreateListAggregationPipeline(ctx, request, now)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("unexpected result: %s", diff)
	}
}

func TestMongoQueryBuilder_CreateListAggregationPipeline_GivenRequestWithCategoryFilter_ShouldBuildQueryWithLookupsBeforeLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockDbClient := createMockDbClient(ctrl)
	request := ListRequestWithPagination{
		Query: pagination.GetDefaultQuery(),
		ListRequest: ListRequest{
			FilterRequest: FilterRequest{
				BaseFilterRequest: BaseFilterRequest{
					Category: "test-category",
				},
			},
		},
	}
	now := types.NewCpsTime()
	expectedDataPipeline := []bson.M{
		{"$sort": bson.D{{Key: "t", Value: -1}, {Key: "_id", Value: 1}}},
		{"$skip": 0},
		{"$limit": 10},
	}
	expectedDataPipeline = append(expectedDataPipeline, getEntityLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getEntityCategoryLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorInfoTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$addFields": getComputedFields(now),
	})
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$project": bson.M{
			"v.steps":             0,
			"pbehavior.comments":  0,
			"pbehavior_info_type": 0,
		},
	})
	expected := []bson.M{
		{"$match": bson.M{"$and": []bson.M{{"v.meta": nil}}}},
	}
	expected = append(expected, getEntityLookup()...)
	expected = append(expected,
		bson.M{"$match": bson.M{"entity.enabled": true}},
		bson.M{"$match": bson.M{"$and": []bson.M{
			{"entity.category": bson.M{"$eq": "test-category"}},
		}}},
		bson.M{"$project": bson.M{
			"entity": 0,
		}},
		bson.M{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		bson.M{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	)

	b := NewMongoQueryBuilder(mockDbClient)
	result, err := b.CreateListAggregationPipeline(ctx, request, now)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("unexpected result: %s", diff)
	}
}

func TestMongoQueryBuilder_CreateListAggregationPipeline_GivenRequestWithInstructionsFilter_ShouldBuildQueryWithLookupsBeforeLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	instructionId := "test-instruction"
	mockCursor := mock_mongo.NewMockCursor(ctrl)
	mockCursor.EXPECT().Next(gomock.Any()).Return(true)
	mockCursor.EXPECT().Next(gomock.Any()).Return(false)
	mockCursor.EXPECT().Close(gomock.Any())
	mockCursor.EXPECT().Decode(gomock.Any()).Do(func(instruction *Instruction) {
		durationCond, err := pattern.NewDurationCondition(pattern.ConditionGT, types.DurationWithUnit{
			Value: 10,
			Unit:  "m",
		})
		if err != nil {
			panic(err)
		}

		*instruction = Instruction{
			ActiveOnPbh:   []string{"maintenance"},
			DisabledOnPbh: []string{"pause"},
			AlarmPatternFields: savedpattern.AlarmPatternFields{
				AlarmPattern: pattern.Alarm{
					{
						{
							Field:     "v.duration",
							Condition: durationCond,
						},
						{
							Field:     "v.infos.info_name",
							FieldType: pattern.FieldTypeInt,
							Condition: pattern.NewIntCondition(pattern.ConditionEqual, 3),
						},
					},
				},
			},
			EntityPatternFields: savedpattern.EntityPatternFields{
				EntityPattern: pattern.Entity{
					{
						{
							Field:     "category",
							Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-category"),
						},
					},
				},
			},
		}
	})
	mockFilterDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockInstructionDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockInstructionDbCollection.EXPECT().Find(gomock.Any(), gomock.Eq(bson.M{
		"_id":    bson.M{"$in": []string{instructionId}},
		"status": bson.M{"$in": bson.A{InstructionStatusApproved, nil}},
	})).Return(mockCursor, nil)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.EXPECT().Collection(gomock.Any()).DoAndReturn(func(name string) mongo.DbCollection {
		switch name {
		case mongo.WidgetFiltersMongoCollection:
			return mockFilterDbCollection
		case mongo.InstructionMongoCollection:
			return mockInstructionDbCollection
		default:
			return nil
		}
	}).AnyTimes()

	hasRunningExecution := true
	request := ListRequestWithPagination{
		Query: pagination.GetDefaultQuery(),
		ListRequest: ListRequest{
			FilterRequest: FilterRequest{
				BaseFilterRequest: BaseFilterRequest{
					HasRunningExecution: &hasRunningExecution,
					IncludeInstructions: []string{instructionId},
				},
			},
		},
	}
	now := types.NewCpsTime()
	expectedDataPipeline := []bson.M{
		{"$sort": bson.D{{Key: "t", Value: -1}, {Key: "_id", Value: 1}}},
		{"$skip": 0},
		{"$limit": 10},
	}
	expectedDataPipeline = append(expectedDataPipeline, getEntityLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getEntityCategoryLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorInfoTypeLookup()...)
	fields := getComputedFields(now)
	durationField := fields["v.duration"]
	delete(fields, "v.duration")
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$addFields": fields,
	})
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$project": bson.M{
			"v.steps":             0,
			"pbehavior.comments":  0,
			"pbehavior_info_type": 0,
		},
	})
	expected := []bson.M{
		{"$addFields": bson.M{
			"v.duration":    durationField,
			"v.infos_array": bson.M{"$objectToArray": "$v.infos"},
		}},
		{"$match": bson.M{"$and": []bson.M{{"v.meta": nil}}}},
	}
	expected = append(expected, getEntityLookup()...)
	expected = append(expected, getInstructionExecutionLookup()...)
	expected = append(expected,
		bson.M{"$match": bson.M{"entity.enabled": true}},
		bson.M{"$match": bson.M{"$and": []bson.M{
			{"instruction_execution": bson.M{"$ne": nil}},
			{"$or": []bson.M{
				{"$and": []bson.M{
					{"$or": []bson.M{{"$and": []bson.M{
						{"v.duration": bson.M{"$gt": 600}},
						{"$and": []bson.M{
							{"v.infos_array.v.info_name": bson.M{"$type": bson.A{"long", "int", "decimal"}}},
							{"v.infos_array.v.info_name": bson.M{"$eq": 3}},
						}},
					}}}},
					{"$or": []bson.M{{"$and": []bson.M{
						{"entity.category": bson.M{"$eq": "test-category"}},
					}}}},
					{"v.pbehavior_info.type": bson.M{"$in": []string{"maintenance"}}},
					{"v.pbehavior_info.type": bson.M{"$nin": []string{"pause"}}},
				}},
			}},
		}}},
		bson.M{"$project": bson.M{
			"entity":                0,
			"instruction_execution": 0,
		}},
		bson.M{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		bson.M{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	)

	b := NewMongoQueryBuilder(mockDbClient)
	result, err := b.CreateListAggregationPipeline(ctx, request, now)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("unexpected result: %s", diff)
	}
}

func TestMongoQueryBuilder_CreateListAggregationPipeline_GivenRequestWithEntitySort_ShouldBuildQueryWithLookupsBeforeLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filterId := "test-filter"
	filter := view.WidgetFilter{
		OldMongoQuery: `{"$and": [
			{"pbehavior._id": "test-pbehavior"},
			{"entity.name": "test-entity"}
		]}`,
	}
	mockDbClient := createMockDbClientWithFilterFetching(ctrl, filterId, filter)
	request := ListRequestWithPagination{
		Query: pagination.GetDefaultQuery(),
		ListRequest: ListRequest{
			FilterRequest: FilterRequest{
				BaseFilterRequest: BaseFilterRequest{
					Filter: filterId,
				},
			},
			SortRequest: SortRequest{
				MultiSort: []string{
					"entity._id,desc",
					"entity.category.name,asc",
				},
			},
		},
	}
	now := types.NewCpsTime()
	expectedDataPipeline := []bson.M{
		{"$sort": bson.D{
			{Key: "entity._id", Value: -1},
			{Key: "entity.category.name", Value: 1},
			{Key: "_id", Value: 1},
		}},
		{"$skip": 0},
		{"$limit": 10},
	}
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorInfoTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$addFields": getComputedFields(now),
	})
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$project": bson.M{
			"v.steps":             0,
			"pbehavior.comments":  0,
			"pbehavior_info_type": 0,
		},
	})
	expected := []bson.M{
		{"$addFields": bson.M{"v.infos_array": bson.M{"$objectToArray": "$v.infos"}}},
		{"$match": bson.M{"$and": []bson.M{{"v.meta": nil}}}},
	}
	expected = append(expected, getEntityLookup()...)
	expected = append(expected, getPbehaviorLookup()...)
	expected = append(expected,
		bson.M{"$match": bson.M{"entity.enabled": true}},
		bson.M{"$match": bson.M{"$and": []bson.M{
			{"pbehavior._id": "test-pbehavior"},
			{"entity.name": "test-entity"},
		}}},
		bson.M{"$project": bson.M{"pbehavior": 0}},
	)
	expected = append(expected, getEntityCategoryLookup()...)
	expected = append(expected,
		bson.M{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		bson.M{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	)

	b := NewMongoQueryBuilder(mockDbClient)
	result, err := b.CreateListAggregationPipeline(ctx, request, now)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("unexpected result: %s", diff)
	}
}

func TestMongoQueryBuilder_CreateListAggregationPipeline_GivenRequestWithDurationSort_ShouldBuildQueryWithAddFieldsBeforeLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filterId := "test-filter"
	durationCond, err := pattern.NewDurationCondition(pattern.ConditionGT, types.DurationWithUnit{
		Value: 10,
		Unit:  "m",
	})
	if err != nil {
		panic(err)
	}
	filter := view.WidgetFilter{
		AlarmPatternFields: savedpattern.AlarmPatternFields{
			AlarmPattern: pattern.Alarm{
				{
					{
						Field:     "v.duration",
						Condition: durationCond,
					},
				},
			},
		},
	}
	mockDbClient := createMockDbClientWithFilterFetching(ctrl, filterId, filter)
	request := ListRequestWithPagination{
		Query: pagination.GetDefaultQuery(),
		ListRequest: ListRequest{
			FilterRequest: FilterRequest{
				BaseFilterRequest: BaseFilterRequest{
					Filter: filterId,
				},
			},
			SortRequest: SortRequest{
				MultiSort: []string{
					"v.duration,desc",
					"v.active_duration,desc",
				},
			},
		},
	}
	now := types.NewCpsTime()
	expectedDataPipeline := []bson.M{
		{"$sort": bson.D{
			{Key: "v.duration", Value: -1},
			{Key: "v.active_duration", Value: -1},
			{Key: "_id", Value: 1},
		}},
		{"$skip": 0},
		{"$limit": 10},
	}
	expectedDataPipeline = append(expectedDataPipeline, getEntityLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getEntityCategoryLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorInfoTypeLookup()...)
	fields := getComputedFields(now)
	durationField := fields["v.duration"]
	activeDurationField := fields["v.active_duration"]
	delete(fields, "v.duration")
	delete(fields, "v.active_duration")
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$addFields": fields,
	})
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$project": bson.M{
			"v.steps":             0,
			"pbehavior.comments":  0,
			"pbehavior_info_type": 0,
		},
	})
	expected := []bson.M{
		{"$addFields": bson.M{"v.duration": durationField}},
		{"$match": bson.M{"$or": []bson.M{{"$and": []bson.M{
			{"v.duration": bson.M{"$gt": 600}},
		}}}}},
		{"$match": bson.M{"$and": []bson.M{{"v.meta": nil}}}},
	}
	expected = append(expected, getEntityLookup()...)
	expected = append(expected, []bson.M{
		{"$match": bson.M{"entity.enabled": true}},
		{"$project": bson.M{"entity": 0}},
		{"$addFields": bson.M{"v.active_duration": activeDurationField}},
		{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	}...)

	b := NewMongoQueryBuilder(mockDbClient)
	result, err := b.CreateListAggregationPipeline(ctx, request, now)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("unexpected result: %s", diff)
	}
}

func TestMongoQueryBuilder_CreateListAggregationPipeline_GivenRequestWithSearch_ShouldBuildQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockDbClient := createMockDbClient(ctrl)
	search := "test-search"
	searchRegexp := primitive.Regex{
		Pattern: fmt.Sprintf(".*%s.*", search),
		Options: "i",
	}
	request := ListRequestWithPagination{
		Query: pagination.GetDefaultQuery(),
		ListRequest: ListRequest{
			FilterRequest: FilterRequest{
				BaseFilterRequest: BaseFilterRequest{
					Search: search,
				},
			},
		},
	}
	now := types.NewCpsTime()
	expectedDataPipeline := []bson.M{
		{"$sort": bson.D{{Key: "t", Value: -1}, {Key: "_id", Value: 1}}},
		{"$skip": 0},
		{"$limit": 10},
	}
	expectedDataPipeline = append(expectedDataPipeline, getEntityLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getEntityCategoryLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorInfoTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$addFields": getComputedFields(now),
	})
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$project": bson.M{
			"v.steps":             0,
			"pbehavior.comments":  0,
			"pbehavior_info_type": 0,
		},
	})
	expected := []bson.M{
		{"$match": bson.M{"$and": []bson.M{
			{"v.meta": nil},
			{"$or": []bson.M{
				{"v.connector": searchRegexp},
				{"v.connector_name": searchRegexp},
				{"v.component": searchRegexp},
				{"v.resource": searchRegexp},
			}},
		}}},
	}
	expected = append(expected, getEntityLookup()...)
	expected = append(expected, []bson.M{
		{"$match": bson.M{"entity.enabled": true}},
		{"$project": bson.M{"entity": 0}},
		{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	}...)

	b := NewMongoQueryBuilder(mockDbClient)
	result, err := b.CreateListAggregationPipeline(ctx, request, now)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("unexpected result: %s", diff)
	}
}

func TestMongoQueryBuilder_CreateListAggregationPipeline_GivenRequestWithSearchAndOnlyParents_ShouldBuildQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockDbClient := createMockDbClient(ctrl)
	search := "test-search"
	searchRegexp := primitive.Regex{
		Pattern: fmt.Sprintf(".*%s.*", search),
		Options: "i",
	}
	request := ListRequestWithPagination{
		Query: pagination.GetDefaultQuery(),
		ListRequest: ListRequest{
			FilterRequest: FilterRequest{
				BaseFilterRequest: BaseFilterRequest{
					Search:      search,
					OnlyParents: true,
				},
			},
		},
	}
	filteredChildrenLookup := getFilteredChildrenLookup(mongo.AlarmMongoCollection, bson.M{
		"$or": []bson.M{
			{"v.connector": searchRegexp},
			{"v.connector_name": searchRegexp},
			{"v.component": searchRegexp},
			{"v.resource": searchRegexp},
		},
	})
	now := types.NewCpsTime()
	expectedDataPipeline := []bson.M{
		{"$sort": bson.D{{Key: "t", Value: -1}, {Key: "_id", Value: 1}}},
		{"$skip": 0},
		{"$limit": 10},
	}
	expectedDataPipeline = append(expectedDataPipeline, getEntityLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getEntityCategoryLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorInfoTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getMetaAlarmRuleLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getChildrenCountLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, filteredChildrenLookup...)
	fields := getComputedFields(now)
	fields["is_meta_alarm"] = getIsMetaAlarmField()
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$addFields": fields,
	})
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$project": bson.M{
			"v.steps":             0,
			"pbehavior.comments":  0,
			"pbehavior_info_type": 0,
			"resolved_children":   0,
		},
	})
	expected := []bson.M{
		{"$match": bson.M{"$and": []bson.M{{"$or": []bson.M{
			{"v.parents": nil},
			{"v.parents": bson.M{"$eq": bson.A{}}},
			{"v.meta": bson.M{"$ne": nil}},
		}}}}},
	}
	expected = append(expected, getEntityLookup()...)
	expected = append(expected, filteredChildrenLookup...)
	expected = append(expected,
		bson.M{"$match": bson.M{"entity.enabled": true}},
		bson.M{"$match": bson.M{"$or": []bson.M{
			{"v.connector": searchRegexp},
			{"v.connector_name": searchRegexp},
			{"v.component": searchRegexp},
			{"v.resource": searchRegexp},
			{"filtered_children": bson.M{"$ne": bson.A{}}},
		}}},
		bson.M{"$project": bson.M{"entity": 0, "filtered_children": 0}},
		bson.M{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		bson.M{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	)

	b := NewMongoQueryBuilder(mockDbClient)
	result, err := b.CreateListAggregationPipeline(ctx, request, now)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("unexpected result: %s", diff)
	}
}

func TestMongoQueryBuilder_CreateListAggregationPipeline_GivenRequestWithSearchExpression_ShouldBuildQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockDbClient := createMockDbClient(ctrl)
	request := ListRequestWithPagination{
		Query: pagination.GetDefaultQuery(),
		ListRequest: ListRequest{
			FilterRequest: FilterRequest{
				BaseFilterRequest: BaseFilterRequest{
					Search: "connector LIKE \"test name\" AND state = 3",
				},
			},
		},
	}
	now := types.NewCpsTime()
	expectedDataPipeline := []bson.M{
		{"$sort": bson.D{{Key: "t", Value: -1}, {Key: "_id", Value: 1}}},
		{"$skip": 0},
		{"$limit": 10},
	}
	expectedDataPipeline = append(expectedDataPipeline, getEntityLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getEntityCategoryLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorInfoTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$addFields": getComputedFields(now),
	})
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$project": bson.M{
			"v.steps":             0,
			"pbehavior.comments":  0,
			"pbehavior_info_type": 0,
		},
	})
	expected := []bson.M{
		{"$match": bson.M{"$and": []bson.M{
			{"v.meta": nil},
			{"$and": []bson.M{
				{"v.connector": bson.M{"$regex": "test name"}},
				{"v.state.val": bson.M{"$eq": 3}},
			}},
		}}},
	}
	expected = append(expected, getEntityLookup()...)
	expected = append(expected, []bson.M{
		{"$match": bson.M{"entity.enabled": true}},
		{"$project": bson.M{"entity": 0}},
		{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	}...)

	b := NewMongoQueryBuilder(mockDbClient)
	result, err := b.CreateListAggregationPipeline(ctx, request, now)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("unexpected result: %s", diff)
	}
}

func TestMongoQueryBuilder_CreateListAggregationPipeline_GivenRequestWithSearchExpression_ShouldBuildQueryWithLookupsBeforeMatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockDbClient := createMockDbClient(ctrl)
	request := ListRequestWithPagination{
		Query: pagination.GetDefaultQuery(),
		ListRequest: ListRequest{
			FilterRequest: FilterRequest{
				BaseFilterRequest: BaseFilterRequest{
					Search: "entity.name LIKE \"test name\" AND v.duration > 100",
				},
			},
		},
	}
	now := types.NewCpsTime()
	expectedDataPipeline := []bson.M{
		{"$sort": bson.D{{Key: "t", Value: -1}, {Key: "_id", Value: 1}}},
		{"$skip": 0},
		{"$limit": 10},
	}
	expectedDataPipeline = append(expectedDataPipeline, getEntityLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getEntityCategoryLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorInfoTypeLookup()...)
	fields := getComputedFields(now)
	durationField := fields["v.duration"]
	delete(fields, "v.duration")
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$addFields": fields,
	})
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$project": bson.M{
			"v.steps":             0,
			"pbehavior.comments":  0,
			"pbehavior_info_type": 0,
		},
	})
	expected := []bson.M{
		{"$addFields": bson.M{
			"v.duration": durationField,
		}},
		{"$match": bson.M{"$and": []bson.M{{"v.meta": nil}}}},
	}
	expected = append(expected, getEntityLookup()...)
	expected = append(expected, []bson.M{
		{"$match": bson.M{"entity.enabled": true}},
		{"$match": bson.M{"$and": []bson.M{
			{"entity.name": bson.M{"$regex": "test name"}},
			{"v.duration": bson.M{"$gt": 100}},
		}}},
		{"$project": bson.M{"entity": 0}},
		{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	}...)

	b := NewMongoQueryBuilder(mockDbClient)
	result, err := b.CreateListAggregationPipeline(ctx, request, now)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("unexpected result: %s", diff)
	}
}

func TestMongoQueryBuilder_CreateListAggregationPipeline_GivenRequestWithSearchExpression_ShouldBuildQueryWithReplaceInfosAlias(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockDbClient := createMockDbClient(ctrl)
	request := ListRequestWithPagination{
		Query: pagination.GetDefaultQuery(),
		ListRequest: ListRequest{
			FilterRequest: FilterRequest{
				BaseFilterRequest: BaseFilterRequest{
					Search: "infos.test1.value LIKE \"test val\"",
				},
			},
		},
	}
	now := types.NewCpsTime()
	expectedDataPipeline := []bson.M{
		{"$sort": bson.D{{Key: "t", Value: -1}, {Key: "_id", Value: 1}}},
		{"$skip": 0},
		{"$limit": 10},
	}
	expectedDataPipeline = append(expectedDataPipeline, getEntityLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getEntityCategoryLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorInfoTypeLookup()...)
	fields := getComputedFields(now)
	infosField := fields["infos"]
	delete(fields, "infos")
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$addFields": fields,
	})
	expectedDataPipeline = append(expectedDataPipeline, bson.M{
		"$project": bson.M{
			"v.steps":             0,
			"pbehavior.comments":  0,
			"pbehavior_info_type": 0,
		},
	})
	expected := []bson.M{
		{"$addFields": bson.M{
			"infos": infosField,
		}},
		{"$match": bson.M{"$and": []bson.M{{"v.meta": nil}}}},
	}
	expected = append(expected, getEntityLookup()...)
	expected = append(expected, []bson.M{
		{"$match": bson.M{"entity.enabled": true}},
		{"$match": bson.M{"entity.infos.test1.value": bson.M{"$regex": "test val"}}},
		{"$project": bson.M{"entity": 0}},
		{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	}...)

	b := NewMongoQueryBuilder(mockDbClient)
	result, err := b.CreateListAggregationPipeline(ctx, request, now)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("unexpected result: %s", diff)
	}
}

func createMockDbClient(ctrl *gomock.Controller) mongo.DbClient {
	mockFilterDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockInstructionDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.EXPECT().Collection(gomock.Any()).DoAndReturn(func(name string) mongo.DbCollection {
		switch name {
		case mongo.WidgetFiltersMongoCollection:
			return mockFilterDbCollection
		case mongo.InstructionMongoCollection:
			return mockInstructionDbCollection
		default:
			return nil
		}
	}).AnyTimes()

	return mockDbClient
}

func createMockDbClientWithFilterFetching(ctrl *gomock.Controller, filterId string, filter view.WidgetFilter) mongo.DbClient {
	mockSingleResult := mock_mongo.NewMockSingleResultHelper(ctrl)
	mockSingleResult.EXPECT().Decode(gomock.Any()).Do(func(v *view.WidgetFilter) {
		*v = filter
	})
	mockFilterDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockFilterDbCollection.EXPECT().FindOne(gomock.Any(), gomock.Eq(bson.M{"_id": filterId})).Return(mockSingleResult)
	mockInstructionDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.EXPECT().Collection(gomock.Any()).DoAndReturn(func(name string) mongo.DbCollection {
		switch name {
		case mongo.WidgetFiltersMongoCollection:
			return mockFilterDbCollection
		case mongo.InstructionMongoCollection:
			return mockInstructionDbCollection
		default:
			return nil
		}
	}).AnyTimes()

	return mockDbClient
}
