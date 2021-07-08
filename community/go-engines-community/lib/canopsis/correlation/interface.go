package correlation

//go:generate mockgen -destination=../../../mocks/lib/canopsis/correlation/metaalarm.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation RulesAdapter

import (
	"context"
	"errors"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

var ErrNoChildren = errors.New("no children")
var ErrChildAlreadyExist = errors.New("child already exists")

const (
	DefaultMetaAlarmComponent     = "metaalarm"
	DefaultMetaAlarmConnector     = "engine"
	DefaultMetaAlarmConnectorName = "correlation"
	DefaultMetaAlarmEntityPrefix  = "meta-alarm-entity-"
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

	ProcessEvent(context.Context, types.Event) ([]types.Event, error)
}

type RuleApplicator interface {
	Apply(context.Context, types.Event, Rule) ([]types.Event, error)
}

type RuleApplicatorContainer interface {
	Get(RuleType) (RuleApplicator, bool)
	Set(RuleType, RuleApplicator)
}

type RuleEntityCounter interface {
	CountTotalEntitiesAmount(ctx context.Context, rule Rule) error
	GetTotalEntitiesAmount(ctx context.Context, rule Rule) (int, error)
}

type ValueGroupEntityCounter interface {
	CountTotalEntitiesAmount(ctx context.Context, rule Rule) error
	CountTotalEntitiesAmountForValuePaths(ctx context.Context, rule Rule, valuePathsMap map[string]string) error
	GetTotalEntitiesAmount(ctx context.Context, ruleId string, valueGroup string) (int, error)
}
