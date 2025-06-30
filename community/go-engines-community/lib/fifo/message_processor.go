package fifo

import (
	"context"
	"errors"
	"fmt"
	"runtime/trace"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/scheduler"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/bsm/redislock"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

const (
	lockBackoff = time.Second
	// lockRetries is the number of retries to obtain or refresh the redis lock.
	lockRetries     = 15
	lockMinDuration = 2 * time.Minute
)

type messageProcessor struct {
	EventFilterService eventfilter.Service
	Scheduler          scheduler.Scheduler
	MetricsSender      metrics.Sender
	Decoder            encoding.Decoder
	Logger             zerolog.Logger

	TechMetricsSender techmetrics.Sender

	FeaturePrintEventOnError bool
	exclusiveProcessor       bool
}

func NewMessageProcessor(
	eventFilterService eventfilter.Service,
	scheduler scheduler.Scheduler,
	metricsSender metrics.Sender,
	decoder encoding.Decoder,
	logger zerolog.Logger,
	techMetricsSender techmetrics.Sender,
	featurePrintEvenotOnError bool,
) *messageProcessor {
	return &messageProcessor{
		EventFilterService:       eventFilterService,
		Scheduler:                scheduler,
		MetricsSender:            metricsSender,
		Decoder:                  decoder,
		Logger:                   logger,
		TechMetricsSender:        techMetricsSender,
		FeaturePrintEventOnError: featurePrintEvenotOnError,
		exclusiveProcessor:       false,
	}
}

func (p *messageProcessor) Process(parentCtx context.Context, d amqp.Delivery) ([]byte, error) {
	if !p.exclusiveProcessor {
		return nil, errors.New("exclusive processor is not active, cannot process message")
	}

	eventMetric := techmetrics.FifoEventMetric{}
	eventMetric.Timestamp = time.Now()

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
		p.TechMetricsSender.SendFifoEvent(eventMetric)
	}()

	event.Format()
	event.ReceivedTimestamp = datetime.NewMicroTime()
	p.MetricsSender.SendMessageRate(time.Now(), event.EventType, event.ConnectorName)

	err = event.InjectExtraInfos(msg)
	if err != nil {
		p.logError(err, "cannot inject extra infos", msg)
		return nil, nil
	}

	if !event.Healthcheck {
		_, _, eventMetric.ExternalRequests, err = p.EventFilterService.ProcessEvent(ctx, &event)
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

func (p *messageProcessor) RefreshExclusiveProcessor(ctx context.Context, refreshInterval, lockDuration time.Duration, l redis.Lock) {
	if l == nil {
		p.Logger.Warn().Msg("exclusive processor lock is nil, cannot refresh exclusive processor")
		return
	}
	p.exclusiveProcessor = true
	go func(ctx context.Context, l redis.Lock) {
		var err error
		defer func() {
			p.Logger.Err(err).Msg("failed to refresh redis lock")
			p.exclusiveProcessor = false
		}()
		ticker := time.NewTicker(refreshInterval)
		defer ticker.Stop()
		var retryTicker *time.Ticker

		retry := newRetryStrategy()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// Refresh doesn't handle retry strategy in v0.9.4, so we need to handle it manually.
				for {
					err = l.Refresh(ctx, lockDuration, nil)
					if err != nil {
						backoff := retry.NextBackoff()
						if backoff < 1 {
							return
						}
						if retryTicker == nil {
							retryTicker = time.NewTicker(backoff)
							defer retryTicker.Stop()
						} else {
							retryTicker.Reset(backoff)
						}
						select {
						case <-ctx.Done():
							return
						case <-retryTicker.C:
						}
						continue
					}
					break
				}
				// on success re-define retry for the next Refresh call
				retry = newRetryStrategy()
			}
		}
	}(ctx, l)
}

// newRetryStrategy is used to create a new retry strategy with reset inner counters.
func newRetryStrategy() redislock.RetryStrategy {
	return redislock.LimitRetry(redislock.LinearBackoff(lockBackoff), lockRetries)
}
