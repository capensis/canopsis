package che

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"runtime/trace"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/contextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
)

type messageProcessor struct {
	FeaturePrintEventOnError bool
	DbClient                 mongo.DbClient
	AlarmConfigProvider      config.AlarmConfigProvider
	EventFilterService       eventfilter.Service
	ContextGraphManager      contextgraph.Manager
	AmqpPublisher            libamqp.Publisher
	Encoder                  encoding.Encoder
	Decoder                  encoding.Decoder
	Logger                   zerolog.Logger
}

func (p *messageProcessor) Process(parentCtx context.Context, d amqp.Delivery) ([]byte, error) {
	t := time.Now()
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

	err = p.DbClient.WithTransaction(ctx, func(tCtx context.Context) error {
		if event.EventType == types.EventTypeRecomputeEntityService {
			updatedEntities, err := p.ContextGraphManager.RecomputeService(ctx, event.GetEID())
			if err != nil {
				return fmt.Errorf("cannot recompute service: %w", err)
			}

			err = p.ContextGraphManager.UpdateEntities(ctx, updatedEntities)
			if err != nil {
				return fmt.Errorf("cannot update entities: %w", err)
			}

			return nil
		}

		eventEntity, err := p.ContextGraphManager.Handle(ctx, event)
		if err != nil {
			return fmt.Errorf("cannot update context graph: %w", err)
		}

		if !eventEntity.Enabled {
			return nil
		}

		event.Entity = &eventEntity

		// Process event by event filters.
		event, err = p.EventFilterService.ProcessEvent(ctx, event)
		if err != nil {
			return fmt.Errorf("cannot apply event filters on event: %w", err)
		}

		eventEntity = *event.Entity

		var updatedEntities []types.Entity
		if eventEntity.IsNew || event.IsEntityUpdated {
			updatedEntities = []types.Entity{eventEntity}
		}

		if event.IsEntityUpdated && eventEntity.Type == types.EntityTypeComponent {
			resources, err := p.ContextGraphManager.FillResourcesWithInfos(ctx, eventEntity)
			if err != nil {
				return fmt.Errorf("cannot update entity infos: %w", err)
			}

			updatedEntities = append(updatedEntities, resources...)
		}

		if len(updatedEntities) > 0 {
			updatedEntities, err = p.ContextGraphManager.CheckServices(ctx, updatedEntities)
			if err != nil {
				return fmt.Errorf("cannot check services: %w", err)
			}

			err = p.ContextGraphManager.UpdateEntities(ctx, updatedEntities)
			if err != nil {
				return fmt.Errorf("cannot update entities: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		if engine.IsConnectionError(err) {
			return nil, err
		}

		p.logError(err, "cannot process event", d.Body)
		return nil, nil
	}

	event.Format()
	body, err := p.Encoder.Encode(event)
	if err != nil {
		p.logError(err, "cannot encode event", d.Body)
		return nil, nil
	}

	fmt.Printf("all = %s\n", time.Since(t))

	return body, nil
}

func (p *messageProcessor) logError(err error, errMsg string, msg []byte) {
	if p.FeaturePrintEventOnError {
		p.Logger.Err(err).Str("event", string(msg)).Msg(errMsg)
	} else {
		p.Logger.Err(err).Msg(errMsg)
	}
}

func (p *messageProcessor) publishToEngineFIFO(event types.Event) error {
	body, err := p.Encoder.Encode(event)
	if err != nil {
		p.Logger.Err(err).Msg("cannot encode event")
		return nil
	}
	return errt.NewIOError(p.AmqpPublisher.Publish(
		"",
		canopsis.FIFOQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json", // this type is mandatory to avoid bad conversions into Python.
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	))
}
