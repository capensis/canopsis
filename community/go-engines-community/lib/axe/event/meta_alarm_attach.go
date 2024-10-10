package event

import (
	"context"
	"errors"
	"fmt"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	libevent "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

func NewMetaAlarmAttachProcessor(
	dbClient mongo.DbClient,
	ruleAdapter correlation.RulesAdapter,
	adapter libalarm.Adapter,
	alarmStatusService alarmstatus.Service,
	templateExecutor template.Executor,
	metricsSender metrics.Sender,
	encoder encoding.Encoder,
	eventGenerator libevent.Generator,
	amqpPublisher libamqp.Publisher,
	logger zerolog.Logger,
) Processor {
	return &metaAlarmAttachProcessor{
		dbClient:           dbClient,
		alarmCollection:    dbClient.Collection(mongo.AlarmMongoCollection),
		ruleAdapter:        ruleAdapter,
		adapter:            adapter,
		alarmStatusService: alarmStatusService,
		templateExecutor:   templateExecutor,
		metricsSender:      metricsSender,
		encoder:            encoder,
		eventGenerator:     eventGenerator,
		amqpPublisher:      amqpPublisher,
		logger:             logger,
	}
}

type metaAlarmAttachProcessor struct {
	dbClient        mongo.DbClient
	alarmCollection mongo.DbCollection

	ruleAdapter        correlation.RulesAdapter
	adapter            libalarm.Adapter
	alarmStatusService alarmstatus.Service

	metricsSender  metrics.Sender
	encoder        encoding.Encoder
	eventGenerator libevent.Generator
	amqpPublisher  libamqp.Publisher
	logger         zerolog.Logger

	templateExecutor template.Executor
}

func (p *metaAlarmAttachProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil {
		return result, nil
	}

	updatedChildrenAlarms, updatedChildrenEvents, err := p.attachChildrenToMetaAlarm(ctx, event)
	if err != nil {
		return result, err
	}

	result.Forward = false

	for _, child := range updatedChildrenAlarms {
		p.metricsSender.SendCorrelation(event.Parameters.Timestamp.Time, child)
	}

	for _, childEvent := range updatedChildrenEvents {
		err = p.sendToFifo(ctx, childEvent)
		if err != nil {
			return result, err
		}
	}

	return result, nil
}

func (p *metaAlarmAttachProcessor) sendToFifo(ctx context.Context, event types.Event) error {
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

func (p *metaAlarmAttachProcessor) attachChildrenToMetaAlarm(ctx context.Context, event rpc.AxeEvent) ([]types.Alarm, []types.Event, error) {
	if len(event.Parameters.MetaAlarmChildren) == 0 {
		return nil, nil, nil
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

		err := p.alarmCollection.FindOne(ctx, bson.M{"d": event.Entity.ID, "v.resolved": nil}).Decode(&metaAlarm)
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

		alarms, err := getAlarmsWithEntityByMatch(ctx, p.alarmCollection, bson.M{
			"d":          bson.M{"$in": event.Parameters.MetaAlarmChildren},
			"v.resolved": nil,
		})
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

			newStep := NewAlarmStep(types.AlarmStepMetaAlarmAttach, event.Parameters, !childAlarm.Alarm.Value.PbehaviorInfo.IsDefaultActive())
			newStep.Message = getMetaAlarmChildStepMsg(rule, metaAlarm, event)
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

		childrenCount, err := p.adapter.GetCountOpenedAlarmsByIDs(ctx, metaAlarm.Value.Children)
		if err != nil {
			return err
		}

		output := ""
		if rule.Type == correlation.RuleTypeManualGroup {
			output = event.Parameters.Output
		} else {
			output, err = executeMetaAlarmOutputTpl(p.templateExecutor, correlation.EventExtraInfosMeta{
				Rule:     rule,
				Count:    childrenCount,
				Children: lastChild,
			})
			if err != nil {
				return err
			}
		}

		var setUpdate, pushUpdate bson.M
		if worstState > metaAlarm.CurrentState() {
			setUpdate, pushUpdate, err = updateMetaAlarmState(&metaAlarm, *event.Entity, event.Parameters.Timestamp,
				worstState, output, p.alarmStatusService)
			if err != nil {
				return err
			}
		}

		if setUpdate == nil {
			setUpdate = bson.M{}
		}

		metaAlarm.Value.Output = output
		setUpdate["v.output"] = output
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
		return nil, nil, err
	}

	if metaAlarm.ID == "" {
		return nil, nil, nil
	}

	childrenEvents := p.applyActionsOnChildren(metaAlarm, updatedChildrenAlarmsWithEntity)
	for _, e := range activateChildEvents {
		err = p.sendToFifo(ctx, e)
		if err != nil {
			return nil, nil, err
		}
	}

	return updatedChildrenAlarms, childrenEvents, nil
}

func (p *metaAlarmAttachProcessor) applyActionsOnChildren(
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

func (p *metaAlarmAttachProcessor) getChildEventByStep(
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

	childEvent.Output = getMetaAlarmChildEventOutput(metaAlarm, metaAlarmStep.Message, metaAlarmStep.Initiator, isTicketStep)

	return childEvent, nil
}
