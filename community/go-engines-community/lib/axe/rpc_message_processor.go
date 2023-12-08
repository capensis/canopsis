package axe

import (
	"context"
	"errors"
	"fmt"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type rpcMessageProcessor struct {
	FeaturePrintEventOnError bool
	RMQChannel               libamqp.Channel
	PbhRpc                   engine.RPCClient
	DynamicInfosRpc          engine.RPCClient
	ServiceRpc               engine.RPCClient
	RemediationRpc           engine.RPCClient
	ActionRpc                engine.RPCClient
	Executor                 operation.Executor
	MetaAlarmEventProcessor  libalarm.MetaAlarmEventProcessor
	Decoder                  encoding.Decoder
	Encoder                  encoding.Encoder
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

	alarmChange, err := p.executeOperation(ctx, event)
	if err != nil {
		if engine.IsConnectionError(err) {
			return nil, err
		}

		p.logError(err, "RPC Message Processor: cannot update alarm", msg)

		return p.getErrRpcEvent(fmt.Errorf("cannot update alarm: %w", err), alarm), nil
	}

	p.sendEventToService(ctx, *alarm, *event.Entity, alarmChange, msg)
	p.sendTriggerEvent(ctx, event, alarmChange, msg)
	p.sendEventToRemediation(ctx, *alarm, *event.Entity, alarmChange, msg)
	p.sendEventToAction(ctx, *alarm, d.CorrelationId, event, alarmChange, msg)

	res := rpc.AxeResultEvent{
		AlarmChangeType: alarmChange.Type,
		Alarm:           alarm,
	}
	err = p.MetaAlarmEventProcessor.ProcessAxeRpc(ctx, event, res)
	if err != nil {
		p.logError(err, "failed to process meta alarm", msg)
	}

	if alarmChange.Type == types.AlarmChangeTypeNone || p.DynamicInfosRpc == nil {
		return p.getRpcEvent(res)
	}

	return p.sendEventToDynamicInfos(ctx, *alarm, *event.Entity, alarmChange, d)
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

func (p *rpcMessageProcessor) executeOperation(ctx context.Context, event rpc.AxeEvent) (types.AlarmChange, error) {
	alarm := event.Alarm
	alarmChange := types.NewAlarmChangeByAlarm(*alarm)

	op := types.Operation{
		Type: event.EventType,
		Parameters: types.OperationParameters{
			Output:            event.Parameters.Output,
			Author:            event.Parameters.Author,
			User:              event.Parameters.User,
			Duration:          event.Parameters.Duration,
			State:             event.Parameters.State,
			Execution:         event.Parameters.Execution,
			Instruction:       event.Parameters.Instruction,
			TicketInfo:        event.Parameters.TicketInfo,
			WebhookRequest:    event.Parameters.WebhookRequest,
			WebhookFailReason: event.Parameters.WebhookFailReason,
		},
	}
	alarmChangeType, err := p.Executor.Exec(ctx, op, alarm, event.Entity,
		types.CpsTime{Time: time.Now()}, "", "", types.InitiatorSystem)
	if err != nil {
		return types.AlarmChange{}, err
	}

	alarmChange.Type = alarmChangeType

	return alarmChange, nil
}

func (p *rpcMessageProcessor) sendEventToService(
	ctx context.Context,
	alarm types.Alarm,
	entity types.Entity,
	alarmChange types.AlarmChange,
	msg []byte,
) {
	switch alarmChange.Type {
	case types.AlarmChangeTypeAck,
		types.AlarmChangeTypeAckremove,
		types.AlarmChangeTypeChangeState:
	default:
		return
	}

	body, err := p.Encoder.Encode(types.RPCServiceEvent{
		Alarm:       &alarm,
		Entity:      &entity,
		AlarmChange: &alarmChange,
	})
	if err != nil {
		p.logError(err, "RPC Message Processor: failed to encode rpc call to engine-service", msg)
		return
	}

	err = p.ServiceRpc.Call(ctx, engine.RPCMessage{
		CorrelationID: utils.NewID(),
		Body:          body,
	})
	if err != nil {
		p.logError(err, "RPC Message Processor: failed to send rpc call to engine-service", msg)
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

func (p *rpcMessageProcessor) sendEventToRemediation(
	ctx context.Context,
	alarm types.Alarm,
	entity types.Entity,
	alarmChange types.AlarmChange,
	msg []byte,
) {
	if p.RemediationRpc == nil {
		return
	}
	if alarmChange.Type != types.AlarmChangeTypeChangeState {
		return
	}

	body, err := p.Encoder.Encode(types.RPCRemediationEvent{
		Alarm:       &alarm,
		Entity:      &entity,
		AlarmChange: alarmChange,
	})
	if err != nil {
		p.logError(err, "RPC Message Processor: failed to encode rpc call to engine-remediation", msg)
		return
	}

	err = p.RemediationRpc.Call(ctx, engine.RPCMessage{
		CorrelationID: alarm.ID,
		Body:          body,
	})
	if err != nil {
		p.logError(err, "RPC Message Processor: failed to send rpc call to engine-remediation", msg)
	}
}

func (p *rpcMessageProcessor) sendTriggerEvent(
	ctx context.Context,
	event rpc.AxeEvent,
	alarmChange types.AlarmChange,
	msg []byte,
) {
	switch alarmChange.Type {
	case types.AlarmChangeTypeAutoInstructionFail,
		types.AlarmChangeTypeAutoInstructionComplete,
		types.AlarmChangeTypeInstructionJobFail,
		types.AlarmChangeTypeInstructionJobComplete:
	case types.AlarmChangeTypeDeclareTicketWebhook:
		if !event.Parameters.EmitTrigger {
			return
		}
	default:
		return
	}

	body, err := p.Encoder.Encode(types.Event{
		EventType:     types.EventTypeTrigger,
		Connector:     event.Alarm.Value.Connector,
		ConnectorName: event.Alarm.Value.ConnectorName,
		Component:     event.Alarm.Value.Component,
		Resource:      event.Alarm.Value.Resource,
		SourceType:    event.Entity.Type,
		AlarmChange:   &alarmChange,
		AlarmID:       event.Alarm.ID,
	})
	if err != nil {
		p.logError(err, "RPC Message Processor: failed to encode a trigger event to engine-fifo", msg)
		return
	}

	err = p.RMQChannel.PublishWithContext(
		ctx,
		"",
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
		p.logError(err, "RPC Message Processor: failed to send a trigger event to engine-fifo", msg)
	}
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
