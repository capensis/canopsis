package event

import (
	"context"
	"errors"

	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice/statecounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
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

	if event.Entity.Enabled {
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
				updatedServiceStates, err = p.stateCountersService.UpdateServiceCounters(ctx, *event.Entity, nil, result.AlarmChange)
			} else {
				updatedServiceStates, err = p.stateCountersService.UpdateServiceCounters(ctx, *event.Entity, &result.Alarm, result.AlarmChange)
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

		entity := *event.Entity
		if beforeAlarm.ID != "" {
			if beforeAlarm.NotAckedMetricSendTime != nil {
				notAckedMetricType = beforeAlarm.NotAckedMetricType
			}

			entity, err = updateEntityOfResolvedAlarm(ctx, p.entityCollection, event.Entity.ID)
			if err != nil || entity.ID == "" {
				return err
			}

			alarm, err := copyAlarmToResolvedCollection(ctx, p.alarmCollection, p.resolvedAlarmCollection, beforeAlarm.ID)
			if err != nil || alarm.ID == "" {
				return err
			}

			alarmChange := types.NewAlarmChange()
			alarmChange.Type = types.AlarmChangeTypeResolve
			result.Forward = true
			result.Alarm = alarm
			result.Entity = entity
			result.AlarmChange = alarmChange
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
