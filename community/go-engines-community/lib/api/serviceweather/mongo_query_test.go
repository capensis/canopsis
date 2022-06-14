package serviceweather

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
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

	mockDbClient := createMockDbClient(ctrl)
	request := ListRequest{
		Query: pagination.Query{
			Page:     2,
			Limit:    10,
			Paginate: true,
		},
	}
	expectedDataPipeline := []bson.M{
		{"$sort": bson.D{{Key: "name", Value: 1}}},
		{"$skip": 10},
		{"$limit": 10},
	}
	expectedDataPipeline = append(expectedDataPipeline, getCategoryLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getAlarmLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorAlarmCountersLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorInfoTypeLookup()...)
	expected := []bson.M{
		{"$match": bson.M{
			"type":    types.EntityTypeService,
			"enabled": true,
		}},
		{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	}

	b := NewMongoQueryBuilder(mockDbClient)
	result, err := b.CreateListAggregationPipeline(ctx, request)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("unexpected result: %s", diff)
	}
}

func TestMongoQueryBuilder_CreateListAggregationPipeline_GivenRequestWithWidgetFilterWithEntityPattern_ShouldBuildQueryWithLookupsAfterLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filterId := "test-filter"
	filter := view.WidgetFilter{
		EntityPatternFields: savedpattern.EntityPatternFields{
			EntityPattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource"),
					},
				},
			},
		},
	}
	mockDbClient := createMockDbClientWithFilterFetching(ctrl, filterId, filter)
	request := ListRequest{
		Query:  pagination.GetDefaultQuery(),
		Filter: filterId,
	}
	expectedDataPipeline := []bson.M{
		{"$sort": bson.D{{Key: "name", Value: 1}}},
		{"$skip": 0},
		{"$limit": 10},
	}
	expectedDataPipeline = append(expectedDataPipeline, getCategoryLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getAlarmLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorAlarmCountersLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorInfoTypeLookup()...)
	expected := []bson.M{
		{"$match": bson.M{
			"type":    types.EntityTypeService,
			"enabled": true,
		}},
		{"$match": bson.M{"$or": []bson.M{{"$and": []bson.M{
			{"name": bson.M{"$eq": "test-resource"}},
		}}}}},
		{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	}

	b := NewMongoQueryBuilder(mockDbClient)
	result, err := b.CreateListAggregationPipeline(ctx, request)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("unexpected result: %s", diff)
	}
}

func TestMongoQueryBuilder_CreateListAggregationPipeline_GivenRequestWithWidgetFilterWithWeatherServicePatternWithStateCond_ShouldBuildQueryWithLookupsBeforeLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filterId := "test-filter"
	filter := view.WidgetFilter{
		WeatherServicePattern: view.WeatherServicePattern{
			{
				{
					Field:     "state",
					Condition: pattern.NewIntCondition(pattern.ConditionEqual, types.AlarmStateMinor),
				},
			},
		},
	}
	mockDbClient := createMockDbClientWithFilterFetching(ctrl, filterId, filter)
	request := ListRequest{
		Query:  pagination.GetDefaultQuery(),
		Filter: filterId,
	}
	expectedDataPipeline := []bson.M{
		{"$sort": bson.D{{Key: "name", Value: 1}}},
		{"$skip": 0},
		{"$limit": 10},
	}
	expectedDataPipeline = append(expectedDataPipeline, getCategoryLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorAlarmCountersLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorInfoTypeLookup()...)
	expected := []bson.M{
		{"$match": bson.M{
			"type":    types.EntityTypeService,
			"enabled": true,
		}},
	}
	expected = append(expected, getAlarmLookup()...)
	expected = append(expected, []bson.M{
		{"$match": bson.M{"$or": []bson.M{{"$and": []bson.M{
			{"state": bson.M{"$eq": types.AlarmStateMinor}},
		}}}}},
		{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	}...)

	b := NewMongoQueryBuilder(mockDbClient)
	result, err := b.CreateListAggregationPipeline(ctx, request)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("unexpected result: %s", diff)
	}
}

func TestMongoQueryBuilder_CreateListAggregationPipeline_GivenRequestWithWidgetFilterWithWeatherServicePatternWithIconCond_ShouldBuildQueryWithLookupsBeforeLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filterId := "test-filter"
	filter := view.WidgetFilter{
		WeatherServicePattern: view.WeatherServicePattern{
			{
				{
					Field:     "icon",
					Condition: pattern.NewStringCondition(pattern.ConditionEqual, "pause"),
				},
			},
		},
	}
	mockDbClient := createMockDbClientWithFilterFetching(ctrl, filterId, filter)
	request := ListRequest{
		Query:  pagination.GetDefaultQuery(),
		Filter: filterId,
	}
	expectedDataPipeline := []bson.M{
		{"$sort": bson.D{{Key: "name", Value: 1}}},
		{"$skip": 0},
		{"$limit": 10},
	}
	expectedDataPipeline = append(expectedDataPipeline, getCategoryLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorInfoTypeLookup()...)
	expected := []bson.M{
		{"$match": bson.M{
			"type":    types.EntityTypeService,
			"enabled": true,
		}},
	}
	expected = append(expected, getAlarmLookup()...)
	expected = append(expected, getPbehaviorAlarmCountersLookup()...)
	expected = append(expected, []bson.M{
		{"$match": bson.M{"$or": []bson.M{{"$and": []bson.M{
			{"icon": bson.M{"$eq": "pause"}},
		}}}}},
		{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	}...)

	b := NewMongoQueryBuilder(mockDbClient)
	result, err := b.CreateListAggregationPipeline(ctx, request)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("unexpected result: %s", diff)
	}
}

func TestMongoQueryBuilder_CreateListAggregationPipeline_GivenRequestWithWidgetFilterWithOldMongoQuery_ShouldBuildQueryWithLookupsBeforeLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filterId := "test-filter"
	filter := view.WidgetFilter{
		OldMongoQuery: `{"$and": [
			{"type": "resource"},
			{"category": "test-category"}
		]}`,
	}
	mockDbClient := createMockDbClientWithFilterFetching(ctrl, filterId, filter)
	request := ListRequest{
		Query:  pagination.GetDefaultQuery(),
		Filter: filterId,
	}
	expectedDataPipeline := []bson.M{
		{"$sort": bson.D{{Key: "name", Value: 1}}},
		{"$skip": 0},
		{"$limit": 10},
	}
	expected := []bson.M{
		{"$match": bson.M{
			"type":    types.EntityTypeService,
			"enabled": true,
		}},
	}
	expected = append(expected, getCategoryLookup()...)
	expected = append(expected, getAlarmLookup()...)
	expected = append(expected, getPbehaviorLookup()...)
	expected = append(expected, getPbehaviorAlarmCountersLookup()...)
	expected = append(expected, getPbehaviorInfoTypeLookup()...)
	expected = append(expected, []bson.M{
		{"$match": bson.M{"$and": []bson.M{
			{"type": "resource"},
			{"category": "test-category"},
		}}},
		{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	}...)

	b := NewMongoQueryBuilder(mockDbClient)
	result, err := b.CreateListAggregationPipeline(ctx, request)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("unexpected result: %s", diff)
	}
}

func TestMongoQueryBuilder_CreateListAggregationPipeline_GivenRequestWithSortByState_ShouldBuildQueryWithLookupsBeforeLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockDbClient := createMockDbClient(ctrl)
	request := ListRequest{
		Query:  pagination.GetDefaultQuery(),
		Sort:   common.SortDesc,
		SortBy: "state",
	}
	expectedDataPipeline := []bson.M{
		{"$sort": bson.D{{Key: "state", Value: -1}, {Key: "name", Value: 1}}},
		{"$skip": 0},
		{"$limit": 10},
	}
	expectedDataPipeline = append(expectedDataPipeline, getCategoryLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorAlarmCountersLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorInfoTypeLookup()...)
	expected := []bson.M{
		{"$match": bson.M{
			"type":    types.EntityTypeService,
			"enabled": true,
		}},
	}
	expected = append(expected, getAlarmLookup()...)
	expected = append(expected, []bson.M{
		{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	}...)

	b := NewMongoQueryBuilder(mockDbClient)
	result, err := b.CreateListAggregationPipeline(ctx, request)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("unexpected result: %s", diff)
	}
}

func createMockDbClient(ctrl *gomock.Controller) mongo.DbClient {
	mockFilterDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.EXPECT().Collection(gomock.Any()).DoAndReturn(func(name string) mongo.DbCollection {
		switch name {
		case mongo.WidgetFiltersMongoCollection:
			return mockFilterDbCollection
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
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.EXPECT().Collection(gomock.Any()).DoAndReturn(func(name string) mongo.DbCollection {
		switch name {
		case mongo.WidgetFiltersMongoCollection:
			return mockFilterDbCollection
		default:
			return nil
		}
	}).AnyTimes()

	return mockDbClient
}
