package idlealarm

import (
	"bytes"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/idlerule"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	mock_alarm "git.canopsis.net/canopsis/go-engines/mocks/lib/canopsis/alarm"
	mock_encoding "git.canopsis.net/canopsis/go-engines/mocks/lib/canopsis/encoding"
	mock_idlerule "git.canopsis.net/canopsis/go-engines/mocks/lib/canopsis/idlerule"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"testing"
	"time"
)

func TestService_Process_GivenLastEventDateRuleAndMatchedAlarm_ShouldReturnEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	ruleAdapterMock := mock_idlerule.NewMockRuleAdapter(ctrl)
	alarmAdapterMock := mock_alarm.NewMockAdapter(ctrl)
	encoderMock := mock_encoding.NewMockEncoder(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	alarm := types.AlarmWithEntity{
		Alarm: types.Alarm{
			Value: types.AlarmValue{
				LastEventDate: types.CpsTime{Time: time.Now().Add(-6 * time.Hour)},
				Connector:     "test-connector",
				ConnectorName: "test-connector-name",
				Component:     "test-component",
				Resource:      "test-resource",
			},
		},
		Entity: types.Entity{},
	}
	author := "test-author"
	output := "test-output"
	rule := idlerule.Rule{
		Type:     idlerule.RuleTypeLastEvent,
		Duration: 0,
		Operation: idlerule.Operation{
			Type: types.ActionTypeAck,
			Parameters: types.OperationParameters{
				Output: output,
				Author: author,
			},
		},
	}

	ruleAdapterMock.
		EXPECT().
		Get().
		Return([]idlerule.Rule{rule}, nil)
	alarmAdapterMock.
		EXPECT().
		GetOpenedAlarmsWithLastEventDateBefore(gomock.Any()).
		Return([]types.AlarmWithEntity{alarm}, nil)

	service := NewService(ruleAdapterMock, alarmAdapterMock, encoderMock, logger)
	events := service.Process()

	if len(events) == 0 {
		t.Errorf("exepected event but got nothing")
		return
	}

	if events[0].EventType != types.EventTypeAck {
		t.Errorf("exepected event type %v but %v", types.EventTypeAck, events[0].EventType)
	}

	if events[0].Connector != alarm.Alarm.Value.Connector {
		t.Errorf("exepected connector %v but %v", alarm.Alarm.Value.Connector, events[0].Connector)
	}

	if events[0].ConnectorName != alarm.Alarm.Value.ConnectorName {
		t.Errorf("exepected connector name %v but %v", alarm.Alarm.Value.ConnectorName, events[0].ConnectorName)
	}

	if events[0].Component != alarm.Alarm.Value.Component {
		t.Errorf("exepected component %v but %v", alarm.Alarm.Value.Component, events[0].Component)
	}

	if events[0].Resource != alarm.Alarm.Value.Resource {
		t.Errorf("exepected resource %v but %v", alarm.Alarm.Value.Resource, events[0].Resource)
	}

	if events[0].Author != author {
		t.Errorf("exepected author %v but %v", author, events[0].Author)
	}

	if events[0].Output != output {
		t.Errorf("exepected output %v but %v", output, events[0].Output)
	}
}

func TestService_Process_GivenLastUpdateDateRuleAndMatchedAlarm_ShouldReturnEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	ruleAdapterMock := mock_idlerule.NewMockRuleAdapter(ctrl)
	alarmAdapterMock := mock_alarm.NewMockAdapter(ctrl)
	encoderMock := mock_encoding.NewMockEncoder(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	alarm := types.AlarmWithEntity{
		Alarm: types.Alarm{
			Value: types.AlarmValue{
				LastUpdateDate: types.CpsTime{Time: time.Now().Add(-6 * time.Hour)},
				Connector:      "test-connector",
				ConnectorName:  "test-connector-name",
				Component:      "test-component",
				Resource:       "test-resource",
			},
		},
		Entity: types.Entity{},
	}
	author := "test-author"
	output := "test-output"
	rule := idlerule.Rule{
		Type:     idlerule.RuleTypeLastUpdate,
		Duration: 0,
		Operation: idlerule.Operation{
			Type: types.ActionTypeAck,
			Parameters: types.OperationParameters{
				Output: output,
				Author: author,
			},
		},
	}

	ruleAdapterMock.
		EXPECT().
		Get().
		Return([]idlerule.Rule{rule}, nil)
	alarmAdapterMock.
		EXPECT().
		GetOpenedAlarmsWithLastUpdateDateBefore(gomock.Any()).
		Return([]types.AlarmWithEntity{alarm}, nil)

	service := NewService(ruleAdapterMock, alarmAdapterMock, encoderMock, logger)
	events := service.Process()

	if len(events) == 0 {
		t.Errorf("exepected event but got nothing")
		return
	}

	if events[0].EventType != types.EventTypeAck {
		t.Errorf("exepected event type %v but %v", types.EventTypeAck, events[0].EventType)
	}

	if events[0].Connector != alarm.Alarm.Value.Connector {
		t.Errorf("exepected connector %v but %v", alarm.Alarm.Value.Connector, events[0].Connector)
	}

	if events[0].ConnectorName != alarm.Alarm.Value.ConnectorName {
		t.Errorf("exepected connector name %v but %v", alarm.Alarm.Value.ConnectorName, events[0].ConnectorName)
	}

	if events[0].Component != alarm.Alarm.Value.Component {
		t.Errorf("exepected component %v but %v", alarm.Alarm.Value.Component, events[0].Component)
	}

	if events[0].Resource != alarm.Alarm.Value.Resource {
		t.Errorf("exepected resource %v but %v", alarm.Alarm.Value.Resource, events[0].Resource)
	}

	if events[0].Author != author {
		t.Errorf("exepected author %v but %v", author, events[0].Author)
	}

	if events[0].Output != output {
		t.Errorf("exepected output %v but %v", output, events[0].Output)
	}
}
