package idlealarm

import (
	"bytes"
	"context"
	"reflect"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlerule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	mock_alarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/alarm"
	mock_encoding "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/encoding"
	mock_engine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/engine"
	mock_entity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/entity"
	mock_idlerule "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/idlerule"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	"github.com/golang/mock/gomock"
	"github.com/kylelemons/godebug/pretty"
	"github.com/rs/zerolog"
)

func TestService_Process_GivenAlarmRuleByLastEventDate_ShouldReturnEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mockRuleAdapter := mock_idlerule.NewMockRuleAdapter(ctrl)
	mockAlarmAdapter := mock_alarm.NewMockAdapter(ctrl)
	mockEntityAdapter := mock_entity.NewMockAdapter(ctrl)
	mockEncoder := mock_encoding.NewMockEncoder(ctrl)
	mockRPCClient := mock_engine.NewMockRPCClient(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	entity := types.Entity{
		Type: types.EntityTypeResource,
	}
	alarm := types.Alarm{
		ID: "test-alarm",
		Value: types.AlarmValue{
			LastEventDate: datetime.CpsTime{Time: time.Now().Add(-6 * time.Hour)},
			Connector:     canopsis.AxeConnector,
			ConnectorName: canopsis.AxeConnector,
			Component:     "test-component",
			Resource:      "test-resource",
		},
	}
	rule := idlerule.Rule{
		Name:           "test-rule-name",
		Type:           idlerule.RuleTypeAlarm,
		AlarmCondition: idlerule.RuleAlarmConditionLastEvent,
		Duration: datetime.DurationWithUnit{
			Value: 10,
			Unit:  "s",
		},
		Operation: &idlerule.Operation{
			Type: types.ActionTypeAck,
			Parameters: idlerule.Parameters{
				Output: "test-output",
			},
		},
		AlarmPatternFields: savedpattern.AlarmPatternFields{
			AlarmPattern: [][]pattern.FieldCondition{
				{
					{
						Field:     "v.resource",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource"),
					},
				},
			},
		},
	}

	mockRuleAdapter.
		EXPECT().
		GetEnabled(gomock.Any()).
		Return([]idlerule.Rule{rule}, nil)
	mockCursor := mock_mongo.NewMockCursor(ctrl)
	firstCall := mockCursor.EXPECT().Next(gomock.Any()).Return(true)
	secondCall := mockCursor.EXPECT().Next(gomock.Any()).Return(false)
	gomock.InOrder(firstCall, secondCall)
	mockCursor.EXPECT().Close(gomock.Any())
	mockCursor.EXPECT().Decode(gomock.Any()).Do(func(v *types.AlarmWithEntity) {
		*v = types.AlarmWithEntity{Alarm: alarm, Entity: entity}
	}).Return(nil)
	mockAlarmAdapter.
		EXPECT().
		GetOpenedAlarmsWithLastDatesBefore(gomock.Any(), gomock.Any()).
		Return(mockCursor, nil)
	mockAlarmAdapter.EXPECT().GetOpenedAlarmsByConnectorIdleRules(gomock.Any()).Return(nil, nil)

	service := NewService(mockRuleAdapter, mockAlarmAdapter, mockEntityAdapter, mockRPCClient, canopsis.AxeConnector, mockEncoder, logger)
	events, err := service.Process(ctx)

	if err != nil {
		t.Errorf("exepected no error but got %v", err)
		return
	}

	if len(events) == 0 {
		t.Errorf("exepected event but got nothing")
		return
	}

	expectedEvent := types.Event{
		EventType:     types.EventTypeAck,
		Connector:     alarm.Value.Connector,
		ConnectorName: alarm.Value.ConnectorName,
		Component:     alarm.Value.Component,
		Resource:      alarm.Value.Resource,
		Author:        canopsis.DefaultEventAuthor,
		Output:        "Idle rule: test-rule-name. Comment: test-output.",
		SourceType:    types.SourceTypeResource,
		Initiator:     types.InitiatorSystem,
		IdleRuleApply: "alarm_last_event",
	}
	event := events[0]
	event.Timestamp = datetime.CpsTime{}

	if diff := pretty.Compare(event, expectedEvent); diff != "" {
		t.Errorf("unexepected event %s", diff)
	}
}

func TestService_Process_GivenAlarmRuleByLastUpdateDate_ShouldReturnEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mockRuleAdapter := mock_idlerule.NewMockRuleAdapter(ctrl)
	mockAlarmAdapter := mock_alarm.NewMockAdapter(ctrl)
	mockEntityAdapter := mock_entity.NewMockAdapter(ctrl)
	mockEncoder := mock_encoding.NewMockEncoder(ctrl)
	mockRPCClient := mock_engine.NewMockRPCClient(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	entity := types.Entity{
		Type: types.EntityTypeResource,
	}
	alarm := types.Alarm{
		ID: "test-alarm",
		Value: types.AlarmValue{
			LastUpdateDate: datetime.CpsTime{Time: time.Now().Add(-6 * time.Hour)},
			Connector:      canopsis.AxeConnector,
			ConnectorName:  canopsis.AxeConnector,
			Component:      "test-component",
			Resource:       "test-resource",
		},
	}
	rule := idlerule.Rule{
		Name:           "test-rule-name",
		Type:           idlerule.RuleTypeAlarm,
		AlarmCondition: idlerule.RuleAlarmConditionLastUpdate,
		Duration: datetime.DurationWithUnit{
			Value: 10,
			Unit:  "s",
		},
		Operation: &idlerule.Operation{
			Type: types.ActionTypeAck,
			Parameters: idlerule.Parameters{
				Output: "test-output",
			},
		},
		AlarmPatternFields: savedpattern.AlarmPatternFields{
			AlarmPattern: [][]pattern.FieldCondition{
				{
					{
						Field:     "v.resource",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource"),
					},
				},
			},
		},
	}
	mockCursor := mock_mongo.NewMockCursor(ctrl)
	firstCall := mockCursor.EXPECT().Next(gomock.Any()).Return(true)
	secondCall := mockCursor.EXPECT().Next(gomock.Any()).Return(false)
	gomock.InOrder(firstCall, secondCall)
	mockCursor.EXPECT().Close(gomock.Any())
	mockCursor.EXPECT().Decode(gomock.Any()).Do(func(v *types.AlarmWithEntity) {
		*v = types.AlarmWithEntity{Alarm: alarm, Entity: entity}
	}).Return(nil)
	mockRuleAdapter.
		EXPECT().
		GetEnabled(gomock.Any()).
		Return([]idlerule.Rule{rule}, nil)
	mockAlarmAdapter.
		EXPECT().
		GetOpenedAlarmsWithLastDatesBefore(gomock.Any(), gomock.Any()).
		Return(mockCursor, nil)
	mockAlarmAdapter.EXPECT().GetOpenedAlarmsByConnectorIdleRules(gomock.Any()).Return(nil, nil)

	service := NewService(mockRuleAdapter, mockAlarmAdapter, mockEntityAdapter, mockRPCClient, canopsis.AxeConnector, mockEncoder, logger)
	events, err := service.Process(ctx)

	if err != nil {
		t.Errorf("exepected no error but got %v", err)
		return
	}

	if len(events) == 0 {
		t.Errorf("exepected event but got nothing")
		return
	}

	expectedEvent := types.Event{
		EventType:     types.EventTypeAck,
		Connector:     alarm.Value.Connector,
		ConnectorName: alarm.Value.ConnectorName,
		Component:     alarm.Value.Component,
		Resource:      alarm.Value.Resource,
		Author:        canopsis.DefaultEventAuthor,
		Output:        "Idle rule: test-rule-name. Comment: test-output.",
		SourceType:    types.SourceTypeResource,
		Initiator:     types.InitiatorSystem,
		IdleRuleApply: "alarm_last_update",
	}
	event := events[0]
	event.Timestamp = datetime.CpsTime{}

	if diff := pretty.Compare(event, expectedEvent); diff != "" {
		t.Errorf("unexepected event %s", diff)
	}
}

func TestService_Process_GivenEntityRule_ShouldReturnEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mockRuleAdapter := mock_idlerule.NewMockRuleAdapter(ctrl)
	mockAlarmAdapter := mock_alarm.NewMockAdapter(ctrl)
	mockEntityAdapter := mock_entity.NewMockAdapter(ctrl)
	mockEncoder := mock_encoding.NewMockEncoder(ctrl)
	mockRPCClient := mock_engine.NewMockRPCClient(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	resource := "test-resource"
	component := "test-component"
	connector := "test-connector"
	connectorName := "test-connector-name"
	entity := types.Entity{
		ID:            resource + "/" + component,
		Type:          types.EntityTypeResource,
		Component:     component,
		Connector:     connector + "/" + connectorName,
		Name:          resource,
		LastEventDate: &datetime.CpsTime{Time: time.Now().Add(-6 * time.Hour)},
	}
	state := types.CpsNumber(types.AlarmStateCritical)
	rule := idlerule.Rule{
		ID:     "test-rule",
		Type:   idlerule.RuleTypeEntity,
		Name:   "test-rule-name",
		Author: "test-author",
		Duration: datetime.DurationWithUnit{
			Value: 10,
			Unit:  "s",
		},
		EntityPatternFields: savedpattern.EntityPatternFields{
			EntityPattern: [][]pattern.FieldCondition{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, resource),
					},
				},
			},
		},
	}

	emptyMockCursor := mock_mongo.NewMockCursor(ctrl)
	emptyMockCursor.EXPECT().Next(gomock.Any()).Return(false)
	emptyMockCursor.EXPECT().Close(gomock.Any()).Return(nil)
	mockCursor := mock_mongo.NewMockCursor(ctrl)
	firstCall := mockCursor.EXPECT().Next(gomock.Any()).Return(true)
	secondCall := mockCursor.EXPECT().Next(gomock.Any()).Return(false)
	gomock.InOrder(firstCall, secondCall)
	mockCursor.EXPECT().Close(gomock.Any())
	mockCursor.EXPECT().Decode(gomock.Any()).Do(func(v *types.Entity) {
		*v = entity
	}).Return(nil)
	mockRuleAdapter.
		EXPECT().
		GetEnabled(gomock.Any()).
		Return([]idlerule.Rule{rule}, nil)
	mockAlarmAdapter.
		EXPECT().
		GetOpenedAlarmsWithLastDatesBefore(gomock.Any(), gomock.Any()).
		Return(emptyMockCursor, nil)
	mockAlarmAdapter.EXPECT().GetLastAlarmByEntityID(gomock.Any(), gomock.Any()).Return(nil, nil)
	mockEntityAdapter.
		EXPECT().
		GetAllWithLastUpdateDateBefore(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockCursor, nil)
	mockAlarmAdapter.EXPECT().GetOpenedAlarmsByConnectorIdleRules(gomock.Any()).Return(nil, nil)

	service := NewService(mockRuleAdapter, mockAlarmAdapter, mockEntityAdapter, mockRPCClient, canopsis.AxeConnector, mockEncoder, logger)
	events, err := service.Process(ctx)

	if err != nil {
		t.Errorf("exepected no error but got %v", err)
		return
	}

	if len(events) == 0 {
		t.Errorf("exepected event but got nothing")
		return
	}

	expectedEvent := types.Event{
		EventType:     types.EventTypeNoEvents,
		State:         state,
		Connector:     canopsis.AxeConnector,
		ConnectorName: canopsis.AxeConnector,
		Component:     component,
		Resource:      resource,
		SourceType:    types.SourceTypeResource,
		Author:        canopsis.DefaultEventAuthor,
		Output:        "Idle rule: test-rule-name",
		Initiator:     types.InitiatorSystem,
		IdleRuleApply: "entity",
	}
	event := events[0]
	event.Timestamp = datetime.CpsTime{}

	if diff := pretty.Compare(event, expectedEvent); diff != "" {
		t.Errorf("unexepected event %s", diff)
	}
}

func TestService_Process_GivenAlarmForConnectorEntity_ShouldReturnEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mockRuleAdapter := mock_idlerule.NewMockRuleAdapter(ctrl)
	mockAlarmAdapter := mock_alarm.NewMockAdapter(ctrl)
	mockEntityAdapter := mock_entity.NewMockAdapter(ctrl)
	mockEncoder := mock_encoding.NewMockEncoder(ctrl)
	mockRPCClient := mock_engine.NewMockRPCClient(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	alarm := types.Alarm{
		ID: "test-alarm",
		Value: types.AlarmValue{
			Connector:     "test-connector",
			ConnectorName: "test-connector-name",
			State: &types.AlarmStep{
				Author:  canopsis.DefaultEventAuthor,
				Message: "test-message",
			},
		},
	}

	mockRuleAdapter.
		EXPECT().
		GetEnabled(gomock.Any()).
		Return(nil, nil)
	mockAlarmAdapter.EXPECT().GetOpenedAlarmsByConnectorIdleRules(gomock.Any()).Return([]types.Alarm{alarm}, nil)

	service := NewService(mockRuleAdapter, mockAlarmAdapter, mockEntityAdapter, mockRPCClient, canopsis.AxeConnector, mockEncoder, logger)
	events, err := service.Process(ctx)

	if err != nil {
		t.Errorf("exepected no error but got %v", err)
		return
	}

	if len(events) == 0 {
		t.Errorf("exepected event but got nothing")
		return
	}

	expectedEvent := types.Event{
		EventType:     types.EventTypeNoEvents,
		State:         types.AlarmStateOK,
		Connector:     alarm.Value.Connector,
		ConnectorName: alarm.Value.ConnectorName,
		Initiator:     types.InitiatorSystem,
		Output:        alarm.Value.State.Message,
		Author:        alarm.Value.State.Author,
		SourceType:    types.SourceTypeConnector,
	}
	event := events[0]
	event.Timestamp = datetime.CpsTime{}

	if !reflect.DeepEqual(event, expectedEvent) {
		t.Errorf("exepected event %+v but got %+v", expectedEvent, event)
	}
}
