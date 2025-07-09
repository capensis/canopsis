package event

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmtag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters/calculator"
	libevent "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
)

func NewContextUpdateProcessor(
	dbClient mongo.DbClient,
	alarmConfigProvider config.AlarmConfigProvider,
	alarmStatusService alarmstatus.Service,
	pbhTypeResolver pbehavior.EntityTypeResolver,
	autoInstructionMatcher AutoInstructionMatcher,
	metaAlarmPostProcessor MetaAlarmPostProcessor,
	metricsSender metrics.Sender,
	remediationRpcClient engine.RPCClient,
	internalTagAlarmMatcher alarmtag.InternalTagAlarmMatcher,
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator,
	componentCountersCalculator calculator.ComponentCountersCalculator,
	eventsSender entitycounters.EventsSender,
	eventGenerator libevent.Generator,
	amqpPublisher amqp.Publisher,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Processor {
	return &contextUpdateProcessor{
		dbClient:                          dbClient,
		alarmCollection:                   dbClient.Collection(mongo.AlarmMongoCollection),
		logger:                            logger,
		componentAndServiceCountersHelper: newComponentAndServiceCountersHelper(entityServiceCountersCalculator, componentCountersCalculator, eventsSender, logger),
		upstreamHelper: newUpstreamHelper(
			dbClient,
			alarmConfigProvider,
			alarmStatusService,
			pbhTypeResolver,
			autoInstructionMatcher,
			metaAlarmPostProcessor,
			metricsSender,
			remediationRpcClient,
			internalTagAlarmMatcher,
			entityServiceCountersCalculator,
			componentCountersCalculator,
			eventsSender,
			eventGenerator,
			amqpPublisher,
			encoder,
			logger,
		),
	}
}

type contextUpdateProcessor struct {
	dbClient                          mongo.DbClient
	alarmCollection                   mongo.DbCollection
	logger                            zerolog.Logger
	componentAndServiceCountersHelper *componentAndServiceCountersHelper
	upstreamHelper                    *upstreamHelper
}

func (p *contextUpdateProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || event.Entity.ID == "" || !event.Entity.Enabled {
		return result, nil
	}

	entity := *event.Entity
	var err error
	if entity.IsUpstreamChanged {
		result, _, err = p.upstreamHelper.Process(ctx, event, true)
		if err != nil {
			return result, err
		}

		return result, nil
	}

	countersRes := componentAndServiceCountersResult{}
	match := getOpenAlarmMatch(event)
	err = p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}
		entity = *event.Entity
		countersRes = componentAndServiceCountersResult{}
		alarm := types.Alarm{}
		err := p.alarmCollection.FindOne(ctx, match).Decode(&alarm)
		if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
			return err
		}

		result.IsCountersUpdated, countersRes, err = p.componentAndServiceCountersHelper.Process(
			ctx,
			&alarm,
			&entity,
			result.AlarmChange,
		)

		return err
	})
	if err != nil {
		return result, err
	}

	go p.componentAndServiceCountersHelper.PostProcess(context.WithoutCancel(ctx), countersRes)

	return result, nil
}
