package eventfilter

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

var failureTypeMapping = map[string]int64{
	externaldata.RefTypeAPI:   FailureTypeExternalDataAPI,
	externaldata.RefTypeTable: FailureTypeExternalDataTable,
}

type enrichmentApplicator struct {
	actionProcessor ActionProcessor
	failureService  FailureService
}

func NewEnrichmentApplicator(
	processor ActionProcessor,
	failureService FailureService,
) RuleApplicator {
	return &enrichmentApplicator{
		actionProcessor: processor,
		failureService:  failureService,
	}
}

func (a *enrichmentApplicator) Apply(
	ctx context.Context,
	rule ParsedRule,
	event *types.Event,
	updatedEntityInfos map[string]UpdatedValue,
	regexMatch RegexMatch,
	externalData map[string]any,
) (RuleResult, error) {
	var err error
	for _, action := range rule.Config.Actions {
		updatedEntityInfos, err = a.actionProcessor.Process(ctx, rule, action, event, updatedEntityInfos, regexMatch, externalData)
		if err != nil {
			return RuleResult{Outcome: rule.Config.OnFailure}, fmt.Errorf("invalid action name=%q type=%q: %w", action.Name, action.Type, err)
		}
	}

	return RuleResult{
		Outcome:            rule.Config.OnSuccess,
		UpdatedEntityInfos: updatedEntityInfos,
	}, nil
}
