package event

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type Processor interface {
	Process(ctx context.Context, event *types.Event) (
		[]types.Entity,
		[]string,
		techmetrics.CheEventMetric,
		error,
	)
}

type ProcessorContainer interface {
	Get(eventType string) (Processor, bool)
	Set(eventType string, p Processor)
	Has(eventType string) bool
}

func NewProcessorContainer() ProcessorContainer {
	return &mapProcessorContainer{
		processors: make(map[string]Processor),
	}
}

type mapProcessorContainer struct {
	processors map[string]Processor
}

func (c *mapProcessorContainer) Get(eventType string) (Processor, bool) {
	p, ok := c.processors[eventType]

	return p, ok
}

func (c *mapProcessorContainer) Set(eventType string, processor Processor) {
	if c.Has(eventType) {
		panic(fmt.Errorf("event processor %q already exists", eventType))
	}

	c.processors[eventType] = processor
}

func (c *mapProcessorContainer) Has(sourceType string) bool {
	_, ok := c.processors[sourceType]

	return ok
}
