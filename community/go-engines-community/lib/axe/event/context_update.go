package event

import (
	"context"
	"errors"
	"fmt"

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
		dbClient:           dbClient,
		alarmCollection:    dbClient.Collection(mongo.AlarmMongoCollection),
		alarmStatusService: alarmStatusService,
		logger:             logger,
		countersHelper:     newCountersHelper(entityServiceCountersCalculator, componentCountersCalculator, eventsSender, logger),
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
	dbClient           mongo.DbClient
	alarmCollection    mongo.DbCollection
	alarmStatusService alarmstatus.Service
	logger             zerolog.Logger
	countersHelper     *countersHelper
	upstreamHelper     *upstreamHelper
}

func (p *contextUpdateProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || event.Entity.ID == "" || !event.Entity.Enabled {
		return result, nil
	}

	entity := *event.Entity
	var err error
	if entity.IsUpstreamChanged {
		if event.Parameters.StateSettingUpdated && entity.StateInfo == nil {
			countersRes := countersResult{}
			match := getOpenAlarmMatch(event)
			err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
				result = Result{}
				entity = *event.Entity
				countersRes = countersResult{}

				alarm := types.Alarm{}
				err := p.alarmCollection.FindOne(ctx, match).Decode(&alarm)
				if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
					return err
				}

				if alarm.ID == "" {
					result, err = p.upstreamHelper.UpdateAlarm(ctx, event, alarm, entity)
					if err != nil {
						return err
					}

					result.IsCountersUpdated, countersRes, err = p.countersHelper.CalculateCounters(
						ctx,
						&result.Alarm,
						&entity,
						result.AlarmChange,
					)

					return err
				}

				newStatus, statusRuleName, err := p.alarmStatusService.ComputeStatusOnStatusChange(ctx, alarm, entity)
				if err != nil {
					return fmt.Errorf("cannot compute alarm status: %w", err)
				}

				var newState types.CpsNumber
				switch newStatus {
				case types.AlarmStatusOngoing:
					// close alarm which were created by state setting method
					newState = types.AlarmStateOK
					newStatus = types.AlarmStatusOff
				case types.AlarmStatusUnknown:
					// override state which were created by state setting method
					newState = types.AlarmStateForUnknown
				}

				currentStatus := alarm.Value.Status.Value
				if newStatus != currentStatus {
					currentState := alarm.Value.State.Value
					if newState == currentState {
						result, err = p.upstreamHelper.UpdateAlarmStatus(ctx, alarm, entity, event, newStatus, statusRuleName)
					} else {
						result, err = p.upstreamHelper.UpdateAlarmStateAndStatus(ctx, alarm, entity, event, newState, newStatus, statusRuleName)
					}

					if err != nil {
						return err
					}

					result.IsCountersUpdated, countersRes, err = p.countersHelper.CalculateCounters(
						ctx,
						&result.Alarm,
						&entity,
						result.AlarmChange,
					)

					return err
				}

				result.IsCountersUpdated, countersRes, err = p.countersHelper.CalculateCounters(
					ctx,
					&alarm,
					&entity,
					result.AlarmChange,
				)

				return err
			})
			if err != nil || result.Alarm.ID == "" {
				return result, err
			}

			go p.upstreamHelper.PostProcess(context.WithoutCancel(ctx), event, result, countersRes)

			return result, nil
		}

		result, _, err = p.upstreamHelper.Process(ctx, event, true)
		if err != nil {
			return result, err
		}

		return result, nil
	}

	countersRes := countersResult{}
	match := getOpenAlarmMatch(event)
	err = p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}
		entity = *event.Entity
		countersRes = countersResult{}
		alarm := types.Alarm{}
		err := p.alarmCollection.FindOne(ctx, match).Decode(&alarm)
		if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
			return err
		}

		result.IsCountersUpdated, countersRes, err = p.countersHelper.CalculateCounters(
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

	go p.countersHelper.UpdateStates(context.WithoutCancel(ctx), countersRes)

	return result, nil
}
