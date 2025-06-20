package event

import (
	"context"
	"errors"
	"fmt"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters/calculator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewResolveDeletedProcessor(
	dbClient mongo.DbClient,
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator,
	componentCountersCalculator calculator.ComponentCountersCalculator,
	eventsSender entitycounters.EventsSender,
	metaAlarmPostProcessor MetaAlarmPostProcessor,
	metaAlarmStatesService correlation.MetaAlarmStateService,
	metricsSender metrics.Sender,
	remediationRpcClient engine.RPCClient,
	eventGenerator event.Generator,
	amqpPublisher libamqp.Publisher,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Processor {
	return &resolveDeletedProcessor{
		dbClient:                        dbClient,
		alarmCollection:                 dbClient.Collection(mongo.AlarmMongoCollection),
		entityCollection:                dbClient.Collection(mongo.EntityMongoCollection),
		resolvedAlarmCollection:         dbClient.Collection(mongo.ResolvedAlarmMongoCollection),
		pbehaviorCollection:             dbClient.Collection(mongo.PbehaviorMongoCollection),
		metaAlarmRuleCollection:         dbClient.Collection(mongo.MetaAlarmRulesMongoCollection),
		closeDelayJobCollection:         dbClient.Collection(mongo.CloseDelayJobCollection),
		entityServiceCountersCalculator: entityServiceCountersCalculator,
		componentCountersCalculator:     componentCountersCalculator,
		eventsSender:                    eventsSender,
		metaAlarmPostProcessor:          metaAlarmPostProcessor,
		metaAlarmStatesService:          metaAlarmStatesService,
		metricsSender:                   metricsSender,
		remediationRpcClient:            remediationRpcClient,
		eventGenerator:                  eventGenerator,
		amqpPublisher:                   amqpPublisher,
		encoder:                         encoder,
		logger:                          logger,
	}
}

type resolveDeletedProcessor struct {
	dbClient                        mongo.DbClient
	alarmCollection                 mongo.DbCollection
	entityCollection                mongo.DbCollection
	resolvedAlarmCollection         mongo.DbCollection
	pbehaviorCollection             mongo.DbCollection
	metaAlarmRuleCollection         mongo.DbCollection
	closeDelayJobCollection         mongo.DbCollection
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator
	componentCountersCalculator     calculator.ComponentCountersCalculator
	eventsSender                    entitycounters.EventsSender
	metaAlarmPostProcessor          MetaAlarmPostProcessor
	metaAlarmStatesService          correlation.MetaAlarmStateService
	metricsSender                   metrics.Sender
	remediationRpcClient            engine.RPCClient
	eventGenerator                  event.Generator
	amqpPublisher                   libamqp.Publisher
	encoder                         encoding.Encoder
	logger                          zerolog.Logger
}

func (p *resolveDeletedProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || event.Entity.SoftDeleted == nil || event.Entity.ResolveDeletedEventProcessed != nil {
		return result, nil
	}

	now := datetime.NewCpsTime()
	match := getOpenAlarmMatch(event)
	var componentStateChanged bool
	var newComponentState int
	var updatedServiceStates map[string]entitycounters.UpdatedServicesInfo
	notAckedMetricType := ""
	err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}
		componentStateChanged = false
		newComponentState = 0
		updatedServiceStates = nil
		notAckedMetricType = ""

		beforeAlarm, err := updateAlarmToResolve(ctx, p.alarmCollection, match, event.Parameters)
		if err != nil {
			return err
		}

		entityUpdate := bson.M{}
		if beforeAlarm.ID != "" {
			_, err = p.closeDelayJobCollection.DeleteOne(ctx, bson.M{"_id": beforeAlarm.ID})
			if err != nil {
				return fmt.Errorf("failed to delete close_delay job on resolve_deleted event: %w", err)
			}

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

		entity := types.Entity{}
		entityUpdate["$set"] = bson.M{"resolve_deleted_event_processed": now}
		err = p.entityCollection.FindOneAndUpdate(ctx, bson.M{"_id": event.Entity.ID}, entityUpdate,
			options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&entity)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		result.Entity = entity
		result.IsCountersUpdated, updatedServiceStates, componentStateChanged, newComponentState, err = processComponentAndServiceCounters(
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

	if result.AlarmChange.Type == types.AlarmChangeTypeResolve {
		go postProcessResolve(
			context.Background(),
			event,
			result,
			updatedServiceStates,
			componentStateChanged,
			newComponentState,
			notAckedMetricType,
			p.eventsSender,
			p.metaAlarmPostProcessor,
			p.metricsSender,
			p.remediationRpcClient,
			p.pbehaviorCollection,
			p.entityCollection,
			p.eventGenerator,
			p.amqpPublisher,
			p.encoder,
			p.logger,
		)
	} else {
		go p.postProcess(context.WithoutCancel(ctx), &event, updatedServiceStates, componentStateChanged, newComponentState)
	}

	return result, nil
}

func (p *resolveDeletedProcessor) postProcess(
	ctx context.Context,
	event *rpc.AxeEvent,
	updatedServiceStates map[string]entitycounters.UpdatedServicesInfo,
	componentStateChanged bool,
	newComponentState int,
) {
	for servID, servInfo := range updatedServiceStates {
		err := p.eventsSender.UpdateServiceState(ctx, servID, servInfo)
		if err != nil {
			p.logger.Err(err).Msg("failed to update service state")
		}
	}

	if componentStateChanged {
		err := p.eventsSender.UpdateComponentState(ctx, event.Entity.Component, newComponentState)
		if err != nil {
			p.logger.Err(err).Msg("failed to update component state")
		}
	}
}
