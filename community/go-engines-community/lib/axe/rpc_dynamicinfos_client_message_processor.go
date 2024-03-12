package axe

import (
	"context"
	"errors"
	"fmt"
	"strings"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type rpcDynamicInfosClientMessageProcessor struct {
	FeaturePrintEventOnError bool
	PublishCh                libamqp.Channel
	Decoder                  encoding.Decoder
	Encoder                  encoding.Encoder
	Logger                   zerolog.Logger
}

func (p *rpcDynamicInfosClientMessageProcessor) Process(ctx context.Context, msg engine.RPCMessage) error {
	data := strings.Split(msg.CorrelationID, "**")
	if len(data) != 2 {
		return fmt.Errorf("bad correlation_id: %s", msg.CorrelationID)
	}

	correlationId := data[0]
	routingKey := data[1]
	var event rpc.DynamicInfosResultEvent
	err := p.Decoder.Decode(msg.Body, &event)
	if err != nil || event.Alarm == nil {
		p.logError(err, "invalid event", msg.Body)

		return p.publishResult(ctx, routingKey, correlationId, p.getErrRpcEvent(fmt.Errorf("invalid event")))
	}

	replyEvent, err := p.getRpcEvent(rpc.AxeResultEvent{
		Alarm:           event.Alarm,
		AlarmChangeType: event.AlarmChangeType,
		Error:           event.Error,
	})
	if err != nil {
		p.logError(err, "failed to create rpc result event", msg.Body)

		replyEvent = p.getErrRpcEvent(errors.New("failed to create rpc result event"))
	}

	err = p.publishResult(ctx, routingKey, correlationId, replyEvent)
	if err != nil {
		p.logError(err, "cannot sent message result back to sender", msg.Body)

		return err
	}

	return nil
}

func (p *rpcDynamicInfosClientMessageProcessor) publishResult(ctx context.Context, routingKey string, correlationID string, event []byte) error {
	return p.PublishCh.PublishWithContext(
		ctx,
		"",
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:   canopsis.JsonContentType,
			CorrelationId: correlationID,
			Body:          event,
			DeliveryMode:  amqp.Persistent,
		},
	)
}

func (p *rpcDynamicInfosClientMessageProcessor) logError(err error, errMsg string, msg []byte) {
	if p.FeaturePrintEventOnError {
		p.Logger.Debug().Err(err).Str("event", string(msg)).Msg(errMsg)
	} else {
		p.Logger.Err(err).Msg(errMsg)
	}
}

func (p *rpcDynamicInfosClientMessageProcessor) getErrRpcEvent(err error) []byte {
	msg, _ := p.getRpcEvent(rpc.AxeResultEvent{Error: &rpc.Error{Error: err}})
	return msg
}

func (p *rpcDynamicInfosClientMessageProcessor) getRpcEvent(event rpc.AxeResultEvent) ([]byte, error) {
	msg, err := p.Encoder.Encode(event)
	if err != nil {
		p.Logger.Err(err).Msg("cannot encode event")
		return nil, err
	}

	return msg, nil
}
