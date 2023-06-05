package eventfilter

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
)

type actionProcessor struct {
	templateExecutor  template.Executor
	techMetricsSender techmetrics.Sender
}

func NewActionProcessor(
	templateExecutor template.Executor,
	sender techmetrics.Sender,
) ActionProcessor {
	return &actionProcessor{
		templateExecutor:  templateExecutor,
		techMetricsSender: sender,
	}
}

func (p *actionProcessor) Process(ctx context.Context, action Action, event types.Event, regexMatchWrapper RegexMatchWrapper, externalData map[string]interface{}) (types.Event, error) {
	switch action.Type {
	case ActionSetField:
		err := event.SetField(action.Name, action.Value)
		return event, err
	case ActionSetFieldFromTemplate:
		templateStr, ok := action.Value.(string)
		if !ok {
			return event, ErrShouldBeAString
		}

		value, err := p.templateExecutor.Execute(
			templateStr,
			Template{
				Event:             event,
				RegexMatchWrapper: regexMatchWrapper,
				ExternalData:      externalData,
			}.GetTemplate(),
		)
		if err != nil {
			return event, err
		}

		err = event.SetField(action.Name, value)

		return event, err
	case ActionSetEntityInfo:
		entityUpdated := false

		if !types.IsInfoValueValid(action.Value) {
			return event, types.ErrInvalidInfoType
		}

		*event.Entity, entityUpdated = p.setEntityInfo(*event.Entity, action.Value, action.Name, action.Description)

		event.IsEntityUpdated = event.IsEntityUpdated || entityUpdated

		return event, nil
	case ActionSetEntityInfoFromTemplate:
		templateStr, ok := action.Value.(string)
		if !ok {
			return event, ErrShouldBeAString
		}

		value, err := p.templateExecutor.Execute(
			templateStr,
			Template{
				Event:             event,
				RegexMatchWrapper: regexMatchWrapper,
				ExternalData:      externalData,
			}.GetTemplate(),
		)
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
			return event, err
		}

		err = event.SetField(action.Name, value)
		if err != nil {
			return event, err
		}

		return event, nil
	case ActionCopyToEntityInfo:
		strValue, ok := action.Value.(string)
		if !ok {
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
			return event, err
		}

		if !types.IsInfoValueValid(value) {
			return event, types.ErrInvalidInfoType
		}

		entityUpdated := false
		*event.Entity, entityUpdated = p.setEntityInfo(*event.Entity, value, action.Name, action.Description)

		event.IsEntityUpdated = event.IsEntityUpdated || entityUpdated

		return event, nil
	}

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
