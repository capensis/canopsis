package event

import (
	"context"
	"errors"
	"fmt"
	"slices"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmtag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters/calculator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func newUpstreamHelper(
	dbClient mongo.DbClient,
	alarmConfigProvider config.AlarmConfigProvider,
	alarmStatusService alarmstatus.Service,
	pbhTypeResolver pbehavior.EntityTypeResolver,
	autoInstructionMatcher AutoInstructionMatcher,
	metaAlarmPostProcessor MetaAlarmPostProcessor,
	metricsSender metrics.Sender,
	remediationRpcClient engine.RPCClient,
	internalTagAlarmMatcher alarmtag.InternalTagAlarmMatcher,
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator,
	componentCountersCalculator calculator.ComponentCountersCalculator,
	eventsSender entitycounters.EventsSender,
	eventGenerator event.Generator,
	amqpPublisher libamqp.Publisher,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) *upstreamHelper {
	return &upstreamHelper{
		dbClient:                dbClient,
		alarmCollection:         dbClient.Collection(mongo.AlarmMongoCollection),
		entityCollection:        dbClient.Collection(mongo.EntityMongoCollection),
		pbehaviorCollection:     dbClient.Collection(mongo.PbehaviorMongoCollection),
		alarmConfigProvider:     alarmConfigProvider,
		alarmStatusService:      alarmStatusService,
		pbhTypeResolver:         pbhTypeResolver,
		autoInstructionMatcher:  autoInstructionMatcher,
		metaAlarmPostProcessor:  metaAlarmPostProcessor,
		metricsSender:           metricsSender,
		remediationRpcClient:    remediationRpcClient,
		internalTagAlarmMatcher: internalTagAlarmMatcher,
		eventGenerator:          eventGenerator,
		amqpPublisher:           amqpPublisher,
		encoder:                 encoder,
		logger:                  logger,
		countersHelper:          newCountersHelper(entityServiceCountersCalculator, componentCountersCalculator, eventsSender, logger),
	}
}

type upstreamHelper struct {
	dbClient                mongo.DbClient
	alarmCollection         mongo.DbCollection
	entityCollection        mongo.DbCollection
	pbehaviorCollection     mongo.DbCollection
	alarmConfigProvider     config.AlarmConfigProvider
	alarmStatusService      alarmstatus.Service
	pbhTypeResolver         pbehavior.EntityTypeResolver
	autoInstructionMatcher  AutoInstructionMatcher
	metaAlarmPostProcessor  MetaAlarmPostProcessor
	metricsSender           metrics.Sender
	remediationRpcClient    engine.RPCClient
	internalTagAlarmMatcher alarmtag.InternalTagAlarmMatcher
	eventGenerator          event.Generator
	amqpPublisher           libamqp.Publisher
	encoder                 encoding.Encoder
	logger                  zerolog.Logger
	countersHelper          *countersHelper
}

func (h *upstreamHelper) SendDownstreamEventsOnOK(ctx context.Context, entity types.Entity) error {
	switch entity.Type {
	case types.EntityTypeResource, types.EntityTypeComponent:
	default:
		return nil
	}

	cursor, err := h.entityCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"upstream": entity.ID,
			"enabled":  true,
		}},
		{"$lookup": bson.M{
			"from":         mongo.AlarmMongoCollection,
			"localField":   "_id",
			"foreignField": "d",
			"as":           "alarm",
			"pipeline": []bson.M{
				{"$match": bson.M{
					"v.resolved":   nil,
					"v.status.val": types.AlarmStatusUnknown,
				}},
				{"$limit": 1},
			},
		}},
		{"$unwind": "$alarm"},
		{"$project": bson.M{"alarm": 0}},
	})
	if err != nil {
		return fmt.Errorf("cannot find downstream entities: %w", err)
	}

	return h.sendEventsForDownstreams(ctx, cursor)
}

func (h *upstreamHelper) SendDownstreamEventsOnKO(ctx context.Context, entity types.Entity) error {
	switch entity.Type {
	case types.EntityTypeResource, types.EntityTypeComponent:
	default:
		return nil
	}

	cursor, err := h.entityCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"upstream": entity.ID,
			"enabled":  true,
		}},
		{"$lookup": bson.M{
			"from":         mongo.AlarmMongoCollection,
			"localField":   "_id",
			"foreignField": "d",
			"as":           "alarm",
			"pipeline": []bson.M{
				{"$match": bson.M{
					"v.resolved": nil,
				}},
				{"$limit": 1},
			},
		}},
		{"$unwind": bson.M{"path": "$alarm", "preserveNullAndEmptyArrays": true}},
		{"$match": bson.M{"$or": []bson.M{
			{"alarm": nil},
			{"alarm.v.status.val": bson.M{"$ne": types.AlarmStatusUnknown}},
		}}},
		{"$project": bson.M{"alarm": 0}},
	})
	if err != nil {
		return fmt.Errorf("cannot find downstream entities: %w", err)
	}

	return h.sendEventsForDownstreams(ctx, cursor)
}

func (h *upstreamHelper) Process(ctx context.Context, event rpc.AxeEvent, forceCountersUpdate bool) (Result, types.Alarm, error) {
	result := Result{}
	alarm := types.Alarm{}
	if event.Entity == nil || event.Entity.ID == "" || !event.Entity.Enabled {
		return result, alarm, nil
	}

	entity := *event.Entity
	countersRes := countersResult{}
	match := getOpenAlarmMatch(event)
	err := h.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		result = Result{}
		alarm = types.Alarm{}
		entity = *event.Entity
		countersRes = countersResult{}

		err := h.alarmCollection.FindOne(ctx, match).Decode(&alarm)
		if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
			return err
		}

		result, err = h.UpdateAlarm(ctx, event, alarm, entity)
		if err != nil {
			return err
		}

		if forceCountersUpdate || result.AlarmChange.Type != "" && result.AlarmChange.Type != types.AlarmChangeTypeUpdateStatus {
			result.IsCountersUpdated, countersRes, err = h.countersHelper.CalculateCounters(
				ctx,
				&result.Alarm,
				&entity,
				result.AlarmChange,
			)
		}

		return err
	})
	if err != nil || result.Alarm.ID == "" {
		return result, alarm, err
	}

	go h.PostProcess(context.WithoutCancel(ctx), event, result, countersRes)

	return result, alarm, nil
}

func (h *upstreamHelper) UpdateAlarm(ctx context.Context, event rpc.AxeEvent, alarm types.Alarm, entity types.Entity) (Result, error) {
	result := Result{}
	newStatus, statusRuleName, err := h.alarmStatusService.ComputeStatusOnStatusChange(ctx, alarm, entity)
	if err != nil {
		return result, fmt.Errorf("cannot compute alarm status: %w", err)
	}

	currentStatus := types.CpsNumber(types.AlarmStatusOff)
	if alarm.ID != "" {
		currentStatus = alarm.Value.Status.Value
	}

	if newStatus == currentStatus {
		if entity.IsUpstreamChanged {
			result.Entity, err = updateEntityByID(ctx, entity.ID, bson.M{"$unset": bson.M{"is_upstream_changed": ""}}, h.entityCollection)
			if err != nil {
				return result, err
			}
		}

		return result, nil
	}

	if alarm.ID == "" {
		if !h.shouldCreateAlarm(newStatus) {
			return result, nil
		}

		var v types.Entity
		err = h.entityCollection.FindOne(ctx, bson.M{"_id": entity.ID}).Decode(&v)
		if err != nil {
			return result, err
		}

		entity = v

		return h.createAlarm(ctx, entity, event, newStatus, statusRuleName)
	}

	if h.shouldCloseAlarm(alarm) {
		// close an alarm or update its status to a status which was shadowed by unknown status
		newState := types.CpsNumber(types.AlarmStateOK)
		if newStatus == types.AlarmStatusOngoing {
			newStatus = types.AlarmStatusOff
		}

		return h.UpdateAlarmStateAndStatus(ctx, alarm, entity, event, newState, newStatus, statusRuleName)
	}

	return h.UpdateAlarmStatus(ctx, alarm, entity, event, newStatus, statusRuleName)
}

func (h *upstreamHelper) PostProcess(
	ctx context.Context,
	event rpc.AxeEvent,
	result Result,
	countersRes countersResult,
) {
	entity := *event.Entity
	if result.Entity.ID != "" {
		entity = result.Entity
	}

	h.metricsSender.SendEventMetrics(
		result.Alarm,
		entity,
		result.AlarmChange,
		event.Parameters.Timestamp.Time,
		event.Parameters.Initiator,
		event.Parameters.User,
		event.Parameters.Instruction,
		"",
	)

	h.countersHelper.UpdateStates(ctx, countersRes)

	err := h.metaAlarmPostProcessor.Process(ctx, event, rpc.AxeResultEvent{
		Alarm:               &result.Alarm,
		AlarmChangeType:     result.AlarmChange.Type,
		AddedExternalTags:   result.AddedExternalTags,
		RemovedExternalTags: result.RemovedExternalTags,
	})
	if err != nil {
		h.logger.Err(err).Msg("cannot process meta alarm")
	}

	err = sendRemediationEvent(ctx, event, result, h.remediationRpcClient, h.encoder)
	if err != nil {
		h.logger.Err(err).Msg("cannot send event to engine-remediation")
	}

	if result.AlarmChange.Type == types.AlarmChangeTypeCreateAndPbhEnter {
		err = updatePbehaviorLastAlarmDate(ctx, h.pbehaviorCollection, result.Alarm.Value.PbehaviorInfo.ID, result.Alarm.Value.PbehaviorInfo.Timestamp)
		if err != nil {
			h.logger.Err(err).Msg("cannot update pbehavior")
		}

		if !result.Alarm.Value.PbehaviorInfo.IsDefaultActive() {
			err = updatePbehaviorAlarmCount(ctx, h.pbehaviorCollection, result.Alarm.Value.PbehaviorInfo.ID, "")
			if err != nil {
				h.logger.Err(err).Msg("cannot update pbehavior")
			}
		}
	}

	switch result.AlarmChange.Type {
	case types.AlarmChangeTypeUpdateStatus:
		alarmStatus := result.Alarm.Value.Status.Value
		prevStatus := result.AlarmChange.PreviousStatus
		if alarmStatus == types.AlarmStatusOff {
			err = h.SendDownstreamEventsOnOK(ctx, entity)
			if err != nil {
				h.logger.Err(err).Msg("cannot send downstream events")
			}
		} else if prevStatus == types.AlarmStatusOff {
			err = h.SendDownstreamEventsOnKO(ctx, entity)
			if err != nil {
				h.logger.Err(err).Msg("cannot send downstream events")
			}
		}
	case types.AlarmChangeTypeCreate, types.AlarmChangeTypeCreateAndPbhEnter:
		err = h.SendDownstreamEventsOnKO(ctx, entity)
		if err != nil {
			h.logger.Err(err).Msg("cannot send downstream events")
		}
	case types.AlarmChangeTypeStateDecrease:
		alarmStatus := result.Alarm.Value.Status.Value
		alarmState := result.Alarm.Value.State.Value
		prevStatus := result.AlarmChange.PreviousStatus
		// if alarm state is ok and status is off/cancelled
		if prevStatus != alarmStatus && alarmStatus == types.AlarmStatusOff ||
			prevStatus == alarmStatus && alarmStatus == types.AlarmStatusCancelled && alarmState == types.AlarmStateOK {
			err = h.SendDownstreamEventsOnOK(ctx, entity)
			if err != nil {
				h.logger.Err(err).Msg("cannot send downstream events")
			}
		}
	}
}

// shouldCreateAlarm
// An alarm can only be created on an entity update if an upstream is changed and an alarm with unknown status should exist.
func (h *upstreamHelper) shouldCreateAlarm(newStatus types.CpsNumber) bool {
	return newStatus == types.AlarmStatusUnknown
}

// shouldCloseAlarm
// An alarm can only be closed on an entity update if an upstream is changed and
// an alarm with unknown status exists and never were updated.
func (h *upstreamHelper) shouldCloseAlarm(alarm types.Alarm) bool {
	prevStatus := alarm.Value.Status.Value

	return prevStatus == types.AlarmStatusUnknown &&
		alarm.Value.InitialStatus == types.AlarmStatusUnknown && // created by upstream change
		alarm.Value.LastEventDate == nil && // no check events
		!alarm.IsStateLocked()
}

func (h *upstreamHelper) UpdateAlarmStatus(
	ctx context.Context,
	alarm types.Alarm,
	entity types.Entity,
	event rpc.AxeEvent,
	newStatus types.CpsNumber,
	statusRuleName string,
) (Result, error) {
	result := Result{}
	alarmStepType := types.AlarmStepStatusIncrease
	prevStatus := alarm.Value.Status.Value
	if prevStatus > newStatus {
		alarmStepType = types.AlarmStepStatusDecrease
	}

	newStepStatusQuery := valStepUpdateQueryWithInPbhInterval(alarmStepType, newStatus, statusRuleName, event.Parameters)
	matchUpdate := getOpenAlarmMatchWithStepsLimit(event)
	update := []bson.M{
		{"$set": bson.M{
			"v.status":                            newStepStatusQuery,
			"v.state_changes_since_status_update": 0,
			"v.last_update_date":                  event.Parameters.Timestamp,
			"v.last_st_upd_dt":                    event.Parameters.Timestamp,
			"v.steps":                             addStepUpdateQuery(newStepStatusQuery),
		}},
	}
	if newStatus == types.AlarmStatusUnknown {
		update = append(update, bson.M{"$unset": bson.A{"v.canceled"}})
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	updatedAlarm := types.Alarm{}
	err := h.alarmCollection.FindOneAndUpdate(ctx, matchUpdate, update, opts).Decode(&updatedAlarm)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return result, nil
		}

		return result, err
	}

	if entity.IsUpstreamChanged {
		result.Entity, err = updateEntityByID(ctx, entity.ID, bson.M{"$unset": bson.M{"is_upstream_changed": ""}}, h.entityCollection)
		if err != nil {
			return result, err
		}
	}

	alarmChange := types.NewAlarmChange()
	alarmChange.PreviousStatus = prevStatus
	alarmChange.Type = types.AlarmChangeTypeUpdateStatus
	result.Forward = true
	result.Alarm = updatedAlarm
	result.AlarmChange = alarmChange

	return result, nil
}

func (h *upstreamHelper) createAlarm(
	ctx context.Context,
	entity types.Entity,
	event rpc.AxeEvent,
	newStatus types.CpsNumber,
	statusRuleName string,
) (Result, error) {
	now := datetime.NewCpsTime()
	result := Result{}
	alarmChange := types.NewAlarmChange()
	var pbehaviorInfo types.PbehaviorInfo
	updateEntityPbhInfo := false
	var err error
	if entity.PbehaviorInfo.IsDefaultActive() {
		updateEntityPbhInfo = true
		pbehaviorInfo, err = resolvePbehaviorInfo(ctx, entity, now, h.pbhTypeResolver)
		if err != nil {
			return result, err
		}
	} else {
		pbehaviorInfo = entity.PbehaviorInfo
		pbehaviorInfo.Timestamp = &now
		alarmChange.PreviousPbehaviorTypeID = entity.PbehaviorInfo.TypeID
		alarmChange.PreviousPbehaviorCannonicalType = entity.PbehaviorInfo.CanonicalType
		alarmChange.PreviousEntityPbehaviorTime = entity.PbehaviorInfo.Timestamp
	}

	alarmConfig := h.alarmConfigProvider.Get()
	alarm, err := h.newAlarm(event.Parameters, entity, now, alarmConfig)
	if err != nil {
		return result, err
	}

	stateStep := NewAlarmStep(types.AlarmStepStateIncrease, event.Parameters, false)
	stateStep.Value = types.AlarmStateForUnknown
	stateStep.Message = statusRuleName
	stateStep.Author = canopsis.DefaultEventAuthor
	stateStep.Initiator = types.InitiatorSystem
	alarm.Value.State = &stateStep
	alarm.Value.MaxState = stateStep.Value
	alarm.Value.InitialState = stateStep.Value
	err = alarm.Value.Steps.Add(stateStep)
	if err != nil {
		return result, fmt.Errorf("cannot add alarm steps: %w", err)
	}

	statusStep := NewAlarmStep(types.AlarmStepStatusIncrease, event.Parameters, false)
	statusStep.Value = newStatus
	statusStep.Message = statusRuleName
	statusStep.Author = canopsis.DefaultEventAuthor
	statusStep.Initiator = types.InitiatorSystem
	alarm.Value.Status = &statusStep
	alarm.Value.InitialStatus = statusStep.Value
	err = alarm.Value.Steps.Add(statusStep)
	if err != nil {
		return result, fmt.Errorf("cannot add alarm steps: %w", err)
	}

	alarm.Value.TotalStateChanges++
	if pbehaviorInfo.IsDefaultActive() {
		alarmChange.Type = types.AlarmChangeTypeCreate
		alarm.NotAckedSince = &alarm.Time
	} else {
		if pbehaviorInfo.IsActive() {
			alarm.NotAckedSince = &alarm.Time
		} else {
			alarm.Value.InactiveStart = &now
		}

		newStep := types.NewPbhAlarmStep(types.AlarmStepPbhEnter, *pbehaviorInfo.Timestamp, pbehaviorInfo.Author,
			pbehaviorInfo.GetStepMessage(), "", "", types.InitiatorSystem, pbehaviorInfo.CanonicalType,
			pbehaviorInfo.IconName, pbehaviorInfo.Color)
		alarm.Value.PbehaviorInfo = pbehaviorInfo
		err := alarm.Value.Steps.Add(newStep)
		if err != nil {
			return result, fmt.Errorf("cannot add alarm steps: %w", err)
		}

		alarmChange.Type = types.AlarmChangeTypeCreateAndPbhEnter
	}

	result.IsInstructionMatched, err = h.autoInstructionMatcher.Match(alarmChange.GetTriggers(), types.AlarmWithEntity{Alarm: alarm, Entity: entity})
	if err != nil {
		return result, err
	}

	if h.alarmConfigProvider.Get().ActivateAlarmAfterAutoRemediation {
		alarm.InactiveAutoInstructionInProgress = result.IsInstructionMatched
	}

	alarm.InternalTags = h.internalTagAlarmMatcher.Match(entity, alarm)
	alarm.InternalTagsUpdated = datetime.NewMicroTime()
	alarm.Tags = slices.Clone(alarm.InternalTags)

	_, err = h.alarmCollection.InsertOne(ctx, types.AlarmWithEntityField{
		Alarm:  alarm,
		Entity: entity,
	})
	if err != nil {
		return result, fmt.Errorf("cannot create alarm: %w", err)
	}

	setEntity, unsetEntity := bson.M{}, bson.M{}
	if alarmChange.Type == types.AlarmChangeTypeCreateAndPbhEnter && updateEntityPbhInfo {
		setEntity["pbehavior_info"] = alarm.Value.PbehaviorInfo
		setEntity["last_pbehavior_date"] = alarm.Value.PbehaviorInfo.Timestamp
	}

	if entity.IsUpstreamChanged {
		unsetEntity["is_upstream_changed"] = ""
	}

	if len(setEntity) > 0 || len(unsetEntity) > 0 {
		update := bson.M{}
		if len(setEntity) > 0 {
			update["$set"] = setEntity
		}

		if len(unsetEntity) > 0 {
			update["$unset"] = unsetEntity
		}

		result.Entity, err = updateEntityByID(ctx, entity.ID, update, h.entityCollection)
		if err != nil {
			return result, err
		}
	}

	result.Forward = true
	result.Alarm = alarm
	result.AlarmChange = alarmChange

	return result, nil
}

func (h *upstreamHelper) UpdateAlarmStateAndStatus(
	ctx context.Context,
	alarm types.Alarm,
	entity types.Entity,
	event rpc.AxeEvent,
	newState, newStatus types.CpsNumber,
	statusRuleName string,
) (Result, error) {
	result := Result{}
	alarmChange := types.NewAlarmChange()
	alarmChange.PreviousState = alarm.Value.State.Value
	alarmChange.PreviousStatusChange = alarm.Value.State.Timestamp
	alarmChange.PreviousStatus = alarm.Value.Status.Value
	stateStep := NewAlarmStep(types.AlarmStepStateIncrease, event.Parameters, !alarm.Value.PbehaviorInfo.IsDefaultActive())
	stateStep.Value = newState
	stateStep.Message = statusRuleName
	alarmChange.Type = types.AlarmChangeTypeStateIncrease
	set := bson.M{}
	if newState < alarmChange.PreviousState {
		alarmChange.Type = types.AlarmChangeTypeStateDecrease
		stateStep.Type = types.AlarmStepStateDecrease
	} else if alarm.Value.MaxState < newState {
		alarm.Value.MaxState = newState
		set["v.max_state"] = newState
	}

	alarm.Value.State = &stateStep
	err := alarm.Value.Steps.Add(stateStep)
	if err != nil {
		return result, fmt.Errorf("cannot add alarm steps: %w", err)
	}

	set["v.state"] = stateStep
	set["v.last_update_date"] = event.Parameters.Timestamp
	set["v.last_st_upd_dt"] = event.Parameters.Timestamp

	statusStep := NewAlarmStep(types.AlarmStepStatusIncrease, event.Parameters, !alarm.Value.PbehaviorInfo.IsDefaultActive())
	statusStep.Value = newStatus
	statusStep.Message = statusRuleName
	if newStatus < alarmChange.PreviousStatus {
		statusStep.Type = types.AlarmStepStatusDecrease
	}

	set["v.status"] = statusStep
	set["v.state_changes_since_status_update"] = 0
	push := bson.M{}
	if stateStep.Type != "" {
		push["v.steps"] = bson.M{"$each": bson.A{stateStep, statusStep}}
	} else {
		push["v.steps"] = statusStep
	}

	match := bson.M{
		"_id":        alarm.ID,
		"v.resolved": nil,
		"$expr":      bson.M{"$lt": bson.A{bson.M{"$size": "$v.steps"}, types.AlarmStepsHardLimit}},
	}
	newAlarm := types.Alarm{}
	err = h.alarmCollection.FindOneAndUpdate(ctx, match, bson.M{
		"$set":  set,
		"$push": push,
	}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&newAlarm)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return result, nil
		}

		return result, fmt.Errorf("cannot update alarm: %w", err)
	}

	if entity.IsUpstreamChanged {
		result.Entity, err = updateEntityByID(ctx, entity.ID, bson.M{"$unset": bson.M{"is_upstream_changed": ""}}, h.entityCollection)
		if err != nil {
			return result, err
		}
	}

	result.Forward = true
	result.Alarm = newAlarm
	result.AlarmChange = alarmChange
	result.IsInstructionMatched, err = h.autoInstructionMatcher.Match(alarmChange.GetTriggers(), types.AlarmWithEntity{Alarm: newAlarm, Entity: entity})

	return result, err
}

func (h *upstreamHelper) newAlarm(
	params rpc.AxeParameters,
	entity types.Entity,
	timestamp datetime.CpsTime,
	alarmConfig config.AlarmConfig,
) (types.Alarm, error) {
	alarm := types.Alarm{
		EntityID: entity.ID,
		ID:       utils.NewID(),
		Time:     timestamp,
		Value: types.AlarmValue{
			CreationDate:                timestamp,
			DisplayName:                 types.GenDisplayName(alarmConfig.DisplayNameScheme),
			InitialOutput:               params.Output,
			Output:                      params.Output,
			InitialLongOutput:           params.LongOutput,
			LongOutput:                  params.LongOutput,
			LongOutputHistory:           []string{params.LongOutput},
			LastUpdateDate:              params.Timestamp,
			LastStateOrStatusUpdateDate: params.Timestamp,
			Parents:                     []string{},
			Children:                    []string{},
			UnlinkedParents:             []string{},
			Infos:                       map[string]map[string]interface{}{},
		},
	}

	connector := canopsis.DefaultSystemAlarmConnector
	connectorName := canopsis.DefaultSystemAlarmConnector
	switch entity.Type {
	case types.EntityTypeResource:
		alarm.Value.Resource = entity.Name
		alarm.Value.Component = entity.Component
		alarm.Value.Connector = connector
		alarm.Value.ConnectorName = connectorName
	case types.EntityTypeComponent, types.EntityTypeService:
		alarm.Value.Component = entity.Name
		alarm.Value.Connector = connector
		alarm.Value.ConnectorName = connectorName
	default:
		return types.Alarm{}, fmt.Errorf("unknown entity type %q", entity.Type)
	}

	return alarm, nil
}

func (h *upstreamHelper) sendEventsForDownstreams(ctx context.Context, cursor mongo.Cursor) error {
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		d := types.Entity{}
		err := cursor.Decode(&d)
		if err != nil {
			return fmt.Errorf("cannot decode downstream entity: %w", err)
		}

		e, err := h.eventGenerator.Generate(d)
		if err != nil {
			return fmt.Errorf("cannot generate downstream event: %w", err)
		}

		e.EventType = types.EventTypeUpdateStatus
		e.Timestamp = datetime.NewCpsTime()
		body, err := h.encoder.Encode(e)
		if err != nil {
			return fmt.Errorf("cannot encode downstream event: %w", err)
		}

		err = h.amqpPublisher.PublishWithContext(
			ctx,
			canopsis.DefaultExchangeName,
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
			return fmt.Errorf("cannot send downstream event: %w", err)
		}
	}

	if err := cursor.Err(); err != nil {
		return fmt.Errorf("cannot fetch downstream entities: %w", err)
	}

	return nil
}
