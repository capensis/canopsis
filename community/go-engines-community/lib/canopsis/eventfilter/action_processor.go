package eventfilter

import (
	"bytes"
	"context"
	"fmt"
	"text/template"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
)

type actionProcessor struct {
	buf                   bytes.Buffer
	techMetricsSender     metrics.TechSender
	metricsConfigProvider config.MetricsConfigProvider
}

func NewActionProcessor(provider config.MetricsConfigProvider, sender metrics.TechSender) ActionProcessor {
	return &actionProcessor{
		buf:                   bytes.Buffer{},
		techMetricsSender:     sender,
		metricsConfigProvider: provider,
	}
}

func (p *actionProcessor) Process(ctx context.Context, action Action, event types.Event, regexMatchWrapper RegexMatchWrapper, externalData map[string]interface{}, cfgTimezone *config.TimezoneConfig) (types.Event, error) {
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
			Template{
				Event:             event,
				RegexMatchWrapper: regexMatchWrapper,
				ExternalData:      externalData,
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

		if !types.IsInfoValueValid(action.Value) {
			return event, types.ErrInvalidInfoType
		}

		*event.Entity, entityUpdated = p.setEntityInfo(ctx, *event.Entity, action.Value, action.Name, action.Description)

		event.IsEntityUpdated = event.IsEntityUpdated || entityUpdated

		return event, nil
	case ActionSetEntityInfoFromTemplate:
		templateStr, ok := action.Value.(string)
		if !ok {
			return event, ErrShouldBeAString
		}

		value, err := p.executeTpl(
			templateStr,
			Template{
				Event:             event,
				RegexMatchWrapper: regexMatchWrapper,
				ExternalData:      externalData,
			},
			cfgTimezone,
		)
		if err != nil {
			return event, err
		}

		entityUpdated := false
		*event.Entity, entityUpdated = p.setEntityInfo(ctx, *event.Entity, value, action.Name, action.Description)

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
		*event.Entity, entityUpdated = p.setEntityInfo(ctx, *event.Entity, value, action.Name, action.Description)

		event.IsEntityUpdated = event.IsEntityUpdated || entityUpdated

		return event, nil
	}

	return event, fmt.Errorf("action type = %s is invalid", action.Type)
}

func (p *actionProcessor) executeTpl(tplText string, params TemplateGetter, cfgTimezone *config.TimezoneConfig) (string, error) {
	tpl, err := template.New("tpl").Option("missingkey=error").Funcs(types.GetTemplateFunc(cfgTimezone)).Parse(tplText)
	if err != nil {
		return "", err
	}

	p.buf.Reset()

	err = tpl.Execute(&p.buf, params.GetTemplate())
	if err != nil {
		return "", err
	}

	return p.buf.String(), nil
}

func (p *actionProcessor) setEntityInfo(ctx context.Context, entity types.Entity, value interface{}, name, description string) (types.Entity, bool) {
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

	if p.metricsConfigProvider.Get().EnableTechMetrics {
		go p.techMetricsSender.SendCheEntityInfo(ctx, time.Now(), name)
	}

	return entity, entityUpdated
}

var ErrShouldBeAString = fmt.Errorf("value should be a string")
