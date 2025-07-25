package event

import (
	"context"
	"errors"
	"fmt"

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

func NewResolveDeletedProcessor(
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
	return &resolveDeletedProcessor{
		dbClient:                dbClient,
		entityCollection:        dbClient.Collection(mongo.EntityMongoCollection),
		closeDelayJobCollection: dbClient.Collection(mongo.CloseDelayJobCollection),
		logger:                  logger,
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
		componentAndServiceCountersHelper: newComponentAndServiceCountersHelper(entityServiceCountersCalculator, componentCountersCalculator, eventsSender, logger),
	}
}

type resolveDeletedProcessor struct {
	helper                            *resolveHelper
	dbClient                          mongo.DbClient
	entityCollection                  mongo.DbCollection
	closeDelayJobCollection           mongo.DbCollection
	logger                            zerolog.Logger
	componentAndServiceCountersHelper *componentAndServiceCountersHelper
}

func (p *resolveDeletedProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || event.Entity.SoftDeleted == nil || event.Entity.ResolveDeletedEventProcessed != nil {
		return result, nil
	}

	now := datetime.NewCpsTime()
	match := getOpenAlarmMatch(event)
	countersRes := componentAndServiceCountersResult{}
	notAckedMetricType := ""
	err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}
		countersRes = componentAndServiceCountersResult{}
		notAckedMetricType = ""

		beforeAlarm, err := p.helper.UpdateAlarmToResolve(ctx, match, event.Parameters)
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
		result.IsCountersUpdated, countersRes, err = p.componentAndServiceCountersHelper.Process(
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
		go p.componentAndServiceCountersHelper.PostProcess(context.WithoutCancel(ctx), countersRes)
	}

	return result, nil
}
