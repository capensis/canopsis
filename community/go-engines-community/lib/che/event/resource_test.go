package event

import (
	"context"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/contextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	mock_eventfilter "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/eventfilter"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/mock/gomock"
)

func TestResourceProcessorAssignsStateSettingToUninitializedComponent(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbClient := mock_mongo.NewMockDbClient(ctrl)
	entityCollection := mock_mongo.NewMockDbCollection(ctrl)
	cursor := mock_mongo.NewMockCursor(ctrl)
	eventFilterService := mock_eventfilter.NewMockService(ctrl)

	dbClient.EXPECT().Collection(gomock.Any()).Return(entityCollection).AnyTimes()
	dbClient.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(context.Context) error) error { return f(ctx) },
	)
	entityCollection.EXPECT().Find(gomock.Any(), gomock.Any()).Return(cursor, nil)
	var resourceFlagsPersisted bool
	entityCollection.EXPECT().BulkWrite(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, models []mongodriver.WriteModel, _ ...options.Lister[options.BulkWriteOptions]) (*mongodriver.BulkWriteResult, error) {
			for _, model := range models {
				updateModel, ok := model.(*mongodriver.UpdateOneModel)
				if !ok {
					continue
				}

				filter, ok := updateModel.Filter.(bson.M)
				if !ok || filter["_id"] != "resource/component" {
					continue
				}

				update, ok := updateModel.Update.(bson.M)
				if !ok {
					continue
				}
				set, ok := update["$set"].(bson.M)
				if ok && set["component_state_settings"] == true && set["component_state_settings_to_add"] == true {
					resourceFlagsPersisted = true
				}
			}

			return &mongodriver.BulkWriteResult{}, nil
		},
	).AnyTimes()
	eventFilterService.EXPECT().ProcessEvent(gomock.Any(), gomock.Any(), gomock.Any()).Return(eventfilter.ServiceResult{}, nil)

	entities := []types.Entity{
		{
			ID:        "resource/component",
			Type:      types.EntityTypeResource,
			Enabled:   true,
			Component: "component",
			Infos:     map[string]types.Info{"env": {Value: "prod"}},
		},
		{
			ID:        "component",
			Type:      types.EntityTypeComponent,
			Enabled:   true,
			StateInfo: nil, // Another concurrent resource event created this component before assigning its rule.
		},
	}
	entityIndex := 0
	cursorIndex := 0
	cursor.EXPECT().Next(gomock.Any()).DoAndReturn(func(context.Context) bool {
		cursorIndex++
		return cursorIndex <= len(entities)
	}).Times(len(entities) + 1)
	cursor.EXPECT().Decode(gomock.Any()).DoAndReturn(func(v any) error {
		*(v.(*types.Entity)) = entities[entityIndex]
		entityIndex++
		return nil
	}).Times(len(entities))
	cursor.EXPECT().Close(gomock.Any()).Return(nil)

	stateSettingAssigner := &testStateSettingAssigner{}
	contextGraphManager := contextgraph.NewManager(nil, dbClient, testEntityServiceStorage{}, stateSettingAssigner, zerolog.Nop())
	processor := &resourceProcessor{
		dbClient:            dbClient,
		dbEntityCollection:  entityCollection,
		dbAlarmCollection:   entityCollection,
		contextGraphManager: contextGraphManager,
		eventFilterService:  eventFilterService,
		logger:              zerolog.Nop(),
	}
	eventEntity := entities[0]
	event := &types.Event{Entity: &eventEntity}
	partialResult := &ProcessorResult{
		// Processing resumes after the context-creation transaction; another resource event created this component first.
		ContextGraphReport: contextgraph.Report{CheckResource: true},
	}

	_, err := processor.Process(context.Background(), event, partialResult)
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	if event.Entity == nil || !event.Entity.ComponentStateSettings || !event.Entity.ComponentStateSettingsToAdd {
		t.Fatalf("resource state-setting flags = %+v, want enabled and queued for add", event.Entity)
	}
	if !resourceFlagsPersisted {
		t.Fatal("inherited state-setting flags were not persisted for the resource")
	}
}

type testStateSettingAssigner struct{}

func (*testStateSettingAssigner) AssignStateSetting(_ context.Context, entity *types.Entity, commRegister mongo.CommandsRegister) (bool, error) {
	entity.StateInfo = &types.StateInfo{
		ID: "state-setting",
		InheritedPattern: pattern.Entity{{
			{
				Field:     "infos.env",
				FieldType: pattern.FieldTypeString,
				Condition: pattern.NewStringCondition(pattern.ConditionEqual, "prod"),
			},
		}},
	}
	commRegister.RegisterUpdate(entity.ID, bson.M{"state_info": entity.StateInfo})
	return true, nil
}

func (testStateSettingAssigner) LoadRules(context.Context) error { return nil }

type testEntityServiceStorage struct{}

func (testEntityServiceStorage) GetAll(context.Context) ([]contextgraph.EntityService, error) {
	return nil, nil
}

func (testEntityServiceStorage) Get(context.Context, string) (entityservice.EntityService, error) {
	return entityservice.EntityService{}, nil
}
