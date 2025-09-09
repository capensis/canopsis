package event

import (
	"context"
	"errors"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewWebhookCompleteProcessor(
	client mongo.DbClient,
	metaAlarmPostProcessor MetaAlarmPostProcessor,
	metricsSender metrics.Sender,
	amqpPublisher libamqp.Publisher,
	eventGenerator event.Generator,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Processor {
	return &webhookCompleteProcessor{
		alarmCollection:         client.Collection(mongo.AlarmMongoCollection),
		resolvedAlarmCollection: client.Collection(mongo.ResolvedAlarmMongoCollection),
		metaAlarmPostProcessor:  metaAlarmPostProcessor,
		metricsSender:           metricsSender,
		amqpPublisher:           amqpPublisher,
		eventGenerator:          eventGenerator,
		encoder:                 encoder,
		logger:                  logger,
	}
}

type webhookCompleteProcessor struct {
	alarmCollection         mongo.DbCollection
	resolvedAlarmCollection mongo.DbCollection
	metaAlarmPostProcessor  MetaAlarmPostProcessor
	metricsSender           metrics.Sender
	amqpPublisher           libamqp.Publisher
	eventGenerator          event.Generator
	encoder                 encoding.Encoder
	logger                  zerolog.Logger
}

func (p *webhookCompleteProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || !event.Entity.Enabled {
		return result, nil
	}

	match := getExactAlarmMatchWithStepsLimit(event)
	if match == nil {
		return result, nil
	}

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
				"v.steps":            addStepUpdateQuery(newStepQuery),
				"v.last_update_date": event.Parameters.Timestamp,
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
				"v.ticket":           newTicketStepQuery,
				"v.tickets":          addTicketUpdateQuery(newTicketStepQuery),
				"v.steps":            addStepUpdateQuery(newStepQuery, newTicketStepQuery),
				"v.last_update_date": event.Parameters.Timestamp,
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

	if alarm.IsResolved() {
		_, err = p.resolvedAlarmCollection.UpdateOne(ctx, match, update)
		if err != nil {
			return result, err
		}
	}

	result.Forward = true
	result.Alarm = alarm
	result.AlarmChange = alarmChange

	go p.postProcess(context.WithoutCancel(ctx), event, result)

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

	err := p.metaAlarmPostProcessor.Process(ctx, event, rpc.AxeResultEvent{
		Alarm:           &result.Alarm,
		AlarmChangeType: result.AlarmChange.Type,
	})
	if err != nil {
		p.logger.Err(err).Msg("cannot process meta alarm")
	}

	sendTriggerEvent(ctx, event, result, p.amqpPublisher, p.encoder, p.eventGenerator, p.logger)
}
