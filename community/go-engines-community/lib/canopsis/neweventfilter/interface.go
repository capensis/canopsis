package neweventfilter

//go:generate mockgen -destination=../../../mocks/lib/canopsis/neweventfilter/neweventfilter.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/neweventfilter RuleApplicator,RuleAdapter,RuleApplicatorContainer

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

// outcome constant values
const (
	OutcomePass = iota
	OutcomeDrop
	OutcomeBreak
)

type RuleApplicator interface {
	// Apply eventfilter rule, the first return value(int) should be one of the outcome constant values
	Apply(context.Context, Rule, types.Event, pattern.EventRegexMatches) (int, types.Event, error)
}

type RuleAdapter interface {
	GetAll(context.Context) ([]Rule, error)
}

type EventFilterService interface {
	ProcessEvent(context.Context, types.Event) (types.Event, error)
	LoadRules(context.Context) error
}

type RuleApplicatorContainer interface {
	Get(string) (RuleApplicator, bool)
	Set(string, RuleApplicator)
}
