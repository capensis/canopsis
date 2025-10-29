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

const (
	TagsNameVar  = "name"
	TagsValueVar = "value"
)

var ErrShouldBeAString = errors.New("value should be a string")

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
	ruleID, ruleDesc string,
	action ParsedAction,
	event *types.Event,
	updatedEntityInfos map[string]UpdatedValue,
	regexMatch RegexMatch,
	externalData map[string]any,
) (map[string]UpdatedValue, error) {
	tplData := Template{
		Event:        event,
		RegexMatch:   regexMatch,
		ExternalData: externalData,
	}
	var err error
	switch action.Type {
	case ActionSetField:
		err = p.setField(ruleID, ruleDesc, action, event)
	case ActionSetFieldFromTemplate:
		err = p.setFieldFromTemplate(ruleID, ruleDesc, action, event, tplData)
	case ActionSetEntityInfo:
		updatedEntityInfos, err = p.setEntityInfo(ruleID, ruleDesc, action, event, updatedEntityInfos)
	case ActionSetEntityInfoFromTemplate:
		updatedEntityInfos, err = p.setEntityInfoFromTemplate(ruleID, ruleDesc, action, event, updatedEntityInfos, tplData)
	case ActionCopy:
		err = p.copy(ruleID, ruleDesc, action, event, tplData)
	case ActionCopyToEntityInfo:
		updatedEntityInfos, err = p.copyEntityInfo(ruleID, ruleDesc, action, event, updatedEntityInfos, tplData)
	case ActionSetTags:
		err = p.setTags(ruleID, ruleDesc, action, event, tplData)
	case ActionSetTagsFromTemplate:
		err = p.setTagsFromTemplate(ruleID, ruleDesc, action, event, tplData)
	case ActionSetEntityInfoFromDictionary:
		updatedEntityInfos, err = p.setEntityInfoFromDictionary(ruleID, ruleDesc, action, event, updatedEntityInfos, tplData)
	default:
		failReason := fmt.Sprintf("action %d has invalid type %q", action.Index, action.Type)
		p.failureService.Add(ruleID, ruleDesc, FailureTypeOther, failReason, event)

		err = fmt.Errorf("action type = %s is invalid", action.Type)
	}

	return updatedEntityInfos, err
}

func (p *actionProcessor) setField(
	ruleID, ruleDesc string,
	action ParsedAction,
	event *types.Event,
) error {
	err := event.SetField(action.Name, action.Value)
	if err != nil {
		failReason := fmt.Sprintf("action %d cannot set %q field: %s", action.Index, action.Name, err.Error())
		p.failureService.Add(ruleID, ruleDesc, FailureTypeOther, failReason, nil)

		return err
	}

	return nil
}

func (p *actionProcessor) setFieldFromTemplate(
	ruleID, ruleDesc string,
	action ParsedAction,
	event *types.Event,
	tplData Template,
) error {
	value, err := p.executeParsedTemplate(action, ruleID, ruleDesc, "field", event, tplData)
	if err != nil {
		return err
	}

	err = event.SetField(action.Name, value)
	if err != nil {
		failReason := fmt.Sprintf("action %d cannot set %q field: %s", action.Index, action.Name, err.Error())
		p.failureService.Add(ruleID, ruleDesc, FailureTypeOther, failReason, nil)

		return err
	}

	return nil
}

func (p *actionProcessor) setEntityInfo(
	ruleID, ruleDesc string,
	action ParsedAction,
	event *types.Event,
	updatedEntityInfos map[string]UpdatedValue,
) (map[string]UpdatedValue, error) {
	if !types.IsInfoValueValid(action.Value) {
		failReason := fmt.Sprintf("action %d cannot set %q entity info: invalid type of %v", action.Index,
			action.Name, action.Value)
		p.failureService.Add(ruleID, ruleDesc, FailureTypeOther, failReason, nil)

		return updatedEntityInfos, types.ErrInvalidInfoType
	}

	updatedEntityInfos = p.updateEntity(ruleID, event.Entity, action.Value, action.Name, action.Description, updatedEntityInfos)

	return updatedEntityInfos, nil
}

func (p *actionProcessor) setEntityInfoFromTemplate(
	ruleID, ruleDesc string,
	action ParsedAction,
	event *types.Event,
	updatedEntityInfos map[string]UpdatedValue,
	tplData Template,
) (map[string]UpdatedValue, error) {
	value, err := p.executeParsedTemplate(action, ruleID, ruleDesc, "entity info", event, tplData)
	if err != nil {
		return updatedEntityInfos, err
	}

	updatedEntityInfos = p.updateEntity(ruleID, event.Entity, value, action.Name, action.Description, updatedEntityInfos)

	return updatedEntityInfos, nil
}

func (p *actionProcessor) copy(
	ruleID, ruleDesc string,
	action ParsedAction,
	event *types.Event,
	tplData Template,
) error {
	strValue, ok := action.Value.(string)
	if !ok {
		failReason := fmt.Sprintf("action %d cannot copy to %q field: value %v must be path to field",
			action.Index, action.Name, action.Value)
		p.failureService.Add(ruleID, ruleDesc, FailureTypeOther, failReason, nil)

		return ErrShouldBeAString
	}

	value, err := utils.GetField(
		tplData,
		strValue,
	)
	if err != nil {
		failReason := fmt.Sprintf("action %d cannot copy from %q to %q: %s", action.Index, strValue,
			action.Name, err.Error())
		p.failureService.Add(ruleID, ruleDesc, FailureTypeOther, failReason, event)

		return err
	}

	err = event.SetField(action.Name, value)
	if err != nil {
		failReason := fmt.Sprintf("action %d cannot copy from %q to %q: %s", action.Index, strValue,
			action.Name, err.Error())
		p.failureService.Add(ruleID, ruleDesc, FailureTypeOther, failReason, event)

		return err
	}

	return nil
}

func (p *actionProcessor) copyEntityInfo(
	ruleID, ruleDesc string,
	action ParsedAction,
	event *types.Event,
	updatedEntityInfos map[string]UpdatedValue,
	tplData Template,
) (map[string]UpdatedValue, error) {
	strValue, ok := action.Value.(string)
	if !ok {
		failReason := fmt.Sprintf("action %d cannot copy to %q entity info: value %v must be path to field",
			action.Index, action.Name, action.Value)
		p.failureService.Add(ruleID, ruleDesc, FailureTypeOther, failReason, nil)

		return updatedEntityInfos, ErrShouldBeAString
	}

	value, err := utils.GetField(
		tplData,
		strValue,
	)
	if err != nil {
		failReason := fmt.Sprintf("action %d cannot copy from %q to %q entity info: %s", action.Index,
			strValue, action.Name, err.Error())
		p.failureService.Add(ruleID, ruleDesc, FailureTypeOther, failReason, event)

		return updatedEntityInfos, err
	}

	if !types.IsInfoValueValid(value) {
		failReason := fmt.Sprintf("action %d cannot copy from %q to %q entity info: invalid type of %v",
			action.Index, strValue, action.Name, value)
		p.failureService.Add(ruleID, ruleDesc, FailureTypeOther, failReason, event)

		return updatedEntityInfos, types.ErrInvalidInfoType
	}

	updatedEntityInfos = p.updateEntity(ruleID, event.Entity, value, action.Name, action.Description, updatedEntityInfos)

	return updatedEntityInfos, nil
}

func (p *actionProcessor) setTags(
	ruleID, ruleDesc string,
	action ParsedAction,
	event *types.Event,
	tplData Template,
) error {
	strValue, ok := action.Value.(string)
	if !ok {
		failReason := fmt.Sprintf("action %d cannot set tags in %q: value %v must be path to field", action.Index,
			action.Name, action.Value)
		p.failureService.Add(ruleID, ruleDesc, FailureTypeOther, failReason, nil)

		return ErrShouldBeAString
	}

	value, err := utils.GetField(tplData, strValue)
	if err != nil {
		failReason := fmt.Sprintf("action %d cannot read source field to set tags in %q: %s", action.Index, action.Name, err)
		p.failureService.Add(ruleID, ruleDesc, FailureTypeOther, failReason, event)

		return err
	}

	if tplData.RegexMatch.MatchedRegexp == nil {
		return nil
	}

	fieldValue, ok := value.(string)
	if !ok {
		failReason := fmt.Sprintf("action %d cannot assert field's type as string to set tags in %q: %s",
			action.Index, action.Name, err)
		p.failureService.Add(ruleID, ruleDesc, FailureTypeOther, failReason, event)

		return ErrShouldBeAString
	}

	matches := utils.FindAllStringSubmatchMapWithRegexExpression(tplData.RegexMatch.MatchedRegexp, fieldValue)
	if len(matches) == 0 {
		return nil
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
		return nil
	}

	return event.SetField(action.Name, tags)
}

func (p *actionProcessor) setTagsFromTemplate(
	ruleID, ruleDesc string,
	action ParsedAction,
	event *types.Event,
	tplData Template,
) error {
	value, err := p.executeParsedTemplate(action, ruleID, ruleDesc, "tags", event, tplData)
	if err != nil {
		return err
	}

	var tags map[string]string
	if tplData.RegexMatch.MatchedRegexp != nil {
		matches := utils.FindAllStringSubmatchMapWithRegexExpression(tplData.RegexMatch.MatchedRegexp, value)
		tags = make(map[string]string, len(matches))
		for i := range matches {
			tagName := matches[i][TagsNameVar]
			tagValue := matches[i][TagsValueVar]
			if tagName != "" && tagValue != "" {
				tags[tagName] = tagValue
			}
		}

		if len(tags) == 0 {
			return nil
		}
	} else {
		tags = map[string]string{
			action.Name: value,
		}
	}

	return event.SetField("Tags", tags)
}

func (p *actionProcessor) setEntityInfoFromDictionary(
	ruleID, ruleDesc string,
	action ParsedAction,
	event *types.Event,
	updatedEntityInfos map[string]UpdatedValue,
	tplData Template,
) (map[string]UpdatedValue, error) {
	strValue, ok := action.Value.(string)
	if !ok {
		failReason := fmt.Sprintf("action %d cannot set entity info in %q: value %v must be path to field", action.Index,
			action.Name, action.Value)
		p.failureService.Add(ruleID, ruleDesc, FailureTypeOther, failReason, nil)

		return updatedEntityInfos, ErrShouldBeAString
	}

	dictValue, err := utils.GetField(tplData, strValue)
	if err != nil {
		if errors.Is(err, utils.ErrFieldNotExist) {
			return updatedEntityInfos, nil
		}

		failReason := fmt.Sprintf("action %d cannot read source field to set entity info in %q: %s",
			action.Index, action.Name, err)
		p.failureService.Add(ruleID, ruleDesc, FailureTypeOther, failReason, event)

		return updatedEntityInfos, err
	}

	dict, ok := dictValue.(map[string]any)
	if !ok {
		failReason := fmt.Sprintf("action %d cannot assert field's type as map to set entity info in %q",
			action.Index, action.Name)
		p.failureService.Add(ruleID, ruleDesc, FailureTypeOther, failReason, event)

		return updatedEntityInfos, errors.New("value should be a map")
	}

	for k, v := range dict {
		if !types.IsInfoValueValid(v) {
			err = fmt.Errorf("value %v for %q is invalid", v, k)
			failReason := fmt.Sprintf("action %d cannot set entity info in %q: %s", action.Index, action.Name, err)
			p.failureService.Add(ruleID, ruleDesc, FailureTypeOther, failReason, event)

			return updatedEntityInfos, err
		}

		updatedEntityInfos = p.updateEntity(ruleID, event.Entity, v, k, action.Description, updatedEntityInfos)
	}

	return updatedEntityInfos, nil
}

func (p *actionProcessor) executeParsedTemplate(action ParsedAction, ruleID, ruleDesc, target string, event *types.Event, tplData Template) (string, error) {
	if action.ParsedValue.Text == "" {
		failReason := fmt.Sprintf("action %d cannot set %q %s: %v must be template", action.Index,
			action.Name, target, action.Value)
		p.failureService.Add(ruleID, ruleDesc, FailureTypeOther, failReason, nil)

		return "", ErrShouldBeAString
	}

	value, err := ExecuteParsedTemplate(ruleID, ruleDesc, "Actions."+strconv.Itoa(action.Index)+".Value",
		action.ParsedValue, tplData, event, p.failureService,
		p.templateExecutor)

	return value, err
}

func (p *actionProcessor) updateEntity(
	ruleID string,
	entity *types.Entity,
	newValue any,
	name, description string,
	changes map[string]UpdatedValue,
) map[string]UpdatedValue {
	enableSorting := p.configProvider.Get().EnableArraySortingInEntityInfos
	if enableSorting {
		if s, ok := utils.IsStringSlice(newValue); ok {
			sort.Strings(s)
			newValue = s
		}
	}

	var oldValue any
	if v, ok := changes[name]; ok {
		oldValue = v.OldValue
		if reflect.DeepEqual(oldValue, newValue) {
			delete(changes, name)
			entity.Infos[name] = types.Info{
				Name:        name,
				Description: description,
				Value:       newValue,
			}

			return changes
		}
	} else if info, ok := entity.Infos[name]; ok {
		oldValue = info.Value
		if enableSorting {
			if s, ok := utils.IsStringSlice(info.Value); ok {
				sort.Strings(s)
				oldValue = s
			}
		}

		if reflect.DeepEqual(oldValue, newValue) {
			return changes
		}
	}

	if changes == nil {
		changes = make(map[string]UpdatedValue, 1)
	}

	changes[name] = UpdatedValue{
		RuleID:   ruleID,
		OldValue: oldValue,
		NewValue: newValue,
	}

	if entity.Infos == nil {
		entity.Infos = make(map[string]types.Info, 1)
	}

	entity.Infos[name] = types.Info{
		Name:        name,
		Description: description,
		Value:       newValue,
	}

	return changes
}
