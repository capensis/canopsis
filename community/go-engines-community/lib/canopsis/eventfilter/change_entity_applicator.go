package eventfilter

import (
	"bytes"
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"text/template"
)

type changeEntityApplicator struct {
	externalDataContainer *ExternalDataContainer
	buf                   bytes.Buffer
}

func NewChangeEntityApplicator(externalDataContainer *ExternalDataContainer) RuleApplicator {
	return &changeEntityApplicator{externalDataContainer: externalDataContainer, buf: bytes.Buffer{}}
}

func (a *changeEntityApplicator) Apply(ctx context.Context, rule Rule, event types.Event, regexMatchWrapper RegexMatchWrapper, cfgTimezone *config.TimezoneConfig) (string, types.Event, error) {
	externalData, err := a.getExternalData(ctx, rule, event, regexMatchWrapper, cfgTimezone)
	if err != nil {
		return OutcomeDrop, event, err
	}

	templateParams := Template{
		Event:             event,
		RegexMatchWrapper: regexMatchWrapper,
		ExternalData:      externalData,
	}

	if rule.Config.Resource != "" {
		event.Resource, err = a.executeTpl(rule.Config.Resource, templateParams, cfgTimezone)
		if err != nil {
			return OutcomeDrop, event, err
		}
	}

	if rule.Config.Component != "" {
		event.Component, err = a.executeTpl(rule.Config.Component, templateParams, cfgTimezone)
		if err != nil {
			return OutcomeDrop, event, err
		}
	}

	if rule.Config.Connector != "" {
		event.Connector, err = a.executeTpl(rule.Config.Connector, templateParams, cfgTimezone)
		if err != nil {
			return OutcomeDrop, event, err
		}
	}

	if rule.Config.ConnectorName != "" {
		event.ConnectorName, err = a.executeTpl(rule.Config.ConnectorName, templateParams, cfgTimezone)
		if err != nil {
			return OutcomeDrop, event, err
		}
	}

	return OutcomePass, event, nil
}

func (a *changeEntityApplicator) getExternalData(ctx context.Context, rule Rule, event types.Event, regexMatchWrapper RegexMatchWrapper, cfgTimezone *config.TimezoneConfig) (map[string]interface{}, error) {
	externalData := make(map[string]interface{})

	for name, parameters := range rule.ExternalData {
		getter, ok := a.externalDataContainer.Get(parameters.Type)
		if !ok {
			return nil, fmt.Errorf("no such data source: %s", parameters.Type)
		}

		data, err := getter.Get(ctx, parameters, Template{
			Event:             event,
			RegexMatchWrapper: regexMatchWrapper,
		}, cfgTimezone)
		if err != nil {
			return externalData, err
		}

		externalData[name] = data
	}

	return externalData, nil
}

func (a *changeEntityApplicator) executeTpl(tplText string, params TemplateGetter, cfgTimezone *config.TimezoneConfig) (string, error) {
	tpl, err := template.New("tpl").Option("missingkey=error").Funcs(types.GetTemplateFunc(cfgTimezone)).Parse(tplText)
	if err != nil {
		return "", err
	}

	a.buf.Reset()

	err = tpl.Execute(&a.buf, params.GetTemplate())
	if err != nil {
		return "", err
	}

	return a.buf.String(), nil
}
