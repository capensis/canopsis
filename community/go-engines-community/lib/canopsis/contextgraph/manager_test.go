package contextgraph_test

import (
	"context"
	"fmt"
	"math/rand"
	"slices"
	"sort"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/contextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	mock_contextgraph "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/contextgraph"
	mock_metrics "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/metrics"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	"github.com/golang/mock/gomock"
	"github.com/kylelemons/godebug/pretty"
)

func TestCheckServices(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	collection := mock_mongo.NewMockDbCollection(ctrl)

	dbClient := mock_mongo.NewMockDbClient(ctrl)
	dbClient.EXPECT().Collection(mongo.EntityMongoCollection).Return(collection).AnyTimes()

	adapter := entity.NewAdapter(dbClient)
	storage := mock_contextgraph.NewMockEntityServiceStorage(ctrl)

	metaUpdater := mock_metrics.NewMockMetaUpdater(ctrl)

	dataSets := []struct {
		services       []entityservice.EntityService
		entities       []types.Entity
		expectedResult []types.Entity
		name           string
	}{
		{
			name: "one entity is added to a single service",
			entities: []types.Entity{
				{
					ID:        "id-1",
					Component: "component-1",
					Enabled:   true,
				},
			},
			services: []entityservice.EntityService{
				{
					Entity: types.Entity{
						ID:      "serv-1",
						Enabled: true,
					},
					EntityPatternFields: savedpattern.EntityPatternFields{
						EntityPattern: [][]pattern.FieldCondition{
							{
								{
									Field:     "component",
									Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
								},
							},
						},
					},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:               "id-1",
					Component:        "component-1",
					Enabled:          true,
					ServicesToAdd:    []string{"serv-1"},
					ServicesToRemove: []string{},
					Services:         []string{"serv-1"},
				},
			},
		},
		{
			name: "one entity is added to multiple services",
			entities: []types.Entity{
				{
					ID:        "id-1",
					Component: "component-1",
					Enabled:   true,
				},
			},
			services: []entityservice.EntityService{
				{
					Entity: types.Entity{
						ID:      "serv-1",
						Enabled: true,
					},
					EntityPatternFields: savedpattern.EntityPatternFields{
						EntityPattern: [][]pattern.FieldCondition{
							{
								{
									Field:     "component",
									Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
								},
							},
						},
					},
				},
				{
					Entity: types.Entity{
						ID:      "serv-2",
						Enabled: true,
					},
					EntityPatternFields: savedpattern.EntityPatternFields{
						EntityPattern: [][]pattern.FieldCondition{
							{
								{
									Field:     "component",
									Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
								},
							},
						},
					},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:               "id-1",
					Component:        "component-1",
					Enabled:          true,
					ServicesToAdd:    []string{"serv-1", "serv-2"},
					ServicesToRemove: []string{},
					Services:         []string{"serv-1", "serv-2"},
				},
			},
		},
		{
			name: "one entity is added to multiple services impacted services to add/remove should be updated",
			entities: []types.Entity{
				{
					ID:               "id-1",
					Component:        "component-1",
					Enabled:          true,
					ServicesToAdd:    []string{"serv-4"},
					ServicesToRemove: []string{"serv-0", "serv-2", "serv-3"},
				},
			},
			services: []entityservice.EntityService{
				{
					Entity: types.Entity{
						ID:      "serv-1",
						Enabled: true,
					},
					EntityPatternFields: savedpattern.EntityPatternFields{
						EntityPattern: [][]pattern.FieldCondition{
							{
								{
									Field:     "component",
									Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
								},
							},
						},
					},
				},
				{
					Entity: types.Entity{
						ID:      "serv-2",
						Enabled: true,
					},
					EntityPatternFields: savedpattern.EntityPatternFields{
						EntityPattern: [][]pattern.FieldCondition{
							{
								{
									Field:     "component",
									Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
								},
							},
						},
					},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:               "id-1",
					Component:        "component-1",
					Enabled:          true,
					ServicesToAdd:    []string{"serv-1", "serv-4"},
					ServicesToRemove: []string{"serv-0", "serv-3"},
					Services:         []string{"serv-1", "serv-2"},
				},
			},
		},
		{
			name: "one entity is removed from a single service",
			entities: []types.Entity{
				{
					ID:        "id-1",
					Component: "component-1",
					Enabled:   true,
					Services:  []string{"serv-1"},
				},
			},
			services: []entityservice.EntityService{
				{
					Entity: types.Entity{
						ID:      "serv-1",
						Enabled: true,
					},
					EntityPatternFields: savedpattern.EntityPatternFields{
						EntityPattern: [][]pattern.FieldCondition{
							{
								{
									Field:     "component",
									Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-2"),
								},
							},
						},
					},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:               "id-1",
					Component:        "component-1",
					Enabled:          true,
					Services:         []string{},
					ServicesToAdd:    []string{},
					ServicesToRemove: []string{"serv-1"},
				},
			},
		},
		{
			name: "one entity is removed from a single service but have this service in ServicesToAdd",
			entities: []types.Entity{
				{
					ID:            "id-1",
					Component:     "component-1",
					Enabled:       true,
					Services:      []string{"serv-1"},
					ServicesToAdd: []string{"serv-1"},
				},
			},
			services: []entityservice.EntityService{
				{
					Entity: types.Entity{
						ID:      "serv-1",
						Enabled: true,
					},
					EntityPatternFields: savedpattern.EntityPatternFields{
						EntityPattern: [][]pattern.FieldCondition{
							{
								{
									Field:     "component",
									Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-2"),
								},
							},
						},
					},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:               "id-1",
					Component:        "component-1",
					Enabled:          true,
					Services:         []string{},
					ServicesToRemove: []string{},
				},
			},
		},
		{
			name: "one entity is removed from multiple services",
			entities: []types.Entity{
				{
					ID:        "id-1",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-1", "serv-2"},
				},
			},
			services: []entityservice.EntityService{
				{
					Entity: types.Entity{
						ID:      "serv-1",
						Enabled: true,
					},
					EntityPatternFields: savedpattern.EntityPatternFields{
						EntityPattern: [][]pattern.FieldCondition{
							{
								{
									Field:     "component",
									Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-2"),
								},
							},
						},
					},
				},
				{
					Entity: types.Entity{
						ID:      "serv-2",
						Enabled: true,
					},
					EntityPatternFields: savedpattern.EntityPatternFields{
						EntityPattern: [][]pattern.FieldCondition{
							{
								{
									Field:     "component",
									Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-2"),
								},
							},
						},
					},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:               "id-1",
					Component:        "component-1",
					Enabled:          true,
					Services:         []string{},
					ServicesToAdd:    []string{},
					ServicesToRemove: []string{"serv-1", "serv-2"},
				},
			},
		},
		{
			name: "one entity is moved from one service to another",
			entities: []types.Entity{
				{
					ID:        "id-1",
					Component: "component-1",
					Enabled:   true,
					Services:  []string{"serv-1", "serv-2"},
				},
			},
			services: []entityservice.EntityService{
				{
					Entity: types.Entity{
						ID:      "serv-1",
						Enabled: true,
					},
					EntityPatternFields: savedpattern.EntityPatternFields{
						EntityPattern: [][]pattern.FieldCondition{
							{
								{
									Field:     "component",
									Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
								},
							},
						},
					},
				},
				{
					Entity: types.Entity{
						ID:      "serv-2",
						Enabled: true,
					},
					EntityPatternFields: savedpattern.EntityPatternFields{
						EntityPattern: [][]pattern.FieldCondition{
							{
								{
									Field:     "component",
									Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-2"),
								},
							},
						},
					},
				},
				{
					Entity: types.Entity{
						ID:      "serv-3",
						Enabled: true,
					},
					EntityPatternFields: savedpattern.EntityPatternFields{
						EntityPattern: [][]pattern.FieldCondition{
							{
								{
									Field:     "component",
									Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
								},
							},
						},
					},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:               "id-1",
					Component:        "component-1",
					Enabled:          true,
					Services:         []string{"serv-1", "serv-3"},
					ServicesToAdd:    []string{"serv-3"},
					ServicesToRemove: []string{"serv-2"},
				},
			},
		},
		{
			name: "multiple firstFindCallEntities is added to a single service",
			entities: []types.Entity{
				{
					ID:        "id-1",
					Component: "component-1",
					Enabled:   true,
				},
				{
					ID:        "id-2",
					Component: "component-1",
					Enabled:   true,
				},
			},
			services: []entityservice.EntityService{
				{
					Entity: types.Entity{
						ID:      "serv-1",
						Enabled: true,
					},
					EntityPatternFields: savedpattern.EntityPatternFields{
						EntityPattern: [][]pattern.FieldCondition{
							{
								{
									Field:     "component",
									Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
								},
							},
						},
					},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:               "id-1",
					Component:        "component-1",
					Enabled:          true,
					ServicesToAdd:    []string{"serv-1"},
					ServicesToRemove: []string{},
					Services:         []string{"serv-1"},
				},
				{
					ID:               "id-2",
					Component:        "component-1",
					Enabled:          true,
					ServicesToAdd:    []string{"serv-1"},
					ServicesToRemove: []string{},
					Services:         []string{"serv-1"},
				},
			},
		},
		{
			name: "multiple firstFindCallEntities is added to multiple services",
			entities: []types.Entity{
				{
					ID:        "id-1",
					Component: "component-1",
					Enabled:   true,
				},
				{
					ID:        "id-2",
					Component: "component-1",
					Enabled:   true,
				},
			},
			services: []entityservice.EntityService{
				{
					Entity: types.Entity{
						ID:      "serv-1",
						Enabled: true,
					},
					EntityPatternFields: savedpattern.EntityPatternFields{
						EntityPattern: [][]pattern.FieldCondition{
							{
								{
									Field:     "component",
									Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
								},
							},
						},
					},
				},
				{
					Entity: types.Entity{
						ID:      "serv-2",
						Enabled: true,
					},
					EntityPatternFields: savedpattern.EntityPatternFields{
						EntityPattern: [][]pattern.FieldCondition{
							{
								{
									Field:     "component",
									Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
								},
							},
						},
					},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:               "id-1",
					Component:        "component-1",
					Enabled:          true,
					ServicesToAdd:    []string{"serv-1", "serv-2"},
					ServicesToRemove: []string{},
					Services:         []string{"serv-1", "serv-2"},
				},
				{
					ID:               "id-2",
					Component:        "component-1",
					Enabled:          true,
					ServicesToAdd:    []string{"serv-1", "serv-2"},
					ServicesToRemove: []string{},
					Services:         []string{"serv-1", "serv-2"},
				},
			},
		},
		{
			name: "multiple firstFindCallEntities is removed from a single service",
			entities: []types.Entity{
				{
					ID:        "id-1",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-1"},
				},
				{
					ID:        "id-2",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-1"},
				},
			},
			services: []entityservice.EntityService{
				{
					Entity: types.Entity{
						ID:      "serv-1",
						Enabled: true,
					},
					EntityPatternFields: savedpattern.EntityPatternFields{
						EntityPattern: [][]pattern.FieldCondition{
							{
								{
									Field:     "component",
									Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-2"),
								},
							},
						},
					},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:               "id-1",
					Enabled:          true,
					Component:        "component-1",
					Services:         []string{},
					ServicesToAdd:    []string{},
					ServicesToRemove: []string{"serv-1"},
				},
				{
					ID:               "id-2",
					Enabled:          true,
					Component:        "component-1",
					Services:         []string{},
					ServicesToAdd:    []string{},
					ServicesToRemove: []string{"serv-1"},
				},
			},
		},
		{
			name: "multiple firstFindCallEntities is removed from multiple services",
			entities: []types.Entity{
				{
					ID:        "id-1",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-1", "serv-2"},
				},
				{
					ID:        "id-2",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-1", "serv-2"},
				},
			},
			services: []entityservice.EntityService{
				{
					Entity: types.Entity{
						ID:      "serv-1",
						Enabled: true,
					},
					EntityPatternFields: savedpattern.EntityPatternFields{
						EntityPattern: [][]pattern.FieldCondition{
							{
								{
									Field:     "component",
									Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-2"),
								},
							},
						},
					},
				},
				{
					Entity: types.Entity{
						ID:      "serv-2",
						Enabled: true,
					},
					EntityPatternFields: savedpattern.EntityPatternFields{
						EntityPattern: [][]pattern.FieldCondition{
							{
								{
									Field:     "component",
									Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-2"),
								},
							},
						},
					},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:               "id-1",
					Enabled:          true,
					Component:        "component-1",
					Services:         []string{},
					ServicesToAdd:    []string{},
					ServicesToRemove: []string{"serv-1", "serv-2"},
				},
				{
					ID:               "id-2",
					Enabled:          true,
					Component:        "component-1",
					Services:         []string{},
					ServicesToAdd:    []string{},
					ServicesToRemove: []string{"serv-1", "serv-2"},
				},
			},
		},
		{
			name: "multiple firstFindCallEntities is moved from one service to another",
			entities: []types.Entity{
				{
					ID:        "id-1",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-1", "serv-2"},
				},
				{
					ID:        "id-2",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-2", "serv-3"},
				},
			},
			services: []entityservice.EntityService{
				{
					Entity: types.Entity{
						ID:      "serv-1",
						Enabled: true,
					},
					EntityPatternFields: savedpattern.EntityPatternFields{
						EntityPattern: [][]pattern.FieldCondition{
							{
								{
									Field:     "component",
									Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
								},
							},
						},
					},
				},
				{
					Entity: types.Entity{
						ID:      "serv-2",
						Enabled: true,
					},
					EntityPatternFields: savedpattern.EntityPatternFields{
						EntityPattern: [][]pattern.FieldCondition{
							{
								{
									Field:     "component",
									Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-2"),
								},
							},
						},
					},
				},
				{
					Entity: types.Entity{
						ID:      "serv-3",
						Enabled: true,
					},
					EntityPatternFields: savedpattern.EntityPatternFields{
						EntityPattern: [][]pattern.FieldCondition{
							{
								{
									Field:     "component",
									Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
								},
							},
						},
					},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:               "id-1",
					Enabled:          true,
					Component:        "component-1",
					Services:         []string{"serv-1", "serv-3"},
					ServicesToAdd:    []string{"serv-3"},
					ServicesToRemove: []string{"serv-2"},
				},
				{
					ID:               "id-2",
					Enabled:          true,
					Component:        "component-1",
					Services:         []string{"serv-1", "serv-3"},
					ServicesToAdd:    []string{"serv-1"},
					ServicesToRemove: []string{"serv-2"},
				},
			},
		},
		{
			name: "no changes",
			entities: []types.Entity{
				{
					ID:        "id-1",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-1", "serv-2"},
				},
				{
					ID:        "id-2",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-1", "serv-2"},
				},
			},
			services: []entityservice.EntityService{
				{
					Entity: types.Entity{
						ID:      "serv-1",
						Enabled: true,
					},
					EntityPatternFields: savedpattern.EntityPatternFields{
						EntityPattern: [][]pattern.FieldCondition{
							{
								{
									Field:     "component",
									Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
								},
							},
						},
					},
				},
				{
					Entity: types.Entity{
						ID:      "serv-2",
						Enabled: true,
					},
					EntityPatternFields: savedpattern.EntityPatternFields{
						EntityPattern: [][]pattern.FieldCondition{
							{
								{
									Field:     "component",
									Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
								},
							},
						},
					},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:        "id-1",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-1", "serv-2"},
				},
				{
					ID:        "id-2",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-1", "serv-2"},
				},
			},
		},
	}

	manager := contextgraph.NewManager(adapter, dbClient, storage, metaUpdater, log.NewLogger(true))

	for _, dataset := range dataSets {
		t.Run(dataset.name, func(t *testing.T) {
			storage.EXPECT().GetAll(gomock.Any()).Return(dataset.services, nil)

			result, err := manager.CheckServices(ctx, dataset.entities)
			if err != nil {
				t.Error(err)
			}

			sort.Slice(result, func(i, j int) bool {
				return result[i].ID < result[j].ID
			})

			for idx := range result {
				sort.Slice(result[idx].Services, func(i, j int) bool {
					return result[idx].Services[i] < result[idx].Services[j]
				})

				sort.Slice(result[idx].ServicesToAdd, func(i, j int) bool {
					return result[idx].ServicesToAdd[i] < result[idx].ServicesToAdd[j]
				})

				sort.Slice(result[idx].ServicesToRemove, func(i, j int) bool {
					return result[idx].ServicesToRemove[i] < result[idx].ServicesToRemove[j]
				})
			}

			for idx := range dataset.expectedResult {
				sort.Slice(dataset.expectedResult[idx].Services, func(i, j int) bool {
					return dataset.expectedResult[idx].Services[i] < dataset.expectedResult[idx].Services[j]
				})

				sort.Slice(dataset.expectedResult[idx].ServicesToAdd, func(i, j int) bool {
					return dataset.expectedResult[idx].ServicesToAdd[i] < dataset.expectedResult[idx].ServicesToAdd[j]
				})

				sort.Slice(dataset.expectedResult[idx].ServicesToRemove, func(i, j int) bool {
					return dataset.expectedResult[idx].ServicesToRemove[i] < dataset.expectedResult[idx].ServicesToRemove[j]
				})
			}

			sort.Slice(dataset.expectedResult, func(i, j int) bool {
				return dataset.expectedResult[i].ID < dataset.expectedResult[j].ID
			})

			for idx := 0; idx < len(result); idx++ {
				if slices.Compare(result[idx].Services, dataset.expectedResult[idx].Services) != 0 {
					t.Errorf("expected Services to be %v, but got %v", dataset.expectedResult[idx].Services, result[idx].Services)
				}

				if slices.Compare(result[idx].ServicesToAdd, dataset.expectedResult[idx].ServicesToAdd) != 0 {
					t.Errorf("expected ServicesToAdd to be %v, but got %v", dataset.expectedResult[idx].ServicesToAdd, result[idx].ServicesToAdd)
				}

				if slices.Compare(result[idx].ServicesToRemove, dataset.expectedResult[idx].ServicesToRemove) != 0 {
					t.Errorf("expected ServicesToRemove to be %v, but got %v", dataset.expectedResult[idx].ServicesToRemove, result[idx].ServicesToRemove)
				}
			}
		})
	}
}

func BenchmarkCenterCheckServices(b *testing.B) {
	ctx := context.Background()
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	collection := mock_mongo.NewMockDbCollection(ctrl)

	dbClient := mock_mongo.NewMockDbClient(ctrl)
	dbClient.EXPECT().Collection(mongo.EntityMongoCollection).Return(collection).AnyTimes()

	adapter := entity.NewAdapter(dbClient)
	storage := mock_contextgraph.NewMockEntityServiceStorage(ctrl)

	metaUpdater := mock_metrics.NewMockMetaUpdater(ctrl)

	services := make([]entityservice.EntityService, 1000)
	for i := 0; i < 1000; i++ {
		services[i] = entityservice.EntityService{
			Entity: types.Entity{ID: fmt.Sprintf("serv-%d", i), Enabled: true},
			EntityPatternFields: savedpattern.EntityPatternFields{
				EntityPattern: [][]pattern.FieldCondition{
					{
						{
							Field:     "component",
							Condition: pattern.NewStringCondition(pattern.ConditionEqual, fmt.Sprintf("component-%d", i)),
						},
					},
				},
			},
		}
	}

	entities := make([]types.Entity, 100)
	for i := 0; i < 100; i++ {
		Services := make([]string, 50)
		for j := 0; j < 50; j++ {
			Services[j] = fmt.Sprintf("serv-%d", rand.Intn(1000))
		}

		entities[i] = types.Entity{
			ID:        fmt.Sprintf("id-%d", i),
			Enabled:   true,
			Component: fmt.Sprintf("component-%d", i),
			Services:  Services,
		}
	}

	storage.EXPECT().GetAll(gomock.Any()).Return(services, nil).AnyTimes()
	manager := contextgraph.NewManager(adapter, dbClient, storage, metaUpdater, log.NewLogger(true))

	for i := 0; i < b.N; i++ {
		_, _ = manager.CheckServices(ctx, entities)
	}
}

func TestRecomputeService(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	collection := mock_mongo.NewMockDbCollection(ctrl)

	dbClient := mock_mongo.NewMockDbClient(ctrl)
	dbClient.EXPECT().Collection(mongo.EntityMongoCollection).Return(collection).AnyTimes()

	adapter := entity.NewAdapter(dbClient)
	storage := mock_contextgraph.NewMockEntityServiceStorage(ctrl)

	metaUpdater := mock_metrics.NewMockMetaUpdater(ctrl)

	dataSets := []struct {
		service                entityservice.EntityService
		firstFindCallEntities  []types.Entity
		secondFindCallEntities []types.Entity
		expectedResult         []types.Entity
		serviceName            string
		name                   string
	}{
		{
			name:        "service deleted",
			service:     entityservice.EntityService{Entity: types.Entity{ID: "service-id"}},
			serviceName: "serv-1",
			firstFindCallEntities: []types.Entity{
				{
					ID:        "id-1",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-1"},
				},
				{
					ID:        "id-2",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-0", "serv-1", "serv-2"},
				},
				{
					ID:            "id-3",
					Enabled:       true,
					Component:     "component-1",
					Services:      []string{"serv-0", "serv-1", "serv-2"},
					ServicesToAdd: []string{"serv-1"},
				},
			},
			expectedResult: []types.Entity{
				{ID: "service-id"},
				{
					ID:        "id-1",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{},
				},
				{
					ID:        "id-2",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-0", "serv-2"},
				},
				{
					ID:            "id-3",
					Enabled:       true,
					Component:     "component-1",
					Services:      []string{"serv-0", "serv-2"},
					ServicesToAdd: []string{},
				},
			},
		},
		{
			name: "service is disabled",
			service: entityservice.EntityService{
				Entity: types.Entity{
					ID: "serv-1",
				},
				EntityPatternFields: savedpattern.EntityPatternFields{
					EntityPattern: [][]pattern.FieldCondition{
						{
							{
								Field:     "component",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
							},
						},
					},
				},
			},
			serviceName: "serv-1",
			firstFindCallEntities: []types.Entity{
				{
					ID:        "id-1",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-1"},
				},
				{
					ID:        "id-2",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-0", "serv-1", "serv-2"},
				},
				{
					ID:            "id-3",
					Enabled:       true,
					Component:     "component-1",
					Services:      []string{"serv-0", "serv-1", "serv-2"},
					ServicesToAdd: []string{"serv-1"},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:        "id-1",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{},
				},
				{
					ID:        "id-2",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-0", "serv-2"},
				},
				{
					ID:            "id-3",
					Enabled:       true,
					Component:     "component-1",
					Services:      []string{"serv-0", "serv-2"},
					ServicesToAdd: []string{},
				},
				{
					ID: "serv-1",
				},
			},
		},
		{
			name: "New service",
			service: entityservice.EntityService{
				Entity: types.Entity{
					ID:      "serv-1",
					Enabled: true,
				},
				EntityPatternFields: savedpattern.EntityPatternFields{
					EntityPattern: [][]pattern.FieldCondition{
						{
							{
								Field:     "component",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
							},
						},
					},
				},
			},
			serviceName:           "serv-1",
			firstFindCallEntities: []types.Entity{},
			secondFindCallEntities: []types.Entity{
				{
					ID:        "id-1",
					Enabled:   true,
					Component: "component-1",
				},
				{
					ID:        "id-2",
					Enabled:   true,
					Component: "component-1",
				},
				{
					ID:        "id-3",
					Enabled:   true,
					Component: "component-1",
				},
			},
			expectedResult: []types.Entity{
				{
					ID:        "id-1",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-1"},
				},
				{
					ID:        "id-2",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-1"},
				},
				{
					ID:        "id-3",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-1"},
				},
				{
					ID:      "serv-1",
					Enabled: true,
				},
			},
		},
		{
			name: "New service(many impacted services)",
			service: entityservice.EntityService{
				Entity: types.Entity{
					ID:      "serv-1",
					Enabled: true,
				},
				EntityPatternFields: savedpattern.EntityPatternFields{
					EntityPattern: [][]pattern.FieldCondition{
						{
							{
								Field:     "component",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
							},
						},
					},
				},
			},
			serviceName:           "serv-1",
			firstFindCallEntities: []types.Entity{},
			secondFindCallEntities: []types.Entity{
				{
					ID:        "id-1",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-0", "serv-2"},
				},
				{
					ID:        "id-2",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-0", "serv-2"},
				},
				{
					ID:        "id-3",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-0", "serv-2"},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:        "id-1",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-1", "serv-0", "serv-2"},
				},
				{
					ID:        "id-2",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-1", "serv-0", "serv-2"},
				},
				{
					ID:        "id-3",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-1", "serv-0", "serv-2"},
				},
				{
					ID:      "serv-1",
					Enabled: true,
				},
			},
		},
		{
			name: "Update service, where some depends are not valid anymore",
			service: entityservice.EntityService{
				Entity: types.Entity{
					ID:      "serv-1",
					Enabled: true,
				},
				EntityPatternFields: savedpattern.EntityPatternFields{
					EntityPattern: [][]pattern.FieldCondition{
						{
							{
								Field:     "component",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
							},
						},
					},
				},
			},
			serviceName: "serv-1",
			firstFindCallEntities: []types.Entity{
				{
					ID:        "id-2",
					Enabled:   true,
					Component: "component-2",
					Services:  []string{"serv-1"},
				},
				{
					ID:        "id-3",
					Enabled:   true,
					Component: "component-2",
					Services:  []string{"serv-1"},
				},
			},
			secondFindCallEntities: []types.Entity{
				{
					ID:        "id-4",
					Enabled:   true,
					Component: "component-1",
				},
			},
			expectedResult: []types.Entity{
				{
					ID:        "id-2",
					Enabled:   true,
					Component: "component-2",
					Services:  []string{},
				},
				{
					ID:        "id-3",
					Enabled:   true,
					Component: "component-2",
					Services:  []string{},
				},
				{
					ID:        "id-4",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-1"},
				},
				{
					ID:      "serv-1",
					Enabled: true,
				},
			},
		},
		{
			name: "Update service, where some depends are not valid anymore(many impacted services)",
			service: entityservice.EntityService{
				Entity: types.Entity{
					ID:      "serv-1",
					Enabled: true,
				},
				EntityPatternFields: savedpattern.EntityPatternFields{
					EntityPattern: [][]pattern.FieldCondition{
						{
							{
								Field:     "component",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
							},
						},
					},
				},
			},
			serviceName: "serv-1",
			firstFindCallEntities: []types.Entity{
				{
					ID:        "id-2",
					Enabled:   true,
					Component: "component-2",
					Services:  []string{"serv-0", "serv-1", "serv-2"},
				},
				{
					ID:        "id-3",
					Enabled:   true,
					Component: "component-2",
					Services:  []string{"serv-0", "serv-1", "serv-2"},
				},
			},
			secondFindCallEntities: []types.Entity{
				{
					ID:        "id-4",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-0", "serv-2"},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:        "id-2",
					Enabled:   true,
					Component: "component-2",
					Services:  []string{"serv-0", "serv-2"},
				},
				{
					ID:        "id-3",
					Enabled:   true,
					Component: "component-2",
					Services:  []string{"serv-0", "serv-2"},
				},
				{
					ID:        "id-4",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-0", "serv-1", "serv-2"},
				},
				{
					ID:      "serv-1",
					Enabled: true,
				},
			},
		},
		{
			name: "Update service where addedTo is exist, if not valid anymore do not add to RemoveFrom",
			service: entityservice.EntityService{
				Entity: types.Entity{
					ID:      "serv-1",
					Enabled: true,
				},
				EntityPatternFields: savedpattern.EntityPatternFields{
					EntityPattern: [][]pattern.FieldCondition{
						{
							{
								Field:     "component",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
							},
						},
					},
				},
			},
			serviceName: "serv-1",
			firstFindCallEntities: []types.Entity{
				{
					ID:            "id-2",
					Enabled:       true,
					Component:     "component-2",
					Services:      []string{"serv-1"},
					ServicesToAdd: []string{"serv-1"},
				},
				{
					ID:            "id-3",
					Enabled:       true,
					Component:     "component-2",
					Services:      []string{"serv-1"},
					ServicesToAdd: []string{"serv-1"},
				},
			},
			secondFindCallEntities: []types.Entity{
				{
					ID:        "id-4",
					Enabled:   true,
					Component: "component-1",
				},
			},
			expectedResult: []types.Entity{
				{
					ID:            "id-2",
					Enabled:       true,
					Component:     "component-2",
					Services:      []string{},
					ServicesToAdd: []string{},
				},
				{
					ID:            "id-3",
					Enabled:       true,
					Component:     "component-2",
					Services:      []string{},
					ServicesToAdd: []string{},
				},
				{
					ID:        "id-4",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-1"},
				},
				{
					ID:      "serv-1",
					Enabled: true,
				},
			},
		},
		{
			name: "Update service where addedTo is exist, if not valid anymore do not add to RemoveFrom(many impacted services)",
			service: entityservice.EntityService{
				Entity: types.Entity{
					ID:      "serv-1",
					Enabled: true,
				},
				EntityPatternFields: savedpattern.EntityPatternFields{
					EntityPattern: [][]pattern.FieldCondition{
						{
							{
								Field:     "component",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
							},
						},
					},
				},
			},
			serviceName: "serv-1",
			firstFindCallEntities: []types.Entity{
				{
					ID:            "id-2",
					Enabled:       true,
					Component:     "component-2",
					Services:      []string{"serv-0", "serv-1", "serv-2"},
					ServicesToAdd: []string{"serv-0", "serv-1", "serv-2"},
				},
				{
					ID:            "id-3",
					Enabled:       true,
					Component:     "component-2",
					Services:      []string{"serv-0", "serv-1", "serv-2", "serv-3"},
					ServicesToAdd: []string{"serv-0", "serv-1", "serv-2"},
				},
			},
			secondFindCallEntities: []types.Entity{
				{
					ID:        "id-4",
					Enabled:   true,
					Component: "component-1",
				},
			},
			expectedResult: []types.Entity{
				{
					ID:            "id-2",
					Enabled:       true,
					Component:     "component-2",
					Services:      []string{"serv-0", "serv-2"},
					ServicesToAdd: []string{"serv-0", "serv-2"},
				},
				{
					ID:            "id-3",
					Enabled:       true,
					Component:     "component-2",
					Services:      []string{"serv-0", "serv-2", "serv-3"},
					ServicesToAdd: []string{"serv-0", "serv-2"},
				},
				{
					ID:        "id-4",
					Enabled:   true,
					Component: "component-1",
					Services:  []string{"serv-1"},
				},
				{
					ID:      "serv-1",
					Enabled: true,
				},
			},
		},
		{
			name: "Update service where RemoveTo is exist, if not valid anymore do not add to AddedTo",
			service: entityservice.EntityService{
				Entity: types.Entity{
					ID:      "serv-1",
					Enabled: true,
				},
				EntityPatternFields: savedpattern.EntityPatternFields{
					EntityPattern: [][]pattern.FieldCondition{
						{
							{
								Field:     "component",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
							},
						},
						{
							{
								Field:     "component",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-2"),
							},
						},
					},
				},
			},
			serviceName:           "serv-1",
			firstFindCallEntities: []types.Entity{},
			secondFindCallEntities: []types.Entity{
				{
					ID:               "id-2",
					Enabled:          true,
					Component:        "component-2",
					Services:         []string{},
					ServicesToRemove: []string{"serv-1"},
				},
				{
					ID:               "id-3",
					Enabled:          true,
					Component:        "component-2",
					Services:         []string{},
					ServicesToRemove: []string{"serv-1"},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:               "id-2",
					Enabled:          true,
					Component:        "component-2",
					Services:         []string{"serv-1"},
					ServicesToRemove: []string{},
				},
				{
					ID:               "id-3",
					Enabled:          true,
					Component:        "component-2",
					Services:         []string{"serv-1"},
					ServicesToRemove: []string{},
				},
				{
					ID:      "serv-1",
					Enabled: true,
				},
			},
		},
		{
			name: "Update service where RemoveTo is exist, if not valid anymore do not add to AddedTo(many services)",
			service: entityservice.EntityService{
				Entity: types.Entity{
					ID:      "serv-1",
					Enabled: true,
				},
				EntityPatternFields: savedpattern.EntityPatternFields{
					EntityPattern: [][]pattern.FieldCondition{
						{
							{
								Field:     "component",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
							},
						},
						{
							{
								Field:     "component",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-2"),
							},
						},
					},
				},
			},
			serviceName:           "serv-1",
			firstFindCallEntities: []types.Entity{},
			secondFindCallEntities: []types.Entity{
				{
					ID:               "id-2",
					Enabled:          true,
					Component:        "component-2",
					Services:         []string{"serv-3"},
					ServicesToRemove: []string{"serv-0", "serv-1", "serv-2"},
				},
				{
					ID:               "id-3",
					Enabled:          true,
					Component:        "component-2",
					Services:         []string{"serv-3", "serv-4"},
					ServicesToRemove: []string{"serv-0", "serv-1", "serv-2"},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:               "id-2",
					Enabled:          true,
					Component:        "component-2",
					Services:         []string{"serv-1", "serv-3"},
					ServicesToRemove: []string{"serv-0", "serv-2"},
				},
				{
					ID:               "id-3",
					Enabled:          true,
					Component:        "component-2",
					Services:         []string{"serv-1", "serv-3", "serv-4"},
					ServicesToRemove: []string{"serv-0", "serv-2"},
				},
				{
					ID:      "serv-1",
					Enabled: true,
				},
			},
		},
	}

	call := 0
	manager := contextgraph.NewManager(adapter, dbClient, storage, metaUpdater, log.NewLogger(true))
	for _, dataset := range dataSets {
		t.Run(dataset.name, func(t *testing.T) {
			storage.EXPECT().GetAll(gomock.Any()).Return(nil, nil).AnyTimes()
			storage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(dataset.service, nil)

			call = 0

			cursor := mock_mongo.NewMockCursor(ctrl)
			cursor.EXPECT().All(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, results interface{}) {
				if call == 0 {
					entities := results.(*[]types.Entity)
					*entities = dataset.firstFindCallEntities
				}

				if call == 1 {
					entities := results.(*[]types.Entity)
					*entities = dataset.secondFindCallEntities
				}

				call++
			}).Return(nil).AnyTimes()
			collection.EXPECT().Find(gomock.Any(), gomock.Any()).Return(cursor, nil).AnyTimes()

			_, result, err := manager.RecomputeService(ctx, dataset.serviceName)
			if err != nil {
				t.Error(err)
			}

			sort.Slice(result, func(i, j int) bool {
				return result[i].ID < result[j].ID
			})

			for idx := range result {
				sort.Slice(result[idx].Services, func(i, j int) bool {
					return result[idx].Services[i] < result[idx].Services[j]
				})

				sort.Slice(result[idx].ServicesToAdd, func(i, j int) bool {
					return result[idx].ServicesToAdd[i] < result[idx].ServicesToAdd[j]
				})

				sort.Slice(result[idx].ServicesToRemove, func(i, j int) bool {
					return result[idx].ServicesToRemove[i] < result[idx].ServicesToRemove[j]
				})
			}

			for idx := range dataset.expectedResult {
				sort.Slice(dataset.expectedResult[idx].Services, func(i, j int) bool {
					return dataset.expectedResult[idx].Services[i] < dataset.expectedResult[idx].Services[j]
				})

				sort.Slice(dataset.expectedResult[idx].ServicesToAdd, func(i, j int) bool {
					return dataset.expectedResult[idx].ServicesToAdd[i] < dataset.expectedResult[idx].ServicesToAdd[j]
				})

				sort.Slice(dataset.expectedResult[idx].ServicesToRemove, func(i, j int) bool {
					return dataset.expectedResult[idx].ServicesToRemove[i] < dataset.expectedResult[idx].ServicesToRemove[j]
				})
			}

			sort.Slice(dataset.expectedResult, func(i, j int) bool {
				return dataset.expectedResult[i].ID < dataset.expectedResult[j].ID
			})

			if diff := pretty.Compare(result, dataset.expectedResult); diff != "" {
				t.Errorf("result is not expected: %s", diff)
			}
		})
	}
}

func BenchmarkRecomputeServicesRemoveAll(b *testing.B) {
	ctx := context.Background()
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	var entities []types.Entity
	for i := 0; i < 1000; i++ {
		eID := fmt.Sprintf("id-%d", i)
		entities = append(entities, types.Entity{
			ID:        eID,
			Enabled:   true,
			Component: "component-1",
			Services:  []string{"serv-1"},
		})
	}

	cursor := mock_mongo.NewMockCursor(ctrl)
	cursor.EXPECT().All(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, results interface{}) {
		ents := results.(*[]types.Entity)
		*ents = append(*ents, entities...)
	}).Return(nil).AnyTimes()

	collection := mock_mongo.NewMockDbCollection(ctrl)
	collection.EXPECT().Find(gomock.Any(), gomock.Any()).Return(cursor, nil).AnyTimes()

	dbClient := mock_mongo.NewMockDbClient(ctrl)
	dbClient.EXPECT().Collection(mongo.EntityMongoCollection).Return(collection).AnyTimes()

	adapter := entity.NewAdapter(dbClient)
	storage := mock_contextgraph.NewMockEntityServiceStorage(ctrl)
	storage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(entityservice.EntityService{
		Entity: types.Entity{
			ID:      "serv-1",
			Enabled: true,
		},
		EntityPatternFields: savedpattern.EntityPatternFields{
			EntityPattern: [][]pattern.FieldCondition{
				{
					{
						Field:     "component",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
					},
				},
			},
		},
	}, nil).AnyTimes()

	metaUpdater := mock_metrics.NewMockMetaUpdater(ctrl)

	manager := contextgraph.NewManager(adapter, dbClient, storage, metaUpdater, log.NewLogger(true))
	for i := 0; i < b.N; i++ {
		_, _, _ = manager.RecomputeService(ctx, "serv-1")
	}
}

func BenchmarkRecomputeServicesAddAll(b *testing.B) {
	ctx := context.Background()
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	var entities []types.Entity
	for i := 0; i < 1000; i++ {
		eID := fmt.Sprintf("id-%d", i)
		entities = append(entities, types.Entity{
			ID:        eID,
			Enabled:   true,
			Component: "component-1",
		})
	}

	call := 0
	cursor := mock_mongo.NewMockCursor(ctrl)
	cursor.EXPECT().All(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, results interface{}) {
		if call == 1 {
			ents := results.(*[]types.Entity)
			*ents = append(*ents, entities...)
		}

		call++
	}).Return(nil).AnyTimes()

	collection := mock_mongo.NewMockDbCollection(ctrl)
	collection.EXPECT().Find(gomock.Any(), gomock.Any()).Return(cursor, nil).AnyTimes()

	dbClient := mock_mongo.NewMockDbClient(ctrl)
	dbClient.EXPECT().Collection(mongo.EntityMongoCollection).Return(collection).AnyTimes()

	adapter := entity.NewAdapter(dbClient)
	storage := mock_contextgraph.NewMockEntityServiceStorage(ctrl)
	storage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(entityservice.EntityService{
		Entity: types.Entity{
			ID:      "serv-1",
			Enabled: true,
		},
		EntityPatternFields: savedpattern.EntityPatternFields{
			EntityPattern: [][]pattern.FieldCondition{
				{
					{
						Field:     "component",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
					},
				},
			},
		},
	}, nil).AnyTimes()

	metaUpdater := mock_metrics.NewMockMetaUpdater(ctrl)

	manager := contextgraph.NewManager(adapter, dbClient, storage, metaUpdater, log.NewLogger(true))
	for i := 0; i < b.N; i++ {
		call = 0
		_, _, _ = manager.RecomputeService(ctx, "serv-1")
	}
}

func BenchmarkRecomputeServicesMixed(b *testing.B) {
	ctx := context.Background()
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	var entitiesToRemove []types.Entity
	for i := 0; i < 500; i++ {
		eID := fmt.Sprintf("id-%d", i)
		entitiesToRemove = append(entitiesToRemove, types.Entity{
			ID:        eID,
			Enabled:   true,
			Component: "component-1",
			Services:  []string{"serv-1"},
		})
	}

	var entitiesToAdd []types.Entity
	for i := 500; i < 1000; i++ {
		eID := fmt.Sprintf("id-%d", i)
		entitiesToAdd = append(entitiesToAdd, types.Entity{
			ID:        eID,
			Enabled:   true,
			Component: "component-1",
		})
	}

	call := 0
	cursor := mock_mongo.NewMockCursor(ctrl)
	cursor.EXPECT().All(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, results interface{}) {
		if call == 0 {
			ents := results.(*[]types.Entity)
			*ents = append(*ents, entitiesToRemove...)
		}

		if call == 1 {
			ents := results.(*[]types.Entity)
			*ents = append(*ents, entitiesToAdd...)
		}

		call++
	}).Return(nil).AnyTimes()

	collection := mock_mongo.NewMockDbCollection(ctrl)
	collection.EXPECT().Find(gomock.Any(), gomock.Any()).Return(cursor, nil).AnyTimes()

	dbClient := mock_mongo.NewMockDbClient(ctrl)
	dbClient.EXPECT().Collection(mongo.EntityMongoCollection).Return(collection).AnyTimes()

	adapter := entity.NewAdapter(dbClient)
	storage := mock_contextgraph.NewMockEntityServiceStorage(ctrl)
	storage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(entityservice.EntityService{
		Entity: types.Entity{
			ID:      "serv-1",
			Enabled: true,
		},
		EntityPatternFields: savedpattern.EntityPatternFields{
			EntityPattern: [][]pattern.FieldCondition{
				{
					{
						Field:     "component",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
					},
				},
			},
		},
	}, nil).AnyTimes()

	metaUpdater := mock_metrics.NewMockMetaUpdater(ctrl)

	manager := contextgraph.NewManager(adapter, dbClient, storage, metaUpdater, log.NewLogger(true))
	for i := 0; i < b.N; i++ {
		call = 0
		_, _, _ = manager.RecomputeService(ctx, "serv-1")
	}
}
