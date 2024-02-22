package event

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type metaAlarmProcessor struct{}

func NewMetaAlarmProcessor() Processor {
	return &metaAlarmProcessor{}
}

func (p *metaAlarmProcessor) Process(_ context.Context, event *types.Event) (
	[]types.Entity,
	[]string,
	techmetrics.CheEventMetric,
	error,
) {
	var eventMetric techmetrics.CheEventMetric
	eventMetric.EventType = event.EventType

	return nil, nil, eventMetric, nil
}
