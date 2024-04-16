package event

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type Processor interface {
	Process(ctx context.Context, event rpc.AxeEvent) (Result, error)
}

type Result struct {
	Alarm       types.Alarm
	Entity      types.Entity
	AlarmChange types.AlarmChange

	Forward              bool
	IsInstructionMatched bool
	IsCountersUpdated    bool
}
