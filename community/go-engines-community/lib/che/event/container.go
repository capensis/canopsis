package event

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/contextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type Processor interface {
	Process(ctx context.Context, event *types.Event, partialRes *ProcessorResult) (ProcessorResult, error)
}

type ProcessorResult struct {
	UpdatedEntitiesForEvent    []types.Entity             `json:"uefe,omitzero"`
	UpdatedEntityIdsForMetrics []string                   `json:"uefm,omitzero"`
	EventMetric                techmetrics.CheEventMetric `json:"em,omitzero"`
	ContextGraphReport         contextgraph.Report        `json:"cgr,omitzero"`

	EventFilterResult eventfilter.ServiceResult `json:"efr,omitzero"`
	PreFilterEvent    *types.Event              `json:"pfe,omitzero"`
}

// IsSuspended reports whether event processing was parked to fetch external data asynchronously.
// A suspended event must be resumed once the webhook RPC has answered instead of being forwarded to the next engine.
func (r ProcessorResult) IsSuspended() bool {
	return r.EventFilterResult.ExternalDataRequest != nil
}

// ResumeWithData turns a suspended result into one carrying the fetched external data,
// ready to be fed back into Processor.Process.
func (r ProcessorResult) ResumeWithData(result []any) ProcessorResult {
	r.EventFilterResult = r.EventFilterResult.ResumeWithData(result)
	return r
}

// ResumeWithError turns a suspended result into one reporting that external data could not be fetched.
// errIndex is the index of the external-data reference that failed.
func (r ProcessorResult) ResumeWithError(err error, errIndex int64) ProcessorResult {
	r.EventFilterResult = r.EventFilterResult.ResumeWithError(err, errIndex)
	return r
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
