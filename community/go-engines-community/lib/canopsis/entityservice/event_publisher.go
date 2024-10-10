package entityservice

import (
	"context"
	"strings"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

// EventPublisher is used to send event to engines' event flow to notify about entity changes.
type EventPublisher interface {
	Publish(ctx context.Context, ch <-chan ChangeEntityMessage)
}

type eventPublisher struct {
	encoder       encoding.Encoder
	amqpPublisher libamqp.Publisher
	// contentType is added to amqp message to support encoded format.
	contentType string
	// exchange, routingKey is amqp queue to publish message.
	exchange, routingKey string
	connector            string
	logger               zerolog.Logger
}

type ChangeEntityMessage struct {
	ID        string
	Name      string
	Component string
	// IsToggled is true if entity is disabled or enabled.
	IsToggled  bool
	IsDeleted  bool
	EntityType string
	// IsPatternChanged defines should service's context graph and state be recomputed.
	IsServicePatternChanged bool
	// Resources are used only when component entity is toggled to toggle dependent resources
	Resources []string
}

func NewEventPublisher(
	publisher libamqp.Publisher,
	encoder encoding.Encoder,
	contentType string,
	exchange, routingKey string,
	connector string,
	logger zerolog.Logger,
) EventPublisher {
	return &eventPublisher{
		amqpPublisher: publisher,
		encoder:       encoder,
		contentType:   contentType,
		exchange:      exchange,
		routingKey:    routingKey,
		connector:     connector,
		logger:        logger,
	}
}

func (p *eventPublisher) Publish(
	ctx context.Context,
	ch <-chan ChangeEntityMessage,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}

			if msg.ID == "" {
				p.logger.Error().Msg("message must contain entity id")
				continue
			}

			if msg.EntityType == types.EntityTypeService {
				go p.publishServiceEvent(ctx, msg)
			} else {
				go p.publishBasicEntityEvent(ctx, msg)
			}
		}
	}
}

func (p *eventPublisher) publishServiceEvent(ctx context.Context, msg ChangeEntityMessage) {
	var eventType string
	if msg.IsServicePatternChanged || msg.IsToggled || msg.IsDeleted {
		eventType = types.EventTypeRecomputeEntityService
	} else {
		eventType = types.EventTypeEntityUpdated
	}

	event := types.Event{
		EventType:     eventType,
		Connector:     p.connector,
		ConnectorName: p.connector,
		Component:     msg.ID,
		Timestamp:     datetime.NewCpsTime(),
		SourceType:    types.SourceTypeService,
		Author:        canopsis.DefaultEventAuthor,
		Initiator:     types.InitiatorSystem,
	}

	err := p.publishEvent(ctx, event)
	if err != nil {
		p.logger.Err(err).Msg("cannot send event to amqp")
		return
	}

	p.logger.Debug().Msgf("publish %s", msg.ID)
}

func (p *eventPublisher) publishBasicEntityEvent(ctx context.Context, msg ChangeEntityMessage) {
	var event types.Event
	switch msg.EntityType {
	case types.EntityTypeConnector:
		event = types.Event{
			Connector:     strings.ReplaceAll(msg.ID, "/"+msg.Name, ""),
			ConnectorName: msg.Name,
			SourceType:    types.SourceTypeConnector,
		}
	case types.EntityTypeComponent:
		event = types.Event{
			Connector:     p.connector,
			ConnectorName: p.connector,
			Component:     msg.ID,
			SourceType:    types.SourceTypeComponent,
		}
	case types.EntityTypeResource:
		event = types.Event{
			Connector:     p.connector,
			ConnectorName: p.connector,
			Component:     msg.Component,
			Resource:      msg.Name,
			SourceType:    types.SourceTypeResource,
		}
	}

	if msg.IsToggled {
		event.EventType = types.EventTypeEntityToggled
	} else if msg.IsDeleted {
		event.EventType = types.EventTypeResolveDeleted
	} else {
		event.EventType = types.EventTypeEntityUpdated
	}

	event.Timestamp = datetime.NewCpsTime()
	event.Author = canopsis.DefaultEventAuthor
	event.Initiator = types.InitiatorSystem
	err := p.publishEvent(ctx, event)
	if err != nil {
		p.logger.Err(err).Str("entity_id", msg.ID).Msg("cannot send event to amqp")
		return
	}

	p.logger.Debug().Msgf("publish %s", msg.ID)

	if msg.IsToggled && msg.EntityType == types.EntityTypeComponent {
		resourceEvent := types.Event{
			EventType:     types.EventTypeEntityToggled,
			Connector:     event.Connector,
			ConnectorName: event.ConnectorName,
			Component:     event.Component,
			SourceType:    types.SourceTypeResource,
			Timestamp:     event.Timestamp,
			Author:        event.Author,
			Initiator:     event.Initiator,
		}

		for _, resID := range msg.Resources {
			resourceEvent.Resource = strings.ReplaceAll(resID, "/"+msg.ID, "")

			err = p.publishEvent(ctx, resourceEvent)
			if err != nil {
				p.logger.Err(err).Str("entity_id", resID).Msg("cannot send event to amqp")
				return
			}

			p.logger.Debug().Msgf("publish %s", resID)
		}
	}
}

func (p *eventPublisher) publishEvent(ctx context.Context, event types.Event) error {
	body, err := p.encoder.Encode(event)
	if err != nil {
		return err
	}

	return p.amqpPublisher.PublishWithContext(
		ctx,
		p.exchange,
		p.routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  p.contentType,
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}
