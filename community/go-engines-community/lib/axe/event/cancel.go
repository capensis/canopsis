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

func NewCancelProcessor(
	client mongo.DbClient,
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor,
	logger zerolog.Logger,
) Processor {
	return &cancelProcessor{
		client:                  client,
		alarmCollection:         client.Collection(mongo.AlarmMongoCollection),
		entityCollection:        client.Collection(mongo.EntityMongoCollection),
		metaAlarmEventProcessor: metaAlarmEventProcessor,
		logger:                  logger,
	}
}

type cancelProcessor struct {
	client                  mongo.DbClient
	alarmCollection         mongo.DbCollection
	entityCollection        mongo.DbCollection
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor
	logger                  zerolog.Logger
}

func (p *cancelProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || !event.Entity.Enabled {
		return result, nil
	}

	entity := *event.Entity
	match := getOpenAlarmMatchWithStepsLimit(event)
	match["v.canceled"] = nil
	output := event.Parameters.Output
	newStep := types.NewAlarmStep(types.AlarmStepCancel, event.Parameters.Timestamp, event.Parameters.Author, output,
		event.Parameters.User, event.Parameters.Role, event.Parameters.Initiator)
	newIncStepStatus := types.NewAlarmStep(types.AlarmStepStatusIncrease, event.Parameters.Timestamp, event.Parameters.Author, output,
		event.Parameters.User, event.Parameters.Role, event.Parameters.Initiator)
	newDecStepStatus := types.NewAlarmStep(types.AlarmStepStatusDecrease, event.Parameters.Timestamp, event.Parameters.Author, output,
		event.Parameters.User, event.Parameters.Role, event.Parameters.Initiator)
	newStatus := types.CpsNumber(types.AlarmStatusCancelled)
	newIncStepStatus.Value = newStatus
	newDecStepStatus.Value = newStatus
	update := []bson.M{
		{"$set": bson.M{
			"v.canceled": newStep,
			"v.status": bson.M{"$cond": bson.M{
				"if":   bson.M{"$gt": bson.A{newStatus, "$v.status.val"}},
				"then": newIncStepStatus,
				"else": newDecStepStatus,
			}},
			"v.steps": bson.M{"$concatArrays": bson.A{
				"$v.steps",
				bson.A{
					newStep,
					bson.M{"$cond": bson.M{
						"if":   bson.M{"$gt": bson.A{newStatus, "$v.status.val"}},
						"then": newIncStepStatus,
						"else": newDecStepStatus,
					}},
				},
			}},
			"v.state_changes_since_status_update": 0,
			"v.last_update_date":                  event.Parameters.Timestamp,
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
		alarmChange.Type = types.AlarmChangeTypeCancel
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

func (p *cancelProcessor) postProcess(
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
