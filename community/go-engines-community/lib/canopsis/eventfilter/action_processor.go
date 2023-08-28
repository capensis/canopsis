package eventfilter

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
)

type actionProcessor struct {
	failureService    FailureService
	templateExecutor  template.Executor
	techMetricsSender techmetrics.Sender
}

func NewActionProcessor(
	failureService FailureService,
	templateExecutor template.Executor,
	sender techmetrics.Sender,
) ActionProcessor {
	return &actionProcessor{
		failureService:    failureService,
		templateExecutor:  templateExecutor,
		techMetricsSender: sender,
	}
}

func (p *actionProcessor) Process(
	_ context.Context,
	ruleID string,
	action ParsedAction,
	event types.Event,
	regexMatchWrapper RegexMatchWrapper,
	externalData map[string]any,
) (types.Event, error) {
	switch action.Type {
	case ActionSetField:
		err := event.SetField(action.Name, action.Value)
		if err != nil {
			failReason := fmt.Sprintf("action %d cannot set %q field: %s", action.Index, action.Name, err.Error())
			p.failureService.Add(ruleID, FailureTypeOther, failReason, nil)
			return event, err
		}

		return event, nil
	case ActionSetFieldFromTemplate:
		if action.ParsedValue.Text == "" {
			failReason := fmt.Sprintf("action %d cannot set %q field: %v must be template", action.Index,
				action.Name, action.Value)
			p.failureService.Add(ruleID, FailureTypeOther, failReason, nil)
			return event, ErrShouldBeAString
		}

		tplData := Template{
			Event:             event,
			RegexMatchWrapper: regexMatchWrapper,
			ExternalData:      externalData,
		}.GetTemplate()
		value, err := ExecuteParsedTemplate(ruleID, "Actions."+strconv.Itoa(action.Index)+".Value",
			action.ParsedValue, tplData, event, p.failureService,
			p.templateExecutor)
		if err != nil {
			return event, err
		}

		err = event.SetField(action.Name, value)
		if err != nil {
			failReason := fmt.Sprintf("action %d cannot set %q field: %s", action.Index, action.Name, err.Error())
			p.failureService.Add(ruleID, FailureTypeOther, failReason, nil)
			return event, err
		}

		return event, nil
	case ActionSetEntityInfo:
		entityUpdated := false

		if !types.IsInfoValueValid(action.Value) {
			failReason := fmt.Sprintf("action %d cannot set %q entity info: invalid type of %v", action.Index,
				action.Name, action.Value)
			p.failureService.Add(ruleID, FailureTypeOther, failReason, nil)
			return event, types.ErrInvalidInfoType
		}

		*event.Entity, entityUpdated = p.setEntityInfo(*event.Entity, action.Value, action.Name, action.Description)

		event.IsEntityUpdated = event.IsEntityUpdated || entityUpdated

		return event, nil
	case ActionSetEntityInfoFromTemplate:
		if action.ParsedValue.Text == "" {
			failReason := fmt.Sprintf("action %d cannot set %q entity info: %v must be template", action.Index,
				action.Name, action.Value)
			p.failureService.Add(ruleID, FailureTypeOther, failReason, nil)
			return event, ErrShouldBeAString
		}

		tplData := Template{
			Event:             event,
			RegexMatchWrapper: regexMatchWrapper,
			ExternalData:      externalData,
		}.GetTemplate()
		value, err := ExecuteParsedTemplate(ruleID, "Actions."+strconv.Itoa(action.Index)+".Value",
			action.ParsedValue, tplData, event, p.failureService,
			p.templateExecutor)
		if err != nil {
			return event, err
		}

		entityUpdated := false
		*event.Entity, entityUpdated = p.setEntityInfo(*event.Entity, value, action.Name, action.Description)

		event.IsEntityUpdated = event.IsEntityUpdated || entityUpdated

		return event, nil
	case ActionCopy:
		strValue, ok := action.Value.(string)
		if !ok {
			failReason := fmt.Sprintf("action %d cannot copy to %q field: value %v must be path to field",
				action.Index, action.Name, action.Value)
			p.failureService.Add(ruleID, FailureTypeOther, failReason, nil)
			return event, ErrShouldBeAString
		}

		t := Template{
			Event:             event,
			RegexMatchWrapper: regexMatchWrapper,
			ExternalData:      externalData,
		}

		value, err := utils.GetField(
			t.GetTemplate(),
			strValue,
		)
		if err != nil {
			failReason := fmt.Sprintf("action %d cannot copy from %q to %q: %s", action.Index, strValue,
				action.Name, err.Error())
			p.failureService.Add(ruleID, FailureTypeOther, failReason, &event)
			return event, err
		}

		err = event.SetField(action.Name, value)
		if err != nil {
			failReason := fmt.Sprintf("action %d cannot copy from %q to %q: %s", action.Index, strValue,
				action.Name, err.Error())
			p.failureService.Add(ruleID, FailureTypeOther, failReason, &event)
			return event, err
		}

		return event, nil
	case ActionCopyToEntityInfo:
		strValue, ok := action.Value.(string)
		if !ok {
			failReason := fmt.Sprintf("action %d cannot copy to %q entity info: value %v must be path to field",
				action.Index, action.Name, action.Value)
			p.failureService.Add(ruleID, FailureTypeOther, failReason, nil)
			return event, ErrShouldBeAString
		}

		t := Template{
			Event:             event,
			RegexMatchWrapper: regexMatchWrapper,
			ExternalData:      externalData,
		}

		value, err := utils.GetField(
			t.GetTemplate(),
			strValue,
		)
		if err != nil {
			failReason := fmt.Sprintf("action %d cannot copy from %q to %q entity info: %s", action.Index,
				strValue, action.Name, err.Error())
			p.failureService.Add(ruleID, FailureTypeOther, failReason, &event)
			return event, err
		}

		if !types.IsInfoValueValid(value) {
			failReason := fmt.Sprintf("action %d cannot copy from %q to %q entity info: invalid type of %v",
				action.Index, strValue, action.Name, value)
			p.failureService.Add(ruleID, FailureTypeOther, failReason, &event)
			return event, types.ErrInvalidInfoType
		}

		entityUpdated := false
		*event.Entity, entityUpdated = p.setEntityInfo(*event.Entity, value, action.Name, action.Description)

		event.IsEntityUpdated = event.IsEntityUpdated || entityUpdated

		return event, nil
	}

	failReason := fmt.Sprintf("action %d has invalid type %q", action.Index, action.Type)
	p.failureService.Add(ruleID, FailureTypeOther, failReason, &event)
	return event, fmt.Errorf("action type = %s is invalid", action.Type)
}

func (p *actionProcessor) setEntityInfo(entity types.Entity, value interface{}, name, description string) (types.Entity, bool) {
	if info, ok := entity.Infos[name]; ok {
		if reflect.DeepEqual(info.Value, value) {
			return entity, false
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

	return entity, true
}

var ErrShouldBeAString = fmt.Errorf("value should be a string")
