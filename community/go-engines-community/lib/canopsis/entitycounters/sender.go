package entitycounters

//go:generate mockgen -destination=../../../mocks/lib/canopsis/entitycounters/sender.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters EventsSender

import (
	"context"
	"fmt"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libamqp "github.com/rabbitmq/amqp091-go"
)

type EventsSender interface {
	UpdateComponentState(ctx context.Context, id, connectorID string, state int) error
	UpdateServiceState(ctx context.Context, serviceID string, serviceInfo UpdatedServicesInfo) error
	RecomputeService(ctx context.Context, serviceID string) error
	RecomputeComponent(ctx context.Context, componentID, connectorID string) error
}

type sender struct {
	encoder         encoding.Encoder
	pubChannel      amqp.Publisher
	pubExchangeName string
	pubQueueName    string
}

func NewEventSender(
	encoder encoding.Encoder,
	pubChannel amqp.Publisher,
	pubExchangeName string,
	pubQueueName string,
) EventsSender {
	return &sender{
		encoder:         encoder,
		pubChannel:      pubChannel,
		pubExchangeName: pubExchangeName,
		pubQueueName:    pubQueueName,
	}
}

func (s *sender) UpdateComponentState(ctx context.Context, id, connectorID string, state int) error {
	connector, connectorName, _ := strings.Cut(connectorID, "/")

	event := types.Event{
		EventType:     types.EventTypeCheck,
		SourceType:    types.SourceTypeComponent,
		Component:     id,
		Connector:     connector,
		ConnectorName: connectorName,
		State:         types.CpsNumber(state),
		Output:        "",
		Timestamp:     datetime.NewCpsTime(),
		Initiator:     types.InitiatorSystem,
	}

	body, err := s.encoder.Encode(event)
	if err != nil {
		return fmt.Errorf("unable to serialize service event: %w", err)
	}

	err = s.pubChannel.PublishWithContext(
		ctx,
		s.pubExchangeName,
		s.pubQueueName,
		false,
		false,
		libamqp.Publishing{
			Body:        body,
			ContentType: "application/json",
		},
	)
	if err != nil {
		return fmt.Errorf("unable to send service event: %w", err)
	}

	return nil
}

func (s *sender) RecomputeService(ctx context.Context, serviceID string) error {
	event := types.Event{
		EventType:     types.EventTypeRecomputeEntityService,
		SourceType:    types.SourceTypeService,
		Component:     serviceID,
		Connector:     types.ConnectorEngineService,
		ConnectorName: types.ConnectorEngineService,
		Timestamp:     datetime.NewCpsTime(),
		Initiator:     types.InitiatorSystem,
	}

	body, err := s.encoder.Encode(event)
	if err != nil {
		return fmt.Errorf("unable to serialize service event: %w", err)
	}

	err = s.pubChannel.PublishWithContext(
		ctx,
		s.pubExchangeName,
		s.pubQueueName,
		false,
		false,
		libamqp.Publishing{
			Body:        body,
			ContentType: "application/json",
		},
	)
	if err != nil {
		return fmt.Errorf("unable to send service event: %w", err)
	}

	return nil
}

func (s *sender) UpdateServiceState(ctx context.Context, serviceID string, serviceInfo UpdatedServicesInfo) error {
	event := types.Event{
		EventType:     types.EventTypeCheck,
		SourceType:    types.SourceTypeService,
		Component:     serviceID,
		Connector:     types.ConnectorEngineService,
		ConnectorName: types.ConnectorEngineService,
		State:         types.CpsNumber(serviceInfo.State),
		Output:        serviceInfo.Output,
		Timestamp:     datetime.NewCpsTime(),
		Initiator:     types.InitiatorSystem,
	}

	body, err := s.encoder.Encode(event)
	if err != nil {
		return fmt.Errorf("unable to serialize service event: %w", err)
	}

	err = s.pubChannel.PublishWithContext(
		ctx,
		s.pubExchangeName,
		s.pubQueueName,
		false,
		false,
		libamqp.Publishing{
			Body:        body,
			ContentType: "application/json",
		},
	)
	if err != nil {
		return fmt.Errorf("unable to send service event: %w", err)
	}

	return nil
}

func (s *sender) RecomputeComponent(ctx context.Context, componentID, connectorID string) error {
	connector, connectorName, _ := strings.Cut(connectorID, "/")

	event := types.Event{
		EventType:     types.EventTypeEntityUpdated,
		SourceType:    types.SourceTypeComponent,
		Component:     componentID,
		Connector:     connector,
		ConnectorName: connectorName,
		Timestamp:     datetime.NewCpsTime(),
		Initiator:     types.InitiatorSystem,
	}

	body, err := s.encoder.Encode(event)
	if err != nil {
		return fmt.Errorf("unable to serialize service event: %w", err)
	}

	err = s.pubChannel.PublishWithContext(
		ctx,
		s.pubExchangeName,
		s.pubQueueName,
		false,
		false,
		libamqp.Publishing{
			Body:        body,
			ContentType: "application/json",
		},
	)
	if err != nil {
		return fmt.Errorf("unable to send service event: %w", err)
	}

	return nil
}
