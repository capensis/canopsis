package neweventfilter

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type dropApplicator struct {}

func NewDropApplicator() RuleApplicator {
	return &breakApplicator{}
}

func (a *dropApplicator) Apply(_ context.Context, _ Rule, event types.Event, _ pattern.EventRegexMatches) (string, types.Event, error) {
	return OutcomeDrop, event, nil
}
