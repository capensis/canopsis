package entityservice

import (
	"context"
	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"time"
)

// EventPublisher is used to send event to engines' event flow to notify about entity changes.
type EventPublisher interface {
	Publish(ctx context.Context, ch <-chan ChangeEntityMessage)
}

type eventPublisher struct {
	alarmAdapter  libalarm.Adapter
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
	IsToggled bool
	IsService bool
	// IsPatternChanged defines should service's context graph and state be recomputed.
	IsServicePatternChanged bool
	// ServiceAlarm is required on entity service delete because alarm is removed from
	// storage but alarm state is required by engine-service.
	ServiceAlarm *types.Alarm
}

func NewEventPublisher(
	alarmAdapter libalarm.Adapter,
	publisher libamqp.Publisher,
	encoder encoding.Encoder,
	contentType string,
	exchange, routingKey string,
	logger zerolog.Logger,
) EventPublisher {
	return &eventPublisher{
		alarmAdapter:  alarmAdapter,
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

			if msg.IsService {
				go p.publishServiceEvent(msg)
			} else {
				go p.publishBasicEntityEvent(ctx, msg)
			}
		}
	}
}

func (p *eventPublisher) publishServiceEvent(msg ChangeEntityMessage) {
	var eventType string
	if msg.IsServicePatternChanged || msg.IsToggled {
		eventType = types.EventTypeRecomputeEntityService
	} else {
		eventType = types.EventTypeEntityUpdated
	}

	event := types.Event{
		EventType:     eventType,
		Connector:     types.ConnectorEngineService,
		ConnectorName: types.ConnectorEngineService,
		Component:     msg.ID,
		Timestamp:     types.CpsTime{Time: time.Now()},
		Author:        canopsis.DefaultEventAuthor,
		SourceType:    types.SourceTypeService,
		Alarm:         msg.ServiceAlarm,
	}

	err := p.publishEvent(event)
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

	if len(alarms) == 0 {
		p.logger.Warn().Str("entity", msg.ID).Msg("no alarm for entity, cannot send event")
		return
	}

	var eventType string
	if msg.IsToggled {
		eventType = types.EventTypeEntityToggled
	} else {
		eventType = types.EventTypeEntityUpdated
	}

	alarm := alarms[0]
	event := types.Event{
		EventType:     eventType,
		Connector:     alarm.Value.Connector,
		ConnectorName: alarm.Value.ConnectorName,
		Component:     alarm.Value.Component,
		Resource:      alarm.Value.Resource,
		Timestamp:     types.CpsTime{Time: time.Now()},
		Author:        canopsis.DefaultEventAuthor,
	}
	event.SourceType = event.DetectSourceType()
	err = p.publishEvent(event)
	if err != nil {
		p.logger.Err(err).Msg("cannot send event to amqp")
		return
	}

	p.logger.Debug().Msgf("publish %s", msg.ID)
}

func (p *eventPublisher) publishEvent(event types.Event) error {
	body, err := p.encoder.Encode(event)
	if err != nil {
		return err
	}

	return p.amqpPublisher.Publish(
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
