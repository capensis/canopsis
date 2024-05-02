package metrics

//go:generate mockgen -destination=../../../mocks/lib/canopsis/metrics/metrics.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics MetaUpdater

import "context"

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

func NewNullMetaUpdater() MetaUpdater {
	return &nullMetaUpdater{}
}

func (u *nullMetaUpdater) UpdateAll(_ context.Context) {
}

func (u *nullMetaUpdater) UpdateById(_ context.Context, _ ...string) {
}

func (u *nullMetaUpdater) DeleteById(_ context.Context, _ ...string) {
}
