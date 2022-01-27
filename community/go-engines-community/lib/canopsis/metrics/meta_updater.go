package metrics

import "context"

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
