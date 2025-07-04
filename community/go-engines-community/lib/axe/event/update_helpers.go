package event

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
		} else if alarm.Value.MaxState < state {
			alarm.Value.MaxState = state
		}

		alarm.Value.State = &newStep
		err := alarm.Value.Steps.Add(newStep)
		if err != nil {
			return nil, nil, err
		}

		alarm.Value.TotalStateChanges++
		alarm.Value.LastUpdateDate = timestamp
		alarm.Value.LastStateOrStatusUpdateDate = timestamp
	}

	newStatus, statusRuleName := service.ComputeStatusOnStateChange(*alarm, entity)
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
				"v.last_st_upd_dt":                    alarm.Value.LastStateOrStatusUpdateDate,
				"v.max_state":                         alarm.Value.MaxState,
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
	alarm.Value.LastStateOrStatusUpdateDate = timestamp

	set := bson.M{
		"v.status":                            alarm.Value.Status,
		"v.state_changes_since_status_update": alarm.Value.StateChangesSinceStatusUpdate,
		"v.last_update_date":                  alarm.Value.LastUpdateDate,
		"v.last_st_upd_dt":                    alarm.Value.LastStateOrStatusUpdateDate,
	}
	newSteps := bson.A{}
	if state != currentState {
		set["v.total_state_changes"] = alarm.Value.TotalStateChanges
		set["v.state"] = alarm.Value.State
		set["v.max_state"] = alarm.Value.MaxState
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

func getMetaAlarmExternalTags(
	filterByLabel []string,
	children []types.AlarmWithEntity,
	existedTags []string,
) []string {
	tagsMap := make(map[string]struct{})
	existedTagsMap := make(map[string]bool, len(existedTags))
	for _, tag := range existedTags {
		existedTagsMap[tag] = true
	}

	f := func(tag string) {
		if existedTagsMap[tag] {
			return
		}

		toCopy := len(filterByLabel) == 0
		for _, label := range filterByLabel {
			if tag == label || strings.HasPrefix(tag, label+":") {
				toCopy = true
				break
			}
		}

		if toCopy {
			tagsMap[tag] = struct{}{}
		}
	}

	for _, child := range children {
		for _, tag := range child.Alarm.ExternalTags {
			f(tag)
		}

		for _, tag := range child.Alarm.ImportTags {
			f(tag)
		}
	}

	tags := make([]string, 0, len(tagsMap))
	for tag := range tagsMap {
		tags = append(tags, tag)
	}

	return tags
}

func getMetaAlarmEntityInfos(
	infoNames []string,
	children []types.AlarmWithEntity,
	existedInfos map[string]types.Info,
) map[string]types.Info {
	if len(infoNames) == 0 {
		return nil
	}

	infos := make(map[string]types.Info)
	for _, child := range children {
		for _, infoName := range infoNames {
			if info, ok := child.Entity.Infos[infoName]; ok {
				if existedInfo, ok := existedInfos[infoName]; ok {
					if reflect.DeepEqual(existedInfo.Value, info.Value) {
						continue
					}
				}

				infos[infoName] = types.Info{
					Name:        infoName,
					Value:       info.Value,
					Description: info.Description,
				}
			}
		}
	}

	return infos
}
