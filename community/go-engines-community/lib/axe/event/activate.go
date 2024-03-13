package event

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewActivateProcessor(
	client mongo.DbClient,
	autoInstructionMatcher AutoInstructionMatcher,
	remediationRpcClient engine.RPCClient,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Processor {
	return &activateProcessor{
		alarmCollection:        client.Collection(mongo.AlarmMongoCollection),
		autoInstructionMatcher: autoInstructionMatcher,
		remediationRpcClient:   remediationRpcClient,
		encoder:                encoder,
		logger:                 logger,
	}
}

type activateProcessor struct {
	alarmCollection        mongo.DbCollection
	autoInstructionMatcher AutoInstructionMatcher
	remediationRpcClient   engine.RPCClient
	encoder                encoding.Encoder
	logger                 zerolog.Logger
}

func (p *activateProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil {
		return result, nil
	}

	match := getOpenAlarmMatch(event)
	match["v.activation_date"] = nil
	newStep := NewAlarmStep(types.AlarmStepActivate, event.Parameters, false)
	update := bson.M{
		"$set":  bson.M{"v.activation_date": event.Parameters.Timestamp},
		"$push": bson.M{"v.steps": newStep},
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
	alarmChange.Type = types.AlarmChangeTypeActivate
	result.Forward = true
	result.Alarm = alarm
	result.AlarmChange = alarmChange
	result.IsInstructionMatched = isInstructionMatched(event, result, p.autoInstructionMatcher, p.logger)
	go p.postProcess(context.Background(), event, result)

	return result, nil
}

func (p *activateProcessor) postProcess(
	ctx context.Context,
	event rpc.AxeEvent,
	result Result,
) {
	err := sendRemediationEvent(ctx, event, result, p.remediationRpcClient, p.encoder)
	if err != nil {
		p.logger.Err(err).Msg("cannot send event to engine-remediation")
	}
}
