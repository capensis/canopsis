package alarmstatus

//go:generate mockgen -destination=../../../mocks/lib/canopsis/alarmstatus/alarmstatus.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus Service

import (
	"context"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/flappingrule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type Service interface {
	Load(ctx context.Context) error
	ComputeStatus(alarm types.Alarm, entity types.Entity) types.CpsNumber
}

func NewService(flappingRuleAdapter flappingrule.Adapter, configProvider config.AlarmConfigProvider) Service {
	return &service{
		flappingRuleAdapter: flappingRuleAdapter,
		configProvider:      configProvider,
	}
}

type service struct {
	flappingRuleAdapter flappingrule.Adapter
	configProvider      config.AlarmConfigProvider

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

	s.flappingRules = rules
	return nil
}

func (s *service) ComputeStatus(alarm types.Alarm, entity types.Entity) types.CpsNumber {
	if alarm.Value.Canceled != nil && alarm.Value.Resolved == nil {
		return types.AlarmStatusCancelled
	}

	if s.isFlapping(alarm, entity) {
		return types.AlarmStatusFlapping
	}

	if s.isStealthy(alarm) {
		return types.AlarmStatusStealthy
	}

	if alarm.Value.State != nil && alarm.Value.State.Value != types.AlarmStateOK {
		return types.AlarmStatusOngoing
	}

	return types.AlarmStatusOff
}

func (s *service) isFlapping(alarm types.Alarm, entity types.Entity) bool {
	s.flappingRulesMx.RLock()
	defer s.flappingRulesMx.RUnlock()

	alarmWithEntity := types.AlarmWithEntity{
		Alarm:  alarm,
		Entity: entity,
	}
	lastStepType := ""
	freq := 0
	for _, rule := range s.flappingRules {
		if rule.Matches(alarmWithEntity) {
			for i := len(alarm.Value.Steps) - 1; i >= 0; i-- {
				step := alarm.Value.Steps[i]
				if time.Since(step.Timestamp.Time) >= rule.Duration.Duration() {
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
					return true
				}
			}

			break
		}
	}

	return false
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
