package event

import (
	"context"
	"errors"

	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewSnoozeProcessor(
	client mongo.DbClient,
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor,
	logger zerolog.Logger,
) Processor {
	return &snoozeProcessor{
		client:                  client,
		alarmCollection:         client.Collection(mongo.AlarmMongoCollection),
		entityCollection:        client.Collection(mongo.EntityMongoCollection),
		metaAlarmEventProcessor: metaAlarmEventProcessor,
		logger:                  logger,
	}
}

type snoozeProcessor struct {
	client                  mongo.DbClient
	alarmCollection         mongo.DbCollection
	entityCollection        mongo.DbCollection
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor
	logger                  zerolog.Logger
}

func (p *snoozeProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || !event.Entity.Enabled || event.Parameters.Duration == nil {
		return result, nil
	}

	entity := *event.Entity
	match := getOpenAlarmMatchWithStepsLimit(event)
	match["v.snooze"] = nil
	stepVal := types.CpsNumber(event.Parameters.Duration.AddTo(event.Parameters.Timestamp).Unix())
	newStepQuery := valStepUpdateQueryWithInPbhInterval(types.AlarmStepSnooze, stepVal, event.Parameters.Output, event.Parameters)
	update := []bson.M{
		{"$set": bson.M{
			"v.snooze":         newStepQuery,
			"v.steps":          addStepUpdateQuery(newStepQuery),
			"v.inactive_start": event.Parameters.Timestamp,
			"v.inactive_duration": bson.M{"$sum": bson.A{
				"$v.inactive_duration",
				bson.M{"$cond": bson.M{
					"if": bson.M{"$gt": bson.A{"$v.inactive_start", 0}},
					"then": bson.M{"$subtract": bson.A{
						event.Parameters.Timestamp,
						"$v.inactive_start",
					}},
					"else": 0,
				}},
			}},
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
		alarmChange.Type = types.AlarmChangeTypeSnooze
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

		return err
	})

	if err != nil || result.Alarm.ID == "" {
		return result, err
	}

	go p.postProcess(context.Background(), event, result)

	return result, nil
}

func (p *snoozeProcessor) postProcess(
	ctx context.Context,
	event rpc.AxeEvent,
	result Result,
) {
	err := p.metaAlarmEventProcessor.ProcessAxeRpc(ctx, event, rpc.AxeResultEvent{
		Alarm:           &result.Alarm,
		AlarmChangeType: result.AlarmChange.Type,
	})
	if err != nil {
		p.logger.Err(err).Msg("cannot process meta alarm")
	}
}
