package eventfilter

import (
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/request"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	libtemplate "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type ParsedRule struct {
	ID           string
	Type         string
	Config       ParsedRuleConfig
	ExternalData map[string]ParsedExternalDataParameters
	Created      *datetime.CpsTime
	Updated      *datetime.CpsTime

	EventPattern pattern.Event
	savedpattern.EntityPatternFields

	ResolvedStart     *datetime.CpsTime
	ResolvedStop      *datetime.CpsTime
	NextResolvedStart *datetime.CpsTime
	NextResolvedStop  *datetime.CpsTime
	ResolvedExdates   []types.Exdate
}

type ParsedRuleConfig struct {
	Resource      libtemplate.ParsedTemplate
	Component     libtemplate.ParsedTemplate
	Connector     libtemplate.ParsedTemplate
	ConnectorName libtemplate.ParsedTemplate

	Actions   []ParsedAction
	OnSuccess string
	OnFailure string
}

type ParsedAction struct {
	Index       int
	Type        string
	Name        string
	Description string
	Value       any
	ParsedValue libtemplate.ParsedTemplate
}

type ParsedExternalDataParameters struct {
	Type string

	Collection string
	Select     map[string]libtemplate.ParsedTemplate
	Regexp     map[string]libtemplate.ParsedTemplate
	SortBy     string
	Sort       string

	RequestParameters *request.ParsedParameters
}

func ParseRule(rule Rule, tplExecutor libtemplate.Executor) ParsedRule {
	parsedActions := make([]ParsedAction, len(rule.Config.Actions))
	for i, action := range rule.Config.Actions {
		var parsedValue libtemplate.ParsedTemplate
		switch action.Type {
		case ActionSetFieldFromTemplate, ActionSetEntityInfoFromTemplate, ActionSetTagsFromTemplate:
			if str, ok := action.Value.(string); ok {
				parsedValue = tplExecutor.Parse(str)
			}
		}

		parsedActions[i] = ParsedAction{
			Index:       i,
			Type:        action.Type,
			Name:        action.Name,
			Description: action.Description,
			Value:       action.Value,
			ParsedValue: parsedValue,
		}
	}

	parsedExternalData := make(map[string]ParsedExternalDataParameters, len(rule.ExternalData))
	for name, params := range rule.ExternalData {
		parsedSelect := make(map[string]libtemplate.ParsedTemplate, len(params.Select))
		for k, v := range params.Select {
			parsedSelect[k] = tplExecutor.Parse(v)
		}

		parsedRegexp := make(map[string]libtemplate.ParsedTemplate, len(params.Regexp))
		for k, v := range params.Regexp {
			parsedRegexp[k] = tplExecutor.Parse(v)
		}

		var parsedRequestParameters *request.ParsedParameters
		if params.RequestParameters != nil {
			parsedRequestParameters = &request.ParsedParameters{
				URL:        tplExecutor.Parse(params.RequestParameters.URL),
				Method:     params.RequestParameters.Method,
				Auth:       params.RequestParameters.Auth,
				Headers:    params.RequestParameters.Headers,
				Payload:    tplExecutor.Parse(params.RequestParameters.Payload),
				SkipVerify: params.RequestParameters.SkipVerify,
				Timeout:    params.RequestParameters.Timeout,
				RetryCount: params.RequestParameters.RetryCount,
				RetryDelay: params.RequestParameters.RetryDelay,
			}
		}

		parsedExternalData[name] = ParsedExternalDataParameters{
			Type:              params.Type,
			Collection:        params.Collection,
			Select:            parsedSelect,
			Regexp:            parsedRegexp,
			SortBy:            params.SortBy,
			Sort:              params.Sort,
			RequestParameters: parsedRequestParameters,
		}
	}

	return ParsedRule{
		ID:   rule.ID,
		Type: rule.Type,
		Config: ParsedRuleConfig{
			Resource:      tplExecutor.Parse(rule.Config.Resource),
			Component:     tplExecutor.Parse(rule.Config.Component),
			Connector:     tplExecutor.Parse(rule.Config.Connector),
			ConnectorName: tplExecutor.Parse(rule.Config.ConnectorName),
			Actions:       parsedActions,
			OnSuccess:     rule.Config.OnSuccess,
			OnFailure:     rule.Config.OnFailure,
		},
		ExternalData:        parsedExternalData,
		Created:             rule.Created,
		Updated:             rule.Updated,
		EventPattern:        rule.EventPattern,
		EntityPatternFields: rule.EntityPatternFields,
		ResolvedStart:       rule.ResolvedStart,
		ResolvedStop:        rule.ResolvedStop,
		NextResolvedStart:   rule.NextResolvedStart,
		NextResolvedStop:    rule.NextResolvedStop,
		ResolvedExdates:     rule.ResolvedExdates,
	}
}

func ExecuteParsedTemplate(
	ruleID string,
	tplName string,
	parsedTpl libtemplate.ParsedTemplate,
	tplData any,
	event *types.Event,
	failureService FailureService,
	templateExecutor libtemplate.Executor,
) (string, error) {
	if parsedTpl.Err != nil {
		failReason := fmt.Sprintf("invalid template %q: %s", tplName, parsedTpl.Err)
		failureService.Add(ruleID, FailureTypeInvalidTemplate, failReason, nil)
		return "", parsedTpl.Err
	}

	if parsedTpl.Tpl != nil {
		res, err := templateExecutor.ExecuteByTpl(parsedTpl.Tpl, tplData)
		if err != nil {
			failReason := fmt.Sprintf("cannot execute template %q for event: %s", tplName, err)
			failureService.Add(ruleID, FailureTypeInvalidTemplate, failReason, event)
			return "", err
		}

		return res, nil
	}

	return "", nil
}
