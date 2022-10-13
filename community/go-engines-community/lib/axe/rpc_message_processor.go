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
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type rpcMessageProcessor struct {
	FeaturePrintEventOnError bool
	RMQChannel               libamqp.Channel
	PbhRpc                   engine.RPCClient
	ServiceRpc               engine.RPCClient
	RemediationRpc           engine.RPCClient
	Executor                 operation.Executor
	MetaAlarmEventProcessor  libalarm.MetaAlarmEventProcessor
	Decoder                  encoding.Decoder
	Encoder                  encoding.Encoder
	Logger                   zerolog.Logger
}

func (p *rpcMessageProcessor) Process(ctx context.Context, d amqp.Delivery) ([]byte, error) {
	msg := d.Body
	var event types.RPCAxeEvent
	err := p.Decoder.Decode(msg, &event)
	if err != nil || event.Alarm == nil {
		p.logError(err, "RPC Message Processor: invalid event", msg)

		return p.getErrRpcEvent(errors.New("invalid event"), nil), nil
	}

	alarm := event.Alarm

	if alarm.IsResolved() {
		p.logError(err, "RPC Message Processor: cannot update resolved alarm", msg)

		return p.getErrRpcEvent(errors.New("cannot update resolved alarm"), alarm), nil
	}

	operationType := event.EventType

	if operationType == types.ActionTypePbehavior {
		body, err := p.Encoder.Encode(types.RPCPBehaviorEvent{
			Alarm:  alarm,
			Entity: event.Entity,
			Params: types.RPCPBehaviorParameters{
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
			p.logError(err, "RPC Message Processor: failed to encode rpc call to pbehavior", msg)

			return p.getErrRpcEvent(fmt.Errorf("cannot encode rpc event : %v", err), alarm), nil
		}

		err = p.PbhRpc.Call(ctx, engine.RPCMessage{
			CorrelationID: fmt.Sprintf("%s**%s", d.CorrelationId, d.ReplyTo),
			Body:          body,
		})
		if err != nil {
			if engine.IsConnectionError(err) {
				return nil, err
			}

			p.logError(err, "RPC Message Processor: failed to send rpc call to pbehavior", msg)

			return p.getErrRpcEvent(fmt.Errorf("failed to send rpc call to pbehavior : %v", err), alarm), nil
		}

		return nil, nil
	}

	alarmChange := types.AlarmChange{
		Type:                            types.AlarmChangeTypeNone,
		PreviousState:                   alarm.Value.State.Value,
		PreviousStateChange:             alarm.Value.State.Timestamp,
		PreviousStatus:                  alarm.Value.Status.Value,
		PreviousStatusChange:            alarm.Value.Status.Timestamp,
		PreviousPbehaviorTypeID:         alarm.Value.PbehaviorInfo.TypeID,
		PreviousPbehaviorCannonicalType: alarm.Value.PbehaviorInfo.CanonicalType,
	}

	op := types.Operation{
		Type: operationType,
		Parameters: types.OperationParameters{
			Output:      event.Parameters.Output,
			Author:      event.Parameters.Author,
			User:        event.Parameters.User,
			Ticket:      event.Parameters.Ticket,
			Duration:    event.Parameters.Duration,
			State:       event.Parameters.State,
			Execution:   event.Parameters.Execution,
			Instruction: event.Parameters.Instruction,
		},
	}
	alarmChangeType, err := p.Executor.Exec(ctx, op, alarm, event.Entity,
		types.CpsTime{Time: time.Now()}, "", "", types.InitiatorSystem)
	if err != nil {
		if engine.IsConnectionError(err) {
			return nil, err
		}

		p.logError(err, "RPC Message Processor: cannot update alarm", msg)
		return p.getErrRpcEvent(fmt.Errorf("cannot update alarm: %v", err), alarm), nil
	}
	alarmChange.Type = alarmChangeType

	if alarmChangeType == types.AlarmChangeTypeAck ||
		alarmChangeType == types.AlarmChangeTypeAckremove ||
		alarmChangeType == types.AlarmChangeTypeChangeState {
		body, err := p.Encoder.Encode(types.RPCServiceEvent{
			Alarm:       alarm,
			Entity:      event.Entity,
			AlarmChange: &alarmChange,
		})
		if err != nil {
			p.logError(err, "RPC Message Processor: failed to encode rpc call to engine-service", msg)
		} else {
			err = p.ServiceRpc.Call(ctx, engine.RPCMessage{
				CorrelationID: utils.NewID(),
				Body:          body,
			})
			if err != nil {
				p.logError(err, "RPC Message Processor: failed to send rpc call to engine-service", msg)
			}
		}
	}

	if event.Entity != nil &&
		alarmChangeType == types.AlarmChangeTypeAutoInstructionFail ||
		alarmChangeType == types.AlarmChangeTypeInstructionJobFail ||
		alarmChangeType == types.AlarmChangeTypeInstructionJobComplete ||
		alarmChangeType == types.AlarmChangeTypeAutoInstructionComplete {
		body, err := p.Encoder.Encode(types.Event{
			EventType:     types.EventTypeTrigger,
			Connector:     alarm.Value.Connector,
			ConnectorName: alarm.Value.ConnectorName,
			Component:     alarm.Value.Component,
			Resource:      alarm.Value.Resource,
			SourceType:    event.Entity.Type,
			AlarmChange:   &alarmChange,
		})
		if err != nil {
			p.logError(err, "RPC Message Processor: failed to encode a trigger event to engine-fifo", msg)
		}

		err = p.RMQChannel.PublishWithContext(
			ctx,
			"",
			canopsis.FIFOQueueName,
			false,
			false,
			amqp.Publishing{
				ContentType:  "application/json",
				Body:         body,
				DeliveryMode: amqp.Persistent,
			},
		)
		if err != nil {
			p.logError(err, "RPC Message Processor: failed to send a trigger event to engine-fifo", msg)
		}
	}

	if p.RemediationRpc != nil && alarmChangeType == types.AlarmChangeTypeChangeState {
		body, err := p.Encoder.Encode(types.RPCRemediationEvent{
			Alarm:       alarm,
			Entity:      event.Entity,
			AlarmChange: alarmChange,
		})
		if err != nil {
			p.logError(err, "RPC Message Processor: failed to encode rpc call to engine-remediation", msg)
		} else {
			err = p.RemediationRpc.Call(ctx, engine.RPCMessage{
				CorrelationID: event.Alarm.ID,
				Body:          body,
			})
			if err != nil {
				p.logError(err, "RPC Message Processor: failed to send rpc call to engine-remediation", msg)
			}
		}
	}

	res := types.RPCAxeResultEvent{
		AlarmChangeType: alarmChangeType,
		Alarm:           alarm,
	}

	err = p.MetaAlarmEventProcessor.ProcessAxeRpc(ctx, event, res)
	if err != nil {
		p.logError(err, "failed to process meta alarm", msg)
	}

	return p.getRpcEvent(res)
}

func (p *rpcMessageProcessor) logError(err error, errMsg string, msg []byte) {
	if p.FeaturePrintEventOnError {
		p.Logger.Debug().Err(err).Str("event", string(msg)).Msg(errMsg)
	} else {
		p.Logger.Err(err).Msg(errMsg)
	}
}

func (p *rpcMessageProcessor) getErrRpcEvent(err error, alarm *types.Alarm) []byte {
	msg, _ := p.getRpcEvent(types.RPCAxeResultEvent{
		Alarm: alarm,
		Error: &types.RPCError{Error: err}},
	)
	return msg
}

func (p *rpcMessageProcessor) getRpcEvent(event types.RPCAxeResultEvent) ([]byte, error) {
	msg, err := p.Encoder.Encode(event)
	if err != nil {
		p.Logger.Err(err).Msg("cannot encode event")
		return nil, nil
	}

	return msg, nil
}
