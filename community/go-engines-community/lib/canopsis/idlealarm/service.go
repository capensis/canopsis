/*
Package idlealarm implements alarm modification on idle alarm.
*/
package idlealarm

//go:generate mockgen -destination=../../../mocks/lib/canopsis/idlealarm/service.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlealarm Service

import (
	"context"
	"fmt"
	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	libentity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	libevent "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlerule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
	"time"
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
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Service {
	return &baseService{
		ruleAdapter:   ruleAdapter,
		alarmAdapter:  alarmAdapter,
		entityAdapter: entityAdapter,
		encoder:       encoder,
		logger:        logger,
	}
}

type baseService struct {
	ruleAdapter   idlerule.RuleAdapter
	alarmAdapter  libalarm.Adapter
	entityAdapter libentity.Adapter
	encoder       encoding.Encoder
	logger        zerolog.Logger
}

func (s *baseService) Process(ctx context.Context) (res []types.Event, resErr error) {
	now := types.NewCpsTime()
	eventGenerator := libevent.NewGenerator(s.entityAdapter)
	rules, err := s.ruleAdapter.GetEnabled(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch idle rules: %w", err)
	}
	ids := make([]string, len(rules))
	for i, rule := range rules {
		ids[i] = rule.ID
	}
	s.logger.Debug().Strs("rules", ids).Msg("load idle rules")

	var allMinDuration, entityMinDuration types.DurationWithUnit
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
	now types.CpsTime,
) (*types.Event, error) {
	lastAlarm := alarm

	for _, rule := range rules {
		switch rule.Type {
		case idlerule.RuleTypeAlarm:
			if alarm != nil && rule.Matches(alarm, &entity, now) && !alarm.Value.PbehaviorInfo.OneOf(rule.DisableDuringPeriods) {
				return s.applyAlarmRule(rule, *alarm, now)
			}
		case idlerule.RuleTypeEntity:
			if !rule.Matches(nil, &entity, now) || entity.PbehaviorInfo.OneOf(rule.DisableDuringPeriods) {
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

			return s.applyEntityRule(ctx, rule, entity, lastAlarm, eventGenerator)
		}
	}

	return nil, nil
}

func (s *baseService) applyAlarmRule(
	rule idlerule.Rule,
	alarm types.Alarm,
	now types.CpsTime,
) (*types.Event, error) {
	event := types.Event{
		Connector:     alarm.Value.Connector,
		ConnectorName: alarm.Value.ConnectorName,
		Component:     alarm.Value.Component,
		Resource:      alarm.Value.Resource,
		Timestamp:     now,
		Initiator:     types.InitiatorSystem,
		IdleRuleApply: fmt.Sprintf("%s_%s", rule.Type, rule.AlarmCondition),
	}
	event.SourceType = event.DetectSourceType()

	switch rule.Operation.Type {
	case types.ActionTypeAck:
		if params, ok := rule.Operation.Parameters.(types.OperationParameters); ok {
			event.EventType = types.EventTypeAck
			event.Output = params.Output
			event.Author = params.Author
			event.UserID = params.User
		}
	case types.ActionTypeAckRemove:
		if params, ok := rule.Operation.Parameters.(types.OperationParameters); ok {
			event.EventType = types.EventTypeAckremove
			event.Output = params.Output
			event.Author = params.Author
			event.UserID = params.User
		}
	case types.ActionTypeCancel:
		if params, ok := rule.Operation.Parameters.(types.OperationParameters); ok {
			event.EventType = types.EventTypeCancel
			event.Output = params.Output
			event.Author = params.Author
			event.UserID = params.User
		}
	case types.ActionTypeAssocTicket:
		if params, ok := rule.Operation.Parameters.(types.OperationAssocTicketParameters); ok {
			event.EventType = types.EventTypeAssocTicket
			event.Ticket = params.Ticket
			event.Output = params.Output
			event.Author = params.Author
			event.UserID = params.User
		}
	case types.ActionTypeChangeState:
		if params, ok := rule.Operation.Parameters.(types.OperationChangeStateParameters); ok {
			event.EventType = types.EventTypeChangestate
			event.State = params.State
			event.Output = params.Output
			event.Author = params.Author
			event.UserID = params.User
		}
	case types.ActionTypePbehavior:
		if params, ok := rule.Operation.Parameters.(types.ActionPBehaviorParameters); ok {
			event.EventType = types.EventTypePbhCreate
			encodedParams, err := s.encoder.Encode(params)
			if err != nil {
				return nil, fmt.Errorf("cannot encode parameters: %w", err)
			}
			event.PbhParameters = string(encodedParams)
			event.Author = params.Author
			event.UserID = params.UserID
		}
	case types.ActionTypeSnooze:
		if params, ok := rule.Operation.Parameters.(types.OperationSnoozeParameters); ok {
			event.EventType = types.EventTypeSnooze
			event.Duration = types.CpsNumber(params.Duration.AddTo(now).Sub(now.Time).Seconds())
			event.Output = params.Output
			event.Author = params.Author
			event.UserID = params.User
		}
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
	ctx context.Context,
	rule idlerule.Rule,
	entity types.Entity,
	alarm *types.Alarm,
	eventGenerator libevent.Generator,
) (*types.Event, error) {
	event := types.Event{}
	if alarm == nil {
		var err error
		event, err = eventGenerator.Generate(ctx, entity)
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
	event.Timestamp = types.CpsTime{Time: time.Now()}
	event.State = types.AlarmStateCritical
	event.Author = rule.Author
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
			Timestamp:     types.CpsTime{Time: time.Now()},
			Initiator:     types.InitiatorSystem,
			Output:        alarm.Value.State.Message,
			Author:        alarm.Value.State.Author,
		}
		events[i].SourceType = events[i].DetectSourceType()
	}

	return events, nil
}
