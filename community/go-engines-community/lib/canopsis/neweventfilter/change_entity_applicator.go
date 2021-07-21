package neweventfilter

import (
	"bytes"
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"text/template"
)

type changeEntityApplicator struct {
	dataSourceFactories map[string]eventfilter.DataSourceFactory
}

func NewChangeEntityApplicator(dataSourceFactories map[string]eventfilter.DataSourceFactory) RuleApplicator {
	return &changeEntityApplicator{dataSourceFactories: dataSourceFactories}
}

func (a *changeEntityApplicator) Apply(ctx context.Context, rule Rule, params ApplicatorParameters) (types.Event, error) {
	var err error
	event := params.Event

	params.ExternalData, err = a.getExternalData(ctx, rule, params)
	if err != nil {
		return params.Event, err
	}

	if rule.Config.Resource != "" {
		event.Resource, err = executeTpl(rule.Config.Resource, params)
		if err != nil {
			return params.Event, err
		}
	}

	if rule.Config.Component != "" {
		event.Component, err = executeTpl(rule.Config.Component, params)
		if err != nil {
			return params.Event, err
		}
	}

	if rule.Config.Connector != "" {
		event.Connector, err = executeTpl(rule.Config.Connector, params)
		if err != nil {
			return params.Event, err
		}
	}

	if rule.Config.ConnectorName != "" {
		event.ConnectorName, err = executeTpl(rule.Config.ConnectorName, params)
		if err != nil {
			return params.Event, err
		}
	}

	return event, nil
}

func (a *changeEntityApplicator) getExternalData(ctx context.Context, rule Rule, params ApplicatorParameters) (map[string]interface{}, error) {
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
			Event:      params.Event,
			RegexMatch: params.RegexMatch,
		}, &eventfilter.Report{})
		if err != nil {
			return externalData, err
		}

		externalData[name] = data
	}

	return externalData, nil
}

func executeTpl(tplText string, params ApplicatorParameters) (string, error) {
	tpl, err := template.New("tpl").Funcs(types.GetTemplateFunc()).Parse(tplText)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, params)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
