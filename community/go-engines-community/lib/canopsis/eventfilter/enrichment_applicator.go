package eventfilter

import (
	"context"
	"errors"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

var failureTypeMapping = map[string]int64{
	externaldata.RefTypeAPI:   FailureTypeExternalDataAPI,
	externaldata.RefTypeTable: FailureTypeExternalDataTable,
}

type enrichmentApplicator struct {
	externalDataContainer *externaldata.GetterContainer
	actionProcessor       ActionProcessor
	failureService        FailureService
}

func NewEnrichmentApplicator(
	externalDataContainer *externaldata.GetterContainer,
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
	updatedEntityInfos map[string]UpdatedValue,
	regexMatch RegexMatch,
) (RuleResult, error) {
	externalData, externalRequestCount, err := getExternalData(ctx, rule, event, regexMatch, a.externalDataContainer, a.failureService)
	if err != nil {
		return RuleResult{Outcome: rule.Config.OnFailure}, err
	}

	for _, action := range rule.Config.Actions {
		updatedEntityInfos, err = a.actionProcessor.Process(ctx, rule, action, event, updatedEntityInfos, regexMatch, externalData)
		if err != nil {
			return RuleResult{Outcome: rule.Config.OnFailure}, fmt.Errorf("invalid action name=%q type=%q: %w", action.Name, action.Type, err)
		}
	}

	return RuleResult{
		Outcome:              rule.Config.OnSuccess,
		UpdatedEntityInfos:   updatedEntityInfos,
		ExternalRequestCount: externalRequestCount,
	}, nil
}

func getExternalData(
	ctx context.Context,
	rule ParsedRule,
	event *types.Event,
	regexMatch RegexMatch,
	externalDataContainer *externaldata.GetterContainer,
	failureService FailureService,
) (map[string]any, map[string]int64, error) {
	externalData := make(map[string]any)
	externalRequestCount := make(map[string]int64)

	for _, parameters := range rule.ExternalData {
		getter, ok := externalDataContainer.Get(parameters.Type)
		if !ok {
			failReason := fmt.Sprintf("external data %q has invalid type %q", parameters.Reference, parameters.Type)
			failureService.Add(rule.ID, rule.Description, rule.Updated, FailureTypeOther, failReason, nil)

			return nil, nil, fmt.Errorf("no such data source: %s", parameters.Type)
		}

		data, err := getter.Get(ctx, parameters, Template{
			Event:      event,
			RegexMatch: regexMatch,
		})
		if err != nil {
			getterTplErr := &externaldata.GetterTplError{}
			getterErr := &externaldata.GetterError{}
			var failureType int64
			failReason := ""
			isParamsInvalid := false
			if errors.As(err, &getterTplErr) {
				failureType = FailureTypeInvalidTemplate
				failReason = getterTplErr.FailReason()
				isParamsInvalid = getterTplErr.IsParamsInvalid()
			} else if errors.As(err, &getterErr) {
				failureType, ok = failureTypeMapping[parameters.Type]
				if !ok {
					failureType = FailureTypeOther
				}

				failReason = getterErr.FailReason()
				isParamsInvalid = getterErr.IsParamsInvalid()
			}

			if failReason != "" {
				if isParamsInvalid {
					failureService.Add(rule.ID, rule.Description, rule.Updated, failureType, failReason, nil)
				} else {
					failureService.Add(rule.ID, rule.Description, rule.Updated, failureType, failReason, event)
				}
			}

			return nil, nil, err
		}

		externalData[parameters.Reference] = data
		externalRequestCount[parameters.Type]++
	}

	return externalData, externalRequestCount, nil
}
