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

func NewUpdateStatusProcessor(
	dbClient mongo.DbClient,
	alarmStatusService alarmstatus.Service,
	metaAlarmPostProcessor MetaAlarmPostProcessor,
	logger zerolog.Logger,
) Processor {
	return &updateStatusProcessor{
		dbClient:               dbClient,
		alarmCollection:        dbClient.Collection(mongo.AlarmMongoCollection),
		alarmStatusService:     alarmStatusService,
		metaAlarmPostProcessor: metaAlarmPostProcessor,
		logger:                 logger,
	}
}

type updateStatusProcessor struct {
	dbClient               mongo.DbClient
	alarmCollection        mongo.DbCollection
	alarmStatusService     alarmstatus.Service
	metaAlarmPostProcessor MetaAlarmPostProcessor
	logger                 zerolog.Logger
}

func (p *updateStatusProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil {
		return result, nil
	}

	match := getOpenAlarmMatch(event)
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

		currentStatus := alarm.Value.Status.Value
		newStatus, statusRuleName := p.alarmStatusService.ComputeStatus(alarm, *event.Entity)
		if newStatus == currentStatus {
			return nil
		}

		alarmStepType := types.AlarmStepStatusIncrease
		if alarm.Value.Status.Value > newStatus {
			alarmStepType = types.AlarmStepStatusDecrease
		}

		statusStepMessage := ConcatOutputAndRuleName(event.Parameters.Output, statusRuleName)
		newStepStatusQuery := valStepUpdateQueryWithInPbhInterval(alarmStepType, newStatus, statusStepMessage, event.Parameters)
		matchUpdate := getOpenAlarmMatchWithStepsLimit(event)
		update := []bson.M{
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
		alarmChange.Type = types.AlarmChangeTypeUpdateStatus
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

func (p *updateStatusProcessor) postProcess(
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
