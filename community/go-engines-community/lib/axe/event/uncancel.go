package event

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewUncancelProcessor(
	dbClient mongo.DbClient,
	alarmStatusService alarmstatus.Service,
	metaAlarmPostProcessor MetaAlarmPostProcessor,
	logger zerolog.Logger,
) Processor {
	return &uncancelProcessor{
		dbClient:               dbClient,
		alarmCollection:        dbClient.Collection(mongo.AlarmMongoCollection),
		alarmStatusService:     alarmStatusService,
		metaAlarmPostProcessor: metaAlarmPostProcessor,
		logger:                 logger,
	}
}

type uncancelProcessor struct {
	dbClient               mongo.DbClient
	alarmCollection        mongo.DbCollection
	alarmStatusService     alarmstatus.Service
	metaAlarmPostProcessor MetaAlarmPostProcessor
	logger                 zerolog.Logger
}

func (p *uncancelProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || !event.Entity.Enabled {
		return result, nil
	}

	match := getOpenAlarmMatch(event)
	match["v.canceled"] = bson.M{"$ne": nil}
	matchUpdate := getOpenAlarmMatchWithStepsLimit(event)
	matchUpdate["v.canceled"] = bson.M{"$ne": nil}
	err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}
		alarm := types.Alarm{}
		err := p.alarmCollection.FindOne(ctx, match).Decode(&alarm)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		alarm.Value.Canceled = nil
		newStatus, statusRuleName := p.alarmStatusService.ComputeStatus(alarm, *event.Entity)
		alarmStepType := types.AlarmStepStatusIncrease
		if alarm.Value.Status.Value > newStatus {
			alarmStepType = types.AlarmStepStatusDecrease
		}

		statusStepMessage := ConcatOutputAndRuleName(event.Parameters.Output, statusRuleName)
		newStepStatusQuery := valStepUpdateQueryWithInPbhInterval(alarmStepType, newStatus, statusStepMessage, event.Parameters)
		update := []bson.M{
			{"$unset": bson.A{
				"v.canceled",
			}},
			{"$set": bson.M{
				"v.status":                            newStepStatusQuery,
				"v.state_changes_since_status_update": 0,
				"v.last_update_date":                  event.Parameters.Timestamp,
				"v.steps":                             addStepUpdateQuery(newStepStatusQuery),
			}},
		}
		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
		updatedAlarm := types.Alarm{}
		err = p.alarmCollection.FindOneAndUpdate(ctx, matchUpdate, update, opts).Decode(&updatedAlarm)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		alarmChange := types.NewAlarmChange()
		alarmChange.Type = types.AlarmChangeTypeUncancel
		result.Forward = true
		result.Alarm = updatedAlarm
		result.AlarmChange = alarmChange
		return nil
	})

	if err != nil || result.Alarm.ID == "" {
		return result, err
	}

	go p.postProcess(context.Background(), event, result)

	return result, nil
}

func (p *uncancelProcessor) postProcess(
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
