package alarm

import (
	"context"
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
	"testing"
)

func TestMongoQueryBuilder_CreateListAggregationPipeline_GivenPaginationRequest_ShouldBuildQueryWithLookupsAfterLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockFilterDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockInstructionDbCollection := mock_mongo.NewMockDbCollection(ctrl)
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
		{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	}

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
	mockSingleResult := mock_mongo.NewMockSingleResultHelper(ctrl)
	mockSingleResult.EXPECT().Decode(gomock.Any()).Do(func(v *view.WidgetFilter) {
		*v = view.WidgetFilter{
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
							Field:     "canonical_type",
							Condition: pattern.NewStringCondition(pattern.ConditionEqual, "pause"),
						},
					},
				},
			},
		}
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
		{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	}

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
	mockSingleResult := mock_mongo.NewMockSingleResultHelper(ctrl)
	mockSingleResult.EXPECT().Decode(gomock.Any()).Do(func(v *view.WidgetFilter) {
		durationCond, err := pattern.NewDurationCondition(pattern.ConditionGT, types.DurationWithUnit{
			Value: 10,
			Unit:  "m",
		})
		if err != nil {
			panic(err)
		}
		*v = view.WidgetFilter{
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
		{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	}

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
	mockSingleResult := mock_mongo.NewMockSingleResultHelper(ctrl)
	mockSingleResult.EXPECT().Decode(gomock.Any()).Do(func(v *view.WidgetFilter) {
		*v = view.WidgetFilter{
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
		{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	}

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
	mockSingleResult := mock_mongo.NewMockSingleResultHelper(ctrl)
	mockSingleResult.EXPECT().Decode(gomock.Any()).Do(func(v *view.WidgetFilter) {
		*v = view.WidgetFilter{
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
							Field:     "canonical_type",
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
	mockSingleResult := mock_mongo.NewMockSingleResultHelper(ctrl)
	mockSingleResult.EXPECT().Decode(gomock.Any()).Do(func(v *view.WidgetFilter) {
		*v = view.WidgetFilter{
			OldMongoQuery: map[string]interface{}{
				"$and": []map[string]interface{}{
					{"v.connector": "test-connector"},
				},
			},
		}
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
		{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	}

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
	mockSingleResult := mock_mongo.NewMockSingleResultHelper(ctrl)
	mockSingleResult.EXPECT().Decode(gomock.Any()).Do(func(v *view.WidgetFilter) {
		*v = view.WidgetFilter{
			OldMongoQuery: map[string]interface{}{
				"$and": []map[string]interface{}{
					{"v.connector": "test-connector"},
					{"v.duration": bson.M{"$gt": 600}},
					{"pbehavior._id": "test-pbehavior"},
					{"entity.category.name": "test-category"},
					{"v.infos.*.info_name": 3},
				},
			},
		}
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

	request := ListRequestWithPagination{
		Query: pagination.GetDefaultQuery(),
		ListRequest: ListRequest{
			FilterRequest: FilterRequest{
				BaseFilterRequest: BaseFilterRequest{
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
	expected = append(expected,
		bson.M{"$match": bson.M{"$and": []bson.M{
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

func TestMongoQueryBuilder_CreateListAggregationPipeline_GivenRequestWithEntitySort_ShouldBuildQueryWithLookupsBeforeLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filterId := "test-filter"
	mockSingleResult := mock_mongo.NewMockSingleResultHelper(ctrl)
	mockSingleResult.EXPECT().Decode(gomock.Any()).Do(func(v *view.WidgetFilter) {
		*v = view.WidgetFilter{
			OldMongoQuery: map[string]interface{}{
				"$and": []map[string]interface{}{
					{"pbehavior._id": "test-pbehavior"},
					{"entity.name": "test-entity"},
				},
			},
		}
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
	mockSingleResult := mock_mongo.NewMockSingleResultHelper(ctrl)
	mockSingleResult.EXPECT().Decode(gomock.Any()).Do(func(v *view.WidgetFilter) {
		durationCond, err := pattern.NewDurationCondition(pattern.ConditionGT, types.DurationWithUnit{
			Value: 10,
			Unit:  "m",
		})
		if err != nil {
			panic(err)
		}
		*v = view.WidgetFilter{
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
		{"$addFields": bson.M{"v.active_duration": activeDurationField}},
		{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	}

	b := NewMongoQueryBuilder(mockDbClient)
	result, err := b.CreateListAggregationPipeline(ctx, request, now)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("unexpected result: %s", diff)
	}
}
