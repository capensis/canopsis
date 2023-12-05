/*
Package idlealarm implements alarm modification on idle alarm.
*/
package idlealarm

//go:generate mockgen -destination=../../../mocks/lib/canopsis/idlealarm/service.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlealarm Service

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	libentity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	libevent "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlerule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
)

// Service interface is used to implement alarms modification by idle rules.
type Service interface {
	Process(ctx context.Context) ([]types.Event, error)
}

// NewService creates new service.
func NewService(
	ruleAdapter idlerule.RuleAdapter,
	alarmAdapter libalarm.Adapter,
	entityAdapter libentity.Adapter,
	pbhRpcClient engine.RPCClient,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Service {
	return &baseService{
		ruleAdapter:   ruleAdapter,
		alarmAdapter:  alarmAdapter,
		entityAdapter: entityAdapter,
		pbhRpcClient:  pbhRpcClient,
		encoder:       encoder,
		logger:        logger,
	}
}

type baseService struct {
	ruleAdapter   idlerule.RuleAdapter
	alarmAdapter  libalarm.Adapter
	entityAdapter libentity.Adapter
	pbhRpcClient  engine.RPCClient
	encoder       encoding.Encoder
	logger        zerolog.Logger
}

func (s *baseService) Process(ctx context.Context) (res []types.Event, resErr error) {
	now := datetime.NewCpsTime()
	eventGenerator := libevent.NewGenerator("engine", "axe")
	rules, err := s.ruleAdapter.GetEnabled(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch idle rules: %w", err)
	}
	ids := make([]string, len(rules))
	for i, rule := range rules {
		ids[i] = rule.ID
	}
	s.logger.Debug().Strs("rules", ids).Msg("load idle rules")

	var allMinDuration, entityMinDuration datetime.DurationWithUnit
	for _, rule := range rules {
		if allMinDuration.Value == 0 || allMinDuration.AddTo(now).After(rule.Duration.AddTo(now)) {
			allMinDuration = rule.Duration
		}

		switch rule.Type {
		case idlerule.RuleTypeAlarm:
			/*do nothing*/
		case idlerule.RuleTypeEntity:
			if entityMinDuration.Value == 0 || entityMinDuration.AddTo(now).After(rule.Duration.AddTo(now)) {
				entityMinDuration = rule.Duration
			}
		default:
			return nil, fmt.Errorf("unknown type of idle rule: %q for id=%s", rule.Type, rule.ID)
		}
	}

	events := make([]types.Event, 0)
	checkedEntities := make([]string, 0)

	if allMinDuration.Value > 0 {
		before := allMinDuration.SubFrom(now)
		cursor, err := s.alarmAdapter.GetOpenedAlarmsWithLastDatesBefore(ctx, before)
		if err != nil {
			return events, fmt.Errorf("cannot fetch alarms: %w", err)
		}

		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			alarm := types.AlarmWithEntity{}
			err := cursor.Decode(&alarm)
			if err != nil {
				return events, fmt.Errorf("cannot decode alarm %w", err)
			}

			event, err := s.applyRules(ctx, rules, alarm.Entity, &alarm.Alarm, eventGenerator, now)
			if err != nil {
				return events, err
			}

			if event != nil {
				events = append(events, *event)
			}

			checkedEntities = append(checkedEntities, alarm.Entity.ID)
		}
	}

	if entityMinDuration.Value > 0 {
		before := entityMinDuration.SubFrom(now)
		cursor, err := s.entityAdapter.GetAllWithLastUpdateDateBefore(ctx, before, checkedEntities)
		if err != nil {
			return events, fmt.Errorf("cannot fetch entities: %w", err)
		}

		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			entity := types.Entity{}
			err := cursor.Decode(&entity)
			if err != nil {
				return events, fmt.Errorf("cannot decode entity : %w", err)
			}

			event, err := s.applyRules(ctx, rules, entity, nil, eventGenerator, now)
			if err != nil {
				return events, err
			}

			if event != nil {
				events = append(events, *event)
			}
		}
	}

	closeEvents, err := s.closeConnectorAlarms(ctx)
	if err != nil {
		return events, err
	}

	events = append(events, closeEvents...)

	return events, nil
}

// applyRules applies rules on entity.
// Alarm rules are skipped if alarm arg is nil.
// Last entity alarm is fetched for entity rules if alarm arg is nil.
func (s *baseService) applyRules(
	ctx context.Context,
	rules []idlerule.Rule,
	entity types.Entity,
	alarm *types.Alarm,
	eventGenerator libevent.Generator,
	now datetime.CpsTime,
) (*types.Event, error) {
	lastAlarm := alarm

	for _, rule := range rules {
		switch rule.Type {
		case idlerule.RuleTypeAlarm:
			if alarm == nil {
				continue
			}

			matched, err := rule.Matches(types.AlarmWithEntity{
				Alarm:  *alarm,
				Entity: entity,
			}, now)
			if err != nil {
				s.logger.Error().Err(err).Str("idle_rule", rule.ID).Msg("match idle rule returned error, skip")
				continue
			}

			if matched && !alarm.Value.PbehaviorInfo.OneOf(rule.DisableDuringPeriods) {
				return s.applyAlarmRule(ctx, rule, *alarm, entity, now)
			}
		case idlerule.RuleTypeEntity:
			matched, err := rule.Matches(types.AlarmWithEntity{Entity: entity}, now)
			if err != nil {
				s.logger.Error().Err(err).Str("idle_rule", rule.ID).Msg("match idle rule returned error, skip")
				continue
			}

			if !matched || entity.PbehaviorInfo.OneOf(rule.DisableDuringPeriods) {
				continue
			}

			if lastAlarm == nil {
				var err error
				lastAlarm, err = s.alarmAdapter.GetLastAlarmByEntityID(ctx, entity.ID)
				if err != nil {
					return nil, fmt.Errorf("cannot fetch alarm: %w", err)
				}

				// If some rule has already been applied on entity.
				if lastAlarm != nil && !lastAlarm.IsResolved() &&
					lastAlarm.Value.Status.Value == types.AlarmStatusNoEvents {
					return nil, nil
				}
			}

			return s.applyEntityRule(rule, entity, lastAlarm, eventGenerator)
		}
	}

	return nil, nil
}

func (s *baseService) applyAlarmRule(
	ctx context.Context,
	rule idlerule.Rule,
	alarm types.Alarm,
	entity types.Entity,
	now datetime.CpsTime,
) (*types.Event, error) {
	idleRuleApply := fmt.Sprintf("%s_%s", rule.Type, rule.AlarmCondition)
	event := types.Event{
		Connector:     alarm.Value.Connector,
		ConnectorName: alarm.Value.ConnectorName,
		Component:     alarm.Value.Component,
		Resource:      alarm.Value.Resource,
		Timestamp:     now,
		Author:        canopsis.DefaultEventAuthor,
		Initiator:     types.InitiatorSystem,
		IdleRuleApply: idleRuleApply,
	}
	event.SourceType = event.DetectSourceType()

	event.Output = rule.Operation.Parameters.Output
	if rule.Operation.Parameters.State != nil {
		event.State = *rule.Operation.Parameters.State
	}

	switch rule.Operation.Type {
	case types.ActionTypeAck:
		event.EventType = types.EventTypeAck
	case types.ActionTypeAckRemove:
		event.EventType = types.EventTypeAckremove
	case types.ActionTypeCancel:
		event.EventType = types.EventTypeCancel
	case types.ActionTypeAssocTicket:
		event.EventType = types.EventTypeAssocTicket
		event.TicketInfo = types.TicketInfo{
			Ticket:           rule.Operation.Parameters.Ticket,
			TicketRuleID:     rule.ID,
			TicketRuleName:   types.TicketRuleNameIdleRulePrefix + rule.Name,
			TicketURL:        rule.Operation.Parameters.TicketURL,
			TicketSystemName: rule.Operation.Parameters.TicketSystemName,
			TicketData:       rule.Operation.Parameters.TicketData,
			TicketComment:    rule.Comment,
		}
	case types.ActionTypeChangeState:
		event.EventType = types.EventTypeChangestate
	case types.ActionTypePbehavior:
		rpcEvent := rpc.PbehaviorEvent{
			Alarm:  &alarm,
			Entity: &entity,
			Params: rpc.PbehaviorParameters{
				Name:           rule.Operation.Parameters.Name,
				Reason:         rule.Operation.Parameters.Reason,
				Type:           rule.Operation.Parameters.Type,
				RRule:          rule.Operation.Parameters.RRule,
				Tstart:         rule.Operation.Parameters.Tstart,
				Tstop:          rule.Operation.Parameters.Tstop,
				StartOnTrigger: rule.Operation.Parameters.StartOnTrigger,
				Duration:       rule.Operation.Parameters.Duration,
			},
		}
		b, err := s.encoder.Encode(rpcEvent)
		if err != nil {
			return nil, fmt.Errorf("cannot encode rpc event: %w", err)
		}

		err = s.pbhRpcClient.Call(ctx, engine.RPCMessage{
			CorrelationID: fmt.Sprintf("%s&&%s", alarm.ID, rule.ID),
			Body:          b,
		})
		if err != nil {
			return nil, fmt.Errorf("cannot sent rpc event: %w", err)
		}

		err = s.entityAdapter.UpdateIdleFields(ctx, entity.ID, entity.IdleSince, idleRuleApply)
		if err != nil {
			return nil, fmt.Errorf("cannot update entity: %w", err)
		}

		return nil, nil
	case types.ActionTypeSnooze:
		event.EventType = types.EventTypeSnooze
		event.Duration = types.CpsNumber(rule.Operation.Parameters.Duration.AddTo(now).Sub(now.Time).Seconds())
	default:
		return nil, fmt.Errorf("unknown idle rule id=%q operation type=%q", rule.ID, rule.Operation.Type)
	}

	if event.EventType == "" {
		return nil, fmt.Errorf("invalid idle rule id=%q operation params %v", rule.ID, rule.Operation.Parameters)
	}

	return &event, nil
}

// applyEntityRule generates no events event.
// It uses alarm arg to fill event field.
// It uses entity type and connector arg to fill event field if alarm arg is nil.
func (s *baseService) applyEntityRule(
	rule idlerule.Rule,
	entity types.Entity,
	alarm *types.Alarm,
	eventGenerator libevent.Generator,
) (*types.Event, error) {
	event := types.Event{}
	if alarm == nil {
		var err error
		event, err = eventGenerator.Generate(entity)
		if err != nil {
			return nil, err
		}
	} else {
		event.Connector = alarm.Value.Connector
		event.ConnectorName = alarm.Value.ConnectorName
		event.Component = alarm.Value.Component
		event.Resource = alarm.Value.Resource
		event.SourceType = event.DetectSourceType()
	}

	event.EventType = types.EventTypeNoEvents
	event.Timestamp = datetime.NewCpsTime()
	event.State = types.AlarmStateCritical
	event.Author = canopsis.DefaultEventAuthor
	event.Initiator = types.InitiatorSystem
	event.Output = fmt.Sprintf("Idle rule %s", rule.Name)
	event.IdleRuleApply = rule.Type

	return &event, nil
}

func (s *baseService) closeConnectorAlarms(ctx context.Context) ([]types.Event, error) {
	alarms, err := s.alarmAdapter.GetOpenedAlarmsByConnectorIdleRules(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch opened alarms: %w", err)
	}

	events := make([]types.Event, len(alarms))
	for i, alarm := range alarms {
		events[i] = types.Event{
			EventType:     types.EventTypeNoEvents,
			State:         types.AlarmStateOK,
			Connector:     alarm.Value.Connector,
			ConnectorName: alarm.Value.ConnectorName,
			Component:     alarm.Value.Component,
			Resource:      alarm.Value.Resource,
			Timestamp:     datetime.NewCpsTime(),
			Initiator:     types.InitiatorSystem,
			Output:        alarm.Value.State.Message,
			Author:        alarm.Value.State.Author,
		}
		events[i].SourceType = events[i].DetectSourceType()
	}

	return events, nil
}
