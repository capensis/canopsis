package event

import (
	"context"
	"errors"

	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice/statecounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewAckRemoveProcessor(
	client mongo.DbClient,
	configProvider config.AlarmConfigProvider,
	stateCountersService statecounters.StateCountersService,
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor,
	metricsSender metrics.Sender,
	logger zerolog.Logger,
) Processor {
	return &ackRemoveProcessor{
		client:                  client,
		alarmCollection:         client.Collection(mongo.AlarmMongoCollection),
		entityCollection:        client.Collection(mongo.EntityMongoCollection),
		configProvider:          configProvider,
		stateCountersService:    stateCountersService,
		metaAlarmEventProcessor: metaAlarmEventProcessor,
		metricsSender:           metricsSender,
		logger:                  logger,
	}
}

type ackRemoveProcessor struct {
	client                  mongo.DbClient
	alarmCollection         mongo.DbCollection
	entityCollection        mongo.DbCollection
	configProvider          config.AlarmConfigProvider
	stateCountersService    statecounters.StateCountersService
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor
	metricsSender           metrics.Sender
	logger                  zerolog.Logger
}

func (p *ackRemoveProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || !event.Entity.Enabled {
		return result, nil
	}

	entity := *event.Entity
	match := getOpenAlarmMatchWithStepsLimit(event)
	match["v.ack"] = bson.M{"$ne": nil}
	conf := p.configProvider.Get()
	output := utils.TruncateString(event.Parameters.Output, conf.OutputLength)
	newStep := types.NewAlarmStep(types.AlarmStepAckRemove, event.Parameters.Timestamp, event.Parameters.Author, output,
		event.Parameters.User, event.Parameters.Role, event.Parameters.Initiator)
	update := bson.M{
		"$set":  bson.M{"not_acked_since": event.Parameters.Timestamp},
		"$push": bson.M{"v.steps": newStep},
		"$unset": bson.M{
			"v.ack": "",
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedServiceStates map[string]statecounters.UpdatedServicesInfo
	firstTry := true
	err := p.client.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}
		updatedServiceStates = nil
		var err error
		if !firstTry {
			entity, err = findEntity(ctx, event.Entity.ID, p.entityCollection)
			if err != nil {
				return err
			}
		}

		firstTry = false
		alarm := types.Alarm{}
		err = p.alarmCollection.FindOneAndUpdate(ctx, match, update, opts).Decode(&alarm)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		alarmChange := types.NewAlarmChange()
		alarmChange.Type = types.AlarmChangeTypeAckremove
		result.Forward = true
		result.Alarm = alarm
		result.AlarmChange = alarmChange

		if event.Parameters.IdleRuleApply != "" {
			result.Entity, err = updateEntityByID(ctx, entity.ID, bson.M{"$set": bson.M{
				"last_idle_rule_apply": event.Parameters.IdleRuleApply,
			}}, p.entityCollection)
			if err != nil {
				return err
			}
		}

		updatedServiceStates, err = p.stateCountersService.UpdateServiceCounters(ctx, entity, &result.Alarm, result.AlarmChange)

		return err
	})

	if err != nil || result.Alarm.ID == "" {
		return result, err
	}

	go p.postProcess(context.Background(), event, result, updatedServiceStates)

	return result, nil
}

func (p *ackRemoveProcessor) postProcess(
	ctx context.Context,
	event rpc.AxeEvent,
	res Result,
	updatedServiceStates map[string]statecounters.UpdatedServicesInfo,
) {
	p.metricsSender.SendEventMetrics(
		res.Alarm,
		*event.Entity,
		res.AlarmChange,
		event.Parameters.Timestamp.Time,
		event.Parameters.Initiator,
		event.Parameters.User,
		event.Parameters.Instruction,
		"",
	)

	for servID, servInfo := range updatedServiceStates {
		err := p.stateCountersService.UpdateServiceState(ctx, servID, servInfo)
		if err != nil {
			p.logger.Err(err).Msg("failed to update service state")
		}
	}

	err := p.metaAlarmEventProcessor.ProcessAxeRpc(ctx, event, rpc.AxeResultEvent{
		Alarm:           &res.Alarm,
		AlarmChangeType: res.AlarmChange.Type,
	})
	if err != nil {
		p.logger.Err(err).Msg("cannot process meta alarm")
	}
}
