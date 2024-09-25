package event

import (
	"context"
	"errors"
	"fmt"

	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice/statecounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewRecomputeEntityServiceProcessor(
	dbClient mongo.DbClient,
	stateCountersService statecounters.StateCountersService,
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor,
	metricsSender metrics.Sender,
	remediationRpcClient engine.RPCClient,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Processor {
	return &recomputeEntityServiceProcessor{
		dbClient:                dbClient,
		alarmCollection:         dbClient.Collection(mongo.AlarmMongoCollection),
		entityCollection:        dbClient.Collection(mongo.EntityMongoCollection),
		resolvedAlarmCollection: dbClient.Collection(mongo.ResolvedAlarmMongoCollection),
		stateCountersService:    stateCountersService,
		metaAlarmEventProcessor: metaAlarmEventProcessor,
		metricsSender:           metricsSender,
		remediationRpcClient:    remediationRpcClient,
		encoder:                 encoder,
		logger:                  logger,
	}
}

type recomputeEntityServiceProcessor struct {
	dbClient                mongo.DbClient
	alarmCollection         mongo.DbCollection
	entityCollection        mongo.DbCollection
	resolvedAlarmCollection mongo.DbCollection
	stateCountersService    statecounters.StateCountersService
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor
	metricsSender           metrics.Sender
	remediationRpcClient    engine.RPCClient
	encoder                 encoding.Encoder
	logger                  zerolog.Logger
}

func (p *recomputeEntityServiceProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil {
		return result, nil
	}

	if event.Entity.Enabled {
		entity := *event.Entity
		var updatedServiceStates map[string]statecounters.UpdatedServicesInfo

		err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
			var err error

			updatedServiceStates, err = p.stateCountersService.RecomputeEntityServiceCounters(ctx, entity)
			return err
		})

		if err != nil {
			return result, err
		}

		for servID, servInfo := range updatedServiceStates {
			err := p.stateCountersService.UpdateServiceState(ctx, servID, servInfo)
			if err != nil {
				return result, fmt.Errorf("failed to update service state: %w", err)
			}
		}

		return result, nil
	}

	now := types.NewCpsTime()
	match := getOpenAlarmMatch(event)
	var updatedServiceStates map[string]statecounters.UpdatedServicesInfo
	notAckedMetricType := ""
	err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}
		updatedServiceStates = nil
		notAckedMetricType = ""

		beforeAlarm, err := updateAlarmToResolve(ctx, p.alarmCollection, match)
		if err != nil {
			return err
		}

		entityUpdate := bson.M{}
		if beforeAlarm.ID != "" {
			if beforeAlarm.NotAckedMetricSendTime != nil {
				notAckedMetricType = beforeAlarm.NotAckedMetricType
			}

			alarm, err := copyAlarmToResolvedCollection(ctx, p.alarmCollection, p.resolvedAlarmCollection, beforeAlarm.ID)
			if err != nil || alarm.ID == "" {
				return err
			}

			entityUpdate = getResolveEntityUpdate()
			alarmChange := types.NewAlarmChange()
			alarmChange.Type = types.AlarmChangeTypeResolve
			result.Forward = true
			result.Alarm = alarm
			result.AlarmChange = alarmChange
		}

		if event.Entity.SoftDeleted != nil && event.Entity.ResolveDeletedEventProcessed == nil {
			entityUpdate["$set"] = bson.M{"resolve_deleted_event_processed": now}
		}

		entity := *event.Entity
		if len(entityUpdate) > 0 {
			entity = types.Entity{}
			err = p.entityCollection.FindOneAndUpdate(ctx, bson.M{"_id": event.Entity.ID}, entityUpdate,
				options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&entity)
			if err != nil {
				if errors.Is(err, mongodriver.ErrNoDocuments) {
					return nil
				}

				return err
			}

			result.Entity = entity
		}

		if result.Alarm.ID == "" {
			updatedServiceStates, err = p.stateCountersService.UpdateServiceCounters(ctx, entity, nil, result.AlarmChange)
		} else {
			updatedServiceStates, err = p.stateCountersService.UpdateServiceCounters(ctx, entity, &result.Alarm, result.AlarmChange)
		}

		return err
	})
	if err != nil {
		return result, err
	}

	go postProcessResolve(context.Background(), event, result, updatedServiceStates, notAckedMetricType, p.stateCountersService, p.metaAlarmEventProcessor, p.metricsSender, p.remediationRpcClient, p.encoder, p.logger)

	return result, nil
}
