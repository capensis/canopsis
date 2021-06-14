package main

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/engine"
	"github.com/rs/zerolog"
)

type rpcWatcherClientMessageProcessor struct {
	Logger zerolog.Logger
}

//We are not waiting for any results from watcher rpc, but we need to read from the queue to keep it clean.
func (p *rpcWatcherClientMessageProcessor) Process(msg engine.RPCMessage) error {
	p.Logger.Debug().Str("RPC Watcher Client: event", string(msg.Body)).Msg("received")

	return nil
}
