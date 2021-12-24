package eventfilter

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
	"sync"
)

type ruleService struct {
	rulesAdapter            RuleAdapter
	ruleApplicatorContainer RuleApplicatorContainer
	rules                   []Rule
	rulesMutex              sync.RWMutex
	timezoneConfigProvider  config.TimezoneConfigProvider
	logger                  zerolog.Logger
}

func (s *ruleService) LoadRules(ctx context.Context, types []string) error {
	s.rulesMutex.Lock()
	defer s.rulesMutex.Unlock()

	var err error

	s.rules, err = s.rulesAdapter.GetByTypes(ctx, types)
	if err != nil {
		return err
	}

	ids := make([]string, len(s.rules))
	for i := 0; i < len(s.rules); i++ {
		ids[i] = s.rules[i].ID
	}

	s.logger.Debug().Strs("rules", ids).Msg("Loading event filter rules")

	return nil
}

func (s *ruleService) ProcessEvent(ctx context.Context, event types.Event) (types.Event, error) {
	s.rulesMutex.RLock()
	defer s.rulesMutex.RUnlock()

	outcome := OutcomePass
	tz := s.timezoneConfigProvider.Get()

	var err error
	for _, rule := range s.rules {
		if outcome != OutcomePass {
			break
		}

		regexMatches, match := rule.Patterns.GetRegexMatches(event)
		if !match {
			if event.Debug {
				s.logger.Info().Str("rule", rule.ID).Str("event_type", event.EventType).Str("entity", event.GetEID()).Msg("Event filter rule service: rule is not matched")
			}

			continue
		}

		if event.Debug {
			s.logger.Info().Str("rule", rule.ID).Str("event_type", event.EventType).Str("entity", event.GetEID()).Msg("Event filter rule service: rule is matched")
		}

		applicator, found := s.ruleApplicatorContainer.Get(rule.Type)
		if !found {
			s.logger.Warn().Str("rule_id", rule.ID).Str("rule_type", rule.Type).Msg("Event filter rule service: RuleApplicator doesn't exist")
			continue
		}

		outcome, event, err = applicator.Apply(ctx, rule, event, regexMatches, &tz)
		if err != nil {
			s.logger.Err(err).Str("rule_id", rule.ID).Str("rule_type", rule.Type).Msg("Event filter rule service: failed to apply")
		}
	}

	if outcome == OutcomeDrop {
		return event, ErrDropOutcome
	}

	return event, nil
}

func NewRuleService(ruleAdapter RuleAdapter, container RuleApplicatorContainer, timezoneConfigProvider config.TimezoneConfigProvider, logger zerolog.Logger) Service {
	return &ruleService{
		rulesMutex:              sync.RWMutex{},
		rulesAdapter:            ruleAdapter,
		ruleApplicatorContainer: container,
		timezoneConfigProvider:  timezoneConfigProvider,
		logger:                  logger,
	}
}
