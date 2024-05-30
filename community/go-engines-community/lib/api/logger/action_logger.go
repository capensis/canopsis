package logger

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const workerPoolSize = 10

type ActionLogger interface {
	LogCreate(ctx context.Context, log ActionCreateLog) error
	LogUpdate(ctx context.Context, log ActionUpdateLog) error
	LogDelete(ctx context.Context, log ActionDeleteLog) error
	Watch(ctx context.Context) error
}

type logger struct {
	dbClient       mongo.DbClient
	pgPoolProvider postgres.PoolProvider
	zLog           zerolog.Logger

	collectionValueTypeMap map[string]string
	watchedCollections     []string

	getObjectQuery          string
	insertObjectQuery       string
	insertObjectChangeQuery string
	updateObjectDeleteQuery string
}

func NewActionLogger(dbClient mongo.DbClient, pgPoolProvider postgres.PoolProvider, zLog zerolog.Logger) ActionLogger {
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

	watchedCollections := make([]string, len(collectionValueTypeMap))
	for k := range collectionValueTypeMap {
		watchedCollections = append(watchedCollections, k)
	}

	return &logger{
		dbClient:       dbClient,
		pgPoolProvider: pgPoolProvider,
		zLog:           zLog,

		collectionValueTypeMap: collectionValueTypeMap,
		watchedCollections:     watchedCollections,

		getObjectQuery: `
			SELECT id FROM action_log_object WHERE value_type = $1 AND value_id = $2
		`,
		insertObjectQuery: `
			INSERT INTO action_log_object (value_type, value_id, initial_value, created, created_by, deleted, deleted_by)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (value_type, value_id)
			DO UPDATE SET initial_value = $3, created = $4, created_by = $5, deleted = $6, deleted_by = $7
			RETURNING id
		`,
		insertObjectChangeQuery: `
			INSERT INTO action_log_object_changes(object_id, author, time, update_description)
			VALUES ($1, $2, $3, $4)
		`,
		updateObjectDeleteQuery: `
			UPDATE action_log_object
			SET deleted = $1, deleted_by = $2 WHERE id = $3
		`,
	}
}

func (l *logger) LogCreate(ctx context.Context, log ActionCreateLog) error {
	pgPool, err := l.pgPoolProvider.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get pgPool: %w", err)
	}

	_, err = pgPool.Exec(ctx, l.insertObjectQuery,
		log.ValueType, log.ValueID, log.InitialValue, log.Timestamp, log.Author, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to insert log: %w", err)
	}

	return nil
}

func (l *logger) LogUpdate(ctx context.Context, log ActionUpdateLog) error {
	pgPool, err := l.pgPoolProvider.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get pgPool: %w", err)
	}

	return pgPool.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		var objectID int

		err := tx.QueryRow(ctx, l.getObjectQuery, log.ValueType, log.ValueID).Scan(&objectID)
		if errors.Is(err, pgx.ErrNoRows) {
			var created time.Time
			var createdBy string

			rawCreated, ok := log.PrevValue["created"]
			if ok {
				if intCreated, ok := rawCreated.(int64); ok {
					created = time.Unix(intCreated, 0)
				}
			} else {
				created = log.Timestamp
			}

			rawCreatedBy, ok := log.PrevValue["author"]
			if ok {
				if strCreatedBy, ok := rawCreatedBy.(string); ok {
					createdBy = strCreatedBy
				}
			} else {
				createdBy = log.Author
			}

			// reconstruct action log object if it doesn't exist
			err = tx.QueryRow(ctx, l.insertObjectQuery,
				log.ValueType, log.ValueID, log.PrevValue, created, createdBy, nil, nil).Scan(&objectID)
		}

		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, l.insertObjectChangeQuery,
			objectID, log.Author, log.Timestamp, log.UpdateDescription)

		return err
	})
}

func (l *logger) LogDelete(ctx context.Context, log ActionDeleteLog) error {
	pgPool, err := l.pgPoolProvider.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get pgPool: %w", err)
	}

	err = pgPool.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		var objectID int

		err := tx.QueryRow(ctx, l.getObjectQuery, log.ValueType, log.ValueID).Scan(&objectID)
		if errors.Is(err, pgx.ErrNoRows) {
			var created time.Time

			rawCreated, ok := log.PrevValue["created"]
			if ok {
				if intCreated, ok := rawCreated.(int64); ok {
					created = time.Unix(intCreated, 0)
				}
			} else {
				created = log.Timestamp
			}

			// reconstruct action log object if it doesn't exist
			// return because the delete data will be inserted
			return tx.QueryRow(ctx, l.insertObjectQuery,
				log.ValueType, log.ValueID, log.PrevValue, created, log.Author, log.Timestamp, log.Author).Scan(&objectID)
		}

		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, l.updateObjectDeleteQuery, log.Timestamp, log.Author, objectID)

		return err
	})

	return err
}

func (l *logger) Watch(ctx context.Context) error {
	eventChan, err := l.runWatcher(ctx)
	if err != nil {
		return fmt.Errorf("failed to run action log change stream watcher: %w", err)
	}

	wg := sync.WaitGroup{}

	for i := 0; i < workerPoolSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			l.runWorker(ctx, eventChan)
		}()
	}

	wg.Wait()

	return nil
}

func (l *logger) runWatcher(ctx context.Context) (<-chan ActionLogEvent, error) {
	stream, err := l.dbClient.Watch(ctx, []bson.M{
		{
			"$match": bson.M{
				"$or": []bson.M{
					{
						"operationType":       "insert",
						"fullDocument.author": bson.M{"$exists": true},
					},
					{
						"operationType": "delete",
					},
					{
						"operationType": "update",
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

	go func() {
		defer func() {
			err := stream.Close(ctx)
			if err != nil {
				l.zLog.Err(err).Msg("failed to close change stream")
			}

			close(eventChan)
		}()

		for stream.Next(ctx) {
			var change struct {
				DocumentID        string         `bson:"document_id"`
				Collection        string         `bson:"collection"`
				OperationType     string         `bson:"operation_type"`
				Document          map[string]any `bson:"document"`
				DocumentBefore    map[string]any `bson:"document_before"`
				UpdateDescription map[string]any `bson:"update_description"`
				ClusterTime       time.Time      `bson:"cluster_time"`
			}

			err = stream.Decode(&change)
			if err != nil {
				l.zLog.Err(err).Msg("failed to decode change stream event")
				continue
			}

			eventChan <- change
		}
	}()

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

		switch event.OperationType {
		case mongo.ChangeStreamTypeInsert:
			log := ActionCreateLog{
				ValueType:    valueType,
				ValueID:      event.DocumentID,
				InitialValue: event.Document,
				Timestamp:    event.ClusterTime,
			}

			rawAuthor, ok := event.Document["author"]
			if ok {
				strAuthor, ok := rawAuthor.(string)
				if ok {
					log.Author = strAuthor
				}
			}

			err := l.LogCreate(ctx, log)
			if err != nil {
				l.zLog.Err(err).Str("value_type", valueType).Str("value_id", log.ValueID).Msg("error on action log create")
			}
		case mongo.ChangeStreamTypeUpdate:
			log := ActionUpdateLog{
				ValueType:         valueType,
				ValueID:           event.DocumentID,
				PrevValue:         event.DocumentBefore,
				UpdateDescription: event.UpdateDescription,
				Timestamp:         event.ClusterTime,
			}

			rawAuthor, ok := event.Document["author"]
			if ok {
				strAuthor, ok := rawAuthor.(string)
				if ok {
					log.Author = strAuthor
				}
			}

			err := l.LogUpdate(ctx, log)
			if err != nil {
				l.zLog.Err(err).Str("value_type", valueType).Str("value_id", log.ValueID).Msg("error on action log update")
			}
		case mongo.ChangeStreamTypeDelete:
			log := ActionDeleteLog{
				ValueType: valueType,
				ValueID:   event.DocumentID,
				PrevValue: event.DocumentBefore,
				Timestamp: event.ClusterTime,
			}

			rawAuthor, ok := event.DocumentBefore["author"]
			if ok {
				strAuthor, ok := rawAuthor.(string)
				if ok {
					log.Author = strAuthor
				}
			}

			err := l.LogDelete(ctx, log)
			if err != nil {
				l.zLog.Err(err).Str("value_type", valueType).Str("value_id", log.ValueID).Msg("error on action log delete")
			}
		}
	}
}
