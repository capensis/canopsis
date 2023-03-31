package eventfilter

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type changeEntityApplicator struct {
	externalDataContainer *ExternalDataContainer
	templateExecutor      template.Executor
}

func NewChangeEntityApplicator(
	externalDataContainer *ExternalDataContainer,
	templateExecutor template.Executor,
) RuleApplicator {
	return &changeEntityApplicator{
		externalDataContainer: externalDataContainer,
		templateExecutor:      templateExecutor,
	}
}

func (a *changeEntityApplicator) Apply(ctx context.Context, rule Rule, event types.Event, regexMatchWrapper RegexMatchWrapper) (string, types.Event, error) {
	externalData, err := a.getExternalData(ctx, rule, event, regexMatchWrapper)
	if err != nil {
		return OutcomeDrop, event, err
	}

	templateParams := Template{
		Event:             event,
		RegexMatchWrapper: regexMatchWrapper,
		ExternalData:      externalData,
	}

	if rule.Config.Resource != "" {
		event.Resource, err = a.templateExecutor.Execute(rule.Config.Resource, templateParams.GetTemplate())
		if err != nil {
			return OutcomeDrop, event, err
		}
	}

	if rule.Config.Component != "" {
		event.Component, err = a.templateExecutor.Execute(rule.Config.Component, templateParams.GetTemplate())
		if err != nil {
			return OutcomeDrop, event, err
		}
	}

	if rule.Config.Connector != "" {
		event.Connector, err = a.templateExecutor.Execute(rule.Config.Connector, templateParams.GetTemplate())
		if err != nil {
			return OutcomeDrop, event, err
		}
	}

	if rule.Config.ConnectorName != "" {
		event.ConnectorName, err = a.templateExecutor.Execute(rule.Config.ConnectorName, templateParams.GetTemplate())
		if err != nil {
			return OutcomeDrop, event, err
		}
	}

	return OutcomePass, event, nil
}

func (a *changeEntityApplicator) getExternalData(ctx context.Context, rule Rule, event types.Event, regexMatchWrapper RegexMatchWrapper) (map[string]interface{}, error) {
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
