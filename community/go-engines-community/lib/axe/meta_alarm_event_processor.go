package axe

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	libaxeevent "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/axe/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	libevent "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/sync/errgroup"
)

const workers = 10

const (
	outputMetaAlarmNamePrefix   = "Display name: "
	outputMetaAlarmEntityPrefix = "Entity: "
	outputMetaAlarmPrefix       = "Meta alarm name: "
)

func NewMetaAlarmEventProcessor(
	dbClient mongo.DbClient,
	adapter libalarm.Adapter,
	ruleAdapter correlation.RulesAdapter,
	alarmStatusService alarmstatus.Service,
	alarmConfigProvider config.AlarmConfigProvider,
	encoder encoding.Encoder,
	amqpPublisher libamqp.Publisher,
	metricsSender metrics.Sender,
	metaAlarmStatesService correlation.MetaAlarmStateService,
	templateExecutor template.Executor,
	eventGenerator libevent.Generator,
	logger zerolog.Logger,
) libalarm.MetaAlarmEventProcessor {
	return &metaAlarmEventProcessor{
		dbClient:               dbClient,
		alarmCollection:        dbClient.Collection(mongo.AlarmMongoCollection),
		metaAlarmStatesService: metaAlarmStatesService,
		adapter:                adapter,
		ruleAdapter:            ruleAdapter,
		alarmStatusService:     alarmStatusService,
		alarmConfigProvider:    alarmConfigProvider,
		encoder:                encoder,
		amqpPublisher:          amqpPublisher,
		metricsSender:          metricsSender,
		templateExecutor:       templateExecutor,
		eventGenerator:         eventGenerator,
		logger:                 logger,
	}
}

type metaAlarmEventProcessor struct {
	dbClient               mongo.DbClient
	alarmCollection        mongo.DbCollection
	metaAlarmStatesService correlation.MetaAlarmStateService
	adapter                libalarm.Adapter
	ruleAdapter            correlation.RulesAdapter

	alarmStatusService  alarmstatus.Service
	alarmConfigProvider config.AlarmConfigProvider

	encoder       encoding.Encoder
	amqpPublisher libamqp.Publisher

	metricsSender metrics.Sender

	eventGenerator libevent.Generator

	logger zerolog.Logger

	templateExecutor template.Executor
}

func (p *metaAlarmEventProcessor) CreateMetaAlarm(
	ctx context.Context,
	event rpc.AxeEvent,
) (*types.Alarm, []types.Alarm, error) {
	if event.Entity == nil {
		return nil, nil, nil
	}

	var updatedChildrenAlarms []types.Alarm
	var metaAlarm types.Alarm
	var activateChildEvents []types.Event
	err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		updatedChildrenAlarms = updatedChildrenAlarms[:0]
		metaAlarm = types.Alarm{}
		activateChildEvents = activateChildEvents[:0]
		rule, err := p.ruleAdapter.GetRule(ctx, event.Parameters.MetaAlarmRuleID)
		if err != nil {
			return fmt.Errorf("cannot fetch meta alarm rule id=%q: %w", event.Parameters.MetaAlarmRuleID, err)
		} else if rule.ID == "" {
			return fmt.Errorf("meta alarm rule id=%q not found", event.Parameters.MetaAlarmRuleID)
		}

		metaAlarm = p.newMetaAlarm(event.Parameters, *event.Entity, p.alarmConfigProvider.Get())
		metaAlarm.Value.Meta = event.Parameters.MetaAlarmRuleID
		metaAlarm.Value.MetaValuePath = event.Parameters.MetaAlarmValuePath
		metaAlarm.Value.LastEventDate = datetime.CpsTime{} // should be empty

		if event.Parameters.DisplayName != "" {
			metaAlarm.Value.DisplayName = event.Parameters.DisplayName
		}

		stateID := rule.ID
		if event.Parameters.MetaAlarmValuePath != "" {
			stateID = fmt.Sprintf("%s&&%s", rule.ID, event.Parameters.MetaAlarmValuePath)
		}

		var childEntityIDs []string
		var archived bool

		if rule.IsManual() {
			childEntityIDs = event.Parameters.MetaAlarmChildren
		} else {
			metaAlarmState, err := p.metaAlarmStatesService.GetMetaAlarmState(ctx, stateID)
			if err != nil {
				return err
			}

			if metaAlarmState.MetaAlarmName != event.Entity.Name {
				// try to get archived state
				metaAlarmState, err = p.metaAlarmStatesService.GetMetaAlarmState(ctx, stateID+"-"+event.Entity.Name)
				if err != nil {
					return err
				}

				if metaAlarmState.ID == "" {
					return fmt.Errorf("meta alarm state for rule id=%q and meta alarm name=%q not found",
						event.Parameters.MetaAlarmRuleID, event.Entity.Name)
				}

				archived = true
			}

			childEntityIDs = metaAlarmState.ChildrenEntityIDs
		}

		var lastChild types.AlarmWithEntity
		worstState := types.CpsNumber(types.AlarmStateMinor)
		eventsCount := types.CpsNumber(0)
		var writeModels []mongodriver.WriteModel

		if len(childEntityIDs) > 0 {
			childAlarms, err := p.getAlarmsWithEntityByEntityIds(ctx, childEntityIDs)
			if err != nil {
				return fmt.Errorf("cannot fetch children alarms: %w", err)
			}

			if len(childAlarms) > 0 {
				lastChild = childAlarms[len(childAlarms)-1]
			}

			writeModels = make([]mongodriver.WriteModel, 0, len(childAlarms))
			for _, childAlarm := range childAlarms {
				if childAlarm.Alarm.Value.State != nil {
					childState := childAlarm.Alarm.Value.State.Value
					if childState > worstState {
						worstState = childState
					}
				}

				if !childAlarm.Alarm.HasParentByEID(metaAlarm.EntityID) {
					metaAlarm.AddChild(childAlarm.Alarm.EntityID)
					childAlarm.Alarm.AddParent(metaAlarm.EntityID)
					newStep := libaxeevent.NewAlarmStep(types.AlarmStepMetaAlarmAttach, event.Parameters, !childAlarm.Alarm.Value.PbehaviorInfo.IsDefaultActive())
					newStep.Message = p.getChildStepMsg(rule, metaAlarm, event)
					err := childAlarm.Alarm.Value.Steps.Add(newStep)
					if err != nil {
						return err
					}

					if childAlarm.Alarm.InactiveDelayMetaAlarmInProgress {
						activateChildEvent, err := p.eventGenerator.Generate(childAlarm.Entity)
						if err != nil {
							return err
						}

						activateChildEvent.EventType = types.EventTypeMetaAlarmChildActivate
						activateChildEvent.Timestamp = datetime.NewCpsTime()
						activateChildEvents = append(activateChildEvents, activateChildEvent)
					}

					writeModels = append(writeModels, mongodriver.NewUpdateOneModel().
						SetFilter(bson.M{
							"_id":       childAlarm.Alarm.ID,
							"v.parents": bson.M{"$ne": metaAlarm.EntityID},
						}).
						SetUpdate(bson.M{
							"$addToSet": bson.M{"v.parents": metaAlarm.EntityID},
							"$push":     bson.M{"v.steps": newStep},
						}))

					updatedChildrenAlarms = append(updatedChildrenAlarms, childAlarm.Alarm)
					eventsCount += childAlarm.Alarm.Value.EventsCount
					if metaAlarm.Value.LastEventDate.Before(childAlarm.Alarm.Value.LastEventDate) {
						metaAlarm.Value.LastEventDate = childAlarm.Alarm.Value.LastEventDate
					}
				}
			}
		}

		_, _, err = updateMetaAlarmState(&metaAlarm, *event.Entity, event.Parameters.Timestamp, worstState,
			event.Parameters.Output, p.alarmStatusService)
		if err != nil {
			return err
		}

		output := ""
		if rule.IsManual() {
			output = event.Parameters.Output
		} else {
			output, err = p.executeOutputTpl(correlation.EventExtraInfosMeta{
				Rule:     rule,
				Count:    int64(len(updatedChildrenAlarms)),
				Children: lastChild,
			})
			if err != nil {
				return err
			}
		}

		metaAlarm.Value.Output = output
		metaAlarm.Value.EventsCount = eventsCount
		writeModels = append(writeModels, mongodriver.NewInsertOneModel().SetDocument(metaAlarm))
		_, err = p.alarmCollection.BulkWrite(ctx, writeModels)
		if err != nil {
			return err
		}

		if !rule.IsManual() && !archived {
			ok, err := p.metaAlarmStatesService.SwitchStateToCreated(ctx, stateID)
			if err != nil || !ok {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	for _, e := range activateChildEvents {
		err = p.sendToFifo(ctx, e)
		if err != nil {
			return nil, nil, err
		}
	}

	return &metaAlarm, updatedChildrenAlarms, nil
}

func (p *metaAlarmEventProcessor) AttachChildrenToMetaAlarm(
	ctx context.Context,
	event rpc.AxeEvent,
) (*types.Alarm, []types.Alarm, []types.Event, error) {
	if len(event.Parameters.MetaAlarmChildren) == 0 {
		return nil, nil, nil, nil
	}

	var updatedChildrenAlarms []types.Alarm
	var updatedChildrenAlarmsWithEntity []types.AlarmWithEntity
	var activateChildEvents []types.Event
	var metaAlarm types.Alarm
	err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		updatedChildrenAlarms = updatedChildrenAlarms[:0]
		updatedChildrenAlarmsWithEntity = updatedChildrenAlarmsWithEntity[:0]
		activateChildEvents = activateChildEvents[:0]
		metaAlarm = types.Alarm{}
		err := p.alarmCollection.FindOne(ctx, bson.M{"d": event.Entity.ID}).Decode(&metaAlarm)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		rule, err := p.ruleAdapter.GetRule(ctx, event.Parameters.MetaAlarmRuleID)
		if err != nil {
			return fmt.Errorf("cannot fetch meta alarm rule id=%q: %w", event.Parameters.MetaAlarmRuleID, err)
		} else if rule.ID == "" {
			return fmt.Errorf("meta alarm rule id=%q not found", event.Parameters.MetaAlarmRuleID)
		}

		alarms, err := p.getAlarmsWithEntityByEntityIds(ctx, event.Parameters.MetaAlarmChildren)
		if err != nil {
			return err
		}

		lastEventDate := metaAlarm.Value.LastEventDate
		worstState := types.CpsNumber(types.AlarmStateOK)
		eventsCount := types.CpsNumber(0)
		writeModels := make([]mongodriver.WriteModel, 0, len(alarms))
		childrenIds := make([]string, 0, len(alarms))
		var lastChild types.AlarmWithEntity
		for _, childAlarm := range alarms {
			if !childAlarm.Alarm.AddParent(metaAlarm.EntityID) {
				continue
			}

			metaAlarm.AddChild(childAlarm.Entity.ID)
			childrenIds = append(childrenIds, childAlarm.Entity.ID)
			eventsCount += childAlarm.Alarm.Value.EventsCount
			if lastEventDate.Before(childAlarm.Alarm.Value.LastEventDate) {
				lastEventDate = childAlarm.Alarm.Value.LastEventDate
			}

			newStep := libaxeevent.NewAlarmStep(types.AlarmStepMetaAlarmAttach, event.Parameters, !childAlarm.Alarm.Value.PbehaviorInfo.IsDefaultActive())
			newStep.Message = p.getChildStepMsg(rule, metaAlarm, event)
			err = childAlarm.Alarm.Value.Steps.Add(newStep)
			if err != nil {
				return err
			}

			if childAlarm.Alarm.Value.State.Value > worstState {
				worstState = childAlarm.Alarm.Value.State.Value
			}

			if childAlarm.Alarm.InactiveDelayMetaAlarmInProgress {
				activateChildEvent, err := p.eventGenerator.Generate(childAlarm.Entity)
				if err != nil {
					return err
				}

				activateChildEvent.EventType = types.EventTypeMetaAlarmChildActivate
				activateChildEvent.Timestamp = datetime.NewCpsTime()
				activateChildEvents = append(activateChildEvents, activateChildEvent)
			}

			writeModels = append(writeModels, mongodriver.NewUpdateOneModel().
				SetFilter(bson.M{
					"_id":       childAlarm.Alarm.ID,
					"v.parents": bson.M{"$ne": metaAlarm.EntityID},
				}).
				SetUpdate(bson.M{
					"$addToSet": bson.M{"v.parents": metaAlarm.EntityID},
					"$push":     bson.M{"v.steps": newStep},
				}))

			updatedChildrenAlarms = append(updatedChildrenAlarms, childAlarm.Alarm)
			updatedChildrenAlarmsWithEntity = append(updatedChildrenAlarmsWithEntity, childAlarm)
			lastChild = childAlarm
		}

		if len(updatedChildrenAlarms) == 0 {
			return nil
		}

		var setUpdate, pushUpdate bson.M
		if worstState > metaAlarm.CurrentState() {
			setUpdate, pushUpdate, err = updateMetaAlarmState(&metaAlarm, *event.Entity, event.Parameters.Timestamp,
				worstState, metaAlarm.Value.Output, p.alarmStatusService)
			if err != nil {
				return err
			}
		}

		if setUpdate == nil {
			setUpdate = bson.M{}
		}

		if metaAlarm.Value.Meta == "" {
			metaAlarm.Value.Meta = event.Parameters.MetaAlarmRuleID
			setUpdate["v.meta"] = event.Parameters.MetaAlarmRuleID
			metaAlarm.Value.MetaValuePath = event.Parameters.MetaAlarmValuePath
			setUpdate["v.meta_value_path"] = event.Parameters.MetaAlarmValuePath
		}

		if metaAlarm.Value.LastEventDate.Unix() != lastEventDate.Unix() {
			metaAlarm.Value.LastEventDate = lastEventDate
			setUpdate["v.last_event_date"] = lastEventDate
		}

		childrenCount, err := p.adapter.GetCountOpenedAlarmsByIDs(ctx, metaAlarm.Value.Children)
		if err != nil {
			return err
		}

		output := ""
		if rule.Type == correlation.RuleTypeManualGroup {
			output = event.Parameters.Output
		} else {
			output, err = p.executeOutputTpl(correlation.EventExtraInfosMeta{
				Rule:     rule,
				Count:    childrenCount,
				Children: lastChild,
			})
			if err != nil {
				return err
			}
		}

		metaAlarm.Value.Output = output
		setUpdate["v.output"] = output

		update := bson.M{
			"$set":      setUpdate,
			"$inc":      bson.M{"v.events_count": eventsCount},
			"$addToSet": bson.M{"v.children": bson.M{"$each": childrenIds}},
		}
		if len(pushUpdate) > 0 {
			update["$push"] = pushUpdate
		}

		writeModels = append(writeModels, mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{"_id": metaAlarm.ID}).
			SetUpdate(update))
		_, err = p.alarmCollection.BulkWrite(ctx, writeModels)

		return err
	})
	if err != nil {
		return nil, nil, nil, err
	}

	if metaAlarm.ID == "" {
		return nil, nil, nil, nil
	}

	childrenEvents := p.applyActionsOnChildren(metaAlarm, updatedChildrenAlarmsWithEntity)
	for _, e := range activateChildEvents {
		err = p.sendToFifo(ctx, e)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	return &metaAlarm, updatedChildrenAlarms, childrenEvents, nil
}

func (p *metaAlarmEventProcessor) DetachChildrenFromMetaAlarm(
	ctx context.Context,
	event rpc.AxeEvent,
) (*types.Alarm, error) {
	if len(event.Parameters.MetaAlarmChildren) == 0 {
		return nil, nil
	}

	var updatedChildrenAlarms []types.Alarm
	var metaAlarm types.Alarm
	err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		updatedChildrenAlarms = updatedChildrenAlarms[:0]
		metaAlarm = types.Alarm{}
		err := p.alarmCollection.FindOne(ctx, bson.M{"d": event.Entity.ID}).Decode(&metaAlarm)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		rule, err := p.ruleAdapter.GetRule(ctx, event.Parameters.MetaAlarmRuleID)
		if err != nil {
			return fmt.Errorf("cannot fetch meta alarm rule id=%q: %w", event.Parameters.MetaAlarmRuleID, err)
		} else if rule.ID == "" {
			return fmt.Errorf("meta alarm rule id=%q not found", event.Parameters.MetaAlarmRuleID)
		}

		alarms, err := p.getAlarmsWithEntityByParentIdAndEntityIds(ctx, metaAlarm.EntityID,
			event.Parameters.MetaAlarmChildren)
		if err != nil {
			return err
		}

		eventsCount := types.CpsNumber(0)
		writeModels := make([]mongodriver.WriteModel, 0, len(alarms))
		childrenIds := make([]string, 0, len(alarms))
		for _, childAlarm := range alarms {
			if childAlarm.Alarm.RemoveParent(event.Entity.ID) {
				metaAlarm.RemoveChild(childAlarm.Entity.ID)
				childrenIds = append(childrenIds, childAlarm.Entity.ID)
				eventsCount -= childAlarm.Alarm.Value.EventsCount
				newStep := libaxeevent.NewAlarmStep(types.AlarmStepMetaAlarmDetach, event.Parameters, !childAlarm.Alarm.Value.PbehaviorInfo.IsDefaultActive())
				newStep.Message = p.getChildStepMsg(rule, metaAlarm, event)
				err = childAlarm.Alarm.Value.Steps.Add(newStep)
				if err != nil {
					return err
				}

				updatedChildrenAlarms = append(updatedChildrenAlarms, childAlarm.Alarm)
				writeModels = append(writeModels, mongodriver.NewUpdateOneModel().
					SetFilter(bson.M{
						"_id":       childAlarm.Alarm.ID,
						"v.parents": metaAlarm.EntityID,
					}).
					SetUpdate(bson.M{
						"$pull": bson.M{"v.parents": event.Entity.ID},
						"$push": bson.M{
							"v.unlinked_parents": event.Entity.ID,
							"v.steps":            newStep,
						},
					}))
			}
		}

		if len(updatedChildrenAlarms) == 0 {
			return nil
		}

		metaAlarmChildren, err := p.getAlarmsWithEntityByParentIdAndEntityIds(ctx, metaAlarm.EntityID,
			metaAlarm.Value.Children)
		if err != nil {
			return err
		}

		var lastEventDate datetime.CpsTime // should be empty
		worstState := types.CpsNumber(types.AlarmStateOK)

		for _, childAlarm := range metaAlarmChildren {
			if childAlarm.Alarm.Value.State.Value > worstState {
				worstState = childAlarm.Alarm.Value.State.Value
			}

			if lastEventDate.Before(childAlarm.Alarm.Value.LastEventDate) {
				lastEventDate = childAlarm.Alarm.Value.LastEventDate
			}
		}

		setUpdate, pushUpdate, err := updateMetaAlarmState(&metaAlarm, *event.Entity, event.Parameters.Timestamp,
			worstState, metaAlarm.Value.Output, p.alarmStatusService)
		if err != nil {
			return err
		}

		if setUpdate == nil {
			setUpdate = bson.M{}
		}

		metaAlarm.Value.LastEventDate = lastEventDate
		setUpdate["v.last_event_date"] = lastEventDate
		infos := correlation.EventExtraInfosMeta{
			Rule:  rule,
			Count: int64(len(metaAlarmChildren)),
		}
		if len(metaAlarmChildren) > 0 {
			infos.Children = metaAlarmChildren[len(metaAlarmChildren)-1]
		}

		output := ""
		if rule.IsManual() {
			output = event.Parameters.Output
		} else {
			output, err = p.executeOutputTpl(infos)
			if err != nil {
				return err
			}
		}

		metaAlarm.Value.Output = output
		setUpdate["v.output"] = output
		update := bson.M{
			"$set":  setUpdate,
			"$inc":  bson.M{"v.events_count": eventsCount},
			"$pull": bson.M{"v.children": bson.M{"$in": childrenIds}},
		}
		if len(pushUpdate) > 0 {
			update["$push"] = pushUpdate
		}

		writeModels = append(writeModels, mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{"_id": metaAlarm.ID}).
			SetUpdate(update))
		_, err = p.alarmCollection.BulkWrite(ctx, writeModels)

		return err
	})
	if err != nil {
		return nil, err
	}

	return &metaAlarm, nil
}

func (p *metaAlarmEventProcessor) applyActionsOnChildren(
	metaAlarm types.Alarm,
	childrenAlarms []types.AlarmWithEntity,
) []types.Event {
	var events []types.Event
	steps := metaAlarm.GetAppliedActions()
	for _, childAlarm := range childrenAlarms {
		for _, step := range steps {
			childEvent, err := p.getChildEventByStep(metaAlarm, childAlarm, step)
			if err != nil {
				p.logger.Err(err).Str("entity", childAlarm.Entity.ID).Msg("cannot process child")
				continue
			}

			if childEvent.EventType != "" {
				events = append(events, childEvent)
			}
		}
	}

	return events
}

func (p *metaAlarmEventProcessor) ProcessAxeRpc(
	ctx context.Context,
	event rpc.AxeEvent,
	eventRes rpc.AxeResultEvent,
) error {
	if eventRes.Alarm == nil {
		return nil
	}

	alarm := eventRes.Alarm

	err := p.processComponentRpc(ctx, event, eventRes)
	if err != nil {
		return err
	}

	if alarm.IsMetaAlarm() {
		return p.processParentRpc(ctx, event, eventRes)
	}

	if alarm.IsMetaChild() {
		return p.processChildRpc(ctx, event, eventRes)
	}

	return nil
}

func (p *metaAlarmEventProcessor) processParentRpc(
	ctx context.Context,
	event rpc.AxeEvent,
	eventRes rpc.AxeResultEvent,
) error {
	if !applyOnChild(eventRes.AlarmChangeType) {
		return nil
	}

	var childAlarms []types.AlarmWithEntity
	err := p.adapter.GetOpenedAlarmsWithEntityByIDs(ctx, eventRes.Alarm.Value.Children, &childAlarms)
	if err != nil {
		return fmt.Errorf("cannot fetch children alarms: %w", err)
	}

	for _, childAlarm := range childAlarms {
		childEvent, err := p.getChildEventByMetaAlarmEvent(*eventRes.Alarm, childAlarm, event, eventRes)
		if err != nil {
			p.logger.Err(err).Str("entity", childAlarm.Entity.ID).Msg("cannot process child")
			continue
		}

		if childEvent.EventType == "" {
			continue
		}

		err = p.sendToFifo(ctx, childEvent)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *metaAlarmEventProcessor) getChildEventByStep(
	metaAlarm types.Alarm,
	childAlarm types.AlarmWithEntity,
	metaAlarmStep types.AlarmStep,
) (types.Event, error) {
	childEvent, err := p.eventGenerator.Generate(childAlarm.Entity)
	if err != nil {
		return childEvent, err
	}

	childEvent.Timestamp = datetime.NewCpsTime()
	childEvent.Author = metaAlarmStep.Author
	childEvent.UserID = metaAlarmStep.UserID
	childEvent.Initiator = metaAlarmStep.Initiator
	childEvent.Role = metaAlarmStep.Role
	isTicketStep := false
	switch metaAlarmStep.Type {
	case types.AlarmStepAck:
		childEvent.EventType = types.EventTypeAck
	case types.AlarmStepSnooze:
		childEvent.EventType = types.EventTypeSnooze
		childEvent.Duration = types.CpsNumber(int64(metaAlarmStep.Value) - childEvent.Timestamp.Unix())
		if childEvent.Duration <= 0 {
			return types.Event{}, nil
		}
	case types.AlarmStepAssocTicket:
		childEvent.EventType = types.EventTypeAssocTicket
		childEvent.TicketInfo = metaAlarmStep.TicketInfo
		childEvent.TicketInfo.TicketMetaAlarmID = metaAlarm.ID
		isTicketStep = true
	case types.AlarmStepDeclareTicket:
		childEvent.EventType = types.EventTypeDeclareTicketWebhook
		childEvent.TicketInfo = metaAlarmStep.TicketInfo
		childEvent.TicketInfo.TicketMetaAlarmID = metaAlarm.ID
		isTicketStep = true
	case types.AlarmStepComment:
		childEvent.EventType = types.EventTypeComment
	default:
		return types.Event{}, nil
	}

	childEvent.Output = p.getChildEventOutput(metaAlarm, metaAlarmStep.Message, metaAlarmStep.Initiator, isTicketStep)

	return childEvent, nil
}

func (p *metaAlarmEventProcessor) getChildEventByMetaAlarmEvent(
	metaAlarm types.Alarm,
	childAlarm types.AlarmWithEntity,
	event rpc.AxeEvent,
	eventRes rpc.AxeResultEvent,
) (types.Event, error) {
	childEvent, err := p.eventGenerator.Generate(childAlarm.Entity)
	if err != nil {
		return childEvent, err
	}

	childEvent.EventType = event.EventType
	childEvent.Timestamp = datetime.NewCpsTime()
	childEvent.Author = event.Parameters.Author
	childEvent.UserID = event.Parameters.User
	childEvent.Initiator = event.Parameters.Initiator
	childEvent.TicketInfo = event.Parameters.TicketInfo
	childEvent.TicketInfo.TicketMetaAlarmID = eventRes.Alarm.ID
	output := event.Parameters.Output
	isTicket := false
	switch eventRes.AlarmChangeType {
	case types.AlarmChangeTypeDeclareTicketWebhook,
		types.AlarmChangeTypeAutoDeclareTicketWebhook:
		childEvent.EventType = types.EventTypeDeclareTicketWebhook
		output = event.Parameters.TicketInfo.GetStepMessage()
		isTicket = true
	case types.AlarmChangeTypeAssocTicket:
		isTicket = true
	}

	if event.Parameters.State != nil {
		childEvent.State = *event.Parameters.State
	}

	if event.Parameters.Duration != nil {
		seconds, err := event.Parameters.Duration.To("s")
		if err == nil {
			childEvent.Duration = types.CpsNumber(seconds.Value)
		}
	}

	childEvent.Output = p.getChildEventOutput(metaAlarm, output, event.Parameters.Initiator, isTicket)

	return childEvent, nil
}

func (p *metaAlarmEventProcessor) getChildEventOutput(
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

func (p *metaAlarmEventProcessor) processComponentRpc(
	ctx context.Context,
	event rpc.AxeEvent,
	eventRes rpc.AxeResultEvent,
) error {
	if !event.Parameters.TicketResources ||
		eventRes.AlarmChangeType != types.AlarmChangeTypeDeclareTicketWebhook ||
		event.Entity.Type != types.EntityTypeComponent {
		return nil
	}

	componentAlarm := eventRes.Alarm
	if componentAlarm == nil || componentAlarm.Value.Ticket == nil {
		return nil
	}

	resources, err := p.adapter.GetAlarmsWithoutTicketByComponent(ctx, eventRes.Alarm.Value.Component)
	if err != nil {
		return fmt.Errorf("cannot fetch alarms: %w", err)
	}

	outputBuilder := strings.Builder{}
	outputBuilder.WriteString(componentAlarm.Value.Ticket.Message)
	outputBuilder.WriteString(" ")
	outputBuilder.WriteString(types.OutputComponentPrefix)
	outputBuilder.WriteString(componentAlarm.EntityID)
	outputBuilder.WriteRune('.')
	output := outputBuilder.String()
	for _, resource := range resources {
		if resource.Entity.Type != types.EntityTypeResource {
			continue
		}

		resourceEvent, err := p.eventGenerator.Generate(resource.Entity)
		if err != nil {
			p.logger.Err(err).Str("entity", resource.Entity.ID).Msg("cannot generate event")
			continue
		}

		resourceEvent.EventType = types.EventTypeDeclareTicketWebhook
		resourceEvent.Timestamp = datetime.NewCpsTime()
		resourceEvent.Output = output
		resourceEvent.TicketInfo = componentAlarm.Value.Ticket.TicketInfo
		resourceEvent.Author = event.Parameters.Author
		resourceEvent.UserID = event.Parameters.User
		resourceEvent.Initiator = event.Parameters.Initiator
		err = p.sendToFifo(ctx, resourceEvent)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *metaAlarmEventProcessor) processChildRpc(
	ctx context.Context, event rpc.AxeEvent, eventRes rpc.AxeResultEvent) error {
	switch eventRes.AlarmChangeType {
	case types.AlarmChangeTypeStateIncrease, types.AlarmChangeTypeStateDecrease, types.AlarmChangeTypeChangeState:
		err := p.updateParentState(ctx, *eventRes.Alarm)
		if err != nil {
			return err
		}
	case types.AlarmChangeTypeResolve:
		err := p.incrementParentEventsCount(ctx, eventRes.Alarm.Value.Parents, -eventRes.Alarm.Value.EventsCount)
		if err != nil {
			return fmt.Errorf("cannot update parent alarms: %w", err)
		}

		err = p.resolveParents(ctx, *eventRes.Alarm, event.Parameters.Timestamp)
		if err != nil {
			return err
		}

		if eventRes.Alarm.Value.State.Value != types.AlarmStateOK {
			err := p.updateParentState(ctx, *eventRes.Alarm)
			if err != nil {
				return err
			}
		}
	}

	if event.EventType == types.EventTypeCheck {
		err := p.incrementParentEventsCount(ctx, eventRes.Alarm.Value.Parents, 1)
		if err != nil {
			return fmt.Errorf("cannot update parent alarms: %w", err)
		}

		err = p.adapter.UpdateLastEventDate(ctx, eventRes.Alarm.Value.Parents, event.Parameters.Timestamp)
		if err != nil {
			return fmt.Errorf("cannot update parent alarms: %w", err)
		}
	}

	return nil
}

func (p *metaAlarmEventProcessor) incrementParentEventsCount(
	ctx context.Context,
	parentIDs []string,
	count types.CpsNumber,
) error {
	_, err := p.alarmCollection.UpdateMany(
		ctx,
		bson.M{
			"d":          bson.M{"$in": parentIDs},
			"v.resolved": nil,
		},
		bson.M{"$inc": bson.M{"v.events_count": count}},
	)

	return err
}

func (p *metaAlarmEventProcessor) resolveParents(
	ctx context.Context,
	childAlarm types.Alarm,
	timestamp datetime.CpsTime,
) error {
	ch := make(chan string)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		defer close(ch)
		for _, p := range childAlarm.Value.Parents {
			select {
			case <-ctx.Done():
				return nil
			case ch <- p:
			}
		}

		return nil
	})

	w := int(math.Min(float64(workers), float64(len(childAlarm.Value.Parents))))
	for i := 0; i < w; i++ {
		g.Go(func() error {
			for parentId := range ch {
				var parentAlarm types.AlarmWithEntity
				err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
					parentAlarm = types.AlarmWithEntity{}
					alarms := make([]types.AlarmWithEntity, 0)
					err := p.adapter.GetOpenedAlarmsWithEntityByIDs(ctx, []string{parentId}, &alarms)
					if err != nil {
						return fmt.Errorf("cannot fetch parent: %w", err)
					}
					if len(alarms) == 0 {
						return nil
					}
					parentAlarm = alarms[0]

					rule, err := p.ruleAdapter.GetRule(ctx, parentAlarm.Alarm.Value.Meta)
					if err != nil {
						return fmt.Errorf("cannot fetch meta alarm rule: %w", err)
					}
					if rule.ID == "" {
						return fmt.Errorf("meta alarm rule %s not found", parentAlarm.Alarm.Value.Meta)
					}
					if !rule.AutoResolve {
						return nil
					}

					resolvedCount, err := p.adapter.CountResolvedAlarm(ctx, parentAlarm.Alarm.Value.Children)
					if err != nil {
						return fmt.Errorf("cannot fetch alarms: %w", err)
					}

					if resolvedCount < len(parentAlarm.Alarm.Value.Children) {
						return nil
					}

					update := resolveMetaAlarm(&parentAlarm.Alarm, timestamp)
					_, err = p.alarmCollection.UpdateOne(ctx, bson.M{"_id": parentAlarm.Alarm.ID}, update)
					if err != nil {
						return fmt.Errorf("cannot update alarm: %w", err)
					}

					return nil
				})
				if err != nil {
					return err
				}

				if parentAlarm.Alarm.IsResolved() {
					err = p.adapter.CopyAlarmToResolvedCollection(ctx, parentAlarm.Alarm)
					if err != nil {
						return fmt.Errorf("cannot update alarm: %w", err)
					}

					p.metricsSender.SendResolve(parentAlarm.Alarm, parentAlarm.Entity, timestamp.Time)
				}
			}

			return nil
		})
	}

	return g.Wait()
}

func (p *metaAlarmEventProcessor) updateParentState(ctx context.Context, childAlarm types.Alarm) error {
	ch := make(chan string)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		defer close(ch)
		for _, p := range childAlarm.Value.Parents {
			select {
			case <-ctx.Done():
				return nil
			case ch <- p:
			}
		}

		return nil
	})

	w := int(math.Min(float64(workers), float64(len(childAlarm.Value.Parents))))
	for i := 0; i < w; i++ {
		g.Go(func() error {
			for parentId := range ch {
				err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
					alarms := make([]types.AlarmWithEntity, 0)
					err := p.adapter.GetOpenedAlarmsWithEntityByIDs(ctx, []string{parentId}, &alarms)
					if err != nil {
						return fmt.Errorf("cannot fetch parent: %w", err)
					}
					if len(alarms) == 0 {
						return nil
					}
					parentAlarm := alarms[0]

					rule, err := p.ruleAdapter.GetRule(ctx, parentAlarm.Alarm.Value.Meta)
					if err != nil {
						return fmt.Errorf("cannot fetch meta alarm rule: %w", err)
					}
					if rule.ID == "" {
						return fmt.Errorf("meta alarm rule %s not found", parentAlarm.Alarm.Value.Meta)
					}

					parentState := parentAlarm.Alarm.Value.State.Value
					childState := childAlarm.Value.State.Value
					if childAlarm.IsResolved() {
						childState = types.AlarmStateOK
					}

					var newState types.CpsNumber
					var newLastEventDate datetime.CpsTime
					if childState > parentState {
						newState = childState
					} else if childState < parentState {
						state, lastEventDate, err := p.adapter.GetWorstAlarmStateAndMaxLastEventDate(ctx,
							parentAlarm.Alarm.Value.Children)
						if err != nil {
							return fmt.Errorf("cannot fetch children state: %w", err)
						}

						newState = types.CpsNumber(state)
						newLastEventDate = datetime.NewCpsTime(lastEventDate)
					} else {
						return nil
					}

					setUpdate, pushUpdate, err := updateMetaAlarmState(&parentAlarm.Alarm, parentAlarm.Entity, childAlarm.Value.LastUpdateDate,
						newState, parentAlarm.Alarm.Value.Output, p.alarmStatusService)
					if err != nil {
						return fmt.Errorf("cannot update parent: %w", err)
					}

					if setUpdate == nil {
						setUpdate = bson.M{}
					}

					parentAlarm.Alarm.Value.LastEventDate = newLastEventDate
					setUpdate["v.last_event_date"] = newLastEventDate
					update := bson.M{
						"$set": setUpdate,
					}
					if len(pushUpdate) > 0 {
						update["$push"] = pushUpdate
					}

					_, err = p.alarmCollection.UpdateOne(ctx, bson.M{"_id": parentAlarm.Alarm.ID}, update)
					if err != nil {
						return fmt.Errorf("cannot update alarm: %w", err)
					}

					return nil
				})
				if err != nil {
					return err
				}
			}

			return nil
		})
	}

	return g.Wait()
}

func (p *metaAlarmEventProcessor) sendToFifo(ctx context.Context, event types.Event) error {
	body, err := p.encoder.Encode(event)
	if err != nil {
		return fmt.Errorf("cannot encode event: %w", err)
	}

	err = p.amqpPublisher.PublishWithContext(
		ctx,
		canopsis.FIFOExchangeName,
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
		return fmt.Errorf("cannot send child event: %w", err)
	}

	return nil
}

func (p *metaAlarmEventProcessor) getAlarmsWithEntityByEntityIds(
	ctx context.Context,
	entityIDs []string,
) ([]types.AlarmWithEntity, error) {
	return p.getAlarmsWithEntityByMatch(ctx, bson.M{
		"d":          bson.M{"$in": entityIDs},
		"v.resolved": nil,
	})
}

func (p *metaAlarmEventProcessor) getAlarmsWithEntityByParentIdAndEntityIds(
	ctx context.Context,
	parentId string,
	entityIDs []string,
) ([]types.AlarmWithEntity, error) {
	return p.getAlarmsWithEntityByMatch(ctx, bson.M{
		"v.parents":  parentId,
		"d":          bson.M{"$in": entityIDs},
		"v.resolved": nil,
	})
}

func (p *metaAlarmEventProcessor) getAlarmsWithEntityByMatch(
	ctx context.Context,
	match bson.M,
) ([]types.AlarmWithEntity, error) {
	var alarms []types.AlarmWithEntity

	cursor, err := p.alarmCollection.Aggregate(ctx, []bson.M{
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

func (p *metaAlarmEventProcessor) executeOutputTpl(data correlation.EventExtraInfosMeta) (string, error) {
	rule := data.Rule
	if rule.OutputTemplate == "" {
		return "", nil
	}

	res, err := p.templateExecutor.Execute(rule.OutputTemplate, data)
	if err != nil {
		return "", fmt.Errorf("unable to execute output template for metaalarm rule %s: %w", rule.ID, err)
	}

	return res, nil
}

func (p *metaAlarmEventProcessor) newMetaAlarm(
	params rpc.AxeParameters,
	entity types.Entity,
	alarmConfig config.AlarmConfig,
) types.Alarm {
	now := datetime.NewCpsTime()
	alarm := types.Alarm{
		EntityID: entity.ID,
		ID:       utils.NewID(),
		Time:     now,
		Value: types.AlarmValue{
			Connector:         canopsis.DefaultSystemAlarmConnector,
			ConnectorName:     canopsis.DefaultSystemAlarmConnector,
			Component:         entity.Component,
			Resource:          entity.Name,
			CreationDate:      now,
			DisplayName:       types.GenDisplayName(alarmConfig.DisplayNameScheme),
			InitialOutput:     params.Output,
			Output:            params.Output,
			InitialLongOutput: params.LongOutput,
			LongOutput:        params.LongOutput,
			LongOutputHistory: []string{params.LongOutput},
			LastUpdateDate:    params.Timestamp,
			LastEventDate:     now,
			Parents:           []string{},
			Children:          []string{},
			UnlinkedParents:   []string{},
			Infos:             map[string]map[string]interface{}{},
			RuleVersion:       map[string]string{},
		},
	}

	return alarm
}

func (p *metaAlarmEventProcessor) getChildStepMsg(
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

func applyOnChild(changeType types.AlarmChangeType) bool {
	switch changeType {
	case types.AlarmChangeTypeAck,
		types.AlarmChangeTypeAckremove,
		types.AlarmChangeTypeAssocTicket,
		types.AlarmChangeTypeCancel,
		types.AlarmChangeTypeChangeState,
		types.AlarmChangeTypeComment,
		types.AlarmChangeTypeSnooze,
		types.AlarmChangeTypeUncancel,
		types.AlarmChangeTypeUpdateStatus,
		types.AlarmChangeTypeDeclareTicketWebhook,
		types.AlarmChangeTypeAutoDeclareTicketWebhook:
		return true
	}

	return false
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
	statusStepMessage := libaxeevent.ConcatOutputAndRuleName(output, statusRuleName)
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

func resolveMetaAlarm(metaAlarm *types.Alarm, timestamp datetime.CpsTime) bson.M {
	metaAlarm.Value.Resolved = &timestamp
	metaAlarm.Value.Duration = int64(timestamp.Sub(metaAlarm.Value.CreationDate.Time).Seconds())
	metaAlarm.Value.CurrentStateDuration = int64(timestamp.Sub(metaAlarm.Value.State.Timestamp.Time).Seconds())
	incUpdate := bson.M{}
	if metaAlarm.Value.Snooze != nil {
		snoozeDuration := int64(timestamp.Sub(metaAlarm.Value.Snooze.Timestamp.Time).Seconds())
		metaAlarm.Value.SnoozeDuration += snoozeDuration
		incUpdate["v.snooze_duration"] = snoozeDuration
	}

	if !metaAlarm.Value.PbehaviorInfo.IsActive() {
		enterTimestamp := datetime.CpsTime{}
		for i := len(metaAlarm.Value.Steps) - 2; i >= 0; i-- {
			if metaAlarm.Value.Steps[i].Type == types.AlarmStepPbhEnter {
				enterTimestamp = metaAlarm.Value.Steps[i].Timestamp
				break
			}
		}

		if !enterTimestamp.IsZero() {
			pbhDuration := int64(timestamp.Sub(enterTimestamp.Time).Seconds())
			metaAlarm.Value.PbehaviorInactiveDuration += pbhDuration
			incUpdate["v.pbh_inactive_duration"] = pbhDuration
		}
	}

	if (metaAlarm.Value.Snooze != nil || !metaAlarm.Value.PbehaviorInfo.IsActive()) && metaAlarm.Value.InactiveStart != nil {
		inactiveDuration := int64(timestamp.Sub(metaAlarm.Value.InactiveStart.Time).Seconds())
		metaAlarm.Value.InactiveDuration += inactiveDuration
		incUpdate["v.inactive_duration"] = inactiveDuration
	}

	metaAlarm.Value.ActiveDuration = metaAlarm.Value.Duration - metaAlarm.Value.InactiveDuration
	newStep := types.NewAlarmStep(types.AlarmStepResolve, timestamp, canopsis.DefaultEventAuthor, "", "", "",
		types.InitiatorSystem, false)
	update := bson.M{
		"$set": bson.M{
			"v.resolved":               metaAlarm.Value.Resolved,
			"v.duration":               metaAlarm.Value.Duration,
			"v.current_state_duration": metaAlarm.Value.CurrentStateDuration,
			"v.active_duration":        metaAlarm.Value.ActiveDuration,
		},
		"$unset": bson.M{
			"not_acked_metric_type":      "",
			"not_acked_metric_send_time": "",
			"not_acked_since":            "",
		},
		"$push": bson.M{
			"v.steps": newStep,
		},
	}
	if len(incUpdate) > 0 {
		update["$inc"] = incUpdate
	}

	return update
}
