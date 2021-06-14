package main

import (
	"context"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/engine"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/watcher"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
	"runtime/trace"
)

type messageProcessor struct {
	FeaturePrintEventOnError bool
	WatcherService           watcher.Service
	Encoder                  encoding.Encoder
	Decoder                  encoding.Decoder
	Logger                   zerolog.Logger
}

func (p *messageProcessor) Process(d amqp.Delivery) ([]byte, error) {
	ctx, task := trace.NewTask(context.Background(), "watcher.MessageProcessor")
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

	if event.EventType == types.EventTypeUpdateWatcher {
		err = p.WatcherService.UpdateWatcher(ctx, event.Component)
	} else if event.AlarmChange != nil && event.AlarmChange.Type == types.AlarmChangeTypeResolve && event.Alarm != nil && event.Entity != nil {
		err = p.WatcherService.ProcessResolvedAlarm(ctx, *event.Alarm, *event.Entity)
	} else {
		err = p.WatcherService.Process(ctx, event)
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
