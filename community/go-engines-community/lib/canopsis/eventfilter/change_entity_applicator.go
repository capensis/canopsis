package eventfilter

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type changeEntityApplicator struct {
	externalDataContainer *externaldata.GetterContainer
	failureService        FailureService
	templateExecutor      template.Executor
}

func NewChangeEntityApplicator(
	externalDataContainer *externaldata.GetterContainer,
	failureService FailureService,
	templateExecutor template.Executor,
) RuleApplicator {
	return &changeEntityApplicator{
		externalDataContainer: externalDataContainer,
		failureService:        failureService,
		templateExecutor:      templateExecutor,
	}
}

func (a *changeEntityApplicator) Apply(ctx context.Context, rule ParsedRule, event *types.Event, regexMatch RegexMatch) (string, bool, map[string]int64, error) {
	externalData, externalRequestCount, err := getExternalData(ctx, rule, event, regexMatch, a.externalDataContainer, a.failureService)
	if err != nil {
		return OutcomeDrop, false, nil, err
	}

	templateParams := Template{
		Event:        event,
		RegexMatch:   regexMatch,
		ExternalData: externalData,
	}

	if rule.Config.Resource.Text != "" {
		resource, err := ExecuteParsedTemplate(rule.ID, "Resource", rule.Config.Resource,
			templateParams, event, a.failureService, a.templateExecutor)
		if err != nil {
			return OutcomeDrop, false, nil, err
		}

		event.Resource = resource
	}

	if rule.Config.Component.Text != "" {
		component, err := ExecuteParsedTemplate(rule.ID, "Component", rule.Config.Component,
			templateParams, event, a.failureService, a.templateExecutor)
		if err != nil {
			return OutcomeDrop, false, nil, err
		}

		event.Component = component
	}

	if rule.Config.Connector.Text != "" {
		connector, err := ExecuteParsedTemplate(rule.ID, "Connector", rule.Config.Connector,
			templateParams, event, a.failureService, a.templateExecutor)
		if err != nil {
			return OutcomeDrop, false, nil, err
		}

		event.Connector = connector
	}

	if rule.Config.ConnectorName.Text != "" {
		connectorName, err := ExecuteParsedTemplate(rule.ID, "ConnectorName", rule.Config.ConnectorName,
			templateParams, event, a.failureService, a.templateExecutor)
		if err != nil {
			return OutcomeDrop, false, nil, err
		}

		event.ConnectorName = connectorName
	}

	if rule.Config.Upstream.Text != "" {
		upstream, err := ExecuteParsedTemplate(rule.ID, "Upstream", rule.Config.Upstream,
			templateParams, event, a.failureService, a.templateExecutor)
		if err != nil {
			return OutcomeDrop, false, nil, err
		}

		event.Upstream = upstream
	}

	return OutcomePass, false, externalRequestCount, nil
}
