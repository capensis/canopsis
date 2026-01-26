package event

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewTicketRemoveProcessor(
	client mongo.DbClient,
	metaAlarmPostProcessor MetaAlarmPostProcessor,
	metricsSender metrics.Sender,
	logger zerolog.Logger,
) Processor {
	return &ticketRemoveProcessor{
		client:                 client,
		alarmCollection:        client.Collection(mongo.AlarmMongoCollection),
		entityCollection:       client.Collection(mongo.EntityMongoCollection),
		metaAlarmPostProcessor: metaAlarmPostProcessor,
		metricsSender:          metricsSender,
		logger:                 logger,
	}
}

type ticketRemoveProcessor struct {
	client                 mongo.DbClient
	alarmCollection        mongo.DbCollection
	entityCollection       mongo.DbCollection
	metaAlarmPostProcessor MetaAlarmPostProcessor
	metricsSender          metrics.Sender
	logger                 zerolog.Logger
}

func (p *ticketRemoveProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || !event.Entity.Enabled {
		return result, nil
	}

	match := getOpenAlarmMatchWithStepsLimit(event)
	match["v.tickets.ticket"] = event.Parameters.Ticket
	newStepQuery := ticketStepUpdateQueryWithInPbhInterval(types.AlarmStepTicketRemove, "",
		event.Parameters.Output, event.Parameters)
	update := []bson.M{
		{"$set": bson.M{
			"v.tickets": bson.M{"$filter": bson.M{
				"input": "$v.tickets",
				"cond":  bson.M{"$ne": bson.A{"$$this.ticket", event.Parameters.Ticket}},
			}},
			"v.steps":            addStepUpdateQuery(newStepQuery),
			"v.last_update_date": event.Parameters.Timestamp,
		}},
		{"$set": bson.M{
			"v.ticket": bson.M{"$arrayElemAt": bson.A{"$v.tickets", -1}},
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
	alarmChange.Type = types.AlarmChangeTypeTicketRemove
	result.Forward = true
	result.Alarm = alarm
	result.AlarmChange = alarmChange

	go p.postProcess(context.WithoutCancel(ctx), event, result)

	return result, nil
}

func (p *ticketRemoveProcessor) postProcess(
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
