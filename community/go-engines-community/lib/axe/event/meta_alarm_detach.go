package event

import (
	"context"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"github.com/rs/zerolog"
)

func NewMetaAlarmDetachProcessor(
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor,
	autoInstructionMatcher AutoInstructionMatcher,
	metricsSender metrics.Sender,
	remediationRpcClient engine.RPCClient,
	encoder encoding.Encoder,
	amqpPublisher libamqp.Publisher,
	logger zerolog.Logger,
) Processor {
	return &metaAlarmDetachProcessor{
		metaAlarmEventProcessor: metaAlarmEventProcessor,
		autoInstructionMatcher:  autoInstructionMatcher,
		metricsSender:           metricsSender,
		remediationRpcClient:    remediationRpcClient,
		encoder:                 encoder,
		amqpPublisher:           amqpPublisher,
		logger:                  logger,
	}
}

type metaAlarmDetachProcessor struct {
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor
	autoInstructionMatcher  AutoInstructionMatcher
	metricsSender           metrics.Sender
	remediationRpcClient    engine.RPCClient
	encoder                 encoding.Encoder
	amqpPublisher           libamqp.Publisher
	logger                  zerolog.Logger
}

func (p *metaAlarmDetachProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil {
		return result, nil
	}

	_, err := p.metaAlarmEventProcessor.DetachChildrenFromMetaAlarm(ctx, event)
	if err != nil {
		return result, err
	}

	result.Forward = false

	return result, nil
}
