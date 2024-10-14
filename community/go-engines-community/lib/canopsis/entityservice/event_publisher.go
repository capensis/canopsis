package entityservice

import (
	"context"
	"strings"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

// EventPublisher is used to send event to engines' event flow to notify about entity changes.
type EventPublisher interface {
	Publish(ctx context.Context, ch <-chan ChangeEntityMessage)
}

type eventPublisher struct {
	alarmAdapter  libalarm.Adapter
	entityAdapter entity.Adapter
	encoder       encoding.Encoder
	amqpPublisher libamqp.Publisher
	// contentType is added to amqp message to support encoded format.
	contentType string
	// exchange, routingKey is amqp queue to publish message.
	exchange, routingKey string
	logger               zerolog.Logger
}

type ChangeEntityMessage struct {
	ID string
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
	alarmAdapter libalarm.Adapter,
	entityAdapter entity.Adapter,
	publisher libamqp.Publisher,
	encoder encoding.Encoder,
	contentType string,
	exchange, routingKey string,
	logger zerolog.Logger,
) EventPublisher {
	return &eventPublisher{
		alarmAdapter:  alarmAdapter,
		entityAdapter: entityAdapter,
		amqpPublisher: publisher,
		encoder:       encoder,
		contentType:   contentType,
		exchange:      exchange,
		routingKey:    routingKey,
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
		Connector:     types.ConnectorEngineService,
		ConnectorName: types.ConnectorEngineService,
		Component:     msg.ID,
		Timestamp:     types.NewCpsTime(),
		Author:        canopsis.DefaultEventAuthor,
		SourceType:    types.SourceTypeService,
	}

	err := p.publishEvent(ctx, event)
	if err != nil {
		p.logger.Err(err).Msg("cannot send event to amqp")
		return
	}

	p.logger.Debug().Msgf("publish %s", msg.ID)
}

func (p *eventPublisher) publishBasicEntityEvent(ctx context.Context, msg ChangeEntityMessage) {
	alarms, err := p.alarmAdapter.GetAlarmsByID(ctx, msg.ID)
	if err != nil {
		p.logger.Err(err).Msg("cannot find alarm")
		return
	}

	var event types.Event

	if len(alarms) == 0 {
		switch msg.EntityType {
		case types.EntityTypeComponent:
			connector, err := p.entityAdapter.FindConnector(ctx, msg.ID)
			if err != nil || connector == nil {
				p.logger.Warn().Str("entity", msg.ID).Msg("cannot generate event to component: no alarms and no connector")
				return
			}
			event = types.Event{
				Connector:     strings.ReplaceAll(connector.ID, "/"+connector.Name, ""),
				ConnectorName: connector.Name,
				Component:     msg.ID,
			}
		case types.EntityTypeResource:
			connector, err := p.entityAdapter.FindConnector(ctx, msg.ID)
			if err != nil || connector == nil {
				p.logger.Warn().Str("entity", msg.ID).Msg("cannot generate event to resource: no alarms and no connector")
				return
			}
			component, err := p.entityAdapter.FindComponent(ctx, msg.ID)
			if err != nil || component == nil {
				p.logger.Warn().Str("entity", msg.ID).Msg("cannot generate event to resource: no alarms and no component")
				return
			}

			event = types.Event{
				Connector:     strings.ReplaceAll(connector.ID, "/"+connector.Name, ""),
				ConnectorName: connector.Name,
				Component:     component.ID,
				Resource:      strings.ReplaceAll(msg.ID, "/"+component.ID, ""),
			}
		}
	} else {
		alarm := alarms[0]
		event = types.Event{
			Connector:     alarm.Value.Connector,
			ConnectorName: alarm.Value.ConnectorName,
			Component:     alarm.Value.Component,
			Resource:      alarm.Value.Resource,
		}
	}

	if msg.IsToggled {
		event.EventType = types.EventTypeEntityToggled
	} else if msg.IsDeleted {
		event.EventType = types.EventTypeResolveDeleted
	} else {
		event.EventType = types.EventTypeEntityUpdated
	}

	event.Timestamp = types.NewCpsTime()
	event.Author = canopsis.DefaultEventAuthor
	event.SourceType = event.DetectSourceType()
	err = p.publishEvent(ctx, event)
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
			Timestamp:     event.Timestamp,
			Author:        event.Author,
			SourceType:    types.SourceTypeResource,
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
