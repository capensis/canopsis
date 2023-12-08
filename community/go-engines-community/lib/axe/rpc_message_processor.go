package axe

import (
	"context"
	"errors"
	"fmt"

	libevent "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/axe/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type rpcMessageProcessor struct {
	FeaturePrintEventOnError bool
	EventProcessor           libevent.Processor
	ActionRpc                engine.RPCClient
	PbhRpc                   engine.RPCClient
	DynamicInfosRpc          engine.RPCClient
	Encoder                  encoding.Encoder
	Decoder                  encoding.Decoder
	Logger                   zerolog.Logger
}

func (p *rpcMessageProcessor) Process(ctx context.Context, d amqp.Delivery) ([]byte, error) {
	msg := d.Body
	var event rpc.AxeEvent
	err := p.Decoder.Decode(msg, &event)
	if err != nil || event.Alarm == nil || event.Alarm.ID == "" || event.Entity == nil || event.Entity.ID == "" {
		p.logError(err, "RPC Message Processor: invalid event", msg)

		return p.getErrRpcEvent(errors.New("invalid event"), nil), nil
	}

	alarm := event.Alarm

	if alarm.IsResolved() {
		switch event.EventType {
		case types.EventTypeAutoWebhookStarted:
			/*do nothing*/
		case types.EventTypeAutoWebhookCompleted,
			types.EventTypeAutoWebhookFailed:
			body, err := p.Encoder.Encode(rpc.AxeResultEvent{
				Alarm:           alarm,
				WebhookHeader:   event.Parameters.WebhookHeader,
				WebhookResponse: event.Parameters.WebhookResponse,
			})
			if err != nil {
				p.logError(err, "RPC Message Processor: failed to encode rpc call to engine-action", msg)
				return nil, nil
			}

			err = p.ActionRpc.Call(ctx, engine.RPCMessage{
				CorrelationID: d.CorrelationId,
				Body:          body,
			})
			if err != nil {
				p.logError(err, "RPC Message Processor: failed to send rpc call to engine-action", msg)
			}

			return nil, nil
		default:
			p.logError(err, "RPC Message Processor: cannot update resolved alarm", msg)
			return p.getErrRpcEvent(errors.New("cannot update resolved alarm"), alarm), nil
		}
	}

	if event.EventType == types.ActionTypePbehavior {
		return p.processPbehaviorEvent(ctx, event, d)
	}

	if event.Parameters.Timestamp.Unix() <= 0 {
		event.Parameters.Timestamp = types.NewCpsTime()
	}

	if event.Parameters.Initiator == "" {
		event.Parameters.Initiator = types.InitiatorSystem
	}

	res, err := p.EventProcessor.Process(ctx, event)
	if err != nil {
		if engine.IsConnectionError(err) {
			return nil, err
		}

		p.logError(err, "RPC Message Processor: cannot update alarm", msg)

		return p.getErrRpcEvent(fmt.Errorf("cannot update alarm: %w", err), alarm), nil
	}

	if res.Alarm.ID != "" {
		alarm = &res.Alarm
	}

	p.sendEventToAction(ctx, *alarm, d.CorrelationId, event, res.AlarmChange, msg)

	if res.AlarmChange.Type == types.AlarmChangeTypeNone || p.DynamicInfosRpc == nil {
		return p.getRpcEvent(rpc.AxeResultEvent{
			AlarmChangeType: res.AlarmChange.Type,
			Alarm:           alarm,
		})
	}

	entity := *event.Entity
	if res.Entity.ID != "" {
		entity = res.Entity
	}

	return p.sendEventToDynamicInfos(ctx, res.Alarm, entity, res.AlarmChange, d)
}

func (p *rpcMessageProcessor) getRpcEvent(event rpc.AxeResultEvent) ([]byte, error) {
	msg, err := p.Encoder.Encode(event)
	if err != nil {
		p.Logger.Err(err).Msg("cannot encode event")
		return nil, nil
	}

	return msg, nil
}

func (p *rpcMessageProcessor) processPbehaviorEvent(ctx context.Context, event rpc.AxeEvent, d amqp.Delivery) ([]byte, error) {
	body, err := p.Encoder.Encode(rpc.PbehaviorEvent{
		Alarm:  event.Alarm,
		Entity: event.Entity,
		Params: rpc.PbehaviorParameters{
			Author:         event.Parameters.Author,
			UserID:         event.Parameters.User,
			Name:           event.Parameters.Name,
			Reason:         event.Parameters.Reason,
			Type:           event.Parameters.Type,
			RRule:          event.Parameters.RRule,
			Tstart:         event.Parameters.Tstart,
			Tstop:          event.Parameters.Tstop,
			StartOnTrigger: event.Parameters.StartOnTrigger,
			Duration:       event.Parameters.Duration,
		},
	})
	if err != nil {
		p.logError(err, "RPC Message Processor: failed to encode rpc call to pbehavior", d.Body)

		return p.getErrRpcEvent(fmt.Errorf("cannot encode rpc event : %w", err), event.Alarm), nil
	}

	err = p.PbhRpc.Call(ctx, engine.RPCMessage{
		CorrelationID: fmt.Sprintf("%s**%s", d.CorrelationId, d.ReplyTo),
		Body:          body,
	})
	if err != nil {
		if engine.IsConnectionError(err) {
			return nil, err
		}

		p.logError(err, "RPC Message Processor: failed to send rpc call to pbehavior", d.Body)

		return p.getErrRpcEvent(fmt.Errorf("failed to send rpc call to pbehavior : %w", err), event.Alarm), nil
	}

	return nil, nil
}

func (p *rpcMessageProcessor) sendEventToAction(
	ctx context.Context,
	alarm types.Alarm,
	correlationID string,
	event rpc.AxeEvent,
	alarmChange types.AlarmChange,
	msg []byte,
) {
	switch event.EventType {
	case types.EventTypeAutoWebhookCompleted,
		types.EventTypeAutoWebhookFailed:
	default:
		return
	}

	body, err := p.Encoder.Encode(rpc.AxeResultEvent{
		Alarm:           &alarm,
		AlarmChangeType: alarmChange.Type,
		WebhookHeader:   event.Parameters.WebhookHeader,
		WebhookResponse: event.Parameters.WebhookResponse,
		Error:           event.Parameters.WebhookError,
	})
	if err != nil {
		p.logError(err, "RPC Message Processor: failed to encode rpc call to engine-action", msg)
		return
	}

	err = p.ActionRpc.Call(ctx, engine.RPCMessage{
		CorrelationID: correlationID,
		Body:          body,
	})
	if err != nil {
		p.logError(err, "RPC Message Processor: failed to send rpc call to engine-action", msg)
	}
}

func (p *rpcMessageProcessor) sendEventToDynamicInfos(
	ctx context.Context,
	alarm types.Alarm,
	entity types.Entity,
	alarmChange types.AlarmChange,
	d amqp.Delivery,
) ([]byte, error) {
	body, err := p.Encoder.Encode(rpc.DynamicInfosEvent{
		Alarm:           &alarm,
		Entity:          &entity,
		AlarmChangeType: alarmChange.Type,
	})
	if err != nil {
		p.logError(err, "failed to encode rpc call to dynamic-infos", d.Body)

		return p.getErrRpcEvent(fmt.Errorf("cannot encode rpc event: %w", err), &alarm), nil
	}

	err = p.DynamicInfosRpc.Call(ctx, engine.RPCMessage{
		CorrelationID: fmt.Sprintf("%s**%s", d.CorrelationId, d.ReplyTo),
		Body:          body,
	})
	if err != nil {
		if engine.IsConnectionError(err) {
			return nil, err
		}

		p.logError(err, "failed to send rpc call to dynamic-infos", d.Body)

		return p.getErrRpcEvent(fmt.Errorf("failed to send rpc call to dynamic-infos: %w", err), &alarm), nil
	}

	return nil, nil
}

func (p *rpcMessageProcessor) logError(err error, errMsg string, msg []byte) {
	if p.FeaturePrintEventOnError {
		p.Logger.Debug().Err(err).Str("event", string(msg)).Msg(errMsg)
	} else {
		p.Logger.Err(err).Msg(errMsg)
	}
}

func (p *rpcMessageProcessor) getErrRpcEvent(err error, alarm *types.Alarm) []byte {
	msg, _ := p.getRpcEvent(rpc.AxeResultEvent{
		Alarm: alarm,
		Error: &rpc.Error{Error: err}},
	)
	return msg
}
