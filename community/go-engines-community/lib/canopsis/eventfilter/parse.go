package eventfilter

import (
	"fmt"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	libtemplate "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

// ParsedExternalDataParameters extends externaldata.ParsedRefParameters with the OnMismatch
// value carried over from ExternalDataParameters.
type ParsedExternalDataParameters struct {
	externaldata.ParsedRefParameters
	OnMismatch int
}

type ParsedRule struct {
	ID           string
	Name         string
	Type         string
	Config       ParsedRuleConfig
	ExternalData []ParsedExternalDataParameters
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
	Upstream      libtemplate.ParsedTemplate

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

	// ExternalDataRefs lists the ExternalData references (Rule.ExternalData[*].Reference)
	// this action's Value depends on. It's used to skip actions relying on external data
	// that couldn't be fetched when OnMismatchSkipActions applies.
	ExternalDataRefs []string
}

func ParseRule(rule Rule, tplExecutor libtemplate.Executor) ParsedRule {
	refs := make([]string, len(rule.ExternalData))
	for i, d := range rule.ExternalData {
		refs[i] = d.Reference
	}

	parsedActions := make([]ParsedAction, len(rule.Config.Actions))
	for i, action := range rule.Config.Actions {
		var parsedValue libtemplate.ParsedTemplate
		var externalDataRefs []string
		switch action.Type {
		case ActionSetFieldFromTemplate, ActionSetEntityInfoFromTemplate, ActionSetTagsFromTemplate:
			if str, ok := action.Value.(string); ok {
				parsedValue = tplExecutor.Parse(str)
				externalDataRefs = templateExternalDataRefs(parsedValue, refs)
			}
		case ActionCopy, ActionCopyToEntityInfo, ActionSetTags, ActionSetEntityInfoFromDictionary:
			if str, ok := action.Value.(string); ok {
				for _, ref := range refs {
					prefix := "ExternalData." + ref
					if str == prefix || strings.HasPrefix(str, prefix+".") {
						externalDataRefs = append(externalDataRefs, ref)
					}
				}
			}
		}

		parsedActions[i] = ParsedAction{
			Index:            i,
			Type:             action.Type,
			Name:             action.Name,
			Description:      action.Description,
			Value:            action.Value,
			ParsedValue:      parsedValue,
			ExternalDataRefs: externalDataRefs,
		}
	}

	refParameters := make([]externaldata.RefParameters, len(rule.ExternalData))
	for i, d := range rule.ExternalData {
		refParameters[i] = d.RefParameters
	}

	parsedRefParameters := externaldata.ParseRefParameters(refParameters, tplExecutor)
	parsedExternalData := make([]ParsedExternalDataParameters, len(parsedRefParameters))
	for i, parsed := range parsedRefParameters {
		parsedExternalData[i] = ParsedExternalDataParameters{
			ParsedRefParameters: parsed,
			OnMismatch:          rule.ExternalData[i].OnMismatch,
		}
	}

	r := ParsedRule{
		ID:   rule.ID,
		Name: rule.Name,
		Type: rule.Type,
		Config: ParsedRuleConfig{
			Resource:      tplExecutor.Parse(rule.Config.Resource),
			Component:     tplExecutor.Parse(rule.Config.Component),
			Connector:     tplExecutor.Parse(rule.Config.Connector),
			ConnectorName: tplExecutor.Parse(rule.Config.ConnectorName),
			Upstream:      tplExecutor.Parse(rule.Config.Upstream),
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
		failureService.Add(rule.ID, rule.Name, rule.Updated, FailureTypeInvalidTemplate, failReason, nil)
		return "", parsedTpl.Err
	}

	if parsedTpl.Tpl != nil {
		res, err := templateExecutor.ExecuteByTpl(parsedTpl.Tpl, tplData)
		if err != nil {
			failReason := fmt.Sprintf("cannot execute template %q for event: %s", tplName, err)
			failureService.Add(rule.ID, rule.Name, rule.Updated, FailureTypeInvalidTemplate, failReason, event)
			return "", err
		}

		return res, nil
	}

	return "", nil
}
