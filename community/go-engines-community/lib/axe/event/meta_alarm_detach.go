package event

import (
	"context"

	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
)

func NewMetaAlarmDetachProcessor(
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor,
) Processor {
	return &metaAlarmDetachProcessor{
		metaAlarmEventProcessor: metaAlarmEventProcessor,
	}
}

type metaAlarmDetachProcessor struct {
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor
}

func (p *metaAlarmDetachProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil {
		return result, nil
	}

	_, err := p.metaAlarmEventProcessor.DetachChildrenFromMetaAlarm(ctx, event)
	if err != nil {
		return result, err
	}

	result.Forward = false

	return result, nil
}
