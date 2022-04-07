package axe

import (
	"context"
	"errors"
	"fmt"
	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice/statecounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
	"time"
)

type rpcMessageProcessor struct {
	DbClient                 mongo.DbClient
	MetricsSender            metrics.Sender
	EntityAdapter            entity.Adapter
	AlarmAdapter             libalarm.Adapter
	PbhRpc                   engine.RPCClient
	RemediationRpc           engine.RPCClient
	Executor                 operation.Executor
	MetaAlarmEventProcessor  libalarm.MetaAlarmEventProcessor
	StateCountersService     statecounters.StateCountersService
	Decoder                  encoding.Decoder
	Encoder                  encoding.Encoder
	Logger                   zerolog.Logger
	FeaturePrintEventOnError bool
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
		PreviousPbehaviorTime:           alarm.Value.PbehaviorInfo.Timestamp,
	}

	now := time.Now()

	var updatedServiceStates map[string]statecounters.UpdatedServicesInfo
	firstTimeTran := true

	err = p.DbClient.WithTransaction(ctx, func(tCtx context.Context) error {
		if !firstTimeTran {
			ent, exist := p.EntityAdapter.Get(tCtx, event.Entity.ID)
			if !exist {
				return fmt.Errorf("entity with id = %s is not found after transaction rollback", event.Entity.ID)
			}

			event.Entity = &ent

			al, err := p.AlarmAdapter.GetOpenedAlarmByAlarmId(tCtx, alarm.ID)
			if _, ok := err.(errt.NotFound); ok {
				return fmt.Errorf("alarm with id = %s is not found after transaction rollback", alarm.ID)
			} else if err != nil {
				return err
			}

			alarm = &al
		}

		firstTimeTran = false

		alarmChange.Type, err = p.Executor.Exec(tCtx, op, alarm, event.Entity,
			types.CpsTime{Time: now}, "", "", types.InitiatorSystem)
		if err != nil {
			return err
		}

		updatedServiceStates, err = p.StateCountersService.UpdateServiceCounters(tCtx, *event.Entity, event.Alarm, alarmChange)
		return err
	})

	if err != nil {
		if engine.IsConnectionError(err) {
			return nil, err
		}

		p.logError(err, "RPC Message Processor: cannot update alarm", msg)
		return p.getErrRpcEvent(fmt.Errorf("cannot update alarm: %v", err), alarm), nil
	}

	// services alarms
	go func() {
		for servID, servInfo := range updatedServiceStates {
			err := p.StateCountersService.UpdateServiceState(servID, servInfo)
			if err != nil {
				p.Logger.Err(err).Msg("failed to update service state")
			}
		}
	}()

	// send metrics
	go func() {
		p.MetricsSender.SendEventMetrics(
			context.Background(),
			*alarm,
			*event.Entity,
			alarmChange,
			now,
			types.InitiatorSystem,
			"",
		)
	}()

	if p.RemediationRpc != nil && alarmChange.Type == types.AlarmChangeTypeChangeState {
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

	res := types.RPCAxeResultEvent{
		AlarmChangeType: alarmChange.Type,
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
