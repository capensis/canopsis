package eventfilter

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/oldpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
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

	for _, rule := range s.rules {
		if outcome != OutcomePass {
			break
		}

		var err error
		var match, backwardCompatibility bool
		var oldRegexMatches oldpattern.EventRegexMatches
		var eventRegexMatches pattern.EventRegexMatches
		var entityRegexMatches pattern.EntityRegexMatches

		if !rule.OldPatterns.IsSet() && len(rule.EntityPattern) == 0 && len(rule.EventPattern) == 0 {
			if event.Debug {
				s.logger.Info().Str("rule", rule.ID).Str("event_type", event.EventType).Str("entity", event.GetEID()).Msg("Event filter rule service: rule is not matched")
			}

			continue
		}

		if len(rule.EventPattern) > 0 || len(rule.EntityPattern) > 0 {
			match, eventRegexMatches, err = rule.EventPattern.Match(event)
			if err != nil {
				s.logger.Err(err).Str("rule_id", rule.ID).Msg("Event filter rule service: invalid event pattern")
				continue
			}

			if !match {
				if event.Debug {
					s.logger.Info().Str("rule", rule.ID).Str("event_type", event.EventType).Str("entity", event.GetEID()).Msg("Event filter rule service: rule is not matched")
				}

				continue
			}

			if event.Entity != nil {
				match, entityRegexMatches, err = rule.EntityPattern.Match(*event.Entity)
				if err != nil {
					s.logger.Err(err).Str("rule_id", rule.ID).Msg("Event filter rule service: invalid entity pattern")
					continue
				}

				if !match {
					if event.Debug {
						s.logger.Info().Str("rule", rule.ID).Str("event_type", event.EventType).Str("entity", event.GetEID()).Msg("Event filter rule service: rule is not matched")
					}

					continue
				}
			}
		} else if rule.OldPatterns.IsSet() {
			backwardCompatibility = true
			if !rule.OldPatterns.IsValid() {
				s.logger.Warn().Msgf("Rule %s has an invalid old event pattern, skip", rule.ID)
				continue
			}

			oldRegexMatches, match = rule.OldPatterns.GetRegexMatches(event)
			if !match {
				if event.Debug {
					s.logger.Info().Str("rule", rule.ID).Str("event_type", event.EventType).Str("entity", event.GetEID()).Msg("Event filter rule service: rule is not matched")
				}

				continue
			}
		}

		if event.Debug {
			s.logger.Info().Str("rule", rule.ID).Str("event_type", event.EventType).Str("entity", event.GetEID()).Msg("Event filter rule service: rule is matched")
		}

		applicator, found := s.ruleApplicatorContainer.Get(rule.Type)
		if !found {
			s.logger.Warn().Str("rule_id", rule.ID).Str("rule_type", rule.Type).Msg("Event filter rule service: RuleApplicator doesn't exist")
			continue
		}

		outcome, event, err = applicator.Apply(ctx, rule, event, RegexMatchWrapper{
			BackwardCompatibility: backwardCompatibility,
			OldRegexMatch:         oldRegexMatches,
			RegexMatch: RegexMatch{
				EventRegexMatches: eventRegexMatches,
				Entity:            entityRegexMatches,
			},
		}, &tz)

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
