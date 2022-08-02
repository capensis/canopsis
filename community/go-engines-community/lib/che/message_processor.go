package che

import (
	"context"
	"errors"
	"fmt"
	"runtime/trace"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/contextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type messageProcessor struct {
	FeaturePrintEventOnError bool
	DbClient                 mongo.DbClient
	AlarmConfigProvider      config.AlarmConfigProvider
	EventFilterService       eventfilter.Service
	MetaUpdater              metrics.MetaUpdater
	ContextGraphManager      contextgraph.Manager
	AmqpPublisher            libamqp.Publisher
	Encoder                  encoding.Encoder
	Decoder                  encoding.Decoder
	Logger                   zerolog.Logger
}

func (p *messageProcessor) Process(parentCtx context.Context, d amqp.Delivery) ([]byte, error) {
	ctx, task := trace.NewTask(parentCtx, "che.WorkerProcess")
	defer task.End()

	trace.Logf(ctx, "event_size", "%d", len(d.Body))
	var event types.Event
	err := p.Decoder.Decode(d.Body, &event)
	if err != nil {
		p.logError(err, "cannot decode event", d.Body)
		return nil, nil
	}

	err = event.InjectExtraInfos(d.Body)
	if err != nil {
		p.logError(err, "cannot inject extra infos", d.Body)
		return nil, nil
	}

	trace.Log(ctx, "event.event_type", event.EventType)
	trace.Log(ctx, "event.timestamp", event.Timestamp.String())
	trace.Log(ctx, "event.source_type", event.SourceType)
	trace.Log(ctx, "event.connector", event.Connector)
	trace.Log(ctx, "event.connector_name", event.ConnectorName)
	trace.Log(ctx, "event.component", event.Component)
	trace.Log(ctx, "event.resource", event.Resource)

	err = event.IsValid()
	if err != nil {
		p.logError(err, "invalid event", d.Body)
		return nil, nil
	}

	alarmConfig := p.AlarmConfigProvider.Get()
	event.Output = utils.TruncateString(event.Output, alarmConfig.OutputLength)
	event.LongOutput = utils.TruncateString(event.LongOutput, alarmConfig.LongOutputLength)

	var updatedEntities []types.Entity

	err = p.DbClient.WithTransaction(ctx, func(tCtx context.Context) error {
		if event.EventType == types.EventManualMetaAlarmGroup ||
			event.EventType == types.EventManualMetaAlarmUngroup ||
			event.EventType == types.EventManualMetaAlarmUpdate {
			return nil
		}

		if event.EventType == types.EventTypeRecomputeEntityService {
			eventEntity, updatedEntities, err := p.ContextGraphManager.RecomputeService(tCtx, event.GetEID())
			if err != nil {
				return fmt.Errorf("cannot recompute service %s: %w", event.Component, err)
			}

			event.Entity = &eventEntity

			_, err = p.ContextGraphManager.UpdateEntities(tCtx, "", updatedEntities)
			if err != nil {
				return fmt.Errorf("cannot update entities %s: %w", event.Component, err)
			}

			return nil
		}

		eventEntity, contextGraphEntities, err := p.ContextGraphManager.HandleEvent(tCtx, event)
		if err != nil {
			return fmt.Errorf("cannot update context graph: %w", err)
		}

		event.Entity = &eventEntity

		// Process event by event filters.
		event, err = p.EventFilterService.ProcessEvent(tCtx, event)
		if err != nil {
			return err
		}

		eventEntity = *event.Entity
		if eventEntity.IsNew ||
			event.IsEntityUpdated ||
			event.EventType == types.EventTypeEntityUpdated ||
			event.EventType == types.EventTypeEntityToggled ||
			event.SourceType == types.SourceTypeService {
			updatedEntities = []types.Entity{eventEntity}
		}

		updatedEntities = append(updatedEntities, contextGraphEntities...)

		if len(updatedEntities) > 0 {
			// if it's a new resource add a component info to check if component is matched by the service
			updatedEntities, err = p.ContextGraphManager.CheckServices(tCtx, updatedEntities)
			if err != nil {
				return fmt.Errorf("cannot check services: %w", err)
			}

			if event.IsEntityUpdated || event.EventType == types.EventTypeEntityUpdated && eventEntity.Type == types.EntityTypeComponent {
				resources, err := p.ContextGraphManager.FillResourcesWithInfos(tCtx, eventEntity)
				if err != nil {
					return fmt.Errorf("cannot update entity infos: %w", err)
				}

				updatedEntities = append(updatedEntities, resources...)
			}

			eventEntity, err = p.ContextGraphManager.UpdateEntities(tCtx, eventEntity.ID, updatedEntities)
			if err != nil {
				return fmt.Errorf("cannot update entities: %w", err)
			}

			event.Entity = &eventEntity
		}

		if !eventEntity.IsNew && !event.IsEntityUpdated && event.Entity.LastEventDate != nil {
			err := p.ContextGraphManager.UpdateLastEventDate(tCtx, event.EventType, event.Entity.ID, *event.Entity.LastEventDate)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, eventfilter.ErrDropOutcome) {
			return nil, nil
		}

		if engine.IsConnectionError(err) {
			return nil, err
		}

		p.logError(err, "cannot process event", d.Body)
		return nil, nil
	}

	go p.postProcessUpdatedEntities(event, updatedEntities)

	event.Format()
	body, err := p.Encoder.Encode(event)
	if err != nil {
		p.logError(err, "cannot encode event", d.Body)
		return nil, nil
	}

	return body, nil
}

func (p *messageProcessor) postProcessUpdatedEntities(event types.Event, updatedEntities []types.Entity) {
	entityIDs := make([]string, len(updatedEntities))

	for idx, ent := range updatedEntities {
		entityIDs[idx] = ent.ID

		if (len(ent.ImpactedServicesToAdd) != 0 || len(ent.ImpactedServicesToRemove) != 0) && ent.ID != event.GetEID() && ent.Type != types.EntityTypeService {
			var updateCountersEvent types.Event

			switch ent.Type {
			case types.EntityTypeResource:
				updateCountersEvent = types.Event{
					EventType:     types.EventTypeUpdateCounters,
					SourceType:    types.SourceTypeResource,
					Connector:     event.Connector,
					ConnectorName: event.ConnectorName,
					Component:     ent.Component,
					Resource:      ent.Name,
					Timestamp:     types.CpsTime{Time: time.Now()},
					Entity:        &ent,
				}
			case types.EntityTypeComponent:
				updateCountersEvent = types.Event{
					EventType:     types.EventTypeUpdateCounters,
					SourceType:    types.SourceTypeComponent,
					Connector:     event.Connector,
					ConnectorName: event.ConnectorName,
					Component:     ent.Component,
					Timestamp:     types.CpsTime{Time: time.Now()},
					Entity:        &ent,
				}
			case types.EntityTypeConnector:
				updateCountersEvent = types.Event{
					EventType:     types.EventTypeUpdateCounters,
					SourceType:    types.SourceTypeConnector,
					Connector:     event.Connector,
					ConnectorName: event.ConnectorName,
					Timestamp:     types.CpsTime{Time: time.Now()},
					Entity:        &ent,
				}
			}

			body, err := p.Encoder.Encode(updateCountersEvent)
			if err != nil {
				p.Logger.Err(err).Msg("unable to serialize event")
			}

			err = p.AmqpPublisher.Publish(
				"",
				canopsis.AxeQueueName,
				false,
				false,
				amqp.Publishing{
					Body:        body,
					ContentType: "application/json",
				},
			)
			if err != nil {
				p.Logger.Err(err).Msg("unable to send service event")
			}
		}
	}

	p.MetaUpdater.UpdateById(context.Background(), entityIDs...)
}

func (p *messageProcessor) logError(err error, errMsg string, msg []byte) {
	if p.FeaturePrintEventOnError {
		p.Logger.Err(err).Str("event", string(msg)).Msg(errMsg)
	} else {
		p.Logger.Err(err).Msg(errMsg)
	}
}
