package event

import (
	"context"
	"fmt"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmtag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters/calculator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func NewEntityToggledProcessor(
	dbClient mongo.DbClient,
	alarmConfigProvider config.AlarmConfigProvider,
	alarmStatusService alarmstatus.Service,
	pbhTypeResolver pbehavior.EntityTypeResolver,
	autoInstructionMatcher AutoInstructionMatcher,
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator,
	componentCountersCalculator calculator.ComponentCountersCalculator,
	eventsSender entitycounters.EventsSender,
	metaAlarmPostProcessor MetaAlarmPostProcessor,
	metaAlarmStatesService correlation.MetaAlarmStateService,
	metricsSender metrics.Sender,
	remediationRpcClient engine.RPCClient,
	internalTagAlarmMatcher alarmtag.InternalTagAlarmMatcher,
	eventGenerator event.Generator,
	amqpPublisher libamqp.Publisher,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Processor {
	return &entityToggledProcessor{
		dbClient:                dbClient,
		alarmCollection:         dbClient.Collection(mongo.AlarmMongoCollection),
		closeDelayJobCollection: dbClient.Collection(mongo.CloseDelayJobCollection),
		logger:                  logger,
		countersHelper:          newCountersHelper(entityServiceCountersCalculator, componentCountersCalculator, eventsSender, logger),
		resolveHelper: newResolveHelper(
			dbClient,
			alarmConfigProvider,
			alarmStatusService,
			pbhTypeResolver,
			autoInstructionMatcher,
			internalTagAlarmMatcher,
			entityServiceCountersCalculator,
			componentCountersCalculator,
			metaAlarmStatesService,
			eventsSender,
			metaAlarmPostProcessor,
			metricsSender,
			remediationRpcClient,
			eventGenerator,
			amqpPublisher,
			encoder,
			logger,
		),
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

type entityToggledProcessor struct {
	dbClient                mongo.DbClient
	alarmCollection         mongo.DbCollection
	closeDelayJobCollection mongo.DbCollection
	logger                  zerolog.Logger
	resolveHelper           *resolveHelper
	countersHelper          *countersHelper
	upstreamHelper          *upstreamHelper
}

func (p *entityToggledProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil {
		return result, nil
	}

	if event.Parameters.Initiator != types.InitiatorSystem {
		return Result{}, fmt.Errorf("unknown initiator %q", event.Parameters.Initiator)
	}

	countersRes := countersResult{}
	var err error
	if event.Entity.Enabled {
		var alarm types.Alarm
		result, alarm, err = p.upstreamHelper.Process(ctx, event, true)
		if err != nil {
			return result, err
		}

		if result.AlarmChange.Type == "" {
			alarmChange := types.NewAlarmChange()
			alarmChange.Type = types.AlarmChangeTypeEnabled
			result.Forward = true
			result.Alarm = alarm
			result.AlarmChange = alarmChange
		}

		go p.upstreamHelper.PostProcess(context.WithoutCancel(ctx), event, result, countersRes)

		return result, nil
	}

	match := getOpenAlarmMatch(event)
	notAckedMetricType := ""
	err = p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}
		countersRes = countersResult{}
		notAckedMetricType = ""

		beforeAlarm, err := p.resolveHelper.UpdateAlarmToResolve(ctx, match, event.Parameters)
		if err != nil {
			return err
		}

		entity := *event.Entity
		if beforeAlarm.ID != "" {
			_, err = p.closeDelayJobCollection.DeleteOne(ctx, bson.M{"_id": beforeAlarm.ID})
			if err != nil {
				return fmt.Errorf("failed to delete close_delay job on entitytoggled event: %w", err)
			}

			if beforeAlarm.NotAckedMetricSendTime != nil {
				notAckedMetricType = beforeAlarm.NotAckedMetricType
			}

			entity, err = p.resolveHelper.UpdateEntityOfResolvedAlarm(ctx, event.Entity.ID)
			if err != nil || entity.ID == "" {
				return err
			}

			alarm, err := p.resolveHelper.CopyAlarmToResolvedCollection(ctx, beforeAlarm.ID)
			if err != nil || alarm.ID == "" {
				return err
			}

			alarmChange := types.NewAlarmChange()
			alarmChange.Type = types.AlarmChangeTypeResolve
			result.Forward = true
			result.Alarm = alarm
			result.Entity = entity
			result.AlarmChange = alarmChange

			err = p.resolveHelper.RemoveMetaAlarmStateOnResolve(ctx, result.Alarm)
			if err != nil {
				return err
			}
		}

		result.IsCountersUpdated, countersRes, err = p.countersHelper.CalculateCounters(
			ctx,
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
		go p.resolveHelper.PostProcess(
			context.WithoutCancel(ctx),
			event,
			result,
			countersRes,
			notAckedMetricType,
		)
	} else {
		go p.countersHelper.UpdateStates(context.WithoutCancel(ctx), countersRes)
	}

	return result, nil
}
