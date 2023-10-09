package event

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
)

func NewCombinedProcessor(container ProcessorContainer) Processor {
	return &combinedProcessor{container: container}
}

type combinedProcessor struct {
	container ProcessorContainer
}

func (e *combinedProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	processor, ok := e.container.Get(event.EventType)
	if !ok {
		return Result{}, nil
	}

	return processor.Process(ctx, event)
}
