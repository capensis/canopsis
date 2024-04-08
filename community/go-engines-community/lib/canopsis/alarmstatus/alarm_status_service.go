package alarmstatus

//go:generate mockgen -destination=../../../mocks/lib/canopsis/alarmstatus/alarmstatus.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus Service

import (
	"context"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/flappingrule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
)

type Service interface {
	Load(ctx context.Context) error
	ComputeStatus(alarm types.Alarm, entity types.Entity) (state types.CpsNumber, ruleName string)
}

func NewService(
	flappingRuleAdapter flappingrule.Adapter,
	configProvider config.AlarmConfigProvider,
	logger zerolog.Logger,
) Service {
	return &service{
		flappingRuleAdapter: flappingRuleAdapter,
		configProvider:      configProvider,
		logger:              logger,
	}
}

type service struct {
	flappingRuleAdapter flappingrule.Adapter
	configProvider      config.AlarmConfigProvider
	logger              zerolog.Logger

	flappingRulesMx sync.RWMutex
	flappingRules   []flappingrule.Rule
}

func (s *service) Load(ctx context.Context) error {
	s.flappingRulesMx.Lock()
	defer s.flappingRulesMx.Unlock()

	rules, err := s.flappingRuleAdapter.Get(ctx)
	if err != nil {
		return err
	}

	ids := make([]string, len(rules))
	for i, rule := range rules {
		ids[i] = rule.ID
	}
	s.logger.Debug().Strs("rules", ids).Msg("load flapping rules")

	s.flappingRules = rules
	return nil
}

func (s *service) ComputeStatus(alarm types.Alarm, entity types.Entity) (types.CpsNumber, string) {
	if alarm.Value.Canceled != nil && alarm.Value.Resolved == nil {
		return types.AlarmStatusCancelled, ""
	}

	if ok, ruleName := s.isFlapping(alarm, entity); ok {
		return types.AlarmStatusFlapping, ruleName
	}

	if s.isStealthy(alarm) {
		return types.AlarmStatusStealthy, ""
	}

	if alarm.Value.State != nil {
		alarmState := alarm.Value.State.Value
		if alarmState != types.AlarmStateOK {
			return types.AlarmStatusOngoing, ""
		}
	}

	return types.AlarmStatusOff, ""
}

func (s *service) isFlapping(alarm types.Alarm, entity types.Entity) (bool, string) {
	s.flappingRulesMx.RLock()
	defer s.flappingRulesMx.RUnlock()

	now := datetime.NewCpsTime()
	alarmWithEntity := types.AlarmWithEntity{
		Alarm:  alarm,
		Entity: entity,
	}
	lastStepType := ""
	freq := 0
	for _, rule := range s.flappingRules {
		matched, err := rule.Matches(alarmWithEntity)
		if err != nil {
			s.logger.Error().Err(err).Str("flapping_rule", rule.ID).Msg("match flapping rule returned error, skip")
			continue
		}

		if matched {
			before := rule.Duration.SubFrom(now)

			for i := len(alarm.Value.Steps) - 1; i >= 0; i-- {
				step := alarm.Value.Steps[i]

				if step.Timestamp.Before(before) {
					break
				}

				if step.Type != lastStepType {
					switch step.Type {
					case types.AlarmStepStateIncrease, types.AlarmStepStateDecrease:
						lastStepType = step.Type
						freq++
					}
				}

				if freq > rule.FreqLimit {
					return true, types.RuleNameRulePrefix + rule.Name
				}
			}

			break
		}
	}

	return false, ""
}

func (s *service) isStealthy(alarm types.Alarm) bool {
	interval := s.configProvider.Get().StealthyInterval

	for i := len(alarm.Value.Steps) - 1; i >= 0; i-- {
		step := alarm.Value.Steps[i]
		if time.Since(step.Timestamp.Time) >= interval {
			break
		}

		if step.Value != types.AlarmStateOK {
			switch step.Type {
			case types.AlarmStepStatusIncrease, types.AlarmStepStateDecrease:
				return true
			default:
				break
			}
		}
	}

	return false
}
