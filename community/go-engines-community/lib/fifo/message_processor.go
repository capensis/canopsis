package fifo

import (
	"context"
	"errors"
	"fmt"
	"runtime/trace"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/ratelimit"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/scheduler"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type messageProcessor struct {
	FeaturePrintEventOnError bool

	EventFilterService eventfilter.Service
	Scheduler          scheduler.Scheduler
	StatsSender        ratelimit.StatsSender
	Decoder            encoding.Decoder
	Logger             zerolog.Logger

	TechMetricsSender techmetrics.Sender
}

func (p *messageProcessor) Process(parentCtx context.Context, d amqp.Delivery) ([]byte, error) {
	eventMetric := techmetrics.EventMetric{
		Timestamp: time.Now(),
	}

	ctx, task := trace.NewTask(parentCtx, "fifo.WorkerProcess")
	defer task.End()

	msg := d.Body
	trace.Logf(ctx, "event_size", "%d", len(msg))

	var event types.Event
	err := p.Decoder.Decode(msg, &event)
	if err != nil {
		p.logError(err, "cannot decode event", msg)
		return nil, nil
	}

	p.Logger.Debug().Msgf("valid input event: %v", string(msg))
	trace.Log(ctx, "event.event_type", event.EventType)
	trace.Log(ctx, "event.timestamp", event.Timestamp.String())
	trace.Log(ctx, "event.source_type", event.SourceType)
	trace.Log(ctx, "event.connector", event.Connector)
	trace.Log(ctx, "event.connector_name", event.ConnectorName)
	trace.Log(ctx, "event.component", event.Component)
	trace.Log(ctx, "event.resource", event.Resource)

	err = event.IsValid()
	if err != nil {
		p.logError(err, "invalid event", msg)
		return nil, nil
	}

	defer func() {
		eventMetric.EventType = event.EventType
		eventMetric.Interval = time.Since(eventMetric.Timestamp)
		p.TechMetricsSender.SendSimpleEvent(techmetrics.FIFOEvent, eventMetric)
	}()

	event.Format()
	p.StatsSender.Add(time.Now().Unix(), true)

	err = event.InjectExtraInfos(msg)
	if err != nil {
		p.logError(err, "cannot inject extra infos", msg)
		return nil, nil
	}

	if !event.Healthcheck {
		_, err = p.EventFilterService.ProcessEvent(ctx, &event)
		if err != nil {
			if errors.Is(err, eventfilter.ErrDropOutcome) {
				return nil, nil
			}

			p.logError(err, "cannot process event by eventfilter service", msg)
			return nil, nil
		}
	}

	p.Logger.Debug().Str("event", fmt.Sprintf("%+v", event)).Msg("sent to scheduler")
	err = p.Scheduler.ProcessEvent(ctx, event)
	if err != nil {
		if engine.IsConnectionError(err) {
			return nil, err
		}

		p.logError(err, "cannot process event", msg)
		return nil, nil
	}

	return nil, nil
}

func (p *messageProcessor) logError(err error, errMsg string, msg []byte) {
	if p.FeaturePrintEventOnError {
		p.Logger.Err(err).Str("event", string(msg)).Msg(errMsg)
	} else {
		p.Logger.Err(err).Msg(errMsg)
	}
}
