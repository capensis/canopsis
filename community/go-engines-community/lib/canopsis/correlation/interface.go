package correlation

import "context"

//go:generate go tool go.uber.org/mock/mockgen -destination=../../../mocks/lib/canopsis/correlation/metaalarm.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation RulesAdapter

type RulesAdapter interface {
	Get(ctx context.Context) ([]Rule, error)
	Save(ctx context.Context, r Rule) error
	GetRule(ctx context.Context, id string) (Rule, error)
}
