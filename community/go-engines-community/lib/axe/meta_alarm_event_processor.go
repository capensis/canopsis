package axe

import (
	"context"
	"fmt"
	"math"
	"strings"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/sync/errgroup"
)

const workers = 10

func NewMetaAlarmEventProcessor(
	dbClient mongo.DbClient,
	adapter libalarm.Adapter,
	ruleAdapter correlation.RulesAdapter,
	alarmStatusService alarmstatus.Service,
	alarmConfigProvider config.AlarmConfigProvider,
	encoder encoding.Encoder,
	amqpPublisher libamqp.Publisher,
	fifoExchange, fifoQueue string,
	metricsSender metrics.Sender,
	logger zerolog.Logger,
) libalarm.MetaAlarmEventProcessor {
	return &metaAlarmEventProcessor{
		dbClient:            dbClient,
		adapter:             adapter,
		ruleAdapter:         ruleAdapter,
		alarmStatusService:  alarmStatusService,
		alarmConfigProvider: alarmConfigProvider,
		encoder:             encoder,
		amqpPublisher:       amqpPublisher,
		fifoExchange:        fifoExchange,
		fifoQueue:           fifoQueue,
		metricsSender:       metricsSender,
		logger:              logger,
	}
}

type metaAlarmEventProcessor struct {
	dbClient    mongo.DbClient
	adapter     libalarm.Adapter
	ruleAdapter correlation.RulesAdapter

	alarmStatusService  alarmstatus.Service
	alarmConfigProvider config.AlarmConfigProvider

	encoder                 encoding.Encoder
	amqpPublisher           libamqp.Publisher
	fifoExchange, fifoQueue string

	metricsSender metrics.Sender

	logger zerolog.Logger
}

func (p *metaAlarmEventProcessor) CreateMetaAlarm(ctx context.Context, event rpc.AxeEvent) (*types.Alarm, error) {
	if event.Entity == nil {
		return nil, nil
	}

	var updatedChildAlarms []types.Alarm
	var metaAlarm types.Alarm

	err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		updatedChildAlarms = make([]types.Alarm, 0)
		metaAlarm = types.Alarm{}

		ruleIdentifier := event.Parameters.MetaAlarmRuleID
		rule, err := p.ruleAdapter.GetRule(ctx, ruleIdentifier)
		if err != nil {
			return fmt.Errorf("cannot fetch meta alarm rule id=%q: %w", ruleIdentifier, err)
		} else if rule.ID == "" {
			return fmt.Errorf("meta alarm rule id=%q not found", ruleIdentifier)
		} else {
			ruleIdentifier = rule.Name
		}

		metaAlarm = p.newAlarm(event.Parameters, *event.Entity, p.alarmConfigProvider.Get())
		metaAlarm.Value.Meta = event.Parameters.MetaAlarmRuleID
		metaAlarm.Value.MetaValuePath = event.Parameters.MetaAlarmValuePath

		if event.Parameters.DisplayName != "" {
			metaAlarm.Value.DisplayName = event.Parameters.DisplayName
		}

		var childAlarms []types.Alarm
		worstState := types.CpsNumber(types.AlarmStateMinor)

		if len(event.Parameters.MetaAlarmChildren) > 0 {
			err := p.adapter.GetOpenedAlarmsByIDs(ctx, event.Parameters.MetaAlarmChildren, &childAlarms)
			if err != nil {
				return fmt.Errorf("cannot fetch children alarms: %w", err)
			}

			newStep := types.NewMetaAlarmAttachStep(metaAlarm, ruleIdentifier)
			for i := range childAlarms {
				if childAlarms[i].Value.State != nil {
					childState := childAlarms[i].Value.State.Value
					if childState > worstState {
						worstState = childState
					}
				}

				if !childAlarms[i].HasParentByEID(metaAlarm.EntityID) {
					childAlarms[i].AddParent(metaAlarm.EntityID)
					err = childAlarms[i].PartialUpdateAddStepWithStep(newStep)
					if err != nil {
						return err
					}
					metaAlarm.AddChild(childAlarms[i].EntityID)
					updatedChildAlarms = append(updatedChildAlarms, childAlarms[i])
				}
			}
		}

		err = UpdateAlarmState(&metaAlarm, *event.Entity, event.Parameters.Timestamp, worstState, event.Parameters.Output, p.alarmStatusService)
		if err != nil {
			return err
		}

		err = p.adapter.Insert(ctx, metaAlarm)
		if err != nil {
			return fmt.Errorf("cannot create alarm: %w", err)
		}

		err = p.adapter.PartialMassUpdateOpen(ctx, childAlarms)
		if err != nil {
			return fmt.Errorf("cannot update children alarms: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	for _, child := range updatedChildAlarms {
		p.metricsSender.SendCorrelation(event.Parameters.Timestamp.Time, child)
	}

	return &metaAlarm, nil
}

func (p *metaAlarmEventProcessor) ProcessAxeRpc(ctx context.Context, event rpc.AxeEvent, eventRes rpc.AxeResultEvent) error {
	if eventRes.Alarm == nil || eventRes.AlarmChangeType == "" {
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

	if alarm.IsMetaChildren() {
		return p.processChildRpc(ctx, event, eventRes)
	}

	return nil
}

func (p *metaAlarmEventProcessor) processParentRpc(ctx context.Context, event rpc.AxeEvent, eventRes rpc.AxeResultEvent) error {
	if !applyOnChild(eventRes.AlarmChangeType) {
		return nil
	}

	childEvent := types.Event{
		EventType:  event.EventType,
		Timestamp:  types.NewCpsTime(),
		Output:     event.Parameters.Output,
		Author:     event.Parameters.Author,
		UserID:     event.Parameters.User,
		Initiator:  types.InitiatorSystem,
		TicketInfo: event.Parameters.TicketInfo,
	}

	childEvent.TicketInfo.TicketMetaAlarmID = eventRes.Alarm.ID

	switch eventRes.AlarmChangeType {
	case types.AlarmChangeTypeDeclareTicketWebhook,
		types.AlarmChangeTypeAutoDeclareTicketWebhook:
		childEvent.EventType = types.EventTypeDeclareTicketWebhook
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

	err := p.sendChildrenEvents(ctx, eventRes.Alarm.Value.Children, childEvent)
	if err != nil {
		return err
	}

	return nil
}

func (p *metaAlarmEventProcessor) processComponentRpc(ctx context.Context, event rpc.AxeEvent, eventRes rpc.AxeResultEvent) error {
	if !event.Parameters.TicketResources ||
		(eventRes.AlarmChangeType != types.AlarmChangeTypeDeclareTicketWebhook && eventRes.AlarmChangeType != types.AlarmChangeTypeAutoDeclareTicketWebhook) ||
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

	for _, resource := range resources {
		if resource.Entity.Type != types.EntityTypeResource {
			continue
		}

		resourceEvent := types.Event{
			EventType:     types.EventTypeDeclareTicketWebhook,
			Connector:     resource.Alarm.Value.Connector,
			ConnectorName: resource.Alarm.Value.ConnectorName,
			Resource:      resource.Alarm.Value.Resource,
			Component:     resource.Alarm.Value.Component,
			Timestamp:     types.NewCpsTime(),
			Output:        componentAlarm.Value.Ticket.Message,
			TicketInfo:    componentAlarm.Value.Ticket.TicketInfo,
			Author:        event.Parameters.Author,
			UserID:        event.Parameters.User,
			Initiator:     types.InitiatorSystem,
		}
		resourceEvent.SourceType = resourceEvent.DetectSourceType()

		err = p.sendToFifo(ctx, resourceEvent)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *metaAlarmEventProcessor) processChildRpc(ctx context.Context, event rpc.AxeEvent, eventRes rpc.AxeResultEvent) error {
	switch eventRes.AlarmChangeType {
	case types.AlarmChangeTypeStateIncrease, types.AlarmChangeTypeStateDecrease, types.AlarmChangeTypeChangeState:
		err := p.updateParentState(ctx, *eventRes.Alarm)
		if err != nil {
			return err
		}
	case types.AlarmChangeTypeResolve:
		err := p.resolveParents(ctx, *eventRes.Alarm, event.Parameters.Timestamp)
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

	if event.EventType == types.EventTypeCheck && p.alarmConfigProvider.Get().EnableLastEventDate {
		err := p.adapter.UpdateLastEventDate(ctx, eventRes.Alarm.Value.Parents, event.Parameters.Timestamp)
		if err != nil {
			return fmt.Errorf("cannot update parent alarms: %w", err)
		}
	}

	return nil
}

func (p *metaAlarmEventProcessor) sendChildrenEvents(ctx context.Context, childrenIds []string, childEvent types.Event) error {
	var alarms []types.Alarm
	err := p.adapter.GetOpenedAlarmsByIDs(ctx, childrenIds, &alarms)
	if err != nil {
		return fmt.Errorf("cannot fetch children alarms: %w", err)
	}

	for _, alarm := range alarms {
		childEvent.Connector = alarm.Value.Connector
		childEvent.ConnectorName = alarm.Value.ConnectorName
		childEvent.Resource = alarm.Value.Resource
		childEvent.Component = alarm.Value.Component
		childEvent.SourceType = childEvent.DetectSourceType()

		err = p.sendToFifo(ctx, childEvent)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *metaAlarmEventProcessor) resolveParents(ctx context.Context, childAlarm types.Alarm, timestamp types.CpsTime) error {
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

					err = parentAlarm.Alarm.PartialUpdateResolve(types.NewCpsTime())
					if err != nil {
						return fmt.Errorf("cannot update alarm: %w", err)
					}

					err = p.adapter.PartialUpdateOpen(ctx, &parentAlarm.Alarm)
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

					if childState > parentState {
						newState = childState
					} else if childState < parentState {
						r, err := p.adapter.GetWorstAlarmState(ctx, parentAlarm.Alarm.Value.Children)
						if err != nil {
							return fmt.Errorf("cannot fetch children state: %w", err)
						}

						newState = types.CpsNumber(r)
					} else {
						return nil
					}

					err = UpdateAlarmState(&parentAlarm.Alarm, parentAlarm.Entity, childAlarm.Value.LastUpdateDate,
						newState, parentAlarm.Alarm.Value.Output, p.alarmStatusService)
					if err != nil {
						return fmt.Errorf("cannot update parent: %w", err)
					}

					err = p.adapter.PartialUpdateOpen(ctx, &parentAlarm.Alarm)
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
		p.fifoExchange,
		p.fifoQueue,
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

func (p *metaAlarmEventProcessor) newAlarm(
	params rpc.AxeParameters,
	entity types.Entity,
	alarmConfig config.AlarmConfig,
) types.Alarm {
	now := types.NewCpsTime()
	alarm := types.Alarm{
		EntityID: entity.ID,
		ID:       utils.NewID(),
		Time:     now,
		Value: types.AlarmValue{
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

	switch entity.Type {
	case types.EntityTypeResource:
		alarm.Value.Resource = entity.Name
		alarm.Value.Component = entity.Component
		alarm.Value.Connector, alarm.Value.ConnectorName, _ = strings.Cut(entity.Connector, "/")
	case types.EntityTypeComponent, types.EntityTypeService:
		alarm.Value.Component = entity.Name
		alarm.Value.Connector, alarm.Value.ConnectorName, _ = strings.Cut(entity.Connector, "/")
	case types.EntityTypeConnector:
		alarm.Value.Connector, alarm.Value.ConnectorName, _ = strings.Cut(entity.ID, "/")
	}

	return alarm
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

func UpdateAlarmState(alarm *types.Alarm, entity types.Entity, timestamp types.CpsTime, state types.CpsNumber, output string,
	service alarmstatus.Service) error {
	var currentState, currentStatus types.CpsNumber
	if alarm.Value.State != nil {
		currentState = alarm.Value.State.Value
		currentStatus = alarm.Value.Status.Value
	}

	author := ""
	if entity.Type != types.EntityTypeService {
		author = strings.Replace(entity.Connector, "/", ".", 1)
	} else {
		author = alarm.Value.Connector + "." + alarm.Value.ConnectorName
	}

	if state != currentState {
		// Event is an OK, so the alarm should be resolved anyway
		if alarm.IsStateLocked() && state != types.AlarmStateOK {
			return nil
		}

		// Create new Step to keep track of the alarm history
		newStep := types.NewAlarmStep(types.AlarmStepStateIncrease, timestamp, author, output, "", "", "")
		newStep.Value = state

		if state < currentState {
			newStep.Type = types.AlarmStepStateDecrease
		}

		alarm.Value.State = &newStep
		err := alarm.Value.Steps.Add(newStep)
		if err != nil {
			return err
		}

		alarm.Value.TotalStateChanges++
		alarm.Value.LastUpdateDate = timestamp
	}

	newStatus := service.ComputeStatus(*alarm, entity)

	if newStatus == currentStatus {
		if state != currentState {
			alarm.Value.StateChangesSinceStatusUpdate++

			alarm.AddUpdate("$set", bson.M{
				"v.state":                             alarm.Value.State,
				"v.state_changes_since_status_update": alarm.Value.StateChangesSinceStatusUpdate,
				"v.total_state_changes":               alarm.Value.TotalStateChanges,
				"v.last_update_date":                  alarm.Value.LastUpdateDate,
			})
			alarm.AddUpdate("$push", bson.M{"v.steps": alarm.Value.State})
		}

		return nil
	}

	// Create new Step to keep track of the alarm history
	newStepStatus := types.NewAlarmStep(types.AlarmStepStatusIncrease, timestamp, author, output, "", "", "")
	newStepStatus.Value = newStatus

	if newStatus < currentStatus {
		newStepStatus.Type = types.AlarmStepStatusDecrease
	}

	alarm.Value.Status = &newStepStatus
	err := alarm.Value.Steps.Add(newStepStatus)
	if err != nil {
		return err
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

	alarm.AddUpdate("$set", set)
	alarm.AddUpdate("$push", bson.M{"v.steps": bson.M{"$each": newSteps}})

	return nil
}
