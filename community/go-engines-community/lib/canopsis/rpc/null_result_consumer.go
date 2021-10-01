package rpc

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"github.com/rs/zerolog"
)

type NullResultConsumer struct {
	FeaturePrintEventOnError bool
	Logger                   zerolog.Logger
	ResultChannel            chan<- action.RpcResult
}

func (p *NullResultConsumer) Process(_ context.Context, _ engine.RPCMessage) error {
	return nil
}
