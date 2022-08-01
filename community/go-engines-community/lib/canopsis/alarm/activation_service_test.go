package alarm

import (
	"bytes"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	mock_amqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/amqp"
	mock_encoding "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/encoding"
	"github.com/golang/mock/gomock"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"testing"
	"time"
)

func TestActivationService_Process_GivenInactiveAlarm_ShouldPublishEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	encoderMock := mock_encoding.NewMockEncoder(ctrl)
	publisherMock := mock_amqp.NewMockPublisher(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	queueName := "testQueue"
	service := NewActivationService(
		encoderMock,
		publisherMock,
		queueName,
		logger,
	)
	alarm := &types.Alarm{}

	eventBody := make([]byte, 1)
	encoderMock.
		EXPECT().
		Encode(gomock.Any()).
		Return(eventBody, nil)

	publisherMock.
		EXPECT().
		Publish(
			gomock.Eq(""),
			gomock.Eq(queueName),
			gomock.Eq(false),
			gomock.Eq(false),
			gomock.Eq(amqp.Publishing{
				ContentType:  "application/json",
				Body:         eventBody,
				DeliveryMode: amqp.Persistent,
			}),
		).
		Times(1)

	_, err := service.Process(alarm)
	if err != nil {
		t.Errorf("exepected not error but got %v", err)
	}
}

func TestActivationService_Process_GivenInactiveAlarm_ShouldPublishActiveEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	encoderMock := mock_encoding.NewMockEncoder(ctrl)
	publisherMock := mock_amqp.NewMockPublisher(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	service := NewActivationService(
		encoderMock,
		publisherMock,
		"testQueue",
		logger,
	)
	alarm := &types.Alarm{
		Value: types.AlarmValue{
			Component:     "testcomp",
			Connector:     "testconn",
			ConnectorName: "testconnname",
			Resource:      "testres",
		},
	}

	eventBody := make([]byte, 1)
	encoderMock.
		EXPECT().
		Encode(gomock.Any()).
		Do(func(event types.Event) {
			if event.EventType != types.EventTypeActivate {
				t.Errorf("expected event type: %v but go %v", types.EventTypeActivate, event.EventType)
			}
			if event.SourceType != types.SourceTypeResource {
				t.Errorf("expected event source type: %v but go %v", types.SourceTypeResource, event.SourceType)
			}
			if event.Connector != alarm.Value.Connector {
				t.Errorf("expected event connector: %v but go %v", alarm.Value.Connector, event.Connector)
			}
			if event.ConnectorName != alarm.Value.ConnectorName {
				t.Errorf("expected event connector name: %v but go %v", alarm.Value.ConnectorName, event.ConnectorName)
			}
			if event.Component != alarm.Value.Component {
				t.Errorf("expected event component: %v but go %v", alarm.Value.Component, event.Component)
			}
			if event.Resource != alarm.Value.Resource {
				t.Errorf("expected event resource: %v but go %v", alarm.Value.Resource, event.Resource)
			}
		}).
		Return(eventBody, nil)

	publisherMock.
		EXPECT().
		Publish(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		)

	_, err := service.Process(alarm)
	if err != nil {
		t.Errorf("exepected not error but got %v", err)
	}
}

func TestActivationService_Process_GivenInactiveAndSnoozedAlarm_ShouldNotPublishEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	encoderMock := mock_encoding.NewMockEncoder(ctrl)
	publisherMock := mock_amqp.NewMockPublisher(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	service := NewActivationService(
		encoderMock,
		publisherMock,
		"testQueue",
		logger,
	)
	alarm := &types.Alarm{
		Value: types.AlarmValue{
			Snooze: &types.AlarmStep{
				Value: types.CpsNumber(time.Now().Unix() + 1000),
			},
		},
	}

	encoderMock.
		EXPECT().
		Encode(gomock.Any()).
		Times(0)

	publisherMock.EXPECT().
		Publish(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).
		Times(0)

	_, err := service.Process(alarm)
	if err != nil {
		t.Errorf("exepected not error but got %v", err)
	}
}

func TestActivationService_Process_GivenInactiveAlarmWithActivePBehavior_ShouldNotPublishEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	encoderMock := mock_encoding.NewMockEncoder(ctrl)
	publisherMock := mock_amqp.NewMockPublisher(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	service := NewActivationService(
		encoderMock,
		publisherMock,
		"testQueue",
		logger,
	)
	alarm := &types.Alarm{
		EntityID: "testID",
		Value: types.AlarmValue{
			PbehaviorInfo: types.PbehaviorInfo{CanonicalType: pbehavior.TypeInactive},
		},
	}

	encoderMock.
		EXPECT().
		Encode(gomock.Any()).
		Times(0)

	publisherMock.EXPECT().
		Publish(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).
		Times(0)

	_, err := service.Process(alarm)
	if err != nil {
		t.Errorf("exepected not error but got %v", err)
	}
}

func TestActivationService_Process_GivenActiveAlarm_ShouldNotPublishEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	encoderMock := mock_encoding.NewMockEncoder(ctrl)
	publisherMock := mock_amqp.NewMockPublisher(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	service := NewActivationService(
		encoderMock,
		publisherMock,
		"testQueue",
		logger,
	)
	alarm := &types.Alarm{
		Value: types.AlarmValue{
			ActivationDate: &types.CpsTime{Time: time.Now()},
		},
	}

	encoderMock.
		EXPECT().
		Encode(gomock.Any()).
		Times(0)

	publisherMock.EXPECT().
		Publish(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).
		Times(0)

	_, err := service.Process(alarm)
	if err != nil {
		t.Errorf("exepected not error but got %v", err)
	}
}
