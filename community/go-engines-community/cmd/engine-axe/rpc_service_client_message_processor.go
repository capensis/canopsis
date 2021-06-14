package main

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/engine"
	"github.com/rs/zerolog"
)

type rpcServiceClientMessageProcessor struct {
	Logger zerolog.Logger
}

//We are not waiting for any results from engine-service rpc, but we need to read from the queue to keep it clean.
func (p *rpcServiceClientMessageProcessor) Process(msg engine.RPCMessage) error {
	p.Logger.Debug().Str("RPC Service Client: event", string(msg.Body)).Msg("received")

	return nil
}
