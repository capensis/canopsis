package entity

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
	"time"
)

type DisabledCleaner interface {
	RunCleanerProcess(ctx context.Context, ch <-chan CleanTask)
}

type worker struct {
	store              Store
	dataStorageAdapter datastorage.Adapter
	logger             zerolog.Logger
}

func NewDisabledCleaner(store Store, adapter datastorage.Adapter, logger zerolog.Logger) DisabledCleaner {
	return &worker{
		store:              store,
		dataStorageAdapter: adapter,
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

				err = w.dataStorageAdapter.UpdateHistoryEntity(ctx, datastorage.EntityHistory{
					Time:             types.CpsTime{Time: time.Now()},
					EntitiesArchived: archived,
				})
				if err != nil {
					w.logger.Err(err).Msg("Failed to update entity history")
					continue
				}

				w.logger.Info().Int64("alarm number", archived).Str("user", task.UserID).Msg("disabled entities have been archived")
				continue
			}

			deleted, err := w.store.DeleteArchivedEntities(ctx)
			if err != nil {
				w.logger.Err(err).Msg("Failed to delete archived entities")
				continue
			}

			err = w.dataStorageAdapter.UpdateHistoryEntity(ctx, datastorage.EntityHistory{
				Time:            types.CpsTime{Time: time.Now()},
				EntitiesDeleted: deleted,
			})
			if err != nil {
				w.logger.Err(err).Msg("Failed to update entity history")
				continue
			}

			w.logger.Info().Int64("alarm number", deleted).Str("user", task.UserID).Msg("archived entities have been deleted")
		}
	}
}
