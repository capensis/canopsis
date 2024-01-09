package axe

import (
	"context"
	"errors"
	"fmt"
	"strings"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	libevent "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/axe/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	libentity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice/statecounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type rpcPBehaviorClientMessageProcessor struct {
	DbClient                 mongo.DbClient
	MetricsSender            metrics.Sender
	PublishCh                libamqp.Channel
	RemediationRpc           engine.RPCClient
	EventProcessor           libevent.Processor
	EntityAdapter            libentity.Adapter
	AlarmAdapter             alarm.Adapter
	PbehaviorAdapter         pbehavior.Adapter
	StateCountersService     statecounters.StateCountersService
	Decoder                  encoding.Decoder
	Encoder                  encoding.Encoder
	Logger                   zerolog.Logger
	FeaturePrintEventOnError bool
}

func (p *rpcPBehaviorClientMessageProcessor) Process(ctx context.Context, msg engine.RPCMessage) error {
	data := strings.Split(msg.CorrelationID, "**")
	if len(data) != 2 {
		return fmt.Errorf("RPC PBehavior Client: bad correlation_id: %s", msg.CorrelationID)
	}

	correlationId := data[0]
	routingKey := data[1]

	var replyEvent []byte
	var event rpc.PbehaviorResultEvent
	err := p.Decoder.Decode(msg.Body, &event)
	if err != nil || event.Alarm == nil || event.Entity == nil {
		p.logError(err, "RPC PBehavior Client: invalid event", msg.Body)

		return p.publishResult(ctx, routingKey, correlationId, p.getErrRpcEvent(fmt.Errorf("invalid event")))
	}

	var alarmChangeType types.AlarmChangeType
	if event.PbhEvent.EventType != "" {
		result, err := p.EventProcessor.Process(ctx, rpc.AxeEvent{
			EventType: event.PbhEvent.EventType,
			Alarm:     event.Alarm,
			Entity:    event.Entity,
			Parameters: rpc.AxeParameters{
				PbehaviorInfo: event.PbhEvent.PbehaviorInfo,
				Author:        event.PbhEvent.Author,
				Initiator:     event.PbhEvent.Initiator,
				Output:        event.PbhEvent.Output,
			},
		})
		if err != nil {
			if engine.IsConnectionError(err) {
				return err
			}

			p.logError(err, "RPC PBehavior Client: cannot update alarm", msg.Body)
			return p.publishResult(ctx, routingKey, correlationId, p.getErrRpcEvent(fmt.Errorf("cannot update alarm: %w", err)))
		}

		alarmChangeType = result.AlarmChange.Type
		if result.Alarm.ID != "" {
			event.Alarm = &result.Alarm
		}

		if result.Entity.ID != "" {
			event.Entity = &result.Entity
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
