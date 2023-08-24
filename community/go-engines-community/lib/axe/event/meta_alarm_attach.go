package event

import (
	"context"
	"fmt"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	amqp "github.com/rabbitmq/amqp091-go"

	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
)

func NewMetaAlarmAttachProcessor(
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor,
	autoInstructionMatcher AutoInstructionMatcher,
	metricsSender metrics.Sender,
	remediationRpcClient engine.RPCClient,
	encoder encoding.Encoder,
	amqpPublisher libamqp.Publisher,
	logger zerolog.Logger,
) Processor {
	return &metaAlarmAttachProcessor{
		metaAlarmEventProcessor: metaAlarmEventProcessor,
		autoInstructionMatcher:  autoInstructionMatcher,
		metricsSender:           metricsSender,
		remediationRpcClient:    remediationRpcClient,
		encoder:                 encoder,
		amqpPublisher:           amqpPublisher,
		logger:                  logger,
	}
}

type metaAlarmAttachProcessor struct {
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor
	autoInstructionMatcher  AutoInstructionMatcher
	metricsSender           metrics.Sender
	remediationRpcClient    engine.RPCClient
	encoder                 encoding.Encoder
	amqpPublisher           libamqp.Publisher
	logger                  zerolog.Logger
}

func (p *metaAlarmAttachProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil {
		return result, nil
	}

	_, updatedChildrenAlarms, updatedChildrenEvents, err := p.metaAlarmEventProcessor.AttachChildrenToMetaAlarm(ctx, event)
	if err != nil {
		return result, err
	}

	result.Forward = false

	for _, child := range updatedChildrenAlarms {
		p.metricsSender.SendCorrelation(event.Parameters.Timestamp.Time, child)
	}

	for _, childEvent := range updatedChildrenEvents {
		err = p.sendToFifo(ctx, childEvent)
		if err != nil {
			return result, err
		}
	}

	go p.postProcess(context.Background(), event, result)

	return result, nil
}

func (p *metaAlarmAttachProcessor) postProcess(
	ctx context.Context,
	event rpc.AxeEvent,
	result Result,
) {
	p.metricsSender.SendEventMetrics(
		result.Alarm,
		*event.Entity,
		result.AlarmChange,
		event.Parameters.Timestamp.Time,
		event.Parameters.Initiator,
		event.Parameters.User,
		event.Parameters.Instruction,
		"",
	)

	err := sendRemediationEvent(ctx, event, result, p.remediationRpcClient, p.encoder)
	if err != nil {
		p.logger.Err(err).Msg("cannot send event to engine-remediation")
	}
}

func (p *metaAlarmAttachProcessor) sendToFifo(ctx context.Context, event types.Event) error {
	body, err := p.encoder.Encode(event)
	if err != nil {
		return fmt.Errorf("cannot encode event: %w", err)
	}

	err = p.amqpPublisher.PublishWithContext(
		ctx,
		canopsis.FIFOExchangeName,
		canopsis.FIFOQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  canopsis.JsonContentType,
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		return fmt.Errorf("cannot send child event: %w", err)
	}

	return nil
}
