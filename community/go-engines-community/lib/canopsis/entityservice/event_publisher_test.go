package entityservice_test

import (
	"context"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	mock_amqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/amqp"
	mock_alarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/alarm"
	mock_encoding "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/encoding"
	mock_entity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/entity"
	"github.com/golang/mock/gomock"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

func TestEventPublisher_Publish_GivenChangedEntity_ShouldSendEvent(t *testing.T) {
	dataSets := map[string]struct {
		Alarm             types.Alarm
		Message           entityservice.ChangeEntityMessage
		ExpectedEventType string
	}{
		"given updated entity should send entityupdated event": {
			Alarm: types.Alarm{
				Value: types.AlarmValue{
					Connector:     "test-connector",
					ConnectorName: "test-connector-name",
					Component:     "test-component",
					Resource:      "test-resource",
				},
			},
			Message: entityservice.ChangeEntityMessage{
				ID: "test-entity",
			},
			ExpectedEventType: types.EventTypeEntityUpdated,
		},
		"given toggled entity should send entitytoggled event": {
			Alarm: types.Alarm{
				Value: types.AlarmValue{
					Connector:     "test-connector",
					ConnectorName: "test-connector-name",
					Component:     "test-component",
					Resource:      "test-resource",
				},
			},
			Message: entityservice.ChangeEntityMessage{
				ID:        "test-entity",
				IsToggled: true,
			},
			ExpectedEventType: types.EventTypeEntityToggled,
		},
	}

	for testCase, data := range dataSets {
		t.Run(testCase, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			contentType := "application/json"
			exchange := ""
			routingKey := canopsis.FIFOQueueName
			logger := zerolog.Logger{}
			body := []byte("test-body")
			mockAlarmAdapter := mock_alarm.NewMockAdapter(ctrl)
			mockAlarmAdapter.EXPECT().GetAlarmsByID(gomock.Any(), gomock.Eq(data.Message.ID)).
				Return([]types.Alarm{data.Alarm}, nil)
			mockEntityAdapter := mock_entity.NewMockAdapter(ctrl)
			mockEncoder := mock_encoding.NewMockEncoder(ctrl)
			mockEncoder.EXPECT().Encode(gomock.Any()).Do(func(event types.Event) {
				if event.EventType != data.ExpectedEventType {
					t.Errorf("expected event type %+v but got %+v", data.ExpectedEventType, event.EventType)
				}
				if event.Connector != data.Alarm.Value.Connector {
					t.Errorf("expected connector %+v but got %+v", data.Alarm.Value.Connector, event.Connector)
				}
				if event.ConnectorName != data.Alarm.Value.ConnectorName {
					t.Errorf("expected connector name %+v but got %+v", data.Alarm.Value.ConnectorName, event.ConnectorName)
				}
				if event.Component != data.Alarm.Value.Component {
					t.Errorf("expected commponent %+v but got %+v", data.Alarm.Value.Component, event.Component)
				}
				if event.Resource != data.Alarm.Value.Resource {
					t.Errorf("expected resource %+v but got %+v", data.Alarm.Value.Resource, event.Resource)
				}
			}).Return(body, nil)
			mockPublisher := mock_amqp.NewMockPublisher(ctrl)
			mockPublisher.EXPECT().PublishWithContext(gomock.Any(), gomock.Eq(exchange), gomock.Eq(routingKey),
				gomock.Any(), gomock.Any(), gomock.Eq(amqp.Publishing{
					ContentType:  contentType,
					Body:         body,
					DeliveryMode: amqp.Persistent,
				})).
				Return(nil)

			eventPublisher := entityservice.NewEventPublisher(mockAlarmAdapter, mockEntityAdapter, mockPublisher, mockEncoder,
				contentType, exchange, routingKey, logger)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			ch := make(chan entityservice.ChangeEntityMessage, 1)
			ch <- data.Message
			defer close(ch)

			go func() {
				<-time.After(10 * time.Millisecond)
				cancel()
			}()

			eventPublisher.Publish(ctx, ch)
		})
	}
}

func TestEventPublisher_Publish_GivenChangedService_ShouldSendEvent(t *testing.T) {
	dataSets := map[string]struct {
		Message           entityservice.ChangeEntityMessage
		ExpectedEventType string
	}{
		"given updated service should send entityupdated event": {
			Message: entityservice.ChangeEntityMessage{
				ID:         "test-service",
				EntityType: types.EntityTypeService,
			},
			ExpectedEventType: types.EventTypeEntityUpdated,
		},
		"given toggled service should send recomputeentityservice event": {
			Message: entityservice.ChangeEntityMessage{
				ID:         "test-service",
				EntityType: types.EntityTypeService,
				IsToggled:  true,
			},
			ExpectedEventType: types.EventTypeRecomputeEntityService,
		},
		"given updated service pattern should send recomputeentityservice event": {
			Message: entityservice.ChangeEntityMessage{
				ID:                      "test-service",
				EntityType:              types.EntityTypeService,
				IsServicePatternChanged: true,
			},
			ExpectedEventType: types.EventTypeRecomputeEntityService,
		},
	}

	for testCase, data := range dataSets {
		t.Run(testCase, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			contentType := "application/json"
			exchange := ""
			routingKey := canopsis.FIFOQueueName
			logger := zerolog.Logger{}
			body := []byte("test-body")
			mockAlarmAdapter := mock_alarm.NewMockAdapter(ctrl)
			mockEntityAdapter := mock_entity.NewMockAdapter(ctrl)
			mockEncoder := mock_encoding.NewMockEncoder(ctrl)
			mockEncoder.EXPECT().Encode(gomock.Any()).Do(func(event types.Event) {
				if event.EventType != data.ExpectedEventType {
					t.Errorf("expected event type %+v but got %+v", data.ExpectedEventType, event.EventType)
				}
				if event.Component != data.Message.ID {
					t.Errorf("expected commponent %+v but got %+v", data.Message.ID, event.Component)
				}
				if event.Resource != "" {
					t.Errorf("expected resource %+v but got %+v", "", event.Resource)
				}
			}).Return(body, nil)
			mockPublisher := mock_amqp.NewMockPublisher(ctrl)
			mockPublisher.EXPECT().PublishWithContext(gomock.Any(), gomock.Eq(exchange), gomock.Eq(routingKey),
				gomock.Any(), gomock.Any(), gomock.Eq(amqp.Publishing{
					ContentType:  contentType,
					Body:         body,
					DeliveryMode: amqp.Persistent,
				})).
				Return(nil)

			eventPublisher := entityservice.NewEventPublisher(mockAlarmAdapter, mockEntityAdapter, mockPublisher, mockEncoder,
				contentType, exchange, routingKey, logger)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			ch := make(chan entityservice.ChangeEntityMessage, 1)
			ch <- data.Message
			defer close(ch)

			go func() {
				<-time.After(10 * time.Millisecond)
				cancel()
			}()

			eventPublisher.Publish(ctx, ch)
		})
	}
}
