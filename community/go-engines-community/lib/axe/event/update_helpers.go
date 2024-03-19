package event

import (
	"context"
	"errors"
	"fmt"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
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
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewAlarmStep(t string, params rpc.AxeParameters, inPbehaviorInterval bool) types.AlarmStep {
	return types.NewAlarmStep(t, params.Timestamp, params.Author, params.Output, params.User, params.Role,
		params.Initiator, inPbehaviorInterval)
}

func resolvePbehaviorInfo(ctx context.Context, entity types.Entity, now datetime.CpsTime, pbhTypeResolver pbehavior.EntityTypeResolver) (types.PbehaviorInfo, error) {
	result, err := pbhTypeResolver.Resolve(ctx, entity, now.Time)
	if err != nil {
		return types.PbehaviorInfo{}, err
	}

	return pbehavior.NewPBehaviorInfo(now, result), nil
}

func sendRemediationEvent(
	ctx context.Context,
	event rpc.AxeEvent,
	result Result,
	remediationRpcClient engine.RPCClient,
	encoder encoding.Encoder,
) error {
	if remediationRpcClient == nil {
		return nil
	}

	switch result.AlarmChange.Type {
	case types.AlarmChangeTypeNone:
		if result.AlarmChange.EventsCount < types.MinimalEventsCountThreshold {
			return nil
		}
	case
		types.AlarmChangeTypeCreate,
		types.AlarmChangeTypeCreateAndPbhEnter,
		types.AlarmChangeTypeStateIncrease,
		types.AlarmChangeTypeStateDecrease,
		types.AlarmChangeTypeChangeState,
		types.AlarmChangeTypeUnsnooze,
		types.AlarmChangeTypeActivate,
		types.AlarmChangeTypePbhEnter,
		types.AlarmChangeTypePbhLeave,
		types.AlarmChangeTypePbhLeaveAndEnter,
		types.AlarmChangeTypeResolve:
	default:
		return nil
	}

	entity := event.Entity
	if result.Entity.ID != "" {
		entity = &result.Entity
	}

	body, err := encoder.Encode(rpc.RemediationEvent{
		Alarm:       &result.Alarm,
		Entity:      entity,
		AlarmChange: result.AlarmChange,
	})
	if err != nil {
		return fmt.Errorf("cannot encode remediation event: %w", err)
	}

	err = remediationRpcClient.Call(ctx, engine.RPCMessage{
		CorrelationID: result.Alarm.ID,
		Body:          body,
	})
	if err != nil {
		return fmt.Errorf("cannot send rpc call to remediation: %w", err)
	}

	return nil
}

func updatePbhLastAlarmDate(ctx context.Context, result Result, pbehaviorCollection mongo.DbCollection) error {
	switch result.AlarmChange.Type {
	case types.AlarmChangeTypeCreateAndPbhEnter,
		types.AlarmChangeTypePbhEnter,
		types.AlarmChangeTypePbhLeaveAndEnter:
	default:
		return nil
	}

	pbhId := ""
	var lastAlarmDate *datetime.CpsTime
	if result.Alarm.ID == "" {
		pbhId = result.Entity.PbehaviorInfo.ID
		lastAlarmDate = result.Entity.PbehaviorInfo.Timestamp
	} else {
		pbhId = result.Alarm.Value.PbehaviorInfo.ID
		lastAlarmDate = result.Alarm.Value.PbehaviorInfo.Timestamp
	}

	_, err := pbehaviorCollection.UpdateOne(ctx,
		bson.M{"_id": pbhId},
		bson.M{"$set": bson.M{"last_alarm_date": lastAlarmDate}},
	)
	if err != nil {
		return fmt.Errorf("cannot update pbehavior: %w", err)
	}

	return nil
}

func isInstructionMatched(event rpc.AxeEvent, result Result, autoInstructionMatcher AutoInstructionMatcher, logger zerolog.Logger) bool {
	triggers := result.AlarmChange.GetTriggers()
	if len(triggers) == 0 {
		return false
	}

	entity := *event.Entity
	if result.Entity.ID != "" {
		entity = result.Entity
	}

	matched, err := autoInstructionMatcher.Match(triggers, types.AlarmWithEntity{Alarm: result.Alarm, Entity: entity})
	if err != nil {
		logger.Err(err).Str("alarm", result.Alarm.ID).Msg("cannot match auto instructions")
		return false
	}

	return matched
}

func updateEntityByID(ctx context.Context, entityID string, update bson.M, entityCollection mongo.DbCollection) (types.Entity, error) {
	newEntity := types.Entity{}
	err := entityCollection.FindOneAndUpdate(ctx, bson.M{"_id": entityID}, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).
		Decode(&newEntity)
	if err != nil {
		return newEntity, fmt.Errorf("cannot update entity: %w", err)
	}

	return newEntity, nil
}

func updateEntity(ctx context.Context, match, update bson.M, entityCollection mongo.DbCollection) (types.Entity, error) {
	newEntity := types.Entity{}
	err := entityCollection.FindOneAndUpdate(ctx, match, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).
		Decode(&newEntity)
	if err != nil {
		return newEntity, fmt.Errorf("cannot update entity: %w", err)
	}

	return newEntity, nil
}

func processComponentAndServiceCounters(
	ctx context.Context,
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator,
	componentCountersCalculator calculator.ComponentCountersCalculator,
	alarm *types.Alarm,
	entity *types.Entity,
	alarmChange types.AlarmChange,
) (map[string]entitycounters.UpdatedServicesInfo, bool, int, error) {
	updatedServiceStates, err := entityServiceCountersCalculator.CalculateCounters(ctx, entity, alarm, alarmChange)
	if err != nil {
		return nil, false, 0, err
	}

	var componentStateChanged bool
	var newComponentState int

	if entity.Type == types.EntityTypeResource {
		componentStateChanged, newComponentState, err = componentCountersCalculator.CalculateCounters(ctx, entity, alarm, alarmChange)
		if err != nil {
			return nil, false, 0, err
		}
	}

	return updatedServiceStates, componentStateChanged, newComponentState, nil
}

func sendTriggerEvent(
	ctx context.Context,
	event rpc.AxeEvent,
	result Result,
	amqpPublisher libamqp.Publisher,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) {
	switch result.AlarmChange.Type {
	case types.AlarmChangeTypeAutoInstructionFail,
		types.AlarmChangeTypeAutoInstructionComplete,
		types.AlarmChangeTypeInstructionJobFail,
		types.AlarmChangeTypeInstructionJobComplete:
	case types.AlarmChangeTypeDeclareTicketWebhook:
		if !event.Parameters.EmitTrigger {
			return
		}
	default:
		return
	}

	body, err := encoder.Encode(types.Event{
		EventType:     types.EventTypeTrigger,
		Connector:     canopsis.AxeConnector,
		ConnectorName: canopsis.AxeConnector,
		Component:     result.Alarm.Value.Component,
		Resource:      result.Alarm.Value.Resource,
		SourceType:    event.Entity.Type,
		AlarmChange:   &result.AlarmChange,
		AlarmID:       result.Alarm.ID,
		Author:        canopsis.DefaultEventAuthor,
		Initiator:     types.InitiatorSystem,
	})
	if err != nil {
		logger.Err(err).Msgf("cannot encode event")
		return
	}

	err = amqpPublisher.PublishWithContext(
		ctx,
		"",
		canopsis.FIFOQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  canopsis.JsonContentType,
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		logger.Err(err).Msgf("cannot send trigger event")
		return
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

func processResolve(
	ctx context.Context,
	match bson.M,
	event rpc.AxeEvent,
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator,
	componentCountersCalculator calculator.ComponentCountersCalculator,
	dbClient mongo.DbClient,
	alarmCollection, entityCollection, resolvedCollection mongo.DbCollection,
) (Result, map[string]entitycounters.UpdatedServicesInfo, string, bool, int, error) {
	update := getResolveAlarmUpdate(datetime.NewCpsTime(), event.Parameters)
	result := Result{}
	var updatedServiceStates map[string]entitycounters.UpdatedServicesInfo
	notAckedMetricType := ""
	var componentStateChanged bool
	var newComponentState int

	err := dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}
		updatedServiceStates = nil
		notAckedMetricType = ""
		componentStateChanged = false
		newComponentState = 0

		beforeAlarm := types.Alarm{}
		opts := options.FindOneAndUpdate().
			SetReturnDocument(options.Before).
			SetProjection(bson.M{
				"not_acked_metric_type":      1,
				"not_acked_metric_send_time": 1,
			})
		err := alarmCollection.FindOneAndUpdate(ctx, match, update, opts).Decode(&beforeAlarm)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		if beforeAlarm.NotAckedMetricSendTime != nil {
			notAckedMetricType = beforeAlarm.NotAckedMetricType
		}

		// extend alarm struct with bookmarks to copy user's bookmarks to a resolved alarm document
		var alarm struct {
			types.Alarm `bson:"inline"`
			Bookmarks   []string `bson:"bookmarks"`
		}
		err = alarmCollection.FindOne(ctx, bson.M{"_id": beforeAlarm.ID}).Decode(&alarm)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}
			return err
		}

		entity := types.Entity{}
		entityUpdate := getResolveEntityUpdate()
		err = entityCollection.FindOneAndUpdate(ctx, bson.M{"_id": event.Entity.ID}, entityUpdate,
			options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&entity)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		alarmChange := types.NewAlarmChange()
		alarmChange.Type = types.AlarmChangeTypeResolve
		result.Forward = true
		result.Alarm = alarm.Alarm
		result.Entity = entity
		result.AlarmChange = alarmChange

		_, err = resolvedCollection.UpdateOne(
			ctx,
			bson.M{"_id": alarm.ID},
			bson.M{"$set": alarm},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			return err
		}

		updatedServiceStates, componentStateChanged, newComponentState, err = processComponentAndServiceCounters(
			ctx,
			entityServiceCountersCalculator,
			componentCountersCalculator,
			&result.Alarm,
			&entity,
			result.AlarmChange,
		)

		return err
	})
	if err != nil || result.Alarm.ID == "" {
		return result, nil, "", false, 0, err
	}

	return result, updatedServiceStates, notAckedMetricType, componentStateChanged, newComponentState, nil
}

func postProcessResolve(
	ctx context.Context,
	event rpc.AxeEvent,
	result Result,
	updatedServiceStates map[string]entitycounters.UpdatedServicesInfo,
	componentChanged bool,
	newComponentState int,
	notAckedMetricType string,
	eventsSender entitycounters.EventsSender,
	metaAlarmEventProcessor libalarm.MetaAlarmEventProcessor,
	metricsSender metrics.Sender,
	remediationRpcClient engine.RPCClient,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) {
	metricsSender.SendEventMetrics(
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
		err := eventsSender.UpdateServiceState(ctx, servID, servInfo)
		if err != nil {
			logger.Err(err).Msg("failed to update service state")
		}
	}

	if componentChanged {
		err := eventsSender.UpdateComponentState(ctx, event.Entity.Component, newComponentState)
		if err != nil {
			logger.Err(err).Msg("failed to update component state")
		}
	}

	err := metaAlarmEventProcessor.ProcessAxeRpc(ctx, event, rpc.AxeResultEvent{
		Alarm:           &result.Alarm,
		AlarmChangeType: result.AlarmChange.Type,
	})
	if err != nil {
		logger.Err(err).Msg("cannot process meta alarm")
	}

	err = sendRemediationEvent(ctx, event, result, remediationRpcClient, encoder)
	if err != nil {
		logger.Err(err).Msg("cannot send event to engine-remediation")
	}
}

func getResolveAlarmUpdate(t datetime.CpsTime, params rpc.AxeParameters) []bson.M {
	newStep := NewAlarmStep(types.AlarmStepResolve, params, false)
	newStep.Timestamp = t

	return []bson.M{
		{"$set": bson.M{
			"v.duration": bson.M{"$subtract": bson.A{
				t,
				"$t",
			}},
			"v.inactive_duration": bson.M{"$sum": bson.A{
				"$v.inactive_duration",
				bson.M{"$cond": bson.M{
					"if": bson.M{"$and": []bson.M{
						{"$ne": bson.A{"$v.inactive_start", nil}},
						{"$or": []bson.M{
							{"$ne": bson.A{"$v.snooze", nil}},
							{"$not": bson.M{"$in": bson.A{"$v.pbehavior_info.canonical_type", bson.A{nil, "", pbehavior.TypeActive}}}},
						}},
					}},
					"then": bson.M{"$subtract": bson.A{
						t,
						"$v.inactive_start",
					}},
					"else": 0,
				}},
			}},
		}},
		{"$set": bson.M{
			"v.resolved": t,
			"v.steps":    bson.M{"$concatArrays": bson.A{"$v.steps", bson.A{newStep}}},
			"v.current_state_duration": bson.M{"$subtract": bson.A{
				t,
				"$v.state.t",
			}},
			"v.active_duration": bson.M{"$subtract": bson.A{
				"$v.duration",
				"$v.inactive_duration",
			}},
			"v.snooze_duration": bson.M{"$sum": bson.A{
				"$v.snooze_duration",
				bson.M{"$cond": bson.M{
					"if": bson.M{"$ne": bson.A{"$v.snooze", nil}},
					"then": bson.M{"$subtract": bson.A{
						t,
						"$v.snooze.t",
					}},
					"else": 0,
				}},
			}},
			"v.pbh_inactive_duration": bson.M{"$sum": bson.A{
				"$v.pbh_inactive_duration",
				bson.M{"$cond": bson.M{
					"if": bson.M{"$not": bson.M{"$in": bson.A{"$v.pbehavior_info.canonical_type", bson.A{nil, "", pbehavior.TypeActive}}}},
					"then": bson.M{"$subtract": bson.A{
						t,
						"$v.pbehavior_info.timestamp",
					}},
					"else": 0,
				}},
			}},
		}},
		{"$unset": bson.A{
			"not_acked_metric_type",
			"not_acked_metric_send_time",
			"not_acked_since",
		}},
	}
}

func getResolveEntityUpdate() bson.M {
	return bson.M{"$unset": bson.M{
		"idle_since":           "",
		"last_idle_rule_apply": "",
	}}
}

func ConcatOutputAndRuleName(output, ruleName string) string {
	if ruleName != "" {
		if output != "" {
			output += "\n"
		}

		output += ruleName
	}

	return output
}
