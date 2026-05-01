package logger

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	libredis "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/bsm/redislock"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/sync/errgroup"
)

const (
	batchSize     = 500
	flushInterval = 100 * time.Millisecond

	redisLockTTLDuration     = 30 * time.Second
	redisLockRefreshDuration = redisLockTTLDuration / 2

	redisLockAcquireRetries  = 5
	redisLockAcquireInterval = redisLockRefreshDuration / redisLockAcquireRetries
)

const (
	logTypeCreate = iota
	logTypeUpdate
	logTypeDelete
)

type ActionLogger interface {
	Watch(ctx context.Context) error
}

type logger struct {
	dbClient       mongo.DbClient
	pgPoolProvider postgres.PoolProvider
	zLog           zerolog.Logger

	redisLockClient libredis.LockClient

	collectionValueTypeMap map[string]string
	watchedCollections     []string

	maxRetries   int
	retryTimeout time.Duration
}

func NewActionLogger(
	dbClient mongo.DbClient,
	redisLockClient libredis.LockClient,
	pgPoolProvider postgres.PoolProvider,
	zLog zerolog.Logger,
	retryCount int,
	retryTimeout time.Duration,
) ActionLogger {
	/*
	   NOTE: Each collection listed below MUST have the MongoDB option
	   `changeStreamPreAndPostImages` enabled.
	   Without this flag, the `fullDocument` and `fullDocumentBeforeChange`
	   fields will be nil for update and delete events, breaking the action logger.
	   Enable it for an existing collection (e.g., in a database migration) with:
	       db.runCommand({collMod: "collection_name", changeStreamPreAndPostImages: {enabled: true}})
	   Or when creating a collection:
	       db.createCollection("collection_name", {changeStreamPreAndPostImages: {enabled: true}})
	*/
	collectionValueTypeMap := map[string]string{
		mongo.AlarmTagCollection:                 ValueTypeAlarmTag,
		mongo.ColorThemeCollection:               ValueTypeColorTheme,
		mongo.LinkRuleMongoCollection:            ValueTypeLinkRule,
		mongo.PatternMongoCollection:             ValueTypePattern,
		mongo.PlaylistMongoCollection:            ValueTypePlayList,
		mongo.RoleCollection:                     ValueTypeRole,
		mongo.ScenarioCollection:                 ValueTypeScenario,
		mongo.UserCollection:                     ValueTypeUser,
		mongo.StateSettingsMongoCollection:       ValueTypeStateSetting,
		mongo.ViewMongoCollection:                ValueTypeView,
		mongo.ViewGroupMongoCollection:           ValueTypeViewGroup,
		mongo.ViewTabMongoCollection:             ValueTypeViewTab,
		mongo.WidgetMongoCollection:              ValueTypeWidget,
		mongo.WidgetFiltersMongoCollection:       ValueTypeWidgetFilter,
		mongo.EntityMongoCollection:              ValueTypeEntity,
		mongo.EntityCategoryMongoCollection:      ValueTypeEntityCategory,
		mongo.BroadcastMessageCollection:         ValueTypeBroadcastMessage,
		mongo.EventFilterRuleCollection:          ValueTypeEventFilter,
		mongo.FlappingRuleMongoCollection:        ValueTypeFlappingRule,
		mongo.IdleRuleMongoCollection:            ValueTypeIdleRule,
		mongo.IconCollection:                     ValueTypeIcon,
		mongo.PbehaviorMongoCollection:           ValueTypePbehavior,
		mongo.PbehaviorTypeMongoCollection:       ValueTypePbehaviorType,
		mongo.PbehaviorReasonMongoCollection:     ValueTypePbehaviorReason,
		mongo.PbehaviorExceptionMongoCollection:  ValueTypePbehaviorException,
		mongo.ResolveRuleMongoCollection:         ValueTypeResolveRule,
		mongo.WidgetTemplateMongoCollection:      ValueTypeWidgetTemplate,
		mongo.DeclareTicketRuleCollection:        ValueTypeDeclareTicketRule,
		mongo.DynamicInfosRulesMongoCollection:   ValueTypeDynamicInfo,
		mongo.JobMongoCollection:                 ValueTypeJob,
		mongo.JobConfigMongoCollection:           ValueTypeJobConfig,
		mongo.KpiFilterMongoCollection:           ValueTypeKpiFilter,
		mongo.MapMongoCollection:                 ValueTypeMap,
		mongo.MetaAlarmRulesMongoCollection:      ValueTypeMetaalarmRule,
		mongo.SnmpRulesCollection:                ValueTypeSnmpRule,
		mongo.InstructionMongoCollection:         ValueTypeInstruction, // TODO: removed {mongo.EventRecordsMongoCollection: ValueTypeEventRecord}; when EventRecordsMongoCollection needs to be logged then its records should be limited by {"m": nil} condition
		mongo.ExternalDataTableCollection:        ValueTypeExternalData,
		mongo.EntityInfosPropertyCollection:      ValueTypeEntityInfosProperty,
		mongo.UserNotificationSettingsCollection: ValueTypeUserNotificationSetting,
		mongo.TemplateTestDataCollection:         ValueTypeTplTestData,
		mongo.TemplateTestCollection:             ValueTypeTplTest,
		mongo.WebhookTokenRuleCollection:         ValueTypeWebhookTokenRule,
		mongo.CommentTemplateMongoCollection:     ValueTypeCommentTemplate,
		mongo.LLMConfigCollection:                ValueTypeLLMConfig,
	}

	watchedCollections := make([]string, 0, len(collectionValueTypeMap))
	for k := range collectionValueTypeMap {
		watchedCollections = append(watchedCollections, k)
	}

	return &logger{
		dbClient:       dbClient,
		pgPoolProvider: pgPoolProvider,
		zLog:           zLog,

		redisLockClient: redisLockClient,

		collectionValueTypeMap: collectionValueTypeMap,
		watchedCollections:     watchedCollections,

		maxRetries:   retryCount,
		retryTimeout: retryTimeout,
	}
}

func (l *logger) obtainLock(ctx context.Context) (libredis.Lock, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, nil
		default:
		}

		lock, err := l.redisLockClient.Obtain(ctx, libredis.ApiActionLogWatchLockKey, redisLockTTLDuration, &redislock.Options{
			RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(redisLockAcquireInterval), redisLockAcquireRetries),
		})
		if err != nil {
			if errors.Is(err, redislock.ErrNotObtained) {
				l.zLog.Debug().Msg("action logger redis lock is not obtained, retry")
				continue
			}

			return nil, fmt.Errorf("cannot obtain lock: %w", err)
		}

		l.zLog.Debug().Msg("action logger redis lock is obtained")

		return lock, nil
	}
}

func (l *logger) startLockRefresher(ctx context.Context, lock libredis.Lock) chan struct{} {
	exitChan := make(chan struct{})

	go func() {
		defer close(exitChan)

		ticker := time.NewTicker(redisLockRefreshDuration)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := lock.Refresh(ctx, redisLockTTLDuration, &redislock.Options{})
				if err != nil {
					l.zLog.Err(err).Msg("failed to refresh lock")
					return
				}

				l.zLog.Debug().Msg("action logger redis lock is refreshed")
			}
		}
	}()

	return exitChan
}

func (l *logger) Watch(ctx context.Context) (err error) {
	ctx, cancel := context.WithCancel(ctx)

	var lock libredis.Lock

	defer func() {
		cancel()

		if lock == nil {
			return
		}

		lockErr := lock.Release(context.WithoutCancel(ctx))
		if lockErr != nil && !errors.Is(lockErr, redislock.ErrLockNotHeld) {
			l.zLog.Err(lockErr).Msg("failed to release lock")
		}

		l.zLog.Debug().Msg("action logger redis lock is released")
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		lock, err = l.obtainLock(ctx)
		if err != nil {
			l.zLog.Warn().Err(err).Msgf("failed to obtain lock for action log watcher, waiting for next attempt")
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(redisLockAcquireInterval):
				l.zLog.Debug().Msg("action logger: retrying to obtain lock")
				continue
			}
		}

		exitChan := l.startLockRefresher(ctx, lock)

		var retryTimeout time.Duration

		for attempt := 0; attempt <= l.maxRetries; attempt++ {
			g, gCtx := errgroup.WithContext(ctx)

			g.Go(func() error {
				select {
				case <-gCtx.Done():
					return nil
				case <-exitChan:
					return libredis.ErrFailedToRefreshLock
				}
			})

			eventChan, err := l.runWatcher(gCtx, g)
			if err == nil {
				// if err = nil, means that stream is created, so we drop counter and timeout to default values.
				attempt = 0
				retryTimeout = l.retryTimeout

				g.Go(func() error {
					return l.runCollector(gCtx, eventChan)
				})

				err = g.Wait()
			}

			if err != nil && !mongo.IsConnectionError(err) {
				if errors.Is(err, libredis.ErrFailedToRefreshLock) {
					// refresh is failed, so the lock is no longer belong to us for whatever reason,
					// do not retry watcher again, break from mongo retries cycle and try to obtain the lock again.
					break
				}

				return err
			}

			if attempt == l.maxRetries || retryTimeout == 0 {
				return fmt.Errorf("action log failed to watch db after %d retries: %w", attempt, err)
			}

			select {
			case <-ctx.Done():
				return nil
			case <-time.After(retryTimeout):
				l.zLog.Warn().Int("attempt", attempt+1).Int("max_attempts", l.maxRetries).Msg("action log watcher: connection retry")
				retryTimeout *= 2
			}
		}
	}
}

func (l *logger) runWatcher(ctx context.Context, g *errgroup.Group) (<-chan ActionLogEvent, error) {
	stream, err := l.dbClient.Watch(ctx, []bson.M{
		{
			"$match": bson.M{
				"$or": []bson.M{
					{
						"operationType":       mongo.ChangeStreamTypeInsert,
						"fullDocument.author": bson.M{"$exists": true},
					},
					{
						"operationType": mongo.ChangeStreamTypeDelete,
					},
					{
						"operationType": mongo.ChangeStreamTypeUpdate,
						"updateDescription.updatedFields.updated": bson.M{"$exists": true},
					},
				},
				"ns.coll":       bson.M{"$in": l.watchedCollections},
				"operationType": bson.M{"$in": []string{mongo.ChangeStreamTypeInsert, mongo.ChangeStreamTypeUpdate, mongo.ChangeStreamTypeDelete}},
			},
		},
		{
			"$project": bson.M{
				"document_id":        "$documentKey._id",
				"collection":         "$ns.coll",
				"operation_type":     "$operationType",
				"document":           "$fullDocument",
				"document_before":    "$fullDocumentBeforeChange",
				"update_description": "$updateDescription",
				"cluster_time":       "$clusterTime",
			},
		},
	}, options.ChangeStream().
		SetFullDocument(options.WhenAvailable).
		SetFullDocumentBeforeChange(options.WhenAvailable),
	)
	if err != nil {
		return nil, err
	}

	eventChan := make(chan ActionLogEvent)

	g.Go(func() error {
		defer func() {
			err := stream.Close(ctx)
			if err != nil {
				l.zLog.Err(err).Msg("failed to close change stream")
			}

			close(eventChan)
		}()

		for stream.Next(ctx) {
			var event ActionLogEvent

			err := stream.Decode(&event)
			if err != nil {
				l.zLog.Err(err).Msg("failed to decode change stream event")
				continue
			}

			eventChan <- event
		}

		return stream.Err()
	})

	return eventChan, nil
}

func (l *logger) resolveEvent(event ActionLogEvent) ActionLog {
	valueType := l.collectionValueTypeMap[event.Collection]

	// The special case for entity services, since they are in the same collection with entities.
	if valueType == ValueTypeEntity {
		var rawType any
		var ok bool

		if event.Document == nil {
			rawType, ok = event.DocumentBefore["type"]
		} else {
			rawType, ok = event.Document["type"]
		}

		if ok {
			strType, ok := rawType.(string)
			if ok && strType == types.EntityTypeService {
				valueType = ValueTypeEntityService
			}
		}
	}

	return ActionLog{
		OperationType:     event.OperationType,
		ValueType:         valueType,
		ValueID:           event.DocumentID,
		Timestamp:         event.ClusterTime,
		CurDocument:       event.Document,
		PrevDocument:      event.DocumentBefore,
		UpdateDescription: event.UpdateDescription,
	}
}

func (l *logger) runCollector(ctx context.Context, eventChan <-chan ActionLogEvent) error {
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	batch := make([]ActionLog, 0, batchSize)

	flush := func() {
		if len(batch) == 0 {
			return
		}

		if err := l.flushBatch(ctx, batch); err != nil {
			l.zLog.Err(err).Msg("failed to flush action log batch")
		}

		batch = batch[:0]
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case event, ok := <-eventChan:
			if !ok {
				flush()
				return nil
			}

			batch = append(batch, l.resolveEvent(event))

			if len(batch) >= batchSize {
				flush()
				ticker.Reset(flushInterval)
			}
		case <-ticker.C:
			flush()
		}
	}
}

func (l *logger) flushBatch(ctx context.Context, batch []ActionLog) error {
	pgPool, err := l.pgPoolProvider.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get pgPool: %w", err)
	}

	var inserts, updates, deletes []ActionLog

	for _, log := range batch {
		switch log.OperationType {
		case mongo.ChangeStreamTypeInsert:
			inserts = append(inserts, log)
		case mongo.ChangeStreamTypeUpdate:
			updates = append(updates, log)
		case mongo.ChangeStreamTypeDelete:
			deletes = append(deletes, log)
		}
	}

	return pgPool.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		if err := l.flushInserts(ctx, tx, inserts); err != nil {
			return err
		}

		if err := l.flushUpdates(ctx, tx, updates); err != nil {
			return err
		}

		return l.flushDeletes(ctx, tx, deletes)
	})
}

func (l *logger) flushInserts(ctx context.Context, tx pgx.Tx, logs []ActionLog) error {
	if len(logs) == 0 {
		return nil
	}

	rows := make([][]any, len(logs))
	for i, log := range logs {
		rows[i] = []any{logTypeCreate, log.ValueType, log.ValueID, log.GetCurAuthor(), log.Timestamp, log.CurDocument}
	}

	_, err := tx.CopyFrom(ctx, pgx.Identifier{"action_log"},
		[]string{"type", "value_type", "value_id", "author", "time", "data"},
		pgx.CopyFromRows(rows),
	)

	return err
}

func (l *logger) flushUpdates(ctx context.Context, tx pgx.Tx, logs []ActionLog) error {
	if len(logs) == 0 {
		return nil
	}

	logByKey := make(map[string]ActionLog, len(logs))
	valueTypes := make([]string, len(logs))
	valueIDs := make([]string, len(logs))

	for i, log := range logs {
		valueTypes[i] = log.ValueType
		valueIDs[i] = log.ValueID
		logByKey[log.ValueType+"\x00"+log.ValueID] = log
	}

	missingPairs, err := l.queryMissingCreateLogs(ctx, tx, valueTypes, valueIDs, logByKey)
	if err != nil {
		return err
	}

	if len(missingPairs) > 0 {
		rows := make([][]any, len(missingPairs))
		for i, log := range missingPairs {
			author := log.GetCurAuthor()
			rows[i] = []any{
				logTypeCreate, log.ValueType, log.ValueID,
				cmp.Or(log.GetPrevAuthor(), author),
				cmp.Or(log.GetPrevCreated(), log.Timestamp),
				log.PrevDocument,
			}
		}

		if _, err := tx.CopyFrom(ctx, pgx.Identifier{"action_log"},
			[]string{"type", "value_type", "value_id", "author", "time", "data"},
			pgx.CopyFromRows(rows),
		); err != nil {
			return fmt.Errorf("failed to insert reconstructed create logs for update: %w", err)
		}
	}

	rows := make([][]any, len(logs))
	for i, log := range logs {
		rows[i] = []any{logTypeUpdate, log.ValueType, log.ValueID, log.GetCurAuthor(), log.Timestamp, log.UpdateDescription}
	}

	_, err = tx.CopyFrom(ctx, pgx.Identifier{"action_log"},
		[]string{"type", "value_type", "value_id", "author", "time", "data"},
		pgx.CopyFromRows(rows),
	)

	return err
}

func (l *logger) flushDeletes(ctx context.Context, tx pgx.Tx, logs []ActionLog) error {
	if len(logs) == 0 {
		return nil
	}

	logByKey := make(map[string]ActionLog, len(logs))
	valueTypes := make([]string, len(logs))
	valueIDs := make([]string, len(logs))

	for i, log := range logs {
		valueTypes[i] = log.ValueType
		valueIDs[i] = log.ValueID
		logByKey[log.ValueType+"\x00"+log.ValueID] = log
	}

	missingPairs, err := l.queryMissingCreateLogs(ctx, tx, valueTypes, valueIDs, logByKey)
	if err != nil {
		return err
	}

	if len(missingPairs) > 0 {
		rows := make([][]any, len(missingPairs))
		for i, log := range missingPairs {
			author := log.GetPrevAuthor()
			created := log.GetPrevCreated()
			if created.IsZero() {
				created = log.Timestamp
			}

			rows[i] = []any{logTypeCreate, log.ValueType, log.ValueID, author, created, log.PrevDocument}
		}

		if _, err := tx.CopyFrom(ctx, pgx.Identifier{"action_log"},
			[]string{"type", "value_type", "value_id", "author", "time", "data"},
			pgx.CopyFromRows(rows),
		); err != nil {
			return fmt.Errorf("failed to insert reconstructed create logs for delete: %w", err)
		}
	}

	rows := make([][]any, len(logs))
	for i, log := range logs {
		rows[i] = []any{logTypeDelete, log.ValueType, log.ValueID, log.GetPrevAuthor(), log.Timestamp, nil}
	}

	_, err = tx.CopyFrom(ctx, pgx.Identifier{"action_log"},
		[]string{"type", "value_type", "value_id", "author", "time", "data"},
		pgx.CopyFromRows(rows),
	)

	return err
}

// queryMissingCreateLogs returns the subset of the input logs for which no logTypeCreate
// row exists yet in action_log within the current transaction.
//
// SELECT query has verified with EXPLAIN ANALYZE on 80k rows
func (l *logger) queryMissingCreateLogs(
	ctx context.Context,
	tx pgx.Tx,
	valueTypes, valueIDs []string,
	logByKey map[string]ActionLog,
) ([]ActionLog, error) {
	const missingPairsQuery = `
		SELECT pairs.value_type, pairs.value_id
		FROM UNNEST($1::text[], $2::text[]) AS pairs(value_type, value_id)
		WHERE NOT EXISTS (
			SELECT 1 FROM action_log al
			WHERE al.type = $3 AND al.value_type = pairs.value_type AND al.value_id = pairs.value_id
		)`

	rows, err := tx.Query(ctx, missingPairsQuery, valueTypes, valueIDs, logTypeCreate)
	if err != nil {
		return nil, fmt.Errorf("failed to query missing create logs: %w", err)
	}
	defer rows.Close()

	var missing []ActionLog

	for rows.Next() {
		var vt, vid string

		if err := rows.Scan(&vt, &vid); err != nil {
			return nil, fmt.Errorf("failed to scan missing pair: %w", err)
		}

		if log, ok := logByKey[vt+"\x00"+vid]; ok {
			missing = append(missing, log)
		}
	}

	return missing, rows.Err()
}
