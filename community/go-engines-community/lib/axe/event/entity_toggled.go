package event

import (
	"context"
	"errors"

	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters/calculator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

func NewEntityToggledProcessor(
	dbClient mongo.DbClient,
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator,
	componentCountersCalculator calculator.ComponentCountersCalculator,
	eventsSender entitycounters.EventsSender,
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor,
	metaAlarmStatesService correlation.MetaAlarmStateService,
	metricsSender metrics.Sender,
	remediationRpcClient engine.RPCClient,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Processor {
	return &entityToggledProcessor{
		dbClient:                        dbClient,
		alarmCollection:                 dbClient.Collection(mongo.AlarmMongoCollection),
		entityCollection:                dbClient.Collection(mongo.EntityMongoCollection),
		resolvedAlarmCollection:         dbClient.Collection(mongo.ResolvedAlarmMongoCollection),
		pbehaviorCollection:             dbClient.Collection(mongo.PbehaviorMongoCollection),
		metaAlarmRuleCollection:         dbClient.Collection(mongo.MetaAlarmRulesMongoCollection),
		entityServiceCountersCalculator: entityServiceCountersCalculator,
		componentCountersCalculator:     componentCountersCalculator,
		eventsSender:                    eventsSender,
		metaAlarmEventProcessor:         metaAlarmEventProcessor,
		metaAlarmStatesService:          metaAlarmStatesService,
		metricsSender:                   metricsSender,
		remediationRpcClient:            remediationRpcClient,
		encoder:                         encoder,
		logger:                          logger,
	}
}

type entityToggledProcessor struct {
	dbClient                        mongo.DbClient
	alarmCollection                 mongo.DbCollection
	entityCollection                mongo.DbCollection
	resolvedAlarmCollection         mongo.DbCollection
	pbehaviorCollection             mongo.DbCollection
	metaAlarmRuleCollection         mongo.DbCollection
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator
	componentCountersCalculator     calculator.ComponentCountersCalculator
	eventsSender                    entitycounters.EventsSender
	metaAlarmEventProcessor         libalarm.MetaAlarmEventProcessor
	metaAlarmStatesService          correlation.MetaAlarmStateService
	metricsSender                   metrics.Sender
	remediationRpcClient            engine.RPCClient
	encoder                         encoding.Encoder
	logger                          zerolog.Logger
}

func (p *entityToggledProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil {
		return result, nil
	}

	var componentStateChanged bool
	var newComponentState int

	if event.Entity.Enabled {
		var updatedServiceStates map[string]entitycounters.UpdatedServicesInfo
		err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
			componentStateChanged = false
			newComponentState = 0
			updatedServiceStates = nil
			result = Result{}

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

			result.IsCountersUpdated, updatedServiceStates, componentStateChanged, newComponentState, err = processComponentAndServiceCounters(
				ctx,
				p.entityServiceCountersCalculator,
				p.componentCountersCalculator,
				&result.Alarm,
				event.Entity,
				result.AlarmChange,
			)

			return err
		})
		if err != nil {
			return result, err
		}

		go p.postProcess(context.Background(), &event, updatedServiceStates, componentStateChanged, newComponentState)

		return result, nil
	}

	match := getOpenAlarmMatch(event)
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

			err = removeMetaAlarmStateOnResolve(ctx, p.metaAlarmRuleCollection, p.metaAlarmStatesService, result.Alarm)
			if err != nil {
				return err
			}
		}

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

	go postProcessResolve(
		context.Background(),
		event,
		result,
		updatedServiceStates,
		componentStateChanged,
		newComponentState,
		notAckedMetricType,
		p.eventsSender,
		p.metaAlarmEventProcessor,
		p.metricsSender,
		p.remediationRpcClient,
		p.pbehaviorCollection,
		p.encoder,
		p.logger,
	)

	return result, nil
}

func (p *entityToggledProcessor) postProcess(
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
