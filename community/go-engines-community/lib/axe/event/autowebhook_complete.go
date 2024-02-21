package event

import (
	"context"
	"errors"

	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewAutoWebhookCompleteProcessor(
	client mongo.DbClient,
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor,
	metricsSender metrics.Sender,
	logger zerolog.Logger,
) Processor {
	return &autoWebhookCompleteProcessor{
		alarmCollection:         client.Collection(mongo.AlarmMongoCollection),
		metaAlarmEventProcessor: metaAlarmEventProcessor,
		metricsSender:           metricsSender,
		logger:                  logger,
	}
}

type autoWebhookCompleteProcessor struct {
	alarmCollection         mongo.DbCollection
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor
	metricsSender           metrics.Sender
	logger                  zerolog.Logger
}

func (p *autoWebhookCompleteProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || !event.Entity.Enabled {
		return result, nil
	}

	match := getOpenAlarmMatchWithStepsLimit(event)
	match["v.steps"] = bson.M{"$not": bson.M{"$elemMatch": bson.M{
		"exec": event.Parameters.Execution,
		"_t":   bson.M{"$in": bson.A{types.AlarmStepWebhookComplete, types.AlarmStepWebhookFail}},
	}}}
	alarmChange := types.NewAlarmChange()
	var update []bson.M
	if event.Parameters.Ticket == "" {
		newStep := types.NewAlarmStep(types.AlarmStepWebhookComplete, event.Parameters.Timestamp, event.Parameters.Author,
			event.Parameters.Output, event.Parameters.User, event.Parameters.Role, event.Parameters.Initiator, false)
		newStep.Execution = event.Parameters.Execution
		newStepQuery := stepUpdateQuery(newStep)
		update = []bson.M{
			{"$set": bson.M{
				"v.steps": addStepUpdateQuery(newStepQuery),
			}},
		}
		alarmChange.Type = types.AlarmChangeTypeAutoWebhookComplete
	} else {
		newStep := types.NewAlarmStep(types.AlarmStepWebhookComplete, event.Parameters.Timestamp, event.Parameters.Author,
			event.Parameters.Output, event.Parameters.User, event.Parameters.Role, event.Parameters.Initiator, false)
		newStep.Execution = event.Parameters.Execution
		newStepQuery := stepUpdateQuery(newStep)
		newTicketStep := types.NewTicketStep(types.AlarmStepDeclareTicket, event.Parameters.Timestamp, event.Parameters.Author,
			event.Parameters.TicketInfo.GetStepMessage(), event.Parameters.User, event.Parameters.Role, event.Parameters.Initiator,
			event.Parameters.TicketInfo, false)
		newTicketStep.Execution = event.Parameters.Execution
		newTicketStepQuery := stepUpdateQuery(newTicketStep)
		update = []bson.M{
			{"$set": bson.M{
				"v.ticket":  newTicketStepQuery,
				"v.tickets": addTicketUpdateQuery(newTicketStepQuery),
				"v.steps":   addStepUpdateQuery(newStepQuery, newTicketStepQuery),
			}},
		}
		alarmChange.Type = types.AlarmChangeTypeAutoDeclareTicketWebhook
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	alarm := types.Alarm{}
	err := p.alarmCollection.FindOneAndUpdate(ctx, match, update, opts).Decode(&alarm)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return result, nil
		}

		return result, err
	}

	result.Forward = true
	result.Alarm = alarm
	result.AlarmChange = alarmChange

	go p.postProcess(context.Background(), event, result)

	return result, nil
}

func (p *autoWebhookCompleteProcessor) postProcess(
	ctx context.Context,
	event rpc.AxeEvent,
	result Result,
) {
	p.metricsSender.SendEventMetrics(
		result.Alarm,
		*event.Entity,
		result.AlarmChange,
		event.Parameters.Timestamp.Time,
		event.Parameters.Initiator,
		event.Parameters.User,
		event.Parameters.Instruction,
		"",
	)

	err := p.metaAlarmEventProcessor.ProcessAxeRpc(ctx, event, rpc.AxeResultEvent{
		Alarm:           &result.Alarm,
		AlarmChangeType: result.AlarmChange.Type,
	})
	if err != nil {
		p.logger.Err(err).Msg("cannot process meta alarm")
	}
}
