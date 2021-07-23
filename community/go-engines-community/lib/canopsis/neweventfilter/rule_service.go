package neweventfilter

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
)

type ruleService struct {
	rulesAdapter            RuleAdapter
	ruleApplicatorContainer RuleApplicatorContainer
	rules                   []Rule
	dataSourceFactories     map[string]eventfilter.DataSourceFactory
	logger                  zerolog.Logger
}

func (s *ruleService) LoadRules(ctx context.Context) error {
	var err error

	s.rules, err = s.rulesAdapter.GetAll(ctx)
	if err != nil {
		return err
	}

	s.logger.Info().Int("number", len(s.rules)).Msg("Successfully loaded eventfilter rules")

	return nil
}

func (s *ruleService) ProcessEvent(ctx context.Context, event types.Event) (types.Event, error) {
	outcome := OutcomePass

	var err error
	for _, rule := range s.rules {
		if outcome != OutcomePass {
			break
		}

		s.logger.Debug().Msgf("Event filter rule service: check rule %s", rule.ID)

		regexMatches, match := rule.Patterns.GetRegexMatches(event)
		if !match {
			s.logger.Debug().Str("rule_id", rule.ID).Msg("Event filter rule service: rule is not matched")
			continue
		}

		s.logger.Debug().Str("rule_id", rule.ID).Msg("Event filter rule service: rule is matched")

		applicator, found := s.ruleApplicatorContainer.Get(rule.Type)
		if !found {
			s.logger.Warn().Str("rule_id", rule.ID).Str("rule_type", rule.Type).Msg("Event filter rule service: RuleApplicator doesn't exist")
			continue
		}

		outcome, event, err = applicator.Apply(ctx, rule, event, regexMatches)
		if err != nil {
			return event, err
		}
	}

	if outcome == OutcomeDrop {
		return event, ErrDropOutcome
	}

	return event, nil
}

func NewRuleService(ruleAdapter RuleAdapter, container RuleApplicatorContainer, logger zerolog.Logger) EventFilterService {
	return &ruleService{
		rulesAdapter:            ruleAdapter,
		ruleApplicatorContainer: container,
		dataSourceFactories:     make(map[string]eventfilter.DataSourceFactory),
		logger:                  logger,
	}
}
