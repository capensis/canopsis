package eventfilter

import (
	"context"
	"fmt"

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

func (a *enrichmentApplicator) Apply(
	ctx context.Context,
	rule Rule,
	event types.Event,
	regexMatchWrapper RegexMatchWrapper,
) (string, types.Event, error) {
	externalData, err := a.getExternalData(ctx, rule, event, regexMatchWrapper)
	if err != nil {
		return rule.Config.OnFailure, event, err
	}

	for _, action := range rule.Config.Actions {
		event, err = a.actionProcessor.Process(ctx, action, event, regexMatchWrapper, externalData)
		if err != nil {
			return rule.Config.OnFailure, event, fmt.Errorf("invalid action name=%q type=%q: %w", action.Name, action.Type, err)
		}
	}

	return rule.Config.OnSuccess, event, nil
}

func (a *enrichmentApplicator) getExternalData(ctx context.Context, rule Rule, event types.Event, regexMatchWrapper RegexMatchWrapper) (map[string]interface{}, error) {
	externalData := make(map[string]interface{})

	for name, parameters := range rule.ExternalData {
		getter, ok := a.externalDataContainer.Get(parameters.Type)
		if !ok {
			return nil, fmt.Errorf("no such data source: %s", parameters.Type)
		}

		data, err := getter.Get(ctx, parameters, Template{
			Event:             event,
			RegexMatchWrapper: regexMatchWrapper,
		})
		if err != nil {
			return externalData, err
		}

		externalData[name] = data
	}

	return externalData, nil
}
