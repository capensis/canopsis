package event

import (
	"context"
	"errors"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewWebhookCompleteProcessor(
	client mongo.DbClient,
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor,
	metricsSender metrics.Sender,
	amqpPublisher libamqp.Publisher,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Processor {
	return &webhookCompleteProcessor{
		alarmCollection:         client.Collection(mongo.AlarmMongoCollection),
		metaAlarmEventProcessor: metaAlarmEventProcessor,
		metricsSender:           metricsSender,
		amqpPublisher:           amqpPublisher,
		encoder:                 encoder,
		logger:                  logger,
	}
}

type webhookCompleteProcessor struct {
	alarmCollection         mongo.DbCollection
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor
	metricsSender           metrics.Sender
	amqpPublisher           libamqp.Publisher
	encoder                 encoding.Encoder
	logger                  zerolog.Logger
}

func (p *webhookCompleteProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
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
		newStepQuery := execStepUpdateQueryWithInPbhInterval(types.AlarmStepWebhookComplete, event.Parameters.RuleExecution,
			event.Parameters.Output, event.Parameters)
		update = []bson.M{
			{"$set": bson.M{
				"v.steps": addStepUpdateQuery(newStepQuery),
			}},
		}
		alarmChange.Type = types.AlarmChangeTypeWebhookComplete
	} else {
		newStepQuery := execStepUpdateQueryWithInPbhInterval(types.AlarmStepWebhookComplete, event.Parameters.RuleExecution,
			event.Parameters.Output, event.Parameters)
		newTicketStepQuery := ticketStepUpdateQueryWithInPbhInterval(types.AlarmStepDeclareTicket,
			event.Parameters.RuleExecution, event.Parameters.TicketInfo.GetStepMessage(), event.Parameters)
		update = []bson.M{
			{"$set": bson.M{
				"v.ticket":  newTicketStepQuery,
				"v.tickets": addTicketUpdateQuery(newTicketStepQuery),
				"v.steps":   addStepUpdateQuery(newStepQuery, newTicketStepQuery),
			}},
		}
		alarmChange.Type = types.AlarmChangeTypeDeclareTicketWebhook
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

func (p *webhookCompleteProcessor) postProcess(
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

	sendTriggerEvent(ctx, event, result, p.amqpPublisher, p.encoder, p.logger)
}
