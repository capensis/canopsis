package event

import (
	"context"
	"errors"

	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
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

func NewEntityToggledProcessor(
	dbClient mongo.DbClient,
	stateCountersService statecounters.StateCountersService,
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor,
	metricsSender metrics.Sender,
	remediationRpcClient engine.RPCClient,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Processor {
	return &entityToggledProcessor{
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

type entityToggledProcessor struct {
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

func (p *entityToggledProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil {
		return result, nil
	}

	entity := *event.Entity
	if entity.Enabled {
		var updatedServiceStates map[string]statecounters.UpdatedServicesInfo
		err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
			updatedServiceStates = nil

			alarm := types.Alarm{}
			err := p.alarmCollection.FindOne(ctx, getOpenAlarmMatch(event)).Decode(&alarm)
			if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
				return err
			}

			alarmChange := types.NewAlarmChange()
			alarmChange.Type = types.AlarmChangeTypeEnabled
			result.Forward = true
			result.Alarm = alarm
			result.AlarmChange = alarmChange
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

		go p.postProcess(context.Background(), updatedServiceStates)

		return result, nil
	}

	match := getOpenAlarmMatch(event)
	update := getResolveAlarmUpdate(datetime.NewCpsTime())
	var updatedServiceStates map[string]statecounters.UpdatedServicesInfo
	notAckedMetricType := ""

	err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}
		updatedServiceStates = nil
		notAckedMetricType = ""

		beforeAlarm := types.Alarm{}
		opts := options.FindOneAndUpdate().
			SetReturnDocument(options.Before).
			SetProjection(bson.M{
				"not_acked_metric_type":      1,
				"not_acked_metric_send_time": 1,
			})
		err := p.alarmCollection.FindOneAndUpdate(ctx, match, update, opts).Decode(&beforeAlarm)
		if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
			return err
		}

		if beforeAlarm.ID != "" {
			if beforeAlarm.NotAckedMetricSendTime != nil {
				notAckedMetricType = beforeAlarm.NotAckedMetricType
			}

			alarm := types.Alarm{}
			err = p.alarmCollection.FindOne(ctx, bson.M{"_id": beforeAlarm.ID}).Decode(&alarm)
			if err != nil {
				if errors.Is(err, mongodriver.ErrNoDocuments) {
					return nil
				}
				return err
			}

			entity = types.Entity{}
			entityUpdate := getResolveEntityUpdate()
			err = p.entityCollection.FindOneAndUpdate(ctx, bson.M{"_id": event.Entity.ID}, entityUpdate,
				options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&entity)
			if err != nil {
				if errors.Is(err, mongodriver.ErrNoDocuments) {
					return nil
				}

				return err
			}

			alarmChange := types.NewAlarmChange()
			alarmChange.Type = types.AlarmChangeTypeResolve
			result.Forward = true
			result.Alarm = alarm
			result.Entity = entity
			result.AlarmChange = alarmChange

			_, err = p.resolvedAlarmCollection.UpdateOne(
				ctx,
				bson.M{"_id": alarm.ID},
				bson.M{"$set": alarm},
				options.Update().SetUpsert(true),
			)
			if err != nil {
				return err
			}
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

func (p *entityToggledProcessor) postProcess(
	ctx context.Context,
	updatedServiceStates map[string]statecounters.UpdatedServicesInfo,
) {
	for servID, servInfo := range updatedServiceStates {
		err := p.stateCountersService.UpdateServiceState(ctx, servID, servInfo)
		if err != nil {
			p.logger.Err(err).Msg("failed to update service state")
		}
	}
}
