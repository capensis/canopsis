package alarm

import (
	"context"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/resolverule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
)

type service struct {
	adapter            Adapter
	resolveRuleAdapter resolverule.Adapter
	alarmStatusService alarmstatus.Service
	logger             zerolog.Logger
}

// NewService gives the correct alarm adapter. Give nil to the redis
// client and it will create a new redis.Client with the dedicated redis
// database for alarms.
func NewService(
	alarmAdapter Adapter,
	resolveRuleAdapter resolverule.Adapter,
	alarmStatusService alarmstatus.Service,
	logger zerolog.Logger,
) Service {
	return &service{
		adapter:            alarmAdapter,
		resolveRuleAdapter: resolveRuleAdapter,
		alarmStatusService: alarmStatusService,
		logger:             logger,
	}
}

func (s *service) ResolveClosed(ctx context.Context) ([]types.AlarmWithEntity, error) {
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

	alarmsToResolve := make([]types.AlarmWithEntity, 0)
	for cursor.Next(ctx) {
		alarmWithEntity := types.AlarmWithEntity{}
		if err := cursor.Decode(&alarmWithEntity); err != nil {
			s.logger.Error().Err(err).Msg("cannot decode alarm with entity")
			continue
		}

		for _, rule := range rules {
			matched, err := rule.Matches(alarmWithEntity)
			if err != nil {
				s.logger.Error().Err(err).Str("resolve_rule", rule.ID).Msg("match resolve rule returned error, skip")
				continue
			}

			if matched {
				alarmState := alarmWithEntity.Alarm.Value.State.Value

				if alarmState == types.AlarmStateOK {
					lastStep := alarmWithEntity.Alarm.Value.Steps[len(alarmWithEntity.Alarm.Value.Steps)-1]
					before := rule.Duration.SubFrom(now)

					if lastStep.Timestamp.Before(before) {
						alarmsToResolve = append(alarmsToResolve, alarmWithEntity)
					}
				}

				break
			}
		}
	}

	return alarmsToResolve, nil
}

func (s *service) ResolveCancels(ctx context.Context, alarmConfig config.AlarmConfig) ([]types.AlarmWithEntity, error) {
	canceledAlarms := make([]types.AlarmWithEntity, 0)

	alarms, err := s.adapter.GetAlarmsWithCancelMark(ctx)
	if err != nil {
		return canceledAlarms, fmt.Errorf("cannot fetch alarms: %w", err)
	}

	for _, alarm := range alarms {
		if time.Since(alarm.Alarm.Value.Canceled.Timestamp.Time) >= alarmConfig.CancelAutosolveDelay {
			canceledAlarms = append(canceledAlarms, alarm)
		}
	}

	return canceledAlarms, nil
}

func (s *service) ResolveSnoozes(ctx context.Context, alarmConfig config.AlarmConfig) ([]types.AlarmWithEntity, error) {
	unsnoozedAlarms := make([]types.AlarmWithEntity, 0)

	alarms, err := s.adapter.GetAlarmsWithSnoozeMark(ctx)
	if err != nil {
		return unsnoozedAlarms, fmt.Errorf("cannot fetch alarms: %w", err)
	}

	for _, alarm := range alarms {
		if !alarm.Alarm.IsSnoozed() && (alarm.Alarm.IsInActivePeriod() || alarmConfig.DisableActionSnoozeDelayOnPbh) {
			unsnoozedAlarms = append(unsnoozedAlarms, alarm)
		}
	}

	return unsnoozedAlarms, nil
}

func (s *service) UpdateFlappingAlarms(ctx context.Context) ([]types.AlarmWithEntity, error) {
	updatedAlarms := make([]types.AlarmWithEntity, 0)

	flappingAlarms, err := s.adapter.GetAlarmsWithFlappingStatus(ctx)
	if err != nil {
		return updatedAlarms, fmt.Errorf("cannot fetch alarms: %w", err)
	}

	for _, alarm := range flappingAlarms {
		currentAlarmStatus := alarm.Alarm.Value.Status.Value
		newStatus := s.alarmStatusService.ComputeStatus(alarm.Alarm, alarm.Entity)

		if newStatus != currentAlarmStatus {
			updatedAlarms = append(updatedAlarms, alarm)
		}
	}

	return updatedAlarms, nil
}
