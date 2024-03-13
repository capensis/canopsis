package action

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"github.com/rs/zerolog"
)

type axeRpcClientMessageProcessor struct {
	FeaturePrintEventOnError bool
	Decoder                  encoding.Decoder
	Logger                   zerolog.Logger
	ResultChannel            chan<- action.RpcResult
}

func (p *axeRpcClientMessageProcessor) Process(_ context.Context, msg engine.RPCMessage) error {
	var event rpc.AxeResultEvent
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
		WebhookHeader:   event.WebhookHeader,
		WebhookResponse: event.WebhookResponse,
		Error:           eventErr,
	}

	return nil
}

func (p *axeRpcClientMessageProcessor) logError(err error, errMsg string, msg []byte) {
	if p.FeaturePrintEventOnError {
		p.Logger.Err(err).Str("event", string(msg)).Msg(errMsg)
	} else {
		p.Logger.Err(err).Msg(errMsg)
	}
}
