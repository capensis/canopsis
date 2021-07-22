package neweventfilter

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
	"io/ioutil"
	"path/filepath"
	"plugin"
	"strings"
)

type ruleService struct {
	rulesAdapter            RuleAdapter
	ruleApplicatorContainer RuleApplicatorContainer
	rules                   []Rule
	dataSourceFactories     map[string]eventfilter.DataSourceFactory
	logger                  zerolog.Logger
}

//TODO: copy from eventfilter package, all mongo plugin feature should be refactored
func (s *ruleService) LoadDataSourceFactories(dataSourceDirectory string) error {
	files, err := ioutil.ReadDir(dataSourceDirectory)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), canopsis.PluginExtension) {
			sourceName := strings.TrimSuffix(file.Name(), canopsis.PluginExtension)
			fileName := filepath.Join(dataSourceDirectory, file.Name())
			s.logger.Info().Msgf("loading data source plugin %s from file %s", sourceName, fileName)

			plug, err := plugin.Open(fileName)
			if err != nil {
				return fmt.Errorf("unable to open plugin: %v", err)
			}

			factorySymbol, err := plug.Lookup("DataSourceFactory")
			if err != nil {
				return fmt.Errorf("unable to load plugin: %v", err)
			}

			factory, isFactory := factorySymbol.(eventfilter.DataSourceFactory)
			if !isFactory {
				return fmt.Errorf("the plugin does not define a valid data source")
			}

			s.dataSourceFactories[sourceName] = factory
		}
	}

	return nil
}

func (s *ruleService) LoadRules(ctx context.Context) error {
	var err error

	s.rules, err = s.rulesAdapter.GetAll(ctx)
	if err != nil {
		return err
	}

	s.logger.Info().Int("number", len(s.rules)).Msg("Successfully loaded meta-alarm rules")

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
