package event

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters/calculator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

func NewResolveCancelProcessor(
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
	return &resolveCancelProcessor{
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

type resolveCancelProcessor struct {
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

func (p *resolveCancelProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil {
		return result, nil
	}

	match := getOpenAlarmMatch(event)
	match["v.canceled"] = bson.M{"$ne": nil}
	result, updatedServiceStates, notAckedMetricType, componentStateChanged, newComponentState, err := processResolve(
		ctx,
		match,
		event,
		p.entityServiceCountersCalculator,
		p.componentCountersCalculator,
		p.metaAlarmStatesService,
		p.dbClient,
		p.alarmCollection,
		p.entityCollection,
		p.resolvedAlarmCollection,
		p.metaAlarmRuleCollection,
	)
	if err != nil || result.Alarm.ID == "" {
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
		p.metaAlarmPostProcessor,
		p.metricsSender,
		p.remediationRpcClient,
		p.pbehaviorCollection,
		p.encoder,
		p.logger,
	)

	return result, nil
}
