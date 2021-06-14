package view

//go:generate mockgen -destination=../../../mocks/lib/canopsis/view/view.go git.canopsis.net/canopsis/go-engines/lib/canopsis/view Adapter

import "context"

type Adapter interface {
	FindJunitWidgets(ctx context.Context) ([]Widget, error)
	FindJunitWidgetsTestSuiteIDs(ctx context.Context,
		widgetIDs []string) (testSuiteIDs []string, err error)
	AddTestSuitesToJunitWidgets(
		ctx context.Context,
		widgetIDs, testSuiteIDs []string,
	) error
	RemoveTestSuitesFromJunitWidgets(
		ctx context.Context,
		testSuiteIDs []string,
	) error
}
