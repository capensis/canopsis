package eventfilter

//go:generate mockgen -destination=../../../mocks/lib/canopsis/eventfilter/eventfilter.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter RuleApplicator,RuleAdapter,RuleApplicatorContainer,ExternalDataGetter,Service,ActionProcessor

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

// outcome constant values
const (
	OutcomePass  = "pass"
	OutcomeDrop  = "drop"
	OutcomeBreak = "break"
)

type ActionProcessor interface {
	Process(ctx context.Context, action Action, event types.Event, regexMatchWrapper RegexMatchWrapper, externalData map[string]interface{}) (types.Event, error)
}

type RuleApplicator interface {
	// Apply eventfilter rule, the first return value(string) should be one of the outcome constant values
	Apply(context.Context, Rule, types.Event, RegexMatchWrapper) (string, types.Event, error)
}

type RuleAdapter interface {
	GetAll(context.Context) ([]Rule, error)
	GetByTypes(context.Context, []string) ([]Rule, error)
}

type Service interface {
	ProcessEvent(context.Context, types.Event) (types.Event, error)
	LoadRules(context.Context, []string) error
}

type RuleApplicatorContainer interface {
	Get(string) (RuleApplicator, bool)
	Set(string, RuleApplicator)
}

type ExternalDataGetter interface {
	Get(ctx context.Context, parameters ExternalDataParameters, templateParameters TemplateGetter) (interface{}, error)
}

type TemplateGetter interface {
	GetTemplate() interface{}
}
