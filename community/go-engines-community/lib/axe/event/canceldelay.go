package event

import (
	"context"
	"errors"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewCancelDelayProcessor(
	client mongo.DbClient,
	metaAlarmPostProcessor MetaAlarmPostProcessor,
	logger zerolog.Logger,
) Processor {
	return &cancelDelayProcessor{
		client:                   client,
		alarmCollection:          client.Collection(mongo.AlarmMongoCollection),
		entityCollection:         client.Collection(mongo.EntityMongoCollection),
		cancelDelayJobCollection: client.Collection(mongo.CancelDelayJobCollection),
		metaAlarmPostProcessor:   metaAlarmPostProcessor,
		logger:                   logger,
	}
}

type cancelDelayProcessor struct {
	client                   mongo.DbClient
	alarmCollection          mongo.DbCollection
	entityCollection         mongo.DbCollection
	cancelDelayJobCollection mongo.DbCollection
	metaAlarmPostProcessor   MetaAlarmPostProcessor
	logger                   zerolog.Logger
}

func (p *cancelDelayProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}

	if event.Entity == nil || !event.Entity.Enabled {
		return result, nil
	}

	match := getOpenAlarmMatch(event)
	match["v.canceled"] = nil
	newStatus := types.CpsNumber(types.AlarmStatusCancelledWithDelay)
	newIncStepStatusQuery := valStepUpdateQueryWithInPbhInterval(types.AlarmStepStatusIncrease, newStatus,
		event.Parameters.Output, event.Parameters)
	newDecStepStatusQuery := valStepUpdateQueryWithInPbhInterval(types.AlarmStepStatusDecrease, newStatus,
		event.Parameters.Output, event.Parameters)
	newStatusStepQuery := bson.M{"$cond": bson.M{
		"if":   bson.M{"$gt": bson.A{newStatus, "$v.status.val"}},
		"then": newIncStepStatusQuery,
		"else": newDecStepStatusQuery,
	}}
	update := []bson.M{
		{"$set": bson.M{
			"v.canceled":                          newStatusStepQuery,
			"v.status":                            newStatusStepQuery,
			"v.steps":                             addStepUpdateQuery(newStatusStepQuery),
			"v.state_changes_since_status_update": 0,
			"v.last_update_date":                  event.Parameters.Timestamp,
			"v.last_st_upd_dt":                    event.Parameters.Timestamp,
		}},
		{"$unset": bson.A{
			"v.cancel_delay_value",
		}},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	alarmChange := types.NewAlarmChange()

	err := p.client.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}
		alarm := types.Alarm{}
		alarmChange.Type = types.AlarmChangeTypeNone

		_, err := p.cancelDelayJobCollection.DeleteOne(ctx, bson.M{"_id": event.AlarmID})
		if err != nil {
			return fmt.Errorf("failed to delete cancel_delay job on cancel_delay event: %w", err)
		}

		err = p.alarmCollection.FindOneAndUpdate(ctx, match, update, opts).Decode(&alarm)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return fmt.Errorf("failed to find alarm on cancel_delay event: %w", err)
		}

		alarmChange.Type = types.AlarmChangeTypeCancel

		result.Forward = true
		result.Alarm = alarm
		result.AlarmChange = alarmChange

		return nil
	})

	if err != nil || result.Alarm.ID == "" {
		return result, err
	}

	go p.postProcess(context.Background(), event, result)

	return result, nil
}

func (p *cancelDelayProcessor) postProcess(
	ctx context.Context,
	event rpc.AxeEvent,
	result Result,
) {
	err := p.metaAlarmPostProcessor.Process(ctx, event, rpc.AxeResultEvent{
		Alarm:           &result.Alarm,
		AlarmChangeType: result.AlarmChange.Type,
	})
	if err != nil {
		p.logger.Err(err).Msg("cannot process meta alarm")
	}
}
