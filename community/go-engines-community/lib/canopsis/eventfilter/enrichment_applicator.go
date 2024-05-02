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
) (string, bool, map[string]int64, error) {
	var entityUpdated bool
	externalData, externalRequestCount, err := getExternalData(ctx, rule, event, regexMatch, a.externalDataContainer, a.failureService)
	if err != nil {
		return rule.Config.OnFailure, false, nil, err
	}

	for _, action := range rule.Config.Actions {
		isUpdated, err := a.actionProcessor.Process(ctx, rule.ID, action, event, regexMatch, externalData)
		if err != nil {
			return rule.Config.OnFailure, false, nil, fmt.Errorf("invalid action name=%q type=%q: %w", action.Name, action.Type, err)
		}

		entityUpdated = entityUpdated || isUpdated
	}

	return rule.Config.OnSuccess, entityUpdated, externalRequestCount, nil
}

func getExternalData(
	ctx context.Context,
	rule ParsedRule,
	event *types.Event,
	regexMatch RegexMatch,
	externalDataContainer *ExternalDataContainer,
	failureService FailureService,
) (map[string]any, map[string]int64, error) {
	externalData := make(map[string]any)
	externalRequestCount := make(map[string]int64)

	for name, parameters := range rule.ExternalData {
		getter, ok := externalDataContainer.Get(parameters.Type)
		if !ok {
			failReason := fmt.Sprintf("external data %q has invalid type %q", name, parameters.Type)
			failureService.Add(rule.ID, FailureTypeOther, failReason, nil)

			return nil, nil, fmt.Errorf("no such data source: %s", parameters.Type)
		}

		data, err := getter.Get(ctx, rule.ID, name, event, parameters, Template{
			Event:      event,
			RegexMatch: regexMatch,
		})
		if err != nil {
			return nil, nil, err
		}

		externalData[name] = data
		externalRequestCount[parameters.Type]++
	}

	return externalData, externalRequestCount, nil
}
