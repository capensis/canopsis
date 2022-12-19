package entity

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
)

type DisabledCleaner interface {
	RunCleanerProcess(ctx context.Context, ch <-chan CleanTask)
}

type worker struct {
	store              Store
	dataStorageAdapter datastorage.Adapter
	metricMetaUpdater  metrics.MetaUpdater
	logger             zerolog.Logger
}

func NewDisabledCleaner(
	store Store,
	adapter datastorage.Adapter,
	metricMetaUpdater metrics.MetaUpdater,
	logger zerolog.Logger,
) DisabledCleaner {
	return &worker{
		store:              store,
		dataStorageAdapter: adapter,
		metricMetaUpdater:  metricMetaUpdater,
		logger:             logger,
	}
}

func (w *worker) RunCleanerProcess(ctx context.Context, ch <-chan CleanTask) {
	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-ch:
			if !ok {
				return
			}

			if *task.Archive {
				archived, err := w.store.ArchiveDisabledEntities(ctx, task.ArchiveDependencies)
				if err != nil {
					w.logger.Err(err).Msg("Failed to archive entities")
					continue
				}

				err = w.dataStorageAdapter.UpdateHistoryEntity(ctx, datastorage.HistoryWithCount{
					Time:     types.CpsTime{Time: time.Now()},
					Archived: archived,
				})
				if err != nil {
					w.logger.Err(err).Msg("Failed to update entity history")
					continue
				}

				if archived > 0 {
					w.metricMetaUpdater.UpdateAll(ctx)
				}

				w.logger.Info().Int64("alarm_number", archived).Str("user", task.UserID).Msg("disabled entities have been archived")
				continue
			}

			deleted, err := w.store.DeleteArchivedEntities(ctx)
			if err != nil {
				w.logger.Err(err).Msg("Failed to delete archived entities")
				continue
			}

			err = w.dataStorageAdapter.UpdateHistoryEntity(ctx, datastorage.HistoryWithCount{
				Time:    types.CpsTime{Time: time.Now()},
				Deleted: deleted,
			})
			if err != nil {
				w.logger.Err(err).Msg("Failed to update entity history")
				continue
			}

			w.logger.Info().Int64("alarm_number", deleted).Str("user", task.UserID).Msg("archived entities have been deleted")
		}
	}
}
