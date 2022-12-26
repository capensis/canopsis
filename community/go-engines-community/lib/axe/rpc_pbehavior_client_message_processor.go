package axe

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	libentity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type rpcPBehaviorClientMessageProcessor struct {
	FeaturePrintEventOnError bool
	PublishCh                libamqp.Channel
	ServiceRpc               engine.RPCClient
	RemediationRpc           engine.RPCClient
	Executor                 operation.Executor
	EntityAdapter            libentity.Adapter
	PbehaviorAdapter         pbehavior.Adapter
	Decoder                  encoding.Decoder
	Encoder                  encoding.Encoder
	Logger                   zerolog.Logger
}

func (p *rpcPBehaviorClientMessageProcessor) Process(ctx context.Context, msg engine.RPCMessage) error {
	data := strings.Split(msg.CorrelationID, "**")
	if len(data) != 2 {
		return fmt.Errorf("RPC PBehavior Client: bad correlation_id: %s", msg.CorrelationID)
	}

	correlationId := data[0]
	routingKey := data[1]

	var replyEvent []byte
	var event types.RPCPBehaviorResultEvent
	err := p.Decoder.Decode(msg.Body, &event)
	if err != nil || event.Alarm == nil || event.Entity == nil {
		p.logError(err, "RPC PBehavior Client: invalid event", msg.Body)

		return p.publishResult(ctx, routingKey, correlationId, p.getErrRpcEvent(fmt.Errorf("invalid event")))
	}

	alarmChangeType := types.AlarmChangeTypeNone

	if event.PbhEvent.EventType != "" {
		alarmChange := types.AlarmChange{
			Type:                            types.AlarmChangeTypeNone,
			PreviousState:                   event.Alarm.Value.State.Value,
			PreviousStateChange:             event.Alarm.Value.State.Timestamp,
			PreviousStatus:                  event.Alarm.Value.Status.Value,
			PreviousStatusChange:            event.Alarm.Value.Status.Timestamp,
			PreviousPbehaviorTypeID:         event.Alarm.Value.PbehaviorInfo.TypeID,
			PreviousPbehaviorCannonicalType: event.Alarm.Value.PbehaviorInfo.CanonicalType,
		}
		alarmChangeType, err = p.Executor.Exec(
			ctx,
			types.Operation{
				Type: event.PbhEvent.EventType,
				Parameters: types.OperationParameters{
					PbehaviorInfo: &event.PbhEvent.PbehaviorInfo,
					Author:        event.PbhEvent.Author,
					Output:        event.PbhEvent.Output,
				},
			},
			event.Alarm,
			event.Entity,
			event.PbhEvent.Timestamp,
			"",
			"",
			types.InitiatorSystem,
		)
		if err != nil {
			if engine.IsConnectionError(err) {
				return err
			}

			p.logError(err, "RPC PBehavior Client: cannot update alarm", msg.Body)
			return p.publishResult(ctx, routingKey, correlationId, p.getErrRpcEvent(fmt.Errorf("cannot update alarm: %v", err)))
		}

		p.updateEntity(ctx, event.PbhEvent.Entity, *event.Alarm, alarmChangeType)
		go p.updatePbhLastAlarmDate(ctx, alarmChangeType, event.Alarm.Value.PbehaviorInfo)

		alarmChange.Type = alarmChangeType
		body, err := p.Encoder.Encode(types.RPCServiceEvent{
			Alarm:       event.Alarm,
			Entity:      event.PbhEvent.Entity,
			AlarmChange: &alarmChange,
		})
		if err != nil {
			p.logError(err, "RPC PBehavior Client: failed to encode rpc call to engine-service", msg.Body)
		} else {
			err = p.ServiceRpc.Call(ctx, engine.RPCMessage{
				CorrelationID: utils.NewID(),
				Body:          body,
			})
			if err != nil {
				if engine.IsConnectionError(err) {
					return err
				}

				p.logError(err, "RPC PBehavior Client: failed to send rpc call to engine-service", msg.Body)
			}
		}
		if p.RemediationRpc != nil {
			body, err = p.Encoder.Encode(types.RPCRemediationEvent{
				Alarm:       event.Alarm,
				Entity:      event.PbhEvent.Entity,
				AlarmChange: alarmChange,
			})
			if err != nil {
				p.logError(err, "RPC PBehavior Client: failed to encode rpc call to engine-remediation", msg.Body)
			} else {
				err = p.RemediationRpc.Call(ctx, engine.RPCMessage{
					CorrelationID: utils.NewID(),
					Body:          body,
				})
				if err != nil {
					if engine.IsConnectionError(err) {
						return err
					}

					p.logError(err, "RPC PBehavior Client: failed to send rpc call to engine-remediation", msg.Body)
				}
			}
		}
	}

	replyEvent, err = p.getRpcEvent(rpc.AxeResultEvent{
		Alarm:           event.Alarm,
		AlarmChangeType: alarmChangeType,
		Error:           nil,
	})
	if err != nil {
		p.logError(err, "RPC PBehavior Client: failed to create rpc result event", msg.Body)

		replyEvent = p.getErrRpcEvent(errors.New("failed to create rpc result event"))
	}

	err = p.publishResult(ctx, routingKey, correlationId, replyEvent)
	if err != nil {
		p.logError(err, "RPC PBehavior Client: cannot sent message result back to sender", msg.Body)

		return err
	}

	return nil
}

func (p *rpcPBehaviorClientMessageProcessor) publishResult(ctx context.Context, routingKey string, correlationID string, event []byte) error {
	return p.PublishCh.PublishWithContext(
		ctx,
		"",         // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: correlationID,
			Body:          event,
			DeliveryMode:  amqp.Persistent,
		},
	)
}

func (p *rpcPBehaviorClientMessageProcessor) logError(err error, errMsg string, msg []byte) {
	if p.FeaturePrintEventOnError {
		p.Logger.Debug().Err(err).Str("event", string(msg)).Msg(errMsg)
	} else {
		p.Logger.Err(err).Msg(errMsg)
	}
}

func (p *rpcPBehaviorClientMessageProcessor) getErrRpcEvent(err error) []byte {
	msg, _ := p.getRpcEvent(rpc.AxeResultEvent{Error: &rpc.Error{Error: err}})
	return msg
}

func (p *rpcPBehaviorClientMessageProcessor) getRpcEvent(event rpc.AxeResultEvent) ([]byte, error) {
	msg, err := p.Encoder.Encode(event)
	if err != nil {
		p.Logger.Err(err).Msg("cannot encode event")
		return nil, err
	}

	return msg, nil
}

func (p *rpcPBehaviorClientMessageProcessor) updateEntity(ctx context.Context, entity *types.Entity, alarm types.Alarm, changeType types.AlarmChangeType) {
	switch changeType {
	case types.AlarmChangeTypeCreateAndPbhEnter, types.AlarmChangeTypePbhEnter,
		types.AlarmChangeTypePbhLeave, types.AlarmChangeTypePbhLeaveAndEnter:
		entity.PbehaviorInfo = alarm.Value.PbehaviorInfo
		err := p.EntityAdapter.UpdatePbehaviorInfo(ctx, entity.ID, entity.PbehaviorInfo)
		if err != nil {
			p.Logger.Err(err).Msg("cannot update entity")
		}
	}
}

func (p *rpcPBehaviorClientMessageProcessor) updatePbhLastAlarmDate(ctx context.Context, changeType types.AlarmChangeType, pbehaviorInfo types.PbehaviorInfo) {
	if changeType != types.AlarmChangeTypeCreateAndPbhEnter &&
		changeType != types.AlarmChangeTypePbhEnter &&
		changeType != types.AlarmChangeTypePbhLeaveAndEnter {
		return
	}

	err := p.PbehaviorAdapter.UpdateLastAlarmDate(ctx, pbehaviorInfo.ID, types.CpsTime{Time: time.Now()})
	if err != nil {
		p.Logger.Err(err).Msg("cannot update pbehavior")
	}
}
