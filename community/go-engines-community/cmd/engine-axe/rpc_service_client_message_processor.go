package main

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"github.com/rs/zerolog"
)

type rpcServiceClientMessageProcessor struct {
	Logger zerolog.Logger
}

//We are not waiting for any results from engine-service rpc, but we need to read from the queue to keep it clean.
func (p *rpcServiceClientMessageProcessor) Process(msg engine.RPCMessage) error {
	p.Logger.Debug().Str("event", string(msg.Body)).Msg("RPC Service Client: received")

	return nil
}
