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
	// rulesAdapter is an rulesAdapter to the rules collection.
	rulesAdapter            RulesAdapter
	ruleApplicatorContainer RuleApplicatorContainer

	// rules is an array of rules.
	rules               []Rule
	dataSourceFactories map[string]eventfilter.DataSourceFactory

	logger zerolog.Logger
}

// copy from eventfilter package
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
	var err error
	newEvent := event

	s.logger.Debug().Msg("Process event by meta-alarm rule service")

	for _, rule := range s.rules {
		s.logger.Debug().Msgf("Meta-alarm rule service: check rule %s", rule.ID)

		regexMatches, match := rule.Patterns.GetRegexMatches(event)
		if !match {
			s.logger.Debug().Msg("Meta-alarm rule service: check rule is not matched")
			continue
		}

		s.logger.Debug().Msg("Meta-alarm rule service: check rule is matched")

		applicator, found := s.ruleApplicatorContainer.Get(rule.Type)
		if !found {
			s.logger.Warn().Msgf("RuleApplicator for %s is not exist", rule.Type)
			continue
		}

		newEvent, err = applicator.Apply(ctx, rule, ApplicatorParameters{
			Event:      event,
			RegexMatch: regexMatches,
		})
		if err != nil {
			return event, err
		}
	}

	return newEvent, nil
}

func NewRuleService(ruleAdapter RulesAdapter, container RuleApplicatorContainer, logger zerolog.Logger) EventFilterService {
	return &ruleService{
		rulesAdapter:            ruleAdapter,
		ruleApplicatorContainer: container,
		dataSourceFactories:     make(map[string]eventfilter.DataSourceFactory),
		logger:                  logger,
	}
}
