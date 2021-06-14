/*
Package idlealarm implements alarm modification on idle alarm.
*/
package idlealarm

//go:generate mockgen -destination=../../../mocks/lib/canopsis/idlealarm/service.go git.canopsis.net/canopsis/go-engines/lib/canopsis/idlealarm Service

import (
	alarmlib "git.canopsis.net/canopsis/go-engines/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/idlerule"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"github.com/rs/zerolog"
	"time"
)

// Service interface is used to implement alarms modification by idle rules.
type Service interface {
	Process() []types.Event
}

// NewService creates new service.
func NewService(
	ruleAdapter idlerule.RuleAdapter,
	alarmAdapter alarmlib.Adapter,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Service {
	return &baseService{
		ruleAdapter:  ruleAdapter,
		alarmAdapter: alarmAdapter,
		encoder:      encoder,
		logger:       logger,
	}
}

type baseService struct {
	ruleAdapter  idlerule.RuleAdapter
	alarmAdapter alarmlib.Adapter
	encoder      encoding.Encoder
	logger       zerolog.Logger
}

func (s *baseService) Process() []types.Event {
	rulesByType := s.getRulesByType()

	if rulesByType == nil {
		return nil
	}

	var events []types.Event
	for ruleType, rules := range rulesByType {
		minDuration := s.getMinDuration(rules)

		if minDuration == nil {
			continue
		}

		var alarms []types.AlarmWithEntity
		var err error
		date := time.Now().Add(-*minDuration)

		switch ruleType {
		case idlerule.RuleTypeLastEvent:
			alarms, err = s.alarmAdapter.GetOpenedAlarmsWithLastEventDateBefore(date)
		case idlerule.RuleTypeLastUpdate:
			alarms, err = s.alarmAdapter.GetOpenedAlarmsWithLastUpdateDateBefore(date)
		default:
			s.logger.Error().Msg("failed to process unknown rule type")
			continue
		}

		if err != nil {
			s.logger.Error().Err(err).Msg("failed to fetch opened alarms")
			continue
		}

		events = append(events, s.processAlarms(alarms, rules)...)
	}

	s.logger.Info().Msg("process idle alarms")

	return events
}

// getRulesByType fetches rules from db and checks if rules are valid.
func (s *baseService) getRulesByType() map[idlerule.RuleType][]idlerule.Rule {
	rules, err := s.ruleAdapter.Get()

	if err != nil {
		s.logger.Error().Err(err).Msg("failed to fetch idle rules")

		return nil
	}

	validRulesByType := make(map[idlerule.RuleType][]idlerule.Rule)

	for _, rule := range rules {
		if !rule.AlarmPatterns.IsValid() {
			s.logger.Warn().Str("idlerule", rule.ID).Msg("failed to parse alarm patterns")
			continue
		}

		if !rule.EntityPatterns.IsValid() {
			s.logger.Warn().Str("idlerule", rule.ID).Msg("failed to parse entity patterns")
			continue
		}

		if _, ok := validRulesByType[rule.Type]; !ok {
			validRulesByType[rule.Type] = make([]idlerule.Rule, 0)
		}

		validRulesByType[rule.Type] = append(validRulesByType[rule.Type], rule)
	}

	return validRulesByType
}

// getMinDuration finds min duration by all rules.
func (s *baseService) getMinDuration(rules []idlerule.Rule) *time.Duration {
	var min *time.Duration

	for _, rule := range rules {
		duration := rule.Duration.Duration()

		if min == nil || duration < *min {
			min = &duration
		}
	}

	return min
}

// processAlarms applies rule operations.
func (s *baseService) processAlarms(
	alarms []types.AlarmWithEntity,
	rules []idlerule.Rule,
) []types.Event {
	var events []types.Event
	for _, alarm := range alarms {
		events = append(events, s.processAlarm(alarm, rules)...)
	}

	return events
}

func (s *baseService) processAlarm(
	alarm types.AlarmWithEntity,
	rules []idlerule.Rule,
) []types.Event {
	var events []types.Event

	for _, rule := range rules {
		// Skip rule on defined periods
		pbhInfo := alarm.Alarm.Value.PbehaviorInfo
		if len(rule.DisableDuringPeriods) > 0 {
			if pbhInfo.OneOf(rule.DisableDuringPeriods) {
				s.logger.
					Debug().
					Str("alarm", alarm.Alarm.ID).
					Str("idlerule", rule.ID).
					Str("pbh_id", pbhInfo.ID).
					Str("pbh_type", pbhInfo.CanonicalType).
					Msg("skip rule")
				continue
			}
		}

		if rule.Matches(&alarm.Alarm, &alarm.Entity) {
			event := types.Event{
				Connector:     alarm.Alarm.Value.Connector,
				ConnectorName: alarm.Alarm.Value.ConnectorName,
				Component:     alarm.Alarm.Value.Component,
				Resource:      alarm.Alarm.Value.Resource,
				Timestamp:     types.CpsTime{Time: time.Now()},
				Initiator:     types.InitiatorSystem,
			}
			if event.Resource == "" {
				event.SourceType = types.SourceTypeComponent
			} else {
				event.SourceType = types.SourceTypeResource
			}

			switch rule.Operation.Type {
			case types.ActionTypeAck, types.ActionTypeAckRemove, types.ActionTypeCancel:
				if params, ok := rule.Operation.Parameters.(types.OperationParameters); ok {
					event.EventType = types.EventTypeAck
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
					output, err := s.encoder.Encode(params)
					if err != nil {
						s.logger.Err(err).Msg("cannot encode parameters")
						continue
					}
					event.Output = string(output)
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
				s.logger.Error().Msgf("unknown operation %v", rule.Operation)
				continue
			}

			if event.EventType == "" {
				s.logger.Error().Msgf("invalid operation %+v", rule.Operation)
				continue
			}

			events = append(events, event)
		}
	}

	return events
}
