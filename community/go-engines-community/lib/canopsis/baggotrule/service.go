package baggotrule

//go:generate mockgen -destination=../../../mocks/lib/canopsis/baggotrule/service.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/baggotrule Service

import (
	"context"
	"fmt"
	"runtime/trace"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
)

// Service interface is used to implement alarms modification by idle baggotRules.
type Service interface {
	Process(ctx context.Context) ([]types.Alarm, error)
}

type service struct {
	ruleAdapter  Adapter
	alarmAdapter alarm.Adapter
	logger       zerolog.Logger
}

// NewService creates new service.
func NewService(
	ruleAdapter Adapter,
	alarmAdapter alarm.Adapter,
	logger zerolog.Logger,
) Service {
	return &service{
		ruleAdapter:  ruleAdapter,
		alarmAdapter: alarmAdapter,
		logger:       logger,
	}
}

func (s *service) match(rules []Rule, alr *types.AlarmWithEntity) *Rule {
	for _, r := range rules {
		if r.Matches(alr) {
			return &r
		}
	}
	return nil
}

func (s *service) Process(ctx context.Context) ([]types.Alarm, error) {
	defer trace.StartRegion(ctx, "alarms.ResolveAlarms").End()

	rules, err := s.ruleAdapter.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("canont fetch baggot rules: %w", err)
	}
	if len(rules) == 0 {
		return nil, nil
	}

	cursor, err := s.alarmAdapter.GetOpenedAlarmsWithEntity(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch open alarms: %w", err)
	}
	defer cursor.Close(ctx)

	updatedAlarms := make([]types.Alarm, 0)
	for cursor.Next(ctx) {
		alrEntity := types.AlarmWithEntity{}
		if err := cursor.Decode(&alrEntity); err != nil {
			s.logger.Error().Err(err).Msg("cannot decode alarm with entity")
			continue
		}

		rule := s.match(rules, &alrEntity)
		if rule != nil && alrEntity.Alarm.Closable(rule.Duration.Duration()) {
			updatedAlarms = append(updatedAlarms, alrEntity.Alarm)
		}
	}

	return updatedAlarms, nil
}
