package che

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/che/event"
)

type ExternalDataCoordinator interface {
	// Dispatch parks the partially processed event and sends its webhook request so
	// the missing external data can be fetched asynchronously.
	// Processing resumes once the webhook RPC response is consumed and the event is fed back through the processor.
	Dispatch(
		ctx context.Context,
		partialEvent types.Event,
		pr event.ProcessorResult,
		restartCount int,
	) error
}

func NewNullExternalDataCoordinator() ExternalDataCoordinator {
	return &nullExternalDataCoordinator{}
}

type nullExternalDataCoordinator struct{}

func (c *nullExternalDataCoordinator) Dispatch(_ context.Context, _ types.Event, _ event.ProcessorResult, _ int) error {
	return errors.New("external data is not supported")
}
