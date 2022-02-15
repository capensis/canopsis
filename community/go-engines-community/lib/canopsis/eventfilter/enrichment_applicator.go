package eventfilter

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type enrichmentApplicator struct {
	externalDataContainer *ExternalDataContainer
	actionProcessor       ActionProcessor
}

func NewEnrichmentApplicator(externalDataContainer *ExternalDataContainer, processor ActionProcessor) RuleApplicator {
	return &enrichmentApplicator{
		externalDataContainer: externalDataContainer,
		actionProcessor:       processor,
	}
}

func (a *enrichmentApplicator) Apply(ctx context.Context, rule Rule, event types.Event, regexMatch pattern.EventRegexMatches, cfgTimezone *config.TimezoneConfig) (string, types.Event, error) {
	externalData, err := a.getExternalData(ctx, rule, event, regexMatch, cfgTimezone)
	if err != nil {
		return rule.Config.OnFailure, event, err
	}

	for _, action := range rule.Config.Actions {
		event, err = a.actionProcessor.Process(action, event, regexMatch, externalData, cfgTimezone)
		if err != nil {
			return rule.Config.OnFailure, event, err
		}
	}

	return rule.Config.OnSuccess, event, nil
}

func (a *enrichmentApplicator) getExternalData(ctx context.Context, rule Rule, event types.Event, regexMatch pattern.EventRegexMatches, cfgTimezone *config.TimezoneConfig) (map[string]interface{}, error) {
	externalData := make(map[string]interface{})

	for name, parameters := range rule.ExternalData {
		getter, ok := a.externalDataContainer.Get(parameters.Type)
		if !ok {
			return nil, fmt.Errorf("no such data source: %s", parameters.Type)
		}

		data, err := getter.Get(ctx, parameters, TemplateParameters{
			Event:      event,
			RegexMatch: regexMatch,
		}, cfgTimezone)
		if err != nil {
			return externalData, err
		}

		externalData[name] = data
	}

	return externalData, nil
}
