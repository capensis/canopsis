package event

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice/statecounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewPbhEnterProcessor(
	client mongo.DbClient,
	autoInstructionMatcher AutoInstructionMatcher,
	stateCountersService statecounters.StateCountersService,
	metricsSender metrics.Sender,
	remediationRpcClient engine.RPCClient,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Processor {
	return &pbhEnterProcessor{
		client:                 client,
		alarmCollection:        client.Collection(mongo.AlarmMongoCollection),
		entityCollection:       client.Collection(mongo.EntityMongoCollection),
		pbehaviorCollection:    client.Collection(mongo.PbehaviorMongoCollection),
		autoInstructionMatcher: autoInstructionMatcher,
		stateCountersService:   stateCountersService,
		metricsSender:          metricsSender,
		remediationRpcClient:   remediationRpcClient,
		encoder:                encoder,
		logger:                 logger,
	}
}

type pbhEnterProcessor struct {
	client                 mongo.DbClient
	alarmCollection        mongo.DbCollection
	entityCollection       mongo.DbCollection
	pbehaviorCollection    mongo.DbCollection
	autoInstructionMatcher AutoInstructionMatcher
	stateCountersService   statecounters.StateCountersService
	metricsSender          metrics.Sender
	remediationRpcClient   engine.RPCClient
	encoder                encoding.Encoder
	logger                 zerolog.Logger
}

func (p *pbhEnterProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || !event.Entity.Enabled || event.Parameters.PbehaviorInfo.IsZero() {
		return result, nil
	}

	match := getOpenAlarmMatchWithStepsLimit(event)
	match["v.pbehavior_info.id"] = bson.M{"$in": bson.A{nil, ""}}
	newStep := types.NewAlarmStep(types.AlarmStepPbhEnter, event.Parameters.Timestamp, event.Parameters.Author, event.Parameters.Output,
		event.Parameters.User, event.Parameters.Role, event.Parameters.Initiator, false)
	newStep.PbehaviorCanonicalType = event.Parameters.PbehaviorInfo.CanonicalType
	var update any
	if event.Parameters.PbehaviorInfo.IsActive() {
		update = bson.M{
			"$set":  bson.M{"v.pbehavior_info": event.Parameters.PbehaviorInfo},
			"$push": bson.M{"v.steps": newStep},
		}
	} else {
		update = []bson.M{
			{"$set": bson.M{
				"v.pbehavior_info": event.Parameters.PbehaviorInfo,
				"v.steps":          bson.M{"$concatArrays": bson.A{"$v.steps", bson.A{newStep}}},
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
				"v.inactive_start": event.Parameters.Timestamp,
			}},
			{"$unset": bson.A{
				"not_acked_metric_type",
				"not_acked_metric_send_time",
				"not_acked_since",
			}},
		}
	}

	var updatedServiceStates map[string]statecounters.UpdatedServicesInfo
	notAckedMetricType := ""
	err := p.client.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}
		updatedServiceStates = nil
		notAckedMetricType = ""

		updatedAlarm := types.Alarm{}
		opts := options.FindOneAndUpdate()
		if event.Parameters.PbehaviorInfo.IsActive() {
			opts = opts.SetReturnDocument(options.After)
		} else {
			opts = opts.SetReturnDocument(options.Before).
				SetProjection(bson.M{
					"not_acked_metric_type":      1,
					"not_acked_metric_send_time": 1,
				})
		}

		err := p.alarmCollection.FindOneAndUpdate(ctx, match, update, opts).Decode(&updatedAlarm)
		if err != nil {
			if !errors.Is(err, mongodriver.ErrNoDocuments) {
				return err
			}

			err := p.alarmCollection.FindOne(ctx, getOpenAlarmMatch(event), options.FindOne().SetProjection(bson.M{"_id": 1})).Err()
			if err == nil || !errors.Is(err, mongodriver.ErrNoDocuments) {
				return err
			}
		}

		alarm := types.Alarm{}
		if updatedAlarm.ID != "" {
			if event.Parameters.PbehaviorInfo.IsActive() {
				alarm = updatedAlarm
			} else {
				if updatedAlarm.NotAckedMetricSendTime != nil {
					notAckedMetricType = updatedAlarm.NotAckedMetricType
				}

				err = p.alarmCollection.FindOne(ctx, bson.M{"_id": updatedAlarm.ID}).Decode(&alarm)
				if err != nil {
					if errors.Is(err, mongodriver.ErrNoDocuments) {
						return nil
					}
					return err
				}
			}
		}

		result.Entity, err = updateEntity(ctx,
			bson.M{
				"_id":               event.Entity.ID,
				"pbehavior_info.id": bson.M{"$in": bson.A{nil, ""}},
			},
			bson.M{"$set": bson.M{
				"pbehavior_info":      event.Parameters.PbehaviorInfo,
				"last_pbehavior_date": event.Parameters.PbehaviorInfo.Timestamp,
			}},
			p.entityCollection,
		)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		alarmChange := types.NewAlarmChange()
		alarmChange.Type = types.AlarmChangeTypePbhEnter
		result.Forward = true
		result.Alarm = alarm
		result.AlarmChange = alarmChange

		if result.Alarm.ID == "" {
			updatedServiceStates, err = p.stateCountersService.UpdateServiceCounters(ctx, result.Entity, nil, result.AlarmChange)
		} else {
			updatedServiceStates, err = p.stateCountersService.UpdateServiceCounters(ctx, result.Entity, &result.Alarm, result.AlarmChange)
		}

		return err
	})

	if err != nil || result.AlarmChange.Type == types.AlarmChangeTypeNone {
		return result, err
	}

	if result.Alarm.ID != "" {
		result.IsInstructionMatched = isInstructionMatched(event, result, p.autoInstructionMatcher, p.logger)
	}

	go p.postProcess(context.Background(), event, result, updatedServiceStates, notAckedMetricType)

	return result, nil
}

func (p *pbhEnterProcessor) postProcess(
	ctx context.Context,
	event rpc.AxeEvent,
	result Result,
	updatedServiceStates map[string]statecounters.UpdatedServicesInfo,
	notAckedMetricType string,
) {
	p.metricsSender.SendEventMetrics(
		result.Alarm,
		*event.Entity,
		result.AlarmChange,
		event.Parameters.Timestamp.Time,
		event.Parameters.Initiator,
		event.Parameters.User,
		event.Parameters.Instruction,
		notAckedMetricType,
	)

	for servID, servInfo := range updatedServiceStates {
		err := p.stateCountersService.UpdateServiceState(ctx, servID, servInfo)
		if err != nil {
			p.logger.Err(err).Msg("failed to update service state")
		}
	}

	if result.Alarm.ID != "" {
		err := sendRemediationEvent(ctx, event, result, p.remediationRpcClient, p.encoder)
		if err != nil {
			p.logger.Err(err).Msg("cannot send event to engine-remediation")
		}
	}

	err := updatePbhLastAlarmDate(ctx, result, p.pbehaviorCollection)
	if err != nil {
		p.logger.Err(err).Msg("cannot update pbehavior")
	}
}
