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

func NewCommentProcessor(
	client mongo.DbClient,
	metaAlarmPostProcessor MetaAlarmPostProcessor,
	logger zerolog.Logger,
) Processor {
	return &commentProcessor{
		alarmCollection:        client.Collection(mongo.AlarmMongoCollection),
		metaAlarmPostProcessor: metaAlarmPostProcessor,
		logger:                 logger,
	}
}

type commentProcessor struct {
	alarmCollection        mongo.DbCollection
	metaAlarmPostProcessor MetaAlarmPostProcessor
	logger                 zerolog.Logger
}

func (p *commentProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || !event.Entity.Enabled {
		return result, nil
	}

	match := getOpenAlarmMatchWithStepsLimit(event)
	newStepQuery := stepUpdateQueryWithInPbhInterval(types.AlarmStepComment, event.Parameters.Output, event.Parameters)
	update := []bson.M{
		{"$set": bson.M{
			"v.last_comment": newStepQuery,
			"v.steps":        addStepUpdateQuery(newStepQuery),
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
	alarmChange.Type = types.AlarmChangeTypeComment
	result.Forward = true
	result.Alarm = alarm
	result.AlarmChange = alarmChange

	go p.postProcess(context.Background(), event, result)

	return result, nil
}

func (p *commentProcessor) postProcess(
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
