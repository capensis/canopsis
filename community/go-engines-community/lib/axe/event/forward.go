package event

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
)

func NewForwardProcessor() Processor {
	return &forwardProcessor{}
}

type forwardProcessor struct {
}

func (p *forwardProcessor) Process(_ context.Context, _ rpc.AxeEvent) (Result, error) {
	result := Result{
		Forward: true,
	}

	return result, nil
}
