package serviceweather

import (
	"context"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entity/dbquery"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	"github.com/golang/mock/gomock"
	"github.com/kylelemons/godebug/pretty"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMongoQueryBuilder_CreateListAggregationPipeline_GivenPaginationRequest_ShouldBuildQueryWithLookupsAfterLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockDbClient := createMockDbClient(ctrl)
	authorProvider := author.NewProvider(mockDbClient, config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop()))
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
	expectedDataPipeline = append(expectedDataPipeline, dbquery.GetCategoryLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getAlarmLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorLookup(authorProvider)...)
	expectedDataPipeline = append(expectedDataPipeline, dbquery.GetPbehaviorInfoTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorAlarmCountersLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, bson.M{"$project": bson.M{
		"services": 0,
	}})
	expectedDataPipeline = append(expectedDataPipeline, bson.M{"$sort": bson.D{{Key: "name", Value: 1}}})
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

	b := NewMongoQueryBuilder(mockDbClient, authorProvider)
	result, err := b.CreateListAggregationPipeline(ctx, request)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if diff := pretty.Compare(author.StripAuthorRandomPrefix(result), author.StripAuthorRandomPrefix(expected)); diff != "" {
		t.Errorf("unexpected result: %s", diff)
	}
}

func TestMongoQueryBuilder_CreateListAggregationPipeline_GivenRequestWithWidgetFilterWithEntityPattern_ShouldBuildQueryWithLookupsAfterLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filter := view.WidgetFilter{
		ID: "test-filter",
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
	mockDbClient := createMockDbClientWithFilterFetching(ctrl, []view.WidgetFilter{filter})
	authorProvider := author.NewProvider(mockDbClient, config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop()))
	request := ListRequest{
		Query:   pagination.GetDefaultQuery(),
		Filters: []string{filter.ID},
	}
	expectedDataPipeline := []bson.M{
		{"$sort": bson.D{{Key: "name", Value: 1}}},
		{"$skip": 0},
		{"$limit": 10},
	}
	expectedDataPipeline = append(expectedDataPipeline, dbquery.GetCategoryLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getAlarmLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorLookup(authorProvider)...)
	expectedDataPipeline = append(expectedDataPipeline, dbquery.GetPbehaviorInfoTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorAlarmCountersLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, bson.M{"$project": bson.M{
		"services": 0,
	}})
	expectedDataPipeline = append(expectedDataPipeline, bson.M{"$sort": bson.D{{Key: "name", Value: 1}}})
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

	b := NewMongoQueryBuilder(mockDbClient, authorProvider)
	result, err := b.CreateListAggregationPipeline(ctx, request)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if diff := pretty.Compare(author.StripAuthorRandomPrefix(result), author.StripAuthorRandomPrefix(expected)); diff != "" {
		t.Errorf("unexpected result: %s", diff)
	}
}

func TestMongoQueryBuilder_CreateListAggregationPipeline_GivenRequestWithWidgetFilterWithWeatherServicePatternWithStateCond_ShouldBuildQueryWithLookupsBeforeLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filter := view.WidgetFilter{
		ID: "test-filter",
		WeatherServicePattern: view.WeatherServicePattern{
			{
				{
					Field:     "state.val",
					Condition: pattern.NewIntCondition(pattern.ConditionEqual, types.AlarmStateMinor),
				},
			},
		},
	}
	mockDbClient := createMockDbClientWithFilterFetching(ctrl, []view.WidgetFilter{filter})
	authorProvider := author.NewProvider(mockDbClient, config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop()))
	request := ListRequest{
		Query:   pagination.GetDefaultQuery(),
		Filters: []string{filter.ID},
	}
	expectedDataPipeline := []bson.M{
		{"$sort": bson.D{{Key: "name", Value: 1}}},
		{"$skip": 0},
		{"$limit": 10},
	}
	expectedDataPipeline = append(expectedDataPipeline, dbquery.GetCategoryLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorLookup(authorProvider)...)
	expectedDataPipeline = append(expectedDataPipeline, dbquery.GetPbehaviorInfoTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorAlarmCountersLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, bson.M{"$project": bson.M{
		"services": 0,
	}})
	expectedDataPipeline = append(expectedDataPipeline, bson.M{"$sort": bson.D{{Key: "name", Value: 1}}})
	expected := []bson.M{
		{"$match": bson.M{
			"type":    types.EntityTypeService,
			"enabled": true,
		}},
	}
	expected = append(expected, getAlarmLookup()...)
	expected = append(expected, []bson.M{
		{"$match": bson.M{"$or": []bson.M{{"$and": []bson.M{
			{"state.val": bson.M{"$eq": types.AlarmStateMinor}},
		}}}}},
		{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	}...)

	b := NewMongoQueryBuilder(mockDbClient, authorProvider)
	result, err := b.CreateListAggregationPipeline(ctx, request)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if diff := pretty.Compare(author.StripAuthorRandomPrefix(result), author.StripAuthorRandomPrefix(expected)); diff != "" {
		t.Errorf("unexpected result: %s", diff)
	}
}

func TestMongoQueryBuilder_CreateListAggregationPipeline_GivenRequestWithWidgetFilterWithWeatherServicePatternWithIconCond_ShouldBuildQueryWithLookupsBeforeLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filter := view.WidgetFilter{
		ID: "test-filter",
		WeatherServicePattern: view.WeatherServicePattern{
			{
				{
					Field:     "icon",
					Condition: pattern.NewStringCondition(pattern.ConditionEqual, "pause"),
				},
			},
		},
	}
	mockDbClient := createMockDbClientWithFilterFetching(ctrl, []view.WidgetFilter{filter})
	authorProvider := author.NewProvider(mockDbClient, config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop()))
	request := ListRequest{
		Query:   pagination.GetDefaultQuery(),
		Filters: []string{filter.ID},
	}
	expectedDataPipeline := []bson.M{
		{"$sort": bson.D{{Key: "name", Value: 1}}},
		{"$skip": 0},
		{"$limit": 10},
	}
	expectedDataPipeline = append(expectedDataPipeline, dbquery.GetCategoryLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorLookup(authorProvider)...)
	expectedDataPipeline = append(expectedDataPipeline, bson.M{"$project": bson.M{
		"services": 0,
	}})
	expectedDataPipeline = append(expectedDataPipeline, bson.M{"$sort": bson.D{{Key: "name", Value: 1}}})
	expected := []bson.M{
		{"$match": bson.M{
			"type":    types.EntityTypeService,
			"enabled": true,
		}},
	}
	expected = append(expected, getAlarmLookup()...)
	expected = append(expected, dbquery.GetPbehaviorInfoTypeLookup()...)
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

	b := NewMongoQueryBuilder(mockDbClient, authorProvider)
	result, err := b.CreateListAggregationPipeline(ctx, request)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if diff := pretty.Compare(author.StripAuthorRandomPrefix(result), author.StripAuthorRandomPrefix(expected)); diff != "" {
		t.Errorf("unexpected result: %s", diff)
	}
}

func TestMongoQueryBuilder_CreateListAggregationPipeline_GivenRequestWithSortByState_ShouldBuildQueryWithLookupsBeforeLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockDbClient := createMockDbClient(ctrl)
	authorProvider := author.NewProvider(mockDbClient, config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop()))
	request := ListRequest{
		Query:  pagination.GetDefaultQuery(),
		Sort:   common.SortDesc,
		SortBy: "state",
	}
	expectedDataPipeline := []bson.M{
		{"$sort": bson.D{{Key: "state.val", Value: -1}, {Key: "name", Value: 1}}},
		{"$skip": 0},
		{"$limit": 10},
	}
	expectedDataPipeline = append(expectedDataPipeline, dbquery.GetCategoryLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorLookup(authorProvider)...)
	expectedDataPipeline = append(expectedDataPipeline, dbquery.GetPbehaviorInfoTypeLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorAlarmCountersLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, bson.M{"$project": bson.M{
		"services": 0,
	}})
	expectedDataPipeline = append(expectedDataPipeline, bson.M{"$sort": bson.D{{Key: "state.val", Value: -1}, {Key: "name", Value: 1}}})
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

	b := NewMongoQueryBuilder(mockDbClient, authorProvider)
	result, err := b.CreateListAggregationPipeline(ctx, request)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if diff := pretty.Compare(author.StripAuthorRandomPrefix(result), author.StripAuthorRandomPrefix(expected)); diff != "" {
		t.Errorf("unexpected result: %s", diff)
	}
}

func TestMongoQueryBuilder_CreateListAggregationPipeline_GivenRequestWithMultipleWidgetFilters_ShouldBuildQueryAllMatches(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filter1 := view.WidgetFilter{
		ID: "test-filter-1",
		EntityPatternFields: savedpattern.EntityPatternFields{
			EntityPattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-service"),
					},
				},
			},
		},
		WeatherServicePattern: view.WeatherServicePattern{
			{
				{
					Field:     "icon",
					Condition: pattern.NewStringCondition(pattern.ConditionEqual, "pause"),
				},
			},
		},
	}
	filter2 := view.WidgetFilter{
		ID: "test-filter-2",
		WeatherServicePattern: view.WeatherServicePattern{
			{
				{
					Field:     "secondary_icon",
					Condition: pattern.NewStringCondition(pattern.ConditionEqual, "pause"),
				},
			},
		},
	}
	mockDbClient := createMockDbClientWithFilterFetching(ctrl, []view.WidgetFilter{filter1, filter2})
	authorProvider := author.NewProvider(mockDbClient, config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop()))
	request := ListRequest{
		Query:   pagination.GetDefaultQuery(),
		Filters: []string{filter1.ID, filter2.ID},
	}
	expectedDataPipeline := []bson.M{
		{"$sort": bson.D{{Key: "name", Value: 1}}},
		{"$skip": 0},
		{"$limit": 10},
	}
	expectedDataPipeline = append(expectedDataPipeline, dbquery.GetCategoryLookup()...)
	expectedDataPipeline = append(expectedDataPipeline, getPbehaviorLookup(authorProvider)...)
	expectedDataPipeline = append(expectedDataPipeline, bson.M{"$project": bson.M{
		"services": 0,
	}})
	expectedDataPipeline = append(expectedDataPipeline, bson.M{"$sort": bson.D{{Key: "name", Value: 1}}})
	expected := []bson.M{
		{"$match": bson.M{
			"type":    types.EntityTypeService,
			"enabled": true,
		}},
		{"$match": bson.M{"$or": []bson.M{{"$and": []bson.M{
			{"name": bson.M{"$eq": "test-service"}},
		}}}}},
	}
	expected = append(expected, getAlarmLookup()...)
	expected = append(expected, dbquery.GetPbehaviorInfoTypeLookup()...)
	expected = append(expected, getPbehaviorAlarmCountersLookup()...)
	expected = append(expected, []bson.M{
		{"$match": bson.M{"$or": []bson.M{{"$and": []bson.M{
			{"icon": bson.M{"$eq": "pause"}},
		}}}}},
		{"$match": bson.M{"$or": []bson.M{{"$and": []bson.M{
			{"secondary_icon": bson.M{"$eq": "pause"}},
		}}}}},
		{"$facet": bson.M{
			"data":        expectedDataPipeline,
			"total_count": []bson.M{{"$count": "count"}},
		}},
		{"$addFields": bson.M{
			"total_count": bson.M{"$sum": "$total_count.count"},
		}},
	}...)

	b := NewMongoQueryBuilder(mockDbClient, authorProvider)
	result, err := b.CreateListAggregationPipeline(ctx, request)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if diff := pretty.Compare(author.StripAuthorRandomPrefix(result), author.StripAuthorRandomPrefix(expected)); diff != "" {
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

func createMockDbClientWithFilterFetching(ctrl *gomock.Controller, filters []view.WidgetFilter) mongo.DbClient {
	mockFilterDbCollection := mock_mongo.NewMockDbCollection(ctrl)

	for _, v := range filters {
		filter := v
		mockSingleResult := mock_mongo.NewMockSingleResultHelper(ctrl)
		mockSingleResult.EXPECT().Decode(gomock.Any()).Do(func(v *view.WidgetFilter) {
			*v = filter
		})
		mockFilterDbCollection.EXPECT().FindOne(gomock.Any(), gomock.Eq(bson.M{"_id": filter.ID})).Return(mockSingleResult)
	}
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
