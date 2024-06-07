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
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
)

const (
	getLogQuery    = "SELECT id FROM action_log WHERE type = $1 AND value_type = $2 AND value_id = $3"
	insertLogQuery = "INSERT INTO action_log (type, value_type, value_id, author, time, data) VALUES ($1, $2, $3, $4, $5, $6)"

	workerPoolSize = 10
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

	collectionValueTypeMap map[string]string
	watchedCollections     []string

	maxRetries   int
	retryTimeout time.Duration
}

func NewActionLogger(
	dbClient mongo.DbClient,
	pgPoolProvider postgres.PoolProvider,
	zLog zerolog.Logger,
	retryCount int,
	retryTimeout time.Duration,
) ActionLogger {
	collectionValueTypeMap := map[string]string{
		mongo.AlarmTagCollection:                ValueTypeAlarmTag,
		mongo.ColorThemeCollection:              ValueTypeColorTheme,
		mongo.LinkRuleMongoCollection:           ValueTypeLinkRule,
		mongo.PatternMongoCollection:            ValueTypePattern,
		mongo.PlaylistMongoCollection:           ValueTypePlayList,
		mongo.RoleCollection:                    ValueTypeRole,
		mongo.ScenarioMongoCollection:           ValueTypeScenario,
		mongo.UserCollection:                    ValueTypeUser,
		mongo.StateSettingsMongoCollection:      ValueTypeStateSetting,
		mongo.ViewMongoCollection:               ValueTypeView,
		mongo.ViewGroupMongoCollection:          ValueTypeViewGroup,
		mongo.ViewTabMongoCollection:            ValueTypeViewTab,
		mongo.WidgetMongoCollection:             ValueTypeWidget,
		mongo.WidgetFiltersMongoCollection:      ValueTypeWidgetFilter,
		mongo.EntityMongoCollection:             ValueTypeEntity,
		mongo.EntityCategoryMongoCollection:     ValueTypeEntityCategory,
		mongo.BroadcastMessageMongoCollection:   ValueTypeBroadcastMessage,
		mongo.EventFilterRuleCollection:         ValueTypeEventFilter,
		mongo.FlappingRuleMongoCollection:       ValueTypeFlappingRule,
		mongo.IdleRuleMongoCollection:           ValueTypeIdleRule,
		mongo.IconCollection:                    ValueTypeIcon,
		mongo.PbehaviorMongoCollection:          ValueTypePbehavior,
		mongo.PbehaviorTypeMongoCollection:      ValueTypePbehaviorType,
		mongo.PbehaviorReasonMongoCollection:    ValueTypePbehaviorReason,
		mongo.PbehaviorExceptionMongoCollection: ValueTypePbehaviorException,
		mongo.ResolveRuleMongoCollection:        ValueTypeResolveRule,
		mongo.WidgetTemplateMongoCollection:     ValueTypeWidgetTemplate,
		mongo.DeclareTicketRuleMongoCollection:  ValueTypeDeclareTicketRule,
		mongo.DynamicInfosRulesMongoCollection:  ValueTypeDynamicInfo,
		mongo.JobMongoCollection:                ValueTypeJob,
		mongo.JobConfigMongoCollection:          ValueTypeJobConfig,
		mongo.KpiFilterMongoCollection:          ValueTypeKpiFilter,
		mongo.MapMongoCollection:                ValueTypeMap,
		mongo.MetaAlarmRulesMongoCollection:     ValueTypeMetaalarmRule,
		mongo.SnmpRulesCollection:               ValueTypeSnmpRule,
		mongo.InstructionMongoCollection:        ValueTypeInstruction,
	}

	watchedCollections := make([]string, 0, len(collectionValueTypeMap))
	for k := range collectionValueTypeMap {
		watchedCollections = append(watchedCollections, k)
	}

	return &logger{
		dbClient:       dbClient,
		pgPoolProvider: pgPoolProvider,
		zLog:           zLog,

		collectionValueTypeMap: collectionValueTypeMap,
		watchedCollections:     watchedCollections,

		maxRetries:   retryCount,
		retryTimeout: retryTimeout,
	}
}

func (l *logger) log(ctx context.Context, log ActionLog) error {
	pgPool, err := l.pgPoolProvider.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get pgPool: %w", err)
	}

	switch log.OperationType {
	case mongo.ChangeStreamTypeInsert:
		_, err = pgPool.Exec(ctx, insertLogQuery, logTypeCreate, log.ValueType, log.ValueID, log.GetCurAuthor(), log.Timestamp, log.CurDocument)
	case mongo.ChangeStreamTypeUpdate:
		err = pgPool.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
			var objectID int
			author := log.GetCurAuthor()

			err := tx.QueryRow(ctx, getLogQuery, logTypeCreate, log.ValueType, log.ValueID).Scan(&objectID)
			if errors.Is(err, pgx.ErrNoRows) {
				// reconstruct action create log if it doesn't exist
				_, err = tx.Exec(ctx, insertLogQuery, logTypeCreate, log.ValueType, log.ValueID,
					cmp.Or(log.GetPrevAuthor(), author), cmp.Or(log.GetPrevCreated(), log.Timestamp), log.PrevDocument)
			}

			if err != nil {
				return err
			}

			_, err = tx.Exec(ctx, insertLogQuery, logTypeUpdate, log.ValueType, log.ValueID, author, log.Timestamp, log.UpdateDescription)

			return err
		})
	case mongo.ChangeStreamTypeDelete:
		err = pgPool.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
			var objectID int
			author := log.GetPrevAuthor()

			err := tx.QueryRow(ctx, getLogQuery, logTypeCreate, log.ValueType, log.ValueID).Scan(&objectID)
			if errors.Is(err, pgx.ErrNoRows) {
				// reconstruct action create log if it doesn't exist
				_, err = tx.Exec(ctx, insertLogQuery, logTypeCreate, log.ValueType, log.ValueID,
					author, cmp.Or(log.GetPrevCreated(), log.Timestamp), log.PrevDocument)
			}

			if err != nil {
				return err
			}

			_, err = tx.Exec(ctx, insertLogQuery, logTypeDelete, log.ValueType, log.ValueID, author, log.Timestamp, nil)

			return err
		})
	}

	return err
}

func (l *logger) Watch(ctx context.Context) error {
	var retryTimeout time.Duration

	for attempt := 0; attempt <= l.maxRetries; attempt++ {
		g := errgroup.Group{}

		eventChan, err := l.runWatcher(ctx, &g)
		if err == nil {
			// if err = nil, means that stream is created, so we drop counter and timeout to default values.
			attempt = 0
			retryTimeout = l.retryTimeout

			for j := 0; j < workerPoolSize; j++ {
				g.Go(func() error {
					l.runWorker(ctx, eventChan)

					return nil
				})
			}

			err = g.Wait()
		}

		if err != nil && !mongo.IsConnectionError(err) {
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

	return nil
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

			err = stream.Decode(&event)
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

func (l *logger) runWorker(ctx context.Context, eventChan <-chan ActionLogEvent) {
	for event := range eventChan {
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

		err := l.log(ctx, ActionLog{
			OperationType:     event.OperationType,
			ValueType:         valueType,
			ValueID:           event.DocumentID,
			Timestamp:         event.ClusterTime,
			CurDocument:       event.Document,
			PrevDocument:      event.DocumentBefore,
			UpdateDescription: event.UpdateDescription,
		})
		if err != nil {
			l.zLog.Err(err).Str("value_type", valueType).Str("value_id", event.DocumentID).Msgf("error on action log %s", event.OperationType)
		}
	}
}
