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
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
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
	dbClient, err := mongo.NewClientWithOptions(ctx, 0, 0, mongo.DefaultServerSelectionTimeout,
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
		runID := utils.NewID()
		taskLogger := w.logger.With().
			Str("task_type", "archive_unlinked").
			Str("run_id", runID).
			Str("user", task.UserID).
			Logger()
		taskStarted := time.Now()

		if task.ArchiveBefore == nil {
			taskLogger.Warn().Msg("archive_before is not set, skip archive_unlinked")
			return
		}

		taskLogger.Info().
			Str("archive_before", task.ArchiveBefore.String()).
			Msg("archive_unlinked started")

		before := task.ArchiveBefore.SubFrom(datetime.NewCpsTime())
		resourcesStarted := time.Now()
		totalArchived, err := arch.ArchiveUnlinkedResources(ctx, before)
		resourcesDuration := time.Since(resourcesStarted)
		if err != nil {
			taskLogger.Err(err).
				Dur("resources_archive_duration", resourcesDuration).
				Msg("failed to archive unlinked resources")
		} else {
			taskLogger.Info().
				Int64("archived_resources", totalArchived).
				Dur("resources_archive_duration", resourcesDuration).
				Msg("archive_unlinked resources archived")
		}

		componentsStarted := time.Now()
		archivedComponents, err := arch.ArchiveUnlinkedComponents(ctx, before)
		componentsDuration := time.Since(componentsStarted)
		if err != nil {
			taskLogger.Err(err).
				Dur("components_archive_duration", componentsDuration).
				Msg("failed to archive unlinked components")
		} else {
			taskLogger.Info().
				Int64("archived_components", archivedComponents).
				Dur("components_archive_duration", componentsDuration).
				Msg("archive_unlinked components archived")
		}

		totalArchived += archivedComponents
		connectorsStarted := time.Now()
		archivedConnectors, err := arch.ArchiveUnlinkedConnectors(ctx, before)
		connectorsDuration := time.Since(connectorsStarted)
		if err != nil {
			taskLogger.Err(err).
				Dur("connectors_archive_duration", connectorsDuration).
				Msg("failed to archive unlinked connectors")
		} else {
			taskLogger.Info().
				Int64("archived_connectors", archivedConnectors).
				Dur("connectors_archive_duration", connectorsDuration).
				Msg("archive_unlinked connectors archived")
		}

		totalArchived += archivedConnectors
		historyStarted := time.Now()
		err = w.dataStorageAdapter.UpdateHistoryEntityUnlinked(ctx, datastorage.HistoryWithCount{
			Time:     datetime.NewCpsTime(),
			Archived: totalArchived,
		})
		historyDuration := time.Since(historyStarted)
		if err != nil {
			taskLogger.Err(err).
				Dur("history_update_duration", historyDuration).
				Msg("failed to update entity history")
			return
		}

		updateAllQueued := false
		var updateAllEnqueueDuration time.Duration
		if totalArchived > 0 {
			enqueueStarted := time.Now()
			w.metricMetaUpdater.UpdateAll(metrics.ContextWithUpdateRunID(ctx, runID))
			updateAllEnqueueDuration = time.Since(enqueueStarted)
			updateAllQueued = true
		}

		taskLogger.Info().
			Int64("entities_number", totalArchived).
			Dur("resources_archive_duration", resourcesDuration).
			Dur("components_archive_duration", componentsDuration).
			Dur("connectors_archive_duration", connectorsDuration).
			Dur("history_update_duration", historyDuration).
			Bool("update_all_queued", updateAllQueued).
			Dur("update_all_enqueue_duration", updateAllEnqueueDuration).
			Dur("task_duration", time.Since(taskStarted)).
			Msg("unlinked entities have been archived")
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
