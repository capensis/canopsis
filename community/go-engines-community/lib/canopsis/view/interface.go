package view

//go:generate go tool go.uber.org/mock/mockgen -destination=../../../mocks/lib/canopsis/view/view.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view Adapter

import "context"

type Adapter interface {
	FindJunitWidgets(ctx context.Context) ([]Widget, error)
	AddTestSuitesToJunitWidgets(
		ctx context.Context,
		widgetIDs, testSuiteIDs []string,
	) error
}
