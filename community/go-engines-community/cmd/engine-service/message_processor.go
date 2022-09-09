package main

import (
	"context"
	"runtime/trace"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type messageProcessor struct {
	TechMetricsSender        techmetrics.Sender
	EntityServiceService     entityservice.Service
	Encoder                  encoding.Encoder
	Decoder                  encoding.Decoder
	Logger                   zerolog.Logger
	FeaturePrintEventOnError bool
}

func (p *messageProcessor) Process(ctx context.Context, d amqp.Delivery) ([]byte, error) {
	eventMetric := techmetrics.EventMetric{}
	eventMetric.Timestamp = time.Now()

	ctx, task := trace.NewTask(ctx, "service.MessageProcessor")
	defer task.End()

	msg := d.Body

	trace.Logf(ctx, "event_size", "%d", len(msg))
	var err error
	var event types.Event
	trace.WithRegion(ctx, "decode-event", func() {
		err = p.Decoder.Decode(msg, &event)
	})
	if err != nil {
		p.logError(err, "cannot decode event", msg)
		return nil, nil
	}

	p.Logger.Debug().Msgf("unmarshaled: %+v", event)
	trace.Log(ctx, "event.event_type", event.EventType)
	trace.Log(ctx, "event.timestamp", event.Timestamp.String())
	trace.Log(ctx, "event.source_type", event.SourceType)
	trace.Log(ctx, "event.connector", event.Connector)
	trace.Log(ctx, "event.connector_name", event.ConnectorName)
	trace.Log(ctx, "event.component", event.Component)
	trace.Log(ctx, "event.resource", event.Resource)

	defer func() {
		eventMetric.EventType = event.EventType
		eventMetric.Interval = time.Since(eventMetric.Timestamp)
		p.TechMetricsSender.SendSimpleEvent(techmetrics.ServiceEvent, eventMetric)
	}()

	if event.EventType == types.EventTypeUpdateEntityService {
		err = p.EntityServiceService.ReloadService(ctx, event.GetEID())
	} else if event.EventType == types.EventTypeRecomputeEntityService {
		err = p.EntityServiceService.UpdateService(ctx, event)
	} else {
		err = p.EntityServiceService.Process(ctx, event)
	}

	if err != nil {
		if engine.IsConnectionError(err) {
			return nil, err
		}

		p.logError(err, "cannot process event", msg)
		return nil, nil
	}

	// Encode and publish the event to the next engine
	var bevent []byte
	trace.WithRegion(ctx, "encode-event", func() {
		bevent, err = p.Encoder.Encode(event)
	})
	if err != nil {
		p.logError(err, "cannot encode event", msg)
		return nil, nil
	}

	return bevent, nil
}

func (p *messageProcessor) logError(err error, errMsg string, msg []byte) {
	if p.FeaturePrintEventOnError {
		p.Logger.Err(err).Str("event", string(msg)).Msg(errMsg)
	} else {
		p.Logger.Err(err).Msg(errMsg)
	}
}
