package axe

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"github.com/rs/zerolog"
)

type rpcServiceClientMessageProcessor struct {
	Logger zerolog.Logger
}

// We are not waiting for any results from engine-service rpc, but we need to read from the queue to keep it clean.
func (p *rpcServiceClientMessageProcessor) Process(_ context.Context, msg engine.RPCMessage) error {
	p.Logger.Debug().Str("RPCServiceClient_event:", string(msg.Body)).Msg("received")

	return nil
}
