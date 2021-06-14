package metaalarm

import (
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
	LoadRules() error

	ProcessEvent(*types.Event) ([]types.Event, error)
}

type RuleApplicator interface {
	Apply(*types.Event, Rule) ([]types.Event, error)
}

type RuleApplicatorContainer interface {
	Get(RuleType) (RuleApplicator, bool)
	Set(RuleType, RuleApplicator)
}

type RuleEntityCounter interface {
	CountTotalEntitiesAmount(Rule) error
	GetTotalEntitiesAmount(Rule) (int, error)
}

type ValueGroupEntityCounter interface {
	CountTotalEntitiesAmount(rule Rule) error
	CountTotalEntitiesAmountForValuePaths(rule Rule, valuePathsMap map[string]string) error
	GetTotalEntitiesAmount(ruleId string, valueGroup string) (int, error)
}