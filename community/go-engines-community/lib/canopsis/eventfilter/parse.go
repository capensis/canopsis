package eventfilter

import (
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	libtemplate "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type ParsedRule struct {
	ID           string
	Description  string
	Type         string
	Config       ParsedRuleConfig
	ExternalData []externaldata.ParsedRefParameters
	Created      datetime.CpsTime
	Updated      datetime.CpsTime

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

	parsedExternalData := externaldata.ParseRefParameters(rule.ExternalData, tplExecutor)

	r := ParsedRule{
		ID:          rule.ID,
		Description: rule.Description,
		Type:        rule.Type,
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
		EventPattern:        rule.EventPattern,
		EntityPatternFields: rule.EntityPatternFields,
		ResolvedStart:       rule.ResolvedStart,
		ResolvedStop:        rule.ResolvedStop,
		NextResolvedStart:   rule.NextResolvedStart,
		NextResolvedStop:    rule.NextResolvedStop,
		ResolvedExdates:     rule.ResolvedExdates,
	}
	if rule.Created != nil {
		r.Created = *rule.Created
	}

	if rule.Updated != nil {
		r.Updated = *rule.Updated
	}

	return r
}

func ExecuteParsedTemplate(
	rule ParsedRule,
	tplName string,
	parsedTpl libtemplate.ParsedTemplate,
	tplData any,
	event *types.Event,
	failureService FailureService,
	templateExecutor libtemplate.Executor,
) (string, error) {
	if parsedTpl.Err != nil {
		failReason := fmt.Sprintf("invalid template %q: %s", tplName, parsedTpl.Err)
		failureService.Add(rule.ID, rule.Description, rule.Updated, FailureTypeInvalidTemplate, failReason, nil)
		return "", parsedTpl.Err
	}

	if parsedTpl.Tpl != nil {
		res, err := templateExecutor.ExecuteByTpl(parsedTpl.Tpl, tplData)
		if err != nil {
			failReason := fmt.Sprintf("cannot execute template %q for event: %s", tplName, err)
			failureService.Add(rule.ID, rule.Description, rule.Updated, FailureTypeInvalidTemplate, failReason, event)
			return "", err
		}

		return res, nil
	}

	return "", nil
}
