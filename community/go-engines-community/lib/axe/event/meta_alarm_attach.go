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
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
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

	metricsSender metrics.Sender
	encoder       encoding.Encoder
	amqpPublisher libamqp.Publisher
	logger        zerolog.Logger

	templateExecutor template.Executor
}

func (p *metaAlarmAttachProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil {
		return result, nil
	}

	_, updatedChildrenAlarms, updatedChildrenEvents, err := p.attachChildrenToMetaAlarm(ctx, event)
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

func (p *metaAlarmAttachProcessor) attachChildrenToMetaAlarm(ctx context.Context, event rpc.AxeEvent) (*types.Alarm, []types.Alarm, []types.Event, error) {
	if len(event.Parameters.MetaAlarmChildren) == 0 {
		return nil, nil, nil, nil
	}

	var updatedChildrenAlarms []types.Alarm
	var metaAlarm types.Alarm
	var err error

	err = p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		updatedChildrenAlarms = make([]types.Alarm, 0)
		var lastChild types.AlarmWithEntity

		err = p.alarmCollection.FindOne(ctx, bson.M{"d": event.Entity.ID, "v.resolved": nil}).Decode(&metaAlarm)
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
		newStep := types.NewMetaAlarmAttachStep(metaAlarm, rule.Name)
		worstState := types.CpsNumber(types.AlarmStateOK)
		eventsCount := types.CpsNumber(0)

		for _, childAlarm := range alarms {
			if !childAlarm.Alarm.AddParent(metaAlarm.EntityID) {
				continue
			}

			metaAlarm.AddChild(childAlarm.Entity.ID)
			eventsCount += childAlarm.Alarm.Value.EventsCount
			if lastEventDate.Before(childAlarm.Alarm.Value.LastEventDate) {
				lastEventDate = childAlarm.Alarm.Value.LastEventDate
			}

			err = childAlarm.Alarm.PartialUpdateAddStepWithStep(newStep)
			if err != nil {
				return err
			}

			if childAlarm.Alarm.Value.State.Value > worstState {
				worstState = childAlarm.Alarm.Value.State.Value
			}

			updatedChildrenAlarms = append(updatedChildrenAlarms, childAlarm.Alarm)
			lastChild = childAlarm
		}

		if len(updatedChildrenAlarms) == 0 {
			return nil
		}

		if metaAlarm.Value.Meta == "" {
			metaAlarm.SetMeta(event.Parameters.MetaAlarmRuleID)
			metaAlarm.SetMetaValuePath(event.Parameters.MetaAlarmValuePath)
		}

		if metaAlarm.Value.LastEventDate.Unix() != lastEventDate.Unix() {
			metaAlarm.PartialUpdateLastEventDate(lastEventDate)
		}

		metaAlarm.IncrementEventsCount(eventsCount)
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

		if output == "" {
			output = metaAlarm.Value.Output
		}

		metaAlarm.UpdateOutput(output)
		if worstState > metaAlarm.CurrentState() {
			err = updateMetaAlarmState(&metaAlarm, *event.Entity, event.Parameters.Timestamp, worstState, output, p.alarmStatusService)
			if err != nil {
				return err
			}
		}

		return p.adapter.PartialMassUpdateOpen(ctx, append([]types.Alarm{metaAlarm}, updatedChildrenAlarms...))
	})
	if err != nil {
		return nil, nil, nil, err
	}

	if metaAlarm.ID == "" {
		return nil, nil, nil, nil
	}

	childrenEvents, err := p.applyActionsOnChildren(metaAlarm, updatedChildrenAlarms)
	if err != nil {
		return nil, nil, nil, err
	}

	return &metaAlarm, updatedChildrenAlarms, childrenEvents, nil
}

func (p *metaAlarmAttachProcessor) applyActionsOnChildren(metaAlarm types.Alarm, childrenAlarms []types.Alarm) ([]types.Event, error) {
	var events []types.Event

	steps := metaAlarm.GetAppliedActions()

	for _, childAlarm := range childrenAlarms {
		childEvent := types.Event{
			Connector:     childAlarm.Value.Connector,
			ConnectorName: childAlarm.Value.ConnectorName,
			Resource:      childAlarm.Value.Resource,
			Component:     childAlarm.Value.Component,
			Timestamp:     types.NewCpsTime(),
		}
		childEvent.SourceType = childEvent.DetectSourceType()

		for _, step := range steps {
			childEvent.Output = step.Message
			childEvent.Author = step.Author
			childEvent.UserID = step.UserID
			childEvent.Initiator = step.Initiator
			childEvent.Role = step.Role
			switch step.Type {
			case types.AlarmStepAck:
				childEvent.EventType = types.EventTypeAck
			case types.AlarmStepSnooze:
				childEvent.EventType = types.EventTypeSnooze
				childEvent.Duration = types.CpsNumber(int64(step.Value) - childEvent.Timestamp.Unix())
				if childEvent.Duration <= 0 {
					continue
				}
			case types.AlarmStepAssocTicket:
				childEvent.EventType = types.EventTypeAssocTicket
				childEvent.TicketInfo = step.TicketInfo
				childEvent.TicketInfo.TicketMetaAlarmID = metaAlarm.ID
			case types.AlarmStepDeclareTicket:
				childEvent.EventType = types.EventTypeDeclareTicketWebhook
				childEvent.TicketInfo = step.TicketInfo
				childEvent.TicketInfo.TicketMetaAlarmID = metaAlarm.ID
			case types.AlarmStepComment:
				childEvent.EventType = types.EventTypeComment
			default:
				continue
			}

			events = append(events, childEvent)
		}
	}

	return events, nil
}
