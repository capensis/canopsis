package event

import (
	"context"
	"errors"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmtag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
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
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewRecomputeEntityServiceProcessor(
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
	return &recomputeEntityServiceProcessor{
		dbClient:                        dbClient,
		entityCollection:                dbClient.Collection(mongo.EntityMongoCollection),
		entityServiceCountersCalculator: entityServiceCountersCalculator,
		eventsSender:                    eventsSender,
		logger:                          logger,
		helper: newResolveHelper(
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
		countersHelper: newCountersHelper(entityServiceCountersCalculator, componentCountersCalculator, eventsSender, logger),
	}
}

type recomputeEntityServiceProcessor struct {
	dbClient                        mongo.DbClient
	entityCollection                mongo.DbCollection
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator
	eventsSender                    entitycounters.EventsSender
	logger                          zerolog.Logger
	helper                          *resolveHelper
	countersHelper                  *countersHelper
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
	countersRes := countersResult{}
	notAckedMetricType := ""
	err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}
		countersRes = countersResult{}
		notAckedMetricType = ""

		beforeAlarm, err := p.helper.UpdateAlarmToResolve(ctx, match, event.Parameters)
		if err != nil {
			return err
		}

		entityUpdate := bson.M{}
		if beforeAlarm.ID != "" {
			if beforeAlarm.NotAckedMetricSendTime != nil {
				notAckedMetricType = beforeAlarm.NotAckedMetricType
			}

			alarm, err := p.helper.CopyAlarmToResolvedCollection(ctx, beforeAlarm.ID)
			if err != nil || alarm.ID == "" {
				return err
			}

			entityUpdate = p.helper.GetResolveEntityUpdate()
			alarmChange := types.NewAlarmChange()
			alarmChange.Type = types.AlarmChangeTypeResolve
			result.Forward = true
			result.Alarm = alarm
			result.AlarmChange = alarmChange

			err = p.helper.RemoveMetaAlarmStateOnResolve(ctx, result.Alarm)
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
		go p.helper.PostProcess(context.WithoutCancel(ctx), event, result, countersRes, notAckedMetricType)
	} else {
		go p.countersHelper.UpdateStates(context.WithoutCancel(ctx), countersRes)
	}

	return result, nil
}
