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

func NewAckProcessor(
	client mongo.DbClient,
	configProvider config.AlarmConfigProvider,
	stateCountersService statecounters.StateCountersService,
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor,
	metricsSender metrics.Sender,
	logger zerolog.Logger,
) Processor {
	return &ackProcessor{
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

type ackProcessor struct {
	client                  mongo.DbClient
	alarmCollection         mongo.DbCollection
	entityCollection        mongo.DbCollection
	configProvider          config.AlarmConfigProvider
	stateCountersService    statecounters.StateCountersService
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor
	metricsSender           metrics.Sender
	logger                  zerolog.Logger
}

func (p *ackProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || !event.Entity.Enabled {
		return result, nil
	}

	entity := *event.Entity
	match := getOpenAlarmMatchWithStepsLimit(event)
	match["v.ack"] = nil
	conf := p.configProvider.Get()
	output := utils.TruncateString(event.Parameters.Output, conf.OutputLength)
	newStep := types.NewAlarmStep(types.AlarmStepAck, event.Parameters.Timestamp, event.Parameters.Author, output,
		event.Parameters.User, event.Parameters.Role, event.Parameters.Initiator, true)
	newStepQuery := stepUpdateQuery(newStep)
	update := []bson.M{
		{"$set": bson.M{
			"v.ack":   newStepQuery,
			"v.steps": addStepUpdateQuery(newStepQuery),
		}},
		{"$unset": bson.A{
			"not_acked_metric_type",
			"not_acked_metric_send_time",
			"not_acked_since",
		}},
	}
	var updatedServiceStates map[string]statecounters.UpdatedServicesInfo
	notAckedMetricType := ""

	err := p.client.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}
		updatedServiceStates = nil
		notAckedMetricType = ""

		doubleAck := false
		beforeAlarm := types.Alarm{}
		opts := options.FindOneAndUpdate().
			SetReturnDocument(options.Before).
			SetProjection(bson.M{
				"not_acked_metric_type":      1,
				"not_acked_metric_send_time": 1,
			})
		err := p.alarmCollection.FindOneAndUpdate(ctx, match, update, opts).Decode(&beforeAlarm)
		if err != nil {
			if !errors.Is(err, mongodriver.ErrNoDocuments) {
				return err
			}

			if !conf.AllowDoubleAck {
				return nil
			}

			doubleAck = true
			delete(match, "v.ack")
			err = p.alarmCollection.FindOneAndUpdate(ctx, match, update, opts).Decode(&beforeAlarm)
			if err != nil {
				if errors.Is(err, mongodriver.ErrNoDocuments) {
					return nil
				}

				return err
			}
		}

		if beforeAlarm.NotAckedMetricSendTime != nil {
			notAckedMetricType = beforeAlarm.NotAckedMetricType
		}

		alarm := types.Alarm{}
		err = p.alarmCollection.FindOne(ctx, bson.M{"_id": beforeAlarm.ID}).Decode(&alarm)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}
			return err
		}

		alarmChange := types.NewAlarmChange()
		if doubleAck {
			alarmChange.Type = types.AlarmChangeTypeDoubleAck
		} else {
			alarmChange.Type = types.AlarmChangeTypeAck
		}

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

	go p.postProcess(context.Background(), event, result, updatedServiceStates, notAckedMetricType)

	return result, nil
}

func (p *ackProcessor) postProcess(
	ctx context.Context,
	event rpc.AxeEvent,
	result Result,
	updatedServiceStates map[string]statecounters.UpdatedServicesInfo,
	notAckedMetricType string,
) {
	p.metricsSender.SendEventMetrics(
		result.Alarm,
		*event.Entity,
		result.AlarmChange,
		event.Parameters.Timestamp.Time,
		event.Parameters.Initiator,
		event.Parameters.User,
		event.Parameters.Instruction,
		notAckedMetricType,
	)

	for servID, servInfo := range updatedServiceStates {
		err := p.stateCountersService.UpdateServiceState(ctx, servID, servInfo)
		if err != nil {
			p.logger.Err(err).Msg("failed to update service state")
		}
	}

	err := p.metaAlarmEventProcessor.ProcessAxeRpc(ctx, event, rpc.AxeResultEvent{
		Alarm:           &result.Alarm,
		AlarmChangeType: result.AlarmChange.Type,
	})
	if err != nil {
		p.logger.Err(err).Msg("cannot process meta alarm")
	}
}

func getOpenAlarmMatch(event rpc.AxeEvent) bson.M {
	if event.Alarm != nil {
		return bson.M{
			"_id":        event.Alarm.ID,
			"v.resolved": nil,
		}
	}

	if event.AlarmID != "" {
		return bson.M{
			"_id":        event.AlarmID,
			"v.resolved": nil,
		}
	}

	return bson.M{
		"d":          event.Entity.ID,
		"v.resolved": nil,
	}
}

func getOpenAlarmMatchWithStepsLimit(event rpc.AxeEvent) bson.M {
	match := getOpenAlarmMatch(event)
	match["$expr"] = bson.M{"$lt": bson.A{bson.M{"$size": "$v.steps"}, types.AlarmStepsHardLimit}}
	return match
}

func stepUpdateQuery(newStep types.AlarmStep) bson.M {
	return bson.M{"$cond": bson.M{
		"if": bson.M{"$and": []bson.M{
			{"$eq": bson.A{bson.M{"$type": "$v.pbehavior_info.id"}, "string"}},
			{"$ne": bson.A{"$v.pbehavior_info.id", ""}},
		}},
		"then": bson.M{"$mergeObjects": bson.A{
			newStep,
			bson.M{"in_pbh": true},
		}},
		"else": newStep,
	}}
}

func addStepUpdateQuery(newStepQueries ...bson.M) bson.M {
	return bson.M{"$concatArrays": bson.A{"$v.steps", newStepQueries}}
}
