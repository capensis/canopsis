package event

import (
	"context"
	"errors"
	"fmt"
	"strings"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters/calculator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	outputMetaAlarmNamePrefix   = "Display name: "
	outputMetaAlarmEntityPrefix = "Entity: "
	outputMetaAlarmPrefix       = "Meta alarm name: "
)

func NewAlarmStep(t string, params rpc.AxeParameters, inPbehaviorInterval bool) types.AlarmStep {
	return types.NewAlarmStep(t, params.Timestamp, params.Author, params.Output, params.User, params.Role,
		params.Initiator, inPbehaviorInterval)
}

func NewPbhAlarmStep(t string, params rpc.AxeParameters, pbehaviorInfo types.PbehaviorInfo) types.AlarmStep {
	return types.NewPbhAlarmStep(t, params.Timestamp, params.Author, params.Output, params.User, params.Role,
		params.Initiator, pbehaviorInfo.CanonicalType, pbehaviorInfo.IconName, pbehaviorInfo.Color)
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

func removeMetaAlarmState(
	ctx context.Context,
	metaAlarm types.Alarm,
	rule correlation.Rule,
	metaAlarmStatesService correlation.MetaAlarmStateService,
) error {
	if rule.IsManual() {
		return nil
	}

	stateID := rule.GetStateID(metaAlarm.Value.MetaValuePath)
	metaAlarmState, err := metaAlarmStatesService.GetMetaAlarmState(ctx, stateID)
	if err != nil {
		return fmt.Errorf("cannot get meta alarm state: %w", err)
	}

	if metaAlarmState.ID == "" {
		return nil
	}

	_, err = metaAlarmStatesService.ArchiveState(ctx, metaAlarmState)
	if err != nil {
		return fmt.Errorf("cannot archive meta alarm state: %w", err)
	}

	_, err = metaAlarmStatesService.DeleteState(ctx, stateID)
	if err != nil {
		return fmt.Errorf("cannot delete meta alarm state: %w", err)
	}

	return nil
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

func updatePbehaviorLastAlarmDate(
	ctx context.Context,
	pbehaviorCollection mongo.DbCollection,
	pbhId string,
	lastAlarmDate *datetime.CpsTime,
) error {
	_, err := pbehaviorCollection.UpdateOne(ctx, bson.M{"_id": pbhId}, bson.M{"$set": bson.M{
		"last_alarm_date": lastAlarmDate,
	}})

	return err
}

func updatePbehaviorAlarmCount(
	ctx context.Context,
	pbehaviorCollection mongo.DbCollection,
	pbhId, prevPbhId string,
) error {
	writeModels := make([]mongodriver.WriteModel, 0, 2)
	if pbhId != "" {
		writeModels = append(writeModels, mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{"_id": pbhId}).
			SetUpdate(bson.M{"$inc": bson.M{
				"alarm_count": 1,
			}}))
	}

	if prevPbhId != "" && pbhId != prevPbhId {
		writeModels = append(writeModels, mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{"_id": prevPbhId}).
			SetUpdate(bson.M{
				"$inc": bson.M{
					"alarm_count": -1,
				},
			}))
	}

	if len(writeModels) > 0 {
		_, err := pbehaviorCollection.BulkWrite(ctx, writeModels)
		if err != nil {
			return fmt.Errorf("cannot update pbehaviors: %w", err)
		}
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
) (bool, map[string]entitycounters.UpdatedServicesInfo, bool, int, error) {
	serviceCountersUpdated, updatedServiceStates, err := entityServiceCountersCalculator.CalculateCounters(ctx, entity, alarm, alarmChange)
	if err != nil {
		return false, nil, false, 0, err
	}

	var componentCountersUpdated bool
	var componentStateChanged bool
	var newComponentState int
	if entity.Type == types.EntityTypeResource {
		componentCountersUpdated, componentStateChanged, newComponentState, err = componentCountersCalculator.CalculateCounters(ctx, entity, alarm, alarmChange)
		if err != nil {
			return false, nil, false, 0, err
		}
	}

	return serviceCountersUpdated || componentCountersUpdated, updatedServiceStates, componentStateChanged, newComponentState, nil
}

func sendTriggerEvent(
	ctx context.Context,
	event rpc.AxeEvent,
	result Result,
	amqpPublisher libamqp.Publisher,
	encoder encoding.Encoder,
	eventGenerator event.Generator,
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

	triggerEvent, err := eventGenerator.Generate(*event.Entity)
	if err != nil {
		logger.Err(err).Msgf("cannot generate event")

		return
	}

	triggerEvent.EventType = types.EventTypeTrigger
	triggerEvent.AlarmChange = &result.AlarmChange
	triggerEvent.AlarmID = result.Alarm.ID
	body, err := encoder.Encode(triggerEvent)
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
	metaAlarmStatesService correlation.MetaAlarmStateService,
	dbClient mongo.DbClient,
	alarmCollection, entityCollection, resolvedCollection, metaAlarmRuleCollection mongo.DbCollection,
) (Result, map[string]entitycounters.UpdatedServicesInfo, string, bool, int, error) {
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

		beforeAlarm, err := updateAlarmToResolve(ctx, alarmCollection, match, event.Parameters)
		if err != nil || beforeAlarm.ID == "" {
			return err
		}

		if beforeAlarm.NotAckedMetricSendTime != nil {
			notAckedMetricType = beforeAlarm.NotAckedMetricType
		}

		entity, err := updateEntityOfResolvedAlarm(ctx, entityCollection, event.Entity.ID)
		if err != nil || entity.ID == "" {
			return err
		}

		alarm, err := copyAlarmToResolvedCollection(ctx, alarmCollection, resolvedCollection, beforeAlarm.ID)
		if err != nil || alarm.ID == "" {
			return err
		}

		alarmChange := types.NewAlarmChange()
		alarmChange.Type = types.AlarmChangeTypeResolve
		result.Forward = true
		result.Alarm = alarm
		result.Entity = entity
		result.AlarmChange = alarmChange

		result.IsCountersUpdated, updatedServiceStates, componentStateChanged, newComponentState, err = processComponentAndServiceCounters(
			ctx,
			entityServiceCountersCalculator,
			componentCountersCalculator,
			&result.Alarm,
			&entity,
			result.AlarmChange,
		)
		if err != nil {
			return err
		}

		return removeMetaAlarmStateOnResolve(ctx, metaAlarmRuleCollection, metaAlarmStatesService, result.Alarm)
	})
	if err != nil || result.Alarm.ID == "" {
		return result, nil, "", false, 0, err
	}

	return result, updatedServiceStates, notAckedMetricType, componentStateChanged, newComponentState, nil
}

func updateAlarmToResolve(ctx context.Context, alarmCollection mongo.DbCollection, match bson.M, params rpc.AxeParameters) (types.Alarm, error) {
	beforeAlarm := types.Alarm{}
	update := getResolveAlarmUpdate(datetime.NewCpsTime(), params)
	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.Before).
		SetProjection(bson.M{
			"not_acked_metric_type":      1,
			"not_acked_metric_send_time": 1,
		})
	err := alarmCollection.FindOneAndUpdate(ctx, match, update, opts).Decode(&beforeAlarm)
	if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
		return beforeAlarm, err
	}

	return beforeAlarm, nil
}

func copyAlarmToResolvedCollection(
	ctx context.Context,
	alarmCollection, resolvedCollection mongo.DbCollection,
	alarmID string,
) (types.Alarm, error) {
	// extend alarm struct with bookmarks to copy user's bookmarks to a resolved alarm document
	var alarm struct {
		types.Alarm `bson:"inline"`
		Bookmarks   []string `bson:"bookmarks"`
	}
	err := alarmCollection.FindOne(ctx, bson.M{"_id": alarmID}).Decode(&alarm)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return alarm.Alarm, nil
		}

		return alarm.Alarm, err
	}

	_, err = resolvedCollection.UpdateOne(
		ctx,
		bson.M{"_id": alarm.ID},
		bson.M{"$set": alarm},
		options.Update().SetUpsert(true),
	)

	return alarm.Alarm, err
}

func updateEntityOfResolvedAlarm(ctx context.Context, entityCollection mongo.DbCollection, entityID string) (types.Entity, error) {
	entity := types.Entity{}
	entityUpdate := getResolveEntityUpdate()
	err := entityCollection.FindOneAndUpdate(ctx, bson.M{"_id": entityID}, entityUpdate,
		options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&entity)
	if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
		return entity, err
	}

	return entity, nil
}

func removeMetaAlarmStateOnResolve(
	ctx context.Context,
	metaAlarmRuleCollection mongo.DbCollection,
	metaAlarmStatesService correlation.MetaAlarmStateService,
	alarm types.Alarm,
) error {
	if !alarm.IsMetaAlarm() {
		return nil
	}

	var rule correlation.Rule
	err := metaAlarmRuleCollection.FindOne(ctx, bson.M{"_id": alarm.Value.Meta}).Decode(&rule)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return fmt.Errorf("meta alarm rule %s not found", alarm.Value.Meta)
		}

		return fmt.Errorf("cannot fetch meta alarm rule: %w", err)
	}

	return removeMetaAlarmState(ctx, alarm, rule, metaAlarmStatesService)
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
	metaAlarmPostProcessor MetaAlarmPostProcessor,
	metricsSender metrics.Sender,
	remediationRpcClient engine.RPCClient,
	pbehaviorCollection mongo.DbCollection,
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

	err := metaAlarmPostProcessor.Process(ctx, event, rpc.AxeResultEvent{
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

	if !result.Alarm.Value.PbehaviorInfo.IsDefaultActive() {
		err = updatePbehaviorAlarmCount(ctx, pbehaviorCollection, "", result.Alarm.Value.PbehaviorInfo.ID)
		if err != nil {
			logger.Err(err).Msg("cannot update pbehavior")
		}
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
						{"$gt": bson.A{"$v.inactive_start", 0}},
						{"$or": []bson.M{
							{"$ne": bson.A{
								bson.M{"$ifNull": bson.A{"$v.snooze", nil}},
								nil,
							}},
							{"$not": bson.M{"$in": bson.A{
								bson.M{"$ifNull": bson.A{"$v.pbehavior_info.canonical_type", nil}},
								bson.A{nil, "", pbehavior.TypeActive},
							}}},
							{"$eq": bson.A{"$auto_instruction_in_progress", true}},
							{"$eq": bson.A{"$inactive_delay_meta_alarm_in_progress", true}},
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
					"if": bson.M{"$ne": bson.A{
						bson.M{"$ifNull": bson.A{"$v.snooze", nil}},
						nil,
					}},
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
					"if": bson.M{"$not": bson.M{"$in": bson.A{
						bson.M{"$ifNull": bson.A{"$v.pbehavior_info.canonical_type", nil}},
						bson.A{nil, "", pbehavior.TypeActive},
					}}},
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

func updateInactiveStart(
	ts datetime.CpsTime,
	withSnoozeCond bool,
	withPbhCond bool,
	withAutoInstructionCond bool,
	withMetaAlarmCond bool,
) bson.M {
	conds := make([]bson.M, 0)
	if withSnoozeCond {
		conds = append(conds, bson.M{"$eq": bson.A{
			bson.M{"$ifNull": bson.A{"$v.snooze", nil}},
			nil,
		}})
	}

	if withPbhCond {
		conds = append(conds, bson.M{"$in": bson.A{
			bson.M{"$ifNull": bson.A{"$v.pbehavior_info.canonical_type", nil}},
			bson.A{nil, "", pbehavior.TypeActive},
		}})
	}

	if withAutoInstructionCond {
		conds = append(conds, bson.M{"$ne": bson.A{"$auto_instruction_in_progress", true}})
	}

	if withMetaAlarmCond {
		conds = append(conds, bson.M{"$ne": bson.A{"$inactive_delay_meta_alarm_in_progress", true}})
	}

	return bson.M{"$cond": bson.M{
		"if":   bson.M{"$and": conds},
		"then": nil,
		"else": ts,
	}}
}

func getAlarmsWithEntityByMatch(ctx context.Context, alarmCollection mongo.DbCollection, match bson.M) ([]types.AlarmWithEntity, error) {
	var alarms []types.AlarmWithEntity

	cursor, err := alarmCollection.Aggregate(ctx, []bson.M{
		{"$match": match},
		{"$project": bson.M{
			"alarm": "$$ROOT",
			"_id":   0,
		}},
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "alarm.d",
			"foreignField": "_id",
			"as":           "entity",
		}},
		{"$unwind": "$entity"},
		{"$sort": bson.M{
			"alarm.v.last_update_date": 1,
		}},
	})
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &alarms)
	if err != nil {
		return nil, err
	}

	return alarms, err
}

func updateMetaAlarmState(
	alarm *types.Alarm,
	entity types.Entity,
	timestamp datetime.CpsTime,
	state types.CpsNumber,
	output string,
	service alarmstatus.Service,
) (bson.M, bson.M, error) {
	var currentState, currentStatus types.CpsNumber
	if alarm.Value.State != nil {
		currentState = alarm.Value.State.Value
		currentStatus = alarm.Value.Status.Value
	}

	author := canopsis.DefaultEventAuthor
	if state != currentState {
		// Event is an Ok, so the alarm should be resolved anyway
		if alarm.IsStateLocked() && state != types.AlarmStateOK {
			return nil, nil, nil
		}

		// Create new Step to keep track of the alarm history
		newStep := types.NewAlarmStep(types.AlarmStepStateIncrease, timestamp, author, output, "", "",
			types.InitiatorSystem, !entity.PbehaviorInfo.IsDefaultActive())
		newStep.Value = state

		if state < currentState {
			newStep.Type = types.AlarmStepStateDecrease
		}

		alarm.Value.State = &newStep
		err := alarm.Value.Steps.Add(newStep)
		if err != nil {
			return nil, nil, err
		}

		alarm.Value.TotalStateChanges++
		alarm.Value.LastUpdateDate = timestamp
	}

	newStatus, statusRuleName := service.ComputeStatus(*alarm, entity)
	statusStepMessage := ConcatOutputAndRuleName(output, statusRuleName)
	if newStatus == currentStatus {
		if state == currentState {
			return nil, nil, nil
		}

		alarm.Value.StateChangesSinceStatusUpdate++

		return bson.M{
				"v.state":                             alarm.Value.State,
				"v.state_changes_since_status_update": alarm.Value.StateChangesSinceStatusUpdate,
				"v.total_state_changes":               alarm.Value.TotalStateChanges,
				"v.last_update_date":                  alarm.Value.LastUpdateDate,
			},
			bson.M{"v.steps": alarm.Value.State},
			nil
	}

	// Create new Step to keep track of the alarm history
	newStepStatus := types.NewAlarmStep(types.AlarmStepStatusIncrease, timestamp, author, statusStepMessage, "", "",
		types.InitiatorSystem, !entity.PbehaviorInfo.IsDefaultActive())
	newStepStatus.Value = newStatus

	if newStatus < currentStatus {
		newStepStatus.Type = types.AlarmStepStatusDecrease
	}

	alarm.Value.Status = &newStepStatus
	err := alarm.Value.Steps.Add(newStepStatus)
	if err != nil {
		return nil, nil, err
	}

	alarm.Value.StateChangesSinceStatusUpdate = 0
	alarm.Value.LastUpdateDate = timestamp

	set := bson.M{
		"v.status":                            alarm.Value.Status,
		"v.state_changes_since_status_update": alarm.Value.StateChangesSinceStatusUpdate,
		"v.last_update_date":                  alarm.Value.LastUpdateDate,
	}
	newSteps := bson.A{}
	if state != currentState {
		set["v.total_state_changes"] = alarm.Value.TotalStateChanges
		set["v.state"] = alarm.Value.State
		newSteps = append(newSteps, alarm.Value.State)
	}

	newSteps = append(newSteps, alarm.Value.Status)

	return set, bson.M{"v.steps": bson.M{"$each": newSteps}}, nil
}

func getMetaAlarmChildStepMsg(
	rule correlation.Rule,
	metaAlarm types.Alarm,
	event rpc.AxeEvent,
) string {
	msgBuilder := strings.Builder{}
	if !rule.IsManual() {
		msgBuilder.WriteString(types.RuleNameRulePrefix)
		msgBuilder.WriteString(rule.Name)
		msgBuilder.WriteString(". ")
	}

	msgBuilder.WriteString(outputMetaAlarmNamePrefix)
	msgBuilder.WriteString(metaAlarm.Value.DisplayName)
	msgBuilder.WriteString(". ")
	msgBuilder.WriteString(outputMetaAlarmEntityPrefix)
	msgBuilder.WriteString(metaAlarm.EntityID)
	msgBuilder.WriteRune('.')

	if event.Parameters.Output != "" {
		msgBuilder.WriteRune(' ')
		msgBuilder.WriteString(types.OutputCommentPrefix)
		msgBuilder.WriteString(event.Parameters.Output)
		msgBuilder.WriteRune('.')
	}

	return msgBuilder.String()
}

func getMetaAlarmChildEventOutput(
	metaAlarm types.Alarm,
	msg string,
	initiator string,
	isTicket bool,
) string {
	outputBuilder := strings.Builder{}
	msgLen := len(msg)
	if msgLen == 0 {
		outputBuilder.WriteString(outputMetaAlarmPrefix)
		outputBuilder.WriteString(metaAlarm.Value.DisplayName)

		return outputBuilder.String()
	}

	outputBuilder.WriteString(msg)
	if initiator == types.InitiatorSystem || initiator == types.InitiatorUser && isTicket {
		if msg[msgLen-1] != '.' {
			outputBuilder.WriteRune('.')
		}

		outputBuilder.WriteRune(' ')
		outputBuilder.WriteString(outputMetaAlarmPrefix)
		outputBuilder.WriteString(metaAlarm.Value.DisplayName)
		outputBuilder.WriteRune('.')
	} else {
		outputBuilder.WriteString("\n")
		outputBuilder.WriteString(outputMetaAlarmPrefix)
		outputBuilder.WriteString(metaAlarm.Value.DisplayName)
	}

	return outputBuilder.String()
}

func executeMetaAlarmOutputTpl(templateExecutor template.Executor, data correlation.EventExtraInfosMeta) (string, error) {
	rule := data.Rule
	if rule.OutputTemplate == "" {
		return "", nil
	}

	res, err := templateExecutor.Execute(rule.OutputTemplate, data)
	if err != nil {
		return "", fmt.Errorf("unable to execute output template for metaalarm rule %s: %w", rule.ID, err)
	}

	return res, nil
}
