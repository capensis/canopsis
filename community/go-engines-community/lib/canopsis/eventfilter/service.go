package eventfilter

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"plugin"
	"strings"
	"sync"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	libcontext "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
)

// Report is a type containing informations about the modifications that were
// made to an event.
type Report struct {
	// EntityUpdate is a boolean indicating whether the event's entity was
	// updated.
	EntityUpdated bool
	// UpdatedEntityServices contains context graph update result if entity has been updated.
	UpdatedEntityServices libcontext.UpdatedEntityServices
}

// service is the service that manages the event filter.
type service struct {
	// adapter is an adapter to the rules collection.
	adapter Adapter

	// rules is an array of rules.
	rules []Rule

	// dataSourceFactories maps the name of a data source to the corresponding
	// DataSourceFactory
	dataSourceFactories map[string]DataSourceFactory

	// rulesMutex is a mutex preventing to overwrite the rules array while
	// it's being used.
	rulesMutex sync.RWMutex

	timezoneConfigProvider config.TimezoneConfigProvider
	logger                 zerolog.Logger
}

// NewService creates an event filter service.
func NewService(adapter Adapter, timezoneConfigProvider config.TimezoneConfigProvider, logger zerolog.Logger) Service {
	s := service{
		adapter:                adapter,
		timezoneConfigProvider: timezoneConfigProvider,
		logger:                 logger,
	}
	return &s
}

// LoadDataSourceFactories loads the data source factories and adds them to the
// service.
func (s *service) LoadDataSourceFactories(enrichmentCenter libcontext.EnrichmentCenter, enrichFields libcontext.EnrichFields, dataSourceDirectory string) error {
	s.dataSourceFactories = make(map[string]DataSourceFactory)

	s.dataSourceFactories["entity"] = NewEntityDataSourceFactory(enrichmentCenter, enrichFields)

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

			factory, isFactory := factorySymbol.(DataSourceFactory)
			if !isFactory {
				return fmt.Errorf("the plugin does not define a valid data source")
			}

			s.dataSourceFactories[sourceName] = factory
		}
	}

	return nil
}

// loadRuleDataSources loads the data source of a rule.
// It uses the type of the data source to get a DataSourceFactory, and uses the
// factory to instanciate a DataSourceGetter.
func (s *service) loadRuleDataSources(rule *Rule) error {
	if rule.Type != RuleTypeEnrichment {
		return nil
	}

	for _, source := range rule.ExternalData {
		factory, success := s.dataSourceFactories[source.Type]
		if !success {
			return fmt.Errorf("no such data source: %s", source.Type)
		}
		getter, err := factory.Create(source.DataSourceBase.Parameters)
		if err != nil {
			return err
		}
		source.DataSourceGetter = getter
	}

	return nil
}

// LoadRules loads the event filter rules from the database, and adds them to
// the service. Note that LoadDataSourceFactories needs to be called before
// calling LoadRules.
func (s *service) LoadRules(ctx context.Context) error {
	allRules, err := s.adapter.List(ctx)
	if err != nil {
		return err
	}

	rules := make([]Rule, 0)
	for _, rule := range allRules {
		err := s.loadRuleDataSources(&rule)
		if err == nil {
			rules = append(rules, rule)
		} else {
			s.logger.Warn().Str("rule", rule.ID).Err(err).Msg("unable to load data source for rule")
		}
	}

	s.rulesMutex.Lock()
	s.rules = rules
	s.rulesMutex.Unlock()
	return nil
}

// ProcessEvent processes an event with the rules of the event filter. It
// returns a DropError if the event should be dropped by the eventfilter.
// Note that LoadRules needs to be called once before ProcessEvent can be used.
func (s *service) ProcessEvent(ctx context.Context, event types.Event) (types.Event, Report, error) {
	s.rulesMutex.RLock()
	defer s.rulesMutex.RUnlock()

	if event.Debug {
		s.logger.Info().Str("event", fmt.Sprintf("%+v", event)).Msg("eventfilter | entering event filter")
	}

	report := Report{}
	outcome := UnsetOutcome
	tz := s.timezoneConfigProvider.Get()
	for _, rule := range s.rules {
		if outcome != UnsetOutcome && outcome != Pass {
			break
		}

		if event.Debug {
			s.logger.Info().Msgf("eventfilter | >>> rule %s", rule.ID)
		}

		regexMatches, match := rule.Patterns.GetRegexMatches(event)
		if match {
			if event.Debug {
				s.logger.Info().Str("regex", fmt.Sprintf("%+v", regexMatches)).Msg("eventfilter | event matches, applying rule with regex matches")
			}

			event, outcome = rule.Apply(ctx, event, regexMatches, &report, &tz, s.logger)

			if event.Debug {
				s.logger.Info().Msgf("eventfilter | outcome: %s", outcome)
				s.logger.Info().Msgf("eventfilter | event: %+v", event)
			}
		} else if event.Debug {
			s.logger.Info().Msg("eventfilter | event does not match")
		}
	}

	if event.Debug {
		s.logger.Info().Msg("eventfilter | leaving event filter")
	}

	if outcome == Drop {
		return event, report, DefaultDropError()
	}

	return event, report, event.IsValid()
}
