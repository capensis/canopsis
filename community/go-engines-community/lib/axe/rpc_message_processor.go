package axe

import (
	"context"
	"errors"
	"fmt"
	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
	"time"
)

type rpcMessageProcessor struct {
	FeaturePrintEventOnError bool
	PbhRpc                   engine.RPCClient
	ServiceRpc               engine.RPCClient
	RemediationRpc           engine.RPCClient
	Executor                 operation.Executor
	AlarmAdapter             libalarm.Adapter
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
	op := types.Operation{
		Type:       operationType,
		Parameters: event.Parameters,
	}

	if operationType == types.ActionTypePbehavior {
		if params, ok := event.Parameters.(types.ActionPBehaviorParameters); ok {
			body, err := p.Encoder.Encode(types.RPCPBehaviorEvent{
				Alarm:  alarm,
				Entity: event.Entity,
				Params: params,
			})
			if err != nil {
				p.logError(err, "RPC Message Processor: failed to encode rpc call to pbehavior", msg)

				return p.getErrRpcEvent(fmt.Errorf("cannot encode rpc event : %v", err), alarm), nil
			}

			err = p.PbhRpc.Call(engine.RPCMessage{
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

		err := errors.New("invalid pbh parameters")
		p.logError(err, "RPC Message Processor: invalid event", msg)

		return p.getErrRpcEvent(err, alarm), nil
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
	alarmChangeType, err := p.Executor.Exec(ctx, op, alarm,
		types.CpsTime{Time: time.Now()}, "", types.InitiatorSystem)
	if err != nil {
		if engine.IsConnectionError(err) {
			return nil, err
		}

		p.logError(err, "RPC Message Processor: cannot update alarm", msg)
		return p.getErrRpcEvent(fmt.Errorf("cannot update alarm: %v", err), alarm), nil
	}
	alarmChange.Type = alarmChangeType
	if alarm.IsMetaAlarm() {
		var childrenAlarms []types.Alarm
		err := p.AlarmAdapter.GetOpenedAlarmsByIDs(ctx, event.Alarm.Value.Children, &childrenAlarms)
		if err != nil {
			p.logError(err, "RPC Message Processor: error getting meta-alarm children", msg)
		} else {
			for _, childAlarm := range childrenAlarms {
				_, err = p.Executor.Exec(ctx, op, &childAlarm, types.CpsTime{Time: time.Now()}, "", types.InitiatorSystem)
				if err != nil {
					p.logError(err, "RPC Message Processor: cannot update child alarm", msg)
				}
			}
		}
	}

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
			err = p.ServiceRpc.Call(engine.RPCMessage{
				CorrelationID: utils.NewID(),
				Body:          body,
			})
			if err != nil {
				p.logError(err, "RPC Message Processor: failed to send rpc call to engine-service", msg)
			}
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
			err = p.RemediationRpc.Call(engine.RPCMessage{
				CorrelationID: event.Alarm.ID,
				Body:          body,
			})
			if err != nil {
				p.logError(err, "RPC Message Processor: failed to send rpc call to engine-remediation", msg)
			}
		}
	}

	return p.getRpcEvent(types.RPCAxeResultEvent{
		AlarmChangeType: alarmChangeType,
		Alarm:           alarm,
	})
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
