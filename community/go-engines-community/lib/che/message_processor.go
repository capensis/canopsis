package che

import (
	"context"
	"errors"
	"runtime/trace"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libevent "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/che/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type messageProcessor struct {
	FeaturePrintEventOnError bool
	AlarmConfigProvider      config.AlarmConfigProvider
	MetricsConfigProvider    config.MetricsConfigProvider
	MetricsSender            metrics.Sender
	TechMetricsSender        techmetrics.Sender
	EntityCollection         mongo.DbCollection
	ExternalData             ExternalDataCoordinator
	PostProcessor            *EventPostProcessor
	Encoder                  encoding.Encoder
	Decoder                  encoding.Decoder
	Logger                   zerolog.Logger

	EventProcessorContainer libevent.ProcessorContainer
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

	event.Format()

	isSuspended := false
	defer func() {
		if !isSuspended {
			if eventMetric.Timestamp.IsZero() {
				eventMetric.Timestamp = start
			}

			eventMetric.Interval = time.Since(eventMetric.Timestamp)
			p.TechMetricsSender.SendCheEvent(eventMetric)
		}
	}()

	alarmConfig := p.AlarmConfigProvider.Get()
	if event.Initiator == types.InitiatorExternal {
		event.Output = utils.TruncateString(event.Output, alarmConfig.OutputLength)
		event.LongOutput = utils.TruncateString(event.LongOutput, alarmConfig.LongOutputLength)
	}

	proc, ok := p.EventProcessorContainer.Get(event.SourceType)
	if !ok {
		p.logError(err, "unsupported source type", d.Body)
		return nil, nil
	}

	pr, err := proc.Process(ctx, &event, nil)
	pr.EventMetric.Timestamp = start
	eventMetric = pr.EventMetric
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

	if pr.IsSuspended() {
		// It parks the partially processed event until the webhook RPC answers.
		// Processing continues when the RPC response is consumed and the event resumed.
		err = p.ExternalData.Dispatch(ctx, event, pr, 1)
		if err != nil {
			if engine.IsConnectionError(err) {
				return nil, err
			}

			p.logError(err, "cannot suspend event for external data", d.Body)

			return nil, nil
		}

		isSuspended = true

		return nil, engine.ErrAckWithoutForward
	}

	p.PostProcessor.Run(ctx, event, pr)
	p.handlePerfData(ctx, &event)

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

func (p *messageProcessor) logError(err error, errMsg string, msg []byte) {
	if p.FeaturePrintEventOnError {
		p.Logger.Err(err).Str("event", string(msg)).Msg(errMsg)
	} else {
		p.Logger.Err(err).Msg(errMsg)
	}
}

func (p *messageProcessor) handlePerfData(ctx context.Context, event *types.Event) {
	if (event.EventType != types.EventTypeCheck && event.EventType != types.EventTypeContextUpdate) || event.Entity == nil {
		return
	}

	perfData := event.GetPerfData(p.MetricsConfigProvider.Get().AllowedPerfDataUnits)
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
