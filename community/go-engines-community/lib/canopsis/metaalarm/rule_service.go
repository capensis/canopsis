package metaalarm

import (
	"context"
	"sync"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
)

type ruleService struct {
	// rulesAdapter is an rulesAdapter to the rules collection.
	rulesAdapter      RulesAdapter
	ruleEntityCounter RuleEntityCounter
	valueGroupEntityCounter ValueGroupEntityCounter

	ruleApplicatorContainer *Container

	// rules is an array of rules.
	rules []Rule

	// countMutex is a mutex preventing to overwrite the rules array while
	// it's being used.
	rulesMutex sync.RWMutex

	logger zerolog.Logger
}

type zeroStateEvent struct {
}

func (z zeroStateEvent) Error() string {
	return "manual_metaalarm_group event has zero state value"
}

func (s *ruleService) LoadRules(ctx context.Context) error {
	rules, err := s.rulesAdapter.Get()
	if err != nil {
		return err
	}

	s.rulesMutex.Lock()
	defer s.rulesMutex.Unlock()

	for _, rule := range rules {
		err := s.ruleEntityCounter.CountTotalEntitiesAmount(ctx, rule)
		if err != nil {
			return err
		}

		err = s.valueGroupEntityCounter.CountTotalEntitiesAmount(ctx, rule)
		if err != nil {
			return err
		}
	}

	s.rules = rules

	s.logger.Info().Int("number", len(s.rules)).Msg("Successfully loaded meta-alarm rules")

	return nil
}

func (s *ruleService) ProcessEvent(ctx context.Context, event types.Event) ([]types.Event, error) {
	var metaAlarmEvents []types.Event

	s.rulesMutex.RLock()
	defer s.rulesMutex.RUnlock()

	s.logger.Debug().Msg("Process event by meta-alarm rule service")

	if event.EventType == types.EventManualMetaAlarmGroup || event.EventType == types.EventManualMetaAlarmUpdate ||
		event.EventType == types.EventManualMetaAlarmUngroup {

		if event.State == 0 && event.EventType == types.EventManualMetaAlarmGroup {
			return metaAlarmEvents, zeroStateEvent{}
		}

		applicator, _ := s.ruleApplicatorContainer.Get(RuleManualGroup)
		events, err := applicator.Apply(ctx, event, Rule{})
		if err != nil {
			return metaAlarmEvents, err
		}
		metaAlarmEvents = append(metaAlarmEvents, events...)
		return metaAlarmEvents, nil
	}

	// avoid state = 0 events
	if event.EventType != types.EventTypeCheck || event.Alarm == nil || !event.Alarm.IsMalfunctioning() {
		return metaAlarmEvents, nil
	}

	for _, rule := range s.rules {
		s.logger.Debug().Msgf("Meta-alarm rule service: check rule %s", rule.ID)

		isMatched := rule.Patterns.IsMatched(event)
		s.logger.Debug().Msgf("Meta-alarm rule service: check rule %t", isMatched)

		if isMatched {
			applicator, found := s.ruleApplicatorContainer.Get(rule.Type)
			if !found {
				s.logger.Warn().Msgf("RuleApplicator for %s is not exist", rule.Type)
			} else {
				metaAlarmEvent, err := applicator.Apply(ctx, event, rule)
				if err != nil {
					s.logger.Err(err).
						Str("rule_id", rule.ID).
						Str("alarm_id", event.Alarm.ID).
						Msgf("failed to apply meta-alarm rule")
				}

				metaAlarmEvents = append(metaAlarmEvents, metaAlarmEvent...)
			}
		}
	}

	return metaAlarmEvents, nil
}

func NewRulesService(ruleAdapter RulesAdapter, ruleEntityCounter RuleEntityCounter, valueGroupEntityCounter ValueGroupEntityCounter, container *Container, logger zerolog.Logger) RulesService {
	return &ruleService{
		rulesAdapter:            ruleAdapter,
		ruleEntityCounter:       ruleEntityCounter,
		valueGroupEntityCounter: valueGroupEntityCounter,
		ruleApplicatorContainer: container,
		rules:                   nil,
		rulesMutex:              sync.RWMutex{},
		logger:                  logger,
	}
}
