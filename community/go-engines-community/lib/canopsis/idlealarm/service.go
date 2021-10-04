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
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlerule"
	libpbehavior "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
	"strings"
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
	pbhTypeResolver libpbehavior.EntityTypeResolver,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Service {
	return &baseService{
		ruleAdapter:     ruleAdapter,
		alarmAdapter:    alarmAdapter,
		entityAdapter:   entityAdapter,
		pbhTypeResolver: pbhTypeResolver,
		encoder:         encoder,
		logger:          logger,

		connectors: make(map[string]types.Entity),
		components: make(map[string]types.Entity),
	}
}

type baseService struct {
	ruleAdapter     idlerule.RuleAdapter
	alarmAdapter    libalarm.Adapter
	entityAdapter   libentity.Adapter
	pbhTypeResolver libpbehavior.EntityTypeResolver
	encoder         encoding.Encoder
	logger          zerolog.Logger

	connectors map[string]types.Entity
	components map[string]types.Entity
}

func (s *baseService) Process(ctx context.Context) (res []types.Event, resErr error) {
	defer func() {
		s.logger.Debug().Msg("process idle alarms")
		s.connectors = make(map[string]types.Entity)
		s.components = make(map[string]types.Entity)
	}()

	rules, err := s.ruleAdapter.GetEnabled(ctx)
	if err != nil {
		return nil, err
	}

	var allMinDuration, entityMinDuration time.Duration
	for _, rule := range rules {
		if allMinDuration == 0 || allMinDuration > rule.Duration.Duration() {
			allMinDuration = rule.Duration.Duration()
		}

		switch rule.Type {
		case idlerule.RuleTypeAlarm:
			/*do nothing*/
		case idlerule.RuleTypeEntity:
			if entityMinDuration == 0 || entityMinDuration > rule.Duration.Duration() {
				entityMinDuration = rule.Duration.Duration()
			}
		default:
			return nil, fmt.Errorf("unknown type of idle rule: %v for id=%s", rule.Type, rule.ID)
		}
	}

	events := make([]types.Event, 0)
	checkedEntities := make([]string, 0)
	now := time.Now()

	if allMinDuration > 0 {
		before := types.CpsTime{Time: now.Add(-allMinDuration)}
		cursor, err := s.alarmAdapter.GetOpenedAlarmsWithLastDatesBefore(ctx, before)
		if err != nil {
			return events, err
		}

		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			alarm := types.AlarmWithEntity{}
			err := cursor.Decode(&alarm)
			if err != nil {
				return events, err
			}

			event, err := s.applyRules(ctx, rules, alarm.Entity, &alarm.Alarm, now)
			if err != nil {
				return events, err
			}

			if event != nil {
				events = append(events, *event)
			}

			checkedEntities = append(checkedEntities, alarm.Entity.ID)
		}
	}

	if entityMinDuration > 0 {
		before := types.CpsTime{Time: now.Add(-entityMinDuration)}
		cursor, err := s.entityAdapter.GetAllWithLastUpdateDateBefore(ctx, before, checkedEntities)
		if err != nil {
			return events, err
		}

		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			entity := types.Entity{}
			err := cursor.Decode(&entity)
			if err != nil {
				return events, err
			}

			event, err := s.applyRules(ctx, rules, entity, nil, now)
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
	now time.Time,
) (*types.Event, error) {
	lastAlarm := alarm

	for _, rule := range rules {
		switch rule.Type {
		case idlerule.RuleTypeAlarm:
			if alarm != nil && rule.Matches(alarm, &entity) && !alarm.Value.PbehaviorInfo.OneOf(rule.DisableDuringPeriods) {
				event, err := s.applyAlarmRule(rule, *alarm)
				if err != nil || event != nil {
					return event, err
				}
			}
		case idlerule.RuleTypeEntity:
			if !rule.Matches(nil, &entity) {
				continue
			}

			if lastAlarm == nil {
				var err error
				lastAlarm, err = s.alarmAdapter.GetLastAlarmByEntityID(ctx, entity.ID)
				if err != nil {
					return nil, err
				}

				// If some rule has already been applied on entity.
				if lastAlarm != nil && !lastAlarm.IsResolved() &&
					lastAlarm.Value.Status.Value == types.AlarmStatusNoEvents {
					return nil, nil
				}
			}

			// If rule is disabled on pbehavior.
			if len(rule.DisableDuringPeriods) > 0 {
				// Check pbehavior period by open alarm.
				if lastAlarm != nil && !lastAlarm.IsResolved() {
					if lastAlarm.Value.PbehaviorInfo.OneOf(rule.DisableDuringPeriods) {
						continue
					}
				} else {
					pbhRes, err := s.pbhTypeResolver.Resolve(ctx, entity.ID, now)
					if err != nil {
						return nil, err
					}
					if pbhResOneOfType(pbhRes, rule.DisableDuringPeriods) {
						continue
					}
				}
			}

			event, err := s.applyEntityRule(ctx, rule, entity, lastAlarm)
			if err != nil || event != nil {
				return event, err
			}
		}
	}

	return nil, nil
}

func (s *baseService) applyAlarmRule(
	rule idlerule.Rule,
	alarm types.Alarm,
) (*types.Event, error) {
	event := types.Event{
		Connector:     alarm.Value.Connector,
		ConnectorName: alarm.Value.ConnectorName,
		Component:     alarm.Value.Component,
		Resource:      alarm.Value.Resource,
		Timestamp:     types.CpsTime{Time: time.Now()},
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
		}
	case types.ActionTypeAckRemove:
		if params, ok := rule.Operation.Parameters.(types.OperationParameters); ok {
			event.EventType = types.EventTypeAckremove
			event.Output = params.Output
			event.Author = params.Author
		}
	case types.ActionTypeCancel:
		if params, ok := rule.Operation.Parameters.(types.OperationParameters); ok {
			event.EventType = types.EventTypeCancel
			event.Output = params.Output
			event.Author = params.Author
		}
	case types.ActionTypeAssocTicket:
		if params, ok := rule.Operation.Parameters.(types.OperationAssocTicketParameters); ok {
			event.EventType = types.EventTypeAssocTicket
			event.Ticket = params.Ticket
			event.Output = params.Output
			event.Author = params.Author
		}
	case types.ActionTypeChangeState:
		if params, ok := rule.Operation.Parameters.(types.OperationChangeStateParameters); ok {
			event.EventType = types.EventTypeChangestate
			event.State = params.State
			event.Output = params.Output
			event.Author = params.Author
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
		}
	case types.ActionTypeSnooze:
		if params, ok := rule.Operation.Parameters.(types.OperationSnoozeParameters); ok {
			event.EventType = types.EventTypeSnooze
			d := types.CpsNumber(params.Duration.Seconds)
			event.Duration = &d
			event.Output = params.Output
			event.Author = params.Author
		}
	default:
		return nil, fmt.Errorf("unknown idle rule operation %v", rule.Operation)
	}

	if event.EventType == "" {
		return nil, fmt.Errorf("invalid idle rule operation %v", rule.Operation)
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
) (*types.Event, error) {
	event := types.Event{
		EventType:     types.EventTypeNoEvents,
		Timestamp:     types.CpsTime{Time: time.Now()},
		State:         types.AlarmStateCritical,
		Author:        rule.Author,
		Initiator:     types.InitiatorSystem,
		Output:        fmt.Sprintf("Idle rule %s", rule.Name),
		IdleRuleApply: rule.Type,
	}
	if alarm != nil {
		event.Connector = alarm.Value.Connector
		event.ConnectorName = alarm.Value.ConnectorName
		event.Component = alarm.Value.Component
		event.Resource = alarm.Value.Resource
	} else {
		switch entity.Type {
		case types.EntityTypeConnector:
			event.Connector = strings.ReplaceAll(entity.ID, "/"+entity.Name, "")
			event.ConnectorName = entity.Name
		case types.EntityTypeComponent:
			connector, err := s.findConnectorForComponent(ctx, entity)
			if err != nil {
				return nil, err
			}
			if connector == nil {
				s.logger.Error().Msgf("cannot generate event for entity %v : not found any alarm and not found linked connector", entity.ID)
				return nil, nil
			}
			event.Connector = strings.ReplaceAll(connector.ID, "/"+connector.Name, "")
			event.ConnectorName = connector.Name
			event.Component = entity.Name
		case types.EntityTypeResource:
			connector, err := s.findConnectorForResource(ctx, entity)
			if err != nil {
				return nil, err
			}
			if connector == nil {
				s.logger.Error().Msgf("cannot generate event for entity %v : not found any alarm and not found linked connector", entity.ID)
				return nil, nil
			}
			event.Connector = strings.ReplaceAll(connector.ID, "/"+connector.Name, "")
			event.ConnectorName = connector.Name
			if entity.Component != "" {
				event.Component = entity.Component
			} else {
				component, err := s.findComponentForResource(ctx, entity)
				if err != nil {
					return nil, err
				}
				if component == nil {
					s.logger.Error().Msgf("cannot generate event for resource %v : not found any alarm and not found linked component", entity.ID)
					return nil, nil
				}
				event.Component = component.ID
			}
			event.Resource = entity.Name
		default:
			return nil, fmt.Errorf("unknown entity type %v", entity.Type)
		}
	}

	event.SourceType = event.DetectSourceType()

	return &event, nil
}

func (s *baseService) closeConnectorAlarms(ctx context.Context) ([]types.Event, error) {
	alarms, err := s.alarmAdapter.GetOpenedAlarmsByConnectorIdleRules(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch opened alarms: %w", err)
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

func (s *baseService) findConnectorForComponent(ctx context.Context, entity types.Entity) (*types.Entity, error) {
	for _, id := range entity.Impacts {
		if connector, ok := s.connectors[id]; ok {
			return &connector, nil
		}
	}

	connector, err := s.entityAdapter.FindConnectorForComponent(ctx, entity.ID)
	if err != nil || connector == nil {
		return nil, err
	}

	s.connectors[connector.ID] = *connector

	return connector, nil
}

func (s *baseService) findConnectorForResource(ctx context.Context, entity types.Entity) (*types.Entity, error) {
	for _, id := range entity.Depends {
		if connector, ok := s.connectors[id]; ok {
			return &connector, nil
		}
	}

	connector, err := s.entityAdapter.FindConnectorForResource(ctx, entity.ID)
	if err != nil || connector == nil {
		return nil, err
	}

	s.connectors[connector.ID] = *connector

	return connector, nil
}

func (s *baseService) findComponentForResource(ctx context.Context, entity types.Entity) (*types.Entity, error) {
	for _, id := range entity.Impacts {
		if component, ok := s.components[id]; ok {
			return &component, nil
		}
	}

	component, err := s.entityAdapter.FindComponentForResource(ctx, entity.ID)
	if err != nil || component == nil {
		return nil, err
	}

	s.components[component.ID] = *component

	return component, nil
}

func pbhResOneOfType(pbhRes libpbehavior.ResolveResult, types []string) bool {
	for _, canonicalType := range types {
		if canonicalType == libpbehavior.TypeActive && pbhRes.ResolvedType == nil ||
			pbhRes.ResolvedType != nil && pbhRes.ResolvedType.Type == canonicalType {
			return true
		}
	}

	return false
}
