package che

import (
	"context"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/che/event"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

// EventPostProcessor performs the work that follows a successfully processed event,
// shared by the primary and the resume (webhook RPC) paths:
// emitting counter-update events for changed component/connector entities and refreshing entity metadata.
type EventPostProcessor struct {
	amqpPublisher libamqp.Publisher
	metaUpdater   metrics.MetaUpdater
	encoder       encoding.Encoder
	logger        zerolog.Logger
}

func NewEventPostProcessor(
	amqpPublisher libamqp.Publisher,
	metaUpdater metrics.MetaUpdater,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) *EventPostProcessor {
	return &EventPostProcessor{
		amqpPublisher: amqpPublisher,
		metaUpdater:   metaUpdater,
		encoder:       encoder,
		logger:        logger,
	}
}

// Run applies the post-processing for a finished event.
func (p *EventPostProcessor) Run(ctx context.Context, event types.Event, pr event.ProcessorResult) {
	go p.postProcessUpdatedEntities(ctx, event, pr.UpdatedEntitiesForEvent, pr.UpdatedEntityIdsForMetrics, pr.ServicesIDsToRecompute)
}

func (p *EventPostProcessor) postProcessUpdatedEntities(
	ctx context.Context,
	event types.Event,
	entitiesForEvent []types.Entity,
	updatedEntityIdsForMetrics []string,
	servicesIDsToRecompute []string,
) {
	now := datetime.NewCpsTime()

	for _, ent := range entitiesForEvent {
		var updateCountersEvent types.Event

		switch ent.Type {
		case types.EntityTypeComponent:
			updateCountersEvent = types.Event{
				EventType:     types.EventTypeUpdateCounters,
				SourceType:    types.SourceTypeComponent,
				Connector:     canopsis.CheConnector,
				ConnectorName: canopsis.CheConnector,
				Component:     ent.Component,
				Timestamp:     now,
				Entity:        &ent,
				Author:        canopsis.DefaultEventAuthor,
				Initiator:     types.InitiatorSystem,
			}
		case types.EntityTypeConnector:
			updateCountersEvent = types.Event{
				EventType:     types.EventTypeUpdateCounters,
				SourceType:    types.SourceTypeConnector,
				Connector:     event.Connector,
				ConnectorName: event.ConnectorName,
				Timestamp:     now,
				Entity:        &ent,
				Author:        canopsis.DefaultEventAuthor,
				Initiator:     types.InitiatorSystem,
			}
		}

		body, err := p.encoder.Encode(updateCountersEvent)
		if err != nil {
			p.logger.Err(err).Msg("unable to serialize event")
			continue
		}

		err = p.amqpPublisher.PublishWithContext(
			ctx,
			canopsis.EngineExchangeName,
			canopsis.AxeSystemQueueName,
			false,
			false,
			amqp.Publishing{
				Body:         body,
				ContentType:  canopsis.JsonContentType,
				DeliveryMode: amqp.Persistent,
			},
		)
		if err != nil {
			p.logger.Err(err).Msg("unable to send service event")
		}
	}

	for _, id := range servicesIDsToRecompute {
		body, err := p.encoder.Encode(types.Event{
			EventType:     types.EventTypeRecomputeEntityService,
			Connector:     canopsis.CheConnector,
			ConnectorName: canopsis.CheConnector,
			Component:     id,
			Timestamp:     datetime.NewCpsTime(),
			SourceType:    types.SourceTypeService,
			Author:        canopsis.DefaultEventAuthor,
			Initiator:     types.InitiatorSystem,
		})
		if err != nil {
			p.logger.Err(err).Msg("unable to serialize event")
			continue
		}

		err = p.amqpPublisher.PublishWithContext(
			ctx,
			canopsis.DefaultExchangeName,
			canopsis.FIFOQueueName,
			false,
			false,
			amqp.Publishing{
				Body:         body,
				ContentType:  canopsis.JsonContentType,
				DeliveryMode: amqp.Persistent,
			},
		)
		if err != nil {
			p.logger.Err(err).Msg("unable to send service event")
		}
	}

	p.metaUpdater.UpdateById(ctx, updatedEntityIdsForMetrics...)
}
