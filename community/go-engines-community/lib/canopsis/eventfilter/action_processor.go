package eventfilter

import (
	"bytes"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"text/template"
)

type actionProcessor struct {
	buf bytes.Buffer
}

func NewActionProcessor() ActionProcessor {
	return &actionProcessor{buf: bytes.Buffer{}}
}

func (p *actionProcessor) Process(action Action, event types.Event, regexMatch pattern.EventRegexMatches, externalData map[string]interface{}, cfgTimezone *config.TimezoneConfig) (types.Event, error) {
	switch action.Type {
	case ActionSetField:
		err := event.SetField(action.Name, action.Value)
		return event, err
	case ActionSetFieldFromTemplate:
		templateStr, ok := action.Value.(string)
		if !ok {
			return event, ErrShouldBeAString
		}

		value, err := p.executeTpl(
			templateStr,
			TemplateParameters{
				Event:        event,
				RegexMatch:   regexMatch,
				ExternalData: externalData,
			},
			cfgTimezone,
		)
		if err != nil {
			return event, err
		}

		err = event.SetField(action.Name, value)

		return event, err
	case ActionSetEntityInfo:
		entityUpdated := false
		*event.Entity, entityUpdated = setEntityInfo(*event.Entity, action.Value, action.Name, action.Description)

		event.IsEntityUpdated = event.IsEntityUpdated || entityUpdated

		return event, nil
	case ActionSetEntityInfoFromTemplate:
		templateStr, ok := action.Value.(string)
		if !ok {
			return event, ErrShouldBeAString
		}

		value, err := p.executeTpl(
			templateStr,
			TemplateParameters{
				Event:        event,
				RegexMatch:   regexMatch,
				ExternalData: externalData,
			},
			cfgTimezone,
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

		value, err := utils.GetField(
			TemplateParameters{
				Event:        event,
				RegexMatch:   regexMatch,
				ExternalData: externalData,
			},
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

		value, err := utils.GetField(
			TemplateParameters{
				Event:        event,
				RegexMatch:   regexMatch,
				ExternalData: externalData,
			},
			strValue,
		)
		if err != nil {
			return event, err
		}

		entityUpdated := false
		*event.Entity, entityUpdated = setEntityInfo(*event.Entity, value, action.Name, action.Description)

		event.IsEntityUpdated = event.IsEntityUpdated || entityUpdated

		return event, nil
	}

	return event, fmt.Errorf("action type = %s is invalid", action.Type)
}

func (p *actionProcessor) executeTpl(tplText string, params TemplateParameters, cfgTimezone *config.TimezoneConfig) (string, error) {
	tpl, err := template.New("tpl").Option("missingkey=error").Funcs(types.GetTemplateFunc(cfgTimezone)).Parse(tplText)
	if err != nil {
		return "", err
	}

	p.buf.Reset()

	err = tpl.Execute(&p.buf, params)
	if err != nil {
		return "", err
	}

	return p.buf.String(), nil
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
