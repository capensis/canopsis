package event

import (
	"context"

	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
)

func NewMetaAlarmProcessor(
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor,
	autoInstructionMatcher AutoInstructionMatcher,
	metricsSender metrics.Sender,
	remediationRpcClient engine.RPCClient,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Processor {
	return &metaAlarmProcessor{
		metaAlarmEventProcessor: metaAlarmEventProcessor,
		autoInstructionMatcher:  autoInstructionMatcher,
		metricsSender:           metricsSender,
		remediationRpcClient:    remediationRpcClient,
		encoder:                 encoder,
		logger:                  logger,
	}
}

type metaAlarmProcessor struct {
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor
	autoInstructionMatcher  AutoInstructionMatcher
	metricsSender           metrics.Sender
	remediationRpcClient    engine.RPCClient
	encoder                 encoding.Encoder
	logger                  zerolog.Logger
}

func (p *metaAlarmProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil {
		return result, nil
	}

	alarm, err := p.metaAlarmEventProcessor.CreateMetaAlarm(ctx, event)
	if err != nil {
		return result, err
	}

	alarmChange := types.NewAlarmChange()
	alarmChange.Type = types.AlarmChangeTypeCreate
	result.Forward = true
	result.Alarm = *alarm
	result.AlarmChange = alarmChange
	result.IsInstructionMatched = isInstructionMatched(event, result, p.autoInstructionMatcher, p.logger)
	go p.postProcess(context.Background(), event, result)

	return result, nil
}

func (p *metaAlarmProcessor) postProcess(
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
