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
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

type messageProcessor struct {
	FeaturePrintEventOnError bool
	DbClient                 mongo.DbClient
	AlarmConfigProvider      config.AlarmConfigProvider
	EventFilterService       eventfilter.Service
	MetricsSender            metrics.Sender
	MetaUpdater              metrics.MetaUpdater
	TechMetricsSender        techmetrics.Sender
	ContextGraphManager      contextgraph.Manager
	AmqpPublisher            libamqp.Publisher
	EntityCollection         mongo.DbCollection
	Encoder                  encoding.Encoder
	Decoder                  encoding.Decoder
	Logger                   zerolog.Logger
}

func (p *messageProcessor) Process(ctx context.Context, d amqp.Delivery) ([]byte, error) {
	var eventMetric techmetrics.CheEventMetric
	start := time.Now()

	ctx, task := trace.NewTask(ctx, "che.WorkerProcess")
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

	defer func() {
		eventMetric.Interval = time.Since(start)
		eventMetric.Timestamp = start

		p.TechMetricsSender.SendCheEvent(eventMetric)
	}()

	alarmConfig := p.AlarmConfigProvider.Get()
	event.Output = utils.TruncateString(event.Output, alarmConfig.OutputLength)
	event.LongOutput = utils.TruncateString(event.LongOutput, alarmConfig.LongOutputLength)
	var updatedEntities []types.Entity
	event, updatedEntities, eventMetric, err = p.handleEvent(ctx, event)
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

	go p.postProcessUpdatedEntities(ctx, event, updatedEntities)

	p.handlePerfData(ctx, event)

	event.Format()
	body, err := p.Encoder.Encode(event)
	if err != nil {
		p.logError(err, "cannot encode event", d.Body)
		return nil, nil
	}

	if event.Healthcheck {
		_, err := p.EntityCollection.DeleteMany(ctx, bson.M{"healthcheck": true})
		if err != nil {
			p.logError(err, "cannot delete temporary entity", d.Body)
		}
	}

	return body, nil
}

func (p *messageProcessor) handleEvent(ctx context.Context, event types.Event) (types.Event, []types.Entity, techmetrics.CheEventMetric, error) {
	var eventMetric techmetrics.CheEventMetric
	eventMetric.EventType = event.EventType

	if event.EventType == types.EventTypeManualMetaAlarmGroup ||
		event.EventType == types.EventTypeManualMetaAlarmUngroup ||
		event.EventType == types.EventTypeManualMetaAlarmUpdate ||
		event.EventType == types.EventTypeMetaAlarmUngroup {
		return event, nil, eventMetric, nil
	}

	if event.EventType == types.EventTypeRecomputeEntityService {
		var updatedEntities []types.Entity
		var eventEntity types.Entity
		err := p.DbClient.WithTransaction(ctx, func(ctx context.Context) error {
			updatedEntities = make([]types.Entity, 0, len(updatedEntities))
			eventEntity = types.Entity{}
			var err error
			eventEntity, updatedEntities, err = p.ContextGraphManager.RecomputeService(ctx, event.GetEID())
			if err != nil {
				return fmt.Errorf("cannot recompute service %s: %w", event.Component, err)
			}

			_, err = p.ContextGraphManager.UpdateEntities(ctx, "", updatedEntities, true)
			if err != nil {
				return fmt.Errorf("cannot update entities %s: %w", event.Component, err)
			}

			return nil
		})
		if err != nil {
			return event, nil, eventMetric, err
		}

		event.Entity = &eventEntity
		eventMetric.EntityType = eventEntity.Type

		return event, updatedEntities, eventMetric, nil
	}

	var eventEntity types.Entity
	var updatedEntityIds []string

	err := p.DbClient.WithTransaction(ctx, func(ctx context.Context) error {
		eventEntity = types.Entity{}
		updatedEntityIds = make([]string, 0, len(updatedEntityIds))
		var contextGraphEntities []types.Entity
		var err error
		eventEntity, contextGraphEntities, err = p.ContextGraphManager.HandleEvent(ctx, event)
		if err != nil {
			return fmt.Errorf("cannot update context graph: %w", err)
		}

		if !eventEntity.IsNew && !eventEntity.IsUpdated && eventEntity.LastEventDate != nil {
			err := p.ContextGraphManager.UpdateLastEventDate(ctx, event.EventType, eventEntity.ID, *eventEntity.LastEventDate)
			if err != nil {
				return err
			}
		}

		updatedEntities := make([]types.Entity, 0, len(contextGraphEntities)+1)
		if eventEntity.IsNew ||
			eventEntity.IsUpdated ||
			event.EventType == types.EventTypeEntityUpdated ||
			event.EventType == types.EventTypeEntityToggled ||
			event.SourceType == types.SourceTypeService {
			updatedEntities = append(updatedEntities, eventEntity)
			updatedEntityIds = append(updatedEntityIds, eventEntity.ID)
		}

		updatedEntities = append(updatedEntities, contextGraphEntities...)
		for _, e := range contextGraphEntities {
			updatedEntityIds = append(updatedEntityIds, e.ID)
		}

		_, err = p.ContextGraphManager.UpdateEntities(ctx, "", updatedEntities, false)
		if err != nil {
			return fmt.Errorf("cannot update entities: %w", err)
		}

		return nil
	})
	if err != nil {
		return event, nil, eventMetric, err
	}

	eventMetric.EntityType = eventEntity.Type
	eventMetric.IsNewEntity = eventEntity.IsNew

	eventEntity.IsNew = false
	eventEntity.IsUpdated = false

	event.Entity = &eventEntity
	// Process event by event filters.
	if event.Entity != nil && !event.Healthcheck {
		if event.Entity.Enabled {
			event, err = p.EventFilterService.ProcessEvent(ctx, event)
			if err != nil {
				return event, nil, eventMetric, err
			}

			if event.Entity.IsUpdated {
				updatedEntityIds = append(updatedEntityIds, event.Entity.ID)

				_, err = p.EntityCollection.UpdateOne(
					ctx,
					bson.M{"_id": event.Entity.ID},
					bson.M{"$set": bson.M{"infos": event.Entity.Infos}},
				)
				if err != nil {
					return event, nil, eventMetric, fmt.Errorf("cannot update entities: %w", err)
				}
			}
		} else if event.IsOnlyServiceUpdate() {
			updatedEntityIds = append(updatedEntityIds, event.Entity.ID)
		}
	}

	if len(updatedEntityIds) == 0 || event.Healthcheck {
		return event, nil, eventMetric, nil
	}

	eventMetric.IsInfosUpdated = event.Entity.IsUpdated

	var serviceUpdatedEntities []types.Entity

	err = p.DbClient.WithTransaction(ctx, func(ctx context.Context) error {
		serviceUpdatedEntities = nil
		eventMetric.IsServicesUpdated = false
		eventEntity = *event.Entity
		cursor, err := p.EntityCollection.Find(ctx, bson.M{"_id": bson.M{"$in": updatedEntityIds}})
		if err != nil {
			return err
		}

		updatedEntities := make([]types.Entity, 0, len(updatedEntityIds))
		err = cursor.All(ctx, &updatedEntities)
		if err != nil {
			return err
		}

		for _, entity := range updatedEntities {
			if entity.ID == eventEntity.ID {
				eventEntity = entity
				break
			}
		}

		// if it's a new resource add a component info to check if component is matched by the service
		serviceUpdatedEntities, err = p.ContextGraphManager.CheckServices(ctx, updatedEntities)
		if err != nil {
			return fmt.Errorf("cannot check services: %w", err)
		}

		if eventEntity.Enabled && (eventMetric.IsInfosUpdated || event.EventType == types.EventTypeEntityUpdated && eventEntity.Type == types.EntityTypeComponent) {
			resources, err := p.ContextGraphManager.FillResourcesWithInfos(ctx, eventEntity)
			if err != nil {
				return fmt.Errorf("cannot update entity infos: %w", err)
			}

			serviceUpdatedEntities = append(serviceUpdatedEntities, resources...)
		}

		updatedEventEntity, err := p.ContextGraphManager.UpdateEntities(ctx, eventEntity.ID, serviceUpdatedEntities, true)
		if err != nil {
			return fmt.Errorf("cannot update entities: %w", err)
		}

		if updatedEventEntity.ID != "" {
			eventEntity = updatedEventEntity
		}

		eventMetric.IsServicesUpdated = len(eventEntity.ServicesToAdd) > 0 || len(eventEntity.ServicesToRemove) > 0
		return nil
	})
	if err != nil {
		return event, nil, eventMetric, err
	}

	event.Entity = &eventEntity

	return event, serviceUpdatedEntities, eventMetric, nil
}

func (p *messageProcessor) postProcessUpdatedEntities(ctx context.Context, event types.Event, updatedEntities []types.Entity) {
	entityIDs := make([]string, len(updatedEntities))

	for idx, ent := range updatedEntities {
		entityIDs[idx] = ent.ID

		if (len(ent.ServicesToAdd) != 0 || len(ent.ServicesToRemove) != 0) && ent.ID != event.GetEID() && ent.Type != types.EntityTypeService {
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
					Initiator:     types.InitiatorSystem,
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
					Initiator:     types.InitiatorSystem,
				}
			case types.EntityTypeConnector:
				updateCountersEvent = types.Event{
					EventType:     types.EventTypeUpdateCounters,
					SourceType:    types.SourceTypeConnector,
					Connector:     event.Connector,
					ConnectorName: event.ConnectorName,
					Timestamp:     types.CpsTime{Time: time.Now()},
					Entity:        &ent,
					Initiator:     types.InitiatorSystem,
				}
			}

			body, err := p.Encoder.Encode(updateCountersEvent)
			if err != nil {
				p.Logger.Err(err).Msg("unable to serialize event")
			}

			err = p.AmqpPublisher.PublishWithContext(
				ctx,
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

func (p *messageProcessor) handlePerfData(ctx context.Context, event types.Event) {
	if event.EventType != types.EventTypeCheck || event.Entity == nil {
		return
	}

	perfData := event.GetPerfData()
	now := time.Now()
	names := make([]string, len(perfData))
	for i, v := range perfData {
		p.MetricsSender.SendPerfData(now, event.Entity.ID, v.Name, v.Value, v.Unit)
		names[i] = v.Name
	}

	if len(names) > 0 {
		go func() {
			_, err := p.EntityCollection.UpdateOne(ctx, bson.M{"_id": event.Entity.ID}, bson.M{
				"$addToSet": bson.M{"perf_data": bson.M{"$each": names}},
				"$set":      bson.M{"perf_data_updated": types.CpsTime{Time: now}},
			})
			if err != nil {
				p.Logger.Err(err).Str("entity", event.Entity.ID).Msg("cannot update entity perf data")
			}
		}()
	}
}
