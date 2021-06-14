package main

import (
	"context"
	"errors"
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/engine"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/watcher"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
)

type rpcServerMessageProcessor struct {
	FeaturePrintEventOnError bool
	WatcherService           watcher.Service
	Decoder                  encoding.Decoder
	Encoder                  encoding.Encoder
	Logger                   zerolog.Logger
}

func (p *rpcServerMessageProcessor) Process(d amqp.Delivery) ([]byte, error) {
	msg := d.Body
	var event types.RPCWatcherEvent
	err := p.Decoder.Decode(msg, &event)
	if err != nil || event.Alarm == nil {
		p.logError(err, "invalid event", msg)

		return p.getErrRpcEvent(errors.New("invalid event")), nil
	}

	err = p.WatcherService.ProcessRpc(context.Background(), types.Event{
		Entity:        event.Entity,
		Alarm:         event.Alarm,
		PbehaviorInfo: event.Alarm.Value.PbehaviorInfo,
	})
	if err != nil {
		if engine.IsConnectionError(err) {
			return nil, err
		}

		p.logError(err, "cannot update watcher", msg)

		return p.getRpcEvent(types.RPCWatcherResultEvent{
			Error: &types.RPCError{Error: fmt.Errorf("cannot update alarm: %v", err)},
		})
	}

	return p.getRpcEvent(types.RPCWatcherResultEvent{})
}

func (p *rpcServerMessageProcessor) logError(err error, errMsg string, msg []byte) {
	if p.FeaturePrintEventOnError {
		p.Logger.Err(err).Str("event", string(msg)).Msg(errMsg)
	} else {
		p.Logger.Err(err).Msg(errMsg)
	}
}

func (p *rpcServerMessageProcessor) getErrRpcEvent(err error) []byte {
	msg, _ := p.getRpcEvent(types.RPCWatcherResultEvent{Error: &types.RPCError{Error: err}})
	return msg
}

func (p *rpcServerMessageProcessor) getRpcEvent(event types.RPCWatcherResultEvent) ([]byte, error) {
	msg, err := p.Encoder.Encode(event)
	if err != nil {
		p.Logger.Err(err).Msg("cannot encode event")
		return nil, nil
	}

	return msg, nil
}
