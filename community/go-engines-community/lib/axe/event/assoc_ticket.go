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

func NewAssocTicketProcessor(
	client mongo.DbClient,
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor,
	metricsSender metrics.Sender,
	logger zerolog.Logger,
) Processor {
	return &assocTicketProcessor{
		client:                  client,
		alarmCollection:         client.Collection(mongo.AlarmMongoCollection),
		entityCollection:        client.Collection(mongo.EntityMongoCollection),
		metaAlarmEventProcessor: metaAlarmEventProcessor,
		metricsSender:           metricsSender,
		logger:                  logger,
	}
}

type assocTicketProcessor struct {
	client                  mongo.DbClient
	alarmCollection         mongo.DbCollection
	entityCollection        mongo.DbCollection
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor
	metricsSender           metrics.Sender
	logger                  zerolog.Logger
}

func (p *assocTicketProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || !event.Entity.Enabled {
		return result, nil
	}

	entity := *event.Entity
	match := getOpenAlarmMatchWithStepsLimit(event)
	newStepQuery := ticketStepUpdateQueryWithInPbhInterval(types.AlarmStepAssocTicket, "",
		event.Parameters.TicketInfo.GetStepMessage(), event.Parameters)
	update := []bson.M{
		{"$set": bson.M{
			"v.ticket":  newStepQuery,
			"v.steps":   addStepUpdateQuery(newStepQuery),
			"v.tickets": addTicketUpdateQuery(newStepQuery),
		}},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := p.client.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}
		alarm := types.Alarm{}
		err := p.alarmCollection.FindOneAndUpdate(ctx, match, update, opts).Decode(&alarm)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		alarmChange := types.NewAlarmChange()
		alarmChange.Type = types.AlarmChangeTypeAssocTicket
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

		return nil
	})

	if err != nil || result.Alarm.ID == "" {
		return result, err
	}

	go p.postProcess(context.Background(), event, result)

	return result, nil
}

func (p *assocTicketProcessor) postProcess(
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
