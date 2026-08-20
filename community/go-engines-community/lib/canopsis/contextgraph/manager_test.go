package contextgraph_test

import (
	"context"
	"fmt"
	"maps"
	"os"
	"slices"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/contextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	mock_contextgraph "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/contextgraph"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	mock_statesetting "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/statesetting"
	"github.com/rs/zerolog"
	"go.uber.org/mock/gomock"
)

func TestCheckServices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	collection := mock_mongo.NewMockDbCollection(ctrl)

	dbClient := mock_mongo.NewMockDbClient(ctrl)
	dbClient.EXPECT().Collection(mongo.EntityMongoCollection).Return(collection).AnyTimes()

	adapter := entity.NewAdapter(dbClient)
	storage := mock_contextgraph.NewMockEntityServiceStorage(ctrl)

	assigner := mock_statesetting.NewMockAssigner(ctrl)

	dataSets := []struct {
		services       []contextgraph.EntityService
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
			services: []contextgraph.EntityService{
				{
					ID:      "serv-1",
					Enabled: true,
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
			services: []contextgraph.EntityService{
				{
					ID:      "serv-1",
					Enabled: true,
					EntityPattern: [][]pattern.FieldCondition{
						{
							{
								Field:     "component",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
							},
						},
					},
				},
				{
					ID:      "serv-2",
					Enabled: true,
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
			services: []contextgraph.EntityService{
				{
					ID:      "serv-1",
					Enabled: true,
					EntityPattern: [][]pattern.FieldCondition{
						{
							{
								Field:     "component",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
							},
						},
					},
				},
				{
					ID:      "serv-2",
					Enabled: true,
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
			services: []contextgraph.EntityService{
				{
					ID:      "serv-1",
					Enabled: true,
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
			services: []contextgraph.EntityService{
				{
					ID:      "serv-1",
					Enabled: true,
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
			services: []contextgraph.EntityService{
				{
					ID:      "serv-1",
					Enabled: true,
					EntityPattern: [][]pattern.FieldCondition{
						{
							{
								Field:     "component",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-2"),
							},
						},
					},
				},
				{
					ID:      "serv-2",
					Enabled: true,
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
			services: []contextgraph.EntityService{
				{
					ID:      "serv-1",
					Enabled: true,
					EntityPattern: [][]pattern.FieldCondition{
						{
							{
								Field:     "component",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
							},
						},
					},
				},
				{
					ID:      "serv-2",
					Enabled: true,
					EntityPattern: [][]pattern.FieldCondition{
						{
							{
								Field:     "component",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-2"),
							},
						},
					},
				},
				{
					ID:      "serv-3",
					Enabled: true,
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
			services: []contextgraph.EntityService{
				{
					ID:      "serv-1",
					Enabled: true,
					EntityPattern: [][]pattern.FieldCondition{
						{
							{
								Field:     "component",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "component-1"),
							},
						},
					},
				},
				{
					ID:      "serv-2",
					Enabled: true,
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
			expectedEntity: types.Entity{
				ID:        "id-1",
				Enabled:   true,
				Component: "component-1",
				Services:  []string{"serv-1", "serv-2"},
			},
		},
	}

	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)
	manager := contextgraph.NewManager(adapter, dbClient, storage, assigner, logger)

	commRegister := mock_mongo.NewMockCommandsRegister(ctrl)
	commRegister.EXPECT().RegisterUpdate(gomock.Any(), gomock.Any()).AnyTimes()
	commRegister.EXPECT().Clear().AnyTimes()

	for _, dataset := range dataSets {
		t.Run(dataset.name, func(t *testing.T) {
			commRegister.Clear()
			storage.EXPECT().GetAll(gomock.Any()).Return(dataset.services, nil)

			err := manager.LoadServices(t.Context())
			if err != nil {
				t.Error(err)
			}

			manager.AssignServices(&dataset.entity, commRegister)

			slices.Sort(dataset.entity.Services)
			slices.Sort(dataset.entity.ServicesToAdd)
			slices.Sort(dataset.entity.ServicesToRemove)
			slices.Sort(dataset.expectedEntity.Services)
			slices.Sort(dataset.expectedEntity.ServicesToAdd)
			slices.Sort(dataset.expectedEntity.ServicesToRemove)

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

func TestCheckServicesByInfos(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	collection := mock_mongo.NewMockDbCollection(ctrl)

	dbClient := mock_mongo.NewMockDbClient(ctrl)
	dbClient.EXPECT().Collection(mongo.EntityMongoCollection).Return(collection).AnyTimes()

	adapter := entity.NewAdapter(dbClient)
	storage := mock_contextgraph.NewMockEntityServiceStorage(ctrl)
	assigner := mock_statesetting.NewMockAssigner(ctrl)

	dataSets := []struct {
		name                     string
		infoUpdates              []string
		componentInfoUpdates     []string
		services                 []contextgraph.EntityService
		entity                   types.Entity
		expectedEntity           types.Entity
		expectedAffectedServices []string
	}{
		{
			name:        "targeted infos update adds matching service and preserves unrelated service",
			infoUpdates: []string{"team"},
			entity: types.Entity{
				ID:      "id-1",
				Enabled: true,
				Infos: map[string]types.Info{
					"team":  {Value: "core"},
					"owner": {Value: "ops"},
				},
				Services: []string{"serv-owner"},
			},
			services: []contextgraph.EntityService{
				{
					ID:      "serv-team",
					Enabled: true,
					EntityPattern: [][]pattern.FieldCondition{{
						{
							Field:     "infos.team",
							FieldType: pattern.FieldTypeString,
							Condition: pattern.NewStringCondition(pattern.ConditionEqual, "core"),
						},
					}},
				},
				{
					ID:      "serv-owner",
					Enabled: true,
					EntityPattern: [][]pattern.FieldCondition{{
						{
							Field:     "infos.owner",
							FieldType: pattern.FieldTypeString,
							Condition: pattern.NewStringCondition(pattern.ConditionEqual, "ops"),
						},
					}},
				},
			},
			expectedEntity: types.Entity{
				ID:      "id-1",
				Enabled: true,
				Infos: map[string]types.Info{
					"team":  {Value: "core"},
					"owner": {Value: "ops"},
				},
				Services:      []string{"serv-owner", "serv-team"},
				ServicesToAdd: []string{"serv-team"},
			},
		},
		{
			name:        "targeted infos update removes only matching indexed service",
			infoUpdates: []string{"team"},
			entity: types.Entity{
				ID:      "id-1",
				Enabled: true,
				Infos: map[string]types.Info{
					"team":  {Value: "dev"},
					"owner": {Value: "ops"},
				},
				Services: []string{"serv-owner", "serv-team"},
			},
			services: []contextgraph.EntityService{
				{
					ID:      "serv-team",
					Enabled: true,
					EntityPattern: [][]pattern.FieldCondition{{
						{
							Field:     "infos.team",
							FieldType: pattern.FieldTypeString,
							Condition: pattern.NewStringCondition(pattern.ConditionEqual, "core"),
						},
					}},
				},
				{
					ID:      "serv-owner",
					Enabled: true,
					EntityPattern: [][]pattern.FieldCondition{{
						{
							Field:     "infos.owner",
							FieldType: pattern.FieldTypeString,
							Condition: pattern.NewStringCondition(pattern.ConditionEqual, "ops"),
						},
					}},
				},
			},
			expectedEntity: types.Entity{
				ID:      "id-1",
				Enabled: true,
				Infos: map[string]types.Info{
					"team":  {Value: "dev"},
					"owner": {Value: "ops"},
				},
				Services:         []string{"serv-owner"},
				ServicesToAdd:    []string{},
				ServicesToRemove: []string{"serv-team"},
			},
		},
		{
			name:        "targeted infos update with no indexed services does nothing",
			infoUpdates: []string{"zone"},
			entity: types.Entity{
				ID:      "id-1",
				Enabled: true,
				Infos: map[string]types.Info{
					"team": {Value: "core"},
				},
				Services: []string{"serv-team"},
			},
			services: []contextgraph.EntityService{
				{
					ID:      "serv-team",
					Enabled: true,
					EntityPattern: [][]pattern.FieldCondition{{
						{
							Field:     "infos.team",
							FieldType: pattern.FieldTypeString,
							Condition: pattern.NewStringCondition(pattern.ConditionEqual, "core"),
						},
					}},
				},
			},
			expectedEntity: types.Entity{
				ID:      "id-1",
				Enabled: true,
				Infos: map[string]types.Info{
					"team": {Value: "core"},
				},
				Services: []string{"serv-team"},
			},
		},
		{
			name:                 "targeted component_infos update adds matching service and preserves unrelated service",
			componentInfoUpdates: []string{"env"},
			entity: types.Entity{
				ID:      "id-1",
				Enabled: true,
				ComponentInfos: map[string]types.Info{
					"env":  {Value: "prod"},
					"zone": {Value: "eu"},
				},
				Services: []string{"serv-zone"},
			},
			services: []contextgraph.EntityService{
				{
					ID:      "serv-env",
					Enabled: true,
					EntityPattern: [][]pattern.FieldCondition{{
						{
							Field:     "component_infos.env",
							FieldType: pattern.FieldTypeString,
							Condition: pattern.NewStringCondition(pattern.ConditionEqual, "prod"),
						},
					}},
				},
				{
					ID:      "serv-zone",
					Enabled: true,
					EntityPattern: [][]pattern.FieldCondition{{
						{
							Field:     "component_infos.zone",
							FieldType: pattern.FieldTypeString,
							Condition: pattern.NewStringCondition(pattern.ConditionEqual, "eu"),
						},
					}},
				},
			},
			expectedEntity: types.Entity{
				ID:      "id-1",
				Enabled: true,
				ComponentInfos: map[string]types.Info{
					"env":  {Value: "prod"},
					"zone": {Value: "eu"},
				},
				Services:      []string{"serv-env", "serv-zone"},
				ServicesToAdd: []string{"serv-env"},
			},
			expectedAffectedServices: []string{"serv-env"},
		},
		{
			name:                 "targeted component_infos update removes only matching indexed service",
			componentInfoUpdates: []string{"env"},
			entity: types.Entity{
				ID:      "id-1",
				Enabled: true,
				ComponentInfos: map[string]types.Info{
					"env":  {Value: "dev"},
					"zone": {Value: "eu"},
				},
				Services: []string{"serv-env", "serv-zone"},
			},
			services: []contextgraph.EntityService{
				{
					ID:      "serv-env",
					Enabled: true,
					EntityPattern: [][]pattern.FieldCondition{{
						{
							Field:     "component_infos.env",
							FieldType: pattern.FieldTypeString,
							Condition: pattern.NewStringCondition(pattern.ConditionEqual, "prod"),
						},
					}},
				},
				{
					ID:      "serv-zone",
					Enabled: true,
					EntityPattern: [][]pattern.FieldCondition{{
						{
							Field:     "component_infos.zone",
							FieldType: pattern.FieldTypeString,
							Condition: pattern.NewStringCondition(pattern.ConditionEqual, "eu"),
						},
					}},
				},
			},
			expectedEntity: types.Entity{
				ID:      "id-1",
				Enabled: true,
				ComponentInfos: map[string]types.Info{
					"env":  {Value: "dev"},
					"zone": {Value: "eu"},
				},
				Services:         []string{"serv-zone"},
				ServicesToAdd:    []string{},
				ServicesToRemove: []string{"serv-env"},
			},
			expectedAffectedServices: []string{"serv-env"},
		},
		{
			name:                 "targeted component_infos update with no indexed services does nothing",
			componentInfoUpdates: []string{"region"},
			entity: types.Entity{
				ID:      "id-1",
				Enabled: true,
				ComponentInfos: map[string]types.Info{
					"env": {Value: "prod"},
				},
				Services: []string{"serv-env"},
			},
			services: []contextgraph.EntityService{
				{
					ID:      "serv-env",
					Enabled: true,
					EntityPattern: [][]pattern.FieldCondition{{
						{
							Field:     "component_infos.env",
							FieldType: pattern.FieldTypeString,
							Condition: pattern.NewStringCondition(pattern.ConditionEqual, "prod"),
						},
					}},
				},
			},
			expectedEntity: types.Entity{
				ID:      "id-1",
				Enabled: true,
				ComponentInfos: map[string]types.Info{
					"env": {Value: "prod"},
				},
				Services: []string{"serv-env"},
			},
			expectedAffectedServices: []string{},
		},
		{
			name:        "service indexed by both infos and component_infos matches only when both conditions match",
			infoUpdates: []string{"team"},
			entity: types.Entity{
				ID:      "id-1",
				Enabled: true,
				Infos: map[string]types.Info{
					"team": {Value: "core"},
				},
				ComponentInfos: map[string]types.Info{
					"env": {Value: "dev"},
				},
			},
			services: []contextgraph.EntityService{
				{
					ID:      "serv-team-prod",
					Enabled: true,
					EntityPattern: [][]pattern.FieldCondition{{
						{
							Field:     "infos.team",
							FieldType: pattern.FieldTypeString,
							Condition: pattern.NewStringCondition(pattern.ConditionEqual, "core"),
						},
						{
							Field:     "component_infos.env",
							FieldType: pattern.FieldTypeString,
							Condition: pattern.NewStringCondition(pattern.ConditionEqual, "prod"),
						},
					}},
				},
			},
			expectedEntity: types.Entity{
				ID:      "id-1",
				Enabled: true,
				Infos: map[string]types.Info{
					"team": {Value: "core"},
				},
				ComponentInfos: map[string]types.Info{
					"env": {Value: "dev"},
				},
			},
		},
		{
			name: "empty updates do not modify services",
			entity: types.Entity{
				ID:      "id-1",
				Enabled: true,
				Infos: map[string]types.Info{
					"team": {Value: "core"},
				},
				Services: []string{"serv-team"},
			},
			services: []contextgraph.EntityService{
				{
					ID:      "serv-team",
					Enabled: true,
					EntityPattern: [][]pattern.FieldCondition{{
						{
							Field:     "infos.team",
							FieldType: pattern.FieldTypeString,
							Condition: pattern.NewStringCondition(pattern.ConditionEqual, "other"),
						},
					}},
				},
			},
			expectedEntity: types.Entity{
				ID:      "id-1",
				Enabled: true,
				Infos: map[string]types.Info{
					"team": {Value: "core"},
				},
				Services: []string{"serv-team"},
			},
		},
		{
			name:                 "component_infos update transitions entity to inherited",
			componentInfoUpdates: []string{"env"},
			entity: types.Entity{
				ID:      "id-1",
				Name:    "my-entity",
				Enabled: true,
				ComponentInfos: map[string]types.Info{
					"env": {Value: "prod"},
				},
			},
			services: []contextgraph.EntityService{
				{
					ID:      "serv-inherited",
					Enabled: true,
					EntityPattern: [][]pattern.FieldCondition{{
						{
							Field:     "name",
							FieldType: pattern.FieldTypeString,
							Condition: pattern.NewStringCondition(pattern.ConditionEqual, "other-entity"),
						},
					}},
					InheritedPattern: [][]pattern.FieldCondition{{
						{
							Field:     "component_infos.env",
							FieldType: pattern.FieldTypeString,
							Condition: pattern.NewStringCondition(pattern.ConditionEqual, "prod"),
						},
					}},
				},
			},
			expectedEntity: types.Entity{
				ID:      "id-1",
				Name:    "my-entity",
				Enabled: true,
				ComponentInfos: map[string]types.Info{
					"env": {Value: "prod"},
				},
			},
			expectedAffectedServices: []string{"serv-inherited"},
		},
		{
			name:                 "component_infos update transitions entity from inherited",
			componentInfoUpdates: []string{"env"},
			entity: types.Entity{
				ID:      "id-1",
				Name:    "my-entity",
				Enabled: true,
				ComponentInfos: map[string]types.Info{
					"env": {Value: "dev"},
				},
				InheritedServices: []string{"serv-inherited"},
			},
			services: []contextgraph.EntityService{
				{
					ID:      "serv-inherited",
					Enabled: true,
					EntityPattern: [][]pattern.FieldCondition{{
						{
							Field:     "name",
							FieldType: pattern.FieldTypeString,
							Condition: pattern.NewStringCondition(pattern.ConditionEqual, "other-entity"),
						},
					}},
					InheritedPattern: [][]pattern.FieldCondition{{
						{
							Field:     "component_infos.env",
							FieldType: pattern.FieldTypeString,
							Condition: pattern.NewStringCondition(pattern.ConditionEqual, "prod"),
						},
					}},
				},
			},
			expectedEntity: types.Entity{
				ID:      "id-1",
				Name:    "my-entity",
				Enabled: true,
				ComponentInfos: map[string]types.Info{
					"env": {Value: "dev"},
				},
				InheritedServices: []string{"serv-inherited"},
			},
			expectedAffectedServices: []string{"serv-inherited"},
		},
		{
			name:                 "component_infos update with stable inherited status does not affect service",
			componentInfoUpdates: []string{"env"},
			entity: types.Entity{
				ID:      "id-1",
				Name:    "my-entity",
				Enabled: true,
				ComponentInfos: map[string]types.Info{
					"env": {Value: "prod"},
				},
				InheritedServices: []string{"serv-inherited"},
			},
			services: []contextgraph.EntityService{
				{
					ID:      "serv-inherited",
					Enabled: true,
					EntityPattern: [][]pattern.FieldCondition{{
						{
							Field:     "name",
							FieldType: pattern.FieldTypeString,
							Condition: pattern.NewStringCondition(pattern.ConditionEqual, "other-entity"),
						},
					}},
					InheritedPattern: [][]pattern.FieldCondition{{
						{
							Field:     "component_infos.env",
							FieldType: pattern.FieldTypeString,
							Condition: pattern.NewStringCondition(pattern.ConditionEqual, "prod"),
						},
					}},
				},
			},
			expectedEntity: types.Entity{
				ID:      "id-1",
				Name:    "my-entity",
				Enabled: true,
				ComponentInfos: map[string]types.Info{
					"env": {Value: "prod"},
				},
				InheritedServices: []string{"serv-inherited"},
			},
			expectedAffectedServices: []string{},
		},
		{
			name:                 "component_infos update with positive match takes priority over inherited path",
			componentInfoUpdates: []string{"env"},
			entity: types.Entity{
				ID:      "id-1",
				Name:    "my-entity",
				Enabled: true,
				ComponentInfos: map[string]types.Info{
					"env":  {Value: "prod"},
					"zone": {Value: "eu"},
				},
			},
			services: []contextgraph.EntityService{
				{
					ID:      "serv-mixed",
					Enabled: true,
					EntityPattern: [][]pattern.FieldCondition{{
						{
							Field:     "component_infos.env",
							FieldType: pattern.FieldTypeString,
							Condition: pattern.NewStringCondition(pattern.ConditionEqual, "prod"),
						},
					}},
					InheritedPattern: [][]pattern.FieldCondition{{
						{
							Field:     "component_infos.zone",
							FieldType: pattern.FieldTypeString,
							Condition: pattern.NewStringCondition(pattern.ConditionEqual, "eu"),
						},
					}},
				},
			},
			expectedEntity: types.Entity{
				ID:      "id-1",
				Name:    "my-entity",
				Enabled: true,
				ComponentInfos: map[string]types.Info{
					"env":  {Value: "prod"},
					"zone": {Value: "eu"},
				},
				Services:      []string{"serv-mixed"},
				ServicesToAdd: []string{"serv-mixed"},
			},
			expectedAffectedServices: []string{"serv-mixed"},
		},
		{
			name:                 "inherited pattern indexed by component_infos triggers recompute when matched via other key",
			componentInfoUpdates: []string{"zone"},
			entity: types.Entity{
				ID:      "id-1",
				Name:    "my-entity",
				Enabled: true,
				ComponentInfos: map[string]types.Info{
					"env":  {Value: "prod"},
					"zone": {Value: "eu"},
				},
			},
			services: []contextgraph.EntityService{
				{
					ID:      "serv-inherited",
					Enabled: true,
					EntityPattern: [][]pattern.FieldCondition{{
						{
							Field:     "name",
							FieldType: pattern.FieldTypeString,
							Condition: pattern.NewStringCondition(pattern.ConditionEqual, "other-entity"),
						},
					}},
					InheritedPattern: [][]pattern.FieldCondition{{
						{
							Field:     "component_infos.env",
							FieldType: pattern.FieldTypeString,
							Condition: pattern.NewStringCondition(pattern.ConditionEqual, "prod"),
						},
						{
							Field:     "component_infos.zone",
							FieldType: pattern.FieldTypeString,
							Condition: pattern.NewStringCondition(pattern.ConditionEqual, "eu"),
						},
					}},
				},
			},
			expectedEntity: types.Entity{
				ID:      "id-1",
				Name:    "my-entity",
				Enabled: true,
				ComponentInfos: map[string]types.Info{
					"env":  {Value: "prod"},
					"zone": {Value: "eu"},
				},
			},
			expectedAffectedServices: []string{"serv-inherited"},
		},
	}

	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)
	manager := contextgraph.NewManager(adapter, dbClient, storage, assigner, logger)

	commRegister := mock_mongo.NewMockCommandsRegister(ctrl)
	commRegister.EXPECT().RegisterUpdate(gomock.Any(), gomock.Any()).AnyTimes()
	commRegister.EXPECT().Clear().AnyTimes()

	for _, dataset := range dataSets {
		t.Run(dataset.name, func(t *testing.T) {
			commRegister.Clear()
			storage.EXPECT().GetAll(gomock.Any()).Return(dataset.services, nil)

			err := manager.LoadServices(t.Context())
			if err != nil {
				t.Error(err)
			}

			manager.AssignServicesByInfoNames(&dataset.entity, dataset.infoUpdates, commRegister)
			affected := manager.AssignServicesByComponentInfoNames(&dataset.entity, dataset.componentInfoUpdates, commRegister)

			slices.Sort(dataset.entity.Services)
			slices.Sort(dataset.entity.ServicesToAdd)
			slices.Sort(dataset.entity.ServicesToRemove)
			slices.Sort(dataset.entity.InheritedServices)
			slices.Sort(dataset.expectedEntity.Services)
			slices.Sort(dataset.expectedEntity.ServicesToAdd)
			slices.Sort(dataset.expectedEntity.ServicesToRemove)
			slices.Sort(dataset.expectedEntity.InheritedServices)

			if slices.Compare(dataset.entity.Services, dataset.expectedEntity.Services) != 0 {
				t.Errorf("expected Services to be %v, but got %v", dataset.expectedEntity.Services, dataset.entity.Services)
			}

			if slices.Compare(dataset.entity.ServicesToAdd, dataset.expectedEntity.ServicesToAdd) != 0 {
				t.Errorf("expected ServicesToAdd to be %v, but got %v", dataset.expectedEntity.ServicesToAdd, dataset.entity.ServicesToAdd)
			}

			if slices.Compare(dataset.entity.ServicesToRemove, dataset.expectedEntity.ServicesToRemove) != 0 {
				t.Errorf("expected ServicesToRemove to be %v, but got %v", dataset.expectedEntity.ServicesToRemove, dataset.entity.ServicesToRemove)
			}

			if slices.Compare(dataset.entity.InheritedServices, dataset.expectedEntity.InheritedServices) != 0 {
				t.Errorf("expected InheritedServices to be %v, but got %v", dataset.expectedEntity.InheritedServices, dataset.entity.InheritedServices)
			}

			if dataset.expectedAffectedServices != nil {
				gotAffected := slices.Sorted(maps.Keys(affected))
				slices.Sort(dataset.expectedAffectedServices)
				if slices.Compare(gotAffected, dataset.expectedAffectedServices) != 0 {
					t.Errorf("expected affected services to be %v, but got %v", dataset.expectedAffectedServices, gotAffected)
				}
			}
		})
	}
}

func BenchmarkRecomputeServicesRemoveAll(b *testing.B) {
	ctx := b.Context()
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
	cursor.EXPECT().All(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, results any) {
		ents := results.(*[]types.Entity)
		*ents = append(*ents, entities...)
	}).Return(nil).AnyTimes()
	cursor.EXPECT().Next(gomock.Any()).Return(false).AnyTimes()
	cursor.EXPECT().Close(gomock.Any()).Return(nil).AnyTimes()
	cursor.EXPECT().Err().Return(nil).AnyTimes()

	collection := mock_mongo.NewMockDbCollection(ctrl)
	collection.EXPECT().Aggregate(gomock.Any(), gomock.Any()).Return(cursor, nil).AnyTimes()

	dbClient := mock_mongo.NewMockDbClient(ctrl)
	dbClient.EXPECT().Collection(mongo.EntityMongoCollection).Return(collection).AnyTimes()

	adapter := entity.NewAdapter(dbClient)
	storage := mock_contextgraph.NewMockEntityServiceStorage(ctrl)
	storage.EXPECT().GetAll(gomock.Any()).Return([]contextgraph.EntityService{}, nil).AnyTimes()
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

	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)
	manager := contextgraph.NewManager(adapter, dbClient, storage, assigner, logger)
	for b.Loop() {
		_, _ = manager.RecomputeService(ctx, "serv-1", commRegister)
	}
}

func BenchmarkRecomputeServicesAddAll(b *testing.B) {
	ctx := b.Context()
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

	addIdx := 0
	cursor := mock_mongo.NewMockCursor(ctrl)
	cursor.EXPECT().All(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	cursor.EXPECT().Next(gomock.Any()).DoAndReturn(func(_ context.Context) bool {
		return addIdx < len(entities)
	}).AnyTimes()
	cursor.EXPECT().Decode(gomock.Any()).DoAndReturn(func(val any) error {
		ent := val.(*types.Entity)
		*ent = entities[addIdx]
		addIdx++
		return nil
	}).AnyTimes()
	cursor.EXPECT().Close(gomock.Any()).Return(nil).AnyTimes()
	cursor.EXPECT().Err().Return(nil).AnyTimes()

	collection := mock_mongo.NewMockDbCollection(ctrl)
	collection.EXPECT().Aggregate(gomock.Any(), gomock.Any()).Return(cursor, nil).AnyTimes()

	dbClient := mock_mongo.NewMockDbClient(ctrl)
	dbClient.EXPECT().Collection(mongo.EntityMongoCollection).Return(collection).AnyTimes()

	adapter := entity.NewAdapter(dbClient)
	storage := mock_contextgraph.NewMockEntityServiceStorage(ctrl)
	storage.EXPECT().GetAll(gomock.Any()).Return([]contextgraph.EntityService{}, nil).AnyTimes()
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

	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)
	manager := contextgraph.NewManager(adapter, dbClient, storage, assigner, logger)
	for b.Loop() {
		addIdx = 0
		_, _ = manager.RecomputeService(ctx, "serv-1", commRegister)
	}
}

func BenchmarkRecomputeServicesMixed(b *testing.B) {
	ctx := b.Context()
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

	addIdx := 0
	cursor := mock_mongo.NewMockCursor(ctrl)
	cursor.EXPECT().All(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, results any) {
		ents := results.(*[]types.Entity)
		*ents = append(*ents, entitiesToRemove...)
	}).Return(nil).AnyTimes()
	cursor.EXPECT().Next(gomock.Any()).DoAndReturn(func(_ context.Context) bool {
		return addIdx < len(entitiesToAdd)
	}).AnyTimes()
	cursor.EXPECT().Decode(gomock.Any()).DoAndReturn(func(val any) error {
		ent := val.(*types.Entity)
		*ent = entitiesToAdd[addIdx]
		addIdx++
		return nil
	}).AnyTimes()
	cursor.EXPECT().Close(gomock.Any()).Return(nil).AnyTimes()
	cursor.EXPECT().Err().Return(nil).AnyTimes()

	collection := mock_mongo.NewMockDbCollection(ctrl)
	collection.EXPECT().Aggregate(gomock.Any(), gomock.Any()).Return(cursor, nil).AnyTimes()

	dbClient := mock_mongo.NewMockDbClient(ctrl)
	dbClient.EXPECT().Collection(mongo.EntityMongoCollection).Return(collection).AnyTimes()

	adapter := entity.NewAdapter(dbClient)
	storage := mock_contextgraph.NewMockEntityServiceStorage(ctrl)
	storage.EXPECT().GetAll(gomock.Any()).Return([]contextgraph.EntityService{}, nil).AnyTimes()
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

	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)
	manager := contextgraph.NewManager(adapter, dbClient, storage, assigner, logger)
	for b.Loop() {
		addIdx = 0
		_, _ = manager.RecomputeService(ctx, "serv-1", commRegister)
	}
}
