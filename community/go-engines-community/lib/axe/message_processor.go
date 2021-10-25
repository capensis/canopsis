package axe

import (
	"context"
	"runtime/trace"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statsng"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
)

type messageProcessor struct {
	FeaturePrintEventOnError bool
	FeatureStatEvents        bool
	EventProcessor           alarm.EventProcessor
	StatsService             statsng.Service
	RemediationRpcClient     engine.RPCClient
	TimezoneConfigProvider   config.TimezoneConfigProvider
	Encoder                  encoding.Encoder
	Decoder                  encoding.Decoder
	Logger                   zerolog.Logger
	PbehaviorAdapter         pbehavior.Adapter
}

func (p *messageProcessor) Process(parentCtx context.Context, d amqp.Delivery) ([]byte, error) {
	ctx, task := trace.NewTask(parentCtx, "axe.WorkerProcess")
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

	alarmChange, err := p.EventProcessor.Process(ctx, &event)
	if err != nil {
		if engine.IsConnectionError(err) {
			return nil, err
		}

		p.logError(err, "cannot process event", msg)
		return nil, nil
	}
	event.AlarmChange = &alarmChange

	err = p.handleRemediation(event, msg)
	if err != nil {
		return nil, err
	}

	p.updatePbhLastAlarmDate(ctx, event)
	p.handleStats(ctx, event, msg)

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

// updatePbhLastAlarmDate updates last time in pbehavior when it was applied on alarm.
func (p *messageProcessor) updatePbhLastAlarmDate(ctx context.Context, event types.Event) {
	if event.AlarmChange.Type != types.AlarmChangeTypeCreateAndPbhEnter &&
		event.AlarmChange.Type != types.AlarmChangeTypePbhEnter &&
		event.AlarmChange.Type != types.AlarmChangeTypePbhLeaveAndEnter {
		return
	}

	go func() {
		err := p.PbehaviorAdapter.UpdateLastAlarmDate(ctx, event.PbehaviorInfo.ID, types.CpsTime{Time: time.Now()})
		if err != nil {
			p.Logger.Err(err).Msg("")
		}
	}()
}

func (p *messageProcessor) handleStats(ctx context.Context, event types.Event, msg []byte) {
	if !p.FeatureStatEvents {
		return
	}

	if event.Alarm == nil {
		p.Logger.Warn().Msg("event.Alarm should not be nil")
		return
	}

	if event.Entity == nil {
		p.Logger.Warn().Msg("event.Entity should not be nil")
		return
	}

	go func() {
		err := p.StatsService.ProcessAlarmChange(
			ctx,
			*event.AlarmChange,
			event.Timestamp,
			*event.Alarm,
			*event.Entity,
			event.Author,
			event.EventType,
			p.TimezoneConfigProvider.Get().Location,
		)
		if err != nil {
			p.logError(err, "cannot update stats", msg)
		}
	}()
}

func (p *messageProcessor) handleRemediation(event types.Event, msg []byte) error {
	if p.RemediationRpcClient == nil {
		return nil
	}

	switch event.AlarmChange.Type {
	case types.AlarmChangeTypeCreate, types.AlarmChangeTypeCreateAndPbhEnter,
		types.AlarmChangeTypeResolve, types.AlarmChangeTypeStateDecrease,
		types.AlarmChangeTypeChangeState:
	default:
		return nil
	}

	body, err := p.Encoder.Encode(types.RPCRemediationEvent{
		Alarm:       event.Alarm,
		Entity:      event.Entity,
		AlarmChange: *event.AlarmChange,
	})
	if err != nil {
		p.logError(err, "failed to encode remediation event", msg)
		return nil
	}

	err = p.RemediationRpcClient.Call(engine.RPCMessage{
		CorrelationID: event.Alarm.ID,
		Body:          body,
	})
	if err != nil {
		if engine.IsConnectionError(err) {
			return err
		}

		p.logError(err, "failed to send rpc call to remediation", msg)
	}

	return nil
}

func (p *messageProcessor) logError(err error, errMsg string, msg []byte) {
	if p.FeaturePrintEventOnError {
		p.Logger.Err(err).Str("event", string(msg)).Msg(errMsg)
	} else {
		p.Logger.Err(err).Msg(errMsg)
	}
}
