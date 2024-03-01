package event

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters/calculator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewPbhLeaveProcessor(
	client mongo.DbClient,
	autoInstructionMatcher AutoInstructionMatcher,
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator,
	componentCountersCalculator calculator.ComponentCountersCalculator,
	eventsSender entitycounters.EventsSender,
	metricsSender metrics.Sender,
	remediationRpcClient engine.RPCClient,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Processor {
	return &pbhLeaveProcessor{
		client:                          client,
		alarmCollection:                 client.Collection(mongo.AlarmMongoCollection),
		entityCollection:                client.Collection(mongo.EntityMongoCollection),
		autoInstructionMatcher:          autoInstructionMatcher,
		entityServiceCountersCalculator: entityServiceCountersCalculator,
		componentCountersCalculator:     componentCountersCalculator,
		eventsSender:                    eventsSender,
		metricsSender:                   metricsSender,
		remediationRpcClient:            remediationRpcClient,
		encoder:                         encoder,
		logger:                          logger,
	}
}

type pbhLeaveProcessor struct {
	client                          mongo.DbClient
	alarmCollection                 mongo.DbCollection
	entityCollection                mongo.DbCollection
	autoInstructionMatcher          AutoInstructionMatcher
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator
	componentCountersCalculator     calculator.ComponentCountersCalculator
	eventsSender                    entitycounters.EventsSender
	metricsSender                   metrics.Sender
	remediationRpcClient            engine.RPCClient
	encoder                         encoding.Encoder
	logger                          zerolog.Logger
}

func (p *pbhLeaveProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || !event.Entity.Enabled || !event.Parameters.PbehaviorInfo.IsZero() {
		return result, nil
	}

	match := getOpenAlarmMatchWithStepsLimit(event)
	match["v.pbehavior_info.id"] = bson.M{"$nin": bson.A{nil, ""}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedServiceStates map[string]entitycounters.UpdatedServicesInfo
	var componentStateChanged bool
	var newComponentState int

	err := p.client.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}
		updatedServiceStates = nil

		alarm := types.Alarm{}
		err := p.alarmCollection.FindOne(ctx, match).Decode(&alarm)
		if err != nil {
			if !errors.Is(err, mongodriver.ErrNoDocuments) {
				return err
			}

			err := p.alarmCollection.FindOne(ctx, getOpenAlarmMatch(event), options.FindOne().SetProjection(bson.M{"_id": 1})).Err()
			if err == nil || !errors.Is(err, mongodriver.ErrNoDocuments) {
				return err
			}
		}

		alarmChange := types.NewAlarmChange()
		if alarm.ID == "" {
			alarmChange.PreviousEntityPbehaviorTime = event.Entity.PbehaviorInfo.Timestamp
			alarmChange.PreviousPbehaviorTypeID = event.Entity.PbehaviorInfo.TypeID
			alarmChange.PreviousPbehaviorCannonicalType = event.Entity.PbehaviorInfo.CanonicalType
		} else {
			alarmChange.PreviousPbehaviorTime = alarm.Value.PbehaviorInfo.Timestamp
			alarmChange.PreviousEntityPbehaviorTime = event.Entity.PbehaviorInfo.Timestamp
			alarmChange.PreviousPbehaviorTypeID = alarm.Value.PbehaviorInfo.TypeID
			alarmChange.PreviousPbehaviorCannonicalType = alarm.Value.PbehaviorInfo.CanonicalType
			newStep := types.NewAlarmStep(types.AlarmStepPbhLeave, event.Parameters.Timestamp, event.Parameters.Author, event.Parameters.Output,
				event.Parameters.User, event.Parameters.Role, event.Parameters.Initiator)
			newStep.PbehaviorCanonicalType = alarm.Value.PbehaviorInfo.CanonicalType
			update := bson.M{
				"$push":  bson.M{"v.steps": newStep},
				"$unset": bson.M{"v.pbehavior_info": ""},
			}

			if !alarm.Value.PbehaviorInfo.IsActive() {
				update["$inc"] = bson.M{
					"v.pbh_inactive_duration": int64(event.Parameters.Timestamp.Sub(alarm.Value.PbehaviorInfo.Timestamp.Time).Seconds()),
					"v.inactive_duration":     int64(event.Parameters.Timestamp.Sub(alarm.Value.InactiveStart.Time).Seconds()),
				}

				var inactiveStart *datetime.CpsTime
				if alarm.Value.Snooze != nil || alarm.InactiveAutoInstructionInProgress {
					inactiveStart = &event.Parameters.Timestamp
				}

				set := bson.M{
					"v.inactive_start": inactiveStart,
					"not_acked_since":  event.Parameters.Timestamp,
				}

				snoozeVal := resolveSnoozeAfterPbhLeave(event.Parameters.Timestamp, alarm)
				if snoozeVal > 0 {
					set["v.snooze.val"] = snoozeVal
				}

				update["$set"] = set
			}

			alarm = types.Alarm{}
			err = p.alarmCollection.FindOneAndUpdate(ctx, match, update, opts).Decode(&alarm)
			if err != nil {
				if errors.Is(err, mongodriver.ErrNoDocuments) {
					return nil
				}

				return err
			}
		}

		result.Entity, err = updateEntity(ctx,
			bson.M{
				"_id":               event.Entity.ID,
				"pbehavior_info.id": bson.M{"$nin": bson.A{nil, ""}},
			},
			bson.M{"$unset": bson.M{
				"pbehavior_info": "",
			}},
			p.entityCollection,
		)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		alarmChange.Type = types.AlarmChangeTypePbhLeave
		result.Forward = true
		result.Alarm = alarm
		result.AlarmChange = alarmChange

		updatedServiceStates, componentStateChanged, newComponentState, err = processComponentAndServiceCounters(
			ctx,
			p.entityServiceCountersCalculator,
			p.componentCountersCalculator,
			&result.Alarm,
			&result.Entity,
			result.AlarmChange,
		)

		return err
	})

	if err != nil || result.AlarmChange.Type == types.AlarmChangeTypeNone {
		return result, err
	}

	if result.Alarm.ID != "" {
		result.IsInstructionMatched = isInstructionMatched(event, result, p.autoInstructionMatcher, p.logger)
	}

	go p.postProcess(context.Background(), event, result, updatedServiceStates, componentStateChanged, newComponentState)

	return result, nil
}

func (p *pbhLeaveProcessor) postProcess(
	ctx context.Context,
	event rpc.AxeEvent,
	result Result,
	updatedServiceStates map[string]entitycounters.UpdatedServicesInfo,
	componentStateChanged bool,
	newComponentState int,
) {
	entity := *event.Entity
	if result.Entity.ID != "" {
		entity = result.Entity
	}

	p.metricsSender.SendEventMetrics(
		result.Alarm,
		entity,
		result.AlarmChange,
		event.Parameters.Timestamp.Time,
		event.Parameters.Initiator,
		event.Parameters.User,
		event.Parameters.Instruction,
		"",
	)

	for servID, servInfo := range updatedServiceStates {
		err := p.eventsSender.UpdateServiceState(ctx, servID, servInfo)
		if err != nil {
			p.logger.Err(err).Msg("failed to update service state")
		}
	}

	if componentStateChanged {
		err := p.eventsSender.UpdateComponentState(ctx, event.Entity.Component, event.Entity.Connector, newComponentState)
		if err != nil {
			p.logger.Err(err).Msg("failed to update component state")
		}
	}

	if result.Alarm.ID != "" {
		err := sendRemediationEvent(ctx, event, result, p.remediationRpcClient, p.encoder)
		if err != nil {
			p.logger.Err(err).Msg("cannot send event to engine-remediation")
		}
	}
}

func resolveSnoozeAfterPbhLeave(timestamp datetime.CpsTime, alarm types.Alarm) int64 {
	if alarm.Value.Snooze == nil || alarm.Value.Snooze.Initiator == types.InitiatorUser {
		return 0
	}

	steps := alarm.Value.Steps
	var snoozeDuration int64
	var snoozeElapsed int64
	var lastEnterTime int64

Loop:
	for i := len(steps) - 1; i >= 0; i-- {
		step := steps[i]
		switch step.Type {
		case types.AlarmStepSnooze:
			// this means, that snooze step is happened after pbh_enter step,
			// it's possible to do with a scenario feature, so if it happens,
			// then elapsed time = 0
			if lastEnterTime == 0 {
				snoozeElapsed = 0
			} else {
				snoozeElapsed += lastEnterTime - step.Timestamp.Unix()
			}

			snoozeDuration = int64(step.Value) - step.Timestamp.Unix()

			break Loop
		case types.AlarmStepPbhEnter:
			if step.PbehaviorCanonicalType != pbehavior.TypeActive {
				lastEnterTime = step.Timestamp.Unix()
			}
		case types.AlarmStepPbhLeave:
			if step.PbehaviorCanonicalType != pbehavior.TypeActive {
				snoozeElapsed += lastEnterTime - step.Timestamp.Unix()
			}
		}
	}

	return timestamp.Unix() + snoozeDuration - snoozeElapsed
}
