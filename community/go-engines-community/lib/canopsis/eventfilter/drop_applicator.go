package eventfilter

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type dropApplicator struct{}

func NewDropApplicator() RuleApplicator {
	return &dropApplicator{}
}

func (a *dropApplicator) Apply(_ context.Context, _ ParsedRule, _ *types.Event, _ RegexMatch) (string, bool, map[string]int64, error) {
	return OutcomeDrop, false, nil, nil
}
