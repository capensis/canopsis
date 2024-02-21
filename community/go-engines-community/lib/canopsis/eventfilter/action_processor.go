package eventfilter

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
)

type actionProcessor struct {
	configProvider    config.AlarmConfigProvider
	failureService    FailureService
	templateExecutor  template.Executor
	techMetricsSender techmetrics.Sender
}

func NewActionProcessor(
	configProvider config.AlarmConfigProvider,
	failureService FailureService,
	templateExecutor template.Executor,
	sender techmetrics.Sender,
) ActionProcessor {
	return &actionProcessor{
		configProvider:    configProvider,
		failureService:    failureService,
		templateExecutor:  templateExecutor,
		techMetricsSender: sender,
	}
}

func (p *actionProcessor) Process(
	_ context.Context,
	ruleID string,
	action ParsedAction,
	event *types.Event,
	regexMatch RegexMatch,
	externalData map[string]any,
) (bool, error) {
	const (
		TagsNameVar  = "name"
		TagsValueVar = "value"
	)
	switch action.Type {
	case ActionSetField:
		err := event.SetField(action.Name, action.Value)
		if err != nil {
			failReason := fmt.Sprintf("action %d cannot set %q field: %s", action.Index, action.Name, err.Error())
			p.failureService.Add(ruleID, FailureTypeOther, failReason, nil)
			return false, err
		}

		return false, nil
	case ActionSetFieldFromTemplate:
		value, err := p.actionExecuteParsedTemplate(action, ruleID, "field", event, regexMatch, externalData)
		if err != nil {
			return false, err
		}

		err = event.SetField(action.Name, value)
		if err != nil {
			failReason := fmt.Sprintf("action %d cannot set %q field: %s", action.Index, action.Name, err.Error())
			p.failureService.Add(ruleID, FailureTypeOther, failReason, nil)
			return false, err
		}

		return false, nil
	case ActionSetEntityInfo:
		if !types.IsInfoValueValid(action.Value) {
			failReason := fmt.Sprintf("action %d cannot set %q entity info: invalid type of %v", action.Index,
				action.Name, action.Value)
			p.failureService.Add(ruleID, FailureTypeOther, failReason, nil)
			return false, types.ErrInvalidInfoType
		}

		entityUpdated := p.setEntityInfo(event.Entity, action.Value, action.Name, action.Description)

		return entityUpdated, nil
	case ActionSetEntityInfoFromTemplate:
		value, err := p.actionExecuteParsedTemplate(action, ruleID, "entity info", event, regexMatch, externalData)
		if err != nil {
			return false, err
		}

		entityUpdated := p.setEntityInfo(event.Entity, value, action.Name, action.Description)

		return entityUpdated, nil
	case ActionCopy:
		strValue, ok := action.Value.(string)
		if !ok {
			failReason := fmt.Sprintf("action %d cannot copy to %q field: value %v must be path to field",
				action.Index, action.Name, action.Value)
			p.failureService.Add(ruleID, FailureTypeOther, failReason, nil)
			return false, ErrShouldBeAString
		}

		t := Template{
			Event:        event,
			RegexMatch:   regexMatch,
			ExternalData: externalData,
		}

		value, err := utils.GetField(
			t,
			strValue,
		)
		if err != nil {
			failReason := fmt.Sprintf("action %d cannot copy from %q to %q: %s", action.Index, strValue,
				action.Name, err.Error())
			p.failureService.Add(ruleID, FailureTypeOther, failReason, event)
			return false, err
		}

		err = event.SetField(action.Name, value)
		if err != nil {
			failReason := fmt.Sprintf("action %d cannot copy from %q to %q: %s", action.Index, strValue,
				action.Name, err.Error())
			p.failureService.Add(ruleID, FailureTypeOther, failReason, event)
			return false, err
		}

		return false, nil
	case ActionCopyToEntityInfo:
		strValue, ok := action.Value.(string)
		if !ok {
			failReason := fmt.Sprintf("action %d cannot copy to %q entity info: value %v must be path to field",
				action.Index, action.Name, action.Value)
			p.failureService.Add(ruleID, FailureTypeOther, failReason, nil)
			return false, ErrShouldBeAString
		}

		t := Template{
			Event:        event,
			RegexMatch:   regexMatch,
			ExternalData: externalData,
		}

		value, err := utils.GetField(
			t,
			strValue,
		)
		if err != nil {
			failReason := fmt.Sprintf("action %d cannot copy from %q to %q entity info: %s", action.Index,
				strValue, action.Name, err.Error())
			p.failureService.Add(ruleID, FailureTypeOther, failReason, event)
			return false, err
		}

		if !types.IsInfoValueValid(value) {
			failReason := fmt.Sprintf("action %d cannot copy from %q to %q entity info: invalid type of %v",
				action.Index, strValue, action.Name, value)
			p.failureService.Add(ruleID, FailureTypeOther, failReason, event)
			return false, types.ErrInvalidInfoType
		}

		entityUpdated := p.setEntityInfo(event.Entity, value, action.Name, action.Description)

		return entityUpdated, nil
	case ActionSetTags:
		strValue, ok := action.Value.(string)
		if !ok {
			failReason := fmt.Sprintf("action %d cannot set tags in %q: value %v must be path to field", action.Index,
				action.Name, action.Value)
			p.failureService.Add(ruleID, FailureTypeOther, failReason, nil)

			return false, ErrShouldBeAString
		}
		t := Template{
			Event:        event,
			RegexMatch:   regexMatch,
			ExternalData: externalData,
		}

		value, err := utils.GetField(t, strValue)
		if err != nil {
			failReason := fmt.Sprintf("action %d cannot read source field to set tags in %q: %s", action.Index, action.Name, err)
			p.failureService.Add(ruleID, FailureTypeOther, failReason, event)
			return false, err
		}
		if regexMatch.MatchedRegexp == nil {
			return false, nil
		}
		fieldValue, ok := value.(string)
		if !ok {
			failReason := fmt.Sprintf("action %d cannot assert field's type as string to set tags in %q: %s",
				action.Index, action.Name, err)
			p.failureService.Add(ruleID, FailureTypeOther, failReason, event)

			return false, ErrShouldBeAString
		}

		matches := utils.FindAllStringSubmatchMapWithRegexExpression(regexMatch.MatchedRegexp, fieldValue)
		if len(matches) == 0 {
			return false, nil
		}

		tags := make(map[string]string, len(matches))
		for i := range matches {
			tagName := matches[i][TagsNameVar]
			tagValue := matches[i][TagsValueVar]
			if tagName != "" && tagValue != "" {
				tags[tagName] = tagValue
			}
		}

		if len(tags) == 0 {
			return false, nil
		}

		err = event.SetField(action.Name, tags)
		if err != nil {
			return false, err
		}

		return false, nil
	case ActionSetTagsFromTemplate:
		value, err := p.actionExecuteParsedTemplate(action, ruleID, "tags", event, regexMatch, externalData)
		if err != nil {
			return false, err
		}

		var tags map[string]string
		if regexMatch.MatchedRegexp != nil {
			matches := utils.FindAllStringSubmatchMapWithRegexExpression(regexMatch.MatchedRegexp, value)
			tags = make(map[string]string, len(matches))
			for i := range matches {
				tagName := matches[i][TagsNameVar]
				tagValue := matches[i][TagsValueVar]
				if tagName != "" && tagValue != "" {
					tags[tagName] = tagValue
				}
			}
			if len(tags) == 0 {
				return false, nil
			}
		} else {
			tags = map[string]string{
				action.Name: value,
			}
		}

		err = event.SetField("Tags", tags)
		if err != nil {
			return false, err
		}

		return false, nil
	}

	failReason := fmt.Sprintf("action %d has invalid type %q", action.Index, action.Type)
	p.failureService.Add(ruleID, FailureTypeOther, failReason, event)

	return false, fmt.Errorf("action type = %s is invalid", action.Type)
}

func (p *actionProcessor) actionExecuteParsedTemplate(action ParsedAction, ruleID, target string, event *types.Event, regexMatch RegexMatch, externalData map[string]any) (string, error) {
	if action.ParsedValue.Text == "" {
		failReason := fmt.Sprintf("action %d cannot set %q %s: %v must be template", action.Index,
			action.Name, target, action.Value)
		p.failureService.Add(ruleID, FailureTypeOther, failReason, nil)
		return "", ErrShouldBeAString
	}

	tplData := Template{
		Event:        event,
		RegexMatch:   regexMatch,
		ExternalData: externalData,
	}
	value, err := ExecuteParsedTemplate(ruleID, "Actions."+strconv.Itoa(action.Index)+".Value",
		action.ParsedValue, tplData, event, p.failureService,
		p.templateExecutor)
	return value, err
}

func (p *actionProcessor) setEntityInfo(entity *types.Entity, value any, name, description string) bool {
	enableSorting := p.configProvider.Get().EnableArraySortingInEntityInfos
	if enableSorting {
		if s, ok := utils.IsStringSlice(value); ok {
			sort.Strings(s)
			value = s
		}
	}

	if info, ok := entity.Infos[name]; ok {
		prev := info.Value
		if enableSorting {
			if s, ok := utils.IsStringSlice(info.Value); ok {
				sort.Strings(s)
				prev = s
			}
		}

		if reflect.DeepEqual(prev, value) {
			return false
		}
	}

	if entity.Infos == nil {
		entity.Infos = make(map[string]types.Info, 1)
	}

	entity.Infos[name] = types.Info{
		Name:        name,
		Description: description,
		Value:       value,
	}

	p.techMetricsSender.SendCheEntityInfo(time.Now(), name)

	return true
}

var ErrShouldBeAString = fmt.Errorf("value should be a string")
