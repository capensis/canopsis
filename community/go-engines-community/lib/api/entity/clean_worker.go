package entity

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	libredis "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

const (
	lockValue          = 1
	lockTickInterval   = time.Minute
	lockExpirationTime = time.Minute + 10*time.Second
)

type DisabledCleaner interface {
	RunCleanerProcess(ctx context.Context, ch <-chan CleanTask)
}

type worker struct {
	redisClient               redis.Cmdable
	dataStorageAdapter        datastorage.Adapter
	dataStorageConfigProvider config.DataStorageConfigProvider
	metricMetaUpdater         metrics.MetaUpdater
	logger                    zerolog.Logger
}

func NewDisabledCleaner(
	redisClient redis.Cmdable,
	adapter datastorage.Adapter,
	dataStorageConfigProvider config.DataStorageConfigProvider,
	metricMetaUpdater metrics.MetaUpdater,
	logger zerolog.Logger,
) DisabledCleaner {
	return &worker{
		redisClient:               redisClient,
		dataStorageAdapter:        adapter,
		dataStorageConfigProvider: dataStorageConfigProvider,
		metricMetaUpdater:         metricMetaUpdater,
		logger:                    logger,
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

			w.processTask(ctx, task)
		}
	}
}

func (w *worker) processTask(ctx context.Context, task CleanTask) {
	res := w.redisClient.SetNX(ctx, libredis.ApiCleanEntitiesLockKey, lockValue, lockExpirationTime)
	if err := res.Err(); err != nil {
		w.logger.Err(err).Msg("cannot set redis lock")
		return
	}
	if !res.Val() {
		return
	}

	defer func() {
		err := w.redisClient.Del(ctx, libredis.ApiCleanEntitiesLockKey).Err()
		if err != nil {
			w.logger.Err(err).Msg("cannot delete redis lock")
			return
		}
	}()

	done := make(chan struct{})
	go func() {
		w.doTask(ctx, task)
		close(done)
	}()

	ticker := time.NewTicker(lockTickInterval)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			err := w.redisClient.SetEx(ctx, libredis.ApiCleanEntitiesLockKey, lockValue, lockExpirationTime).Err()
			if err != nil {
				w.logger.Err(err).Msg("cannot update redis lock")
			}
		}
	}
}

func (w *worker) doTask(ctx context.Context, task CleanTask) {
	dbClient, err := mongo.NewClientWithOptions(ctx, 0, 0, 0,
		w.dataStorageConfigProvider.Get().MongoClientTimeout, w.logger)
	if err != nil {
		w.logger.Err(err).Msg("cannot connect to mongo")
		return
	}

	defer func() {
		err = dbClient.Disconnect(ctx)
		if err != nil {
			w.logger.Err(err).Msg("cannot disconnect from mongo")
		}
	}()

	arch := NewArchiver(dbClient)
	switch task.Type {
	case CleanTaskTypeArchiveDisabled:
		archived, err := arch.ArchiveDisabledEntities(ctx, task.ArchiveWithDependencies)
		if err != nil {
			w.logger.Err(err).Msg("failed to archive disabled entities")
			return
		}

		err = w.dataStorageAdapter.UpdateHistoryEntityDisabled(ctx, datastorage.HistoryWithCount{
			Time:     datetime.NewCpsTime(),
			Archived: archived,
		})
		if err != nil {
			w.logger.Err(err).Msg("failed to update entity history")
			return
		}

		if archived > 0 {
			w.metricMetaUpdater.UpdateAll(ctx)
		}

		w.logger.Info().Int64("entities_number", archived).Str("user", task.UserID).Msg("disabled entities have been archived")
	case CleanTaskTypeArchiveUnlinked:
		if task.ArchiveBefore == nil {
			return
		}

		before := task.ArchiveBefore.SubFrom(datetime.NewCpsTime())
		totalArchived, err := arch.ArchiveUnlinkedResources(ctx, before)
		if err != nil {
			w.logger.Err(err).Msg("failed to archive unlinked resources")
		}

		archivedComponents, err := arch.ArchiveUnlinkedComponents(ctx, before)
		if err != nil {
			w.logger.Err(err).Msg("failed to archive unlinked components")
		}

		totalArchived += archivedComponents
		archivedConnectors, err := arch.ArchiveUnlinkedConnectors(ctx, before)
		if err != nil {
			w.logger.Err(err).Msg("failed to archive unlinked connectors")
		}

		totalArchived += archivedConnectors
		err = w.dataStorageAdapter.UpdateHistoryEntityUnlinked(ctx, datastorage.HistoryWithCount{
			Time:     datetime.NewCpsTime(),
			Archived: totalArchived,
		})
		if err != nil {
			w.logger.Err(err).Msg("failed to update entity history")
			return
		}

		if totalArchived > 0 {
			w.metricMetaUpdater.UpdateAll(ctx)
		}

		w.logger.Info().Int64("entities_number", totalArchived).Str("user", task.UserID).Msg("unlinked entities have been archived")
	case CleanTaskTypeCleanArchived:
		deleted, err := arch.DeleteArchivedEntities(ctx)
		if err != nil {
			w.logger.Err(err).Msg("failed to delete archived entities")
			return
		}

		err = w.dataStorageAdapter.UpdateHistoryEntityCleaned(ctx, datastorage.HistoryWithCount{
			Time:    datetime.NewCpsTime(),
			Deleted: deleted,
		})
		if err != nil {
			w.logger.Err(err).Msg("failed to update entity history")
			return
		}

		w.logger.Info().Int64("alarm_number", deleted).Str("user", task.UserID).Msg("archived entities have been deleted")
	default:
		w.logger.Error().Msgf("unknown task type %d", task.Type)
	}
}
