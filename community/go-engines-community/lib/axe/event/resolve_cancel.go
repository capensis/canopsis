package event

import (
	"context"

	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice/statecounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

func NewResolveCancelProcessor(
	dbClient mongo.DbClient,
	stateCountersService statecounters.StateCountersService,
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor,
	metaAlarmStatesService correlation.MetaAlarmStateService,
	metricsSender metrics.Sender,
	remediationRpcClient engine.RPCClient,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Processor {
	return &resolveCancelProcessor{
		dbClient:                dbClient,
		alarmCollection:         dbClient.Collection(mongo.AlarmMongoCollection),
		entityCollection:        dbClient.Collection(mongo.EntityMongoCollection),
		resolvedAlarmCollection: dbClient.Collection(mongo.ResolvedAlarmMongoCollection),
		metaAlarmRuleCollection: dbClient.Collection(mongo.MetaAlarmRulesMongoCollection),
		stateCountersService:    stateCountersService,
		metaAlarmEventProcessor: metaAlarmEventProcessor,
		metaAlarmStatesService:  metaAlarmStatesService,
		metricsSender:           metricsSender,
		remediationRpcClient:    remediationRpcClient,
		encoder:                 encoder,
		logger:                  logger,
	}
}

type resolveCancelProcessor struct {
	dbClient                mongo.DbClient
	alarmCollection         mongo.DbCollection
	entityCollection        mongo.DbCollection
	resolvedAlarmCollection mongo.DbCollection
	metaAlarmRuleCollection mongo.DbCollection
	stateCountersService    statecounters.StateCountersService
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor
	metaAlarmStatesService  correlation.MetaAlarmStateService
	metricsSender           metrics.Sender
	remediationRpcClient    engine.RPCClient
	encoder                 encoding.Encoder
	logger                  zerolog.Logger
}

func (p *resolveCancelProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil {
		return result, nil
	}

	match := getOpenAlarmMatch(event)
	match["v.canceled"] = bson.M{"$ne": nil}
	result, updatedServiceStates, notAckedMetricType, err := processResolve(ctx, match, event, p.stateCountersService, p.metaAlarmStatesService, p.dbClient, p.alarmCollection, p.entityCollection, p.resolvedAlarmCollection, p.metaAlarmRuleCollection)
	if err != nil || result.Alarm.ID == "" {
		return result, err
	}

	go postProcessResolve(context.Background(), event, result, updatedServiceStates, notAckedMetricType, p.stateCountersService, p.metaAlarmEventProcessor, p.metricsSender, p.remediationRpcClient, p.encoder, p.logger)

	return result, nil
}
