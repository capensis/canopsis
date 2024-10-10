package event

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters/calculator"
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
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator,
	componentCountersCalculator calculator.ComponentCountersCalculator,
	eventsSender entitycounters.EventsSender,
	metaAlarmPostProcessor MetaAlarmPostProcessor,
	metaAlarmStatesService correlation.MetaAlarmStateService,
	metricsSender metrics.Sender,
	remediationRpcClient engine.RPCClient,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Processor {
	return &recomputeEntityServiceProcessor{
		dbClient:                        dbClient,
		alarmCollection:                 dbClient.Collection(mongo.AlarmMongoCollection),
		entityCollection:                dbClient.Collection(mongo.EntityMongoCollection),
		resolvedAlarmCollection:         dbClient.Collection(mongo.ResolvedAlarmMongoCollection),
		pbehaviorCollection:             dbClient.Collection(mongo.PbehaviorMongoCollection),
		metaAlarmRuleCollection:         dbClient.Collection(mongo.MetaAlarmRulesMongoCollection),
		entityServiceCountersCalculator: entityServiceCountersCalculator,
		componentCountersCalculator:     componentCountersCalculator,
		eventsSender:                    eventsSender,
		metaAlarmPostProcessor:          metaAlarmPostProcessor,
		metaAlarmStatesService:          metaAlarmStatesService,
		metricsSender:                   metricsSender,
		remediationRpcClient:            remediationRpcClient,
		encoder:                         encoder,
		logger:                          logger,
	}
}

type recomputeEntityServiceProcessor struct {
	dbClient                        mongo.DbClient
	alarmCollection                 mongo.DbCollection
	entityCollection                mongo.DbCollection
	resolvedAlarmCollection         mongo.DbCollection
	pbehaviorCollection             mongo.DbCollection
	metaAlarmRuleCollection         mongo.DbCollection
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator
	componentCountersCalculator     calculator.ComponentCountersCalculator
	eventsSender                    entitycounters.EventsSender
	metaAlarmPostProcessor          MetaAlarmPostProcessor
	metaAlarmStatesService          correlation.MetaAlarmStateService
	metricsSender                   metrics.Sender
	remediationRpcClient            engine.RPCClient
	encoder                         encoding.Encoder
	logger                          zerolog.Logger
}

func (p *recomputeEntityServiceProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil {
		return result, nil
	}

	if event.Entity.Enabled {
		entity := *event.Entity
		var updatedServiceStates map[string]entitycounters.UpdatedServicesInfo

		err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
			var err error

			updatedServiceStates, err = p.entityServiceCountersCalculator.RecomputeCounters(ctx, &entity)
			return err
		})

		if err != nil {
			return result, err
		}

		for servID, servInfo := range updatedServiceStates {
			err := p.eventsSender.UpdateServiceState(ctx, servID, servInfo)
			if err != nil {
				p.logger.Err(err).Msg("failed to update service state")
			}
		}

		return result, nil
	}

	now := datetime.NewCpsTime()
	match := getOpenAlarmMatch(event)
	var updatedServiceStates map[string]entitycounters.UpdatedServicesInfo
	notAckedMetricType := ""
	err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}
		updatedServiceStates = nil
		notAckedMetricType = ""

		beforeAlarm, err := updateAlarmToResolve(ctx, p.alarmCollection, match, event.Parameters)
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

			err = removeMetaAlarmStateOnResolve(ctx, p.metaAlarmRuleCollection, p.metaAlarmStatesService, result.Alarm)
			if err != nil {
				return err
			}
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

		result.IsCountersUpdated, updatedServiceStates, _, _, err = processComponentAndServiceCounters(
			ctx,
			p.entityServiceCountersCalculator,
			p.componentCountersCalculator,
			&result.Alarm,
			&entity,
			result.AlarmChange,
		)

		return err
	})
	if err != nil {
		return result, err
	}

	go postProcessResolve(
		context.Background(),
		event,
		result,
		updatedServiceStates,
		false,
		0,
		notAckedMetricType,
		p.eventsSender,
		p.metaAlarmPostProcessor,
		p.metricsSender,
		p.remediationRpcClient,
		p.pbehaviorCollection,
		p.encoder,
		p.logger,
	)

	return result, nil
}
