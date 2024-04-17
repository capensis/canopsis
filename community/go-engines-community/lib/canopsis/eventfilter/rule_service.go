package eventfilter

import (
	"context"
	"sync"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
)

type ruleService struct {
	rulesAdapter            RuleAdapter
	ruleApplicatorContainer RuleApplicatorContainer
	eventCounter            EventCounter
	failureService          FailureService
	templateExecutor        template.Executor
	rules                   []ParsedRule
	rulesMutex              sync.RWMutex
	logger                  zerolog.Logger
}

func (s *ruleService) LoadRules(ctx context.Context, types []string) error {
	rules, err := s.rulesAdapter.GetByTypes(ctx, types)
	if err != nil {
		return err
	}

	parsedRules := make([]ParsedRule, len(rules))
	ids := make([]string, len(rules))
	for i := 0; i < len(rules); i++ {
		ids[i] = rules[i].ID
		parsedRules[i] = ParseRule(rules[i], s.templateExecutor)
	}

	s.rulesMutex.Lock()
	defer s.rulesMutex.Unlock()
	s.rules = parsedRules
	s.logger.Debug().Strs("rules", ids).Msg("Loading event filter rules")

	return nil
}

func (s *ruleService) ProcessEvent(ctx context.Context, event *types.Event) (bool, int64, map[string]int64, error) {
	s.rulesMutex.RLock()
	defer s.rulesMutex.RUnlock()

	outcome := OutcomePass
	now := datetime.NewCpsTime()
	entityUpdated := false
	var executedEnrichRuleCount int64
	externalRequestCount := make(map[string]int64)
	for _, rule := range s.rules {
		if outcome != OutcomePass {
			break
		}

		if rule.ResolvedStart != nil && rule.ResolvedStop != nil {
			inExDate := false
			for _, exdate := range rule.ResolvedExdates {
				if now.After(exdate.Begin) && now.Before(exdate.End) {
					inExDate = true
					break
				}
			}

			if inExDate {
				continue
			}

			if now.Before(*rule.ResolvedStart) {
				continue
			}

			if now.After(*rule.ResolvedStop) {
				if rule.NextResolvedStart == nil || rule.NextResolvedStop == nil {
					continue
				}

				if now.Before(*rule.NextResolvedStart) || now.After(*rule.NextResolvedStop) {
					continue
				}
			}
		}

		var err error
		var matched bool
		var eventRegexMatches match.EventRegexMatches
		var entityRegexMatches match.EntityRegexMatches

		if len(rule.EntityPattern) == 0 && len(rule.EventPattern) == 0 {
			if event.Debug {
				s.logger.Info().Str("rule", rule.ID).Str("event_type", event.EventType).Str("entity", event.GetEID()).Msg("Event filter rule service: rule is not matched")
			}

			continue
		}

		if len(rule.EventPattern) > 0 || len(rule.EntityPattern) > 0 {
			matched, eventRegexMatches, err = match.MatchEventPatternWithRegexMatches(rule.EventPattern, event)
			if err != nil {
				s.logger.Err(err).Str("rule_id", rule.ID).Msg("Event filter rule service: invalid event pattern")
				s.failureService.Add(rule.ID, FailureTypeInvalidPattern, "invalid event pattern: "+err.Error(), nil)
				continue
			}

			if !matched {
				if event.Debug {
					s.logger.Info().Str("rule", rule.ID).Str("event_type", event.EventType).Str("entity", event.GetEID()).Msg("Event filter rule service: rule is not matched")
				}

				continue
			}

			if event.Entity != nil {
				matched, entityRegexMatches, err = match.MatchEntityPatternWithRegexMatches(rule.EntityPattern, event.Entity)
				if err != nil {
					s.logger.Err(err).Str("rule_id", rule.ID).Msg("Event filter rule service: invalid entity pattern")
					s.failureService.Add(rule.ID, FailureTypeInvalidPattern, "invalid entity pattern: "+err.Error(), nil)
					continue
				}

				if !matched {
					if event.Debug {
						s.logger.Info().Str("rule", rule.ID).Str("event_type", event.EventType).Str("entity", event.GetEID()).Msg("Event filter rule service: rule is not matched")
					}

					continue
				}
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

		var isUpdated bool
		var ruleExternalRequestCount map[string]int64
		outcome, isUpdated, ruleExternalRequestCount, err = applicator.Apply(ctx, rule, event, RegexMatch{
			EventRegexMatches: eventRegexMatches,
			Entity:            entityRegexMatches,
		})
		if err != nil {
			s.logger.Err(err).Str("rule_id", rule.ID).Str("rule_type", rule.Type).Msg("Event filter rule service: failed to apply")
			continue
		}

		entityUpdated = entityUpdated || isUpdated
		if rule.Updated != nil {
			s.eventCounter.Add(rule.ID, *rule.Updated)
		}

		if rule.Type == RuleTypeEnrichment {
			executedEnrichRuleCount++
		}

		for k, v := range ruleExternalRequestCount {
			externalRequestCount[k] += v
		}
	}

	if outcome == OutcomeDrop {
		return false, executedEnrichRuleCount, externalRequestCount, ErrDropOutcome
	}

	return entityUpdated, executedEnrichRuleCount, externalRequestCount, nil
}

func NewRuleService(
	ruleAdapter RuleAdapter,
	container RuleApplicatorContainer,
	eventCounter EventCounter,
	failureService FailureService,
	templateExecutor template.Executor,
	logger zerolog.Logger,
) Service {
	return &ruleService{
		rulesMutex:              sync.RWMutex{},
		rulesAdapter:            ruleAdapter,
		ruleApplicatorContainer: container,
		eventCounter:            eventCounter,
		failureService:          failureService,
		templateExecutor:        templateExecutor,
		logger:                  logger,
	}
}
