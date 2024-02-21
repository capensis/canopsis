package event

import (
	"context"
	"errors"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
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

func NewDeclareTicketWebhookProcessor(
	client mongo.DbClient,
	metricsSender metrics.Sender,
	amqpPublisher libamqp.Publisher,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Processor {
	return &declareTicketWebhookProcessor{
		alarmCollection: client.Collection(mongo.AlarmMongoCollection),
		metricsSender:   metricsSender,
		amqpPublisher:   amqpPublisher,
		encoder:         encoder,
		logger:          logger,
	}
}

type declareTicketWebhookProcessor struct {
	alarmCollection mongo.DbCollection
	metricsSender   metrics.Sender
	amqpPublisher   libamqp.Publisher
	encoder         encoding.Encoder
	logger          zerolog.Logger
}

func (p *declareTicketWebhookProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || !event.Entity.Enabled {
		return result, nil
	}

	match := getOpenAlarmMatchWithStepsLimit(event)
	newTicketStep := types.NewTicketStep(types.AlarmStepDeclareTicket, event.Parameters.Timestamp, event.Parameters.Author,
		event.Parameters.TicketInfo.GetStepMessage(), event.Parameters.User, event.Parameters.Role, event.Parameters.Initiator,
		event.Parameters.TicketInfo, false)
	newTicketStepQuery := stepUpdateQuery(newTicketStep)
	update := []bson.M{
		{"$set": bson.M{
			"v.ticket":  newTicketStepQuery,
			"v.tickets": addTicketUpdateQuery(newTicketStepQuery),
			"v.steps":   addStepUpdateQuery(newTicketStepQuery),
		}},
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

	alarmChange := types.NewAlarmChange()
	alarmChange.Type = types.AlarmChangeTypeDeclareTicketWebhook
	result.Forward = true
	result.Alarm = alarm
	result.AlarmChange = alarmChange

	go p.postProcess(context.Background(), event, result)

	return result, nil
}

func (p *declareTicketWebhookProcessor) postProcess(
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

	sendTriggerEvent(ctx, event, result, p.amqpPublisher, p.encoder, p.logger)
}
