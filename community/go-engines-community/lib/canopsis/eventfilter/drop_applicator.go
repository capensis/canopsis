package eventfilter

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type dropApplicator struct{}

func NewDropApplicator() RuleApplicator {
	return &dropApplicator{}
}

func (a *dropApplicator) Apply(_ context.Context, _ Rule, event types.Event, _ RegexMatchWrapper) (string, types.Event, error) {
	return OutcomeDrop, event, nil
}
