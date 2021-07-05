package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metaalarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metaalarm/service"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	mock_alarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/alarm"
	mock_config "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/config"
	"github.com/golang/mock/gomock"
)

func TestCreateMetaAlarm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	maService := service.NewMetaAlarmService(
		mock_alarm.NewMockAdapter(ctrl),
		mock_config.NewMockAlarmConfigProvider(ctrl),
		log.NewTestLogger(),
	)

	event := types.Event{
		Timestamp: types.NewCpsTime(time.Now().Unix()),
		Author:    "test",
		State:     types.CpsNumber(1),
	}

	children := []types.AlarmWithEntity{
		{
			Alarm: types.Alarm{
				ID: "alarm-1",
			},
			Entity: types.Entity{
				ID: "entity-1",
			},
		},
		{
			Alarm: types.Alarm{
				ID: "alarm-2",
			},
			Entity: types.Entity{
				ID: "entity-2",
			},
		},
	}

	rule := metaalarm.Rule{ID: "rule-1", OutputTemplate: "Number of children: {{ .Count }}"}

	metaAlarmEvent, err := maService.CreateMetaAlarm(event, children, rule)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	if metaAlarmEvent.Timestamp != event.Timestamp {
		t.Errorf("expected %v, but got %v", event.Timestamp, metaAlarmEvent.Timestamp)
	}

	if metaAlarmEvent.Author != event.Author {
		t.Errorf("expected %v, but got %v", event.Author, metaAlarmEvent.Author)
	}

	if metaAlarmEvent.State != event.State {
		t.Errorf("expected %v, but got %v", event.State, metaAlarmEvent.State)
	}

	if metaAlarmEvent.Component != metaalarm.DefaultMetaAlarmComponent {
		t.Errorf("expected %v, but got %v", metaalarm.DefaultMetaAlarmComponent, metaAlarmEvent.Component)
	}

	if metaAlarmEvent.Connector != metaalarm.DefaultMetaAlarmConnector {
		t.Errorf("expected %v, but got %v", metaalarm.DefaultMetaAlarmConnector, metaAlarmEvent.Connector)
	}

	if metaAlarmEvent.ConnectorName != metaalarm.DefaultMetaAlarmConnectorName {
		t.Errorf("expected %v, but got %v", metaalarm.DefaultMetaAlarmConnectorName, metaAlarmEvent.ConnectorName)
	}

	if !strings.HasPrefix(metaAlarmEvent.Resource, metaalarm.DefaultMetaAlarmEntityPrefix) {
		t.Errorf("%v should have prefix %v", metaAlarmEvent.Resource, metaalarm.DefaultMetaAlarmEntityPrefix)
	}

	if metaAlarmEvent.SourceType != types.SourceTypeResource {
		t.Errorf("expected %v, but got %v", types.SourceTypeResource, metaAlarmEvent.SourceType)
	}

	if metaAlarmEvent.EventType != types.EventTypeMetaAlarm {
		t.Errorf("expected %v, but got %v", types.EventTypeMetaAlarm, metaAlarmEvent.EventType)
	}

	if metaAlarmEvent.MetaAlarmRuleID != rule.ID {
		t.Errorf("expected %v, but got %v", rule.ID, metaAlarmEvent.MetaAlarmRuleID)
	}

	if metaAlarmEvent.Output != fmt.Sprintf("Number of children: %d", len(children)) {
		t.Errorf("expected %v, but got %v", fmt.Sprintf("Number of children: %d", len(children)), metaAlarmEvent.Output)
	}

	meta, ok := metaAlarmEvent.ExtraInfos["Meta"]
	if !ok {
		t.Error("expected to have Meta key in ExtraInfos map, but it doesn't")
	}

	extraInfosMeta, ok := meta.(service.EventExtraInfosMeta)
	if !ok {
		t.Error("expected ExtraInfos[\"Meta\"] value to be EventExtraInfosMeta")
	}

	if extraInfosMeta.Count != int64(len(children)) {
		t.Errorf("expected %v, but got %v", len(children), extraInfosMeta.Count)
	}

	if extraInfosMeta.Rule.ID != rule.ID {
		t.Errorf("expected %v, but got %v", rule.ID, extraInfosMeta.Rule.ID)
	}

	if extraInfosMeta.Children.Alarm.ID != children[len(children)-1].Alarm.ID {
		t.Errorf("expected %v, but got %v", children[len(children)-1].Alarm.ID, extraInfosMeta.Children.Alarm.ID)
	}

	if extraInfosMeta.Children.Entity.ID != children[len(children)-1].Entity.ID {
		t.Errorf("expected %v, but got %v", children[len(children)-1].Entity.ID, extraInfosMeta.Children.Entity.ID)
	}

	for i := 0; i < len(children); i++ {
		if children[i].Entity.ID != (*metaAlarmEvent.MetaAlarmChildren)[i] {
			t.Errorf("expected %v, but got %v", children[i].Entity.ID, (*metaAlarmEvent.MetaAlarmChildren)[i])
		}
	}

	_, err = maService.CreateMetaAlarm(event, []types.AlarmWithEntity{}, rule)
	if !errors.Is(err, metaalarm.ErrNoChildren) {
		t.Fatalf("expected err %v but got %v", metaalarm.ErrNoChildren, err)
	}
}

func TestAddChildToMetaAlarm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	alarmAdapterMock := mock_alarm.NewMockAdapter(ctrl)
	alarmConfigProviderMock := mock_config.NewMockAlarmConfigProvider(ctrl)

	s := service.NewMetaAlarmService(alarmAdapterMock, alarmConfigProviderMock, log.NewTestLogger())

	rule := metaalarm.Rule{ID: "test-rule", OutputTemplate: "Number of children: {{ .Count }}"}
	event := types.Event{
		Timestamp: types.NewCpsTime(time.Now().Unix()),
		Author:    "test",
		State:     types.CpsNumber(1),
	}

	existedChild := types.Alarm{EntityID: "existed-child"}
	existedChild.Value.Parents = []string{"metaalarm"}

	metaAlarm := types.Alarm{
		EntityID: "metaalarm",
		Value: types.AlarmValue{
			Children: []string{"existed-child"},
			Connector: "meta-alarm-connector",
			ConnectorName: "meta-alarm-connector-name",
			Component: "meta-alarm-component",
			Resource: "meta-alarm-resource",
			Meta: rule.ID,
		},
	}

	newChild := types.Alarm{EntityID: "new-child"}
	newChildWithEntity := types.AlarmWithEntity{
		Alarm: newChild,
		Entity: types.Entity{ID: newChild.EntityID},
	}

	expectedChildren := []types.AlarmWithEntity{
		{
			Alarm: existedChild,
			Entity: types.Entity{ID: existedChild.EntityID},
		},
		{
			Alarm: newChild,
			Entity: types.Entity{ID: newChild.EntityID},
		},
	}

	expectedUpdates := []types.Alarm{newChild, metaAlarm}
	alarmAdapterMock.
		EXPECT().
		PartialMassUpdateOpen(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, alarms []types.Alarm) error {
			for idx, alarm := range alarms {
				if alarm.EntityID != expectedUpdates[idx].EntityID {
					t.Errorf("expected %v, but got %v", expectedUpdates[idx].EntityID, alarm.EntityID)
				}
			}

			return nil
		})
	alarmAdapterMock.EXPECT().GetCountOpenedAlarmsByIDs(gomock.Any()).Return(int64(2), nil)

	metaAlarmEvent, err := s.AddChildToMetaAlarm(context.Background(), event, metaAlarm, newChildWithEntity, rule)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	if metaAlarmEvent.Timestamp != event.Timestamp {
		t.Errorf("expected %v, but got %v", event.Timestamp, metaAlarmEvent.Timestamp)
	}

	if metaAlarmEvent.Author != event.Author {
		t.Errorf("expected %v, but got %v", event.Author, metaAlarmEvent.Author)
	}

	if metaAlarmEvent.Component != metaAlarm.Value.Component {
		t.Errorf("expected %v, but got %v", metaAlarm.Value.Component, metaAlarmEvent.Component)
	}

	if metaAlarmEvent.Connector != metaAlarm.Value.Connector {
		t.Errorf("expected %v, but got %v", metaAlarm.Value.Connector, metaAlarmEvent.Connector)
	}

	if metaAlarmEvent.ConnectorName != metaAlarm.Value.ConnectorName {
		t.Errorf("expected %v, but got %v", metaAlarm.Value.ConnectorName, metaAlarmEvent.ConnectorName)
	}

	if metaAlarmEvent.Resource != metaAlarm.Value.Resource {
		t.Errorf("expected %v, but got %v", metaAlarm.Value.Resource, metaAlarmEvent.ConnectorName)
	}

	if metaAlarmEvent.SourceType != types.SourceTypeResource {
		t.Errorf("expected %v, but got %v", types.SourceTypeResource, metaAlarmEvent.SourceType)
	}

	if metaAlarmEvent.EventType != types.EventTypeMetaAlarmUpdated {
		t.Errorf("expected %v, but got %v", types.EventTypeMetaAlarmUpdated, metaAlarmEvent.EventType)
	}

	if metaAlarmEvent.MetaAlarmRuleID != rule.ID {
		t.Errorf("expected %v, but got %v", rule.ID, metaAlarmEvent.MetaAlarmRuleID)
	}

	if metaAlarmEvent.Output != fmt.Sprintf("Number of children: %d", len(expectedChildren)) {
		t.Errorf("expected %v, but got %v", fmt.Sprintf("Number of children: %d", len(expectedChildren)), metaAlarmEvent.Output)
	}

	meta, ok := metaAlarmEvent.ExtraInfos["Meta"]
	if !ok {
		t.Error("expected to have Meta key in ExtraInfos map, but it doesn't")
	}

	extraInfosMeta, ok := meta.(service.EventExtraInfosMeta)
	if !ok {
		t.Error("expected ExtraInfos[\"Meta\"] value to be EventExtraInfosMeta")
	}

	if extraInfosMeta.Count != int64(len(expectedChildren)) {
		t.Errorf("expected %v, but got %v", len(expectedChildren), extraInfosMeta.Count)
	}

	if extraInfosMeta.Rule.ID != rule.ID {
		t.Errorf("expected %v, but got %v", rule.ID, extraInfosMeta.Rule.ID)
	}

	if extraInfosMeta.Children.Alarm.ID != expectedChildren[len(expectedChildren)-1].Alarm.ID {
		t.Errorf("expected %v, but got %v", expectedChildren[len(expectedChildren)-1].Alarm.ID, extraInfosMeta.Children.Alarm.ID)
	}

	if extraInfosMeta.Children.Entity.ID != expectedChildren[len(expectedChildren)-1].Entity.ID {
		t.Errorf("expected %v, but got %v", expectedChildren[len(expectedChildren)-1].Entity.ID, extraInfosMeta.Children.Entity.ID)
	}

	for i := 0; i < len(expectedChildren); i++ {
		if expectedChildren[i].Entity.ID != (*metaAlarmEvent.MetaAlarmChildren)[i] {
			t.Errorf("expected %v, but got %v", expectedChildren[i].Entity.ID, (*metaAlarmEvent.MetaAlarmChildren)[i])
		}
	}
}

func TestAddChildToMetaAlarmComponent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	alarmAdapterMock := mock_alarm.NewMockAdapter(ctrl)
	alarmConfigProviderMock := mock_config.NewMockAlarmConfigProvider(ctrl)

	s := service.NewMetaAlarmService(alarmAdapterMock, alarmConfigProviderMock, log.NewTestLogger())

	rule := metaalarm.Rule{ID: "test-rule", OutputTemplate: "Number of children: {{ .Count }}"}
	event := types.Event{
		Timestamp: types.NewCpsTime(time.Now().Unix()),
		Author:    "test",
		State:     types.CpsNumber(1),
	}

	existedChild := types.Alarm{EntityID: "existed-child"}
	existedChild.Value.Parents = []string{"metaalarm"}

	metaAlarm := types.Alarm{
		EntityID: "metaalarm",
		Value: types.AlarmValue{
			Children: []string{"existed-child"},
			Connector: "meta-alarm-connector",
			ConnectorName: "meta-alarm-connector-name",
			Component: "meta-alarm-component",
			Meta: rule.ID,
		},
	}

	newChild := types.Alarm{EntityID: "new-child"}
	newChildWithEntity := types.AlarmWithEntity{
		Alarm: newChild,
		Entity: types.Entity{ID: newChild.EntityID},
	}

	expectedUpdates := []types.Alarm{newChild, metaAlarm}
	alarmAdapterMock.
		EXPECT().
		PartialMassUpdateOpen(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, alarms []types.Alarm) error {
			for idx, alarm := range alarms {
				if alarm.EntityID != expectedUpdates[idx].EntityID {
					t.Errorf("expected %v, but got %v", expectedUpdates[idx].EntityID, alarm.EntityID)
				}
			}

			return nil
		})
	alarmAdapterMock.EXPECT().GetCountOpenedAlarmsByIDs(gomock.Any()).Return(int64(2), nil)

	metaAlarmEvent, err := s.AddChildToMetaAlarm(context.Background(), event, metaAlarm, newChildWithEntity, rule)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	if metaAlarmEvent.SourceType != types.SourceTypeComponent {
		t.Errorf("expected %v, but got %v", types.SourceTypeComponent, metaAlarmEvent.SourceType)
	}
}

func TestAddMultipleChildrenToMetaAlarm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	alarmAdapterMock := mock_alarm.NewMockAdapter(ctrl)
	alarmConfigProviderMock := mock_config.NewMockAlarmConfigProvider(ctrl)

	s := service.NewMetaAlarmService(alarmAdapterMock, alarmConfigProviderMock, log.NewTestLogger())

	rule := metaalarm.Rule{ID: "test-rule", OutputTemplate: "Number of children: {{ .Count }}"}
	event := types.Event{
		Timestamp: types.NewCpsTime(time.Now().Unix()),
		Author:    "test",
		State:     types.CpsNumber(1),
	}

	existedChild := types.Alarm{EntityID: "existed-child"}
	existedChild.Value.Parents = []string{"metaalarm"}

	metaAlarm := types.Alarm{
		EntityID: "metaalarm",
		Value: types.AlarmValue{
			Children: []string{"existed-child"},
			Connector: "meta-alarm-connector",
			ConnectorName: "meta-alarm-connector-name",
			Component: "meta-alarm-component",
			Resource: "meta-alarm-resource",
			Meta: rule.ID,
		},
	}

	newChild1 := types.Alarm{EntityID: "new-child-1"}
	newChild2 := types.Alarm{EntityID: "new-child-2"}
	newChildrenWithEntity := []types.AlarmWithEntity{
		{
			Alarm: newChild1,
			Entity: types.Entity{ID: newChild1.EntityID},
		},
		{
			Alarm: newChild2,
			Entity: types.Entity{ID: newChild2.EntityID},
		},
	}

	expectedChildren := []types.AlarmWithEntity{
		{
			Alarm: existedChild,
			Entity: types.Entity{ID: existedChild.EntityID},
		},
		{
			Alarm: newChild1,
			Entity: types.Entity{ID: newChild1.EntityID},
		},
		{
			Alarm: newChild2,
			Entity: types.Entity{ID: newChild2.EntityID},
		},
	}

	expectedUpdates := []types.Alarm{newChild1, newChild2, metaAlarm}
	alarmAdapterMock.
		EXPECT().
		PartialMassUpdateOpen(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, alarms []types.Alarm) error {
			for idx, alarm := range alarms {
				if alarm.EntityID != expectedUpdates[idx].EntityID {
					t.Errorf("expected %v, but got %v", expectedUpdates[idx].EntityID, alarm.EntityID)
				}
			}

			return nil
		})
	alarmAdapterMock.EXPECT().GetCountOpenedAlarmsByIDs(gomock.Any()).Return(int64(3), nil)

	metaAlarmEvent, err := s.AddMultipleChildsToMetaAlarm(context.Background(), event, metaAlarm, newChildrenWithEntity, rule)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	if metaAlarmEvent.Timestamp != event.Timestamp {
		t.Errorf("expected %v, but got %v", event.Timestamp, metaAlarmEvent.Timestamp)
	}

	if metaAlarmEvent.Author != event.Author {
		t.Errorf("expected %v, but got %v", event.Author, metaAlarmEvent.Author)
	}

	if metaAlarmEvent.Component != metaAlarm.Value.Component {
		t.Errorf("expected %v, but got %v", metaAlarm.Value.Component, metaAlarmEvent.Component)
	}

	if metaAlarmEvent.Connector != metaAlarm.Value.Connector {
		t.Errorf("expected %v, but got %v", metaAlarm.Value.Connector, metaAlarmEvent.Connector)
	}

	if metaAlarmEvent.ConnectorName != metaAlarm.Value.ConnectorName {
		t.Errorf("expected %v, but got %v", metaAlarm.Value.ConnectorName, metaAlarmEvent.ConnectorName)
	}

	if metaAlarmEvent.Resource != metaAlarm.Value.Resource {
		t.Errorf("expected %v, but got %v", metaAlarm.Value.Resource, metaAlarmEvent.ConnectorName)
	}

	if metaAlarmEvent.SourceType != types.SourceTypeResource {
		t.Errorf("expected %v, but got %v", types.SourceTypeResource, metaAlarmEvent.SourceType)
	}

	if metaAlarmEvent.EventType != types.EventTypeMetaAlarmUpdated {
		t.Errorf("expected %v, but got %v", types.EventTypeMetaAlarmUpdated, metaAlarmEvent.EventType)
	}

	if metaAlarmEvent.MetaAlarmRuleID != rule.ID {
		t.Errorf("expected %v, but got %v", rule.ID, metaAlarmEvent.MetaAlarmRuleID)
	}

	if metaAlarmEvent.Output != fmt.Sprintf("Number of children: %d", len(expectedChildren)) {
		t.Errorf("expected %v, but got %v", fmt.Sprintf("Number of children: %d", len(expectedChildren)), metaAlarmEvent.Output)
	}

	meta, ok := metaAlarmEvent.ExtraInfos["Meta"]
	if !ok {
		t.Error("expected to have Meta key in ExtraInfos map, but it doesn't")
	}

	extraInfosMeta, ok := meta.(service.EventExtraInfosMeta)
	if !ok {
		t.Error("expected ExtraInfos[\"Meta\"] value to be EventExtraInfosMeta")
	}

	if extraInfosMeta.Count != int64(len(expectedChildren)) {
		t.Errorf("expected %v, but got %v", len(expectedChildren), extraInfosMeta.Count)
	}

	if extraInfosMeta.Rule.ID != rule.ID {
		t.Errorf("expected %v, but got %v", rule.ID, extraInfosMeta.Rule.ID)
	}

	if extraInfosMeta.Children.Alarm.ID != expectedChildren[len(expectedChildren)-1].Alarm.ID {
		t.Errorf("expected %v, but got %v", expectedChildren[len(expectedChildren)-1].Alarm.ID, extraInfosMeta.Children.Alarm.ID)
	}

	if extraInfosMeta.Children.Entity.ID != expectedChildren[len(expectedChildren)-1].Entity.ID {
		t.Errorf("expected %v, but got %v", expectedChildren[len(expectedChildren)-1].Entity.ID, extraInfosMeta.Children.Entity.ID)
	}

	for i := 0; i < len(expectedChildren); i++ {
		if expectedChildren[i].Entity.ID != (*metaAlarmEvent.MetaAlarmChildren)[i] {
			t.Errorf("expected %v, but got %v", expectedChildren[i].Entity.ID, (*metaAlarmEvent.MetaAlarmChildren)[i])
		}
	}
}

func TestAddMultipleChildrenToMetaAlarmComponent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	alarmAdapterMock := mock_alarm.NewMockAdapter(ctrl)
	alarmConfigProviderMock := mock_config.NewMockAlarmConfigProvider(ctrl)

	s := service.NewMetaAlarmService(alarmAdapterMock, alarmConfigProviderMock, log.NewTestLogger())

	rule := metaalarm.Rule{ID: "test-rule", OutputTemplate: "Number of children: {{ .Count }}"}
	event := types.Event{
		Timestamp: types.NewCpsTime(time.Now().Unix()),
		Author:    "test",
		State:     types.CpsNumber(1),
	}

	existedChild := types.Alarm{EntityID: "existed-child"}
	existedChild.Value.Parents = []string{"metaalarm"}

	metaAlarm := types.Alarm{
		EntityID: "metaalarm",
		Value: types.AlarmValue{
			Children: []string{"existed-child"},
			Connector: "meta-alarm-connector",
			ConnectorName: "meta-alarm-connector-name",
			Component: "meta-alarm-component",
			Meta: rule.ID,
		},
	}

	newChild1 := types.Alarm{EntityID: "new-child-1"}
	newChild2 := types.Alarm{EntityID: "new-child-2"}
	newChildrenWithEntity := []types.AlarmWithEntity{
		{
			Alarm: newChild1,
			Entity: types.Entity{ID: newChild1.EntityID},
		},
		{
			Alarm: newChild2,
			Entity: types.Entity{ID: newChild2.EntityID},
		},
	}

	expectedUpdates := []types.Alarm{newChild1, newChild2, metaAlarm}
	alarmAdapterMock.
		EXPECT().
		PartialMassUpdateOpen(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, alarms []types.Alarm) error {
			for idx, alarm := range alarms {
				if alarm.EntityID != expectedUpdates[idx].EntityID {
					t.Errorf("expected %v, but got %v", expectedUpdates[idx].EntityID, alarm.EntityID)
				}
			}

			return nil
		})
	alarmAdapterMock.EXPECT().GetCountOpenedAlarmsByIDs(gomock.Any()).Return(int64(3), nil)

	metaAlarmEvent, err := s.AddMultipleChildsToMetaAlarm(context.Background(), event, metaAlarm, newChildrenWithEntity, rule)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	if metaAlarmEvent.SourceType != types.SourceTypeComponent {
		t.Errorf("expected %v, but got %v", types.SourceTypeComponent, metaAlarmEvent.SourceType)
	}
}

func TestRemoveMultipleChildToMetaAlarm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	alarmAdapterMock := mock_alarm.NewMockAdapter(ctrl)
	alarmConfigProviderMock := mock_config.NewMockAlarmConfigProvider(ctrl)

	s := service.NewMetaAlarmService(alarmAdapterMock, alarmConfigProviderMock, log.NewTestLogger())

	rule := metaalarm.Rule{ID: "test-rule", OutputTemplate: "Number of children: {{ .Count }}"}
	event := types.Event{
		Timestamp: types.NewCpsTime(time.Now().Unix()),
		Author:    "test",
		State:     types.CpsNumber(1),
	}

	childAlarm1, childAlarm2, childAlarm3, childAlarm4 := types.Alarm{EntityID: "child-1"}, types.Alarm{EntityID: "child-2"}, types.Alarm{EntityID: "child-3"}, types.Alarm{EntityID: "child-4"}
	metaAlarm := types.Alarm{
		EntityID: "parent-1",
		Value: types.AlarmValue{
			Connector: "meta-alarm-connector",
			ConnectorName: "meta-alarm-connector-name",
			Component: "meta-alarm-component",
			Resource: "meta-alarm-resource",
			Meta: rule.ID,
		},
	}
	metaAlarm.AddChild(childAlarm1.EntityID)
	metaAlarm.AddChild(childAlarm2.EntityID)
	metaAlarm.AddChild(childAlarm3.EntityID)
	metaAlarm.AddChild(childAlarm4.EntityID)
	childAlarm1.AddParent(metaAlarm.EntityID)
	childAlarm2.AddParent(metaAlarm.EntityID)
	childAlarm3.AddParent(metaAlarm.EntityID)
	childAlarm4.AddParent(metaAlarm.EntityID)

	removeChildrenWithEntity := []types.AlarmWithEntity{
		{
			Alarm: childAlarm2,
			Entity: types.Entity{ID: childAlarm2.EntityID},
		},
		{
			Alarm: childAlarm3,
			Entity: types.Entity{ID: childAlarm3.EntityID},
		},
	}

	expectedChildren := []types.AlarmWithEntity{
		{
			Alarm: childAlarm1,
			Entity: types.Entity{ID: childAlarm1.EntityID},
		},
		{
			Alarm: childAlarm4,
			Entity: types.Entity{ID: childAlarm4.EntityID},
		},
	}

	expectedUpdates := []types.Alarm{childAlarm2, childAlarm3, metaAlarm}
	alarmAdapterMock.
		EXPECT().
		PartialMassUpdateOpen(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, alarms []types.Alarm) error {
			for idx, alarm := range alarms {
				if alarm.EntityID != expectedUpdates[idx].EntityID {
					t.Errorf("expected %v, but got %v", expectedUpdates[idx].EntityID, alarm.EntityID)
				}
			}

			return nil
		})
	alarmAdapterMock.EXPECT().GetOpenedAlarmsWithEntityByIDs(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ids []string, alarms *[]types.AlarmWithEntity) error {
			*alarms = expectedChildren

			return nil
		})

	metaAlarmEvent, err := s.RemoveMultipleChildToMetaAlarm(context.Background(), event, metaAlarm, removeChildrenWithEntity, rule)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	if metaAlarmEvent.Timestamp != event.Timestamp {
		t.Errorf("expected %v, but got %v", event.Timestamp, metaAlarmEvent.Timestamp)
	}

	if metaAlarmEvent.Author != event.Author {
		t.Errorf("expected %v, but got %v", event.Author, metaAlarmEvent.Author)
	}

	if metaAlarmEvent.Component != metaAlarm.Value.Component {
		t.Errorf("expected %v, but got %v", metaAlarm.Value.Component, metaAlarmEvent.Component)
	}

	if metaAlarmEvent.Connector != metaAlarm.Value.Connector {
		t.Errorf("expected %v, but got %v", metaAlarm.Value.Connector, metaAlarmEvent.Connector)
	}

	if metaAlarmEvent.ConnectorName != metaAlarm.Value.ConnectorName {
		t.Errorf("expected %v, but got %v", metaAlarm.Value.ConnectorName, metaAlarmEvent.ConnectorName)
	}

	if metaAlarmEvent.Resource != metaAlarm.Value.Resource {
		t.Errorf("expected %v, but got %v", metaAlarm.Value.Resource, metaAlarmEvent.ConnectorName)
	}

	if metaAlarmEvent.SourceType != types.SourceTypeResource {
		t.Errorf("expected %v, but got %v", types.SourceTypeResource, metaAlarmEvent.SourceType)
	}

	if metaAlarmEvent.EventType != types.EventTypeMetaAlarmUpdated {
		t.Errorf("expected %v, but got %v", types.EventTypeMetaAlarmUpdated, metaAlarmEvent.EventType)
	}

	if metaAlarmEvent.MetaAlarmRuleID != rule.ID {
		t.Errorf("expected %v, but got %v", rule.ID, metaAlarmEvent.MetaAlarmRuleID)
	}

	if metaAlarmEvent.Output != fmt.Sprintf("Number of children: %d", len(expectedChildren)) {
		t.Errorf("expected %v, but got %v", fmt.Sprintf("Number of children: %d", len(expectedChildren)), metaAlarmEvent.Output)
	}

	meta, ok := metaAlarmEvent.ExtraInfos["Meta"]
	if !ok {
		t.Error("expected to have Meta key in ExtraInfos map, but it doesn't")
	}

	extraInfosMeta, ok := meta.(service.EventExtraInfosMeta)
	if !ok {
		t.Error("expected ExtraInfos[\"Meta\"] value to be EventExtraInfosMeta")
	}

	if extraInfosMeta.Count != int64(len(expectedChildren)) {
		t.Errorf("expected %v, but got %v", len(expectedChildren), extraInfosMeta.Count)
	}

	if extraInfosMeta.Rule.ID != rule.ID {
		t.Errorf("expected %v, but got %v", rule.ID, extraInfosMeta.Rule.ID)
	}

	if extraInfosMeta.Children.Alarm.ID != expectedChildren[len(expectedChildren)-1].Alarm.ID {
		t.Errorf("expected %v, but got %v", expectedChildren[len(expectedChildren)-1].Alarm.ID, extraInfosMeta.Children.Alarm.ID)
	}

	if extraInfosMeta.Children.Entity.ID != expectedChildren[len(expectedChildren)-1].Entity.ID {
		t.Errorf("expected %v, but got %v", expectedChildren[len(expectedChildren)-1].Entity.ID, extraInfosMeta.Children.Entity.ID)
	}

	for i := 0; i < len(expectedChildren); i++ {
		if expectedChildren[i].Entity.ID != (*metaAlarmEvent.MetaAlarmChildren)[i] {
			t.Errorf("expected %v, but got %v", expectedChildren[i].Entity.ID, (*metaAlarmEvent.MetaAlarmChildren)[i])
		}
	}
}

func TestRemoveMultipleChildToMetaAlarmComponent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	alarmAdapterMock := mock_alarm.NewMockAdapter(ctrl)
	alarmConfigProviderMock := mock_config.NewMockAlarmConfigProvider(ctrl)

	s := service.NewMetaAlarmService(alarmAdapterMock, alarmConfigProviderMock, log.NewTestLogger())

	rule := metaalarm.Rule{ID: "test-rule", OutputTemplate: "Number of children: {{ .Count }}"}
	event := types.Event{
		Timestamp: types.NewCpsTime(time.Now().Unix()),
		Author:    "test",
		State:     types.CpsNumber(1),
	}

	childAlarm1, childAlarm2, childAlarm3, childAlarm4 := types.Alarm{EntityID: "child-1"}, types.Alarm{EntityID: "child-2"}, types.Alarm{EntityID: "child-3"}, types.Alarm{EntityID: "child-4"}
	metaAlarm := types.Alarm{
		EntityID: "parent-1",
		Value: types.AlarmValue{
			Children: []string{"existed-child"},
			Connector: "meta-alarm-connector",
			ConnectorName: "meta-alarm-connector-name",
			Component: "meta-alarm-component",
			Meta: rule.ID,
		},
	}
	metaAlarm.AddChild(childAlarm1.EntityID)
	metaAlarm.AddChild(childAlarm2.EntityID)
	metaAlarm.AddChild(childAlarm3.EntityID)
	metaAlarm.AddChild(childAlarm4.EntityID)
	childAlarm1.AddParent(metaAlarm.EntityID)
	childAlarm2.AddParent(metaAlarm.EntityID)
	childAlarm3.AddParent(metaAlarm.EntityID)
	childAlarm4.AddParent(metaAlarm.EntityID)

	removeChildrenWithEntity := []types.AlarmWithEntity{
		{
			Alarm: childAlarm2,
			Entity: types.Entity{ID: childAlarm2.EntityID},
		},
		{
			Alarm: childAlarm3,
			Entity: types.Entity{ID: childAlarm3.EntityID},
		},
	}

	expectedChildren := []types.AlarmWithEntity{
		{
			Alarm: childAlarm1,
			Entity: types.Entity{ID: childAlarm1.EntityID},
		},
		{
			Alarm: childAlarm4,
			Entity: types.Entity{ID: childAlarm4.EntityID},
		},
	}

	expectedUpdates := []types.Alarm{childAlarm2, childAlarm3, metaAlarm}
	alarmAdapterMock.
		EXPECT().
		PartialMassUpdateOpen(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, alarms []types.Alarm) error {
			for idx, alarm := range alarms {
				if alarm.EntityID != expectedUpdates[idx].EntityID {
					t.Errorf("expected %v, but got %v", expectedUpdates[idx].EntityID, alarm.EntityID)
				}
			}

			return nil
		})
	alarmAdapterMock.EXPECT().GetOpenedAlarmsWithEntityByIDs(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ids []string, alarms *[]types.AlarmWithEntity) error {
			*alarms = expectedChildren

			return nil
		})

	metaAlarmEvent, err := s.RemoveMultipleChildToMetaAlarm(context.Background(), event, metaAlarm, removeChildrenWithEntity, rule)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	if metaAlarmEvent.SourceType != types.SourceTypeComponent {
		t.Errorf("expected %v, but got %v", types.SourceTypeComponent, metaAlarmEvent.SourceType)
	}
}

func TestAddChildToMetaAlarmWorstState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	alarmAdapterMock := mock_alarm.NewMockAdapter(ctrl)
	alarmConfigProviderMock := mock_config.NewMockAlarmConfigProvider(ctrl)
	alarmConfigProviderMock.EXPECT().Get().Return(config.AlarmConfig{})

	const expectedState = types.AlarmStateCritical

	s := service.NewMetaAlarmService(alarmAdapterMock, alarmConfigProviderMock, log.NewTestLogger())

	alarmJSON := []byte(`{
		"_id" : "40b8aaef-a975-491a-a143-2dfdabe059c8",
		"t" : 1594203725,
		"d" : "meta-alarm-entity-24cc5fff-6d0e-4f64-ae91-7bd4a0cdc9b9/metaalarm",
		"v" : {
			"state" : {
				"_t" : "stateinc",
				"t" : 1594203725,
				"a" : "engine.correlation",
				"m" : "testvalue alarm component",
				"val" : 2
			},
			"status" : {
				"_t" : "statusinc",
				"t" : 1594203725,
				"a" : "engine.correlation",
				"m" : "testvalue alarm component",
				"val" : 1
			},
			"steps" : [ 
				{
					"_t" : "statusinc",
					"t" : 1594203725,
					"a" : "engine.correlation",
					"m" : "testvalue alarm component",
					"val" : 2
				}
			],
			"component" : "metaalarm",
			"connector" : "engine",
			"connector_name" : "correlation",
			"creation_date" : 1594203725,
			"display_name" : "WP-NZ-KV",
			"initial_output" : "testvalue alarm component",
			"output" : "testvalue alarm component",
			"last_update_date" : 1594203725,
			"last_event_date" : 1594203725,
			"resource" : "meta-alarm-entity-24cc5fff-6d0e-4f64-ae91-7bd4a0cdc9b9",
			"meta" : "testRule",
			"parents" : [],
			"children" : [ 
				"testvaluecomponent922-110", 
				"testvaluecomponent922-112", 
				"testvaluecomponent922-A311"
			],
			"total_state_changes" : 0,
			"extra" : {},
			"infos" : {},
			"infos_rule_version" : {}
		}
	}`)
	var metaAlarm types.Alarm
	err := json.Unmarshal(alarmJSON, &metaAlarm)
	if err != nil {
		t.Fatalf("MetaAlarm unmarshal error %s", err)
	}

	alarmJSON = []byte(`{
		"_id" : "1e6e41ce-3823-4a44-9ee5-c31b7a9fd812",
		"t" : 1594204048,
		"d" : "testvalueressource311/testvaluecomponent922-A311",
		"v" : {
			"state" : {
				"_t" : "stateinc",
				"t" : 1594204048,
				"a" : "testvalueconnector.testvalueconnectorname",
				"m" : "testvalue alarm",
				"val" : 3
			},
			"status" : {
				"_t" : "statusinc",
				"t" : 1594204048,
				"a" : "testvalueconnector.testvalueconnectorname",
				"m" : "testvalue alarm",
				"val" : 1
			},
			"steps" : [ 
				{
					"_t" : "stateinc",
					"t" : 1594204048,
					"a" : "testvalueconnector.testvalueconnectorname",
					"m" : "testvalue alarm",
					"val" : 3
				}, 
				{
					"_t" : "statusinc",
					"t" : 1594204048,
					"a" : "testvalueconnector.testvalueconnectorname",
					"m" : "testvalue alarm",
					"val" : 1
				}
			],
			"component" : "testvaluecomponent922-A311",
			"connector" : "testvalueconnector",
			"connector_name" : "testvalueconnectorname",
			"creation_date" : 1594204048,
			"display_name" : "VX-LV-EA",
			"initial_output" : "testvalue alarm",
			"output" : "testvalue alarm",
			"initial_long_output" : "",
			"last_update_date" : 1594204048,
			"last_event_date" : 1594204048,
			"resource" : "testvalueressource311",
			"parents" : [ 
				"testvaluecomponent922-A311"
			],
			"children" : [],
			"total_state_changes" : 1,
			"extra" : {},
			"infos" : {},
			"infos_rule_version" : {}
		}
	}`)

	var child types.AlarmWithEntity
	err = json.Unmarshal(alarmJSON, &child.Alarm)
	if err != nil {
		t.Fatalf("Child Alarm unmarshal error %s", err)
	}

	rule := metaalarm.Rule{}
	event := types.Event{}

	var updatedAlarms *[]types.Alarm

	alarmAdapterMock.
		EXPECT().
		PartialMassUpdateOpen(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, alarms []types.Alarm) error {
			updatedAlarms = &alarms
			return nil
		})
	alarmAdapterMock.EXPECT().GetCountOpenedAlarmsByIDs(gomock.Any()).Return(int64(1), nil)

	updateEvent, err := s.AddChildToMetaAlarm(context.Background(), event, metaAlarm, child, rule)
	if err != nil {
		t.Fatalf("AddChildToMetaAlarm error %s", err)
	}

	if state := (*updatedAlarms)[0].CurrentState(); state != expectedState {
		t.Fatalf("wrong state value %d", state)
	}

	if updateEvent.EventType != types.EventTypeMetaAlarmUpdated {
		t.Fatalf("AddChildToMetaAlarm event %v", updateEvent)
	}
}

func TestAddMultipleChildrenToMetaAlarmWorstState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	alarmAdapterMock := mock_alarm.NewMockAdapter(ctrl)
	alarmConfigProviderMock := mock_config.NewMockAlarmConfigProvider(ctrl)
	alarmConfigProviderMock.EXPECT().Get().Return(config.AlarmConfig{})

	const expectedState = types.AlarmStateCritical

	s := service.NewMetaAlarmService(alarmAdapterMock, alarmConfigProviderMock, log.NewTestLogger())

	alarmJSON := []byte(`{
		"_id" : "40b8aaef-a975-491a-a143-2dfdabe059c8",
		"t" : 1594203725,
		"d" : "meta-alarm-entity-24cc5fff-6d0e-4f64-ae91-7bd4a0cdc9b9/metaalarm",
		"v" : {
			"state" : {
				"_t" : "stateinc",
				"t" : 1594203725,
				"a" : "engine.correlation",
				"m" : "testvalue alarm component",
				"val" : 1
			},
			"status" : {
				"_t" : "statusinc",
				"t" : 1594203725,
				"a" : "engine.correlation",
				"m" : "testvalue alarm component",
				"val" : 1
			},
			"steps" : [ 
				{
					"_t" : "statusinc",
					"t" : 1594203725,
					"a" : "engine.correlation",
					"m" : "testvalue alarm component",
					"val" : 2
				}
			],
			"component" : "metaalarm",
			"connector" : "engine",
			"connector_name" : "correlation",
			"creation_date" : 1594203725,
			"display_name" : "WP-NZ-KV",
			"initial_output" : "testvalue alarm component",
			"output" : "testvalue alarm component",
			"last_update_date" : 1594203725,
			"last_event_date" : 1594203725,
			"resource" : "meta-alarm-entity-24cc5fff-6d0e-4f64-ae91-7bd4a0cdc9b9",
			"meta" : "testRule",
			"parents" : [],
			"children" : [ 
				"testvaluecomponent922-112", 
				"testvaluecomponent922-A311"
			],
			"total_state_changes" : 0,
			"extra" : {},
			"infos" : {},
			"infos_rule_version" : {}
		}
	}`)
	var metaAlarm types.Alarm
	err := json.Unmarshal(alarmJSON, &metaAlarm)
	if err != nil {
		t.Fatalf("MetaAlarm unmarshal error %s", err)
	}

	alarmJSON = []byte(`{
		"_id" : "1e6e41ce-3823-4a44-9ee5-c31b7a9fd812",
		"t" : 1594204048,
		"d" : "testvalueressource311/testvaluecomponent922-A311",
		"v" : {
			"state" : {
				"_t" : "stateinc",
				"t" : 1594204048,
				"a" : "testvalueconnector.testvalueconnectorname",
				"m" : "testvalue alarm",
				"val" : 2
			},
			"status" : {
				"_t" : "statusinc",
				"t" : 1594204048,
				"a" : "testvalueconnector.testvalueconnectorname",
				"m" : "testvalue alarm",
				"val" : 1
			},
			"steps" : [ 
				{
					"_t" : "stateinc",
					"t" : 1594204048,
					"a" : "testvalueconnector.testvalueconnectorname",
					"m" : "testvalue alarm",
					"val" : 2
				}, 
				{
					"_t" : "statusinc",
					"t" : 1594204048,
					"a" : "testvalueconnector.testvalueconnectorname",
					"m" : "testvalue alarm",
					"val" : 1
				}
			],
			"component" : "testvaluecomponent922-A311",
			"connector" : "testvalueconnector",
			"connector_name" : "testvalueconnectorname",
			"creation_date" : 1594204048,
			"display_name" : "VX-LV-EA",
			"initial_output" : "testvalue alarm",
			"output" : "testvalue alarm",
			"initial_long_output" : "",
			"last_update_date" : 1594204048,
			"last_event_date" : 1594204048,
			"resource" : "testvalueressource311",
			"parents" : [ 
				"testvaluecomponent922-A311"
			],
			"children" : [],
			"total_state_changes" : 1,
			"extra" : {},
			"infos" : {},
			"infos_rule_version" : {}
		}
	}`)

	var child types.AlarmWithEntity
	err = json.Unmarshal(alarmJSON, &child.Alarm)
	if err != nil {
		t.Fatalf("Child Alarm unmarshal error %s", err)
	}

	children := []types.AlarmWithEntity{child}

	alarmJSON = []byte(`{
		"_id" : "9f327dd2-763f-4997-9a09-8f3a4b94e155",
		"t" : 1594203725,
		"d" : "testvaluecomponent922-110",
		"v" : {
			"state" : {
				"_t" : "stateinc",
				"t" : 1594203725,
				"a" : "testvalueconnector.testvalueconnectorname",
				"m" : "testvalue alarm component",
				"val" : 3
			},
			"status" : {
				"_t" : "statusinc",
				"t" : 1594203725,
				"a" : "testvalueconnector.testvalueconnectorname",
				"m" : "testvalue alarm component",
				"val" : 1
			},
			"steps" : [ 
				{
					"_t" : "stateinc",
					"t" : 1594203725,
					"a" : "testvalueconnector.testvalueconnectorname",
					"m" : "testvalue alarm component",
					"val" : 3
				}, 
				{
					"_t" : "statusinc",
					"t" : 1594203725,
					"a" : "testvalueconnector.testvalueconnectorname",
					"m" : "testvalue alarm component",
					"val" : 1
				}
			],
			"component" : "testvaluecomponent922-110",
			"connector" : "testvalueconnector",
			"connector_name" : "testvalueconnectorname",
			"creation_date" : 1594203725,
			"display_name" : "PW-DK-XZ",
			"initial_output" : "testvalue alarm component",
			"output" : "testvalue alarm component",
			"last_update_date" : 1594203725,
			"last_event_date" : 1594203725,
			"total_state_changes" : 1,
			"extra" : {},
			"infos" : {},
			"infos_rule_version" : {}
		}
	}`)

	err = json.Unmarshal(alarmJSON, &child.Alarm)
	if err != nil {
		t.Fatalf("Child Alarm unmarshal error %s", err)
	}

	children = append(children, child)
	rule := metaalarm.Rule{}
	event := types.Event{}

	var updatedAlarms *[]types.Alarm

	alarmAdapterMock.
		EXPECT().
		PartialMassUpdateOpen(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, alarms []types.Alarm) error {
			updatedAlarms = &alarms
			return nil
		})
	alarmAdapterMock.EXPECT().GetCountOpenedAlarmsByIDs(gomock.Any()).Return(int64(1), nil)

	updateEvent, err := s.AddMultipleChildsToMetaAlarm(context.Background(), event, metaAlarm, children, rule)
	if err != nil {
		t.Fatalf("AddChildToMetaAlarm error %s", err)
	}

	if state := (*updatedAlarms)[0].CurrentState(); state != expectedState {
		t.Errorf("wrong state value %d", state)
	}

	if updateEvent.EventType != types.EventTypeMetaAlarmUpdated {
		t.Fatalf("AddChildToMetaAlarm event %v", updateEvent)
	}
}

func TestChildInheritMetaAlarmActions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	alarmAdapterMock := mock_alarm.NewMockAdapter(ctrl)
	alarmConfigProviderMock := mock_config.NewMockAlarmConfigProvider(ctrl)
	alarmConfigProviderMock.EXPECT().Get().AnyTimes().Return(config.AlarmConfig{})

	s := service.NewMetaAlarmService(alarmAdapterMock, alarmConfigProviderMock, log.NewTestLogger())

	var dataSets = []struct {
		testName      string
		metaalarm     []byte
		child         []byte
		inheritAck    bool
		inheritTicket bool
		inheritSnooze bool
	}{
		{
			testName:  "test inherit all",
			metaalarm: getMetaalarm(),
			child: []byte(`{
				"_id" : "1e6e41ce-3823-4a44-9ee5-c31b7a9fd812",
				"t" : 1594204048,
				"d" : "testvalueressource311/testvaluecomponent922-A311",
				"v" : {
					"state" : {
						"_t" : "stateinc",
						"t" : 1594204048,
						"a" : "testvalueconnector.testvalueconnectorname",
						"m" : "testvalue child",
						"val" : 2
					},
					"status" : {
						"_t" : "statusinc",
						"t" : 1594204048,
						"a" : "testvalueconnector.testvalueconnectorname",
						"m" : "testvalue child",
						"val" : 1
					},
					"steps" : [ 
						{
							"_t" : "stateinc",
							"t" : 1594204048,
							"a" : "testvalueconnector.testvalueconnectorname",
							"m" : "testvalue child",
							"val" : 2
						}, 
						{
							"_t" : "statusinc",
							"t" : 1594204048,
							"a" : "testvalueconnector.testvalueconnectorname",
							"m" : "testvalue child",
							"val" : 1
						}
					],
					"component" : "testvaluecomponent922-A311",
					"connector" : "testvalueconnector",
					"connector_name" : "testvalueconnectorname",
					"creation_date" : 1594204048,
					"display_name" : "VX-LV-EA",
					"initial_output" : "testvalue child",
					"output" : "testvalue child",
					"initial_long_output" : "",
					"last_update_date" : 1594204048,
					"last_event_date" : 1594204048,
					"resource" : "testvalueressource311",
					"parents" : [ 
						"testvaluecomponent922-A311"
					],
					"children" : [],
					"total_state_changes" : 1,
					"extra" : {},
					"infos" : {},
					"infos_rule_version" : {}
				}
			}`),
			inheritAck:    true,
			inheritSnooze: true,
			inheritTicket: true,
		},
		{
			testName:  "test inherit none",
			metaalarm: getMetaalarm(),
			child: []byte(`{
				"_id" : "1e6e41ce-3823-4a44-9ee5-c31b7a9fd812",
				"t" : 1594204048,
				"d" : "testvalueressource311/testvaluecomponent922-A311",
				"v" : {
					"ack": {
						"_t": "ack",
						"t" : 1594204000,
						"a" : "root",
						"m" : "Ack original alarm",
						"role": "admin",
						"val" : 0
					},
					"ticket": {
						"_t": "assocticket",
						"t" : 1594204000,
						"a" : "root",
						"m" : "Ticket original alarm",
						"role": "admin",
						"val" : "Ticket original alarm"
					},
					"state" : {
						"_t" : "stateinc",
						"t" : 1594204048,
						"a" : "testvalueconnector.testvalueconnectorname",
						"m" : "testvalue child",
						"val" : 2
					},
					"status" : {
						"_t" : "statusinc",
						"t" : 1594204048,
						"a" : "testvalueconnector.testvalueconnectorname",
						"m" : "testvalue child",
						"val" : 1
					},
					"snooze": {
						"_t": "snooze",
						"t" : 1594204000,
						"a" : "root",
						"m" : "Snooze original alarm",
						"role": "admin",
						"val" : 2594204000
					},
					"steps" : [ 
						{
							"_t" : "stateinc",
							"t" : 1594204048,
							"a" : "testvalueconnector.testvalueconnectorname",
							"m" : "testvalue child",
							"val" : 2
						}, 
						{
							"_t" : "statusinc",
							"t" : 1594204048,
							"a" : "testvalueconnector.testvalueconnectorname",
							"m" : "testvalue child",
							"val" : 1
						},
						{
							"_t": "assocticket",
							"t" : 1594204000,
							"a" : "root",
							"m" : "Ticket original alarm",
							"role": "admin",
							"val" : 0
						},
						{
							"_t": "ack",
							"t" : 1594204000,
							"a" : "root",
							"m" : "Ack original alarm",
							"role": "admin",
							"val" : 0
						},
						{
							"_t": "snooze",
							"t" : 1594204000,
							"a" : "root",
							"m" : "Snooze original alarm",
							"role": "admin",
							"val" : 2594204000
						}
					],
					"component" : "testvaluecomponent922-A311",
					"connector" : "testvalueconnector",
					"connector_name" : "testvalueconnectorname",
					"creation_date" : 1594204048,
					"display_name" : "VX-LV-EA",
					"initial_output" : "testvalue child",
					"output" : "testvalue child",
					"initial_long_output" : "",
					"last_update_date" : 1594204048,
					"last_event_date" : 1594204048,
					"resource" : "testvalueressource311",
					"parents" : [ 
						"testvaluecomponent922-A311"
					],
					"children" : [],
					"total_state_changes" : 1,
					"extra" : {},
					"infos" : {},
					"infos_rule_version" : {}
				}
			}`),
			inheritAck:    false,
			inheritSnooze: false,
			inheritTicket: false,
		},
		{
			testName:  "test inherit some steps",
			metaalarm: getMetaalarm(),
			child: []byte(`{
				"_id" : "1e6e41ce-3823-4a44-9ee5-c31b7a9fd812",
				"t" : 1594204048,
				"d" : "testvalueressource311/testvaluecomponent922-A311",
				"v" : {
					"ack": {
						"_t": "ack",
						"t" : 1594204000,
						"a" : "root",
						"m" : "Ack original alarm",
						"role": "admin",
						"val" : 0
					},
					"state" : {
						"_t" : "stateinc",
						"t" : 1594204048,
						"a" : "testvalueconnector.testvalueconnectorname",
						"m" : "testvalue child",
						"val" : 2
					},
					"status" : {
						"_t" : "statusinc",
						"t" : 1594204048,
						"a" : "testvalueconnector.testvalueconnectorname",
						"m" : "testvalue child",
						"val" : 1
					},
					"steps" : [ 
						{
							"_t" : "stateinc",
							"t" : 1594204048,
							"a" : "testvalueconnector.testvalueconnectorname",
							"m" : "testvalue child",
							"val" : 2
						}, 
						{
							"_t" : "statusinc",
							"t" : 1594204048,
							"a" : "testvalueconnector.testvalueconnectorname",
							"m" : "testvalue child",
							"val" : 1
						},
						{
							"_t": "ack",
							"t" : 1594204000,
							"a" : "root",
							"m" : "Ack original alarm",
							"role": "admin",
							"val" : 0
						}
					],
					"component" : "testvaluecomponent922-A311",
					"connector" : "testvalueconnector",
					"connector_name" : "testvalueconnectorname",
					"creation_date" : 1594204048,
					"display_name" : "VX-LV-EA",
					"initial_output" : "testvalue child",
					"output" : "testvalue child",
					"initial_long_output" : "",
					"last_update_date" : 1594204048,
					"last_event_date" : 1594204048,
					"resource" : "testvalueressource311",
					"parents" : [ 
						"testvaluecomponent922-A311"
					],
					"children" : [],
					"total_state_changes" : 1,
					"extra" : {},
					"infos" : {},
					"infos_rule_version" : {}
				}
			}`),
			inheritAck:    false,
			inheritSnooze: true,
			inheritTicket: true,
		},
	}

	for _, dataset := range dataSets {
		t.Run(dataset.testName, func(t *testing.T) {
			var metaAlarm types.Alarm
			err := json.Unmarshal(dataset.metaalarm, &metaAlarm)
			if err != nil {
				t.Fatalf("MetaAlarm unmarshal error %s", err)
			}

			var child types.Alarm
			err = json.Unmarshal(dataset.child, &child)
			if err != nil {
				t.Fatalf("Child Alarm unmarshal error %s", err)
			}

			children := []types.AlarmWithEntity{
				{Alarm: child},
			}

			var updatedAlarms *[]types.Alarm

			alarmAdapterMock.
				EXPECT().
				PartialMassUpdateOpen(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, alarms []types.Alarm) error {
					updatedAlarms = &alarms
					return nil

				})
			alarmAdapterMock.EXPECT().GetCountOpenedAlarmsByIDs(gomock.Any()).Return(int64(1), nil)

			event := types.Event{}
			rule := metaalarm.Rule{}
			updateEvent, err := s.AddMultipleChildsToMetaAlarm(context.Background(), event, metaAlarm, children, rule)
			if err != nil {
				t.Fatalf("AddMultipleChildsToMetaAlarm error %s", err)
			}

			if updateEvent.EventType != types.EventTypeMetaAlarmUpdated {
				t.Fatalf("AddMultipleChildsToMetaAlarm event %v", updateEvent)
			}

			for _, alarm := range *updatedAlarms {
				if alarm.EntityID != metaAlarm.EntityID {
					alarmAck := alarm.Value.ACK
					if alarmAck == nil {
						t.Fatalf("invalid child's ACK %s", alarm.EntityID)
					}

					if dataset.inheritAck {
						if alarmAck.Message != metaAlarm.Value.ACK.Message {
							t.Errorf("child's ACK message mismatch %s != %s", alarmAck.Message, metaAlarm.Value.ACK.Message)
						}
					} else {
						if alarmAck.Message == metaAlarm.Value.ACK.Message {
							t.Errorf("child's ACK message shouldn't be inherited")
						}
					}

					alarmTicket := alarm.Value.Ticket
					if alarmTicket == nil {
						t.Fatalf("invalid child's ACK %s", alarm.EntityID)
					}

					if dataset.inheritTicket {
						if alarmTicket.Message != metaAlarm.Value.Ticket.Message {
							t.Errorf("child's assocticket message mismatch %s != %s", alarmTicket.Message, metaAlarm.Value.Ticket.Message)
						}
					} else {
						if alarmTicket.Message == metaAlarm.Value.Ticket.Message {
							t.Errorf("child's assocticket message shouldn't be inherited")
						}
					}

					alarmSnooze := alarm.Value.Snooze
					if alarmSnooze == nil {
						t.Fatalf("invalid child's snooze %s", alarm.EntityID)
					}
					if dataset.inheritSnooze {
						if alarmSnooze.Message != metaAlarm.Value.Snooze.Message {
							t.Errorf("child's snooze message mismatch %s != %s", alarmSnooze.Message, metaAlarm.Value.Snooze.Message)
						}
					} else {
						if alarmSnooze.Message == metaAlarm.Value.Snooze.Message {
							t.Errorf("child's snooze message shouldn't be inherited")
						}
					}

					for _, step := range alarm.Value.Steps {
						switch step.Type {
						case alarmAck.Type:
							if dataset.inheritAck {
								if step.Message != metaAlarm.Value.ACK.Message {
									t.Errorf("step's ACK message mismatch %s != %s", step.Message, metaAlarm.Value.ACK.Message)
								}
								if step.Author != "correlation" {
									t.Errorf("invalid author value %s", step.Author)
								}
							} else {
								if step.Message == metaAlarm.Value.ACK.Message {
									t.Errorf("step's ACK message shouldn't be inherited")
								}
								if step.Author == "correlation" {
									t.Errorf("ack step's author shouldn't be %s", step.Author)
								}
							}
						case alarmTicket.Type:
							if dataset.inheritTicket {
								if step.Message != metaAlarm.Value.Ticket.Message {
									t.Errorf("step's assocticket message mismatch %s != %s", step.Message, metaAlarm.Value.Ticket.Message)
								}
								if step.Author != "correlation" {
									t.Errorf("invalid author value %s", step.Author)
								}
							} else {
								if step.Message == metaAlarm.Value.Ticket.Message {
									t.Errorf("step's assocticket message shouldn't be inherited")
								}
								if step.Author == "correlation" {
									t.Errorf("assocticket step's author shouldn't be %s", step.Author)
								}
							}
						case alarmSnooze.Type:
							if dataset.inheritSnooze {
								if step.Message != metaAlarm.Value.Snooze.Message {
									t.Errorf("step's snooze message mismatch %s != %s", step.Message, metaAlarm.Value.Snooze.Message)
								}
								if step.Author != "correlation" {
									t.Errorf("invalid author value %s", step.Author)
								}
							} else {
								if step.Message == metaAlarm.Value.Snooze.Message {
									t.Errorf("step's snooze message shouldn't be inherited")
								}
								if step.Author == "correlation" {
									t.Errorf("snooze step's author shouldn't be %s", step.Author)
								}
							}
						}
					}
				}
			}
		})
	}
}

func getMetaalarm() []byte {
	return []byte(`{
		"_id" : "40b8aaef-a975-491a-a143-2dfdabe059c8",
		"t" : 1594203725,
		"d" : "meta-child-entity-24cc5fff-6d0e-4f64-ae91-7bd4a0cdc9b9/metaalarm",
		"v" : {
			"ack": {
				"_t": "ack",
				"t" : 1594204000,
				"a" : "root",
				"m" : "Ack meta-child",
				"role": "admin",
				"val" : 0
			},
			"ticket": {
				"_t": "assocticket",
				"t" : 1594204000,
				"a" : "root",
				"m" : "Ticket meta-child",
				"role": "admin",
				"val" : "Ticket meta-child"
			},
			"state" : {
				"_t" : "stateinc",
				"t" : 1594203725,
				"a" : "engine.correlation",
				"m" : "testvalue child component",
				"val" : 1
			},
			"status" : {
				"_t" : "statusinc",
				"t" : 1594203725,
				"a" : "engine.correlation",
				"m" : "testvalue child component",
				"val" : 1
			},
			"snooze": {
				"_t": "snooze",
				"t" : 1594204000,
				"a" : "root",
				"m" : "Snooze meta-child",
				"role": "admin",
				"val" : 2594204000
			},
			"steps" : [ 
				{
					"_t" : "statusinc",
					"t" : 1594203725,
					"a" : "engine.correlation",
					"m" : "testvalue child component",
					"val" : 2
				},
				{
					"_t": "assocticket",
					"t" : 1594204000,
					"a" : "root",
					"m" : "Ticket meta-child",
					"role": "admin",
					"val" : 0
				},
				{
					"_t": "ack",
					"t" : 1594204000,
					"a" : "root",
					"m" : "Ack meta-child",
					"role": "admin",
					"val" : 0
				},
				{
					"_t": "snooze",
					"t" : 1594204000,
					"a" : "root",
					"m" : "Snooze meta-child",
					"role": "admin",
					"val" : 2594204000
				}
			],
			"component" : "metaalarm",
			"connector" : "engine",
			"connector_name" : "correlation",
			"creation_date" : 1594203725,
			"display_name" : "WP-NZ-KV",
			"initial_output" : "testvalue child component",
			"output" : "testvalue child component",
			"last_update_date" : 1594203725,
			"last_event_date" : 1594203725,
			"resource" : "meta-child-entity-24cc5fff-6d0e-4f64-ae91-7bd4a0cdc9b9",
			"meta" : "testRule",
			"parents" : [],
			"children" : [
				"testvaluecomponent922-A311"
			],
			"total_state_changes" : 0,
			"extra" : {},
			"infos" : {},
			"infos_rule_version" : {}
		}
	}`)
}
