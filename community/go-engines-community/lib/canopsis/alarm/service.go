package alarm

import (
	"context"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	libevent "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/resolverule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
)

type service struct {
	adapter            Adapter
	resolveRuleAdapter resolverule.Adapter
	alarmStatusService alarmstatus.Service
	eventGenerator     libevent.Generator
	logger             zerolog.Logger
}

// NewService gives the correct alarm adapter. Give nil to the redis
// client and it will create a new redis.Client with the dedicated redis
// database for alarms.
func NewService(
	alarmAdapter Adapter,
	resolveRuleAdapter resolverule.Adapter,
	alarmStatusService alarmstatus.Service,
	eventGenerator libevent.Generator,
	logger zerolog.Logger,
) Service {
	return &service{
		adapter:            alarmAdapter,
		resolveRuleAdapter: resolveRuleAdapter,
		alarmStatusService: alarmStatusService,
		eventGenerator:     eventGenerator,
		logger:             logger,
	}
}

func (s *service) ResolveClosed(ctx context.Context) ([]types.Event, error) {
	now := datetime.NewCpsTime()
	rules, err := s.resolveRuleAdapter.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("canont fetch resolve rules: %w", err)
	}

	ids := make([]string, len(rules))
	for i, rule := range rules {
		ids[i] = rule.ID
	}

	s.logger.Debug().Strs("rules", ids).Msg("load resolve rules")
	if len(rules) == 0 {
		return nil, nil
	}

	cursor, err := s.adapter.GetOpenedAlarmsWithEntity(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch open alarms: %w", err)
	}

	defer cursor.Close(ctx)
	events := make([]types.Event, 0)
	for cursor.Next(ctx) {
		alarm := types.AlarmWithEntity{}
		if err := cursor.Decode(&alarm); err != nil {
			s.logger.Error().Err(err).Msg("cannot decode alarm with entity")
			continue
		}

		for _, rule := range rules {
			matched, err := rule.Matches(alarm)
			if err != nil {
				s.logger.Error().Err(err).Str("resolve_rule", rule.ID).Msg("match resolve rule returned error, skip")
				continue
			}

			if matched {
				alarmState := alarm.Alarm.Value.State.Value
				if alarmState == types.AlarmStateOK {
					lastStep := alarm.Alarm.Value.Steps[len(alarm.Alarm.Value.Steps)-1]
					before := rule.Duration.SubFrom(now)

					if lastStep.Timestamp.Before(before) {
						event, err := s.eventGenerator.Generate(alarm.Entity)
						if err != nil {
							s.logger.Err(err).Msg("cannot generate event")
							continue
						}

						event.EventType = types.EventTypeResolveClose
						event.Timestamp = datetime.NewCpsTime()
						event.Output = types.RuleNameRulePrefix + rule.Name
						events = append(events, event)
					}
				}

				break
			}
		}
	}

	return events, nil
}

func (s *service) ResolveCancels(ctx context.Context, alarmConfig config.AlarmConfig) ([]types.Event, error) {
	events := make([]types.Event, 0)
	alarms, err := s.adapter.GetAlarmsWithCancelMark(ctx)
	if err != nil {
		return events, fmt.Errorf("cannot fetch alarms: %w", err)
	}

	for _, alarm := range alarms {
		if time.Since(alarm.Alarm.Value.Canceled.Timestamp.Time) >= alarmConfig.CancelAutosolveDelay {
			event, err := s.eventGenerator.Generate(alarm.Entity)
			if err != nil {
				s.logger.Err(err).Msg("cannot generate event")
				continue
			}

			event.EventType = types.EventTypeResolveCancel
			event.Timestamp = datetime.NewCpsTime()
			events = append(events, event)
		}
	}

	return events, nil
}

func (s *service) ResolveSnoozes(ctx context.Context, alarmConfig config.AlarmConfig) ([]types.Event, error) {
	events := make([]types.Event, 0)
	alarms, err := s.adapter.GetAlarmsWithSnoozeMark(ctx)
	if err != nil {
		return events, fmt.Errorf("cannot fetch alarms: %w", err)
	}

	for _, alarm := range alarms {
		if !alarm.Alarm.IsSnoozed() && (alarm.Alarm.IsInActivePeriod() || alarmConfig.DisableActionSnoozeDelayOnPbh) {
			event, err := s.eventGenerator.Generate(alarm.Entity)
			if err != nil {
				s.logger.Err(err).Msg("cannot generate event")
				continue
			}

			event.EventType = types.EventTypeUnsnooze
			event.Timestamp = datetime.NewCpsTime()
			events = append(events, event)
		}
	}

	return events, nil
}

func (s *service) UpdateFlappingAlarms(ctx context.Context) ([]types.Event, error) {
	events := make([]types.Event, 0)
	flappingAlarms, err := s.adapter.GetAlarmsWithFlappingStatus(ctx)
	if err != nil {
		return events, fmt.Errorf("cannot fetch alarms: %w", err)
	}

	for _, alarm := range flappingAlarms {
		currentAlarmStatus := alarm.Alarm.Value.Status.Value
		newStatus, _ := s.alarmStatusService.ComputeStatus(alarm.Alarm, alarm.Entity)
		if newStatus != currentAlarmStatus {
			event, err := s.eventGenerator.Generate(alarm.Entity)
			if err != nil {
				s.logger.Err(err).Msg("cannot generate event")
				continue
			}

			event.EventType = types.EventTypeUpdateStatus
			event.Timestamp = datetime.NewCpsTime()
			events = append(events, event)
		}
	}

	return events, nil
}
