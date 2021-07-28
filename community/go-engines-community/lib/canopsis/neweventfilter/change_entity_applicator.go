package neweventfilter

import (
	"bytes"
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"text/template"
)

type changeEntityApplicator struct {
	dataSourceFactories map[string]eventfilter.DataSourceFactory
	buf bytes.Buffer
}

func NewChangeEntityApplicator(dataSourceFactories map[string]eventfilter.DataSourceFactory) RuleApplicator {
	return &changeEntityApplicator{dataSourceFactories: dataSourceFactories, buf: bytes.Buffer{}}
}

func (a *changeEntityApplicator) Apply(ctx context.Context, rule Rule, event types.Event, regexMatch pattern.EventRegexMatches) (int, types.Event, error) {
	externalData, err := a.getExternalData(ctx, rule, event, regexMatch)
	if err != nil {
		return OutcomeDrop, event, err
	}

	templateParams := TemplateParameters{
		Event:        event,
		RegexMatch:   regexMatch,
		ExternalData: externalData,
	}

	if rule.Config.Resource != "" {
		event.Resource, err = a.executeTpl(rule.Config.Resource, templateParams)
		if err != nil {
			return OutcomeDrop, event, err
		}
	}

	if rule.Config.Component != "" {
		event.Component, err = a.executeTpl(rule.Config.Component, templateParams)
		if err != nil {
			return OutcomeDrop, event, err
		}
	}

	if rule.Config.Connector != "" {
		event.Connector, err = a.executeTpl(rule.Config.Connector, templateParams)
		if err != nil {
			return OutcomeDrop, event, err
		}
	}

	if rule.Config.ConnectorName != "" {
		event.ConnectorName, err = a.executeTpl(rule.Config.ConnectorName, templateParams)
		if err != nil {
			return OutcomeDrop, event, err
		}
	}

	return OutcomePass, event, nil
}

//TODO: copy from eventfilter package, all mongo plugin feature should be refactored
func (a *changeEntityApplicator) getExternalData(ctx context.Context, rule Rule, event types.Event, regexMatch pattern.EventRegexMatches) (map[string]interface{}, error) {
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

func (a *changeEntityApplicator) executeTpl(tplText string, params TemplateParameters) (string, error) {
	tpl, err := template.New("tpl").Funcs(types.GetTemplateFunc()).Parse(tplText)
	if err != nil {
		return "", err
	}

	a.buf.Reset()

	err = tpl.Execute(&a.buf, params)
	if err != nil {
		return "", err
	}

	return a.buf.String(), nil
}
