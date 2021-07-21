package neweventfilter

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type RuleApplicator interface {
	Apply(_ context.Context, rule Rule, params ApplicatorParameters) (types.Event, error)
}

type RulesAdapter interface {
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
