package main

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/action"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/engine"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"github.com/rs/zerolog"
)

type webhookRpcClientMessageProcessor struct {
	FeaturePrintEventOnError bool
	Decoder                  encoding.Decoder
	Logger                   zerolog.Logger
	ResultChannel            chan<- action.RpcResult
}

func (p *webhookRpcClientMessageProcessor) Process(msg engine.RPCMessage) error {
	var event types.RPCWebhookResultEvent
	err := p.Decoder.Decode(msg.Body, &event)
	if err != nil {
		p.logError(err, "cannot decode event", msg.Body)
		return nil
	}

	var eventErr error
	if event.Error != nil {
		eventErr = event.Error.Error
	}

	p.ResultChannel <- action.RpcResult{
		CorrelationID:   msg.CorrelationID,
		Alarm:           event.Alarm,
		AlarmChangeType: event.AlarmChangeType,
		Header:          event.Header,
		Response:        event.Response,
		Error:           eventErr,
	}

	return nil
}

func (p *webhookRpcClientMessageProcessor) logError(err error, errMsg string, msg []byte) {
	if p.FeaturePrintEventOnError {
		p.Logger.Err(err).Str("event", string(msg)).Msg(errMsg)
	} else {
		p.Logger.Err(err).Msg(errMsg)
	}
}
