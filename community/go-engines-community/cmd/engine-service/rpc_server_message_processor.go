package main

import (
	"context"
	"errors"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type rpcServerMessageProcessor struct {
	FeaturePrintEventOnError bool
	EntityServiceService     entityservice.Service
	Decoder                  encoding.Decoder
	Encoder                  encoding.Encoder
	Logger                   zerolog.Logger
}

func (p *rpcServerMessageProcessor) Process(ctx context.Context, d amqp.Delivery) ([]byte, error) {
	msg := d.Body
	var event types.RPCServiceEvent
	err := p.Decoder.Decode(msg, &event)
	if err != nil || event.Alarm == nil {
		p.logError(err, "invalid event", msg)

		return p.getErrRpcEvent(errors.New("invalid event")), nil
	}

	err = p.EntityServiceService.ProcessRpc(ctx, types.Event{
		Entity:        event.Entity,
		Alarm:         event.Alarm,
		PbehaviorInfo: event.Alarm.Value.PbehaviorInfo,
		AlarmChange:   event.AlarmChange,
	})
	if err != nil {
		if engine.IsConnectionError(err) {
			return nil, err
		}

		p.logError(err, "cannot update entity service", msg)

		return p.getRpcEvent(types.RPCServiceResultEvent{
			Error: &types.RPCError{Error: fmt.Errorf("cannot update alarm: %v", err)},
		})
	}

	return p.getRpcEvent(types.RPCServiceResultEvent{})
}

func (p *rpcServerMessageProcessor) logError(err error, errMsg string, msg []byte) {
	if p.FeaturePrintEventOnError {
		p.Logger.Err(err).Str("event", string(msg)).Msg(errMsg)
	} else {
		p.Logger.Err(err).Msg(errMsg)
	}
}

func (p *rpcServerMessageProcessor) getErrRpcEvent(err error) []byte {
	msg, _ := p.getRpcEvent(types.RPCServiceResultEvent{Error: &types.RPCError{Error: err}})
	return msg
}

func (p *rpcServerMessageProcessor) getRpcEvent(event types.RPCServiceResultEvent) ([]byte, error) {
	msg, err := p.Encoder.Encode(event)
	if err != nil {
		p.Logger.Err(err).Msg("cannot encode event")
		return nil, nil
	}

	return msg, nil
}
