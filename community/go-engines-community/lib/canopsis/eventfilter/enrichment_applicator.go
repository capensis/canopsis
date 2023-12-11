package eventfilter

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type enrichmentApplicator struct {
	externalDataContainer *ExternalDataContainer
	actionProcessor       ActionProcessor
	failureService        FailureService
}

func NewEnrichmentApplicator(
	externalDataContainer *ExternalDataContainer,
	processor ActionProcessor,
	failureService FailureService,
) RuleApplicator {
	return &enrichmentApplicator{
		externalDataContainer: externalDataContainer,
		actionProcessor:       processor,
		failureService:        failureService,
	}
}

func (a *enrichmentApplicator) Apply(
	ctx context.Context,
	rule ParsedRule,
	event *types.Event,
	regexMatch RegexMatch,
) (string, bool, error) {
	var entityUpdated bool

	externalData, err := a.getExternalData(ctx, rule, event, regexMatch)
	if err != nil {
		return rule.Config.OnFailure, false, err
	}

	for _, action := range rule.Config.Actions {
		isUpdated, err := a.actionProcessor.Process(ctx, rule.ID, action, event, regexMatch, externalData)
		if err != nil {
			return rule.Config.OnFailure, false, fmt.Errorf("invalid action name=%q type=%q: %w", action.Name, action.Type, err)
		}

		entityUpdated = entityUpdated || isUpdated
	}

	return rule.Config.OnSuccess, entityUpdated, nil
}

func (a *enrichmentApplicator) getExternalData(ctx context.Context, rule ParsedRule, event *types.Event, regexMatch RegexMatch) (map[string]interface{}, error) {
	externalData := make(map[string]interface{})

	for name, parameters := range rule.ExternalData {
		getter, ok := a.externalDataContainer.Get(parameters.Type)
		if !ok {
			failReason := fmt.Sprintf("external data %q has invalid type %q", name, parameters.Type)
			a.failureService.Add(rule.ID, FailureTypeOther, failReason, nil)
			return nil, fmt.Errorf("no such data source: %s", parameters.Type)
		}

		data, err := getter.Get(ctx, rule.ID, name, event, parameters, Template{
			Event:      event,
			RegexMatch: regexMatch,
		})
		if err != nil {
			return externalData, err
		}

		externalData[name] = data
	}

	return externalData, nil
}
