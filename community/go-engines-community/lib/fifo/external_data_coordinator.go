package fifo

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type ExternalDataCoordinator interface {
	// Dispatch parks the partially processed event and sends its webhook request so
	// the missing external data can be fetched asynchronously.
	// Processing resumes once the webhook RPC response is consumed and the event is fed back through the processor.
	Dispatch(
		ctx context.Context,
		partialEvent types.Event,
		pr eventfilter.ServiceResult,
		preFilterEvent types.Event,
		eventMetric techmetrics.FifoEventMetric,
		restartCount int,
	) error
}

func NewNullExternalDataCoordinator() ExternalDataCoordinator {
	return &nullExternalDataCoordinator{}
}

type nullExternalDataCoordinator struct{}

func (nullExternalDataCoordinator) Dispatch(_ context.Context, _ types.Event, _ eventfilter.ServiceResult, _ types.Event, _ techmetrics.FifoEventMetric, _ int) error {
	return errors.New("external data is not supported")
}
