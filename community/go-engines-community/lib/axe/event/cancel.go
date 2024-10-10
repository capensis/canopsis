package event

import (
	"context"
	"errors"

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
	metaAlarmPostProcessor MetaAlarmPostProcessor,
	logger zerolog.Logger,
) Processor {
	return &cancelProcessor{
		client:                 client,
		alarmCollection:        client.Collection(mongo.AlarmMongoCollection),
		entityCollection:       client.Collection(mongo.EntityMongoCollection),
		metaAlarmPostProcessor: metaAlarmPostProcessor,
		logger:                 logger,
	}
}

type cancelProcessor struct {
	client                 mongo.DbClient
	alarmCollection        mongo.DbCollection
	entityCollection       mongo.DbCollection
	metaAlarmPostProcessor MetaAlarmPostProcessor
	logger                 zerolog.Logger
}

func (p *cancelProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || !event.Entity.Enabled {
		return result, nil
	}

	entity := *event.Entity
	match := getOpenAlarmMatch(event)
	match["v.canceled"] = nil
	newStatus := types.CpsNumber(types.AlarmStatusCancelled)
	newIncStepStatusQuery := valStepUpdateQueryWithInPbhInterval(types.AlarmStepStatusIncrease, newStatus,
		event.Parameters.Output, event.Parameters)
	newDecStepStatusQuery := valStepUpdateQueryWithInPbhInterval(types.AlarmStepStatusDecrease, newStatus,
		event.Parameters.Output, event.Parameters)
	newStepQuery := bson.M{"$cond": bson.M{
		"if":   bson.M{"$gt": bson.A{newStatus, "$v.status.val"}},
		"then": newIncStepStatusQuery,
		"else": newDecStepStatusQuery,
	}}
	update := []bson.M{
		{"$set": bson.M{
			"v.canceled":                          newStepQuery,
			"v.status":                            newStepQuery,
			"v.steps":                             addStepUpdateQuery(newStepQuery),
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
	err := p.metaAlarmPostProcessor.Process(ctx, event, rpc.AxeResultEvent{
		Alarm:           &result.Alarm,
		AlarmChangeType: result.AlarmChange.Type,
	})
	if err != nil {
		p.logger.Err(err).Msg("cannot process meta alarm")
	}
}
