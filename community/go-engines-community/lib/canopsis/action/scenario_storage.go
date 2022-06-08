package action

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
	"strings"
	"sync"
)

type scenarioStorage struct {
	adapter                Adapter
	delayedScenarioManager DelayedScenarioManager
	scenariosMx            sync.RWMutex
	scenarios              []Scenario
	logger                 zerolog.Logger
}

func NewScenarioStorage(
	actionAdapter Adapter,
	delayedScenarioManager DelayedScenarioManager,
	logger zerolog.Logger,
) ScenarioStorage {
	return &scenarioStorage{
		adapter:                actionAdapter,
		delayedScenarioManager: delayedScenarioManager,
		logger:                 logger,
	}
}

func (s *scenarioStorage) ReloadScenarios(ctx context.Context) error {
	s.scenariosMx.Lock()
	defer s.scenariosMx.Unlock()

	var err error
	scenarios, err := s.adapter.GetEnabled(ctx)

	if err != nil {
		return err
	}

	validScenarios := make([]Scenario, 0)
	validScenariosIDs := make([]string, 0)

	for _, scenario := range scenarios {
		valid := true
		for i, action := range scenario.Actions {
			if !action.OldEntityPatterns.IsSet() && !action.OldAlarmPatterns.IsSet() &&
				len(action.AlarmPattern) == 0 && len(action.EntityPattern) == 0 {
				s.logger.Warn().Str("scenario", scenario.ID).Int("action number", i).Msg("action doesn't have patterns")
				valid = false

				break
			}

			if !action.OldEntityPatterns.IsValid() {
				s.logger.Warn().Str("scenario", scenario.ID).Int("action number", i).Msg("failed to parse entity patterns")
				valid = false

				break
			}

			if !action.OldAlarmPatterns.IsValid() {
				s.logger.Warn().Str("scenario", scenario.ID).Int("action number", i).Msg("failed to parse alarm patterns")
				valid = false

				break
			}
		}

		if valid {
			validScenarios = append(validScenarios, scenario)
			validScenariosIDs = append(validScenariosIDs, scenario.ID)
		}
	}

	s.scenarios = validScenarios
	s.logger.Debug().Str("scenarios", strings.Join(validScenariosIDs, ", ")).Int("number", len(s.scenarios)).Msg("Successfully loaded scenarios")

	return nil
}

func (s *scenarioStorage) GetTriggeredScenarios(
	triggers []string,
	alarm types.Alarm,
) ([]Scenario, error) {
	s.scenariosMx.RLock()
	defer s.scenariosMx.RUnlock()

	triggeredScenarios := make([]Scenario, 0)

	for _, scenario := range s.scenarios {
		if alarm.Value.PbehaviorInfo.OneOf(scenario.DisableDuringPeriods) {
			continue
		}

		if !scenario.IsTriggered(triggers) {
			continue
		}

		if scenario.Delay != nil && scenario.Delay.Value > 0 {
			continue
		}

		triggeredScenarios = append(triggeredScenarios, scenario)
	}

	return triggeredScenarios, nil
}

func (s *scenarioStorage) RunDelayedScenarios(
	ctx context.Context,
	triggers []string,
	alarm types.Alarm,
	entity types.Entity,
	additionalData AdditionalData,
) error {
	s.scenariosMx.RLock()
	defer s.scenariosMx.RUnlock()

	for _, scenario := range s.scenarios {
		if alarm.Value.PbehaviorInfo.OneOf(scenario.DisableDuringPeriods) {
			continue
		}

		if !scenario.IsTriggered(triggers) {
			continue
		}

		if scenario.Delay != nil && scenario.Delay.Value > 0 {
			// Check if at least on action matches alarm.
			matched := false
			var err error

			for idx, action := range scenario.Actions {
				if action.OldAlarmPatterns.IsSet() {
					if !action.OldAlarmPatterns.IsValid() {
						s.logger.Warn().Msgf("Action %d from scenario %s has an invalid old alarm pattern, skip", idx, scenario.ID)
						continue
					}

					matched = action.OldAlarmPatterns.Matches(&alarm)
				} else {
					matched, err = action.AlarmPattern.Match(alarm)
					if err != nil {
						s.logger.Err(err).Msgf("Action %d from scenario %s alarm pattern match returned error", idx, scenario.ID)
						continue
					}
				}

				if !matched {
					if action.DropScenarioIfNotMatched {
						break
					}

					continue
				}

				if action.OldEntityPatterns.IsSet() {
					if !action.OldEntityPatterns.IsValid() {
						s.logger.Warn().Msgf("Action %d from scenario %s has an invalid old entity pattern, skip", idx, scenario.ID)
						continue
					}

					matched = action.OldEntityPatterns.Matches(&entity)
				} else {
					matched, _, err = action.EntityPattern.Match(entity)
					if err != nil {
						s.logger.Err(err).Msgf("Action %d from scenario %s entity pattern match returned error", idx, scenario.ID)
						continue
					}
				}

				if matched || action.DropScenarioIfNotMatched {
					break
				}
			}

			if matched {
				err := s.delayedScenarioManager.AddDelayedScenario(ctx, alarm, scenario, additionalData)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (s *scenarioStorage) GetScenario(id string) *Scenario {
	s.scenariosMx.RLock()
	defer s.scenariosMx.RUnlock()

	for _, scenario := range s.scenarios {
		if scenario.ID == id {
			return &scenario
		}
	}

	return nil
}
