package entityservice_test

import (
	"context"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	mock_amqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/amqp"
	mock_encoding "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/encoding"
	"github.com/golang/mock/gomock"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

func TestEventPublisher_Publish_GivenChangedEntity_ShouldSendEvent(t *testing.T) {
	dataSets := map[string]struct {
		Message       entityservice.ChangeEntityMessage
		ExpectedEvent types.Event
	}{
		"given updated connector should send entityupdated event": {
			Message: entityservice.ChangeEntityMessage{
				ID:         "test-connector/test-connector-name",
				EntityType: types.EntityTypeConnector,
				Name:       "test-connector-name",
			},
			ExpectedEvent: types.Event{
				EventType:     types.EventTypeEntityUpdated,
				Connector:     "test-connector",
				ConnectorName: "test-connector-name",
				SourceType:    types.SourceTypeConnector,
			},
		},
		"given updated component should send entityupdated event": {
			Message: entityservice.ChangeEntityMessage{
				ID:         "test-component",
				EntityType: types.EntityTypeComponent,
				Name:       "test-component",
				Component:  "test-component",
			},
			ExpectedEvent: types.Event{
				EventType:     types.EventTypeEntityUpdated,
				Connector:     types.ConnectorApi,
				ConnectorName: types.ConnectorApi,
				Component:     "test-component",
				SourceType:    types.SourceTypeComponent,
			},
		},
		"given updated resource should send entityupdated event": {
			Message: entityservice.ChangeEntityMessage{
				ID:         "test-resource/test-component",
				EntityType: types.EntityTypeResource,
				Name:       "test-resource",
				Component:  "test-component",
			},
			ExpectedEvent: types.Event{
				EventType:     types.EventTypeEntityUpdated,
				Connector:     types.ConnectorApi,
				ConnectorName: types.ConnectorApi,
				Component:     "test-component",
				Resource:      "test-resource",
				SourceType:    types.SourceTypeResource,
			},
		},
		"given toggled entity should send entitytoggled event": {
			Message: entityservice.ChangeEntityMessage{
				ID:         "test-resource/test-component",
				EntityType: types.EntityTypeResource,
				Name:       "test-resource",
				Component:  "test-component",
				IsToggled:  true,
			},
			ExpectedEvent: types.Event{
				EventType:     types.EventTypeEntityToggled,
				Connector:     types.ConnectorApi,
				ConnectorName: types.ConnectorApi,
				Component:     "test-component",
				Resource:      "test-resource",
				SourceType:    types.SourceTypeResource,
			},
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
			mockEncoder := mock_encoding.NewMockEncoder(ctrl)
			mockEncoder.EXPECT().Encode(gomock.Any()).Do(func(event types.Event) {
				if event.EventType != data.ExpectedEvent.EventType {
					t.Errorf("expected event type %+v but got %+v", data.ExpectedEvent.EventType, event.EventType)
				}
				if event.Connector != data.ExpectedEvent.Connector {
					t.Errorf("expected connector %+v but got %+v", data.ExpectedEvent.Connector, event.Connector)
				}
				if event.ConnectorName != data.ExpectedEvent.ConnectorName {
					t.Errorf("expected connector name %+v but got %+v", data.ExpectedEvent.ConnectorName, event.ConnectorName)
				}
				if event.Component != data.ExpectedEvent.Component {
					t.Errorf("expected commponent %+v but got %+v", data.ExpectedEvent.Component, event.Component)
				}
				if event.Resource != data.ExpectedEvent.Resource {
					t.Errorf("expected resource %+v but got %+v", data.ExpectedEvent.Resource, event.Resource)
				}
				if event.SourceType != data.ExpectedEvent.SourceType {
					t.Errorf("expected source type %+v but got %+v", data.ExpectedEvent.SourceType, event.SourceType)
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

			eventPublisher := entityservice.NewEventPublisher(mockPublisher, mockEncoder, contentType, exchange, routingKey, logger)

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
		Message       entityservice.ChangeEntityMessage
		ExpectedEvent types.Event
	}{
		"given updated service should send entityupdated event": {
			Message: entityservice.ChangeEntityMessage{
				ID:         "test-service",
				EntityType: types.EntityTypeService,
			},
			ExpectedEvent: types.Event{
				EventType:     types.EventTypeEntityUpdated,
				Connector:     types.ConnectorApi,
				ConnectorName: types.ConnectorApi,
				Component:     "test-service",
				SourceType:    types.SourceTypeService,
			},
		},
		"given toggled service should send recomputeentityservice event": {
			Message: entityservice.ChangeEntityMessage{
				ID:         "test-service",
				EntityType: types.EntityTypeService,
				IsToggled:  true,
			},
			ExpectedEvent: types.Event{
				EventType:     types.EventTypeRecomputeEntityService,
				Connector:     types.ConnectorApi,
				ConnectorName: types.ConnectorApi,
				Component:     "test-service",
				SourceType:    types.SourceTypeService,
			},
		},
		"given updated service pattern should send recomputeentityservice event": {
			Message: entityservice.ChangeEntityMessage{
				ID:                      "test-service",
				EntityType:              types.EntityTypeService,
				IsServicePatternChanged: true,
			},
			ExpectedEvent: types.Event{
				EventType:     types.EventTypeRecomputeEntityService,
				Connector:     types.ConnectorApi,
				ConnectorName: types.ConnectorApi,
				Component:     "test-service",
				SourceType:    types.SourceTypeService,
			},
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
			mockEncoder := mock_encoding.NewMockEncoder(ctrl)
			mockEncoder.EXPECT().Encode(gomock.Any()).Do(func(event types.Event) {
				if event.EventType != data.ExpectedEvent.EventType {
					t.Errorf("expected event type %+v but got %+v", data.ExpectedEvent.EventType, event.EventType)
				}
				if event.Connector != data.ExpectedEvent.Connector {
					t.Errorf("expected connector %+v but got %+v", data.ExpectedEvent.Connector, event.Connector)
				}
				if event.ConnectorName != data.ExpectedEvent.ConnectorName {
					t.Errorf("expected connector name %+v but got %+v", data.ExpectedEvent.ConnectorName, event.ConnectorName)
				}
				if event.Component != data.ExpectedEvent.Component {
					t.Errorf("expected commponent %+v but got %+v", data.ExpectedEvent.Component, event.Component)
				}
				if event.Resource != data.ExpectedEvent.Resource {
					t.Errorf("expected resource %+v but got %+v", data.ExpectedEvent.Resource, event.Resource)
				}
				if event.SourceType != data.ExpectedEvent.SourceType {
					t.Errorf("expected source type %+v but got %+v", data.ExpectedEvent.SourceType, event.SourceType)
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

			eventPublisher := entityservice.NewEventPublisher(mockPublisher, mockEncoder, contentType, exchange, routingKey, logger)

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
