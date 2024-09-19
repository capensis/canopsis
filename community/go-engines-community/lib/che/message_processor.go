package che

import (
	"context"
	"errors"
	"runtime/trace"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/che/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

type messageProcessor struct {
	FeaturePrintEventOnError bool
	AlarmConfigProvider      config.AlarmConfigProvider
	MetricsSender            metrics.Sender
	MetaUpdater              metrics.MetaUpdater
	TechMetricsSender        techmetrics.Sender
	AmqpPublisher            libamqp.Publisher
	EntityCollection         mongo.DbCollection
	Encoder                  encoding.Encoder
	Decoder                  encoding.Decoder
	Logger                   zerolog.Logger

	EventProcessorContainer event.ProcessorContainer
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
	if event.Initiator == types.InitiatorExternal {
		event.Output = utils.TruncateString(event.Output, alarmConfig.OutputLength)
		event.LongOutput = utils.TruncateString(event.LongOutput, alarmConfig.LongOutputLength)
	}

	var updatedEntitiesForEvent []types.Entity
	var updatedEntityIdsForMetrics []string

	proc, ok := p.EventProcessorContainer.Get(event.SourceType)
	if !ok {
		p.logError(err, "unsupported source type", d.Body)
		return nil, nil
	}

	updatedEntitiesForEvent, updatedEntityIdsForMetrics, eventMetric, err = proc.Process(ctx, &event)
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

	go p.postProcessUpdatedEntities(ctx, event, updatedEntitiesForEvent, updatedEntityIdsForMetrics)

	p.handlePerfData(ctx, &event)

	event.Format()

	body, err := p.Encoder.Encode(&event)
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

func (p *messageProcessor) postProcessUpdatedEntities(
	ctx context.Context,
	event types.Event,
	entitiesForEvent []types.Entity,
	updatedEntityIdsForMetrics []string,
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

	p.MetaUpdater.UpdateById(ctx, updatedEntityIdsForMetrics...)
}

func (p *messageProcessor) logError(err error, errMsg string, msg []byte) {
	if p.FeaturePrintEventOnError {
		p.Logger.Err(err).Str("event", string(msg)).Msg(errMsg)
	} else {
		p.Logger.Err(err).Msg(errMsg)
	}
}

func (p *messageProcessor) handlePerfData(ctx context.Context, event *types.Event) {
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
				"$set":      bson.M{"perf_data_updated": datetime.CpsTime{Time: now}},
			})
			if err != nil {
				p.Logger.Err(err).Str("entity", event.Entity.ID).Msg("cannot update entity perf data")
			}
		}()
	}
}
