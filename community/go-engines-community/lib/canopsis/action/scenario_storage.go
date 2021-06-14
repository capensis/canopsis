package action

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
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

func (s *scenarioStorage) ReloadScenarios() error {
	s.scenariosMx.Lock()
	defer s.scenariosMx.Unlock()

	var err error
	scenarios, err := s.adapter.GetEnabled()

	if err != nil {
		return err
	}

	validScenarios := make([]Scenario, 0)
	validScenariosIDs := make([]string, 0)

	for _, scenario := range scenarios {
		valid := true
		for i, action := range scenario.Actions {
			if !action.EntityPatterns.IsValid() {
				s.logger.Warn().Str("scenario", scenario.ID).Int("action number", i).Msg("failed to parse entity patterns")
				valid = false

				break
			}

			if !action.AlarmPatterns.IsValid() {
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

		if scenario.Delay != nil && scenario.Delay.Seconds > 0 {
			continue
		}

		triggeredScenarios = append(triggeredScenarios, scenario)
	}

	return triggeredScenarios, nil
}

func (s *scenarioStorage) RunDelayedScenarios(
	triggers []string,
	alarm types.Alarm,
	entity types.Entity,
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

		if scenario.Delay != nil && scenario.Delay.Seconds > 0 {
			// Check if at least on action matches alarm.
			matched := false
			for _, action := range scenario.Actions {
				if action.AlarmPatterns.Matches(&alarm) && action.EntityPatterns.Matches(&entity) {
					matched = true
					break
				} else if action.DropScenarioIfNotMatched {
					break
				}
			}

			if matched {
				err := s.delayedScenarioManager.AddDelayedScenario(alarm, scenario)
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
