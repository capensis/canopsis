package event

import "fmt"

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

func (c *mapProcessorContainer) Has(eventType string) bool {
	_, ok := c.processors[eventType]

	return ok
}
