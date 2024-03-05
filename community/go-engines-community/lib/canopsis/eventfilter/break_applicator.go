package eventfilter

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type breakApplicator struct{}

func NewBreakApplicator() RuleApplicator {
	return &breakApplicator{}
}

func (a *breakApplicator) Apply(_ context.Context, _ ParsedRule, _ *types.Event, _ RegexMatch) (string, bool, error) {
	return OutcomeBreak, false, nil
}
