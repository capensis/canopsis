package metrics

//go:generate go tool go.uber.org/mock/mockgen -destination=../../../mocks/lib/canopsis/metrics/metrics.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics MetaUpdater

import "context"

type metaUpdaterContextKey string

const updateRunIDContextKey metaUpdaterContextKey = "update_run_id"

type AsyncMetaUpdater interface {
	MetaUpdater
	Run(ctx context.Context)
}

type MetaUpdater interface {
	UpdateAll(ctx context.Context)
	UpdateById(ctx context.Context, ids ...string)
	DeleteById(ctx context.Context, ids ...string)
}

type nullMetaUpdater struct{}

func ContextWithUpdateRunID(ctx context.Context, runID string) context.Context {
	if runID == "" {
		return ctx
	}

	return context.WithValue(ctx, updateRunIDContextKey, runID)
}

func UpdateRunIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	runID, _ := ctx.Value(updateRunIDContextKey).(string)

	return runID
}

func NewNullMetaUpdater() MetaUpdater {
	return &nullMetaUpdater{}
}

func (u *nullMetaUpdater) UpdateAll(_ context.Context) {
}

func (u *nullMetaUpdater) UpdateById(_ context.Context, _ ...string) {
}

func (u *nullMetaUpdater) DeleteById(_ context.Context, _ ...string) {
}
