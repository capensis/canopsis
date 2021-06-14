package metaalarm

import (
	"context"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
)

type RulesAdapter interface {
	// Get read meta-alarm rules from db
	Get() ([]Rule, error)

	Save(Rule) error
	GetManualRule() (Rule, error)
	GetRule(id string) (Rule, error)
}

type RulesService interface {
	// LoadRules loads the meta-alarm rules from the database, and adds them to the ruleService.
	LoadRules(ctx context.Context) error

	ProcessEvent(context.Context, *types.Event) ([]types.Event, error)
}

type RuleApplicator interface {
	Apply(context.Context, *types.Event, Rule) ([]types.Event, error)
}

type RuleApplicatorContainer interface {
	Get(RuleType) (RuleApplicator, bool)
	Set(RuleType, RuleApplicator)
}

type RuleEntityCounter interface {
	CountTotalEntitiesAmount(context.Context, Rule) error
	GetTotalEntitiesAmount(context.Context, Rule) (int, error)
}

type ValueGroupEntityCounter interface {
	CountTotalEntitiesAmount(ctx context.Context, rule Rule) error
	CountTotalEntitiesAmountForValuePaths(ctx context.Context, rule Rule, valuePathsMap map[string]string) error
	GetTotalEntitiesAmount(ctx context.Context, ruleId string, valueGroup string) (int, error)
}