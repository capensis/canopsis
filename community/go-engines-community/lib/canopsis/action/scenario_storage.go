package action

import (
	"context"
	"strings"
	"sync"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
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
	scenarios, err := s.adapter.GetEnabled(ctx)
	if err != nil {
		return err
	}

	validScenarios := make([]Scenario, 0)
	validScenariosIDs := make([]string, 0)

	for _, scenario := range scenarios {
		valid := true
		for i, action := range scenario.Actions {
			if len(action.AlarmPattern) == 0 && len(action.EntityPattern) == 0 {
				s.logger.Warn().Str("scenario", scenario.ID).Int("action_number", i).Msg("action doesn't have patterns")
				valid = false

				break
			}
		}

		if valid {
			validScenarios = append(validScenarios, scenario)
			validScenariosIDs = append(validScenariosIDs, scenario.ID)
		}
	}

	s.scenariosMx.Lock()
	defer s.scenariosMx.Unlock()

	s.scenarios = validScenarios
	s.logger.Debug().Str("scenarios", strings.Join(validScenariosIDs, ", ")).Int("number", len(s.scenarios)).Msg("Successfully loaded scenarios")

	return nil
}

func (s *scenarioStorage) GetTriggeredScenarios(
	triggers []string,
	alarm types.Alarm,
) (map[string][]Scenario, error) {
	s.scenariosMx.RLock()
	defer s.scenariosMx.RUnlock()

	triggeredScenarios := make(map[string][]Scenario, 0)

	for _, scenario := range s.scenarios {
		if alarm.Value.PbehaviorInfo.OneOf(scenario.DisableDuringPeriods) {
			continue
		}

		trigger := scenario.IsTriggered(triggers)
		if trigger == "" {
			continue
		}

		if scenario.Delay != nil && scenario.Delay.Value > 0 {
			continue
		}

		triggeredScenarios[trigger] = append(triggeredScenarios[trigger], scenario)
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

		trigger := scenario.IsTriggered(triggers)
		if trigger == "" {
			continue
		}

		if scenario.Delay != nil && scenario.Delay.Value > 0 {
			// Check if at least on action matches alarm.
			matched := false
			var err error

			for idx, action := range scenario.Actions {
				matched, err = action.Match(entity, alarm)
				if err != nil {
					s.logger.Err(err).Msgf("match action %d from scenario %s returned error", idx, scenario.ID)
					break
				}

				if matched || action.DropScenarioIfNotMatched {
					break
				}
			}

			if matched {
				additionalData.Trigger = trigger
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
