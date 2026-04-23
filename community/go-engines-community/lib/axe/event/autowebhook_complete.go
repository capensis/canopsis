package event

import (
	"context"
	"errors"

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

func NewAutoWebhookCompleteProcessor(
	client mongo.DbClient,
	metaAlarmPostProcessor MetaAlarmPostProcessor,
	metricsSender metrics.Sender,
	checkTicketStatusService webhook.CheckTicketStatusService,
	logger zerolog.Logger,
) Processor {
	return &autoWebhookCompleteProcessor{
		client:                   client,
		alarmCollection:          client.Collection(mongo.AlarmMongoCollection),
		resolvedAlarmCollection:  client.Collection(mongo.ResolvedAlarmMongoCollection),
		metaAlarmPostProcessor:   metaAlarmPostProcessor,
		checkTicketStatusService: checkTicketStatusService,
		metricsSender:            metricsSender,
		logger:                   logger,
	}
}

type autoWebhookCompleteProcessor struct {
	client                   mongo.DbClient
	alarmCollection          mongo.DbCollection
	resolvedAlarmCollection  mongo.DbCollection
	metaAlarmPostProcessor   MetaAlarmPostProcessor
	checkTicketStatusService webhook.CheckTicketStatusService
	metricsSender            metrics.Sender
	logger                   zerolog.Logger
}

func (p *autoWebhookCompleteProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	if event.Entity == nil || !event.Entity.Enabled {
		return Result{}, nil
	}

	match := getExactAlarmMatchWithStepsLimit(event)
	if match == nil {
		return Result{}, nil
	}

	match["v.steps"] = bson.M{"$not": bson.M{"$elemMatch": bson.M{
		"exec": event.Parameters.Execution,
		"_t":   bson.M{"$in": bson.A{types.AlarmStepWebhookComplete, types.AlarmStepWebhookFail}},
	}}}

	alarmChange := types.NewAlarmChange()
	var update []bson.M

	if event.Parameters.Ticket == "" {
		newStepQuery := execStepUpdateQueryWithInPbhInterval(types.AlarmStepWebhookComplete, event.Parameters.Execution,
			event.Parameters.Output, event.Parameters)
		update = []bson.M{
			{"$set": bson.M{
				"v.steps":            addStepUpdateQuery(newStepQuery),
				"v.last_update_date": event.Parameters.Timestamp,
			}},
		}

		alarmChange.Type = types.AlarmChangeTypeAutoWebhookComplete
	}

	var result Result

	err := p.client.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}

		if event.Parameters.Ticket != "" {
			if event.Parameters.CheckTicketStatus {
				job, err := p.checkTicketStatusService.CreateCheckStatusJob(ctx, event.Parameters.Execution, event.Parameters.TicketStatus, event.Parameters.TicketSourceStatus)
				if err != nil {
					return err
				}

				event.Parameters.TicketStatus = job.TicketStatus
				event.Parameters.TicketSourceStatus = job.TicketSourceStatus
				event.Parameters.TicketCheckStatusJobID = job.ID
			}

			newStepQuery := execStepUpdateQueryWithInPbhInterval(types.AlarmStepWebhookComplete, event.Parameters.Execution,
				event.Parameters.Output, event.Parameters)
			newTicketStepQuery := ticketStepUpdateQueryWithInPbhInterval(types.AlarmStepDeclareTicket,
				event.Parameters.Execution, event.Parameters.TicketInfo.GetStepMessage(), event.Parameters)
			update = []bson.M{
				{"$set": bson.M{
					"v.ticket":           newTicketStepQuery,
					"v.tickets":          addTicketUpdateQuery(newTicketStepQuery),
					"v.steps":            addStepUpdateQuery(newStepQuery, newTicketStepQuery),
					"v.last_update_date": event.Parameters.Timestamp,
				}},
				{"$unset": bson.A{
					"v.failed_ticket",
				}},
			}

			alarmChange.Type = types.AlarmChangeTypeAutoDeclareTicketWebhook
		}

		alarm := types.Alarm{}
		err := p.alarmCollection.FindOneAndUpdate(ctx, match, update,
			options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&alarm)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		if alarm.IsResolved() {
			_, err = p.resolvedAlarmCollection.UpdateOne(ctx, match, update)
			if err != nil {
				return err
			}
		}

		result.Alarm = alarm

		return nil
	})
	if err != nil || result.Alarm.ID == "" {
		return Result{}, err
	}

	result.Forward = true
	result.AlarmChange = alarmChange

	go p.postProcess(context.WithoutCancel(ctx), event, result)

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

	err := p.metaAlarmPostProcessor.Process(ctx, event, rpc.AxeResultEvent{
		Alarm:           &result.Alarm,
		AlarmChangeType: result.AlarmChange.Type,
	})
	if err != nil {
		p.logger.Err(err).Msg("cannot process meta alarm")
	}
}
