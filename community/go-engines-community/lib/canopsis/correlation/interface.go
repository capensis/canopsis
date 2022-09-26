package correlation

import "context"

//go:generate mockgen -destination=../../../mocks/lib/canopsis/correlation/metaalarm.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation RulesAdapter

type RulesAdapter interface {
	// Get read meta-alarm rules from db
	Get(ctx context.Context) ([]Rule, error)

	Save(ctx context.Context, r Rule) error
	GetManualRule(ctx context.Context) (Rule, error)
	GetRule(ctx context.Context, id string) (Rule, error)
}
