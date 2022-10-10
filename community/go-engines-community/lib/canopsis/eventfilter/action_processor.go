package eventfilter

import (
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
)

type actionProcessor struct {
	templateExecutor *template.Executor
}

func NewActionProcessor(timezoneConfigProvider config.TimezoneConfigProvider) ActionProcessor {
	return &actionProcessor{templateExecutor: template.NewExecutor(timezoneConfigProvider)}
}

func (p *actionProcessor) Process(action Action, event types.Event, regexMatchWrapper RegexMatchWrapper, externalData map[string]interface{}) (types.Event, error) {
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

		*event.Entity, entityUpdated = setEntityInfo(*event.Entity, action.Value, action.Name, action.Description)

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
		*event.Entity, entityUpdated = setEntityInfo(*event.Entity, value, action.Name, action.Description)

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
		*event.Entity, entityUpdated = setEntityInfo(*event.Entity, value, action.Name, action.Description)

		event.IsEntityUpdated = event.IsEntityUpdated || entityUpdated

		return event, nil
	}

	return event, fmt.Errorf("action type = %s is invalid", action.Type)
}

func setEntityInfo(entity types.Entity, value interface{}, name, description string) (types.Entity, bool) {
	info, ok := entity.Infos[name]

	entityUpdated := false
	valueChanged := !ok || info.Value != value
	if valueChanged {
		entityUpdated = true
	}

	info.Name = name
	info.Description = description
	info.Value = value

	if entity.Infos == nil {
		entity.Infos = make(map[string]types.Info)
	}

	entity.Infos[name] = info

	return entity, entityUpdated
}

var ErrShouldBeAString = fmt.Errorf("value should be a string")
