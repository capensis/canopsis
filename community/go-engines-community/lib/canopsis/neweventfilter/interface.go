package neweventfilter

//go:generate mockgen -destination=../../../mocks/lib/canopsis/neweventfilter/neweventfilter.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/neweventfilter RuleApplicator,RuleAdapter,RuleApplicatorContainer,ExternalDataGetter

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

// outcome constant values
const (
	OutcomePass  = "pass"
	OutcomeDrop  = "drop"
	OutcomeBreak = "break"
)

type ActionProcessor interface {
	Process(action Action, event types.Event, regexMatch pattern.EventRegexMatches, externalData map[string]interface{}) (types.Event, error)
}

type RuleApplicator interface {
	// Apply eventfilter rule, the first return value(string) should be one of the outcome constant values
	Apply(context.Context, Rule, types.Event, pattern.EventRegexMatches) (string, types.Event, error)
}

type RuleAdapter interface {
	GetAll(context.Context) ([]Rule, error)
	GetByType(context.Context, string) ([]Rule, error)
}

type EventFilterService interface {
	ProcessEvent(context.Context, types.Event) (types.Event, error)
	LoadRules(context.Context) error
}

type RuleApplicatorContainer interface {
	Get(string) (RuleApplicator, bool)
	Set(string, RuleApplicator)
}

type ExternalDataGetter interface {
	Get(ctx context.Context, parameters ExternalDataParameters, templateParameters TemplateParameters) (interface{}, error)
}
