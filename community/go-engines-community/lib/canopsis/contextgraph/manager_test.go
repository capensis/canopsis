package contextgraph_test

import (
	"context"
	"fmt"
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
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	mock_statesetting "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/statesetting"
	"github.com/golang/mock/gomock"
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

	assigner := mock_statesetting.NewMockAssigner(ctrl)

	dataSets := []struct {
		services       []entityservice.EntityService
		entity         types.Entity
		expectedEntity types.Entity
		name           string
	}{
		{
			name: "one entity is added to a single service",
			entity: types.Entity{
				ID:        "id-1",
				Component: "component-1",
				Enabled:   true,
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
			expectedEntity: types.Entity{
				ID:               "id-1",
				Component:        "component-1",
				Enabled:          true,
				ServicesToAdd:    []string{"serv-1"},
				ServicesToRemove: []string{},
				Services:         []string{"serv-1"},
			},
		},
		{
			name: "one entity is added to multiple services",
			entity: types.Entity{
				ID:        "id-1",
				Component: "component-1",
				Enabled:   true,
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
			expectedEntity: types.Entity{
				ID:               "id-1",
				Component:        "component-1",
				Enabled:          true,
				ServicesToAdd:    []string{"serv-1", "serv-2"},
				ServicesToRemove: []string{},
				Services:         []string{"serv-1", "serv-2"},
			},
		},
		{
			name: "one entity is added to multiple services impacted services to add/remove should be updated",
			entity: types.Entity{
				ID:               "id-1",
				Component:        "component-1",
				Enabled:          true,
				ServicesToAdd:    []string{"serv-4"},
				ServicesToRemove: []string{"serv-0", "serv-2", "serv-3"},
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
			expectedEntity: types.Entity{
				ID:               "id-1",
				Component:        "component-1",
				Enabled:          true,
				ServicesToAdd:    []string{"serv-1", "serv-4"},
				ServicesToRemove: []string{"serv-0", "serv-3"},
				Services:         []string{"serv-1", "serv-2"},
			},
		},
		{
			name: "one entity is removed from a single service",
			entity: types.Entity{
				ID:        "id-1",
				Component: "component-1",
				Enabled:   true,
				Services:  []string{"serv-1"},
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
			expectedEntity: types.Entity{
				ID:               "id-1",
				Component:        "component-1",
				Enabled:          true,
				Services:         []string{},
				ServicesToAdd:    []string{},
				ServicesToRemove: []string{"serv-1"},
			},
		},
		{
			name: "one entity is removed from a single service but have this service in ServicesToAdd",
			entity: types.Entity{
				ID:            "id-1",
				Component:     "component-1",
				Enabled:       true,
				Services:      []string{"serv-1"},
				ServicesToAdd: []string{"serv-1"},
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
			expectedEntity: types.Entity{
				ID:               "id-1",
				Component:        "component-1",
				Enabled:          true,
				Services:         []string{},
				ServicesToRemove: []string{},
			},
		},
		{
			name: "one entity is removed from multiple services",
			entity: types.Entity{
				ID:        "id-1",
				Enabled:   true,
				Component: "component-1",
				Services:  []string{"serv-1", "serv-2"},
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
			expectedEntity: types.Entity{
				ID:               "id-1",
				Component:        "component-1",
				Enabled:          true,
				Services:         []string{},
				ServicesToAdd:    []string{},
				ServicesToRemove: []string{"serv-1", "serv-2"},
			},
		},
		{
			name: "one entity is moved from one service to another",
			entity: types.Entity{
				ID:        "id-1",
				Component: "component-1",
				Enabled:   true,
				Services:  []string{"serv-1", "serv-2"},
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
			expectedEntity: types.Entity{
				ID:               "id-1",
				Component:        "component-1",
				Enabled:          true,
				Services:         []string{"serv-1", "serv-3"},
				ServicesToAdd:    []string{"serv-3"},
				ServicesToRemove: []string{"serv-2"},
			},
		},
		{
			name: "no changes",
			entity: types.Entity{
				ID:        "id-1",
				Enabled:   true,
				Component: "component-1",
				Services:  []string{"serv-1", "serv-2"},
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
			expectedEntity: types.Entity{
				ID:        "id-1",
				Enabled:   true,
				Component: "component-1",
				Services:  []string{"serv-1", "serv-2"},
			},
		},
	}

	manager := contextgraph.NewManager(adapter, dbClient, storage, assigner, log.NewLogger(true))

	commRegister := mock_mongo.NewMockCommandsRegister(ctrl)
	commRegister.EXPECT().RegisterUpdate(gomock.Any(), gomock.Any()).AnyTimes()
	commRegister.EXPECT().Clear().AnyTimes()

	for _, dataset := range dataSets {
		t.Run(dataset.name, func(t *testing.T) {
			commRegister.Clear()
			storage.EXPECT().GetAll(gomock.Any()).Return(dataset.services, nil)

			err := manager.LoadServices(ctx)
			if err != nil {
				t.Error(err)
			}

			manager.AssignServices(&dataset.entity, commRegister)

			sort.Slice(dataset.entity.Services, func(i, j int) bool {
				return dataset.entity.Services[i] < dataset.entity.Services[j]
			})

			sort.Slice(dataset.entity.ServicesToAdd, func(i, j int) bool {
				return dataset.entity.ServicesToAdd[i] < dataset.entity.ServicesToAdd[j]
			})

			sort.Slice(dataset.entity.ServicesToRemove, func(i, j int) bool {
				return dataset.entity.ServicesToRemove[i] < dataset.entity.ServicesToRemove[j]
			})

			sort.Slice(dataset.expectedEntity.Services, func(i, j int) bool {
				return dataset.expectedEntity.Services[i] < dataset.expectedEntity.Services[j]
			})

			sort.Slice(dataset.expectedEntity.ServicesToAdd, func(i, j int) bool {
				return dataset.expectedEntity.ServicesToAdd[i] < dataset.expectedEntity.ServicesToAdd[j]
			})

			sort.Slice(dataset.expectedEntity.ServicesToRemove, func(i, j int) bool {
				return dataset.expectedEntity.ServicesToRemove[i] < dataset.expectedEntity.ServicesToRemove[j]
			})

			if slices.Compare(dataset.entity.Services, dataset.expectedEntity.Services) != 0 {
				t.Errorf("expected Services to be %v, but got %v", dataset.expectedEntity.Services, dataset.entity.Services)
			}

			if slices.Compare(dataset.entity.ServicesToAdd, dataset.expectedEntity.ServicesToAdd) != 0 {
				t.Errorf("expected ServicesToAdd to be %v, but got %v", dataset.expectedEntity.ServicesToAdd, dataset.entity.ServicesToAdd)
			}

			if slices.Compare(dataset.entity.ServicesToRemove, dataset.expectedEntity.ServicesToRemove) != 0 {
				t.Errorf("expected ServicesToRemove to be %v, but got %v", dataset.expectedEntity.ServicesToRemove, dataset.entity.ServicesToRemove)
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
	storage.EXPECT().GetAll(gomock.Any()).Return([]entityservice.EntityService{}, nil).AnyTimes()
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

	assigner := mock_statesetting.NewMockAssigner(ctrl)
	assigner.EXPECT().AssignStateSetting(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil).AnyTimes()

	commRegister := mock_mongo.NewMockCommandsRegister(ctrl)
	commRegister.EXPECT().RegisterUpdate(gomock.Any(), gomock.Any()).AnyTimes()

	manager := contextgraph.NewManager(adapter, dbClient, storage, assigner, log.NewLogger(true))
	for i := 0; i < b.N; i++ {
		_, _ = manager.RecomputeService(ctx, "serv-1", commRegister)
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
	storage.EXPECT().GetAll(gomock.Any()).Return([]entityservice.EntityService{}, nil).AnyTimes()
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

	assigner := mock_statesetting.NewMockAssigner(ctrl)
	assigner.EXPECT().AssignStateSetting(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil).AnyTimes()

	commRegister := mock_mongo.NewMockCommandsRegister(ctrl)
	commRegister.EXPECT().RegisterUpdate(gomock.Any(), gomock.Any()).AnyTimes()

	manager := contextgraph.NewManager(adapter, dbClient, storage, assigner, log.NewLogger(true))
	for i := 0; i < b.N; i++ {
		call = 0
		_, _ = manager.RecomputeService(ctx, "serv-1", commRegister)
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
	storage.EXPECT().GetAll(gomock.Any()).Return([]entityservice.EntityService{}, nil).AnyTimes()
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

	assigner := mock_statesetting.NewMockAssigner(ctrl)
	assigner.EXPECT().AssignStateSetting(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil).AnyTimes()

	commRegister := mock_mongo.NewMockCommandsRegister(ctrl)
	commRegister.EXPECT().RegisterUpdate(gomock.Any(), gomock.Any()).AnyTimes()

	manager := contextgraph.NewManager(adapter, dbClient, storage, assigner, log.NewLogger(true))
	for i := 0; i < b.N; i++ {
		call = 0
		_, _ = manager.RecomputeService(ctx, "serv-1", commRegister)
	}
}
