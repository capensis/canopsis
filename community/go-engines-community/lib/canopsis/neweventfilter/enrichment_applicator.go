package neweventfilter

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type enrichmentApplicator struct {
	dataSourceFactories map[string]eventfilter.DataSourceFactory
	actionProcessor ActionProcessor
}

func NewEnrichmentApplicator(dataSourceFactories map[string]eventfilter.DataSourceFactory, processor ActionProcessor) RuleApplicator {
	return &enrichmentApplicator{
		dataSourceFactories: dataSourceFactories,
		actionProcessor: processor,
	}
}

func (a *enrichmentApplicator) Apply(ctx context.Context, rule Rule, event types.Event, regexMatch pattern.EventRegexMatches) (string, types.Event, error) {
	externalData, err := a.getExternalData(ctx, rule, event, regexMatch)
	if err != nil {
		return rule.Config.OnFailure, event, err
	}

	for _, action := range rule.Config.Actions {
		event, err = a.actionProcessor.Process(action, event, regexMatch, externalData)
		if err != nil {
			return rule.Config.OnFailure, event, err
		}
	}

	return rule.Config.OnSuccess, event, nil
}

//TODO: copy from eventfilter package, all mongo plugin feature should be refactored
func (a *enrichmentApplicator) getExternalData(ctx context.Context, rule Rule, event types.Event, regexMatch pattern.EventRegexMatches) (map[string]interface{}, error) {
	externalData := make(map[string]interface{})

	for name, source := range rule.ExternalData {
		factory, success := a.dataSourceFactories[source.Type]
		if !success {
			return nil, fmt.Errorf("no such data source: %s", source.Type)
		}
		getter, err := factory.Create(source.DataSourceBase.Parameters)
		if err != nil {
			return nil, err
		}

		data, err := getter.Get(ctx, eventfilter.DataSourceGetterParameters{
			Event:      event,
			RegexMatch: regexMatch,
		}, &eventfilter.Report{})
		if err != nil {
			return externalData, err
		}

		externalData[name] = data
	}

	return externalData, nil
}
