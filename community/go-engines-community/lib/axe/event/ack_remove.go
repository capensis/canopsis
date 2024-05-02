package event

import (
	"context"
	"errors"

	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters/calculator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewAckRemoveProcessor(
	client mongo.DbClient,
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator,
	eventsSender entitycounters.EventsSender,
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor,
	metricsSender metrics.Sender,
	logger zerolog.Logger,
) Processor {
	return &ackRemoveProcessor{
		client:                          client,
		alarmCollection:                 client.Collection(mongo.AlarmMongoCollection),
		entityCollection:                client.Collection(mongo.EntityMongoCollection),
		entityServiceCountersCalculator: entityServiceCountersCalculator,
		eventsSender:                    eventsSender,
		metaAlarmEventProcessor:         metaAlarmEventProcessor,
		metricsSender:                   metricsSender,
		logger:                          logger,
	}
}

type ackRemoveProcessor struct {
	client                          mongo.DbClient
	alarmCollection                 mongo.DbCollection
	entityCollection                mongo.DbCollection
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator
	eventsSender                    entitycounters.EventsSender
	metaAlarmEventProcessor         libalarm.MetaAlarmEventProcessor
	metricsSender                   metrics.Sender
	logger                          zerolog.Logger
}

func (p *ackRemoveProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || !event.Entity.Enabled {
		return result, nil
	}

	entity := *event.Entity
	match := getOpenAlarmMatchWithStepsLimit(event)
	match["v.ack"] = bson.M{"$ne": nil}
	newStepQuery := stepUpdateQueryWithInPbhInterval(types.AlarmStepAckRemove, event.Parameters.Output, event.Parameters)
	update := []bson.M{
		{"$set": bson.M{
			"not_acked_since": event.Parameters.Timestamp,
			"v.steps":         addStepUpdateQuery(newStepQuery),
		}},
		{"$unset": bson.A{
			"v.ack",
		}},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedServiceStates map[string]entitycounters.UpdatedServicesInfo

	err := p.client.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}
		updatedServiceStates = nil

		alarm := types.Alarm{}
		err := p.alarmCollection.FindOneAndUpdate(ctx, match, update, opts).Decode(&alarm)
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

		result.IsCountersUpdated, updatedServiceStates, err = p.entityServiceCountersCalculator.CalculateCounters(ctx, &entity, &result.Alarm, result.AlarmChange)

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
	updatedServiceStates map[string]entitycounters.UpdatedServicesInfo,
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
		err := p.eventsSender.UpdateServiceState(ctx, servID, servInfo)
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
