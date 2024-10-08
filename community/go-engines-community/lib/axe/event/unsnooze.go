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

func NewUnsnoozeProcessor(
	client mongo.DbClient,
	autoInstructionMatcher AutoInstructionMatcher,
	remediationRpcClient engine.RPCClient,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Processor {
	return &unsnoozeProcessor{
		alarmCollection:        client.Collection(mongo.AlarmMongoCollection),
		autoInstructionMatcher: autoInstructionMatcher,
		remediationRpcClient:   remediationRpcClient,
		encoder:                encoder,
		logger:                 logger,
	}
}

type unsnoozeProcessor struct {
	alarmCollection        mongo.DbCollection
	autoInstructionMatcher AutoInstructionMatcher
	remediationRpcClient   engine.RPCClient
	encoder                encoding.Encoder
	logger                 zerolog.Logger
}

func (p *unsnoozeProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || !event.Entity.Enabled {
		return Result{}, nil
	}

	match := getOpenAlarmMatch(event)
	match["v.snooze"] = bson.M{"$ne": nil}
	update := []bson.M{
		{"$set": bson.M{
			"v.snooze_duration": bson.M{"$sum": bson.A{
				"$v.snooze_duration",
				bson.M{"$subtract": bson.A{
					event.Parameters.Timestamp,
					"$v.snooze.t",
				}},
			}},
			"v.inactive_duration": bson.M{"$sum": bson.A{
				"$v.inactive_duration",
				bson.M{"$subtract": bson.A{
					event.Parameters.Timestamp,
					"$v.inactive_start",
				}},
			}},
			"v.inactive_start": updateInactiveStart(event.Parameters.Timestamp, false, true, true),
		}},
		{"$unset": "v.snooze"},
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
	alarmChange.Type = types.AlarmChangeTypeUnsnooze
	result.Forward = true
	result.Alarm = alarm
	result.AlarmChange = alarmChange
	result.IsInstructionMatched = isInstructionMatched(event, result, p.autoInstructionMatcher, p.logger)
	go p.postProcess(context.Background(), event, result)

	return result, nil
}

func (p *unsnoozeProcessor) postProcess(
	ctx context.Context,
	event rpc.AxeEvent,
	result Result,
) {
	err := sendRemediationEvent(ctx, event, result, p.remediationRpcClient, p.encoder)
	if err != nil {
		p.logger.Err(err).Msg("cannot send event to engine-remediation")
	}
}
