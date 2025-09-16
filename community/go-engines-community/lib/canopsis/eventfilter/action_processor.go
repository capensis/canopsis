package eventfilter

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
)

type actionProcessor struct {
	configProvider   config.AlarmConfigProvider
	failureService   FailureService
	templateExecutor template.Executor
}

func NewActionProcessor(
	configProvider config.AlarmConfigProvider,
	failureService FailureService,
	templateExecutor template.Executor,
) ActionProcessor {
	return &actionProcessor{
		configProvider:   configProvider,
		failureService:   failureService,
		templateExecutor: templateExecutor,
	}
}

func (p *actionProcessor) Process(
	_ context.Context,
	ruleID string,
	action ParsedAction,
	event *types.Event,
	regexMatch RegexMatch,
	externalData map[string]any,
) (map[string]UpdatedValue, error) {
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

			return nil, err
		}

		return nil, nil
	case ActionSetFieldFromTemplate:
		value, err := p.actionExecuteParsedTemplate(action, ruleID, "field", event, regexMatch, externalData)
		if err != nil {
			return nil, err
		}

		err = event.SetField(action.Name, value)
		if err != nil {
			failReason := fmt.Sprintf("action %d cannot set %q field: %s", action.Index, action.Name, err.Error())
			p.failureService.Add(ruleID, FailureTypeOther, failReason, nil)
			return nil, err
		}

		return nil, nil
	case ActionSetEntityInfo:
		if !types.IsInfoValueValid(action.Value) {
			failReason := fmt.Sprintf("action %d cannot set %q entity info: invalid type of %v", action.Index,
				action.Name, action.Value)
			p.failureService.Add(ruleID, FailureTypeOther, failReason, nil)

			return nil, types.ErrInvalidInfoType
		}

		if !p.setEntityInfo(event.Entity, action.Value, action.Name, action.Description) {
			return nil, nil
		}

		return map[string]UpdatedValue{action.Name: {RuleID: ruleID, Value: action.Value}}, nil
	case ActionSetEntityInfoFromTemplate:
		value, err := p.actionExecuteParsedTemplate(action, ruleID, "entity info", event, regexMatch, externalData)
		if err != nil {
			return nil, err
		}

		if !p.setEntityInfo(event.Entity, value, action.Name, action.Description) {
			return nil, nil
		}

		return map[string]UpdatedValue{action.Name: {RuleID: ruleID, Value: value}}, nil
	case ActionCopy:
		strValue, ok := action.Value.(string)
		if !ok {
			failReason := fmt.Sprintf("action %d cannot copy to %q field: value %v must be path to field",
				action.Index, action.Name, action.Value)
			p.failureService.Add(ruleID, FailureTypeOther, failReason, nil)

			return nil, ErrShouldBeAString
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

			return nil, err
		}

		err = event.SetField(action.Name, value)
		if err != nil {
			failReason := fmt.Sprintf("action %d cannot copy from %q to %q: %s", action.Index, strValue,
				action.Name, err.Error())
			p.failureService.Add(ruleID, FailureTypeOther, failReason, event)

			return nil, err
		}

		return nil, nil
	case ActionCopyToEntityInfo:
		strValue, ok := action.Value.(string)
		if !ok {
			failReason := fmt.Sprintf("action %d cannot copy to %q entity info: value %v must be path to field",
				action.Index, action.Name, action.Value)
			p.failureService.Add(ruleID, FailureTypeOther, failReason, nil)

			return nil, ErrShouldBeAString
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

			return nil, err
		}

		if !types.IsInfoValueValid(value) {
			failReason := fmt.Sprintf("action %d cannot copy from %q to %q entity info: invalid type of %v",
				action.Index, strValue, action.Name, value)
			p.failureService.Add(ruleID, FailureTypeOther, failReason, event)

			return nil, types.ErrInvalidInfoType
		}

		if !p.setEntityInfo(event.Entity, value, action.Name, action.Description) {
			return nil, nil
		}

		return map[string]UpdatedValue{action.Name: {RuleID: ruleID, Value: value}}, nil
	case ActionSetTags:
		strValue, ok := action.Value.(string)
		if !ok {
			failReason := fmt.Sprintf("action %d cannot set tags in %q: value %v must be path to field", action.Index,
				action.Name, action.Value)
			p.failureService.Add(ruleID, FailureTypeOther, failReason, nil)

			return nil, ErrShouldBeAString
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

			return nil, err
		}

		if regexMatch.MatchedRegexp == nil {
			return nil, nil
		}

		fieldValue, ok := value.(string)
		if !ok {
			failReason := fmt.Sprintf("action %d cannot assert field's type as string to set tags in %q: %s",
				action.Index, action.Name, err)
			p.failureService.Add(ruleID, FailureTypeOther, failReason, event)

			return nil, ErrShouldBeAString
		}

		matches := utils.FindAllStringSubmatchMapWithRegexExpression(regexMatch.MatchedRegexp, fieldValue)
		if len(matches) == 0 {
			return nil, nil
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
			return nil, nil
		}

		err = event.SetField(action.Name, tags)
		if err != nil {
			return nil, err
		}

		return nil, nil
	case ActionSetTagsFromTemplate:
		value, err := p.actionExecuteParsedTemplate(action, ruleID, "tags", event, regexMatch, externalData)
		if err != nil {
			return nil, err
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
				return nil, nil
			}
		} else {
			tags = map[string]string{
				action.Name: value,
			}
		}

		err = event.SetField("Tags", tags)
		if err != nil {
			return nil, err
		}

		return nil, nil
	case ActionSetEntityInfoFromDictionary:
		strValue, ok := action.Value.(string)
		if !ok {
			failReason := fmt.Sprintf("action %d cannot set entity info in %q: value %v must be path to field", action.Index,
				action.Name, action.Value)
			p.failureService.Add(ruleID, FailureTypeOther, failReason, nil)

			return nil, ErrShouldBeAString
		}

		t := Template{
			Event:        event,
			RegexMatch:   regexMatch,
			ExternalData: externalData,
		}
		value, err := utils.GetField(t, strValue)
		if err != nil {
			if errors.Is(err, utils.ErrFieldNotExist) {
				return nil, nil
			}

			failReason := fmt.Sprintf("action %d cannot read source field to set entity info in %q: %s",
				action.Index, action.Name, err)
			p.failureService.Add(ruleID, FailureTypeOther, failReason, event)

			return nil, err
		}

		dict, ok := value.(map[string]any)
		if !ok {
			failReason := fmt.Sprintf("action %d cannot assert field's type as map to set entity info in %q",
				action.Index, action.Name)
			p.failureService.Add(ruleID, FailureTypeOther, failReason, event)

			return nil, errors.New("value should be a map")
		}

		updatedValues, ok, err := p.setInfosFromDict(event.Entity, dict, ruleID, action.Description)
		if err != nil {
			failReason := fmt.Sprintf("action %d cannot set entity info in %q: %s", action.Index, action.Name, err)
			p.failureService.Add(ruleID, FailureTypeOther, failReason, event)

			return nil, err
		}

		if !ok {
			return nil, nil
		}

		return updatedValues, nil
	}

	failReason := fmt.Sprintf("action %d has invalid type %q", action.Index, action.Type)
	p.failureService.Add(ruleID, FailureTypeOther, failReason, event)

	return nil, fmt.Errorf("action type = %s is invalid", action.Type)
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

	return true
}

func (p *actionProcessor) setInfosFromDict(entity *types.Entity, dict map[string]any, ruleID, description string) (map[string]UpdatedValue, bool, error) {
	updatedValues := make(map[string]UpdatedValue)
	updated := false
	for name, value := range dict {
		if !types.IsInfoValueValid(value) {
			return nil, false, fmt.Errorf("value %v for %q is invalid", value, name)
		}

		if p.setEntityInfo(entity, value, name, description) {
			updated = true
			updatedValues[name] = UpdatedValue{RuleID: ruleID, Value: value}
		}
	}

	return updatedValues, updated, nil
}

var ErrShouldBeAString = errors.New("value should be a string")
