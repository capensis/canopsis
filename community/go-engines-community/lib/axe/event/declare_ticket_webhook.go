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
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/webhook"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewDeclareTicketWebhookProcessor(
	client mongo.DbClient,
	metricsSender metrics.Sender,
	amqpPublisher libamqp.Publisher,
	eventGenerator event.Generator,
	checkTicketStatusService webhook.CheckTicketStatusService,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Processor {
	return &declareTicketWebhookProcessor{
		dbClient:                 client,
		alarmCollection:          client.Collection(mongo.AlarmMongoCollection),
		metricsSender:            metricsSender,
		amqpPublisher:            amqpPublisher,
		eventGenerator:           eventGenerator,
		checkTicketStatusService: checkTicketStatusService,
		encoder:                  encoder,
		logger:                   logger,
	}
}

type declareTicketWebhookProcessor struct {
	dbClient                 mongo.DbClient
	alarmCollection          mongo.DbCollection
	metricsSender            metrics.Sender
	amqpPublisher            libamqp.Publisher
	eventGenerator           event.Generator
	encoder                  encoding.Encoder
	checkTicketStatusService webhook.CheckTicketStatusService
	logger                   zerolog.Logger
}

func (p *declareTicketWebhookProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	if event.Entity == nil || !event.Entity.Enabled {
		return Result{}, nil
	}

	var result Result

	err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}

		if event.Parameters.TicketCheckStatusJobID != "" {
			job, err := p.checkTicketStatusService.AddAlarmToCheckStatusJob(ctx, event.Parameters.TicketCheckStatusJobID, event.AlarmID)
			if err != nil {
				return err
			}

			if job.ID != "" {
				event.Parameters.TicketStatus = job.TicketStatus
				event.Parameters.TicketPrevStatus = job.PrevTicketStatus
				event.Parameters.TicketSourceStatus = job.TicketSourceStatus
				event.Parameters.TicketPrevSourceStatus = job.PrevTicketSourceStatus
				event.Parameters.TicketLastCheckTime = job.CheckedAt
			}
		}

		match := getOpenAlarmMatchWithStepsLimit(event)
		newTicketStepQuery := ticketStepUpdateQueryWithInPbhInterval(types.AlarmStepDeclareTicket, "",
			event.Parameters.Output, event.Parameters)
		update := []bson.M{
			{"$set": bson.M{
				"v.ticket":           newTicketStepQuery,
				"v.tickets":          addTicketUpdateQuery(newTicketStepQuery),
				"v.steps":            addStepUpdateQuery(newTicketStepQuery),
				"v.last_update_date": event.Parameters.Timestamp,
			}},
			{"$unset": bson.A{
				"v.failed_ticket",
			}},
		}

		alarm := types.Alarm{}
		err := p.alarmCollection.FindOneAndUpdate(ctx, match, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&alarm)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		result.Alarm = alarm

		return nil
	})
	if err != nil || result.Alarm.ID == "" {
		return Result{}, err
	}

	alarmChange := types.NewAlarmChange()
	alarmChange.Type = types.AlarmChangeTypeDeclareTicketWebhook
	result.Forward = true
	result.AlarmChange = alarmChange

	go p.postProcess(context.WithoutCancel(ctx), event, result)

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

	sendTriggerEvent(ctx, event, result, p.amqpPublisher, p.encoder, p.eventGenerator, p.logger)
}
