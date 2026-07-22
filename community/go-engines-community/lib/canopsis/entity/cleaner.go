package entity

import (
	"context"
	"errors"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	libredis "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/bsm/redislock"
	"github.com/rs/zerolog"
)

const (
	CleanTaskTypeArchiveDisabled = iota
	CleanTaskTypeArchiveUnlinked
	CleanTaskTypeCleanArchived
)

type CleanTask struct {
	Type                    int
	ArchiveWithDependencies bool
	ArchiveBefore           datetime.CpsTime
	UserID                  string
	Limit                   int
}

const (
	redisLockTTLDuration     = 30 * time.Second
	redisLockRefreshDuration = redisLockTTLDuration / 2

	redisLockAcquireRetries  = 5
	redisLockAcquireInterval = redisLockRefreshDuration / redisLockAcquireRetries
)

type Cleaner interface {
	datastorage.Cleaner
	RunCleanerProcess(ctx context.Context, ch <-chan CleanTask)
}

type worker struct {
	redisLockClient           libredis.LockClient
	dataStorageAdapter        datastorage.Adapter
	dataStorageConfigProvider config.DataStorageConfigProvider
	metricMetaUpdater         metrics.MetaUpdater
	logger                    zerolog.Logger
}

func NewCleaner(
	redisLockClient libredis.LockClient,
	adapter datastorage.Adapter,
	dataStorageConfigProvider config.DataStorageConfigProvider,
	metricMetaUpdater metrics.MetaUpdater,
	logger zerolog.Logger,
) Cleaner {
	return &worker{
		dataStorageAdapter:        adapter,
		dataStorageConfigProvider: dataStorageConfigProvider,
		metricMetaUpdater:         metricMetaUpdater,
		logger:                    logger,
		redisLockClient:           redisLockClient,
	}
}

func (w *worker) IsEnabled(conf datastorage.Config) bool {
	return datetime.IsDurationEnabledAndValid(conf.Entity.ArchiveAfter)
}

func (w *worker) Clean(ctx context.Context, _ mongo.DbClient, conf datastorage.Config, t datetime.CpsTime, limit int) (datastorage.CleanResult, error) {
	res := datastorage.CleanResult{}
	if !w.IsEnabled(conf) {
		return res, nil
	}

	return w.processTask(ctx, CleanTask{
		Type:          CleanTaskTypeArchiveUnlinked,
		ArchiveBefore: t,
		UserID:        canopsis.DefaultEventAuthor,
		Limit:         limit,
	}, true)
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

			_, err := w.processTask(ctx, task, false)
			if err != nil {
				w.logger.Err(err).Msg("failed to process clean task")
			}
		}
	}
}

func (w *worker) processTask(ctx context.Context, task CleanTask, retryLock bool) (datastorage.CleanResult, error) {
	ctx, cancel := context.WithCancel(ctx)

	var lock libredis.Lock

	defer func() {
		cancel()

		if lock == nil {
			return
		}

		lockErr := lock.Release(context.WithoutCancel(ctx))
		if lockErr != nil && !errors.Is(lockErr, redislock.ErrLockNotHeld) {
			w.logger.Err(lockErr).Msg("entity cleaner: failed to release lock")
		}

		w.logger.Debug().Msg("entity cleaner: redis lock is released")
	}()

	for {
		select {
		case <-ctx.Done():
			return datastorage.CleanResult{}, nil
		default:
		}

		var err error

		lock, err = w.redisLockClient.Obtain(ctx, libredis.ApiCleanEntitiesLockKey, redisLockTTLDuration, &redislock.Options{
			RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(redisLockAcquireInterval), redisLockAcquireRetries),
		})
		if err != nil {
			if errors.Is(err, redislock.ErrNotObtained) {
				if retryLock {
					w.logger.Debug().Msg("entity cleaner: redis lock is not obtained, retry")

					continue
				} else {
					w.logger.Debug().Msg("entity cleaner: redis lock is not obtained, abort")

					return datastorage.CleanResult{}, nil
				}
			}

			return datastorage.CleanResult{}, fmt.Errorf("entity cleaner: cannot obtain lock: %w", err)
		}

		w.logger.Debug().Msg("entity cleaner: redis lock is obtained")

		break
	}

	go func() {
		ticker := time.NewTicker(redisLockRefreshDuration)

		defer ticker.Stop()
		defer cancel()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := lock.Refresh(ctx, redisLockTTLDuration, &redislock.Options{})
				if err != nil {
					w.logger.Err(err).Msg("entity cleaner: failed to refresh lock")
					return
				}

				w.logger.Debug().Msg("entity cleaner: redis lock is refreshed")
			}
		}
	}()

	return w.doTask(ctx, task)
}

func (w *worker) doTask(ctx context.Context, task CleanTask) (datastorage.CleanResult, error) {
	dbClient, err := mongo.NewClient(ctx, mongo.ClientOptions{
		ClientTimeout: w.dataStorageConfigProvider.Get().MongoClientTimeout,
	})
	if err != nil {
		return datastorage.CleanResult{}, fmt.Errorf("cannot connect to mongo: %w", err)
	}

	defer func() {
		err = dbClient.Disconnect(ctx)
		if err != nil {
			w.logger.Err(err).Msg("cannot disconnect from mongo")
		}
	}()

	limit := w.dataStorageConfigProvider.Get().MaxUpdates

	arch := NewArchiver(dbClient)
	switch task.Type {
	case CleanTaskTypeArchiveDisabled:
		archived, err := arch.ArchiveDisabledEntities(ctx, task.ArchiveWithDependencies, limit)
		if err != nil {
			return datastorage.CleanResult{}, fmt.Errorf("failed to archive disabled entities: %w", err)
		}

		err = w.dataStorageAdapter.UpdateHistoryEntityDisabled(ctx, datastorage.HistoryWithCount{
			Time:     datetime.NewCpsTime(),
			Archived: archived,
		})
		if err != nil {
			return datastorage.CleanResult{Archived: archived}, fmt.Errorf("failed to update entity history: %w", err)
		}

		if archived > 0 {
			w.metricMetaUpdater.UpdateAll(ctx)
		}

		w.logger.Info().Int64("entities_number", archived).Str("user", task.UserID).Msg("disabled entities have been archived")

		return datastorage.CleanResult{Archived: archived}, nil
	case CleanTaskTypeArchiveUnlinked:
		runID := utils.NewID()
		taskLogger := w.logger.With().
			Str("task_type", "archive_unlinked").
			Str("run_id", runID).
			Str("user", task.UserID).
			Logger()

		taskLogger.Info().
			Str("archive_before", task.ArchiveBefore.String()).
			Msg("archive_unlinked started")

		archivedEntityIDs := make([]string, 0, canopsis.DefaultBulkSize)
		taskStarted := time.Now()
		archivedIDs, totalArchived, err := arch.ArchiveUnlinkedResources(ctx, task.ArchiveBefore, limit)
		resourcesDuration := time.Since(taskStarted)
		if err != nil {
			taskLogger.Err(err).
				Dur("resources_archive_duration", resourcesDuration).
				Msg("failed to archive unlinked resources")
			return datastorage.CleanResult{}, fmt.Errorf("failed to archive unlinked resources: %w", err)
		}
		archivedEntityIDs = append(archivedEntityIDs, archivedIDs...)
		taskLogger.Info().
			Int64("archived_resources", totalArchived).
			Dur("resources_archive_duration", resourcesDuration).
			Msg("archive_unlinked resources archived")

		started := time.Now()
		archivedIDs, archivedComponents, err := arch.ArchiveUnlinkedComponents(ctx, task.ArchiveBefore, limit)
		componentsDuration := time.Since(started)
		if err != nil {
			taskLogger.Err(err).
				Dur("components_archive_duration", componentsDuration).
				Msg("failed to archive unlinked components")
			return datastorage.CleanResult{Archived: totalArchived}, fmt.Errorf("failed to archive unlinked components: %w", err)
		}
		archivedEntityIDs = append(archivedEntityIDs, archivedIDs...)
		taskLogger.Info().
			Int64("archived_components", archivedComponents).
			Dur("components_archive_duration", componentsDuration).
			Msg("archive_unlinked components archived")

		totalArchived += archivedComponents
		started = time.Now()
		archivedIDs, archivedConnectors, err := arch.ArchiveUnlinkedConnectors(ctx, task.ArchiveBefore, limit)
		connectorsDuration := time.Since(started)
		if err != nil {
			taskLogger.Err(err).
				Dur("connectors_archive_duration", connectorsDuration).
				Msg("failed to archive unlinked connectors")
			return datastorage.CleanResult{Archived: totalArchived}, fmt.Errorf("failed to archive unlinked connectors: %w", err)
		}
		archivedEntityIDs = append(archivedEntityIDs, archivedIDs...)
		taskLogger.Info().
			Int64("archived_connectors", archivedConnectors).
			Dur("connectors_archive_duration", connectorsDuration).
			Msg("archive_unlinked connectors archived")

		totalArchived += archivedConnectors
		started = time.Now()
		err = w.dataStorageAdapter.UpdateHistoryEntityUnlinked(ctx, datastorage.HistoryWithCount{
			Time:     datetime.NewCpsTime(),
			Archived: totalArchived,
		})
		historyDuration := time.Since(started)
		if err != nil {
			taskLogger.Err(err).
				Dur("history_update_duration", historyDuration).
				Msg("failed to update entity history")
			return datastorage.CleanResult{Archived: totalArchived}, fmt.Errorf("failed to update entity history: %w", err)
		}

		deletedMetaRowsQueued := 0
		deleteMetaQueued := false
		var deleteMetaEnqueueDuration time.Duration
		if totalArchived > 0 {
			started = time.Now()
			w.metricMetaUpdater.DeleteById(metrics.ContextWithRunID(ctx, runID), archivedEntityIDs...)
			deletedMetaRowsQueued = len(archivedEntityIDs)
			deleteMetaEnqueueDuration = time.Since(started)
			deleteMetaQueued = deletedMetaRowsQueued > 0
		}

		taskLogger.Info().
			Int64("entities_number", totalArchived).
			Int("archived_entity_ids", len(archivedEntityIDs)).
			Dur("resources_archive_duration", resourcesDuration).
			Dur("components_archive_duration", componentsDuration).
			Dur("connectors_archive_duration", connectorsDuration).
			Dur("history_update_duration", historyDuration).
			Bool("delete_meta_queued", deleteMetaQueued).
			Int("delete_meta_rows_queued", deletedMetaRowsQueued).
			Dur("delete_meta_enqueue_duration", deleteMetaEnqueueDuration).
			Dur("task_duration", time.Since(taskStarted)).
			Msg("unlinked entities have been archived")

		return datastorage.CleanResult{Archived: totalArchived}, nil
	case CleanTaskTypeCleanArchived:
		deleted, err := arch.DeleteArchivedEntities(ctx, limit)
		if err != nil {
			return datastorage.CleanResult{}, fmt.Errorf("failed to delete archived entities: %w", err)
		}

		err = w.dataStorageAdapter.UpdateHistoryEntityCleaned(ctx, datastorage.HistoryWithCount{
			Time:    datetime.NewCpsTime(),
			Deleted: deleted,
		})
		if err != nil {
			return datastorage.CleanResult{Deleted: deleted}, fmt.Errorf("failed to update entity history: %w", err)
		}

		w.logger.Info().Int64("entities_number", deleted).Str("user", task.UserID).Msg("archived entities have been deleted")

		return datastorage.CleanResult{Deleted: deleted}, nil
	default:
		w.logger.Error().Msgf("unknown task type %d", task.Type)
	}

	return datastorage.CleanResult{}, nil
}
