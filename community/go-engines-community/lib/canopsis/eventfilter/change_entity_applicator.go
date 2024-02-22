package eventfilter

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type changeEntityApplicator struct {
	externalDataContainer *ExternalDataContainer
	failureService        FailureService
	templateExecutor      template.Executor
}

func NewChangeEntityApplicator(
	externalDataContainer *ExternalDataContainer,
	failureService FailureService,
	templateExecutor template.Executor,
) RuleApplicator {
	return &changeEntityApplicator{
		externalDataContainer: externalDataContainer,
		failureService:        failureService,
		templateExecutor:      templateExecutor,
	}
}

func (a *changeEntityApplicator) Apply(ctx context.Context, rule ParsedRule, event *types.Event, regexMatch RegexMatch) (string, bool, error) {
	externalData, err := a.getExternalData(ctx, rule, event, regexMatch)
	if err != nil {
		return OutcomeDrop, false, err
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
			return OutcomeDrop, false, err
		}

		event.Resource = resource
	}

	if rule.Config.Component.Text != "" {
		component, err := ExecuteParsedTemplate(rule.ID, "Component", rule.Config.Component,
			templateParams, event, a.failureService, a.templateExecutor)
		if err != nil {
			return OutcomeDrop, false, err
		}

		event.Component = component
	}

	if rule.Config.Connector.Text != "" {
		connector, err := ExecuteParsedTemplate(rule.ID, "Connector", rule.Config.Connector,
			templateParams, event, a.failureService, a.templateExecutor)
		if err != nil {
			return OutcomeDrop, false, err
		}

		event.Connector = connector
	}

	if rule.Config.ConnectorName.Text != "" {
		connectorName, err := ExecuteParsedTemplate(rule.ID, "ConnectorName", rule.Config.ConnectorName,
			templateParams, event, a.failureService, a.templateExecutor)
		if err != nil {
			return OutcomeDrop, false, err
		}

		event.ConnectorName = connectorName
	}

	return OutcomePass, false, nil
}

func (a *changeEntityApplicator) getExternalData(ctx context.Context, rule ParsedRule, event *types.Event, regexMatch RegexMatch) (map[string]interface{}, error) {
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
