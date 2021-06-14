package main

import (
	"errors"
	"fmt"
	libamqp "git.canopsis.net/canopsis/go-engines/lib/amqp"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/engine"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/operation"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/utils"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
	"strings"
)

type rpcPBehaviorClientMessageProcessor struct {
	FeaturePrintEventOnError bool
	PublishCh                libamqp.Channel
	WatcherRpc               engine.RPCClient
	Executor                 operation.Executor
	Decoder                  encoding.Decoder
	Encoder                  encoding.Encoder
	Logger                   zerolog.Logger
}

func (p *rpcPBehaviorClientMessageProcessor) Process(msg engine.RPCMessage) error {
	data := strings.Split(msg.CorrelationID, "**")
	if len(data) != 2 {
		return fmt.Errorf("RPC PBehavior Client: bad correlation_id: %s", msg.CorrelationID)
	}

	correlationId := data[0]
	routingKey := data[1]

	var replyEvent []byte
	var event types.RPCPBehaviorResultEvent
	err := p.Decoder.Decode(msg.Body, &event)
	if err != nil || event.Alarm == nil {
		p.logError(err, "RPC PBehavior Client: invalid event", msg.Body)

		return p.publishResult(routingKey, correlationId, p.getErrRpcEvent(fmt.Errorf("invalid event")))
	}

	if event.PbhEvent.EventType != "" {
		_, err = p.Executor.Exec(
			types.Operation{
				Type: event.PbhEvent.EventType,
				Parameters: types.OperationPbhParameters{
					PbehaviorInfo: event.PbhEvent.PbehaviorInfo,
					Author:        event.PbhEvent.Author,
					Output:        event.PbhEvent.Output,
				},
			},
			event.Alarm,
			event.PbhEvent.Timestamp,
			"",
			types.InitiatorSystem,
		)
		if err != nil {
			if engine.IsConnectionError(err) {
				return err
			}

			p.logError(err, "RPC PBehavior Client: cannot update alarm", msg.Body)
			return p.publishResult(routingKey, correlationId, p.getErrRpcEvent(fmt.Errorf("cannot update alarm: %v", err)))
		}

		body, err := p.Encoder.Encode(types.RPCWatcherEvent{
			Alarm:  event.Alarm,
			Entity: event.PbhEvent.Entity,
		})
		if err != nil {
			p.logError(err, "RPC PBehavior Client: failed to encode rpc call to watcher", msg.Body)
		} else {
			err = p.WatcherRpc.Call(engine.RPCMessage{
				CorrelationID: utils.NewID(),
				Body:          body,
			})
			if err != nil {
				if engine.IsConnectionError(err) {
					return err
				}

				p.logError(err, "RPC PBehavior Client: failed to send rpc call to watcher", msg.Body)
			}
		}
	}

	replyEvent, err = p.getRpcEvent(types.RPCAxeResultEvent{
		Alarm:           event.Alarm,
		AlarmChangeType: types.AlarmChangeType(event.PbhEvent.EventType),
		Error:           nil,
	})
	if err != nil {
		p.logError(err, "RPC PBehavior Client: failed to create rpc result event", msg.Body)

		replyEvent = p.getErrRpcEvent(errors.New("failed to create rpc result event"))
	}

	err = p.publishResult(routingKey, correlationId, replyEvent)
	if err != nil {
		p.logError(err, "RPC PBehavior Client: cannot sent message result back to sender", msg.Body)

		return err
	}

	return nil
}

func (p *rpcPBehaviorClientMessageProcessor) publishResult(routingKey string, correlationID string, event []byte) error {
	return p.PublishCh.Publish(
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
	msg, _ := p.getRpcEvent(types.RPCAxeResultEvent{Error: &types.RPCError{Error: err}})
	return msg
}

func (p *rpcPBehaviorClientMessageProcessor) getRpcEvent(event types.RPCAxeResultEvent) ([]byte, error) {
	msg, err := p.Encoder.Encode(event)
	if err != nil {
		p.Logger.Err(err).Msg("cannot encode event")
		return nil, err
	}

	return msg, nil
}
